package ctrlutils

import (
	"context"
	"fmt"
	"net"
	"sort"
	"strings"

	"github.com/go-logr/logr"
	"github.com/kong/deck/file"
	"github.com/kong/kubernetes-ingress-controller/pkg/sendconfig"
	"github.com/prometheus/common/log"

	"sync"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	knative "knative.dev/networking/pkg/apis/networking/v1alpha1"
	knativeversioned "knative.dev/networking/pkg/client/clientset/versioned"
	knativeApis "knative.dev/pkg/apis"
	"knative.dev/pkg/network"
)

func PullConfigUpdate(kongConfig sendconfig.Kong, log logr.Logger, ctx context.Context, kubeConfig *rest.Config, stopCh <-chan struct{}) {
	log.Info("Launching Customer Resource Update thread.")
	var wg sync.WaitGroup
	for {
		select {
		case updateDone := <-kongConfig.ConfigDone:
			log.Info("receive configuration information. Update ingress status \n%v\n \n", &updateDone)
			wg.Add(1)
			go UpdateIngress(&updateDone, log, ctx, kubeConfig, wg)
		case <-stopCh:
			log.Info("stop status update channel.")
			return
		}
	}
}

// update ingress status according to generated rules and specs
func UpdateIngress(targetContent *file.Content, log logr.Logger, ctx context.Context, kubeconfig *rest.Config, wg sync.WaitGroup) error {
	defer wg.Done()
	for _, svc := range targetContent.Services {
		for _, plugin := range svc.Plugins {
			log.Info("\n service host %s name %s plugin enablement %v\n", *svc.Service.Host, *svc.Service.Name, *svc.Plugins[0].Enabled)
			if *plugin.Enabled == true {
				// filter the plugins (tcp/udp/tls/http) here
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

var ingressCondSet = knativeApis.NewLivingConditionSet()

func UpdateKnativeIngress(ctx context.Context, logger logr.Logger, svc file.FService, kubeCfg *rest.Config) error {
	// retrieve cr name/namespace from route
	var name, namespace string
	routeInf := strings.Split(*(svc.Routes[0].Name), ".")
	namespace = routeInf[0]
	name = routeInf[1]
	log.Info("svc namespace %s name %s", namespace, name)
	if len(namespace) == 0 || len(name) == 0 {
		return fmt.Errorf("configured route information is not completed which should not happen.")
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
	log.Info("retrieving existing CR <%v> ", *curIng)

	// check if CR current status already updated
	var status []apiv1.LoadBalancerIngress
	sort.SliceStable(status, lessLoadBalancerIngress(status))
	curIPs := toCoreLBStatus(curIng.Status.PublicLoadBalancer)
	ips, err := RunningAddresses(ctx, kubeCfg)
	if err != nil {
		return fmt.Errorf("failed to retrieve cluster loadbalancer.")
	}
	status = SliceToStatus(ips)
	if IngressSliceEqual(status, curIPs) &&
		curIng.Status.ObservedGeneration == curIng.GetObjectMeta().GetGeneration() {
		log.Info("no change in status, update skipped")
		return nil
	}

	// updating current custom status
	log.Info("attempting to update Knative Ingress status")
	lbStatus := toKnativeLBStatus(status)
	clusterDomain := network.GetClusterDomainName()
	log.Info("cluster domain %s\n", clusterDomain)
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

func RunningAddresses(ctx context.Context, kubeCfg *rest.Config) ([]string, error) {
	addrs := []string{}

	namespace := "kong-system"
	CoreClient, _ := clientset.NewForConfig(kubeCfg)
	svc, err := CoreClient.CoreV1().Services(namespace).Get(ctx, "ingress-controller-kong-proxy", metav1.GetOptions{})
	if err != nil {
		log.Info("err %v", err)
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

func IngressSliceEqual(lhs, rhs []apiv1.LoadBalancerIngress) bool {
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
