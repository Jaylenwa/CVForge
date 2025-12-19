import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useLanguage } from '../../../contexts/LanguageContext';
import { sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateCNAcademic: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const titleMap: Record<string, string> = {
    [ResumeSectionType.Summary]: '研究兴趣与自我评价',
    [ResumeSectionType.Education]: '教育背景',
    [ResumeSectionType.Experience]: '学术与工作经历',
    [ResumeSectionType.Skills]: '专业技能',
    [ResumeSectionType.Projects]: '科研项目',
    [ResumeSectionType.Custom]: '其他信息'
  };

  return (
    <div className={`w-full bg-[#fdfbf7] text-gray-900 h-full min-h-[1123px] ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto`} style={{ fontFamily: '"Times New Roman", "Songti SC", serif', lineHeight: '1.8' }}>
      <div className="px-16 py-16">
        {/* 头部 */}
        <div className="text-center mb-12 border-b-4 border-double border-gray-300 pb-8">
           <h1 className="text-4xl font-bold mb-4 tracking-wider text-gray-800">{data.personalInfo.fullName}</h1>
           <p className="text-xl text-gray-600 italic mb-6">{data.personalInfo.jobTitle}</p>
           
           <div className="flex justify-center flex-wrap gap-6 text-sm text-gray-600 font-medium">
              <span>{data.personalInfo.phone}</span>
              <span>•</span>
              <span>{data.personalInfo.email}</span>
              {data.personalInfo.age && (
                <>
                  <span>•</span>
                  <span>{data.personalInfo.age}</span>
                </>
              )}
           </div>
           
           {data.personalInfo.avatarUrl && (
             <img
               src={data.personalInfo.avatarUrl}
               alt={t('a11y.avatarAlt')}
               className="w-24 h-24 rounded-md object-cover mx-auto mt-6 border border-gray-300"
             />
           )}
        </div>

        {/* 内容 */}
        <div className="space-y-10">
           {(data.sections || []).filter(s => s.isVisible).map(section => (
              <div key={section.id}>
                 <h3 className="text-lg font-bold text-gray-900 uppercase tracking-widest border-b border-gray-400 pb-1 mb-6 flex items-center">
                    <span className="bg-gray-800 w-4 h-4 mr-3 block"></span>
                    {titleMap[section.type] || section.title}
                 </h3>

                 {section.type === ResumeSectionType.Skills ? (
                    <div className="grid grid-cols-2 gap-4">
                       {section.items.map(item => (
                          <div key={item.id} className="flex items-start">
                             <span className="mr-2 text-gray-500">❧</span>
                             <span className="text-sm text-gray-800 italic">{item.description}</span>
                          </div>
                       ))}
                    </div>
                 ) : (
                    <div className="space-y-8">
                       {section.items.map(item => (
                          <div key={item.id}>
                             <div className="flex justify-between items-baseline mb-1">
                                <h4 className="text-xl font-bold text-gray-900">{item.title}</h4>
                                {item.dateRange && <span className="text-sm text-gray-500 font-serif italic">{item.dateRange}</span>}
                             </div>
                             {item.subtitle && <div className="text-base text-gray-700 font-semibold mb-2">{item.subtitle}</div>}
                             {item.description && (
                                <div className="text-sm text-gray-600 leading-relaxed text-justify" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
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
  );
};
