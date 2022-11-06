package bufferqueue

import (
	"io"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBBQEverythingEmpty(t *testing.T) {
	bbq := New()

	bbq.MarkEOF()

	buf := make([]byte, 1)

	n, err := bbq.Read(buf)
	assert.Equal(t, io.EOF, err)
	assert.EqualValues(t, 0, n)
	assert.Equal(t, []byte{0}, buf)
}

func TestBBQSimpleZeroSizedRead(t *testing.T) {
	bbq := New()

	bbq.QueueBuffer([]byte{1})
	bbq.MarkEOF()

	zeroSizedBuf := make([]byte, 0)
	buf := make([]byte, 1)

	n, err := bbq.Read(zeroSizedBuf)
	assert.NoError(t, err)
	assert.EqualValues(t, 0, n)

	n, err = bbq.Read(buf)
	assert.Equal(t, io.EOF, err)
	assert.EqualValues(t, 1, n)
	assert.Equal(t, []byte{1}, buf)

	n, err = bbq.Read(buf)
	assert.Equal(t, io.EOF, err)
	assert.EqualValues(t, 0, n)

	n, err = bbq.Read(zeroSizedBuf)
	assert.Equal(t, io.EOF, err)
	assert.EqualValues(t, 0, n)
}

func TestBBQSimpleTrivial(t *testing.T) {
	bbq := New()

	bbq.QueueBuffer([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	bbq.MarkEOF()

	buf := make([]byte, 10)

	n, err := bbq.Read(buf)
	assert.Equal(t, io.EOF, err)
	assert.EqualValues(t, 10, n)
	assert.Equal(t, []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, buf)

	n, err = bbq.Read(buf)
	assert.Equal(t, io.EOF, err)
	assert.EqualValues(t, 0, n)
}

func TestBBQSimpleTrivial2(t *testing.T) {
	bbq := New()

	bbq.QueueBuffer([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	bbq.MarkEOF()

	buf := make([]byte, 5)

	n, err := bbq.Read(buf)
	assert.NoError(t, err)
	assert.EqualValues(t, 5, n)
	assert.Equal(t, []byte{0, 1, 2, 3, 4}, buf)

	n, err = bbq.Read(buf)
	assert.Equal(t, io.EOF, err)
	assert.EqualValues(t, 5, n)
	assert.Equal(t, []byte{5, 6, 7, 8, 9}, buf)

	n, err = bbq.Read(buf)
	assert.Equal(t, io.EOF, err)
	assert.EqualValues(t, 0, n)
}

func TestBBQSimpleTrivial3(t *testing.T) {
	bbq := New()

	bbq.QueueBuffer([]byte{0, 1, 2, 3, 4})
	bbq.QueueBuffer([]byte{5, 6, 7, 8, 9})
	bbq.MarkEOF()

	buf := make([]byte, 5)

	n, err := bbq.Read(buf)
	assert.NoError(t, err)
	assert.EqualValues(t, 5, n)
	assert.Equal(t, []byte{0, 1, 2, 3, 4}, buf)

	n, err = bbq.Read(buf)
	assert.Equal(t, io.EOF, err)
	assert.EqualValues(t, 5, n)
	assert.Equal(t, []byte{5, 6, 7, 8, 9}, buf)

	n, err = bbq.Read(buf)
	assert.Equal(t, io.EOF, err)
	assert.EqualValues(t, 0, n)
}

func TestBBQSimpleTrivial4(t *testing.T) {
	bbq := New()

	bbq.QueueBuffer([]byte{0, 1, 2, 3, 4})
	bbq.QueueBuffer([]byte{5, 6, 7})
	bbq.MarkEOF()

	buf := make([]byte, 5)

	n, err := bbq.Read(buf)
	assert.NoError(t, err)
	assert.EqualValues(t, 5, n)
	assert.Equal(t, []byte{0, 1, 2, 3, 4}, buf)

	n, err = bbq.Read(buf)
	assert.Equal(t, io.EOF, err)
	assert.EqualValues(t, 3, n)
	assert.Equal(t, []byte{5, 6, 7, 3, 4}, buf)

	n, err = bbq.Read(buf)
	assert.Equal(t, io.EOF, err)
	assert.EqualValues(t, 0, n)
}

func TestBBQSimpleShortRead(t *testing.T) {
	bbq := New()

	bbq.QueueBuffer([]byte{0, 1, 2, 3})
	bbq.MarkEOF()

	buf := make([]byte, 10)

	n, err := bbq.Read(buf)
	assert.Equal(t, io.EOF, err)
	assert.EqualValues(t, 4, n)
	assert.Equal(t, []byte{0, 1, 2, 3, 0, 0, 0, 0, 0, 0}, buf)

	n, err = bbq.Read(buf)
	assert.Equal(t, io.EOF, err)
	assert.EqualValues(t, 0, n)
}

func TestBBQSomeEmptyBuffers(t *testing.T) {
	bbq := New()

	bbq.QueueBuffer([]byte{0, 1, 2, 3})
	bbq.QueueBuffer(nil)
	bbq.QueueBuffer([]byte{4})
	bbq.QueueBuffer([]byte{})
	bbq.QueueBuffer([]byte{5, 6, 7})
	bbq.QueueBuffer(nil)
	bbq.QueueBuffer([]byte{})
	bbq.QueueBuffer([]byte{8, 9})
	bbq.MarkEOF()

	buf := make([]byte, 5)

	n, err := bbq.Read(buf)
	assert.NoError(t, err)
	assert.EqualValues(t, 5, n)
	assert.Equal(t, []byte{0, 1, 2, 3, 4}, buf)

	n, err = bbq.Read(buf)
	assert.Equal(t, io.EOF, err)
	assert.EqualValues(t, 5, n)
	assert.Equal(t, []byte{5, 6, 7, 8, 9}, buf)

	n, err = bbq.Read(buf)
	assert.Equal(t, io.EOF, err)
	assert.EqualValues(t, 0, n)

}

func TestBBQSimpleQueueEverythingFirst(t *testing.T) {
	bbq := New()

	bbq.QueueBuffer([]byte{0, 1, 2, 3, 4})
	bbq.QueueBuffer([]byte{5, 6, 7})
	bbq.QueueBuffer([]byte{8, 9})
	bbq.MarkEOF()

	buf := make([]byte, 5)

	n, err := bbq.Read(buf)
	assert.NoError(t, err)
	assert.EqualValues(t, 5, n)
	assert.Equal(t, []byte{0, 1, 2, 3, 4}, buf)

	n, err = bbq.Read(buf)
	assert.Equal(t, io.EOF, err)
	assert.EqualValues(t, 5, n)
	assert.Equal(t, []byte{5, 6, 7, 8, 9}, buf)

	n, err = bbq.Read(buf)
	assert.Equal(t, io.EOF, err)
	assert.EqualValues(t, 0, n)
}

func TestBBQSeparateGoroutines(t *testing.T) {
	bbq := New()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		time.Sleep(5 * time.Millisecond)
		bbq.QueueBuffer([]byte{0, 1, 2, 3, 4})
		time.Sleep(4 * time.Millisecond)
		bbq.QueueBuffer([]byte{5, 6, 7})
		time.Sleep(3 * time.Millisecond)
		bbq.QueueBuffer([]byte{8, 9})
		time.Sleep(2 * time.Millisecond)
		bbq.MarkEOF()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		buf := make([]byte, 5)

		n, err := bbq.Read(buf)
		assert.NoError(t, err)
		assert.EqualValues(t, 5, n)
		assert.Equal(t, []byte{0, 1, 2, 3, 4}, buf)

		n, _ = bbq.Read(buf)
		// here err might be either nil or EOF, due to timing, so it's not
		// asserted
		assert.EqualValues(t, 5, n)
		assert.Equal(t, []byte{5, 6, 7, 8, 9}, buf)

		n, err = bbq.Read(buf)
		assert.Equal(t, io.EOF, err)
		assert.EqualValues(t, 0, n)
	}()

	wg.Wait()
}
