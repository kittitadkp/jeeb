import { useState } from "react";
import { useNavigate, useParams } from "react-router-dom";

import { FlashcardTool } from "@/components/study/FlashcardTool";
import { RecallTool } from "@/components/study/RecallTool";
import { Badge, StatusDot } from "@/components/ui/Badge";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { StatCard } from "@/components/ui/StatCard";
import { ErrorState, LoadingState } from "@/components/ui/States";
import { useAllItems } from "@/hooks/useItems";
import { useResetProgress, useTopicProgress, useUpsertProgress } from "@/hooks/useProgress";
import { useTopic } from "@/hooks/useTopics";
import { C, R, S, T, W } from "@/lib/design";

type Tab = "browse" | "flashcard" | "recall" | "progress";

const TABS: { key: Tab; label: string }[] = [
  { key: "browse", label: "Browse" },
  { key: "flashcard", label: "Flashcard" },
  { key: "recall", label: "Recall" },
  { key: "progress", label: "Progress" },
];

export function Topic() {
  const { id = "" } = useParams();
  const nav = useNavigate();
  const [tab, setTab] = useState<Tab>("browse");
  const [search, setSearch] = useState("");
  const [filterCategory, setFilterCategory] = useState("");

  const { data: topic, isLoading: topicLoading, error: topicError } = useTopic(id);
  const { data: items, isLoading: itemsLoading } = useAllItems(id);
  const { data: progressMap = {} } = useTopicProgress(id);
  const upsert = useUpsertProgress();
  const reset = useResetProgress(id);

  if (topicLoading || itemsLoading) return <LoadingState />;
  if (topicError || !topic) return <ErrorState message="Topic not found" />;

  const allItems = items ?? [];
  const categories = [...new Set(allItems.map((i) => i.category))];

  const mastered = allItems.filter((i) => progressMap[i.id] === "mastered").length;
  const learning = allItems.filter((i) => progressMap[i.id] === "learning").length;
  const notStarted = allItems.length - mastered - learning;
  const pct = allItems.length > 0 ? (mastered / allItems.length) * 100 : 0;

  const filtered = allItems.filter((item) => {
    const matchSearch = !search || item.term.toLowerCase().includes(search.toLowerCase()) || item.meaning.toLowerCase().includes(search.toLowerCase());
    const matchCat = !filterCategory || item.category === filterCategory;
    return matchSearch && matchCat;
  });

  function handleUpsert(itemId: string, status: "learning" | "mastered") {
    upsert.mutate({ itemId, topicId: id, status });
  }

  return (
    <div style={{ padding: `${S[6]}px ${S[8]}px`, maxWidth: 900, margin: "0 auto" }}>
      <button onClick={() => nav("/")} style={{ background: "none", border: "none", color: C.text2, cursor: "pointer", fontSize: T.sm, marginBottom: S[4], padding: 0 }}>
        ← Back
      </button>

      <div style={{ display: "flex", alignItems: "center", gap: S[3], marginBottom: S[6] }}>
        <span style={{ fontSize: 36 }}>{topic.icon}</span>
        <div>
          <h1 style={{ margin: 0, fontSize: T["2xl"], fontWeight: W.bold, color: C.text }}>{topic.name}</h1>
          <div style={{ fontSize: T.sm, color: C.text2 }}>{topic.description}</div>
        </div>
      </div>

      <div style={{ display: "flex", gap: 2, marginBottom: S[5], borderBottom: `1px solid ${C.border}`, paddingBottom: 0 }}>
        {TABS.map((t) => (
          <button
            key={t.key}
            onClick={() => setTab(t.key)}
            style={{
              background: "none",
              border: "none",
              borderBottom: tab === t.key ? `2px solid ${C.primary}` : "2px solid transparent",
              padding: `${S[2]}px ${S[4]}px`,
              fontSize: T.base,
              fontWeight: tab === t.key ? W.semi : W.normal,
              color: tab === t.key ? C.primary : C.text2,
              cursor: "pointer",
              marginBottom: -1,
              transition: "color .15s",
            }}
          >
            {t.label}
          </button>
        ))}
      </div>

      {tab === "browse" && (
        <div>
          <div style={{ display: "flex", gap: S[3], marginBottom: S[4], flexWrap: "wrap" }}>
            <input
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              placeholder="Search term or meaning…"
              style={{
                flex: "1 1 200px",
                padding: "7px 12px",
                fontSize: T.base,
                background: C.surface2,
                border: `1px solid ${C.border}`,
                borderRadius: R.md,
                color: C.text,
                outline: "none",
              }}
            />
            <div style={{ display: "flex", gap: S[2], flexWrap: "wrap" }}>
              <button
                onClick={() => setFilterCategory("")}
                style={{
                  padding: "4px 12px",
                  borderRadius: R.full,
                  border: `1px solid ${!filterCategory ? C.primary : C.border}`,
                  background: !filterCategory ? `${C.primary}20` : "transparent",
                  color: !filterCategory ? C.primary : C.text2,
                  fontSize: T.sm,
                  cursor: "pointer",
                }}
              >
                All
              </button>
              {categories.map((cat) => (
                <button
                  key={cat}
                  onClick={() => setFilterCategory(filterCategory === cat ? "" : cat)}
                  style={{
                    padding: "4px 12px",
                    borderRadius: R.full,
                    border: `1px solid ${filterCategory === cat ? C.primary : C.border}`,
                    background: filterCategory === cat ? `${C.primary}20` : "transparent",
                    color: filterCategory === cat ? C.primary : C.text2,
                    fontSize: T.sm,
                    cursor: "pointer",
                  }}
                >
                  {cat}
                </button>
              ))}
            </div>
          </div>

          <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fill, minmax(200px, 1fr))", gap: S[3] }}>
            {filtered.map((item) => (
              <Card key={item.id} style={{ display: "flex", flexDirection: "column", gap: S[2] }}>
                <div style={{ display: "flex", justifyContent: "space-between", alignItems: "flex-start" }}>
                  <span style={{ fontSize: T["2xl"], fontWeight: W.bold, fontFamily: "monospace", color: C.text }}>{item.term}</span>
                  <StatusDot status={progressMap[item.id]} />
                </div>
                <div style={{ fontSize: T.sm, color: C.text2 }}>{item.meaning}</div>
                {item.example && <div style={{ fontSize: T.xs, color: C.text3 }}>"{item.example}"</div>}
                <Badge>{item.category}</Badge>
              </Card>
            ))}
          </div>
        </div>
      )}

      {tab === "flashcard" && (
        <FlashcardTool
          items={allItems}
          topicId={id}
          progressMap={progressMap}
          onUpsert={handleUpsert}
        />
      )}

      {tab === "recall" && (
        <RecallTool
          items={allItems}
          topicId={id}
          progressMap={progressMap}
          onUpsert={handleUpsert}
        />
      )}

      {tab === "progress" && (
        <div style={{ display: "flex", flexDirection: "column", gap: S[6] }}>
          <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fill, minmax(150px, 1fr))", gap: S[3] }}>
            <StatCard label="Mastered" value={mastered} color="#16A34A" />
            <StatCard label="Learning" value={learning} color="#D97706" />
            <StatCard label="Not started" value={notStarted} color={C.text3} />
            <StatCard label="Total" value={allItems.length} />
          </div>

          <div>
            <div style={{ display: "flex", justifyContent: "space-between", fontSize: T.sm, color: C.text2, marginBottom: S[2] }}>
              <span>Progress</span>
              <span>{Math.round(pct)}%</span>
            </div>
            <div style={{ height: 8, background: C.surface3, borderRadius: R.full }}>
              <div style={{ height: "100%", width: `${pct}%`, background: "#16A34A", borderRadius: R.full, transition: "width .4s" }} />
            </div>
          </div>

          {mastered > 0 && (
            <div>
              <div style={{ fontSize: T.sm, color: C.text2, marginBottom: S[3], fontWeight: W.medium }}>Mastered items</div>
              <div style={{ display: "flex", gap: S[2], flexWrap: "wrap" }}>
                {allItems.filter((i) => progressMap[i.id] === "mastered").map((i) => (
                  <Badge key={i.id} color="#16A34A" style={{ fontFamily: "monospace", fontSize: T.sm }}>
                    {i.term}
                  </Badge>
                ))}
              </div>
            </div>
          )}

          <div>
            <Button
              variant="secondary"
              style={{ borderColor: C.danger, color: C.danger }}
              onClick={() => {
                if (confirm("Reset all progress for this topic?")) {
                  reset.mutate();
                }
              }}
            >
              Reset Progress
            </Button>
          </div>
        </div>
      )}
    </div>
  );
}
