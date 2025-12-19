import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useLanguage } from '../../../contexts/LanguageContext';
import { sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateCNMagazine: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const titleMap: Record<string, string> = {
    [ResumeSectionType.Summary]: 'INTRO',
    [ResumeSectionType.Education]: 'EDU',
    [ResumeSectionType.Experience]: 'WORK',
    [ResumeSectionType.Skills]: 'SKILL',
    [ResumeSectionType.Projects]: 'PROJ',
    [ResumeSectionType.Custom]: 'MISC'
  };

  return (
    <div className={`w-full bg-white text-black h-full min-h-[1123px] ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto overflow-hidden`} style={{ fontFamily: '"Bodoni MT", "Didot", serif' }}>
       <div className="grid grid-cols-12 h-full">
          {/* 左侧大图区 */}
          <div className="col-span-5 bg-gray-100 relative overflow-hidden flex flex-col justify-end p-8">
             {data.personalInfo.avatarUrl ? (
                <div className="absolute inset-0 z-0">
                   <img src={data.personalInfo.avatarUrl} alt="Cover" className="w-full h-full object-cover grayscale contrast-125" />
                   <div className="absolute inset-0 bg-gradient-to-t from-black/80 via-transparent to-transparent"></div>
                </div>
             ) : (
                <div className="absolute inset-0 bg-gray-300 z-0"></div>
             )}
             
             <div className="relative z-10 text-white">
                <h1 className="text-6xl font-bold leading-none mb-4 tracking-tighter mix-blend-difference">{data.personalInfo.fullName}</h1>
                <p className="text-xl font-light tracking-widest uppercase mb-8 border-l-4 border-white pl-4">{data.personalInfo.jobTitle}</p>
                
                <div className="space-y-2 text-sm font-sans font-light opacity-80">
                   <p>{data.personalInfo.phone}</p>
                   <p>{data.personalInfo.email}</p>
                   {data.personalInfo.city && <p>{data.personalInfo.city}</p>}
                </div>
             </div>
          </div>

          {/* 右侧内容区 */}
          <div className="col-span-7 p-10 flex flex-col h-full overflow-hidden">
             <div className="flex-grow space-y-10">
                {(data.sections || []).filter(s => s.isVisible).map((section, idx) => (
                   <div key={section.id} className="relative">
                      <h3 className="text-5xl font-bold text-gray-100 absolute -top-6 -left-4 -z-10 select-none">{titleMap[section.type] || section.title}</h3>
                      <h4 className="text-lg font-bold uppercase tracking-widest border-b-2 border-black mb-4 inline-block">{section.title}</h4>

                      <div className="space-y-6">
                         {section.type === ResumeSectionType.Skills ? (
                            <div className="flex flex-wrap gap-x-6 gap-y-2 font-sans">
                               {section.items.map(item => (
                                  <span key={item.id} className="text-sm font-bold uppercase">{item.description}</span>
                               ))}
                            </div>
                         ) : (
                            section.items.map(item => (
                               <div key={item.id} className="font-sans">
                                  <div className="flex justify-between items-baseline mb-1">
                                     <h5 className="font-bold text-base">{item.title}</h5>
                                     {item.dateRange && <span className="text-xs text-gray-500">{item.dateRange}</span>}
                                  </div>
                                  {item.subtitle && <div className="text-sm text-gray-600 mb-2 italic serif">{item.subtitle}</div>}
                                  {item.description && (
                                     <div className="text-sm text-gray-800 leading-snug text-justify" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
                                  )}
                               </div>
                            ))
                         )}
                      </div>
                   </div>
                ))}
             </div>
             
             {/* 底部条形码装饰 */}
             <div className="mt-auto pt-8 border-t border-gray-200 flex justify-between items-end">
                <div className="h-8 w-48 bg-black" style={{ maskImage: 'repeating-linear-gradient(90deg, black, black 1px, transparent 1px, transparent 3px)' }}></div>
                <div className="text-xs font-sans text-gray-400">ISSUE 01 • {new Date().getFullYear()}</div>
             </div>
          </div>
       </div>
    </div>
  );
};
