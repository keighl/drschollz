package main

import (
	"github.com/keighl/drschollz"
	"errors"
	"fmt"
	"net/http"
)

var (
	ds *drschollz.Queue
)

func init() {
	drschollz.Conf.AppName        = "DR_SCHOLLZ_EXAMPLE"
	drschollz.Conf.MandrillAPIKey = "XXXXXXXX"
	drschollz.Conf.EmailsTo       = []string{"devs@example.com"}
	drschollz.Conf.EmailFrom      = "errors@example.com"
}

func main() {
	// Start up a Dr Schollz queue with 3 workers
	ds, _ = drschollz.Start(3)
	defer ds.Stop()

	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		err := doSomething()
		if err != nil {
			fmt.Fprintf(w, err.Error())
		} else {
			fmt.Fprintf(w, "All good!")
		}
	})

    http.ListenAndServe(":5000", nil)
}

func doSomething() error {
	err := errors.New("Uh, there was a problemz")
	// Wrapping the error in ds.Error() asynchronously sends an
	// email with the err, a backtrace, and an abitrary list of
	// other stuff ([]interface{})
	// The ds.Error() method returns immediately.
	return ds.Error(err, "some", "debugging", "info")
}
