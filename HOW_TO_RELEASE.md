## Release
Use [goreleaser](https://goreleaser.com/)

set Github token `GITHUB_TOKEN`

local build

```bash
gorelease build --snapshot --rm-dist
```

Create tag

```bash
git tag -a vX.X.X -m "comment"
git push origin vX.X.X
```

Release

```bash
goreleaser release --rm-dist
```
