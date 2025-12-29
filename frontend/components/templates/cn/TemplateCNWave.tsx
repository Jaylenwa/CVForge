import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useLanguage } from '../../../contexts/LanguageContext';
import { sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateCNWave: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const titleMap: Record<string, string> = {
    [ResumeSectionType.Summary]: '自我评价',
    [ResumeSectionType.Education]: '教育经历',
    [ResumeSectionType.Experience]: '工作经历',
    [ResumeSectionType.Skills]: '技能特长',
    [ResumeSectionType.Projects]: '项目经验',
    [ResumeSectionType.Custom]: '其他信息'
  };

  return (
    <div className={`w-full bg-white text-slate-800 h-full min-h-[1123px] ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto`} style={{ fontFamily: '"PingFang SC", "Microsoft YaHei", sans-serif' }}>
      {/* 顶部波浪背景 */}
      <div className="bg-gradient-to-r from-teal-400 to-blue-500 text-white pb-20 pt-12 px-10 relative">
          <div className="relative z-10 flex justify-between items-center">
              <div>
                  <h1 className="text-5xl font-bold mb-2 tracking-wide">{data.personalInfo.fullName}</h1>
                  <p className="text-xl opacity-90">{data.personalInfo.jobTitle}</p>
              </div>
              <div className="text-right text-sm opacity-90">
                  {data.personalInfo.avatarUrl && (
                    <img
                      src={data.personalInfo.avatarUrl}
                      alt={t('a11y.avatarAlt')}
                      className="w-20 h-20 object-cover rounded-md border border-white/50 ml-auto mb-2"
                    />
                  )}
                  <div>{data.personalInfo.phone}</div>
                  <div>{data.personalInfo.email}</div>
                  <div>{data.personalInfo.city}</div>
              </div>
          </div>
          
          <div className="absolute bottom-0 left-0 right-0">
              <svg viewBox="0 0 1440 320" className="w-full h-24 text-white fill-current block">
                  <path fillOpacity="1" d="M0,160L48,170.7C96,181,192,203,288,202.7C384,203,480,181,576,165.3C672,149,768,139,864,154.7C960,171,1056,213,1152,218.7C1248,224,1344,192,1392,176L1440,160L1440,320L1392,320C1344,320,1248,320,1152,320C1056,320,960,320,864,320C768,320,672,320,576,320C480,320,384,320,288,320C192,320,96,320,48,320L0,320Z"></path>
              </svg>
          </div>
      </div>

      <div className="px-10 py-6">
        <div className="space-y-10">
            {(data.sections || []).filter(s => s.isVisible).map((section, index) => (
                <div key={section.id} className="relative">
                    <div className="flex items-center mb-6">
                        <div className="bg-teal-500 text-white rounded-full p-2 mr-3 shadow-lg transform -rotate-12">
                            <span className="font-bold text-lg px-2">{index + 1}</span>
                        </div>
                        <h3 className="text-2xl font-bold text-teal-600">{titleMap[section.type] || section.title}</h3>
                        <div className="flex-grow h-px bg-teal-100 ml-4"></div>
                    </div>

                    <div className="pl-4">
                        {section.type === ResumeSectionType.Skills ? (
                            <div className="flex flex-wrap gap-3">
                                {section.items.map(item => (
                                    <div key={item.id} className="bg-teal-50 text-teal-700 px-4 py-2 rounded-full border border-teal-200 text-sm font-medium">
                                        {item.description}
                                    </div>
                                ))}
                            </div>
                        ) : (
                            <div className="space-y-8">
                                {section.items.map(item => (
                                    <div key={item.id} className="group">
                                        <div className="flex justify-between items-baseline mb-1 border-b border-gray-100 pb-2 group-hover:border-teal-200 transition-colors">
                                            <h4 className="text-lg font-bold text-gray-800">{item.title}</h4>
                                            {item.dateRange && <span className="text-sm text-gray-500">{item.dateRange}</span>}
                                        </div>
                                        {item.subtitle && <div className="text-base text-teal-500 mb-2 font-medium">{item.subtitle}</div>}
                                        {item.description && (
                                            <div className="text-sm text-gray-600 leading-relaxed" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
                                        )}
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>
                </div>
            ))}
        </div>
      </div>
    </div>
  );
};
