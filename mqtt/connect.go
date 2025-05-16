package mqtt

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"mqTT/logger"
	"os"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func (m *mqttStruct) newTLSConfig() *tls.Config {
	caCert, err := os.ReadFile(m.opts.Path)
	if err != nil {
		logger.Level("fatal", "mqtt-newTLSConfig", "failed loading CA certificate:"+err.Error())
	}

	// Create a CA certificate pool and add the certificate
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caCert) {
		logger.Level("fatal", "mqtt-newTLSConfig", "failed to append CA certificate")
	}

	// Configure TLS
	return &tls.Config{
		RootCAs:            certPool, // Use the CA certificate pool
		MinVersion:         tls.VersionTLS12,
		ServerName:         m.opts.ServerName,
		InsecureSkipVerify: false, // Set to true if you don't want to verify the server certificate (not recommended)
	}
}

func (m *mqttStruct) Connect(ctx context.Context) error {
	//init chanel
	m.msg = make(chan []byte)

	// MQTT Client Options
	opts := MQTT.NewClientOptions()
	opts.AddBroker(m.opts.Host).SetUsername(m.opts.User).SetPassword(m.opts.Pass).SetClientID(m.opts.ClientID)

	// Set TLS Config
	opts.SetTLSConfig(m.newTLSConfig())

	// Keep-alive and reconnect options
	opts.SetKeepAlive(15 * time.Second)
	opts.SetCleanSession(true).SetConnectTimeout(15 * time.Second)
	opts.SetConnectRetry(true).SetConnectRetryInterval(15 * time.Second)
	opts.SetAutoReconnect(true).SetMaxReconnectInterval(15 * time.Second)
	opts.SetWill(m.opts.ClientID+"/wills", "good-bye!", 0, false)

	connection := false
	opts.OnConnect = func(client MQTT.Client) {
		logger.Level("info", "mqtt-Connect", "succesed connected to mqtt")
		connection = true

		subs := client.Subscribe(m.opts.Topic, byte(m.opts.Qos), m.versi1)
		if subs.WaitTimeout(5*time.Second) && subs.Error() != nil {
			logger.Level("error", "mqtt-Connect", "timeout subs")
		}
	}

	opts.OnConnectionLost = func(client MQTT.Client, err error) {
		logger.Level("warning", "mqtt-Connect", "connection lost . . .")
	}

	m.client = MQTT.NewClient(opts)
	token := m.client.Connect()
	if token.WaitTimeout(15*time.Second) && token.Error() != nil {
		logger.Level("error", "mqtt-Connect", "timeout")
		return token.Error()
	}
	logger.Level("debug", "mqtt-Connect", fmt.Sprintf("mqtt IsConnected ? %v", m.client.IsConnected()))

	time.Sleep(1 * time.Second)

	if !connection {
		newErr := errors.New("mqtt not connect")
		return newErr
	}

	logger.Level("debug", "mqtt-Connect", "done")

	return nil
}
