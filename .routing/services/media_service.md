# Media Service (`media_service.md`)

> **Domain service file — self-contained.** Combines the media domain rules (`RN12`, extracted from the former `DOMAIN.md`, where they are referenced per field) with the technical contract of `image-service` (from the former `INTERFACES.md`, Module 8).
> **Service:** `image-service` (Go 1.24, `:8002`) · **Business Rule:** `RN12` (Image sanitization and NSFW filter) · **Related:** `RNF-SEG-04` (EXIF stripping), `RN13` (see [`arena_service.md`](arena_service.md)) · **Routes:** backend API `/api/v1/images/upload`, `/uploads/*` (see section 2.8.2); frontend routes: `TBD`
> **Global rules:** the contract-first rules, SOA architecture, and domain-error-to-HTTP mapping shared by every service live in [`DOMAIN.md`](../DOMAIN.md) (sections 3 and 4).
> **Canonical persistence:** `image-service` owns no database tables; sanitized files live in its local volume and consuming tables store only the resulting URL. · **Shared DTOs:** [`pkg/models/`](../../pkg/models)

---

## Part A — Domain Model

### RN12 — Image Sanitization and NSFW Filter

`RN12` has no standalone section in the former `DOMAIN.md`; it is referenced on every entity field that stores an image URL. All such URLs must be produced by `image-service` (magic-bytes validation, heuristic NSFW filter, EXIF stripping, WebP output) before being persisted. The consolidated rule references are:

| Entity (Table) | Field | Source rule (verbatim from `DOMAIN.md`) | Owning service |
| :--- | :--- | :--- | :--- |
| `User` (`users`) | `avatar_url` | URL to WebP image optimized by `image-service` (**RN12**). | [`identity_service.md`](identity_service.md) |
| `EditorRequest` (`editor_requests`) | `credentials_url` | Documentary evidence or academic verification link (**RN12**). | [`governance_service.md`](governance_service.md) |
| `Word` (`words`) | `image_url` | URL to conceptual or representative image generated and sanitized by `image-service` (**RN12**). | [`dictionary_service.md`](dictionary_service.md) |
| `Sense` (`senses`) | `image_url` | URL to image specific to the concrete dialectal sense processed by `image-service` (**RN12**). | [`dictionary_service.md`](dictionary_service.md) |
| `SenseApproval` (`sense_approvals`) | trigger `trg_auto_approve_sense` (**RN06**) | Automatic promotion validates both the textual definition and the suitability of the associated image (**RN12**). | [`curation_service.md`](curation_service.md) |

Related media rules referenced elsewhere in the domain:
* **Contract fields:** `UpdateProfileRequest.avatar_url` (identity), `CreateWordProposalRequest.image_url` and `CreateSenseRequest.image_url` (dictionary) are documented as *"URL processed by image-service (RN12)"*.
* **`RN13` (Quorum of Visual Distractors):** the `WORD_TO_IMAGE` Arena modality requires at least 4 unique visual assets (`image_url IS NOT NULL`); see [`arena_service.md`](arena_service.md).
* **Upload use cases:** the upload endpoint is referenced from `CU01`, `CU02`, and `CU06` (`avatar`, `credentials`, `word`, `sense`).

---

## Part B — Technical Contracts

> [!IMPORTANT]
> Contract-first rules apply: no invented JSON fields, strict `snake_case`, and `response.Error(w, status, message)` / `response.ErrorWithCode(w, status, code, message)` (`{ "code", "error" }` in English) for every error. See [`DOMAIN.md`](../DOMAIN.md#3-global-contract-first-rules-cross-cutting).

### MODULE 8: Multimedia Processing (`image-service` in Go 1.24 — `:8002`)

#### 2.8.1. Data Structures (`pkg/models/media.go`)

```go
package models

// ImageUploadResponse represents the standardized response contract
// when uploading a media artifact in image-service (CU01, CU02, CU06, RNF-SEG-04).
type ImageUploadResponse struct {
	URL      string `json:"url"`
	Category string `json:"category"`
	MIMEType string `json:"mime_type"`
	Size     int64  `json:"size"`
	Bytes    int    `json:"bytes,omitempty"` // Compatibility numeric alias
}
```

#### 2.8.2. HTTP Endpoints Specification (`image-service`)

| Method | Path | Middleware | Use Case | Status | Successful Response | Description |
| :---: | :--- | :---: | :---: | :---: | :--- | :--- |
| `GET` | `/health` | Public | Infra | ✅ Ready | `200 OK` $\to$ `{"status":"ok","service":"image-service","timestamp":"..."}` | Operational health probe. |
| `GET` | `/` | Public | Infra | ✅ Ready | `200 OK` $\to$ `{"message":"Bienvenido a Palabrario Image Service (Golang)"}` | Service welcome message. |
| `POST` | `/api/v1/images/upload` | `RequireAuth` (`pkg/auth` via `auth.ExtractBearerToken`) | CU01 / CU02 / CU06 (RN12) | ✅ Ready | `201 Created` $\to$ `models.ImageUploadResponse` | Receives `multipart/form-data` (`file` max 5 MB + `category`: `avatar`, `credentials`, `word`, or `sense`). Validates magic bytes (`image/jpeg`, `image/png`, `image/webp`), applies heuristic skin-entropy safety filter in pure Go, strips EXIF metadata (`RNF-SEG-04`), stores in local volume, and returns `models.ImageUploadResponse`. Returns `401` English `response.ErrorWithCode` (`MISSING_TOKEN`, `INVALID_TOKEN_FORMAT`, `INVALID_TOKEN`) without leaking internal JWT verification errors in the HTTP body; `400` on invalid category/file/empty, `413` on file too large (> 5 MB), `415` on unsupported MIME, `422` on NSFW rejection. CORS uses `CORS_ALLOWED_ORIGINS` with `AllowCredentials: false`. |
| `GET` | `/uploads/*` | Public | Infra / Media | ✅ Ready | `200 OK` (binary stream) | Statically serves sanitized images from local volume of `image-service` with `X-Content-Type-Options: nosniff`. |
