# Contributing to anytype-go

Thank you for your interest in contributing to the anytype-go project! This document provides guidelines and instructions for contributing to this project.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Creating a New Release](#creating-a-new-release)
- [Submitting Pull Requests](#submitting-pull-requests)
- [Coding Standards](#coding-standards)

### Commit Message Format

We follow conventional commits for our commit messages:

- `feat:` - A new feature
- `fix:` - A bug fix
- `docs:` - Documentation only changes
- `style:` - Changes that do not affect the meaning of the code
- `refactor:` - A code change that neither fixes a bug nor adds a feature
- `perf:` - A code change that improves performance
- `test:` - Adding missing tests or correcting existing tests
- `chore:` - Changes to the build process or auxiliary tools

## Creating a New Release

This section documents the process for creating and publishing a new tagged release of the anytype-go SDK.

### Release Process

1. **Determine the new version number**

   Following [Semantic Versioning](https://semver.org/), determine the new version based on the changes since the last release:
   - MAJOR version for incompatible API changes
   - MINOR version for backward-compatible new functionality
   - PATCH version for backward-compatible bug fixes
   - For pre-release versions, append `-alpha.X`, `-beta.X`, etc.

2. **Check the current version and latest tag**

   ```bash
   # Check the current version in the code
   grep "Version =" pkg/anytype/version.go
   
   # List all existing tags
   git tag -l
   
   # Check changes since the last tag
   git log LAST_TAG..HEAD --oneline
   ```

3. **Update version in the codebase**

   Edit the version number in `pkg/anytype/version.go` to match the new version you've determined:

   ```go
   Version = "X.Y.Z-alpha.N"  // Replace with your new version
   ```

4. **Update the CHANGELOG.md**

   Create or update the `CHANGELOG.md` file with a new section for the release:
   - Use the format `## [X.Y.Z-alpha.N] - YYYY-MM-DD`
   - Categorize changes under `### Added`, `### Changed`, `### Fixed`, etc.
   - Include all notable changes since the last release

5. **Commit the version changes**

   ```bash
   git add pkg/anytype/version.go CHANGELOG.md
   git commit -m "chore: prepare release vX.Y.Z-alpha.N"
   ```

6. **Create and push the release tag**

   ```bash
   # Create an annotated tag
   git tag -a vX.Y.Z-alpha.N -m "Release vX.Y.Z-alpha.N"
   
   # Push the tag to the remote repository
   git push origin vX.Y.Z-alpha.N
   ```

### Example: Creating v0.2.0-alpha.2

For example, here's how we created the v0.2.0-alpha.2 release:

1. We determined that a minor version bump was appropriate due to new features
2. We checked the latest tag (v0.1.0-alpha.1) and reviewed changes since then
3. We updated the version in pkg/anytype/version.go from "0.2.0-alpha" to "0.2.0-alpha.2"
4. We created a CHANGELOG.md documenting all changes since v0.1.0-alpha.1
5. We committed these changes with "chore: prepare release v0.2.0-alpha.2"
6. We created and pushed a new tag: "v0.2.0-alpha.2"

## Submitting Pull Requests

When submitting a Pull Request, please:

1. Reference any related issues in the PR description
2. Provide a clear description of the changes
3. Update documentation as needed
4. Add or update tests to cover your changes
5. Ensure all tests pass

## Coding Standards

- Follow standard Go coding conventions
- Use meaningful variable and function names
- Write comments for complex logic
- Keep functions focused and small
- Add appropriate error handling
- Write unit tests for your code
