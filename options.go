package gocgi

import (
	"flag"
	"fmt"
	"strings"
)

type FlagStrings []string

func (f *FlagStrings) String() string {
	return strings.Join(*f, ",")
}

func (f *FlagStrings) Set(s string) error {
	if s == "" {
		return nil
	}
	*f = strings.Split(s, ",")
	return nil
}

type FlagMap map[string]string

func (f *FlagMap) String() string {
	ss := make([]string, 0, len(*f))
	for k, v := range *f {
		ss = append(ss, k+"="+v)
	}
	return strings.Join(ss, ",")
}

func (f *FlagMap) Set(s string) error {
	kvs := strings.Split(s, ",")
	for _, kv := range kvs {
		if kv == "" {
			continue
		}
		ss := strings.SplitN(kv, "=", 2)
		if len(ss) != 2 {
			return fmt.Errorf("invalid kv pair %s", kv)
		}
		(*f)[ss[0]] = ss[1]
	}
	return nil
}

type GoCGIOptions struct {
	Addr       string
	Path       string
	Root       string
	Dir        string
	Env        FlagStrings
	InheritEnv FlagStrings
	Args       FlagStrings
	Stderr     string
	StaticMap  FlagMap
	Users      FlagStrings
}

func (g *GoCGIOptions) BindFlags(flagSet *flag.FlagSet) {
	flagSet.StringVar(&g.Addr, "addr", "127.0.0.1:6080", "listen addr")
	flagSet.StringVar(&g.Path, "path", "", "path to the CGI executable")
	flagSet.StringVar(&g.Root, "root", "", "root URI prefix of handler or empty for \"/\"")
	flagSet.StringVar(&g.Dir, "dir", "", "Dir specifies the CGI executable's working directory. If Dir is empty, the base directory of Path is used. If Path has no base directory, the current working  directory is used.")
	flagSet.Var(&g.Env, "env", `extra environment variables to set, if any, as "key=value"`)
	flagSet.Var(&g.InheritEnv, "inherit-env", `environment variables to inherit from host, as "key"`)
	flagSet.Var(&g.Args, "args", "optional arguments to pass to child process")
	flagSet.StringVar(&g.Stderr, "stderr", "", "redirect child process stderr to file; empty means stderr")
	flagSet.Var(&g.StaticMap, "static-map", `map local path as static path as path=localpath`)
	flagSet.Var(&g.Users, "users", "http basic auth")
}

func NewGoCGIOptions() *GoCGIOptions {
	return &GoCGIOptions{StaticMap: make(map[string]string)}
}
