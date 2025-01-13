# GO-TEMP

## Get started

```bash
# Init Project
$ go mod vendor
# Or
$ make init

# HTTP/1.1 Server
$ go run . http
# Or
$ make http

# GRPC Server
$ go run . grpc
# Or
$ make grpc
```

## Structure

- app
  - console
  - grpc
    - pb
  - modules
    - \[module]
      - dto
        - \[module].dto.go
      - ent
        - \[module].ent.go
      - int
        - \[module].inf.go
      - \[module].mod.go
      - \[module].ctl.go
      - \[module].mid.go
      - \[module].svc.go
- config
  - i18n
- internal
  - cmd
  - collections
  - database
  - encoding
  - logger
  - math
    - rand
  - otel
    - collector
  - ssl
- router
- storage

## Commitlint

### Commit Types

https://gist.github.com/parmentf/035de27d6ed1dce0b36a
Commonly used commit types from [Conventional Commit Types](https://github.com/commitizen/conventional-commit-types)

| Type     | Description                                                                      |
| :------- | :------------------------------------------------------------------------------- |
| feat     | A new feature                                                                    |
| fix      | A bug fix                                                                        |
| docs     | Documentation only changes                                                       |
| style    | Changes that do not affect the meaning of the code (white-space, formatting etc) |
| refactor | A code change that neither fixes a bug nor adds a feature                        |
| perf     | A code change that improves performance                                          |
| test     | Adding missing tests or correcting existing tests                                |
| build    | Changes that affect the build system or external dependencies                    |
| ci       | Changes to our CI configuration files and scripts                                |
| chore    | Other changes that don't modify src or test files                                |
| revert   | Reverts a previous commit                                                        |
