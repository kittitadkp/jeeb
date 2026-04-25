# Jeeb App — UI Kit

A high-fidelity click-through prototype of the Jeeb personal management web app.

## Usage

Open `index.html` in a browser. Navigate between all major sections using the sidebar.

## Pages included

| Page | File |
|---|---|
| Dashboard | `Dashboard.jsx` |
| Workouts | `Workouts.jsx` |
| Study (with live timer) | `Study.jsx` |
| Sleep | `Sleep.jsx` |
| Finance + Budget | `Finance.jsx` |
| Calendar | `Calendar.jsx` |
| Settings | inline in `index.html` |

## Components

- `Layout.jsx` — `AppLayout`, `Header`, `Sidebar`, `Icon` helpers
- Each page file exports its page component to `window`

## Design tokens used

- Primary: `#2563EB` (blue-600)
- Neutrals: slate scale
- Radius: 8px cards, 12px modals
- Shadow: `0 1px 2px rgb(0 0 0/.05)` for cards
- Font: Inter (400/500/600/700) + JetBrains Mono for timers/numbers
- Icons: inline SVG (Lucide-style, 1.5px stroke)

## Notes

- No backend — all data is local state; forms work interactively
- localStorage persists the active page across reloads
- The study timer is fully functional
