import react from "@vitejs/plugin-react";
import fs from "fs";
import path from "path";
import { defineConfig } from "vite";

function loadAppEnv(appEnv: string): Record<string, string> {
  const envFile = path.resolve(__dirname, "env", `.env.${appEnv}`);
  if (!fs.existsSync(envFile)) return {};
  return Object.fromEntries(
    fs
      .readFileSync(envFile, "utf-8")
      .split("\n")
      .flatMap((line) => {
        const m = line.match(/^\s*([^#=\s][^=]*?)\s*=\s*(.*?)\s*$/);
        return m ? [[m[1], m[2]]] : [];
      })
  );
}

export default defineConfig(() => {
  const appEnv = process.env.APP_ENV ?? "local";
  const env = loadAppEnv(appEnv);

  const clientVars = Object.fromEntries(
    Object.entries(env)
      .filter(([k]) => k.startsWith("VITE_"))
      .map(([k, v]) => [`import.meta.env.${k}`, JSON.stringify(v)])
  );

  return {
    plugins: [react()],
    define: clientVars,
    resolve: {
      alias: { "@": path.resolve(__dirname, "src") },
    },
    server: {
      port: env.PORT ? parseInt(env.PORT) : 3001,
    },
  };
});
