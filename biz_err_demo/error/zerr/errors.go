package zerr

import (
	"errors"
	"fmt"
	"io"
)

func New(message string) error {
	return &fundamental{
		msg:   message,
		stack: callers(),
	}
}

func Errorf(format string, args ...interface{}) error {
	return &fundamental{
		msg:   fmt.Sprintf(format, args...),
		stack: callers(),
	}
}

type fundamental struct {
	msg string
	*stack
}

func (f *fundamental) Error() string { return f.msg }

func (f *fundamental) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, f.msg)
			f.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, f.msg)
	case 'q':
		fmt.Fprintf(s, "%q", f.msg)
	}
}

func WithStack(err error) error {
	if err == nil {
		return nil
	}
	return &withStack{
		err,
		callers(),
	}
}

type withStack struct {
	error
	*stack
}

func (w *withStack) Cause() error { return w.error }

func (w *withStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v", w.Cause())
			w.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, w.Error())
	case 'q':
		fmt.Fprintf(s, "%q", w.Error())
	}
}

func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	err = &withMessage{
		cause: err,
		msg:   message,
	}
	return &withStack{
		err,
		callers(),
	}
}

func Trace(err error) error {
	return Wrapf(err, "")
}

func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	err = &withMessage{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}
	return &withStack{
		err,
		callers(),
	}
}

func WithMessage(err error, message string) error {
	if err == nil {
		return nil
	}
	return &withMessage{
		cause: err,
		msg:   message,
	}
}

func WithMessagef(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &withMessage{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}
}

type withMessage struct {
	cause error
	msg   string
}

func (w *withMessage) Error() string {
	return w.msg + ": " + w.cause.Error()
}

func (w *withMessage) Cause() error {
	return w.cause
}

func (w *withMessage) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", w.Cause())
			io.WriteString(s, w.msg)
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, w.Error())
	}
}

func Cause(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}

func WithCode(err error, code string) error {
	if err == nil {
		return nil
	}
	return &ErrWrap{
		cause: err,
		code:  code,
	}
}

func WithCodef(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &ErrWrap{
		cause: err,
		code:  fmt.Sprintf(format, args...),
	}
}

type ErrWrap struct {
	cause error
	code  string
	vars  []string
}

func (w *ErrWrap) Vars() []string {
	return w.vars
}

func (w *ErrWrap) Code() string {
	return w.code
}

func (w *ErrWrap) Error() string {
	var msg string
	if w.cause != nil {
		msg += w.cause.Error()
	}
	return msg
}

func (w *ErrWrap) Cause() error {
	return w.cause
}

// Format rewrite format
func (w *ErrWrap) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", w.Cause())
			io.WriteString(s, "BizCode=["+string(w.code)+"]")
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, w.Error())
	}
}

func BizWrap(err error, code string, message string, vars ...string) error {
	if err == nil {
		return nil
	}
	codeErr := &ErrWrap{
		cause: err,
		code:  code,
		vars:  vars,
	}
	err = &withMessage{
		cause: codeErr,
		msg:   message,
	}
	return &withStack{
		err,
		callers(),
	}
}

func DefaultBizWrap(code string, vars ...string) error {
	err := errors.New("")
	codeErr := &ErrWrap{
		cause: err,
		code:  code,
		vars:  vars,
	}
	err = &withMessage{
		cause: codeErr,
		//msg:   strings.Join(message, ","),
	}
	return &withStack{
		err,
		callers(),
	}
}
