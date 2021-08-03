package ctrlutils

import (
	"context"
	"fmt"
	"net"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/blang/semver/v4"
	"github.com/go-logr/logr"
	"github.com/kong/deck/file"
	"github.com/prometheus/common/log"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	knative "knative.dev/networking/pkg/apis/networking/v1alpha1"
	knativeversioned "knative.dev/networking/pkg/client/clientset/versioned"
	knativeApis "knative.dev/pkg/apis"
	"knative.dev/pkg/network"

	"github.com/kong/kubernetes-ingress-controller/internal/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/internal/util"
	kicclientset "github.com/kong/kubernetes-ingress-controller/pkg/clientset"
)

const (
	statusUpdateRetry    = 3
	statusUpdateWaitTick = time.Second
)

// PullConfigUpdate is a dedicated function that process ingress/customer resource status update after configuration is updated within kong.
func PullConfigUpdate(
	ctx context.Context,
	kongConfig sendconfig.Kong,
	log logr.Logger,
	kubeConfig *rest.Config,
	publishService string,
	publishAddresses []string,
) {
	ips, hostname, err := RunningAddresses(ctx, kubeConfig, publishService, publishAddresses)
	if err != nil {
		log.Error(err, "failed to determine kong proxy external ips/hostnames.")
		return
	}

	cli, err := clientset.NewForConfig(kubeConfig)
	if err != nil {
		log.Error(err, "failed to create k8s client.")
		return
	}

	versionInfo, err := cli.ServerVersion()
	if err != nil {
		log.Error(err, "failed to retrieve cluster version")
		return
	}

	kubernetesVersion, err := semver.Parse(strings.TrimPrefix(versionInfo.String(), "v"))
	if err != nil {
		log.Error(err, "could not parse cluster version")
		return
	}

	kiccli, err := kicclientset.NewForConfig(kubeConfig)
	if err != nil {
		log.Error(err, "failed to create kong ingress client.")
		return
	}

	log.Info("Launching Ingress Status Update Thread.")

	var wg sync.WaitGroup
	for {
		select {
		case updateDone := <-kongConfig.ConfigDone:
			log.V(4).Info("receive configuration information. Update ingress status %v \n", updateDone)
			wg.Add(1)
			go func() {
				if err := UpdateStatuses(ctx, &updateDone, log, cli, kiccli, &wg, ips, hostname, kubeConfig, kubernetesVersion); err != nil {
					log.Error(err, "failed to update resource statuses")
				}
			}()
		case <-ctx.Done():
			log.Info("stop status update channel.")
			wg.Wait()
			return
		}
	}
}

// UpdateStatuses update resources statuses according to generated rules and specs
func UpdateStatuses(
	ctx context.Context,
	targetContent *file.Content,
	log logr.Logger,
	cli *clientset.Clientset,
	kiccli *kicclientset.Clientset,
	wg *sync.WaitGroup,
	ips []string,
	hostname string,
	kubeConfig *rest.Config,
	kubernetesVersion semver.Version,
) error {
	defer wg.Done()

	for _, svc := range targetContent.Services {
		for _, plugin := range svc.Plugins {
			log.V(5).Info("\n service host %s name %s plugin enablement %v\n", *svc.Service.Host, *svc.Service.Name, *svc.Plugins[0].Enabled)
			if *plugin.Enabled {
				if config, ok := plugin.Config["add"]; ok {
					for _, header := range config.(map[string]interface{})["headers"].([]interface{}) {
						if strings.HasPrefix(header.(string), "Knative-Serving-") {
							log.Info("knative service updated. update knative CR condition and status...")
							if err := UpdateKnativeIngress(ctx, log, svc, kubeConfig, ips, hostname); err != nil {
								return fmt.Errorf("failed to update knative ingress err %v", err)
							}
						}
					}
				}
			}
		}

		switch proto := *svc.Protocol; proto {
		case "tcp":
			if err := UpdateTCPIngress(ctx, log, svc, kiccli, ips); err != nil {
				return fmt.Errorf("failed to update tcp ingress. err %v", err)
			}
		case "udp":
			if err := UpdateUDPIngress(ctx, log, svc, kiccli, ips); err != nil {
				return fmt.Errorf("failed to update udp ingress. err %v", err)
			}

		case "http":
			// if the cluster is on a very old version, we fall back to legacy Ingress support
			// for compatibility with clusters older than v1.19.x.
			// TODO: this can go away once we drop support for Kubernetes older than v1.19
			if kubernetesVersion.Major >= uint64(1) && kubernetesVersion.Minor > uint64(18) {
				if err := UpdateIngress(ctx, log, svc, cli, ips); err != nil {
					return fmt.Errorf("failed to update ingressv1. err %v", err)
				}
			} else {
				if err := UpdateIngressLegacy(ctx, log, svc, cli, ips); err != nil {
					return fmt.Errorf("failed to update ingressv1. err %v", err)
				}
			}
		default:
			log.Info("protocol " + proto + " is not supported")
		}
	}

	return nil
}

func toKnativeLBStatus(coreLBStatus []apiv1.LoadBalancerIngress) []knative.LoadBalancerIngressStatus {
	var res []knative.LoadBalancerIngressStatus
	for _, status := range coreLBStatus {
		res = append(res, knative.LoadBalancerIngressStatus{
			IP:     status.IP,
			Domain: status.Hostname,
		})
	}
	return res
}

// UpdateIngress networking v1 ingress status
func UpdateIngress(
	ctx context.Context,
	logger logr.Logger,
	svc file.FService,
	cli *clientset.Clientset,
	ips []string,
) error {
	for _, route := range svc.Routes {
		routeInf := strings.Split(*((*route).Name), ".")
		namespace := routeInf[0]
		name := routeInf[1]
		log.Debugf("updating status for v1.Ingress route name %s  namespace %s", name, namespace)

		ingCli := cli.NetworkingV1().Ingresses(namespace)
		retry := 0
		for retry < statusUpdateRetry {
			curIng, err := ingCli.Get(ctx, name, metav1.GetOptions{})
			if err != nil || curIng == nil {
				if errors.IsNotFound(err) {
					log.Debugf("failed to retrieve Ingress V1: the object (%s/%s) is gone, status update stopped.", namespace, name)
					return nil
				}

				log.Errorf("failed to fetch Ingress %v/%v due to error: %v. retrying...", namespace, name, err)
				retry++
				time.Sleep(statusUpdateWaitTick)
				continue
			}

			var status []apiv1.LoadBalancerIngress
			sort.SliceStable(status, lessLoadBalancerIngress(status))
			curIPs := curIng.Status.LoadBalancer.Ingress

			status = SliceToStatus(ips)
			if ingressSliceEqual(status, curIPs) {
				log.Debugf("no change in status, update ingress v1 skipped")
				return nil
			}

			curIng.Status.LoadBalancer.Ingress = status

			_, err = ingCli.UpdateStatus(ctx, curIng, metav1.UpdateOptions{})
			if err == nil {
				break
			}
			if errors.IsNotFound(err) {
				log.Debugf("failed to update Ingress V1 status because the object (%s/%s) is gone, status update stopped.", namespace, name)
				return nil
			}
			if errors.IsConflict(err) {
				log.Debugf("failed to update Ingress V1 status because the object (%s/%s) changed: %v retrying...", namespace, name, err)
			} else {
				log.Errorf("failed to update Ingress V1 status. %v. retrying...", err)
			}
			time.Sleep(statusUpdateWaitTick)
			retry++
		}
	}

	log.Debugf("successfully updated networkingv1 Ingress status")
	return nil
}

// UpdateIngressLegacy networking v1beta1 ingress status
// TODO: this can be removed once we no longer support old kubernetes < v1.19
func UpdateIngressLegacy(
	ctx context.Context,
	logger logr.Logger,
	svc file.FService,
	cli *clientset.Clientset,
	ips []string,
) error {
	for _, route := range svc.Routes {
		routeInf := strings.Split(*((*route).Name), ".")
		namespace := routeInf[0]
		name := routeInf[1]
		log.Debugf("updating status for v1beta1.Ingress route name %s  namespace %s", name, namespace)

		ingCli := cli.NetworkingV1beta1().Ingresses(namespace)
		retry := 0
		for retry < statusUpdateRetry {
			curIng, err := ingCli.Get(ctx, name, metav1.GetOptions{})
			if err != nil || curIng == nil {
				if errors.IsNotFound(err) {
					log.Debugf("failed to retrieve Legacy Ingress: the object (%s/%s) is gone, status update stopped.", namespace, name)
					return nil
				}

				log.Errorf("failed to fetch Ingress %v/%v due to error: %v. retrying...", namespace, name, err)
				retry++
				time.Sleep(statusUpdateWaitTick)
				continue
			}

			var status []apiv1.LoadBalancerIngress
			sort.SliceStable(status, lessLoadBalancerIngress(status))
			curIPs := curIng.Status.LoadBalancer.Ingress

			status = SliceToStatus(ips)
			if ingressSliceEqual(status, curIPs) {
				log.Debugf("no change in status, update ingress v1beta1 skipped")
				return nil
			}

			curIng.Status.LoadBalancer.Ingress = status

			_, err = ingCli.UpdateStatus(ctx, curIng, metav1.UpdateOptions{})
			if err == nil {
				break
			}
			if errors.IsNotFound(err) {
				log.Debugf("failed to update Ingress V1 status because the object (%s/%s) is gone, status update stopped.", namespace, name)
				return nil
			}
			if errors.IsConflict(err) {
				log.Debugf("failed to update Ingress V1Beta1 status because the object (%s/%s) changed: %v retrying...", namespace, name, err)
			} else {
				log.Errorf("failed to update Ingress V1Beta1 status. %v. retrying...", err)
			}
			time.Sleep(statusUpdateWaitTick)
			retry++
		}
	}

	log.Debugf("successfully updated networkingv1beta1 Ingress status")
	return nil
}

// UpdateUDPIngress update udp ingress status
func UpdateUDPIngress(ctx context.Context, logger logr.Logger, svc file.FService, kiccli *kicclientset.Clientset,
	ips []string) error {
	for _, route := range svc.Routes {
		routeInf := strings.Split(*((*route).Name), ".")
		namespace := routeInf[0]
		name := routeInf[1]
		log.Debugf("updating UDP ingress route name %s namespace %s", name, namespace)
		ingCli := kiccli.ConfigurationV1beta1().UDPIngresses(namespace)
		retry := 0
		for retry < statusUpdateRetry {
			curIng, err := ingCli.Get(ctx, name, metav1.GetOptions{})
			if err != nil || curIng == nil {
				if errors.IsNotFound(err) {
					log.Debugf("failed to retrieve UDPIngress: the object (%s/%s) is gone, status update stopped.", namespace, name)
					return nil
				}

				log.Errorf("failed to fetch UDP Ingress %v/%v due to error: %v", namespace, name, err)
				time.Sleep(statusUpdateWaitTick)
				retry++
				continue
			}

			var status []apiv1.LoadBalancerIngress
			sort.SliceStable(status, lessLoadBalancerIngress(status))
			curIPs := curIng.Status.LoadBalancer.Ingress

			status = SliceToStatus(ips)
			if ingressSliceEqual(status, curIPs) {
				log.Debugf("no change in status, update udp ingress skipped")
				return nil
			}

			curIng.Status.LoadBalancer.Ingress = status

			_, err = ingCli.UpdateStatus(ctx, curIng, metav1.UpdateOptions{})
			if err == nil {
				break
			}
			if errors.IsNotFound(err) {
				log.Debugf("failed to update UDPIngress status because the object (%s/%s) is gone, status update stopped.", namespace, name)
				return nil
			}
			if errors.IsConflict(err) {
				log.Debugf("failed to update UDPIngress status because the object (%s/%s) changed: %v retrying...", namespace, name, err)
			} else {
				log.Errorf("failed to update UDPIngress status. %v. retrying...", err)
			}
			time.Sleep(statusUpdateWaitTick)
			retry++
		}
	}
	log.Debugf("successfully updated UDPIngress status")
	return nil
}

// UpdateTCPIngress TCP ingress status
func UpdateTCPIngress(ctx context.Context, logger logr.Logger, svc file.FService, kiccli *kicclientset.Clientset,
	ips []string) error {
	for _, route := range svc.Routes {
		routeInf := strings.Split(*((*route).Name), ".")
		namespace := routeInf[0]
		name := routeInf[1]
		log.Debugf("Updating TCP ingress route name %s namespace %s", name, namespace)

		ingCli := kiccli.ConfigurationV1beta1().TCPIngresses(namespace)

		retry := 0
		for retry < statusUpdateRetry {
			curIng, err := ingCli.Get(ctx, name, metav1.GetOptions{})
			if err != nil || curIng == nil {
				if errors.IsNotFound(err) {
					log.Debugf("failed to retrieve TCPIngress: the object (%s/%s) is gone, status update stopped.", namespace, name)
					return nil
				}

				log.Errorf("failed to fetch TCPIngress %v/%v due to error: %v", namespace, name, err)
				time.Sleep(statusUpdateWaitTick)
				retry++
				continue
			}

			curIPs := curIng.Status.LoadBalancer.Ingress
			status := SliceToStatus(ips)
			if ingressSliceEqual(status, curIPs) {
				log.Debugf("no change in status, update tcp ingress skipped")
				return nil
			}

			curIng.Status.LoadBalancer.Ingress = status
			_, err = ingCli.UpdateStatus(ctx, curIng, metav1.UpdateOptions{})
			if err == nil {
				break
			}
			if errors.IsNotFound(err) {
				log.Debugf("failed to update TCPIngress status because the object (%s/%s) is gone, status update stopped.", namespace, name)
				return nil
			}
			if errors.IsConflict(err) {
				log.Debugf("failed to update TCPIngress status because the object (%s/%s) changed: %v retrying...", namespace, name, err)
			} else {
				log.Errorf("failed to update TCPIngress status. %v. retrying...", err)
			}
			time.Sleep(statusUpdateWaitTick)
			retry++
		}
	}

	log.Debugf("Successfully updated TCPIngress status")
	return nil
}

var ingressCondSet = knativeApis.NewLivingConditionSet()

// UpdateKnativeIngress update knative ingress status
func UpdateKnativeIngress(ctx context.Context, logger logr.Logger, svc file.FService, kubeCfg *rest.Config,
	ips []string, hostname string) error {
	for _, route := range svc.Routes {
		routeInf := strings.Split(*((*route).Name), ".")
		namespace := routeInf[0]
		name := routeInf[1]
		log.Debugf("Updating Knative route name %s namespace %s", name, namespace)

		knativeCli, err := knativeversioned.NewForConfig(kubeCfg)
		if err != nil {
			return fmt.Errorf("failed to generate knative client. err %v", err)
		}
		ingClient := knativeCli.NetworkingV1alpha1().Ingresses(namespace)

		retry := 0
		for retry < statusUpdateRetry {
			curIng, err := ingClient.Get(ctx, name, metav1.GetOptions{})
			if err != nil || curIng == nil {
				if errors.IsNotFound(err) {
					log.Debugf("failed to retrieve Knative Ingress: the object (%s/%s) is gone, status update stopped.", namespace, name)
					return nil
				}

				log.Errorf("failed to fetch Knative Ingress %v/%v due to error: %v", namespace, name, err)
				time.Sleep(statusUpdateWaitTick)
				retry++
				continue
			}

			// check if CR current status already updated
			var status []apiv1.LoadBalancerIngress
			sort.SliceStable(status, lessLoadBalancerIngress(status))
			curIPs := toCoreLBStatus(curIng.Status.PublicLoadBalancer)
			status = SliceToStatus(ips)
			if ingressSliceEqual(status, curIPs) &&
				curIng.Status.ObservedGeneration == curIng.GetObjectMeta().GetGeneration() {
				log.Debugf("no change in status, update knative ingress skipped")
				return nil
			}

			// updating current custom status
			lbStatus := toKnativeLBStatus(status)

			for i := 0; i < len(lbStatus); i++ {
				lbStatus[i].DomainInternal = hostname
			}

			curIng.Status.MarkLoadBalancerReady(lbStatus, lbStatus)
			ingressCondSet.Manage(&curIng.Status).MarkTrue(knative.IngressConditionReady)
			ingressCondSet.Manage(&curIng.Status).MarkTrue(knative.IngressConditionNetworkConfigured)
			curIng.Status.ObservedGeneration = curIng.GetObjectMeta().GetGeneration()

			_, err = ingClient.UpdateStatus(ctx, curIng, metav1.UpdateOptions{})
			if err == nil {
				break
			}
			if errors.IsNotFound(err) {
				log.Debugf("failed to update Knative Ingress status because the object (%s/%s) is gone, status update stopped.", namespace, name)
				return nil
			}
			if errors.IsConflict(err) {
				log.Debugf("failed to update Knative Ingress status because the object (%s/%s) changed: %v. retrying...", namespace, name, err)
			} else {
				log.Errorf("failed to update Knative Ingress status. %v. retrying...", err)
			}
			time.Sleep(statusUpdateWaitTick)
			retry++
		}
	}

	log.Debugf("successfully updated knative ingress status")
	return nil
}

// RunningAddresses retrieve cluster loader balance IP or hostaddress using networking
func RunningAddresses(ctx context.Context, kubeCfg *rest.Config, publishService string,
	publishAddresses []string) ([]string, string, error) {
	addrs := []string{}
	if len(publishAddresses) > 0 {
		addrs = append(addrs, publishAddresses...)
		return addrs, "", nil
	}
	namespace, name, err := util.ParseNameNS(publishService)
	if err != nil {
		return nil, "", fmt.Errorf("unable to retrieve service for status: %w", err)
	}

	CoreClient, _ := clientset.NewForConfig(kubeCfg)
	svc, err := CoreClient.CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		log.Infof("running address err %v", err)
		return nil, "", err
	}

	clusterDomain := network.GetClusterDomainName()
	hostname := fmt.Sprintf("%s.%s.svc.%s", name, namespace, clusterDomain)

	//nolint:exhaustive
	switch svc.Spec.Type {
	case apiv1.ServiceTypeLoadBalancer:
		for _, ip := range svc.Status.LoadBalancer.Ingress {
			if ip.IP == "" {
				addrs = append(addrs, ip.Hostname)
			} else {
				addrs = append(addrs, ip.IP)
			}
		}

		addrs = append(addrs, svc.Spec.ExternalIPs...)
		return addrs, hostname, nil
	default:
		return addrs, hostname, nil
	}
}

// SliceToStatus converts a slice of IP and/or hostnames to LoadBalancerIngress
func SliceToStatus(endpoints []string) []apiv1.LoadBalancerIngress {
	lbi := []apiv1.LoadBalancerIngress{}
	for _, ep := range endpoints {
		if net.ParseIP(ep) == nil {
			lbi = append(lbi, apiv1.LoadBalancerIngress{Hostname: ep})
		} else {
			lbi = append(lbi, apiv1.LoadBalancerIngress{IP: ep})
		}
	}

	sort.SliceStable(lbi, func(a, b int) bool {
		return lbi[a].IP < lbi[b].IP
	})

	return lbi
}

// InSlice checks whether a string is present in a list of strings.
func InSlice(e string, arr []string) bool {
	for _, v := range arr {
		if v == e {
			return true
		}
	}

	return false
}

func ingressSliceEqual(lhs, rhs []apiv1.LoadBalancerIngress) bool {
	if len(lhs) != len(rhs) {
		return false
	}

	for i := range lhs {
		if lhs[i].IP != rhs[i].IP {
			return false
		}
		if lhs[i].Hostname != rhs[i].Hostname {
			return false
		}
	}
	return true
}

func toCoreLBStatus(knativeLBStatus *knative.LoadBalancerStatus) []apiv1.LoadBalancerIngress {
	var res []apiv1.LoadBalancerIngress
	if knativeLBStatus == nil {
		return res
	}
	for _, status := range knativeLBStatus.Ingress {
		res = append(res, apiv1.LoadBalancerIngress{
			IP:       status.IP,
			Hostname: status.Domain,
		})
	}
	return res
}

func lessLoadBalancerIngress(addrs []apiv1.LoadBalancerIngress) func(int, int) bool {
	return func(a, b int) bool {
		switch strings.Compare(addrs[a].Hostname, addrs[b].Hostname) {
		case -1:
			return true
		case 1:
			return false
		}
		return addrs[a].IP < addrs[b].IP
	}
}
