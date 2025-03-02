// package executor declares the interface for a request executor and provides a factory
// for instantiating whichever executor you need.
package executor

import (
	"errors"
	"fmt"
	"strings"

	"postal/config"
	"postal/executor/curl"
	"postal/executor/native"
)

// enum to identify which executor you want
const (
	NativeExecutor ExecutorId = iota
	CurlExecutor
)

const (
	StrNativeExecutor = "native"
	StrCurlExecutor   = "curl"
)

const (
	errMsgInvalid = "invalid executor"
)

type ExecutorId int

type Executor interface {
	Execute(config *config.Config) error
}

func ToExecutorId(name string) (ExecutorId, error) {
	switch strings.ToLower(name) {
	case StrNativeExecutor:
		return NativeExecutor, nil
	case StrCurlExecutor:
		return CurlExecutor, nil
	default:
		return -1, fmt.Errorf("%s name '%s'", errMsgInvalid, name)
	}
}

// 'factory' function that instantiates the executor, executes it, and returns the result.
// TODO: figure out what the result struct needs to be.
func Run(id ExecutorId, cfg *config.Config) error {
	// validate the config?
	// - I may need to take it in, write it to a toml document (string), resolve it, convert it back to config,
	//   and then validate the result...
	//   BECAUSE some things may not validate with the tags left unresolved (e.g. urls)
	// instantiate the executor
	// call it
	// get the result
	// return the result
	switch id {
	case NativeExecutor:
		e := native.NewExecutor()
		return e.Execute(cfg)
	case CurlExecutor:
		e := curl.NewExecutor()
		return e.Execute(cfg)
	default:
		return errors.New(errMsgInvalid)
	}
}

func RunNamed(name string, cfg *config.Config) (err error) {
	var id ExecutorId
	if id, err = ToExecutorId(name); err != nil {
		return err
	}
	return Run(id, cfg)
}
