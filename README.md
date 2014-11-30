# nv

The commandline tool for downloading videos and mylist at nicovideo.

## Installation

```session
$ cp ./nv /usr/local/bin/nv
```

## Usage

### Setup

```session
$ nv config email john@example.com
$ nv config password pAsSwoRd
```

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

# Test address

- http://www.nicovideo.jp/watch/sm22538737
- http://www.nicovideo.jp/watch/1341379235?group_id=4139928&ref=my_mylist_s1_p1_n451
- http://www.nicovideo.jp/mylist/33435425
