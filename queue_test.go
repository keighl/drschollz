package drschollz

import (
	"errors"
	"testing"
)

func Test_Start_Success(t *testing.T) {
	q, err := Start(1)
	q.Log = false
	expect(t, err, nil)
	refute(t, q, nil)
	expect(t, q.Running, true)
}

func Test_Start_NotEnoughWorkers(t *testing.T) {
	q, err := Start(0)
	refute(t, err, nil)
	expect(t, q, (*Queue)(nil))
}

func Test_Start_AlreadyRunning(t *testing.T) {
	q, err := Start(1)
	expect(t, err, nil)
	refute(t, q, nil)

	err = q.Start(1)
	refute(t, err, nil)
}

func Test_Stop(t *testing.T) {
	q, err := Start(1)
	q.Log = false
	expect(t, err, nil)
	refute(t, q, nil)

	q.Stop()
	expect(t, q.Running, false)
}

func Test_Stop_NotRunning(t *testing.T) {
	q, err := Start(1)
	q.Log = false
	expect(t, err, nil)
	refute(t, q, nil)

	q.Stop()
	expect(t, q.Running, false)

	q.Stop()
	expect(t, q.Running, false)
}

func Test_Queue_Logging(t *testing.T) {
	q, err := Start(1)
	expect(t, err, nil)
	q.Log = true
	q.Println("Kenny Loggins")

	q.Log = false
	q.Stop()
}

func Test_Worker_DeliveryError(t *testing.T) {
	q, err := Start(1)
	q.Log = false

	Conf.MandrillAPIKey = "SANDBOX_ERROR"
	err = errors.New("Cheese!")
	errOut := q.Error(err)
	expect(t, errOut, err)

	q.Stop()
	expect(t, q.Running, false)
	Conf.MandrillAPIKey = "SANDBOX_SUCCESS"
}

func Test_JobSubject_Short(t *testing.T) {
	j := &Job{
		Err: errors.New("Lorem ipsum dolor"),
	}
	expect(t, j.subject(), j.Error())
}

func Test_JobSubject_Long(t *testing.T) {
	j := &Job{
		Err: errors.New("Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."),
	}

	expect(t, j.subject(), j.Error()[0:60]+"...")
}
