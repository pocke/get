package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	if err := Main(os.Args); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

func Main(args []string) error {
	c, err := ParseCmdArg(args)
	if err != nil {
		return err
	}

	fn, ok := Getters[c.Type]
	if !ok {
		return fmt.Errorf("Type %s doesn't exist", c.Type)
	}

	return fn(c.Args)
}

type CmdArg struct {
	Name    string   // get
	Options []string // [--debug -d]
	Type    string   // go or ghq
	Args    []string // [-u github.com/pocke/get]
}

func ParseCmdArg(args []string) (*CmdArg, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("Too few arguments. Please `get [-d] TYPE ADDR ...`")
	}
	cmdArg := new(CmdArg)
	cmdArg.Name = args[0]

	typeIdx := 0
	for idx, v := range args[1:] {
		if strings.HasPrefix(v, "-") {
			continue
		}
		typeIdx = idx + 1
		break
	}
	// TODO: check index
	cmdArg.Options = args[1:typeIdx]
	cmdArg.Type = args[typeIdx]
	cmdArg.Args = args[typeIdx+1:]

	return cmdArg, nil
}

var Getters = map[string]func(addrs []string) error{
	"ghq": func(addrs []string) error {
		args := []string{"get"}
		for _, a := range addrs {
			addr, err := ParseAddr(a)
			if err != nil {
				args = append(args, a)
			} else {
				args = append(args, addr.ToSSH())
			}
		}

		c := exec.Command("ghq", args...)
		c.Stdin = os.Stdin
		c.Stderr = os.Stderr
		c.Stdout = os.Stdout
		return c.Run()
	},
	"go": func(addrs []string) error {
		args := []string{"get"}
		for _, a := range addrs {
			addr, err := ParseAddr(a)
			if err != nil {
				args = append(args, a)
			} else {
				args = append(args, addr.ToGoStyle())
			}
		}
		c := exec.Command("go", args...)
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
	// Dir string
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
