---
applyTo: "*.go,go.mod,go.sum"
---

# Go Guidelines

## Code Style

- Keep functions short and focused on a single task.
- Use meaningful variable and function names that describe their purpose.

## Project Structure

- Follow the standard Go project layout:
  - `/cmd` - Main applications
  - `/pkg` - Library code that's ok to use by external applications
  - `/core` - Private application and library code
  - `/api` - API definitions (protobuf, OpenAPI specs)
  - `/scripts` - Scripts for build, setup, etc.

## Error Handling

- Return errors rather than using panic.
- Don't discard errors using `_` unless there's a good reason.
- Add context to errors: `fmt.Errorf("opening config file: %w", err)`
- Consider using custom error types for specific error conditions.

## Testing

- Write code that is testable: avoid hidden dependencies, use interfaces, and make functions pure when possible.
- Write tests before or alongside code
- Write tests for all exported functions.
- Use table-driven tests when applicable.
- Aim for high test coverage, especially for critical paths.

## Concurrency

- Use goroutines judiciously - they are lightweight but not free.
- Always use proper synchronization when sharing data between goroutines.
- Consider using sync.WaitGroup for managing groups of goroutines.
- Use channels for communication between goroutines, not for synchronization.
- Consider using contexts for cancellation and timeouts.

## Dependencies

- Minimize external dependencies when possible.
- Use Go modules for dependency management.
- Pin dependencies to specific versions in go.mod.
- Regularly update dependencies to get security fixes.

## Documentation

- Document all exported functions, types, and constants.
- Include examples in documentation when helpful.
- Follow the godoc conventions for comments.

## Performance

- Profile before optimizing.
- Consider performance implications of string concatenation, memory allocations, etc.
- Use sync.Pool for frequently allocated temporary objects.
- Be aware of escape analysis and stack vs. heap allocations.
