import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useLanguage } from '../../../contexts/LanguageContext';
import { sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateCNPixel: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const titleMap: Record<string, string> = {
    [ResumeSectionType.Summary]: 'PLAYER BIO',
    [ResumeSectionType.Education]: 'TUTORIALS',
    [ResumeSectionType.Experience]: 'QUESTS',
    [ResumeSectionType.Skills]: 'ABILITIES',
    [ResumeSectionType.Projects]: 'ACHIEVEMENTS',
    [ResumeSectionType.Custom]: 'INVENTORY'
  };

  return (
    <div className={`w-full bg-[#e0d5c1] text-[#4a412a] h-full min-h-[1123px] ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto font-mono`} style={{ fontFamily: '"Courier New", monospace', imageRendering: 'pixelated' }}>
      <div className="p-8 border-8 border-double border-[#8b7355] h-full m-4">
        {/* Header */}
        <div className="text-center mb-12 border-b-4 border-[#8b7355] pb-8 bg-[#fdf5e6] p-4 shadow-[4px_4px_0px_0px_#8b7355]">
            <h1 className="text-4xl font-bold mb-2 tracking-widest uppercase">{data.personalInfo.fullName}</h1>
            <div className="text-xl mb-4 text-[#8b4513] font-bold">Lvl.99 {data.personalInfo.jobTitle}</div>
            
            <div className="flex justify-center gap-4 text-xs font-bold uppercase">
                <span className="bg-[#8b7355] text-[#fdf5e6] px-2 py-1">HP: 100/100</span>
                <span className="bg-[#8b7355] text-[#fdf5e6] px-2 py-1">MP: 100/100</span>
            </div>
            
            <div className="mt-4 text-sm flex justify-center gap-6 font-bold">
                <span>☎ {data.personalInfo.phone}</span>
                <span>✉ {data.personalInfo.email}</span>
            </div>
            
            {data.personalInfo.avatarUrl && (
              <img
                src={data.personalInfo.avatarUrl}
                alt={t('a11y.avatarAlt')}
                className="w-24 h-24 object-cover mx-auto rounded border-2 border-[#8b7355] mt-4"
              />
            )}
        </div>

        <div className="grid grid-cols-12 gap-8">
            <div className="col-span-12 space-y-10">
                {(data.sections || []).filter(s => s.isVisible).map(section => (
                    <div key={section.id} className="relative">
                        <div className="bg-[#8b7355] text-[#fdf5e6] px-4 py-2 inline-block mb-4 border-2 border-[#4a412a] shadow-[4px_4px_0px_0px_#4a412a]">
                            <h3 className="text-xl font-bold uppercase">{titleMap[section.type] || section.title}</h3>
                        </div>

                        {section.type === ResumeSectionType.Skills ? (
                            <div className="bg-[#fdf5e6] border-2 border-[#8b7355] p-4 grid grid-cols-2 md:grid-cols-3 gap-4">
                                {section.items.map(item => (
                                    <div key={item.id} className="flex items-center gap-2">
                                        <span className="text-xl">⚔</span>
                                        <span className="font-bold">{item.description}</span>
                                    </div>
                                ))}
                            </div>
                        ) : (
                            <div className="space-y-6 bg-[#fdf5e6] border-2 border-[#8b7355] p-6">
                                {section.items.map(item => (
                                    <div key={item.id} className="border-b-2 border-dashed border-[#8b7355] pb-4 last:border-0 last:pb-0">
                                        <div className="flex justify-between items-start mb-2">
                                            <h4 className="text-lg font-bold uppercase">{item.title}</h4>
                                            {item.dateRange && <span className="text-xs bg-[#4a412a] text-[#e0d5c1] px-2 py-1">{item.dateRange}</span>}
                                        </div>
                                        {item.subtitle && <div className="text-sm font-bold mb-2 text-[#8b4513]">&gt; {item.subtitle}</div>}
                                        {item.description && (
                                            <div className="text-sm leading-relaxed" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
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
