#!/bin/bash
# This script bumps the version number of Distru. The variable "Version"
# should be declared as a constant in the file distru.go.
#
# To use: bump.sh <level> [--tag]
#
# <level> is the depth of the bump, determined by number of periods
# to the left of the number. For example, bumping v1 to v2 is a level
# 0 bump. Bumping v1.1 to v1.2 is a level 1 bump. By default, no more
# than two levels are allowed, just due to versioning clutter.
#
# If --tag is specified, then this script will pass the following
# argument to `git commit`. This could be nothing, which will open
# an editor to add the change normally, or another argument like
# --amend, to as to append the change to the most recent commit. This
# script will then create a git tag with the version number, prepended
# with "v".
#
# If --tag is enabled, and <level> is 0 (bumping base version) then
# this script will create a GPG-signed tag by default.

# DEFAULTS
SIGNBASE=TRUE
MAXLEVEL=2
FILE=$(dirname $0)/distru.go

# SCRIPT
LEVEL=$1
GIT=FALSE

if [[ -z $LEVEL ]]
then
	echo No depth given.
	exit 1
fi

if [[ $LEVEL > $MAXLEVEL ]]
then
	echo No more than two sub versions, please.
	exit 1
fi

if [[ $2 == "--tag" ]]
then
	GIT=TRUE
fi

VER=$(grep -P "Version =" $FILE \
	| grep -o -P "[0-9]+(\.[0-9]+)*")

NEWVER=

div=
i=0
for SUB in $(echo $VER | tr "\." "\n")
do
	if [[ $i == $LEVEL ]]
	then
		NEWVER=$NEWVER$div$(expr $SUB + 1)
		i=$(expr $i + 1)
		break
	fi
	
	NEWVER=$NEWVER$div$SUB
	div=.
	i=$(expr $i + 1)
done

if [[ $(expr $i - 1) != $LEVEL ]]
then
	NEWVER=$NEWVER$div\1
fi

sed -i "s/Version = \"$VER\"/Version = \"$NEWVER\"/" $FILE


if [[ $GIT == TRUE ]]
then
	# We're just going to add the file with
	# the version constant.
	git add $FILE

	# Then amend it to the most recent commit.
	git commit $3

	if [[ $? != 0 ]]
	then
		echo Commit and tagging aborted.
		exit 1
	fi

	if [[ $LEVEL == 0 ]] && [[ $SIGNBASE == TRUE ]]
	then
		SIGNING="-s"
	fi

	# Now we make the tag.
	git tag $SIGNING -m v$NEWVER 

	echo Tagged as v$NEWVER
	exit 0
fi

echo Bumped to v$NEWVER
