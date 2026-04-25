// Jeeb UI Kit — Sleep Page v2

const { useState, useEffect, useRef } = React;

const SLEEP_DATA = [
  { id: 1, date: 'Apr 22', bedtime: '11:00 PM', wake: '6:30 AM', duration: '7h 30m', quality: 4, note: 'Felt refreshed' },
  { id: 2, date: 'Apr 21', bedtime: '11:30 PM', wake: '7:00 AM', duration: '7h 30m', quality: 3, note: '' },
  { id: 3, date: 'Apr 20', bedtime: '10:30 PM', wake: '6:00 AM', duration: '7h 30m', quality: 5, note: 'Great sleep!' },
  { id: 4, date: 'Apr 19', bedtime: '12:00 AM', wake: '7:30 AM', duration: '7h 30m', quality: 2, note: 'Restless night' },
  { id: 5, date: 'Apr 18', bedtime: '10:00 PM', wake: '6:30 AM', duration: '8h 30m', quality: 5, note: '' },
  { id: 6, date: 'Apr 17', bedtime: '11:45 PM', wake: '6:45 AM', duration: '7h 0m',  quality: 3, note: '' },
  { id: 7, date: 'Apr 16', bedtime: '10:15 PM', wake: '5:45 AM', duration: '7h 30m', quality: 4, note: '' },
];
const CHART_HOURS = [7.5, 7.5, 7.5, 7.5, 8.5, 7.0, 7.5];
const DAYS = ['Mon','Tue','Wed','Thu','Fri','Sat','Sun'];

function SleepChart() {
  const t = window.useTheme();
  const maxH = Math.max(...CHART_HOURS);
  return (
    <div style={{ display: 'flex', alignItems: 'flex-end', gap: 6, height: 80, paddingTop: 8 }}>
      {CHART_HOURS.map((h, i) => (
        <div key={i} style={{ flex: 1, display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 4 }}>
          <div style={{ fontSize: 9, color: t.fg3 }}>{h}h</div>
          <div style={{ width: '100%', borderRadius: '4px 4px 0 0', background: i === CHART_HOURS.length - 1 ? '#2563EB' : (t.dark ? '#1E3A8A' : '#BFDBFE'), height: `${(h / maxH) * 52}px` }} />
          <div style={{ fontSize: 9, color: t.fg3 }}>{DAYS[i]}</div>
        </div>
      ))}
    </div>
  );
}

function Stars({ count }) {
  return <span style={{ fontSize: 12 }}>{[1,2,3,4,5].map(i => <span key={i} style={{ color: i <= count ? '#F59E0B' : '#CBD5E1' }}>★</span>)}</span>;
}

function SleepLogCard({ s }) {
  const t = window.useTheme();
  return (
    <div style={{ background: t.surface, border: `1px solid ${t.border}`, borderRadius: 8, padding: '14px 16px', display: 'flex', alignItems: 'center', gap: 12 }}>
      <div style={{ width: 36, height: 36, borderRadius: '50%', background: t.primarySubtle, display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: 18, flexShrink: 0 }}>😴</div>
      <div style={{ flex: 1 }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 2 }}>
          <span style={{ fontSize: 14, fontWeight: 600, color: t.fg1 }}>{s.date}</span>
          <Stars count={s.quality} />
        </div>
        <div style={{ fontSize: 12, color: t.fg2 }}>{s.bedtime} → {s.wake} · {s.duration}</div>
        {s.note && <div style={{ fontSize: 12, color: t.fg3, marginTop: 2, fontStyle: 'italic' }}>{s.note}</div>}
      </div>
    </div>
  );
}

function SleepForm({ onClose, onAdd }) {
  const t = window.useTheme();
  const [form, setForm] = useState({ bedtime: '23:00', wake: '06:30', quality: 4, note: '' });
  const set = (k, v) => setForm(f => ({ ...f, [k]: v }));
  const inputStyle = { width: '100%', border: `1px solid ${t.border}`, borderRadius: 8, padding: '8px 12px', fontSize: 14, fontFamily: 'inherit', outline: 'none', background: t.inputBg, color: t.fg1 };
  return (
    <div style={{ position: 'fixed', inset: 0, background: 'rgb(0 0 0/.5)', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 200 }}>
      <div style={{ background: t.surface, borderRadius: 12, width: 400, boxShadow: t.modalShadow }}>
        <div style={{ padding: '16px 20px', borderBottom: `1px solid ${t.border}`, display: 'flex', justifyContent: 'space-between' }}>
          <h2 style={{ fontSize: 16, fontWeight: 600, color: t.fg1 }}>Log Sleep</h2>
          <button onClick={onClose} style={{ background: 'none', border: 'none', cursor: 'pointer', color: t.fg3, fontSize: 20 }}>×</button>
        </div>
        <div style={{ padding: 20, display: 'flex', flexDirection: 'column', gap: 14 }}>
          {[['Bedtime','bedtime'],['Wake time','wake']].map(([label, key]) => (
            <div key={key}>
              <label style={{ fontSize: 13, fontWeight: 500, color: t.fg1, display: 'block', marginBottom: 4 }}>{label}</label>
              <input type="time" value={form[key]} onChange={e => set(key, e.target.value)} style={inputStyle} />
            </div>
          ))}
          <div>
            <label style={{ fontSize: 13, fontWeight: 500, color: t.fg1, display: 'block', marginBottom: 6 }}>Quality</label>
            <div style={{ display: 'flex', gap: 4 }}>
              {[1,2,3,4,5].map(n => <span key={n} onClick={() => set('quality', n)} style={{ fontSize: 24, cursor: 'pointer', color: n <= form.quality ? '#F59E0B' : t.border }}>★</span>)}
            </div>
          </div>
          <div>
            <label style={{ fontSize: 13, fontWeight: 500, color: t.fg1, display: 'block', marginBottom: 4 }}>Notes</label>
            <input type="text" placeholder="How did you sleep?" value={form.note} onChange={e => set('note', e.target.value)} style={inputStyle} />
          </div>
        </div>
        <div style={{ padding: '14px 20px', borderTop: `1px solid ${t.border}`, display: 'flex', justifyContent: 'flex-end', gap: 8 }}>
          <button onClick={onClose} style={{ padding: '8px 16px', border: `1px solid ${t.border}`, borderRadius: 8, background: t.surface, fontSize: 14, cursor: 'pointer', color: t.fg2 }}>Cancel</button>
          <button onClick={() => { onAdd({ id: Date.now(), date: 'Today', bedtime: form.bedtime, wake: form.wake, duration: '?h', quality: form.quality, note: form.note }); onClose(); }}
            style={{ padding: '8px 18px', border: 'none', borderRadius: 8, background: '#2563EB', color: '#fff', fontSize: 14, fontWeight: 500, cursor: 'pointer' }}>Save</button>
        </div>
      </div>
    </div>
  );
}

function Sleep() {
  const t = window.useTheme();
  const [logs, setLogs] = useState(SLEEP_DATA);
  const [showForm, setShowForm] = useState(false);
  const width = window.useWindowWidth();
  const isMobile = width < 768;

  return (
    <div style={{ maxWidth: 760 }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
        <h1 style={{ fontSize: 22, fontWeight: 700, color: t.fg1 }}>Sleep</h1>
        <button onClick={() => setShowForm(true)} style={{ padding: '8px 16px', background: '#2563EB', color: '#fff', border: 'none', borderRadius: 8, fontSize: 14, fontWeight: 500, cursor: 'pointer' }}>+ Log Sleep</button>
      </div>
      <div style={{ display: 'grid', gridTemplateColumns: isMobile ? '1fr' : '1fr 1fr', gap: 12, marginBottom: 16 }}>
        <div style={{ background: t.surface, border: `1px solid ${t.border}`, borderRadius: 8, padding: 16 }}>
          <div style={{ fontSize: 13, fontWeight: 600, color: t.fg1, marginBottom: 4 }}>Weekly Sleep</div>
          <SleepChart />
        </div>
        <div style={{ background: t.surface, border: `1px solid ${t.border}`, borderRadius: 8, padding: 16, display: 'flex', flexDirection: 'column', gap: 10 }}>
          {[['7.2h','Weekly average'],['8.5h','Best this week'],['4.2 / 5','Avg quality']].map(([v, l]) => (
            <div key={l} style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <span style={{ fontSize: 13, color: t.fg2 }}>{l}</span>
              <span style={{ fontSize: 15, fontWeight: 600, color: t.fg1 }}>{v}</span>
            </div>
          ))}
        </div>
      </div>
      <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
        {logs.map(s => <SleepLogCard key={s.id} s={s} />)}
      </div>
      {showForm && <SleepForm onClose={() => setShowForm(false)} onAdd={s => setLogs(ls => [s, ...ls])} />}
    </div>
  );
}

Object.assign(window, { Sleep });
