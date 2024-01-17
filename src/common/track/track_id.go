package track

import (
	"context"
	"go-kit/src/common/logger"
)

const Key = "track_id"

func GetTrackId(ctx context.Context) string {
	v, _ := ctx.Value(Key).(string)
	return v
}

func InjectLogger(ctx context.Context, trackId string) context.Context {
	ctxWithTrackId := context.WithValue(ctx, Key, trackId)
	loggerCtx := logger.WithContextValue(ctxWithTrackId, map[string]interface{}{Key: trackId})
	return loggerCtx
}
