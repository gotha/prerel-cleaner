# PREREL_CLEANER

A tool that can delete releases, pre-releases and tags from github on bulk.


## Build

Get the source code and run

```sh
go build
```

## Install

You need to have git installed.

Copy the binary in your exec path

Ex:
```cp
cp prerel-cleaner /usr/local/bin
```

Tested on OSX, will probably work on Linux, we'll never know if it works on Windows.


## Usage

Go to your git repository and type 

```
prerel-cleaner
```

It will open a file in your favorite editor that looks like this:

```
keep - [PRERELEASE] (1.0.2-test-rc1) 1.0.1-test
keep - (1.0.1-test) 1.0.1-test
keep - (1.0.0) 1.0.0
keep - (0.1.0-test) 0.1.0-test
keep - (0.0.1) 0.0.1

# ---------------------------------
# keep| k - keep release and tag
# del| d - delete both release and tag
# rel| r - delete release but keep the tag
# rel| r - delete release but keep the tag
```

Change the file like this:

```
del - [PRERELEASE](1.0.2-test-rc1) 1.0.1-test
keep - (1.0.1-test) 1.0.1-test
keep - (1.0.0) 1.0.0
keep - (0.1.0-test) 0.1.0-test
keep - (0.0.1) 0.0.1
```

and the marked releases or prereleases will be deleted (with their tags if you have chosen to)



## Disclaimer

This is something that I wrote in an hour on a lazy afternoon. Use at your own risk.
The project uses code from [hub](https://github.com/github/hub).
