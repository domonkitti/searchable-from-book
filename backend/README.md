# DemoSearch Backend (Go + Gin)

## Place files
- Put your main Excel at: `backend/data.xlsx`
  Columns: A หมวด, B หมวดย่อย, C รายการ, D หน้า, E ลำดับ, F เงื่อนไขพิเศษ, G การใช้งบ, H อำนาจเขต

- Optional kits at: `backend/kits.xlsx`
  Columns: A ชุดพัสดุ, B รายการ, C จำนวน

## Run
```powershell
cd backend
go mod tidy
go run ./cmd/server
```

Health:
- http://127.0.0.1:8080/api/health

## Tune (optional env)
- MIN_SCORE (e.g. 0.15)
- TITLE_BOOST (e.g. 3)
- NGRAM_MIN / NGRAM_MAX
