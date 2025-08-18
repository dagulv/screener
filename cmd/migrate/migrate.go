package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/dagulv/screener/internal/env"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := start(ctx, os.Args[1:]); err != nil {
		panic(err)
	}
}

func start(_ context.Context, args []string) (err error) {
	if len(args) < 1 {
		return fmt.Errorf("too few arguments: expected 1, got %d", len(args))
	}

	env, err := env.Load(true)

	if err != nil {
		return
	}

	var path string

	// If go.mod is found
	if path, err = findDir("go.mod"); err == nil {
		path = filepath.Join(path, "internal", "adapter", "postgres", "migrations")
	} else {
		// Otherwise, migrations might exist alongside the executable
		path, err = os.Executable()
	}

	if err != nil {
		return
	}

	path = filepath.Join(filepath.Dir(path), "migrations")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return errors.New("migration files not found")
	}

	m, err := migrate.New("file://"+path, env.PostgresConn)

	if err != nil {
		return
	}

	defer m.Close()

	if len(args) > 0 {
		switch args[0] {
		case "up":
			var steps int

			if len(args) > 1 {
				steps, _ = strconv.Atoi(args[1])
			}

			if steps <= 0 {
				return m.Up()
			}

			return m.Steps(steps)

		case "down":
			var steps int

			if len(args) > 1 {
				steps, _ = strconv.Atoi(args[1])
			}

			if steps <= 0 {
				return m.Down()
			}

			return m.Steps(-steps)

		case "force":
			var v int

			if len(args) > 1 {
				v, _ = strconv.Atoi(args[1])
			}

			if v <= 0 {
				return
			}

			return m.Force(v)

		case "create":
			var name string

			if len(args) > 1 {
				exp := regexp.MustCompile(`[^a-z0-9]+`)
				name = exp.ReplaceAllString(strings.ToLower(args[1]), "_")
				name = strings.Trim(name, "_")
			}

			if len(name) == 0 {
				return
			}

			files, e := os.ReadDir(path)

			if e != nil {
				return e
			}

			var filename string
			var nr int

			if len(files) > 0 {
				for _, file := range files {
					if file.IsDir() {
						continue
					}

					if file.Name() > filename {
						filename = file.Name()
					}
				}

				parts := strings.Split(filename, "_")

				if nr, e = strconv.Atoi(parts[0]); e != nil {
					return e
				}
			}

			os.WriteFile(fmt.Sprintf("%s/%06d_%s.up.sql", path, nr+1, name), nil, 0644)
			os.WriteFile(fmt.Sprintf("%s/%06d_%s.down.sql", path, nr+1, name), nil, 0644)

			return
		}
	}

	return m.Up()
}

func findDir(filename string) (dir string, err error) {
	dir, err = os.Getwd()

	if err != nil {
		return
	}

	for dir != "." {
		if _, err = os.Stat(filepath.Join(dir, filename)); os.IsNotExist(err) {
			dir = filepath.Dir(dir)
			continue
		}

		return
	}

	err = fmt.Errorf("no %s file found", filename)

	return
}
