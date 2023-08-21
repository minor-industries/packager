package build

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
)

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		panic(err)
	}
	return true
}

func Cd(
	dir string,
	callback func() error,
) error {
	pwd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "getwd")
	}

	defer func() {
		if err := os.Chdir(pwd); err != nil {
			panic(errors.Wrap(err, "chdir2"))
		}
	}()

	fmt.Println("cd", dir)
	if err := os.Chdir(dir); err != nil {
		return errors.Wrap(err, "chdir")
	}

	if err := callback(); err != nil {
		return errors.Wrap(err, "callback")
	}

	return nil
}
