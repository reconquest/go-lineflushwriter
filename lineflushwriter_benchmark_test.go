package lineflushwriter

import (
	"bufio"
	"bytes"
	"io"
	"sync"
	"testing"
)

func BenchmarkWriter_Write_Parallel(b *testing.B) {
	data := []byte("partial\nwrite")

	buffer := &bytes.Buffer{}
	writer := New(nopCloser{buffer}, &sync.Mutex{}, true)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			writer.Write(data)
		}
	})
}

func BenchmarkWriter_Write_WritesLineAtomically(b *testing.B) {
	testLinePartA := "12345"
	testLinePartB := "6789"

	var (
		mutex  = &sync.Mutex{}
		buffer = &bytes.Buffer{}
		wg     = sync.WaitGroup{}
	)

	for i := 0; i < b.N; i++ {
		wg.Add(1)

		writer := New(nopCloser{buffer}, mutex, true)

		go func(writer io.Writer) {
			io.WriteString(writer, testLinePartA)
			io.WriteString(writer, testLinePartB+"\n")

			wg.Done()
		}(writer)
	}

	wg.Wait()

	scanner := bufio.NewScanner(buffer)
	for i := 0; i < b.N; i++ {
		if !scanner.Scan() {
			b.Fatalf("unexpected end of buffer at line %d", i)
		}

		if scanner.Text() != testLinePartA+testLinePartB {
			b.Fatalf("unexpected string at line %d: '%s'", i, scanner.Text())
		}
	}
}
