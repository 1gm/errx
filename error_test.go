package errx_test

import (
	"errors"
	"testing"

	"fmt"

	"github.com/1gm/errx"
)

func asError(err error) *errx.Error {
	e, _ := err.(*errx.Error)
	return e
}

func Test_New(t *testing.T) {
	expectedMsg := "My error!"
	e := errx.New("My error!")
	if e.Error() != expectedMsg {
		t.Fatalf("expected message to be %v but was %v", expectedMsg, e.Error())
	}

	err := asError(e)
	if err.Inner != nil {
		t.Fatalf("expected inner error to be nil but was %v", err.Inner)
	}

	if err.StackTrace == nil {
		t.Fatalf("expected stack trace but was nil")
	}
}

func Test_Errorf(t *testing.T) {
	expectedMsg := fmt.Sprintf("My error is %s", "awesome.")
	e := errx.Errorf("My error is %s", "awesome.")
	if e.Error() != expectedMsg {
		t.Fatalf("expected message to be %v but was %v", expectedMsg, e.Error())
	}

	err := asError(e)
	if err.Inner != nil {
		t.Fatalf("expected inner error to be nil but was %v", err.Inner)
	}

	if err.StackTrace == nil {
		t.Fatalf("expected stack trace but was nil")
	}
}

func Test_Wrapped(t *testing.T) {
	td := []struct {
		inner                                                         error
		outerMessage, expectedErrorMessage, expectedInnerErrorMessage string
	}{
		{
			inner:                     errx.New("inner"),
			outerMessage:              "outer",
			expectedErrorMessage:      "outer: inner",
			expectedInnerErrorMessage: "inner",
		},
		{
			inner:                     errx.New("zxcvzxcv"),
			outerMessage:              "",
			expectedErrorMessage:      "zxcvzxcv",
			expectedInnerErrorMessage: "zxcvzxcv",
		},
		{
			inner:                     errx.New(""),
			outerMessage:              "",
			expectedErrorMessage:      "",
			expectedInnerErrorMessage: "",
		},
		{
			inner:                     errors.New("inner"),
			outerMessage:              "outer",
			expectedErrorMessage:      "outer: inner",
			expectedInnerErrorMessage: "inner",
		},
		{
			inner:                     errors.New("zxcvzxcv"),
			outerMessage:              "",
			expectedErrorMessage:      "zxcvzxcv",
			expectedInnerErrorMessage: "zxcvzxcv",
		},
		{
			inner:                     errors.New(""),
			outerMessage:              "",
			expectedErrorMessage:      "",
			expectedInnerErrorMessage: "",
		},
	}

	for i, test := range td {
		e := errx.Wrap(test.inner, test.outerMessage)
		if e.Error() != test.expectedErrorMessage {
			t.Fatalf("[%d]  expected e.Error() to be %s but was %s", i, test.expectedErrorMessage, e.Error())
		}

		err := asError(e)
		if err.Inner == nil {
			t.Fatalf("[%d] expected inner to not be nil", i)
		}

		if err.Inner.Error() != test.expectedInnerErrorMessage {
			t.Fatalf("[%d] expected message to be %s but was %s", i, test.expectedInnerErrorMessage, err.Inner.Error())
		}

		e = errx.Wrapf(test.inner, "%s", test.outerMessage)
		if e.Error() != test.expectedErrorMessage {
			t.Fatalf("[%d]  expected e.Error() to be %s but was %s", i, test.expectedErrorMessage, e.Error())
		}

		err = asError(e)
		if err.Inner == nil {
			t.Fatalf("[%d] expected inner to not be nil", i)
		}

		if err.Inner.Error() != test.expectedInnerErrorMessage {
			t.Fatalf("[%d] expected message to be %s but was %s", i, test.expectedInnerErrorMessage, err.Inner.Error())
		}
	}
}
