# 032 — Spec: Harden Wails BFF Session Persistence, Auto-Refresh, Error Codes, and Tests

- [x] `app/app.go` defines `APIError{Status, Code, Message}`, stores `accessToken` in RAM, and persists `refreshToken` to `os.UserConfigDir()/palabrario/session.json` (`0600`).
- [x] `app/app.go` implements single-flight auto-refresh on `401` (`executeRequest`) and clears RAM + deletes `session.json` + emits `auth:expired` on `403 ACCOUNT_SUSPENDED` or failed refresh.
- [x] `app/app_media.go` uses `executeRequest` (zero duplicated Bearer/JSON-error code) and English error strings.
- [x] `app/app_arena.go` sends WebSocket auth via `Authorization: Bearer` header (`?token=` removed) and `FindMatchRoom` inspects `APIError.Status == http.StatusNotFound`.
- [x] All 14 duplicate alias methods (`GetMyProfile`, `CreateArenaMatch`, `GetMyMatchHistory`, `GetMatchReplayDetail`, `ConnectArena`, `LeaveArena`, `ListPendingEditorRequests`, `GenerateInvitationTickets`, `ListUsersForGovernance`, `UpdateUserRole`, `UpdateUserStatus`, `CreateCommunityReport`, `ListCommunityReports`, `ResolveCommunityReport`) are removed and `HasActiveSession` is removed from `app/`.
- [x] `.routing/DOMAIN.md` (Sections 3 & 4) reflects `app/` (no `local-app/`), dual-token BFF storage (RAM access + host `0600` `session.json` refresh + single-flight refresh + `auth:expired`), and WS `Authorization: Bearer` header transport.
- [x] `wails3 task bindings` (in `app/`) and `go build .` (in `app/`) exit with code `0`.
- [x] `go test -race -count=1 ./...` (in `app/`) runs `app_auth_test.go` and exits with code `0`.
- [x] `git status --short` confirms no files outside the declared `Scope` were created or modified.
- [ ] No unauthorized `git commit`, `git push`, or `git merge` was run. (If committed after explicit user authorization: English, Conventional Commits, detailed body.)
