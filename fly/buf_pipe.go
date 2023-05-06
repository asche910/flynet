package fly

import "sync"

type BufferedPipe struct {
	mu   sync.Mutex
	cond *sync.Cond
	buf  []byte
}

func NewBufferedPipe(size int) *BufferedPipe {
	p := &BufferedPipe{
		buf: make([]byte, 0, size),
	}
	p.cond = sync.NewCond(&p.mu)
	return p
}

func (p *BufferedPipe) Write(b []byte) (int, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for len(p.buf)+len(b) > cap(p.buf) {
		p.cond.Wait()
	}
	p.buf = append(p.buf, b...)
	p.cond.Broadcast()
	return len(b), nil
}

func (p *BufferedPipe) Read(b []byte) (int, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for len(p.buf) == 0 {
		p.cond.Wait()
	}
	n := copy(b, p.buf)
	p.buf = p.buf[n:]
	p.cond.Broadcast()
	return n, nil
}

func (p *BufferedPipe) Size() int {
	return len(p.buf)
}
