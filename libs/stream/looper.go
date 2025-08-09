package stream

import (
	"io"
	"time"
)

type loop struct {
	start, end int64
}

type Looper struct {
	src   io.ReadSeeker
	srcSr float64

	loops   []*loop
	current *loop

	cursor int64
}

func NewLooper(sampleRate int, stream io.ReadSeeker) *Looper {
	return &Looper{
		src:   stream,
		srcSr: float64(sampleRate),
	}
}

func (l *Looper) Read(p []byte) (int, error) {
	// No loop
	if l.current == nil {
		n, err := l.src.Read(p)
		l.cursor += int64(n)
		return n, err
	}

	// End of loop reached
	if l.cursor >= l.current.end {
		l.cursor = l.current.start
		_, err := l.src.Seek(l.cursor, io.SeekStart)
		if err != nil {
			return 0, err
		}
	}

	// Read data from the source
	n, err := l.src.Read(p)
	if err != nil {
		return n, err
	}

	// Trim if went past the end
	l.cursor += int64(n)
	if l.cursor > l.current.end {
		n -= int(l.cursor - l.current.end)
	}

	return n, nil
}

func (l *Looper) AddLoop(start, end time.Duration) int {
	// 2 bytes per sample, 2 channels
	const bytePerFrame = 4

	s := l.srcSr * bytePerFrame * start.Seconds()
	e := l.srcSr * bytePerFrame * end.Seconds()

	// Stop right on frame
	sb := int64(s) + int64(s)%bytePerFrame
	eb := int64(e) - int64(e)%bytePerFrame

	l.loops = append(l.loops, &loop{start: sb, end: eb})

	return len(l.loops) - 1
}

func (l *Looper) Play(index int) {
	if index < 0 || index >= len(l.loops) {
		return
	}

	l.current = l.loops[index]
}

func (l *Looper) Stop() {
	l.current = nil
}
