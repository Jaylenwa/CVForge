import React, { useState } from 'react';

interface MarkdownEditorProps {
  value: string;
  onChange: (val: string) => void;
  onOptimize?: () => void;
  isOptimizing?: boolean;
}

const MarkdownEditor: React.FC<MarkdownEditorProps> = ({ value, onChange, onOptimize, isOptimizing = false }) => {
  const [isPreview, setIsPreview] = useState(false);

  const insertText = (before: string, after: string) => {
    const textarea = document.getElementById('md-textarea') as HTMLTextAreaElement;
    if (!textarea) return;
    const start = textarea.selectionStart;
    const end = textarea.selectionEnd;
    const text = textarea.value;
    const beforeText = text.substring(0, start);
    const afterText = text.substring(end, text.length);
    const selectedText = text.substring(start, end);
    const newValue = beforeText + before + selectedText + after + afterText;
    onChange(newValue);
    setTimeout(() => {
      textarea.focus();
      textarea.setSelectionRange(start + before.length, end + before.length);
    }, 0);
  };

  const toolbarActions = [
    { icon: 'B', label: 'Bold', action: () => insertText('**', '**') },
    { icon: 'I', label: 'Italic', action: () => insertText('*', '*') },
    { icon: 'U', label: 'Underline', action: () => insertText('<u>', '</u>') },
    { icon: '🔗', label: 'Link', action: () => insertText('[', '](url)') },
    { icon: '•', label: 'Bullets', action: () => insertText('\n• ', '') },
    { icon: '1.', label: 'Numbered', action: () => insertText('\n1. ', '') },
    { icon: '⌸', label: 'Table', action: () => insertText('\n| col | col |\n|---|---|\n| row | row |', '') },
  ];

  return (
    <div className="border rounded-lg bg-white overflow-hidden flex flex-col h-full shadow-sm">
      <div className="bg-slate-50 border-b p-1.5 flex items-center justify-between flex-wrap gap-2">
        <div className="flex items-center space-x-1">
          {toolbarActions.map((btn, i) => (
            <button
              key={i}
              onClick={btn.action}
              className="p-1 px-2 hover:bg-slate-200 rounded text-slate-700 font-semibold text-xs transition-colors border border-transparent hover:border-slate-300"
              title={btn.label}
            >
              {btn.icon}
            </button>
          ))}
          <div className="w-px h-6 bg-slate-300 mx-1"></div>
          <button 
            onClick={() => setIsPreview(!isPreview)}
            className={`text-[11px] px-2 py-1 rounded transition-colors ${isPreview ? 'bg-teal-100 text-teal-700' : 'bg-slate-200 text-slate-600'}`}
          >
            {isPreview ? '编辑模式' : '预览模式'}
          </button>
        </div>
        <span className="text-[10px] text-slate-400">功能提示：选中文字后再使用工具栏效果更佳</span>
      </div>
      
      <div className="relative flex-1">
        {isPreview ? (
          <div className="p-3 prose prose-sm max-w-none h-40 overflow-y-auto whitespace-pre-wrap text-slate-700">
            {value || <span className="text-slate-300 italic">尚未输入内容...</span>}
          </div>
        ) : (
          <textarea
            id="md-textarea"
            className="w-full h-40 p-3 text-sm focus:outline-none resize-none text-slate-700 placeholder:text-slate-300 leading-relaxed"
            placeholder="请输入工作详细内容，使用 Markdown 格式..."
            value={value}
            onChange={(e) => onChange(e.target.value)}
          />
        )}
        
        <div className="absolute right-3 bottom-3 flex flex-col space-y-2">
          <button 
             onClick={onOptimize}
             disabled={isOptimizing}
             className="bg-[#00c594] hover:bg-[#00b085] disabled:bg-slate-300 text-white px-3 py-1.5 rounded-full text-xs font-medium shadow-md transition-all flex items-center justify-center min-w-[88px]"
          >
            {isOptimizing ? (
              <>
                <svg className="animate-spin h-4 w-4 mr-2 text-white" viewBox="0 0 24 24"><circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle><path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>
                AI 优化中...
              </>
            ) : 'AI 撰写/优化'}
          </button>
        </div>
      </div>
    </div>
  );
};

export default MarkdownEditor;
