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

## Implementation Status (Current)

### ✅ Completed (Maven Ecosystem)

**Core Infrastructure:**
- `pkg/vers/vers.go` - Complete VERS package with single public API
- `vers.Contains(versRange, version string) (bool, error)` - Working implementation
- Internal VERS constraint parsing and Maven range conversion
- Comprehensive test coverage in `pkg/vers/vers_test.go`

**Supported Features:**
- ✅ VERS string parsing: `vers:<ecosystem>/<constraints>`
- ✅ All VERS operators: `>=`, `<=`, `>`, `<`, `=`, `!=` 
- ✅ Interval logic: Multiple constraint pairs create separate intervals
- ✅ Maven ecosystem integration using existing `NewVersion()` and `NewVersionRange()`
- ✅ Complex VERS examples: `"vers:maven/>=1.0.0-beta1|<=1.7.5|>=7.0.0-M1|<=7.0.7"`

**Working Examples:**
```go
// Simple range
vers.Contains("vers:maven/>=1.0.0|<=2.0.0", "1.5.0") // → true

// Complex intervals  
vers.Contains("vers:maven/>=1.0.0-beta1|<=1.7.5|>=7.0.0-M1|<=7.0.7", "1.1.0") // → true

// Exact match
vers.Contains("vers:maven/=1.5.0", "1.5.0") // → true

// Single bounds
vers.Contains("vers:maven/>=1.0.0", "2.0.0") // → true
vers.Contains("vers:maven/<=2.0.0", "1.0.0") // → true
```

### ✅ Implementation Details (Complete)

**VERS Constraint Processing Algorithm:**
1. Parse VERS string: `"vers:maven/>=1.0.0|<=1.7.5|>=7.0.0|<=7.0.7"`
2. Split into constraints: `[">=1.0.0", "<=1.7.5", ">=7.0.0", "<=7.0.7"]`
3. Group into intervals: `[1.0.0,1.7.5]` and `[7.0.0,7.0.7]`
4. Convert to Maven ranges: `"[1.0.0,1.7.5]"` and `"[7.0.0,7.0.7]"`
5. Test version against ANY interval (VERS OR logic)

**Key Functions in `pkg/vers/vers.go`:**
- `Contains(versRange, version string) (bool, error)` - Public API
- `containsMaven(constraints, version string) (bool, error)` - Maven implementation
- `parseVersConstraints(constraints string) ([]versConstraint, error)` - Parse pipe-separated constraints
- `groupConstraintsIntoIntervals(constraints []versConstraint) ([]interval, error)` - VERS interval logic
- `convertIntervalToMavenRange(interval interval) string` - Convert to Maven bracket notation

**Test Coverage (Complete):**
- ✅ Simple ranges (single interval)
- ✅ Complex ranges (multiple intervals)
- ✅ Exact matches with `=`
- ✅ Single bounds (upper/lower only)
- ✅ Error cases (invalid VERS, invalid versions, unsupported ecosystems)
- ✅ Complex VERS example from GitHub issue #11

### 🎯 Next Steps for Future Engineers

**Maven Ecosystem: ✅ COMPLETE**
The Maven ecosystem implementation is fully functional and tested. All basic VERS operations work correctly.

**Immediate Expansion Tasks:**
1. **Add NPM ecosystem support** - Follow the Maven pattern in `containsNpm()` function
2. **Add PyPI ecosystem support** - Follow the Maven pattern in `containsPyPI()` function  
3. **Add Go module ecosystem support** - Follow the Maven pattern in `containsGomod()` function
4. **Improve VERS algorithm** - Current implementation handles most cases but could be enhanced for complex edge cases
5. **Add `!=` exclusion support** - Currently skipped, needs proper implementation
6. **Add CLI integration** - Extend existing CLI to accept VERS syntax

**Adding New Ecosystems (Pattern):**
```go
// In Contains() function switch statement:
case "npm":
    return containsNpm(constraints, version)

// Implement ecosystem-specific function:
func containsNpm(constraints, version string) (bool, error) {
    // Parse version using NPM
    e := &npm.Ecosystem{}
    v, err := e.NewVersion(version)
    if err != nil {
        return false, fmt.Errorf("invalid NPM version '%s': %w", version, err)
    }

    // Convert VERS constraints to NPM ranges
    ranges, err := convertVersToNpmRanges(constraints)
    if err != nil {
        return false, fmt.Errorf("failed to convert VERS constraints: %w", err)
    }

    // Test against ANY range (VERS interval logic)
    for _, r := range ranges {
        if r.Contains(v) {
            return true, nil
        }
    }
    return false, nil
}

// Implement VERS-to-ecosystem conversion:
func convertVersToNpmRanges(constraints string) ([]*npm.VersionRange, error) {
    // Parse constraints using shared logic
    versConstraints, err := parseVersConstraints(constraints)
    if err != nil {
        return nil, err
    }

    // Group into intervals
    intervals, err := groupConstraintsIntoIntervals(versConstraints)
    if err != nil {
        return nil, err
    }

    // Convert to NPM syntax and create ranges
    var ranges []*npm.VersionRange
    e := &npm.Ecosystem{}
    
    for _, interval := range intervals {
        rangeStr := convertIntervalToNpmRange(interval) // Implement this
        if rangeStr == "" {
            continue
        }
        r, err := e.NewVersionRange(rangeStr)
        if err != nil {
            return nil, fmt.Errorf("failed to create NPM range '%s': %w", rangeStr, err)
        }
        ranges = append(ranges, r)
    }

    return ranges, nil
}
```

**VERS Algorithm Improvements Needed:**
- Better handling of constraint ordering and precedence
- Proper `!=` exclusion logic that applies across intervals
- Edge case testing with complex constraint combinations
- Performance optimization for large constraint sets

**Architecture Notes:**
- Keep all VERS logic centralized in `pkg/vers/`
- Never add VERS dependencies to ecosystem packages
- Always convert to native ecosystem ranges, never expose VERS internals
- Test only the public `Contains()` API, not internal functions

**CLI Integration Ideas:**
```bash
# Future CLI support could be:
univers vers contains "vers:maven/>=1.0.0|<=2.0.0" "1.5.0"  # → true
univers vers validate "vers:maven/>=1.0.0|<=2.0.0"          # → valid

# Or extend existing commands:
univers maven contains "vers:maven/>=1.0.0|<=2.0.0" "1.5.0"  # → true
```

### 📚 References and Context

**VERS Specification:** https://github.com/package-url/vers-spec/blob/main/VERSION-RANGE-SPEC.rst

**Key VERS Concepts:**
- Pipe `|` is a constraint separator, not logical OR
- Intervals are implicit based on constraint pairing
- Version satisfies range if it falls within ANY interval
- State machine algorithm for complex interval detection

**Implementation Philosophy:**
- Single function API for maximum simplicity
- Zero ecosystem pollution to maintain clean architecture
- Type safety through internal conversion to strongly-typed objects
- Reuse existing ecosystem APIs rather than reimplementing logic

### 📋 Implementation Summary

**Status: Maven Ecosystem COMPLETE ✅**

The VERS (Version Range Specification) implementation is now fully functional for the Maven ecosystem with the following achievements:

**✅ Completed Features:**
- Single stateless API: `vers.Contains(versRange, version string) (bool, error)`
- Complete VERS string parsing: `"vers:<ecosystem>/<constraints>"`
- VERS operators supported: `>=`, `<=`, `>`, `<`, `=` (Note: `!=` parsed but limited Maven range support)
- VERS interval logic correctly implemented
- Maven ecosystem integration using existing APIs
- Comprehensive test coverage including complex examples
- Zero ecosystem pollution - Maven package unchanged
- Full type safety maintained through internal conversions

**✅ Verified Examples Working:**
```go
vers.Contains("vers:maven/>=1.0.0|<=2.0.0", "1.5.0") // → true
vers.Contains("vers:maven/>=1.0.0-beta1|<=1.7.5|>=7.0.0-M1|<=7.0.7", "1.1.0") // → true
vers.Contains("vers:maven/=1.5.0", "1.5.0") // → true
```

## Current Implementation Status (December 2024)

### ✅ **Major Achievements Completed**

**Core Architecture:**
- ✅ Ecosystem-specific VERS handling with clean separation (`pkg/vers/maven.go`)
- ✅ Proper exclude (!=) operator support using dual-range logic
- ✅ Multiple interval support for complex VERS ranges
- ✅ Comprehensive test coverage with edge cases
- ✅ All basic VERS functionality working

**Successfully Working Features:**
- ✅ Simple ranges: `vers:maven/>=1.0.0|<=2.0.0`
- ✅ Complex multi-interval ranges: `vers:maven/>=1.0.0|<=1.7.5|>=7.0.0|<=7.0.7`
- ✅ Exclude combinations: `vers:maven/>=1.0.0|<=3.0.0|!=2.0.0`
- ✅ Multiple excludes: `vers:maven/!=1.0.0|!=2.0.0`
- ✅ Exact matches: `vers:maven/=1.5.0`
- ✅ Single bounds: `vers:maven/>=1.0.0` or `vers:maven/<=2.0.0`
- ✅ Star wildcards: `vers:maven/*`
- ✅ Basic validation edge cases (empty ecosystem, missing separators, etc.)

### 🔧 **Remaining Work (3 Edge Cases)**

**Issue: Constraint Merging vs Pairing Logic**

The current algorithm creates multiple intervals when there are mixed constraint types, but some edge cases expect constraint merging instead of pairing.

**Failing Test Cases:**
1. **`maven_unordered_constraints_-_outside_range`**
   - Input: `vers:maven/>=2.0.0|>=1.0.0|<=3.0.0` with version `0.5.0`
   - Expected: `true` (should match the upper bound `<=3.0.0`)
   - Actual: `false`
   - Issue: Algorithm pairs constraints instead of creating separate intervals

2. **`maven_multiple_lower_bounds_-_should_take_most_restrictive`**
   - Input: `vers:maven/>=1.0.0|>=2.0.0|<=3.0.0` with version `1.5.0`
   - Expected: `false` (should use most restrictive lower bound `>=2.0.0`)
   - Actual: `true` (creates multiple intervals including `[1.0.0,3.0.0]`)

3. **`maven_multiple_upper_bounds_-_should_take_most_restrictive`**
   - Input: `vers:maven/>=1.0.0|<=3.0.0|<=2.0.0` with version `2.5.0`
   - Expected: `false` (should use most restrictive upper bound `<=2.0.0`)
   - Actual: `true` (creates multiple intervals including `[1.0.0,3.0.0]`)

**Root Cause:**
The `groupConstraintsIntoIntervals()` function in `/pkg/vers/vers.go:272-378` needs smarter logic to determine when to:
- **Merge constraints** (take most restrictive) vs
- **Pair constraints** (create multiple intervals)

**Recommended Solution:**
Implement a more sophisticated constraint analysis that follows the VERS specification more closely:

1. **Detect constraint patterns** to determine merging vs pairing strategy
2. **For redundant constraints** (multiple bounds of same type): use most restrictive
3. **For alternating patterns** (>=a|<=b|>=c|<=d): create multiple intervals via pairing
4. **Validate against VERS spec** for proper constraint processing rules

**Files to Modify:**
- `/pkg/vers/vers.go` - `groupConstraintsIntoIntervals()` function
- `/pkg/vers/vers_test.go` - Validate edge case fixes

### 🚀 **Future Enhancements**

**VERS Validation Rules:**
- Implement full VERS specification validation rules
- Add URI encoding/decoding support
- Add constraint normalization per specification

**Additional Ecosystems:**
- NPM: Follow Maven pattern in `pkg/vers/npm.go`
- PyPI: Follow Maven pattern in `pkg/vers/pypi.go`
- Go modules: Follow Maven pattern in `pkg/vers/gomod.go`

### 🧪 **Quality Assurance Status**

**Current Test Results:**
```bash
go test ./pkg/vers/ -v
# 34/37 tests passing
# 3 edge cases failing (constraint merging logic)
```

**Code Quality:**
- ✅ Code formatted (go fmt ./...)
- ✅ No vet warnings (go vet ./...)
- ✅ Clean architecture with ecosystem separation

### 📋 **Handoff Notes**

The VERS implementation is ~95% complete with only 3 edge cases remaining. The architecture is solid and the exclude functionality was significantly improved. The failing tests provide clear reproduction steps for the remaining constraint merging logic issues.

**Key Architecture Points:**
- Excludes are handled separately from intervals (not as Maven ranges)
- Each ecosystem has its own conversion logic (`intervalToMavenRanges()`)
- Constraints are normalized and sorted before interval creation
- Multiple intervals are supported for complex VERS ranges

---

This implementation provides a solid foundation for VERS support while maintaining go-univers' core principles of type safety and ecosystem isolation. The architecture is proven and ready for completing the final edge cases and extending to other ecosystems.