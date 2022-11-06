# bufferqueue

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/xen0n/bufferqueue)
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/xen0n/bufferqueue)
![GitHub branch checks state](https://img.shields.io/github/checks-status/xen0n/bufferqueue/develop)
[![GoDoc](https://pkg.go.dev/badge/github.com/xen0n/bufferqueue)](https://pkg.go.dev/github.com/xen0n/bufferqueue)

```go
import (
    "github.com/xen0n/bufferqueue" // package bufferqueue
)
```

A simple single-producer single-consumer buffer queue that only copies its
incoming data once, for every `Read()` from it. Zero-copy is not possible
given the semantic of `io.Reader`.

**Beware**: This implementation might be extremely buggy, do not put into
production yet without much testing.

## License

* [MIT](./LICENSE)
