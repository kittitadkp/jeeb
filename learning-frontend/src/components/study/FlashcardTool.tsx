import { useMemo, useRef, useState } from "react";

import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { StatCard } from "@/components/ui/StatCard";
import { C, R, S, T, W } from "@/lib/design";
import type { Item } from "@/types";

interface Props {
  items: Item[];
  topicId: string;
  progressMap: Record<string, string>;
  onUpsert: (itemId: string, status: "learning" | "mastered") => void;
}

export function FlashcardTool({ items, progressMap, onUpsert }: Props) {
  const queue = useMemo(() => {
    const unmastered = items.filter((i) => progressMap[i.id] !== "mastered");
    const pool = unmastered.length > 0 ? unmastered : items;
    return [...pool].sort(() => Math.random() - 0.5);
  }, [items, progressMap]);

  const [idx, setIdx] = useState(0);
  const [revealed, setRevealed] = useState(false);
  const [done, setDone] = useState(false);
  const [correct, setCorrect] = useState(0);
  const startTime = useRef(Date.now());
  const [elapsed, setElapsed] = useState(0);

  const item = queue[idx];
  const total = queue.length;

  function advance() {
    setRevealed(false);
    if (idx + 1 >= total) {
      setElapsed(Math.round((Date.now() - startTime.current) / 1000));
      setDone(true);
    } else {
      setIdx((i) => i + 1);
    }
  }

  function handleKnow() {
    onUpsert(item.id, "mastered");
    setCorrect((c) => c + 1);
    advance();
  }

  function handleLearning() {
    onUpsert(item.id, "learning");
    advance();
  }

  function formatTime(s: number) {
    const m = Math.floor(s / 60);
    const sec = s % 60;
    return m > 0 ? `${m}m ${sec}s` : `${sec}s`;
  }

  if (done) {
    return (
      <div style={{ display: "flex", flexDirection: "column", alignItems: "center", gap: S[6], padding: S[8] }}>
        <div style={{ fontSize: T["2xl"], fontWeight: W.bold, color: C.text }}>Session complete!</div>
        <div style={{ display: "flex", gap: S[6] }}>
          <StatCard label="Marked known" value={correct} color="#16A34A" />
          <StatCard label="Time taken" value={formatTime(elapsed)} />
        </div>
        <Button onClick={() => { setIdx(0); setRevealed(false); setDone(false); setCorrect(0); startTime.current = Date.now(); }}>
          Study again
        </Button>
      </div>
    );
  }

  return (
    <div style={{ display: "flex", flexDirection: "column", alignItems: "center", gap: S[5], padding: S[4] }}>
      <div style={{ width: "100%", maxWidth: 480 }}>
        <div style={{ display: "flex", justifyContent: "space-between", marginBottom: S[3], fontSize: T.sm, color: C.text2 }}>
          <span>{idx + 1} / {total}</span>
          <span style={{ color: C.text3 }}>{item.category}</span>
        </div>

        <div style={{ height: 4, background: C.surface3, borderRadius: R.full, marginBottom: S[5] }}>
          <div style={{ height: "100%", width: `${((idx + 1) / total) * 100}%`, background: C.primary, borderRadius: R.full, transition: "width .3s" }} />
        </div>

        <Card style={{ minHeight: 180, display: "flex", flexDirection: "column", alignItems: "center", justifyContent: "center", textAlign: "center", gap: S[3] }}>
          <div style={{ fontSize: T["3xl"], fontWeight: W.bold, fontFamily: "monospace", color: C.text, letterSpacing: 2 }}>
            {item.term}
          </div>
          {revealed && (
            <>
              <div style={{ width: 40, height: 1, background: C.border2 }} />
              <div style={{ fontSize: T.md, color: C.text2, fontWeight: W.medium }}>{item.meaning}</div>
              {item.example && (
                <div style={{ fontSize: T.sm, color: C.text3 }}>Example: <em>"{item.example}"</em></div>
              )}
            </>
          )}
        </Card>

        {!revealed ? (
          <Button fullWidth style={{ marginTop: S[4] }} onClick={() => setRevealed(true)}>
            Reveal
          </Button>
        ) : (
          <div style={{ display: "flex", gap: S[3], marginTop: S[4] }}>
            <Button variant="secondary" fullWidth onClick={handleLearning} style={{ borderColor: C.danger, color: C.danger }}>
              ✗ Still learning
            </Button>
            <Button fullWidth onClick={handleKnow} style={{ background: "#16A34A" }}>
              ✓ Know it
            </Button>
          </div>
        )}
      </div>
    </div>
  );
}
