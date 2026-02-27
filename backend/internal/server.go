package internal

import (
    "net/http"
    "net/http/httputil"
    "net/url"
    "os"
    "strconv"
    "strings"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

var (
    engine  *Engine
    docs    []Doc
    kits    []KitDetail
    ruleCfg RuleConfig
)

func getenvStr(key, def string) string {
    v := os.Getenv(key)
    if v == "" {
        return def
    }
    return v
}

func getenvInt(key string, def int) int {
    v := os.Getenv(key)
    if v == "" {
        return def
    }
    n, err := strconv.Atoi(v)
    if err != nil {
        return def
    }
    return n
}

func getenvFloat(key string, def float64) float64 {
    v := os.Getenv(key)
    if v == "" {
        return def
    }
    f, err := strconv.ParseFloat(v, 64)
    if err != nil {
        return def
    }
    return f
}

// RunServer loads Excel once on startup.
func RunServer() {
    itemsPath := getenvStr("ITEMS_XLSX", "data.xlsx")
    kitsPath := getenvStr("KITS_XLSX", "kits.xlsx")
    rulesPath := getenvStr("RULES_JSON", "data/rules.json")

    titleBoost := getenvInt("TITLE_BOOST", 3)
    minScore := getenvFloat("MIN_SCORE", 0.0)
    nMin := getenvInt("NGRAM_MIN", 3)
    nMax := getenvInt("NGRAM_MAX", 6)

    d, err := LoadDocsFromExcel(itemsPath, titleBoost)
    if err != nil {
        docs = []Doc{}
    } else {
        docs = d
    }
    cfg := DefaultEngineConfig()
    cfg.MinScore = minScore
    cfg.NMin = nMin
    cfg.NMax = nMax
    engine = NewEngine(docs, cfg)

    kits, _ = LoadKitsFromExcel(kitsPath)
    ruleCfg = LoadRuleConfig(rulesPath)

    r := gin.Default()

    r.Use(cors.New(cors.Config{
        AllowOrigins: []string{"http://localhost:3000", "http://127.0.0.1:3000"},
        AllowMethods: []string{"GET", "POST", "OPTIONS"},
        AllowHeaders: []string{"Origin", "Content-Type"},
    }))

    r.GET("/api/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "ok":       true,
            "docs":     len(docs),
            "kits":     len(kits),
            "minScore": engine.Cfg.MinScore,
        })
    })

    r.GET("/api/search", func(c *gin.Context) {
        q := c.Query("q")
        k := 20
        if kk := c.Query("k"); kk != "" {
            if n, err := strconv.Atoi(kk); err == nil && n >= 1 && n <= 50 {
                k = n
            }
        }
        res := engine.Search(q, k)
        c.JSON(200, gin.H{"query": q, "results": res})
    })

    r.GET("/api/doc/:id", func(c *gin.Context) {
        id := c.Param("id")
        if d, ok := engine.GetByID(id); ok {
            c.JSON(200, d)
            return
        }
        c.JSON(http.StatusNotFound, gin.H{"detail": "Not Found"})
    })
    r.GET("/api/subcategory", func(c *gin.Context) {
        main := c.Query("main")
        sub := c.Query("sub")

        items := make([]gin.H, 0, 200)

        for _, d := range docs {
            m := d.Meta // map[string]any

            cm, _ := m["categoryMain"].(string)
            cs, _ := m["categorySub"].(string)

            if strings.TrimSpace(cm) == strings.TrimSpace(main) &&
                strings.TrimSpace(cs) == strings.TrimSpace(sub) {
                items = append(items, gin.H{
                    "id":    d.ID,
                    "title": d.Title,
                    "page":  m["page"],
                    "row":   m["row"],
                })
            }
        }

        c.JSON(200, gin.H{
            "main":  main,
            "sub":   sub,
            "count": len(items),
            "items": items,
        })
    })
    r.GET("/api/kits", func(c *gin.Context) {
        c.JSON(200, gin.H{"kits": kits})
    })

    // detail by kitId (counter)
    r.GET("/api/kits/:kitId", func(c *gin.Context) {
        kitId := c.Param("kitId")

        var found *KitDetail
        for i := range kits {
            if kits[i].KitID == kitId {
                found = &kits[i]
                break
            }
        }
        if found == nil {
            c.JSON(404, gin.H{"detail": "Not Found"})
            return
        }
        c.JSON(200, gin.H{"kit": found})
    })
    r.GET("/api/rules/config", func(c *gin.Context) {
        c.JSON(200, ruleCfg)
    })

    r.POST("/api/rules/eval", func(c *gin.Context) {
        var payload map[string]any
        if err := c.BindJSON(&payload); err != nil {
            c.JSON(400, gin.H{"error": "bad json"})
            return
        }

        inputs := map[string]float64{}
        for _, inp := range ruleCfg.Inputs {
            v, ok := payload[inp.Key]
            if !ok {
                continue
            }
            switch t := v.(type) {
            case float64:
                inputs[inp.Key] = t
            case string:
                if f, err := strconv.ParseFloat(t, 64); err == nil {
                    inputs[inp.Key] = f
                }
            default:
                // ignore
            }
        }

        budget, allTrue, conditions := EvalRules(ruleCfg, inputs)
        c.JSON(200, gin.H{
            "budgetType": budget,
            "allTrue":    allTrue,
            "conditions": conditions,
            "logicNote":  ruleCfg.LogicNote,
        })
    })
    nextURL, _ := url.Parse("http://127.0.0.1:3000")
    proxy := httputil.NewSingleHostReverseProxy(nextURL)

    r.NoRoute(func(c *gin.Context) {
        // กันไม่ให้ /api โดน proxy
        if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
            c.JSON(404, gin.H{"detail": "Not Found"})
            return
        }

        proxy.ServeHTTP(c.Writer, c.Request)
    })
    _ = r.Run(":8080")
}
