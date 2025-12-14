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

export const RichTextEditor: React.FC<RichTextEditorProps> = ({ value, onChange, label, aiContext, className = '', minRows = 4, maxHeight = 300 }) => {
  const editorRef = useRef<HTMLDivElement>(null);
  const [isEnhancing, setIsEnhancing] = useState(false);
  const [htmlContent, setHtmlContent] = useState(textToHtml(value || ''));
  const { t } = useLanguage();

  useEffect(() => {
    const nextHtml = textToHtml(value || '');
    if (editorRef.current && editorRef.current.innerHTML !== nextHtml) {
      editorRef.current.innerHTML = nextHtml;
      setHtmlContent(nextHtml);
    }
  }, [value]);

  const handleInput = () => {
    if (editorRef.current) {
      const nextHtml = editorRef.current.innerHTML;
      setHtmlContent(nextHtml);
      onChange(htmlToText(nextHtml));
    }
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
      const improved = await polishText(htmlToText(htmlContent));
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
          <button type="button" onClick={(e) => { e.preventDefault(); execCommand('bold'); }} className="p-1.5 rounded hover:bg-slate-200 text-slate-600 transition-colors" title="加粗">
            <Bold className="w-4 h-4" />
          </button>
          <button type="button" onClick={(e) => { e.preventDefault(); execCommand('italic'); }} className="p-1.5 rounded hover:bg-slate-200 text-slate-600 transition-colors" title="斜体">
            <Italic className="w-4 h-4" />
          </button>
          <button type="button" onClick={(e) => { e.preventDefault(); execCommand('insertUnorderedList'); }} className="p-1.5 rounded hover:bg-slate-200 text-slate-600 transition-colors" title="项目符号">
            <List className="w-4 h-4" />
          </button>
          <div className="w-px h-4 bg-slate-300 mx-2" />
          <button type="button" onClick={(e) => { e.preventDefault(); execCommand('undo'); }} className="p-1.5 rounded hover:bg-slate-200 text-slate-600 transition-colors" title="撤销">
            <Undo className="w-4 h-4" />
          </button>
          <button type="button" onClick={(e) => { e.preventDefault(); execCommand('redo'); }} className="p-1.5 rounded hover:bg-slate-200 text-slate-600 transition-colors" title="重做">
            <Redo className="w-4 h-4" />
          </button>
        </div>
        <div
          ref={editorRef}
          className="p-3 outline-none prose prose-sm max-w-none text-slate-700 overflow-y-auto"
          style={{ minHeight: `${minRows * 24}px`, maxHeight }}
          contentEditable
          onInput={handleInput}
          suppressContentEditableWarning
          dangerouslySetInnerHTML={{ __html: htmlContent }}
        />
      </div>
    </div>
  );
};

export default RichTextEditor;
