import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// 输出到 dist/，由后端 go:embed 打包。
export default defineConfig({
  plugins: [react()],
  base: '/',
  build: {
    outDir: 'dist',
  },
  server: {
    // 开发时把 /api 代理到后端，前端可独立热更新
    proxy: {
      '/api': 'http://localhost:8080',
    },
  },
})