import React from 'react';
import { ResumeSection, ResumeSectionType } from '../../../types';
import { ExamScoreTable } from '../ExamScoreTable';
import { getOrderedItems } from './templateTokens';

export type ExamSectionProps = {
  section: ResumeSection;
  color: string;
  t: (key: string) => string;
};

export const ExamSection: React.FC<ExamSectionProps> = ({ section, color, t }) => {
  if (section.type !== ResumeSectionType.Exam) return null;

  const items = getOrderedItems(section.items || []);
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
};

export default ExamSection;

