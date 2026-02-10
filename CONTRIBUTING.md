# Contributing

## Pull Requests

1. Fork the repository and create a feature branch
2. Submit a PR with a clear description of the changes

## Creating a Release

1. Update the version in `version.go`
2. Commit: `git commit -m "chore: prepare release vX.Y.Z"`
3. Tag and push:

   ```bash
   git tag -a vX.Y.Z -m "Release vX.Y.Z"
   git push origin vX.Y.Z
   ```

This project follows [Semantic Versioning](https://semver.org/).
