---
name: golang-project-layout
description: Go project layout rules based on golang-standards/project-layout
---

# Go Project Layout Rules

When organizing or creating new files in this Go repository, you MUST follow the standard Go project layout conventions.

## Key Directories

- `/cmd`: Main applications for this project. Keep `main.go` small and invoke code from `/internal` or `/pkg`. Each app gets a subdirectory (e.g., `/cmd/myapp`).
- `/internal`: Private application and library code. This code cannot be imported by external applications. Structure further into `/internal/app` and `/internal/pkg` if needed.
- `/pkg`: Library code that is safe to use by external applications.
- `/api`: OpenAPI/Swagger specs, JSON schemas, protocol definition files.
- `/configs`: Configuration file templates or default configs.
- `/scripts`: Scripts to perform build, install, analysis, etc.
- `/build`: Packaging and Continuous Integration configurations.
- `/deployments`: IaaS, PaaS, system, and container orchestration deployment configurations (docker-compose, kubernetes, etc.).
- `/test`: Additional external test apps and test data.
- `/docs`: Design and user documents.
- `/tools`: Supporting tools for this project.
- `/examples`: Examples for your applications and/or public libraries.

## Strict Prohibitions

- **DO NOT** use a `src/` directory for project source files. This is a Java pattern and strongly discouraged in Go.
- **DO NOT** dump all code into the `/cmd` directory. Code that is reusable goes to `/pkg`, and private code goes into `/internal`.
- **DO NOT** put reusable library code directly inside `/cmd`.
