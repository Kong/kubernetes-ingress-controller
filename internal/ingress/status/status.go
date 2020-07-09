/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package status

import (
	"context"
	"fmt"
	"net"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/kong/kubernetes-ingress-controller/internal/ingress/task"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/utils"
	configurationv1beta1 "github.com/kong/kubernetes-ingress-controller/pkg/apis/configuration/v1beta1"
	configClientSet "github.com/kong/kubernetes-ingress-controller/pkg/client/configuration/clientset/versioned"
	pool "gopkg.in/go-playground/pool.v3"
	apiv1 "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/wait"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/leaderelection"
	knativeApis "knative.dev/pkg/apis"
	"knative.dev/pkg/network"
	knative "knative.dev/serving/pkg/apis/networking/v1alpha1"
	knativeClientSet "knative.dev/serving/pkg/client/clientset/versioned"
)

const (
	updateInterval = 60 * time.Second
)

// Sync ...
type Sync interface {
	Run()
	Shutdown(l bool)

	Callbacks() leaderelection.LeaderCallbacks
}

type ingressLister interface {
	// ListIngresses returns the list of Ingresses
	ListIngresses() []*networking.Ingress
	ListTCPIngresses() ([]*configurationv1beta1.TCPIngress, error)
	ListKnativeIngresses() ([]*knative.Ingress, error)
}

// Config ...
type Config struct {
	CoreClient       clientset.Interface
	KongConfigClient configClientSet.Interface
	KnativeClient    knativeClientSet.Interface

	OnStartedLeading func()

	PublishService string

	PublishStatusAddress string

	UpdateStatusOnShutdown bool

	IngressLister ingressLister

	UseNetworkingV1beta1 bool

	Logger logrus.FieldLogger
}

// statusSync keeps the status IP in each Ingress rule updated executing a periodic check
// in all the defined rules. To simplify the process leader election is used so the update
// is executed only in one node (Ingress controllers can be scaled to more than one)
// If the controller is running with the flag --publish-service (with a valid service)
// the IP address behind the service is used, if it is running with the flag
// --publish-status-address, the address specified in the flag is used, if neither of the
// two flags are set, the source is the IP/s of the node/s
type statusSync struct {
	Config
	// pod contains runtime information about this pod
	pod *utils.PodInfo

	// workqueue used to keep in sync the status IP/s
	// in the Ingress rules
	syncQueue *task.Queue
	callbacks leaderelection.LeaderCallbacks

	Logger logrus.FieldLogger
}

// Run starts the loop to keep the status in sync.
func (s statusSync) Run() {
}

func (s statusSync) Callbacks() leaderelection.LeaderCallbacks {
	return s.callbacks
}

// Shutdown stop the sync.
//
// When this instance is the leader it will remove its current IP address if no other instances are
// running.
func (s statusSync) Shutdown(isLeader bool) {
	go s.syncQueue.Shutdown()
	if !isLeader {
		return
	}
	logger := s.Logger.WithField("context", "shutdown")

	if !s.UpdateStatusOnShutdown {
		logger.WithField("UpdateStatusOnShutDown",
			s.UpdateStatusOnShutdown).Infof("update of ingress status skipped")
		return
	}

	// Remove our IP address from all Ingress "status" subresources that mention it.
	s.Logger.Infof("updating status of Ingress rules (remove)")

	addrs, err := s.runningAddresses()
	if err != nil {
		logger.Errorf("failed to fetch IP addresses of running ingress controllers: %v", err)
		return
	}

	if len(addrs) > 1 {
		// Leave the job to the next leader.
		logger.Infof("leaving status update for next leader (%d other candidates)", len(addrs)-1)
		return
	}

	if s.isRunningMultiplePods() {
		logger.Infof("skipping Ingress status update (multiple pods running; another one will be elected as leader)")
		return
	}

	logger.Infof("removing address from Ingress status (%v)", addrs)
	s.updateStatus([]apiv1.LoadBalancerIngress{})
}

func (s *statusSync) sync(key interface{}) error {
	if s.syncQueue.IsShuttingDown() {
		s.Logger.Debugf("shutdown in progress, skipping ingress status update")
		return nil
	}

	addrs, err := s.runningAddresses()
	if err != nil {
		return err
	}
	s.updateStatus(sliceToStatus(addrs))

	return nil
}

func (s statusSync) keyfunc(input interface{}) (interface{}, error) {
	return input, nil
}

// NewStatusSyncer returns a new Sync instance
func NewStatusSyncer(config Config) (Sync, error) {
	pod, err := utils.GetPodDetails(config.CoreClient)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pod information: %v", err)
	}

	st := statusSync{
		pod:    pod,
		Config: config,
		Logger: config.Logger,
	}
	st.syncQueue = task.NewCustomTaskQueue(st.sync, st.keyfunc,
		config.Logger.WithField("component", "status-queue"))

	st.callbacks = leaderelection.LeaderCallbacks{
		OnStartedLeading: func(ctx context.Context) {
			st.Logger.Infof("started leading")
			if st.Config.OnStartedLeading != nil {
				st.Config.OnStartedLeading()
			}
			go st.syncQueue.Run(time.Second, ctx.Done())
			err := wait.PollUntil(updateInterval, func() (bool, error) {
				// send a dummy object to the queue to force a sync
				st.syncQueue.Enqueue("sync status")
				return false, nil
			}, ctx.Done())
			if err != nil {
				st.Logger.Errorf("polling failed :%v", err)
			}
		},
		OnStoppedLeading: func() {
			st.Logger.Infof("stopped leading")
		},
		OnNewLeader: func(identity string) {
			st.Logger.WithField("leader", identity).Infof("leadership changed")
		},
	}

	return st, nil
}

// runningAddresses returns a list of IP addresses and/or FQDN where the
// ingress controller is currently running
func (s *statusSync) runningAddresses() ([]string, error) {
	addrs := []string{}

	if s.PublishStatusAddress != "" {
		addrs = append(addrs, s.PublishStatusAddress)
		return addrs, nil
	}

	ns, name, _ := utils.ParseNameNS(s.PublishService)
	svc, err := s.CoreClient.CoreV1().Services(ns).Get(name, metav1.GetOptions{})
	if err != nil {
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
		// get information about all the pods running the ingress controller
		pods, err := s.CoreClient.CoreV1().Pods(s.pod.Namespace).List(metav1.ListOptions{
			LabelSelector: labels.SelectorFromSet(s.pod.Labels).String(),
		})
		if err != nil {
			return nil, err
		}

		for _, pod := range pods.Items {
			// only Running pods are valid
			if pod.Status.Phase != apiv1.PodRunning {
				continue
			}

			name := utils.GetNodeIPOrName(s.CoreClient, pod.Spec.NodeName)
			if !inSlice(name, addrs) {
				addrs = append(addrs, name)
			}
		}

		return addrs, nil
	}
}

func inSlice(e string, arr []string) bool {
	for _, v := range arr {
		if v == e {
			return true
		}
	}

	return false
}

func (s *statusSync) isRunningMultiplePods() bool {
	pods, err := s.CoreClient.CoreV1().Pods(s.pod.Namespace).List(metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(s.pod.Labels).String(),
	})
	if err != nil {
		return false
	}

	return len(pods.Items) > 1
}

// sliceToStatus converts a slice of IP and/or hostnames to LoadBalancerIngress
func sliceToStatus(endpoints []string) []apiv1.LoadBalancerIngress {
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

// updateStatus changes the status information of Ingress rules
func (s *statusSync) updateStatus(newIngressPoint []apiv1.LoadBalancerIngress) {
	ings := s.IngressLister.ListIngresses()
	tcpIngresses, err := s.IngressLister.ListTCPIngresses()
	if err != nil {
		s.Logger.Errorf("failed to list TPCIngresses: %v", err)
	}
	knativeIngresses, err := s.IngressLister.ListKnativeIngresses()
	if err != nil {
		s.Logger.Errorf("failed to list Knative Ingresses: %v", err)
	}

	p := pool.NewLimited(10)
	defer p.Close()

	batch := p.Batch()

	for _, ing := range ings {
		batch.Queue(s.runUpdate(ing, newIngressPoint, s.CoreClient))
	}
	for _, ing := range tcpIngresses {
		batch.Queue(s.runUpdateTCPIngress(ing, newIngressPoint, s.KongConfigClient))
	}
	for _, ing := range knativeIngresses {
		batch.Queue(s.runUpdateKnativeIngress(ing, newIngressPoint, s.KnativeClient))
	}

	batch.QueueComplete()
	batch.WaitAll()
}

func (s *statusSync) runUpdate(ing *networking.Ingress, status []apiv1.LoadBalancerIngress,
	client clientset.Interface) pool.WorkFunc {
	return func(wu pool.WorkUnit) (interface{}, error) {
		if wu.IsCancelled() {
			return nil, nil
		}

		logger := s.Logger.WithFields(logrus.Fields{
			"ingress_namespace": ing.Namespace,
			"ingress_name":      ing.Name,
		})
		sort.SliceStable(status, lessLoadBalancerIngress(status))

		curIPs := ing.Status.LoadBalancer.Ingress
		sort.SliceStable(curIPs, lessLoadBalancerIngress(curIPs))

		if ingressSliceEqual(status, curIPs) {
			logger.Debugf("no change in status, update skipped")
			return true, nil
		}

		if s.UseNetworkingV1beta1 {
			ingClient := client.NetworkingV1beta1().Ingresses(ing.Namespace)

			currIng, err := ingClient.Get(ing.Name, metav1.GetOptions{})
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("failed to fetch Ingress %v/%v", ing.Namespace, ing.Name))
			}

			logger.WithField("ingress_status", status).Debugf("attempting to update ingress status")
			currIng.Status.LoadBalancer.Ingress = status
			_, err = ingClient.UpdateStatus(currIng)
			if err != nil {
				// TODO return this error?
				logger.Errorf("failed to update ingress status: %v", err)
			} else {
				logger.WithField("ingress_status", status).Debugf("successfully updated ingress status")
			}

		} else {
			ingClient := client.ExtensionsV1beta1().Ingresses(ing.Namespace)

			currIng, err := ingClient.Get(ing.Name, metav1.GetOptions{})
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("failed to fetch Ingress %v/%v", ing.Namespace, ing.Name))
			}

			logger.WithField("ingress_status", status).Debugf("attempting to update ingress status")
			currIng.Status.LoadBalancer.Ingress = status
			_, err = ingClient.UpdateStatus(currIng)
			if err != nil {
				// TODO return this error?
				logger.Errorf("failed to update ingress status: %v", err)
			} else {
				logger.WithField("ingress_status", status).Debugf("successfully updated ingress status")
			}
		}
		return true, nil
	}
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

func (s *statusSync) runUpdateKnativeIngress(ing *knative.Ingress,
	status []apiv1.LoadBalancerIngress,
	client knativeClientSet.Interface) pool.WorkFunc {
	return func(wu pool.WorkUnit) (interface{}, error) {
		if wu.IsCancelled() {
			return nil, nil
		}

		logger := s.Logger.WithFields(logrus.Fields{
			"ingress_namespace": ing.Namespace,
			"ingress_name":      ing.Name,
		})
		sort.SliceStable(status, lessLoadBalancerIngress(status))
		curIPs := toCoreLBStatus(ing.Status.PublicLoadBalancer)
		sort.SliceStable(curIPs, lessLoadBalancerIngress(curIPs))

		ingClient := client.NetworkingV1alpha1().Ingresses(ing.Namespace)

		currIng, err := ingClient.Get(ing.Name, metav1.GetOptions{})
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to fetch Knative Ingress %v/%v", ing.Namespace, ing.Name))
		}

		if ingressSliceEqual(status, curIPs) &&
			currIng.Status.ObservedGeneration == currIng.GetObjectMeta().GetGeneration() {
			logger.Debugf("no change in status, update skipped")
			return true, nil
		}

		logger.WithField("ingress_status", status).Debugf("attempting to update Knative Ingress status")
		lbStatus := toKnativeLBStatus(status)

		// TODO: handle the case when s.PublishService is empty
		namespace, svcName, err := utils.ParseNameNS(s.PublishService)
		clusterDomain := network.GetClusterDomainName()
		if err != nil {
			return false, err
		}

		for i := 0; i < len(lbStatus); i++ {
			lbStatus[i].DomainInternal = fmt.Sprintf("%s.%s.svc.%s",
				svcName, namespace, clusterDomain)
		}

		currIng.Status.MarkLoadBalancerReady(lbStatus, lbStatus, lbStatus)
		ingressCondSet.Manage(&currIng.Status).MarkTrue(knative.IngressConditionReady)
		ingressCondSet.Manage(&currIng.Status).MarkTrue(knative.IngressConditionNetworkConfigured)
		currIng.Status.ObservedGeneration = currIng.GetObjectMeta().GetGeneration()

		_, err = ingClient.UpdateStatus(currIng)
		if err != nil {
			logger.Errorf("failed to update ingress status: %v", err)
		} else {
			logger.WithField("ingress_status", status).Debugf("successfully updated ingress status")
		}
		return true, nil
	}
}

func (s *statusSync) runUpdateTCPIngress(ing *configurationv1beta1.TCPIngress,
	status []apiv1.LoadBalancerIngress,
	client configClientSet.Interface) pool.WorkFunc {
	return func(wu pool.WorkUnit) (interface{}, error) {
		if wu.IsCancelled() {
			return nil, nil
		}

		logger := s.Logger.WithFields(logrus.Fields{
			"ingress_namespace": ing.Namespace,
			"ingress_name":      ing.Name,
		})
		sort.SliceStable(status, lessLoadBalancerIngress(status))

		curIPs := ing.Status.LoadBalancer.Ingress
		sort.SliceStable(curIPs, lessLoadBalancerIngress(curIPs))

		if ingressSliceEqual(status, curIPs) {
			logger.Debugf("no change in status, update skipped")
			return true, nil
		}

		ingClient := client.ConfigurationV1beta1().TCPIngresses(ing.Namespace)

		currIng, err := ingClient.Get(ing.Name, metav1.GetOptions{})
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to fetch TCPIngress %v/%v", ing.Namespace, ing.Name))
		}

		logger.WithField("ingress_status", status).Debugf("attempting to update TCPIngress status")
		currIng.Status.LoadBalancer.Ingress = status
		_, err = ingClient.UpdateStatus(currIng)
		if err != nil {
			logger.Errorf("failed to update TCPIngress status: %v", err)
		} else {
			logger.WithField("ingress_status", status).Debugf("successfully updated TCPIngress status")
		}
		return true, nil
	}
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
