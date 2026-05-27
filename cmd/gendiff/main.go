package main

import (
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
			path1 := cmd.Args().Get(0)
			path2 := cmd.Args().Get(1)

			s, err := storage.NewStorage(path1, path2)
			if err != nil {
				return err
			}

			res, err := parser.Parse(s, cmd.StringArg("format"))
			if err != nil {
				return err
			}
			fmt.Println(res)
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
