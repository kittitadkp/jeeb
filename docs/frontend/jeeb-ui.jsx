// Shared UI components + charts for Jeeb redesign

const C = {
  bg: '#090b10',
  surface: '#0f1117',
  surface2: '#15171f',
  surface3: '#1c1f2e',
  border: 'rgba(255,255,255,0.07)',
  border2: 'rgba(255,255,255,0.13)',
  text: '#e8eaf6',
  text2: '#7b7f9e',
  text3: '#3a3d56',
  primary: '#7c6ef5',
  sections: {
    dashboard: '#7c6ef5',
    workouts: '#f43f5e',
    study: '#f59e0b',
    sleep: '#a855f7',
    finance: '#10b981',
    calendar: '#14b8a6',
    settings: '#6b7280',
    goals: '#f97316',
    events: '#06b6d4',
  }
};

// ── Period data ───────────────────────────────────────────────────────────────

const PERIOD_DATA = {
  workouts: {
    day:   { value:1,  unit:'workout',  change:'45 min today',      spark:[0,0,45,0,0,0,0], chart:[{label:'6AM',value:45},{label:'10AM',value:0},{label:'2PM',value:0},{label:'6PM',value:0},{label:'10PM',value:0}] },
    week:  { value:4,  unit:'workouts', change:'6-day streak 🔥',   spark:[45,0,60,30,0,75,45], chart:['Mon','Tue','Wed','Thu','Fri','Sat','Sun'].map((l,i)=>({label:l,value:[45,0,60,30,0,75,45][i]})) },
    month: { value:14, unit:'workouts', change:'best: Wk3',         spark:[3,4,4,3], chart:[{label:'Wk1',value:3},{label:'Wk2',value:4},{label:'Wk3',value:4},{label:'Wk4',value:3}] },
    year:  { value:52, unit:'workouts', change:'+18% vs last year', spark:[12,10,14,16], chart:[{label:'Jan',value:12},{label:'Feb',value:10},{label:'Mar',value:14},{label:'Apr',value:16}] },
  },
  study: {
    day:   { value:'2h',    unit:'today',       change:'Mathematics',       spark:[0,0,2,0,0,0,0], chart:[{label:'9AM',value:1},{label:'12PM',value:0.5},{label:'3PM',value:0.5},{label:'6PM',value:0},{label:'9PM',value:0}] },
    week:  { value:'14.5h', unit:'this week',   change:'52h this month',    spark:[2,3.5,1,4,0,2.5,1.5], chart:['Mon','Tue','Wed','Thu','Fri','Sat','Sun'].map((l,i)=>({label:l,value:[2,3.5,1,4,0,2.5,1.5][i]})) },
    month: { value:'52h',   unit:'this month',  change:'18 sessions',       spark:[11,14,15,12], chart:[{label:'Wk1',value:11},{label:'Wk2',value:14},{label:'Wk3',value:15},{label:'Wk4',value:12}] },
    year:  { value:'180h',  unit:'this year',   change:'+32% vs last year', spark:[40,35,52,53], chart:[{label:'Jan',value:40},{label:'Feb',value:35},{label:'Mar',value:52},{label:'Apr',value:53}] },
  },
  sleep: {
    day:   { value:'7.5h', unit:'last night',  change:'quality 4/5',       spark:[0,0,7.5,0,0,0,0], chart:[{label:'11PM',value:7.5}] },
    week:  { value:'7.2h', unit:'avg this wk', change:'quality 3.8/5',     spark:[6.5,7,8,6,7.5,8.5,7], chart:['Mon','Tue','Wed','Thu','Fri','Sat','Sun'].map((l,i)=>({label:l,value:[6.5,7,8,6,7.5,8.5,7][i]})) },
    month: { value:'7.1h', unit:'avg / night', change:'24 nights logged',  spark:[7,6.8,7.2,7.4], chart:[{label:'Wk1',value:7},{label:'Wk2',value:6.8},{label:'Wk3',value:7.2},{label:'Wk4',value:7.4}] },
    year:  { value:'7.0h', unit:'avg / night', change:'consistent sleeper', spark:[7.1,6.9,7.0,7.2], chart:[{label:'Jan',value:7.1},{label:'Feb',value:6.9},{label:'Mar',value:7.0},{label:'Apr',value:7.2}] },
  },
  finance: {
    day:   { value:'$128', unit:'spent today',   change:'-$42 transport',  spark:[0,85,0,0,43,0,0], chart:[{label:'9AM',value:85},{label:'12PM',value:0},{label:'3PM',value:43},{label:'6PM',value:0},{label:'9PM',value:0}] },
    week:  { value:'$840', unit:'this week',     change:'on track',        spark:[120,340,80,450,200,560,90], chart:['Mon','Tue','Wed','Thu','Fri','Sat','Sun'].map((l,i)=>({label:l,value:[120,340,80,450,200,560,90][i]})) },
    month: { value:'$2,840',unit:'this month',   change:'net +$1,660',     spark:[640,280,390,180], chart:[{label:'Wk1',value:640},{label:'Wk2',value:720},{label:'Wk3',value:890},{label:'Wk4',value:590}] },
    year:  { value:'$9.2k', unit:'this year',    change:'avg $2.3k/mo',    spark:[2200,2400,2840,1800], chart:[{label:'Jan',value:2200},{label:'Feb',value:2400},{label:'Mar',value:2840},{label:'Apr',value:1800}] },
  },
};

const MOCK = {
  user: { name: 'Alex' },
  workouts: {
    thisWeek: 4, thisMonth: 14, streak: 6,
    weeklyMins: [45, 0, 60, 30, 0, 75, 45],
    recent: [
      { id:1, type:'strength', duration:45, exercises:['Bench Press 3×8 @ 80kg','Back Squat 3×10 @ 70kg','Romanian DL 3×8'], date:'2026-04-25' },
      { id:2, type:'cardio', duration:30, exercises:['5km run — 26:40'], date:'2026-04-23' },
      { id:3, type:'flexibility', duration:20, exercises:['Morning yoga flow'], date:'2026-04-22' },
    ]
  },
  study: {
    thisWeek: 14.5, thisMonth: 52, sessions: 18,
    weeklyHours: [2, 3.5, 1, 4, 0, 2.5, 1.5],
    recent: [
      { id:1, subject:'Mathematics', duration:120, date:'2026-04-25' },
      { id:2, subject:'Physics', duration:90, date:'2026-04-24' },
      { id:3, subject:'Computer Science', duration:150, date:'2026-04-22' },
    ]
  },
  sleep: {
    avgDuration: 7.2, avgQuality: 3.8, thisMonth: 24,
    weeklyHours: [6.5, 7, 8, 6, 7.5, 8.5, 7],
    weeklyQuality: [3, 4, 5, 3, 4, 5, 4],
    recent: [
      { id:1, start:'2026-04-25T23:30', end:'2026-04-26T07:00', quality:4, notes:'Felt rested' },
      { id:2, start:'2026-04-24T23:00', end:'2026-04-25T06:30', quality:3, notes:'' },
      { id:3, start:'2026-04-23T22:00', end:'2026-04-24T07:00', quality:5, notes:'Perfect night' },
    ]
  },
  finance: {
    income: 4500, expense: 2840, net: 1660,
    byCategory: { 'Food & Dining':640, 'Transport':280, 'Entertainment':390, 'Utilities':180, 'Healthcare':150, 'Shopping':1200 },
    weeklyExpense: [120, 340, 80, 450, 200, 560, 90],
    transactions: [
      { id:1, type:'income', amount:4500, category:'Salary', notes:'Monthly salary', date:'2026-04-01' },
      { id:2, type:'expense', amount:85.50, category:'Food & Dining', notes:'Grocery run', date:'2026-04-25' },
      { id:3, type:'expense', amount:42, category:'Transport', notes:'Uber rides', date:'2026-04-24' },
      { id:4, type:'expense', amount:120, category:'Entertainment', notes:'Netflix + Spotify', date:'2026-04-20' },
      { id:5, type:'expense', amount:180, category:'Utilities', notes:'Electricity + internet', date:'2026-04-15' },
    ]
  },
  events: [
    { id:1, type:'workout', title:'Morning Run', start:'2026-04-26T07:00' },
    { id:2, type:'study', title:'Physics Exam Prep', start:'2026-04-26T14:00' },
    { id:3, type:'custom', title:'Team Meeting', start:'2026-04-27T10:00' },
    { id:4, type:'finance', title:'Pay Rent', start:'2026-04-28T09:00' },
    { id:5, type:'sleep', title:'Sleep Goal Check', start:'2026-04-29T22:00' },
  ]
};

const NOTIFICATIONS = [
  { id:1, emoji:'🔥', title:'6-day workout streak', desc:'Keep it up — one more day for your best!', time:'5m ago', unread:true, acc:C.sections.workouts },
  { id:2, emoji:'📚', title:'Study goal 72% complete', desc:'14.5 of 20h logged this week', time:'2h ago', unread:true, acc:C.sections.study },
  { id:3, emoji:'😴', title:'Great sleep last night', desc:'7.5h · quality 4/5 — above your average', time:'8h ago', unread:true, acc:C.sections.sleep },
  { id:4, emoji:'💰', title:'Budget at 95%', desc:'$160 left of your $3,000 monthly budget', time:'Today', unread:false, acc:'#f43f5e' },
  { id:5, emoji:'📅', title:'Morning Run tomorrow', desc:'Scheduled for 7:00 AM', time:'Yesterday', unread:false, acc:C.sections.calendar },
];

// ── Charts ────────────────────────────────────────────────────────────────────

function Sparkline({ data, color = C.primary, width = 100, height = 36 }) {
  const w = width, h = height;
  const min = Math.min(...data);
  const max = Math.max(...data);
  const range = max - min || 1;
  const pts = data.map((v, i) => [
    (i / (data.length - 1)) * w,
    h - 4 - ((v - min) / range) * (h - 8)
  ]);
  const line = pts.map((p, i) => `${i === 0 ? 'M' : 'L'}${p[0].toFixed(1)},${p[1].toFixed(1)}`).join(' ');
  const area = `${line} L${w},${h} L0,${h} Z`;
  const id = `sp${color.replace(/[^a-z0-9]/gi,'')}${Math.round(width)}`;
  return (
    <svg width={w} height={h} viewBox={`0 0 ${w} ${h}`} preserveAspectRatio="none" style={{overflow:'visible'}}>
      <defs>
        <linearGradient id={id} x1="0" y1="0" x2="0" y2="1">
          <stop offset="0%" stopColor={color} stopOpacity="0.25" />
          <stop offset="100%" stopColor={color} stopOpacity="0" />
        </linearGradient>
      </defs>
      <path d={area} fill={`url(#${id})`} />
      <path d={line} fill="none" stroke={color} strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" />
      <circle cx={pts[pts.length-1][0]} cy={pts[pts.length-1][1]} r="2.5" fill={color} />
    </svg>
  );
}

function BarChart({ data, color = C.primary, height = 80 }) {
  const vals = data.map(d => d.value);
  const max = Math.max(...vals) || 1;
  const n = data.length;
  // Fixed virtual coordinate space — SVG scales to fill container
  const BAR_W = 28, GAP = 10;
  const VW = n * (BAR_W + GAP) - GAP;
  const VH = height + 20;
  return (
    <div style={{ width:'100%', overflowX:'hidden' }}>
      <svg width="100%" height={VH} viewBox={`0 0 ${VW} ${VH}`}
        preserveAspectRatio="none" style={{ display:'block' }}>
        {data.map((d, i) => {
          const barH = Math.max(3, (d.value / max) * height);
          const x = i * (BAR_W + GAP);
          const y = height - barH;
          const isLast = i === n - 1;
          // label font size in viewBox coords — stays legible after scale
          const fs = Math.min(10, VW / n * 0.38);
          return (
            <g key={i}>
              <rect x={x} y={y} width={BAR_W} height={barH} rx={5}
                fill={isLast ? color : `${color}45`} style={{transition:'all 0.3s'}} />
              <text x={x + BAR_W / 2} y={height + 15} textAnchor="middle"
                fontSize={fs} fill={C.text2} fontFamily="Inter,sans-serif">{d.label}</text>
            </g>
          );
        })}
      </svg>
    </div>
  );
}

function DonutChart({ segments, size = 120 }) {
  const total = segments.reduce((s,x) => s + x.value, 0) || 1;
  const r = 44, cx = size/2, cy = size/2;
  const circ = 2 * Math.PI * r;
  let offset = 0;
  const slices = segments.map(seg => {
    const pct = seg.value / total;
    const dash = pct * circ;
    const slice = {...seg, dash, offset};
    offset += dash;
    return slice;
  });
  return (
    <svg width={size} height={size} viewBox={`0 0 ${size} ${size}`}>
      <circle cx={cx} cy={cy} r={r} fill="none" stroke={C.surface3} strokeWidth="18" />
      {slices.map((s,i) => (
        <circle key={i} cx={cx} cy={cy} r={r} fill="none"
          stroke={s.color} strokeWidth="18"
          strokeDasharray={`${s.dash} ${circ - s.dash}`}
          strokeDashoffset={-s.offset}
          style={{transform:`rotate(-90deg)`,transformOrigin:`${cx}px ${cy}px`}} />
      ))}
      <text x={cx} y={cy-4} textAnchor="middle" fontSize="13" fontWeight="700"
        fill={C.text} fontFamily="Space Grotesk,sans-serif">
        ${(segments.reduce((s,x)=>s+x.value,0)/1000).toFixed(1)}k
      </text>
      <text x={cx} y={cy+14} textAnchor="middle" fontSize="9" fill={C.text2} fontFamily="Inter,sans-serif">spent</text>
    </svg>
  );
}

// ── Period Selector ───────────────────────────────────────────────────────────

function PeriodSelector({ value, onChange, accent }) {
  const periods = ['day','week','month','year'];
  return (
    <div style={{ display:'flex', gap:2, background:C.surface2, borderRadius:20, padding:3 }}>
      {periods.map(p => (
        <button key={p} onClick={() => onChange(p)} style={{
          padding:'4px 11px', borderRadius:14, fontSize:11, fontWeight:600, cursor:'pointer',
          background: value === p ? (accent || C.primary) : 'transparent',
          color: value === p ? '#fff' : C.text2,
          border: 'none', transition:'all 0.15s', textTransform:'capitalize', letterSpacing:'0.02em'
        }}>{p}</button>
      ))}
    </div>
  );
}

// ── UI Primitives ─────────────────────────────────────────────────────────────

function Card({ children, style, accent, onClick }) {
  const [hov, setHov] = React.useState(false);
  return (
    <div onClick={onClick}
      onMouseEnter={() => setHov(true)} onMouseLeave={() => setHov(false)}
      style={{
        background: C.surface, borderRadius: 14, overflow:'hidden',
        border: `1px solid ${hov ? C.border2 : C.border}`,
        transition:'border-color 0.15s, box-shadow 0.15s',
        boxShadow: hov ? `0 8px 32px rgba(0,0,0,0.4)` : '0 1px 3px rgba(0,0,0,0.3)',
        cursor: onClick ? 'pointer' : 'default', position:'relative',
        ...(accent ? { borderTop:`2px solid ${accent}` } : {}), ...style
      }}>
      {children}
    </div>
  );
}

function Btn({ children, color, outline, size='md', onClick, style, disabled }) {
  const [hov, setHov] = React.useState(false);
  const col = color || C.primary;
  const pad = size==='sm' ? '6px 14px' : size==='lg' ? '12px 24px' : '9px 18px';
  const fs = size==='sm' ? 12 : size==='lg' ? 15 : 13;
  return (
    <button onClick={onClick} disabled={disabled}
      onMouseEnter={() => setHov(true)} onMouseLeave={() => setHov(false)}
      style={{
        background: outline ? (hov ? `${col}18` : 'transparent') : (hov ? `${col}dd` : col),
        border: `1px solid ${outline ? `${col}60` : 'transparent'}`,
        color: outline ? col : '#fff', padding: pad, borderRadius:9,
        fontSize:fs, fontWeight:600, display:'inline-flex', alignItems:'center', gap:6,
        transition:'all 0.15s', opacity: disabled ? 0.5 : 1,
        cursor: disabled ? 'not-allowed' : 'pointer', letterSpacing:'0.01em', ...style
      }}>
      {children}
    </button>
  );
}

function Badge({ children, color }) {
  const col = color || C.primary;
  return (
    <span style={{ background:`${col}20`, color:col, border:`1px solid ${col}40`,
      borderRadius:6, padding:'2px 8px', fontSize:11, fontWeight:600,
      letterSpacing:'0.02em', textTransform:'capitalize' }}>
      {children}
    </span>
  );
}

function StatCard({ emoji, label, value, change, trend, sparkData, accent }) {
  const trendColor = trend==='up' ? '#10b981' : trend==='down' ? '#f43f5e' : C.text2;
  const arrow = trend==='up' ? '↑' : trend==='down' ? '↓' : '';
  return (
    <Card accent={accent} style={{ padding:'20px 20px 16px' }}>
      <div style={{ display:'flex', justifyContent:'space-between', alignItems:'flex-start' }}>
        <div style={{ flex:1 }}>
          <div style={{ fontSize:22, marginBottom:8 }}>{emoji}</div>
          <div style={{ fontSize:28, fontWeight:700, color:C.text, lineHeight:1, fontFamily:'Space Grotesk,sans-serif' }}>{value}</div>
          <div style={{ fontSize:11, color:C.text2, marginTop:5, textTransform:'uppercase', letterSpacing:'0.06em' }}>{label}</div>
          {change && <div style={{ fontSize:11, fontWeight:600, color:trendColor, marginTop:6 }}>{arrow} {change}</div>}
        </div>
        {sparkData && (
          <div style={{ alignSelf:'flex-end', marginBottom:4, opacity:0.9 }}>
            <Sparkline data={sparkData} color={accent || C.primary} width={80} height={32} />
          </div>
        )}
      </div>
    </Card>
  );
}

function SectionLabel({ color, children }) {
  return (
    <div style={{ display:'flex', alignItems:'center', gap:8, marginBottom:20 }}>
      <div style={{ width:3, height:20, borderRadius:2, background:color }} />
      <h1 style={{ fontSize:22, fontWeight:700, color:C.text }}>{children}</h1>
    </div>
  );
}

function Input({ label, id, ...props }) {
  const [foc, setFoc] = React.useState(false);
  return (
    <div>
      {label && <label htmlFor={id} style={{ display:'block', fontSize:11, color:C.text2, marginBottom:5, textTransform:'uppercase', letterSpacing:'0.06em' }}>{label}</label>}
      <input id={id} {...props}
        onFocus={e => { setFoc(true); props.onFocus?.(e); }}
        onBlur={e => { setFoc(false); props.onBlur?.(e); }}
        style={{ width:'100%', background:C.surface2, border:`1px solid ${foc ? C.primary : C.border2}`,
          borderRadius:9, padding:'9px 12px', fontSize:13, color:C.text, outline:'none',
          transition:'border-color 0.15s', ...props.style }} />
    </div>
  );
}

function Select({ label, id, children, ...props }) {
  return (
    <div>
      {label && <label htmlFor={id} style={{ display:'block', fontSize:11, color:C.text2, marginBottom:5, textTransform:'uppercase', letterSpacing:'0.06em' }}>{label}</label>}
      <select id={id} {...props} style={{ width:'100%', background:C.surface2, border:`1px solid ${C.border2}`,
        borderRadius:9, padding:'9px 12px', fontSize:13, color:C.text, outline:'none', appearance:'none' }}>
        {children}
      </select>
    </div>
  );
}

function EmptyState({ emoji, title, desc }) {
  return (
    <div style={{ textAlign:'center', padding:'48px 24px' }}>
      <div style={{ fontSize:36, marginBottom:12 }}>{emoji}</div>
      <div style={{ fontSize:15, fontWeight:600, color:C.text, marginBottom:6 }}>{title}</div>
      <div style={{ fontSize:13, color:C.text2 }}>{desc}</div>
    </div>
  );
}

Object.assign(window, {
  C, MOCK, PERIOD_DATA, NOTIFICATIONS,
  Sparkline, BarChart, DonutChart, PeriodSelector,
  Card, Btn, Badge, StatCard, SectionLabel, Input, Select, EmptyState
});
