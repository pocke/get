package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
)

func main() {
	if err := Main(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func Main(args []string) error {
	configArgs, err := LoadConfig()
	if err != nil {
		return err
	}

	args2 := []string{args[0]}
	args2 = append(args2, configArgs...)
	args2 = append(args2, args[1:]...)

	c, err := ParseCmdArg(args2)
	if err != nil {
		return err
	}

	fn, ok := Getters[c.Type]
	if !ok {
		return fmt.Errorf("Type %s doesn't exist", c.Type)
	}

	return fn(c)
}

type CmdArg struct {
	Name      string // get
	Debug     bool
	Shallow   bool
	Unshallow bool
	Type      string   // go or ghq
	Args      []string // [-u github.com/pocke/get]
}

func LoadConfig() ([]string, error) {
	confPath, err := homedir.Expand("~/.config/get/args")
	if err != nil {
		return nil, err
	}
	if !FileExists(confPath) {
		return []string{}, nil
	}
	b, err := ioutil.ReadFile(confPath)
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimRight(string(b), "\n"), " "), nil

}

func ParseCmdArg(args []string) (*CmdArg, error) {
	cmdArg := new(CmdArg)
	fs := flag.NewFlagSet("get", flag.ContinueOnError)
	fs.BoolVar(&cmdArg.Shallow, "shallow", false, "Shallow clone")
	fs.BoolVar(&cmdArg.Unshallow, "unshallow", false, "Make the repository unshallow after cloned")
	fs.BoolVar(&cmdArg.Debug, "debug", false, "Debug mode")
	fs.Parse(args[1:])
	cmdArg.Name = args[0]

	parsedArgs := fs.Args()

	if len(parsedArgs) == 0 {
		return nil, fmt.Errorf("Too few arguments. Please `get [-d] TYPE ADDR ...`")
	}

	cmdArg.Type = parsedArgs[0]
	cmdArg.Args = parsedArgs[1:]

	return cmdArg, nil
}

var Getters = map[string]func(opt *CmdArg) error{
	"ghq": func(opt *CmdArg) error {
		args := []string{"get"}
		if opt.Shallow {
			args = append(args, "--shallow")
		}
		addrs := make([]*Addr, 0, len(opt.Args))
		for _, a := range opt.Args {
			addr, err := ParseAddr(a)
			if err != nil {
				args = append(args, a)
			} else {
				addrs = append(addrs, addr)
				args = append(args, addr.ToSSH())
			}
		}

		c := exec.Command("ghq", args...)
		c.Stdin = os.Stdin
		c.Stderr = os.Stderr
		c.Stdout = os.Stdout
		if opt.Debug {
			fmt.Println(strings.Join(c.Args, " "))
		}
		err := c.Run()
		if opt.Unshallow {
			for _, addr := range addrs {
				b, err := exec.Command("ghq", "list", "-e", "-p", addr.ToGoStyle()).Output()
				if err != nil {
					return err
				}

				unshallowCmd := exec.Command("git", "fetch", "--unshallow")
				unshallowCmd.Dir = strings.TrimRight(string(b), "\n")
				err = unshallowCmd.Start()
				if err != nil {
					return err
				}
			}
		}
		return err
	},
	"go": func(opt *CmdArg) error {
		args := []string{"get"}
		for _, a := range opt.Args {
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
		if opt.Debug {
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

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
