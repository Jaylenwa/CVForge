import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useLanguage } from '../../../contexts/LanguageContext';
import { sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateCNMedical: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const titleMap: Record<string, string> = {
    [ResumeSectionType.Summary]: '执业概述',
    [ResumeSectionType.Education]: '医学教育',
    [ResumeSectionType.Experience]: '临床经历',
    [ResumeSectionType.Skills]: '专业技能',
    [ResumeSectionType.Projects]: '科研成果',
    [ResumeSectionType.Custom]: '其他资质'
  };

  return (
    <div className={`w-full bg-white text-slate-800 h-full min-h-[1123px] ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto`} style={{ fontFamily: '"Segoe UI", "Microsoft YaHei", sans-serif' }}>
      {/* 顶部蓝色条 */}
      <div className="h-3 bg-[#005eb8]"></div>
      
      <div className="p-10">
        {/* 头部 */}
        <div className="flex justify-between items-start mb-12 border-b border-gray-200 pb-8">
            <div className="flex items-center gap-6">
                {data.personalInfo.avatarUrl && (
                    <div className="w-28 h-36 bg-gray-100 border border-gray-200 p-1">
                        <img src={data.personalInfo.avatarUrl} className="w-full h-full object-cover" alt="Doctor Avatar" />
                    </div>
                )}
                <div>
                    <h1 className="text-4xl font-bold text-[#005eb8] mb-2">{data.personalInfo.fullName}</h1>
                    <div className="text-xl text-gray-600 font-medium mb-4">{data.personalInfo.jobTitle}</div>
                    
                    <div className="flex flex-wrap gap-x-6 gap-y-2 text-sm text-gray-500">
                        <span className="flex items-center gap-2">
                            <svg className="w-4 h-4 text-[#005eb8]" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"></path></svg>
                            {data.personalInfo.phone}
                        </span>
                        <span className="flex items-center gap-2">
                            <svg className="w-4 h-4 text-[#005eb8]" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"></path></svg>
                            {data.personalInfo.email}
                        </span>
                    </div>
                </div>
            </div>
            
            {/* 十字装饰 */}
            <div className="hidden md:block text-[#e6f0fa]">
                <svg className="w-32 h-32" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
                </svg>
            </div>
        </div>

        <div className="grid grid-cols-12 gap-10">
            {/* 左侧栏 - 技能与证书 */}
            <div className="col-span-4 space-y-8">
                {(data.sections || []).filter(s => (s.type === ResumeSectionType.Skills || s.type === ResumeSectionType.Custom) && s.isVisible).map(section => (
                    <div key={section.id} className="bg-[#f0f7ff] p-6 rounded-lg">
                        <h3 className="text-lg font-bold text-[#003087] mb-4 flex items-center gap-2">
                            <span className="text-2xl">+</span> {titleMap[section.type] || section.title}
                        </h3>
                        
                        {section.type === ResumeSectionType.Skills ? (
                            <ul className="space-y-2">
                                {section.items.map(item => (
                                    <li key={item.id} className="text-sm text-gray-700 flex items-start">
                                        <span className="mr-2 text-[#005eb8]">•</span>
                                        {item.description}
                                    </li>
                                ))}
                            </ul>
                        ) : (
                             <div className="text-sm text-gray-700 space-y-2">
                                {section.items.map(item => (
                                    <div key={item.id}>
                                        <div className="font-bold text-[#005eb8]">{item.title}</div>
                                        <div dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description || '') }} />
                                    </div>
                                ))}
                             </div>
                        )}
                    </div>
                ))}

                {/* 个人信息卡片 */}
                <div className="bg-gray-50 p-6 rounded-lg border border-gray-100">
                    <h3 className="text-sm font-bold text-gray-500 uppercase mb-4">基本信息</h3>
                    <div className="space-y-2 text-sm">
                        {data.personalInfo.gender && <div className="flex justify-between"><span className="text-gray-500">性别</span><span>{data.personalInfo.gender}</span></div>}
                        {data.personalInfo.age && <div className="flex justify-between"><span className="text-gray-500">年龄</span><span>{data.personalInfo.age}</span></div>}
                        {data.personalInfo.politicalStatus && <div className="flex justify-between"><span className="text-gray-500">政治面貌</span><span>{data.personalInfo.politicalStatus}</span></div>}
                    </div>
                </div>
            </div>

            {/* 右侧栏 - 经历与教育 */}
            <div className="col-span-8 space-y-10">
                {(data.sections || []).filter(s => s.type !== ResumeSectionType.Skills && s.type !== ResumeSectionType.Custom && s.isVisible).map(section => (
                    <div key={section.id}>
                        <h3 className="text-xl font-bold text-[#003087] border-b-2 border-[#005eb8] pb-2 mb-6 flex items-center gap-3">
                            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
                            {titleMap[section.type] || section.title}
                        </h3>

                        <div className="space-y-8">
                            {section.items.map(item => (
                                <div key={item.id} className="relative pl-6 border-l-2 border-gray-100">
                                    <div className="absolute -left-[5px] top-2 w-2 h-2 rounded-full bg-[#005eb8]"></div>
                                    <div className="flex justify-between items-baseline mb-1">
                                        <h4 className="text-lg font-bold text-gray-900">{item.title}</h4>
                                        {item.dateRange && <span className="text-sm font-medium text-gray-500">{item.dateRange}</span>}
                                    </div>
                                    {item.subtitle && <div className="text-[#005eb8] font-medium mb-2">{item.subtitle}</div>}
                                    {item.description && (
                                        <div className="text-sm text-gray-700 leading-relaxed text-justify" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
                                    )}
                                </div>
                            ))}
                        </div>
                    </div>
                ))}
            </div>
        </div>
      </div>
      
      {/* 底部心电图装饰 */}
      <div className="absolute bottom-0 left-0 right-0 h-12 overflow-hidden opacity-10 pointer-events-none">
          <svg className="w-full h-full text-[#005eb8]" preserveAspectRatio="none" viewBox="0 0 100 10">
             <path d="M0 5 L10 5 L12 2 L14 8 L16 5 L20 5 L22 0 L24 10 L26 5 L30 5 L100 5" fill="none" stroke="currentColor" vectorEffect="non-scaling-stroke" />
          </svg>
      </div>
    </div>
  );
};

