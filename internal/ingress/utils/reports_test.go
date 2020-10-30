package utils

import (
	"bytes"
	"net"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	reportsHost = "localhost"
	pingInterval = 1
	os.Exit(m.Run())
}

func TestReporterOnce(t *testing.T) {
	assert := assert.New(t)
	info := Info{
		KubernetesVersion: "k8s.version",
		KongVersion:       "kong.version",
		KICVersion:        "kic.version",
		Hostname:          "example.local",
		KongDB:            "off",
		ID:                "6acb7447-eedf-4815-a193-d714c5108f7b",
	}
	reporter := Reporter{
		Info:   info,
		Logger: logrus.New(),
	}
	assert.Nil(reporter.once())
	want := "v=kic.version;k8sv=k8s.version;kv=kong.version;db=off;" +
		"id=6acb7447-eedf-4815-a193-d714c5108f7b;hn=example.local;"
	assert.Equal(want, reporter.serializedInfo)
}

func TestReporterSendStart(t *testing.T) {
	assert := assert.New(t)
	info := Info{
		KubernetesVersion: "k8s.version",
		KongVersion:       "kong.version",
		KICVersion:        "kic.version",
		Hostname:          "example.local",
		KongDB:            "off",
		ID:                "6acb7447-eedf-4815-a193-d714c5108f7b",
	}
	reporter := Reporter{
		Info:   info,
		Logger: logrus.New(),
	}
	assert.Nil(reporter.once())
	addr, err := net.ResolveUDPAddr("udp", reportsHost+
		":"+strconv.Itoa(reportsPort))
	assert.Nil(err)
	conn, err := net.ListenUDP("udp", addr)
	assert.Nil(err)
	defer conn.Close()

	reporter.sendStart()

	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buffer)
	serialized := "<14>signal=kic-start;uptime=0;v=kic.version;" +
		"k8sv=k8s.version;kv=kong.version;db=off;" +
		"id=6acb7447-eedf-4815-a193-d714c5108f7b;hn=example.local;"
	assert.Equal(len(serialized), n)
	assert.Nil(err)
	assert.Equal(serialized, string(bytes.Trim(buffer, "\x00")))
}

func TestReporterSendPing(t *testing.T) {
	assert := assert.New(t)
	info := Info{
		KubernetesVersion: "k8s.version",
		KongVersion:       "kong.version",
		KICVersion:        "kic.version",
		Hostname:          "example.local",
		KongDB:            "off",
		ID:                "6acb7447-eedf-4815-a193-d714c5108f7b",
	}
	reporter := Reporter{
		Info:   info,
		Logger: logrus.New(),
	}
	assert.Nil(reporter.once())
	addr, err := net.ResolveUDPAddr("udp", reportsHost+
		":"+strconv.Itoa(reportsPort))
	assert.Nil(err)
	conn, err := net.ListenUDP("udp", addr)
	assert.Nil(err)
	defer conn.Close()

	reporter.sendPing(42)

	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buffer)
	serialized := "<14>signal=kic-ping;uptime=42;v=kic.version;" +
		"k8sv=k8s.version;kv=kong.version;db=off;" +
		"id=6acb7447-eedf-4815-a193-d714c5108f7b;hn=example.local;"
	assert.Equal(len(serialized), n)
	assert.Nil(err)
	assert.Equal(serialized, string(bytes.Trim(buffer, "\x00")))
}

func TestReporterRun(t *testing.T) {
	assert := assert.New(t)
	info := Info{
		KubernetesVersion: "k8s.version",
		KongVersion:       "kong.version",
		KICVersion:        "kic.version",
		Hostname:          "example.local",
		KongDB:            "off",
		ID:                "6acb7447-eedf-4815-a193-d714c5108f7b",
	}
	reporter := Reporter{
		Info:   info,
		Logger: logrus.New(),
	}
	assert.Nil(reporter.once())
	addr, err := net.ResolveUDPAddr("udp", reportsHost+
		":"+strconv.Itoa(reportsPort))
	assert.Nil(err)
	conn, err := net.ListenUDP("udp", addr)
	assert.Nil(err)
	defer conn.Close()
	done := make(chan struct{})

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		reporter.Run(done)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		serializedContent := []string{
			"<14>signal=kic-start;uptime=0;v=kic.version;k8sv=k8s.version;" +
				"kv=kong.version;db=off;" +
				"id=6acb7447-eedf-4815-a193-d714c5108f7b;hn=example.local;",
			"<14>signal=kic-ping;uptime=1;v=kic.version;k8sv=k8s.version;" +
				"kv=kong.version;db=off;" +
				"id=6acb7447-eedf-4815-a193-d714c5108f7b;hn=example.local;",
			"<14>signal=kic-ping;uptime=2;v=kic.version;k8sv=k8s.version;" +
				"kv=kong.version;db=off;" +
				"id=6acb7447-eedf-4815-a193-d714c5108f7b;hn=example.local;",
		}
		for _, expect := range serializedContent {
			buffer := make([]byte, 1024)
			n, _, err := conn.ReadFromUDP(buffer)
			assert.Equal(len(expect), n)
			assert.Nil(err)
			assert.Equal(expect, string(bytes.Trim(buffer, "\x00")))

		}
		close(done)
	}()
	wg.Wait()
}
