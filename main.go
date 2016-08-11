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
		fmt.Fprintln(os.Stderr, err)
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

	return fn(c.Args, c.Debug)
}

type CmdArg struct {
	Name  string // get
	Debug bool
	Type  string   // go or ghq
	Args  []string // [-u github.com/pocke/get]
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

	opts := args[1:typeIdx]
	cmdArg.Debug = includeString(opts, "-d") || includeString(opts, "--debug")

	cmdArg.Type = args[typeIdx]
	cmdArg.Args = args[typeIdx+1:]

	return cmdArg, nil
}

var Getters = map[string]func(addrs []string, debug bool) error{
	"ghq": func(addrs []string, debug bool) error {
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
		if debug {
			fmt.Println(strings.Join(c.Args, " "))
		}
		return c.Run()
	},
	"go": func(addrs []string, debug bool) error {
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
		if debug {
			fmt.Println(strings.Join(c.Args, " "))
		}
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
	mats := []AddrMatcher{
		// https://github.com/pocke/get
		// https://github.com/pocke/get.git
		{`^https://([^/]+)/([^/]+)/([^/]+?)(?:\.git)?$`, 1, 2, 3},
		// git@github.com:pocke/get.git
		{`^git\@([^:]+):([^/]+)/(.+)\.git$`, 1, 2, 3},
		// github.com/pocke/get
		{`^([^/]+)/([^/]+)/([^/]+)$`, 1, 2, 3},
		// pocke/get
		{`^([^/]+)/([^/]+)$`, -1, 1, 2},
	}
	for _, m := range mats {
		addr := m.Parse(addrStr)
		if addr == nil {
			continue
		}
		return addr, nil
	}

	return nil, fmt.Errorf("Can't parse %s as address", addrStr)
}

func includeString(s []string, t string) bool {
	for _, v := range s {
		if v == t {
			return true
		}
	}
	return false
}

type AddrMatcher struct {
	Pattern     string
	HostIdx     int
	UserIdx     int
	RepoNameIdx int
}

func (m *AddrMatcher) Parse(addrStr string) *Addr {
	re := regexp.MustCompile(m.Pattern)
	ma := re.FindStringSubmatch(addrStr)
	if len(ma) == 0 {
		return nil
	}

	addr := &Addr{
		Host: "github.com",
	}

	if m.HostIdx > 0 {
		addr.Host = ma[m.HostIdx]
	}
	if m.UserIdx > 0 {
		addr.User = ma[m.UserIdx]
	}
	if m.RepoNameIdx > 0 {
		addr.RepoName = ma[m.RepoNameIdx]
	}

	return addr
}
