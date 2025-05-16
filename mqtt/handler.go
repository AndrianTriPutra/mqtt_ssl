package mqtt

import (
	"encoding/json"
	"mqTT/logger"
	"strings"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func (m *mqttStruct) versi1(client MQTT.Client, msg MQTT.Message) {
	topics := string(msg.Topic())
	logger.Trace("topic   :", topics)
	topic := strings.Trim(m.opts.Topic, "#")
	subtopic := strings.TrimPrefix(topics, topic)

	var js interface{}
	valid := json.Unmarshal(msg.Payload(), &js)
	if valid == nil {
		logger.Trace("subtopic:", subtopic)
		logger.Trace("payload :", string(msg.Payload()))
		m.msg <- msg.Payload()
	}
}

func (m *mqttStruct) OnMsg() chan []byte {
	return m.msg
}
