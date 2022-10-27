package bufferqueue

import (
	"io"
)

// BBQ is a byte buffer queue that only copies data when it is read.
// The zero value of a BBQ is ready to use.
// This type is not designed for multi-producer or multi-consumer use.
type BBQ struct {
	bufs       [][]byte
	eof        bool
	avail      uint64
	notifyChan chan struct{}
}

var _ io.Reader = (*BBQ)(nil)

func New() *BBQ {
	return &BBQ{
		notifyChan: make(chan struct{}, 10),
	}
}

func (b *BBQ) QueueBuffer(buf []byte) {
	b.bufs = append(b.bufs, buf)
	b.avail += uint64(len(buf))
	b.notifyChan <- struct{}{}
}

func (b *BBQ) MarkEOF() {
	b.eof = true
	close(b.notifyChan)
}

func (b *BBQ) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		var err error
		if b.eof && b.avail == 0 {
			err = io.EOF
		}
		return 0, err
	}

	if !b.canImmediatelyRead(uint64(cap(p))) {
		// block until enough data is available or EOF
		for {
			<-b.notifyChan
			if b.canImmediatelyRead(uint64(cap(p))) {
				break
			}
		}
	}

	return b.immediatelyRead(p)
}

func (b *BBQ) canImmediatelyRead(numBytes uint64) bool {
	return b.avail >= numBytes || b.eof
}

// either p can be completely filled, or b.eof == true
func (b *BBQ) immediatelyRead(p []byte) (int, error) {
	if b.avail == 0 && b.eof {
		return 0, io.EOF
	}

	numBytesToRead := cap(p)
	numBytesRead := 0
	remainingBuf := p
	var err error
	for numBytesRead < numBytesToRead {
		copied := copy(remainingBuf, b.bufs[0])
		numBytesRead += copied
		remainingBuf = remainingBuf[copied:]
		if copied == len(b.bufs[0]) {
			// b.bufs[0] is fully consumed
			b.bufs = b.bufs[1:]
			if len(b.bufs) == 0 && b.eof {
				// eof is reached
				err = io.EOF
				break
			}
		} else {
			b.bufs[0] = b.bufs[0][copied:]
		}
	}

	b.avail -= uint64(numBytesRead)
	return numBytesRead, err
}
