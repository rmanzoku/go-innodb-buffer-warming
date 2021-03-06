package main

import (
	"flag"
	"fmt"
	"io"

	"github.com/rmanzoku/friction"
)

const (
	ExitCodeOK = iota
	ExitCodeParseFlagError
)

type CLI struct {
	outStream, errStream io.Writer
}

func (c *CLI) Run(args []string) int {

	var version bool

	var user string
	var password string
	var host string
	var port string

	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(c.errStream)
	flags.BoolVar(&version, "version", false, "Print version information and quit")

	flags.StringVar(&host, "host", "127.0.0.1", "MySQL host")
	flags.StringVar(&port, "port", "3306", "MySQL port")
	flags.StringVar(&user, "user", "root", "MySQL User")
	flags.StringVar(&password, "password", "", "MySQL Password")

	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeParseFlagError
	}

	if version {
		fmt.Fprintf(c.errStream, "%s version %s\n", Name, Version)
		return ExitCodeOK
	}

	dbname := flags.Arg(0)

	if dbname == "" {
		fmt.Fprint(c.outStream, "Set DB name")
		return ExitCodeParseFlagError
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		user,
		password,
		host,
		port,
		dbname,
	)

	fmt.Println(dsn)

	db := friction.InitDB(dsn)

	tables, _ := friction.ShowTables(db)

	fmt.Fprint(c.outStream, tables)
	for _, t := range tables {
		columns, err := friction.GetIndexColumns(db, t)
		if err != nil {
			panic(err)
		}

		for _, c := range columns {
			friction.WarmUp(db, t, c, 3)
		}
	}

	return ExitCodeOK
}
