import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useLanguage } from '../../../contexts/LanguageContext';
import { sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateCNRetro: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const titleMap: Record<string, string> = {
     [ResumeSectionType.Summary]: '人物特写',
     [ResumeSectionType.Education]: '教育专栏',
     [ResumeSectionType.Experience]: '职场风云',
     [ResumeSectionType.Skills]: '专业技能',
     [ResumeSectionType.Projects]: '特别报道',
     [ResumeSectionType.Custom]: '副刊'
  };

  return (
    <div className={`w-full bg-[#f0e6d2] text-[#1a1a1a] h-full min-h-[1123px] ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto overflow-hidden`} style={{ fontFamily: '"Georgia", "Times New Roman", serif' }}>
       {/* 纸张纹理 */}
       <div className="absolute inset-0 opacity-20 pointer-events-none" style={{ filter: 'noise(0.5)' }}></div>
       
       <div className="p-10 h-full flex flex-col">
          {/* 报头 */}
          <div className="border-b-4 border-double border-black pb-4 mb-8 text-center relative">
             <div className="text-xs font-bold border-b border-black pb-1 mb-2 flex justify-between uppercase tracking-widest">
                <span>The Daily Resume</span>
                <span>Vol. {new Date().getFullYear()} No. 1</span>
                <span>Price: 1 Interview</span>
             </div>
             <h1 className="text-6xl font-black uppercase tracking-tight transform scale-y-110 mb-2">{data.personalInfo.fullName}</h1>
             <div className="text-xl font-bold italic font-serif border-t border-black pt-2 inline-block px-8">{data.personalInfo.jobTitle}</div>
             
             {/* 联系方式栏 */}
             <div className="absolute top-1/2 right-0 transform -translate-y-1/2 text-right hidden md:block text-xs font-bold w-32 leading-tight border-l border-black pl-2">
                <div>{data.personalInfo.phone}</div>
                <div className="break-words mt-1">{data.personalInfo.email}</div>
             </div>
          </div>

          {/* 多栏布局 */}
          <div className="flex-grow grid grid-cols-12 gap-6" style={{ columnRule: '1px solid #000' }}>
             {/* 左侧栏 - 头条 */}
             <div className="col-span-4 border-r border-black pr-6 flex flex-col">
                {data.personalInfo.avatarUrl && (
                   <div className="mb-6 border-2 border-black p-1 bg-white">
                      <img src={data.personalInfo.avatarUrl} alt="Profile" className="w-full object-cover" />
                      <div className="text-[10px] text-center mt-1 font-sans uppercase">Fig 1. Candidate Portrait</div>
                   </div>
                )}
                
                <div className="flex-grow">
                   {(data.sections || []).filter(s => s.type === ResumeSectionType.Summary || s.type === ResumeSectionType.Skills).map(section => (
                      <div key={section.id} className="mb-8">
                         <h3 className="text-lg font-bold uppercase border-b-2 border-black mb-3">{titleMap[section.type] || section.title}</h3>
                         {section.type === ResumeSectionType.Skills ? (
                            <ul className="list-disc pl-4 text-sm leading-relaxed">
                               {section.items.map(item => (
                                  <li key={item.id} className="pl-1 mb-1">{item.description}</li>
                               ))}
                            </ul>
                         ) : (
                            <div className="text-sm leading-relaxed text-justify first-letter:text-3xl first-letter:font-bold first-letter:float-left first-letter:mr-1 first-letter:leading-none">
                               {section.items.map(item => (
                                  <div key={item.id} dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description || '') }} />
                               ))}
                            </div>
                         )}
                      </div>
                   ))}
                </div>
             </div>

             {/* 右侧栏 - 详细报道 */}
             <div className="col-span-8">
                {(data.sections || []).filter(s => s.type !== ResumeSectionType.Summary && s.type !== ResumeSectionType.Skills && s.isVisible).map(section => (
                   <div key={section.id} className="mb-8 last:mb-0">
                      <div className="text-center mb-4 relative">
                         <span className="bg-[#f0e6d2] px-4 relative z-10 font-bold text-lg uppercase">{titleMap[section.type] || section.title}</span>
                         <div className="absolute top-1/2 left-0 w-full h-px bg-black -z-0"></div>
                      </div>

                      <div className="space-y-6">
                         {section.items.map(item => (
                            <div key={item.id}>
                               <div className="flex justify-between items-baseline border-b border-dotted border-black mb-1">
                                  <h4 className="text-xl font-bold">{item.title}</h4>
                                  {item.dateRange && <span className="text-sm font-bold">{item.dateRange}</span>}
                               </div>
                               {item.subtitle && <div className="text-sm font-bold italic mb-2">{item.subtitle}</div>}
                               {item.description && (
                                  <div className="text-sm leading-relaxed text-justify columns-1 md:columns-2 gap-4" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
                               )}
                            </div>
                         ))}
                      </div>
                   </div>
                ))}
             </div>
          </div>
          
          {/* 页脚 */}
          <div className="border-t-2 border-black mt-4 pt-2 text-center text-[10px] uppercase font-bold tracking-widest">
             Printed in {new Date().getFullYear()} • OpenResume Press • All Rights Reserved
          </div>
       </div>
    </div>
  );
};
