import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useLanguage } from '../../../contexts/LanguageContext';
import { sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateCNProduct: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const titleMap: Record<string, string> = {
    [ResumeSectionType.Summary]: '用户画像 (Summary)',
    [ResumeSectionType.Education]: '知识库 (Education)',
    [ResumeSectionType.Experience]: '迭代历程 (Experience)',
    [ResumeSectionType.Skills]: '技能栈 (Tech Stack)',
    [ResumeSectionType.Projects]: '核心功能 (Projects)',
    [ResumeSectionType.Custom]: '需求池 (Backlog)'
  };

  return (
    <div className={`w-full bg-white text-gray-800 h-full min-h-[1123px] ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto`} style={{ fontFamily: '"Comic Sans MS", "Chalkboard SE", sans-serif' }}>
      {/* 背景网格 */}
      <div className="absolute inset-0 opacity-10 pointer-events-none" 
           style={{ backgroundImage: 'linear-gradient(#ccc 1px, transparent 1px), linear-gradient(90deg, #ccc 1px, transparent 1px)', backgroundSize: '20px 20px' }}>
      </div>

      <div className="p-10 relative z-10">
        {/* 头部 - 原型图风格 */}
        <div className="border-4 border-dashed border-gray-800 p-6 mb-8 transform -rotate-1 bg-white shadow-[4px_4px_0px_0px_rgba(0,0,0,0.2)]">
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-3xl font-bold mb-2">{data.personalInfo.fullName}</h1>
                    <div className="text-xl bg-yellow-200 inline-block px-2 transform rotate-1">{data.personalInfo.jobTitle}</div>
                </div>
                <div className="text-right text-sm font-mono bg-gray-100 p-2 rounded border border-gray-300">
                    <div>PM_ID: {Math.floor(Math.random() * 10000)}</div>
                    <div>{data.personalInfo.phone}</div>
                    <div>{data.personalInfo.email}</div>
                </div>
            </div>
        </div>

        <div className="grid grid-cols-12 gap-8">
            {/* 左侧 - 侧边栏组件 */}
            <div className="col-span-4 space-y-8">
                <div className="bg-gray-50 border-2 border-gray-800 p-4 rounded-xl">
                    {data.personalInfo.avatarUrl && (
                        <div className="mb-4 text-center">
                            <div className="w-24 h-24 mx-auto bg-gray-200 rounded-full border-2 border-dashed border-gray-400 flex items-center justify-center overflow-hidden">
                                <img src={data.personalInfo.avatarUrl} alt="Avatar" className="w-full h-full object-cover" />
                            </div>
                        </div>
                    )}
                    
                    <div className="space-y-2 text-sm font-mono">
                         <div className="flex items-center">
                             <div className="w-3 h-3 rounded-full bg-green-400 mr-2 border border-black"></div>
                             <span>Status: Available</span>
                         </div>
                         <div className="flex items-center">
                             <div className="w-3 h-3 rounded-full bg-blue-400 mr-2 border border-black"></div>
                             <span>Loc: {data.personalInfo.city || 'Remote'}</span>
                         </div>
                    </div>
                </div>

                {(data.sections || []).filter(s => s.type === ResumeSectionType.Skills && s.isVisible).map(section => (
                    <div key={section.id} className="bg-white border-2 border-gray-800 p-4 shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]">
                        <h3 className="text-lg font-bold border-b-2 border-dashed border-gray-400 pb-2 mb-3">{titleMap[section.type] || section.title}</h3>
                        <div className="flex flex-wrap gap-2">
                            {section.items.map(item => (
                                <span key={item.id} className="bg-blue-100 px-2 py-1 text-sm border border-blue-300 rounded-lg transform hover:scale-105 transition-transform cursor-default">
                                    {item.description}
                                </span>
                            ))}
                        </div>
                    </div>
                ))}
            </div>

            {/* 右侧 - 主要内容 */}
            <div className="col-span-8 space-y-8">
                {(data.sections || []).filter(s => s.type !== ResumeSectionType.Skills && s.isVisible).map(section => (
                    <div key={section.id} className="bg-white p-6 border-2 border-gray-200 rounded-lg hover:border-gray-400 transition-colors">
                        <h3 className="text-xl font-bold mb-4 flex items-center">
                            <span className="w-8 h-8 bg-gray-800 text-white rounded flex items-center justify-center mr-3 text-sm">#</span>
                            {titleMap[section.type] || section.title}
                        </h3>

                        <div className="space-y-6">
                            {section.items.map(item => (
                                <div key={item.id} className="relative pl-6 border-l-2 border-dashed border-gray-300">
                                    <div className="absolute -left-[5px] top-1.5 w-2 h-2 rounded-full bg-gray-400"></div>
                                    <div className="flex justify-between items-start mb-1">
                                        <h4 className="text-lg font-bold">{item.title}</h4>
                                        {item.dateRange && <span className="text-xs bg-gray-100 px-2 py-1 rounded border border-gray-200">{item.dateRange}</span>}
                                    </div>
                                    {item.subtitle && <div className="text-sm text-gray-600 mb-2">{item.subtitle}</div>}
                                    {item.description && (
                                        <div className="text-sm leading-relaxed" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
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

