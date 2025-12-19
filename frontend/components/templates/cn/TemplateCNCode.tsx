import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useLanguage } from '../../../contexts/LanguageContext';
import { sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateCNCode: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const titleMap: Record<string, string> = {
    [ResumeSectionType.Summary]: 'README.md',
    [ResumeSectionType.Education]: 'education.json',
    [ResumeSectionType.Experience]: 'experience.log',
    [ResumeSectionType.Skills]: 'package.json',
    [ResumeSectionType.Projects]: 'projects.ts',
    [ResumeSectionType.Custom]: 'config.yml'
  };

  return (
    <div className={`w-full bg-[#1e1e1e] text-[#d4d4d4] h-full min-h-[1123px] ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto font-mono text-sm`} style={{ fontFamily: '"Menlo", "Consolas", "Courier New", monospace' }}>
      <div className="flex h-full">
        {/* Activity Bar */}
        <div className="w-12 bg-[#333333] flex flex-col items-center py-4 space-y-6 text-[#858585]">
            <div className="w-6 h-6 border-2 border-[#858585]"></div>
            <div className="w-6 h-6 border-2 border-[#858585] rounded-full"></div>
            <div className="w-6 h-6 border-2 border-[#858585] transform rotate-45"></div>
        </div>

        {/* Explorer */}
        <div className="w-48 bg-[#252526] text-[#cccccc] p-4 hidden md:block">
            <div className="text-xs font-bold uppercase mb-4 pl-2">Explorer</div>
            <div className="space-y-1">
                <div className="flex items-center gap-1 font-bold text-blue-400">
                    <span>▼</span>
                    <span>{data.personalInfo.fullName?.replace(/\s+/g, '_') || 'USER'}</span>
                </div>
                {(data.sections || []).filter(s => s.isVisible).map(section => (
                    <div key={section.id} className="pl-5 flex items-center gap-2 hover:bg-[#37373d] cursor-pointer py-1">
                        <span className="text-yellow-400">JS</span>
                        <span className="truncate">{titleMap[section.type] || section.title}</span>
                    </div>
                ))}
            </div>
        </div>

        {/* Editor Area */}
        <div className="flex-grow bg-[#1e1e1e] p-6 overflow-hidden">
            {/* Tabs */}
            <div className="flex bg-[#2d2d2d] mb-4 text-sm">
                <div className="px-4 py-2 bg-[#1e1e1e] border-t-2 border-blue-400 text-white flex items-center gap-2">
                    <span className="text-yellow-400">TS</span>
                    resume.tsx
                    <span className="ml-2 hover:bg-[#444] rounded-full w-4 h-4 flex items-center justify-center">×</span>
                </div>
            </div>
            {data.personalInfo.avatarUrl && (
              <div className="flex justify-end mb-3">
                <img
                  src={data.personalInfo.avatarUrl}
                  alt={t('a11y.avatarAlt')}
                  className="w-16 h-16 rounded-md object-cover border border-[#37373d]"
                />
              </div>
            )}

            <div className="space-y-1">
                {/* Header Code */}
                <div className="flex">
                    <div className="w-8 text-right text-[#858585] mr-4 select-none">1</div>
                    <div className="text-[#569cd6]">const <span className="text-[#4fc1ff]">candidate</span> = <span className="text-[#dcdcaa]">{'{'}</span></div>
                </div>
                <div className="flex">
                    <div className="w-8 text-right text-[#858585] mr-4 select-none">2</div>
                    <div className="pl-4"><span className="text-[#9cdcfe]">name</span>: <span className="text-[#ce9178]">"{data.personalInfo.fullName}"</span>,</div>
                </div>
                <div className="flex">
                    <div className="w-8 text-right text-[#858585] mr-4 select-none">3</div>
                    <div className="pl-4"><span className="text-[#9cdcfe]">role</span>: <span className="text-[#ce9178]">"{data.personalInfo.jobTitle}"</span>,</div>
                </div>
                <div className="flex">
                    <div className="w-8 text-right text-[#858585] mr-4 select-none">4</div>
                    <div className="pl-4"><span className="text-[#9cdcfe]">contact</span>: <span className="text-[#dcdcaa]">{'{'}</span></div>
                </div>
                <div className="flex">
                    <div className="w-8 text-right text-[#858585] mr-4 select-none">5</div>
                    <div className="pl-8"><span className="text-[#9cdcfe]">email</span>: <span className="text-[#ce9178]">"{data.personalInfo.email}"</span>,</div>
                </div>
                <div className="flex">
                    <div className="w-8 text-right text-[#858585] mr-4 select-none">6</div>
                    <div className="pl-8"><span className="text-[#9cdcfe]">phone</span>: <span className="text-[#ce9178]">"{data.personalInfo.phone}"</span></div>
                </div>
                <div className="flex">
                    <div className="w-8 text-right text-[#858585] mr-4 select-none">7</div>
                    <div className="pl-4"><span className="text-[#dcdcaa]">{'}'}</span>,</div>
                </div>

                {/* Sections */}
                {(data.sections || []).filter(s => s.isVisible).map((section, index) => (
                    <div key={section.id}>
                        <div className="flex mt-4">
                            <div className="w-8 text-right text-[#858585] mr-4 select-none">{8 + index * 10}</div>
                            <div className="text-[#6a9955]">// {titleMap[section.type] || section.title}</div>
                        </div>
                        
                        {section.type === ResumeSectionType.Skills ? (
                            <div className="flex">
                                <div className="w-8 text-right text-[#858585] mr-4 select-none">{9 + index * 10}</div>
                                <div className="pl-0">
                                    <span className="text-[#c586c0]">export const</span> <span className="text-[#dcdcaa]">skills</span> = [
                                    {section.items.map(item => <span key={item.id} className="text-[#ce9178]">"{item.description}"</span>).reduce((prev, curr) => [prev, ', ', curr])}
                                    ];
                                </div>
                            </div>
                        ) : (
                            <div className="space-y-1">
                                {section.items.map((item, i) => (
                                    <div key={item.id}>
                                        <div className="flex">
                                            <div className="w-8 text-right text-[#858585] mr-4 select-none">{9 + index * 10 + i}</div>
                                            <div className="pl-0">
                                                <span className="text-[#dcdcaa]">{item.title?.replace(/\s+/g, '') || `entry${i + 1}`}</span>() <span className="text-[#dcdcaa]">{'{'}</span>
                                            </div>
                                        </div>
                                        <div className="flex">
                                             <div className="w-8 text-right text-[#858585] mr-4 select-none"></div>
                                             <div className="pl-4 text-[#6a9955]">/* {item.dateRange} | {item.subtitle} */</div>
                                        </div>
                                        <div className="flex">
                                            <div className="w-8 text-right text-[#858585] mr-4 select-none"></div>
                                            <div className="pl-4 text-[#ce9178] opacity-80" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description || '') }} />
                                        </div>
                                        <div className="flex">
                                            <div className="w-8 text-right text-[#858585] mr-4 select-none"></div>
                                            <div className="pl-0 text-[#dcdcaa]">{'}'}</div>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>
                ))}

                <div className="flex mt-4">
                     <div className="w-8 text-right text-[#858585] mr-4 select-none">99</div>
                     <div className="text-[#dcdcaa]">{'}'}</div>
                </div>
            </div>
        </div>
      </div>
    </div>
  );
};
