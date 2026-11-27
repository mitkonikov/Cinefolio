import adapter from '@sveltejs/adapter-static';
import { sveltekit } from '@sveltejs/kit/vite';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';
import { defineConfig } from 'vitest/config';

export default defineConfig({
	plugins: [sveltekit({ adapter: adapter(), preprocess: vitePreprocess() })],
	server: {
		host: '127.0.0.1',
		strictPort: true,
		proxy: {
			'/api': 'http://127.0.0.1:5002',
			'/media': 'http://127.0.0.1:5002'
		}
	},
	test: {
		include: ['src/**/*.{test,spec}.{js,ts}']
	}
});
