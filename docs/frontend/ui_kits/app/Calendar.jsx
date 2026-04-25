// Jeeb UI Kit — Calendar Page v2

const { useState } = React;

const EVENTS = {
  '2026-04-23': [{ type: 'workout', label: '💪 Gym session', time: '9:00 AM' }, { type: 'study', label: '📚 Study Math', time: '2:00 PM' }],
  '2026-04-25': [{ type: 'finance', label: '💰 Electricity bill', time: 'Due' }],
  '2026-04-28': [{ type: 'custom',  label: '📅 Doctor appt', time: '11:00 AM' }],
  '2026-04-20': [{ type: 'workout', label: '💪 Cardio', time: '7:30 AM' }],
  '2026-04-17': [{ type: 'study',   label: '📚 Physics', time: '3:00 PM' }],
  '2026-04-15': [{ type: 'finance', label: '💰 Rent due', time: 'Due' }],
};

function Calendar() {
  const t = window.useTheme();
  const width = window.useWindowWidth();
  const isMobile = width < 768;
  const [year, setYear] = useState(2026);
  const [month, setMonth] = useState(3);
  const [selected, setSelected] = useState('2026-04-23');

  const firstDay = new Date(year, month, 1).getDay();
  const daysInMonth = new Date(year, month + 1, 0).getDate();
  const monthName = new Date(year, month, 1).toLocaleString('default', { month: 'long' });
  const dateKey = d => `${year}-${String(month+1).padStart(2,'0')}-${String(d).padStart(2,'0')}`;
  const selectedEvents = EVENTS[selected] || [];

  const typeColors = {
    workout: { bg: t.primarySubtle,   text: t.navActiveText },
    study:   { bg: t.successSubtle,   text: t.successText },
    finance: { bg: t.warningSubtle,   text: t.warningText },
    custom:  { bg: t.surfaceActive,   text: t.fg2 },
  };

  const cells = [];
  for (let i = 0; i < firstDay; i++) cells.push(null);
  for (let d = 1; d <= daysInMonth; d++) cells.push(d);

  return (
    <div style={{ maxWidth: 760 }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
        <h1 style={{ fontSize: 22, fontWeight: 700, color: t.fg1 }}>Calendar</h1>
        <button style={{ padding: '8px 16px', background: '#2563EB', color: '#fff', border: 'none', borderRadius: 8, fontSize: 14, fontWeight: 500, cursor: 'pointer' }}>+ Event</button>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: isMobile ? '1fr' : '1fr 240px', gap: 16 }}>
        {/* Calendar grid */}
        <div style={{ background: t.surface, border: `1px solid ${t.border}`, borderRadius: 8, padding: 16 }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 14 }}>
            <button onClick={() => { if (month===0) { setMonth(11); setYear(y=>y-1); } else setMonth(m=>m-1); }}
              style={{ border: 'none', background: t.surfaceActive, borderRadius: 6, width: 30, height: 30, cursor: 'pointer', fontSize: 16, color: t.fg2 }}>‹</button>
            <span style={{ fontSize: 15, fontWeight: 600, color: t.fg1 }}>{monthName} {year}</span>
            <button onClick={() => { if (month===11) { setMonth(0); setYear(y=>y+1); } else setMonth(m=>m+1); }}
              style={{ border: 'none', background: t.surfaceActive, borderRadius: 6, width: 30, height: 30, cursor: 'pointer', fontSize: 16, color: t.fg2 }}>›</button>
          </div>
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(7,1fr)', gap: 2, marginBottom: 4 }}>
            {['Sun','Mon','Tue','Wed','Thu','Fri','Sat'].map(d => (
              <div key={d} style={{ textAlign: 'center', fontSize: 11, fontWeight: 600, color: t.fg3, padding: '4px 0' }}>{d}</div>
            ))}
          </div>
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(7,1fr)', gap: 2 }}>
            {cells.map((d, i) => {
              if (!d) return <div key={`e${i}`} />;
              const key = dateKey(d);
              const hasEvents = !!EVENTS[key];
              const isSelected = key === selected;
              const isToday = key === '2026-04-23';
              return (
                <div key={key} onClick={() => setSelected(key)} style={{
                  padding: '6px 4px', textAlign: 'center', borderRadius: 6, cursor: 'pointer', position: 'relative',
                  background: isSelected ? '#2563EB' : isToday ? t.primarySubtle : 'transparent',
                  color: isSelected ? '#fff' : isToday ? '#2563EB' : t.fg1,
                  fontWeight: (isToday || isSelected) ? 600 : 400, fontSize: 13,
                }}
                onMouseEnter={e => { if (!isSelected) e.currentTarget.style.background = t.surfaceActive; }}
                onMouseLeave={e => { if (!isSelected) e.currentTarget.style.background = isToday ? t.primarySubtle : 'transparent'; }}>
                  {d}
                  {hasEvents && <div style={{ width: 4, height: 4, borderRadius: '50%', background: isSelected ? '#fff' : '#2563EB', margin: '2px auto 0' }} />}
                </div>
              );
            })}
          </div>
        </div>

        {/* Event panel */}
        <div style={{ background: t.surface, border: `1px solid ${t.border}`, borderRadius: 8, padding: 16 }}>
          <div style={{ fontSize: 13, fontWeight: 600, color: t.fg1, marginBottom: 12 }}>
            {new Date(selected + 'T12:00:00').toLocaleDateString('default', { month: 'short', day: 'numeric', year: 'numeric' })}
          </div>
          {selectedEvents.length === 0
            ? <div style={{ textAlign: 'center', padding: '24px 0', color: t.fg3, fontSize: 13 }}>No events</div>
            : <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
                {selectedEvents.map((ev, i) => {
                  const c = typeColors[ev.type] || typeColors.custom;
                  return (
                    <div key={i} style={{ background: c.bg, border: `1px solid ${t.border}`, borderRadius: 6, padding: '10px 12px' }}>
                      <div style={{ fontSize: 13, fontWeight: 500, color: t.fg1 }}>{ev.label}</div>
                      <div style={{ fontSize: 11, color: t.fg2, marginTop: 2 }}>{ev.time}</div>
                    </div>
                  );
                })}
              </div>
          }
          <button style={{ marginTop: 12, width: '100%', padding: 8, border: `1px dashed ${t.borderStrong}`, borderRadius: 8, background: 'transparent', fontSize: 13, color: t.fg3, cursor: 'pointer' }}>
            + Add event
          </button>
        </div>
      </div>
    </div>
  );
}

Object.assign(window, { Calendar });
