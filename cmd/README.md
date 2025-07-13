# CLI Research and Design

## Research: Existing Version/Range CLIs

### 1. npm semver (Node.js)
**Usage**: `semver [options] <version> [<version> [...]]`

**Key Features**:
- Sort versions: `semver 1.2.3 2.0.0 0.1.0`
- Range filtering: `semver -r ">=1.2.3" 1.2.4 2.0.0`
- Version increment: `semver -i minor 1.2.3`
- Coercion: `semver -c "1.2"` → `1.2.0`

**Flags**:
- `-r, --range <range>`: Filter versions by range
- `-i, --increment [level]`: Increment version (major/minor/patch/prerelease)
- `-c, --coerce`: Convert to valid semver
- `-p, --include-prerelease`: Include prereleases in range matching
- `-l, --loose`: Loose parsing

### 2. python-semver (pysemver)
**Subcommands**:
- `pysemver compare v1 v2`: Compare two versions (-1/0/1)
- `pysemver bump major 1.2.3`: Increment version parts

### 3. PHLAK SemVer-CLI (PHP)
**Command Structure**: `semver <action>:<target>`

**Key Patterns**:
- Initialization: `semver initialize [version]`
- Setters: `semver set:major 2`, `semver set:version 1.2.3`
- Getters: `semver get:version`, `semver get:major`
- Incrementers: `semver increment:minor`

### 4. Go semver-cli Tools
**Common Patterns**:
- Project initialization: `semver init`
- Version retrieval: `semver get`
- Version bumping: `semver up major|minor|patch|alpha|beta|rc`

## Key Insights

### 1. Command Patterns
- **Single command + flags** (npm semver): Simple, Unix-like
- **Subcommands** (pysemver, PHLAK): More discoverable, extensible
- **Action:target syntax** (PHLAK): Explicit but verbose

### 2. Core Operations
All tools support these fundamental operations:
- **Parse/validate** versions
- **Compare** versions
- **Increment** versions (major/minor/patch)
- **Sort** multiple versions
- **Range matching** (check if version satisfies range)

### 3. Extensibility Patterns
- **Ecosystem flags**: `--ecosystem npm|pypi` for multi-ecosystem support
- **Output formats**: JSON, plain text, structured data
- **Batch operations**: Handle multiple versions/ranges at once

### 4. User Experience
- **Single-purpose commands** work well for scripting
- **Interactive modes** help with discovery
- **Validation feedback** prevents user errors
- **Consistent exit codes** enable shell scripting

### 5. Multi-Ecosystem Considerations
Unlike existing tools that focus on single ecosystems, go-univers needs:
- Ecosystem selection mechanism
- Consistent behavior across ecosystems
- Clear error messages when mixing ecosystems
- Ecosystem-specific validation

## Proposed CLI API Design

### Core Philosophy
- **Ecosystem-first**: Users must specify ecosystem (npm, pypi) explicitly to leverage go-univers' type safety
- **Composable**: Small, focused commands that work well together and in scripts
- **Extensible**: Easy to add new ecosystems and operations without breaking existing workflows

### Command Structure
```bash
univers <ecosystem> <command> [args]
```

### Core Commands

#### Compare versions
```bash
univers npm compare "1.2.3" "1.2.4"    # Output: -1
univers pypi compare "1.0.0" "1.0.0"    # Output: 0
univers npm compare "2.0.0" "1.9.9"     # Output: 1
```

#### Sort versions
```bash
univers npm sort "2.0.0" "1.2.3" "1.10.0"
# Output: 1.2.3, 1.10.0, 2.0.0

univers pypi sort "1.0.0a1" "1.0.0" "0.9.0"
# Output: 0.9.0, 1.0.0a1, 1.0.0
```

#### Check if version satisfies range
```bash
univers npm satisfies "1.2.5" "^1.2.0"  # Exit 0 (true)
univers npm satisfies "2.0.0" "^1.2.0"  # Exit 1 (false)
univers pypi satisfies "1.2.5" "~=1.2.0" # Exit 0 (true)
```

### Utility Commands
```bash
univers help              # Show general help
univers help npm          # Show npm-specific help
univers version           # Show tool version
```

### Key Design Decisions

1. **Ecosystem as first argument**: Forces users to be explicit about which version system they're using, preventing cross-ecosystem confusion

2. **Minimal command set**: Start with core operations that cover the most common use cases

3. **Commands over flags**: `help` and `version` are commands for consistency and discoverability

4. **Shell-friendly**: Proper exit codes for scripting (satisfies command returns 0/1)

5. **Simple output**: Plain text output that's easy to parse and human-readable

This design leverages go-univers' type safety while providing a practical CLI that can grow with future ecosystem additions.

## Implementation Research: Go CLI Best Practices

### Standard Library CLI Patterns

#### 1. Project Structure (No External Dependencies)
**Recommended Pattern**:
- **main package**: Single line `main()` function that calls CLI and uses `os.Exit()`
- **app package**: Core CLI logic, flag parsing, and command execution
- **Avoid**: Global variables, `init()` functions, package-level state

**Key Benefits**:
- Go compiles to single, standalone binary with no external dependencies
- Standard library (`flag`, `os`, `fmt`, `errors`) provides all necessary CLI tools
- Easy distribution and deployment across different systems

#### 2. Testable Design Patterns
**Core Principle**: "Can I at least pass in a dummy `os.Args`?"

**Recommended Architecture**:
```go
type appEnv struct {
    // Configuration and state
}

func (env *appEnv) CLI(args []string) int {
    // Parse flags from args (not os.Args directly)
    // Execute command logic
    // Return exit code
}

func main() {
    os.Exit((&appEnv{}).CLI(os.Args[1:]))
}
```

**Testing Benefits**:
- CLI logic accepts `args []string` parameter instead of using `os.Args` directly
- Methods can be tested with mock arguments
- Exit codes can be verified without actually exiting

#### 3. Flag Package Usage
**Standard Library Approach**:
- Use `flag.NewFlagSet()` for subcommands instead of global `flag` package
- Create separate flag sets for each command
- Parse flags within command methods, not globally

**Pattern**:
```go
func (env *appEnv) parseFlags(args []string) error {
    fs := flag.NewFlagSet("command", flag.ContinueOnError)
    // Define flags on fs, not global flag
    return fs.Parse(args)
}
```

### Architectural Layers

#### Three-Layer CLI Design
1. **Input Layer**: Command-line argument parsing and validation
2. **Logic Layer**: Core command execution (calls into pkg/ libraries)
3. **Output Layer**: Formatting and displaying results

#### Command Structure Pattern
**Flexible Command Design**:
```go
type Command struct {
    name    string
    flagSet *flag.FlagSet
    run     func(args []string) error
}

func (c *Command) Execute(args []string) error {
    if err := c.flagSet.Parse(args); err != nil {
        return err
    }
    return c.run(c.flagSet.Args())
}
```

### Error Handling Best Practices
- **Return errors instead of panicking**: CLI should handle errors gracefully
- **Meaningful exit codes**: 0 for success, non-zero for different error types
- **User-friendly messages**: Clear error descriptions for end users
- **Consistent error format**: Standardized error output across commands

### Key Design Principles
1. **Zero external dependencies**: Use only Go standard library
2. **Single binary distribution**: No runtime dependencies required
3. **Testable by design**: Accept arguments as parameters, return error codes
4. **Extensible architecture**: Easy to add new commands and ecosystems
5. **Shell-friendly**: Proper exit codes and parseable output formats