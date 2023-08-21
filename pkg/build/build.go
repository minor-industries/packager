package build

import (
	"fmt"
	"github.com/minor-industries/packager/pkg/packager"
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

func Build(cfg *runCfg, exe string, args ...string) error {
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

func BuildSingle(req *packager.BuildRequest) error {
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

	if err := Build(&runCfg{
		Env: []string{"GOOS=linux", "GOARCH=arm64"},
		Dir: filepath.Join("cmd", req.Name),
	},
		"go",
		"build",
		"-o", filepath.Join(tmp, "bin")+"/",
		".",
	); err != nil {
		return err
	}

	outputFile := fmt.Sprintf("%s/%s_%s_%s.tar.gz", req.Arch, req.Name, req.Version, req.Arch)
	fullPathOutputFile := filepath.Join(req.Folder, "builds", outputFile)

	fmt.Println("output file:", fullPathOutputFile)

	if Exists(fullPathOutputFile) {
		return errors.New("output file exists")
	}

	if err := Build(&runCfg{Dir: tmp},
		"tar",
		"-czv",
		"-f", fullPathOutputFile,
		"bin",
	); err != nil {
		return errors.Wrap(err, "tar")
	}

	return nil
}
