# Ruby Gems Versioning Research

## Overview

Ruby Gems follow semantic versioning principles with some unique features and constraints specific to the Ruby ecosystem.

## Version Format

### Basic Structure
Ruby Gem versions typically follow the format: `MAJOR.MINOR.PATCH[.BUILD][-PRERELEASE][+BUILD]`

Examples:
- `1.2.3` - Basic semantic version
- `1.2.3.4` - With build number
- `1.2.3-alpha` - With prerelease identifier
- `1.2.3.pre` - Common prerelease format
- `2.0.0.rc1` - Release candidate format

### Key Characteristics
1. **Numeric segments**: Major, minor, patch are always integers
2. **Prerelease identifiers**: Can contain letters (a-z) and are case-insensitive
3. **Build numbers**: Additional numeric segment after patch
4. **Flexible format**: Can have varying number of segments

### Prerelease Handling
- Prerelease versions contain alphabetic characters
- Common prerelease identifiers: `alpha`, `beta`, `rc`, `pre`
- Prerelease versions sort lower than release versions
- Example ordering: `1.0 > 1.0.b1 > 1.0.a.2 > 0.9`

## Version Constraints

### Constraint Operators

1. **Exact match**: `= 1.2.3` or just `1.2.3`
2. **Inequality**: `!= 1.2.3`
3. **Comparison**: `>= 1.2.3`, `> 1.2.3`, `<= 1.2.3`, `< 1.2.3`
4. **Pessimistic (Twiddle-wakka)**: `~> 1.2.3`

### Pessimistic Constraint (~>)

The `~>` operator is the most important Ruby-specific constraint:

- `~> 1.2.3` means `>= 1.2.3` AND `< 1.3.0`
- `~> 1.2` means `>= 1.2.0` AND `< 2.0.0`
- `~> 0.1.0` means `>= 0.1.0` AND `< 0.2.0`

Rules:
- Only the **last** digit specified can increment
- Prevents major/minor version updates that might break compatibility
- Guards against potential bugs in future releases

### Multiple Constraints
Can combine multiple constraints with commas:
- `~> 2.2, >= 2.2.1` - Pessimistic with minimum version
- `~> 2, != 2.2.1` - Pessimistic with exclusion

## Comparison Algorithm

### Segment-by-Segment Comparison
1. Split version into segments by dots
2. Compare numerically when both segments are numbers
3. Compare lexically when segments contain letters
4. Numeric segments have lower precedence than alphabetic ones
5. Missing segments are treated as having lower precedence

### Special Rules
1. **Trailing zeros**: Remove trailing zeros from version segments
2. **Canonicalization**: Add dots between numeric and non-numeric segments
3. **Case handling**: Prerelease identifiers are case-insensitive
4. **Prerelease precedence**: `1.0.0 > 1.0.0-alpha`

## Implementation References

### Python Implementation (aboutcode-org/univers)
- Uses segment-based parsing with canonical form
- Implements sophisticated comparison with prerelease handling
- Supports version manipulation methods (bump, release)

### Go Implementation (google/deps.dev)
- Uses `gemExtension` struct for version components
- Complex element-based comparison with string/numeric handling
- Canonicalizes version strings for consistent comparison

### Go Implementation (google/osv-scalibr)
- Focuses on canonicalization and segment grouping
- Removes trailing zeros and handles mixed numeric/string segments
- Implements robust component comparison algorithm

## Key Implementation Considerations

1. **Type Safety**: Separate types for gem versions to prevent cross-ecosystem mixing
2. **Parsing Robustness**: Handle various version formats and edge cases
3. **Prerelease Support**: Proper ordering of alpha, beta, rc, pre versions
4. **Constraint Implementation**: Full support for pessimistic operator
5. **Performance**: Efficient parsing and comparison algorithms
6. **Error Handling**: Clear error messages for invalid version strings

## Testing Requirements

1. **Basic versions**: Standard semantic versions
2. **Prerelease versions**: Alpha, beta, rc, pre formats
3. **Complex versions**: Multiple segments with mixed types
4. **Edge cases**: Trailing zeros, case sensitivity, malformed input
5. **Constraint testing**: All operators including pessimistic
6. **Comparison matrix**: Comprehensive version ordering tests

## Examples for Testing

### Valid Versions
```
1.0.0
1.2.3.4
1.0.0-alpha
2.0.0.rc1
1.2.3.pre
1.0.0.beta.1
```

### Version Constraints
```
~> 1.2.3    # >= 1.2.3, < 1.3.0
~> 1.2      # >= 1.2.0, < 2.0.0
>= 1.0.0    # Greater than or equal
!= 1.5.0    # Not equal
```

### Comparison Examples
```
1.0.0 > 1.0.0.rc1
1.0.0.rc1 > 1.0.0.beta
1.0.0.beta > 1.0.0.alpha
2.0.0 > 1.9.9
```