

# Contributing

# Builds
gitall currently uses github workflows to run goreleaser for cross-platform builds.

We are also experimenting with xgo using https://github.com/crazy-max/ghaction-xgo

For more on xgo... https://github.com/karalabe/xgo

The workflows is triggered on a new tag push, so to make a release...

```
> git tag v7.7.7
> git push
> git push --tags
```

