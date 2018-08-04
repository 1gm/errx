package errx_test

import (
	"testing"

	"os"
	"path/filepath"

	"github.com/1gm/errx"
)

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

func genError() *errx.Error {
	return asError(errx.New("e"))
}

func getAbsoluteStackTestFilePath(t *testing.T) string {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("%v: failed to get cwd", t.Name())
	}
	return filepath.Join(cwd, "stack_test.go")
}

func TestStackTrace_TopFrame(t *testing.T) {
	// absolute path of the this file (stack_test.go)
	stackTestFilePath := getAbsoluteStackTestFilePath(t)

	var td = []struct {
		e     *errx.Error
		frame errx.StackFrame
	}{
		{
			genError(),
			errx.StackFrame{
				FunctionName:    "github.com/1gm/errx_test.genError",
				FileName:        stackTestFilePath,
				TrimmedFileName: "github.com/1gm/errx/stack_test.go",
				Line:            25,
			},
		},
		{
			asError(fakeErrFunction()),
			errx.StackFrame{
				FunctionName:    "github.com/1gm/errx_test.fakeErrFunction",
				FileName:        stackTestFilePath,
				TrimmedFileName: "github.com/1gm/errx/stack_test.go",
				Line:            13,
			},
		},
		{
			asError(anonymousFuncGenError()),
			errx.StackFrame{
				FunctionName:    "github.com/1gm/errx_test.anonymousFuncGenError.func1",
				FileName:        stackTestFilePath,
				TrimmedFileName: "github.com/1gm/errx/stack_test.go",
				Line:            19,
			},
		},
	}

	for i, d := range td {
		if d.e.StackTrace == nil {
			t.Errorf("[%d] expected stack trace but was nil", i)
			continue
		}

		if d.frame != d.e.StackTrace[0] {
			t.Errorf("[%d] expected frame to be %v but was %v", i, d.frame, d.e.StackTrace[0])
		}
	}
}
