import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useSectionTitle } from '../../../hooks/useSectionTitle';
import { useLanguage } from '../../../contexts/LanguageContext';
import { hasExtraPersonalInfo, sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateCNModern: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const titleMap: Record<string, string> = {
    [ResumeSectionType.Summary]: '自我评价',
    [ResumeSectionType.Education]: '教育背景',
    [ResumeSectionType.Experience]: '工作经历',
    [ResumeSectionType.Skills]: '专业技能',
    [ResumeSectionType.Projects]: '项目经验',
    [ResumeSectionType.Custom]: '其他信息'
  };

  return (
    <div className={`w-full bg-white h-auto ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none relative`} style={{ fontFamily: styles.fontFamily, lineHeight: parseFloat(styles.spacingMultiplier) * 1.5 }}>
      {/* 顶部斜切背景 */}
      <div className="absolute top-0 left-0 right-0 h-48 bg-slate-800" style={{ clipPath: 'polygon(0 0, 100% 0, 100% 60%, 0 100%)' }}></div>
      <div className="absolute top-0 right-0 w-64 h-64 bg-slate-700 opacity-20 rounded-full -translate-y-1/2 translate-x-1/4"></div>

      <div className="relative px-12 pt-12">
        {/* 头部信息 */}
        <div className="flex justify-between items-start mb-12 text-white">
          <div className="pt-4">
            <h1 className="text-4xl font-bold mb-2 tracking-wide">{data.personalInfo.fullName}</h1>
            <p className="text-xl text-slate-200 font-medium">{data.personalInfo.jobTitle}</p>
          </div>
          {data.personalInfo.avatarUrl && (
            <img 
              src={data.personalInfo.avatarUrl} 
              alt={t('a11y.avatarAlt')} 
              className="w-32 h-32 rounded-lg object-cover border-4 border-white/20 shadow-lg bg-white"
            />
          )}
        </div>

        {/* 联系信息条 */}
        <div className="bg-white shadow-md rounded-lg p-5 -mt-8 mb-10 relative z-10 border-l-4 border-slate-600">
          <div className="flex flex-wrap gap-y-2 gap-x-8 text-sm text-gray-600 justify-center md:justify-start">
             <div className="flex items-center">
               <span className="w-2 h-2 rounded-full bg-slate-400 mr-2"></span>
               <span className="font-semibold text-gray-900 mr-1">电话:</span> {data.personalInfo.phone}
             </div>
             <div className="flex items-center">
               <span className="w-2 h-2 rounded-full bg-slate-400 mr-2"></span>
               <span className="font-semibold text-gray-900 mr-1">邮箱:</span> <span className="break-all">{data.personalInfo.email}</span>
             </div>
             {data.personalInfo.gender && (
               <div className="flex items-center">
                 <span className="w-2 h-2 rounded-full bg-slate-400 mr-2"></span>
                 <span className="font-semibold text-gray-900 mr-1">性别:</span> {data.personalInfo.gender}
               </div>
             )}
             {data.personalInfo.age && (
               <div className="flex items-center">
                 <span className="w-2 h-2 rounded-full bg-slate-400 mr-2"></span>
                 <span className="font-semibold text-gray-900 mr-1">年龄:</span> {data.personalInfo.age}
               </div>
             )}
          </div>
        </div>

        {/* 两栏布局 */}
        <div className="grid grid-cols-1 md:grid-cols-12 gap-8">
          {/* 左侧窄栏 */}
          <div className="md:col-span-4 space-y-8">
            {/* 技能部分 */}
            {(data.sections || []).filter(s => s.type === ResumeSectionType.Skills && s.isVisible).map(section => (
              <div key={section.id}>
                <h3 className="text-lg font-bold text-slate-800 mb-4 pb-2 border-b-2 border-slate-200 inline-block">
                  {titleMap[section.type] || section.title}
                </h3>
                <div className="flex flex-col gap-2">
                  {section.items.map(item => (
                    <div key={item.id} className="bg-slate-50 px-3 py-2 rounded border-l-2 border-slate-400 text-sm text-slate-700">
                      {item.description}
                    </div>
                  ))}
                </div>
              </div>
            ))}
            
            {/* 补充信息 */}
            <div className="bg-slate-50 p-4 rounded-lg">
                <h4 className="font-bold text-slate-800 mb-2 text-sm">更多信息</h4>
                <div className="space-y-2 text-sm text-slate-600">
                     {data.personalInfo.politicalStatus && <div><span className="font-medium">政治面貌：</span>{data.personalInfo.politicalStatus}</div>}
                     {data.personalInfo.birthplace && <div><span className="font-medium">籍贯：</span>{data.personalInfo.birthplace}</div>}
                     {data.personalInfo.ethnicity && <div><span className="font-medium">民族：</span>{data.personalInfo.ethnicity}</div>}
                </div>
            </div>
          </div>

          {/* 右侧宽栏 */}
          <div className="md:col-span-8 space-y-10">
            {(data.sections || []).filter(s => s.type !== ResumeSectionType.Skills && s.isVisible).map(section => (
              <div key={section.id}>
                <div className="flex items-center mb-6">
                   <div className="w-1.5 h-6 bg-slate-700 mr-3"></div>
                   <h3 className="text-xl font-bold text-slate-800 tracking-wide">
                     {titleMap[section.type] || section.title}
                   </h3>
                   <div className="flex-grow h-px bg-slate-200 ml-4"></div>
                </div>
                
                <div className="space-y-6">
                  {section.items.map(item => (
                    <div key={item.id} className="relative pl-4 border-l border-dashed border-slate-300">
                      <div className="absolute -left-1.5 top-1.5 w-3 h-3 rounded-full bg-slate-400 border-2 border-white"></div>
                      <div className="flex justify-between items-baseline mb-1">
                        <h4 className="font-bold text-slate-900 text-lg">{item.title}</h4>
                        {item.dateRange && (
                          <span className="text-sm font-medium text-slate-500 bg-slate-100 px-2 py-0.5 rounded">
                            {item.dateRange}
                          </span>
                        )}
                      </div>
                      {item.subtitle && (
                        <div className="text-slate-600 font-medium mb-2">{item.subtitle}</div>
                      )}
                      {item.description && (
                        <div 
                          className="text-gray-600 text-sm leading-relaxed whitespace-pre-wrap"
                          dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }}
                        />
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
