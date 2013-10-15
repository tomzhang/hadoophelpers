package main

import (
	"errors"
	//"github.com/davecgh/go-spew/spew"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/elazarl/hadoophelpers/go/lib/hadoopconf"
	"github.com/jessevdk/go-flags"
	"github.com/elazarl/hadoophelpers/go/lib/table"
	//"github.com/wsxiaoys/terminal"
)

type getOpts struct {
	Verbose bool `short:"v" long:"verbose" description:"Show verbose debug information"`
}

type setOpts struct {
	Verbose bool `short:"v" long:"verbose" description:"Show verbose debug information"`
}

func (o getOpts) Execute(args []string) error {
	if len(args) == 0 {
		return errors.New("get must have nonzero number arguments")
	}
	t := table.New(4)
	c := opt.getConf()
	keys := []string{}
	for _, key := range c.Keys() {
		for _, arg := range args {
			if ok, _ := filepath.Match(arg, key); ok {
				keys = append(keys, key)
				break
			}
		}
	}
	for _, arg := range keys {
		v, src := c.SourceGet(arg)
		if v == "" && src == hadoopconf.NoSource {
			t.Add("", arg, "", "no property")
		} else {
			t.Add(filepath.Base(src.Source), arg, "=", v)
		}
	}
	fmt.Print(t.String())
	return nil
}

func (o setOpts) Execute(args []string) error {
	if len(args) == 0 {
		return errors.New("get must have nonzero number arguments")
	}
	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			return errors.New("set accepts arguments of the form x=y, no '=' in " + arg)
		}
		opt.getConf().SetIfExist(parts[0], parts[1])
	}
	opt.getConf().Save()
	return nil
}

type gOpts struct {
	Get getOpts `command:"get"`
	Set setOpts `command:"set"`
	ConfPath string `short:"c" long:"conf" description:"Set hadoop configuration dir"`
	conf *hadoopconf.HadoopConf
}

func (opt *gOpts) getConf() *hadoopconf.HadoopConf {
	if opt.conf != nil {
		return opt.conf
	}
	var err error
	var p = "."
	if opt.ConfPath != "" {
		p = opt.ConfPath
	} else if os.Getenv("HADOOP_CONF") != "" {
		p = os.Getenv("HADOOP_CONF")
	}
	opt.conf, err = hadoopconf.New(p)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return opt.conf
}

var opt gOpts
var conf *hadoopconf.HadoopConf

func main() {
	if args, err := flags.ParseArgs(&opt, os.Args[1:]); err != nil || len(args) > 0 {
		os.Exit(1)
	}
}