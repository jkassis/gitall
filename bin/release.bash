#!/usr/bin/env bash
set -e # exit when any command fails

# make sure the user runs `gh auth login` first
echo "You must have a current authorized session with github cli before running this."
echo 'Have you run `gh auth login` lately?'
read -p "Y/N [Y]: " YESNO
YESNO=${YESNO:-Y}
if [ "$YESNO" != "Y" ]; then
  echo "No problem. Let's do it"
  gh auth login
fi

# check for local changes
echo "git: checking for changes"
if [ -n "$(git status --porcelain)" ]; then
  echo "git: there are uncommitted changes to this repo. Commit changes and build with bin/build.sh first."
  exit 1
else
  echo "git: no changes"
fi

# check that branches are in sync
echo "git: checking that local branch is in sync with origin"
BRANCH=`git rev-parse --abbrev-ref HEAD`
if [ x"$(git rev-parse $BRANCH)" != x"$(git rev-parse origin/$BRANCH)" ]
then
  echo "$BRANCH is not in sync with origin/$BRANCH. You need to rebase or push first."
  exit 1
fi

# clean up the dist directory
echo "cleaning dist dir"
mkdir -p dist
DIST=`ls dist`
for i in $DIST; do rm dist/$i; done

# copy the execs
echo "copying executables to dist dir"
find build -type f -exec cp {} dist \;

# tar the execs
echo "taring executables in dist"
DIST=`ls dist`
for i in $DIST; do
  tar -czvf dist/$i.tar.gz dist/$i
  rm dist/$i
done

# bump the minor release version
echo "bumping minor release version"
semver up release
git add .semver.yaml
VERSION=`semver get release`
git commit -m $VERSION
git push

# tag the release
echo "tagging release"
git tag $VERSION
git push --tags

# create the github release
echo "creating github release"
gh release create $VERSION dist/*

# finish
echo "all done"
gh release view $VERSION
