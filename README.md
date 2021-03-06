# PREREL_CLEANER

A tool that can delete releases, pre-releases and tags from github on bulk.


## Build

Get the source code and run

```sh
go build
```

## Install

### Prerequisites

You need to have git installed.

### Auth

The easiest authentication method is to create file `~/.config/hub` with the following content

```
github.com:
- user: gotha
  oauth_token: <YOUR_GITHUB_TOKEN>
  protocol: https
```

Otherwise you will be asked for username and password every time.


### Get the binary

You can either install with:

```sh
go install github.com/gotha/prerel-cleaner@latest
```

or get the code, build it yourself and copy the binary in your exec path

```sh
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
