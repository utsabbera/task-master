---
applyTo: "_"
---

# Go Testing Guidelines

## Best Practices

- Keep tests small and focused
- Use descriptive test names (start with should)
- Avoid testing private functions directly
- Aim for high test coverage, especially for critical paths
- Use subtests for better organization
- Don't modify global state without restoring it
- Avoid race conditions (run tests with `go test -race` to catch race conditions)
