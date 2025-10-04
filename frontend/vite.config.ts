import { defineConfig, loadEnv } from "vite";
import react from "@vitejs/plugin-react";
import path from "path";

type ProjectKey = "sso" | "portal";

interface ProjectConfig {
  root: string;
  port: number;
  outDir: string;
}

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), "");

  const projects: Record<ProjectKey, ProjectConfig> = {
    sso: {
      root: path.resolve(__dirname, "packages/sso"),
      port: 5173,
      outDir: "dist/sso",
    },
    portal: {
      root: path.resolve(__dirname, "packages/portal"),
      port: 5174,
      outDir: "dist/portal",
    },
  };

  const project = (process.env.VITE_PROJECT || "sso") as ProjectKey;

  if (!(project in projects)) {
    throw new Error(`Invalid VITE_PROJECT: ${project}. Must be 'sso' or 'portal'.`);
  }

  const aliases = {
    "@common": path.resolve(__dirname, "packages/common/src"),
    "@": path.resolve(__dirname, `packages/${project}/src`),
  };

  return {
    plugins: [react()],
    resolve: {
      alias: aliases,
    },
    root: projects[project].root,
    base: "/",
    server: {
      port: projects[project].port,
      open: true,
      proxy: {
        "/api": {
          target: env.VITE_API_SERVER,
          changeOrigin: true,
        },
      },
    },
    envDir: "../..",
    build: {
      outDir: projects[project].outDir,
      rollupOptions: {
        input: {
          main: path.resolve(projects[project].root, "index.html"),
        },
      },
    },
  };
});
