package main

import "testing"

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
