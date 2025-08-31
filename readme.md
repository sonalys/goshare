# GoShare

[![Tests](https://github.com/sonalys/goshare/actions/workflows/test.yml/badge.svg)](https://github.com/sonalys/goshare/actions/workflows/test.yml)
[![Linter](https://github.com/sonalys/goshare/actions/workflows/lint.yml/badge.svg)](https://github.com/sonalys/goshare/actions/workflows/lint.yml)
[![codecov](https://codecov.io/gh/sonalys/goshare/graph/badge.svg?token=YE75TAB5BF)](https://codecov.io/gh/sonalys/goshare)

A showcase of my current preferences for the development of DDD RESTful APIs in Golang.  
GoShare is a ledger/expenses sharing api.

## Project's Architecture

It follows a clean hexagonal and domain driven design.

```
/goshare
├── .config/                # Configuration files and deployment variables
├── .github/                # Github CI/CD automation files
├── .tools/                 # Tools and their dependency tree
├── .vscode/                # Scripts for quick automations around debugging / launching
├── cmd/
│   ├── migration           # Migration deployment entrypoint
│   └── server              # HTTP Server entrypoint
│
├── internal/
│   ├── domain              # Entities, aggregates
│   ├── application         # Controllers, usecases
│   ├── ports               # Interfaces
│   ├── mocks               # Code generated mocks for dependency injection
│   └── infrastructure/
│       ├── repositories    # Database adapters
│       ├── postgres        # Postgres connection management
│       └── http/           # HTTP adapter (handlers, router, server, open api specification)
│           ├── middlewares # Recovery, Logging, Security
│           ├── router      # Handlers for specific endpoints
│           └── server      # Code generated server from open api specification
└── pkg/                    # Shared and general purpose packages and definitions
```