# go-edu
The Go-Edu is a collection of texts, links and examples to help to learn the Go (Golang).

## General
Hi Everyone! This is a repository to help people to learn Golang. It is based on my experience of learning Go. I've been learning Go as a second language since 2016 and it's been my favorite language since then. If Golang is your first programming language then maybe my approach will not help you much.

Go is a simple language by design. There are many mythologies about Go like "Go is a poor/simple language because it was developed by engineers from the 70s/80s who don't know modern things". It is just another myth of IT. To get the right understanding of the ideas behind the design of Golang read the article by Go author Rob Pike, it explains a lot of things. ([Go at Google: Language Design in the Service of Software Engineering](https://go.dev/talks/2012/splash.article), this article location was changed several times, original version was based on "talk given by Rob Pike at the SPLASH 2012 conference in Tucson, Arizona, on October 25, 2012").

As I said Go is a simple language ([Go Spec](https://go.dev/ref/spec) less than `c` lang spec). The power of Go in simplicity and implementation of concurrency. If Golang your second language and you have good understanding of basics of parallelism and concurrency (os/user level threads, race conditions, data synchronisations, segfaults, etc) to learn Go enough the official documentations [Go Docs](https://go.dev/doc/).

## Go Learn Roadmap:
* The [Go Tour](https://go.dev/tour/list) takes 1-2 hours and gives you knowledge of 20% of language to solve 80% of problems (yes I am a little optimistic :). After that I started to write my first production code. Minimalism of language gives you this option.
* Best way to learn is practice. One of the good task it is write concurrent HTTP Server from scratch just using [stdlib's net package (tcp server)](https://pkg.go.dev/net#example-Listener). It gives you understanding how to work with bytes, read/write from file/stream, work with concurrency, net, http, etc.
* After, implement HTTP Server using [stdlib's http package (http server)](https://pkg.go.dev/net/http). You will use it in practice (Go's stdlib good enough).
* To escape Golang [gotchas](https://en.wikipedia.org/wiki/Gotcha_(programming)) search articles "golang gotchas" to understand some  counter-intuitive things. For example:
    1. https://divan.dev/posts/avoid_gotchas/
    1. https://100go.co/
* To learn Go just write code. See stdlib docs for code examples or if something is not clear when writing code. It is minimalistic, clear, easy to read and provides good examples (consider it as best practices examples). [Go Stdlib](https://pkg.go.dev/std).
* [Effective Go](https://go.dev/doc/effective_go) code style.
* If you do not understand why this code do not print indexes you need also look articles explaining how work Go's goroutines scheduler ([for example](https://www.ardanlabs.com/blog/2018/08/scheduling-in-go-part2.html)).
```golang
package main
import (
	"fmt"
	"runtime"
)
// if do not understand why this code will not print loop's indexes
// you probably need read how work go scheduler
func main() {
	// try comment this line and see difference
	runtime.GOMAXPROCS(1)
	for i := 0; i < 10; i++ {
		// run concurrent goroutines to print indexes
		go func(i int) {
			fmt.Println(i)
		}(i)
	}
	cnt := 0
	// emulate some computation
	// (it should not be very big otherwise you will not reproduce same result
	// because since version Go1.14 they add some solution improving this situation)
	for i := 0; i < 1000_000; i++ {
		cnt++
	}
}
```

## The repo examples of HTTP server:
* From scratch using `net` tcp. [httpserver v1](./internal/httpserverv1/)
* Using `net/http`. [httpserver v2](./internal/httpserverv2/)
* Using popular third party lib.

## Other examples of Go code:
* Simple revers tcp proxy with support mTLS Auth and balancing. [Simple TCP Proxy](https://github.com/radisvaliullin/proxy)

## Run example code
```
go run cmd/httpserver/main.go
```

Thank you.
