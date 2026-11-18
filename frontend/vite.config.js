import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'
import {defineConfig} from 'vite'

export default defineConfig({
    base: "/",
    plugins: [react(), tailwindcss()],
    preview: {
        port: 8080,
        strictPort: true,
    },
    server: {
        port: 8080,
        strictPort: true,
        host: true,
    },
})
