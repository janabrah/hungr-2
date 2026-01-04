# Migration Plan: Next.js to Vite + React + Go

## Overview

Migrate from Next.js monolith to:
- **Frontend**: Vite + React + React Router
- **Backend**: Go HTTP server (chi router)
- **Storage**: Supabase Storage (replacing Vercel Blob)
- **Database**: Supabase PostgreSQL (unchanged)

## Current Architecture

```
Next.js App
├── / (home page)
├── /upload_recipe (upload form)
├── /show_recipe (browse recipes)
└── /api/recipe/upload (API route)
    ├── POST - upload recipe + images
    └── GET - fetch recipes for user
```

## Target Architecture

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   Vite + React  │────▶│    Go Server    │────▶│    Supabase     │
│   (Frontend)    │     │    (Backend)    │     │  (DB + Storage) │
└─────────────────┘     └─────────────────┘     └─────────────────┘
     :5173                   :8080
```

---

## Phase 1: Set Up Go Backend

### 1.1 Initialize Go Module

```bash
mkdir -p backend
cd backend
go mod init github.com/yourusername/hungr
```

### 1.2 Dependencies

```bash
go get github.com/go-chi/chi/v5
go get github.com/go-chi/cors
go get github.com/supabase-community/supabase-go
go get github.com/joho/godotenv
```

### 1.3 Project Structure

```
backend/
├── main.go              # Entry point, server setup
├── handlers/
│   └── recipe.go        # Recipe CRUD handlers
├── storage/
│   └── supabase.go      # Supabase client + storage helpers
├── models/
│   └── recipe.go        # Data structures
└── .env                 # Environment variables
```

### 1.4 API Endpoints to Implement

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/recipes` | Upload recipe with images |
| GET | `/api/recipes?user_id=X` | Get all recipes for user |
| GET | `/api/recipes/:id` | Get single recipe (future) |

### 1.5 Go Data Models

```go
// models/recipe.go
package models

type Recipe struct {
    ID        int       `json:"id"`
    Filename  string    `json:"filename"`
    UserID    int       `json:"user_id"`
    TagString string    `json:"tag_string"`
    CreatedAt time.Time `json:"created_at"`
}

type File struct {
    ID    int    `json:"id"`
    URL   string `json:"url"`
    Image bool   `json:"image"`
}

type FileRecipe struct {
    FileID     int `json:"file_id"`
    RecipeID   int `json:"recipe_id"`
    PageNumber int `json:"page_number"`
}

type Tag struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

type RecipeTag struct {
    RecipeID int `json:"recipe_id"`
    TagID    int `json:"tag_id"`
}

// Response types
type RecipesResponse struct {
    RecipeData  []Recipe     `json:"recipeData"`
    FileData    []File       `json:"fileData"`
    MappingData []FileRecipe `json:"mappingData"`
}

type UploadResponse struct {
    Success bool   `json:"success"`
    Recipe  Recipe `json:"recipe"`
    Tags    []Tag  `json:"tags"`
}
```

### 1.6 Key Handler Logic (POST /api/recipes)

The upload handler needs to:

1. Parse multipart form data
2. Extract `filename` and `tagString` from query params
3. For each file:
   - Upload to Supabase Storage bucket
   - Get public URL
4. Insert recipe into `recipes` table
5. Insert file records into `files` table
6. Create `file_recipes` mappings
7. Parse tags, create hash-based IDs, upsert into `tags` table
8. Create `recipe_tags` mappings
9. Return success response

```go
// Tag ID generation (matches existing JS logic)
func createTagID(tag string) int {
    h := sha256.New()
    h.Write([]byte(tag))
    hash := hex.EncodeToString(h.Sum(nil))
    id, _ := strconv.ParseInt(hash[:8], 16, 64)
    return int(id)
}
```

### 1.7 Supabase Storage Setup

Create a storage bucket in Supabase dashboard:
- Bucket name: `recipe-images`
- Public: Yes (for direct image access)

```go
// Upload file to Supabase Storage
func UploadFile(bucket, filename string, file io.Reader) (string, error) {
    // Use Supabase Storage API
    // Returns public URL
}
```

---

## Phase 2: Set Up Vite + React Frontend

### 2.1 Initialize Vite Project

```bash
npm create vite@latest frontend -- --template react-ts
cd frontend
npm install
npm install react-router-dom
npm install -D tailwindcss postcss autoprefixer
npx tailwindcss init -p
```

### 2.2 Project Structure

```
frontend/
├── src/
│   ├── main.tsx           # Entry point with router
│   ├── App.tsx            # Root component
│   ├── pages/
│   │   ├── Home.tsx
│   │   ├── UploadRecipe.tsx
│   │   └── ShowRecipe.tsx
│   ├── components/
│   │   └── HomeButton.tsx
│   ├── types/
│   │   └── recipe.ts
│   ├── api/
│   │   └── recipes.ts     # API client functions
│   └── index.css          # Tailwind imports
├── tailwind.config.js
├── vite.config.ts
└── .env
```

### 2.3 Router Setup

```tsx
// src/main.tsx
import { BrowserRouter, Routes, Route } from 'react-router-dom'
import Home from './pages/Home'
import UploadRecipe from './pages/UploadRecipe'
import ShowRecipe from './pages/ShowRecipe'

ReactDOM.createRoot(document.getElementById('root')!).render(
  <BrowserRouter>
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/upload_recipe" element={<UploadRecipe />} />
      <Route path="/show_recipe" element={<ShowRecipe />} />
    </Routes>
  </BrowserRouter>
)
```

### 2.4 Component Changes Required

| File | Changes |
|------|---------|
| `Home.tsx` | Remove Next.js `Link`, use React Router `Link` |
| `UploadRecipe.tsx` | Remove `@vercel/blob` types, update API URL |
| `ShowRecipe.tsx` | Replace `next/image` with `<img>`, update API URL |
| `HomeButton.tsx` | Use React Router `Link` |

### 2.5 API Client

```typescript
// src/api/recipes.ts
const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080'

export async function uploadRecipe(
  files: File[],
  filename: string,
  tagString: string
): Promise<UploadResponse> {
  const formData = new FormData()
  files.forEach(file => formData.append('file', file))

  const response = await fetch(
    `${API_BASE}/api/recipes?filename=${filename}&tagString=${tagString}`,
    { method: 'POST', body: formData }
  )
  return response.json()
}

export async function getRecipes(userId: number): Promise<RecipesResponse> {
  const response = await fetch(
    `${API_BASE}/api/recipes?user_id=${userId}`
  )
  return response.json()
}
```

### 2.6 Environment Variables

```bash
# frontend/.env
VITE_API_URL=http://localhost:8080
```

### 2.7 Vite Proxy (Development)

```typescript
// vite.config.ts
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/api': 'http://localhost:8080'
    }
  }
})
```

---

## Phase 3: Migration Steps (Ordered)

### Step 1: Backend Foundation
- [ ] Create `backend/` directory with Go module
- [ ] Set up chi router with CORS middleware
- [ ] Create Supabase client connection
- [ ] Implement health check endpoint

### Step 2: Supabase Storage
- [ ] Create `recipe-images` bucket in Supabase dashboard
- [ ] Implement file upload to Supabase Storage
- [ ] Test file upload returns public URL

### Step 3: GET Endpoint
- [ ] Implement `GET /api/recipes?user_id=X`
- [ ] Query recipes, files, file_recipes tables
- [ ] Return JSON matching current response format
- [ ] Test with curl/Postman

### Step 4: POST Endpoint
- [ ] Implement multipart form parsing
- [ ] Implement file upload loop
- [ ] Implement database inserts (recipes, files, file_recipes)
- [ ] Implement tag parsing and hash ID generation
- [ ] Implement tag upsert and recipe_tags linking
- [ ] Test full upload flow

### Step 5: Frontend Scaffold
- [ ] Create Vite project with React + TypeScript
- [ ] Install and configure Tailwind CSS
- [ ] Set up React Router
- [ ] Create page components (empty shells)

### Step 6: Migrate Pages
- [ ] Migrate `Home.tsx` (simplest)
- [ ] Migrate `HomeButton.tsx` component
- [ ] Migrate `ShowRecipe.tsx` (replace next/image with img)
- [ ] Migrate `UploadRecipe.tsx` (remove Vercel Blob types)

### Step 7: Integration Testing
- [ ] Test upload flow end-to-end
- [ ] Test browse/display flow end-to-end
- [ ] Test tag filtering
- [ ] Fix any CORS issues

### Step 8: Cleanup
- [ ] Remove old `hungr/` Next.js directory
- [ ] Update CLAUDE.md
- [ ] Update deployment configuration

---

## Phase 4: Deployment Options

### Option A: Single VPS (Simplest)
- Run Go binary + serve static files
- Use nginx as reverse proxy
- Deploy to DigitalOcean/Linode/Fly.io

### Option B: Separate Services
- Frontend: Vercel/Netlify/Cloudflare Pages (free static hosting)
- Backend: Fly.io/Railway/Render (Go hosting)

### Option C: Containerized
- Dockerfile for Go backend
- Dockerfile for frontend (nginx + static files)
- Docker Compose for local dev
- Deploy to any container platform

---

## Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| File upload handling differences | Medium | Medium | Test with various file sizes/types |
| CORS issues | High | Low | Configure chi/cors properly |
| Supabase Storage API differences | Low | Medium | Read Supabase Go SDK docs |
| Image display issues | Low | Low | Standard img tags work everywhere |
| Tag ID hash mismatch | Medium | High | Port exact same hash logic |

---

## Environment Variables Summary

### Backend (.env)
```
SUPABASE_URL=https://xxx.supabase.co
SUPABASE_SERVICE_KEY=eyJ...
PORT=8080
```

### Frontend (.env)
```
VITE_API_URL=http://localhost:8080
VITE_SUPABASE_URL=https://xxx.supabase.co  # Only if needed client-side
```

---

## Estimated Effort

| Phase | Effort |
|-------|--------|
| Phase 1: Go Backend | 2-3 hours |
| Phase 2: Vite Frontend | 1-2 hours |
| Phase 3: Integration | 1-2 hours |
| Phase 4: Deployment | 1-2 hours |
| **Total** | **5-9 hours** |

---

## Quick Start Commands

```bash
# Terminal 1: Backend
cd backend
go run main.go

# Terminal 2: Frontend
cd frontend
npm run dev

# Open http://localhost:5173
```
