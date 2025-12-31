import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useLanguage } from '../../../contexts/LanguageContext';
import { sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateCNOrigami: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const titleMap: Record<string, string> = {
    [ResumeSectionType.Summary]: '简介',
    [ResumeSectionType.Education]: '教育',
    [ResumeSectionType.Experience]: '经历',
    [ResumeSectionType.Skills]: '技能',
    [ResumeSectionType.Projects]: '项目',
    [ResumeSectionType.Custom]: '其他'
  };

  return (
    <div className={`w-full bg-[#f0f2f5] text-slate-800 h-auto ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none`} style={{ fontFamily: '"Microsoft YaHei", sans-serif' }}>
      <div className="p-8">
        {/* 头部卡片 */}
        <div className="bg-white p-8 rounded-lg shadow-sm mb-8 relative overflow-hidden">
            <div className="absolute top-0 right-0 w-32 h-32 bg-blue-500 transform rotate-45 translate-x-16 -translate-y-16"></div>
            
            <div className="flex gap-6 items-center relative z-10">
                {data.personalInfo.avatarUrl && (
                    <div className="w-24 h-24 bg-white p-1 shadow-md transform -rotate-2">
                        <img src={data.personalInfo.avatarUrl} className="w-full h-full object-cover" alt="Avatar" />
                    </div>
                )}
                <div>
                    <h1 className="text-3xl font-bold text-slate-800 mb-1">{data.personalInfo.fullName}</h1>
                    <p className="text-blue-600 font-medium text-lg mb-4">{data.personalInfo.jobTitle}</p>
                    <div className="flex gap-4 text-sm text-slate-500">
                        <span>{data.personalInfo.phone}</span>
                        <span>{data.personalInfo.email}</span>
                        {data.personalInfo.city && <span>{data.personalInfo.city}</span>}
                    </div>
                </div>
            </div>
        </div>

        <div className="grid grid-cols-12 gap-6">
            <div className="col-span-4 space-y-6">
                 {(data.sections || []).filter(s => s.type === ResumeSectionType.Skills && s.isVisible).map(section => (
                    <div key={section.id} className="bg-white rounded-lg shadow-sm overflow-hidden">
                        {/* 折纸标题效果 */}
                        <div className="bg-blue-600 text-white py-2 px-4 relative shadow-md">
                            <div className="absolute top-full left-0 w-2 h-2 bg-blue-800" style={{ clipPath: 'polygon(0 0, 100% 0, 100% 100%)' }}></div>
                            <h3 className="font-bold">{titleMap[section.type] || section.title}</h3>
                        </div>
                        <div className="p-4 pt-6">
                            <div className="space-y-3">
                                {section.items.map(item => (
                                    <div key={item.id}>
                                        <div className="text-sm font-medium mb-1 flex justify-between">
                                            <span>{item.description}</span>
                                            <span className="text-blue-400">●●●●○</span>
                                        </div>
                                        <div className="h-1.5 w-full bg-gray-100 rounded-full overflow-hidden">
                                            <div className="h-full bg-blue-400 w-4/5 rounded-full"></div>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        </div>
                    </div>
                 ))}

                 {/* 个人信息 */}
                 <div className="bg-white rounded-lg shadow-sm overflow-hidden">
                    <div className="bg-slate-700 text-white py-2 px-4 relative shadow-md">
                        <h3 className="font-bold">基本信息</h3>
                    </div>
                    <div className="p-4 text-sm space-y-2 text-slate-600">
                        {data.personalInfo.gender && <div className="flex justify-between border-b border-gray-100 pb-1"><span>性别</span><span>{data.personalInfo.gender}</span></div>}
                        {data.personalInfo.age && <div className="flex justify-between border-b border-gray-100 pb-1"><span>年龄</span><span>{data.personalInfo.age}</span></div>}
                        {data.personalInfo.politicalStatus && <div className="flex justify-between border-b border-gray-100 pb-1"><span>政治面貌</span><span>{data.personalInfo.politicalStatus}</span></div>}
                    </div>
                 </div>
            </div>

            <div className="col-span-8 space-y-6">
                {(data.sections || []).filter(s => s.type !== ResumeSectionType.Skills && s.isVisible).map(section => (
                    <div key={section.id} className="bg-white rounded-lg shadow-sm p-6 relative">
                        {/* 侧边标签 */}
                        <div className="absolute -left-2 top-6 bg-blue-600 text-white py-1 px-3 shadow-md z-10">
                             <div className="absolute top-full left-0 w-2 h-2 bg-blue-800" style={{ clipPath: 'polygon(0 0, 100% 0, 100% 100%)' }}></div>
                             <h3 className="font-bold">{titleMap[section.type] || section.title}</h3>
                        </div>

                        <div className="ml-4 space-y-8 mt-4">
                            {section.items.map(item => (
                                <div key={item.id} className="relative border-l-2 border-gray-100 pl-6 pb-2">
                                    <div className="absolute -left-[7px] top-1.5 w-3 h-3 bg-white border-2 border-blue-400 rounded-full"></div>
                                    <div className="flex justify-between items-baseline mb-1">
                                        <h4 className="text-lg font-bold text-slate-800">{item.title}</h4>
                                        {item.dateRange && <span className="text-sm text-slate-500 bg-slate-50 px-2 py-0.5 rounded">{item.dateRange}</span>}
                                    </div>
                                    {item.subtitle && <div className="text-blue-600 font-medium mb-2 text-sm">{item.subtitle}</div>}
                                    {item.description && (
                                        <div className="text-sm text-slate-600 leading-relaxed" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
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
