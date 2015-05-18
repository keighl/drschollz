package drschollz

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/keighl/mandrill"
	"runtime"
	"text/template"
	"time"
)

var (
	// Conf is the global Config object all queues referece
	Conf = &Config{
		AppName:       "[Go Error]",
		EmailTemplate: DefaultTmpl,
	}
)

// Config is well... a config object
type Config struct {
	AppName        string
	MandrillAPIKey string
	EmailFrom      string
	EmailsTo       []string
	EmailTemplate  string
}

func init() {
	Conf = &Config{
		AppName:       "APP_NAME",
		EmailTemplate: DefaultTmpl,
	}
}

// Error wraps an error value to be sent asynchronously on the queue.
// The method immediately returns the error value.
func (q *Queue) Error(err error, x ...interface{}) error {
	q.Println("Received err", err)

	if !q.Running {
		return err
	}

	q.jobQueue <- Job{
		Err:       err,
		BackTrace: backTrace(2),
		Extras:    x,
		Time:      time.Now(),
	}
	return err
}

// Deliver synchronously sends an error notification via Mandrill
func Deliver(j Job) error {
	if Conf.MandrillAPIKey == "" {
		return errors.New("No Mandrill API Key!")
	}

	if len(Conf.EmailsTo) == 0 {
		return errors.New("Need some `to` email addresses!")
	}

	if Conf.EmailFrom == "" {
		return errors.New("Need a `from` email address!")
	}

	tmpl, _ := template.New("DefaultTmpl").Parse(Conf.EmailTemplate)

	var content bytes.Buffer
	err := tmpl.Execute(&content, j)
	if err != nil {
		return err
	}

	message := &mandrill.Message{}
	for _, email := range Conf.EmailsTo {
		message.AddRecipient(email, "", "to")
	}
	message.FromEmail = Conf.EmailFrom
	message.FromName = "Dr. Schollz"
	message.Subject = fmt.Sprintf("[%s] %s", Conf.AppName, j.subject())
	message.Text = content.String()

	client := mandrill.ClientWithKey(Conf.MandrillAPIKey)
	_, err = client.MessagesSend(message)
	return err
}

func backTrace(skip int) (backTrace string) {
	for skip := skip; ; skip++ {
		pc, file, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		if file[len(file)-1] == 'c' {
			continue
		}
		f := runtime.FuncForPC(pc)
		backTrace += fmt.Sprintf("%s:%d %s()\n", file, line, f.Name())
	}
	return backTrace
}
