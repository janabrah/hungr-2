# Hungr

A recipe tracking webapp. Live at https://hungr.dev

## Features

- Upload recipe images with names and tags
- Browse and search recipes by tag
- View and edit structured recipe steps with ingredients
- Import recipes from URLs using AI extraction (OpenAI)
- Unit conversion for ingredients (cups, tsp, grams, etc.)

## Development

### Prerequisites

- Node.js 18+
- Go 1.21+
- PostgreSQL
- OpenAI API key (for recipe extraction feature)

### Running the Dev Servers

**Frontend:**
```bash
cd frontend
npm install
npm run dev
# Runs at http://localhost:5173
```

**Backend:**
```bash
cd backend
# Set up .env with DATABASE_URL and OPENAI_API_KEY
make run
# Runs at http://localhost:8080
```

Or use `go run .` directly:
```bash
cd backend
DATABASE_URL="postgres://localhost/hungr_db" go run .
```

### Linters and Tests

**Frontend:**
```bash
cd frontend
npm run lint         # ESLint
npm run format:check # Prettier check
npm run build        # TypeScript type check
```

**Backend:**
```bash
cd backend
go build ./...       # Compile check
DATABASE_URL="postgres://localhost/hungr_db" go test ./...
```

### Database Migrations

Migrations use Goose and are run locally via make scripts:

```bash
cd backend
make migrate-up       # Run migrations on local database
make migrate-prod-up  # Run migrations on production database
```

## Deployment

### Frontend (Vercel)
- Deploys automatically on push to main
- Connected to the `frontend/` directory

### Backend (Render)
- Manual deploys via Render dashboard
- Requires `DATABASE_URL`, `OPENAI_API_KEY`, and `OPENAI_MODEL` env vars

### Database Migrations
- Must be run manually before deploying backend changes that require schema updates
- Use `make migrate-prod-up` from the backend directory
