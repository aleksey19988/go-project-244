package main

import (
	"code/internal/diff"
	"code/internal/formatters"
	"code/internal/storage"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "gendiff",
		Usage: "Compares two configuration files and shows a difference.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Value:   "stylish",
				Usage:   "output format",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.NArg() < 2 {
				return cli.Exit("Not enough arguments (must be 2 arguments)", 1)
			}

			path1 := cmd.Args().Get(0)
			path2 := cmd.Args().Get(1)

			s, err := storage.NewStorage(path1, path2)
			if err != nil {
				return err
			}

			err = s.LoadFiles()
			if err != nil {
				return err
			}

			maps, err := s.CreateMapsFromData()
			if err != nil {
				return err
			}

			fields, err := diff.Diff(maps, 1)
			if err != nil {
				return err
			}

			f, err := formatters.NewFormatter(cmd.String("format"))
			if err != nil {
				return err
			}

			result, err := f.Format(fields)
			if err != nil {
				return err
			}

			fmt.Println(result)
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
