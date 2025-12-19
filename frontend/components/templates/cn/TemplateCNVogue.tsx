import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useLanguage } from '../../../contexts/LanguageContext';
import { sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateCNVogue: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const titleMap: Record<string, string> = {
    [ResumeSectionType.Summary]: 'EDITOR\'S NOTE',
    [ResumeSectionType.Education]: 'ACADEMIC',
    [ResumeSectionType.Experience]: 'CAREER',
    [ResumeSectionType.Skills]: 'EXPERTISE',
    [ResumeSectionType.Projects]: 'FEATURED',
    [ResumeSectionType.Custom]: 'EXTRAS'
  };

  return (
    <div className={`w-full bg-white text-black h-full min-h-[1123px] ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto`} style={{ fontFamily: '"Didot", "Bodoni MT", "Playfair Display", serif' }}>
      {/* 顶部大标题 */}
      <div className="border-b-4 border-black p-8 text-center mb-8">
          <h1 className="text-7xl font-bold uppercase tracking-tighter mb-2">{data.personalInfo.fullName?.split(' ')[0] || 'NAME'}</h1>
          <h1 className="text-7xl font-bold uppercase tracking-tighter text-gray-400 mb-6">{data.personalInfo.fullName?.split(' ').slice(1).join(' ') || 'RESUME'}</h1>
          
          <div className="flex justify-between items-center border-t border-black pt-2 text-xs font-sans tracking-widest uppercase">
              <span>VOL. {new Date().getFullYear()}</span>
              <span className="font-bold">{data.personalInfo.jobTitle}</span>
              <span>ISSUE NO. 1</span>
          </div>
      </div>

      <div className="px-8 grid grid-cols-12 gap-8">
          {/* 左侧栏 - 窄 */}
          <div className="col-span-4 border-r border-black pr-8 flex flex-col h-full">
              {data.personalInfo.avatarUrl && (
                  <div className="mb-8 filter grayscale contrast-125">
                      <img src={data.personalInfo.avatarUrl} className="w-full h-auto object-cover" alt="Portrait" />
                      <div className="text-[10px] font-sans text-right mt-1 uppercase text-gray-500">Photo by Self</div>
                  </div>
              )}

              <div className="font-sans text-xs mb-8 space-y-1">
                  <div className="font-bold uppercase mb-2 text-base font-serif">Contact</div>
                  <div>{data.personalInfo.phone}</div>
                  <div>{data.personalInfo.email}</div>
                  <div>{data.personalInfo.city}</div>
              </div>

              {(data.sections || []).filter(s => s.type === ResumeSectionType.Skills && s.isVisible).map(section => (
                  <div key={section.id} className="mb-8">
                      <h3 className="text-2xl font-bold uppercase italic mb-4 border-b border-black pb-1">{titleMap[section.type] || section.title}</h3>
                      <div className="text-sm font-sans leading-relaxed text-justify">
                          {section.items.map(item => item.description).join(' • ')}
                      </div>
                  </div>
              ))}
          </div>

          {/* 右侧栏 - 宽 */}
          <div className="col-span-8">
              {(data.sections || []).filter(s => s.type !== ResumeSectionType.Skills && s.isVisible).map(section => (
              <div key={section.id} className="mb-10">
                  <h3 className="text-4xl font-bold uppercase mb-6 flex items-center">
                      <span className="text-6xl mr-2 leading-none">{section.title?.charAt(0) || ''}</span>
                      <span className="tracking-widest text-2xl border-b-2 border-black pb-1 w-full">{titleMap[section.type]?.slice(1) || section.title?.slice(1) || ''}</span>
                  </h3>

                  <div className="space-y-8">
                          {section.items.map(item => (
                              <div key={item.id}>
                                  <div className="flex justify-between items-baseline mb-2 font-sans">
                                      <h4 className="text-xl font-bold uppercase">{item.title}</h4>
                                      {item.dateRange && <span className="text-xs text-gray-500">{item.dateRange}</span>}
                                  </div>
                                  {item.subtitle && <div className="text-lg italic text-gray-600 mb-2 font-serif">{item.subtitle}</div>}
                                  {item.description && (
                                      <div className="text-sm font-sans leading-relaxed text-justify columns-1" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
                                  )}
                              </div>
                          ))}
                      </div>
                  </div>
              ))}
          </div>
      </div>
      
      {/* 底部条形码装饰 */}
      <div className="fixed bottom-8 right-8 w-32 opacity-50 hidden print:hidden md:block">
          <div className="h-10 bg-black" style={{ maskImage: 'repeating-linear-gradient(90deg, black, black 1px, transparent 1px, transparent 3px)' }}></div>
          <div className="text-center text-[10px] font-mono mt-1 tracking-[0.5em]">{data.personalInfo.phone?.replace(/\D/g, '').slice(0, 12) || '0000000000'}</div>
      </div>
    </div>
  );
};

