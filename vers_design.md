# VERS Implementation Design

This document outlines the selected design for implementing VERS (Version Range Specification) parsing support in go-univers, as requested in [GitHub Issue #11](https://github.com/alowayed/go-univers/issues/11).

## Background

VERS provides a universal notation for expressing version ranges across different package ecosystems using the syntax: `vers:<ecosystem>/<constraints>`

**Key Design Principle**: Maintain type safety and ecosystem isolation - "spend as little time as possible in a generic vers world and convert to strongly-typed ecosystem ranges as quickly as possible."

**Implementation Scope**: Parse VERS strings into strongly-typed ecosystem-specific version ranges. No conversion back to VERS format is required.

## Selected Approach: Ecosystem-Specific Parsing Methods

**User API:**
```go
// Each ecosystem implements its own VERS parsing
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
// In each ecosystem package (e.g., pkg/ecosystem/maven/):
func ParseVers(versString string) (*VersionRange, error)

// Existing method remains unchanged for backward compatibility:
func (e *Ecosystem) NewVersionRange(s string) (*VersionRange, error)

// Shared VERS parsing logic in pkg/vers/:
func ParseVersConstraints(ecosystem, constraints string) ([]Constraint, error)
```

**Why This Approach:**
- ✅ Maintains type safety (no cross-ecosystem mixing possible)
- ✅ Follows existing patterns in the codebase
- ✅ Each ecosystem can optimize parsing for its specific syntax
- ✅ Clear, intuitive API that matches existing `NewVersion()` pattern
- ✅ Aligns with project philosophy of ecosystem isolation
- ✅ Simple scope - just parsing, no conversion complexity

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

### Phase 3: Testing
3. **Comprehensive testing** across all ecosystems with VERS parsing

### Phase 4: Documentation
4. **Update documentation** with VERS parsing examples and usage patterns

## CLI Integration Examples

```bash
# Parse VERS range and check version containment using existing commands
univers maven contains "[1.0.0,2.0.0]" "1.5.0"  # Current functionality

# Future: Could extend existing commands to accept VERS syntax
univers maven contains "vers:maven/>=1.0.0|<=2.0.0" "1.5.0"  # → true
```

## Success Criteria

- [ ] All supported ecosystems can parse VERS notation into strongly-typed version ranges
- [ ] Type safety maintained - no generic types enabling cross-ecosystem mixing  
- [ ] Comprehensive test coverage with edge cases for VERS parsing
- [ ] Documentation with VERS parsing examples and usage patterns
- [ ] Backward compatibility with existing APIs

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