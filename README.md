# nv

[![Build Status](https://travis-ci.org/uetchy/nv.svg?branch=master)](https://travis-ci.org/uetchy/nv)

The commandline tool for downloading videos and mylist at nicovideo.

## Installation

### Stable version

```session
curl -Ls https://gist.githubusercontent.com/uetchy/b285401a11134d6c3688945b1037cd81/raw/install.sh | bash
```

### HEAD version

```session
$ go get -d github.com/uetchy/nv
```

## Usage

Default config file will be put on **~/.config/nv/config.yml**

## Download

```session
$ nv get http://www.nicovideo.jp/watch/sm22538737
$ nv get http://www.nicovideo.jp/mylist/33435425
```

You also can use more shorten addresses.

```session
$ nv get sm9
$ nv get mylist/33435425
```

### Options

#### Download comments

```session
$ nv get sm9 --with-comments
```

## Show info

```session
$ nv info sm9
$ nv info mylist/33435425
```

## Open video on nicovideo.jp

```session
$ nv browse "./Cat Movie [sm00000].mp4"
```

## Contribution

1.  Fork (<https://github.com/uetchy/nv/fork>)
2.  Create a feature branch
3.  Commit your changes
4.  Rebase your local changes against the master branch
5.  Run test suite with the `go test ./...` command and confirm that it passes
6.  Run `gofmt -s`
7.  Create a new Pull Request
