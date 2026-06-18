package main

import (
	"code/internal/compare"
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

			var filesData []map[string]any
			for _, f := range s.Files {
				fileData, err := f.CreateMapFromData()
				if err != nil {
					return err
				}
				filesData = append(filesData, fileData)
			}

			fds, err := compare.GenDiff(filesData, 1)
			if err != nil {
				return err
			}

			f, err := formatters.NewFormatter(cmd.String("format"))
			if err != nil {
				return err
			}

			result, err := f.Format(fds)
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
