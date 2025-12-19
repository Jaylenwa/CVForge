import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useLanguage } from '../../../contexts/LanguageContext';
import { sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateCNGlitch: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const titleMap: Record<string, string> = {
    [ResumeSectionType.Summary]: 'SYSTEM.LOG',
    [ResumeSectionType.Education]: 'DATA.UPLOAD',
    [ResumeSectionType.Experience]: 'RUN.EXEC',
    [ResumeSectionType.Skills]: 'CORE.MODULES',
    [ResumeSectionType.Projects]: 'BETA.BUILDS',
    [ResumeSectionType.Custom]: 'EXTRA.BITS'
  };

  return (
    <div className={`w-full bg-black text-[#00ff00] h-full min-h-[1123px] ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto overflow-hidden`} style={{ fontFamily: '"Courier New", Courier, monospace' }}>
       {/* Glitch Background Effect */}
       <div className="absolute inset-0 opacity-10 pointer-events-none" style={{ 
           backgroundImage: 'linear-gradient(transparent 50%, rgba(0, 255, 0, 0.05) 50%)',
           backgroundSize: '100% 4px'
       }}></div>
       
       <div className="p-10 relative z-10">
          <div className="border-b-2 border-[#00ff00] pb-6 mb-10 relative">
             <div className="absolute top-0 left-0 w-full h-full opacity-50 pointer-events-none" style={{ textShadow: '2px 0 #ff00ff, -2px 0 #00ffff' }}></div>
             <h1 className="text-5xl font-bold mb-2 uppercase tracking-tighter" style={{ textShadow: '2px 2px #ff00ff' }}>{data.personalInfo.fullName}</h1>
             <div className="text-xl text-[#00ff00] uppercase tracking-widest flex items-center">
                 <span className="animate-pulse mr-2">_</span>
                 {data.personalInfo.jobTitle}
             </div>
             
             <div className="mt-4 flex flex-wrap gap-6 text-sm font-bold text-[#00ff00]">
                <span>{'>'} {data.personalInfo.phone}</span>
                <span>{'>'} {data.personalInfo.email}</span>
                {data.personalInfo.city && <span>{'>'} {data.personalInfo.city}</span>}
             </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-12 gap-8">
             <div className="md:col-span-4 space-y-8 border-r-2 border-[#00ff00] border-dashed pr-6">
                {data.personalInfo.avatarUrl && (
                   <div className="mb-8 relative group">
                      <img src={data.personalInfo.avatarUrl} alt="Avatar" className="w-full object-cover transition-all" />
                      <div className="absolute top-0 left-0 w-full h-1 bg-[#ff00ff] opacity-50"></div>
                      <div className="absolute bottom-0 right-0 w-full h-1 bg-[#00ffff] opacity-50"></div>
                   </div>
                )}
                
                {(data.sections || []).filter(s => s.type === ResumeSectionType.Skills && s.isVisible).map(section => (
                   <div key={section.id}>
                      <h3 className="text-xl font-bold mb-4 bg-[#00ff00] text-black px-2 inline-block transform -skew-x-12">
                         {titleMap[section.type] || section.title}
                      </h3>
                      <div className="space-y-2">
                         {section.items.map(item => (
                            <div key={item.id} className="flex items-center text-sm">
                               <span className="text-[#ff00ff] mr-2">[{item.description}]</span>
                            </div>
                         ))}
                      </div>
                   </div>
                ))}
             </div>

             <div className="md:col-span-8 space-y-10">
                {(data.sections || []).filter(s => s.type !== ResumeSectionType.Skills && s.isVisible).map(section => (
                   <div key={section.id}>
                      <h3 className="text-2xl font-bold mb-6 flex items-center border-b border-[#00ff00] pb-1">
                         <span className="mr-2 animate-pulse">#</span>
                         {titleMap[section.type] || section.title}
                      </h3>

                      <div className="space-y-8">
                         {section.items.map(item => (
                            <div key={item.id} className="relative pl-4 border-l border-[#00ff00]">
                               <div className="absolute -left-[5px] top-1.5 w-2 h-2 bg-[#ff00ff]"></div>
                               <div className="flex justify-between items-baseline mb-1">
                                  <h4 className="text-xl font-bold text-white">{item.title}</h4>
                                  {item.dateRange && <span className="text-sm font-mono bg-[#003300] px-2 text-[#00ff00]">{item.dateRange}</span>}
                               </div>
                               {item.subtitle && <div className="text-sm text-[#00ffff] mb-2 uppercase">{item.subtitle}</div>}
                               {item.description && (
                                  <div className="text-sm leading-relaxed text-[#ccffcc] font-mono" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
                                )}
                            </div>
                         ))}
                      </div>
                   </div>
                ))}
             </div>
          </div>
          
          <div className="mt-12 pt-4 border-t border-[#00ff00] text-center text-xs text-[#003300]">
             SYSTEM_ID: {Math.random().toString(36).substr(2, 9).toUpperCase()} // END_OF_FILE
          </div>
       </div>
    </div>
  );
};
