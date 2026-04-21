# mytets

## Dependency Rationale

This project introduces `github.com/spf13/cobra` for CLI command routing.
The dependency is used to provide stable subcommand and flag behavior,
consistent `--help` output, and clean extensibility for future commands.

The standard library flag parser was not selected for this feature because it
does not provide subcommand ergonomics and help output consistency at the same
level without additional custom command-routing scaffolding.