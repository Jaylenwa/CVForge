import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateCNTimeline: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const titleMap: Record<string, string> = {
     [ResumeSectionType.Summary]: '概况',
     [ResumeSectionType.Education]: '求学',
     [ResumeSectionType.Experience]: '履历',
     [ResumeSectionType.Skills]: '技能',
     [ResumeSectionType.Projects]: '项目',
     [ResumeSectionType.Custom]: '其他'
  };

  return (
    <div className={`w-full bg-slate-50 text-slate-800 h-full min-h-[1123px] ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto overflow-hidden`} style={{ fontFamily: '"Segoe UI", "Roboto", sans-serif' }}>
       <div className="bg-slate-800 text-white p-12 pb-24 relative overflow-hidden">
          <div className="absolute top-0 right-0 w-64 h-64 bg-slate-700 rounded-full mix-blend-overlay opacity-50 -mr-16 -mt-16"></div>
          <div className="relative z-10 flex justify-between items-end">
             <div>
                <h1 className="text-4xl font-light tracking-wide mb-2">{data.personalInfo.fullName}</h1>
                <p className="text-slate-400 uppercase tracking-widest text-sm">{data.personalInfo.jobTitle}</p>
             </div>
             <div className="text-right text-sm text-slate-400 space-y-1">
                {data.personalInfo.avatarUrl && (
                  <img
                    src={data.personalInfo.avatarUrl}
                    alt="Avatar"
                    className="w-20 h-20 object-cover rounded-md border border-slate-500 ml-auto mb-2"
                  />
                )}
                <div>{data.personalInfo.phone}</div>
                <div>{data.personalInfo.email}</div>
             </div>
          </div>
       </div>

       <div className="px-12 pb-12 -mt-12 relative z-20">
          <div className="bg-white rounded-lg shadow-lg p-8 min-h-[800px] relative">
             {/* 中心线 */}
             <div className="absolute left-[40px] top-8 bottom-8 w-0.5 bg-slate-200"></div>

             <div className="space-y-12">
                {(data.sections || []).filter(s => s.isVisible).map(section => (
                   <div key={section.id} className="relative pl-20">
                      {/* 图标节点 */}
                      <div className="absolute left-6 top-0 w-10 h-10 -ml-5 bg-slate-800 rounded-full border-4 border-white shadow-sm flex items-center justify-center text-white text-xs font-bold z-10">
                         {titleMap[section.type]?.[0] || section.title?.[0]}
                      </div>

                      <h3 className="text-xl font-bold text-slate-800 mb-6">{titleMap[section.type] || section.title}</h3>

                      <div className="space-y-8">
                         {section.type === ResumeSectionType.Skills ? (
                            <div className="flex flex-wrap gap-2">
                               {section.items.map(item => (
                                  <span key={item.id} className="px-3 py-1 bg-slate-100 text-slate-700 rounded text-sm">{item.description}</span>
                               ))}
                            </div>
                         ) : (
                            section.items.map(item => (
                               <div key={item.id} className="relative border-l-2 border-slate-100 pl-4 ml-1">
                                  <div className="flex justify-between items-start mb-2">
                                     <div>
                                        <h4 className="text-lg font-semibold text-slate-700">{item.title}</h4>
                                        {item.subtitle && <div className="text-sm text-slate-500">{item.subtitle}</div>}
                                     </div>
                                     {item.dateRange && <div className="text-xs font-mono text-slate-400 bg-slate-50 px-2 py-1 rounded">{item.dateRange}</div>}
                                  </div>
                                  {item.description && (
                                     <div className="text-sm text-slate-600 leading-relaxed" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
                                  )}
                               </div>
                            ))
                         )}
                      </div>
                   </div>
                ))}
             </div>
          </div>
       </div>
    </div>
  );
};
