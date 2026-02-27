# DemoSearch (Next.js + Go/Gin) — Full Project (fixed)

## Put Excel files
- `backend/data.xlsx`  (A-H columns)
- optional `backend/kits.xlsx`

## Run backend
```powershell
cd backend
go mod tidy
go run ./cmd/server
```

## Run frontend
```powershell
cd frontend
copy .env.example .env.local
npm install
npm run dev
```

Open: http://localhost:3000
Backend health: http://127.0.0.1:8080/api/health
