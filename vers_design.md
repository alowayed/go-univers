# VERS Implementation Design

This document outlines the selected design for implementing VERS (Version Range Specification) parsing support in go-univers, as requested in [GitHub Issue #11](https://github.com/alowayed/go-univers/issues/11).

## Background

VERS provides a universal notation for expressing version ranges across different package ecosystems using the syntax: `vers:<ecosystem>/<constraints>`

**Key Design Principle**: Maintain type safety and ecosystem isolation - "spend as little time as possible in a generic vers world and convert to strongly-typed ecosystem ranges as quickly as possible."

**Implementation Scope**: Parse VERS strings into strongly-typed ecosystem-specific version ranges. No conversion back to VERS format is required.

## Selected Approach: Stateless VERS-Only API

**User API:**
```go
// Single stateless API - automatic ecosystem detection and VERS processing
contains, err := vers.Contains("vers:maven/>=1.0.0|<=1.7.5|>=7.0.0|<=7.0.7", "1.1.0")

// That's it - no ecosystem-specific VERS methods needed
```

**Implementation Structure:**
```go
// Only public API in pkg/vers/:
func Contains(versRange, version string) (bool, error)

// No changes to any ecosystem packages - they remain completely clean

// Internal implementation converts VERS to native ecosystem ranges:
func Contains(versRange, version string) (bool, error) {
    ecosystem, constraints := parseVersString(versRange)
    
    switch ecosystem {
    case "maven":
        // Convert VERS constraints to native Maven ranges
        ranges := convertVersToMavenRanges(constraints)
        
        // Use existing Maven API to check containment
        version, err := (&maven.Ecosystem{}).NewVersion(version)
        if err != nil { return false, err }
        
        // VERS interval logic: version satisfies range if it's in ANY interval
        for _, range := range ranges {
            if range.Contains(version) {
                return true
            }
        }
        return false
    // ... similar for all ecosystems
    }
}

// VERS-specific logic converts constraints to native ecosystem ranges
func convertVersToMavenRanges(constraints string) []*maven.VersionRange {
    // Example: ">=1.0.0|<=1.7.5|>=7.0.0|<=7.0.7"
    // Becomes: Maven ranges "[1.0.0,1.7.5]" and "[7.0.0,7.0.7]" 
}
```

**Why This Approach:**
- ✅ **Simplest user experience** - Single function for all VERS operations
- ✅ **Matches VERS philosophy** - Truly universal interface across all ecosystems
- ✅ **No ecosystem pollution** - Ecosystems remain completely unchanged
- ✅ **Minimal API surface** - Only one function: `Contains()`
- ✅ **Reuses existing ecosystem APIs** - `NewVersion()`, `NewVersionRange()`, `Contains()`
- ✅ **Centralized VERS logic** - All VERS complexity in one place, no duplication
- ✅ **Backward compatible** - No changes to any existing ecosystem APIs
- ✅ **Type safety maintained** - VERS package handles conversion to strongly-typed objects

## Other Approaches Considered

Several alternative approaches were evaluated but rejected:

1. **Ecosystem-Specific ParseVers Methods** - Would add `ParseVers()` to each ecosystem, creating API bloat
2. **Hybrid Stateless + Ecosystem APIs** - Two different ways to do the same thing, confusing  
3. **Generic Factory with Type Parameters** - Complex generics and less intuitive API
4. **Registry-Based Factory** - Breaks type safety with runtime type assertions

The selected approach was chosen because it:
- Provides the simplest possible user experience (single function call)
- Keeps all ecosystems completely clean and unchanged
- Centralizes all VERS complexity in one place
- Maintains type safety through internal conversion to strongly-typed ecosystem objects

## Implementation Plan

### Phase 1: VERS Package Core
1. **Create `pkg/vers/` package** with stateless API:
   - `Contains(versRange, version string) (bool, error)`
   - Internal VERS parsing and constraint conversion logic

### Phase 2: Ecosystem Integration
2. **Add VERS-to-native conversion for each ecosystem**:
   - `convertVersToMavenRanges(constraints string) []*maven.VersionRange`
   - `convertVersToNpmRanges(constraints string) []*npm.VersionRange`
   - No changes to ecosystem packages themselves

### Phase 3: Testing
3. **Comprehensive testing** of VERS functionality:
   - Test stateless API functions with all supported ecosystems
   - Test VERS constraint parsing and conversion logic
   - Test complex VERS examples from the specification

### Phase 4: Documentation
4. **Update documentation** with VERS usage examples and patterns

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

- [ ] Stateless API function (`vers.Contains()`) works with all ecosystems
- [ ] VERS constraint parsing and conversion to native ecosystem ranges works correctly
- [ ] Type safety maintained - all conversions use strongly-typed ecosystem objects
- [ ] Complex VERS examples (intervals, exclusions) work as specified
- [ ] Comprehensive test coverage with edge cases for VERS functionality
- [ ] Documentation with VERS usage examples and patterns
- [ ] Zero changes to existing ecosystem APIs (complete backward compatibility)

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