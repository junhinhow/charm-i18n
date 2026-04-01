# Contributing to charm-i18n

Thanks for your interest in contributing!

## How to contribute

1. Fork the repository
2. Create a feature branch (`git checkout -b feat/my-feature`)
3. Write tests for your changes
4. Ensure all tests pass (`go test ./...`)
5. Commit your changes using [conventional commits](https://www.conventionalcommits.org/)
6. Push to your branch and open a Pull Request

## Development

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detector
go test -race ./...

# Run linter
go vet ./...
```

## Code style

- Follow standard Go conventions (`gofmt`, `go vet`)
- Write descriptive variable and function names
- Add comments in Portuguese for internal documentation
- All public APIs must be thread-safe
- Keep functions small and focused

## Adding features

- New locale detection sources should follow the existing priority chain
- All catalog methods must be protected with the mutex
- JSON format must remain backward-compatible

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
