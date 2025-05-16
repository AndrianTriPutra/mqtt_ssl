package mqtt

import (
	"context"
	"errors"
	"time"
)

func (m *mqttStruct) Publish(ctx context.Context, topic string, message interface{}) error {
	sending := m.client.Publish(topic, byte(m.opts.Qos), false, message)
	if !sending.WaitTimeout(1 * time.Second) {
		err := errors.New("Publish timeout")
		return err
	}

	return nil
}
