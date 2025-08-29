import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import path from "path";

type ProjectKey = "sso" | "portal";

interface ProjectConfig {
	root: string;
	port: number;
	outDir: string;
}

export default defineConfig(() => {
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
		throw new Error(
			`Invalid VITE_PROJECT: ${project}. Must be 'sso' or 'portal'.`,
		);
	}

	return {
		plugins: [react()],
		resolve: {
			alias: {
				"@frontend/shared": path.resolve(
					__dirname,
					"packages/shared/src",
				),
			},
		},
		root: projects[project].root,
		base: "/",
		server: {
			port: projects[project].port,
			open: true,
		},
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
