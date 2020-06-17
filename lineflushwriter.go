package lineflushwriter

import (
	"bytes"
	"io"
	"sync"
)

// Writer implements writer, that will proxy to specified `backend` writer only
// complete lines, e.g. that ends in newline. This writer is thread-safe.
type Writer struct {
	lock    sync.Locker
	backend io.WriteCloser
	buffer  []byte

	ensureNewline bool
}

// New returns new Writer, that will proxy data to the `backend` writer,
// thread-safety is guaranteed via `lock`. Optionally, writer can ensure, that
// last line of output ends with newline, if `ensureNewline` is true.
func New(
	writer io.WriteCloser,
	lock sync.Locker,
	ensureNewline bool,
) *Writer {
	return &Writer{
		backend: writer,
		lock:    lock,

		ensureNewline: ensureNewline,
	}
}

// Writer writes data into Writer.
//
// Signature matches with io.Writer's Write().
func (writer *Writer) Write(data []byte) (int, error) {
	writer.lock.Lock()
	writer.buffer = append(writer.buffer, data...)
	defer writer.lock.Unlock()

	var last = bytes.LastIndexByte(writer.buffer, '\n') + 1

	if last > 0 {
		written, err := writer.backend.Write(writer.buffer[:last])
		if err != nil {
			return written, err
		}

		writer.buffer = writer.buffer[last:]
	}

	return len(data), nil
}

// Close flushes all remaining data and closes underlying backend writer.
// If `ensureNewLine` was specified and remaining data does not ends with
// newline, then newline will be added.
//
// Signature matches with io.WriteCloser's Close().
func (writer *Writer) Close() error {
	if writer.ensureNewline && len(writer.buffer) > 0 {
		if writer.buffer[len(writer.buffer)-1] != '\n' {
			writer.buffer = append(writer.buffer, '\n')
		}
	}

	if len(writer.buffer) > 0 {
		_, err := writer.backend.Write(writer.buffer)
		if err != nil {
			return err
		}
	}

	return writer.backend.Close()
}
