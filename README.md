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

## Main features

* Upload and store images
* Easily link to images

    URL structure:

        //example.com/<hash-of-image>/[width]

* Resize images on the fly
* Search among the uploaded images

## Install

1. Download the [latest release](https://github.com/kir-dev/warp-drive/releases/latest).
2. Extract the archive
3. Create or upgrade the database
    * for a clean install follow the instructions [here](#setting-up-the-database).
    * for an upgrade use the `config/upgrade.sql` if any.
4. start your application:

        $ ./warp
        # for more information on the command-line arguments
        $ ./warp --help

Using it with a reverse HTTP proxy is probably the easiest way to provide SSL
support. (Currently the application only support HTTP.)

## Development

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

Open http://localhost:8080 in your browser.

### Modifying the database schema

When modifying the database schema create a new `.sql` file in `scripts/sql`
directory. The name should contain the date and the purpose of the modification.
Eg: `2014-04-09-created-users-table.sql`.

**DO NOT** forget to add the the modifications to the `schema.sql` file as well.
The schema file should always contain the full schema for a clean install.

## Configuration

Rename or create a copy of the `config.json.dist`. `config/config.json` is
recognized by default, use the `-config` option to provide an alternate config
file.

Currently the config file contains the following options:

* `uploadPath`: base directory for image uploads
* `serverAddress`: the host (including port if necessary)
* `user`, `password`: the user name and password for uploading images
* `db`: connection information for the database

## Contributing

1. Fork it.
2. Create a branch (`git checkout -b my_awesome_path`)
3. Commit your changes (`git commit -am "Awesome stuff added"`)
4. Push to the branch (`git push origin my_awesome_path`)
5. Open a [Pull Request][1]
6. Enjoy a refreshing Diet Coke and wait

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

## Releasing

For versioning we try to follow the rules of [semver](http://semver.org/). This
means that every release has a `MAJOR.MINOR.PATCH` versioning scheme.

1. Create a new distribution archive

        $ ./scripts/dist.sh MAJOR.MINOR.PATCH

2. Tag the new release

        $ git tag -a vMAJOR.MINOR.PATCH

    In the tag message describe the new release briefly.

3. Prepare a new [github release][2]. Describe the new release in detail and
upload the distribution archive for the release.

[1]: https://github.com/kir-dev/warp-drive/pulls
[2]: https://github.com/blog/1547-release-your-software
