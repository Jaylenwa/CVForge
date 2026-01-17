
import React from 'react';
import { Template } from '../types';
import { Download, Eye } from 'lucide-react';

interface ResumeCardProps {
  template: Template;
}

const ResumeCard: React.FC<ResumeCardProps> = ({ template }) => {
  return (
    <div className="group relative bg-white rounded-xl border border-slate-200 overflow-hidden transition-all duration-300 hover:shadow-xl hover:-translate-y-1">
      {/* Image Preview */}
      <div className="relative aspect-[3/4] overflow-hidden bg-slate-100">
        <img 
          src={template.thumbnail} 
          alt={template.title}
          className="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
        />
        {/* Overlay Buttons */}
        <div className="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity duration-300 flex items-center justify-center gap-4">
          <button className="bg-blue-600 text-white p-3 rounded-full hover:bg-blue-700 transition-colors shadow-lg">
            <Eye size={20} />
          </button>
          <button className="bg-white text-slate-900 p-3 rounded-full hover:bg-slate-100 transition-colors shadow-lg">
            <Download size={20} />
          </button>
        </div>
      </div>

      {/* Info Section */}
      <div className="p-5">
        <h3 className="text-lg font-bold text-slate-800 mb-1 line-clamp-1">{template.title}</h3>
        <div className="flex items-center justify-between mt-3">
          <span className="text-sm text-slate-500">{template.usageCount.toLocaleString()} 使用次数</span>
          <span className="px-2.5 py-1 text-xs font-medium bg-blue-50 text-blue-600 rounded-md border border-blue-100">
            {template.tag}
          </span>
        </div>
      </div>
    </div>
  );
};

export default ResumeCard;
