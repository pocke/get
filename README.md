get
============

A wrapper of ghq and go.

[![Build Status](https://travis-ci.org/pocke/get.svg?branch=master)](https://travis-ci.org/pocke/get)
[![Coverage Status](https://coveralls.io/repos/github/pocke/get/badge.svg?branch=master)](https://coveralls.io/github/pocke/get?branch=master)


Installation
-----------

```sh
go get github.com/pocke/get
```

<!-- Or download a binary from [Latest release](https://github.com/pocke/get/releases/latest). -->


Usage
-----------

```sh
$ get TYPE ADDRESS
```

`get` supports 2 types.

- [go](https://golang.org/cmd/go/#hdr-Download_and_install_packages_and_dependencies)
- [ghq](https://github.com/motemen/ghq)

`get` supports 4 styles of address.

- `https://github.com/pocke/get`
- `https://github.com/pocke/get.git`
- `github.com/pocke/get.git`
- `git@github.com:pocke/get.git`

### For example

```sh
$ get go https://github.com/pocke/get
$ get ghq github.com/pocke/get.git
```

Advanced Usage
---

`Get` supports `-shallow` and `-unshallow` options. They work with `ghq`. `Get` just ignores them with `go`.

When `-shallow` option is given, `get` clones the specified repository shallowly.
When `-unshallow` option is given, `get` executes `git fetch --unshallow` asynchronously.

They improves cloning speed. If you specify `-shallow` and `-unshallow`, you can clone repository faster, and get whole repository after a while.
For example:

```bash
$ get -shallow -unshallow ghq https://github.com/pocke/get
```

If you'd like to enable this feature by default, put a config file to `~/.config/get/args` with the below content.

```
-shallow -unshallow
```

Links
-------

- [go get / ghq get でのアドレス形式の違いから人類を解放した - pockestrap](http://pocke.hatenablog.com/entry/2016/08/22/170516)

License
-------

These codes are licensed under CC0.

[![CC0](http://i.creativecommons.org/p/zero/1.0/88x31.png "CC0")](http://creativecommons.org/publicdomain/zero/1.0/deed.en)
