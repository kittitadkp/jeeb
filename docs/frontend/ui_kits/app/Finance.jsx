// Jeeb UI Kit — Finance Page v2

const { useState } = React;

const TRANSACTIONS = [
  { id: 1, type: 'expense', category: 'Food & Dining', name: 'Groceries',      amount: -85.50, date: 'Apr 22' },
  { id: 2, type: 'income',  category: 'Salary',        name: 'Monthly salary', amount: 5000,   date: 'Apr 20' },
  { id: 3, type: 'expense', category: 'Transport',     name: 'Grab ride',      amount: -18.00, date: 'Apr 20' },
  { id: 4, type: 'expense', category: 'Utilities',     name: 'Electricity bill',amount: -120.00,date: 'Apr 19' },
  { id: 5, type: 'expense', category: 'Food & Dining', name: 'Restaurant',     amount: -62.00, date: 'Apr 18' },
  { id: 6, type: 'expense', category: 'Entertainment', name: 'Netflix',        amount: -15.99, date: 'Apr 15' },
  { id: 8, type: 'income',  category: 'Freelance',     name: 'Design project', amount: 800,    date: 'Apr 10' },
];
const BUDGETS = [
  { category: 'Food & Dining', spent: 450, limit: 600 },
  { category: 'Transport',     spent: 80,  limit: 150 },
  { category: 'Utilities',     spent: 220, limit: 250 },
  { category: 'Entertainment', spent: 60,  limit: 100 },
];
const CAT_EMOJI = { 'Food & Dining':'🛒', Salary:'💼', Transport:'🚗', Utilities:'💡', Entertainment:'🎬', Freelance:'💻' };

function TransactionItem({ t: tx, onDelete }) {
  const t = window.useTheme();
  const isExpense = tx.type === 'expense';
  return (
    <div style={{ display: 'flex', alignItems: 'center', gap: 12, padding: '12px 16px', background: t.surface, border: `1px solid ${t.border}`, borderRadius: 8 }}>
      <div style={{ width: 36, height: 36, borderRadius: '50%', background: isExpense ? t.dangerSubtle : t.successSubtle, display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: 16, flexShrink: 0 }}>
        {CAT_EMOJI[tx.category] || '💳'}
      </div>
      <div style={{ flex: 1 }}>
        <div style={{ fontSize: 14, fontWeight: 500, color: t.fg1 }}>{tx.name}</div>
        <div style={{ fontSize: 12, color: t.fg2 }}>{tx.category} · {tx.date}</div>
      </div>
      <div style={{ fontSize: 15, fontWeight: 600, color: isExpense ? t.dangerText : t.successText }}>
        {isExpense ? '-' : '+'}${Math.abs(tx.amount).toFixed(2)}
      </div>
      <button onClick={() => onDelete(tx.id)} style={{ background: 'none', border: 'none', cursor: 'pointer', color: t.fg3, fontSize: 18 }}
        onMouseEnter={e => e.currentTarget.style.color = t.dangerText}
        onMouseLeave={e => e.currentTarget.style.color = t.fg3}>×</button>
    </div>
  );
}

function BudgetBar({ category, spent, limit }) {
  const t = window.useTheme();
  const pct = Math.min((spent / limit) * 100, 100);
  const barColor = pct > 90 ? '#DC2626' : pct > 75 ? '#F59E0B' : '#2563EB';
  return (
    <div style={{ marginBottom: 14 }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 4 }}>
        <span style={{ fontSize: 13, fontWeight: 500, color: t.fg1 }}>{category}</span>
        <span style={{ fontSize: 12, color: t.fg2 }}>${spent} / ${limit}</span>
      </div>
      <div style={{ height: 8, background: t.surfaceActive, borderRadius: 9999 }}>
        <div style={{ height: '100%', borderRadius: 9999, background: barColor, width: `${pct}%` }} />
      </div>
      <div style={{ fontSize: 11, color: t.fg3, marginTop: 2 }}>${limit - spent} remaining</div>
    </div>
  );
}

function AddTransactionForm({ onClose, onAdd }) {
  const t = window.useTheme();
  const [form, setForm] = useState({ type: 'expense', amount: '', category: 'Food & Dining', note: '' });
  const set = (k, v) => setForm(f => ({ ...f, [k]: v }));
  const inputStyle = { width: '100%', border: `1px solid ${t.border}`, borderRadius: 8, padding: '8px 12px', fontSize: 14, fontFamily: 'inherit', outline: 'none', background: t.inputBg, color: t.fg1 };
  return (
    <div style={{ position: 'fixed', inset: 0, background: 'rgb(0 0 0/.5)', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 200 }}>
      <div style={{ background: t.surface, borderRadius: 12, width: 400, boxShadow: t.modalShadow }}>
        <div style={{ padding: '16px 20px', borderBottom: `1px solid ${t.border}`, display: 'flex', justifyContent: 'space-between' }}>
          <h2 style={{ fontSize: 16, fontWeight: 600, color: t.fg1 }}>New Transaction</h2>
          <button onClick={onClose} style={{ background: 'none', border: 'none', cursor: 'pointer', color: t.fg3, fontSize: 20 }}>×</button>
        </div>
        <div style={{ padding: 20, display: 'flex', flexDirection: 'column', gap: 14 }}>
          <div style={{ display: 'flex', gap: 8 }}>
            {['expense','income'].map(type => (
              <label key={type} style={{ flex: 1, display: 'flex', alignItems: 'center', gap: 6, padding: '8px 12px', border: `1px solid ${form.type === type ? '#2563EB' : t.border}`, borderRadius: 8, cursor: 'pointer', background: form.type === type ? t.primarySubtle : t.surface }}>
                <input type="radio" name="type" checked={form.type === type} onChange={() => set('type', type)} style={{ accentColor: '#2563EB' }} />
                <span style={{ fontSize: 13, fontWeight: 500, color: form.type === type ? '#2563EB' : t.fg2, textTransform: 'capitalize' }}>{type}</span>
              </label>
            ))}
          </div>
          <div>
            <label style={{ fontSize: 13, fontWeight: 500, color: t.fg1, display: 'block', marginBottom: 4 }}>Amount <span style={{ color: '#DC2626' }}>*</span></label>
            <input type="number" placeholder="0.00" value={form.amount} onChange={e => set('amount', e.target.value)} style={inputStyle} />
          </div>
          <div>
            <label style={{ fontSize: 13, fontWeight: 500, color: t.fg1, display: 'block', marginBottom: 4 }}>Category</label>
            <select value={form.category} onChange={e => set('category', e.target.value)} style={inputStyle}>
              <option>Food &amp; Dining</option><option>Transport</option><option>Utilities</option><option>Entertainment</option><option>Salary</option><option>Freelance</option>
            </select>
          </div>
          <div>
            <label style={{ fontSize: 13, fontWeight: 500, color: t.fg1, display: 'block', marginBottom: 4 }}>Note</label>
            <input type="text" placeholder="What was this for?" value={form.note} onChange={e => set('note', e.target.value)} style={inputStyle} />
          </div>
        </div>
        <div style={{ padding: '14px 20px', borderTop: `1px solid ${t.border}`, display: 'flex', justifyContent: 'flex-end', gap: 8 }}>
          <button onClick={onClose} style={{ padding: '8px 16px', border: `1px solid ${t.border}`, borderRadius: 8, background: t.surface, fontSize: 14, cursor: 'pointer', color: t.fg2 }}>Cancel</button>
          <button onClick={() => { if (form.amount) { onAdd({ id: Date.now(), type: form.type, category: form.category, name: form.note || form.category, amount: form.type === 'expense' ? -parseFloat(form.amount) : parseFloat(form.amount), date: 'Today' }); onClose(); } }}
            style={{ padding: '8px 18px', border: 'none', borderRadius: 8, background: '#2563EB', color: '#fff', fontSize: 14, fontWeight: 500, cursor: 'pointer' }}>Save</button>
        </div>
      </div>
    </div>
  );
}

function Finance() {
  const t = window.useTheme();
  const [transactions, setTransactions] = useState(TRANSACTIONS);
  const [showForm, setShowForm] = useState(false);
  const [view, setView] = useState('transactions');
  const income   = transactions.filter(t => t.type === 'income').reduce((s, t) => s + t.amount, 0);
  const expenses = transactions.filter(t => t.type === 'expense').reduce((s, t) => s + Math.abs(t.amount), 0);
  const balance  = income - expenses;

  return (
    <div style={{ maxWidth: 760 }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
        <h1 style={{ fontSize: 22, fontWeight: 700, color: t.fg1 }}>Finance</h1>
        <div style={{ display: 'flex', gap: 8 }}>
          <button onClick={() => setView(v => v === 'budget' ? 'transactions' : 'budget')} style={{ padding: '8px 14px', border: `1px solid ${t.border}`, borderRadius: 8, background: t.surface, fontSize: 13, cursor: 'pointer', color: t.fg2 }}>
            {view === 'budget' ? 'Transactions' : 'Budget'}
          </button>
          <button onClick={() => setShowForm(true)} style={{ padding: '8px 16px', background: '#2563EB', color: '#fff', border: 'none', borderRadius: 8, fontSize: 14, fontWeight: 500, cursor: 'pointer' }}>+ Add</button>
        </div>
      </div>
      <div style={{ display: 'flex', gap: 12, marginBottom: 16 }}>
        {[['Balance', balance, balance >= 0 ? t.successText : t.dangerText, balance >= 0 ? '+' : ''],
          ['Income', income, t.successText, '+'],
          ['Expenses', expenses, t.dangerText, '-']].map(([label, val, color, prefix]) => (
          <div key={label} style={{ flex: 1, background: t.surface, border: `1px solid ${t.border}`, borderRadius: 8, padding: '14px 16px' }}>
            <div style={{ fontSize: 12, color: t.fg2, marginBottom: 4 }}>{label}</div>
            <div style={{ fontSize: 20, fontWeight: 700, color }}>{prefix}${Math.abs(val).toFixed(2)}</div>
          </div>
        ))}
      </div>
      {view === 'transactions' ? (
        <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
          {transactions.map(tx => <TransactionItem key={tx.id} t={tx} onDelete={id => setTransactions(ts => ts.filter(t => t.id !== id))} />)}
        </div>
      ) : (
        <div style={{ background: t.surface, border: `1px solid ${t.border}`, borderRadius: 8, padding: 20 }}>
          <div style={{ fontSize: 15, fontWeight: 600, color: t.fg1, marginBottom: 4 }}>Budget — April 2024</div>
          <div style={{ fontSize: 13, color: t.fg2, marginBottom: 16 }}>Total: ${expenses.toFixed(0)} spent</div>
          {BUDGETS.map(b => <BudgetBar key={b.category} {...b} />)}
        </div>
      )}
      {showForm && <AddTransactionForm onClose={() => setShowForm(false)} onAdd={tx => setTransactions(ts => [tx, ...ts])} />}
    </div>
  );
}

Object.assign(window, { Finance });
