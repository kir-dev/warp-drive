#!/bin/bash
# usage ./dist.sh

TEMPDIR=/tmp/warp-dist

DISTHASH=$(git rev-parse HEAD | cut -c1-10)
DISTNAME="warp-dist-$DISTHASH.tar.gz"
DISTPATH="dist"

# test first, failing test means no dist
make -s test > /dev/null
if [ $? -ne 0 ]; then
    echo -e '\e[31mtests failed, fix the tests before creating a distribution pacakge\e[0m'
    exit 1
fi

# remove temp directory
if [ -d $TEMPDIR ]; then
    rm -r $TEMPDIR
fi

make -s warp

# create temp dir
mkdir -p $TEMPDIR/config

# copy resources to temp dir
cp -r warp static/ template/ $TEMPDIR
cp config/config.json.dist $TEMPDIR/config/
cp scripts/sql/*.sql $TEMPDIR/config/

# create dist directory
mkdir -p $DISTPATH
# create dist archive
tar -C $TEMPDIR -czf $DISTPATH/$DISTNAME .

rm -r $TEMPDIR
