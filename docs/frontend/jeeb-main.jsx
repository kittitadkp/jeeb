// Calendar + Settings + App Shell + QuickLogModal + NotificationsPanel

// ── Quick Log Modal ───────────────────────────────────────────────────────────

const QL_TYPES = [
  { id:'workout', emoji:'💪', label:'Workout', accent: C.sections.workouts },
  { id:'study',   emoji:'📚', label:'Study',   accent: C.sections.study },
  { id:'sleep',   emoji:'😴', label:'Sleep',   accent: C.sections.sleep },
  { id:'finance', emoji:'💰', label:'Finance', accent: C.sections.finance },
];

function QLForm({ type, onClose }) {
  const info = QL_TYPES.find(t => t.id === type);
  const acc = info?.accent || C.primary;
  const today = new Date().toISOString().split('T')[0];
  const yesterday = new Date(Date.now()-86400000).toISOString().split('T')[0];

  // Workout form
  const [wForm, setWForm] = React.useState({ date:today, type:'strength', duration:30, notes:'' });
  // Study form
  const [sForm, setSForm] = React.useState({ date:today, subject:'', duration:60, notes:'' });
  // Sleep form
  const [slForm, setSlForm] = React.useState({ date:today, bedtime:'23:00', wake:'07:00', quality:4 });
  // Finance form
  const [fForm, setFForm] = React.useState({ type:'expense', amount:'', category:'', date:today, notes:'' });

  const inputStyle = { width:'100%', background:C.surface3, border:`1px solid ${C.border2}`,
    borderRadius:9, padding:'9px 12px', fontSize:13, color:C.text, outline:'none', fontFamily:'inherit' };
  const labelStyle = { display:'block', fontSize:11, color:C.text2, marginBottom:5,
    textTransform:'uppercase', letterSpacing:'0.06em' };
  const grid2 = { display:'grid', gridTemplateColumns:'1fr 1fr', gap:12 };
  const grid3 = { display:'grid', gridTemplateColumns:'1fr 1fr 1fr', gap:12 };

  function Field({ label, children }) {
    return <div><label style={labelStyle}>{label}</label>{children}</div>;
  }

  // Date quick-pick chips
  function DateChips({ value, onChange }) {
    const opts = [
      { label:'Today', val: today },
      { label:'Yesterday', val: yesterday },
      { label:'Custom', val: 'custom' },
    ];
    const isCustom = value !== today && value !== yesterday;
    return (
      <div>
        <label style={labelStyle}>Date</label>
        <div style={{ display:'flex', gap:6, marginBottom: isCustom ? 8 : 0 }}>
          {opts.map(o => (
            <button key={o.val} onClick={() => onChange(o.val === 'custom' ? yesterday : o.val)}
              type="button" style={{ flex:1, padding:'7px 0', borderRadius:8, fontSize:12, fontWeight:600, cursor:'pointer',
                background: (o.val==='custom' ? isCustom : value===o.val) ? acc : C.surface3,
                color: (o.val==='custom' ? isCustom : value===o.val) ? '#fff' : C.text2,
                border:`1px solid ${(o.val==='custom' ? isCustom : value===o.val) ? acc : C.border2}`,
                transition:'all 0.15s' }}>{o.label}</button>
          ))}
        </div>
        {isCustom && (
          <input type="date" value={value} onChange={e => onChange(e.target.value)} style={inputStyle} />
        )}
      </div>
    );
  }

  return (
    <div>
      <div style={{ display:'flex', alignItems:'center', gap:10, marginBottom:20, paddingBottom:16, borderBottom:`1px solid ${C.border}` }}>
        <div style={{ width:36, height:36, borderRadius:10, background:`${acc}20`,
          display:'flex', alignItems:'center', justifyContent:'center', fontSize:18 }}>{info?.emoji}</div>
        <span style={{ fontSize:15, fontWeight:700, color:C.text }}>Log {info?.label}</span>
      </div>

      {type === 'workout' && (
        <div style={{ display:'flex', flexDirection:'column', gap:12 }}>
          <DateChips value={wForm.date} onChange={v=>setWForm(f=>({...f,date:v}))} />
          <div style={grid2}>
            <Field label="Type">
              <select value={wForm.type} onChange={e=>setWForm(f=>({...f,type:e.target.value}))} style={inputStyle}>
                <option value="strength">Strength</option>
                <option value="cardio">Cardio</option>
                <option value="flexibility">Flexibility</option>
              </select>
            </Field>
            <Field label="Duration (min)">
              <input type="number" min={1} value={wForm.duration}
                onChange={e=>setWForm(f=>({...f,duration:+e.target.value}))} style={inputStyle} />
            </Field>
          </div>
          <Field label="Notes">
            <input placeholder="Optional" value={wForm.notes}
              onChange={e=>setWForm(f=>({...f,notes:e.target.value}))} style={inputStyle} />
          </Field>
        </div>
      )}

      {type === 'study' && (
        <div style={{ display:'flex', flexDirection:'column', gap:12 }}>
          <DateChips value={sForm.date} onChange={v=>setSForm(f=>({...f,date:v}))} />
          <div style={grid2}>
            <Field label="Subject">
              <input placeholder="Mathematics" value={sForm.subject}
                onChange={e=>setSForm(f=>({...f,subject:e.target.value}))} style={inputStyle} />
            </Field>
            <Field label="Duration (min)">
              <input type="number" min={5} value={sForm.duration}
                onChange={e=>setSForm(f=>({...f,duration:+e.target.value}))} style={inputStyle} />
            </Field>
          </div>
          <Field label="Notes">
            <input placeholder="Optional" value={sForm.notes}
              onChange={e=>setSForm(f=>({...f,notes:e.target.value}))} style={inputStyle} />
          </Field>
        </div>
      )}

      {type === 'sleep' && (
        <div style={{ display:'flex', flexDirection:'column', gap:12 }}>
          <DateChips value={slForm.date} onChange={v=>setSlForm(f=>({...f,date:v}))} />
          <div style={grid2}>
            <Field label="Bedtime">
              <input type="time" value={slForm.bedtime}
                onChange={e=>setSlForm(f=>({...f,bedtime:e.target.value}))} style={inputStyle} />
            </Field>
            <Field label="Wake Up">
              <input type="time" value={slForm.wake}
                onChange={e=>setSlForm(f=>({...f,wake:e.target.value}))} style={inputStyle} />
            </Field>
          </div>
          <Field label="Quality">
            <div style={{ display:'flex', gap:8 }}>
              {[1,2,3,4,5].map(n => (
                <button key={n} onClick={()=>setSlForm(f=>({...f,quality:n}))} style={{
                  width:42, height:42, borderRadius:9, fontSize:13, fontWeight:600, cursor:'pointer',
                  background: slForm.quality===n ? acc : C.surface3,
                  color: slForm.quality===n ? '#fff' : C.text2,
                  border:`1px solid ${slForm.quality===n ? acc : C.border2}`, transition:'all 0.15s'
                }}>{n}</button>
              ))}
            </div>
          </Field>
        </div>
      )}

      {type === 'finance' && (
        <div style={{ display:'flex', flexDirection:'column', gap:12 }}>
          <div style={grid2}>
            <Field label="Type">
              <select value={fForm.type} onChange={e=>setFForm(f=>({...f,type:e.target.value}))} style={inputStyle}>
                <option value="expense">Expense</option>
                <option value="income">Income</option>
              </select>
            </Field>
            <Field label="Amount ($)">
              <input type="number" min="0.01" step="0.01" placeholder="0.00" value={fForm.amount}
                onChange={e=>setFForm(f=>({...f,amount:e.target.value}))} style={inputStyle} />
            </Field>
          </div>
          <Field label="Category">
            <input placeholder="Food & Dining" value={fForm.category}
              onChange={e=>setFForm(f=>({...f,category:e.target.value}))} style={inputStyle} />
          </Field>
          <DateChips value={fForm.date} onChange={v=>setFForm(f=>({...f,date:v}))} />
          <Field label="Notes">
            <input placeholder="Optional" value={fForm.notes}
              onChange={e=>setFForm(f=>({...f,notes:e.target.value}))} style={inputStyle} />
          </Field>
        </div>
      )}

      <div style={{ display:'flex', justifyContent:'flex-end', gap:8, marginTop:20, paddingTop:16, borderTop:`1px solid ${C.border}` }}>
        <Btn outline color={C.text2} size="sm" onClick={onClose}>Cancel</Btn>
        <Btn color={acc} onClick={onClose}>Save {info?.label}</Btn>
      </div>
    </div>
  );
}

function QuickLogModal({ open, initialType, onClose }) {
  const [type, setType] = React.useState(initialType || null);
  React.useEffect(() => { setType(initialType || null); }, [initialType, open]);
  if (!open) return null;
  return (
    <div style={{ position:'fixed', inset:0, zIndex:200, display:'flex', alignItems:'center', justifyContent:'center' }}>
      <div onClick={onClose} style={{ position:'absolute', inset:0, background:'rgba(0,0,0,0.75)',
        backdropFilter:'blur(6px)', WebkitBackdropFilter:'blur(6px)' }} />
      <div style={{ position:'relative', width:'100%', maxWidth:460, margin:'0 16px',
        background:C.surface, borderRadius:18, border:`1px solid ${C.border2}`,
        boxShadow:'0 24px 80px rgba(0,0,0,0.6)', overflow:'hidden',
        animation:'qlIn 0.2s cubic-bezier(0.34,1.56,0.64,1)' }}>
        <style>{`@keyframes qlIn{from{opacity:0;transform:scale(0.94) translateY(8px)}to{opacity:1;transform:scale(1) translateY(0)}}`}</style>
        <div style={{ padding:'20px 24px' }}>
          {!type ? (
            <>
              <div style={{ display:'flex', justifyContent:'space-between', alignItems:'center', marginBottom:20 }}>
                <span style={{ fontSize:16, fontWeight:700, color:C.text, fontFamily:'Space Grotesk,sans-serif' }}>Quick Log</span>
                <button onClick={onClose} style={{ width:28, height:28, borderRadius:8, background:C.surface2,
                  border:`1px solid ${C.border}`, color:C.text2, cursor:'pointer', fontSize:16, display:'flex',
                  alignItems:'center', justifyContent:'center' }}>×</button>
              </div>
              <div style={{ display:'grid', gridTemplateColumns:'1fr 1fr', gap:10 }}>
                {QL_TYPES.map(t => (
                  <button key={t.id} onClick={()=>setType(t.id)} style={{
                    padding:'18px 16px', borderRadius:12, border:`1px solid ${C.border2}`,
                    background:C.surface2, cursor:'pointer', textAlign:'left', transition:'all 0.15s' }}
                    onMouseEnter={e=>{e.currentTarget.style.borderColor=t.accent;e.currentTarget.style.background=`${t.accent}12`;}}
                    onMouseLeave={e=>{e.currentTarget.style.borderColor=C.border2;e.currentTarget.style.background=C.surface2;}}>
                    <div style={{ fontSize:24, marginBottom:8 }}>{t.emoji}</div>
                    <div style={{ fontSize:13, fontWeight:600, color:C.text }}>{t.label}</div>
                  </button>
                ))}
              </div>
            </>
          ) : (
            <>
              <div style={{ display:'flex', justifyContent:'space-between', alignItems:'center', marginBottom:4 }}>
                {!initialType && (
                  <button onClick={()=>setType(null)} style={{ fontSize:12, color:C.text2, cursor:'pointer',
                    background:'none', border:'none', display:'flex', alignItems:'center', gap:4 }}>‹ Back</button>
                )}
                <div style={{ flex:1 }} />
                <button onClick={onClose} style={{ width:28, height:28, borderRadius:8, background:C.surface2,
                  border:`1px solid ${C.border}`, color:C.text2, cursor:'pointer', fontSize:16, display:'flex',
                  alignItems:'center', justifyContent:'center' }}>×</button>
              </div>
              <QLForm type={type} onClose={onClose} />
            </>
          )}
        </div>
      </div>
    </div>
  );
}

// ── Notifications Panel ───────────────────────────────────────────────────────

function NotificationsPanel({ open, onClose }) {
  const [notifs, setNotifs] = React.useState(NOTIFICATIONS);
  const unreadCount = notifs.filter(n => n.unread).length;
  return (
    <>
      {open && <div onClick={onClose} style={{ position:'fixed', inset:0, zIndex:150 }} />}
      <div style={{
        position:'fixed', top:56, right:0, bottom:0, width:340, zIndex:160,
        background:C.surface, borderLeft:`1px solid ${C.border}`,
        transform: open ? 'translateX(0)' : 'translateX(100%)',
        transition:'transform 0.25s cubic-bezier(0.4,0,0.2,1)',
        display:'flex', flexDirection:'column',
        boxShadow: open ? '-16px 0 48px rgba(0,0,0,0.4)' : 'none',
      }}>
        <div style={{ padding:'18px 20px', borderBottom:`1px solid ${C.border}`,
          display:'flex', alignItems:'center', justifyContent:'space-between', flexShrink:0 }}>
          <div style={{ display:'flex', alignItems:'center', gap:8 }}>
            <span style={{ fontSize:14, fontWeight:700, color:C.text, fontFamily:'Space Grotesk,sans-serif' }}>Notifications</span>
            {unreadCount > 0 && (
              <span style={{ background:C.primary, color:'#fff', borderRadius:10, padding:'1px 7px', fontSize:11, fontWeight:700 }}>{unreadCount}</span>
            )}
          </div>
          <div style={{ display:'flex', gap:8 }}>
            {unreadCount > 0 && (
              <button onClick={()=>setNotifs(ns=>ns.map(n=>({...n,unread:false})))} style={{ fontSize:11, color:C.primary,
                cursor:'pointer', background:'none', border:'none', fontWeight:600 }}>Mark all read</button>
            )}
            <button onClick={onClose} style={{ width:26, height:26, borderRadius:7, background:C.surface2,
              border:`1px solid ${C.border}`, color:C.text2, cursor:'pointer', fontSize:15,
              display:'flex', alignItems:'center', justifyContent:'center' }}>×</button>
          </div>
        </div>

        <div style={{ flex:1, overflowY:'auto' }}>
          {notifs.map(n => (
            <div key={n.id} onClick={()=>setNotifs(ns=>ns.map(x=>x.id===n.id?{...x,unread:false}:x))}
              style={{ display:'flex', gap:12, padding:'14px 20px', borderBottom:`1px solid ${C.border}`,
                background: n.unread ? `${n.acc}08` : 'transparent',
                cursor:'pointer', transition:'background 0.15s',
                position:'relative' }}>
              {n.unread && <div style={{ position:'absolute', left:8, top:'50%', transform:'translateY(-50%)',
                width:4, height:4, borderRadius:'50%', background:n.acc }} />}
              <div style={{ width:36, height:36, borderRadius:10, background:`${n.acc}18`, flexShrink:0,
                display:'flex', alignItems:'center', justifyContent:'center', fontSize:16 }}>{n.emoji}</div>
              <div style={{ flex:1, minWidth:0 }}>
                <div style={{ fontSize:13, fontWeight:600, color:C.text, marginBottom:2 }}>{n.title}</div>
                <div style={{ fontSize:11, color:C.text2, lineHeight:1.4 }}>{n.desc}</div>
                <div style={{ fontSize:10, color:C.text3, marginTop:5 }}>{n.time}</div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </>
  );
}

// ── Calendar ──────────────────────────────────────────────────────────────────

function Calendar() {
  const acc = C.sections.calendar;
  const today = new Date();
  const [current, setCurrent] = React.useState(new Date(today.getFullYear(), today.getMonth(), 1));
  const year = current.getFullYear(), month = current.getMonth();
  const monthName = current.toLocaleDateString('en-US',{month:'long',year:'numeric'});
  const firstDay = new Date(year,month,1).getDay();
  const daysInMonth = new Date(year,month+1,0).getDate();
  const cells = Array.from({length:firstDay+daysInMonth},(_,i)=>i-firstDay+1>0?i-firstDay+1:null);
  const eventAccent = t => ({workout:C.sections.workouts,study:C.sections.study,sleep:C.sections.sleep,finance:C.sections.finance,custom:acc})[t]||acc;
  const eventEmoji = t => ({workout:'💪',study:'📚',sleep:'😴',finance:'💰',custom:'📅'})[t]||'📅';
  const eventDays = MOCK.events.map(e=>new Date(e.start).getDate());
  const upcoming = MOCK.events.filter(e=>new Date(e.start)>new Date()).sort((a,b)=>new Date(a.start)-new Date(b.start));

  return (
    <div style={{ maxWidth:860, display:'flex', flexDirection:'column', gap:20 }}>
      <SectionLabel color={acc}>Calendar</SectionLabel>
      <div style={{ display:'grid', gridTemplateColumns:'1fr 320px', gap:16 }}>
        <Card style={{ padding:24 }}>
          <div style={{ display:'flex', alignItems:'center', justifyContent:'space-between', marginBottom:20 }}>
            <button onClick={()=>setCurrent(d=>new Date(d.getFullYear(),d.getMonth()-1,1))} style={{ width:32,height:32,borderRadius:8,background:C.surface2,border:`1px solid ${C.border2}`,color:C.text2,fontSize:14,cursor:'pointer',display:'flex',alignItems:'center',justifyContent:'center' }}>‹</button>
            <span style={{ fontWeight:700, fontSize:16, color:C.text, fontFamily:'Space Grotesk,sans-serif' }}>{monthName}</span>
            <button onClick={()=>setCurrent(d=>new Date(d.getFullYear(),d.getMonth()+1,1))} style={{ width:32,height:32,borderRadius:8,background:C.surface2,border:`1px solid ${C.border2}`,color:C.text2,fontSize:14,cursor:'pointer',display:'flex',alignItems:'center',justifyContent:'center' }}>›</button>
          </div>
          <div style={{ display:'grid', gridTemplateColumns:'repeat(7,1fr)', gap:2, marginBottom:8 }}>
            {['Su','Mo','Tu','We','Th','Fr','Sa'].map(d=>(
              <div key={d} style={{ textAlign:'center', fontSize:11, color:C.text2, fontWeight:600, padding:'4px 0', textTransform:'uppercase', letterSpacing:'0.05em' }}>{d}</div>
            ))}
          </div>
          <div style={{ display:'grid', gridTemplateColumns:'repeat(7,1fr)', gap:3 }}>
            {cells.map((day,i)=>{
              const isToday=day===today.getDate()&&month===today.getMonth()&&year===today.getFullYear();
              const hasEvent=day&&eventDays.includes(day);
              return (
                <div key={i} style={{ aspectRatio:'1',display:'flex',flexDirection:'column',alignItems:'center',justifyContent:'center',borderRadius:9,cursor:day?'pointer':'default',background:isToday?acc:'transparent',border:`1px solid ${isToday?acc:day?C.border:'transparent'}`,position:'relative' }}>
                  {day&&(<>
                    <span style={{ fontSize:13, fontWeight:isToday?700:400, color:isToday?'#fff':C.text }}>{day}</span>
                    {hasEvent&&!isToday&&<div style={{ width:4,height:4,borderRadius:'50%',background:acc,position:'absolute',bottom:4 }} />}
                  </>)}
                </div>
              );
            })}
          </div>
        </Card>

        <Card style={{ padding:0, display:'flex', flexDirection:'column' }}>
          <div style={{ padding:'18px 20px 12px', borderBottom:`1px solid ${C.border}` }}>
            <span style={{ fontSize:13, fontWeight:600, color:C.text }}>Upcoming Events</span>
          </div>
          <div style={{ flex:1, padding:'8px 20px', overflowY:'auto' }}>
            {upcoming.map(e=>{
              const d=new Date(e.start); const ea=eventAccent(e.type);
              const isToday=d.toDateString()===today.toDateString();
              const isTom=d.toDateString()===new Date(today.getTime()+86400000).toDateString();
              const dayLabel=isToday?'Today':isTom?'Tomorrow':d.toLocaleDateString('en-US',{month:'short',day:'numeric'});
              return (
                <div key={e.id} style={{ display:'flex',gap:12,padding:'12px 0',borderBottom:`1px solid ${C.border}` }}>
                  <div style={{ width:38,height:38,borderRadius:10,background:`${ea}18`,display:'flex',alignItems:'center',justifyContent:'center',fontSize:17,flexShrink:0 }}>{eventEmoji(e.type)}</div>
                  <div style={{ flex:1 }}>
                    <div style={{ fontSize:13,fontWeight:600,color:C.text }}>{e.title}</div>
                    <div style={{ fontSize:11,color:C.text2,marginTop:3 }}>{dayLabel} · {d.toLocaleTimeString('en-US',{hour:'numeric',minute:'2-digit'})}</div>
                  </div>
                  <div style={{ marginLeft:'auto',width:3,borderRadius:2,background:ea,flexShrink:0,alignSelf:'stretch' }} />
                </div>
              );
            })}
          </div>
          
        </Card>
      </div>
    </div>
  );
}

// ── Settings ──────────────────────────────────────────────────────────────────

function Settings() {
  const [name, setName] = React.useState('Alex');
  const [email, setEmail] = React.useState('alex@example.com');
  const [theme, setTheme] = React.useState('dark');
  const [accent, setAccent] = React.useState('#7c6ef5');
  const [notifs, setNotifs] = React.useState({
    workout: { on:true, time:'07:00' },
    study:   { on:true, time:'09:00' },
    sleep:   { on:true, time:'22:30' },
    budget:  { on:false, threshold:'80' },
  });
  const [goalDefs, setGoalDefs] = React.useState({ workouts:5, study:20, sleep:8, budget:3000 });

  const ACCENTS = ['#7c6ef5','#f43f5e','#f59e0b','#10b981','#14b8a6','#06b6d4'];
  const inputSt = { background:C.surface2, border:`1px solid ${C.border2}`, borderRadius:8, padding:'7px 12px', fontSize:13, color:C.text, outline:'none', fontFamily:'inherit' };

  function Toggle({ on, onChange, color }) {
    const col = color || C.primary;
    return (
      <button onClick={() => onChange(!on)} style={{ width:44, height:24, borderRadius:12, background:on?col:C.surface3,
        border:`1px solid ${on?col:C.border2}`, position:'relative', cursor:'pointer', transition:'all 0.2s', flexShrink:0 }}>
        <div style={{ position:'absolute', top:3, left:on?22:3, width:16, height:16, borderRadius:8, background:'#fff', transition:'left 0.2s' }} />
      </button>
    );
  }

  function Stepper({ value, onChange, min=1, max=100 }) {
    return (
      <div style={{ display:'flex', alignItems:'center', gap:0, background:C.surface2, borderRadius:8, border:`1px solid ${C.border2}`, overflow:'hidden' }}>
        <button onClick={()=>onChange(Math.max(min,value-1))} style={{ width:32, height:34, color:C.text2, fontSize:16, cursor:'pointer', background:'none', border:'none', borderRight:`1px solid ${C.border2}` }}>−</button>
        <span style={{ width:40, textAlign:'center', fontSize:13, fontWeight:600, color:C.text }}>{value}</span>
        <button onClick={()=>onChange(Math.min(max,value+1))} style={{ width:32, height:34, color:C.text2, fontSize:16, cursor:'pointer', background:'none', border:'none', borderLeft:`1px solid ${C.border2}` }}>+</button>
      </div>
    );
  }

  function SectionHead({ label }) {
    return <div style={{ padding:'16px 0 10px', borderBottom:`1px solid ${C.border}`, marginBottom:4 }}>
      <span style={{ fontSize:11, fontWeight:700, color:C.text2, textTransform:'uppercase', letterSpacing:'0.08em' }}>{label}</span>
    </div>;
  }

  function Row({ label, desc, children }) {
    return (
      <div style={{ display:'flex', alignItems:'center', justifyContent:'space-between', padding:'14px 0', borderBottom:`1px solid ${C.border}` }}>
        <div style={{ flex:1, marginRight:16 }}>
          <div style={{ fontSize:13, fontWeight:500, color:C.text }}>{label}</div>
          {desc && <div style={{ fontSize:11, color:C.text2, marginTop:2 }}>{desc}</div>}
        </div>
        {children}
      </div>
    );
  }

  return (
    <div style={{ maxWidth:640, display:'flex', flexDirection:'column', gap:20 }}>
      <SectionLabel color={C.sections.settings}>Settings</SectionLabel>

      {/* Profile */}
      <Card style={{ padding:'0 24px' }}>
        <SectionHead label="Profile" />
        <div style={{ padding:'16px 0 8px', display:'flex', alignItems:'center', gap:16 }}>
          <div style={{ width:56, height:56, borderRadius:16, background:`linear-gradient(135deg, ${C.primary}, #a78bfa)`,
            display:'flex', alignItems:'center', justifyContent:'center', fontSize:22, fontWeight:700, color:'#fff', flexShrink:0 }}>
            {name[0]?.toUpperCase()}
          </div>
          <div style={{ flex:1, display:'grid', gridTemplateColumns:'1fr 1fr', gap:10 }}>
            <div>
              <label style={{ display:'block', fontSize:10, color:C.text2, marginBottom:4, textTransform:'uppercase', letterSpacing:'0.06em' }}>Name</label>
              <input value={name} onChange={e=>setName(e.target.value)} style={{...inputSt, width:'100%'}} />
            </div>
            <div>
              <label style={{ display:'block', fontSize:10, color:C.text2, marginBottom:4, textTransform:'uppercase', letterSpacing:'0.06em' }}>Email</label>
              <input value={email} onChange={e=>setEmail(e.target.value)} style={{...inputSt, width:'100%'}} />
            </div>
          </div>
        </div>
        <Row label="Password" desc="Last changed 3 months ago"><Btn outline color={C.text2} size="sm">Update</Btn></Row>
      </Card>

      {/* Appearance */}
      <Card style={{ padding:'0 24px' }}>
        <SectionHead label="Appearance" />
        <Row label="Theme" desc="Color scheme for the interface">
          <div style={{ display:'flex', gap:4, background:C.surface2, borderRadius:10, padding:3 }}>
            {['dark','light','system'].map(t => (
              <button key={t} onClick={()=>setTheme(t)} style={{ padding:'5px 12px', borderRadius:7, fontSize:11, fontWeight:600, cursor:'pointer', textTransform:'capitalize',
                background: theme===t ? C.primary : 'transparent', color: theme===t ? '#fff' : C.text2, border:'none', transition:'all 0.15s' }}>{t}</button>
            ))}
          </div>
        </Row>
        <Row label="Accent Color" desc="Primary highlight color across the app">
          <div style={{ display:'flex', gap:6 }}>
            {ACCENTS.map(col => (
              <button key={col} onClick={()=>setAccent(col)} style={{ width:22, height:22, borderRadius:'50%', background:col, border:`2px solid ${accent===col ? '#fff' : 'transparent'}`, cursor:'pointer', transition:'all 0.15s', outline: accent===col ? `2px solid ${col}` : 'none', outlineOffset:1 }} />
            ))}
          </div>
        </Row>
        <Row label="Week starts on" desc="First day shown in Calendar">
          <select style={{...inputSt}}><option>Monday</option><option>Sunday</option></select>
        </Row>
      </Card>

      {/* Notifications */}
      <Card style={{ padding:'0 24px' }}>
        <SectionHead label="Notifications" />
        {[
          { key:'workout', emoji:'💪', label:'Workout reminder', desc:'Daily prompt to log a session', col:C.sections.workouts },
          { key:'study',   emoji:'📚', label:'Study session',    desc:'Start your focus time',         col:C.sections.study },
          { key:'sleep',   emoji:'😴', label:'Bedtime reminder', desc:'Wind-down notification',        col:C.sections.sleep },
          { key:'budget',  emoji:'💰', label:'Budget alerts',    desc:'Warn when spending threshold hit', col:C.sections.finance },
        ].map(n => (
          <Row key={n.key} label={<span>{n.emoji} {n.label}</span>} desc={n.desc}>
            <div style={{ display:'flex', alignItems:'center', gap:10 }}>
              {notifs[n.key].on && (
                <input type={n.key==='budget'?'text':'time'} value={n.key==='budget'?`>${notifs[n.key].threshold}%`:notifs[n.key].time}
                  onChange={e => setNotifs(ns=>({...ns,[n.key]:{...ns[n.key],[n.key==='budget'?'threshold':'time']:e.target.value}}))}
                  style={{...inputSt, width:n.key==='budget'?64:88, padding:'5px 8px', fontSize:12}} />
              )}
              <Toggle on={notifs[n.key].on} onChange={v=>setNotifs(ns=>({...ns,[n.key]:{...ns[n.key],on:v}}))} color={n.col} />
            </div>
          </Row>
        ))}
      </Card>

      {/* Goal defaults */}
      <Card style={{ padding:'0 24px' }}>
        <SectionHead label="Default Goal Targets" />
        {[
          { key:'workouts', label:'💪 Workouts per week', unit:'sessions', max:14 },
          { key:'study',    label:'📚 Study per week',    unit:'hours',    max:60 },
          { key:'sleep',    label:'😴 Sleep target',      unit:'hours',    max:12 },
        ].map(g => (
          <Row key={g.key} label={g.label} desc={`Target: ${goalDefs[g.key]} ${g.unit}`}>
            <Stepper value={goalDefs[g.key]} onChange={v=>setGoalDefs(gd=>({...gd,[g.key]:v}))} max={g.max} />
          </Row>
        ))}
        <Row label="💰 Monthly budget" desc="Spending limit for Finance goals">
          <div style={{ display:'flex', alignItems:'center', gap:6 }}>
            <span style={{ fontSize:13, color:C.text2 }}>$</span>
            <input type="number" value={goalDefs.budget} onChange={e=>setGoalDefs(gd=>({...gd,budget:+e.target.value}))}
              style={{...inputSt, width:90, textAlign:'right'}} />
          </div>
        </Row>
      </Card>

      {/* Data & Privacy */}
      <Card style={{ padding:'0 24px' }}>
        <SectionHead label="Data & Privacy" />
        <Row label="Export Data" desc="Download everything as JSON">
          <Btn outline color={C.primary} size="sm">Export JSON</Btn>
        </Row>
        <Row label="Backup" desc="Last synced 2 hours ago">
          <Btn outline color={C.text2} size="sm">Sync Now</Btn>
        </Row>
        <Row label="Delete Account" desc="Permanently remove all data">
          <Btn outline color="#f43f5e" size="sm">Delete Account</Btn>
        </Row>
      </Card>
    </div>
  );
}

// ── Sidebar ───────────────────────────────────────────────────────────────────

const NAV = [
  {id:'dashboard',label:'Dashboard',emoji:'📊'},
  {id:'workouts', label:'Workouts', emoji:'💪'},
  {id:'study',    label:'Study',    emoji:'📚'},
  {id:'sleep',    label:'Sleep',    emoji:'😴'},
  {id:'finance',  label:'Finance',  emoji:'💰'},
  {id:'calendar', label:'Calendar', emoji:'📅'},
  {id:'goals',    label:'Goals',    emoji:'🎯'},
  {id:'events',   label:'Events',   emoji:'📌'},
];

function NavItem({item,active,onClick}){
  const [hov,setHov]=React.useState(false);
  const acc=C.sections[item.id]||C.primary;
  const isActive=active===item.id;
  return(
    <button onClick={()=>onClick(item.id)} onMouseEnter={()=>setHov(true)} onMouseLeave={()=>setHov(false)}
      style={{ width:'100%',display:'flex',alignItems:'center',gap:10,padding:'9px 12px',borderRadius:10,cursor:'pointer',
        background:isActive?`${acc}15`:hov?C.border:'transparent',
        border:`1px solid ${isActive?`${acc}30`:'transparent'}`,
        borderLeft:`3px solid ${isActive?acc:'transparent'}`,
        transition:'all 0.15s',textAlign:'left' }}>
      <span style={{ fontSize:16 }}>{item.emoji}</span>
      <span style={{ fontSize:13,fontWeight:600,color:isActive?acc:hov?C.text:C.text2,transition:'color 0.15s' }}>{item.label}</span>
    </button>
  );
}

function Sidebar({page,navigate,onBulkLog}){
  return(
    <aside style={{ width:220,flexShrink:0,display:'flex',flexDirection:'column',gap:4,padding:'16px 12px',
      overflowY:'auto',borderRight:`1px solid ${C.border}`,background:C.surface }}>
      <div style={{ padding:'4px 8px 16px' }}>
        <span style={{ fontSize:20,fontWeight:800,letterSpacing:'-0.03em',fontFamily:'Space Grotesk,sans-serif',
          background:`linear-gradient(135deg, ${C.primary}, #a78bfa)`,WebkitBackgroundClip:'text',WebkitTextFillColor:'transparent' }}>Jeeb</span>
      </div>
      <div style={{ display:'flex',flexDirection:'column',gap:2 }}>
        {NAV.map(item=><NavItem key={item.id} item={item} active={page} onClick={navigate} />)}
      </div>
      <div style={{ height:1,background:C.border,margin:'8px 0' }} />
      <NavItem item={{id:'settings',label:'Settings',emoji:'⚙️'}} active={page} onClick={navigate} />
      {/* Bulk Log button */}
      <div style={{ margin:'8px 0 4px', padding:'0 4px' }}>
        <button onClick={onBulkLog} style={{
          width:'100%', display:'flex', alignItems:'center', gap:9, padding:'9px 12px',
          borderRadius:10, cursor:'pointer', transition:'all 0.15s', textAlign:'left',
          background:`linear-gradient(135deg, ${C.primary}18, #a78bfa18)`,
          border:`1px solid ${C.primary}30`,
        }}
          onMouseEnter={e=>{e.currentTarget.style.background=`linear-gradient(135deg, ${C.primary}28, #a78bfa28)`;}}
          onMouseLeave={e=>{e.currentTarget.style.background=`linear-gradient(135deg, ${C.primary}18, #a78bfa18)`;}}>
          <span style={{ fontSize:15 }}>📋</span>
          <span style={{ fontSize:13, fontWeight:600, background:`linear-gradient(135deg, ${C.primary}, #a78bfa)`, WebkitBackgroundClip:'text', WebkitTextFillColor:'transparent' }}>Bulk Log</span>
        </button>
      </div>
      <div style={{ marginTop:'auto',padding:'12px 8px 4px' }}>
        <div style={{ display:'flex',alignItems:'center',gap:10 }}>
          <div style={{ width:32,height:32,borderRadius:10,background:`linear-gradient(135deg, ${C.primary}, #a78bfa)`,
            display:'flex',alignItems:'center',justifyContent:'center',fontSize:13,fontWeight:700,color:'#fff' }}>{MOCK.user.name[0]}</div>
          <div>
            <div style={{ fontSize:12,fontWeight:600,color:C.text }}>{MOCK.user.name}</div>
            <div style={{ fontSize:10,color:C.text2 }}>alex@example.com</div>
          </div>
        </div>
      </div>
    </aside>
  );
}

// ── Header ────────────────────────────────────────────────────────────────────

function Header({page,onNotifClick,notifOpen,unreadCount}){
  return(
    <header style={{ height:56,display:'flex',alignItems:'center',gap:12,padding:'0 20px',
      borderBottom:`1px solid ${C.border}`,background:C.surface,flexShrink:0,position:'sticky',top:0,zIndex:100 }}>
      <div style={{ display:'flex',alignItems:'center',gap:8,flex:1,maxWidth:280,background:C.surface2,
        border:`1px solid ${C.border}`,borderRadius:10,padding:'7px 12px',cursor:'text' }}>
        <span style={{ fontSize:12,color:C.text3 }}>🔍</span>
        <span style={{ fontSize:13,color:C.text2 }}>Search…</span>
        <span style={{ marginLeft:'auto',fontSize:10,color:C.text3,background:C.surface3,padding:'2px 6px',borderRadius:4,border:`1px solid ${C.border2}` }}>⌘K</span>
      </div>
      <div style={{ marginLeft:'auto',display:'flex',alignItems:'center',gap:8 }}>
        <button onClick={onNotifClick} style={{ width:34,height:34,borderRadius:9,background:notifOpen?`${C.primary}20`:C.surface2,
          border:`1px solid ${notifOpen?C.primary:C.border}`,display:'flex',alignItems:'center',justifyContent:'center',
          cursor:'pointer',fontSize:14,position:'relative',transition:'all 0.15s' }}>
          🔔
          {unreadCount>0&&<div style={{ position:'absolute',top:6,right:6,width:6,height:6,borderRadius:'50%',background:'#f43f5e' }} />}
        </button>
        <div style={{ width:34,height:34,borderRadius:10,background:`linear-gradient(135deg, ${C.primary}, #a78bfa)`,
          display:'flex',alignItems:'center',justifyContent:'center',fontSize:13,fontWeight:700,color:'#fff',cursor:'pointer' }}>{MOCK.user.name[0]}</div>
      </div>
    </header>
  );
}

// ── FAB ───────────────────────────────────────────────────────────────────────

function FAB({ onClick }) {
  const [hov,setHov]=React.useState(false);
  return (
    <button onClick={onClick} onMouseEnter={()=>setHov(true)} onMouseLeave={()=>setHov(false)}
      style={{ position:'fixed',bottom:32,right:32,width:52,height:52,borderRadius:'50%',zIndex:90,
        background:`linear-gradient(135deg, ${C.primary}, #a78bfa)`,border:'none',cursor:'pointer',
        display:'flex',alignItems:'center',justifyContent:'center',fontSize:22,color:'#fff',
        boxShadow:hov?`0 8px 32px ${C.primary}60, 0 0 0 4px ${C.primary}20`:`0 4px 20px rgba(0,0,0,0.4)`,
        transform:hov?'scale(1.08)':'scale(1)', transition:'all 0.2s cubic-bezier(0.34,1.56,0.64,1)' }}>
      ＋
    </button>
  );
}

// ── App ───────────────────────────────────────────────────────────────────────

const PAGES = {dashboard:Dashboard,workouts:Workouts,sleep:Sleep,finance:Finance,study:Study,calendar:Calendar,settings:Settings,goals:Goals,events:Events};

function App(){
  const [page,setPage]=React.useState('dashboard');
  const [notifOpen,setNotifOpen]=React.useState(false);
  const [quickLog,setQuickLog]=React.useState({open:false,type:null});
  const [bulkOpen,setBulkOpen]=React.useState(false);
  const PageComponent=PAGES[page]||Dashboard;
  const unreadCount=NOTIFICATIONS.filter(n=>n.unread).length;

  function openQuickLog(type=null){ setQuickLog({open:true,type}); }
  function closeQuickLog(){ setQuickLog({open:false,type:null}); }

  // Expose bulk log opener globally for vanilla injector
  React.useEffect(() => {
    window.__openBulkLog = () => setBulkOpen(true);
    document.addEventListener('jeeb:openBulkLog', () => setBulkOpen(true));
    return () => { delete window.__openBulkLog; };
  }, []);

  return(
    <div style={{ display:'flex',height:'100vh',flexDirection:'column',background:C.bg,overflow:'hidden' }}>
      <Header page={page} onNotifClick={()=>setNotifOpen(o=>!o)} notifOpen={notifOpen} unreadCount={unreadCount} />
      <div style={{ display:'flex',flex:1,overflow:'hidden' }}>
        <Sidebar page={page} navigate={p=>{setPage(p);setNotifOpen(false);}} onBulkLog={()=>setBulkOpen(true)} />
        <main style={{ flex:1,overflowY:'auto',padding:'28px 32px',position:'relative' }}>
          <PageComponent navigate={setPage} openQuickLog={openQuickLog} />
        </main>
      </div>
      <NotificationsPanel open={notifOpen} onClose={()=>setNotifOpen(false)} />
      <QuickLogModal open={quickLog.open} initialType={quickLog.type} onClose={closeQuickLog} />
      <BulkLogModal open={bulkOpen} onClose={()=>setBulkOpen(false)} />
      <FAB onClick={()=>openQuickLog(null)} />
    </div>
  );
}

ReactDOM.createRoot(document.getElementById('root')).render(<App />);
