package batching

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBatcher(t *testing.T) {
	t.Run("honors the batch size when writing", func(t *testing.T) {
		writer := &spyWriter{}
		batcher := NewBatcher(2, 1*time.Microsecond, writer)

		batcher.Write("data")
		batcher.Write("data2")
		assert.Equal(t, len(writer.batch), 2)
		assert.Equal(t, writer.batch[0], "data")
		assert.Equal(t, writer.batch[1], "data2")
	})

	t.Run("resets the internal buffer after writes", func(t *testing.T) {
		writer := &spyWriter{}
		batcher := NewBatcher(1, 1*time.Minute, writer)

		batcher.Write("data")
		assert.Equal(t, len(writer.batch), 1)
		assert.Equal(t, writer.batch[0], "data")

		batcher.Write("other-data")
		assert.Equal(t, len(writer.batch), 1)
		assert.Equal(t, writer.batch[0], "other-data")
	})

	t.Run("honors the time interval when writing", func(t *testing.T) {
		writer := &spyWriter{}
		b := NewBatcher(2, time.Minute, writer)

		b.Write("item")

		assert.Equal(t, len(writer.batch), 0)
	})

	t.Run("writes a partial batch when the interval has lapsed", func(t *testing.T) {
		writer := &spyWriter{}
		b := NewBatcher(10, 1*time.Millisecond, writer)

		b.Write("item")

		// Wait for interval to lapse
		time.Sleep(10 * time.Millisecond)
		b.Write("other-item")
		time.Sleep(2 * time.Millisecond)

		assert.Equal(t, len(writer.batch), 1)
		assert.Equal(t, writer.called, 2)
	})

	t.Run("honors the time interval when flushing", func(t *testing.T) {
		writer := &spyWriter{}
		b := NewBatcher(2, time.Minute, writer)

		b.Write("item")
		b.Flush()
		assert.Equal(t, writer.called, 1)
	})

	t.Run("avoids calling the writer with an empty batch", func(t *testing.T) {
		writer := &spyWriter{}
		b := NewBatcher(2, time.Nanosecond, writer)

		// Wait for interval to lapse
		time.Sleep(time.Millisecond)

		b.Flush()

		assert.Equal(t, writer.called, 0)
	})

	t.Run("writing all and clear timer when stopping", func(t *testing.T) {
		writer := &spyWriter{}
		b := NewBatcher(2, time.Second, writer)

		b.Write("item")
		b.Stop()
		b.Write("item data")
		b.Write("item data 2")
		b.Write("item data 3")
		assert.Equal(t, writer.called, 1)
	})

	t.Run("provides a means to guarantee a flushed write", func(t *testing.T) {
		writer := &spyWriter{}
		b := NewBatcher(2, time.Second, writer)

		b.Write("item")
		b.Flush()

		assert.Equal(t, writer.called, 1)
	})

	t.Run("avoids writing an empty batch on forced flush", func(t *testing.T) {
		writer := &spyWriter{}
		b := NewBatcher(2, time.Second, writer)

		b.Flush()

		assert.Equal(t, writer.called, 0)
	})
}

type spyWriter struct {
	batch  []interface{}
	called int
}

func (w *spyWriter) Write(batch []interface{}) {
	w.batch = batch
	fmt.Println(w.batch)
	w.called++
	fmt.Println(w.called)
}
