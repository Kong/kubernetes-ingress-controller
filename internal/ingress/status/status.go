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

	"github.com/golang/glog"
	"github.com/pkg/errors"

	configurationv1beta1 "github.com/kong/kubernetes-ingress-controller/internal/apis/configuration/v1beta1"
	configurationClientSet "github.com/kong/kubernetes-ingress-controller/internal/client/configuration/clientset/versioned"
	pool "gopkg.in/go-playground/pool.v3"
	apiv1 "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/wait"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/leaderelection"

	"github.com/kong/kubernetes-ingress-controller/internal/ingress/task"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/utils"
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
}

// Config ...
type Config struct {
	CoreClient       clientset.Interface
	KongConfigClient configurationClientSet.Interface

	OnStartedLeading func()

	PublishService string

	PublishStatusAddress string

	ElectionID string

	UpdateStatusOnShutdown bool

	IngressLister ingressLister

	UseNetworkingV1beta1 bool
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

	electionID string
	// workqueue used to keep in sync the status IP/s
	// in the Ingress rules
	syncQueue *task.Queue
	callbacks leaderelection.LeaderCallbacks
}

// Run starts the loop to keep the status in sync
func (s statusSync) Run() {
}

func (s statusSync) Callbacks() leaderelection.LeaderCallbacks {
	return s.callbacks
}

// Shutdown stop the sync. In case the instance is the leader it will remove the current IP
// if there is no other instances running.
func (s statusSync) Shutdown(isLeader bool) {
	go s.syncQueue.Shutdown()
	// remove IP from Ingress
	if !isLeader {
		return
	}

	// on shutdown we remove information about the leader election to
	// avoid up to 30 seconds of delay in start the synchronization process
	c, err := s.CoreClient.CoreV1().ConfigMaps(s.pod.Namespace).Get(s.electionID, metav1.GetOptions{})
	if err == nil {
		c.Annotations = map[string]string{}
		s.CoreClient.CoreV1().ConfigMaps(s.pod.Namespace).Update(c)
	}

	if !s.UpdateStatusOnShutdown {
		glog.Warningf("skipping update of status of Ingress rules")
		return
	}

	glog.Infof("updating status of Ingress rules (remove)")

	addrs, err := s.runningAddresses()
	if err != nil {
		glog.Errorf("error obtaining running IPs: %v", addrs)
		return
	}

	if len(addrs) > 1 {
		// leave the job to the next leader
		glog.Infof("leaving status update for next leader (%v)", len(addrs))
		return
	}

	if s.isRunningMultiplePods() {
		glog.V(2).Infof("skipping Ingress status update (multiple pods running - another one will be elected as master)")
		return
	}

	glog.Infof("removing address from ingress status (%v)", addrs)
	s.updateStatus([]apiv1.LoadBalancerIngress{})
}

func (s *statusSync) sync(key interface{}) error {
	if s.syncQueue.IsShuttingDown() {
		glog.V(2).Infof("skipping Ingress status update (shutting down in progress)")
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
func NewStatusSyncer(config Config) Sync {
	pod, err := utils.GetPodDetails(config.CoreClient)
	if err != nil {
		glog.Fatalf("unexpected error obtaining pod information: %v", err)
	}

	st := statusSync{
		pod: pod,

		Config: config,
	}
	st.syncQueue = task.NewCustomTaskQueue(st.sync, st.keyfunc)

	st.callbacks = leaderelection.LeaderCallbacks{
		OnStartedLeading: func(ctx context.Context) {
			glog.V(2).Infof("I am the new status update leader")
			if st.Config.OnStartedLeading != nil {
				st.Config.OnStartedLeading()
			}
			go st.syncQueue.Run(time.Second, ctx.Done())
			wait.PollUntil(updateInterval, func() (bool, error) {
				// send a dummy object to the queue to force a sync
				st.syncQueue.Enqueue("sync status")
				return false, nil
			}, ctx.Done())
		},
		OnStoppedLeading: func() {
			glog.V(2).Infof("I am not status update leader anymore")
		},
		OnNewLeader: func(identity string) {
			glog.Infof("new leader elected: %v", identity)
		},
	}

	return st
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
		glog.Errorf("error listing TPCIngresses for status update")
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

	batch.QueueComplete()
	batch.WaitAll()
}

func (s *statusSync) runUpdate(ing *networking.Ingress, status []apiv1.LoadBalancerIngress,
	client clientset.Interface) pool.WorkFunc {
	return func(wu pool.WorkUnit) (interface{}, error) {
		if wu.IsCancelled() {
			return nil, nil
		}

		sort.SliceStable(status, lessLoadBalancerIngress(status))

		curIPs := ing.Status.LoadBalancer.Ingress
		sort.SliceStable(curIPs, lessLoadBalancerIngress(curIPs))

		if ingressSliceEqual(status, curIPs) {
			glog.V(3).Infof("skipping update of Ingress %v/%v (no change)", ing.Namespace, ing.Name)
			return true, nil
		}

		if s.UseNetworkingV1beta1 {
			ingClient := client.NetworkingV1beta1().Ingresses(ing.Namespace)

			currIng, err := ingClient.Get(ing.Name, metav1.GetOptions{})
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("unexpected error searching Ingress %v/%v", ing.Namespace, ing.Name))
			}

			glog.Infof("updating Ingress %v/%v status to %v", currIng.Namespace, currIng.Name, status)
			currIng.Status.LoadBalancer.Ingress = status
			_, err = ingClient.UpdateStatus(currIng)
			if err != nil {
				glog.Warningf("error updating ingress rule: %v", err)
			}
		} else {
			ingClient := client.ExtensionsV1beta1().Ingresses(ing.Namespace)

			currIng, err := ingClient.Get(ing.Name, metav1.GetOptions{})
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("unexpected error searching Ingress %v/%v", ing.Namespace, ing.Name))
			}

			glog.Infof("updating Ingress %v/%v status to %v", currIng.Namespace, currIng.Name, status)
			currIng.Status.LoadBalancer.Ingress = status
			_, err = ingClient.UpdateStatus(currIng)
			if err != nil {
				glog.Warningf("error updating ingress rule: %v", err)
			}
		}
		return true, nil
	}
}

func (s *statusSync) runUpdateTCPIngress(ing *configurationv1beta1.TCPIngress,
	status []apiv1.LoadBalancerIngress,
	client configurationClientSet.Interface) pool.WorkFunc {
	return func(wu pool.WorkUnit) (interface{}, error) {
		if wu.IsCancelled() {
			return nil, nil
		}

		sort.SliceStable(status, lessLoadBalancerIngress(status))

		curIPs := ing.Status.LoadBalancer.Ingress
		sort.SliceStable(curIPs, lessLoadBalancerIngress(curIPs))

		if ingressSliceEqual(status, curIPs) {
			glog.V(3).Infof("skipping update of TCPIngress %v/%v (no change)", ing.Namespace, ing.Name)
			return true, nil
		}

		ingClient := client.ConfigurationV1beta1().TCPIngresses(ing.Namespace)

		currIng, err := ingClient.Get(ing.Name, metav1.GetOptions{})
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("unexpected error searching TCPIngress %v/%v", ing.Namespace, ing.Name))
		}

		glog.Infof("updating TCPIngress %v/%v status to %v", currIng.Namespace, currIng.Name, status)
		currIng.Status.LoadBalancer.Ingress = status
		_, err = ingClient.UpdateStatus(currIng)
		if err != nil {
			glog.Warningf("error updating status of TCPIngress: %v", err)
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
