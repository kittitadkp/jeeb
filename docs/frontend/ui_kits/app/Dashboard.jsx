// Jeeb UI Kit — Dashboard Page v2

const { useState } = React;

function Card({ children, style }) {
  const t = window.useTheme();
  return (
    <div style={{ background: t.surface, borderRadius: 8, boxShadow: t.shadow, border: `1px solid ${t.border}`, ...style }}>
      {children}
    </div>
  );
}

function StatCard({ emoji, label, value, change, trend }) {
  const t = window.useTheme();
  const trendColor = trend === 'up' ? t.successText : trend === 'down' ? t.dangerText : t.fg3;
  const arrow = trend === 'up' ? '↑' : trend === 'down' ? '↓' : '';
  return (
    <Card style={{ padding: '18px 20px', flex: 1, minWidth: 140 }}>
      <div style={{ fontSize: 22, marginBottom: 8 }}>{emoji}</div>
      <div style={{ fontSize: 28, fontWeight: 700, color: t.fg1, lineHeight: 1 }}>{value}</div>
      <div style={{ fontSize: 12, color: t.fg2, marginTop: 4 }}>{label}</div>
      {change && <div style={{ fontSize: 12, fontWeight: 500, color: trendColor, marginTop: 6 }}>{arrow} {change}</div>}
    </Card>
  );
}

function ActivityItem({ emoji, text, time }) {
  const t = window.useTheme();
  return (
    <div style={{ display: 'flex', alignItems: 'center', gap: 10, padding: '8px 0', borderBottom: `1px solid ${t.border}` }}>
      <div style={{ width: 32, height: 32, borderRadius: '50%', background: t.primarySubtle, display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: 14, flexShrink: 0 }}>{emoji}</div>
      <div style={{ flex: 1, fontSize: 13, color: t.fg1, fontWeight: 500 }}>{text}</div>
      <div style={{ fontSize: 11, color: t.fg3, whiteSpace: 'nowrap' }}>{time}</div>
    </div>
  );
}

function UpcomingItem({ emoji, text, time }) {
  const t = window.useTheme();
  return (
    <div style={{ display: 'flex', alignItems: 'center', gap: 10, padding: '8px 0', borderBottom: `1px solid ${t.border}` }}>
      <div style={{ width: 32, height: 32, borderRadius: 6, background: t.surfaceActive, border: `1px solid ${t.border}`, display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: 14, flexShrink: 0 }}>{emoji}</div>
      <div style={{ flex: 1 }}>
        <div style={{ fontSize: 13, color: t.fg1, fontWeight: 500 }}>{text}</div>
        <div style={{ fontSize: 11, color: t.fg3 }}>{time}</div>
      </div>
    </div>
  );
}

function Dashboard({ onNav }) {
  const t = window.useTheme();
  const hour = new Date().getHours();
  const greeting = hour < 12 ? 'Good morning' : hour < 17 ? 'Good afternoon' : 'Good evening';
  const width = window.useWindowWidth();
  const isMobile = width < 768;

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 20, maxWidth: 960 }}>
      <div>
        <h1 style={{ fontSize: isMobile ? 20 : 24, fontWeight: 700, color: t.fg1 }}>{greeting}, John 👋</h1>
        <p style={{ fontSize: 14, color: t.fg2, marginTop: 2 }}>Here's your summary for today, Apr 23</p>
      </div>

      <div style={{ display: 'flex', gap: 12, flexWrap: 'wrap' }}>
        <StatCard emoji="💪" label="workouts this week" value="5" change="+3 from last week" trend="up" />
        <StatCard emoji="📚" label="study hours this week" value="12h" change="+2h from last week" trend="up" />
        <StatCard emoji="😴" label="avg sleep this week" value="7.2h" change="−0.3h from last week" trend="down" />
        <StatCard emoji="💰" label="spending this month" value="$1,250" change="of $2,000 budget" trend="neutral" />
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: isMobile ? '1fr' : '1fr 1fr', gap: 16 }}>
        <Card style={{ padding: '16px 20px' }}>
          <h2 style={{ fontSize: 15, fontWeight: 600, color: t.fg1, marginBottom: 12 }}>Recent Activity</h2>
          <ActivityItem emoji="💪" text="Upper Body workout — 45 min" time="2h ago" />
          <ActivityItem emoji="😴" text="Sleep logged — 7h 30m, quality 4/5" time="8h ago" />
          <ActivityItem emoji="💰" text="Groceries — −$85.50" time="Yesterday" />
          <ActivityItem emoji="📚" text="Studied Mathematics — 2h 30m" time="Yesterday" />
          <ActivityItem emoji="💪" text="Cardio — 30 min" time="2 days ago" />
        </Card>
        <Card style={{ padding: '16px 20px' }}>
          <h2 style={{ fontSize: 15, fontWeight: 600, color: t.fg1, marginBottom: 12 }}>Upcoming</h2>
          <UpcomingItem emoji="📚" text="Study Mathematics" time="Tomorrow · 2:00 PM" />
          <UpcomingItem emoji="💪" text="Gym session" time="Tomorrow · 9:00 AM" />
          <UpcomingItem emoji="💰" text="Electricity bill due" time="Apr 25" />
          <UpcomingItem emoji="📅" text="Doctor appointment" time="Apr 28 · 11:00 AM" />
        </Card>
      </div>

      <div style={{ display: 'flex', gap: 10, flexWrap: 'wrap' }}>
        {[['💪','+ Log Workout','workouts'],['📚','+ Start Study','study'],['😴','+ Log Sleep','sleep'],['💰','+ Add Transaction','finance']].map(([emoji, label, page]) => (
          <button key={page} onClick={() => onNav(page)} style={{
            display: 'flex', alignItems: 'center', gap: 6, padding: '8px 16px',
            border: `1px solid ${t.border}`, borderRadius: 8, background: t.surface,
            cursor: 'pointer', fontSize: 13, fontWeight: 500, color: t.fg1,
          }} onMouseEnter={e => e.currentTarget.style.background = t.surfaceHover}
             onMouseLeave={e => e.currentTarget.style.background = t.surface}>
            <span>{emoji}</span> {label}
          </button>
        ))}
      </div>
    </div>
  );
}

Object.assign(window, { Dashboard, Card, StatCard });
