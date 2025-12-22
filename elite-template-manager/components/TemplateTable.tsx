
import React from 'react';
import { Template } from '../types';

interface TableProps {
  templates: Template[];
  onDelete: (id: string) => void;
  onUpdate: (id: string) => void;
  onPreview: (id: string) => void;
}

export const TemplateTable: React.FC<TableProps> = ({ templates, onDelete, onUpdate, onPreview }) => {
  if (templates.length === 0) {
    return (
      <div className="py-20 text-center">
        <div className="w-16 h-16 bg-slate-100 rounded-full flex items-center justify-center mx-auto mb-4 text-slate-400">
          <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><rect width="18" height="18" x="3" y="3" rx="2"/><path d="M9 17h6"/><path d="M10 11h.01"/><path d="M14 11h.01"/></svg>
        </div>
        <h3 className="text-lg font-medium text-slate-900">No templates found</h3>
        <p className="text-slate-500">Try adjusting your filters or search keywords.</p>
      </div>
    );
  }

  return (
    <div className="overflow-x-auto">
      <table className="w-full text-left border-collapse">
        <thead>
          <tr className="bg-slate-50/80 border-b border-slate-200">
            <th className="px-6 py-4 text-xs font-bold text-slate-500 uppercase tracking-wider">ID</th>
            <th className="px-6 py-4 text-xs font-bold text-slate-500 uppercase tracking-wider">Template Name</th>
            <th className="px-6 py-4 text-xs font-bold text-slate-500 uppercase tracking-wider">Industry</th>
            <th className="px-6 py-4 text-xs font-bold text-slate-500 uppercase tracking-wider">Popularity</th>
            <th className="px-6 py-4 text-xs font-bold text-slate-500 uppercase tracking-wider">Status</th>
            <th className="px-6 py-4 text-xs font-bold text-slate-500 uppercase tracking-wider">Tags</th>
            <th className="px-6 py-4 text-xs font-bold text-slate-500 uppercase tracking-wider text-right">Actions</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-slate-100">
          {templates.map((t) => (
            <tr key={t.id} className="hover:bg-indigo-50/30 transition-colors group">
              <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-slate-400">
                #{t.id}
              </td>
              <td className="px-6 py-4 whitespace-nowrap">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-12 bg-slate-100 rounded-md overflow-hidden flex-shrink-0 border border-slate-200">
                    <img src={`https://picsum.photos/seed/${t.id}/40/60`} alt="" className="w-full h-full object-cover grayscale group-hover:grayscale-0 transition-all" />
                  </div>
                  <div>
                    <div className="text-sm font-semibold text-slate-900">{t.name}</div>
                    <div className="text-[10px] text-slate-400 uppercase tracking-tight">Last updated: {t.lastUpdated}</div>
                  </div>
                </div>
              </td>
              <td className="px-6 py-4 whitespace-nowrap">
                <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border ${
                  t.industry === 'IT' ? 'bg-blue-50 text-blue-700 border-blue-100' :
                  t.industry === 'Finance' ? 'bg-emerald-50 text-emerald-700 border-emerald-100' :
                  t.industry === 'Creative' ? 'bg-purple-50 text-purple-700 border-purple-100' :
                  'bg-slate-50 text-slate-700 border-slate-100'
                }`}>
                  {t.industry}
                </span>
              </td>
              <td className="px-6 py-4 whitespace-nowrap">
                <div className="flex items-center gap-2">
                  <div className="flex-1 w-16 bg-slate-100 rounded-full h-1.5 overflow-hidden">
                    <div 
                      className="bg-indigo-500 h-full rounded-full transition-all duration-1000"
                      style={{ width: `${t.popularity}%` }}
                    />
                  </div>
                  <span className="text-sm font-semibold text-slate-700">{t.popularity}%</span>
                </div>
              </td>
              <td className="px-6 py-4 whitespace-nowrap">
                {t.isPremium ? (
                  <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-bold bg-amber-50 text-amber-700 border border-amber-200 shadow-sm">
                    <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="currentColor" className="text-amber-500"><path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/></svg>
                    PREMIUM
                  </span>
                ) : (
                  <span className="text-xs font-medium text-slate-400">Standard</span>
                )}
              </td>
              <td className="px-6 py-4">
                <div className="flex flex-wrap gap-1.5 max-w-[200px]">
                  {t.tags.map((tag, idx) => (
                    <span key={idx} className="px-2 py-0.5 bg-white border border-slate-200 rounded text-[11px] text-slate-600 font-medium hover:border-indigo-300 hover:bg-indigo-50 transition-colors">
                      {tag}
                    </span>
                  ))}
                </div>
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-right">
                <div className="flex items-center justify-end gap-1">
                  <button 
                    onClick={() => onPreview(t.id)}
                    className="p-2 text-slate-400 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg transition-all"
                    title="Preview"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7Z"/><circle cx="12" cy="12" r="3"/></svg>
                  </button>
                  <button 
                    onClick={() => onUpdate(t.id)}
                    className="p-2 text-slate-400 hover:text-emerald-600 hover:bg-emerald-50 rounded-lg transition-all"
                    title="Edit"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M17 3a2.85 2.83 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5Z"/><path d="m15 5 4 4"/></svg>
                  </button>
                  <button 
                    onClick={() => onDelete(t.id)}
                    className="p-2 text-slate-400 hover:text-rose-600 hover:bg-rose-50 rounded-lg transition-all"
                    title="Delete"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/><line x1="10" x2="10" y1="11" y2="17"/><line x1="14" x2="14" y1="11" y2="17"/></svg>
                  </button>
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};
