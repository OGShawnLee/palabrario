# 020 — Spec checklist

Run by the router before accepting the ticket. All items must pass.

- [x] All 8 `.routing/services/*_service.md` files audited against `pkg/models/`, `core-service/internal/service/`, and `core-service/internal/api/`.
- [x] Every Go struct in each service spec matches its corresponding `pkg/models/` definition:
  - All field names and types match.
  - All `json:"..."` tags match `snake_case`.
  - Embedded structs have their definitions present in the same file or clearly cross-linked.
- [x] Every SOA interface in each service spec matches the Go interface in `core-service/internal/service/` verbatim (method names, parameter types, return signatures).
- [x] Endpoint tables in each service spec match routes registered in `core-service/internal/api/routes.go` (and satellite service routes).
- [x] No application code, Go backend code, or `.context/` files modified (documentation edits only).
- [x] Comprehensive audit report generated in `.tickets/020-audit-service-specs-with-backend/changelog.md` detailing findings per service.
