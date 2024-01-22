package batching

import (
	"context"
	"time"
)

// Batcher will accept messages and invoke the Writer when the batch
// requirements have been fulfilled (either batch size or interval have been
// exceeded). Batcher should be created with NewBatcher().
type Batcher struct {
	w          Writer
	interval   time.Duration
	timer      *time.Timer
	ctx        context.Context
	cancel     context.CancelFunc
	size       int
	batch      []interface{}
	isShutdown bool
}

// Writer is used to submit the completed batch. The batch may be partial if
// the interval lapsed instead of filling the batch.
type Writer interface {
	// Write submits the batch.
	Write(batch []interface{})
}

// WriterFunc is an adapter to allow ordinary functions to be a Writer.
type WriterFunc func(batch []interface{})

// Write implements Writer.
func (f WriterFunc) Write(batch []interface{}) {
	f(batch)
}

// NewBatcher creates a new Batcher.
func NewBatcher(size int, interval time.Duration, writer Writer) *Batcher {
	ctx, cancel := context.WithCancel(context.Background())

	batcher := &Batcher{
		size:       size,
		interval:   interval,
		w:          writer,
		ctx:        ctx,
		cancel:     cancel,
		isShutdown: false,
	}

	timer := time.NewTimer(interval)
	batcher.timer = timer
	batcher.listenTimer()
	return batcher
}

// Write stores data to the batch. It will not submit the batch to the writer
// until either the batch has been filled, or the interval has lapsed. NOTE:
// Write is *not* thread safe and should be called by the same goroutine that
// calls Flush.
func (b *Batcher) Write(data interface{}) {
	if b.isShutdown {
		return
	}
	b.batch = append(b.batch, data)
	if b.partialBatch() {
		return
	}

	b.writeBatch()
}

// Flush bypasses the batch interval and batch size checks and writes
// immediately.
func (b *Batcher) Flush() {
	b.writeBatch()
}

// writeBatch writes the batch (if any) to the writer and resets the batch and
// interval.
func (b *Batcher) writeBatch() {
	if len(b.batch) > 0 {
		b.w.Write(b.batch)
		b.batch = nil
	}

	b.timer.Reset(b.interval)
}

func (b *Batcher) partialBatch() bool {
	return len(b.batch) < b.size
}

func (b *Batcher) listenTimer() {
	go func() {
		for {
			select {
			case <-b.ctx.Done():
				b.timer.Stop()
				return
			case <-b.timer.C:
				b.writeBatch()
			}
		}
	}()
}

func (b *Batcher) Stop() {
	b.writeBatch()
	b.isShutdown = true
	b.cancel()
}
