# VERS: Version Range Specification

This document explains how the [VERS (Version Range Specification)](https://github.com/package-url/vers-spec/blob/main/VERSION-RANGE-SPEC.rst) represents and treats version ranges.

## Overview

VERS is a universal notation for expressing version ranges across different package ecosystems. It provides a compact, minimalist syntax that can represent complex version constraints.

## Basic Syntax

```
vers:<ecosystem>/<constraints>
```

Where:
- `<ecosystem>` is the package type (e.g., `maven`, `npm`, `pypi`)
- `<constraints>` are version constraints separated by pipe (`|`) characters

Example: `vers:maven/>=1.0.0-beta1|<=1.7.5|>=7.0.0-M1|<=7.0.7|>=7.1.0|<=7.1.2|>=8.0.0-M1|<=8.0.1`

## The Pipe Operator (`|`)

**Key Understanding**: The pipe (`|`) is NOT simply an OR operator. It's a constraint separator that defines intervals in the version timeline.

The pipe character has no special semantic meaning beyond being a separator. The actual logic for determining version containment follows a state machine algorithm.

## Version Containment Algorithm

The VERS specification defines a precise algorithm for checking if a version is contained within a range:

### 1. Initial Checks

- If the range is `*`, any version is IN the range
- Check for exact matches with `=` constraints
- Check for explicit exclusions with `!=` constraints

### 2. Constraint Processing

The algorithm processes constraints as intervals:

1. **Sort constraints** by version in ascending order
2. **Ensure uniqueness** - no duplicate versions in constraints
3. **Iterate through constraint pairs** to find valid intervals

### 3. State Machine Logic

For each pair of constraints, the algorithm checks:

- **First constraint `<` or `<=`**: If the tested version is less than this constraint, it's IN the range
- **Last constraint `>` or `>=`**: If the tested version is greater than this constraint, it's IN the range
- **Interval checking**: If current constraint is `>` or `>=` and next is `<` or `<=`, and the tested version falls between them, it's IN the range

### 4. Final Determination

If no valid interval is found after complete iteration, the version is NOT in the range.

## Example Analysis

Let's analyze: `vers:maven/>=1.0.0-beta1|<=1.7.5|>=7.0.0-M1|<=7.0.7|>=7.1.0|<=7.1.2|>=8.0.0-M1|<=8.0.1`

This creates the following intervals:
1. `[1.0.0-beta1, 1.7.5]` - versions from 1.0.0-beta1 up to and including 1.7.5
2. `[7.0.0-M1, 7.0.7]` - versions from 7.0.0-M1 up to and including 7.0.7
3. `[7.1.0, 7.1.2]` - versions from 7.1.0 up to and including 7.1.2
4. `[8.0.0-M1, 8.0.1]` - versions from 8.0.0-M1 up to and including 8.0.1

A version satisfies this range if it falls within ANY of these intervals.

## Common Patterns

### Simple Range
```
vers:npm/>=1.0.0|<2.0.0
```
Matches versions from 1.0.0 (inclusive) to 2.0.0 (exclusive).

### Multiple Intervals
```
vers:pypi/>=1.0.0|<=1.3.0|>=2.0.0|<=2.5.0
```
Matches versions in [1.0.0, 1.3.0] OR [2.0.0, 2.5.0].

### Exclusions
```
vers:npm/>=1.0.0|<2.0.0|!=1.5.0
```
Matches versions from 1.0.0 to 2.0.0, excluding 1.5.0.

### Exact Version
```
vers:maven/=1.2.3
```
Matches only version 1.2.3.

## Implementation Considerations

When implementing VERS support:

1. **Parse constraints** into structured objects with operator and version
2. **Sort constraints** by version to ensure proper interval detection
3. **Implement the state machine** as described in the specification
4. **Handle edge cases** like first/last constraints and exclusions
5. **Respect ecosystem-specific** version comparison rules

## Key Differences from Other Range Notations

Unlike other version range syntaxes:
- Pipes are separators, not logical operators
- Intervals are implicit based on constraint ordering
- The algorithm is universal across ecosystems
- Complex ranges can be expressed compactly

## Ecosystem Support

VERS is designed to work with any versioning scheme, including:
- npm (Semantic Versioning)
- PyPI (PEP 440)
- Maven (Maven versioning)
- RubyGems (Ruby versioning)
- And many others

Each ecosystem uses its own version comparison rules while following the universal VERS interval logic.

## References

- [VERS Specification](https://github.com/package-url/vers-spec/blob/main/VERSION-RANGE-SPEC.rst)
- [Reference Implementation](https://github.com/aboutcode-org/univers)
- [Package URL Specification](https://github.com/package-url/purl-spec)