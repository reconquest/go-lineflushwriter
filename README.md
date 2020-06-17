# Line Flush Writer

Simple writer, that write only complete lines (e.g. that ends with newline) to
the specified backend writer.

It's guaranteed, that concurent writes to the same buffer using same mutex
will be safe and writes from one routine-go will not interleave output from
another.

# Behavior

```
New     Returns Writer With Specified Values
Writer  Writes Nothing At Empty Data
Writer  Writes Single New Line
Writer  Writes Nothing If Line Is Not Complete
Writer  Writes Line If Line Is Complete
Writer  Writes Only Complete Lines
Writer  Writes Complete Lines
Writer  Bufferize Line Until Complete
Writer  Flushes Buffer On Close
Writer  Do Not Append New Line If Nothing Written
Writer  Can Ensure Newline At End Of The String On Close
Writer  Not Appends Newlines Twice On Close
Writer  Call Backend Write Only Once Per Original Call
```

Generated with [loverage](https://github.com/kovetskiy/loverage).

# Benchmark

```
goos: linux
goarch: amd64
pkg: github.com/reconquest/lineflushwriter-go
BenchmarkWriter_Write_Parallel-12                       12476222              99.6 ns/op
BenchmarkWriter_Write_WritesLineAtomically-12            1400367               750 ns/op
PASS
ok      github.com/reconquest/lineflushwriter-go        3.315s
```

# Reference

See reference at [godoc.org](https://godoc.org/github.com/reconquest/lineflushwriter-go).
