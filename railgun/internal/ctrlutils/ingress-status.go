package ctrlutils

import (
	"context"
	"fmt"
	"net"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"github.com/kong/deck/file"
	"github.com/mitchellh/hashstructure/v2"
	"github.com/prometheus/common/log"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	knative "knative.dev/networking/pkg/apis/networking/v1alpha1"
	knativeversioned "knative.dev/networking/pkg/client/clientset/versioned"
	knativeApis "knative.dev/pkg/apis"
	"knative.dev/pkg/network"

	"github.com/kong/kubernetes-ingress-controller/pkg/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/pkg/util"
	kicclientset "github.com/kong/kubernetes-ingress-controller/railgun/pkg/clientset"
)

const (
	statusUpdateRetry = 3
)

// PullConfigUpdate is a dedicated function that process ingress/customer resource status update after configuration is updated within kong.
func PullConfigUpdate(ctx context.Context, kongConfig sendconfig.Kong, log logr.Logger, kubeConfig *rest.Config,
	publishService string, publishAddresses []string) {
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

	kiccli, err := kicclientset.NewForConfig(kubeConfig)
	if err != nil {
		log.Error(err, "failed to create kong ingress client.")
		return
	}

	if err := util.InitCache(); err != nil {
		log.Error(err, "failed to initialize memory cache.")
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
				if err := UpdateIngress(ctx, &updateDone, log, cli, kiccli, &wg, ips, hostname, kubeConfig); err != nil {
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

// UpdateIngress update ingress status according to generated rules and specs
func UpdateIngress(ctx context.Context, targetContent *file.Content, log logr.Logger, cli *clientset.Clientset,
	kiccli *kicclientset.Clientset,
	wg *sync.WaitGroup, ips []string, hostname string,
	kubeConfig *rest.Config) error {
	defer wg.Done()

	for _, svc := range targetContent.Services {
		for _, plugin := range svc.Plugins {
			log.V(5).Info("service host %s name %s plugin enablement %v", *svc.Service.Host, *svc.Service.Name, *svc.Plugins[0].Enabled)
			if *plugin.Enabled {
				if config, ok := plugin.Config["add"]; ok {
					for _, header := range config.(map[string]interface{})["headers"].([]interface{}) {
						if strings.HasPrefix(header.(string), "Knative-Serving-") {
							log.Info("knative service updated. update knative CR condition and status...")
							if err := UpdateKnativeIngress(ctx, log, svc, kubeConfig, ips, hostname); err != nil {
								panic(err)
							}
						}
					}
				}
			}
		}

		switch proto := *svc.Protocol; proto {
		case "tcp":
			if err := UpdateTCPIngress(ctx, log, svc, kiccli, ips); err != nil {
				panic(err)
			}
		case "udp":
			if err := UpdateUDPIngress(ctx, log, svc, kiccli, ips); err != nil {
				panic(err)
			}
		case "http":
			if err := UpdateIngressV1(ctx, log, svc, cli, ips); err != nil {
				panic(err)
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

// UpdateIngressV1 networking v1 ingress status
func UpdateIngressV1(ctx context.Context, logger logr.Logger, svc file.FService, cli *clientset.Clientset,
	ips []string) error {
	for _, route := range svc.Routes {
		routeInf := strings.Split(*((*route).Name), ".")
		namespace := routeInf[0]
		name := routeInf[1]
		log.Infof("updating status for v1.Ingress route: name %s namespace %s", name, namespace)

		ingresKey := fmt.Sprintf("%s-%s", namespace, name)
		log.Info("Updating Networking V1 Ingress " + ingresKey + " status.")

		ingCli := cli.NetworkingV1().Ingresses(namespace)
		retry := 0
		for retry < statusUpdateRetry {
			curIng, err := ingCli.Get(ctx, name, metav1.GetOptions{})
			if err != nil || curIng == nil {
				log.Errorf("failed to fetch Ingress %v/%v: %v. retrying...", namespace, name, err)
				retry++
				time.Sleep(time.Second)
				continue
			}

			existingHash, err := util.GetValue(ingresKey)
			if err != nil {
				log.Infof("v1ingress %s not processed yet. ", ingresKey)
			} else {
				curHash, err := hashstructure.Hash(*curIng, hashstructure.FormatV2, nil)
				if err != nil {
					panic(err)
				}
				if existingHash == curHash {
					log.Info("no change in v1 ingress %s, skip updating.", ingresKey)
					return nil
				}
			}

			status := SliceToStatus(ips)
			curIng.Status.LoadBalancer.Ingress = status
			configuredV1Ingress, err := ingCli.UpdateStatus(ctx, curIng, metav1.UpdateOptions{})
			if err == nil {
				log.Infof("configured v1ingress %v ", *configuredV1Ingress)
				hash, err := hashstructure.Hash(*configuredV1Ingress, hashstructure.FormatV2, nil)
				if err != nil {
					panic(err)
				}
				if err = util.SetValue(ingresKey, hash); err != nil {
					log.Errorf("failed to persist ingress v1 %s status into mem cache. err %v", ingresKey, err)
					time.Sleep(time.Second)
					retry++
					continue
				}
				log.Info("successfully updated " + ingresKey + " status")
				break
			} else {
				time.Sleep(time.Second)
				retry++
				continue
			}
		}
	}
	log.Info("successfully updated networkingv1 ingresses status")
	return nil
}

// UpdateUDPIngress update udp ingress status
func UpdateUDPIngress(ctx context.Context, logger logr.Logger, svc file.FService, kiccli *kicclientset.Clientset,
	ips []string) error {
	for _, route := range svc.Routes {
		routeInf := strings.Split(*((*route).Name), ".")
		namespace := routeInf[0]
		name := routeInf[1]

		ingresKey := fmt.Sprintf("%s-%s", namespace, name)
		log.Infof("updating UDP ingress route namespace-name %s", ingresKey)
		ingCli := kiccli.ConfigurationV1beta1().UDPIngresses(namespace)
		retry := 0
		for retry < statusUpdateRetry {
			curIng, err := ingCli.Get(ctx, name, metav1.GetOptions{})
			if err != nil || curIng == nil {
				log.Errorf("failed to fetch UDP Ingress %v/%v: %v", namespace, name, err)
				time.Sleep(time.Second)
				retry++
				continue
			}

			existingHash, err := util.GetValue(ingresKey)
			if err != nil {
				log.Infof("udp ingress %s not processed yet. ", ingresKey)
			} else {
				curHash, err := hashstructure.Hash(*curIng, hashstructure.FormatV2, nil)
				if err != nil {
					panic(err)
				}
				if existingHash == curHash {
					log.Info("no change in udp ingress %s, skip updating.", ingresKey)
					return nil
				}
			}

			status := SliceToStatus(ips)
			curIng.Status.LoadBalancer.Ingress = status
			configuredUdpIngress, err := ingCli.UpdateStatus(ctx, curIng, metav1.UpdateOptions{})
			if err == nil {
				hash, err := hashstructure.Hash(*configuredUdpIngress, hashstructure.FormatV2, nil)
				if err != nil {
					panic(err)
				}
				if err = util.SetValue(ingresKey, hash); err != nil {
					log.Errorf("failed to persist udp ingress %s status into mem cache. err %v", ingresKey, err)
					time.Sleep(time.Second)
					retry++
					continue
				}
				log.Info("successfully updated " + ingresKey + " status.")
				break
			} else {
				log.Errorf("failed to update UDPIngress status: %v. retry...", err)
				time.Sleep(time.Second)
				retry++
				continue
			}
		}
	}
	log.Info("successfully update udp ingresses status.")
	return nil
}

// UpdateTCPIngress TCP ingress status
func UpdateTCPIngress(ctx context.Context, logger logr.Logger, svc file.FService, kiccli *kicclientset.Clientset,
	ips []string) error {

	for _, route := range svc.Routes {
		routeInf := strings.Split(*((*route).Name), ".")
		namespace := routeInf[0]
		name := routeInf[1]
		log.Infof("Updating TCP ingress route name %s namespace %s", name, namespace)

		ingresKey := fmt.Sprintf("%s-%s", namespace, name)
		log.Info("Updating TCP Ingress " + ingresKey + " status.")

		retry := 0
		for retry < statusUpdateRetry {
			ingCli := kiccli.ConfigurationV1beta1().TCPIngresses(namespace)
			curIng, err := ingCli.Get(ctx, name, metav1.GetOptions{})
			if err != nil || curIng == nil {
				log.Error("failed to fetch TCP Ingress %v/%v: %w", namespace, name, err)
				time.Sleep(time.Second)
				retry++
				continue
			}

			existingHash, err := util.GetValue(ingresKey)
			if err != nil {
				log.Infof("tcp ingress %s not processed yet. ", ingresKey)
			} else {
				curHash, err := hashstructure.Hash(*curIng, hashstructure.FormatV2, nil)
				if err != nil {
					panic(err)
				}
				if existingHash == curHash {
					log.Info("no change in tcp ingress %s, skip updating.", ingresKey)
					return nil
				}
			}

			status := SliceToStatus(ips)
			curIng.Status.LoadBalancer.Ingress = status
			configuredTCPIngress, err := ingCli.UpdateStatus(ctx, curIng, metav1.UpdateOptions{})
			if err != nil {
				log.Errorf("failed to update TCPIngress status: %v", err)
				time.Sleep(time.Second)
				retry++
				continue
			}
			log.Infof("configured tcp ingress %v", *configuredTCPIngress)
			hash, err := hashstructure.Hash(*configuredTCPIngress, hashstructure.FormatV2, nil)
			if err != nil {
				log.Errorf("failed to generate ingress %s hash. err %v", ingresKey, err)
				time.Sleep(time.Second)
				retry++
				continue
			}
			if err = util.SetValue(ingresKey, hash); err != nil {
				log.Errorf("failed to persist ingress %s status into cache. err %v", ingresKey, err)
				time.Sleep(time.Second)
				retry++
				continue
			}
			log.Info("Successfully updated TCPIngress " + ingresKey + " status.")
			break
		}

	}

	log.Info("Successfully updated TCP ingresses status.")
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
		ingresKey := fmt.Sprintf("%s-%s", namespace, name)
		log.Infof("Updating Knative route namespace-name %s", ingresKey)

		knativeCli, err := knativeversioned.NewForConfig(kubeCfg)
		if err != nil {
			return fmt.Errorf("failed to generate knative client. err %v", err)
		}
		ingClient := knativeCli.NetworkingV1alpha1().Ingresses(namespace)

		retry := 0
		for retry < statusUpdateRetry {
			curIng, err := ingClient.Get(ctx, name, metav1.GetOptions{})
			if err != nil || curIng == nil {
				log.Errorf("failed to fetch Knative Ingress %v/%v: %v", namespace, name, err)
				time.Sleep(time.Second)
				retry++
				continue
			}
			ingClient := knativeCli.NetworkingV1alpha1().Ingresses(namespace)

			existingHash, err := util.GetValue(ingresKey)
			if err != nil {
				log.Infof("knative ingress %s not processed yet. ", ingresKey)
			} else {
				curHash, err := hashstructure.Hash(*curIng, hashstructure.FormatV2, nil)
				if err != nil {
					panic(err)
				}
				if existingHash == curHash {
					log.Info("no change in knative ingress %s, skip updating.", ingresKey)
					return nil
				}
			}

			// check if CR current status already updated
			status := SliceToStatus(ips)

			// updating current custom status
			lbStatus := toKnativeLBStatus(status)

			for i := 0; i < len(lbStatus); i++ {
				lbStatus[i].DomainInternal = hostname
			}

			curIng.Status.MarkLoadBalancerReady(lbStatus, lbStatus)
			ingressCondSet.Manage(&curIng.Status).MarkTrue(knative.IngressConditionReady)
			ingressCondSet.Manage(&curIng.Status).MarkTrue(knative.IngressConditionNetworkConfigured)
			curIng.Status.ObservedGeneration = curIng.GetObjectMeta().GetGeneration()

			configuredKnativeIngress, err := ingClient.UpdateStatus(ctx, curIng, metav1.UpdateOptions{})
			if err == nil {
				hash, err := hashstructure.Hash(*configuredKnativeIngress, hashstructure.FormatV2, nil)
				if err != nil {
					panic(err)
				}
				if err = util.SetValue(ingresKey, hash); err != nil {
					log.Errorf("failed to persist knative ingress %s status into cache. err %v", ingresKey, err)
					time.Sleep(time.Second)
					retry++
					continue
				}
				logger.Info("successfully updated knative ingress" + ingresKey + " status")
				break
			} else {
				log.Errorf("failed to update Knative Ingress %v/%v: %v", namespace, name, err)
				time.Sleep(time.Second)
				retry++
				continue
			}
		}
	}
	logger.Info("successfully updated knative ingresses status")
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
