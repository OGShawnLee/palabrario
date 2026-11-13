# 032 — Harden Wails BFF Session Persistence, Auto-Refresh, Error Codes, and Tests

| Field | Value |
| :--- | :--- |
| **ID** | 032 |
| **Status** | completed |
| **Model** | Gemini 4 Hydrogen, `thinking_level: medium` (`gemini-4-hydrogen`) |
| **Role** | Desktop BFF Implementer (`implementation`) |
| **Spec** | [spec.md](spec.md) |

> **Specification Note:**  
> Fourth ticket of the end-to-end auth refactor (depends on `031`). Upgrades `app/` to store short-lived access tokens in RAM while persisting the opaque refresh token cleanly on the host filesystem (`os.UserConfigDir()/palabrario/session.json`, `0600`), implements single-flight auto-refresh on `401 TOKEN_EXPIRED` + automatic cleanup & `auth:expired` Wails event emission on `403 ACCOUNT_SUSPENDED` / failed refresh, standardizes typed `APIError{Status, Code, Message}` in English, passes Arena WebSocket auth via `Authorization` header instead of `?token=`, removes the 14 duplicate BFF alias methods, regenerates Wails TS bindings, and adds `app/app_auth_test.go`.

## Scope

May create or edit:
- `app/app.go`
- `app/app_auth.go`
- `app/app_arena.go`
- `app/app_governance.go`
- `app/app_media.go`
- `app/app_auth_test.go` (create)
- `app/frontend/bindings/` (exclusively regenerated via `wails3 task bindings` or `wails3 task check` — note: frontend TypeScript consumers are updated in Ticket `033`)
- `.routing/DOMAIN.md` (Sections 3 & 4: `local-app/` → `app/`, dual-token RAM + disk `session.json` persistence, single-flight refresh, `auth:expired` eviction, WS `Authorization: Bearer` header edge)
- `.tickets/032-harden-wails-bff-auth/changelog.md` (create)

Must not touch:
- `SCHEMA_FINAL.sql`, `sql/`, `pkg/`, `core-service/`, `realtime-service/`, `image-service/`, `app/frontend/src/`, `.context/`, `.routing/services/`, `.routing/API.md`, `.routing/CONTEXT.md`, `.routing/DESIGN.md`, `.routing/PROGRESS.md`, `.routing/ROUTES.md`, `.routing/INFRASTRUCTURE.md`
- Any file outside the above scope
- **No role (Router, Investigator, Implementer, Translator, or Reviewer) may commit, push, or merge without explicit user authorization.** Do not run `git commit`, `git push`, or `git merge`.

## Inputs

- [.tickets/028-investigate-auth-workflow/report.md](../028-investigate-auth-workflow/report.md): Findings `BFF-01`..`BFF-05`, `SIB-03`.
- [app/app.go](../../app/app.go), [app/app_auth.go](../../app/app_auth.go), [app/app_arena.go](../../app/app_arena.go), [app/app_governance.go](../../app/app_governance.go), [app/app_media.go](../../app/app_media.go)
- [pkg/models/auth.go](../../pkg/models/auth.go) and [pkg/response/json.go](../../pkg/response/json.go)
- [.routing/DOMAIN.md](../../.routing/DOMAIN.md)

## Task

1. **Typed `APIError`, Single-Flight Refresh & Host Refresh-Token Persistence ([app/app.go](../../app/app.go)):**
   - Define `APIError`:
     ```go
     type APIError struct {
         Status  int    `json:"status"`
         Code    string `json:"code"`
         Message string `json:"message"`
     }
     func (e *APIError) Error() string { return e.Message }
     ```
   - Replace `sessionToken string` on `App` with `accessToken string`, `refreshToken string`, `sessionFilePath string` (defaulting in `NewApp()` to `filepath.Join(cfgDir, "palabrario", "session.json")` via `os.UserConfigDir()` with fallback to `os.TempDir()` only if `UserConfigDir` fails), and `refreshMu sync.Mutex`.
   - Load persisted `{ "refresh_token": "..." }` from `sessionFilePath` on `NewApp()` startup (ignore missing/corrupt file gracefully). Persist on `setTokens(access, refresh)` (`os.MkdirAll(..., 0700)` + atomic write `0600`) and remove file on `clearSession()`. Add an unexported test constructor `newAppWithConfig(coreURL, realtimeURL, imageURL, sessionFilePath string)` so unit tests stay isolated.
   - Replace `doRequest` with `executeRequest(ctx context.Context, buildReq func() (*http.Request, error), target any, allowRetry bool) error` + `doRequest(ctx, method, targetURL, body, target)` wrapper:
     - Inject `Authorization: Bearer <accessToken>` if non-empty and `Accept: application/json`.
     - If response is `401 Unauthorized` (`allowRetry && tokenWasSent && a.getRefreshToken() != ""`), run `a.refreshSession(ctx)` (synchronized under `refreshMu` so concurrent `401`s trigger one `/api/v1/auth/refresh` call) and retry `executeRequest` once with `allowRetry = false`.
     - If response is `403` with `Code == "ACCOUNT_SUSPENDED"`, or `401` after failed/absent refresh (for authenticated calls), call `a.clearSession()` and emit Wails event `"auth:expired"` with `map[string]string{"code": apiErr.Code, "reason": apiErr.Message}` (safe nil check on `application.Get()` for Go tests).
     - Wrap network failures in English without leaking internal transport stack noise (`"Unable to reach the Palabrario server. Please check your connection."`).

2. **Update Auth Methods & Remove `HasActiveSession` Bloat ([app/app.go](../../app/app.go), [app/app_auth.go](../../app/app_auth.go), [app/app_governance.go](../../app/app_governance.go)):**
   - Remove `HasActiveSession()` from `app.go` (note: Ticket `033` updates `app/frontend/src/`, so `wails3 task bindings` in `032` regenerates `app.ts` cleanly, while `wails3 task check` waits for Ticket `033`—do NOT run `vue-tsc` expecting it to pass if frontend still references removed aliases/methods until Ticket `033`; instead verify `go build .` + `wails3 task bindings`!).
   - In `SignIn` & `SignUp`: save both `authResp.Token` and `authResp.RefreshToken` via `setTokens`, and return `*models.UserProfile`.
   - In `SignOut`: best-effort `POST /api/v1/auth/sign-out` with `{RefreshToken}` (ignore network errors), call `clearSession()`, and `LeaveArenaRoom()`.
   - In `GetMe`: if `accessToken == ""` and `refreshToken != ""`, call `refreshSession` first (which returns the `UserProfile` cleanly on app startup!); if neither token exists, return `(nil, nil)` (`null` in JS, matching `MockAuthAPI`) rather than firing an unauthenticated `/auth/me` call; on `401`/`403` expired/suspended session, return `(nil, nil)` (or return `APIError` for transient `500`/network error so UI distinguishes server outage from logged-out).
   - In `UpdateProfile` (`app_auth.go`) and `RedeemInvitationTicket` (`app_governance.go`): after receiving the updated `*models.UserProfile`, call `a.refreshSession(a.context())` best-effort (ignoring transient error or updating tokens) so `accessToken` immediately reflects the new `AnchorDialect` or promoted `Role`.
   - Keep `UploadImage` in `app_media.go` calling `a.executeRequest(...)` so it shares Bearer injection, `401` auto-refresh, `auth:expired`, and `APIError`. Translate `app_media.go` category error to English.

3. **WebSocket Header Auth, Remove 14 Duplicate BFF Aliases & Sync `DOMAIN.md` (`app_auth.go`, `app_arena.go`, `app_governance.go`, `.routing/DOMAIN.md`):**
   - In `app_arena.go`:
     - In `FindMatchRoom`, check `var apiErr *APIError; errors.As(err, &apiErr) && apiErr.Status == http.StatusNotFound` instead of substring matching (`"not found"` / `"no encontrada"`).
     - In `ConnectArenaRoom`: if `accessToken == ""` and `refreshToken != ""`, call `refreshSession` first; pass `http.Header{"Authorization": {"Bearer " + token}}` to `websocket.DefaultDialer.DialContext` instead of `?token=...` query string. Translate error messages in `app_arena.go` to English.
   - Remove the 14 duplicate alias methods verified unused in `app/frontend/src/`:
     - `GetMyProfile` (`app_auth.go`)
     - `CreateArenaMatch`, `GetMyMatchHistory`, `GetMatchReplayDetail`, `ConnectArena`, `LeaveArena` (`app_arena.go`)
     - `ListPendingEditorRequests`, `GenerateInvitationTickets`, `ListUsersForGovernance`, `UpdateUserRole`, `UpdateUserStatus`, `CreateCommunityReport`, `ListCommunityReports`, `ResolveCommunityReport` (`app_governance.go`)
   - Regenerate bindings with `wails3 task bindings` in `app/`.
   - Update [.routing/DOMAIN.md](../../.routing/DOMAIN.md) Sections 3 and 4: replace stale `local-app/` and `./local-app` (`go.work`) references with `app/` (`app/app.go`, `app/app_auth.go`), document the dual-token BFF architecture (short-lived access JWT in RAM + persistent host `session.json` `0600` refresh token + single-flight auto-refresh + `auth:expired` eviction), and update the Mermaid WebSocket edge from `?token={jwt}` to `Authorization: Bearer {jwt}` header.

4. **Add Comprehensive `app/app_auth_test.go`:**
   - Using `net/http/httptest.Server`, `t.TempDir()`, and `newAppWithConfig`, test:
     - `SignIn` persists refresh token to disk (`0600`) and clears on `SignOut`.
     - Restart recovery: new `App` instance reading the persisted `session.json` exchanges refresh token via `/api/v1/auth/refresh` on `GetMe()`.
     - Transparent auto-refresh on `401 TOKEN_EXPIRED` during `doRequest` and `UploadImage` (retries original request with new token; concurrent `401` calls trigger only 1 refresh request).
     - `403 ACCOUNT_SUSPENDED` and failed refresh clear RAM + delete `session.json`.
     - `UpdateProfile` and `RedeemInvitationTicket` rotate/refresh `accessToken`.
     - `ConnectArenaRoom` sends `Authorization: Bearer ...` header and no `?token=` query parameter.

## Verification

Run every check in [spec.md](spec.md).

## Change log

Write `.tickets/032-harden-wails-bff-auth/changelog.md` recording session file path/format, single-flight refresh flow, removed aliases, regenerated bindings output, and `go test -race -count=1 ./...` in `app/`.
