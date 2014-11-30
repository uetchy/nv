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

Default config file will be put on __$HOME/.config/nv__

### Download

```session
$ nv dl http://www.nicovideo.jp/watch/sm22538737
$ nv dl http://www.nicovideo.jp/mylist/33435425
```

Also you can use more easy way.

```session
$ nv dl sm9
$ nv dl mylist/33435425
```

### Audit

```session
$ nv info sm9
$ nv info mylist/33435425
```
