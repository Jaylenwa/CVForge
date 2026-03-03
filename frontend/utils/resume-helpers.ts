import { ResumeData, ResumeItem, ResumeSection, ResumeSectionType } from '../types';

export const getFontStack = (fontId: string): string => {
    switch (fontId) {
        case 'notosans': return '"Noto Sans SC", "Noto Sans CJK SC", "Microsoft YaHei", "PingFang SC", "Hiragino Sans GB", sans-serif';
        default: return '"Noto Sans SC", "Noto Sans CJK SC", "Microsoft YaHei", "PingFang SC", "Hiragino Sans GB", sans-serif';
    }
};

export const getThemeStyles = (theme?: { Color?: string; Font?: string; Spacing?: string; FontSize?: string }) => {
    const fontFamily = getFontStack(theme?.Font || 'notosans');

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

const withUniqueItemIDs = (items: ResumeItem[], prefix: string): ResumeItem[] => {
  const seen = new Set<string>();
  return (items || []).map((it, idx) => {
    let id = it?.id == null ? '' : String(it.id);
    if (!id || seen.has(id)) id = `${prefix}${idx}`;
    seen.add(id);
    return { ...it, id, orderNum: idx };
  });
};

export const normalizeResumeDataForRender = (data: ResumeData): ResumeData => {
  const sections = Array.isArray(data.sections) ? (data.sections as ResumeSection[]) : [];
  const selfIndex = sections.findIndex((s) => s?.type === ResumeSectionType.SelfEvaluation);
  if (selfIndex < 0) return data;

  const self = sections[selfIndex];
  const mergedItems = withUniqueItemIDs((Array.isArray(self.items) ? (self.items as ResumeItem[]) : []), 'selfEvaluation-item-');
  const nextSections = sections.map((s, idx) => idx === selfIndex ? { ...self, items: mergedItems } : s);
  return { ...data, sections: nextSections };
};
