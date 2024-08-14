# go-edu
The Go-Edu is texts, links, examples, playground to help to learn the Go/Golang.

## General
Hi Everyone! This is a repo to help people to learn Golang. It is based on my experience of learning the Go. I've been learning Go as a second language since 2016. If Golang is your first programming language maybe you need to start first from sources explaining basics of programming in general.

Go is a simple language by design. There are many mythologies about Go like "Go is a poor/simple language because it was invented by engineers from the 70s/80s who don't know modern things". I guess it is just another myth of IT. To get right understanding things/ideas behind of design of Golang read article by Go author Rob Pike ([Go at Google: Language Design in the Service of Software Engineering](https://go.dev/talks/2012/splash.article), this article location was changed several times, original version was based on "talk given by Rob Pike at the SPLASH 2012 conference in Tucson, Arizona, on October 25, 2012"). It explains a lot of things.

As I said Go is a simple language (Go spec less than `c` lang [Go Spec](https://go.dev/ref/spec)). The power of Go in implementation of concurrency. If Golang your second language and you have good understanding of basics of parallelism and concurrency (os/user level threads, race conditions, data synchronisations, segfaults, etc) to learn Go enough official documentations [Go Docs](https://go.dev/doc/).

## Go Learn Roadmap:
* [Go Tour] (https://go.dev/tour/list). Go tour gives you knowledge of 20% of language to solve 80% of problems (yes I am a little optimistic :).
* Best way to learn is practice. One of the good task it is write concurrent HTTP Server from scratch just using [stdlib tcp package](https://pkg.go.dev/net#example-Listener). It gives you understanding how to work with bytes, read/write from file/stream, work with concurrency, net, http, etc.
* After, implement HTTP Server using [stdlib http package](https://pkg.go.dev/net/http). You will use it in practice (Go stdlib good enough).
* To escape Golang [gotchas](https://en.wikipedia.org/wiki/Gotcha_(programming)) search articles "golang gotchas" to understand some  counter-intuitive things. For example:
    1. https://divan.dev/posts/avoid_gotchas/
    1. https://100go.co/
* See stdlib docs if something is not clear when you write code. It is minimalistic, clear, easy to read and provides good examples (consider it as best practices examples). [Go Stdlib](https://pkg.go.dev/std).
* [Effective Go](https://go.dev/doc/effective_go) code style.

## This repo examples of HTTP server:
* From scratch using `net` tcp. [httpserver v1](./internal/httpserverv1/)
* Using `net/http`. [httpserver v2](./internal/httpserverv2/)
* Using popular third party lib.

Thank you.
