package errx_test

import (
	"testing"

	"github.com/1gm/errx"
)

func TestStackTrace_TopFrame(t *testing.T) {
	var td = []struct {
		e     *errx.Error
		frame errx.StackFrame
	}{
		{
			asError(errx.New("e")),
			errx.StackFrame{
				FunctionName:    "github.com/1gm/errx_test.TestStackTrace_TopFrame",
				FileName:        "/home/george/go/src/github.com/1gm/errx/stack_test.go",
				TrimmedFileName: "github.com/1gm/errx/stack_test.go",
				Line:            15,
			},
		},
		{
			asError(fakeErrFunction()),
			errx.StackFrame{
				FunctionName:    "github.com/1gm/errx_test.fakeErrFunction",
				FileName:        "/home/george/go/src/github.com/1gm/errx/stack_test.go",
				TrimmedFileName: "github.com/1gm/errx/stack_test.go",
				Line:            55,
			},
		},
		{
			asError(anonymousFuncGenError()),
			errx.StackFrame{
				FunctionName:    "github.com/1gm/errx_test.anonymousFuncGenError.func1",
				FileName:        "/home/george/go/src/github.com/1gm/errx/stack_test.go",
				TrimmedFileName: "github.com/1gm/errx/stack_test.go",
				Line:            61,
			},
		},
	}

	for i, d := range td {
		if d.e.StackTrace == nil {
			t.Fatalf("[%d] expected stack trace but was nil", i)
		}

		if d.frame != d.e.StackTrace[0] {
			t.Fatalf("[%d] expected frame to be %v but was %v", i, d.frame, d.e.StackTrace[0])
		}
	}
}

func fakeErrFunction() error {
	e := errx.Wrap(errx.New("inner"), "outer")
	return asError(asError(e).Inner)
}

func anonymousFuncGenError() error {
	funcGen := func() error {
		return errx.New("foobar")
	}
	return funcGen()
}
