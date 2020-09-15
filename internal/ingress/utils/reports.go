package utils

import (
	"net"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	reportsHost  = "kong-hf.konghq.com"
	reportsPort  = 61829
	pingInterval = 3600
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
	conn           *net.UDPConn

	Logger logrus.FieldLogger
}

func (r *Reporter) once() error {
	var serializedInfo string
	serializedInfo = serializedInfo + "v=" + r.Info.KICVersion + ";"
	serializedInfo = serializedInfo + "k8sv=" + r.Info.KubernetesVersion + ";"
	serializedInfo = serializedInfo + "kv=" + r.Info.KongVersion + ";"
	serializedInfo = serializedInfo + "db=" + r.Info.KongDB + ";"
	serializedInfo = serializedInfo + "id=" + r.Info.ID + ";"
	serializedInfo = serializedInfo + "hn=" + r.Info.Hostname + ";"
	r.serializedInfo = serializedInfo

	addr, err := net.ResolveUDPAddr("udp", reportsHost+":"+
		strconv.Itoa(reportsPort))
	if err != nil {
		return err
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	r.conn = conn
	return nil
}

// Run starts the reporter. It will send reports until done is closed.
func (r Reporter) Run(done <-chan struct{}) {
	err := r.once()
	if err != nil {
		r.Logger.Errorf("failed to initialize reporter: %s", err)
		return
	}
	defer r.conn.Close()

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

func (r Reporter) sendStart() {
	signal := prd + "-start"
	r.send(signal, 0)
}

func (r Reporter) sendPing(uptime int) {
	signal := prd + "-ping"
	r.send(signal, uptime)
}

func (r Reporter) send(signal string, uptime int) {
	message := "<14>signal=" + signal + ";uptime=" +
		strconv.Itoa(uptime) + ";" + r.serializedInfo
	_, err := r.conn.Write([]byte(message))
	if err != nil {
		r.Logger.Errorf("failed to send report: %s", err)
	}
}
