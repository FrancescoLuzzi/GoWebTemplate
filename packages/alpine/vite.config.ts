import { defineConfig } from 'vite';
import { fileURLToPath } from 'url';
import { dirname, resolve } from 'path';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

// https://vitejs.dev/config/
export default defineConfig({
  root: '.',
  build: {
    lib: {
      entry: resolve(__dirname, './src/index.ts'),
      name: 'alpine',
      fileName: 'alpine',
      formats: ['es'],
    },
    outDir: 'public/assets',
    copyPublicDir: false,
    emptyOutDir: false,
  },
  resolve: {
    alias: {
      '@': resolve(__dirname, './src'),
    },
  },
});
