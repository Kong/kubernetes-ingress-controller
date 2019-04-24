package election

import (
	"os"
	"time"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/golang/glog"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/utils"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/client-go/tools/record"
)

// Elector elects a leader in a cluster.
type Elector interface {
	Run()
	IsLeader() bool
}

// Config holds the configuration for a leader election
type Config struct {
	Client     clientset.Interface
	ElectionID string
	Callbacks  leaderelection.LeaderCallbacks
}

type elector struct {
	Config
	elector *leaderelection.LeaderElector
}

func (e elector) Run() {
	e.elector.Run()
}

func (e elector) IsLeader() bool {
	return e.elector.IsLeader()
}

// NewElector returns an instance of Elector based on config.
func NewElector(config Config) Elector {
	pod, err := utils.GetPodDetails(config.Client)
	if err != nil {
		glog.Fatalf("unexpected error obtaining pod information: %v", err)
	}

	es := elector{
		Config: config,
	}

	broadcaster := record.NewBroadcaster()
	hostname, _ := os.Hostname()

	recorder := broadcaster.NewRecorder(scheme.Scheme, apiv1.EventSource{
		Component: "ingress-leader-elector",
		Host:      hostname,
	})

	lock := resourcelock.ConfigMapLock{
		ConfigMapMeta: metav1.ObjectMeta{Namespace: pod.Namespace,
			Name: config.ElectionID},
		Client: config.Client.CoreV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity:      pod.Name,
			EventRecorder: recorder,
		},
	}

	ttl := 30 * time.Second
	le, err := leaderelection.NewLeaderElector(
		leaderelection.LeaderElectionConfig{
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
