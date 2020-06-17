package errorsutil

import (
	"errors"
	"fmt"
	"strings"
)

// Append adds additional text to an existing error.
func Append(err error, str string) error {
	return errors.New(fmt.Sprint(err) + str)
}

func Join(inclNils bool, errs ...error) error {
	strs := []string{}
	for _, err := range errs {
		if err == nil {
			if inclNils {
				strs = append(strs, "nil")
			}
		} else {
			strs = append(strs, err.Error())
		}
	}
	if len(strs) > 0 {
		return errors.New(strings.Join(strs, ";"))
	}
	return nil
}

// PanicOnErr is a syntactic sugar function to panic on error.
func PanicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func ErrorsToStrings(errs []error) (int, []string) {
	strs := []string{}
	count := 0
	for _, err := range errs {
		if err != nil {
			strs = append(strs, err.Error())
			count += 1
		} else {
			strs = append(strs, "")
		}
	}
	return count, strs
}

type ErrorInfo struct {
	Error       error
	ErrorString string // must match Error
	Code        string
	Display     string
	Input       string
	Correct     string
}

type ErrorInfos []*ErrorInfo

func (eis ErrorInfos) Inflate() {
	for i, ei := range eis {
		if ei != nil && ei.Error != nil {
			ei.ErrorString = ei.Error.Error()
		}
		eis[i] = ei
	}
}

func (eis ErrorInfos) GoodInputs() []string {
	inputs := []string{}
	for _, ei := range eis {
		if ei.Error == nil {
			inputs = append(inputs, ei.Input)
		}
	}
	return inputs
}

func (eis ErrorInfos) GoodCorrects() []string {
	inputs := []string{}
	for _, ei := range eis {
		if ei.Error == nil {
			inputs = append(inputs, ei.Correct)
		}
	}
	return inputs
}

func (eis ErrorInfos) ErrorsString() []string {
	estrings := []string{}
	for _, ei := range eis {
		if ei.Error != nil {
			estrings = append(estrings, ei.Error.Error())
		}
	}
	return estrings
}

func (eis ErrorInfos) Filter(isError bool) ErrorInfos {
	filtered := ErrorInfos{}
	for _, ei := range eis {
		if isError {
			if ei.Error != nil {
				filtered = append(filtered, ei)
			} else {
				filtered = append(filtered, nil)
			}
		} else if !isError {
			if ei.Error == nil {
				filtered = append(filtered, ei)
			} else {
				filtered = append(filtered, nil)
			}
		}
	}
	return filtered
}
