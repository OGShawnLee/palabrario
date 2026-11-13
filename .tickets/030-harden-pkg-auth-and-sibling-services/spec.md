# 030 — Spec: Harden pkg/auth, pkg/response, Sibling Services, and Flaky Tests

- [x] `pkg/response/json.go` defines `ErrorBody{Code, Error}` and `ErrorWithCode(w, status, code, message)`.
- [x] `pkg/auth/auth.go` enforces bcrypt `<= 72` bytes, defaults access TTL to `15 * time.Minute`, adds `TokenVersion` to `Claims`, provides `GenerateRefreshToken`/`HashRefreshToken`, `ExtractBearerToken`, and `ParseCORSAllowedOrigins`, and uses English error messages.
- [x] Neither `realtime-service/internal/config/config.go` nor `image-service/internal/config/config.go` contains any hardcoded fallback `JWT_SECRET` (grep for `super_secret` and `super-secret` in `realtime-service/` and `image-service/` returns nothing).
- [x] Both `realtime-service/internal/config/config.go` and `image-service/internal/config/config.go` reject missing/short (`< 32`) `JWT_SECRET` and invalid `ACCESS_TOKEN_TTL_MINUTES` / `JWT_EXPIRATION_HOURS`.
- [x] `image-service/internal/api/middleware/auth_middleware.go` uses `auth.ExtractBearerToken`, returns English `ErrorWithCode`, and does not append `err.Error()` in `401` responses.
- [x] `realtime-service` and `image-service` routers set `AllowCredentials: false` and load `AllowedOrigins` from config; `realtime-service/internal/service/hub.go` validates non-empty WebSocket `Origin` headers against allowed origins.
- [x] `realtime-service/internal/service/service_test.go` uses a dedicated WS reader channel instead of repeating `ReadJSON` after timeout errors, removes the fragile `50ms`/`5ms` sleeps, and passes reliably.
- [x] `.routing/services/arena_service.md` (Sections `2.7.1`, `2.7.5`, Part B banner) and `.routing/services/media_service.md` (Part B banner, Section `2.8.2`) reflect `Authorization: Bearer` primary WS auth, origin checking, `TokenVersion` / `15m` TTL, and English `ErrorWithCode` (`401` without leaking internal errors).
- [x] `go test -race -count=3 ./pkg/... ./realtime-service/... ./image-service/...` exits with code `0`.
- [x] `git status --short` confirms no files outside the declared `Scope` were created or modified.
- [x] No unauthorized `git commit`, `git push`, or `git merge` was run. (If committed after explicit user authorization: English, Conventional Commits, detailed body.)
