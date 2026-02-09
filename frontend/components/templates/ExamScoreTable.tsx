import React from 'react';

type ExamSubjectScore = { subject: string; score: string };

export interface ExamScoreTableProps {
  color: string;
  schoolLabel: string;
  majorLabel: string;
  scoreLabel: string;
  school: string;
  major: string;
  items: ExamSubjectScore[];
  maxCols?: number;
}

const chunk = <T,>(arr: T[], size: number): T[][] => {
  if (size <= 0) return [arr];
  const out: T[][] = [];
  for (let i = 0; i < arr.length; i += size) out.push(arr.slice(i, i + size));
  return out;
};

export const ExamScoreTable: React.FC<ExamScoreTableProps> = ({
  color,
  schoolLabel,
  majorLabel,
  scoreLabel,
  school,
  major,
  items,
  maxCols = 6,
}) => {
  const normalized: ExamSubjectScore[] = (items || [])
    .map(it => ({ subject: String(it.subject || '').trim(), score: String(it.score || '').trim() }))
    .filter(it => it.subject || it.score);

  const blocks: ExamSubjectScore[][] = chunk<ExamSubjectScore>(normalized, maxCols);

  return (
    <div className="w-full border border-slate-200 rounded-md overflow-hidden bg-white" style={{ borderColor: color }}>
      <div className="grid grid-cols-2">
        <div className="px-4 py-3 border-b border-r border-slate-200">
          <div className="text-xs text-slate-500 mb-1">{schoolLabel}</div>
          <div className="text-sm font-semibold text-slate-900">{school}</div>
        </div>
        <div className="px-4 py-3 border-b border-slate-200">
          <div className="text-xs text-slate-500 mb-1">{majorLabel}</div>
          <div className="text-sm font-semibold text-slate-900">{major}</div>
        </div>
      </div>

      {blocks.map((block, bi) => {
        const cols = Math.max(1, block.length);
        const leftText = bi === 0 ? scoreLabel : '';
        return (
          <div key={`exam-block-${bi}`} className={bi === 0 ? '' : 'border-t border-slate-200'}>
            <div className="grid" style={{ gridTemplateColumns: `120px repeat(${cols}, minmax(72px, 1fr))` }}>
              <div className="row-span-2 px-4 py-4 bg-white text-sm font-semibold text-slate-700 flex items-center justify-center border-r border-slate-200 text-center">
                {leftText}
              </div>
              {Array.from({ length: cols }).map((_, i) => (
                <div
                  key={`exam-h-${bi}-${i}`}
                  className={`px-3 py-2 text-sm font-semibold text-slate-800 border-b border-slate-200 text-center ${i === cols - 1 ? '' : 'border-r border-slate-200'}`}
                >
                  {block[i]?.subject || ''}
                </div>
              ))}
              {Array.from({ length: cols }).map((_, i) => (
                <div
                  key={`exam-v-${bi}-${i}`}
                  className={`px-3 py-2 text-sm text-slate-700 text-center ${i === cols - 1 ? '' : 'border-r border-slate-200'}`}
                >
                  {block[i]?.score || ''}
                </div>
              ))}
            </div>
          </div>
        );
      })}
    </div>
  );
};

export default ExamScoreTable;
