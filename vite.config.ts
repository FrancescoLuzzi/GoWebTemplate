import { defineConfig } from 'vite';
import path from 'path';
import { fileURLToPath } from 'url';
import { dirname, resolve } from 'path';
import viteCompression from 'vite-plugin-compression';

// if in ESM context
const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

// https://vitejs.dev/config/
export default defineConfig({
  root: 'frontend_lib',
  build: {
    lib: {
      entry: resolve(__dirname, 'frontend_lib/main.ts'),
      name: 'fe',
      fileName: 'index',
      formats: ['es'],
    },
    outDir: '../public/assets',
    emptyOutDir: false,
  },
  plugins: [viteCompression({ deleteOriginFile: true })],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './frontend_lib'),
    },
  },
});
