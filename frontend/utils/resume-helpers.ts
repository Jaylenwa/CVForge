import { ResumeData } from '../types';

export const getFontStack = (fontId: string): string => {
    switch (fontId) {
        case 'inter': return '"Inter", ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif';
        case 'roboto': return '"Roboto", "Helvetica Neue", Arial, sans-serif';
        case 'merriweather': return '"Merriweather", "Georgia", serif';
        case 'playfair': return '"Playfair Display", "Times New Roman", serif';
        case 'mono': return '"Roboto Mono", "Courier New", monospace';
        
        // Chinese Fonts
        case 'yahei': return '"Microsoft YaHei", "Heiti SC", "PingFang SC", sans-serif';
        case 'notosans': return '"Noto Sans SC", "Microsoft YaHei", sans-serif';
        case 'simsun': return '"SimSun", "Songti SC", "Noto Serif SC", serif';
        case 'kaiti': return '"KaiTi", "STKaiti", "Kai", serif';
        
        default: return '"Inter", sans-serif';
    }
};

export const getThemeStyles = (theme?: { Color?: string; Font?: string; Spacing?: string; FontSize?: string }) => {
    const fontFamily = getFontStack(theme?.Font || 'inter');

    const spacingMultiplier = {
        'compact': '0.85',
        'normal': '1',
        'spacious': '1.25'
    }[theme?.Spacing || 'normal'];

  const fontSize = (() => {
    const n = Number(theme?.FontSize);
    if (!Number.isFinite(n)) return 13;
    return Math.min(16, Math.max(12, Math.round(n)));
  })();

  return { fontFamily, spacingMultiplier, fontSize };
};

export const hasExtraPersonalInfo = (data: ResumeData) => {
    const p = (data.Personal || {}) as NonNullable<ResumeData['Personal']>;
    const customArr = (() => {
        try {
            const raw = p?.CustomInfo;
            if (raw) {
                const parsed = JSON.parse(raw);
                if (Array.isArray(parsed)) return parsed;
            }
        } catch {}
        return [];
    })();
    return !!(p?.Gender || p?.Age || p?.Degree || (customArr && customArr.length > 0));
};

export const sanitizeHtml = (html: string) => {
  const parser = new DOMParser();
  const doc = parser.parseFromString(html || '', 'text/html');
  const allowedTags = new Set(['b','strong','i','em','u','br','p','div','ul','ol','li','span','a']);
  const escapeText = (s: string) => s.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;');
  const wrap = (content: string, tags: string[]) => {
    let out = content;
    tags.forEach(tag => { out = `<${tag}>${out}</${tag}>`; });
    return out;
  };
  const sanitizeNode = (node: Node): string => {
    if (node.nodeType === Node.TEXT_NODE) {
      return escapeText(node.textContent || '');
    }
    if (node.nodeType !== Node.ELEMENT_NODE) return '';
    const el = node as HTMLElement;
    const tag = el.tagName.toLowerCase();
    if (!allowedTags.has(tag)) {
      let s = '';
      el.childNodes.forEach(child => { s += sanitizeNode(child); });
      return s;
    }
    let attrs = '';
    if (tag === 'a') {
      const raw = el.getAttribute('href') || '';
      const lower = raw.trim().toLowerCase();
      const hasProtocol = /^[a-z][a-z0-9+.-]*:/.test(lower);
      let ok = false;
      if (hasProtocol) {
        const proto = lower.split(':')[0];
        ok = ['http', 'https', 'mailto'].includes(proto);
      } else {
        ok = true;
      }
      if (ok) {
        attrs = ` href="${escapeText(raw)}" rel="noopener noreferrer nofollow"`;
      } else {
        let s = '';
        el.childNodes.forEach(child => { s += sanitizeNode(child); });
        return s;
      }
    }
    if (tag === 'br') return '<br/>';
    let content = '';
    el.childNodes.forEach(child => { content += sanitizeNode(child); });
    if (tag === 'span') {
      const style = (el.getAttribute('style') || '').toLowerCase();
      const tagsToWrap: string[] = [];
      if (/\bfont-weight\s*:\s*(bold|700|800|900)\b/.test(style)) tagsToWrap.push('strong');
      if (/\bfont-style\s*:\s*italic\b/.test(style)) tagsToWrap.push('em');
      if (/\btext-decoration\s*:\s*underline\b/.test(style)) tagsToWrap.push('u');
      if (tagsToWrap.length) return wrap(content, tagsToWrap);
      return content;
    }
    return `<${tag}${attrs}>${content}</${tag}>`;
  };
  let out = '';
  doc.body.childNodes.forEach(n => { out += sanitizeNode(n); });
  return out;
};
