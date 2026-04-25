// Dashboard + Workouts pages

const DAYS = ['Mon','Tue','Wed','Thu','Fri','Sat','Sun'];

// ── Goals Section ─────────────────────────────────────────────────────────────

function GoalRow({ emoji, label, current, target, displayCurrent, displayTarget, accent, warn }) {
  const pct = Math.min(100, (current / target) * 100);
  const barColor = warn ? '#f43f5e' : accent;
  return (
    <div>
      <div style={{ display:'flex', justifyContent:'space-between', alignItems:'center', marginBottom:6 }}>
        <span style={{ fontSize:12, color:C.text2 }}>{emoji} {label}</span>
        <span style={{ fontSize:11, fontWeight:700, color: warn ? '#f43f5e' : C.text }}>
          {displayCurrent} <span style={{ color:C.text3, fontWeight:400 }}>/ {displayTarget}</span>
        </span>
      </div>
      <div style={{ height:5, borderRadius:6, background:C.surface3, overflow:'hidden' }}>
        <div style={{ height:'100%', width:`${pct}%`, borderRadius:6,
          background:`linear-gradient(90deg, ${barColor}60, ${barColor})`,
          transition:'width 0.6s cubic-bezier(0.4,0,0.2,1)' }} />
      </div>
      <div style={{ textAlign:'right', marginTop:3 }}>
        <span style={{ fontSize:10, color: warn ? '#f43f5e' : C.text3 }}>{Math.round(pct)}%</span>
      </div>
    </div>
  );
}

function GoalsSection() {
  const acc = C.sections;
  const goals = [
    { emoji:'💪', label:'Workouts / week', current:4, target:5, displayCurrent:'4', displayTarget:'5 sessions', accent:acc.workouts },
    { emoji:'📚', label:'Study / week', current:14.5, target:20, displayCurrent:'14.5h', displayTarget:'20h', accent:acc.study },
    { emoji:'😴', label:'Avg sleep', current:7.2, target:8, displayCurrent:'7.2h', displayTarget:'8h', accent:acc.sleep },
    { emoji:'💰', label:'Monthly budget', current:2840, target:3000, displayCurrent:'$2,840', displayTarget:'$3,000', accent:acc.finance, warn:true },
  ];
  const streaks = [
    { emoji:'🔥', label:'Workout', value:6, accent:acc.workouts },
    { emoji:'📚', label:'Study', value:12, accent:acc.study },
    { emoji:'😴', label:'Sleep', value:8, accent:acc.sleep },
  ];
  return (
    <Card style={{ padding:'20px 24px' }}>
      <div style={{ display:'flex', justifyContent:'space-between', alignItems:'center', marginBottom:18 }}>
        <span style={{ fontSize:13, fontWeight:600, color:C.text }}>Goals & Streaks</span>
        <Btn outline color={C.text2} size="sm">Edit Goals</Btn>
      </div>
      <div style={{ display:'grid', gridTemplateColumns:'1fr 1fr', gap:'14px 32px', marginBottom:20 }}>
        {goals.map(g => <GoalRow key={g.label} {...g} />)}
      </div>
      <div style={{ height:1, background:C.border, margin:'4px 0 16px' }} />
      <div style={{ display:'flex', gap:10, flexWrap:'wrap' }}>
        {streaks.map(s => (
          <div key={s.label} style={{ display:'flex', alignItems:'center', gap:8, padding:'7px 14px',
            background:C.surface2, borderRadius:20, border:`1px solid ${s.accent}25` }}>
            <span style={{ fontSize:15 }}>{s.emoji}</span>
            <span style={{ fontSize:20, fontWeight:700, color:s.accent, fontFamily:'Space Grotesk,sans-serif', lineHeight:1 }}>{s.value}</span>
            <div>
              <div style={{ fontSize:10, fontWeight:600, color:C.text, lineHeight:1.3 }}>{s.label}</div>
              <div style={{ fontSize:9, color:C.text2 }}>day streak</div>
            </div>
          </div>
        ))}
      </div>
    </Card>
  );
}

// ── Dashboard ─────────────────────────────────────────────────────────────────

function ActivityItem({ emoji, text, time, accent }) {
  return (
    <div style={{ display:'flex', alignItems:'center', gap:12, padding:'10px 0', borderBottom:`1px solid ${C.border}` }}>
      <div style={{ width:34, height:34, borderRadius:10, flexShrink:0,
        background:`${accent || C.primary}18`, display:'flex', alignItems:'center', justifyContent:'center', fontSize:15 }}>{emoji}</div>
      <span style={{ flex:1, fontSize:13, fontWeight:500, color:C.text }}>{text}</span>
      <span style={{ fontSize:11, color:C.text2, whiteSpace:'nowrap' }}>{time}</span>
    </div>
  );
}

function UpcomingItem({ emoji, title, time, accent }) {
  return (
    <div style={{ display:'flex', alignItems:'center', gap:12, padding:'9px 0', borderBottom:`1px solid ${C.border}` }}>
      <div style={{ width:34, height:34, borderRadius:8, flexShrink:0,
        border:`1px solid ${C.border2}`, background:C.surface2,
        display:'flex', alignItems:'center', justifyContent:'center', fontSize:14 }}>{emoji}</div>
      <div style={{ flex:1 }}>
        <div style={{ fontSize:13, fontWeight:500, color:C.text }}>{title}</div>
        <div style={{ fontSize:11, color:C.text2, marginTop:2 }}>{time}</div>
      </div>
      <div style={{ width:6, height:6, borderRadius:'50%', background:accent || C.primary, flexShrink:0 }} />
    </div>
  );
}

function QuickLogBtn({ emoji, label, accent, onClick }) {
  const [hov, setHov] = React.useState(false);
  return (
    <button onClick={onClick} onMouseEnter={() => setHov(true)} onMouseLeave={() => setHov(false)}
      style={{ display:'flex', alignItems:'center', gap:8, padding:'9px 16px',
        background: hov ? `${accent}18` : C.surface,
        border:`1px solid ${hov ? accent+'60' : C.border}`,
        borderRadius:10, fontSize:13, fontWeight:600,
        color: hov ? accent : C.text2, transition:'all 0.15s', cursor:'pointer' }}>
      <span>{emoji}</span> {label}
    </button>
  );
}

function Dashboard({ navigate, openQuickLog }) {
  const hour = new Date().getHours();
  const greeting = hour < 12 ? 'Good morning' : hour < 17 ? 'Good afternoon' : 'Good evening';
  const today = new Date().toLocaleDateString('en-US',{weekday:'long',month:'long',day:'numeric'});
  const [period, setPeriod] = React.useState('week');
  const acc = C.sections;

  const eventAccent = t => ({ workout:acc.workouts, study:acc.study, sleep:acc.sleep, finance:acc.finance, custom:acc.calendar })[t] || acc.dashboard;
  const eventEmoji = t => ({ workout:'💪', study:'📚', sleep:'😴', finance:'💰', custom:'📅' })[t] || '📅';

  const PD = PERIOD_DATA;
  const periodLabel = { day:'today', week:'this week', month:'this month', year:'this year' }[period];

  const stats = [
    { emoji:'💪', label:`workouts ${periodLabel}`, value:PD.workouts[period].value, change:PD.workouts[period].change, trend:'up', spark:PD.workouts[period].spark, accent:acc.workouts },
    { emoji:'📚', label:`study ${periodLabel}`, value:PD.study[period].value, change:PD.study[period].change, trend:'up', spark:PD.study[period].spark, accent:acc.study },
    { emoji:'😴', label:`avg sleep ${periodLabel}`, value:PD.sleep[period].value, change:PD.sleep[period].change, trend:'neutral', spark:PD.sleep[period].spark, accent:acc.sleep },
    { emoji:'💰', label:`spending ${periodLabel}`, value:PD.finance[period].value, change:PD.finance[period].change, trend:'down', spark:PD.finance[period].spark, accent:acc.finance },
  ];

  const activity = [
    { key:'w1', emoji:'💪', text:'Strength workout — 45 min', time:'2h ago', accent:acc.workouts },
    { key:'s1', emoji:'📚', text:'Studied Mathematics — 2h', time:'5h ago', accent:acc.study },
    { key:'sl1', emoji:'😴', text:'Sleep logged — 7h 30m · quality 4/5', time:'8h ago', accent:acc.sleep },
    { key:'f1', emoji:'💸', text:'Food & Dining — −$85.50', time:'Yesterday', accent:acc.finance },
    { key:'f2', emoji:'💼', text:'Salary — +$4,500', time:'Apr 1', accent:acc.finance },
  ];

  const upcomingEvents = MOCK.events.filter(e => new Date(e.start) > new Date()).sort((a,b) => new Date(a.start)-new Date(b.start));

  return (
    <div style={{ maxWidth:880, display:'flex', flexDirection:'column', gap:24 }}>
      <div style={{ display:'flex', alignItems:'flex-end', justifyContent:'space-between' }}>
        <div>
          <h1 style={{ fontSize:26, fontWeight:700, color:C.text }}>{greeting}, {MOCK.user.name} 👋</h1>
          <p style={{ fontSize:13, color:C.text2, marginTop:4 }}>{today}</p>
        </div>
        <PeriodSelector value={period} onChange={setPeriod} accent={C.primary} />
      </div>

      <div style={{ display:'grid', gridTemplateColumns:'repeat(4,1fr)', gap:12 }}>
        {stats.map(s => <StatCard key={s.label} {...s} sparkData={s.spark} />)}
      </div>

      <GoalsSection />

      <div style={{ display:'grid', gridTemplateColumns:'1fr 1fr', gap:16 }}>
        <Card>
          <div style={{ padding:'18px 20px 4px', borderBottom:`1px solid ${C.border}` }}>
            <span style={{ fontSize:13, fontWeight:600, color:C.text }}>Recent Activity</span>
          </div>
          <div style={{ padding:'4px 20px 8px' }}>
            {activity.map(a => <ActivityItem key={a.key} {...a} />)}
          </div>
        </Card>
        <Card>
          <div style={{ padding:'18px 20px 4px', borderBottom:`1px solid ${C.border}` }}>
            <span style={{ fontSize:13, fontWeight:600, color:C.text }}>Upcoming</span>
          </div>
          <div style={{ padding:'4px 20px 8px' }}>
            {upcomingEvents.map(e => (
              <UpcomingItem key={e.id} emoji={eventEmoji(e.type)} title={e.title} accent={eventAccent(e.type)}
                time={(() => {
                  const d = new Date(e.start);
                  const t = d.toLocaleTimeString('en-US',{hour:'numeric',minute:'2-digit'});
                  const tod = new Date(); const tom = new Date(tod); tom.setDate(tod.getDate()+1);
                  if (d.toDateString()===tod.toDateString()) return `Today · ${t}`;
                  if (d.toDateString()===tom.toDateString()) return `Tomorrow · ${t}`;
                  return `${d.toLocaleDateString('en-US',{month:'short',day:'numeric'})} · ${t}`;
                })()} />
            ))}
          </div>
        </Card>
      </div>

    </div>
  );
}

// ── Workouts ──────────────────────────────────────────────────────────────────

const TYPE_COLOR = { strength:C.sections.workouts, cardio:'#fb923c', flexibility:'#a855f7' };
const FILTERS = ['All','Strength','Cardio','Flexibility'];

function WorkoutCard({ workout, onDelete }) {
  const col = TYPE_COLOR[workout.type] || C.sections.workouts;
  return (
    <Card accent={col}>
      <div style={{ padding:'14px 18px', display:'flex', alignItems:'center', justifyContent:'space-between', borderBottom:`1px solid ${C.border}` }}>
        <div style={{ display:'flex', alignItems:'center', gap:10 }}>
          <div style={{ width:36, height:36, borderRadius:10, background:`${col}20`,
            display:'flex', alignItems:'center', justifyContent:'center', fontSize:17 }}>💪</div>
          <div>
            <div style={{ fontSize:14, fontWeight:600, color:C.text, textTransform:'capitalize' }}>{workout.type} workout</div>
            <div style={{ fontSize:11, color:C.text2, marginTop:2 }}>{workout.duration} min</div>
          </div>
        </div>
        <Badge color={col}>{workout.type}</Badge>
      </div>
      {workout.exercises?.length > 0 && (
        <div style={{ padding:'12px 18px', borderBottom:`1px solid ${C.border}` }}>
          {workout.exercises.map((e,i) => (
            <div key={i} style={{ fontSize:12, color:C.text2, lineHeight:1.7 }}>· {e}</div>
          ))}
        </div>
      )}
      <div style={{ padding:'10px 18px', display:'flex', alignItems:'center', justifyContent:'space-between' }}>
        <span style={{ fontSize:11, color:C.text3 }}>
          {new Date(workout.date).toLocaleDateString('en-US',{month:'short',day:'numeric',year:'numeric'})}
        </span>
        <div style={{ display:'flex', gap:8 }}>
          <Btn outline color={C.text2} size="sm">✏️ Edit</Btn>
          <Btn outline color="#f43f5e" size="sm" onClick={onDelete}>🗑 Delete</Btn>
        </div>
      </div>
    </Card>
  );
}

function Workouts({ openQuickLog }) {
  const [filter, setFilter] = React.useState('All');
  const [workouts, setWorkouts] = React.useState(MOCK.workouts.recent);
  const [period, setPeriod] = React.useState('week');
  const acc = C.sections.workouts;
  const pd = PERIOD_DATA.workouts[period];
  const filtered = filter==='All' ? workouts : workouts.filter(w => w.type===filter.toLowerCase());

  return (
    <div style={{ maxWidth:760, display:'flex', flexDirection:'column', gap:20 }}>
      <div style={{ display:'flex', alignItems:'center', justifyContent:'space-between' }}>
        <SectionLabel color={acc}>Workouts</SectionLabel>
        <Btn color={acc} onClick={() => openQuickLog('workout')}>+ Add Workout</Btn>
      </div>

      <div style={{ display:'flex', gap:12 }}>
        {[[pd.value, { day:'today', week:'this week', month:'this month', year:'this year' }[period]],
          [MOCK.workouts.thisMonth,'All time total'],
          [`🔥 ${MOCK.workouts.streak}`,'Day streak']].map(([v,l]) => (
          <Card key={String(l)} style={{ flex:1, padding:'16px 20px', textAlign:'center' }}>
            <div style={{ fontSize:24, fontWeight:700, color:C.text, fontFamily:'Space Grotesk,sans-serif' }}>{v}</div>
            <div style={{ fontSize:11, color:C.text2, marginTop:4 }}>{l}</div>
          </Card>
        ))}
      </div>

      <Card style={{ padding:'20px 24px' }}>
        <div style={{ display:'flex', justifyContent:'space-between', alignItems:'center', marginBottom:16 }}>
          <span style={{ fontSize:12, fontWeight:600, color:C.text2, textTransform:'uppercase', letterSpacing:'0.06em' }}>
            Activity — Minutes
          </span>
          <PeriodSelector value={period} onChange={setPeriod} accent={acc} />
        </div>
        <BarChart data={pd.chart} color={acc} height={80} />
      </Card>

      <div style={{ display:'flex', gap:8, flexWrap:'wrap' }}>
        {FILTERS.map(f => (
          <button key={f} onClick={() => setFilter(f)} style={{
            padding:'6px 16px', borderRadius:20, fontSize:12, fontWeight:600, cursor:'pointer',
            background: filter===f ? acc : 'transparent',
            color: filter===f ? '#fff' : C.text2,
            border:`1px solid ${filter===f ? acc : C.border2}`, transition:'all 0.15s'
          }}>{f}</button>
        ))}
      </div>

      <div style={{ display:'flex', flexDirection:'column', gap:12 }}>
        {filtered.length===0
          ? <Card><EmptyState emoji="💪" title="No workouts yet" desc="Start tracking your fitness journey" /></Card>
          : filtered.map(w => <WorkoutCard key={w.id} workout={w} onDelete={() => setWorkouts(ws => ws.filter(x => x.id!==w.id))} />)
        }
      </div>
    </div>
  );
}

Object.assign(window, { Dashboard, Workouts, GoalsSection });
