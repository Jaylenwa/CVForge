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
        <span className="text-sm text-gray-600 font-medium">
          {(item.timeStart || item.timeEnd || '').replace('-', '.')}
          {' ~ '}
          {item.today ? t('common.toPresent') : (item.timeEnd || '').replace('-', '.')}
        </span>
      );
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

  return (
    <div className={`w-full bg-white text-gray-900 h-auto ${disableShadow ? 'shadow-none' : 'shadow-lg'} print:shadow-none`} style={{ fontFamily: styles.fontFamily, lineHeight: parseFloat(styles.spacingMultiplier) * 1.5 }}>
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
                  className="w-28 h-28 rounded-md object-cover ring-2 ring-white/70"
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

      <div className="px-10 pt-6 pb-10 space-y-8">
        {sectionsOrdered.map((section, idx) => (
          <div key={section.id} className="flex gap-6">
            <div className="w-10 relative">
              <div className="absolute left-1/2 -translate-x-1/2 top-0 bottom-0 border-l-2" style={{ borderColor: color }} />
              <div className="absolute left-1/2 -translate-x-1/2 -top-1.5 w-7 h-7 rounded-full flex items-center justify-center text-white shadow" style={{ backgroundColor: color }}>
                <SectionIcon type={section.type} />
              </div>
            </div>
            <div className="flex-1">
              <h3 className="text-lg font-bold mb-3 flex items-center justify-between" style={{ color }}>
                <span>{getSectionTitle(section)}</span>
              </h3>

              {section.type === ResumeSectionType.Skills ? (
                <div className="space-y-3">
                  {section.items.map(item => (
                    <div key={item.id}>
                      {item.description && (
                        <div className="resume-rich-content text-gray-700 text-sm leading-relaxed" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
                      )}
                    </div>
                  ))}
                </div>
              ) : (
                <div className="space-y-5">
                  {section.items.map(item => (
                    <div key={item.id}>
                      <div className="flex flex-col md:flex-row md:justify-between md:items-baseline mb-1">
                        <div className="font-semibold text-gray-900 text-base">
                          {item.title}
                          {item.subtitle && <span className="text-gray-700 font-medium ml-2"> {item.subtitle}</span>}
                          {section.type === ResumeSectionType.Education && (item.major || item.degree) && (
                            <span className="text-gray-700 text-sm ml-2">
                              {item.major}{item.major && item.degree ? ' • ' : ''}{item.degree}
                            </span>
                          )}
                        </div>
                        {renderItemTime(item)}
                      </div>
                      {item.description && (
                        <div className="resume-rich-content text-gray-700 text-sm leading-relaxed" dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
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
  );
};

export default TemplateMintTimeline;
