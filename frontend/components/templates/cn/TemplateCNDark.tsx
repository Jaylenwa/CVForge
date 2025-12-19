import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateCNDark: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const titleMap: Record<string, string> = {
    [ResumeSectionType.Summary]: 'SUMMARY',
    [ResumeSectionType.Education]: 'EDUCATION',
    [ResumeSectionType.Experience]: 'EXPERIENCE',
    [ResumeSectionType.Skills]: 'SKILLS',
    [ResumeSectionType.Projects]: 'PROJECTS',
    [ResumeSectionType.Custom]: 'OTHERS'
  };

  return (
    <div className={`w-full bg-[#1a1a1a] text-gray-300 h-full min-h-[1123px] ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto overflow-hidden`} style={{ fontFamily: '"Didot", "Playfair Display", serif' }}>
       {/* 金色纹理 */}
       <div className="absolute top-0 left-0 right-0 h-1 bg-gradient-to-r from-[#1a1a1a] via-[#d4af37] to-[#1a1a1a]"></div>
       <div className="absolute bottom-0 left-0 right-0 h-1 bg-gradient-to-r from-[#1a1a1a] via-[#d4af37] to-[#1a1a1a]"></div>

       <div className="p-16">
          <div className="text-center mb-16 relative">
             <div className="inline-block border border-[#d4af37] p-1 rounded-full mb-6">
                {data.personalInfo.avatarUrl && <img src={data.personalInfo.avatarUrl} alt="Avatar" className="w-28 h-28 rounded-full object-cover grayscale brightness-75 contrast-125" />}
             </div>
             <h1 className="text-4xl text-white font-normal tracking-[0.3em] uppercase mb-4">{data.personalInfo.fullName}</h1>
             <div className="text-[#d4af37] text-sm tracking-[0.2em] uppercase mb-8">{data.personalInfo.jobTitle}</div>
             
             <div className="flex justify-center gap-8 text-xs text-gray-500 tracking-widest uppercase border-t border-gray-800 pt-8 max-w-lg mx-auto">
                <span>{data.personalInfo.phone}</span>
                <span>•</span>
                <span>{data.personalInfo.email}</span>
             </div>
          </div>

          <div className="space-y-12 max-w-3xl mx-auto">
             {(data.sections || []).filter(s => s.isVisible).map(section => (
                <div key={section.id}>
                   <h3 className="text-center text-[#d4af37] text-sm tracking-[0.3em] uppercase mb-8 flex items-center justify-center gap-4">
                      <span className="w-12 h-px bg-gray-800"></span>
                      {titleMap[section.type] || section.title}
                      <span className="w-12 h-px bg-gray-800"></span>
                   </h3>

                   {section.type === ResumeSectionType.Skills ? (
                      <div className="text-center space-x-6">
                         {section.items.map(item => (
                            <span key={item.id} className="text-sm text-gray-400 inline-block py-1 border-b border-gray-800 hover:border-[#d4af37] transition-colors">{item.description}</span>
                         ))}
                      </div>
                   ) : (
                      <div className="space-y-10">
                         {section.items.map(item => (
                            <div key={item.id} className="relative group">
                               <div className="text-center mb-4">
                                  <h4 className="text-xl text-white font-normal mb-1">{item.title}</h4>
                                  {item.subtitle && <div className="text-sm text-[#d4af37] italic mb-1">{item.subtitle}</div>}
                                  {item.dateRange && <div className="text-xs text-gray-600 font-sans">{item.dateRange}</div>}
                               </div>
                               {item.description && (
                                  <div className="text-sm text-gray-400 leading-loose text-center font-sans font-light max-w-xl mx-auto" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
                               )}
                            </div>
                         ))}
                      </div>
                   )}
                </div>
             ))}
          </div>
       </div>
    </div>
  );
};
