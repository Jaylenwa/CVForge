import React from 'react';
import { ResumeData, ResumeSectionType } from '../../../types';
import { useSectionTitle } from '../../../hooks/useSectionTitle';
import { useLanguage } from '../../../contexts/LanguageContext';
import { hasExtraPersonalInfo, sanitizeHtml } from '../../../utils/resume-helpers';

export const TemplateCNBlue: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const { t } = useLanguage();
  const getSectionTitle = useSectionTitle();
  const SectionHeader: React.FC<{ title: string }> = ({ title }) => (
    <div className="bg-blue-100 text-blue-900 font-bold px-4 py-2 rounded">{title}</div>
  );
  const SkillBar: React.FC<{ text: string }> = ({ text }) => {
    const raw = text || '';
    const name = raw.includes(':') ? raw.split(':')[0] : (raw.includes(' ') ? raw.split(' ')[0] : raw);
    const lower = raw.toLowerCase();
    let level = 3;
    let label = '良好';
    if (raw.includes('精通') || lower.includes('expert') || lower.includes('advanced')) { level = 5; label = '精通'; }
    else if (raw.includes('熟练') || lower.includes('proficient')) { level = 4; label = '熟练'; }
    else if (raw.includes('一般') || lower.includes('beginner')) { level = 2; label = '一般'; }
    const pct = `${Math.min(100, Math.max(0, level * 20))}%`;
    return (
      <div className="mt-2">
        <div className="flex justify-between items-center text-sm">
          <span className="font-medium text-gray-800">{name}</span>
          <span className="text-gray-500">{label}</span>
        </div>
        <div className="h-1.5 bg-gray-200 rounded mt-1 overflow-hidden">
          <div className="h-full bg-blue-500" style={{ width: pct }} />
        </div>
      </div>
    );
  };
  return (
    <div className={`w-full bg-white h-full min-h-[1123px] ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none print:min-h-0 print:h-auto`} style={{ fontFamily: styles.fontFamily, lineHeight: parseFloat(styles.spacingMultiplier) * 1.5 }}>
      <div className="hidden print:block fixed left-0 top-0 bottom-0 w-[12mm] -z-10" style={{ backgroundColor: '#4b5563' }}></div>
      <div className="px-10 pt-10">
        <div className="text-center mb-6">
          <h1 className="text-3xl font-bold text-blue-700 tracking-wide">个人简历</h1>
          <p className="text-gray-500 mt-1 text-sm">{data.personalInfo.jobTitle || '努力超越自己，每天进步一点点'}</p>
        </div>
        <div className="mb-8">
          <SectionHeader title="基本信息" />
          <div className="grid grid-cols-1 md:grid-cols-12 gap-6 mt-4 items-start">
            <div className="md:col-span-9 grid grid-cols-2 gap-y-2 text-sm">
              <div><span className="text-gray-500 mr-2">姓名：</span><span className="text-gray-900 font-medium">{data.personalInfo.fullName}</span></div>
              <div><span className="text-gray-500 mr-2">电话：</span><span className="text-gray-900">{data.personalInfo.phone}</span></div>
              <div className="col-span-2"><span className="text-gray-500 mr-2">邮箱：</span><span className="text-gray-900 break-all">{data.personalInfo.email}</span></div>
              {/* website removed */}
              {data.personalInfo.gender && <div><span className="text-gray-500 mr-2">性别：</span><span className="text-gray-900">{data.personalInfo.gender}</span></div>}
              {data.personalInfo.age && <div><span className="text-gray-500 mr-2">年龄：</span><span className="text-gray-900">{data.personalInfo.age}</span></div>}
              {data.personalInfo.maritalStatus && <div><span className="text-gray-500 mr-2">婚姻：</span><span className="text-gray-900">{data.personalInfo.maritalStatus}</span></div>}
              {data.personalInfo.politicalStatus && <div><span className="text-gray-500 mr-2">政治面貌：</span><span className="text-gray-900">{data.personalInfo.politicalStatus}</span></div>}
              {data.personalInfo.birthplace && <div><span className="text-gray-500 mr-2">籍贯：</span><span className="text-gray-900">{data.personalInfo.birthplace}</span></div>}
              {data.personalInfo.ethnicity && <div><span className="text-gray-500 mr-2">民族：</span><span className="text-gray-900">{data.personalInfo.ethnicity}</span></div>}
              {data.personalInfo.height && <div><span className="text-gray-500 mr-2">身高：</span><span className="text-gray-900">{data.personalInfo.height}</span></div>}
              {data.personalInfo.weight && <div><span className="text-gray-500 mr-2">体重：</span><span className="text-gray-900">{data.personalInfo.weight}</span></div>}
              {(data.personalInfo.customInfo || []).map((ci, idx) => (
                <div key={idx} className="col-span-2">
                  <span className="text-gray-500 mr-2">{ci.label}：</span><span className="text-gray-900 break-all">{ci.value}</span>
                </div>
              ))}
            </div>
            {data.personalInfo.avatarUrl && (
              <div className="md:col-span-3 flex md:justify-end">
                <img
                  src={data.personalInfo.avatarUrl}
                  alt={t('a11y.avatarAlt')}
                  className="w-28 h-28 rounded-md object-cover border-2 shadow-sm"
                  style={{ borderColor: data.themeConfig?.color || '#93c5fd' }}
                />
              </div>
            )}
          </div>
        </div>
        {data.sections.filter(s => s.isVisible).map(section => (
          <div key={section.id} className="mb-6">
            <SectionHeader title={getSectionTitle(section)} />
            <div className="pt-4 space-y-4">
              {section.type === ResumeSectionType.Skills
                ? (
                  <div className="space-y-2">
                    {section.items.map(it => (
                      <div key={it.id} className="text-sm text-gray-700 leading-relaxed whitespace-pre-wrap">
                        {it.description}
                      </div>
                    ))}
                  </div>
                )
                : section.items.map(item => (
                  <div key={item.id} className="border-b border-gray-200 pb-3">
                    <div className="flex justify-between items-baseline">
                      <h4 className="font-bold text-gray-900">{item.title}</h4>
                      {item.dateRange && <span className="text-sm text-gray-600">{item.dateRange}</span>}
                    </div>
                    {item.subtitle && (
                      <div className="mt-0.5 text-sm font-medium text-gray-700">
                        {item.subtitle}
                      </div>
                    )}
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

