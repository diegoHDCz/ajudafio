# AjudaFio

AjudaFio is a health management platform designed to help users track and manage their health-related information and activities.

## Architecture Documentation

Architecture decision records and diagrams are located in the [`docs/`](./docs) folder:

- [`docs/architecture/`](./docs/architecture) — C4 context diagrams and architectural overviews
- [`docs/adr/`](./docs/adr) — Architecture Decision Records (ADRs)

## Tech Stack

- **Language:** Go
- **Identity Provider:** Keycloak
- **Database:** PostgreSQL
- **Architecture:** Hexagonal (Ports & Adapters) with vertical slicing
