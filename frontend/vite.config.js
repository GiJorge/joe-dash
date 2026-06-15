import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
  base: './', // <--- Add this line
  plugins: [vue(), tailwindcss()],
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:10000', // Matches your Go server port
        changeOrigin: true,
        secure: false,
      }
    }
  }
})