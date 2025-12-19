import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useLanguage } from '../../../contexts/LanguageContext';
import { sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateCNNewspaper: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const titleMap: Record<string, string> = {
    [ResumeSectionType.Summary]: 'HEADLINES',
    [ResumeSectionType.Education]: 'EDUCATION',
    [ResumeSectionType.Experience]: 'CAREER HISTORY',
    [ResumeSectionType.Skills]: 'EXPERTISE',
    [ResumeSectionType.Projects]: 'FEATURED PROJECTS',
    [ResumeSectionType.Custom]: 'EXTRAS'
  };

  return (
    <div className={`w-full bg-[#fcfae1] text-[#2f2f2f] h-full min-h-[1123px] ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto`} style={{ fontFamily: '"Playfair Display", "Times New Roman", serif' }}>
      <div className="p-8">
        {/* Newspaper Header */}
        <div className="border-b-4 border-double border-black pb-4 mb-8 text-center">
            <div className="flex justify-between border-b border-black pb-1 mb-2 text-xs font-sans font-bold uppercase">
                <span>The Daily Resume</span>
                <span>{new Date().toDateString()}</span>
                <span>Price: Hired</span>
            </div>
            <h1 className="text-6xl font-black uppercase tracking-tighter mb-4">{data.personalInfo.fullName}</h1>
            <div className="border-t border-b border-black py-2 flex justify-center gap-8 text-sm font-bold font-sans uppercase">
                <span>{data.personalInfo.jobTitle}</span>
                <span>•</span>
                <span>{data.personalInfo.city}</span>
                <span>•</span>
                <span>{data.personalInfo.email}</span>
            </div>
        </div>

        <div className="grid grid-cols-12 gap-8">
            {/* Left Column (Main Article) */}
            <div className="col-span-8 space-y-8 pr-6 border-r border-black">
                {(data.sections || []).filter(s => s.type !== ResumeSectionType.Skills && s.isVisible).map(section => (
                    <div key={section.id}>
                        <h3 className="text-2xl font-bold uppercase border-b-2 border-black mb-4">{titleMap[section.type] || section.title}</h3>
                        <div className="space-y-6">
                            {section.items.map(item => (
                                <div key={item.id}>
                                    <h4 className="text-xl font-bold italic mb-1">{item.title}</h4>
                                    <div className="flex justify-between text-xs font-sans font-bold text-gray-600 mb-2 border-b border-gray-400 pb-1">
                                        <span>{item.subtitle || 'Contributor'}</span>
                                        <span>{item.dateRange}</span>
                                    </div>
                                    {item.description && (
                                        <div className="text-sm leading-snug text-justify font-serif" style={{ columnCount: item.description.length > 300 ? 2 : 1, columnGap: '1.5rem' }} dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
                                    )}
                                </div>
                            ))}
                        </div>
                    </div>
                ))}
            </div>

            {/* Right Column (Sidebar/Ads) */}
            <div className="col-span-4 space-y-8">
                {data.personalInfo.avatarUrl && (
                    <div className="border-2 border-black p-2 bg-white transform rotate-1 shadow-lg">
                        <img src={data.personalInfo.avatarUrl} className="w-full object-cover" alt="Profile" />
                        <div className="text-center text-xs font-sans mt-2 font-bold uppercase">Figure 1. The Candidate</div>
                    </div>
                )}

                {(data.sections || []).filter(s => s.type === ResumeSectionType.Skills && s.isVisible).map(section => (
                    <div key={section.id} className="border-2 border-black p-4 bg-white">
                        <h3 className="text-xl font-bold uppercase text-center border-b-2 border-black mb-4">{titleMap[section.type] || section.title}</h3>
                        <ul className="list-disc pl-4 space-y-1 font-sans text-sm">
                            {section.items.map(item => (
                                <li key={item.id}>
                                    <span className="font-bold">{item.description}</span>
                                </li>
                            ))}
                        </ul>
                    </div>
                ))}

                <div className="border border-black p-4 text-center bg-gray-100">
                    <h4 className="font-bold uppercase mb-2">Contact Now</h4>
                    <p className="text-sm font-sans">{data.personalInfo.phone}</p>
                    <p className="text-xs mt-2 italic">Available for immediate start</p>
                </div>
            </div>
        </div>
      </div>
    </div>
  );
};
