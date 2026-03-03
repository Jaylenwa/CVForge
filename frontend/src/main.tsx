import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from '../App';
import { toast } from '../components/ui/Toast';
import { tGlobal } from '../contexts/LanguageContext';

const w = window as unknown as { fetch: typeof window.fetch; __cvforgeFetchPatched?: boolean };
if (!w.__cvforgeFetchPatched) {
  w.__cvforgeFetchPatched = true;
  const originalFetch = w.fetch.bind(window);
  let lastToastAt = 0;
  w.fetch = async (input: RequestInfo | URL, init?: RequestInit) => {
    const res = await originalFetch(input, init);
    if (res.status === 429) {
      const now = Date.now();
      if (now - lastToastAt > 1500) {
        lastToastAt = now;
        let retryAfterSeconds: number | null = null;
        const h = res.headers.get('Retry-After');
        if (h) {
          const n = Number.parseInt(h, 10);
          if (Number.isFinite(n) && n > 0) retryAfterSeconds = n;
        }
        if (retryAfterSeconds == null) {
          try {
            const data = (await res.clone().json()) as any;
            const v = Number(data?.retryAfterSeconds);
            if (Number.isFinite(v) && v > 0) retryAfterSeconds = v;
          } catch {
          }
        }
        toast(retryAfterSeconds ? tGlobal('toast.rateLimitedRetry', { seconds: retryAfterSeconds }) : tGlobal('toast.rateLimited'), 'error');
      }
    }
    return res;
  };
}

const rootElement = document.getElementById('root');
if (!rootElement) {
  throw new Error("Could not find root element to mount to");
}

const root = ReactDOM.createRoot(rootElement);
root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
