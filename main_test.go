package main

import (
	"reflect"
	"testing"
)

func TestParseCmdArg_Simple(t *testing.T) {
	args := []string{"get", "go", "github.com/pocke/get"}
	c, err := ParseCmdArg(args)
	if err != nil {
		t.Fatal(err)
	}

	if c.Name != "get" {
		t.Errorf("Name should be `get`, but got %s", c.Name)
	}
	if c.Debug {
		t.Errorf("Options should be empty, but got %v", c.Debug)
	}
	if c.Type != "go" {
		t.Errorf("Type should be `go`, but got %s", c.Type)
	}
	if !reflect.DeepEqual(c.Args, []string{"github.com/pocke/get"}) {
		t.Errorf("Args should have an addr, but got %v", c.Args)
	}
}

func TestParseCmdArg_DebugOption(t *testing.T) {
	args := []string{"get", "--debug", "go", "github.com/pocke/get"}
	c, err := ParseCmdArg(args)
	if err != nil {
		t.Fatal(err)
	}

	if c.Name != "get" {
		t.Errorf("Name should be `get`, but got %s", c.Name)
	}
	if !c.Debug {
		t.Errorf("Options should be empty, but got %v", c.Debug)
	}
	if c.Type != "go" {
		t.Errorf("Type should be `go`, but got %s", c.Type)
	}
	if !reflect.DeepEqual(c.Args, []string{"github.com/pocke/get"}) {
		t.Errorf("Args should have an addr, but got %v", c.Args)
	}
}

func TestParseCmdArg_UpdateOption(t *testing.T) {
	args := []string{"get", "--debug", "go", "-u", "github.com/pocke/get"}
	c, err := ParseCmdArg(args)
	if err != nil {
		t.Fatal(err)
	}

	if c.Name != "get" {
		t.Errorf("Name should be `get`, but got %s", c.Name)
	}
	if !c.Debug {
		t.Errorf("Options should be empty, but got %v", c.Debug)
	}
	if c.Type != "go" {
		t.Errorf("Type should be `go`, but got %s", c.Type)
	}
	if !reflect.DeepEqual(c.Args, []string{"-u", "github.com/pocke/get"}) {
		t.Errorf("Args should have some args, but got %v", c.Args)
	}
}

func TestParseAddr_HTTPS(t *testing.T) {
	addrStr := "https://github.com/pocke/get"
	addr, err := ParseAddr(addrStr)
	if err != nil {
		t.Fatal(err)
	}

	if addr.Host != "github.com" {
		t.Errorf("Host should be `github.com`, but got `%s`", addr.Host)
	}
	if addr.User != "pocke" {
		t.Errorf("User should be `pocke`, but got `%s`", addr.User)
	}
	if addr.RepoName != "get" {
		t.Errorf("User should be `get`, but got `%s`", addr.RepoName)
	}
}

func TestParseAddr_HTTPSWithGit(t *testing.T) {
	addrStr := "https://github.com/pocke/get.git"
	addr, err := ParseAddr(addrStr)
	if err != nil {
		t.Fatal(err)
	}

	if addr.Host != "github.com" {
		t.Errorf("Host should be `github.com`, but got `%s`", addr.Host)
	}
	if addr.User != "pocke" {
		t.Errorf("User should be `pocke`, but got `%s`", addr.User)
	}
	if addr.RepoName != "get" {
		t.Errorf("User should be `get`, but got `%s`", addr.RepoName)
	}
}

func TestParseAddr_SSH(t *testing.T) {
	addrStr := "git@github.com:pocke/get.git"
	addr, err := ParseAddr(addrStr)
	if err != nil {
		t.Fatal(err)
	}

	if addr.Host != "github.com" {
		t.Errorf("Host should be `github.com`, but got `%s`", addr.Host)
	}
	if addr.User != "pocke" {
		t.Errorf("User should be `pocke`, but got `%s`", addr.User)
	}
	if addr.RepoName != "get" {
		t.Errorf("User should be `get`, but got `%s`", addr.RepoName)
	}
}

func TestParseAddr_GoStyle(t *testing.T) {
	addrStr := "github.com/pocke/get"
	addr, err := ParseAddr(addrStr)
	if err != nil {
		t.Fatal(err)
	}

	if addr.Host != "github.com" {
		t.Errorf("Host should be `github.com`, but got `%s`", addr.Host)
	}
	if addr.User != "pocke" {
		t.Errorf("User should be `pocke`, but got `%s`", addr.User)
	}
	if addr.RepoName != "get" {
		t.Errorf("User should be `get`, but got `%s`", addr.RepoName)
	}
}

func TestParseAddr_GitHub(t *testing.T) {
	addrStr := "pocke/get"
	addr, err := ParseAddr(addrStr)
	if err != nil {
		t.Fatal(err)
	}

	if addr.Host != "github.com" {
		t.Errorf("Host should be `github.com`, but got `%s`", addr.Host)
	}
	if addr.User != "pocke" {
		t.Errorf("User should be `pocke`, but got `%s`", addr.User)
	}
	if addr.RepoName != "get" {
		t.Errorf("User should be `get`, but got `%s`", addr.RepoName)
	}
}
