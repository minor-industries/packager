package main

import (
	"database/sql"
	"fmt"
	"github.com/bitfield/script"
	"github.com/google/uuid"
	"github.com/jessevdk/go-flags"
	"github.com/minor-industries/packager/pkg/db"
	"github.com/pkg/errors"
	"os"
	"strings"
)

var opts struct {
	Minor        bool   `long:"minor"`
	AllowDirty   bool   `long:"allow-dirty"`
	New          bool   `long:"new"`
	Arch         string `long:"arch" required:"true"`
	SharedFolder string `long:"shared-folder" default:"$HOME/shared"`
}

func run() error {
	args, err := flags.Parse(&opts)
	if err != nil {
		return err
	}

	if !opts.Minor {
		return errors.New("only minor releases supported for now")
	}

	switch opts.Arch {
	case "arm64":
	default:
		return errors.New("unknown arch")
	}

	opts.SharedFolder = os.ExpandEnv(opts.SharedFolder)

	fmt.Println(args)

	dbmap, err := db.Get("localhost", 3306, "cloud_config", db.DBMapInit)
	if err != nil {
		return errors.Wrap(err, "get db")
	}

	ref, err := script.Exec("git describe --tags --dirty").String()
	if err != nil {
		return errors.Wrap(err, "git describe")
	}

	fmt.Println(ref)

	if strings.Contains(ref, "-dirty") && !opts.AllowDirty {
		return errors.New("repo is dirty")
	}

	if ref == "" {
		return errors.New("couldn't determine git ref")
	}

	name := args[0]

	count, err := dbmap.SelectInt("select count(*) from Package where name = :name", map[string]any{"name": name})
	if err != nil {
		return errors.Wrap(err, "count")
	}

	// TODO: make sure we don't have a matching git tag already for this package (via unique index?)

	latest := new(db.Package)
	if count > 0 {
		if opts.New {
			return fmt.Errorf("packages exist for %s but --new specified", name)
		}
		err = dbmap.SelectOne(
			latest,
			"select * from Package where name = :name order by major desc, minor desc, patch desc limit 1",
			map[string]any{
				"name": name,
			},
		)
	} else {
		if !opts.New {
			return fmt.Errorf("no packages exist for %s. Use --new to create one", name)
		}
	}

	fmt.Println(count)

	if err != nil {
		return errors.Wrap(err, "select")
	}
	fmt.Println(latest.Name)

	newPkg := &db.Package{
		ID:    uuid.New().String(),
		Name:  name,
		Major: latest.Major,
		Minor: latest.Minor + 1,
		Patch: 0,
		Arch:  opts.Arch,
		OS:    "linux",
		GitRef: sql.NullString{
			String: strings.TrimSpace(ref),
			Valid:  true,
		},
	}

	version := fmt.Sprintf("%d.%d.%d", newPkg.Major, newPkg.Minor, newPkg.Patch)
	if err := buildSingle(name, version, opts.Arch, opts.SharedFolder); err != nil {
		return errors.Wrap(err, "build single")
	}

	// TODO: prevent duplicate inserts (with a unique index)
	err = dbmap.Insert(newPkg)
	if err != nil {
		return errors.Wrap(err, "insert")
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
