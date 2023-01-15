# Changelog

Changelog is used to keep track of version changes. The versioning scheme used is [SemVer](https://semver.org/). First integer is used for breaking change, second integer is used for major patches, and third integer is used for minor bug fixes.

## Version 1.0.5 (15/01/2023)

- Fix a bug where the username is compared again, thus resulting in a failed comparison because the username is matched with the password in the default map function.
- Upgrade CI version in GitHub Actions.

## Version 1.0.4 (09/05/2022)

- Improve CI by introducing three new checks: `go mod verify`, `go vet ./...`, and `go build -v ./...`.
- Additional check for race conditions in unit-tests: `go test -race -v -cover ./... ./...`.
- Create new CI job specific for paralellized linting purposes with [`golangci-lint`](https://github.com/golangci/golangci-lint).
- Revise `example` project to check for uncommon errors (courtesy of `golangci-lint`).

## Version 1.0.3 (26/04/2022)

- Fix typos in `basic.go` and `README.md`.
- Ensures consistency in code comments in `basic.go`.
- Remove unused code in `.github/workflows/tag.yml`.

## Version 1.0.2 (24/03/2022)

- Create GitHub action to automate releases in `workflow_dispatch`.
- Elaborate `Users` attribute, which is a one-to-one mapping of usernames and passwords.

## Version 1.0.1 (24/03/2022)

- Fix possible timing attack in the default `Authenticator` function.

## Version 1.0.0 (24/03/2022)

- Official initial release of the library.
