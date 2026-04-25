// Sleep + Finance + Study pages

// ── Sleep ─────────────────────────────────────────────────────────────────────

function Stars({ count }) {
  return (
    <span>{[1,2,3,4,5].map(i => (
      <span key={i} style={{ fontSize:12, color: i<=count ? '#f59e0b' : C.text3 }}>★</span>
    ))}</span>
  );
}

function SleepBar({ hours, maxHours=10, color }) {
  const pct = Math.min(100, (hours/maxHours)*100);
  return (
    <div style={{ position:'relative', height:6, borderRadius:6, background:C.surface3, overflow:'hidden' }}>
      <div style={{ position:'absolute', left:0, top:0, height:'100%', width:`${pct}%`, borderRadius:6,
        background:`linear-gradient(90deg, ${color}80, ${color})`, transition:'width 0.4s ease' }} />
      <div style={{ position:'absolute', left:'80%', top:-2, bottom:-2, width:1, background:C.text3 }} />
    </div>
  );
}

function Sleep({ openQuickLog }) {
  const [records, setRecords] = React.useState(MOCK.sleep.recent);
  const [period, setPeriod] = React.useState('week');
  const acc = C.sections.sleep;
  const pd = PERIOD_DATA.sleep[period];
  const DAYS7 = ['Mon','Tue','Wed','Thu','Fri','Sat','Sun'];

  function dur(start, end) {
    const m = (new Date(end)-new Date(start))/60000;
    const h = Math.floor(m/60), mn = Math.round(m%60);
    return mn>0 ? `${h}h ${mn}m` : `${h}h`;
  }

  const chartIsWeek = period==='week';

  return (
    <div style={{ maxWidth:760, display:'flex', flexDirection:'column', gap:20 }}>
      <div style={{ display:'flex', alignItems:'center', justifyContent:'space-between' }}>
        <SectionLabel color={acc}>Sleep</SectionLabel>
        <Btn color={acc} onClick={() => openQuickLog('sleep')}>+ Log Sleep</Btn>
      </div>

      <div style={{ display:'flex', gap:12, flexWrap:'wrap' }}>
        <StatCard emoji="😴" label={pd.unit} value={pd.value} change={pd.change} trend="neutral" accent={acc} sparkData={pd.spark} />
        <StatCard emoji="⭐" label="avg quality" value={`${MOCK.sleep.avgQuality}/5`} change="all time" accent={acc} sparkData={MOCK.sleep.weeklyQuality} />
        <StatCard emoji="📅" label="logged this month" value={MOCK.sleep.thisMonth} accent={acc} />
      </div>

      <Card style={{ padding:'20px 24px' }}>
        <div style={{ display:'flex', justifyContent:'space-between', alignItems:'center', marginBottom:16 }}>
          <span style={{ fontSize:12, fontWeight:600, color:C.text2, textTransform:'uppercase', letterSpacing:'0.06em' }}>
            Sleep Duration {period==='day' ? '— Last Night' : `— ${period.charAt(0).toUpperCase()+period.slice(1)}`}
          </span>
          <PeriodSelector value={period} onChange={setPeriod} accent={acc} />
        </div>
        {chartIsWeek ? (
          <div style={{ display:'flex', flexDirection:'column', gap:10 }}>
            {DAYS7.map((day,i) => (
              <div key={day} style={{ display:'flex', alignItems:'center', gap:12 }}>
                <span style={{ fontSize:11, color:C.text2, width:28 }}>{day}</span>
                <div style={{ flex:1 }}><SleepBar hours={MOCK.sleep.weeklyHours[i]} color={acc} /></div>
                <div style={{ display:'flex', alignItems:'center', gap:6, width:60 }}>
                  <span style={{ fontSize:12, color:C.text, fontWeight:500 }}>{MOCK.sleep.weeklyHours[i]}h</span>
                  <Stars count={MOCK.sleep.weeklyQuality[i]} />
                </div>
              </div>
            ))}
            <div style={{ marginTop:4, fontSize:11, color:C.text3 }}>⎸ = 8h recommended</div>
          </div>
        ) : (
          <BarChart data={pd.chart} color={acc} height={80} />
        )}
      </Card>

      <div style={{ display:'flex', flexDirection:'column', gap:12 }}>
        {records.length===0
          ? <Card><EmptyState emoji="😴" title="No sleep records" desc="Log your sleep to track patterns" /></Card>
          : records.map(r => (
            <Card key={r.id} accent={acc}>
              <div style={{ padding:'14px 18px', display:'flex', alignItems:'center', gap:12 }}>
                <div style={{ width:36, height:36, borderRadius:10, background:`${acc}18`,
                  display:'flex', alignItems:'center', justifyContent:'center', fontSize:17, flexShrink:0 }}>😴</div>
                <div style={{ flex:1 }}>
                  <div style={{ display:'flex', alignItems:'center', gap:8, marginBottom:2 }}>
                    <span style={{ fontSize:14, fontWeight:600, color:C.text }}>
                      {new Date(r.start).toLocaleDateString('en-US',{month:'short',day:'numeric'})}
                    </span>
                    <Stars count={r.quality} />
                  </div>
                  <div style={{ fontSize:12, color:C.text2 }}>
                    {new Date(r.start).toLocaleTimeString([],{hour:'2-digit',minute:'2-digit'})} →{' '}
                    {new Date(r.end).toLocaleTimeString([],{hour:'2-digit',minute:'2-digit'})} · {dur(r.start,r.end)}
                  </div>
                  {r.notes && <div style={{ fontSize:11, color:C.text3, marginTop:3, fontStyle:'italic' }}>{r.notes}</div>}
                </div>
                <Btn outline color="#f43f5e" size="sm" onClick={() => setRecords(rs => rs.filter(x=>x.id!==r.id))}>🗑</Btn>
              </div>
            </Card>
          ))
        }
      </div>
    </div>
  );
}

// ── Finance ───────────────────────────────────────────────────────────────────

const CAT_COLORS = ['#10b981','#7c6ef5','#f59e0b','#f43f5e','#14b8a6','#a855f7'];
const CAT_EMOJI = { 'Food & Dining':'🛒', Salary:'💼', Transport:'🚗', Utilities:'💡', Entertainment:'🎬', Freelance:'💻', Healthcare:'🏥', Shopping:'🛍', default:'💳' };

function TransactionItem({ tx, onDelete }) {
  const isExp = tx.type==='expense';
  const col = isExp ? '#f43f5e' : '#10b981';
  return (
    <div style={{ display:'flex', alignItems:'center', gap:12, padding:'12px 0', borderBottom:`1px solid ${C.border}` }}>
      <div style={{ width:38, height:38, borderRadius:10, flexShrink:0,
        background:`${col}15`, display:'flex', alignItems:'center', justifyContent:'center', fontSize:17 }}>
        {CAT_EMOJI[tx.category]||CAT_EMOJI.default}
      </div>
      <div style={{ flex:1, minWidth:0 }}>
        <div style={{ fontSize:13, fontWeight:500, color:C.text, overflow:'hidden', textOverflow:'ellipsis', whiteSpace:'nowrap' }}>{tx.notes||tx.category}</div>
        <div style={{ fontSize:11, color:C.text2, marginTop:2 }}>{tx.category} · {new Date(tx.date).toLocaleDateString('en-US',{month:'short',day:'numeric'})}</div>
      </div>
      <div style={{ fontSize:14, fontWeight:700, color:col, flexShrink:0 }}>
        {isExp ? '−' : '+'}${tx.amount.toLocaleString('en-US',{minimumFractionDigits:2,maximumFractionDigits:2})}
      </div>
      <button onClick={onDelete} style={{ color:C.text3, fontSize:18, padding:'2px 4px', cursor:'pointer', background:'none', border:'none', lineHeight:1 }}
        onMouseEnter={e=>e.target.style.color='#f43f5e'} onMouseLeave={e=>e.target.style.color=C.text3}>×</button>
    </div>
  );
}

function Finance({ openQuickLog }) {
  const [txns, setTxns] = React.useState(MOCK.finance.transactions);
  const [period, setPeriod] = React.useState('month');
  const acc = C.sections.finance;
  const pd = PERIOD_DATA.finance[period];
  const { income, expense, net, byCategory } = MOCK.finance;
  const catEntries = Object.entries(byCategory).sort(([,a],[,b])=>b-a);
  const donutSegments = catEntries.map(([label,value],i)=>({label,value,color:CAT_COLORS[i%CAT_COLORS.length]}));

  return (
    <div style={{ maxWidth:760, display:'flex', flexDirection:'column', gap:20 }}>
      <div style={{ display:'flex', alignItems:'center', justifyContent:'space-between' }}>
        <SectionLabel color={acc}>Finance</SectionLabel>
        <Btn color={acc} onClick={() => openQuickLog('finance')}>+ Add Transaction</Btn>
      </div>

      <div style={{ display:'flex', gap:12, flexWrap:'wrap' }}>
        <StatCard emoji="💰" label="income this month" value={`$${income.toLocaleString()}`} trend="up" accent={acc} sparkData={MOCK.finance.weeklyExpense} />
        <StatCard emoji="💸" label="expenses this month" value={`$${expense.toLocaleString()}`} trend="down" accent="#f43f5e" sparkData={MOCK.finance.weeklyExpense} />
        <StatCard emoji="📊" label="net this month" value={`+$${net.toLocaleString()}`} trend="up" accent={acc} />
      </div>

      <div style={{ display:'grid', gridTemplateColumns:'1fr auto', gap:16, alignItems:'start' }}>
        <Card style={{ padding:'20px 24px' }}>
          <div style={{ display:'flex', justifyContent:'space-between', alignItems:'center', marginBottom:16 }}>
            <span style={{ fontSize:12, fontWeight:600, color:C.text2, textTransform:'uppercase', letterSpacing:'0.06em' }}>Spending</span>
            <PeriodSelector value={period} onChange={setPeriod} accent={acc} />
          </div>
          <div style={{ fontSize:22, fontWeight:700, color:C.text, fontFamily:'Space Grotesk,sans-serif', marginBottom:4 }}>{pd.value}</div>
          <div style={{ fontSize:11, color:C.text2, marginBottom:16 }}>{pd.change}</div>
          <BarChart data={pd.chart} color={acc} height={80} />
        </Card>

        <Card style={{ padding:'20px 24px', display:'flex', flexDirection:'column', alignItems:'center', gap:14 }}>
          <span style={{ fontSize:12, fontWeight:600, color:C.text2, textTransform:'uppercase', letterSpacing:'0.06em' }}>By Category</span>
          <DonutChart segments={donutSegments} size={130} />
          <div style={{ display:'flex', flexDirection:'column', gap:5, width:'100%' }}>
            {donutSegments.slice(0,4).map((s,i) => (
              <div key={i} style={{ display:'flex', alignItems:'center', gap:8 }}>
                <div style={{ width:8, height:8, borderRadius:2, background:s.color, flexShrink:0 }} />
                <span style={{ fontSize:11, color:C.text2, flex:1 }}>{s.label}</span>
                <span style={{ fontSize:11, color:C.text, fontWeight:600 }}>${s.value}</span>
              </div>
            ))}
          </div>
        </Card>
      </div>

      <Card style={{ padding:'0 20px' }}>
        <div style={{ padding:'16px 0 12px', borderBottom:`1px solid ${C.border}`, display:'flex', justifyContent:'space-between', alignItems:'center' }}>
          <span style={{ fontSize:13, fontWeight:600, color:C.text }}>Transactions</span>
          <span style={{ fontSize:11, color:C.text2 }}>{txns.length} records</span>
        </div>
        {txns.length===0
          ? <EmptyState emoji="💰" title="No transactions" desc="Add income or expenses to get started" />
          : txns.map(tx => <TransactionItem key={tx.id} tx={tx} onDelete={() => setTxns(ts=>ts.filter(x=>x.id!==tx.id))} />)
        }
      </Card>
    </div>
  );
}

// ── Study ─────────────────────────────────────────────────────────────────────

function Study({ openQuickLog }) {
  const [sessions, setSessions] = React.useState(MOCK.study.recent);
  const [period, setPeriod] = React.useState('week');
  const acc = C.sections.study;
  const pd = PERIOD_DATA.study[period];

  return (
    <div style={{ maxWidth:760, display:'flex', flexDirection:'column', gap:20 }}>
      <div style={{ display:'flex', alignItems:'center', justifyContent:'space-between' }}>
        <SectionLabel color={acc}>Study</SectionLabel>
        <Btn color={acc} onClick={() => openQuickLog('study')}>+ Log Session</Btn>
      </div>

      <div style={{ display:'flex', gap:12, flexWrap:'wrap' }}>
        <StatCard emoji="📚" label={pd.unit} value={pd.value} change={pd.change} trend="up" accent={acc} sparkData={pd.spark} />
        <StatCard emoji="🗓" label="sessions logged" value={MOCK.study.sessions} accent={acc} />
        <StatCard emoji="🎯" label="subjects tracked" value={3} accent={acc} />
      </div>

      <Card style={{ padding:'20px 24px' }}>
        <div style={{ display:'flex', justifyContent:'space-between', alignItems:'center', marginBottom:16 }}>
          <span style={{ fontSize:12, fontWeight:600, color:C.text2, textTransform:'uppercase', letterSpacing:'0.06em' }}>Study Hours</span>
          <PeriodSelector value={period} onChange={setPeriod} accent={acc} />
        </div>
        <div style={{ fontSize:22, fontWeight:700, color:C.text, fontFamily:'Space Grotesk,sans-serif', marginBottom:4 }}>{pd.value}</div>
        <div style={{ fontSize:11, color:C.text2, marginBottom:16 }}>{pd.change}</div>
        <BarChart data={pd.chart} color={acc} height={80} />
      </Card>

      <div style={{ display:'flex', flexDirection:'column', gap:10 }}>
        {sessions.length===0
          ? <Card><EmptyState emoji="📚" title="No sessions yet" desc="Start logging your study time" /></Card>
          : sessions.map(s => (
            <Card key={s.id} accent={acc}>
              <div style={{ padding:'14px 18px', display:'flex', alignItems:'center', gap:12 }}>
                <div style={{ width:36, height:36, borderRadius:10, background:`${acc}18`,
                  display:'flex', alignItems:'center', justifyContent:'center', fontSize:17, flexShrink:0 }}>📚</div>
                <div style={{ flex:1 }}>
                  <div style={{ fontSize:14, fontWeight:600, color:C.text }}>{s.subject}</div>
                  <div style={{ fontSize:12, color:C.text2, marginTop:2 }}>
                    {Math.floor(s.duration/60)}h {s.duration%60>0 ? `${s.duration%60}m` : ''} ·{' '}
                    {new Date(s.date).toLocaleDateString('en-US',{month:'short',day:'numeric'})}
                  </div>
                </div>
                <Btn outline color="#f43f5e" size="sm" onClick={() => setSessions(ss=>ss.filter(x=>x.id!==s.id))}>🗑</Btn>
              </div>
            </Card>
          ))
        }
      </div>
    </div>
  );
}

Object.assign(window, { Sleep, Finance, Study });
