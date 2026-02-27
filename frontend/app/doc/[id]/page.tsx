"use client";

import { useEffect, useMemo, useState } from "react";
import { useParams } from "next/navigation";
import Navbar from "../../components/Navbar";
import BackButton from "../../components/BackButton";

const API = process.env.NEXT_PUBLIC_API_BASE || "";
fetch(`${API}/api/kits`);

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

  useEffect(() => {
    if (!id) return;
    setLoading(true);

    fetch(`${API}/api/doc/${encodeURIComponent(id)}`, { cache: "no-store" })
      .then((r) => (r.ok ? r.json() : null))
      .then((d) => setDoc(d))
      .finally(() => setLoading(false));
  }, [id]);

  const m = useMemo(() => doc?.meta || {}, [doc]);

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
            <div className="title" style={{ marginTop: 10 }}>
              {esc(doc.title)}
            </div>

            {/* meta ตาม excel ใหม่ของคุณ */}
            <div className="small" style={{ marginTop: 8 }}>
            <div>
              <b>หมวด:</b> {esc(m.categoryMain || "-")}
            </div>

              {m.categorySub && (
                <div style={{ marginTop: 4 }}>
                  <b>หมวดย่อย:</b> {esc(m.categorySub)}
                </div>
              )}
            </div>

            <div className="small" style={{ marginTop: 6 }}>
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
                {/* ให้ขึ้นบรรทัดตาม \n */}
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