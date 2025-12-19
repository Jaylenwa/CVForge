import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useLanguage } from '../../../contexts/LanguageContext';
import { sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateCNCloud: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const titleMap: Record<string, string> = {
    [ResumeSectionType.Summary]: '自我介绍',
    [ResumeSectionType.Education]: '教育背景',
    [ResumeSectionType.Experience]: '工作经历',
    [ResumeSectionType.Skills]: '擅长技能',
    [ResumeSectionType.Projects]: '项目展示',
    [ResumeSectionType.Custom]: '其他'
  };

  return (
    <div className={`w-full bg-blue-50 text-slate-700 h-full min-h-[1123px] ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto`} style={{ fontFamily: '"Nunito", "Microsoft YaHei", sans-serif' }}>
      {/* 云朵背景装饰 */}
      <div className="absolute top-0 left-0 w-full h-64 bg-white rounded-b-[50%] shadow-sm overflow-hidden">
          <div className="absolute top-10 left-10 w-20 h-20 bg-blue-100 rounded-full opacity-50"></div>
          <div className="absolute top-20 right-20 w-32 h-32 bg-blue-100 rounded-full opacity-50"></div>
      </div>

      <div className="p-10 relative z-10">
        <div className="text-center mb-12">
            <div className="inline-block p-2 bg-white rounded-full shadow-md mb-4">
                {data.personalInfo.avatarUrl && (
                    <img src={data.personalInfo.avatarUrl} className="w-32 h-32 rounded-full object-cover border-4 border-blue-50" alt="Avatar" />
                )}
            </div>
            <h1 className="text-3xl font-bold text-slate-800 mb-2">{data.personalInfo.fullName}</h1>
            <div className="inline-block bg-blue-100 text-blue-600 px-4 py-1 rounded-full text-sm font-bold mb-4">{data.personalInfo.jobTitle}</div>
            
            <div className="flex justify-center gap-6 text-sm text-slate-500">
                <span className="flex items-center gap-1">
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"></path></svg>
                    {data.personalInfo.phone}
                </span>
                <span className="flex items-center gap-1">
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"></path></svg>
                    {data.personalInfo.email}
                </span>
            </div>
        </div>

        <div className="space-y-8">
            {(data.sections || []).filter(s => s.isVisible).map(section => (
                <div key={section.id} className="bg-white p-8 rounded-3xl shadow-sm border border-blue-50">
                    <h3 className="text-xl font-bold text-blue-500 mb-6 flex items-center">
                        <span className="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center mr-3 text-blue-500">
                            {section.title?.charAt(0) || ''}
                        </span>
                        {titleMap[section.type] || section.title}
                    </h3>

                    {section.type === ResumeSectionType.Skills ? (
                        <div className="flex flex-wrap gap-3">
                            {section.items.map(item => (
                                <span key={item.id} className="bg-blue-50 text-blue-600 px-4 py-2 rounded-xl text-sm font-medium border border-blue-100">
                                    {item.description}
                                </span>
                            ))}
                        </div>
                    ) : (
                        <div className="space-y-8">
                            {section.items.map(item => (
                                <div key={item.id}>
                                    <div className="flex justify-between items-baseline mb-2">
                                        <h4 className="text-lg font-bold text-slate-800">{item.title}</h4>
                                        {item.dateRange && <span className="text-sm text-slate-400 bg-slate-50 px-3 py-1 rounded-full">{item.dateRange}</span>}
                                    </div>
                                    {item.subtitle && <div className="text-base text-blue-400 mb-2 font-medium">{item.subtitle}</div>}
                                    {item.description && (
                                        <div className="text-sm text-slate-600 leading-relaxed" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
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
