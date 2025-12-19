import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useSectionTitle } from '../../../hooks/useSectionTitle';
import { useLanguage } from '../../../contexts/LanguageContext';
import { hasExtraPersonalInfo, sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateCNZiyuan: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const titleMap: Record<string, string> = {
    [ResumeSectionType.Summary]: '自我评价',
    [ResumeSectionType.Education]: '教育背景',
    [ResumeSectionType.Experience]: '工作经验',
    [ResumeSectionType.Skills]: '职业技能',
    [ResumeSectionType.Projects]: '项目经历',
    [ResumeSectionType.Custom]: '自定义模块'
  };
  const eduSection = (data.sections || []).find(s => s.type === ResumeSectionType.Education && s.isVisible);
  const eduItem = eduSection && eduSection.items && eduSection.items[0];
  const degreeText = eduItem ? (eduItem.subtitle || eduItem.title || '') : '';
  const findCustom = (label: string) => {
    const list = data.personalInfo.customInfo || [];
    const item = list.find(ci => ci.label.toLowerCase().includes(label.toLowerCase()));
    return item ? item.value : '';
  };
  const mmBg = '#e5eef8';
  const primary = '#567ca8';
  return (
    <div className={`w-full bg-white h-full min-h-[1123px] ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto`} style={{ fontFamily: styles.fontFamily, lineHeight: parseFloat(styles.spacingMultiplier) * 1.5 }}>
      <div className="hidden print:block fixed left-0 top-0 bottom-0 w-[12mm] -z-10" style={{ backgroundColor: '#e5e7eb' }}></div>
      <div className="px-10 pt-10">
        <div className="flex items-end">
          <div className="inline-flex items-center rounded-full px-6 py-4" style={{ backgroundColor: primary }}>
            <span className="text-white text-2xl font-bold mr-4">求职简历</span>
            <div className="text-white/90 text-sm font-medium">PERSONAL RESUME</div>
          </div>
          <div className="ml-4 text-gray-500 text-sm">做一生的努力!</div>
        </div>
        <div className="mt-6 h-px w-full" style={{ backgroundColor: '#c0cbd8' }} />
        <div className="mt-10">
          <div className="inline-block rounded-full px-4 py-2 text-sm font-bold" style={{ backgroundColor: mmBg, color: primary }}>基本信息</div>
          <div className="grid grid-cols-12 gap-6 mt-4 items-start">
            <div className="col-span-12 md:col-span-9 bg-white">
              <div className="grid grid-cols-2 gap-y-2 text-sm">
                <div><span className="text-gray-500 mr-2">姓 名：</span><span className="text-gray-900 font-medium">{data.personalInfo.fullName}</span></div>
                {data.personalInfo.age && <div><span className="text-gray-500 mr-2">年 龄：</span><span className="text-gray-900">{data.personalInfo.age}</span></div>}
                <div><span className="text-gray-500 mr-2">学 历：</span><span className="text-gray-900">{degreeText || ' '}</span></div>
                <div><span className="text-gray-500 mr-2">求职意向：</span><span className="text-gray-900">{data.personalInfo.jobTitle}</span></div>
                <div><span className="text-gray-500 mr-2">手 机：</span><span className="text-gray-900">{data.personalInfo.phone}</span></div>
                <div className="col-span-2"><span className="text-gray-500 mr-2">邮 箱：</span><span className="text-gray-900 break-all">{data.personalInfo.email}</span></div>
                {findCustom('微信') && <div className="col-span-2"><span className="text-gray-500 mr-2">微 信：</span><span className="text-gray-900 break-all">{findCustom('微信')}</span></div>}
                {findCustom('地址') && <div className="col-span-2"><span className="text-gray-500 mr-2">地 址：</span><span className="text-gray-900 break-all">{findCustom('地址')}</span></div>}
              </div>
            </div>
            {data.personalInfo.avatarUrl && (
              <div className="col-span-12 md:col-span-3 flex md:justify-end">
                <img src={data.personalInfo.avatarUrl} alt={t('a11y.avatarAlt')} className="w-28 h-28 rounded-md object-cover border-2 shadow-sm" style={{ borderColor: primary }} />
              </div>
            )}
          </div>
        </div>
        {(data.sections || []).filter(s => s.isVisible).map(section => (
          <div key={section.id} className="mt-8">
            <div className="inline-block rounded-full px-4 py-2 text-sm font-bold" style={{ backgroundColor: mmBg, color: primary }}>
              {titleMap[section.type] || section.title}
            </div>
            <div className="pt-4 space-y-4">
              {section.type === ResumeSectionType.Skills ? (
                <div className="space-y-2">
                  {section.items.map(it => (
                    <div key={it.id} className="text-sm text-gray-700 leading-relaxed whitespace-pre-wrap">
                      {it.description}
                    </div>
                  ))}
                </div>
              ) : section.items.map(item => (
                <div key={item.id} className="border-b border-gray-200 pb-3">
                  <div className="flex justify-between items-baseline">
                    <h4 className="font-bold text-gray-900">{item.title}</h4>
                    {item.dateRange && <span className="text-sm text-gray-600">{item.dateRange}</span>}
                  </div>
                  {item.subtitle && <div className="mt-0.5 text-sm font-medium text-gray-700">{item.subtitle}</div>}
                  <div className="text-sm text-gray-700 leading-relaxed whitespace-pre-wrap mt-1" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description || '') }} />
                </div>
              ))}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

