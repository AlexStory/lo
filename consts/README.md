# Consts

The `consts` package provides globally shared singleton objects for the `lo` interpreter.

This helps avoid unnecessary allocations for common immutable values such as:

- `TrueBool`: The boolean `true` value.
- `FalseBool`: The boolean `false` value.
- `Nil`: The `nil` value.
