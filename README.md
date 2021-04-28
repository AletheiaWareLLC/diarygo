![Diary Icon](./cmd/diaryfyne/icon.svg)

diarygo
=======

Diary is a personal journal written in Go that stores entries in S P A C E.

# GUI

![Diary Locked GUI Screenshot](./img/gui-locked.png)

![Diary Add GUI Screenshot](./img/gui-add.png)

## Installation

```
go get fyne.io/fyne/v2/fyne
fyne install aletheiaware.com/diarygo/cmd/diaryfyne
```

## Usage

```
diaryfyne
```

# CLI

## Installation

```
go get aletheiaware.com/diarygo/cmd/diary
```

## Usage

```
Diary Usage:
    diary - display usage
    diary add - add a new diary entry from stdin
    diary add [file] - add a new diary entry from file
    diary list - display all entries
    diary show - display latest entry
    diary show [id] - display entry with given id
```