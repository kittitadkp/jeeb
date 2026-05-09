import { defineConfig } from "tsup";

export default defineConfig({
  entry: {
    index: "src/index.ts",
    "ui/index": "src/ui/index.ts",
    "charts/index": "src/charts/index.ts",
    "auth/index": "src/auth/index.ts",
    "utils/index": "src/utils/index.ts",
  },
  format: ["esm", "cjs"],
  dts: true,
  sourcemap: true,
  external: ["react", "react-dom", "tailwindcss"],
  treeshake: true,
  clean: true,
});
