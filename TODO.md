## Fix:

- [x] Fix `make run`

## Review:

- [x] Review the usages of pointers in service and repository

## Features:

- [ ] Add stage mode to show preview before commititng the changes
- [ ] Add multi action support through prompt
- [ ] Integrate database to persist data
- [ ] Add API request validation mechanism
- [ ] Add more detail to README.md

- [x] Add feature to create new task in any status
- [x] Implement handlers - call service methods + Add tests

## Refactor:

- [x] Make the UpdatedAt and CreatedAt testable
- [x] Rename service and repository methods to make it consistent all accross

## Tech Enhancements:

- [ ] Add more test coverage
- [ ] Integrate better logging mechanism
- [ ] Use enummer for better enum management and parsing

- [ ] Setup remote pipeline
- [ ] Add hook to check coverage
- [ ] Configure golangci-lint for better linting

- [ ] Enhance swagger and bruno docs

- [x] Add api tests to test e2e
- [x] Restructure hooks - optimise duplicate runners for push and commit
