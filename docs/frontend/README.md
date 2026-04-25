# Jeeb Design System

## About Jeeb

**Jeeb** is a personal management web application that helps individuals track and manage their daily life across four core domains: **workouts**, **study sessions**, **sleep**, and **personal finance** — with a unified calendar and notification layer tying everything together.

It is a single-product, single-surface web app (React + Vite SPA), styled with **Shadcn/ui** on top of **Tailwind CSS**, and backed by a Go API + MongoDB, with Keycloak for authentication.

### Products / Surfaces

| Surface | Description |
|---|---|
| **Web App** | Primary product — a desktop/tablet-first React SPA with responsive mobile fallback |

### Source materials used

- **Codebase**: `jeeb/` (mounted locally) — Go backend + docs; no actual frontend source code was present; design was reconstructed from detailed frontend documentation
- **Docs**: `jeeb/docs/frontend/spec.md`, `components.md`, `pages.md`

---

## CONTENT FUNDAMENTALS

### Tone & Voice
- **Personal and direct** — the app is just for the user; copy speaks directly to them (e.g. "Good morning, John 👋")
- **Practical and minimal** — no marketing fluff; every label earns its place
- **Friendly but not chatty** — short, clear labels; no long explanations in the UI
- **First-person welcome** — greetings use "Good morning/evening"; never robotic

### Casing
- **Title case** for page headings and nav items: "Workout List", "Sleep Log"
- **Sentence case** for descriptions and helper text: "Start tracking your fitness journey"
- **ALL CAPS** avoided entirely

### Copy examples
- "Good morning, John 👋" (dashboard greeting)
- "+ Add New" / "+ Log" (action buttons — brief and action-first)
- "No workouts yet" (empty state headline — plain, lowercase-second-word)
- "Start tracking your fitness journey" (empty state description)
- "Showing 1–10 of 50" (pagination — plain and informational)
- "This week", "This month" (stat card context labels)

### Emoji usage
- Used sparingly for **section icons** in sidebar and stat cards only (💪 📚 😴 💰 📅)
- Not used decoratively in body text or buttons
- Star ratings use ⭐ for sleep quality

### Numbers & Data
- Currency: `$1,250` format (dollar sign, comma-separated)
- Duration: "7h 30m" or "45 min" (abbreviated)
- Dates: "Apr 20, 2024" or relative ("2h ago", "Yesterday", "Tomorrow 2PM")

---

## VISUAL FOUNDATIONS

### Color System
- **Primary**: Blue-600 (`#2563EB`) — interactive elements, active states, CTA buttons
- **Secondary**: Slate-600 (`#475569`) — secondary text, subdued labels
- **Success**: Green-600 (`#16A34A`) — positive trends, income, completed states
- **Warning**: Amber-500 (`#F59E0B`) — caution states, budget approaching limit
- **Danger**: Red-600 (`#DC2626`) — destructive actions, errors, expenses
- **Background (light)**: Slate-50 (`#F8FAFC`) — page background
- **Background (dark)**: Slate-950 (`#020617`) — dark mode page
- **Surface (light)**: White (`#FFFFFF`) — cards and panels
- **Surface (dark)**: Slate-900 (`#0F172A`) — dark mode cards

### Typography
- **Font family**: Inter (system fallback — no custom font files in codebase; Google Fonts substitute used)
- **Scale**: Standard Tailwind sizing (sm=14px, base=16px, lg=18px, xl=20px, 2xl=24px, 3xl=30px)
- **Weight usage**: 400 for body, 500 for labels/meta, 600 for card headings, 700 for page titles
- **Mono**: `ui-monospace` / JetBrains Mono — used for timer displays (e.g. "01:45:30"), numeric values

### Spacing
- **Base unit**: 4px (Tailwind default)
- **Component padding**: p-4 (16px) for cards; p-6 (24px) for page content
- **Stack gap**: gap-4 (16px) standard; gap-6 (24px) between sections
- **Inline gap**: gap-2 (8px) within form rows, tag groups

### Backgrounds
- Flat color backgrounds only — no gradients, no textures, no illustrations
- Page bg: slate-50 in light mode
- Cards: white with subtle `shadow-sm`
- No full-bleed images in the app UI

### Cards
- **Border radius**: `rounded-lg` = 8px (all cards, modals, inputs)
- **Shadow**: `shadow-sm` for cards (soft, barely-there); `shadow-lg` for modals/popovers
- **Border**: no explicit border by default; sometimes `border border-slate-200` for definition
- **Padding**: p-4 (16px) inner padding standard

### Borders & Dividers
- `border-slate-200` for light mode dividers and card borders
- `border-slate-700` for dark mode
- Dividers use `<hr>` or `border-t border-slate-200`

### Hover & Focus States
- Buttons: darken (e.g. `hover:bg-blue-700`)
- Ghost/outline buttons: subtle bg fill on hover (`hover:bg-slate-100`)
- Nav items: `hover:bg-slate-100` background highlight
- Active nav item: `bg-blue-50 text-blue-600` with `font-medium`
- Focus rings: standard browser focus outline (Tailwind `focus-visible:ring-2 ring-blue-500`)

### Press States
- Buttons subtly compress (`active:scale-[0.98]`)
- No bounce animations

### Animations
- **Minimal** — Shadcn/ui default transitions (150–200ms ease-out)
- Modals/dialogs: fade + subtle scale-in
- No bounces, no spring physics in the UI itself
- Toast notifications: slide-in from bottom

### Corner Radii
- `rounded-lg` (8px) — default for cards, inputs, buttons, badges
- `rounded-full` — avatar, circular icon buttons, pill badges
- `rounded-xl` (12px) — modals only

### Shadow System
- `shadow-sm` — cards, stat boxes (subtle depth)
- `shadow-md` — dropdowns, select menus
- `shadow-lg` — modals, dialogs, popovers

### Use of Transparency & Blur
- No frosted glass / backdrop blur
- Overlays: `bg-black/50` backdrop behind modals

### Layout
- **Sidebar**: fixed left sidebar (240px wide on desktop), collapses to icon-only on tablet
- **Header**: fixed top bar with logo, search, notifications, avatar
- **Content**: flex-1, p-6 padding
- **Mobile**: sidebar becomes a bottom tab bar

### Color vibe of imagery
- No decorative imagery in the app shell
- Data visualizations use brand colors (blue for primary data, slate for secondary)
- Charts use blue-600 bars/lines by default; categorical colors from the semantic palette

---

## ICONOGRAPHY

Jeeb uses **Lucide React** as its icon system exclusively.

- **Style**: Stroke-based, 24px default, 1.5px stroke weight, rounded line caps
- **Usage**: Nav items (Dashboard, Workouts, etc.), action buttons (edit, delete, add), status indicators
- **Color**: Inherits text color — typically `text-slate-500` for decorative; `text-blue-600` for active
- **Emoji as icons**: Only used alongside sidebar nav labels and stat cards to add personality (💪 📚 😴 💰 📅 📅)
- **No custom icon font or SVG sprite** — all icons are inline React components from `lucide-react`
- **CDN**: For the design system, Lucide is loaded from `https://unpkg.com/lucide@latest/dist/umd/lucide.min.js`

Key icons in use:
| Context | Icon |
|---|---|
| Dashboard | `LayoutDashboard` |
| Workouts | `Dumbbell` |
| Study | `BookOpen` |
| Sleep | `Moon` |
| Finance | `Wallet` |
| Calendar | `Calendar` |
| Settings | `Settings` |
| Notifications | `Bell` |
| Add | `Plus` |
| Edit | `Pencil` |
| Delete | `Trash2` |
| Search | `Search` |
| User | `User` |

No logos, brand illustrations, or background images were found in the codebase. The app uses text-based identity ("Jeeb") with blue accent.

---

## File Index

```
README.md                    ← This file
SKILL.md                     ← Agent skill descriptor
colors_and_type.css          ← CSS custom properties for colors & typography
preview/
  colors-primary.html        ← Blue color scale
  colors-neutral.html        ← Slate neutral scale
  colors-semantic.html       ← Success / warning / danger
  type-scale.html            ← Heading type scale
  type-body.html             ← Body, mono, caption specimens
  spacing.html               ← Spacing & border radius tokens
  shadows.html               ← Shadow system
  components-buttons.html    ← Button variants & states
  components-cards.html      ← Card, StatCard components
  components-forms.html      ← Form fields, selects, toggles
  components-badges.html     ← Badges, status indicators
  components-nav.html        ← Sidebar & header navigation
assets/                      ← (no external assets found in codebase)
ui_kits/
  app/
    README.md                ← UI kit usage notes
    index.html               ← Interactive app prototype (Dashboard → all features)
    Layout.jsx               ← AppLayout, Header, Sidebar
    Dashboard.jsx            ← Dashboard page
    Workouts.jsx             ← Workouts list + form
    Study.jsx                ← Study sessions + timer
    Sleep.jsx                ← Sleep log + chart
    Finance.jsx              ← Finance + budget
    Calendar.jsx             ← Calendar view
```
