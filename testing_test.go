package rotatelogs

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"path"
	"reflect"
	"runtime"
	"testing"
)

func assertEqual(t testing.TB, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	if err := validateEqualArgs(expected, actual); err != nil {
		return failTest(t, 1, fmt.Sprintf("Equal: invalid operation `%#v` == `%#v` (%v)", expected, actual, err), msgAndArgs...)
	}

	if !isObjectEqual(expected, actual) {
		return failTest(t, 1, fmt.Sprintf("Equal: expected `%#v`, actual `%#v`", expected, actual), msgAndArgs...)
	}

	return true
}

func assertNoError(t testing.TB, err error, msgAndArgs ...interface{}) bool {
	if err != nil {
		return failTest(t, 1, fmt.Sprintf("Received unexpected error: `%v`", err), msgAndArgs...)
	}

	return true
}

func assertError(t testing.TB, err error, msgAndArgs ...interface{}) bool {
	if err == nil {
		return failTest(t, 1, "An error is expected but got nil.", msgAndArgs...)
	}

	return true
}

func assertTrue(t testing.TB, value bool, msgAndArgs ...interface{}) bool {
	if !value {
		return failTest(t, 1, "Should be true", msgAndArgs...)
	}

	return true
}

func assertFalse(t testing.TB, value bool, msgAndArgs ...interface{}) bool {
	if value {
		return failTest(t, 1, "Should be false", msgAndArgs...)
	}

	return true
}

func assertStringLen(t testing.TB, str string, length int, msgAndArgs ...interface{}) bool {
	l := len(str)
	if l != length {
		return failTest(t, 1, fmt.Sprintf("\"%s\" should have %d item(s), but has %d", str, length, l), msgAndArgs...)
	}

	return true
}

func assertFileExists(t testing.TB, path string, msgAndArgs ...interface{}) bool {
	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return failTest(t, 1, fmt.Sprintf("unable to find file %q", path), msgAndArgs...)
		}
		return failTest(t, 1, fmt.Sprintf("error when running os.Lstat(%q): %s", path, err), msgAndArgs...)
	}
	if info.IsDir() {
		return failTest(t, 1, fmt.Sprintf("%q is a directory", path), msgAndArgs...)
	}
	return true
}

func validateEqualArgs(expected, actual interface{}) error {
	if expected == nil || actual == nil {
		return nil
	}

	expectedKind := reflect.TypeOf(expected).Kind()
	actualKind := reflect.TypeOf(actual).Kind()
	if expectedKind == reflect.Func || actualKind == reflect.Func {
		return errors.New("cannot take func type as argument")
	}

	return nil
}

func isObjectEqual(expected, actual interface{}) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}

	return reflect.DeepEqual(expected, actual)
}

func failTest(t testing.TB, skip int, msg string, msgAndArgs ...interface{}) bool {
	flag := ""

	if skip >= 0 {
		_, file, line, _ := runtime.Caller(skip + 1)
		addition := ""
		if len(msgAndArgs) == 1 {
			addition = fmt.Sprintf("%v", msgAndArgs[1])
		} else if len(msgAndArgs) > 1 {
			addition = fmt.Sprintf(fmt.Sprintf("%v", msgAndArgs[0]), msgAndArgs[1:]...)
		}

		if addition == "" {
			fmt.Printf("%s%s:%d %s\n", flag, path.Base(file), line, msg)
		} else {
			fmt.Printf("%s%s:%d %s [%s]\n", flag, path.Base(file), line, msg, addition)
		}
	}
	t.Fail()

	return false
}
