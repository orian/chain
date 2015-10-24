# Chain (extended [Alice](https://github.com/justinas/alice))

[![Build Status](https://travis-ci.org/orian/chain.svg?branch=master)](https://travis-ci.org/orian/chain)

This document is a refactored version of the original Alice's README. Thanks for [Justinas](https://github.com/justinas) for a great library.

Chain provides a convenient way to chain
your context aware HTTP middleware functions and the app handler.

In short, it transforms

    Middleware1(Middleware2(Middleware3(App)))

to

    chain.New(Middleware1, Middleware2, Middleware3).Then(App).

### Why?

None of the other middleware chaining solutions behaves exactly like Alice.
Alice is as minimal as it gets:
in essence, it's just a for loop that does the wrapping for you.

Check out [this blog post](http://justinas.org/alice-painless-middleware-chaining-for-go/)
for explanation how Alice is different from other chaining solutions.

### Usage

Your middleware constructors should have the form of

```go
    // github.com/orian/wctx
    // Handler has one method: func () ServeHTTP(context.Context, http.ResponseWriter, http.Request)

    func (wctx.Handler) wctx.Handler
```

Some middleware provide this out of the box.
For ones that don't, it's trivial to write one yourself.

```go
func modifyContext(h wctx.Handler) wctx.Handler {
    return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
        c := context.WithValue(ctx, "key", "value")
        h.ServeHTTP(c, w, r)
    }
}
```

This complete example shows the full power of Alice.

```go
package main

import (
    "net/http"
    "time"

    "github.com/PuerkitoBio/throttled"
    "github.com/justinas/alice"
    "github.com/justinas/nosurf"
)

func timeoutHandler(h http.Handler) http.Handler {
    return http.TimeoutHandler(h, 1*time.Second, "timed out")
}

func myApp(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello world!"))
}

func main() {
    th := throttled.Interval(throttled.PerSec(10), 1, &throttled.VaryBy{Path: true}, 50)
    myHandler := http.HandlerFunc(myApp)

    chain := alice.New(th.Throttle, timeoutHandler, nosurf.NewPure).Then(myHandler)
    http.ListenAndServe(":8000", chain)
}
```

Here, the request will pass [throttled](https://github.com/PuerkitoBio/throttled) first,
then an http.TimeoutHandler we've set up,
then [nosurf](https://github.com/justinas/nosurf)
and will finally reach our handler.

Note that Alice makes **no guarantees** for
how one or another piece of  middleware will behave.
It executes all middleware sequentially so that if a
piece of middleware were to stop the chain,
the request will not reach the inner handlers.
This is intentional behavior.

Alice works with Go 1.0 and higher,
but running tests requires at least Go 1.1.

### Contributing

0. Find an issue that bugs you / open a new one.
1. Discuss.
2. Branch off, commit, test.
3. Make a pull request / attach the commits to the issue.
