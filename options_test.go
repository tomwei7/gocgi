package gocgi

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoCGIOptions(t *testing.T) {
	flagSet := flag.NewFlagSet("", flag.PanicOnError)
	opts := NewGoCGIOptions()
	opts.BindFlags(flagSet)

	flagSet.Set("addr", "127.0.0.1:8080")
	flagSet.Set("path", "cgi-bin")
	flagSet.Set("root", "root")
	flagSet.Set("dir", "dir")
	flagSet.Set("env", "k1=v1,k2=v2")
	flagSet.Set("inherit-env", "k3,k4")
	flagSet.Set("args", "a1,a2")
	flagSet.Set("stderr", "error.log")
	flagSet.Set("static-map", "static=/img,js=/js")
	flagSet.Set("users", "admin:admin,root:toor")

	assert.Equal(t, &GoCGIOptions{
		Addr:       "127.0.0.1:8080",
		Path:       "cgi-bin",
		Root:       "root",
		Dir:        "dir",
		Env:        []string{"k1=v1", "k2=v2"},
		InheritEnv: []string{"k3", "k4"},
		Args:       []string{"a1", "a2"},
		Stderr:     "error.log",
		StaticMap:  map[string]string{"static": "/img", "js": "/js"},
		Users:      []string{"admin:admin", "root:toor"},
	}, opts)
}
