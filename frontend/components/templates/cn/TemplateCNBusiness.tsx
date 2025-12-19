import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useSectionTitle } from '../../../hooks/useSectionTitle';
import { useLanguage } from '../../../contexts/LanguageContext';
import { hasExtraPersonalInfo, sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateCNBusiness: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const titleMap: Record<string, string> = {
    [ResumeSectionType.Summary]: '个人概述',
    [ResumeSectionType.Education]: '教育经历',
    [ResumeSectionType.Experience]: '职业经历',
    [ResumeSectionType.Skills]: '核心能力',
    [ResumeSectionType.Projects]: '重点项目',
    [ResumeSectionType.Custom]: '附加信息'
  };

  const primaryColor = '#0f172a'; // slate-900

  return (
    <div className={`w-full bg-white h-full min-h-[1123px] ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto`} style={{ fontFamily: styles.fontFamily, lineHeight: parseFloat(styles.spacingMultiplier) * 1.5 }}>
      {/* 顶部几何背景 */}
      <div className="bg-slate-900 text-white p-10 relative overflow-hidden">
        <div className="absolute top-0 right-0 w-64 h-64 bg-white opacity-5 transform rotate-45 translate-x-32 -translate-y-32"></div>
        <div className="absolute bottom-0 left-0 w-48 h-48 bg-white opacity-5 transform rotate-45 -translate-x-24 translate-y-24"></div>
        
        <div className="relative z-10 flex flex-col md:flex-row items-center md:items-end justify-between gap-6">
           <div className="text-center md:text-left">
              <h1 className="text-4xl font-bold mb-2 tracking-wider">{data.personalInfo.fullName}</h1>
              <p className="text-lg text-slate-300 font-light uppercase tracking-widest mb-4">{data.personalInfo.jobTitle}</p>
              
              <div className="flex flex-wrap gap-4 text-xs text-slate-400 justify-center md:justify-start">
                  <span>{data.personalInfo.phone}</span>
                  <span>|</span>
                  <span>{data.personalInfo.email}</span>
                  {data.personalInfo.gender && (
                      <><span>|</span><span>{data.personalInfo.gender}</span></>
                  )}
                  {data.personalInfo.age && (
                      <><span>|</span><span>{data.personalInfo.age}</span></>
                  )}
              </div>
           </div>
           
           {data.personalInfo.avatarUrl && (
            <img 
              src={data.personalInfo.avatarUrl} 
              alt={t('a11y.avatarAlt')} 
              className="w-28 h-28 object-cover border-2 border-slate-700 shadow-xl"
            />
          )}
        </div>
      </div>

      <div className="p-10">
        <div className="grid grid-cols-1 gap-10">
            {(data.sections || []).filter(s => s.isVisible).map(section => (
                <div key={section.id}>
                    {/* 丝带标题样式 */}
                    <div className="flex items-center mb-6">
                        <div className="bg-slate-800 text-white py-1.5 px-6 font-bold text-lg relative shadow-md">
                            {titleMap[section.type] || section.title}
                            <div className="absolute top-0 right-0 -mr-2 w-0 h-0 border-t-[15px] border-t-slate-800 border-r-[10px] border-r-transparent border-b-[15px] border-b-slate-800"></div>
                        </div>
                        <div className="h-0.5 bg-slate-200 flex-grow ml-4"></div>
                    </div>

                    {section.type === ResumeSectionType.Skills ? (
                        <div className="bg-slate-50 p-6 border border-slate-100 rounded">
                             <div className="grid grid-cols-1 md:grid-cols-2 gap-x-8 gap-y-3">
                                {section.items.map(item => (
                                    <div key={item.id} className="flex items-start">
                                        <div className="w-1.5 h-1.5 bg-slate-800 mt-2 mr-2 rounded-full flex-shrink-0"></div>
                                        <span className="text-sm text-slate-700">{item.description}</span>
                                    </div>
                                ))}
                             </div>
                        </div>
                    ) : (
                        <div className="space-y-6">
                            {section.items.map((item, idx) => (
                                <div key={item.id} className={`p-4 ${idx % 2 === 0 ? 'bg-slate-50/50' : 'bg-white'} border-l-4 border-slate-200 hover:border-slate-800 transition-colors duration-300`}>
                                    <div className="flex flex-col md:flex-row md:justify-between md:items-baseline mb-2">
                                        <h4 className="font-bold text-slate-900 text-lg">{item.title}</h4>
                                        {item.dateRange && (
                                            <span className="text-sm text-slate-500 font-mono mt-1 md:mt-0">{item.dateRange}</span>
                                        )}
                                    </div>
                                    {item.subtitle && (
                                        <div className="text-slate-700 font-medium mb-3 italic">{item.subtitle}</div>
                                    )}
                                    {item.description && (
                                        <div 
                                            className="text-slate-600 text-sm leading-relaxed whitespace-pre-wrap"
                                            dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }}
                                        />
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

