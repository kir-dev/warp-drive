#!/bin/bash
# usage ./dist.sh [version]
# if version is not set, using the latest commit's abbreviated hash

DISTHASH=$(git rev-parse HEAD | cut -c1-10)

if [ -n "$1" ]; then
    DISTVER=$1
else
    DISTVER=$DISTHASH
fi

DISTNAME="warp-$DISTVER.tar.gz"
DISTPATH="dist"

APPTEMPDIR=warp-$DISTVER
TEMPDIR="/tmp/$APPTEMPDIR"

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
# copy schema.sql for clean installs
cp scripts/sql/schema.sql $TEMPDIR/config/

# concatenate & copy schema changes for this release
UPGRADESQL=$TEMPDIR/config/upgrade.sql
git diff --name-status $(git describe --abbrev=0) HEAD | grep -E "^A.*\.sql$" | cut -c3- | xargs cat > $UPGRADESQL
# delete if empty
if [ ! -s $UPGRADESQL ]; then
    rm $UPGRADESQL
fi

# create dist directory
mkdir -p $DISTPATH
# create dist archive
tar -C /tmp -czf $DISTPATH/$DISTNAME $APPTEMPDIR

rm -r $TEMPDIR
