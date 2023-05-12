package fly

import "sync"

type BufferedPipe struct {
	mu    sync.Mutex
	cond  *sync.Cond
	buf   []byte
	start int
	size  int
}

func NewBufferedPipe(size int) *BufferedPipe {
	p := &BufferedPipe{
		buf:   make([]byte, size),
		start: 0,
		size:  0,
	}
	p.cond = sync.NewCond(&p.mu)
	return p
}

func (p *BufferedPipe) Write(b []byte) (int, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for p.start+p.size+len(b) > cap(p.buf) {
		logger.Printf("Write pipe wait --- > start:%d, size:%d, len(b):%d, cap(buf):%d \n", p.start, p.size, len(b), cap(p.buf))
		p.cond.Wait()
	}
	//p.buf = append(p.buf, b...)
	copy(p.buf[p.start+p.size:p.start+p.size+len(b)], b)
	p.size += len(b)
	p.cond.Broadcast()
	return len(b), nil
}

func (p *BufferedPipe) Read(b []byte) (readSize int, err error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for p.size == 0 {
		p.cond.Wait()
	}
	if len(b) >= p.size {
		copy(b[0:p.size], p.buf[p.start:p.start+p.size])
		readSize = p.size
		p.start = 0
		p.size = 0
	} else {
		copy(b, p.buf[p.start:p.start+len(b)])
		readSize = len(b)
		p.start += readSize
		p.size -= readSize
		if p.start+p.size > cap(p.buf)/2 {
			newBuf := make([]byte, cap(p.buf))
			copy(newBuf[0:p.size], p.buf[p.start:p.start+p.size])
			p.buf = newBuf
			p.start = 0
		}
	}

	//n := copy(b, p.buf)
	//p.buf = p.buf[n:]
	p.cond.Broadcast()
	//return n, nil
	return
}

func (p *BufferedPipe) Size() int {
	return p.size
}
