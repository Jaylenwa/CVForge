import React, { useEffect, useMemo, useRef, useState } from 'react';
import ReactDOM from 'react-dom';
import { Check, ChevronDown, Search, Square, CheckSquare, X } from 'lucide-react';

export type MultiSelectOption = { value: string; label: string };

export const MultiSelect: React.FC<{
  options: MultiSelectOption[];
  value: string[];
  onChange: (value: string[]) => void;
  placeholder?: string;
  searchPlaceholder?: string;
  selectAllLabel?: string;
  unselectAllLabel?: string;
  availableSuffix?: string;
  noResultsLabel?: (q: string) => string;
  disabled?: boolean;
}> = ({
  options,
  value,
  onChange,
  placeholder = 'Select...',
  searchPlaceholder = 'Search...',
  selectAllLabel = 'Select all',
  unselectAllLabel = 'Unselect all',
  availableSuffix = 'available',
  noResultsLabel = (q) => `No results for "${q}"`,
  disabled,
}) => {
  const [isOpen, setIsOpen] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const containerRef = useRef<HTMLDivElement>(null);
  const triggerRef = useRef<HTMLButtonElement>(null);
  const dropdownRef = useRef<HTMLDivElement>(null);
  const [dropdownStyle, setDropdownStyle] = useState<React.CSSProperties>({});

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      const target = event.target as Node;
      if (containerRef.current?.contains(target)) return;
      if (dropdownRef.current?.contains(target)) return;
      setIsOpen(false);
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  useEffect(() => {
    if (!isOpen) return;
    const updatePosition = () => {
      const rect = triggerRef.current?.getBoundingClientRect();
      if (!rect) return;

      const gap = 8;
      const viewportH = window.innerHeight;
      const desiredTop = rect.bottom + gap;
      const maxPanelH = 360;
      const openUp = desiredTop + maxPanelH > viewportH && rect.top > viewportH / 2;
      const top = openUp ? Math.max(8, rect.top - gap) : desiredTop;

      setDropdownStyle({
        position: 'fixed',
        left: Math.max(8, rect.left),
        top,
        width: rect.width,
        zIndex: 1000,
        transform: openUp ? 'translateY(-100%)' : undefined,
      });
    };

    updatePosition();
    const onScroll = () => updatePosition();
    const onResize = () => updatePosition();
    document.addEventListener('scroll', onScroll, true);
    window.addEventListener('resize', onResize);
    return () => {
      document.removeEventListener('scroll', onScroll, true);
      window.removeEventListener('resize', onResize);
    };
  }, [isOpen]);

  const selectedLabels = useMemo(() => {
    const labelByValue = new Map(options.map((o) => [o.value, o.label]));
    return value.map((v) => labelByValue.get(v)).filter(Boolean) as string[];
  }, [options, value]);

  const filteredOptions = useMemo(() => {
    const q = searchTerm.trim().toLowerCase();
    if (!q) return options;
    return options.filter((opt) => opt.label.toLowerCase().includes(q));
  }, [options, searchTerm]);

  const isAllSelected = filteredOptions.length > 0 && filteredOptions.every((opt) => value.includes(opt.value));

  const toggleOption = (optionValue: string) => {
    if (value.includes(optionValue)) {
      onChange(value.filter((v) => v !== optionValue));
      return;
    }
    onChange([...value, optionValue]);
  };

  const handleSelectAll = () => {
    const filteredValues = filteredOptions.map((o) => o.value);
    if (filteredValues.length === 0) return;
    if (isAllSelected) {
      onChange(value.filter((v) => !filteredValues.includes(v)));
      return;
    }
    onChange(Array.from(new Set([...value, ...filteredValues])));
  };

  const removeValue = (e: React.MouseEvent, val: string) => {
    e.stopPropagation();
    onChange(value.filter((v) => v !== val));
  };

  const clearAll = (e: React.MouseEvent) => {
    e.stopPropagation();
    onChange([]);
  };

  const displayLimit = 3;

  return (
    <div ref={containerRef} className="relative w-full">
      <button
        ref={triggerRef}
        type="button"
        disabled={!!disabled}
        onClick={() => setIsOpen((p) => !p)}
        className={`min-h-[42px] w-full flex items-center justify-between px-3 py-1.5 bg-white border rounded-lg transition-all duration-200 ${
          disabled
            ? 'bg-slate-50 text-slate-400 border-slate-200 cursor-not-allowed'
            : isOpen
              ? 'border-blue-500 ring-2 ring-blue-500/20 cursor-pointer'
              : 'border-slate-200 hover:border-slate-300 cursor-pointer'
        }`}
      >
        <div className="flex flex-wrap gap-1.5 overflow-hidden">
          {value.length === 0 ? (
            <span className="text-slate-400 text-sm">{placeholder}</span>
          ) : (
            <>
              {selectedLabels.slice(0, displayLimit).map((label, idx) => (
                <span
                  key={label}
                  className="inline-flex items-center px-2 py-0.5 rounded bg-blue-50 text-blue-700 text-xs font-medium border border-blue-100"
                >
                  <span className="max-w-[220px] truncate" title={label}>
                    {label}
                  </span>
                  <span
                    role="button"
                    tabIndex={0}
                    onClick={(e) => removeValue(e as any, value[idx])}
                    className="ml-1 hover:text-blue-900"
                  >
                    <X size={12} />
                  </span>
                </span>
              ))}
              {selectedLabels.length > displayLimit ? (
                <span className="text-xs font-medium text-slate-500 py-0.5">+{selectedLabels.length - displayLimit}</span>
              ) : null}
            </>
          )}
        </div>
        <div className={`flex items-center ml-2 ${disabled ? 'text-slate-300' : 'text-slate-400'}`}>
          {value.length > 0 && !disabled ? (
            <span role="button" tabIndex={0} className="mr-1 hover:text-slate-600" onClick={clearAll}>
              <X size={16} />
            </span>
          ) : null}
          <ChevronDown size={18} className={`transition-transform duration-200 ${isOpen ? 'rotate-180' : ''}`} />
        </div>
      </button>

      {isOpen && !disabled
        ? ReactDOM.createPortal(
            <div ref={dropdownRef} style={dropdownStyle} className="bg-white border border-slate-200 rounded-xl shadow-xl overflow-hidden">
              <div className="flex items-center px-3 py-2 border-b border-slate-100 bg-slate-50">
                <Search size={16} className="text-slate-400 mr-2" />
                <input
                  type="text"
                  className="w-full bg-transparent border-none outline-none text-sm text-slate-700 placeholder-slate-400"
                  placeholder={searchPlaceholder}
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  onClick={(e) => e.stopPropagation()}
                />
              </div>

              <div className="px-3 py-2 flex items-center justify-between border-b border-slate-50">
                <button
                  type="button"
                  onClick={(e) => {
                    e.stopPropagation();
                    handleSelectAll();
                  }}
                  className="flex items-center text-xs font-semibold text-blue-600 hover:text-blue-800 transition-colors"
                >
                  {isAllSelected ? <CheckSquare size={14} className="mr-1" /> : <Square size={14} className="mr-1" />}
                  {isAllSelected ? unselectAllLabel : selectAllLabel}
                </button>
                <span className="text-[10px] text-slate-400 uppercase tracking-wider font-bold">
                  {filteredOptions.length} {availableSuffix}
                </span>
              </div>

              <div className="max-h-60 overflow-y-auto py-1">
                {filteredOptions.length === 0 ? (
                  <div className="px-4 py-8 text-center text-slate-500 text-sm">{noResultsLabel(searchTerm)}</div>
                ) : (
                  filteredOptions.map((option) => {
                    const isSelected = value.includes(option.value);
                    return (
                      <div
                        key={option.value}
                        onClick={(e) => {
                          e.stopPropagation();
                          toggleOption(option.value);
                        }}
                        className={`px-4 py-2.5 flex items-center justify-between cursor-pointer transition-colors ${
                          isSelected ? 'bg-blue-50 text-blue-700' : 'hover:bg-slate-50 text-slate-700'
                        }`}
                      >
                        <span className="text-sm font-medium">{option.label}</span>
                        {isSelected ? <Check size={16} className="text-blue-600" /> : null}
                      </div>
                    );
                  })
                )}
              </div>
            </div>,
            document.body
          )
        : null}
    </div>
  );
};
