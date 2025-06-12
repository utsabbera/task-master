# task-master

AI powered task manager

## Development

### Development Environment

This project uses [Devbox](https://www.jetpack.io/devbox/) to manage development dependencies and [direnv](https://direnv.net/) to automatically load the environment.

#### Prerequisites

- Install Devbox by following the instructions at [jetpack.io/devbox/docs/installing_devbox/](https://www.jetpack.io/devbox/docs/installing_devbox/)
- Install direnv by following the instructions at [direnv.net](https://direnv.net/docs/installation.html)

#### Getting Started

1. Clone this repository
2. Set up direnv in your shell:

   - For zsh, add `eval "$(direnv hook zsh)"` to your `~/.zshrc`
   - For bash, add `eval "$(direnv hook bash)"` to your `~/.bashrc` or `~/.bash_profile`

3. Allow direnv to use the project's `.envrc` file:

```bash
cd /path/to/task-master
direnv allow
```

4. The environment will be automatically loaded whenever you enter the project directory

5. Use the Makefile commands:

- Format code: `make fmt`
- Lint code: `make lint`
- Run tests: `make test`
- Build binary: `make build`
- Show all available commands: `make help`

All necessary tools (Go, golangci-lint) will be automatically installed in the isolated environment.

#### Manual Activation (alternative to direnv)

If you prefer not to use direnv, you can manually start the devbox shell:

```bash
devbox shell
```

### Git Hooks with Lefthook

This project uses [Lefthook](https://github.com/evilmartians/lefthook) to manage Git hooks, automatically running tests and linting before commits and pushes.

#### Setup

After cloning the repository, install the Git hooks:

```bash
make hooks
```

#### Available Hooks

- **pre-commit**: Runs linting and tests before each commit
- **pre-push**: Runs comprehensive tests with race detection before pushing

#### Run a specific hook:

```bash
lefthook run pre-commit
lefthook run pre-push
```

#### Skip Hooks

In cases where you need to skip hooks (not recommended for normal workflow):

```bash
git commit -m "message" --no-verify
git push --no-verify
```
