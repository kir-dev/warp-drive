# warp-drive

Warp-drive is a self-hosted image hosting service. It was born out of the need
to store images for [our blog](http://kir-dev.sch.bme.hu). Storing a lot of
images in a git repo is not a lot of fun.

## Overview

It stores images in the file system, while metadata is in a database. We are
using postgres. There is a simple page for uploading and searching images. Every
image can be linked by its unique url and an optional width can be provided. If
a width is provided then the image will be resized and the aspect ratio will be
kept. If the width is greater than the original with then the original size will
be used.

## Features

URL structure:

    /<hash-of-image>/[width]

TODO

## Installing

### Prerequisites

1. install the latest [PostgreSQL](http://www.postgresql.org/download/)
2. install [go1.2+](http://golang.org/doc/install#download)

### Setting up the database

Assuming a running and functional postgresql server.

    $ sudo su - postgres
    $ createuser -l -E -P -R -S -D warp
    $ createdb -O warp -E utf8 warp
    $ psql -U warp -d warp -h localhost -f /path/to/warp/scripts/sql/schema.sql

### Setting up your go environment

Assuming you have go installed, create your workspace.

    $ mkdir -p /path/to/your/workspace/src
    $ export GOPATH=$GOPATH:/path/to/your/workspace

The src folder at the end is mandatory.

### Get the code

    $ cd /path/to/your/workspace/src
    $ git clone https://github.com/kir-dev/warp-drive.git warp-drive

### Setup your config file

    $ cp config/config.json.dist config/config.json
    $ vim config/config.json

For detailed information about the configuration look at the
[configuration](#configuration) section.

### Build & run

We are using [godep](https://github.com/tools/godep) for managing dependencies,
so you must have `godep` installed in your `PATH`.

To build the bot itself just run

    $ make

It creates a new executable named `warp`. To run it simply:

    $ ./warp

## Configuration

TODO

## Contributing

When committing go code **always** use the `go fmt` tool first. Possibly one could
set up a pre-commit git-hook to automate this.

Or you can do it manually:

    $ go fmt ./...
    # or
    $ make fmt

## Adding a new dependency

We are using [godep](https://github.com/tools/godep) for managing dependencies,
so you must have `godep` installed in your `PATH`.

Use the [godep workflow](https://github.com/tools/godep#add-or-update-a-dependency)
and use the `-copy=false` option on save.
