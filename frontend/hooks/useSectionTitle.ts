import { useLanguage } from '../contexts/LanguageContext';

export const useSectionTitle = () => {
  const { t } = useLanguage();

  const DEFAULTS: Record<string, string[]> = {
    personal: ['Personal Information', '个人信息'],
    experience: ['Work Experience', '工作经历'],
    education: ['Education', '教育背景'],
    skills: ['Skills', '技能特长'],
    projects: ['Projects', '项目经历'],
    internships: ['Internship Experience', '实习经验'],
    portfolio: ['Personal Works', '个人作品'],
    awards: ['Honors & Certificates', '荣誉证书'],
    interests: ['Interests & Hobbies', '兴趣爱好'],
    jobApplication: ['Job Application', '求职岗位'],
    exam: ['Exam Information', '报考信息'],
    selfEvaluation: ['Self Evaluation', '自我评价'],
    custom: ['Custom Section', '自定义模块'],
  };

  return (section: any) => {
    const typeKey = String(section.type);
    const title = String(section.title || '');
    const synonyms = DEFAULTS[typeKey] || [];
    const translated = t(`section.${typeKey}`);
    if (!title || title === translated || synonyms.includes(title)) {
      return translated;
    }
    return title;
  };
};
