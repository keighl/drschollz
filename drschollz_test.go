package drschollz

import (
	"errors"
	"reflect"
	"testing"
)

func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func init() {
	Conf.MandrillAPIKey = "SANDBOX_SUCCESS"
	Conf.EmailsTo = []string{"tester@example.com"}
	Conf.EmailFrom = "errors@example.com"
}

//// ERROR ///////

func Test_Error_Success(t *testing.T) {
	q, _ := Start(2)
	q.Log = false
	err := errors.New("Cheese!")
	errOut := q.Error(err)
	expect(t, errOut, err)
	q.Stop()
}

func Test_Error_NotRunning(t *testing.T) {
	q, _ := Start(1)
	refute(t, q, nil)
	q.Log = false
	q.Stop()

	err := errors.New("Cheese!")
	errOut := q.Error(err)
	expect(t, errOut, err)
}

//// DELIVER ///////

func Test_Deliver_Success(t *testing.T) {
	Conf.MandrillAPIKey = "SANDBOX_SUCCESS"
	err := errors.New("Cheese!")
	j := Job{
		BackTrace: "XXXXXX",
		Err:       err,
		Extras:    []interface{}{"XXXXXXX", 4, false},
	}

	deliverErr := Deliver(j)
	expect(t, deliverErr, nil)
}

func Test_Deliver_Fail(t *testing.T) {
	Conf.MandrillAPIKey = "SANDBOX_ERROR"
	err := errors.New("Cheese!")
	j := Job{
		BackTrace: "XXXXXX",
		Err:       err,
		Extras:    []interface{}{"XXXXXXX", 4, false},
	}

	deliverErr := Deliver(j)
	refute(t, deliverErr, nil)
	Conf.MandrillAPIKey = "SANDBOX_SUCCESS"
}

func Test_Deliver_FailNoAPIKey(t *testing.T) {
	Conf.MandrillAPIKey = ""
	deliverErr := Deliver(Job{})
	refute(t, deliverErr, nil)
	Conf.MandrillAPIKey = "SANDBOX_SUCCESS"
}

func Test_Deliver_FailNoEmailsTo(t *testing.T) {
	Conf.EmailsTo = []string{}
	deliverErr := Deliver(Job{})
	refute(t, deliverErr, nil)
	Conf.EmailsTo = []string{"tester@example.com"}
}

func Test_Deliver_FailNoEmailFrom(t *testing.T) {
	Conf.EmailFrom = ""
	deliverErr := Deliver(Job{})
	refute(t, deliverErr, nil)
	Conf.EmailFrom = "errors@example.com"
}

func Test_Deliver_ContentFail(t *testing.T) {
	Conf.EmailTemplate = `{{.NOOOOOOOOO}}`

	err := errors.New("Cheese!")
	j := Job{
		BackTrace: "XXXXXX",
		Err:       err,
		Extras:    []interface{}{"XXXXXXX", 4, false},
	}

	deliverErr := Deliver(j)
	refute(t, deliverErr, nil)
	Conf.EmailTemplate = DefaultTmpl
}
