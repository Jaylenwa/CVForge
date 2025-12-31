import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useSectionTitle } from '../../../hooks/useSectionTitle';
import { useLanguage } from '../../../contexts/LanguageContext';
import { hasExtraPersonalInfo, sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateCNCreative: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const titleMap: Record<string, string> = {
    [ResumeSectionType.Summary]: '关于我',
    [ResumeSectionType.Education]: '教育轨迹',
    [ResumeSectionType.Experience]: '职业旅程',
    [ResumeSectionType.Skills]: '技能树',
    [ResumeSectionType.Projects]: '作品集',
    [ResumeSectionType.Custom]: '更多信息'
  };

  // 辅助函数：生成波浪线 SVG
  const WaveLine = () => (
    <svg className="w-24 h-2 text-pink-400 mb-6" viewBox="0 0 100 10" preserveAspectRatio="none">
       <path d="M0 5 Q 25 10 50 5 T 100 5" fill="none" stroke="currentColor" strokeWidth="3" />
    </svg>
  );

  return (
    <div className={`w-full bg-white h-auto ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none`} style={{ fontFamily: styles.fontFamily, lineHeight: parseFloat(styles.spacingMultiplier) * 1.5 }}>
      {/* 顶部不规则背景 */}
      <div className="relative pt-12 pb-16 px-10 bg-gradient-to-r from-pink-50 to-purple-50">
          <div className="absolute top-0 right-0 w-64 h-64 bg-purple-200 rounded-full mix-blend-multiply filter blur-2xl opacity-30 animate-blob"></div>
          <div className="absolute top-0 left-0 w-64 h-64 bg-pink-200 rounded-full mix-blend-multiply filter blur-2xl opacity-30 animate-blob animation-delay-2000"></div>
          
          <div className="relative z-10 flex flex-col items-center">
              {data.personalInfo.avatarUrl && (
                <div className="mb-6 relative">
                    <div className="absolute inset-0 bg-pink-300 rounded-full blur opacity-50 transform translate-x-1 translate-y-1"></div>
                    <img 
                        src={data.personalInfo.avatarUrl} 
                        alt={t('a11y.avatarAlt')} 
                        className="w-32 h-32 rounded-full object-cover border-4 border-white relative z-10"
                    />
                </div>
              )}
              <h1 className="text-4xl font-black text-transparent bg-clip-text bg-gradient-to-r from-purple-600 to-pink-600 mb-2">
                  {data.personalInfo.fullName}
              </h1>
              <p className="text-lg text-gray-600 font-medium mb-6">{data.personalInfo.jobTitle}</p>
              
              <div className="flex flex-wrap justify-center gap-4 text-sm bg-white/60 backdrop-blur-sm py-2 px-6 rounded-full shadow-sm">
                  <span className="text-purple-700">{data.personalInfo.phone}</span>
                  <span className="text-gray-300">•</span>
                  <span className="text-purple-700">{data.personalInfo.email}</span>
                  {data.personalInfo.age && (
                      <>
                        <span className="text-gray-300">•</span>
                        <span className="text-purple-700">{data.personalInfo.age}</span>
                      </>
                  )}
              </div>
          </div>
          
          {/* 底部波浪边缘 */}
          <div className="absolute bottom-0 left-0 right-0 h-12 overflow-hidden">
             <svg viewBox="0 0 1200 120" preserveAspectRatio="none" className="absolute bottom-0 w-full h-full text-white fill-current">
                 <path d="M321.39,56.44c58-10.79,114.16-30.13,172-41.86,82.39-16.72,168.19-17.73,250.45-.39C823.78,31,906.67,72,985.66,92.83c70.05,18.48,146.53,26.09,214.34,3V0H0V27.35A600.21,600.21,0,0,0,321.39,56.44Z"></path>
             </svg>
          </div>
      </div>

      <div className="px-10 py-6">
        <div className="grid grid-cols-1 md:grid-cols-12 gap-10">
            {/* 左侧栏 - 技能与简介 */}
            <div className="md:col-span-4 space-y-8">
                 {(data.sections || []).filter(s => (s.type === ResumeSectionType.Skills || s.type === ResumeSectionType.Summary) && s.isVisible).map(section => (
                     <div key={section.id} className="relative group">
                         {/* 背景装饰Blob */}
                         <div className="absolute -inset-2 bg-pink-50 rounded-2xl transform rotate-1 group-hover:rotate-2 transition-transform duration-300 opacity-50"></div>
                         <div className="relative p-4">
                             <h3 className="text-xl font-bold text-gray-800 mb-2">{titleMap[section.type] || section.title}</h3>
                             <WaveLine />
                             
                             {section.type === ResumeSectionType.Skills ? (
                                 <div className="flex flex-wrap gap-2">
                                     {section.items.map(item => (
                                         <span key={item.id} className="px-3 py-1 bg-white border border-pink-100 rounded-full text-xs font-medium text-pink-600 shadow-sm">
                                             {item.description}
                                         </span>
                                     ))}
                                 </div>
                             ) : (
                                 <div className="text-sm text-gray-600 leading-relaxed">
                                     {section.items.map(item => (
                                         <div key={item.id} dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description || '') }} />
                                     ))}
                                 </div>
                             )}
                         </div>
                     </div>
                 ))}
                 
                 {/* 其他自定义信息 */}
                 {(data.sections || []).filter(s => s.type === ResumeSectionType.Custom && s.isVisible).map(section => (
                     <div key={section.id} className="mt-8">
                        <h3 className="text-lg font-bold text-gray-800 mb-2">{section.title}</h3>
                        <div className="w-10 h-1 bg-purple-400 rounded mb-4"></div>
                        {section.items.map(item => (
                            <div key={item.id} className="mb-3">
                                <div className="font-medium text-gray-700">{item.title}</div>
                                <div className="text-xs text-gray-500" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description || '') }} />
                            </div>
                        ))}
                     </div>
                 ))}
            </div>

            {/* 右侧栏 - 经历 */}
            <div className="md:col-span-8 space-y-10">
                {(data.sections || []).filter(s => ![ResumeSectionType.Skills, ResumeSectionType.Summary, ResumeSectionType.Custom].includes(s.type as any) && s.isVisible).map(section => (
                    <div key={section.id}>
                        <div className="flex items-center mb-6">
                            <span className="text-3xl mr-3 opacity-20">❝</span>
                            <h3 className="text-2xl font-bold text-gray-800">{titleMap[section.type] || section.title}</h3>
                        </div>
                        
                        <div className="space-y-8 border-l-2 border-purple-100 pl-8 ml-3">
                            {section.items.map(item => (
                                <div key={item.id} className="relative">
                                    <div className="absolute -left-[39px] top-1 w-5 h-5 bg-white border-4 border-purple-300 rounded-full"></div>
                                    
                                    <div className="bg-white p-6 rounded-2xl shadow-sm border border-gray-100 hover:shadow-md transition-shadow">
                                        <div className="flex justify-between items-start mb-2">
                                            <h4 className="font-bold text-lg text-gray-900">{item.title}</h4>
                                            {item.dateRange && (
                                                <span className="text-xs font-bold text-white bg-gradient-to-r from-purple-400 to-pink-400 px-3 py-1 rounded-full shadow-sm">
                                                    {item.dateRange}
                                                </span>
                                            )}
                                        </div>
                                        {item.subtitle && (
                                            <div className="text-purple-600 font-medium mb-3">{item.subtitle}</div>
                                        )}
                                        {item.description && (
                                            <div 
                                                className="text-gray-600 text-sm leading-relaxed"
                                                dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }}
                                            />
                                        )}
                                    </div>
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
