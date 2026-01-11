# Repository Guidelines

## Project Structure & Module Organization
- `backend/` contains the Go API server, with request handlers in `backend/handlers/`, DB access in `backend/storage/`, domain types in `backend/models/`, and migrations in `backend/migrations/`.
- `backend/cmd/` hosts CLI tools (`seed`, `extract-recipe`). `backend/units/` and `backend/logger/` are shared packages.
- `frontend/` is a Vite + React + TypeScript app. UI pages live in `frontend/src/pages/`, shared components in `frontend/src/components/`, API helpers in `frontend/src/api.ts`, and generated types in `frontend/src/types.gen.ts`.
- Tests are co-located: Go tests are `*_test.go` alongside packages; frontend E2E tests are in `frontend/e2e/`. Static assets live in `frontend/public/`.

## Build, Test, and Development Commands
- Backend (from `backend/`):
  - `make run` starts the API server.
  - `make test` runs Go tests using the test DB and migrations.
  - `make fmt` / `make vet` apply Go formatting and vet checks.
  - `make migrate-up` / `make migrate-down` apply DB migrations.
  - `make seed` seeds local data.
  - `make types` regenerates `frontend/src/types.gen.ts` via tygo.
- Frontend (from `frontend/`):
  - `npm run dev` starts the Vite dev server.
  - `npm run build` type-checks and builds production assets.
  - `npm run lint` runs ESLint.
  - `npx playwright test` runs E2E tests in `frontend/e2e/`.

## Coding Style & Naming Conventions
- Go code is formatted with `gofmt`; keep handler logic in `handlers/` and DB access in `storage/`.
- TypeScript uses `.tsx` for components; prefer PascalCase for component files (e.g., `RecipeSteps.tsx`).
- Treat `frontend/src/types.gen.ts` as generated; update via `make types` after changing backend models.

## Testing Guidelines
- Backend tests rely on `DATABASE_URL_TEST` (see `backend/Makefile`); `make test` prepares a test DB.
- Add tests for new handlers and storage operations (`*_test.go`).
- E2E tests live in `frontend/e2e/` and should cover critical UI flows.

## Commit & Pull Request Guidelines
- Commit messages are short, imperative, and capitalized (e.g., “Fixing failing test”, “Updating gitignore”).
- PRs should include a brief summary, test commands run, and linked issues. Add screenshots or clips for UI changes and note any migration steps.

## Configuration & Secrets
- Backend uses `DATABASE_URL`, `OPENAI_API_KEY`, `OPENAI_MODEL`, and `PORT` in environment config.
- Frontend expects `VITE_API_BASE` for the API URL.
