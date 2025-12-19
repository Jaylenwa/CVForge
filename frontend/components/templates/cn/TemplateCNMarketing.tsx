import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useLanguage } from '../../../contexts/LanguageContext';
import { sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateCNMarketing: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const titleMap: Record<string, string> = {
    [ResumeSectionType.Summary]: '个人优势',
    [ResumeSectionType.Education]: '学习经历',
    [ResumeSectionType.Experience]: '实战经验',
    [ResumeSectionType.Skills]: '技能工具',
    [ResumeSectionType.Projects]: '项目成果',
    [ResumeSectionType.Custom]: '其他展示'
  };

  return (
    <div className={`w-full bg-[#ffde59] h-full min-h-[1123px] ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto overflow-hidden border-8 border-black`} style={{ fontFamily: '"Arial Black", "Impact", sans-serif' }}>
       {/* 斜纹背景装饰 */}
       <div className="absolute top-0 right-0 w-[500px] h-[500px] bg-black opacity-5 transform rotate-45 translate-x-1/2 -translate-y-1/2 pointer-events-none" style={{ backgroundImage: 'repeating-linear-gradient(45deg, transparent, transparent 10px, #000 10px, #000 12px)' }}></div>

       <div className="p-12 relative z-10">
          <div className="bg-white border-4 border-black shadow-[8px_8px_0px_0px_rgba(0,0,0,1)] p-8 mb-12 transform -rotate-1">
             <div className="flex justify-between items-center">
                <div>
                   <h1 className="text-5xl font-black mb-2 text-black uppercase tracking-tighter">{data.personalInfo.fullName}</h1>
                   <div className="text-2xl font-bold bg-black text-white px-4 py-1 inline-block transform rotate-2">{data.personalInfo.jobTitle}</div>
                </div>
                {data.personalInfo.avatarUrl && (
                   <div className="w-32 h-32 border-4 border-black rounded-full overflow-hidden shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-yellow-300">
                      <img src={data.personalInfo.avatarUrl} alt="Avatar" className="w-full h-full object-cover grayscale contrast-125" />
                   </div>
                )}
             </div>
             
             <div className="mt-6 flex flex-wrap gap-4 text-sm font-bold border-t-4 border-black pt-4">
                <span className="flex items-center"><span className="bg-black text-white rounded-full w-6 h-6 flex items-center justify-center mr-2">T</span>{data.personalInfo.phone}</span>
                <span className="flex items-center"><span className="bg-black text-white rounded-full w-6 h-6 flex items-center justify-center mr-2">E</span>{data.personalInfo.email}</span>
                {data.personalInfo.age && <span className="flex items-center"><span className="bg-black text-white rounded-full w-6 h-6 flex items-center justify-center mr-2">A</span>{data.personalInfo.age}</span>}
             </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
             {(data.sections || []).filter(s => s.isVisible).map((section, idx) => (
                <div key={section.id} className={`${(idx % 3 === 0) ? 'md:col-span-2' : ''} bg-white border-4 border-black p-6 shadow-[6px_6px_0px_0px_rgba(0,0,0,1)] hover:translate-x-1 hover:translate-y-1 hover:shadow-none transition-all`}>
                   <h3 className="text-2xl font-black mb-4 border-b-4 border-black pb-2 flex justify-between items-center">
                      {titleMap[section.type] || section.title}
                      <span className="text-4xl leading-none opacity-20">0{idx + 1}</span>
                   </h3>

                   {section.type === ResumeSectionType.Skills ? (
                      <div className="flex flex-wrap gap-3">
                         {section.items.map(item => (
                            <span key={item.id} className="bg-black text-white px-3 py-1 font-bold text-sm transform hover:scale-110 transition-transform cursor-default">{item.description}</span>
                         ))}
                      </div>
                   ) : (
                      <div className="space-y-6">
                         {section.items.map(item => (
                            <div key={item.id}>
                               <div className="flex flex-col md:flex-row md:justify-between md:items-start mb-2">
                                  <h4 className="text-xl font-bold">{item.title}</h4>
                                  {item.dateRange && <span className="bg-yellow-300 px-2 border-2 border-black text-xs font-bold self-start md:self-auto">{item.dateRange}</span>}
                               </div>
                               {item.subtitle && <div className="font-bold text-gray-600 mb-2">{item.subtitle}</div>}
                               {item.description && (
                                  <div className="text-sm font-medium text-gray-800 leading-snug" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
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
