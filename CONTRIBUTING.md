

# Contributing

# Builds

gitall currently uses github workflows to run xgo for cross-platform builds.

For more on all of this, see `.github/workflows/xgo.yml` and the following projects that perform checkout, build, and release.

* https://github.com/actions/checkout
* https://github.com/crazy-max/ghaction-xgo
* https://github.com/karalabe/xgo
* https://github.com/softprops/action-gh-release



For more on xgo... https://github.com/karalabe/xgo

The workflows triggers on a new tag push, so to make a release...

```
> git tag v7.7.7
> git push
> git push --tags
```

