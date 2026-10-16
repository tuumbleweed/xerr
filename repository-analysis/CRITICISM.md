# Review findings

## Remaining Issues

- **Standard error interoperability (confirmed):** `type.go` defines no `Error() string` or `Unwrap() error` methods. The record cannot be returned as an `error` or used directly with `errors.Is`/`errors.As`. Add these methods if conventional error wrapping is intended, with defined nil behavior.
- **Call-site attribution (confirmed):** `NewError` uses a fixed caller skip of two. `NewErrorEC`, its single/multiline wrappers, and helpers in `init.go` therefore record a package-internal location. Centralize construction with a caller-depth parameter and verify each public entry point.
- **Cause/text synchronization (confirmed edge cases):** printing reads `ErrStr` directly, so replacing the exported `Err` can leave stale terminal output. `AfterFind` does not clear an existing `Err` when `ErrStr` is empty. JSON encoding clears `ErrStr` when `Err` is nil, including on records decoded from JSON. Define a single synchronization policy and add round-trip and mutation tests.
- **Truncation boundaries (confirmed):** `helpers.go` slices text by byte index, potentially splitting a UTF-8 character. Its notice increases the final length beyond 64 KiB, and formatting happens before the cap. Choose UTF-8-safe truncation; document output limits separately from resource limits. Avoid describing arbitrary Stringer execution or JSON formatting as universally safe.
- **Missing regression coverage (confirmed):** no test files or CI configuration are present. Prioritize context branches, serialization hooks, caller locations, and exit behavior in subprocess tests. The demo's “no exit” comment contradicts its positive stop code (`examples/demo.go`).

## Simplification Opportunities

- `print.go` constructs all three color maps on every print call. Build only the selected theme or share a small theme constructor to reduce repetitive palette definitions. Keep tintlog as the provider of color and logging behavior.
- `type.go` imports GORM solely for hook signatures, making an ORM dependency mandatory for terminal/JSON-only consumers. If such consumers matter, consider an optional adapter package. This is a tradeoff, not evidence that the present integration is unnecessary.
