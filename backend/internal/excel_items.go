package internal

import (
    "fmt"
    "os"
    "strings"

    "github.com/xuri/excelize/v2"
)

// Column mapping for data.xlsx (0-index): A=0, B=1, ...
const (
    ColCategoryMain = 0 // A หมวด
    ColCategorySub  = 1 // B หมวดย่อย
    ColTitle        = 2 // C รายการ
    ColPage         = 3 // D หน้า
    ColOrder        = 4 // E ลำดับ
    ColSpecial      = 5 // F เงื่อนไขพิเศษ
    ColBudgetUse    = 6 // G การใช้งบ
    ColAuthority    = 7 // H อำนาจเขต
)

var headerKeywords = []string{"หมวด", "หมวดย่อย", "รายการ", "หน้า", "ลำดับ", "เงื่อนไข", "การใช้งบ", "อำนาจ", "category", "title", "page", "order"}

func looksLikeHeader(row []string) bool {
    joined := strings.ToLower(strings.TrimSpace(strings.Join(row, " ")))
    if joined == "" {
        return false
    }
    for _, k := range headerKeywords {
        if strings.Contains(joined, strings.ToLower(k)) {
            return true
        }
    }
    return false
}

func getCol(row []string, idx int) string {
    if idx < 0 || idx >= len(row) {
        return ""
    }
    return strings.TrimSpace(row[idx])
}

func defaultDash(s string) string {
    if strings.TrimSpace(s) == "" {
        return "-"
    }
    return strings.TrimSpace(s)
}

func nonEmpty(xs []string) []string {
    out := make([]string, 0, len(xs))
    for _, x := range xs {
        x = strings.TrimSpace(x)
        if x != "" {
            out = append(out, x)
        }
    }
    return out
}

// LoadDocsFromExcel reads ALL sheets from Excel. Each non-empty row becomes 1 doc.
func LoadDocsFromExcel(path string, titleBoost int) ([]Doc, error) {
    if _, err := os.Stat(path); err != nil {
        return nil, err
    }

    f, err := excelize.OpenFile(path)
    if err != nil {
        return nil, err
    }
    defer func() { _ = f.Close() }()

    sheets := f.GetSheetList()
    docs := make([]Doc, 0, 2048)

    counter := 0

    for _, sh := range sheets {
        rows, err := f.GetRows(sh)
        if err != nil {
            continue
        }

        for rIdx, row := range rows {
            // Skip header if first row looks like header
            if rIdx == 0 && looksLikeHeader(row) {
                continue
            }

            title := getCol(row, ColTitle)
            if title == "" || title == "รายการ" || title == "หมวด" {
                continue
            }

            catMain := getCol(row, ColCategoryMain)
            catSub := getCol(row, ColCategorySub)
            page := getCol(row, ColPage)
            orderNo := getCol(row, ColOrder)
            special := getCol(row, ColSpecial)
            budgetUse := getCol(row, ColBudgetUse)
            authority := getCol(row, ColAuthority)

            joined := strings.Join(nonEmpty(row), " | ")
            boosted := strings.TrimSpace(strings.Repeat(title+" ", titleBoost))
            fullText := fmt.Sprintf("%s| %s", boosted, joined)

            counter++
            id := fmt.Sprintf("%d", counter) // ✅ URL-safe ID

            meta := map[string]any{
                "source":       "excel",
                "categoryMain": defaultDash(catMain),
                "categorySub":  strings.TrimSpace(catSub),
                "page":         defaultDash(page),
                "row":          defaultDash(orderNo),
                "budgetUse":    strings.TrimSpace(budgetUse),
                "authority":    strings.TrimSpace(authority),
                "special":      strings.ReplaceAll(special, "\r\n", "\n"),
            }

            docs = append(docs, Doc{ID: id, Title: title, Text: fullText, Meta: meta})
        }
    }

    return docs, nil
}
