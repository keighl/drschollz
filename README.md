# Dr. Schollz

[![Build Status](https://travis-ci.org/keighl/drschollz.svg)](https://travis-ci.org/keighl/drschollz)

Easily send asynchronous error notifications from your Golang app via [Mandrill](http://mandrillapp.com).

### Install

    go get github.com/keighl/drschollz

### Usage

Configure it, and start it up

```go
var ds *drschollz.Queue

func init() {
    drschollz.Conf.AppName        = "MY_APP"
    drschollz.Conf.MandrillAPIKey = "XXXXXXXX"
    drschollz.Conf.EmailsTo       = []string{"devs@example.com"}
    drschollz.Conf.EmailFrom      = "errors@example.com"
}

func main() {
    // Start up a Dr Schollz queue with 3 workers
    ds, _ = drschollz.Start(3)
    defer ds.Stop()

    // your app...
}
```

Deliver error messages using `ds.Error()`

```go
func Something(w http.ResponseWriter, r *http.Request) {
    // ...
    err := doSomething()
    if err != nil {
        // ds.Error() asynchronously sends an email with the
        // err, and a backtrace. the method returns immediately.
        // You can also pass arbitrary extra info to be included in
        // the email
        ds.Error(err, user.ID, user.Age)
        fmt.Fprintf(w, err.Error())
    } else {
        fmt.Fprintf(w, "All good!")
    }
})
```

Or wrap return errors using `ds.Error()`

```go
func doSomething() error {
    // ...

    err := errors.New("Nothing is working!")

    // Error() returns the original error immediately,
    // so it can be used to wrap errors before they are
    // returned.
    return ds.Error(err)
})
```

See a working example in `./example`

### Example Email

```
-------------------------------
Error:
-------------------------------

    Unable to connect to the database

-------------------------------
Info:
-------------------------------

    - 2015-05-18 15:07:09.565792214 -0400 EDT

    - SELECT * FROM `posts`

    - UserID 401


-------------------------------
Backtrace:
-------------------------------

/Users/keighl/go/src/github.com/keighl/drschollz/example/main.go:44 main.Posts()
/Users/keighl/go/src/github.com/keighl/drschollz/example/main.go:27 main.func·001()
/usr/local/go/src/net/http/server.go:1265 net/http.HandlerFunc.ServeHTTP()
/usr/local/go/src/net/http/server.go:1541 net/http.(*ServeMux).ServeHTTP()
/usr/local/go/src/net/http/server.go:1703 net/http.serverHandler.ServeHTTP()
/usr/local/go/src/net/http/server.go:1204 net/http.(*conn).serve()
/usr/local/go/src/runtime/asm_amd64.s:2232 runtime.goexit()
```

### TODO

- Some better documentation
- Other email delivery platforms

