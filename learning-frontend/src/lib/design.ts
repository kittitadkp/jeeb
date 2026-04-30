export const C = {
  bg:           "var(--c-bg)",
  surface:      "var(--c-surface)",
  surface2:     "var(--c-surface2)",
  surface3:     "var(--c-surface3)",
  border:       "var(--c-border)",
  border2:      "var(--c-border2)",
  text:         "var(--c-text)",
  text2:        "var(--c-text2)",
  text3:        "var(--c-text3)",
  primary:      "#0ea5e9",
  shadow:       "var(--c-shadow)",
  shadowMd:     "var(--c-shadow-md)",
  modalShadow:  "var(--c-modal-shadow)",
  success:      "#16A34A",
  danger:       "#DC2626",
  warning:      "#D97706",
  dangerBg:     "var(--c-danger-bg)",
  dangerBorder: "var(--c-danger-border)",
} as const;

export const T = {
  xs:   11,
  sm:   12,
  base: 13,
  md:   15,
  lg:   20,
  xl:   22,
  "2xl": 28,
  "3xl": 32,
} as const;

export const W = {
  normal:  400,
  medium:  500,
  semi:    600,
  bold:    700,
} as const;

export const R = {
  sm:   6,
  md:   8,
  lg:   10,
  card: 14,
  full: 9999,
} as const;

export const S = {
  1:  4,
  2:  8,
  3:  12,
  4:  16,
  5:  20,
  6:  24,
  8:  32,
  10: 40,
} as const;
