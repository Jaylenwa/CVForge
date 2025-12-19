import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useLanguage } from '../../../contexts/LanguageContext';
import { sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateCNArchitect: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const titleMap: Record<string, string> = {
    [ResumeSectionType.Summary]: 'DESIGN CONCEPT',
    [ResumeSectionType.Education]: 'FOUNDATION',
    [ResumeSectionType.Experience]: 'CONSTRUCTION',
    [ResumeSectionType.Skills]: 'TOOLS',
    [ResumeSectionType.Projects]: 'PORTFOLIO',
    [ResumeSectionType.Custom]: 'SPECS'
  };

  return (
    <div className={`w-full bg-[#003366] text-white h-full min-h-[1123px] ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto font-mono`} style={{ fontFamily: '"Courier New", monospace' }}>
      {/* Blueprint Grid */}
      <div className="absolute inset-0 opacity-20 pointer-events-none" 
           style={{ backgroundImage: 'linear-gradient(rgba(255,255,255,0.3) 1px, transparent 1px), linear-gradient(90deg, rgba(255,255,255,0.3) 1px, transparent 1px)', backgroundSize: '40px 40px' }}>
      </div>
      
      <div className="p-8 border-[12px] border-white/10 h-full m-0 relative">
          {/* Corner Decorations */}
          <div className="absolute top-0 left-0 w-16 h-16 border-r border-b border-white/30"></div>
          <div className="absolute bottom-0 right-0 w-64 h-32 border-l border-t border-white/30 bg-[#002244] p-4 flex flex-col justify-end">
              <div className="border border-white/50 p-2 text-[10px] space-y-1">
                  <div className="flex justify-between border-b border-white/20"><span>SCALE</span><span>1:1</span></div>
                  <div className="flex justify-between border-b border-white/20"><span>DATE</span><span>{new Date().getFullYear()}</span></div>
                  <div className="flex justify-between"><span>DWG NO.</span><span>A-101</span></div>
              </div>
          </div>

          {/* Header Block */}
          <div className="border-b-2 border-white/50 pb-8 mb-12 flex justify-between items-end relative z-10">
              <div>
                  <h1 className="text-5xl font-bold mb-2 tracking-[0.2em]">{data.personalInfo.fullName}</h1>
                  <div className="text-xl tracking-widest opacity-80">{data.personalInfo.jobTitle}</div>
              </div>
              <div className="text-right text-sm space-y-2 opacity-70">
                   <div className="flex items-center justify-end gap-2">
                       <span className="w-2 h-2 border border-white rounded-full"></span>
                       {data.personalInfo.phone}
                   </div>
                   <div className="flex items-center justify-end gap-2">
                       <span className="w-2 h-2 border border-white rounded-full"></span>
                       {data.personalInfo.email}
                   </div>
              </div>
          </div>

          <div className="grid grid-cols-12 gap-12 relative z-10">
              <div className="col-span-12 space-y-12">
                  {(data.sections || []).filter(s => s.isVisible).map(section => (
                      <div key={section.id}>
                          <div className="flex items-center mb-6">
                              <div className="w-8 h-8 border border-white rounded-full flex items-center justify-center mr-4 text-xs font-bold">
                                  {section.title?.charAt(0) || ''}
                              </div>
                              <h3 className="text-2xl font-bold tracking-widest border-b border-white/30 pb-1 pr-12">
                                  {titleMap[section.type] || section.title}
                              </h3>
                          </div>

                          {section.type === ResumeSectionType.Skills ? (
                              <div className="grid grid-cols-4 gap-4">
                                  {section.items.map(item => (
                                      <div key={item.id} className="border border-white/30 p-2 text-center text-sm hover:bg-white/10 transition-colors">
                                          {item.description}
                                      </div>
                                  ))}
                              </div>
                          ) : (
                              <div className="space-y-8 pl-12">
                                  {section.items.map(item => (
                                      <div key={item.id} className="relative">
                                          {/* Measurement Line */}
                                          <div className="absolute -left-8 top-2 bottom-2 w-4 border-l border-white/20 flex flex-col justify-between items-center">
                                              <div className="w-2 h-px bg-white/20"></div>
                                              <div className="w-2 h-px bg-white/20"></div>
                                          </div>

                                          <div className="flex justify-between items-baseline mb-2">
                                              <h4 className="text-xl font-bold text-white/90">{item.title}</h4>
                                              {item.dateRange && <span className="text-xs border border-white/30 px-2 py-1">{item.dateRange}</span>}
                                          </div>
                                          {item.subtitle && <div className="text-lg opacity-70 mb-2 italic">{item.subtitle}</div>}
                                          {item.description && (
                                              <div className="text-sm opacity-80 leading-relaxed max-w-[90%]" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
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
    </div>
  );
};
