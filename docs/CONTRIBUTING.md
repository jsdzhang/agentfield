# Contributing Guide

Thank you for your interest in contributing to AgentField! This guide outlines how to propose changes, report issues, and participate in the community.

## Ground Rules

- Be kind and respectful. See `CODE_OF_CONDUCT.md`.
- Create issues before large changes to align on direction.
- Keep pull requests focused and small. Large refactors should be split.
- Follow existing coding style and conventions.
- Ensure tests pass locally before opening a pull request.

## Development Environment

1. Fork the repository and clone your fork.
2. Install dependencies:
   ```bash
   ./scripts/install.sh
   ```
3. Create a feature branch:
   ```bash
   git checkout -b feat/my-feature
   ```

## Commit Guidelines

- Use [Conventional Commits](https://www.conventionalcommits.org) when possible (`feat:`, `fix:`, `chore:`, etc.).
- Keep commit messages concise yet descriptive.
- Reference related issues with `Fixes #<id>` or `Refs #<id>` when applicable.

## Pull Requests

Before submitting:

1. Run `./scripts/test-all.sh`.
2. Run `make fmt tidy` to keep code formatted and dependencies tidy.
3. Update documentation and changelog entries where relevant.
4. Ensure CI workflows pass.

When opening a PR:

- Provide context in the description.
- Highlight user-facing changes and migration steps.
- Include screenshots for UI changes.
- Link to the issue being resolved.

## Issue Reporting

- Search existing issues to avoid duplicates.
- Use the provided issue templates (`bug`, `feature`, `question`).
- Include reproduction steps, logs, or stack traces when possible.

## Documentation

- Keep docs precise and actionable.
- Update `docs/DEVELOPMENT.md` for tooling or workflow changes.
- Update `docs/ARCHITECTURE.md` for structural changes.

## Release Workflow

Releases are automated via GitHub Actions. The workflow now uses
[git-cliff](https://github.com/orhun/git-cliff) to render changelog entries from the
commit history as part of the version bump step. Maintainers typically only need to:

1. Ensure commits follow the conventional prefixes (`feat:`, `fix:`, etc.) so they are
   categorized correctly.
2. Trigger the `Release` workflow with the desired SemVer component/channel.
3. Let the workflow bump versions, rebuild SDKs, update `CHANGELOG.md`, and publish
   artifacts automatically.

To preview the generated changelog locally install `git-cliff`
(`cargo install git-cliff` or grab a release binary) and run:

```bash
python scripts/update_changelog.py --version 0.1.8 --dry-run
```

If you install via Cargo make sure `~/.cargo/bin` is on your `PATH` so `git cliff`
invokes the plugin binary correctly.

## Questions?

Open a `question` issue or start a discussion in the repository. We’re excited to build with you!
