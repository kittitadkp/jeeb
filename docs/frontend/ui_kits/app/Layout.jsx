// Jeeb UI Kit — Layout Components v2
// ThemeContext, Header, Sidebar, BottomNav, AppLayout — exported to window

const { useState, useContext, createContext, useEffect } = React;

// ── Theme ───────────────────────────────────────────────────
const LIGHT = {
  bg: '#F8FAFC', surface: '#ffffff', surfaceHover: '#F8FAFC',
  surfaceActive: '#F1F5F9', border: '#E2E8F0', borderStrong: '#CBD5E1',
  fg1: '#0F172A', fg2: '#475569', fg3: '#94A3B8',
  primarySubtle: '#EFF6FF', primaryMuted: '#DBEAFE',
  navActive: '#EFF6FF', navActiveText: '#2563EB',
  inputBg: '#F8FAFC', shadow: '0 1px 2px 0 rgb(0 0 0/.05)',
  modalShadow: '0 20px 40px rgb(0 0 0/.15)',
  successSubtle: '#F0FDF4', warningSubtle: '#FFFBEB', dangerSubtle: '#FEF2F2',
  dangerBorder: '#FEE2E2', dangerText: '#DC2626',
  successText: '#16A34A', warningText: '#D97706',
  dark: false,
};
const DARK = {
  bg: '#020617', surface: '#0F172A', surfaceHover: '#1E293B',
  surfaceActive: '#1E293B', border: '#1E293B', borderStrong: '#334155',
  fg1: '#F1F5F9', fg2: '#94A3B8', fg3: '#475569',
  primarySubtle: '#172554', primaryMuted: '#1E3A8A',
  navActive: '#172554', navActiveText: '#60A5FA',
  inputBg: '#1E293B', shadow: '0 1px 2px 0 rgb(0 0 0/.3)',
  modalShadow: '0 20px 40px rgb(0 0 0/.5)',
  successSubtle: '#052e16', warningSubtle: '#431407', dangerSubtle: '#450a0a',
  dangerBorder: '#7f1d1d', dangerText: '#f87171',
  successText: '#4ade80', warningText: '#fbbf24',
  dark: true,
};

const ThemeCtx = createContext(LIGHT);
function useTheme() { return useContext(ThemeCtx); }

// ── Shared helpers ──────────────────────────────────────────
function useWindowWidth() {
  const [w, setW] = useState(window.innerWidth);
  useEffect(() => {
    const fn = () => setW(window.innerWidth);
    window.addEventListener('resize', fn);
    return () => window.removeEventListener('resize', fn);
  }, []);
  return w;
}

// ── Nav items ───────────────────────────────────────────────
const NAV_ITEMS = [
  { id: 'dashboard', label: 'Dashboard', emoji: '📊' },
  { id: 'workouts',  label: 'Workouts',  emoji: '💪' },
  { id: 'study',     label: 'Study',     emoji: '📚' },
  { id: 'sleep',     label: 'Sleep',     emoji: '😴' },
  { id: 'finance',   label: 'Finance',   emoji: '💰' },
  { id: 'calendar',  label: 'Calendar',  emoji: '📅' },
];

// ── Header ──────────────────────────────────────────────────
function Header({ onNav, darkMode, toggleDark }) {
  const t = useTheme();
  const [showMenu, setShowMenu] = useState(false);

  return (
    <header style={{
      height: 56, background: t.surface, borderBottom: `1px solid ${t.border}`,
      display: 'flex', alignItems: 'center', padding: '0 16px', gap: 12,
      position: 'sticky', top: 0, zIndex: 50, flexShrink: 0,
    }}>
      <div style={{ fontWeight: 700, fontSize: 18, color: '#2563EB', letterSpacing: '-0.02em', marginRight: 4 }}>Jeeb</div>

      <div style={{
        display: 'flex', alignItems: 'center', gap: 6, background: t.inputBg,
        border: `1px solid ${t.border}`, borderRadius: 8, padding: '6px 12px',
        flex: 1, maxWidth: 260, color: t.fg3, fontSize: 13, cursor: 'text',
      }}>
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
        <span>Search…</span>
        <span style={{ marginLeft: 'auto', background: t.surfaceActive, borderRadius: 4, padding: '1px 5px', fontSize: 11, color: t.fg2 }}>⌘K</span>
      </div>

      <div style={{ marginLeft: 'auto', display: 'flex', alignItems: 'center', gap: 6 }}>
        {/* Dark mode toggle */}
        <button onClick={toggleDark} title={darkMode ? 'Light mode' : 'Dark mode'} style={{
          width: 34, height: 34, borderRadius: 8, background: t.surfaceActive,
          border: 'none', display: 'flex', alignItems: 'center', justifyContent: 'center',
          cursor: 'pointer', color: t.fg2, fontSize: 16,
        }}>{darkMode ? '☀️' : '🌙'}</button>

        {/* Bell */}
        <div style={{
          width: 34, height: 34, borderRadius: 8, background: t.surfaceActive,
          display: 'flex', alignItems: 'center', justifyContent: 'center',
          cursor: 'pointer', color: t.fg2, position: 'relative',
        }}>
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"/><path d="M13.73 21a2 2 0 0 1-3.46 0"/></svg>
          <span style={{ position: 'absolute', top: 6, right: 6, width: 7, height: 7, background: '#DC2626', borderRadius: '50%', border: `2px solid ${t.surfaceActive}` }} />
        </div>

        {/* Avatar */}
        <div style={{ position: 'relative' }}>
          <div onClick={() => setShowMenu(!showMenu)} style={{
            display: 'flex', alignItems: 'center', gap: 8, cursor: 'pointer',
            padding: '4px 8px', borderRadius: 8,
          }}>
            <div style={{ width: 30, height: 30, borderRadius: '50%', background: '#2563EB', color: '#fff', fontSize: 13, fontWeight: 600, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>J</div>
            <span style={{ fontSize: 13, fontWeight: 500, color: t.fg1 }}>John</span>
          </div>
          {showMenu && (
            <div style={{
              position: 'absolute', right: 0, top: '100%', marginTop: 4,
              background: t.surface, border: `1px solid ${t.border}`, borderRadius: 8,
              boxShadow: t.modalShadow, width: 150, zIndex: 100, padding: '4px 0',
            }}>
              {['Profile', 'Settings', '—', 'Log out'].map(label =>
                label === '—'
                  ? <div key={label} style={{ height: 1, background: t.border, margin: '4px 0' }} />
                  : <div key={label} onClick={() => { setShowMenu(false); label === 'Settings' && onNav('settings'); }} style={{ padding: '8px 14px', fontSize: 13, color: t.fg1, cursor: 'pointer' }}
                      onMouseEnter={e => e.currentTarget.style.background = t.surfaceHover}
                      onMouseLeave={e => e.currentTarget.style.background = 'transparent'}>{label}</div>
              )}
            </div>
          )}
        </div>
      </div>
    </header>
  );
}

// ── Sidebar ─────────────────────────────────────────────────
function Sidebar({ active, onNav }) {
  const t = useTheme();
  return (
    <aside style={{
      width: 220, flexShrink: 0, background: t.surface,
      borderRight: `1px solid ${t.border}`, padding: '12px 8px',
      display: 'flex', flexDirection: 'column', gap: 2, overflowY: 'auto',
    }}>
      {NAV_ITEMS.map(item => {
        const isActive = active === item.id;
        return (
          <div key={item.id} onClick={() => onNav(item.id)} style={{
            display: 'flex', alignItems: 'center', gap: 10, padding: '8px 10px',
            borderRadius: 6, cursor: 'pointer',
            background: isActive ? t.navActive : 'transparent',
            color: isActive ? t.navActiveText : t.fg2,
            fontWeight: isActive ? 600 : 500, fontSize: 14,
          }}
          onMouseEnter={e => { if (!isActive) e.currentTarget.style.background = t.surfaceActive; }}
          onMouseLeave={e => { if (!isActive) e.currentTarget.style.background = 'transparent'; }}>
            <span style={{ fontSize: 16 }}>{item.emoji}</span>
            <span>{item.label}</span>
          </div>
        );
      })}
      <div style={{ height: 1, background: t.border, margin: '6px 4px' }} />
      <div onClick={() => onNav('settings')} style={{
        display: 'flex', alignItems: 'center', gap: 10, padding: '8px 10px',
        borderRadius: 6, cursor: 'pointer',
        background: active === 'settings' ? t.navActive : 'transparent',
        color: active === 'settings' ? t.navActiveText : t.fg2,
        fontWeight: 500, fontSize: 14,
      }}
      onMouseEnter={e => { if (active !== 'settings') e.currentTarget.style.background = t.surfaceActive; }}
      onMouseLeave={e => { if (active !== 'settings') e.currentTarget.style.background = 'transparent'; }}>
        <span style={{ fontSize: 16 }}>⚙️</span><span>Settings</span>
      </div>
    </aside>
  );
}

// ── Bottom Nav (mobile) ─────────────────────────────────────
function BottomNav({ active, onNav }) {
  const t = useTheme();
  const items = [...NAV_ITEMS.slice(0, 5), { id: 'settings', label: 'Settings', emoji: '⚙️' }];
  return (
    <nav style={{
      position: 'fixed', bottom: 0, left: 0, right: 0,
      background: t.surface, borderTop: `1px solid ${t.border}`,
      display: 'flex', zIndex: 50, paddingBottom: 'env(safe-area-inset-bottom)',
    }}>
      {items.map(item => {
        const isActive = active === item.id;
        return (
          <div key={item.id} onClick={() => onNav(item.id)} style={{
            flex: 1, display: 'flex', flexDirection: 'column', alignItems: 'center',
            padding: '8px 0 6px', cursor: 'pointer', gap: 2,
            color: isActive ? '#2563EB' : t.fg3,
          }}>
            <span style={{ fontSize: 20 }}>{item.emoji}</span>
            <span style={{ fontSize: 9, fontWeight: isActive ? 600 : 400 }}>{item.label}</span>
          </div>
        );
      })}
    </nav>
  );
}

// ── AppLayout ────────────────────────────────────────────────
function AppLayout({ page, onNav, children, darkMode, toggleDark }) {
  const t = useTheme();
  const width = useWindowWidth();
  const isMobile = width < 768;

  return (
    <div style={{ display: 'flex', flexDirection: 'column', height: '100vh', fontFamily: "'Inter', sans-serif", background: t.bg }}>
      <Header onNav={onNav} darkMode={darkMode} toggleDark={toggleDark} />
      <div style={{ display: 'flex', flex: 1, overflow: 'hidden' }}>
        {!isMobile && <Sidebar active={page} onNav={onNav} />}
        <main style={{ flex: 1, overflowY: 'auto', padding: isMobile ? '16px 12px 80px' : 24 }}>
          {children}
        </main>
      </div>
      {isMobile && <BottomNav active={page} onNav={onNav} />}
    </div>
  );
}

Object.assign(window, { AppLayout, Sidebar, Header, BottomNav, NAV_ITEMS, ThemeCtx, useTheme, LIGHT, DARK, useWindowWidth });
