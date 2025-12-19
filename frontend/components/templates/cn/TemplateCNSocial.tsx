import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useLanguage } from '../../../contexts/LanguageContext';
import { sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateCNSocial: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const titleMap: Record<string, string> = {
    [ResumeSectionType.Summary]: 'Intro',
    [ResumeSectionType.Education]: 'Education',
    [ResumeSectionType.Experience]: 'Work',
    [ResumeSectionType.Skills]: 'Skills',
    [ResumeSectionType.Projects]: 'Highlights',
    [ResumeSectionType.Custom]: 'More'
  };

  return (
    <div className={`w-full bg-white text-gray-900 h-full min-h-[1123px] ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto font-sans`} style={{ fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif' }}>
      <div className="max-w-2xl mx-auto p-8">
        {/* Profile Header */}
        <div className="flex gap-8 mb-10 border-b border-gray-200 pb-8">
            <div className="flex-shrink-0">
                <div className="w-32 h-32 rounded-full p-[2px] bg-gradient-to-tr from-yellow-400 via-red-500 to-purple-500">
                    <div className="bg-white p-[2px] rounded-full w-full h-full">
                         {data.personalInfo.avatarUrl ? (
                             <img src={data.personalInfo.avatarUrl} className="w-full h-full rounded-full object-cover" alt="Profile" />
                         ) : (
                             <div className="w-full h-full rounded-full bg-gray-200"></div>
                         )}
                    </div>
                </div>
            </div>
            
            <div className="flex-grow pt-2">
                <div className="flex justify-between items-start mb-4">
                    <h1 className="text-2xl font-semibold">{data.personalInfo.fullName}</h1>
                    <div className="bg-blue-500 text-white px-4 py-1.5 rounded-lg font-semibold text-sm">Follow</div>
                </div>
                
                <div className="grid grid-cols-3 gap-4 mb-4 text-sm">
                    <div className="text-center"><span className="font-bold block">10+</span>Years Exp</div>
                    <div className="text-center"><span className="font-bold block">50+</span>Projects</div>
                    <div className="text-center"><span className="font-bold block">1k+</span>Skills</div>
                </div>
                
                <div className="text-sm">
                    <div className="font-semibold">{data.personalInfo.jobTitle}</div>
                    <div className="text-gray-600">📍 {data.personalInfo.city || 'Remote'}</div>
                    <div className="text-blue-900">🔗 {data.personalInfo.email}</div>
                </div>
            </div>
        </div>

        {/* Highlights/Skills */}
        {(data.sections || []).filter(s => s.type === ResumeSectionType.Skills && s.isVisible).map(section => (
            <div key={section.id} className="flex gap-4 overflow-x-auto pb-6 mb-4 no-scrollbar">
                {section.items.map(item => (
                    <div key={item.id} className="flex-shrink-0 flex flex-col items-center gap-1 w-16">
                        <div className="w-14 h-14 rounded-full border border-gray-300 bg-gray-50 flex items-center justify-center text-xs text-center p-1 overflow-hidden">
                            {item.description.slice(0, 4)}
                        </div>
                        <span className="text-xs truncate w-full text-center">{item.description}</span>
                    </div>
                ))}
            </div>
        ))}

        {/* Tab Navigation */}
        <div className="flex border-t border-gray-200 mb-6">
            <div className="flex-1 text-center border-t-2 border-black py-3 text-xs font-bold uppercase tracking-widest">Feed</div>
            <div className="flex-1 text-center border-t border-transparent py-3 text-xs font-bold uppercase tracking-widest text-gray-400">Reels</div>
            <div className="flex-1 text-center border-t border-transparent py-3 text-xs font-bold uppercase tracking-widest text-gray-400">Tagged</div>
        </div>

        {/* Feed Content */}
        <div className="space-y-8">
            {(data.sections || []).filter(s => s.type !== ResumeSectionType.Skills && s.isVisible).map(section => (
                <div key={section.id}>
                    <div className="flex items-center gap-2 mb-3">
                        {data.personalInfo.avatarUrl && <img src={data.personalInfo.avatarUrl} className="w-8 h-8 rounded-full" />}
                        <span className="font-semibold text-sm">{data.personalInfo.fullName}</span>
                        <span className="text-gray-500 text-xs">• {titleMap[section.type] || section.title}</span>
                    </div>

                    <div className="bg-white border border-gray-200 rounded-lg overflow-hidden mb-2">
                        <div className="bg-gray-50 p-4 border-b border-gray-100">
                             <h3 className="font-bold text-lg">{titleMap[section.type] || section.title}</h3>
                        </div>
                        <div className="p-4 space-y-6">
                            {section.items.map(item => (
                                <div key={item.id}>
                                    <div className="flex justify-between font-semibold text-sm mb-1">
                                        <span>{item.title}</span>
                                        <span className="text-gray-500 font-normal">{item.dateRange}</span>
                                    </div>
                                    {item.subtitle && <div className="text-sm text-gray-600 mb-2">{item.subtitle}</div>}
                                    {item.description && (
                                        <div className="text-sm text-gray-800" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
                                    )}
                                </div>
                            ))}
                        </div>
                    </div>
                    
                    <div className="flex gap-4 text-2xl mb-2">
                        <span>♡</span><span>💬</span><span>➢</span>
                    </div>
                    <div className="text-sm font-semibold mb-6">Liked by Recruiter and others</div>
                </div>
            ))}
        </div>
      </div>
    </div>
  );
};

