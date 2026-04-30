import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Route, BrowserRouter as Router, Routes } from "react-router-dom";

import { AuthProvider } from "@/lib/auth";
import { Home } from "@/pages/Home";
import { Topic } from "@/pages/Topic";
import { useDarkMode } from "@/store/theme";

const qc = new QueryClient({
  defaultOptions: { queries: { retry: 1, staleTime: 30_000 } },
});

function Shell() {
  const { dark, toggleDark } = useDarkMode();

  return (
    <div style={{ minHeight: "100vh", background: "var(--c-bg)", color: "var(--c-text)" }}>
      <header style={{
        position: "sticky",
        top: 0,
        zIndex: 10,
        background: "var(--c-surface)",
        borderBottom: "1px solid var(--c-border)",
        padding: "0 32px",
        height: 52,
        display: "flex",
        alignItems: "center",
        justifyContent: "space-between",
      }}>
        <div style={{ display: "flex", alignItems: "center", gap: 8, fontWeight: 600, fontSize: 15, color: "#0ea5e9" }}>
          🎓 Jeeb Learning
        </div>
        <button
          onClick={toggleDark}
          style={{ background: "none", border: "none", cursor: "pointer", color: "var(--c-text2)", fontSize: 18 }}
          title="Toggle dark mode"
        >
          {dark ? "☀️" : "🌙"}
        </button>
      </header>

      <main style={{ paddingTop: 0 }}>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/topics/:id" element={<Topic />} />
        </Routes>
      </main>
    </div>
  );
}

export default function App() {
  return (
    <QueryClientProvider client={qc}>
      <Router>
        <AuthProvider>
          <Shell />
        </AuthProvider>
      </Router>
    </QueryClientProvider>
  );
}
