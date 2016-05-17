# nv

The commandline tool for downloading videos and mylist at nicovideo.

## Installation

```session
$ go get -d github.com/uetchy/nv
```

## Usage

### Setup

```session
$ nv config email john@example.com
$ nv config password pAsSwoRd
```

# Default config file will be put on **~/.config/nv/config.yml**

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

#### Directory

```session
$ nv get sm9 --with-dir
$ nv get mylist/33435425 --without-dir
```

#### Comments

```session
$ nv get sm9 --with-comments
```

## Audit

```session
$ nv info sm9
$ nv info mylist/33435425
```

## Contribution

1. Fork (<https://github.com/uetchy/nv/fork>)
2. Create a feature branch
3. Commit your changes
4. Rebase your local changes against the master branch
5. Run test suite with the `go test ./...` command and confirm that it passes
6. Run `gofmt -s`
7. Create a new Pull Request
