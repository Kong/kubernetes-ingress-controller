package election

import (
	"fmt"
	"os"
	"time"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/golang/glog"
	"github.com/kong/kubernetes-ingress-controller/internal/k8s"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/client-go/tools/record"
)

type Elector interface {
	Run()
	IsLeader() bool
}

type Config struct {
	Client              clientset.Interface
	ElectionID          string
	IngressClass        string
	DefaultIngressClass string
	Callbacks           leaderelection.LeaderCallbacks
}

type electorStatus struct {
	Config
	electionID string
	elector    *leaderelection.LeaderElector
}

func (e electorStatus) Run() {
	e.elector.Run()
}

func (e electorStatus) IsLeader() bool {
	return e.elector.IsLeader()
}

func NewElector(config Config) Elector {
	pod, err := k8s.GetPodDetails(config.Client)
	if err != nil {
		glog.Fatalf("unexpected error obtaining pod information: %v", err)
	}

	es := electorStatus{
		Config: config,
	}

	// we need to use the defined ingress class to allow multiple leaders
	// in order to update information about ingress status
	electionID := fmt.Sprintf("%v-%v", config.ElectionID, config.DefaultIngressClass)
	if config.IngressClass != "" {
		electionID = fmt.Sprintf("%v-%v", config.ElectionID, config.IngressClass)
	}

	es.electionID = electionID

	broadcaster := record.NewBroadcaster()
	hostname, _ := os.Hostname()

	recorder := broadcaster.NewRecorder(scheme.Scheme, apiv1.EventSource{
		Component: "ingress-leader-elector",
		Host:      hostname,
	})

	lock := resourcelock.ConfigMapLock{
		ConfigMapMeta: metav1.ObjectMeta{Namespace: pod.Namespace, Name: electionID},
		Client:        config.Client.CoreV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity:      pod.Name,
			EventRecorder: recorder,
		},
	}

	ttl := 30 * time.Second
	le, err := leaderelection.NewLeaderElector(leaderelection.LeaderElectionConfig{
		Lock:          &lock,
		LeaseDuration: ttl,
		RenewDeadline: ttl / 2,
		RetryPeriod:   ttl / 4,
		Callbacks:     config.Callbacks,
	})

	if err != nil {
		glog.Fatalf("unexpected error starting leader election: %v", err)
	}

	es.elector = le
	return es
}
