// Standalone mount for BulkLogModal — separate React root, bypasses cache

(function() {
  // Create portal container
  const portalDiv = document.createElement('div');
  portalDiv.id = 'bulk-portal';
  document.body.appendChild(portalDiv);

  function BulkApp() {
    const [open, setOpen] = React.useState(false);
    const [triggerReady, setTriggerReady] = React.useState(false);

    React.useEffect(() => {
      // Expose opener globally
      window.__openBulkLog = () => setOpen(true);
      document.addEventListener('jeeb:openBulkLog', () => setOpen(true));
      setTriggerReady(true);

      // Inject sidebar button once aside is ready
      function injectSidebarBtn() {
        const aside = document.querySelector('aside');
        if (!aside || aside.querySelector('[data-bulk-btn]')) return;
        const userRow = aside.lastElementChild;
        const sep = document.createElement('div');
        sep.style.cssText = 'height:1px;background:rgba(255,255,255,0.07);margin:6px 8px';
        const wrap = document.createElement('div');
        wrap.setAttribute('data-bulk-btn','1');
        wrap.style.cssText = 'padding:0 4px;margin-bottom:4px';
        const btn = document.createElement('button');
        btn.style.cssText = [
          'width:100%', 'display:flex', 'align-items:center', 'gap:9px',
          'padding:9px 12px', 'border-radius:10px', 'cursor:pointer',
          'transition:all 0.15s', 'text-align:left', 'font-family:inherit',
          'background:linear-gradient(135deg,rgba(124,110,245,0.12),rgba(167,139,250,0.12))',
          'border:1px solid rgba(124,110,245,0.28)'
        ].join(';');
        btn.innerHTML = '<span style="font-size:15px">📋</span>'
          + '<span style="font-size:13px;font-weight:600;background:linear-gradient(135deg,#7c6ef5,#a78bfa);-webkit-background-clip:text;-webkit-text-fill-color:transparent;letter-spacing:0.01em">Bulk Log</span>';
        btn.onmouseenter = () => { btn.style.background = 'linear-gradient(135deg,rgba(124,110,245,0.24),rgba(167,139,250,0.24))'; };
        btn.onmouseleave = () => { btn.style.background = 'linear-gradient(135deg,rgba(124,110,245,0.12),rgba(167,139,250,0.12))'; };
        btn.onclick = () => setOpen(true);
        wrap.appendChild(btn);
        aside.insertBefore(sep, userRow);
        aside.insertBefore(wrap, userRow);
      }

      // Try injecting with retries — React re-renders can wipe injected nodes
      let tries = 0;
      function tryInject() {
        injectSidebarBtn();
        if (++tries < 10) setTimeout(tryInject, 500);
      }
      setTimeout(tryInject, 800);

      // Re-inject on any nav click (React re-renders the aside)
      document.addEventListener('click', () => {
        setTimeout(injectSidebarBtn, 200);
      });

      return () => { delete window.__openBulkLog; };
    }, []);

    return React.createElement(React.Fragment, null,
      React.createElement(BulkLogModal, { open, onClose: () => setOpen(false) })
    );
  }

  ReactDOM.createRoot(portalDiv).render(React.createElement(BulkApp));
})();
