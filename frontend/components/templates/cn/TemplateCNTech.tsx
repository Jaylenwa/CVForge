import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useLanguage } from '../../../contexts/LanguageContext';
import { sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateCNTech: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const titleMap: Record<string, string> = {
    [ResumeSectionType.Summary]: 'SYSTEM.SUMMARY',
    [ResumeSectionType.Education]: 'DATA.EDUCATION',
    [ResumeSectionType.Experience]: 'EXEC.EXPERIENCE',
    [ResumeSectionType.Skills]: 'CORE.SKILLS',
    [ResumeSectionType.Projects]: 'PROJECT.LOGS',
    [ResumeSectionType.Custom]: 'EXTRA.INFO'
  };

  return (
    <div className={`w-full bg-slate-900 text-cyan-400 h-full min-h-[1123px] ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto overflow-hidden font-mono`} style={{ fontFamily: '"Space Mono", "Courier New", monospace', lineHeight: '1.6' }}>
      {/* 背景装饰网格 */}
      <div className="absolute inset-0 opacity-10 pointer-events-none" 
           style={{ backgroundImage: 'linear-gradient(rgba(6,182,212,0.1) 1px, transparent 1px), linear-gradient(90deg, rgba(6,182,212,0.1) 1px, transparent 1px)', backgroundSize: '20px 20px' }}>
      </div>

      <div className="p-10 relative z-10">
        {/* 头部 */}
        <div className="border-b-2 border-cyan-500 pb-6 mb-8 flex justify-between items-end">
           <div>
              <h1 className="text-4xl font-bold mb-2 text-white glitch-effect">{data.personalInfo.fullName}</h1>
              <div className="text-xl text-cyan-300">{'>>'} {data.personalInfo.jobTitle}</div>
           </div>
           <div className="text-right text-xs text-cyan-600">
              <div>ID: {Math.random().toString(36).substr(2, 9).toUpperCase()}</div>
              <div>STATUS: ONLINE</div>
           </div>
        </div>

        <div className="grid grid-cols-12 gap-8">
           {/* 左侧栏 */}
           <div className="col-span-4 space-y-8 border-r border-cyan-900 pr-6">
              {data.personalInfo.avatarUrl && (
                <div className="relative mb-6 group">
                   <div className="absolute -inset-1 bg-gradient-to-r from-cyan-400 to-purple-600 rounded-lg blur opacity-25 group-hover:opacity-75 transition duration-1000 group-hover:duration-200"></div>
                   <img src={data.personalInfo.avatarUrl} alt="Avatar" className="relative w-full rounded-lg transition-all duration-500 border border-cyan-500/30" />
                </div>
              )}
              
              <div className="space-y-4 text-sm">
                 <div className="flex items-center space-x-2">
                    <span className="text-purple-400">[TEL]</span>
                    <span className="text-gray-300">{data.personalInfo.phone}</span>
                 </div>
                 <div className="flex items-center space-x-2">
                    <span className="text-purple-400">[MAIL]</span>
                    <span className="text-gray-300 break-all">{data.personalInfo.email}</span>
                 </div>
                 {data.personalInfo.gender && (
                   <div className="flex items-center space-x-2">
                      <span className="text-purple-400">[SEX]</span>
                      <span className="text-gray-300">{data.personalInfo.gender}</span>
                   </div>
                 )}
              </div>

              {(data.sections || []).filter(s => s.type === ResumeSectionType.Skills && s.isVisible).map(section => (
                 <div key={section.id}>
                    <h3 className="text-lg font-bold text-white mb-4 border-l-4 border-purple-500 pl-3">{titleMap[section.type] || section.title}</h3>
                    <div className="space-y-3">
                       {section.items.map(item => (
                          <div key={item.id}>
                             <div className="flex justify-between text-xs mb-1 text-cyan-200">
                                <span>{item.description}</span>
                                <span>LOADED</span>
                             </div>
                             <div className="h-1 w-full bg-slate-800 rounded overflow-hidden">
                                <div className="h-full bg-cyan-600 w-full animate-pulse"></div>
                             </div>
                          </div>
                       ))}
                    </div>
                 </div>
              ))}
           </div>

           {/* 右侧内容 */}
           <div className="col-span-8 space-y-10">
              {(data.sections || []).filter(s => s.type !== ResumeSectionType.Skills && s.isVisible).map(section => (
                 <div key={section.id}>
                    <div className="flex items-center mb-6">
                       <span className="text-purple-500 mr-2">{'{'}</span>
                       <h3 className="text-xl font-bold text-white tracking-wider">{titleMap[section.type] || section.title}</h3>
                       <span className="text-purple-500 ml-2">{'}'}</span>
                       <div className="flex-grow h-px bg-cyan-900/50 ml-4"></div>
                    </div>

                    <div className="space-y-8 pl-4 border-l border-cyan-900/30">
                       {section.items.map(item => (
                          <div key={item.id} className="relative">
                             <div className="absolute -left-[21px] top-1.5 w-2 h-2 bg-slate-900 border border-cyan-500 rotate-45"></div>
                             <div className="mb-2 flex justify-between items-baseline">
                                <h4 className="text-lg font-bold text-cyan-100">{item.title}</h4>
                                {item.dateRange && <span className="text-xs text-cyan-600 border border-cyan-900 px-2 py-0.5 rounded">{item.dateRange}</span>}
                             </div>
                             {item.subtitle && <div className="text-sm text-purple-300 mb-2 font-semibold">@ {item.subtitle}</div>}
                             {item.description && (
                                <div className="text-sm text-gray-400 leading-relaxed" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
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
