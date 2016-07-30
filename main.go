package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

func main() {
	if err := Main(os.Args); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

func Main(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("Too few arguments. Please `get TYPE ADDR`")
	}
	t := args[1]
	addr := args[2]

	fn, ok := Getters[t]
	if !ok {
		return fmt.Errorf("Type %s doesn't exist", t)
	}

	return fn(addr)
}

var Getters = map[string]func(addr string) error{
	"ghq": func(addrStr string) error {
		addr, err := ParseAddr(addrStr)
		if err != nil {
			return err
		}
		c := exec.Command("ghq", "get", addr.ToSSH())
		c.Stdin = os.Stdin
		c.Stderr = os.Stderr
		c.Stdout = os.Stdout
		return c.Run()
	},
	"go": func(addrStr string) error {
		addr, err := ParseAddr(addrStr)
		if err != nil {
			return err
		}
		c := exec.Command("go", "get", addr.ToGoStyle())
		c.Stdin = os.Stdin
		c.Stderr = os.Stderr
		c.Stdout = os.Stdout
		return c.Run()
	},
}

type Addr struct {
	Host     string
	User     string
	RepoName string
	// TODO: Parse dir. e.g.) https://github.com/pocke/lemonade/tree/master/client
	Dir string
}

func (a *Addr) ToSSH() string {
	return fmt.Sprintf("git@%s:%s/%s.git", a.Host, a.User, a.RepoName)
}

func (a *Addr) ToGoStyle() string {
	return fmt.Sprintf("%s/%s/%s", a.Host, a.User, a.RepoName)
}

func ParseAddr(addrStr string) (*Addr, error) {
	if addr := matchHTTPS(addrStr); addr != nil {
		return addr, nil
	}
	if addr := matchSSH(addrStr); addr != nil {
		return addr, nil
	}
	if addr := matchGoStyle(addrStr); addr != nil {
		return addr, nil
	}
	return nil, fmt.Errorf("Can't parse %s as address", addrStr)
}

// TODO: DRY

// https://github.com/pocke/get
// https://github.com/pocke/get.git
func matchHTTPS(addrStr string) *Addr {
	re := regexp.MustCompile(`^https://([^/]+)/([^/]+)/([^/]+?)(?:\.git)?$`)
	ma := re.FindStringSubmatch(addrStr)
	if len(ma) == 0 {
		return nil
	}

	return &Addr{
		Host:     ma[1],
		User:     ma[2],
		RepoName: ma[3],
	}
}

// git@github.com:pocke/get.git
func matchSSH(addrStr string) *Addr {
	re := regexp.MustCompile(`^git\@([^:]+):([^/]+)/(.+)\.git$`)
	ma := re.FindStringSubmatch(addrStr)
	if len(ma) == 0 {
		return nil
	}

	return &Addr{
		Host:     ma[1],
		User:     ma[2],
		RepoName: ma[3],
	}
}

// github.com/pocke/get
func matchGoStyle(addrStr string) *Addr {
	re := regexp.MustCompile(`^([^/]+)/([^/]+)/([^/]+)$`)
	ma := re.FindStringSubmatch(addrStr)
	if len(ma) == 0 {
		return nil
	}

	return &Addr{
		Host:     ma[1],
		User:     ma[2],
		RepoName: ma[3],
	}
}
