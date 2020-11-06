// Package pcsinit 初始化配置包
package pcsinit

import (
	"fmt"
	"github.com/urfave/cli"
	_ "unsafe" // for go:linkname
)

//go:linkname helpCommand1 github.com/iikira/BaiduPCS-Go/vendor/github.com/urfave/cli.helpCommand
//go:linkname helpCommand2 github.com/urfave/cli.helpCommand
var (
	helpCommand1 cli.Command
	helpCommand2 cli.Command
)

func init() {
	cli.AppHelpTemplate = `----
	{{.Name}}{{if .Usage}} - {{.Usage}}{{end}}

USAGE:
	{{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Version}}{{if not .HideVersion}}

VERSION:
	{{.Version}}{{end}}{{end}}{{if .Description}}

DESCRIPTION:
	{{.Description}}{{end}}{{if len .Authors}}

AUTHOR{{with $length := len .Authors}}{{if ne 1 $length}}S{{end}}{{end}}:
	{{range $index, $author := .Authors}}{{if $index}}
	{{end}}{{$author}}{{end}}{{end}}{{if .VisibleCommands}}

COMMANDS:{{range .VisibleCategories}}{{if .Name}}
	{{.Name}}:{{end}}{{range .VisibleCommands}}
		{{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{end}}{{end}}{{if .VisibleFlags}}

GLOBAL OPTIONS:
	{{range $index, $option := .VisibleFlags}}{{if $index}}
	{{end}}{{$option}}{{end}}{{end}}{{if .Copyright}}

COPYRIGHT:
	{{.Copyright}}{{end}}
`

	cli.CommandHelpTemplate = `----
	{{.HelpName}} - {{.Usage}}

USAGE:
	{{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}}{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Category}}

CATEGORY:
	{{.Category}}{{end}}{{if .Description}}

DESCRIPTION:
	{{.Description}}{{end}}{{if .VisibleFlags}}

OPTIONS:
	{{range .VisibleFlags}}{{.}}
	{{end}}{{end}}
`

	cli.SubcommandHelpTemplate = `----
	{{.HelpName}} - {{.Usage}}

USAGE:
	{{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} command{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Description}}

DESCRIPTION:
	{{.Description}}{{end}}

COMMANDS:{{range .VisibleCategories}}{{if .Name}}
	{{.Name}}:{{end}}{{range .VisibleCommands}}
		{{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}
{{end}}{{if .VisibleFlags}}
OPTIONS:
	{{range .VisibleFlags}}{{.}}
	{{end}}{{end}}
`

	helpCommand1.Aliases = append(helpCommand1.Aliases, "?", "？")
	helpCommand1.Action = func(c *cli.Context) error {
		args := c.Args()
		if args.Present() {
			err := cli.ShowCommandHelp(c, args.First())
			if err != nil {
				fmt.Printf("%s\n", err)
			}
			return nil
		}

		cli.ShowAppHelp(c)
		return nil
	}

	helpCommand2.Aliases = helpCommand1.Aliases
	helpCommand2.Action = helpCommand1.Action
}
