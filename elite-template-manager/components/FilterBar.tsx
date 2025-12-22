
import React from 'react';
import { Industry, FilterState } from '../types';

interface FilterBarProps {
  filters: FilterState;
  setFilters: (filters: FilterState) => void;
  onAdd: () => void;
}

const INDUSTRIES: Industry[] = ['All', 'General', 'IT', 'Finance', 'Creative', 'Medical', 'Education'];

export const FilterBar: React.FC<FilterBarProps> = ({ filters, setFilters, onAdd }) => {
  return (
    <div className="p-4 sm:p-6 border-b border-slate-200 flex flex-col md:flex-row gap-4 items-stretch md:items-center justify-between bg-white">
      <div className="flex flex-col sm:flex-row gap-4 flex-1">
        <div className="relative group flex-1">
          <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-slate-400 group-focus-within:text-indigo-500 transition-colors">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><circle cx="11" cy="11" r="8"/><path d="m21 21-4.3-4.3"/></svg>
          </div>
          <input 
            type="text" 
            placeholder="Search templates, tags..." 
            value={filters.keyword}
            onChange={(e) => setFilters({ ...filters, keyword: e.target.value })}
            className="block w-full pl-10 pr-3 py-2.5 bg-slate-50 border border-slate-200 rounded-xl text-sm placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all"
          />
        </div>
        
        <div className="relative">
          <select 
            value={filters.industry}
            onChange={(e) => setFilters({ ...filters, industry: e.target.value as Industry })}
            className="appearance-none block w-full sm:w-48 pl-4 pr-10 py-2.5 bg-slate-50 border border-slate-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all cursor-pointer"
          >
            {INDUSTRIES.map(ind => (
              <option key={ind} value={ind}>{ind} Industry</option>
            ))}
          </select>
          <div className="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none text-slate-400">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="m6 9 6 6 6-6"/></svg>
          </div>
        </div>
      </div>

      <button 
        onClick={onAdd}
        className="flex items-center justify-center gap-2 px-6 py-2.5 bg-indigo-600 hover:bg-indigo-700 text-white rounded-xl text-sm font-semibold shadow-md shadow-indigo-200 hover:shadow-lg hover:shadow-indigo-300 transition-all active:scale-95"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M5 12h14"/><path d="M12 5v14"/></svg>
        New Template
      </button>
    </div>
  );
};
