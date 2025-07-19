# Maven Versioning Research

## Maven Version Format

Maven versions follow a complex format that goes beyond simple semantic versioning:

### Components
- **MajorVersion.MinorVersion.IncrementalVersion-BuildNumber-Qualifier**
- Examples: `1.2.3`, `1.2.3-4`, `1.2.3-SNAPSHOT`, `1.2.3-alpha-1`

### Qualifier Ordering (from lowest to highest precedence)
1. `alpha` or `a`
2. `beta` or `b` 
3. `milestone` or `m`
4. `rc` or `cr` (release candidate)
5. `snapshot`
6. (empty string) or `ga` or `final` or `release`
7. `sp` (service pack)

### Key Comparison Rules

1. **Versions with qualifiers are always older than release versions**
   - `1.2-beta` < `1.2`
   - `1.2-SNAPSHOT` < `1.2`

2. **Numeric components compared numerically**
   - `1.10` > `1.9`

3. **String components compared case-insensitively by qualifier precedence**
   - `1.2-alpha` < `1.2-beta` < `1.2-rc` < `1.2`

4. **Separators**: Both `.` and `-` are separators, with transitions between chars and digits also acting as separators
   - `1.0alpha1` → `[1, 0, alpha, 1]`

5. **Normalization**:
   - Trailing zeros, empty strings, "final", "ga" are trimmed
   - `1.0.0` = `1.0` = `1`
   - `1.0-ga` = `1.0`

## Version Range Syntax

Maven supports version ranges for dependency management:

### Range Types
- `[1.0]` - Exact version 1.0
- `(,1.0]` - Any version ≤ 1.0 (exclusive lower bound)
- `[1.0,)` - Any version ≥ 1.0 (inclusive lower bound)
- `[1.0,2.0]` - Any version between 1.0 and 2.0 (inclusive)
- `(1.0,2.0)` - Any version between 1.0 and 2.0 (exclusive)
- `[1.0,2.0)` - Any version ≥ 1.0 and < 2.0

### Special Keywords
- `LATEST` - Latest released or snapshot version
- `RELEASE` - Latest non-snapshot release
- `SNAPSHOT` - Development/unreleased version

## Implementation References

### Maven ComparableVersion Algorithm
Maven's official implementation in `org.apache.maven.artifact.versioning.ComparableVersion`:
- Handles unlimited version components
- Supports mixing of separators (`.` and `-`)
- Case-insensitive qualifier comparison
- Complex normalization rules

### Existing Go Implementations
1. **google/osv-scalibr**: Sophisticated token-based parsing with Maven-specific rules
2. **google/deps.dev**: Custom element parsing with qualifier normalization
3. **aboutcode-org/univers**: Python reference implementation with comprehensive range support

## Special Cases and Edge Cases

### SNAPSHOT Versions
- Represents unreleased/development versions
- Always considered older than release versions
- Used in CI/CD for dynamic updates

### Normalization Edge Cases
- `1.0.0-ga` = `1.0.0` = `1.0` = `1`
- `1.0-final` = `1.0`
- `1.0.0.0` = `1.0`

### Qualifier Shortcuts
- `a` → `alpha`
- `b` → `beta`
- `m` → `milestone`
- `cr` → `rc`

## Testing Strategy

Key test cases to implement:
1. Basic semantic versioning: `1.0.0` vs `2.0.0`
2. Qualifier precedence: `1.0-alpha` < `1.0-beta` < `1.0-rc` < `1.0`
3. Normalization: `1.0` = `1.0.0` = `1.0-ga`
4. Complex versions: `1.0.0-alpha-1` vs `1.0.0-alpha-2`
5. SNAPSHOT handling: `1.0-SNAPSHOT` < `1.0`
6. Range parsing and containment tests
7. Edge cases with unusual separators and formats

## Architecture Notes

Following the go-univers pattern:
- Separate `Version` and `VersionRange` types
- Implement `univers.Version` and `univers.VersionRange` interfaces
- Private parsing logic, public API only
- Comprehensive table-driven tests
- CLI integration following existing patterns