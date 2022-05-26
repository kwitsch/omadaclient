package log

import (
	"fmt"

	"github.com/kwitsch/omadaclient/utils"
)

type Log struct {
	prefix  *string
	verbose *bool
}

func New(prefix string, verbose bool) *Log {
	result := &Log{
		prefix:  &prefix,
		verbose: &verbose,
	}

	return result
}

func (l *Log) M(params ...interface{}) {
	fmt.Println(*l.prefix, utils.ToString(params...))
}

func (l *Log) V(params ...interface{}) {
	if *l.verbose {
		l.M(params...)
	}
}

func (l *Log) E(param interface{}) error {
	var err error
	if e, ok := param.(error); ok {
		err = e
	} else {
		err = utils.NewError(param)
	}
	l.M("Error:", err)
	return err
}

func (l *Log) Return(params ...interface{}) {
	if *l.verbose {
		l.M("Returns:", utils.ToString(params...))
	}
}

func (l *Log) ReturnSuccess() {
	l.Return("Success")
}
