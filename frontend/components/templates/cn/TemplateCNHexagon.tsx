import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useLanguage } from '../../../contexts/LanguageContext';
import { sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateCNHexagon: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const titleMap: Record<string, string> = {
    [ResumeSectionType.Summary]: '个人简介',
    [ResumeSectionType.Education]: '教育背景',
    [ResumeSectionType.Experience]: '工作经历',
    [ResumeSectionType.Skills]: '技能专长',
    [ResumeSectionType.Projects]: '项目经验',
    [ResumeSectionType.Custom]: '其他'
  };

  return (
    <div className={`w-full bg-slate-50 text-slate-800 h-full min-h-[1123px] ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto`} style={{ fontFamily: '"Microsoft YaHei", sans-serif' }}>
      <div className="flex h-full">
        {/* 左侧深色栏 */}
        <div className="w-1/3 bg-slate-800 text-white p-8 relative overflow-hidden">
            {/* 六边形背景装饰 */}
            <div className="absolute top-0 left-0 w-full h-full opacity-5 pointer-events-none">
                <svg width="100%" height="100%" patternUnits="userSpaceOnUse">
                    <defs>
                        <pattern id="hexagons" width="50" height="43.4" patternUnits="userSpaceOnUse" patternTransform="scale(2)">
                            <path d="M25 0 L50 14.4 L50 43.3 L25 57.7 L0 43.3 L0 14.4 Z" fill="none" stroke="currentColor" strokeWidth="1"/>
                        </pattern>
                    </defs>
                    <rect width="100%" height="100%" fill="url(#hexagons)" />
                </svg>
            </div>

            <div className="relative z-10 text-center">
                {data.personalInfo.avatarUrl && (
                    <div className="w-40 h-40 mx-auto mb-6 relative">
                        <div className="absolute inset-0 bg-amber-400 transform rotate-6" style={{ clipPath: 'polygon(50% 0%, 100% 25%, 100% 75%, 50% 100%, 0% 75%, 0% 25%)' }}></div>
                        <img src={data.personalInfo.avatarUrl} className="w-full h-full object-cover relative z-10" style={{ clipPath: 'polygon(50% 0%, 100% 25%, 100% 75%, 50% 100%, 0% 75%, 0% 25%)' }} alt="Avatar" />
                    </div>
                )}
                
                <h1 className="text-2xl font-bold mb-2">{data.personalInfo.fullName}</h1>
                <p className="text-amber-400 mb-8 font-medium">{data.personalInfo.jobTitle}</p>

                <div className="space-y-6 text-left text-sm text-slate-300">
                    <div className="flex items-center gap-3 bg-slate-700/50 p-3 rounded-lg border border-slate-600">
                        <span className="text-amber-400 font-bold">P:</span> {data.personalInfo.phone}
                    </div>
                    <div className="flex items-center gap-3 bg-slate-700/50 p-3 rounded-lg border border-slate-600">
                        <span className="text-amber-400 font-bold">E:</span> {data.personalInfo.email}
                    </div>
                    {data.personalInfo.city && (
                        <div className="flex items-center gap-3 bg-slate-700/50 p-3 rounded-lg border border-slate-600">
                            <span className="text-amber-400 font-bold">L:</span> {data.personalInfo.city}
                        </div>
                    )}
                </div>

                <div className="mt-12 space-y-8 text-left">
                    {(data.sections || []).filter(s => s.type === ResumeSectionType.Skills && s.isVisible).map(section => (
                        <div key={section.id}>
                            <h3 className="text-lg font-bold text-amber-400 border-b border-slate-600 pb-2 mb-4">{titleMap[section.type] || section.title}</h3>
                            <div className="flex flex-wrap gap-2">
                                {section.items.map(item => (
                                    <span key={item.id} className="border border-amber-400/50 text-amber-100 px-3 py-1 text-sm rounded-full bg-amber-400/10">
                                        {item.description}
                                    </span>
                                ))}
                            </div>
                        </div>
                    ))}
                </div>
            </div>
        </div>

        {/* 右侧内容 */}
        <div className="w-2/3 p-10 bg-slate-50">
            <div className="space-y-10">
                {(data.sections || []).filter(s => s.type !== ResumeSectionType.Skills && s.isVisible).map(section => (
                    <div key={section.id}>
                        <div className="flex items-center gap-4 mb-6">
                            <div className="w-10 h-10 bg-slate-800 flex items-center justify-center text-amber-400 font-bold text-xl" style={{ clipPath: 'polygon(50% 0%, 100% 25%, 100% 75%, 50% 100%, 0% 75%, 0% 25%)' }}>
                                {(section.title?.charAt(0) || '').toUpperCase()}
                            </div>
                            <h3 className="text-2xl font-bold text-slate-800">{titleMap[section.type] || section.title}</h3>
                            <div className="flex-grow h-px bg-slate-200"></div>
                        </div>

                        <div className="space-y-8 pl-4 border-l-2 border-slate-200 ml-5">
                            {section.items.map(item => (
                                <div key={item.id} className="relative pl-6">
                                    <div className="absolute -left-[9px] top-1.5 w-4 h-4 bg-amber-400" style={{ clipPath: 'polygon(50% 0%, 100% 25%, 100% 75%, 50% 100%, 0% 75%, 0% 25%)' }}></div>
                                    <div className="flex justify-between items-baseline mb-2">
                                        <h4 className="text-lg font-bold text-slate-700">{item.title}</h4>
                                        {item.dateRange && <span className="text-sm text-slate-500 font-mono bg-slate-100 px-2 rounded">{item.dateRange}</span>}
                                    </div>
                                    {item.subtitle && <div className="text-base text-amber-600 mb-2 font-medium">{item.subtitle}</div>}
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
