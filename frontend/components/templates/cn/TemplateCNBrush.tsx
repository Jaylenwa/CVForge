import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useLanguage } from '../../../contexts/LanguageContext';
import { sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateCNBrush: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const titleMap: Record<string, string> = {
    [ResumeSectionType.Summary]: '自我评价',
    [ResumeSectionType.Education]: '教育经历',
    [ResumeSectionType.Experience]: '工作经历',
    [ResumeSectionType.Skills]: '技能特长',
    [ResumeSectionType.Projects]: '项目经验',
    [ResumeSectionType.Custom]: '其他'
  };

  return (
    <div className={`w-full bg-[#fdfbf7] text-[#2c3e50] h-auto ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none`} style={{ fontFamily: '"KaiTi", "STKaiti", serif' }}>
      {/* 水墨背景装饰 */}
      <div className="absolute top-0 right-0 w-64 h-64 opacity-10 pointer-events-none" style={{ backgroundImage: 'radial-gradient(circle at center, black 0%, transparent 70%)', filter: 'blur(20px)' }}></div>
      <div className="absolute bottom-0 left-0 w-96 h-96 opacity-5 pointer-events-none" style={{ backgroundImage: 'radial-gradient(circle at center, black 0%, transparent 70%)', filter: 'blur(40px)' }}></div>

      <div className="p-12 relative z-10">
        {/* 竖排标题 */}
        <div className="absolute top-12 right-12 writing-vertical-rl text-4xl font-bold tracking-[0.2em] text-gray-800 border-l-2 border-red-800 pl-4 h-64 flex items-center">
            {data.personalInfo.fullName}
        </div>
        
        {/* 印章效果 */}
        <div className="absolute top-[340px] right-14 w-12 h-12 border-2 border-red-700 rounded text-red-700 flex items-center justify-center text-xs font-bold transform rotate-3 opacity-80" style={{ boxShadow: 'inset 0 0 5px rgba(185, 28, 28, 0.5)' }}>
            {data.personalInfo.fullName?.slice(0, 2) || '简历'}
        </div>

        <div className="mr-32"> {/* 为右侧竖排留空 */}
            <div className="mb-16 pt-8">
                <div className="text-2xl font-bold text-gray-600 mb-6">{data.personalInfo.jobTitle}</div>
                {data.personalInfo.avatarUrl && (
                  <div className="mb-4">
                    <img
                      src={data.personalInfo.avatarUrl}
                      alt={t('a11y.avatarAlt')}
                      className="w-24 h-24 rounded-md object-cover border-2 border-red-800 shadow-sm"
                    />
                  </div>
                )}
                <div className="flex flex-wrap gap-6 text-sm text-gray-500 font-sans border-t border-b border-gray-300 py-3">
                    <span>{data.personalInfo.phone}</span>
                    <span>{data.personalInfo.email}</span>
                    {data.personalInfo.city && <span>{data.personalInfo.city}</span>}
                    {data.personalInfo.gender && <span>{data.personalInfo.gender}</span>}
                    {data.personalInfo.age && <span>{data.personalInfo.age}</span>}
                </div>
            </div>

            <div className="space-y-12">
                {(data.sections || []).filter(s => s.isVisible).map(section => (
                    <div key={section.id}>
                        <div className="relative mb-6">
                            {/* 笔刷风格标题背景 */}
                            <div className="absolute -left-4 -top-2 w-32 h-10 bg-gray-200 transform -skew-x-12 opacity-50 rounded-sm"></div>
                            <h3 className="text-2xl font-bold relative z-10 flex items-center gap-2">
                                <span className="w-2 h-8 bg-red-800 block"></span>
                                {titleMap[section.type] || section.title}
                            </h3>
                        </div>

                        {section.type === ResumeSectionType.Skills ? (
                            <div className="flex flex-wrap gap-4 px-4">
                                {section.items.map(item => (
                                    <div key={item.id} className="flex items-center gap-2">
                                        <span className="w-1.5 h-1.5 rounded-full bg-red-800"></span>
                                        <span className="text-base">{item.description}</span>
                                    </div>
                                ))}
                            </div>
                        ) : (
                            <div className="space-y-8 px-4 border-l border-gray-200 ml-1">
                                {section.items.map(item => (
                                    <div key={item.id}>
                                        <div className="flex justify-between items-baseline mb-2">
                                            <h4 className="text-xl font-bold text-gray-800">{item.title}</h4>
                                            {item.dateRange && <span className="text-sm text-gray-500 font-sans">{item.dateRange}</span>}
                                        </div>
                                        {item.subtitle && <div className="text-lg text-gray-600 mb-2 italic">{item.subtitle}</div>}
                                        {item.description && (
                                            <div className="text-base text-gray-700 leading-loose text-justify" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
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
    </div>
  );
};
