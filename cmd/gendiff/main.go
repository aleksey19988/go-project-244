package main

import (
	"code/internal/compare"
	fls "code/internal/files"
	"code/internal/formatters"
	"code/internal/parser"
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

			files, err := storage.GetFilesData(path1, path2)
			if err != nil {
				return err
			}

			err = fls.Validate(files)
			if err != nil {
				return err
			}

			var filesData []map[string]any
			for _, f := range files {
				fileData, err := parser.ParseData(*f)
				if err != nil {
					return err
				}
				filesData = append(filesData, fileData)
			}

			fields := compare.GenDiff(filesData[0], filesData[1])

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
