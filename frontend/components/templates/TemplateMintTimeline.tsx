import React from 'react';
import { ResumeData, ResumeSectionType } from '../../types';
import { useSectionTitle } from '../../hooks/useSectionTitle';
import { useLanguage } from '../../contexts/LanguageContext';
import { hasExtraPersonalInfo, sanitizeHtml } from '../../utils/resume-helpers';
import { Mail, Phone, MapPin, GraduationCap, User, Briefcase, Wrench, Layers, BookOpen, Award, Heart, Image as ImageIcon } from 'lucide-react';
import { ExamScoreTable } from './ExamScoreTable';

export const TemplateMintTimeline: React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }> = ({ data, styles, disableShadow }) => {
  const getSectionTitle = useSectionTitle();
  const { t } = useLanguage();
  const color = data.Theme?.Color || '#14b8a6';

  const findJobTarget = () => {
    return data.Personal?.Job || '';
  };
  const findCity = () => {
    return data.Personal?.City || '';
  };
  const findMoney = () => {
    return data.Personal?.Money || '';
  };
  const findJoinTime = () => {
    return data.Personal?.JoinTime || '';
  };
  const findDegree = () => {
    return data.Personal?.Degree || '';
  };

  const SectionIcon: React.FC<{ type: ResumeSectionType }> = ({ type }) => {
    switch (type) {
      case ResumeSectionType.Experience: return <Briefcase size={18} />;
      case ResumeSectionType.Education: return <GraduationCap size={18} />;
      case ResumeSectionType.Skills: return <Wrench size={18} />;
      case ResumeSectionType.Projects: return <Layers size={18} />;
      case ResumeSectionType.Internships: return <Briefcase size={18} />;
      case ResumeSectionType.Portfolio: return <ImageIcon size={18} />;
      case ResumeSectionType.Awards: return <Award size={18} />;
      case ResumeSectionType.Interests: return <Heart size={18} />;
      case ResumeSectionType.Exam: return <BookOpen size={18} />;
      case ResumeSectionType.SelfEvaluation: return <User size={18} />;
      default: return <User size={18} />;
    }
  };

  const renderItemTime = (item: any) => {
    if (item?.timeStart || item?.timeEnd || item?.today) {
      return (
        <span className="text-sm text-gray-600 font-bold" style={{ fontSize: styles.fontSize }}>
          {(item.timeStart || item.timeEnd || '').replace('-', '.')}
          {' - '}
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
  
  const showMarkerForSection = (sectionType: ResumeSectionType) => {
    return [
      ResumeSectionType.Experience,
      ResumeSectionType.Education,
      ResumeSectionType.Projects,
      ResumeSectionType.Internships
    ].includes(sectionType);
  };

  const sectionsOrdered = React.useMemo(() => {
    const visible = (data.sections || []).filter(s => s.isVisible);
    return visible.slice().sort((a, b) => (a.orderNum ?? 0) - (b.orderNum ?? 0));
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
              <h1 className="text-3xl font-bold tracking-wide text-white">{data.Personal?.FullName}</h1>
              {findJobTarget() && <p className="text-base mt-2 text-white/90">{t('editor.fields.jobApplication')}：{findJobTarget()}</p>}
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
                {data.Personal?.Gender && (
                  <span className="inline-flex items-center gap-1">
                    <User size={16} /> {data.Personal.Gender}
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
                {findMoney() && (
                  <span className="inline-flex items-center gap-1">
                    {t('editor.fields.expectedSalary')}：{findMoney()}
                  </span>
                )}
                {findJoinTime() && (
                  <span className="inline-flex items-center gap-1">
                    {t('editor.fields.joinTime')}：{findJoinTime()}
                  </span>
                )}
                {(() => {
                  try {
                    const raw = data.Personal?.CustomInfo;
                    if (raw) {
                      const arr = JSON.parse(raw);
                      if (Array.isArray(arr) && arr.length > 0) {
                        return arr.map((ci: any, idx: number) => (
                          <span key={idx} className="inline-flex items-center gap-1">
                            {ci.label}: {ci.value}
                          </span>
                        ));
                      }
                    }
                  } catch {}
                  return null;
                })()}
              </div>
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
        {/* <div className="w-full" style={{ backgroundColor: color }}>
          <svg viewBox="0 0 1200 60" preserveAspectRatio="none" className="w-full h-6 block">
            <path d="M0 30 Q 300 60 600 30 T 1200 30 L 1200 60 L 0 60 Z" fill="white"></path>
          </svg>
        </div> */}
      </div>

      <div className="px-10 pt-6 pb-10 space-y-1">
                {sectionsOrdered.map((section, idx) => (
                  (() => {
                    const items = (section.items || []).slice().sort((a, b) => (a.orderNum ?? 0) - (b.orderNum ?? 0));
                    const isLast = idx === sectionsOrdered.length - 1;
                    const dashTop = 24;
                    const dashBottom = isLast ? 2 : -55;
                    return (
                      <div key={section.id} className="flex gap-6 items-start">
                        <div className="relative mt-4 -translate-y-1/2" style={{ width: 140 }}>
                          <div className="text-base font-bold tracking-wide pl-3" style={{ color }}>
                            {(() => {
                              const title = String(getSectionTitle(section) || '').trim();
                              const parts = title.match(/^([A-Za-z]+)\s+([A-Za-z]+)$/);
                              if (parts) {
                                return (
                                  <span className="inline-block leading-tight">
                                    <span className="block">{parts[1]}</span>
                                    <span className="block">{parts[2]}</span>
                                  </span>
                                );
                              }
                              return title;
                            })()}
                          </div>
                        </div>
                        <div className="flex-1">
                          <div className="mt-4 border-t-2 pt-4 pb-4 relative" style={{ borderColor: color }}>
                    <div className="absolute border-l border-dashed" style={{ borderColor: color, left: -32, top: dashTop, bottom: dashBottom }} />
                    <div className="absolute left-[-48px] w-8 h-8 rounded-full flex items-center justify-center text-white ring-2 ring-white/60 print:ring-0 z-10" style={{ backgroundColor: color, top: -17 }}>
                      <SectionIcon type={section.type} />
                    </div>
                    {section.type === ResumeSectionType.Exam ? (
                      (() => {
                        const meta = items[0];
                        const scores = items.slice(1);
                        return (
                          <ExamScoreTable
                            color={color}
                            schoolLabel={t('exam.school')}
                            majorLabel={t('exam.major')}
                            scoreLabel={(meta?.description && String(meta.description).trim()) ? String(meta.description).trim() : t('exam.scoreLabel')}
                            school={meta?.title || ''}
                            major={meta?.subtitle || ''}
                            items={scores.map(s => ({ subject: s.title || '', score: s.subtitle || '' }))}
                          />
                        );
                      })()
                    ) : section.type === ResumeSectionType.Skills ? (
                      <div className="space-y-3">
                        {items.map(item => (
                          <div key={item.id} className="relative">
                            {showMarkerForSection(section.type) && (
                              <div className="absolute -translate-x-1/2 top-2 w-3 h-3 rotate-45" style={{ backgroundColor: color, left: -32 }} />
                            )}
                            {item.description && (
                          <div className="resume-rich-content text-gray-700 text-sm leading-relaxed" style={{ fontSize: styles.fontSize }} dangerouslySetInnerHTML={{ __html: sanitizeHtml(item.description) }} />
                        )}
                      </div>
                    ))}
                  </div>
                ) : (
                      <div className={items.length > 1 ? 'space-y-5' : ''}>
                        {items.map(item => (
                          <div key={item.id} className="relative">
                            <div className="relative flex justify-between items-baseline mb-1">
                              {showMarkerForSection(section.type) && (
                                <div className="absolute -translate-x-1/2 top-1/2 -translate-y-1/2 w-3 h-3 rotate-45" style={{ backgroundColor: color, left: -32 }} />
                              )}
                              {renderItemTime(item) && (
                                <div className="absolute -translate-y-1/2 text-gray-600 text-sm whitespace-nowrap" style={{ left: -164, top: '50%' }}>
                                  {renderItemTime(item)}
                                </div>
                              )}
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
