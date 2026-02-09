import path from 'path';
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig(() => {
    return {
      server: {
        port: 3000,
        host: '0.0.0.0',
        proxy: {
          '/api': {
            target: 'http://localhost:8080',
            changeOrigin: true,
          },
          '/public': {
            target: 'http://localhost:8080',
            changeOrigin: true,
          },
        },
        allowedHosts: ['host.docker.internal'],
      },
      optimizeDeps: {
        exclude: ['framer-motion'],
      },
      plugins: [react()],
      resolve: {
        alias: {
          '@': path.resolve(__dirname, '.'),
        },
        dedupe: ['react', 'react-dom']
      },
      build: {
        chunkSizeWarningLimit: 1200,
        rollupOptions: {
          output: {
            manualChunks: {
              react: ['react', 'react-dom', 'react-router-dom'],
              editor: ['@tinymce/tinymce-react', 'tinymce', '@wangeditor/editor', '@wangeditor/editor-for-react'],
              dnd: ['@dnd-kit/core', '@dnd-kit/sortable', '@dnd-kit/utilities'],
              motion: ['framer-motion'],
              icons: ['lucide-react'],
              charts: ['recharts'],
              genai: ['@google/genai'],
            },
          },
        },
      }
    };
});
