import { ResumeData, ThemeConfig } from '../types';

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

export const getThemeStyles = (config?: ThemeConfig) => {
    const fontFamily = getFontStack(config?.fontFamily || 'inter');

    const spacingMultiplier = {
        'compact': '0.85',
        'normal': '1',
        'spacious': '1.25'
    }[config?.spacing || 'normal'];

  return { fontFamily, spacingMultiplier };
};

export const hasExtraPersonalInfo = (data: ResumeData) => {
    const p = data.personalInfo || ({} as ResumeData['personalInfo']);
    return !!(p.gender || p.age || p.maritalStatus || p.politicalStatus || p.birthplace || p.ethnicity || p.height || p.weight || (p.customInfo && p.customInfo.length > 0));
};

export const sanitizeHtml = (html: string) => {
  const parser = new DOMParser();
  const doc = parser.parseFromString(html || '', 'text/html');
  const allowedTags = new Set(['b','strong','i','em','u','br','p','div','ul','ol','li','span','a']);
  const escapeText = (s: string) => s.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;');
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
      try {
        const u = new URL(raw, window.location.origin);
        const proto = u.protocol.replace(':','');
        if (['http','https','mailto'].includes(proto)) {
          attrs = ` href="${escapeText(raw)}" rel="noopener noreferrer nofollow"`;
        }
      } catch {}
    }
    if (tag === 'br') return '<br/>';
    let content = '';
    el.childNodes.forEach(child => { content += sanitizeNode(child); });
    return `<${tag}${attrs}>${content}</${tag}>`;
  };
  let out = '';
  doc.body.childNodes.forEach(n => { out += sanitizeNode(n); });
  return out;
};
