# cropper

A simple image croper. Makes use of the supersmart [smartcrop]("github.com/muesli/smartcrop") lib.


## usage
```
cropper -wh 500x300 /path/
or
cropper -wh 500x300 /path/to/image/file
```


## install
```
go get -u github.com/denkhaus/cropper
```

## options

```
NAME:
   cropper - A simple image cropper

USAGE:
   cropper [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --wh value     width and hight of new image in format <with>x<height> (default: "580x434")
   --help, -h     show help
   --version, -v  print the version

```