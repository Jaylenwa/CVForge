import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useLanguage } from '../../../contexts/LanguageContext';
import { sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateCNBlock: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const titleMap: Record<string, string> = {
    [ResumeSectionType.Summary]: '概览',
    [ResumeSectionType.Education]: '教育',
    [ResumeSectionType.Experience]: '经历',
    [ResumeSectionType.Skills]: '技能',
    [ResumeSectionType.Projects]: '项目',
    [ResumeSectionType.Custom]: '其他'
  };

  return (
    <div className={`w-full bg-gray-100 text-gray-800 h-full min-h-[1123px] ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto`} style={{ fontFamily: '"Microsoft YaHei", sans-serif' }}>
      <div className="grid grid-cols-12 min-h-full">
        {/* 侧边栏 */}
        <div className="col-span-4 bg-[#2d3436] text-white p-8 flex flex-col">
            <div className="mb-12 text-center">
                {data.personalInfo.avatarUrl && (
                    <div className="w-32 h-32 mx-auto bg-white rounded-full p-1 mb-4">
                        <img src={data.personalInfo.avatarUrl} className="w-full h-full rounded-full object-cover" alt="Avatar" />
                    </div>
                )}
                <h1 className="text-2xl font-bold mb-2">{data.personalInfo.fullName}</h1>
                <p className="text-[#dfe6e9]">{data.personalInfo.jobTitle}</p>
            </div>

            <div className="space-y-8 flex-grow">
                <div className="bg-[#636e72] p-4 rounded-lg bg-opacity-30">
                    <div className="text-sm space-y-2 opacity-90">
                        <div>{data.personalInfo.phone}</div>
                        <div>{data.personalInfo.email}</div>
                        <div>{data.personalInfo.city}</div>
                    </div>
                </div>

                {(data.sections || []).filter(s => s.type === ResumeSectionType.Skills && s.isVisible).map(section => (
                    <div key={section.id}>
                        <h3 className="text-lg font-bold border-b border-gray-500 pb-2 mb-4 text-[#00cec9]">{titleMap[section.type] || section.title}</h3>
                        <div className="flex flex-wrap gap-2">
                            {section.items.map(item => (
                                <span key={item.id} className="bg-[#0984e3] px-3 py-1 rounded text-sm">{item.description}</span>
                            ))}
                        </div>
                    </div>
                ))}
            </div>
        </div>

        {/* 主内容 */}
        <div className="col-span-8 p-8 space-y-8">
            {(data.sections || []).filter(s => s.type !== ResumeSectionType.Skills && s.isVisible).map((section, index) => (
                <div key={section.id} className={`bg-white p-6 rounded-lg shadow-sm border-l-4 ${index % 2 === 0 ? 'border-[#00b894]' : 'border-[#e17055]'}`}>
                    <h3 className={`text-xl font-bold mb-6 ${index % 2 === 0 ? 'text-[#00b894]' : 'text-[#e17055]'}`}>
                        {titleMap[section.type] || section.title}
                    </h3>

                    <div className="space-y-6">
                        {section.items.map(item => (
                            <div key={item.id} className="relative">
                                <div className="flex justify-between items-baseline mb-2">
                                    <h4 className="text-lg font-bold text-gray-800">{item.title}</h4>
                                    {item.dateRange && <span className="text-sm text-gray-500 bg-gray-100 px-2 py-1 rounded">{item.dateRange}</span>}
                                </div>
                                {item.subtitle && <div className="text-sm font-medium text-gray-600 mb-2">{item.subtitle}</div>}
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
  );
};

