import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useLanguage } from '../../../contexts/LanguageContext';
import { sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateCNFinance: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const titleMap: Record<string, string> = {
    [ResumeSectionType.Summary]: 'PROFESSIONAL PROFILE',
    [ResumeSectionType.Education]: 'EDUCATION',
    [ResumeSectionType.Experience]: 'WORK EXPERIENCE',
    [ResumeSectionType.Skills]: 'CORE COMPETENCIES',
    [ResumeSectionType.Projects]: 'KEY PROJECTS',
    [ResumeSectionType.Custom]: 'ADDITIONAL INFO'
  };

  return (
    <div className={`w-full bg-white text-[#0a192f] h-full min-h-[1123px] ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto`} style={{ fontFamily: '"Calibri", "Arial", sans-serif' }}>
      {/* 顶部金线 */}
      <div className="h-2 bg-[#b8860b]"></div>
      
      <div className="p-10">
        {/* 头部 */}
        <div className="flex justify-between items-end border-b-2 border-[#0a192f] pb-6 mb-10">
            <div>
                <h1 className="text-4xl font-bold text-[#0a192f] uppercase mb-1">{data.personalInfo.fullName}</h1>
                <p className="text-lg text-[#b8860b] font-bold uppercase tracking-wider">{data.personalInfo.jobTitle}</p>
            </div>
            <div className="text-right text-sm text-gray-600 leading-snug">
                {data.personalInfo.avatarUrl && (
                  <img
                    src={data.personalInfo.avatarUrl}
                    alt={t('a11y.avatarAlt')}
                    className="w-20 h-20 object-cover rounded-md border border-[#b8860b] ml-auto mb-2"
                  />
                )}
                <div>{data.personalInfo.phone}</div>
                <div>{data.personalInfo.email}</div>
                {data.personalInfo.city && <div>{data.personalInfo.city}</div>}
            </div>
        </div>

        {/* 内容 */}
        <div className="space-y-8">
            {(data.sections || []).filter(s => s.isVisible).map(section => (
                <div key={section.id}>
                    <h3 className="bg-[#f0f4f8] text-[#0a192f] font-bold text-sm uppercase px-4 py-2 mb-4 border-l-4 border-[#b8860b]">
                        {titleMap[section.type] || section.title}
                    </h3>

                    <div className="px-2">
                        {section.type === ResumeSectionType.Skills ? (
                            <div className="grid grid-cols-3 gap-4">
                                {section.items.map(item => (
                                    <div key={item.id} className="flex items-center border-b border-gray-100 pb-1">
                                        <div className="w-2 h-2 bg-[#b8860b] rounded-full mr-2"></div>
                                        <span className="text-sm font-medium">{item.description}</span>
                                    </div>
                                ))}
                            </div>
                        ) : (
                            <div className="space-y-6">
                                {section.items.map(item => (
                                    <div key={item.id}>
                                        <div className="flex justify-between font-bold text-[#0a192f] mb-1">
                                            <div className="text-base">{item.title}</div>
                                            <div className="text-sm text-[#b8860b]">{item.dateRange}</div>
                                        </div>
                                        {item.subtitle && <div className="text-sm font-bold text-gray-600 mb-2 italic">{item.subtitle}</div>}
                                        {item.description && (
                                            <div className="text-sm text-gray-700 leading-snug" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
                                        )}
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>
                </div>
            ))}
        </div>
      </div>
    </div>
  );
};
