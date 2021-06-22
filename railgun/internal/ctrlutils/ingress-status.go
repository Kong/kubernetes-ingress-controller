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
	fmt.Printf("Launching Customer Resource Update thread.\n")
	for {
		select {
		case updateDone := <-kongConfig.ConfigDone:
			fmt.Printf("receive configuration information. Update ingress status \n%v\n \n", &updateDone)
			go UpdateIngress(&updateDone, log, ctx, kubeConfig)
		case <-stopCh:
			fmt.Printf("stop status update channel.")
			return
		}
	}
}

// update ingress status according to generated rules and specs
func UpdateIngress(targetContent *file.Content, log logr.Logger, ctx context.Context, kubeconfig *rest.Config) error {
	for _, svc := range targetContent.Services {
		fmt.Printf("\n service host %s name %s\n ", *svc.Service.Host, *svc.Service.Name)
		for _, plugin := range svc.Plugins {
			fmt.Printf("\n plugin enablement %v\n", *svc.Plugins[0].Enabled)
			if *plugin.Enabled == true {
				if config, ok := plugin.Config["add"]; ok {
					for _, header := range config.(map[string]interface{})["headers"].([]interface{}) {
						fmt.Printf("header %s", header.(string))
						if strings.HasPrefix(header.(string), "Knative-Serving-") {
							fmt.Printf("\n knative service updated. update knative CR condition and status..\n")
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

	fmt.Printf("svc namespace <%s> name <%s>", namespace, name)
	if len(namespace) == 0 || len(name) == 0 {
		return fmt.Errorf("configured route information is not completed which should not happen.")
	}
	fmt.Println("create knative cli.")
	knativeCli, err := knativeversioned.NewForConfig(kubeCfg)
	if err != nil {
		return fmt.Errorf("failed to generate knative client. err %v", err)
	}

	fmt.Println("create lister to retrieve CR.")
	ingClient := knativeCli.NetworkingV1alpha1().Ingresses(namespace)
	curIng, err := ingClient.Get(ctx, name, metav1.GetOptions{})
	if err != nil || curIng == nil {
		return fmt.Errorf("failed to fetch Knative Ingress %v/%v: %w", namespace, name, err)
	}
	fmt.Printf("\n able to retrieve existing CR <%v> \n", *curIng)

	// check if CR current status already updated
	var status []apiv1.LoadBalancerIngress
	sort.SliceStable(status, lessLoadBalancerIngress(status))
	curIPs := toCoreLBStatus(curIng.Status.PublicLoadBalancer)
	ips, err := RunningAddresses(ctx, kubeCfg)
	if err != nil {
		return fmt.Errorf("failed to retrieve cluster loadbalancer.")
	}

	fmt.Printf("\n king ip %v cluster ip %v \n", curIPs, ips)
	status = SliceToStatus(ips)
	if IngressSliceEqual(status, curIPs) &&
		curIng.Status.ObservedGeneration == curIng.GetObjectMeta().GetGeneration() {
		fmt.Printf("no change in status, update skipped")
		return nil
	}
	fmt.Printf("\n convert to corev1 lb status %v \n", status)

	// updating current custom status
	fmt.Printf("attempting to update Knative Ingress status")
	lbStatus := toKnativeLBStatus(status)

	//
	//proxyNS, svcName, err := util.ParseNameNS("kong-system/ingress-controller-kong-proxy")
	clusterDomain := network.GetClusterDomainName()
	fmt.Printf("cluster domain %s\n", clusterDomain)
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
		fmt.Printf("err %v", err)
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
