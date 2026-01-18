
import React, { useState, useRef, useEffect, useMemo } from 'react';
import { ChevronDown, X, Check, Search, CheckSquare, Square } from 'lucide-react';
import { Option } from '../types';

interface MultiSelectProps {
  options: Option[];
  value: string[];
  onChange: (value: string[]) => void;
  placeholder?: string;
  label?: string;
}

const MultiSelect: React.FC<MultiSelectProps> = ({
  options,
  value,
  onChange,
  placeholder = "Select options...",
  label
}) => {
  const [isOpen, setIsOpen] = useState(false);
  const [searchTerm, setSearchTerm] = useState("");
  const containerRef = useRef<HTMLDivElement>(null);

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const filteredOptions = useMemo(() => {
    return options.filter(opt => 
      opt.label.toLowerCase().includes(searchTerm.toLowerCase())
    );
  }, [options, searchTerm]);

  const toggleOption = (optionValue: string) => {
    const newValue = value.includes(optionValue)
      ? value.filter(v => v !== optionValue)
      : [...value, optionValue];
    onChange(newValue);
  };

  const isAllSelected = filteredOptions.length > 0 && 
    filteredOptions.every(opt => value.includes(opt.value));

  const handleSelectAll = () => {
    if (isAllSelected) {
      // Unselect only the filtered ones
      const filteredValues = filteredOptions.map(o => o.value);
      onChange(value.filter(v => !filteredValues.includes(v)));
    } else {
      // Select all filtered ones
      const newValues = Array.from(new Set([...value, ...filteredOptions.map(o => o.value)]));
      onChange(newValues);
    }
  };

  const removeValue = (e: React.MouseEvent, val: string) => {
    e.stopPropagation();
    onChange(value.filter(v => v !== val));
  };

  // Logic to handle display limit in the input area
  const displayLimit = 3;
  const selectedLabels = value.map(v => options.find(o => o.value === v)?.label).filter(Boolean);

  return (
    <div className="relative w-full" ref={containerRef}>
      {label && <label className="block text-sm font-semibold text-slate-700 mb-1.5">{label}</label>}
      
      <div 
        onClick={() => setIsOpen(!isOpen)}
        className={`
          min-h-[42px] w-full flex items-center justify-between px-3 py-1.5
          bg-white border rounded-lg cursor-pointer transition-all duration-200
          ${isOpen ? 'border-indigo-500 ring-2 ring-indigo-100' : 'border-slate-300 hover:border-slate-400'}
        `}
      >
        <div className="flex flex-wrap gap-1.5 overflow-hidden">
          {value.length === 0 ? (
            <span className="text-slate-400">{placeholder}</span>
          ) : (
            <>
              {selectedLabels.slice(0, displayLimit).map((label, idx) => (
                <span 
                  key={idx} 
                  className="inline-flex items-center px-2 py-0.5 rounded bg-indigo-50 text-indigo-700 text-xs font-medium border border-indigo-100"
                >
                  {label}
                  <button 
                    onClick={(e) => removeValue(e, value[idx])}
                    className="ml-1 hover:text-indigo-900 focus:outline-none"
                  >
                    <X size={12} />
                  </button>
                </span>
              ))}
              {selectedLabels.length > displayLimit && (
                <span className="text-xs font-medium text-slate-500 py-0.5">
                  +{selectedLabels.length - displayLimit} more
                </span>
              )}
            </>
          )}
        </div>
        <div className="flex items-center text-slate-400 ml-2">
          {value.length > 0 && (
            <X 
              size={16} 
              className="mr-1 hover:text-slate-600" 
              onClick={(e) => { e.stopPropagation(); onChange([]); }}
            />
          )}
          <ChevronDown size={18} className={`transition-transform duration-200 ${isOpen ? 'rotate-180' : ''}`} />
        </div>
      </div>

      {isOpen && (
        <div className="absolute z-50 w-full mt-2 bg-white border border-slate-200 rounded-xl shadow-xl overflow-hidden animate-in fade-in zoom-in duration-200 origin-top">
          {/* Search Bar */}
          <div className="flex items-center px-3 py-2 border-b border-slate-100 bg-slate-50">
            <Search size={16} className="text-slate-400 mr-2" />
            <input 
              type="text" 
              className="w-full bg-transparent border-none outline-none text-sm text-slate-700 placeholder-slate-400"
              placeholder="Search items..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              onClick={(e) => e.stopPropagation()}
            />
          </div>

          {/* Action Header */}
          <div className="px-3 py-2 flex items-center justify-between border-b border-slate-50">
            <button 
              onClick={(e) => { e.stopPropagation(); handleSelectAll(); }}
              className="flex items-center text-xs font-semibold text-indigo-600 hover:text-indigo-800 transition-colors"
            >
              {isAllSelected ? <CheckSquare size={14} className="mr-1" /> : <Square size={14} className="mr-1" />}
              {isAllSelected ? "Unselect All" : "Select All"}
            </button>
            <span className="text-[10px] text-slate-400 uppercase tracking-wider font-bold">
              {filteredOptions.length} available
            </span>
          </div>

          {/* Options List */}
          <div className="max-h-60 overflow-y-auto py-1 custom-scrollbar">
            {filteredOptions.length === 0 ? (
              <div className="px-4 py-8 text-center text-slate-500 text-sm">
                No results found for "{searchTerm}"
              </div>
            ) : (
              filteredOptions.map((option) => {
                const isSelected = value.includes(option.value);
                return (
                  <div
                    key={option.value}
                    onClick={(e) => { e.stopPropagation(); toggleOption(option.value); }}
                    className={`
                      px-4 py-2.5 flex items-center justify-between cursor-pointer transition-colors
                      ${isSelected ? 'bg-indigo-50 text-indigo-700' : 'hover:bg-slate-50 text-slate-700'}
                    `}
                  >
                    <span className="text-sm font-medium">{option.label}</span>
                    {isSelected && <Check size={16} className="text-indigo-600" />}
                  </div>
                );
              })
            )}
          </div>
        </div>
      )}
    </div>
  );
};

export default MultiSelect;
