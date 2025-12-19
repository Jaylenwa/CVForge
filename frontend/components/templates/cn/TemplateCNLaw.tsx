import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useLanguage } from '../../../contexts/LanguageContext';
import { sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateCNLaw: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const titleMap: Record<string, string> = {
    [ResumeSectionType.Summary]: '专业概述',
    [ResumeSectionType.Education]: '教育背景',
    [ResumeSectionType.Experience]: '法律实务',
    [ResumeSectionType.Skills]: '执业资格',
    [ResumeSectionType.Projects]: '代表案例',
    [ResumeSectionType.Custom]: '社会兼职'
  };

  return (
    <div className={`w-full bg-[#fdfdfd] text-[#1a1a1a] h-full min-h-[1123px] ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto border-t-8 border-[#2c3e50]`} style={{ fontFamily: '"Times New Roman", "Songti SC", serif' }}>
      <div className="p-12 max-w-[900px] mx-auto">
        {/* 头部 */}
        <div className="text-center mb-16 relative">
            <h1 className="text-4xl font-bold tracking-widest text-[#2c3e50] mb-4 uppercase">{data.personalInfo.fullName}</h1>
            <div className="flex items-center justify-center gap-4 mb-6">
                <div className="h-px bg-gray-400 w-12"></div>
                <div className="text-lg italic text-gray-600">{data.personalInfo.jobTitle}</div>
                <div className="h-px bg-gray-400 w-12"></div>
            </div>
            
            <div className="flex justify-center gap-6 text-sm text-[#505050] font-sans">
                <span>{data.personalInfo.phone}</span>
                <span>|</span>
                <span>{data.personalInfo.email}</span>
                {data.personalInfo.city && (
                    <>
                        <span>|</span>
                        <span>{data.personalInfo.city}</span>
                    </>
                )}
            </div>
            
            {/* 天平背景水印 */}
            <div className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 opacity-[0.03] pointer-events-none">
                <svg width="300" height="300" viewBox="0 0 24 24" fill="currentColor">
                    <path d="M12 2L1 12h22L12 2zm0 3.5L18.5 10H5.5L12 5.5zM12 13c-3.31 0-6 2.69-6 6h12c0-3.31-2.69-6-6-6z"/>
                </svg>
            </div>
        </div>

        {/* 内容 */}
        <div className="space-y-10">
            {(data.sections || []).filter(s => s.isVisible).map(section => (
                <div key={section.id}>
                    <div className="flex items-center mb-6">
                        <h3 className="text-xl font-bold text-[#2c3e50] uppercase tracking-wider border-b-2 border-[#2c3e50] pb-1 pr-8">
                            {titleMap[section.type] || section.title}
                        </h3>
                    </div>

                    {section.type === ResumeSectionType.Skills ? (
                        <div className="grid grid-cols-2 gap-4 bg-gray-50 p-6 border border-gray-200">
                            {section.items.map(item => (
                                <div key={item.id} className="flex items-center gap-2">
                                    <div className="w-1.5 h-1.5 bg-[#2c3e50] transform rotate-45"></div>
                                    <span className="text-sm font-medium">{item.description}</span>
                                </div>
                            ))}
                        </div>
                    ) : (
                        <div className="space-y-8">
                            {section.items.map(item => (
                                <div key={item.id}>
                                    <div className="flex justify-between items-baseline mb-1">
                                        <h4 className="text-lg font-bold text-[#2c3e50]">{item.title}</h4>
                                        {item.dateRange && <span className="text-sm font-sans text-gray-500">{item.dateRange}</span>}
                                    </div>
                                    {item.subtitle && <div className="text-base text-gray-700 italic mb-2 font-medium">{item.subtitle}</div>}
                                    {item.description && (
                                        <div className="text-sm leading-relaxed text-justify text-gray-800" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
                                    )}
                                </div>
                            ))}
                        </div>
                    )}
                </div>
            ))}
        </div>
        
        {/* 页脚 */}
        <div className="mt-16 pt-8 border-t border-gray-200 text-center text-xs text-gray-400 font-sans">
            CONFIDENTIAL • {new Date().getFullYear()}
        </div>
      </div>
    </div>
  );
};

