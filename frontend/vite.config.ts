import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    react(),
    tailwindcss(),
  ],
  build: {
    outDir: '../dist',
    emptyOutDir: true,
  },
  server: {
    port: 5175, // Run PPID on port 5175 to avoid conflict with port 5174
    proxy: {
      "/api": {
        target: "http://localhost:3003",
        changeOrigin: true,
      },
      "/uploads": {
        target: "http://localhost:3003",
        changeOrigin: true,
      },
    },
  }
})
