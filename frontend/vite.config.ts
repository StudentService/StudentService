import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tsconfigPaths from 'vite-tsconfig-paths'

export default defineConfig({
  plugins: [
    react(),
    tsconfigPaths(),          // ← синхронизирует пути из tsconfig.json
  ],
  assetsInclude: ['**/*.html'],  // ← можно оставить, но обычно не нужно
})
