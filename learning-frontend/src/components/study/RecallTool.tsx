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

type CheckState = "idle" | "correct" | "incorrect";

export function RecallTool({ items, progressMap, onUpsert }: Props) {
  const queue = useMemo(() => {
    const unmastered = items.filter((i) => progressMap[i.id] !== "mastered");
    const pool = unmastered.length > 0 ? unmastered : items;
    return [...pool].sort(() => Math.random() - 0.5);
  }, [items, progressMap]);

  const [idx, setIdx] = useState(0);
  const [input, setInput] = useState("");
  const [check, setCheck] = useState<CheckState>("idle");
  const [done, setDone] = useState(false);
  const [correct, setCorrect] = useState(0);
  const inputRef = useRef<HTMLInputElement>(null);
  const startTime = useRef(Date.now());
  const [elapsed, setElapsed] = useState(0);

  const item = queue[idx];
  const total = queue.length;

  function normalize(s: string) {
    return s.trim().toLowerCase().replace(/\s+/g, " ");
  }

  function handleSubmit() {
    if (!input.trim() || check !== "idle") return;
    const isCorrect = normalize(input) === normalize(item.term);
    onUpsert(item.id, isCorrect ? "mastered" : "learning");
    if (isCorrect) setCorrect((c) => c + 1);
    setCheck(isCorrect ? "correct" : "incorrect");
  }

  function formatTime(s: number) {
    const m = Math.floor(s / 60);
    const sec = s % 60;
    return m > 0 ? `${m}m ${sec}s` : `${sec}s`;
  }

  function handleNext() {
    setInput("");
    setCheck("idle");
    if (idx + 1 >= total) {
      setElapsed(Math.round((Date.now() - startTime.current) / 1000));
      setDone(true);
    } else {
      setIdx((i) => i + 1);
      setTimeout(() => inputRef.current?.focus(), 50);
    }
  }

  if (done) {
    return (
      <div style={{ display: "flex", flexDirection: "column", alignItems: "center", gap: S[6], padding: S[8] }}>
        <div style={{ fontSize: T["2xl"], fontWeight: W.bold, color: C.text }}>Session complete!</div>
        <div style={{ display: "flex", gap: S[6] }}>
          <StatCard label="Correct" value={`${correct} / ${total}`} color="#16A34A" />
          <StatCard label="Time taken" value={formatTime(elapsed)} />
        </div>
        <Button onClick={() => { setIdx(0); setInput(""); setCheck("idle"); setDone(false); setCorrect(0); startTime.current = Date.now(); }}>
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

        <Card style={{ textAlign: "center", gap: S[2], display: "flex", flexDirection: "column" }}>
          <div style={{ fontSize: T.md, color: C.text, fontWeight: W.medium }}>{item.meaning}</div>
          {item.example && (
            <div style={{ fontSize: T.sm, color: C.text2 }}>Example: <em>"{item.example}"</em></div>
          )}
        </Card>

        <div style={{ marginTop: S[4], display: "flex", flexDirection: "column", gap: S[3] }}>
          <label style={{ fontSize: T.sm, color: C.text2 }}>Type the IPA symbol:</label>
          <div style={{ display: "flex", gap: S[2] }}>
            <input
              ref={inputRef}
              value={input}
              onChange={(e) => setInput(e.target.value)}
              onKeyDown={(e) => e.key === "Enter" && (check === "idle" ? handleSubmit() : handleNext())}
              disabled={check !== "idle"}
              placeholder="e.g. /p/"
              style={{
                flex: 1,
                padding: "8px 12px",
                fontSize: T.md,
                fontFamily: "monospace",
                background: C.surface2,
                border: `1px solid ${check === "correct" ? "#16A34A" : check === "incorrect" ? C.danger : C.border}`,
                borderRadius: R.md,
                color: C.text,
                outline: "none",
              }}
            />
            {check === "idle" && (
              <Button onClick={handleSubmit}>Submit</Button>
            )}
          </div>

          {check !== "idle" && (
            <div style={{ display: "flex", alignItems: "center", justifyContent: "space-between" }}>
              <div style={{ fontSize: T.sm, fontWeight: W.medium, color: check === "correct" ? "#16A34A" : C.danger }}>
                {check === "correct"
                  ? `✓ Correct!`
                  : `✗ Incorrect — answer was ${item.term}`}
              </div>
              <Button size="sm" variant="secondary" onClick={handleNext}>
                Next →
              </Button>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
