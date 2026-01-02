import React, { useEffect, useMemo, useRef, useState } from 'react';
import { Wand2, Loader2 } from 'lucide-react';
import { Editor, Toolbar } from '@wangeditor/editor-for-react';
import '@wangeditor/editor/dist/css/style.css';
import { polishText } from '../../services/geminiService';
import { useLanguage } from '../../contexts/LanguageContext';
import { sanitizeHtml } from '../../utils/resume-helpers';

interface RichTextEditorProps {
  value: string;
  onChange: (value: string) => void;
  label?: string;
  aiContext?: string;
  className?: string;
  minRows?: number;
  maxHeight?: number;
}

const htmlToText = (html: string) => {
  const d = document.createElement('div');
  d.innerHTML = html || '';
  return (d.textContent || '').replace(/\u00a0/g, ' ').trim();
};

export const RichTextEditor: React.FC<RichTextEditorProps> = ({ value, onChange, label, aiContext, className = '', minRows = 4, maxHeight = 300 }) => {
  const editorRef = useRef<any>(null);
  const [isEnhancing, setIsEnhancing] = useState(false);
  const { t } = useLanguage();
  const [html, setHtml] = useState<string>(value || '');
  const [editor, setEditor] = useState<any>(null);

  useEffect(() => {
    setHtml(value || '');
  }, [value]);

  const toolbarConfig = useMemo(() => ({
    toolbarKeys: ['bold', 'italic', 'underline', 'bulletedList', 'numberedList', 'link', 'undo', 'redo'],
  }), []);

  const editorConfig = useMemo(() => ({
    placeholder: '',
    autoFocus: false,
  }), []);

  const handleAiEnhance = async () => {
    if (!aiContext) return;
    setIsEnhancing(true);
    try {
      const currentText = editorRef.current?.getText?.() ?? htmlToText(html || '');
      const improved = await polishText(currentText);
      const nextHtml = sanitizeHtml(improved.split(/\r?\n/).map(l => `<p>${l || '<br>'}</p>`).join(''));
      editorRef.current?.setHtml?.(nextHtml);
      setHtml(nextHtml);
      onChange(nextHtml);
    } catch {
    } finally {
      setIsEnhancing(false);
    }
  };

  return (
    <div className={`mb-2 ${className}`}>
      {(label || aiContext) && (
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
      )}
      <div className="border border-slate-300 rounded-lg overflow-hidden focus-within:ring-2 focus-within:ring-indigo-500 focus-within:border-indigo-500 transition-all bg-white">
        <Toolbar editor={editor} defaultConfig={toolbarConfig} className="border-b border-slate-200 bg-slate-50" />
        <Editor
          defaultConfig={editorConfig}
          value={html}
          onCreated={(editor) => {
            editorRef.current = editor;
            setEditor(editor);
          }}
          onChange={(editor) => {
            const content = editor.getHtml();
            const safe = sanitizeHtml(content || '');
            setHtml(safe);
            onChange(safe);
          }}
          mode="default"
          style={{ minHeight: `${minRows * 48}px`, maxHeight, overflowY: 'auto', padding: '12px' }}
        />
      </div>
    </div>
  );
};

export default RichTextEditor;
