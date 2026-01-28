import React, { useEffect, useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import { ResumeData } from '../../types';
import { API_BASE } from '../../config';
import { ResumePreview } from '../editor/ResumePreview';
import { INITIAL_RESUME } from '../../services/mockData';
import { useLanguage } from '../../contexts/LanguageContext';
import { fetchContentPresetData } from '../../services/catalogService';
import { applyTemplateDefaultsToResumeData } from '../../utils/template-defaults';

export const PrintResume: React.FC = () => {
  const [searchParams] = useSearchParams();
  const [data, setData] = useState<ResumeData | null>(null);
  const { setLanguage } = useLanguage();

  useEffect(() => {
    const id = searchParams.get('id');
    const template = searchParams.get('template');
    const presetId = searchParams.get('presetId') || '';

    if (template && !id) {
      let cancelled = false;

      const resolveSeed = () => {
        if (!presetId) return Promise.resolve({ ...(INITIAL_RESUME as any) });
        return fetchContentPresetData(Number(presetId))
          .then((parsed) => parsed && typeof parsed === 'object' ? ({ ...(INITIAL_RESUME as any), ...(parsed as any) }) : (INITIAL_RESUME as any))
          .catch(() => INITIAL_RESUME as any);
      };

      resolveSeed().then((seed: ResumeData) => {
        if (cancelled) return;
        const demo: ResumeData = applyTemplateDefaultsToResumeData({ ...(seed as any), templateId: template } as ResumeData);
        setData(demo);
        setLanguage(demo.language);
      });

      return () => { cancelled = true; };
    }
    if (!id) return;
    const token = localStorage.getItem('token');
    fetch(`${API_BASE}/resumes/${id}`, {
      headers: token ? { Authorization: `Bearer ${token}` } : {}
    })
      .then(r => r.json())
      .then((res: any) => {
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
            orderNum: s.OrderNum,
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
              orderNum: i.OrderNum
            })).sort((a: any, b: any) => (Number.isFinite(b.orderNum) || Number.isFinite(a.orderNum)) ? ((a.orderNum ?? 0) - (b.orderNum ?? 0)) : 0)
          })).sort((a: any, b: any) => (Number.isFinite(b.orderNum) || Number.isFinite(a.orderNum)) ? ((a.orderNum ?? 0) - (b.orderNum ?? 0)) : 0)
        };
        setData(mapped);
        setLanguage(mapped.language);
      });
  }, [searchParams]);

  return (
    <div className="bg-white print:bg-white min-h-screen p-0 m-0 overflow-visible">
      {data && <ResumePreview data={data} scale={1} scrollInside={false} />}
    </div>
  );
};
