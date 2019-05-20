package request

import "context"

// ConsumePubSub decodes and execute instructions
func ConsumePubSub(ctx context.Context, m PubSubMessage) error {
	i, err := decodeInstruction(m.Data)
	if err != nil {
		return err
	}
	_, err = execute(ctx, i)
	return err
}
