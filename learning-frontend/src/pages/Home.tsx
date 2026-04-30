import { useNavigate } from "react-router-dom";

import { SectionLabel } from "@/components/ui/SectionLabel";
import { StatCard } from "@/components/ui/StatCard";
import { ErrorState, LoadingState } from "@/components/ui/States";
import { useTopics } from "@/hooks/useTopics";
import { useStats } from "@/hooks/useProgress";
import { C, R, S, T, W } from "@/lib/design";

export function Home() {
  const { data: topics, isLoading, error } = useTopics();
  const { data: stats } = useStats();
  const nav = useNavigate();

  const totalMastered = stats?.reduce((sum, s) => sum + s.mastered, 0) ?? 0;
  const activeTopics = stats?.filter((s) => s.learning > 0 || s.mastered > 0).length ?? 0;

  if (isLoading) return <LoadingState />;
  if (error) return <ErrorState message={String(error)} />;

  return (
    <div style={{ padding: `${S[6]}px ${S[8]}px`, maxWidth: 900, margin: "0 auto" }}>
      <div style={{ marginBottom: S[6] }}>
        <SectionLabel>🎓 Learning</SectionLabel>
        <h1 style={{ margin: `${S[2]}px 0 0`, fontSize: T["2xl"], fontWeight: W.bold, color: C.text }}>
          Your Topics
        </h1>
      </div>

      <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fill, minmax(160px, 1fr))", gap: S[3], marginBottom: S[8] }}>
        <StatCard label="Total mastered" value={totalMastered} color="#16A34A" />
        <StatCard label="Active topics" value={activeTopics} color={C.primary} />
        <StatCard label="Topics" value={topics?.length ?? 0} />
      </div>

      <div style={{ marginBottom: S[4] }}>
        <SectionLabel>Topics</SectionLabel>
      </div>

      <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fill, minmax(260px, 1fr))", gap: S[4] }}>
        {topics?.map((topic) => {
          const s = stats?.find((x) => x.topic_id === topic.id);
          const mastered = s?.mastered ?? 0;
          const total = s?.total ?? 0;
          const pct = total > 0 ? (mastered / total) * 100 : 0;

          return (
            <div
              key={topic.id}
              onClick={() => nav(`/topics/${topic.id}`)}
              style={{
                background: C.surface,
                border: `1px solid ${C.border}`,
                borderRadius: R.card,
                padding: S[5],
                cursor: "pointer",
                transition: "border-color .15s",
                display: "flex",
                flexDirection: "column",
                gap: S[3],
              }}
              onMouseEnter={(e) => (e.currentTarget.style.borderColor = C.primary)}
              onMouseLeave={(e) => (e.currentTarget.style.borderColor = C.border)}
            >
              <div style={{ display: "flex", alignItems: "center", gap: S[2] }}>
                <span style={{ fontSize: 28 }}>{topic.icon}</span>
                <div>
                  <div style={{ fontSize: T.md, fontWeight: W.semi, color: C.text }}>{topic.name}</div>
                  <div style={{ fontSize: T.xs, color: C.text3 }}>{topic.category}</div>
                </div>
              </div>

              <div style={{ fontSize: T.sm, color: C.text2, lineHeight: 1.5 }}>{topic.description}</div>

              {total > 0 && (
                <div>
                  <div style={{ display: "flex", justifyContent: "space-between", fontSize: T.xs, color: C.text3, marginBottom: 4 }}>
                    <span>{mastered} / {total} mastered</span>
                    <span>{Math.round(pct)}%</span>
                  </div>
                  <div style={{ height: 4, background: C.surface3, borderRadius: R.full }}>
                    <div style={{ height: "100%", width: `${pct}%`, background: "#16A34A", borderRadius: R.full, transition: "width .3s" }} />
                  </div>
                </div>
              )}

              <div style={{ fontSize: T.sm, color: C.primary, fontWeight: W.medium, marginTop: "auto" }}>
                Study →
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
