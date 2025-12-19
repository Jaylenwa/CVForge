import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useLanguage } from '../../../contexts/LanguageContext';
import { sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateCNGame: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const titleMap: Record<string, string> = {
    [ResumeSectionType.Summary]: '角色背景',
    [ResumeSectionType.Education]: '新手村试炼',
    [ResumeSectionType.Experience]: '主线任务',
    [ResumeSectionType.Skills]: '技能树',
    [ResumeSectionType.Projects]: '副本成就',
    [ResumeSectionType.Custom]: '背包物品'
  };

  return (
    <div className={`w-full bg-[#2b2b2b] text-[#dcdcdc] h-full min-h-[1123px] ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto`} style={{ fontFamily: '"Microsoft YaHei", sans-serif' }}>
      <div className="p-2 h-full border-4 border-[#4a4a4a] m-1 rounded-lg">
        <div className="h-full border-2 border-[#808080] p-6 rounded bg-[#333]">
            {/* 角色面板头部 */}
            <div className="flex gap-6 mb-8 bg-[#222] p-4 rounded border border-[#555] shadow-inner">
                <div className="w-32 h-32 bg-[#111] border-2 border-[#666] relative">
                    {data.personalInfo.avatarUrl ? (
                        <img src={data.personalInfo.avatarUrl} className="w-full h-full object-cover" alt="Avatar" />
                    ) : (
                        <div className="w-full h-full flex items-center justify-center text-[#444]">NO IMAGE</div>
                    )}
                    <div className="absolute -bottom-3 -right-3 w-8 h-8 bg-yellow-600 rounded-full border-2 border-[#fff] text-white flex items-center justify-center font-bold text-xs">99</div>
                </div>
                
                <div className="flex-grow">
                    <h1 className="text-3xl font-bold text-yellow-500 mb-1">{data.personalInfo.fullName}</h1>
                    <div className="text-[#aaa] mb-4 text-sm">职业: <span className="text-white">{data.personalInfo.jobTitle}</span></div>
                    
                    <div className="grid grid-cols-2 gap-x-8 gap-y-2 text-xs font-mono">
                        <div className="flex justify-between"><span className="text-[#888]">HP</span> <div className="w-32 h-3 bg-[#400] border border-[#600]"><div className="h-full bg-red-600 w-full"></div></div></div>
                        <div className="flex justify-between"><span className="text-[#888]">MP</span> <div className="w-32 h-3 bg-[#004] border border-[#006]"><div className="h-full bg-blue-600 w-full"></div></div></div>
                        <div className="flex justify-between"><span className="text-[#888]">EXP</span> <div className="w-32 h-3 bg-[#440] border border-[#660]"><div className="h-full bg-yellow-500 w-[80%]"></div></div></div>
                    </div>
                </div>
            </div>

            <div className="grid grid-cols-12 gap-6">
                <div className="col-span-12 space-y-6">
                    {(data.sections || []).filter(s => s.isVisible).map(section => (
                        <div key={section.id} className="bg-[#2a2a2a] border border-[#444] p-4 rounded">
                            <h3 className="text-lg font-bold text-yellow-500 mb-4 border-b border-[#444] pb-2 flex items-center gap-2">
                                <span className="text-[#888]">◆</span>
                                {titleMap[section.type] || section.title}
                            </h3>

                            {section.type === ResumeSectionType.Skills ? (
                                <div className="grid grid-cols-2 gap-4">
                                    {section.items.map(item => (
                                        <div key={item.id} className="flex items-center gap-3 bg-[#222] p-2 rounded border border-[#333]">
                                            <div className="w-8 h-8 bg-[#444] rounded border border-[#666] flex items-center justify-center text-xs">SKILL</div>
                                            <div className="flex-grow">
                                                <div className="text-sm font-bold text-[#ddd]">{item.description}</div>
                                                <div className="h-1 bg-[#333] w-full mt-1"><div className="h-full bg-green-500 w-[90%]"></div></div>
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            ) : (
                                <div className="space-y-4">
                                    {section.items.map(item => (
                                        <div key={item.id} className="bg-[#222] p-3 rounded border border-[#333] hover:border-[#666] transition-colors">
                                            <div className="flex justify-between items-start mb-1">
                                                <h4 className="text-base font-bold text-[#eee]">{item.title}</h4>
                                                {item.dateRange && <span className="text-xs text-[#888]">{item.dateRange}</span>}
                                            </div>
                                            {item.subtitle && <div className="text-sm text-yellow-600 mb-2">{item.subtitle}</div>}
                                            {item.description && (
                                                <div className="text-xs text-[#bbb] leading-relaxed" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
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
    </div>
  );
};

