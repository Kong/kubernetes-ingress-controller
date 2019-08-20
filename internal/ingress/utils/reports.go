package utils

import (
	"net"
	"strconv"
	"time"

	"github.com/golang/glog"
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
	info Info

	serializedInfo string
	conn           *net.UDPConn
}

// NewReporter creates a reporter based on the info provided.
// Use the Run() method to start reporting.
func NewReporter(info Info) Reporter {
	return Reporter{info: info}
}

func (r *Reporter) once() error {
	var serializedInfo string
	serializedInfo = serializedInfo + "v=" + r.info.KICVersion + ";"
	serializedInfo = serializedInfo + "k8sv=" + r.info.KubernetesVersion + ";"
	serializedInfo = serializedInfo + "kv=" + r.info.KongVersion + ";"
	serializedInfo = serializedInfo + "db=" + r.info.KongDB + ";"
	serializedInfo = serializedInfo + "id=" + r.info.ID + ";"
	serializedInfo = serializedInfo + "hn=" + r.info.Hostname + ";"
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
		glog.Errorf("error initializing reports: %s", err)
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
		glog.Errorf("error sending report: %s", err)
	}
}
