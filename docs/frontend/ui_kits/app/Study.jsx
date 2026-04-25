// Jeeb UI Kit — Study Page v2

const { useState, useEffect } = React;

const SESSIONS = [
  { id: 1, subject: 'Mathematics', topic: 'Calculus Chapter 5', duration: '2h 30m', date: 'Apr 22', notes: 'Completed 3 practice tests' },
  { id: 2, subject: 'Physics',     topic: 'Thermodynamics',    duration: '1h 45m', date: 'Apr 21', notes: '' },
  { id: 3, subject: 'English',     topic: 'Essay writing',     duration: '1h 0m',  date: 'Apr 20', notes: 'Drafted intro and body' },
  { id: 4, subject: 'Mathematics', topic: 'Integration',       duration: '3h 0m',  date: 'Apr 18', notes: '' },
];
const SUBJECTS = [
  { name: 'Mathematics', hours: 30, goal: 40, color: '#2563EB' },
  { name: 'Physics',     hours: 12, goal: 30, color: '#16A34A' },
  { name: 'English',     hours: 8,  goal: 20, color: '#D97706' },
];

function StudyTimer({ subject, onStop }) {
  const t = window.useTheme();
  const [secs, setSecs] = useState(0);
  const [paused, setPaused] = useState(false);
  useEffect(() => {
    const id = setInterval(() => { if (!paused) setSecs(s => s + 1); }, 1000);
    return () => clearInterval(id);
  }, [paused]);
  const fmt = s => {
    const h = String(Math.floor(s / 3600)).padStart(2,'0');
    const m = String(Math.floor((s % 3600) / 60)).padStart(2,'0');
    const sec = String(s % 60).padStart(2,'0');
    return `${h}:${m}:${sec}`;
  };
  return (
    <div style={{ background: t.surface, border: `1px solid ${t.border}`, borderRadius: 8, padding: '28px 20px', textAlign: 'center', marginBottom: 16 }}>
      <div style={{ fontSize: 12, color: t.fg2, marginBottom: 4 }}>Studying: <strong style={{ color: t.fg1 }}>{subject}</strong></div>
      <div style={{ fontFamily: "'JetBrains Mono', monospace", fontSize: 48, fontWeight: 500, color: t.fg1, letterSpacing: '-0.02em', margin: '12px 0' }}>{fmt(secs)}</div>
      <div style={{ display: 'flex', gap: 10, justifyContent: 'center' }}>
        <button onClick={() => setPaused(p => !p)} style={{ padding: '8px 20px', border: `1px solid ${t.border}`, borderRadius: 8, background: t.surfaceActive, fontSize: 14, fontWeight: 500, cursor: 'pointer', color: t.fg2 }}>
          {paused ? 'Resume' : 'Pause'}
        </button>
        <button onClick={() => onStop(Math.floor(secs / 60))} style={{ padding: '8px 20px', border: 'none', borderRadius: 8, background: '#DC2626', fontSize: 14, fontWeight: 500, cursor: 'pointer', color: '#fff' }}>
          Stop & Save
        </button>
      </div>
      <div style={{ fontSize: 11, color: t.fg3, marginTop: 10 }}>
        Started at {new Date(Date.now() - secs * 1000).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
      </div>
    </div>
  );
}

function SubjectProgress({ s }) {
  const t = window.useTheme();
  const pct = Math.min((s.hours / s.goal) * 100, 100);
  return (
    <div style={{ marginBottom: 12 }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 4 }}>
        <span style={{ fontSize: 13, fontWeight: 500, color: t.fg1 }}>{s.name}</span>
        <span style={{ fontSize: 12, color: t.fg2 }}>{s.hours}h / {s.goal}h</span>
      </div>
      <div style={{ height: 8, background: t.surfaceActive, borderRadius: 9999 }}>
        <div style={{ height: '100%', borderRadius: 9999, background: s.color, width: `${pct}%` }} />
      </div>
    </div>
  );
}

function SessionCard({ s }) {
  const t = window.useTheme();
  return (
    <div style={{ background: t.surface, border: `1px solid ${t.border}`, borderRadius: 8, padding: '14px 16px', display: 'flex', gap: 12, alignItems: 'flex-start' }}>
      <div style={{ fontSize: 20, marginTop: 2 }}>📚</div>
      <div style={{ flex: 1 }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <span style={{ fontSize: 14, fontWeight: 600, color: t.fg1 }}>{s.subject}</span>
          <span style={{ fontSize: 13, fontWeight: 600, color: '#2563EB' }}>{s.duration}</span>
        </div>
        <div style={{ fontSize: 12, color: t.fg2, margin: '2px 0' }}>{s.topic}</div>
        {s.notes && <div style={{ fontSize: 12, color: t.fg3 }}>{s.notes}</div>}
        <div style={{ fontSize: 11, color: t.fg3, marginTop: 4 }}>{s.date}</div>
      </div>
    </div>
  );
}

function StartTimerModal({ onStart, onClose }) {
  const t = window.useTheme();
  const [subject, setSubject] = useState('Mathematics');
  return (
    <div style={{ position: 'fixed', inset: 0, background: 'rgb(0 0 0/.5)', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 200 }}>
      <div style={{ background: t.surface, borderRadius: 12, width: 360, boxShadow: t.modalShadow }}>
        <div style={{ padding: '16px 20px', borderBottom: `1px solid ${t.border}`, display: 'flex', justifyContent: 'space-between' }}>
          <h2 style={{ fontSize: 16, fontWeight: 600, color: t.fg1 }}>Start Study Session</h2>
          <button onClick={onClose} style={{ background: 'none', border: 'none', cursor: 'pointer', color: t.fg3, fontSize: 20 }}>×</button>
        </div>
        <div style={{ padding: 20 }}>
          <label style={{ fontSize: 13, fontWeight: 500, color: t.fg1, display: 'block', marginBottom: 6 }}>Subject</label>
          <select value={subject} onChange={e => setSubject(e.target.value)} style={{ width: '100%', border: `1px solid ${t.border}`, borderRadius: 8, padding: '8px 12px', fontSize: 14, fontFamily: 'inherit', outline: 'none', background: t.inputBg, color: t.fg1 }}>
            <option>Mathematics</option><option>Physics</option><option>English</option>
          </select>
        </div>
        <div style={{ padding: '14px 20px', borderTop: `1px solid ${t.border}`, display: 'flex', justifyContent: 'flex-end', gap: 8 }}>
          <button onClick={onClose} style={{ padding: '8px 16px', border: `1px solid ${t.border}`, borderRadius: 8, background: t.surface, fontSize: 14, cursor: 'pointer', color: t.fg2 }}>Cancel</button>
          <button onClick={() => onStart(subject)} style={{ padding: '8px 18px', border: 'none', borderRadius: 8, background: '#2563EB', color: '#fff', fontSize: 14, fontWeight: 500, cursor: 'pointer' }}>Start Timer</button>
        </div>
      </div>
    </div>
  );
}

function Study() {
  const t = window.useTheme();
  const [sessions, setSessions] = useState(SESSIONS);
  const [activeTimer, setActiveTimer] = useState(null);
  const [showModal, setShowModal] = useState(false);
  const width = window.useWindowWidth();
  const isMobile = width < 768;

  return (
    <div style={{ maxWidth: 760 }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
        <h1 style={{ fontSize: 22, fontWeight: 700, color: t.fg1 }}>Study</h1>
        <div style={{ display: 'flex', gap: 8 }}>
          {!activeTimer && <button onClick={() => setShowModal(true)} style={{ padding: '8px 14px', border: '1px solid #2563EB', borderRadius: 8, background: t.primarySubtle, fontSize: 13, fontWeight: 500, cursor: 'pointer', color: '#2563EB' }}>▶ Start Timer</button>}
          <button style={{ padding: '8px 16px', background: '#2563EB', color: '#fff', border: 'none', borderRadius: 8, fontSize: 14, fontWeight: 500, cursor: 'pointer' }}>+ Log Session</button>
        </div>
      </div>
      {activeTimer && (
        <StudyTimer subject={activeTimer} onStop={mins => {
          setSessions(ss => [{ id: Date.now(), subject: activeTimer, topic: 'Session', duration: `${Math.floor(mins/60)}h ${mins%60}m`, date: 'Today', notes: '' }, ...ss]);
          setActiveTimer(null);
        }} />
      )}
      <div style={{ display: 'grid', gridTemplateColumns: isMobile ? '1fr' : '1fr 1fr', gap: 16 }}>
        <div style={{ background: t.surface, border: `1px solid ${t.border}`, borderRadius: 8, padding: '16px 20px' }}>
          <div style={{ fontSize: 14, fontWeight: 600, color: t.fg1, marginBottom: 14 }}>Subject Progress</div>
          {SUBJECTS.map(s => <SubjectProgress key={s.name} s={s} />)}
        </div>
        <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
          {sessions.slice(0, 4).map(s => <SessionCard key={s.id} s={s} />)}
        </div>
      </div>
      {showModal && <StartTimerModal onStart={s => { setActiveTimer(s); setShowModal(false); }} onClose={() => setShowModal(false)} />}
    </div>
  );
}

Object.assign(window, { Study });
