package util

import (
	"crypto/tls"
	"net"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	reportsHost  = "kong-hf.konghq.com"
	reportsPort  = 61833
	pingInterval = 3600
	tlsConf      = tls.Config{MinVersion: tls.VersionTLS13}
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
}

// Reporter sends anonymous reports of runtime properties and
// errors in Kong.
type Reporter struct {
	Info Info

	serializedInfo string

	Logger logrus.FieldLogger
}

func (r *Reporter) once() {
	var serializedInfo string
	serializedInfo = serializedInfo + "v=" + r.Info.KICVersion + ";"
	serializedInfo = serializedInfo + "k8sv=" + r.Info.KubernetesVersion + ";"
	serializedInfo = serializedInfo + "kv=" + r.Info.KongVersion + ";"
	serializedInfo = serializedInfo + "db=" + r.Info.KongDB + ";"
	serializedInfo = serializedInfo + "id=" + r.Info.ID + ";"
	serializedInfo = serializedInfo + "hn=" + r.Info.Hostname + ";"
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
	conn, err := tls.DialWithDialer(&dialer, "tcp", net.JoinHostPort(reportsHost,
		strconv.FormatUint(uint64(reportsPort), 10)), &tlsConf)
	if err != nil {
		r.Logger.Errorf("failed to connect to reporting server: %s", err)
		return
	}
	err = conn.SetDeadline(time.Now().Add(time.Minute))
	if err != nil {
		r.Logger.Errorf("failed to set report connection deadline: %s", err)
		return
	}
	defer conn.Close()
	_, err = conn.Write([]byte(message))
	if err != nil {
		r.Logger.Errorf("failed to send report: %s", err)
	}
}
