## Description

Please include a summary of the changes and which issue is fixed. Include relevant motivation and context.

Fixes #(issue)

## Type of Change

Please delete options that are not relevant:

- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update
- [ ] Refactoring (no functional changes)
- [ ] Performance improvement
- [ ] Translation (new language or translation updates)
- [ ] CI/CD or tooling changes

## Changes Made

Please provide a clear list of what you've changed:

- Change 1
- Change 2
- Change 3

## Testing

Please describe the tests you ran to verify your changes. Provide instructions so reviewers can reproduce:

**Test Configuration**:
- OS: [e.g., Ubuntu 22.04]
- Go version: [e.g., 1.23.0]
- Jellyfin version: [e.g., 10.8.13]

**Test Steps**:
1. Step 1
2. Step 2
3. Expected result

## Checklist

Please check all that apply:

- [ ] My code follows the Go code style guidelines (`go fmt`)
- [ ] I have run `go test ./...` and all tests pass
- [ ] I have run `golangci-lint run` (if available) and addressed any issues
- [ ] I have performed a self-review of my own code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes
- [ ] Any dependent changes have been merged and published

## Translation Updates

If your changes include user-facing text:

- [ ] I have updated `locales/active.en.toml` with new translation keys
- [ ] I have updated `locales/active.fa.toml` with Persian translations
- [ ] I have tested the changes in both English and Persian
- [ ] N/A - no user-facing text changes

## Breaking Changes

If this PR includes breaking changes, please describe:

**What breaks?**
- Description of what will no longer work

**Migration path:**
- How should users update their configuration/usage?

**Deprecation notice:**
- Should we deprecate old behavior first?

## Screenshots (if applicable)

If your changes include UI changes or new features, please add screenshots:

**Before:**
[Screenshot]

**After:**
[Screenshot]

## Additional Notes

Add any other context about the pull request here:

- Performance considerations
- Security implications
- Dependencies added or updated
- Future improvements planned
