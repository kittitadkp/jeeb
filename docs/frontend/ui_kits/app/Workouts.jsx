// Jeeb UI Kit — Workouts Page v3
// Includes: Set/Rest tracker with tap-to-log sets

const { useState, useEffect, useRef } = React;

const WORKOUT_DATA = [
  { id: 1, type: 'Strength', name: 'Upper Body', duration: 45, exercises: ['Bench Press 3×10', 'Pull Ups 3×8', 'Shoulder Press 3×12'], date: 'Apr 22, 2024' },
  { id: 2, type: 'Cardio',   name: 'Morning Run', duration: 35, exercises: ['5km run'], date: 'Apr 21, 2024' },
  { id: 3, type: 'Strength', name: 'Lower Body', duration: 50, exercises: ['Squats 4×8', 'Deadlifts 3×6', 'Lunges 3×12'], date: 'Apr 20, 2024' },
  { id: 4, type: 'Flexibility', name: 'Yoga Flow', duration: 30, exercises: ['Sun salutation', 'Warrior sequence'], date: 'Apr 18, 2024' },
];

const TYPE_COLORS = {
  Strength:    { bg: '#EFF6FF', color: '#2563EB', darkBg: '#172554', darkColor: '#60A5FA' },
  Cardio:      { bg: '#F0FDF4', color: '#16A34A', darkBg: '#052e16', darkColor: '#4ade80' },
  Flexibility: { bg: '#FFFBEB', color: '#D97706', darkBg: '#431407', darkColor: '#fbbf24' },
};

const DEFAULT_REST = 60; // seconds

// ── Helpers ─────────────────────────────────────────────────
function fmtTime(s) {
  const m = String(Math.floor(s / 60)).padStart(2,'0');
  const sec = String(s % 60).padStart(2,'0');
  return `${m}:${sec}`;
}

// ── Workout Tracker Modal ────────────────────────────────────
function WorkoutTracker({ onClose, onSave }) {
  const t = window.useTheme();

  // exercise list for this session
  const [exercises, setExercises] = useState([
    { name: 'Bench Press', sets: 3, reps: 10 },
    { name: 'Pull Ups', sets: 3, reps: 8 },
    { name: 'Shoulder Press', sets: 3, reps: 12 },
  ]);
  const [exIdx, setExIdx] = useState(0);          // current exercise index
  const [phase, setPhase] = useState('ready');    // ready | set | rest | done
  const [setsDone, setSetsDone] = useState([]);   // [{exName, setNum, elapsed}]
  const [setNum, setSetNum] = useState(1);        // current set number
  const [elapsed, setElapsed] = useState(0);      // timer (up for set, down for rest)
  const [totalSecs, setTotalSecs] = useState(0);  // total workout time
  const [restSecs, setRestSecs] = useState(DEFAULT_REST);
  const timerRef = useRef(null);

  const curEx = exercises[exIdx];

  // Total workout timer
  useEffect(() => {
    const id = setInterval(() => setTotalSecs(s => s + 1), 1000);
    return () => clearInterval(id);
  }, []);

  // Phase timer
  useEffect(() => {
    clearInterval(timerRef.current);
    if (phase === 'set') {
      timerRef.current = setInterval(() => setElapsed(s => s + 1), 1000);
    } else if (phase === 'rest') {
      timerRef.current = setInterval(() => {
        setElapsed(s => {
          if (s <= 1) {
            clearInterval(timerRef.current);
            setPhase('set');
            setElapsed(0);
            return 0;
          }
          return s - 1;
        });
      }, 1000);
    }
    return () => clearInterval(timerRef.current);
  }, [phase]);

  function handleDoneSet() {
    const log = { exName: curEx.name, setNum, elapsed };
    setSetsDone(prev => [...prev, log]);
    const isLastSet = setNum >= curEx.sets;
    const isLastEx = exIdx >= exercises.length - 1;

    if (isLastSet && isLastEx) {
      setPhase('done');
    } else {
      // go to rest
      setElapsed(restSecs);
      setPhase('rest');
      if (isLastSet) {
        setSetNum(1);
        setExIdx(i => i + 1);
      } else {
        setSetNum(n => n + 1);
      }
    }
  }

  function handleSkipRest() {
    clearInterval(timerRef.current);
    setElapsed(0);
    setPhase('set');
  }

  // ── Circular progress for rest ───────────────────────────
  const restPct = phase === 'rest' ? elapsed / restSecs : 0;
  const circumference = 2 * Math.PI * 54;
  const strokeDash = circumference * restPct;

  // ── Layout ───────────────────────────────────────────────
  const bg = t.dark ? '#0F172A' : '#fff';
  const overlay = 'rgb(0 0 0 / .6)';

  if (phase === 'done') {
    const totalMins = Math.round(totalSecs / 60);
    return (
      <div style={{ position: 'fixed', inset: 0, background: overlay, display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 300 }}>
        <div style={{ background: bg, borderRadius: 16, width: 380, padding: '32px 24px', textAlign: 'center', boxShadow: t.modalShadow }}>
          <div style={{ fontSize: 48, marginBottom: 12 }}>🎉</div>
          <h2 style={{ fontSize: 22, fontWeight: 700, color: t.fg1, marginBottom: 4 }}>Workout Complete!</h2>
          <p style={{ fontSize: 14, color: t.fg2, marginBottom: 24 }}>{totalMins} min · {setsDone.length} sets logged</p>
          <div style={{ background: t.surfaceActive, borderRadius: 10, padding: 16, marginBottom: 20, textAlign: 'left' }}>
            {setsDone.map((s, i) => (
              <div key={i} style={{ display: 'flex', justifyContent: 'space-between', padding: '5px 0', borderBottom: i < setsDone.length - 1 ? `1px solid ${t.border}` : 'none' }}>
                <span style={{ fontSize: 13, color: t.fg1 }}>{s.exName} — Set {s.setNum}</span>
                <span style={{ fontSize: 13, color: t.fg2, fontFamily: 'monospace' }}>{fmtTime(s.elapsed)}</span>
              </div>
            ))}
          </div>
          <div style={{ display: 'flex', gap: 10 }}>
            <button onClick={onClose} style={{ flex: 1, padding: '10px', border: `1px solid ${t.border}`, borderRadius: 10, background: t.surface, fontSize: 14, cursor: 'pointer', color: t.fg2 }}>Discard</button>
            <button onClick={() => { onSave(totalMins, setsDone.length); onClose(); }} style={{ flex: 2, padding: '10px', border: 'none', borderRadius: 10, background: '#2563EB', color: '#fff', fontSize: 14, fontWeight: 600, cursor: 'pointer' }}>Save Workout</button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div style={{ position: 'fixed', inset: 0, background: overlay, display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 300 }}>
      <div style={{ background: bg, borderRadius: 16, width: 400, boxShadow: t.modalShadow, overflow: 'hidden' }}>

        {/* Header */}
        <div style={{ background: '#2563EB', padding: '14px 20px', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div>
            <div style={{ fontSize: 12, color: 'rgba(255,255,255,.7)', marginBottom: 1 }}>Workout in progress</div>
            <div style={{ fontFamily: 'monospace', fontSize: 20, fontWeight: 600, color: '#fff', letterSpacing: '0.05em' }}>{fmtTime(totalSecs)}</div>
          </div>
          <button onClick={onClose} style={{ background: 'rgba(255,255,255,.2)', border: 'none', borderRadius: 8, padding: '6px 12px', color: '#fff', fontSize: 13, cursor: 'pointer' }}>End</button>
        </div>

        {/* Exercise progress pills */}
        <div style={{ padding: '14px 20px 0', display: 'flex', gap: 6 }}>
          {exercises.map((ex, i) => (
            <div key={i} style={{ flex: 1, height: 4, borderRadius: 9999, background: i < exIdx ? '#2563EB' : i === exIdx ? (t.dark ? '#1E3A8A' : '#BFDBFE') : t.border }}>
              {i === exIdx && (
                <div style={{ height: '100%', borderRadius: 9999, background: '#2563EB', width: `${(setNum - 1) / curEx.sets * 100}%`, transition: 'width .3s' }} />
              )}
            </div>
          ))}
        </div>

        {/* Current exercise */}
        <div style={{ padding: '16px 20px 0', textAlign: 'center' }}>
          <div style={{ fontSize: 12, color: t.fg3, marginBottom: 4 }}>Exercise {exIdx + 1} of {exercises.length}</div>
          <div style={{ fontSize: 20, fontWeight: 700, color: t.fg1 }}>{curEx.name}</div>
          <div style={{ fontSize: 13, color: t.fg2, marginTop: 2 }}>{curEx.reps} reps per set</div>
        </div>

        {/* Main timer area */}
        <div style={{ padding: '20px', display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 16 }}>

          {phase === 'set' && (
            <>
              {/* Set indicator */}
              <div style={{ display: 'flex', gap: 8 }}>
                {Array.from({ length: curEx.sets }, (_, i) => (
                  <div key={i} style={{ width: 32, height: 32, borderRadius: '50%', border: `2px solid ${i < setNum - 1 ? '#2563EB' : i === setNum - 1 ? '#2563EB' : t.border}`, background: i < setNum - 1 ? '#2563EB' : 'transparent', display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: 13, fontWeight: 600, color: i < setNum - 1 ? '#fff' : i === setNum - 1 ? '#2563EB' : t.fg3 }}>
                    {i < setNum - 1 ? '✓' : i + 1}
                  </div>
                ))}
              </div>

              {/* Set timer */}
              <div style={{ fontFamily: "'JetBrains Mono', monospace", fontSize: 52, fontWeight: 500, color: t.fg1, letterSpacing: '-0.02em', lineHeight: 1 }}>
                {fmtTime(elapsed)}
              </div>
              <div style={{ fontSize: 13, color: t.fg3 }}>Set {setNum} of {curEx.sets} — tap when done</div>

              {/* Done set button */}
              <button onClick={handleDoneSet} style={{ width: '100%', padding: '18px', background: '#2563EB', border: 'none', borderRadius: 12, fontSize: 17, fontWeight: 700, color: '#fff', cursor: 'pointer', letterSpacing: '0.01em', transition: 'transform .1s', boxShadow: '0 4px 12px rgb(37 99 235 / .35)' }}
                onMouseDown={e => e.currentTarget.style.transform = 'scale(0.97)'}
                onMouseUp={e => e.currentTarget.style.transform = 'scale(1)'}>
                ✓ Done Set {setNum}
              </button>
            </>
          )}

          {phase === 'rest' && (
            <>
              {/* Circular rest countdown */}
              <div style={{ position: 'relative', width: 140, height: 140 }}>
                <svg width="140" height="140" style={{ transform: 'rotate(-90deg)' }}>
                  <circle cx="70" cy="70" r="54" fill="none" stroke={t.border} strokeWidth="8" />
                  <circle cx="70" cy="70" r="54" fill="none" stroke="#2563EB" strokeWidth="8"
                    strokeDasharray={`${strokeDash} ${circumference}`}
                    strokeLinecap="round" style={{ transition: 'stroke-dasharray .9s linear' }} />
                </svg>
                <div style={{ position: 'absolute', inset: 0, display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center' }}>
                  <div style={{ fontFamily: "'JetBrains Mono', monospace", fontSize: 32, fontWeight: 600, color: t.fg1, lineHeight: 1 }}>{elapsed}</div>
                  <div style={{ fontSize: 11, color: t.fg3, marginTop: 2 }}>seconds</div>
                </div>
              </div>

              <div style={{ fontSize: 15, fontWeight: 600, color: t.fg1 }}>Rest time</div>
              <div style={{ fontSize: 13, color: t.fg2 }}>
                {exIdx < exercises.length - 1 && setNum === 1
                  ? `Next: ${exercises[exIdx].name}`
                  : `Next: Set ${setNum} of ${curEx.name}`}
              </div>

              <button onClick={handleSkipRest} style={{ width: '100%', padding: '14px', background: t.surfaceActive, border: `1px solid ${t.border}`, borderRadius: 12, fontSize: 15, fontWeight: 600, color: t.fg1, cursor: 'pointer' }}>
                Skip Rest →
              </button>
            </>
          )}

          {phase === 'ready' && (
            <button onClick={() => setPhase('set')} style={{ width: '100%', padding: '18px', background: '#2563EB', border: 'none', borderRadius: 12, fontSize: 17, fontWeight: 700, color: '#fff', cursor: 'pointer', boxShadow: '0 4px 12px rgb(37 99 235 / .35)' }}>
              ▶ Start First Set
            </button>
          )}
        </div>

        {/* Log of completed sets */}
        {setsDone.length > 0 && (
          <div style={{ borderTop: `1px solid ${t.border}`, padding: '12px 20px', maxHeight: 120, overflowY: 'auto' }}>
            <div style={{ fontSize: 11, fontWeight: 600, color: t.fg3, textTransform: 'uppercase', letterSpacing: '.05em', marginBottom: 8 }}>Completed</div>
            {setsDone.map((s, i) => (
              <div key={i} style={{ display: 'flex', justifyContent: 'space-between', padding: '4px 0', fontSize: 13 }}>
                <span style={{ color: t.fg2 }}>{s.exName} — Set {s.setNum}</span>
                <span style={{ color: t.fg3, fontFamily: 'monospace' }}>{fmtTime(s.elapsed)}</span>
              </div>
            ))}
          </div>
        )}

        {/* Rest duration setting */}
        <div style={{ borderTop: `1px solid ${t.border}`, padding: '10px 20px', display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
          <span style={{ fontSize: 12, color: t.fg3 }}>Rest duration</span>
          <div style={{ display: 'flex', gap: 6 }}>
            {[30, 60, 90, 120].map(s => (
              <button key={s} onClick={() => setRestSecs(s)} style={{ padding: '3px 8px', border: `1px solid ${restSecs === s ? '#2563EB' : t.border}`, borderRadius: 6, background: restSecs === s ? t.primarySubtle : t.surface, color: restSecs === s ? '#2563EB' : t.fg3, fontSize: 11, cursor: 'pointer', fontWeight: restSecs === s ? 600 : 400 }}>{s}s</button>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}

// ── TypeBadge ────────────────────────────────────────────────
function TypeBadge({ type }) {
  const t = window.useTheme();
  const c = TYPE_COLORS[type] || { bg: '#F1F5F9', color: '#475569', darkBg: '#1E293B', darkColor: '#94A3B8' };
  return (
    <span style={{ background: t.dark ? c.darkBg : c.bg, color: t.dark ? c.darkColor : c.color, borderRadius: 9999, padding: '2px 9px', fontSize: 11, fontWeight: 500 }}>{type}</span>
  );
}

// ── WorkoutCard ──────────────────────────────────────────────
function WorkoutCard({ w, onDelete }) {
  const t = window.useTheme();
  return (
    <div style={{ background: t.surface, border: `1px solid ${t.border}`, borderRadius: 8, boxShadow: t.shadow, overflow: 'hidden' }}>
      <div style={{ padding: '14px 16px', display: 'flex', justifyContent: 'space-between', alignItems: 'center', borderBottom: `1px solid ${t.border}` }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
          <span style={{ fontSize: 18 }}>💪</span>
          <div>
            <div style={{ fontSize: 14, fontWeight: 600, color: t.fg1 }}>{w.name}</div>
            <div style={{ fontSize: 12, color: t.fg2 }}>{w.duration} min</div>
          </div>
        </div>
        <TypeBadge type={w.type} />
      </div>
      <div style={{ padding: '10px 16px', fontSize: 13, color: t.fg2, lineHeight: 1.7 }}>
        {w.exercises.map((e, i) => <div key={i}>• {e}</div>)}
      </div>
      <div style={{ padding: '10px 16px', display: 'flex', justifyContent: 'space-between', alignItems: 'center', borderTop: `1px solid ${t.border}` }}>
        <span style={{ fontSize: 11, color: t.fg3 }}>{w.date}</span>
        <div style={{ display: 'flex', gap: 6 }}>
          <button style={{ padding: '4px 10px', border: `1px solid ${t.border}`, borderRadius: 6, background: t.surface, fontSize: 12, color: t.fg2, cursor: 'pointer' }}>Edit</button>
          <button onClick={() => onDelete(w.id)} style={{ padding: '4px 10px', border: `1px solid ${t.dangerBorder}`, borderRadius: 6, background: t.dangerSubtle, fontSize: 12, color: t.dangerText, cursor: 'pointer' }}>Delete</button>
        </div>
      </div>
    </div>
  );
}

// ── Workouts page ────────────────────────────────────────────
function Workouts() {
  const t = window.useTheme();
  const [workouts, setWorkouts] = useState(WORKOUT_DATA);
  const [filter, setFilter] = useState('All');
  const [showTracker, setShowTracker] = useState(false);
  const filtered = filter === 'All' ? workouts : workouts.filter(w => w.type === filter);

  return (
    <div style={{ maxWidth: 760 }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
        <h1 style={{ fontSize: 22, fontWeight: 700, color: t.fg1 }}>Workouts</h1>
        <button onClick={() => setShowTracker(true)} style={{ display: 'flex', alignItems: 'center', gap: 6, padding: '8px 18px', background: '#2563EB', color: '#fff', border: 'none', borderRadius: 8, fontSize: 14, fontWeight: 600, cursor: 'pointer', boxShadow: '0 2px 8px rgb(37 99 235/.3)' }}>
          ▶ Start Workout
        </button>
      </div>

      {/* Stats */}
      <div style={{ display: 'flex', gap: 10, marginBottom: 16 }}>
        {[['5','This week'],['18','This month'],['🔥 7','Day streak']].map(([v,l]) => (
          <div key={l} style={{ background: t.surface, border: `1px solid ${t.border}`, borderRadius: 8, padding: '12px 16px', flex: 1, textAlign: 'center' }}>
            <div style={{ fontSize: 20, fontWeight: 700, color: t.fg1 }}>{v}</div>
            <div style={{ fontSize: 12, color: t.fg2, marginTop: 2 }}>{l}</div>
          </div>
        ))}
      </div>

      {/* Filters */}
      <div style={{ display: 'flex', gap: 8, marginBottom: 16, flexWrap: 'wrap' }}>
        {['All','Strength','Cardio','Flexibility'].map(f => (
          <button key={f} onClick={() => setFilter(f)} style={{
            padding: '6px 14px', borderRadius: 9999, fontSize: 13, fontWeight: 500, cursor: 'pointer',
            background: filter === f ? '#2563EB' : t.surface,
            color: filter === f ? '#fff' : t.fg2,
            border: filter === f ? '1px solid #2563EB' : `1px solid ${t.border}`,
          }}>{f}</button>
        ))}
      </div>

      {/* List */}
      <div style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
        {filtered.map(w => <WorkoutCard key={w.id} w={w} onDelete={id => setWorkouts(ws => ws.filter(w => w.id !== id))} />)}
      </div>

      {showTracker && (
        <WorkoutTracker
          onClose={() => setShowTracker(false)}
          onSave={(mins, sets) => {
            setWorkouts(ws => [{ id: Date.now(), type: 'Strength', name: 'Tracked Session', duration: mins, exercises: [`${sets} sets logged`], date: 'Today' }, ...ws]);
          }}
        />
      )}
    </div>
  );
}

Object.assign(window, { Workouts });
