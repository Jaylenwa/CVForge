import React, { useEffect, useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import { ResumeData } from '../../types';
import { ResumePreview } from '../editor/ResumePreview';
import { useLanguage } from '../../contexts/LanguageContext';
import { fetchContentPresetData } from '../../services/catalogService';
import { applyTemplateDefaultsToResumeData } from '../../utils/template-defaults';
import { apiJson } from '../../services/apiClient';

export const PrintResume: React.FC = () => {
  const [searchParams] = useSearchParams();
  const [data, setData] = useState<ResumeData | null>(null);
  const { language, setLanguage } = useLanguage();

  useEffect(() => {
    const id = searchParams.get('id');
    const template = searchParams.get('template');
    const presetId = searchParams.get('presetId') || '';
    const toOrderNum = (v: any): number | undefined => {
      const n = typeof v === 'number' ? v : Number(v);
      return Number.isFinite(n) ? n : undefined;
    };

    if (template && !id) {
      let cancelled = false;

      const resolveSeed = () => {
        if (!presetId) return Promise.resolve(null);
        return fetchContentPresetData(Number(presetId), undefined, language)
          .then((parsed) => (parsed && typeof parsed === 'object') ? parsed : null)
          .catch(() => null);
      };

      resolveSeed().then((seed: any) => {
        if (cancelled) return;
        const normalize = (raw: string): 'en' | 'zh' => (String(raw || '').trim().toLowerCase() === 'en' ? 'en' : 'zh');
        const baseLang = normalize(String((seed as any)?.language || language || ''));
        const fallback: ResumeData = applyTemplateDefaultsToResumeData({
          id: 'preview',
          title: baseLang === 'en' ? 'Resume' : '我的简历',
          templateId: template,
          lastModified: Date.now(),
          language: baseLang,
          Personal: {
            FullName: baseLang === 'en' ? 'Alex Chen' : '陈小明',
            Email: 'alex@example.com',
            Phone: baseLang === 'en' ? '+1 555 0100' : '13800000000',
            City: baseLang === 'en' ? 'Shanghai' : '上海',
            Job: baseLang === 'en' ? 'Software Engineer' : '软件工程师',
          },
          Theme: {},
          sections: [
            {
              id: 'selfEvaluation',
              type: 'selfEvaluation' as any,
              title: baseLang === 'en' ? 'Self Evaluation' : '自我评价',
              isVisible: true,
              items: [{ id: 's1', description: baseLang === 'en' ? 'A concise summary goes here.' : '这里是一段简短的个人简介。' }],
            },
            {
              id: 'exp',
              type: 'experience' as any,
              title: baseLang === 'en' ? 'Experience' : '工作经历',
              isVisible: true,
              items: [
                {
                  id: 'e1',
                  title: baseLang === 'en' ? 'Company A' : '公司 A',
                  subtitle: baseLang === 'en' ? 'Software Engineer' : '软件工程师',
                  timeStart: '2023-01',
                  timeEnd: '2024-12',
                  description: baseLang === 'en' ? 'Built features and improved performance.' : '负责功能开发与性能优化。',
                },
              ],
            },
          ],
        });
        const demo: ResumeData = seed
          ? applyTemplateDefaultsToResumeData({ ...(seed as any), templateId: template } as ResumeData)
          : fallback;
        setData(demo);
        const nextLang = normalize(String((demo as any)?.language || baseLang));
        setLanguage(nextLang);
      });

      return () => { cancelled = true; };
    }
    if (!id) return;
    const token = localStorage.getItem('token');
    (async () => {
      try {
        const res = await apiJson<any>(`/resumes/${id}`, token ? { auth: true } : undefined);
        const mapped: ResumeData = {
          id: res.ID || Number(id),
          title: res.Title,
          templateId: res.TemplateID,
          language: (res.Language || '') === 'en' ? 'en' : 'zh',
          Theme: res.Theme,
          lastModified: res.LastModified,
          Personal: res.Personal,
          sections: (res.Sections || []).map((s: any) => ({
            id: s.ID,
            type: s.Type,
            title: s.Title,
            isVisible: s.IsVisible,
            orderNum: toOrderNum(s.OrderNum),
            items: (s.Items || []).map((i: any) => ({
              id: i.ID,
              title: i.Title,
              subtitle: i.Subtitle,
              major: i.Major,
              degree: i.Degree,
              timeStart: i.TimeStart,
              timeEnd: i.TimeEnd,
              today: !!i.Today,
              description: i.Description,
              orderNum: toOrderNum(i.OrderNum)
            })).sort((a: any, b: any) => (Number.isFinite(b.orderNum) || Number.isFinite(a.orderNum)) ? ((a.orderNum ?? 0) - (b.orderNum ?? 0)) : 0)
          })).sort((a: any, b: any) => (Number.isFinite(b.orderNum) || Number.isFinite(a.orderNum)) ? ((a.orderNum ?? 0) - (b.orderNum ?? 0)) : 0)
        };
        setData(mapped);
        setLanguage(mapped.language);
      } catch {}
    })();
  }, [searchParams]);

  return (
    <div className="bg-white print:bg-white min-h-screen p-0 m-0 overflow-visible">
      {data && <ResumePreview data={data} scale={1} scrollInside={false} />}
    </div>
  );
};
