package appcli

import (
	"fmt"
	"sse_demo/database"

	"github.com/urfave/cli/v2"
)

//StartFunc 开启节点
func MigrateFunc(c *cli.Context) error {
	_ = database.GetDatabase()
	fmt.Println("do migrate suc")
	return nil
}

func Migrate() *cli.Command {
	return &cli.Command{
		Name:  "migrate",
		Usage: fmt.Sprintf("start %s server", ""),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "user",
				Aliases: []string{"u"},
				Usage:   "eosc",
			},
			&cli.StringFlag{
				Name:    "group",
				Aliases: []string{"g"},
				Usage:   "eosc",
			},
		},
		Action: MigrateFunc,
	}
}
