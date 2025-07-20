# Contribute to go-univers

## Environment setup

1. Install `go`.
2. Configure git.
    ```
    $ git config --global user.name "John Doe"
    $ git config --global user.email "john.doe@example.com"
    ```

## Contribution steps

1. Fork the repository.
2. Create a feature branch.
3. Add your changes.
4. Add tests.
5. Ensure all tests pass: `go test ./...`.
6. [Sign](#sign-your-work) and commit your changes.
7. Submit a pull request.

## Sign your work

The `sign-off` is a line at the end of the explanation for the patch. Your
signature certifies that you wrote the patch or otherwise have the right to pass
it on as an open-source patch. The rules are pretty simple: if you can certify
the below (from [developercertificate.org](http://developercertificate.org/)):

```
Developer Certificate of Origin
Version 1.1

Copyright (C) 2004, 2006 The Linux Foundation and its contributors.
1 Letterman Drive
Suite D4700
San Francisco, CA, 94129

Everyone is permitted to copy and distribute verbatim copies of this
license document, but changing it is not allowed.

Developer's Certificate of Origin 1.1

By making a contribution to this project, I certify that:

(a) The contribution was created in whole or in part by me and I
    have the right to submit it under the open source license
    indicated in the file; or

(b) The contribution is based upon previous work that, to the best
    of my knowledge, is covered under an appropriate open source
    license and I have the right under that license to submit that
    work with modifications, whether created in whole or in part
    by me, under the same open source license (unless I am
    permitted to submit under a different license), as indicated
    in the file; or

(c) The contribution was provided directly to me by some other
    person who certified (a), (b) or (c) and I have not modified
    it.

(d) I understand and agree that this project and the contribution
    are public and that a record of the contribution (including all
    personal information I submit with it, including my sign-off) is
    maintained indefinitely and may be redistributed consistent with
    this project or the open source license(s) involved.
```

Then you just add a line to every git commit message:

```
Signed-off-by: Joe Smith <joe.smith@email.com>
```

Use your real name (sorry, no pseudonyms or anonymous contributions.)

If you set your `user.name` and `user.email` git configs, you can sign your
commit automatically with:

```
$ git commit -s -m "this is a commit message"
```

To double-check that the commit was signed-off, look at the log output:

```
$ git log -1
commit 4ec3560ff087e0f2526b2bd9d32ba50949ccf943 (HEAD -> issue-35, origin/main, main)
Author: John Doe <john.doe@example.com>
Date:   Mon Aug 1 11:22:33 2020 -0400

    this is a commit message

    Signed-off-by: John Doe <john.doe@example.com>
```

## Commit Message Format

This project uses [Conventional Commits](https://www.conventionalcommits.org/) for automated versioning and changelog generation.

### Format
```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Types
- **feat**: A new feature (triggers minor version bump)
- **fix**: A bug fix (triggers patch version bump)
- **docs**: Documentation only changes
- **style**: Changes that do not affect the meaning of the code
- **refactor**: A code change that neither fixes a bug nor adds a feature
- **perf**: A code change that improves performance
- **test**: Adding missing tests or correcting existing tests
- **build**: Changes that affect the build system or external dependencies
- **ci**: Changes to our CI configuration files and scripts
- **chore**: Other changes that don't modify src or test files

### Breaking Changes
Add `BREAKING CHANGE:` in the footer or `!` after the type to trigger a major version bump:
```
feat!: remove deprecated API endpoints

BREAKING CHANGE: The old v1 API has been removed. Use v2 endpoints instead.
```

### Examples
```
feat(gem): add Ruby Gems ecosystem support
fix(npm): handle caret ranges with zero versions correctly
docs: update README with gem ecosystem examples
```

## Release Process

Releases are automated via GitHub Actions:
1. Commit using conventional commit format
2. Push to main branch (via PR)
3. GitHub Actions will automatically:
   - Run tests
   - Determine version bump
   - Create GitHub release
   - Update CHANGELOG.md

## Common contributions

### Adding a new ecosystem

When adding new ecosystems:
1. Create a new package under `ecosystem/`.
2. Implement the core interfaces defined in `univers.go`.
3. Add table-driven unit tests.
4. Update the README.