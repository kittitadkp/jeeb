// Bulk Insert Log Modal

const BULK_TYPES = [
  { id:'workout', emoji:'💪', label:'Workout',     accent:C.sections.workouts },
  { id:'study',   emoji:'📚', label:'Study',       accent:C.sections.study },
  { id:'sleep',   emoji:'😴', label:'Sleep',       accent:C.sections.sleep },
  { id:'finance', emoji:'💰', label:'Finance',     accent:C.sections.finance },
];

function pastDates(n) {
  return Array.from({length:n}, (_,i) => {
    const d = new Date(); d.setDate(d.getDate() - i);
    return d.toISOString().split('T')[0];
  });
}

function fmtDate(d) {
  const today = new Date().toISOString().split('T')[0];
  const yest  = new Date(Date.now()-86400000).toISOString().split('T')[0];
  if (d === today) return 'Today';
  if (d === yest)  return 'Yesterday';
  return new Date(d+'T12:00').toLocaleDateString('en-US',{weekday:'short',month:'short',day:'numeric'});
}

// Default row factories per type
function makeRow(type, date) {
  const id = Date.now() + Math.random();
  if (type==='workout') return { id, date, on:true, wtype:'strength', duration:45, notes:'' };
  if (type==='study')   return { id, date, on:true, subject:'', duration:60, notes:'' };
  if (type==='sleep')   return { id, date, on:true, bedtime:'23:00', wake:'07:00', quality:4 };
  if (type==='finance') return { id, date, on:true, ftype:'expense', amount:'', category:'', notes:'' };
}

// ── Row editors ───────────────────────────────────────────────────────────────

const cellSt = (accent) => ({
  background:'transparent', border:`1px solid ${C.border}`, borderRadius:7,
  padding:'6px 9px', fontSize:12, color:C.text, outline:'none', fontFamily:'inherit',
  width:'100%', transition:'border-color 0.15s'
});

const selSt = (accent) => ({ ...cellSt(accent), appearance:'none', cursor:'pointer' });

function WorkoutRow({ row, onChange }) {
  const acc = C.sections.workouts;
  const cs = cellSt(acc);
  return (
    <>
      <td style={{ padding:'4px 6px' }}>
        <select value={row.wtype} onChange={e=>onChange({wtype:e.target.value})} style={{...selSt(acc),width:100}}>
          <option value="strength">Strength</option>
          <option value="cardio">Cardio</option>
          <option value="flexibility">Flexibility</option>
        </select>
      </td>
      <td style={{ padding:'4px 6px' }}>
        <input type="number" min={1} value={row.duration} onChange={e=>onChange({duration:+e.target.value})} style={{...cs,width:72}} placeholder="min" />
      </td>
      <td style={{ padding:'4px 6px' }}>
        <input value={row.notes} onChange={e=>onChange({notes:e.target.value})} style={{...cs,width:160}} placeholder="Notes (optional)" />
      </td>
    </>
  );
}

function StudyRow({ row, onChange }) {
  const acc = C.sections.study;
  const cs = cellSt(acc);
  return (
    <>
      <td style={{ padding:'4px 6px' }}>
        <input value={row.subject} onChange={e=>onChange({subject:e.target.value})} style={{...cs,width:140}} placeholder="Subject" />
      </td>
      <td style={{ padding:'4px 6px' }}>
        <input type="number" min={5} value={row.duration} onChange={e=>onChange({duration:+e.target.value})} style={{...cs,width:72}} placeholder="min" />
      </td>
      <td style={{ padding:'4px 6px' }}>
        <input value={row.notes} onChange={e=>onChange({notes:e.target.value})} style={{...cs,width:160}} placeholder="Notes (optional)" />
      </td>
    </>
  );
}

function SleepRow({ row, onChange }) {
  const acc = C.sections.sleep;
  const cs = cellSt(acc);
  const dur = (() => {
    try {
      const [bh,bm] = row.bedtime.split(':').map(Number);
      let [wh,wm] = row.wake.split(':').map(Number);
      if (wh < bh) wh += 24;
      const mins = (wh*60+wm) - (bh*60+bm);
      const h = Math.floor(mins/60), m = mins%60;
      return m>0 ? `${h}h ${m}m` : `${h}h`;
    } catch { return '—'; }
  })();
  return (
    <>
      <td style={{ padding:'4px 6px' }}>
        <input type="time" value={row.bedtime} onChange={e=>onChange({bedtime:e.target.value})} style={{...cs,width:96}} />
      </td>
      <td style={{ padding:'4px 6px' }}>
        <input type="time" value={row.wake} onChange={e=>onChange({wake:e.target.value})} style={{...cs,width:96}} />
      </td>
      <td style={{ padding:'4px 6px' }}>
        <span style={{ fontSize:12, color:C.text2, paddingLeft:4 }}>{dur}</span>
      </td>
      <td style={{ padding:'4px 6px' }}>
        <div style={{ display:'flex', gap:3 }}>
          {[1,2,3,4,5].map(n => (
            <button key={n} onClick={()=>onChange({quality:n})} style={{
              width:26, height:26, borderRadius:6, fontSize:11, fontWeight:700, cursor:'pointer',
              background: row.quality===n ? acc : C.surface2,
              color: row.quality===n ? '#fff' : C.text2,
              border:`1px solid ${row.quality===n ? acc : C.border2}`, transition:'all 0.1s'
            }}>{n}</button>
          ))}
        </div>
      </td>
    </>
  );
}

function FinanceRow({ row, onChange }) {
  const acc = C.sections.finance;
  const cs = cellSt(acc);
  const isExp = row.ftype === 'expense';
  return (
    <>
      <td style={{ padding:'4px 6px' }}>
        <select value={row.ftype} onChange={e=>onChange({ftype:e.target.value})} style={{...selSt(acc),width:88}}>
          <option value="expense">Expense</option>
          <option value="income">Income</option>
        </select>
      </td>
      <td style={{ padding:'4px 6px' }}>
        <div style={{ position:'relative' }}>
          <span style={{ position:'absolute', left:8, top:'50%', transform:'translateY(-50%)', fontSize:12, color:C.text2 }}>$</span>
          <input type="number" min="0.01" step="0.01" value={row.amount} onChange={e=>onChange({amount:e.target.value})}
            style={{...cs, paddingLeft:20, width:90}} placeholder="0.00" />
        </div>
      </td>
      <td style={{ padding:'4px 6px' }}>
        <input value={row.category} onChange={e=>onChange({category:e.target.value})} style={{...cs,width:120}} placeholder="Category" />
      </td>
      <td style={{ padding:'4px 6px' }}>
        <input value={row.notes} onChange={e=>onChange({notes:e.target.value})} style={{...cs,width:130}} placeholder="Notes" />
      </td>
    </>
  );
}

// ── Column headers ────────────────────────────────────────────────────────────

function ColHeaders({ type }) {
  const thSt = { fontSize:10, fontWeight:700, color:C.text2, textTransform:'uppercase', letterSpacing:'0.06em', padding:'8px 6px', textAlign:'left', whiteSpace:'nowrap' };
  if (type==='workout') return <><th style={thSt}>Type</th><th style={thSt}>Duration</th><th style={thSt}>Notes</th></>;
  if (type==='study')   return <><th style={thSt}>Subject</th><th style={thSt}>Duration</th><th style={thSt}>Notes</th></>;
  if (type==='sleep')   return <><th style={thSt}>Bedtime</th><th style={thSt}>Wake Up</th><th style={thSt}>Duration</th><th style={thSt}>Quality</th></>;
  if (type==='finance') return <><th style={thSt}>Type</th><th style={thSt}>Amount</th><th style={thSt}>Category</th><th style={thSt}>Notes</th></>;
  return null;
}

// ── Main Modal ────────────────────────────────────────────────────────────────

function BulkLogModal({ open, onClose }) {
  const [step, setStep] = React.useState('pick'); // 'pick' | 'log'
  const [type, setType] = React.useState(null);
  const [rows, setRows] = React.useState([]);
  const [saved, setSaved] = React.useState(false);

  React.useEffect(() => {
    if (!open) { setTimeout(() => { setStep('pick'); setType(null); setRows([]); setSaved(false); }, 300); }
  }, [open]);

  function pickType(t) {
    setType(t);
    const dates = pastDates(7);
    setRows(dates.map(d => makeRow(t, d)));
    setStep('log');
  }

  function updateRow(id, patch) {
    setRows(rs => rs.map(r => r.id===id ? {...r,...patch} : r));
  }

  function toggleRow(id) {
    setRows(rs => rs.map(r => r.id===id ? {...r, on:!r.on} : r));
  }

  function addRow() {
    const d = new Date(); d.setDate(d.getDate() - rows.length);
    setRows(rs => [...rs, makeRow(type, d.toISOString().split('T')[0])]);
  }

  function toggleAll() {
    const allOn = rows.every(r => r.on);
    setRows(rs => rs.map(r => ({...r, on:!allOn})));
  }

  function handleSave() {
    setSaved(true);
    setTimeout(() => { onClose(); }, 1400);
  }

  const info = BULK_TYPES.find(t => t.id === type);
  const acc = info?.accent || C.primary;
  const selectedCount = rows.filter(r => r.on).length;

  if (!open) return null;

  const thSt = { fontSize:10, fontWeight:700, color:C.text2, textTransform:'uppercase', letterSpacing:'0.06em', padding:'8px 6px', textAlign:'left', whiteSpace:'nowrap' };

  return (
    <div style={{ position:'fixed', inset:0, zIndex:200, display:'flex', alignItems:'center', justifyContent:'center' }}>
      <div onClick={onClose} style={{ position:'absolute', inset:0, background:'rgba(0,0,0,0.8)', backdropFilter:'blur(8px)', WebkitBackdropFilter:'blur(8px)' }} />
      <div style={{
        position:'relative', width:'100%', maxWidth: step==='pick' ? 480 : 860,
        margin:'0 16px', background:C.surface, borderRadius:20, border:`1px solid ${C.border2}`,
        boxShadow:'0 32px 100px rgba(0,0,0,0.7)', maxHeight:'90vh',
        display:'flex', flexDirection:'column',
        transition:'max-width 0.3s cubic-bezier(0.4,0,0.2,1)',
        animation:'qlIn 0.22s cubic-bezier(0.34,1.4,0.64,1)'
      }}>
        {/* Header */}
        <div style={{ padding:'20px 24px 16px', borderBottom:`1px solid ${C.border}`, display:'flex', alignItems:'center', justifyContent:'space-between', flexShrink:0 }}>
          <div style={{ display:'flex', alignItems:'center', gap:10 }}>
            {step==='log' && (
              <button onClick={()=>setStep('pick')} style={{ width:28, height:28, borderRadius:8, background:C.surface2, border:`1px solid ${C.border}`, color:C.text2, cursor:'pointer', fontSize:14, display:'flex', alignItems:'center', justifyContent:'center' }}>‹</button>
            )}
            <div>
              <div style={{ fontSize:16, fontWeight:700, color:C.text, fontFamily:'Space Grotesk,sans-serif', lineHeight:1.2 }}>
                {step==='pick' ? 'Bulk Log' : `${info?.emoji} Bulk ${info?.label} Log`}
              </div>
              {step==='log' && <div style={{ fontSize:11, color:C.text2, marginTop:2 }}>{selectedCount} of {rows.length} entries selected</div>}
            </div>
          </div>
          <button onClick={onClose} style={{ width:28, height:28, borderRadius:8, background:C.surface2, border:`1px solid ${C.border}`, color:C.text2, cursor:'pointer', fontSize:16, display:'flex', alignItems:'center', justifyContent:'center' }}>×</button>
        </div>

        {/* Body */}
        <div style={{ flex:1, overflowY:'auto', padding:'20px 24px' }}>
          {step==='pick' && (
            <div>
              <p style={{ fontSize:13, color:C.text2, marginBottom:20 }}>Pick a category to bulk-log multiple entries at once — great for catching up on missed days.</p>
              <div style={{ display:'grid', gridTemplateColumns:'1fr 1fr', gap:10 }}>
                {BULK_TYPES.map(t => (
                  <button key={t.id} onClick={()=>pickType(t.id)} style={{
                    padding:'20px 18px', borderRadius:14, border:`1px solid ${C.border2}`,
                    background:C.surface2, cursor:'pointer', textAlign:'left', transition:'all 0.15s' }}
                    onMouseEnter={e=>{e.currentTarget.style.borderColor=t.accent;e.currentTarget.style.background=`${t.accent}12`;}}
                    onMouseLeave={e=>{e.currentTarget.style.borderColor=C.border2;e.currentTarget.style.background=C.surface2;}}>
                    <div style={{ fontSize:28, marginBottom:8 }}>{t.emoji}</div>
                    <div style={{ fontSize:14, fontWeight:700, color:C.text }}>{t.label}</div>
                    <div style={{ fontSize:11, color:C.text2, marginTop:3 }}>Log multiple {t.label.toLowerCase()} entries</div>
                  </button>
                ))}
              </div>
            </div>
          )}

          {step==='log' && !saved && (
            <div>
              {/* Quick actions */}
              <div style={{ display:'flex', gap:8, marginBottom:14, flexWrap:'wrap' }}>
                <button onClick={toggleAll} style={{ padding:'5px 12px', borderRadius:8, fontSize:12, fontWeight:600, cursor:'pointer', background:C.surface2, border:`1px solid ${C.border2}`, color:C.text2 }}>
                  {rows.every(r=>r.on) ? '☐ Deselect all' : '☑ Select all'}
                </button>
                <button onClick={()=>setRows(pastDates(7).map(d=>makeRow(type,d)))} style={{ padding:'5px 12px', borderRadius:8, fontSize:12, fontWeight:600, cursor:'pointer', background:C.surface2, border:`1px solid ${C.border2}`, color:C.text2 }}>
                  ↺ Last 7 days
                </button>
                <button onClick={()=>setRows(pastDates(14).map(d=>makeRow(type,d)))} style={{ padding:'5px 12px', borderRadius:8, fontSize:12, fontWeight:600, cursor:'pointer', background:C.surface2, border:`1px solid ${C.border2}`, color:C.text2 }}>
                  ↺ Last 14 days
                </button>
                <button onClick={addRow} style={{ padding:'5px 12px', borderRadius:8, fontSize:12, fontWeight:600, cursor:'pointer', background:`${acc}18`, border:`1px solid ${acc}40`, color:acc }}>
                  + Add row
                </button>
              </div>

              {/* Table */}
              <div style={{ overflowX:'auto', borderRadius:12, border:`1px solid ${C.border}` }}>
                <table style={{ width:'100%', borderCollapse:'collapse', fontSize:13 }}>
                  <thead>
                    <tr style={{ borderBottom:`1px solid ${C.border}`, background:C.surface2 }}>
                      <th style={{ ...thSt, width:28, paddingLeft:12 }}>
                        <input type="checkbox" checked={rows.every(r=>r.on)} onChange={toggleAll}
                          style={{ accentColor:acc, width:14, height:14, cursor:'pointer' }} />
                      </th>
                      <th style={{ ...thSt, minWidth:130 }}>Date</th>
                      <ColHeaders type={type} />
                      <th style={{ ...thSt, width:32 }}></th>
                    </tr>
                  </thead>
                  <tbody>
                    {rows.map((row, i) => (
                      <tr key={row.id} style={{
                        borderBottom: i < rows.length-1 ? `1px solid ${C.border}` : 'none',
                        background: row.on ? 'transparent' : `${C.surface2}80`,
                        opacity: row.on ? 1 : 0.45, transition:'all 0.15s'
                      }}>
                        <td style={{ padding:'4px 6px 4px 12px' }}>
                          <input type="checkbox" checked={row.on} onChange={()=>toggleRow(row.id)}
                            style={{ accentColor:acc, width:14, height:14, cursor:'pointer' }} />
                        </td>
                        <td style={{ padding:'4px 6px', whiteSpace:'nowrap' }}>
                          <div style={{ display:'flex', alignItems:'center', gap:6 }}>
                            <div style={{ fontSize:11, fontWeight:600, color: row.on ? acc : C.text2 }}>
                              {fmtDate(row.date)}
                            </div>
                            <input type="date" value={row.date} onChange={e=>updateRow(row.id,{date:e.target.value})}
                              style={{ background:'transparent', border:'none', color:C.text3, fontSize:10, outline:'none', cursor:'pointer', fontFamily:'inherit' }} />
                          </div>
                        </td>
                        {type==='workout' && <WorkoutRow row={row} onChange={p=>updateRow(row.id,p)} />}
                        {type==='study'   && <StudyRow   row={row} onChange={p=>updateRow(row.id,p)} />}
                        {type==='sleep'   && <SleepRow   row={row} onChange={p=>updateRow(row.id,p)} />}
                        {type==='finance' && <FinanceRow row={row} onChange={p=>updateRow(row.id,p)} />}
                        <td style={{ padding:'4px 8px' }}>
                          <button onClick={()=>setRows(rs=>rs.filter(r=>r.id!==row.id))}
                            style={{ color:C.text3, fontSize:15, cursor:'pointer', background:'none', border:'none', lineHeight:1, padding:'2px 4px' }}
                            onMouseEnter={e=>e.target.style.color='#f43f5e'}
                            onMouseLeave={e=>e.target.style.color=C.text3}>×</button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          )}

          {saved && (
            <div style={{ textAlign:'center', padding:'40px 0' }}>
              <div style={{ fontSize:48, marginBottom:16, animation:'qlIn 0.3s ease' }}>✅</div>
              <div style={{ fontSize:18, fontWeight:700, color:C.text, fontFamily:'Space Grotesk,sans-serif' }}>
                {selectedCount} {info?.label} entr{selectedCount===1?'y':'ies'} saved!
              </div>
              <div style={{ fontSize:13, color:C.text2, marginTop:6 }}>Your logs have been added.</div>
            </div>
          )}
        </div>

        {/* Footer */}
        {step==='log' && !saved && (
          <div style={{ padding:'16px 24px', borderTop:`1px solid ${C.border}`, display:'flex', justifyContent:'space-between', alignItems:'center', flexShrink:0 }}>
            <span style={{ fontSize:12, color:C.text2 }}>
              {selectedCount === 0 ? 'Select at least one row to save' : `Saving ${selectedCount} entr${selectedCount===1?'y':'ies'}`}
            </span>
            <div style={{ display:'flex', gap:8 }}>
              <Btn outline color={C.text2} size="sm" onClick={onClose}>Cancel</Btn>
              <Btn color={acc} disabled={selectedCount===0} onClick={handleSave}>
                💾 Save {selectedCount} {selectedCount===1?'Entry':'Entries'}
              </Btn>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

Object.assign(window, { BulkLogModal });
