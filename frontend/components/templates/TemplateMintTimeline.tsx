import React from 'react';
import { ResumeData, ResumeSectionType } from '../../types';
import { useSectionTitle } from '../../hooks/useSectionTitle';
import { useLanguage } from '../../contexts/LanguageContext';
import { hasExtraPersonalInfo, sanitizeHtml } from '../../utils/resume-helpers';
import { Mail, Phone, MapPin, GraduationCap, User, Briefcase, Wrench, Layers } from 'lucide-react';

export const TemplateMintTimeline: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const getSectionTitle = useSectionTitle();
  const { t } = useLanguage();
  const color = data.Theme?.Color || '#14b8a6';

  const findJobTarget = () => {
    const jobSection = (data.sections || []).find(s => s.type === ResumeSectionType.JobApplication);
    const first = jobSection?.items?.[0];
    return first?.title || data.Job?.Job || data.Personal?.JobTitle || '';
  };
  const findCity = () => {
    const jobSection = (data.sections || []).find(s => s.type === ResumeSectionType.JobApplication);
    const first = jobSection?.items?.[0];
    return data.Job?.City || first?.subtitle || '';
  };
  const findDegree = () => {
    const edu = (data.sections || []).find(s => s.type === ResumeSectionType.Education);
    const first = edu?.items?.[0];
    return first?.degree || '';
  };

  const SectionIcon: React.FC<{ type: ResumeSectionType }> = ({ type }) => {
    switch (type) {
      case ResumeSectionType.Experience: return <Briefcase size={18} />;
      case ResumeSectionType.Education: return <GraduationCap size={18} />;
      case ResumeSectionType.Skills: return <Wrench size={18} />;
      case ResumeSectionType.Projects: return <Layers size={18} />;
      default: return <User size={18} />;
    }
  };

  const renderItemTime = (item: any) => {
    if (item?.timeStart || item?.timeEnd || item?.today) {
      return (
        <span className="text-sm text-gray-600 font-medium" style={{ fontSize: styles.fontSize }}>
          {(item.timeStart || item.timeEnd || '').replace('-', '.')}
          {' ~ '}
          {item.today ? t('common.toPresent') : (item.timeEnd || '').replace('-', '.')}
        </span>
      );
    }
    return null;
  };

  const getRightTag = (sectionType: ResumeSectionType, item: any) => {
    if (sectionType === ResumeSectionType.Experience) {
      return item?.subtitle || null;
    }
    if (sectionType === ResumeSectionType.Internships) {
      return item?.subtitle || null;
    }
    if (sectionType === ResumeSectionType.Education) {
      if (item?.major || item?.degree) {
        const major = item?.major || '';
        const degree = item?.degree || '';
        if (major && degree) return `${major}（${degree}）`;
        return major || degree || null;
      }
      return null;
    }
    if (sectionType === ResumeSectionType.Projects) {
      return item?.subtitle || null;
    }
    return null;
  };

  const sectionsOrdered = React.useMemo(() => {
    const order = [
      ResumeSectionType.Experience,
      ResumeSectionType.Education,
      ResumeSectionType.Skills,
      ResumeSectionType.Projects
    ];
    const visible = (data.sections || []).filter(s => s.isVisible);
    const main = visible.filter(s => order.includes(s.type)).sort((a, b) => order.indexOf(a.type) - order.indexOf(b.type));
    const others = visible.filter(s => !order.includes(s.type));
    return [...main, ...others];
  }, [data.sections]);

  const hasMeaningfulContent = (item: any, type: ResumeSectionType) => {
    if (!item) return false;
    if (type === ResumeSectionType.Skills) {
      return !!(item.description && String(item.description).trim());
    }
    return !!(
      (item.title && String(item.title).trim()) ||
      (item.subtitle && String(item.subtitle).trim()) ||
      (item.major && String(item.major).trim()) ||
      (item.degree && String(item.degree).trim()) ||
      (item.description && String(item.description).trim()) ||
      (item.timeStart && String(item.timeStart).trim()) ||
      (item.timeEnd && String(item.timeEnd).trim()) ||
      item.today
    );
  };

  return (
    <div className={`w-full bg-white text-gray-900 h-auto ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none`} style={{ fontFamily: styles.fontFamily, lineHeight: parseFloat(styles.spacingMultiplier) * 1.5, fontSize: styles.fontSize }}>
      <div className="relative">
        <div className="bg-teal-500" style={{ backgroundColor: color }}>
          <div className="px-10 pt-8 pb-6 flex items-start gap-6">
            <div className="flex-1">
              <h1 className="text-4xl font-bold tracking-wide text-white">{data.Personal?.FullName}</h1>
              {findJobTarget() && <p className="text-lg mt-2 text-white/90">求职目标：{findJobTarget()}</p>}
              <div className="mt-4 flex flex-wrap gap-4 text-sm text-white/90">
                {data.Personal?.Age && (
                  <span className="inline-flex items-center gap-1">
                    <User size={16} /> {data.Personal.Age}
                  </span>
                )}
                {data.Personal?.Phone && (
                  <span className="inline-flex items-center gap-1">
                    <Phone size={16} /> {data.Personal.Phone}
                  </span>
                )}
                {data.Personal?.Email && (
                  <span className="inline-flex items-center gap-1">
                    <Mail size={16} /> {data.Personal.Email}
                  </span>
                )}
                {findDegree() && (
                  <span className="inline-flex items-center gap-1">
                    <GraduationCap size={16} /> {findDegree()}
                  </span>
                )}
                {findCity() && (
                  <span className="inline-flex items-center gap-1">
                    <MapPin size={16} /> {findCity()}
                  </span>
                )}
              </div>
              {hasExtraPersonalInfo(data) && (
                <div className="mt-2 flex flex-wrap gap-x-4 gap-y-1 text-xs text-white/90">
                  {data.Personal?.Gender && <span>{t('editor.fields.gender')}: {data.Personal.Gender}</span>}
                  {data.Personal?.MaritalStatus && <span>{t('editor.fields.maritalStatus')}: {data.Personal.MaritalStatus}</span>}
                  {data.Personal?.PoliticalStatus && <span>{t('editor.fields.politicalStatus')}: {data.Personal.PoliticalStatus}</span>}
                  {data.Personal?.Birthplace && <span>{t('editor.fields.birthplace')}: {data.Personal.Birthplace}</span>}
                </div>
              )}
            </div>
            {data.Personal?.AvatarURL && (
              <div className="flex-shrink-0">
                <img 
                  src={data.Personal.AvatarURL} 
                  alt={t('a11y.avatarAlt')} 
                  className="w-[105px] h-[147px] rounded-md object-cover ring-2 ring-white/70"
                />
              </div>
            )}
          </div>
        </div>
        <div className="w-full" style={{ backgroundColor: color }}>
          <svg viewBox="0 0 1200 60" preserveAspectRatio="none" className="w-full h-6 block">
            <path d="M0 30 Q 300 60 600 30 T 1200 30 L 1200 60 L 0 60 Z" fill="white"></path>
          </svg>
        </div>
      </div>

      <div className="px-10 pt-6 pb-10 space-y-10">
                {sectionsOrdered.map((section, idx) => (
                  (() => {
                    const items = (section.items || []).filter(it => hasMeaningfulContent(it, section.type));
                    if (items.length === 0) return null;
                    const isLast = idx === sectionsOrdered.length - 1;
                    const dashTop = 24;
                    const dashBottom = isLast ? 2 : -55;
                    return (
                      <div key={section.id} className="flex gap-6 items-start">
                        <div className="relative" style={{ width: 140 }}>
                          {/* 标题字体大小 */}
                          <div className="text-xl font-bold tracking-wide" style={{ color }}>
                            {getSectionTitle(section)}
                          </div>                  
                        </div>
                        <div className="flex-1">
                          <div className="mt-4 border-t-2 pt-4 pb-4 relative" style={{ borderColor: color }}>
                    <div className="absolute border-l border-dashed" style={{ borderColor: color, left: -32, top: dashTop, bottom: dashBottom }} />
                    <div className="absolute left-[-32px] -translate-x-1/2 -translate-y-1/2 w-8 h-8 rounded-full flex items-center justify-center text-white shadow z-10" style={{ backgroundColor: color, top: -1 }}>
                      <SectionIcon type={section.type} />
                    </div>
                    {section.type === ResumeSectionType.Skills ? (
                      <div className="space-y-3">
                        {items.map(item => (
                          <div key={item.id} className="relative pl-1">
                            <div className="absolute -translate-x-1/2 top-2 w-3 h-3 rotate-45" style={{ backgroundColor: color, left: -32 }} />
                            {item.description && (
                          <div className="resume-rich-content text-gray-700 text-sm leading-relaxed" style={{ fontSize: styles.fontSize }} dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
                        )}
                      </div>
                    ))}
                  </div>
                ) : (
                      <div className={items.length > 1 ? 'space-y-5' : ''}>
                        {items.map(item => (
                          <div key={item.id} className="relative pl-1">
                            <div className="absolute -translate-x-1/2 top-2 w-3 h-3 rotate-45" style={{ backgroundColor: color, left: -32 }} />
                            {renderItemTime(item) && (
                              <div className="absolute -translate-y-1/2 text-gray-600 text-sm whitespace-nowrap" style={{ left: -164, top: 14 }}>
                                {renderItemTime(item)}
                              </div>
                            )}
                            <div className="flex justify-between items-baseline mb-1">
                              <div className="font-semibold text-gray-900 text-sm">
                                {item.title}
                              </div>
                              {getRightTag(section.type, item) && (
                                <span className="text-gray-900 font-medium">{getRightTag(section.type, item)}</span>
                              )}
                            </div>
                            {item.description && (
                              <div className="resume-rich-content text-gray-700 text-sm leading-relaxed mt-1" style={{ fontSize: styles.fontSize }} dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
                            )}
                          </div>
                        ))}
                      </div>
                    )}
                  </div>
                </div>
              </div>
            );
          })()
        ))}
      </div>
    </div>
  );
};

export default TemplateMintTimeline;
