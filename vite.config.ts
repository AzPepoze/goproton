import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import electron from "vite-plugin-electron/simple";
import path from "path";
import { fileURLToPath } from "url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));

export default defineConfig({
	root: ".",
	publicDir: "src/public",
	build: {
		outDir: "build/frontend",
		target: "node20",
	},
	plugins: [
		svelte(),
		electron({
			main: {
				entry: "src/electron/main.ts",
				vite: {
					build: {
						outDir: "build/electron",
					},
				},
			},
			preload: {
				input: "src/electron/preload.ts",
				vite: {
					build: {
						outDir: "build/electron",
						rollupOptions: {
							output: {
								format: "cjs",
								entryFileNames: "[name].js",
							},
						},
					},
				},
			},
			renderer: {},
		}),
	],
	resolve: {
		alias: {
			"@": path.resolve(__dirname, "./src/electron/frontend"),
		},
	},
	css: {
		preprocessorOptions: {
			scss: {
				api: "modern-compiler",
			},
		},
	},
});
