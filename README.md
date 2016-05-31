# Line Flush Writer

Simple writer, that write only complete lines (e.g. that ends with newline) to
the specified backend writer.

It's guaranteed, that concurent writes to the same buffer using same mutex
will be safe and writes from one go-routine will not interleave output from
another.

# Behavior

```
New     Returns Writer With Specified Values
Writer  Writes Nothing At Empty Data
Writer  Writes Nothing If Line Is Not Complete
Writer  Writes Line If Line Is Complete
Writer  Writes Only Complete Lines
Writer  Writes Complete Lines
Writer  Bufferize Line Until Complete
Writer  Flushes Buffer On Close
Writer  Can Ensure Newline At End Of The String On Close
Writer  Not Appends Newlines Twice On Close
```

Generated with [loverage](https://github.com/kovetskiy/loverage)

# Reference

See reference at [godoc.org](https://godoc.org/github.com/reconquest/go-lineflushwriter)
