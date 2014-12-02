# nv

The commandline tool for downloading videos and mylist at nicovideo.

## Installation

```session
$ gem install nv
```

## Usage

### Setup

```session
$ nv config email john@example.com
$ nv config password pAsSwoRd
```

Default config file will be put on __~/.config/nv__

### Download

```session
$ nv dl http://www.nicovideo.jp/watch/sm22538737
$ nv dl http://www.nicovideo.jp/mylist/33435425
```

You also can use more shorten addresses.

```session
$ nv dl sm9
$ nv dl mylist/33435425
```

#### Options

##### Directory

```session
$ nv dl sm9 --with-dir
$ nv dl mylist/33435425 --without-dir
```

##### Comments

```session
$ nv dl sm9 --with-comments
```

### Audit

```session
$ nv info sm9
$ nv info mylist/33435425
```
