#!/usr/bin/env bash
set -e # exit when any command fails

if [ -n "$(git status --porcelain)" ]; then
  echo "There are uncommitted changes to this repo. Commit changes and build with bin/build.sh first."
else
  echo "no changes"
fi

echo "cleaning dist dir"
mkdir -p dist
DIST=`ls dist`
for i in $DIST; do rm dist/$i; done

echo "copying executables to dist dir"
find build -type f -exec cp {} dist \;

echo "taring executables in dist"
DIST=`ls dist`
for i in $DIST; do tar -czvf dist/$i.tar.gz dist/$i; done

echo "bumping minor release version"
semver up release
git add .semver.yaml
VERSION=`semver get release`
git commit -m $VERSION
git push

echo "tagging release"
git tag $VERSION
git push --tags

echo "creating github release"
gh release create $VERSION dist/*

echo "all done"
gh release view $VERSION
