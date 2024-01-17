package pubsub

import (
	"context"
	"encoding/json"
	"go-kit/src/common/track"
	"io"
)

func ConvertPayload[T Payload](ctx context.Context, message io.Reader) (context.Context, *T, error) {
	var payload T
	err := json.NewDecoder(message).Decode(&payload)
	if err != nil {
		return ctx, nil, err
	}
	ctx = track.InjectLogger(ctx, payload.GetTrackId())
	return ctx, &payload, nil
}
