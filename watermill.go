package activitylog

import (
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
)

func NewActivityLogMiddlewareWatermill(publisher message.Publisher, cfg ActivityLogConfig) message.HandlerMiddleware {
	return func(h message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {

			trx := &Transaction{
				ActorType:        cfg.ActorType,
				ActorEmail:       cfg.ActorEmail,
				Service:          cfg.ServiceName,
				IsHtppMiddleware: false,
				Publisher:        publisher,
				RequestBody:      parseMessagePayload(msg),
			}
			trx.Start()
			msg.SetContext(NewContext(msg.Context(), trx))

			defer PublishLog(msg.Context(), publisher, trx, cfg.TopicName)
			defer trx.End()

			return h(msg)
		}
	}
}

func parseMessagePayload(msg *message.Message) map[string]interface{} {
	payload := make(map[string]interface{})
	err := json.Unmarshal(msg.Payload, &payload)
	if err != nil {
		return nil
	}
	return payload
}
