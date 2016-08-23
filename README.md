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

Links
-------

- [go get / ghq get でのアドレス形式の違いから人類を解放した - pockestrap](http://pocke.hatenablog.com/entry/2016/08/22/170516)

License
-------

These codes are licensed under CC0.

[![CC0](http://i.creativecommons.org/p/zero/1.0/88x31.png "CC0")](http://creativecommons.org/publicdomain/zero/1.0/deed.en)
