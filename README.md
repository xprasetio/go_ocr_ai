# OCR AI

Fullstack OCR receipt analyzer.

## Backend
- Go Echo + DDD + DI Sarulabs + Raw SQL + Viper + Cobra
- Cloudflare D1 for persistence
- Gemini for OCR parsing
- Supabase bucket `ocr_ai_receipt` for file storage

## Run
1. Copy `env.example` to `.env`
2. Fill Cloudflare, Supabase, and Gemini credentials
3. Run `go mod tidy`
4. Run `go run ./cmd/api`

Migration runs automatically using `CREATE TABLE IF NOT EXISTS`.

## Frontend
- Next.js + TypeScript + Tailwind CSS ada di folder `web`

Run frontend:
1. `cd web`
2. `cp .env.example .env.local`
3. `npm install`
4. `npm run dev`
