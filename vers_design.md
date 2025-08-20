# VERS Implementation Design

This document outlines the selected design for implementing VERS (Version Range Specification) parsing support in go-univers, as requested in [GitHub Issue #11](https://github.com/alowayed/go-univers/issues/11).

## Background

VERS provides a universal notation for expressing version ranges across different package ecosystems using the syntax: `vers:<ecosystem>/<constraints>`

**Key Design Principle**: Maintain type safety and ecosystem isolation - "spend as little time as possible in a generic vers world and convert to strongly-typed ecosystem ranges as quickly as possible."

**Implementation Scope**: Parse VERS strings into strongly-typed ecosystem-specific version ranges. No conversion back to VERS format is required.

## Selected Approach: Hybrid Stateless + Ecosystem-Specific Methods

**User API:**
```go
// Convenient stateless API (primary usage) - automatic ecosystem detection
contains, err := vers.Contains("vers:maven/>=1.0.0|<=1.7.5", "1.1.0")
compare, err := vers.Compare("vers:npm/^1.2.0", "1.3.0", "1.4.0")

// Type-safe ecosystem API (power users) - explicit ecosystem selection
range, err := maven.ParseVers("vers:maven/>=1.0.0-beta1|<=1.7.5|>=7.0.0-M1|<=7.0.7")
version, _ := maven.NewVersion("1.1.0")
contains := range.Contains(version) // Returns bool

// Works with all ecosystems
npmRange, err := npm.ParseVers("vers:npm/>=1.2.0|<2.0.0")
pypiRange, err := pypi.ParseVers("vers:pypi/~=1.2.3|!=1.2.5")
goRange, err := gomod.ParseVers("vers:go/>=v1.2.0|<v2.0.0")
alpineRange, err := alpine.ParseVers("vers:alpine/>=1.2.0")
```

**Implementation Structure:**
```go
// Stateless API in pkg/vers/ (primary usage):
func Contains(versRange, version string) (bool, error)
func Compare(versRange, version1, version2 string) (int, error)

// In each ecosystem package (e.g., pkg/ecosystem/maven/):
func ParseVers(versString string) (*VersionRange, error)

// Existing method remains unchanged for backward compatibility:
func (e *Ecosystem) NewVersionRange(s string) (*VersionRange, error)

// Shared VERS parsing logic in pkg/vers/:
func ParseVersConstraints(ecosystem, constraints string) ([]Constraint, error)

// Internal implementation of stateless API:
func Contains(versRange, version string) (bool, error) {
    ecosystem, _ := parseVersString(versRange)
    switch ecosystem {
    case "maven":
        range, err := maven.ParseVers(versRange)  // Reuse existing ParseVers
        if err != nil { return false, err }
        
        ver, err := (&maven.Ecosystem{}).NewVersion(version)
        if err != nil { return false, err }
        
        return range.Contains(ver), nil
    // ... similar for all ecosystems
    }
}
```

**Why This Hybrid Approach:**
- ✅ **Best user experience** - Stateless API for convenience, ecosystem API for power users
- ✅ **Matches VERS philosophy** - Universal interface that feels truly universal
- ✅ **Maintains type safety** - Ecosystem-specific methods prevent cross-ecosystem mixing
- ✅ **Minimal API surface** - Only adds `ParseVers()` to ecosystems and stateless functions
- ✅ **Reuses existing methods** - `NewVersion()`, `Contains()`, `Compare()` 
- ✅ **Mirrors CLI architecture** - Automatic ecosystem detection like the CLI does
- ✅ **Backward compatible** - Existing `NewVersionRange()` remains unchanged
- ✅ **Simple scope** - Just parsing, no conversion complexity

## Other Approaches Considered

Several alternative approaches were evaluated but rejected:

1. **Generic Factory with Type Parameters** - Complex generics and less intuitive API
2. **Registry-Based Factory** - Breaks type safety with runtime type assertions  
3. **Interface-Based with Ecosystem Injection** - More verbose API with extra objects
4. **Hybrid Package-Level Functions** - Larger API surface with potential confusion

- Maintains the core principle of type safety and ecosystem isolation
- Minimizes time spent in generic VERS world  
- Prevents cross-ecosystem version mixing at compile time
- Matches the pattern of `ecosystem.NewVersion()` and `ecosystem.NewVersionRange()`
- Provides intuitive API: `maven.ParseVers()` clearly indicates what you get back
- Allows each ecosystem to optimize for its specific constraints
- Enables shared common logic in `pkg/vers/` for constraint parsing

## Implementation Plan

### Phase 1: Core Infrastructure
1. **Create `pkg/vers/` package** with shared VERS parsing utilities:
   - `ParseVersConstraints(ecosystem, constraints string) ([]Constraint, error)`
   - VERS state machine implementation for version containment checking

### Phase 2: Ecosystem Integration  
2. **Add VERS parsing method to each ecosystem**:
   - `ParseVers(versString string) (*VersionRange, error)`
   - Integration with existing constraint types

### Phase 3: Stateless API
3. **Create stateless API in `pkg/vers/`**:
   - `Contains(versRange, version string) (bool, error)`
   - `Compare(versRange, version1, version2 string) (int, error)`
   - Ecosystem routing logic with switch statement

### Phase 4: Testing
4. **Comprehensive testing** across all ecosystems:
   - Test ecosystem-specific `ParseVers()` methods
   - Test stateless API functions
   - Test ecosystem routing logic

### Phase 5: Documentation
5. **Update documentation** with both API patterns and usage examples

## CLI Integration Examples

```bash
# Current functionality (unchanged)
univers maven contains "[1.0.0,2.0.0]" "1.5.0"  # → true

# Future: Could add new vers commands using stateless API
univers vers contains "vers:maven/>=1.0.0|<=2.0.0" "1.5.0"  # → true
univers vers compare "vers:npm/^1.2.0" "1.3.0" "1.4.0"      # → -1

# Or extend existing commands to accept VERS syntax
univers maven contains "vers:maven/>=1.0.0|<=2.0.0" "1.5.0"  # → true
```

## Success Criteria

- [ ] All supported ecosystems can parse VERS notation into strongly-typed version ranges via `ParseVers()`
- [ ] Stateless API functions (`vers.Contains()`, `vers.Compare()`) work with all ecosystems  
- [ ] Type safety maintained - no generic types enabling cross-ecosystem mixing
- [ ] Ecosystem routing logic correctly dispatches to appropriate parsers
- [ ] Comprehensive test coverage with edge cases for both API patterns
- [ ] Documentation with examples for both stateless and ecosystem-specific usage
- [ ] Backward compatibility with existing APIs (`NewVersionRange()` unchanged)

## Constructor Naming Convention Analysis

### Go Standard Library Patterns

Based on analysis of Go standard library packages, the following constructor naming patterns are established:

#### 1. **Parse*** for String Input
```go
// time package
time.Parse(layout, value string)
time.ParseInLocation(layout, value string, loc *Location)
time.ParseDuration(s string)

// url package  
url.Parse(rawURL string)
url.ParseRequestURI(rawURL string)
url.ParseQuery(query string)

// strconv package
strconv.ParseInt(s string, base int, bitSize int)
strconv.ParseFloat(s string, bitSize int)
```

#### 2. **New*** for Object Creation
```go
// encoding/json
json.NewEncoder(w io.Writer)
json.NewDecoder(r io.Reader)

// time package
time.NewTimer(d Duration)
time.NewTicker(d Duration)
```

#### 3. **Unmarshal*** for Deserialization
```go
// encoding/json
json.Unmarshal(data []byte, v any)
```

#### 4. **Specific Constructors** (Direct Names)
```go
// time package
time.Now()
time.Unix(sec, nsec int64)
time.Date(year, month, day, hour, min, sec, nsec, loc)
```

### Decision: ParseVers vs Alternatives

**Selected Approach:**
```go
// Current (native format parsing) - KEEP AS IS
func (e *Ecosystem) NewVersionRange(s string) (*VersionRange, error)

// New (VERS format parsing) - ADD THIS  
func ParseVers(versString string) (*VersionRange, error)
```

**Why `ParseVers()` is best:**

1. **Follows Go conventions**: `Parse*` for string-to-type conversion is the standard pattern
2. **Clear distinction**: `NewVersionRange()` for native format, `ParseVers()` for VERS format
3. **No breaking changes**: Existing code continues to work
4. **Intuitive**: Developers immediately understand `ParseVers()` works with VERS strings

**Alternative considered but rejected:**
```go
// This would require breaking changes and is less clear
func NewVersionRangeFromNative(s string) (*VersionRange, error)
func NewVersionRangeFromVers(s string) (*VersionRange, error)
```

The `Parse*` prefix is the idiomatic Go way to handle string parsing, as seen throughout the standard library.

---

This design positions go-univers as a leading implementation of the emerging VERS standard while maintaining its core architectural principles of type safety and ecosystem isolation. The focus on parsing-only keeps the implementation simple and achieves the core goal of supporting VERS notation.