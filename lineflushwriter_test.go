package lineflushwriter

import (
	"bytes"
	"io"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type nopCloser struct {
	*bytes.Buffer
}

func (closer nopCloser) Close() error {
	return nil
}

func TestNew_ReturnsWriterWithSpecifiedValues(t *testing.T) {
	test := assert.New(t)

	mutex := &sync.Mutex{}
	writer := New(nil, mutex, true)

	test.Equal(mutex, writer.mutex)
	test.Equal(true, writer.ensureNewline)
}

func TestWriter_WritesNothingAtEmptyData(t *testing.T) {
	testWriter(t, nil, false, "", "")
}

func TestWriter_WritesNothingIfLineIsNotComplete(t *testing.T) {
	testWriter(t, nil, false, "123", "")
}

func TestWriter_WritesLineIfLineIsComplete(t *testing.T) {
	testWriter(t, nil, false, "123\n", "123\n")
}

func TestWriter_WritesOnlyCompleteLines(t *testing.T) {
	testWriter(t, nil, false, "123\n456", "123\n")
}

func TestWriter_WritesCompleteLines(t *testing.T) {
	var writer io.WriteCloser

	writer = testWriter(t, nil, false, "123\n456\n", "123\n456\n")
	_ = testWriter(t, writer, false, "7\n", "123\n456\n7\n")
}

func TestWriter_BufferizeLineUntilComplete(t *testing.T) {
	var writer io.WriteCloser

	writer = testWriter(t, nil, false, "123", "")
	writer = testWriter(t, writer, false, "456", "")
	_ = testWriter(t, writer, false, "7\n", "1234567\n")
}

func TestWriter_FlushesBufferOnClose(t *testing.T) {
	writer := testWriter(t, nil, false, "123", "")
	testWriterClose(t, writer, "123")
}

func TestWriter_CanEnsureNewlineAtEndOfTheStringOnClose(t *testing.T) {
	writer := testWriter(t, nil, true, "123", "")
	testWriterClose(t, writer, "123\n")
}

func TestWriter_NotAppendsNewlinesTwiceOnClose(t *testing.T) {
	writer := testWriter(t, nil, true, "123\n", "123\n")
	testWriterClose(t, writer, "123\n")
}

func testWriterClose(
	t *testing.T,
	writer io.WriteCloser,
	expected string,
) {
	test := assert.New(t)

	writer.Close()

	buffer := writer.(*Writer).backend.(nopCloser).Buffer
	test.Equal(expected, buffer.String())
}

func testWriter(
	t *testing.T,
	writer io.WriteCloser,
	ensureNewline bool,
	data string,
	expected string,
) io.WriteCloser {
	test := assert.New(t)

	if writer == nil {
		writer = New(nopCloser{&bytes.Buffer{}}, &sync.Mutex{}, ensureNewline)
	}

	written, err := writer.Write([]byte(data))
	test.Nil(err)
	test.Equal(len(data), written)

	buffer := writer.(*Writer).backend.(nopCloser).Buffer

	test.Equal(expected, buffer.String())

	return writer
}
