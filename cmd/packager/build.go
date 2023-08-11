package main

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type runCfg struct {
	Env []string
	Dir string
}

func Run(cfg *runCfg, exe string, args ...string) error {
	if cfg == nil {
		cfg = &runCfg{}
	}

	cmdInfo := strings.TrimSpace(strings.Join([]string{
		strings.Join(cfg.Env, " "),
		exe,
		strings.Join(args, " "),
	}, " "))

	if cfg.Dir != "" {
		cmdInfo = fmt.Sprintf("(cd '%s' && %s)", cfg.Dir, cmdInfo)
	}

	fmt.Println(cmdInfo)

	cmd := exec.Command(exe, args...)

	cmd.Env = append(os.Environ(), cfg.Env...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = cfg.Dir
	err := cmd.Run()
	return errors.Wrap(err, "run")
}

func buildSingle(
	name string,
	version string,
	arch string,
	sharedFolder string,
) error {
	tmp, err := os.MkdirTemp("", "")
	if err != nil {
		return errors.Wrap(err, "mkdirtemp")
	}
	defer func() {
		err := os.RemoveAll(tmp)
		if err != nil {
			panic(err)
		}
	}()

	if err := Run(&runCfg{
		Env: []string{"GOOS=linux", "GOARCH=arm64"},
		Dir: filepath.Join("cmd", name),
	},
		"go",
		"build",
		"-o", filepath.Join(tmp, "bin")+"/",
		".",
	); err != nil {
		return err
	}

	outputFile := fmt.Sprintf("%s/%s_%s_%s.tar.gz", arch, name, version, arch)
	fullPathOutputFile := filepath.Join(sharedFolder, "builds", outputFile)

	fmt.Println("output file:", fullPathOutputFile)

	if Exists(fullPathOutputFile) {
		return errors.New("output file exists")
	}

	if err := Run(&runCfg{Dir: tmp},
		"tar",
		"-czv",
		"-f", fullPathOutputFile,
		"bin",
	); err != nil {
		return errors.Wrap(err, "tar")
	}

	return nil
}
