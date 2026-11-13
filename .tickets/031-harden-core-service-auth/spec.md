# 031 — Spec: Harden Core Service Auth, Short-Lived JWTs, Validation, and English Errors

- [x] `pkg/models/auth.go` includes `RefreshToken` on `AuthResponse`, `RefreshTokenRequest`, and `SignOutRequest`.
- [x] `core-service/internal/config/config.go` defaults `AccessTokenTTL` to `15 * time.Minute`, `RefreshTokenTTL` to `168 * time.Hour` (or `JWT_EXPIRATION_HOURS`), and configures `AllowedOrigins`.
- [x] `core-service/internal/service/errors.go` maps `22001`/`23514`/`23503` to `400`, distinguishes `email` vs `user_name` `23505` conflicts, and emits English messages + machine-readable `code`.
- [x] `core-service/internal/service/identity_service.go` implements `RefreshToken` (atomic rotation), `SignOut` (token revocation), symmetric `SignUp` validation, `72`-byte password guard, `GetMe` `401`/`403` on missing/suspended accounts, and `ChangePassword` token-version bump + refresh revocation.
- [x] `core-service/internal/service/governance_service.go` bumps `token_version` on role changes and revokes refresh tokens + bumps `token_version` on account suspension (`is_active = false`).
- [x] `core-service/internal/api/middleware/auth_middleware.go` verifies 15m JWTs on all protected routes, enforces DB `is_active`, `token_version`, and live `role` on all `POST`/`PUT`/`DELETE` routes plus `/api/v1/auth/*` and `/api/v1/governance/*` (Option B), and returns `401` in `OptionalAuth` when a malformed/expired `Authorization` header is present.
- [x] `core-service/internal/api/router.go` exposes `POST /api/v1/auth/refresh` and `POST /api/v1/auth/sign-out` and sets `AllowCredentials: false` with configured `AllowedOrigins`.
- [x] `.routing/DOMAIN.md` (Section 3, 4, 4.2), `.routing/services/identity_service.md` (Part A & Part B), `governance_service.md`, `dictionary_service.md`, `curation_service.md`, `studio_service.md`, and `mastery_service.md` reflect the English `{ "code", "error" }` contract, Option B middleware, token invalidation rules, and refresh/sign-out endpoints.
- [x] With `postgres-test (:5433)` running (`task test:up`), `go test -race -count=1 ./core-service/...` passes with exit code `0`.
- [x] `git status --short` confirms no files outside the declared `Scope` were created or modified.
- [ ] No unauthorized `git commit`, `git push`, or `git merge` was run. (If committed after explicit user authorization: English, Conventional Commits, detailed body.)
