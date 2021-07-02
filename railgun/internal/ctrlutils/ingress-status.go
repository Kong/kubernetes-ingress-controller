package ctrlutils

import (
	"context"
	"fmt"
	"net"
	"sort"
	"strings"
	"sync"

	"github.com/go-logr/logr"
	"github.com/kong/deck/file"
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
	kicclientset "github.com/kong/kubernetes-ingress-controller/railgun/pkg/clientset"
)

// dedicated function that process ingress/customer resource status update after configuration is updated within kong.
func PullConfigUpdate(kongConfig sendconfig.Kong, log logr.Logger, ctx context.Context, kubeConfig *rest.Config, stopCh <-chan struct{}) {
	log.Info("Launching Ingress Status Update Thread.")
	var wg sync.WaitGroup
	for {
		select {
		case updateDone := <-kongConfig.ConfigDone:
			log.Info("receive configuration information. Update ingress status %v \n", updateDone)
			wg.Add(1)
			go UpdateIngress(&updateDone, log, ctx, kubeConfig, &wg)
			wg.Wait()
		case <-stopCh:
			log.Info("stop status update channel.")
			return
		}
	}
}

// update ingress status according to generated rules and specs
func UpdateIngress(targetContent *file.Content, log logr.Logger, ctx context.Context, kubeconfig *rest.Config, wg *sync.WaitGroup) error {
	defer wg.Done()

	for _, svc := range targetContent.Services {

		for _, plugin := range svc.Plugins {
			log.Info("\n service host %s name %s plugin enablement %v\n", *svc.Service.Host, *svc.Service.Name, *svc.Plugins[0].Enabled)
			if *plugin.Enabled == true {
				if config, ok := plugin.Config["add"]; ok {
					for _, header := range config.(map[string]interface{})["headers"].([]interface{}) {
						if strings.HasPrefix(header.(string), "Knative-Serving-") {
							log.Info("knative service updated. update knative CR condition and status...")
							err := UpdateKnativeIngress(ctx, log, svc, kubeconfig)
							return fmt.Errorf("failed to update knative ingress err %v", err)
						}
					}
				}
			}
		}

		switch proto := *svc.Protocol; proto {
		case "tcp":
			err := UpdateTCPIngress(ctx, log, svc, kubeconfig)
			return fmt.Errorf("failed to update tcp ingress. err %v", err)
		case "udp":
			err := UpdateUDPIngress(ctx, log, svc, kubeconfig)
			return fmt.Errorf("failed to update udp ingress. err %v", err)
		case "http":
			err := UpdateIngressV1(ctx, log, svc, kubeconfig)
			return fmt.Errorf("failed to update ingressv1. err %v", err)
		default:
			log.Info("other 3rd party ingress not supported yet.")
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

func retrieveNSAndNM(svc file.FService) (string, string, error) {
	var name, namespace string
	routeInf := strings.Split(*(svc.Routes[0].Name), ".")
	namespace = routeInf[0]
	name = routeInf[1]
	if len(namespace) == 0 || len(name) == 0 {
		return "", "", fmt.Errorf("configured route information is not completed which should not happen.")
	}
	return namespace, name, nil
}

// update networking v1 ingress status
func UpdateIngressV1(ctx context.Context, logger logr.Logger, svc file.FService, kubeCfg *rest.Config) error {
	namespace, name, err := retrieveNSAndNM(svc)
	if err != nil {
		return fmt.Errorf("failed to generate UDP client. err %v", err)
	}

	cli, err := clientset.NewForConfig(kubeCfg)
	if err != nil {
		return fmt.Errorf("failed to generate UDP client. err %v", err)
	}
	ingCli := cli.NetworkingV1().Ingresses(namespace)
	cli.NetworkingV1beta1()
	if err != nil {
		return fmt.Errorf("failed to generate UDP client. err %v", err)
	}
	curIng, err := ingCli.Get(ctx, name, metav1.GetOptions{})
	if err != nil || curIng == nil {
		return fmt.Errorf("failed to fetch Ingress %v/%v: %w", namespace, name, err)
	}

	var status []apiv1.LoadBalancerIngress
	sort.SliceStable(status, lessLoadBalancerIngress(status))
	curIPs := curIng.Status.LoadBalancer.Ingress
	ips, err := RunningAddresses(ctx, kubeCfg)
	if err != nil {
		return fmt.Errorf("failed to retrieve cluster loadbalancer.")
	}

	status = SliceToStatus(ips)
	if ingressSliceEqual(status, curIPs) {
		logger.Info("no change in status, update skipped")
		return nil
	}

	curIng.Status.LoadBalancer.Ingress = status
	_, err = ingCli.UpdateStatus(ctx, curIng, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update UDPIngress status: %v", err)
	} else {
		return fmt.Errorf("ingress_status successfully updated UDPIngress status")
	}

}

// updagte udp ingress status
func UpdateUDPIngress(ctx context.Context, logger logr.Logger, svc file.FService, kubeCfg *rest.Config) error {
	namespace, name, err := retrieveNSAndNM(svc)
	if err != nil {
		return fmt.Errorf("failed to generate UDP client. err %v", err)
	}

	cli, err := kicclientset.NewForConfig(kubeCfg)
	if err != nil {
		return fmt.Errorf("failed to generate UDP client. err %v", err)
	}
	ingCli := cli.ConfigurationV1beta1().UDPIngresses(namespace)
	curIng, err := ingCli.Get(ctx, name, metav1.GetOptions{})
	if err != nil || curIng == nil {
		return fmt.Errorf("failed to fetch UDP Ingress %v/%v: %w", namespace, name, err)
	}

	var status []apiv1.LoadBalancerIngress
	sort.SliceStable(status, lessLoadBalancerIngress(status))
	curIPs := curIng.Status.LoadBalancer.Ingress
	ips, err := RunningAddresses(ctx, kubeCfg)
	if err != nil {
		return fmt.Errorf("failed to retrieve cluster loadbalancer.")
	}

	status = SliceToStatus(ips)
	if ingressSliceEqual(status, curIPs) {
		log.Debugf("no change in status, update skipped")
		return nil
	}

	curIng.Status.LoadBalancer.Ingress = status
	_, err = ingCli.UpdateStatus(ctx, curIng, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update UDPIngress status: %v", err)
	} else {
		return fmt.Errorf("ingress_status successfully updated UDPIngress status")
	}
}

// update TCP ingress status
func UpdateTCPIngress(ctx context.Context, logger logr.Logger, svc file.FService, kubeCfg *rest.Config) error {
	namespace, name, err := retrieveNSAndNM(svc)
	if err != nil {
		return fmt.Errorf("failed to generate UDP client. err %v", err)
	}

	cli, err := kicclientset.NewForConfig(kubeCfg)
	if err != nil {
		return fmt.Errorf("failed to generate TCP client. err %v", err)
	}
	ingCli := cli.ConfigurationV1beta1().TCPIngresses(namespace)
	curIng, err := ingCli.Get(ctx, name, metav1.GetOptions{})
	if err != nil || curIng == nil {
		return fmt.Errorf("failed to fetch TCP Ingress %v/%v: %w", namespace, name, err)
	}

	var status []apiv1.LoadBalancerIngress
	sort.SliceStable(status, lessLoadBalancerIngress(status))
	curIPs := curIng.Status.LoadBalancer.Ingress
	ips, err := RunningAddresses(ctx, kubeCfg)
	if err != nil {
		return fmt.Errorf("failed to retrieve cluster loadbalancer.")
	}

	status = SliceToStatus(ips)
	if ingressSliceEqual(status, curIPs) {
		log.Debugf("no change in status, update skipped")
		return nil
	}

	curIng.Status.LoadBalancer.Ingress = status
	_, err = ingCli.UpdateStatus(ctx, curIng, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update TCPIngress status: %v", err)
	} else {
		return fmt.Errorf("ingress_status successfully updated TCPIngress status")
	}
}

var ingressCondSet = knativeApis.NewLivingConditionSet()

// update knative ingress status
func UpdateKnativeIngress(ctx context.Context, logger logr.Logger, svc file.FService, kubeCfg *rest.Config) error {
	namespace, name, err := retrieveNSAndNM(svc)
	if err != nil {
		return fmt.Errorf("failed to generate UDP client. err %v", err)
	}

	knativeCli, err := knativeversioned.NewForConfig(kubeCfg)
	if err != nil {
		return fmt.Errorf("failed to generate knative client. err %v", err)
	}
	ingClient := knativeCli.NetworkingV1alpha1().Ingresses(namespace)
	curIng, err := ingClient.Get(ctx, name, metav1.GetOptions{})
	if err != nil || curIng == nil {
		return fmt.Errorf("failed to fetch Knative Ingress %v/%v: %w", namespace, name, err)
	}

	// check if CR current status already updated
	var status []apiv1.LoadBalancerIngress
	sort.SliceStable(status, lessLoadBalancerIngress(status))
	curIPs := toCoreLBStatus(curIng.Status.PublicLoadBalancer)
	ips, err := RunningAddresses(ctx, kubeCfg)
	if err != nil {
		return fmt.Errorf("failed to retrieve cluster loadbalancer.")
	}
	status = SliceToStatus(ips)
	if ingressSliceEqual(status, curIPs) &&
		curIng.Status.ObservedGeneration == curIng.GetObjectMeta().GetGeneration() {
		log.Debugf("no change in status, update skipped")
		return nil
	}

	// updating current custom status
	lbStatus := toKnativeLBStatus(status)
	clusterDomain := network.GetClusterDomainName()
	if err != nil {
		return err
	}

	for i := 0; i < len(lbStatus); i++ {
		lbStatus[i].DomainInternal = fmt.Sprintf("%s.%s.svc.%s",
			"ingress-controller-kong-proxy", "kong-system", clusterDomain)
	}

	curIng.Status.MarkLoadBalancerReady(lbStatus, lbStatus)
	ingressCondSet.Manage(&curIng.Status).MarkTrue(knative.IngressConditionReady)
	ingressCondSet.Manage(&curIng.Status).MarkTrue(knative.IngressConditionNetworkConfigured)
	curIng.Status.ObservedGeneration = curIng.GetObjectMeta().GetGeneration()

	_, err = ingClient.UpdateStatus(ctx, curIng, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update ingress status: %v", err)
	} else {
		logger.Info("successfully updated ingress status")
	}
	return nil
}

// retrieve cluster loader balance IP or hostaddress using networking
func RunningAddresses(ctx context.Context, kubeCfg *rest.Config) ([]string, error) {
	addrs := []string{}
	// loading ns from environment variable when https://github.com/Kong/kubernetes-ingress-controller/issues/1480
	// is resolved
	namespace := "kong-system"
	CoreClient, _ := clientset.NewForConfig(kubeCfg)
	svc, err := CoreClient.CoreV1().Services(namespace).Get(ctx, "ingress-controller-kong-proxy", metav1.GetOptions{})
	if err != nil {
		log.Infof("running address err %v", err)
		return nil, err
	}

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
		return addrs, nil
	default:
		return addrs, nil
	}
}

// sliceToStatus converts a slice of IP and/or hostnames to LoadBalancerIngress
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
