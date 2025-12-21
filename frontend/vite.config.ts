import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	build: {
		minify: 'esbuild',
		sourcemap: false,
		rollupOptions: {
			output: {
				manualChunks: undefined
			}
		}
	},
	plugins: [tailwindcss(), sveltekit()],
	server: {
		proxy: {
			'/api/v1': {
				target: 'http://localhost:80',
				changeOrigin: true
			}
		}
	}
});
