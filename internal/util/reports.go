package util

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/go-logr/logr"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/meshdetect"
)

var (
	reportsHost  = "kong-hf.konghq.com"
	reportsPort  = 61833
	pingInterval = 3600
	tlsConf      = tls.Config{MinVersion: tls.VersionTLS12, MaxVersion: tls.VersionTLS12} // nolint:gosec
	dialer       = net.Dialer{Timeout: time.Second * 30}
)

const (
	prd = "kic"
)

// Info holds the metadata to be sent as part of a report.
type Info struct {
	KubernetesVersion string
	KongVersion       string
	KICVersion        string
	Hostname          string
	KongDB            string
	ID                string
	FeatureGates      map[string]bool
}

// Reporter sends anonymous reports of runtime properties and
// errors in Kong.
type Reporter struct {
	Info Info

	serializedInfo string

	Logger logr.Logger

	MeshDetectionEnabled bool
	MeshDetector         *meshdetect.Detector
}

func (r *Reporter) once() {
	var serializedInfo string
	serializedInfo = serializedInfo + "v=" + r.Info.KICVersion + ";"
	serializedInfo = serializedInfo + "k8sv=" + r.Info.KubernetesVersion + ";"
	serializedInfo = serializedInfo + "kv=" + r.Info.KongVersion + ";"
	serializedInfo = serializedInfo + "db=" + r.Info.KongDB + ";"
	serializedInfo = serializedInfo + "id=" + r.Info.ID + ";"
	serializedInfo = serializedInfo + "hn=" + r.Info.Hostname + ";"

	for feature, enabled := range r.Info.FeatureGates {
		serializedInfo = fmt.Sprintf("%sfeature-%s=%t;", serializedInfo, strings.ToLower(feature), enabled)
	}

	r.serializedInfo = serializedInfo
}

// Run starts the reporter. It will send reports until done is closed.
func (r Reporter) Run(done <-chan struct{}) {
	r.once()

	r.sendStart()
	ticker := time.NewTicker(time.Duration(pingInterval) * time.Second)
	i := 1
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			r.sendPing(i * pingInterval)
			i++
		}
	}
}

func (r *Reporter) sendStart() {
	signal := prd + "-start"
	r.send(signal, 0)
}

func (r *Reporter) sendPing(uptime int) {
	signal := prd + "-ping"
	r.send(signal, uptime)
}

func (r *Reporter) send(signal string, uptime int) {
	message := "<14>signal=" + signal + ";uptime=" +
		strconv.Itoa(uptime) + ";" + r.serializedInfo
	// run mesh detection if enabled.
	if r.MeshDetectionEnabled {
		meshMessage, err := r.getMeshMessages(context.Background())
		if err != nil {
			// log the error if mesh detection fails,
			// but still send the messages without mesh detection results.
			r.Logger.V(DebugLevel).Info("failed to run mesh detection", "error", err)
		} else {
			// append results from mesh detection to reported message if returned any.
			if meshMessage != "" {
				if !strings.HasSuffix(message, ";") {
					message = message + ";"
				}
				message = message + meshMessage
			}
		}
	}

	conn, err := tls.DialWithDialer(&dialer, "tcp", net.JoinHostPort(reportsHost,
		strconv.FormatUint(uint64(reportsPort), 10)), &tlsConf)
	if err != nil {
		r.Logger.V(DebugLevel).Info("failed to connect to reporting server", "error", err)
		return
	}
	err = conn.SetDeadline(time.Now().Add(time.Minute))
	if err != nil {
		r.Logger.V(DebugLevel).Info("failed to set report connection deadline", "error", err)
		return
	}
	defer conn.Close()
	_, err = conn.Write([]byte(message))
	if err != nil {
		r.Logger.V(DebugLevel).Info("failed to send report", "error", err)
	}
}

// getMeshMessages returns serialized messages of results of mesh detection.
func (r *Reporter) getMeshMessages(ctx context.Context) (string, error) {
	if r.MeshDetector == nil {
		return "", fmt.Errorf("no mesh detector")
	}

	meshMessages := []string{}
	deploymentResults := r.MeshDetector.DetectMeshDeployment(ctx)
	meshMessages = append(meshMessages, serializeMeshDeploymentResults(deploymentResults))

	runUnderResults := r.MeshDetector.DetectRunUnder(ctx)
	meshMessages = append(meshMessages, serializeMeshRunUnderResults(runUnderResults))

	serviceDistributionResults, detectErr := r.MeshDetector.DetectServiceDistribution(ctx)
	if detectErr != nil {
		return "", fmt.Errorf("failed to detect service distribution under meshes: %w", detectErr)
	}
	meshMessages = append(meshMessages, serializeMeshServiceDistribution(serviceDistributionResults))
	return strings.Join(meshMessages, ";"), nil
}
