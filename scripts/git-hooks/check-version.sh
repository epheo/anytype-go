#!/bin/bash
# pre-commit hook script to check if version needs updating
# Place this file in .git/hooks/pre-commit and make it executable

VERSION_FILE="pkg/anytype/version.go"
BRANCH=$(git branch --show-current)

# Only run this check when on the master branch
if [ "$BRANCH" != "master" ]; then
    exit 0
fi

# Get the current version from the version.go file
CURRENT_VERSION=$(grep 'Version = "' "$VERSION_FILE" | sed 's/.*Version = "//g' | sed 's/".*//g')

# Get the last release tag from git
LAST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")

# Check if the commit message indicates a version bump should occur
commit_msg=$(cat "$1" 2>/dev/null || git log -1 --pretty=%B)

if echo "$commit_msg" | grep -qE "fix:|bug:|patch:"; then
    echo "⚠️ This commit appears to be a bug fix."
    echo "🔄 Consider bumping the patch version before merging to master."
    echo "📝 Current version: $CURRENT_VERSION"
elif echo "$commit_msg" | grep -qE "feat:|feature:|minor:"; then
    echo "⚠️ This commit adds new features."
    echo "🔄 Consider bumping the minor version before merging to master."
    echo "📝 Current version: $CURRENT_VERSION"
elif echo "$commit_msg" | grep -qE "BREAKING|breaking:"; then
    echo "⚠️ This commit contains BREAKING CHANGES."
    echo "🔄 Consider bumping the major version before merging to master."
    echo "📝 Current version: $CURRENT_VERSION"
fi

# This hook is advisory only and doesn't block the commit
exit 0
