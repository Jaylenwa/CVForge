import React, { useRef, useEffect, useState } from 'react';
import { Bold, Italic, List, Wand2, Loader2, Undo, Redo } from 'lucide-react';
import { polishText } from '../../services/geminiService';
import { useLanguage } from '../../contexts/LanguageContext';

interface RichTextEditorProps {
  value: string;
  onChange: (value: string) => void;
  label?: string;
  aiContext?: string;
  className?: string;
  minRows?: number;
  maxHeight?: number;
  valueFormat?: 'text' | 'html';
  outputFormat?: 'text' | 'html';
}

const textToHtml = (text: string) => {
  const escaped = text.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
  const lines = escaped.split(/\r?\n/);
  return lines.map(l => `<div>${l || '<br>'}</div>`).join('');
};

const htmlToText = (html: string) => {
  const temp = document.createElement('div');
  temp.innerHTML = html;
  const out: string[] = [];
  const walk = (node: Node) => {
    if (node.nodeType === Node.TEXT_NODE) {
      out.push((node as Text).nodeValue || '');
      return;
    }
    const el = node as HTMLElement;
    const name = (el.tagName || '').toLowerCase();
    if (name === 'br') {
      out.push('\n');
      return;
    }
    if (name === 'li') {
      out.push('• ');
    }
    Array.from(node.childNodes).forEach(walk);
    if (name === 'div' || name === 'p' || name === 'li') {
      out.push('\n');
    }
  };
  Array.from(temp.childNodes).forEach(walk);
  return out.join('').replace(/\n{3,}/g, '\n\n').trim();
};

export const RichTextEditor: React.FC<RichTextEditorProps> = ({ value, onChange, label, aiContext, className = '', minRows = 4, maxHeight = 300, valueFormat = 'text', outputFormat = 'text' }) => {
  const editorRef = useRef<HTMLDivElement>(null);
  const [isEnhancing, setIsEnhancing] = useState(false);
  const [isComposing, setIsComposing] = useState(false);
  const { t } = useLanguage();

  useEffect(() => {
    const nextHtml = valueFormat === 'html' ? (value || '') : textToHtml(value || '');
    if (editorRef.current && editorRef.current.innerHTML !== nextHtml) {
      editorRef.current.innerHTML = nextHtml;
    }
  }, [value, valueFormat]);

  const handleInput = () => {
    if (editorRef.current) {
      const nextHtml = editorRef.current.innerHTML;
      if (isComposing) return;
      onChange(outputFormat === 'html' ? nextHtml : htmlToText(nextHtml));
    }
  };

  const handleCompositionStart = () => {
    setIsComposing(true);
  };

  const handleCompositionEnd = () => {
    setIsComposing(false);
    if (editorRef.current) {
      const nextHtml = editorRef.current.innerHTML;
      onChange(outputFormat === 'html' ? nextHtml : htmlToText(nextHtml));
    }
  };

  const sanitizeHtml = (html: string) => {
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

  const handlePaste = (e: React.ClipboardEvent<HTMLDivElement>) => {
    e.preventDefault();
    const html = e.clipboardData.getData('text/html');
    const text = e.clipboardData.getData('text/plain');
    const source = html || (text ? text.replace(/\r?\n/g, '<br/>') : '');
    const safe = sanitizeHtml(source);
    document.execCommand('insertHTML', false, safe);
    handleInput();
  };

  const execCommand = (command: string) => {
    document.execCommand(command, false);
    handleInput();
    editorRef.current?.focus();
  };

  const handleAiEnhance = async () => {
    if (!aiContext) return;
    setIsEnhancing(true);
    try {
      const currentHtml = editorRef.current?.innerHTML ?? (valueFormat === 'html' ? (value || '') : textToHtml(value || ''));
      const improved = await polishText(htmlToText(currentHtml));
      const nextHtml = textToHtml(improved);
      if (editorRef.current) {
        editorRef.current.innerHTML = nextHtml;
        handleInput();
      }
    } catch {
    } finally {
      setIsEnhancing(false);
    }
  };

  return (
    <div className={`mb-2 ${className}`}>
      <div className="flex justify-between items-center mb-1">
        {label && <label className="block text-sm font-medium text-slate-700">{label}</label>}
        {aiContext && (
          <button
            type="button"
            onClick={handleAiEnhance}
            disabled={isEnhancing}
            className="flex items-center text-xs text-indigo-600 hover:text-indigo-800 disabled:opacity-50 transition-colors"
          >
            {isEnhancing ? <Loader2 className="w-3 h-3 mr-1 animate-spin" /> : <Wand2 className="w-3 h-3 mr-1" />}
            {isEnhancing ? t('editor.ai.polishing') : t('editor.ai_polish')}
          </button>
        )}
      </div>
      <div className="border border-slate-300 rounded-lg overflow-hidden focus-within:ring-2 focus-within:ring-indigo-500 focus-within:border-indigo-500 transition-all bg-white">
        <div className="flex items-center space-x-1 border-b border-slate-200 bg-slate-50 p-1">
          <button type="button" onClick={(e) => { e.preventDefault(); execCommand('bold'); }} className="p-1.5 rounded hover:bg-slate-200 text-slate-600 transition-colors" title={t('rte.tooltip.bold')}>
            <Bold className="w-4 h-4" />
          </button>
          <button type="button" onClick={(e) => { e.preventDefault(); execCommand('italic'); }} className="p-1.5 rounded hover:bg-slate-200 text-slate-600 transition-colors" title={t('rte.tooltip.italic')}>
            <Italic className="w-4 h-4" />
          </button>
          <button type="button" onClick={(e) => { e.preventDefault(); execCommand('insertUnorderedList'); }} className="p-1.5 rounded hover:bg-slate-200 text-slate-600 transition-colors" title={t('rte.tooltip.bullets')}>
            <List className="w-4 h-4" />
          </button>
          <div className="w-px h-4 bg-slate-300 mx-2" />
          <button type="button" onClick={(e) => { e.preventDefault(); execCommand('undo'); }} className="p-1.5 rounded hover:bg-slate-200 text-slate-600 transition-colors" title={t('rte.tooltip.undo')}>
            <Undo className="w-4 h-4" />
          </button>
          <button type="button" onClick={(e) => { e.preventDefault(); execCommand('redo'); }} className="p-1.5 rounded hover:bg-slate-200 text-slate-600 transition-colors" title={t('rte.tooltip.redo')}>
            <Redo className="w-4 h-4" />
          </button>
        </div>
        <div
          ref={editorRef}
          className="p-3 outline-none prose prose-sm max-w-none text-slate-700 overflow-y-auto"
          style={{ minHeight: `${minRows * 24}px`, maxHeight }}
          contentEditable
          onInput={handleInput}
          onCompositionStart={handleCompositionStart}
          onCompositionEnd={handleCompositionEnd}
          onPaste={handlePaste}
          suppressContentEditableWarning
        />
      </div>
    </div>
  );
};

export default RichTextEditor;
