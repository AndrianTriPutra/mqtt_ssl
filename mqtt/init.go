package mqtt

import (
	"context"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type Setting struct {
	Path string

	ServerName string
	Host       string
	User       string
	Pass       string

	ClientID string
	Qos      int
	Topic    string
}

type mqttStruct struct {
	client MQTT.Client
	opts   Setting

	msg chan []byte
}

type MqttI interface {
	Connect(context.Context) error
	OnMsg() chan []byte

	Publish(ctx context.Context, topic string, message interface{}) error
}

func NewMQTT(opts Setting) MqttI {
	return &mqttStruct{
		opts: opts,
	}
}
