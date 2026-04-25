// Goals Management + Events/Todo pages

const GOAL_ACC = '#f97316';
const EVENT_ACC = '#06b6d4';

// ── Goals ─────────────────────────────────────────────────────────────────────

const GOAL_CATS = [
  { id:'workout', emoji:'💪', label:'Workout', accent:C.sections.workouts },
  { id:'study',   emoji:'📚', label:'Study',   accent:C.sections.study },
  { id:'sleep',   emoji:'😴', label:'Sleep',   accent:C.sections.sleep },
  { id:'finance', emoji:'💰', label:'Finance', accent:C.sections.finance },
  { id:'custom',  emoji:'🎯', label:'Custom',  accent:C.primary },
];

const INIT_GOALS = [
  { id:1, emoji:'💪', title:'Workouts per week', cat:'workout', current:4, target:5, unit:'sessions', reverse:false },
  { id:2, emoji:'📚', title:'Study hours per week', cat:'study', current:14.5, target:20, unit:'hours', reverse:false },
  { id:3, emoji:'😴', title:'Average sleep', cat:'sleep', current:7.2, target:8, unit:'hours', reverse:false },
  { id:4, emoji:'💰', title:'Monthly spending limit', cat:'finance', current:2840, target:3000, unit:'$', reverse:true },
  { id:5, emoji:'🏃', title:'Weekly running distance', cat:'workout', current:12, target:20, unit:'km', reverse:false },
];

function GoalFormModal({ open, goal, onClose, onSave }) {
  const blank = { emoji:'🎯', title:'', cat:'custom', current:'', target:'', unit:'', reverse:false };
  const [form, setForm] = React.useState(goal || blank);
  React.useEffect(() => { setForm(goal || blank); }, [open]);
  if (!open) return null;
  const set = k => v => setForm(f => ({...f,[k]:v}));
  const catInfo = GOAL_CATS.find(c => c.id === form.cat) || GOAL_CATS[4];
  const inputSt = { width:'100%', background:C.surface3, border:`1px solid ${C.border2}`, borderRadius:9, padding:'9px 12px', fontSize:13, color:C.text, outline:'none', fontFamily:'inherit' };
  const labelSt = { display:'block', fontSize:11, color:C.text2, marginBottom:5, textTransform:'uppercase', letterSpacing:'0.06em' };
  return (
    <div style={{ position:'fixed', inset:0, zIndex:200, display:'flex', alignItems:'center', justifyContent:'center' }}>
      <div onClick={onClose} style={{ position:'absolute', inset:0, background:'rgba(0,0,0,0.75)', backdropFilter:'blur(6px)', WebkitBackdropFilter:'blur(6px)' }} />
      <div style={{ position:'relative', width:'100%', maxWidth:440, margin:'0 16px', background:C.surface, borderRadius:18, border:`1px solid ${C.border2}`, boxShadow:'0 24px 80px rgba(0,0,0,0.6)', padding:'24px' }}>
        <div style={{ display:'flex', justifyContent:'space-between', alignItems:'center', marginBottom:20 }}>
          <span style={{ fontSize:15, fontWeight:700, color:C.text, fontFamily:'Space Grotesk,sans-serif' }}>{goal ? 'Edit Goal' : 'New Goal'}</span>
          <button onClick={onClose} style={{ width:28, height:28, borderRadius:8, background:C.surface2, border:`1px solid ${C.border}`, color:C.text2, cursor:'pointer', fontSize:16, display:'flex', alignItems:'center', justifyContent:'center' }}>×</button>
        </div>

        <div style={{ display:'flex', flexDirection:'column', gap:12 }}>
          <div style={{ display:'grid', gridTemplateColumns:'72px 1fr', gap:12 }}>
            <div>
              <label style={labelSt}>Emoji</label>
              <input value={form.emoji} onChange={e=>set('emoji')(e.target.value)} style={{...inputSt, textAlign:'center', fontSize:20}} maxLength={2} />
            </div>
            <div>
              <label style={labelSt}>Goal Title</label>
              <input placeholder="e.g. Weekly pushups" value={form.title} onChange={e=>set('title')(e.target.value)} style={inputSt} />
            </div>
          </div>
          <div>
            <label style={labelSt}>Category</label>
            <div style={{ display:'flex', gap:6, flexWrap:'wrap' }}>
              {GOAL_CATS.map(c => (
                <button key={c.id} onClick={()=>set('cat')(c.id)} style={{ padding:'6px 12px', borderRadius:20, fontSize:12, fontWeight:600, cursor:'pointer',
                  background: form.cat===c.id ? c.accent : C.surface2,
                  color: form.cat===c.id ? '#fff' : C.text2,
                  border:`1px solid ${form.cat===c.id ? c.accent : C.border2}`, transition:'all 0.15s' }}>
                  {c.emoji} {c.label}
                </button>
              ))}
            </div>
          </div>
          <div style={{ display:'grid', gridTemplateColumns:'1fr 1fr 1fr', gap:12 }}>
            <div>
              <label style={labelSt}>Current</label>
              <input type="number" placeholder="0" value={form.current} onChange={e=>set('current')(e.target.value)} style={inputSt} />
            </div>
            <div>
              <label style={labelSt}>Target</label>
              <input type="number" placeholder="10" value={form.target} onChange={e=>set('target')(e.target.value)} style={inputSt} />
            </div>
            <div>
              <label style={labelSt}>Unit</label>
              <input placeholder="sessions" value={form.unit} onChange={e=>set('unit')(e.target.value)} style={inputSt} />
            </div>
          </div>
          <div style={{ display:'flex', alignItems:'center', gap:10 }}>
            <button onClick={()=>set('reverse')(!form.reverse)} style={{ width:36, height:20, borderRadius:10, background:form.reverse?catInfo.accent:C.surface3, border:`1px solid ${form.reverse?catInfo.accent:C.border2}`, position:'relative', cursor:'pointer', transition:'all 0.2s', flexShrink:0 }}>
              <div style={{ position:'absolute', top:2, left:form.reverse?18:2, width:14, height:14, borderRadius:7, background:'#fff', transition:'left 0.2s' }} />
            </button>
            <span style={{ fontSize:12, color:C.text2 }}>Lower is better (e.g. spending, weight)</span>
          </div>
        </div>

        <div style={{ display:'flex', justifyContent:'flex-end', gap:8, marginTop:20, paddingTop:16, borderTop:`1px solid ${C.border}` }}>
          <Btn outline color={C.text2} size="sm" onClick={onClose}>Cancel</Btn>
          <Btn color={catInfo.accent} onClick={() => { onSave({...form, id: goal?.id || Date.now(), current:+form.current, target:+form.target}); onClose(); }}>
            {goal ? 'Save Changes' : 'Add Goal'}
          </Btn>
        </div>
      </div>
    </div>
  );
}

function GoalCard({ goal, onEdit, onDelete }) {
  const catInfo = GOAL_CATS.find(c => c.id === goal.cat) || GOAL_CATS[4];
  const acc = catInfo.accent;
  const pct = Math.min(100, (goal.current / goal.target) * 100);
  const warn = goal.reverse && pct > 90;
  const barColor = warn ? '#f43f5e' : acc;
  const displayCurrent = goal.unit === '$' ? `$${Number(goal.current).toLocaleString()}` : `${goal.current}${goal.unit ? ' ' + goal.unit : ''}`;
  const displayTarget = goal.unit === '$' ? `$${Number(goal.target).toLocaleString()}` : `${goal.target}${goal.unit ? ' ' + goal.unit : ''}`;

  return (
    <Card accent={acc}>
      <div style={{ padding:'16px 20px' }}>
        <div style={{ display:'flex', alignItems:'flex-start', justifyContent:'space-between', marginBottom:14 }}>
          <div style={{ display:'flex', alignItems:'center', gap:10 }}>
            <div style={{ width:38, height:38, borderRadius:10, background:`${acc}18`, display:'flex', alignItems:'center', justifyContent:'center', fontSize:18, flexShrink:0 }}>{goal.emoji}</div>
            <div>
              <div style={{ fontSize:14, fontWeight:600, color:C.text }}>{goal.title}</div>
              <Badge color={acc}>{catInfo.emoji} {catInfo.label}</Badge>
            </div>
          </div>
          <div style={{ display:'flex', gap:6, flexShrink:0, marginLeft:12 }}>
            <Btn outline color={C.text2} size="sm" onClick={onEdit}>✏️</Btn>
            <Btn outline color="#f43f5e" size="sm" onClick={onDelete}>🗑</Btn>
          </div>
        </div>

        <div style={{ display:'flex', justifyContent:'space-between', alignItems:'center', marginBottom:8 }}>
          <span style={{ fontSize:12, color:C.text2 }}>Progress</span>
          <span style={{ fontSize:12, fontWeight:700, color: warn ? '#f43f5e' : C.text }}>
            {displayCurrent} <span style={{ color:C.text3, fontWeight:400 }}>/ {displayTarget}</span>
          </span>
        </div>
        <div style={{ height:8, borderRadius:8, background:C.surface3, overflow:'hidden', marginBottom:6 }}>
          <div style={{ height:'100%', width:`${pct}%`, borderRadius:8,
            background:`linear-gradient(90deg, ${barColor}70, ${barColor})`,
            transition:'width 0.6s cubic-bezier(0.4,0,0.2,1)' }} />
        </div>
        <div style={{ display:'flex', justifyContent:'space-between' }}>
          <span style={{ fontSize:11, color: warn ? '#f43f5e' : pct >= 100 ? '#10b981' : C.text2 }}>
            {pct >= 100 ? '✓ Goal reached!' : warn ? `⚠️ ${Math.round(pct)}% — over budget` : `${Math.round(pct)}% complete`}
          </span>
          {pct < 100 && <span style={{ fontSize:11, color:C.text3 }}>
            {goal.unit === '$'
              ? `$${(goal.target - goal.current).toLocaleString()} remaining`
              : `${(goal.target - goal.current).toFixed(1).replace(/\.0$/,'')} ${goal.unit} to go`}
          </span>}
        </div>
      </div>
    </Card>
  );
}

function Goals() {
  const [goals, setGoals] = React.useState(INIT_GOALS);
  const [modal, setModal] = React.useState({ open:false, goal:null });
  const completed = goals.filter(g => (g.current / g.target) >= 1).length;

  return (
    <div style={{ maxWidth:760, display:'flex', flexDirection:'column', gap:20 }}>
      <div style={{ display:'flex', alignItems:'center', justifyContent:'space-between' }}>
        <SectionLabel color={GOAL_ACC}>Goals</SectionLabel>
        <Btn color={GOAL_ACC} onClick={() => setModal({open:true,goal:null})}>+ Add Goal</Btn>
      </div>

      <div style={{ display:'flex', gap:12 }}>
        {[[goals.length,'Total goals'],[goals.length - completed,'In progress'],[completed,'Completed']].map(([v,l]) => (
          <Card key={l} style={{ flex:1, padding:'16px 20px', textAlign:'center' }}>
            <div style={{ fontSize:24, fontWeight:700, color:C.text, fontFamily:'Space Grotesk,sans-serif' }}>{v}</div>
            <div style={{ fontSize:11, color:C.text2, marginTop:4 }}>{l}</div>
          </Card>
        ))}
      </div>

      <div style={{ display:'flex', flexDirection:'column', gap:12 }}>
        {goals.length === 0
          ? <Card><EmptyState emoji="🎯" title="No goals yet" desc="Add your first goal to start tracking progress" /></Card>
          : goals.map(g => (
            <GoalCard key={g.id} goal={g}
              onEdit={() => setModal({open:true, goal:g})}
              onDelete={() => setGoals(gs => gs.filter(x => x.id !== g.id))} />
          ))
        }
      </div>

      <GoalFormModal open={modal.open} goal={modal.goal} onClose={() => setModal({open:false,goal:null})}
        onSave={saved => setGoals(gs => modal.goal ? gs.map(g => g.id===saved.id ? saved : g) : [...gs, saved])} />
    </div>
  );
}

// ── Events / Todo ─────────────────────────────────────────────────────────────

const EV_CATS = [
  { id:'workout', emoji:'💪', label:'Workout', accent:C.sections.workouts },
  { id:'study',   emoji:'📚', label:'Study',   accent:C.sections.study },
  { id:'sleep',   emoji:'😴', label:'Sleep',   accent:C.sections.sleep },
  { id:'finance', emoji:'💰', label:'Finance', accent:C.sections.finance },
  { id:'custom',  emoji:'📌', label:'General', accent:EVENT_ACC },
];

const INIT_EVENTS_TODO = [
  { id:1, done:true,  title:'Morning meditation',    date:'2026-04-25', time:'08:00', cat:'custom' },
  { id:2, done:true,  title:'Grocery shopping',      date:'2026-04-25', time:'',      cat:'custom' },
  { id:3, done:false, title:'Morning Run',            date:'2026-04-26', time:'07:00', cat:'workout' },
  { id:4, done:false, title:'Physics Exam Prep',      date:'2026-04-26', time:'14:00', cat:'study' },
  { id:5, done:false, title:'Team Meeting',           date:'2026-04-27', time:'10:00', cat:'custom' },
  { id:6, done:false, title:'Pay Rent',               date:'2026-04-28', time:'09:00', cat:'finance' },
  { id:7, done:false, title:'Log sleep data',         date:'2026-04-29', time:'22:00', cat:'sleep' },
  { id:8, done:false, title:'Review weekly goals',    date:'2026-04-30', time:'18:00', cat:'custom' },
];

function EventFormModal({ open, onClose, onSave }) {
  const today = new Date().toISOString().split('T')[0];
  const [form, setForm] = React.useState({ title:'', date:today, time:'', cat:'custom' });
  React.useEffect(() => { if (open) setForm({ title:'', date:today, time:'', cat:'custom' }); }, [open]);
  if (!open) return null;
  const set = k => v => setForm(f => ({...f,[k]:v}));
  const catInfo = EV_CATS.find(c => c.id === form.cat) || EV_CATS[4];
  const inputSt = { width:'100%', background:C.surface3, border:`1px solid ${C.border2}`, borderRadius:9, padding:'9px 12px', fontSize:13, color:C.text, outline:'none', fontFamily:'inherit' };
  const labelSt = { display:'block', fontSize:11, color:C.text2, marginBottom:5, textTransform:'uppercase', letterSpacing:'0.06em' };
  return (
    <div style={{ position:'fixed', inset:0, zIndex:200, display:'flex', alignItems:'center', justifyContent:'center' }}>
      <div onClick={onClose} style={{ position:'absolute', inset:0, background:'rgba(0,0,0,0.75)', backdropFilter:'blur(6px)', WebkitBackdropFilter:'blur(6px)' }} />
      <div style={{ position:'relative', width:'100%', maxWidth:420, margin:'0 16px', background:C.surface, borderRadius:18, border:`1px solid ${C.border2}`, boxShadow:'0 24px 80px rgba(0,0,0,0.6)', padding:'24px' }}>
        <div style={{ display:'flex', justifyContent:'space-between', alignItems:'center', marginBottom:20 }}>
          <span style={{ fontSize:15, fontWeight:700, color:C.text, fontFamily:'Space Grotesk,sans-serif' }}>New Event</span>
          <button onClick={onClose} style={{ width:28, height:28, borderRadius:8, background:C.surface2, border:`1px solid ${C.border}`, color:C.text2, cursor:'pointer', fontSize:16, display:'flex', alignItems:'center', justifyContent:'center' }}>×</button>
        </div>
        <div style={{ display:'flex', flexDirection:'column', gap:12 }}>
          <div><label style={labelSt}>Title</label><input placeholder="What needs to be done?" value={form.title} onChange={e=>set('title')(e.target.value)} style={inputSt} /></div>
          <div style={{ display:'grid', gridTemplateColumns:'1fr 1fr', gap:12 }}>
            <div><label style={labelSt}>Date</label><input type="date" value={form.date} onChange={e=>set('date')(e.target.value)} style={inputSt} /></div>
            <div><label style={labelSt}>Time</label><input type="time" value={form.time} onChange={e=>set('time')(e.target.value)} style={inputSt} /></div>
          </div>
          <div>
            <label style={labelSt}>Category</label>
            <div style={{ display:'flex', gap:6, flexWrap:'wrap' }}>
              {EV_CATS.map(c => (
                <button key={c.id} onClick={()=>set('cat')(c.id)} style={{ padding:'5px 10px', borderRadius:16, fontSize:11, fontWeight:600, cursor:'pointer',
                  background: form.cat===c.id ? c.accent : C.surface2,
                  color: form.cat===c.id ? '#fff' : C.text2,
                  border:`1px solid ${form.cat===c.id ? c.accent : C.border2}`, transition:'all 0.15s' }}>
                  {c.emoji} {c.label}
                </button>
              ))}
            </div>
          </div>
        </div>
        <div style={{ display:'flex', justifyContent:'flex-end', gap:8, marginTop:20, paddingTop:16, borderTop:`1px solid ${C.border}` }}>
          <Btn outline color={C.text2} size="sm" onClick={onClose}>Cancel</Btn>
          <Btn color={catInfo.accent} onClick={() => { if(form.title) { onSave({...form, id:Date.now(), done:false}); onClose(); } }}>Add Event</Btn>
        </div>
      </div>
    </div>
  );
}

function TodoItem({ ev, onToggle, onDelete }) {
  const catInfo = EV_CATS.find(c => c.id === ev.cat) || EV_CATS[4];
  const acc = catInfo.accent;
  const [hov, setHov] = React.useState(false);
  const today = new Date().toISOString().split('T')[0];
  const tomorrow = new Date(Date.now()+86400000).toISOString().split('T')[0];
  function dateStr(d) {
    if (d === today) return 'Today';
    if (d === tomorrow) return 'Tomorrow';
    return new Date(d+'T12:00').toLocaleDateString('en-US',{month:'short',day:'numeric'});
  }
  const timeStr = ev.time
    ? new Date(`${ev.date}T${ev.time}`).toLocaleTimeString('en-US',{hour:'numeric',minute:'2-digit'})
    : null;
  const subtitle = [dateStr(ev.date), timeStr].filter(Boolean).join(' · ');
  return (
    <div onMouseEnter={()=>setHov(true)} onMouseLeave={()=>setHov(false)}
      style={{ display:'flex', alignItems:'center', gap:12, padding:'11px 0', borderBottom:`1px solid ${C.border}` }}>
      <button onClick={() => onToggle(ev.id)} style={{ width:22, height:22, borderRadius:'50%', flexShrink:0, cursor:'pointer', transition:'all 0.2s',
        background: ev.done ? acc : 'transparent',
        border: `2px solid ${ev.done ? acc : C.text3}`,
        display:'flex', alignItems:'center', justifyContent:'center', color:'#fff', fontSize:11 }}>
        {ev.done ? '✓' : ''}
      </button>
      <div style={{ flex:1, minWidth:0 }}>
        <div style={{ fontSize:13, fontWeight:ev.done ? 400 : 500, color: ev.done ? C.text3 : C.text,
          textDecoration: ev.done ? 'line-through' : 'none', transition:'all 0.2s' }}>{ev.title}</div>
        <div style={{ fontSize:11, color:C.text3, marginTop:3 }}>{subtitle}</div>
      </div>
      <div style={{ display:'flex', alignItems:'center', gap:8, flexShrink:0 }}>
        <Badge color={acc}>{catInfo.emoji} {catInfo.label}</Badge>
        {hov && (
          <button onClick={() => onDelete(ev.id)} style={{ color:C.text3, fontSize:16, cursor:'pointer', background:'none', border:'none', padding:'0 2px', lineHeight:1 }}
            onMouseEnter={e=>e.target.style.color='#f43f5e'} onMouseLeave={e=>e.target.style.color=C.text3}>×</button>
        )}
      </div>
    </div>
  );
}

const EV_FILTERS = ['All','Today','Upcoming','Completed'];

function Events() {
  const [events, setEvents] = React.useState(INIT_EVENTS_TODO);
  const [filter, setFilter] = React.useState('All');
  const [showForm, setShowForm] = React.useState(false);
  const today = new Date().toISOString().split('T')[0];
  const tomorrow = new Date(Date.now()+86400000).toISOString().split('T')[0];

  function onToggle(id) { setEvents(es => es.map(e => e.id===id ? {...e,done:!e.done} : e)); }
  function onDelete(id) { setEvents(es => es.filter(e => e.id!==id)); }
  function onAdd(ev) { setEvents(es => [...es, ev]); }

  const filtered = events.filter(e => {
    if (filter==='Today') return e.date===today && !e.done;
    if (filter==='Upcoming') return e.date>today && !e.done;
    if (filter==='Completed') return e.done;
    return true;
  });

  // Group by date
  const groups = filtered.reduce((acc, ev) => {
    const key = ev.date;
    if (!acc[key]) acc[key] = [];
    acc[key].push(ev);
    return acc;
  }, {});
  const sortedDates = Object.keys(groups).sort();

  function dateLabel(d) {
    if (d===today) return 'Today';
    if (d===tomorrow) return 'Tomorrow';
    return new Date(d+'T12:00').toLocaleDateString('en-US',{weekday:'long',month:'short',day:'numeric'});
  }

  const total = events.length;
  const done = events.filter(e=>e.done).length;
  const pct = total ? Math.round((done/total)*100) : 0;

  return (
    <div style={{ maxWidth:680, display:'flex', flexDirection:'column', gap:20 }}>
      <div style={{ display:'flex', alignItems:'center', justifyContent:'space-between' }}>
        <SectionLabel color={EVENT_ACC}>Events</SectionLabel>
        <Btn color={EVENT_ACC} onClick={() => setShowForm(true)}>+ Add Event</Btn>
      </div>

      {/* Progress summary */}
      <Card style={{ padding:'18px 24px' }}>
        <div style={{ display:'flex', justifyContent:'space-between', alignItems:'center', marginBottom:10 }}>
          <span style={{ fontSize:13, fontWeight:600, color:C.text }}>{done} of {total} complete</span>
          <span style={{ fontSize:13, fontWeight:700, color:EVENT_ACC }}>{pct}%</span>
        </div>
        <div style={{ height:6, borderRadius:6, background:C.surface3, overflow:'hidden' }}>
          <div style={{ height:'100%', width:`${pct}%`, borderRadius:6,
            background:`linear-gradient(90deg, ${EVENT_ACC}80, ${EVENT_ACC})`,
            transition:'width 0.5s ease' }} />
        </div>
        <div style={{ display:'flex', gap:16, marginTop:12 }}>
          {[['📌', events.filter(e=>e.date===today&&!e.done).length, 'due today'],
            ['⏳', events.filter(e=>e.date>today&&!e.done).length, 'upcoming'],
            ['✅', done, 'completed']].map(([icon,val,label]) => (
            <div key={label} style={{ display:'flex', alignItems:'center', gap:6 }}>
              <span style={{ fontSize:14 }}>{icon}</span>
              <span style={{ fontSize:13, fontWeight:700, color:C.text }}>{val}</span>
              <span style={{ fontSize:11, color:C.text2 }}>{label}</span>
            </div>
          ))}
        </div>
      </Card>

      {/* Filters */}
      <div style={{ display:'flex', gap:8 }}>
        {EV_FILTERS.map(f => (
          <button key={f} onClick={() => setFilter(f)} style={{ padding:'6px 16px', borderRadius:20, fontSize:12, fontWeight:600, cursor:'pointer',
            background: filter===f ? EVENT_ACC : 'transparent',
            color: filter===f ? '#fff' : C.text2,
            border:`1px solid ${filter===f ? EVENT_ACC : C.border2}`, transition:'all 0.15s' }}>{f}</button>
        ))}
      </div>

      {/* Grouped list */}
      {sortedDates.length === 0
        ? <Card><EmptyState emoji="📌" title="Nothing here" desc="All done — or try a different filter" /></Card>
        : sortedDates.map(date => (
          <Card key={date} style={{ padding:'0 20px' }}>
            <div style={{ padding:'14px 0 10px', borderBottom:`1px solid ${C.border}`, display:'flex', alignItems:'center', gap:10 }}>
              <span style={{ fontSize:12, fontWeight:700, color:EVENT_ACC, textTransform:'uppercase', letterSpacing:'0.06em' }}>{dateLabel(date)}</span>
              <span style={{ fontSize:11, color:C.text3 }}>· {groups[date].length} item{groups[date].length!==1?'s':''}</span>
            </div>
            {groups[date].map(ev => (
              <TodoItem key={ev.id} ev={ev} onToggle={onToggle} onDelete={onDelete} />
            ))}
          </Card>
        ))
      }

      <EventFormModal open={showForm} onClose={() => setShowForm(false)} onSave={onAdd} />
    </div>
  );
}

Object.assign(window, { Goals, Events });
