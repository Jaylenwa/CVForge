
import React from 'react';
import { JobSubCategory } from '../types';

interface JobMegaMenuProps {
  subCategories: JobSubCategory[];
  isVisible: boolean;
}

const JobMegaMenu: React.FC<JobMegaMenuProps> = ({ subCategories, isVisible }) => {
  if (!isVisible || subCategories.length === 0) return null;

  return (
    <div className="absolute left-full top-0 ml-4 w-[600px] bg-white rounded-2xl shadow-2xl border border-slate-100 z-50 p-8 animate-in fade-in slide-in-from-left-2 duration-200">
      <div className="grid grid-cols-1 gap-10">
        {subCategories.map((sub, index) => (
          <div key={index} className="space-y-4">
            <h4 className="text-base font-bold text-slate-900 flex items-center">
              {sub.title}
              <div className="h-px flex-1 bg-slate-100 ml-4"></div>
            </h4>
            <div className="flex flex-wrap gap-x-6 gap-y-3">
              {sub.roles.map((role, idx) => (
                <button 
                  key={idx} 
                  className="text-sm text-slate-600 hover:text-blue-600 transition-colors"
                >
                  {role}
                </button>
              ))}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default JobMegaMenu;
