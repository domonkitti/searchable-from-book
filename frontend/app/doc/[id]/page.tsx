"use client";

import { useEffect, useMemo, useState } from "react";
import { useParams } from "next/navigation";
import Navbar from "../../components/Navbar";
import BackButton from "../../components/BackButton";

const API = process.env.NEXT_PUBLIC_API_BASE || "";

type Doc = {
  id: string;
  title: string;
  meta?: any;
};

function esc(s: any) {
  return (s ?? "").toString();
}

export default function DocDetailPage() {
  const params = useParams();
  const id = (params?.id as string) || "";

  const [doc, setDoc] = useState<Doc | null>(null);
  const [loading, setLoading] = useState(true);

  // ---- subcategory expand ----
  const [subOpen, setSubOpen] = useState(false);
  const [subLoading, setSubLoading] = useState(false);
  const [subItems, setSubItems] = useState<any[]>([]);
  const [subLoadedKey, setSubLoadedKey] = useState<string>(""); // กันโหลดซ้ำ

  useEffect(() => {
    if (!id) return;
    setLoading(true);

    fetch(`${API}/api/doc/${encodeURIComponent(id)}`, { cache: "no-store" })
      .then((r) => (r.ok ? r.json() : null))
      .then((d) => setDoc(d))
      .finally(() => setLoading(false));
  }, [id]);

  const m = useMemo(() => doc?.meta || {}, [doc]);

  async function toggleSubcategory(main: string, sub: string) {
    if (!sub) return;

    // toggle ปิด
    if (subOpen) {
      setSubOpen(false);
      return;
    }

    // toggle เปิด
    setSubOpen(true);

    // ถ้าเคยโหลดของหมวดย่อยนี้แล้ว ไม่ต้องโหลดซ้ำ
    const key = `${main}||${sub}`;
    if (subLoadedKey === key && subItems.length > 0) return;

    setSubLoading(true);
    setSubItems([]);
    setSubLoadedKey(key);

    try {
      const res = await fetch(
        `${API}/api/subcategory?main=${encodeURIComponent(main)}&sub=${encodeURIComponent(sub)}`,
        { cache: "no-store" }
      );
      const data = await res.json();
      setSubItems(data.items || []);
    } finally {
      setSubLoading(false);
    }
  }

  const catMain = esc(m.categoryMain || "-");
  const catSub = esc(m.categorySub || "");

  return (
    <main className="wrap">
      <Navbar />
      <BackButton />
      <div className="card">

        {loading && (
          <div className="small" style={{ marginTop: 10 }}>
            กำลังโหลด…
          </div>
        )}

        {!loading && !doc && (
          <div className="small" style={{ marginTop: 10 }}>
            ไม่พบข้อมูล
          </div>
        )}

        {doc && (
          <>
            <div className="title no-hover" style={{ marginTop: 10 }}>
              {esc(doc.title)}
            </div>

            <div className="small" style={{ marginTop: 8 }}>
              <div>
                <b>หมวด:</b> {catMain}
              </div>

              {catSub && (
                <div style={{ marginTop: 6 }}>
                  <b>หมวดย่อย:</b> {catSub}{" "}

                  <button
                    type="button"
                    onClick={() => toggleSubcategory(catMain, catSub)}
                    aria-expanded={subOpen}
                    style={{
                      background: "none",
                      border: "none",
                      padding: 0,
                      marginLeft: 8,
                      cursor: "pointer",
                      font: "inherit",
                      fontSize: 13,
                      opacity: 0.7,
                      color: "#007BFF"
                    }}
                  >
                    {subOpen ? "ซ่อนรายการ ▲" : "ดูรายการ ▼"}
                  </button>

                  {subOpen && (
                    <div style={{ marginTop: 8, paddingLeft: 22 }}>
                      {subLoading && <div className="small">กำลังโหลด…</div>}

                      {!subLoading && subItems.length === 0 && (
                        <div className="small">ไม่พบรายการในหมวดย่อยนี้</div>
                      )}

                      {!subLoading && subItems.length > 0 && (
                        <ul style={{ margin: 0, paddingLeft: 18 }}>
                          {subItems.map((it: any, idx: number) => (
                            <li key={it.id ?? idx} style={{ marginTop: 6 }}>
                              <a
                                href={`/doc/${encodeURIComponent(it.id)}`}
                                style={{ textDecoration: "none" }}
                              >
                                {esc(it.title)}
                              </a>
                            </li>
                          ))}
                        </ul>
                      )}
                    </div>
                  )}
                </div>
              )}
            </div>

            <div className="small" style={{ marginTop: 10 }}>
              <b>อ้างอิง:</b> หน้า {esc(m.page || "-")} ลำดับ {esc(m.row || "-")}
            </div>

            {!!esc(m.budgetUse).trim() && (
              <div className="small" style={{ marginTop: 6 }}>
                <b>การใช้งบ:</b> {esc(m.budgetUse)}
              </div>
            )}

            {!!esc(m.authority).trim() && (
              <div className="small" style={{ marginTop: 6 }}>
                <b>อำนาจเขต:</b> {esc(m.authority)}
              </div>
            )}

            {!!esc(m.special).trim() && (
              <div className="card" style={{ marginTop: 12 }}>
                <div className="small">
                  <b>เงื่อนไขพิเศษ</b>
                </div>
                <div className="small" style={{ marginTop: 8, whiteSpace: "pre-wrap" }}>
                  {esc(m.special)}
                </div>
              </div>
            )}
          </>
        )}
      </div>
    </main>
  );
}