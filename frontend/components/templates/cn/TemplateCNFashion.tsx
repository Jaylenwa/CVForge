import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useLanguage } from '../../../contexts/LanguageContext';
import { sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateCNFashion: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const titleMap: Record<string, string> = {
    [ResumeSectionType.Summary]: 'PROFILE',
    [ResumeSectionType.Education]: 'EDUCATION',
    [ResumeSectionType.Experience]: 'EXPERIENCE',
    [ResumeSectionType.Skills]: 'SKILLS',
    [ResumeSectionType.Projects]: 'PROJECTS',
    [ResumeSectionType.Custom]: 'INFO'
  };

  return (
    <div className={`w-full bg-[#f8f5f2] text-black h-full min-h-[1123px] ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto`} style={{ fontFamily: '"Didot", "Bodoni MT", serif' }}>
      <div className="p-16 h-full border-[20px] border-white">
        <div className="flex justify-between items-end mb-24 relative">
             <div className="absolute -top-10 left-0 text-[10px] tracking-[0.5em] uppercase text-gray-400">Portfolio / {new Date().getFullYear()}</div>
             
             <div>
                 <h1 className="text-6xl font-normal uppercase tracking-widest leading-none mb-4">{data.personalInfo.fullName}</h1>
                 <p className="text-xl italic text-gray-500 font-serif">{data.personalInfo.jobTitle}</p>
             </div>

             <div className="text-right text-xs font-sans tracking-widest uppercase text-gray-400 space-y-1">
                 {data.personalInfo.avatarUrl && (
                   <img
                     src={data.personalInfo.avatarUrl}
                     alt={t('a11y.avatarAlt')}
                     className="w-20 h-20 object-cover rounded-md border border-gray-300 ml-auto mb-2"
                   />
                 )}
                 <div>{data.personalInfo.phone}</div>
                 <div>{data.personalInfo.email}</div>
                 <div>{data.personalInfo.city}</div>
             </div>
        </div>

        <div className="grid grid-cols-12 gap-12">
            <div className="col-span-12 space-y-16">
                {(data.sections || []).filter(s => s.isVisible).map(section => (
                    <div key={section.id} className="grid grid-cols-12 gap-8">
                        <div className="col-span-3">
                            <h3 className="text-sm font-bold uppercase tracking-[0.3em] border-t border-black pt-2">{titleMap[section.type] || section.title}</h3>
                        </div>
                        <div className="col-span-9 space-y-10">
                            {section.type === ResumeSectionType.Skills ? (
                                <div className="flex flex-wrap gap-x-8 gap-y-2 text-sm font-sans uppercase tracking-widest text-gray-600">
                                    {section.items.map(item => (
                                        <span key={item.id}>{item.description}</span>
                                    ))}
                                </div>
                            ) : (
                                section.items.map(item => (
                                    <div key={item.id}>
                                        <div className="flex justify-between items-baseline mb-2">
                                            <h4 className="text-xl font-normal uppercase tracking-wider">{item.title}</h4>
                                            {item.dateRange && <span className="text-xs font-sans text-gray-400">{item.dateRange}</span>}
                                        </div>
                                        {item.subtitle && <div className="text-sm italic text-gray-500 mb-4 font-serif">{item.subtitle}</div>}
                                        {item.description && (
                                            <div className="text-sm font-sans text-gray-600 leading-loose text-justify" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
                                        )}
                                    </div>
                                ))
                            )}
                        </div>
                    </div>
                ))}
            </div>
        </div>
      </div>
    </div>
  );
};
