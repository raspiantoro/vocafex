package buffer

import "errors"

type CircularBuffer struct {
	buffer            []float32
	readPos, writePos int
	isResizing        bool
}

func NewCircularBuffer(size int, startingWritePos int) *CircularBuffer {
	return &CircularBuffer{
		buffer:   make([]float32, size),
		readPos:  0,
		writePos: startingWritePos,
	}
}

func (c *CircularBuffer) Read() (data float32, err error) {
	if c.readPos == len(c.buffer)-1 {
		c.readPos = 0
	}

	if c.readPos == c.writePos {
		err = errors.New("buffer is empty")
	}

	if c.isResizing {
		err = errors.New("resizing buffer")
	}

	data = c.buffer[c.readPos]
	c.readPos++

	return
}

func (c *CircularBuffer) Write(data float32) (err error) {
	if c.writePos == len(c.buffer)-1 {
		c.writePos = 0
	}

	if c.readPos == c.writePos {
		err = errors.New("buffer is full")
	}

	if c.isResizing {
		err = errors.New("resizing buffer")
	}

	c.buffer[c.writePos] = data
	c.writePos++

	return
}

func (c *CircularBuffer) Resize(size int) {
	newbuffer := make([]float32, size)
	copy(newbuffer, c.buffer)

	c.isResizing = true
	if c.writePos > (size - 1) {
		c.writePos = size - 1
	}

	if c.readPos > (size - 1) {
		c.readPos = size - 1
	}

	c.buffer = newbuffer
	c.isResizing = false
}
