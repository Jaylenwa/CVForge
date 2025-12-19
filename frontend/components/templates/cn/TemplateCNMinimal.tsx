import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useLanguage } from '../../../contexts/LanguageContext';
import { sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateCNMinimal: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const titleMap: Record<string, string> = {
    [ResumeSectionType.Summary]: '关于',
    [ResumeSectionType.Education]: '教育',
    [ResumeSectionType.Experience]: '经历',
    [ResumeSectionType.Skills]: '技能',
    [ResumeSectionType.Projects]: '项目',
    [ResumeSectionType.Custom]: '其他'
  };

  return (
    <div className={`w-full bg-white text-gray-900 h-full min-h-[1123px] ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto`} style={{ fontFamily: '"Helvetica Neue", Arial, sans-serif', lineHeight: '1.6' }}>
      <div className="p-16">
         {/* 头部 - 极度简化 */}
         <div className="mb-20">
            <h1 className="text-6xl font-light text-gray-900 mb-6 tracking-tighter">{data.personalInfo.fullName}</h1>
            <div className="text-xl text-gray-400 font-light mb-8">{data.personalInfo.jobTitle}</div>
            {data.personalInfo.avatarUrl && (
              <div className="mb-6">
                <img
                  src={data.personalInfo.avatarUrl}
                  alt={t('a11y.avatarAlt')}
                  className="w-24 h-24 rounded-md object-cover border border-gray-200"
                />
              </div>
            )}
            
            <div className="flex flex-col space-y-1 text-sm text-gray-500 font-mono">
               <div>{data.personalInfo.phone}</div>
               <div>{data.personalInfo.email}</div>
               {data.personalInfo.age && <div>{data.personalInfo.age} / {data.personalInfo.gender}</div>}
            </div>
         </div>

         {/* 内容 - 非对称网格 */}
         <div className="space-y-16">
            {(data.sections || []).filter(s => s.isVisible).map(section => (
               <div key={section.id} className="grid grid-cols-12 gap-8 group">
                  <div className="col-span-3">
                     <h3 className="text-sm font-bold uppercase tracking-[0.2em] text-gray-900 pt-1 group-hover:text-black transition-colors">
                        {titleMap[section.type] || section.title}
                     </h3>
                  </div>
                  
                  <div className="col-span-9 space-y-10 border-l border-gray-100 pl-8 -ml-8">
                     {section.type === ResumeSectionType.Skills ? (
                        <div className="flex flex-wrap gap-x-8 gap-y-2">
                           {section.items.map(item => (
                              <span key={item.id} className="text-sm text-gray-600 border-b border-gray-200 pb-1">{item.description}</span>
                           ))}
                        </div>
                     ) : (
                        section.items.map(item => (
                           <div key={item.id}>
                              <div className="flex justify-between items-baseline mb-2">
                                 <h4 className="text-lg font-medium text-gray-900">{item.title}</h4>
                                 {item.dateRange && <span className="text-xs text-gray-400 font-mono">{item.dateRange}</span>}
                              </div>
                              {item.subtitle && <div className="text-sm text-gray-500 mb-3">{item.subtitle}</div>}
                              {item.description && (
                                 <div className="text-sm text-gray-600 leading-relaxed font-light" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
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
