package requester

import (
	"context"
	"log"
)

type PubSubMessage struct {
	Data []byte `json:"data"`
}

// ConsumePubSub instruction and execute it
func ConsumePubSub(ctx context.Context, m PubSubMessage) error {
	i := decodeInstruction(m.Data)
	log.Printf("Received instruction: %+v", i)
	return nil
}
