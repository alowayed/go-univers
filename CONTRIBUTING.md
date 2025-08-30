# Contribute to go-univers

## Environment setup

1. **Install Go 1.24+**
   ```bash
   # Check version
   go version
   ```

2. **Install golangci-lint** (for code quality checks)
   ```bash
   # Install latest version
   curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
   ```

3. **Install markdown-link-check** (for documentation link validation)
   ```bash
   # Install via npm (requires Node.js)
   npm install -g markdown-link-check
   ```

4. **Configure git**
   ```bash
   git config --global user.name "John Doe"
   git config --global user.email "john.doe@example.com"
   ```

## Contribution steps

1. **Fork and clone**
   ```bash
   git clone https://github.com/YOUR_USERNAME/go-univers.git
   cd go-univers
   ```

2. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make your changes and add tests**
   ```bash
   # Run tests frequently during development
   go test ./...
   ```

4. **Verify code quality**
   ```bash
   # Format code
   go fmt ./...
   
   # Run linters
   golangci-lint run
   
   # Check documentation links
   markdown-link-check . --config mlc_config.json
   
   # Ensure dependencies are clean
   go mod tidy
   ```

5. **[Sign](#sign-your-work) and commit your changes**
   ```bash
   git commit -s -m "feat: add new feature"
   ```

6. **Push and create pull request**
   ```bash
   git push origin feature/your-feature-name
   # Then create PR on GitHub
   ```

The CI pipeline will automatically test your changes on multiple platforms, verify code quality, and validate documentation links.

## Helpful Scripts

### Branch Cleanup
Clean up local branches that have been squash-merged:
```bash
./scripts/cleanup-merged-branches.sh
```

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

## Common contributions

## Adding a new ecosystem

See [CLAUDE.md](./CLAUDE.md) for detailed guidance on adding new ecosystems. The process involves:

1. Create package under `pkg/ecosystem/<ecosystem>/`
2. Implement `Version` and `VersionRange` types
3. Add comprehensive table-driven tests
4. Extend CLI support in `cmd/cli/commands.go`
5. Add the new ecosystem to the 'Supported Ecosystems' table in README.md

Refer to existing ecosystems like `cargo/` or `nuget/` for implementation patterns.

## Architecture

go-univers uses a **type-safe, ecosystem-isolated architecture** that prevents accidental cross-ecosystem version mixing. Each ecosystem (npm, pypi, go, etc.) has its own `Version` and `VersionRange` types, eliminating the common bug of accidentally comparing versions from different package managers.
