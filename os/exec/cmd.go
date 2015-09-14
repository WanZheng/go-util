package exec

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

func Run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

type ParseOutputFn func(io.Reader) error

func GetOutputFn(fn ParseOutputFn, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	out, err := cmd.StdoutPipe()
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to create pipe to command %v: %v", name, err))
	}
	if err := cmd.Start(); err != nil {
		return errors.New(fmt.Sprintf("Failed to start %v: %v", name, err))
	}
	retErr := fn(out)
	if err := cmd.Wait(); err != nil {
		return errors.New(fmt.Sprintf("Failed to wait %v: %v", name, err))
	}
	return retErr
}

func GetOutput(name string, args ...string) ([]byte, error) {
	var buf []byte
	fn := func(r io.Reader) (err error) {
		buf, err = ioutil.ReadAll(r)
		return err
	}
	err := GetOutputFn(fn, name, args...)
	return buf, err
}
