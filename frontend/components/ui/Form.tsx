import React from 'react';

export const Label: React.FC<{ children: React.ReactNode; htmlFor?: string; required?: boolean; right?: React.ReactNode }> = ({ children, htmlFor, required, right }) => (
  <div className="flex items-center justify-between mb-1.5">
    <label htmlFor={htmlFor} className="block text-sm font-medium text-slate-700 flex items-center">
      {children}
      {required && <span className="text-rose-500 ml-1">*</span>}
    </label>
    {right ? <div className="text-xs text-slate-500">{right}</div> : null}
  </div>
);

export const Input: React.FC<React.InputHTMLAttributes<HTMLInputElement>> = (props) => (
  <input
    {...props}
    className={`w-full px-3 py-2 bg-white border border-slate-200 rounded-lg text-slate-900 text-sm ring-offset-white focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all duration-200 disabled:bg-slate-50 disabled:text-slate-400 ${props.className || ''}`}
  />
);

export const Textarea: React.FC<React.TextareaHTMLAttributes<HTMLTextAreaElement>> = (props) => (
  <textarea
    {...props}
    className={`w-full px-3 py-2 bg-white border border-slate-200 rounded-lg text-slate-900 text-sm ring-offset-white focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all duration-200 disabled:bg-slate-50 disabled:text-slate-400 ${props.className || ''}`}
  />
);

export const Select: React.FC<
  React.SelectHTMLAttributes<HTMLSelectElement> & { options: { label: string; value: string; disabled?: boolean; hidden?: boolean }[] }
> = ({ options, ...props }) => (
  <div className="relative">
    <select
      {...props}
      className={`w-full px-3 py-2 bg-white border border-slate-200 rounded-lg text-slate-900 text-sm appearance-none focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all duration-200 ${props.className || ''}`}
    >
      {options.map((opt) => (
        <option key={opt.value} value={opt.value} disabled={!!opt.disabled} hidden={!!opt.hidden}>
          {opt.label}
        </option>
      ))}
    </select>
    <div className="absolute inset-y-0 right-0 flex items-center pr-2 pointer-events-none text-slate-400">
      <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 9l-7 7-7-7" /></svg>
    </div>
  </div>
);

export const Checkbox: React.FC<{ label: React.ReactNode; checked: boolean; onChange: (checked: boolean) => void }> = ({ label, checked, onChange }) => (
  <label className="group flex items-center space-x-3 cursor-pointer select-none">
    <div className={`relative w-5 h-5 flex items-center justify-center rounded border transition-all duration-200 ${checked ? 'bg-blue-600 border-blue-600' : 'bg-white border-slate-300 group-hover:border-blue-400'}`}>
      {checked && (
        <svg className="w-3.5 h-3.5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="3" d="M5 13l4 4L19 7" /></svg>
      )}
      <input type="checkbox" className="absolute opacity-0 w-full h-full cursor-pointer" checked={checked} onChange={(e) => onChange(e.target.checked)} />
    </div>
    <span className="text-sm font-medium text-slate-700">{label}</span>
  </label>
);

export const TagInput: React.FC<{ tags: string[]; onChange: (tags: string[]) => void; placeholder?: string }> = ({ tags, onChange, placeholder }) => {
  const [inputValue, setInputValue] = React.useState('');

  const addTag = () => {
    const trimmed = inputValue.trim();
    if (trimmed && !tags.includes(trimmed)) {
      onChange([...tags, trimmed]);
      setInputValue('');
    } else if (!trimmed) {
      setInputValue('');
    }
  };

  const removeTag = (tagToRemove: string) => {
    onChange(tags.filter(t => t !== tagToRemove));
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' || e.key === ',') {
      e.preventDefault();
      addTag();
    }
    if (e.key === 'Backspace' && !inputValue && tags.length > 0) {
      onChange(tags.slice(0, -1));
    }
  };

  return (
    <div className="w-full min-h-[42px] px-2 py-1.5 bg-white border border-slate-200 rounded-lg flex flex-wrap gap-2 focus-within:ring-2 focus-within:ring-blue-500/20 focus-within:border-blue-500 transition-all duration-200">
      {tags.map((tag) => (
        <span key={tag} className="inline-flex items-center px-2 py-0.5 rounded-md bg-blue-50 text-blue-700 text-xs font-medium border border-blue-100">
          <span className="max-w-[160px] truncate" title={tag}>{tag}</span>
          <button type="button" onClick={() => removeTag(tag)} className="ml-1.5 text-blue-400 hover:text-blue-600 focus:outline-none">
            <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12" /></svg>
          </button>
        </span>
      ))}
      <input
        type="text"
        value={inputValue}
        onChange={(e) => setInputValue(e.target.value)}
        onKeyDown={handleKeyDown}
        onBlur={addTag}
        placeholder={tags.length === 0 ? (placeholder || '输入标签后按回车') : ''}
        className="flex-grow min-w-[60px] bg-transparent text-sm text-slate-900 focus:outline-none"
      />
    </div>
  );
};
