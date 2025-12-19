import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useLanguage } from '../../../contexts/LanguageContext';
import { sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateCNGeometric: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const titleMap: Record<string, string> = {
    [ResumeSectionType.Summary]: '概况',
    [ResumeSectionType.Education]: '教育',
    [ResumeSectionType.Experience]: '经历',
    [ResumeSectionType.Skills]: '技能',
    [ResumeSectionType.Projects]: '项目',
    [ResumeSectionType.Custom]: '其他'
  };

  return (
    <div className={`w-full bg-white h-full min-h-[1123px] ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto border-4 border-black`} style={{ fontFamily: '"Helvetica Neue", Arial, sans-serif' }}>
      {/* 蒙德里安风格网格布局 */}
      <div className="grid grid-cols-12 h-full border-b-4 border-black">
        {/* 头部区域 */}
        <div className="col-span-8 border-r-4 border-black p-8 bg-white flex flex-col justify-center">
             <h1 className="text-5xl font-black uppercase tracking-tighter mb-2">{data.personalInfo.fullName}</h1>
             <p className="text-xl font-bold bg-yellow-400 inline-block px-2 border-2 border-black transform -rotate-1">{data.personalInfo.jobTitle}</p>
        </div>
        <div className="col-span-4 bg-red-600 p-8 flex flex-col justify-center text-white border-black">
             {data.personalInfo.avatarUrl && (
                 <img src={data.personalInfo.avatarUrl} className="w-24 h-24 border-4 border-black bg-white object-cover mb-4 filter grayscale" alt="Avatar" />
             )}
             <div className="text-xs font-bold space-y-1">
                 <div className="bg-black px-2 py-1 inline-block">{data.personalInfo.phone}</div>
                 <div className="bg-black px-2 py-1 inline-block">{data.personalInfo.email}</div>
             </div>
        </div>
      </div>

      <div className="grid grid-cols-12 h-full flex-grow min-h-[900px]">
        {/* 左侧栏 - 蓝色块 */}
        <div className="col-span-4 border-r-4 border-black flex flex-col">
            <div className="bg-blue-600 h-32 border-b-4 border-black"></div>
            <div className="p-6 space-y-8 flex-grow bg-gray-50">
                {(data.sections || []).filter(s => (s.type === ResumeSectionType.Skills || s.type === ResumeSectionType.Summary) && s.isVisible).map(section => (
                    <div key={section.id}>
                        <h3 className="text-2xl font-black border-b-4 border-black mb-4 inline-block bg-white px-2 shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]">
                            {titleMap[section.type] || section.title}
                        </h3>
                        {section.type === ResumeSectionType.Skills ? (
                            <div className="flex flex-wrap gap-2">
                                {section.items.map(item => (
                                    <span key={item.id} className="border-2 border-black px-2 py-1 text-sm font-bold bg-white">{item.description}</span>
                                ))}
                            </div>
                        ) : (
                            <div className="text-sm font-medium leading-relaxed border-l-4 border-black pl-3">
                                {section.items.map(item => (
                                    <div key={item.id} dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description || '') }} />
                                ))}
                            </div>
                        )}
                    </div>
                ))}
                
                {/* 额外信息 */}
                <div className="bg-yellow-400 border-4 border-black p-4 shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]">
                    <div className="font-bold text-sm">个人信息</div>
                    <div className="text-xs mt-2 space-y-1 font-mono">
                        {data.personalInfo.gender && <div>性别: {data.personalInfo.gender}</div>}
                        {data.personalInfo.age && <div>年龄: {data.personalInfo.age}</div>}
                        {data.personalInfo.city && <div>城市: {data.personalInfo.city}</div>}
                    </div>
                </div>
            </div>
        </div>

        {/* 右侧主内容 */}
        <div className="col-span-8 p-8 bg-white">
            {(data.sections || []).filter(s => s.type !== ResumeSectionType.Skills && s.type !== ResumeSectionType.Summary && s.isVisible).map(section => (
                <div key={section.id} className="mb-10">
                    <div className="flex items-center mb-6">
                        <div className="w-4 h-4 bg-red-600 border-2 border-black mr-3"></div>
                        <h3 className="text-3xl font-black uppercase tracking-wide">{titleMap[section.type] || section.title}</h3>
                        <div className="flex-grow h-1 bg-black ml-4"></div>
                    </div>

                    <div className="space-y-6">
                        {section.items.map(item => (
                            <div key={item.id} className="group">
                                <div className="flex justify-between items-baseline border-b-2 border-gray-200 pb-2 mb-2 group-hover:border-black transition-colors">
                                    <h4 className="text-xl font-bold">{item.title}</h4>
                                    {item.dateRange && <span className="font-mono text-sm bg-gray-100 px-2 border border-black">{item.dateRange}</span>}
                                </div>
                                {item.subtitle && <div className="font-bold text-gray-600 mb-2">{item.subtitle}</div>}
                                {item.description && (
                                    <div className="text-sm leading-relaxed text-justify" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
                                )}
                            </div>
                        ))}
                    </div>
                </div>
            ))}
        </div>
      </div>
      
      {/* 底部装饰 */}
      <div className="h-4 bg-black w-full flex">
          <div className="w-1/3 bg-blue-600 h-full"></div>
          <div className="w-1/3 bg-yellow-400 h-full"></div>
          <div className="w-1/3 bg-red-600 h-full"></div>
      </div>
    </div>
  );
};

