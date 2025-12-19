import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useLanguage } from '../../../contexts/LanguageContext';
import { sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateCNArtistic: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const titleMap: Record<string, string> = {
    [ResumeSectionType.Summary]: '自述',
    [ResumeSectionType.Education]: '师承',
    [ResumeSectionType.Experience]: '历练',
    [ResumeSectionType.Skills]: '技艺',
    [ResumeSectionType.Projects]: '造物',
    [ResumeSectionType.Custom]: '补遗'
  };

  return (
    <div className={`w-full bg-[#f4f1ea] text-[#2c2c2c] h-full min-h-[1123px] ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto overflow-hidden`} style={{ fontFamily: '"KaiTi", "STKaiti", "Baskerville", serif' }}>
       {/* 水墨背景 */}
       <div className="absolute top-0 left-0 w-full h-64 bg-gradient-to-b from-[#e6e2d8] to-transparent opacity-50 pointer-events-none"></div>
       <div className="absolute bottom-0 right-0 w-96 h-96 bg-[radial-gradient(circle_at_center,_var(--tw-gradient-stops))] from-gray-200 to-transparent opacity-40 pointer-events-none rounded-full blur-3xl"></div>

       <div className="p-16 relative z-10">
          <div className="flex flex-col items-center mb-16 relative">
             {data.personalInfo.avatarUrl && (
                <div className="mb-6 p-1 border border-gray-400 rounded-full">
                   <img src={data.personalInfo.avatarUrl} alt="Avatar" className="w-24 h-24 rounded-full object-cover grayscale opacity-90" />
                </div>
             )}
             
             {/* 竖排名字装饰 */}
             <div className="absolute right-0 top-0 border border-red-800 p-1 hidden md:block opacity-80">
                <div className="border border-red-800 w-6 h-16 flex flex-col items-center justify-center text-xs text-red-900 font-bold writing-vertical-rl">
                   <span>求职</span>
                </div>
             </div>

             <h1 className="text-4xl font-normal mb-3 tracking-[0.2em] text-gray-800">{data.personalInfo.fullName}</h1>
             <div className="w-16 h-px bg-gray-400 mb-4"></div>
             <p className="text-lg text-gray-600 mb-4">{data.personalInfo.jobTitle}</p>
             
             <div className="flex gap-6 text-sm text-gray-500 font-light">
                <span>{data.personalInfo.phone}</span>
                <span>{data.personalInfo.email}</span>
             </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-x-16 gap-y-12">
             {(data.sections || []).filter(s => s.isVisible).map(section => (
                <div key={section.id} className={`${section.type === ResumeSectionType.Summary || section.type === ResumeSectionType.Projects ? 'md:col-span-2' : ''}`}>
                   <div className="flex items-center mb-6">
                      <div className="w-8 h-8 rounded-full border border-gray-300 flex items-center justify-center text-sm font-serif text-gray-500 mr-3 shadow-sm bg-[#faf9f6]">
                         {titleMap[section.type]?.[0] || section.title?.[0]}
                      </div>
                      <h3 className="text-xl font-normal tracking-widest text-gray-800 border-b border-gray-300 pb-1 flex-grow">
                         {titleMap[section.type] || section.title}
                      </h3>
                   </div>

                   <div className="space-y-6 pl-4 border-l border-gray-200 ml-4">
                      {section.type === ResumeSectionType.Skills ? (
                         <div className="flex flex-wrap gap-4">
                            {section.items.map(item => (
                               <div key={item.id} className="relative pl-4">
                                  <div className="absolute left-0 top-1.5 w-1.5 h-1.5 rounded-full bg-red-800 opacity-60"></div>
                                  <span className="text-gray-700">{item.description}</span>
                               </div>
                            ))}
                         </div>
                      ) : (
                         section.items.map(item => (
                            <div key={item.id} className="relative group">
                               <div className="absolute -left-[21px] top-1.5 w-2 h-2 rounded-full bg-white border border-gray-300 group-hover:bg-gray-200 transition-colors"></div>
                               <div className="flex justify-between items-baseline mb-2">
                                  <h4 className="text-lg font-bold text-gray-800">{item.title}</h4>
                                  {item.dateRange && <span className="text-xs text-gray-400 font-serif">{item.dateRange}</span>}
                               </div>
                               {item.subtitle && <div className="text-sm text-gray-600 mb-2 italic">{item.subtitle}</div>}
                               {item.description && (
                                  <div className="text-sm text-gray-600 leading-relaxed text-justify" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
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
  );
};
