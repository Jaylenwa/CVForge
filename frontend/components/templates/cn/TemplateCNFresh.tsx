import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useLanguage } from '../../../contexts/LanguageContext';
import { sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateCNFresh: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const titleMap: Record<string, string> = {
    [ResumeSectionType.Summary]: '自我介绍',
    [ResumeSectionType.Education]: '教育背景',
    [ResumeSectionType.Experience]: '工作经历',
    [ResumeSectionType.Skills]: '技能特长',
    [ResumeSectionType.Projects]: '项目经验',
    [ResumeSectionType.Custom]: '其他'
  };

  return (
    <div className={`w-full bg-[#f0f9f4] text-[#2d3748] h-full min-h-[1123px] ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto`} style={{ fontFamily: '"Nunito", "Microsoft YaHei", sans-serif' }}>
       {/* 装饰性背景 */}
       <div className="absolute top-0 right-0 w-64 h-64 bg-[#c6f6d5] rounded-bl-full opacity-50"></div>
       <div className="absolute bottom-0 left-0 w-48 h-48 bg-[#9ae6b4] rounded-tr-full opacity-40"></div>
       
       <div className="p-12 relative z-10">
          <div className="bg-white rounded-3xl shadow-sm p-8 mb-10 flex items-center gap-8 relative overflow-hidden">
             <div className="absolute top-0 left-0 w-2 h-full bg-[#48bb78]"></div>
             {data.personalInfo.avatarUrl && (
                <div className="w-24 h-24 rounded-2xl overflow-hidden flex-shrink-0 shadow-md ring-4 ring-[#f0f9f4]">
                   <img src={data.personalInfo.avatarUrl} alt="Avatar" className="w-full h-full object-cover" />
                </div>
             )}
             <div className="flex-grow">
                <h1 className="text-3xl font-bold text-[#276749] mb-2">{data.personalInfo.fullName}</h1>
                <div className="text-lg text-[#38a169] font-medium mb-3">{data.personalInfo.jobTitle}</div>
                <div className="flex flex-wrap gap-4 text-sm text-gray-500">
                   <span className="flex items-center gap-1"><span className="w-2 h-2 rounded-full bg-[#48bb78]"></span>{data.personalInfo.phone}</span>
                   <span className="flex items-center gap-1"><span className="w-2 h-2 rounded-full bg-[#48bb78]"></span>{data.personalInfo.email}</span>
                   {data.personalInfo.city && <span className="flex items-center gap-1"><span className="w-2 h-2 rounded-full bg-[#48bb78]"></span>{data.personalInfo.city}</span>}
                </div>
             </div>
          </div>

          <div className="grid grid-cols-3 gap-8">
             {/* 左侧 - 技能与摘要 */}
             <div className="col-span-1 space-y-8">
                {(data.sections || []).filter(s => (s.type === ResumeSectionType.Summary || s.type === ResumeSectionType.Skills) && s.isVisible).map(section => (
                   <div key={section.id} className="bg-white rounded-2xl p-6 shadow-sm">
                      <h3 className="text-lg font-bold text-[#2f855a] mb-4 flex items-center">
                         <span className="w-8 h-1 bg-[#48bb78] rounded-full mr-2"></span>
                         {titleMap[section.type] || section.title}
                      </h3>
                      {section.type === ResumeSectionType.Skills ? (
                         <div className="flex flex-wrap gap-2">
                            {section.items.map(item => (
                               <span key={item.id} className="px-3 py-1 bg-[#f0fff4] text-[#276749] text-sm rounded-full border border-[#c6f6d5]">{item.description}</span>
                            ))}
                         </div>
                      ) : (
                         <div className="text-sm text-gray-600 leading-relaxed">
                            {section.items.map(item => (
                               <div key={item.id} dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description || '') }} />
                            ))}
                         </div>
                      )}
                   </div>
                ))}
             </div>

             {/* 右侧 - 主要经历 */}
             <div className="col-span-2 space-y-8">
                {(data.sections || []).filter(s => s.type !== ResumeSectionType.Summary && s.type !== ResumeSectionType.Skills && s.isVisible).map(section => (
                   <div key={section.id} className="bg-white/60 rounded-2xl p-6 hover:bg-white transition-colors duration-300">
                      <h3 className="text-xl font-bold text-[#2f855a] mb-6 border-b border-[#c6f6d5] pb-2 inline-block pr-12 relative">
                         {titleMap[section.type] || section.title}
                         <span className="absolute bottom-0 right-0 w-8 h-1 bg-[#48bb78] rounded-t-full"></span>
                      </h3>

                      <div className="space-y-8">
                         {section.items.map(item => (
                            <div key={item.id} className="relative pl-6 border-l-2 border-[#e2e8f0]">
                               <div className="absolute -left-[9px] top-1.5 w-4 h-4 bg-[#c6f6d5] rounded-full border-2 border-white"></div>
                               <div className="flex justify-between items-baseline mb-2">
                                  <h4 className="text-lg font-bold text-gray-800">{item.title}</h4>
                                  {item.dateRange && <span className="text-sm text-[#38a169] bg-[#f0fff4] px-2 py-0.5 rounded-md">{item.dateRange}</span>}
                               </div>
                               {item.subtitle && <div className="text-sm text-gray-600 font-medium mb-2">{item.subtitle}</div>}
                               {item.description && (
                                  <div className="text-sm text-gray-600 leading-relaxed" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
                               )}
                            </div>
                         ))}
                      </div>
                   </div>
                ))}
             </div>
          </div>
       </div>
    </div>
  );
};
