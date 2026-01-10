import React, { useEffect, useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import { ResumeData } from '../../types';
import { API_BASE } from '../../config';
import { ResumePreview } from '../editor/ResumePreview';
import { INITIAL_RESUME } from '../../services/mockData';
import { useLanguage } from '../../contexts/LanguageContext';

export const PrintResume: React.FC = () => {
  const [searchParams] = useSearchParams();
  const [data, setData] = useState<ResumeData | null>(null);
  const { setLanguage } = useLanguage();

  useEffect(() => {
    const id = searchParams.get('id');
    const template = searchParams.get('template');
    if (template && !id) {
      const demo: ResumeData = { ...INITIAL_RESUME, templateId: template };
      setData(demo);
      return;
    }
    if (!id) return;
    const token = localStorage.getItem('token');
    fetch(`${API_BASE}/resumes/${id}`, {
      headers: token ? { Authorization: `Bearer ${token}` } : {}
    })
      .then(r => r.json())
      .then((res: any) => {
        const mapped: ResumeData = {
          id: res.ExternalID || id,
          title: res.Title,
          templateId: res.TemplateID,
          language: (res.Language || '') === 'en' ? 'en' : 'zh',
          Theme: res.Theme,
          lastModified: res.LastModified,
          Personal: res.Personal,
          sections: (res.Sections || []).map((s: any) => ({
            id: s.ExternalID || s.ID,
            type: s.Type,
            title: s.Title,
            isVisible: s.IsVisible,
            orderNum: s.OrderNum,
            items: (s.Items || []).map((i: any) => ({
              id: i.ExternalID || i.ID,
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
        try {
          if (mapped.language) setLanguage(mapped.language);
        } catch {}
      });
  }, [searchParams]);

  return (
    <div className="bg-white print:bg-white min-h-screen p-0 m-0 overflow-visible">
      {data && <ResumePreview data={data} scale={1} scrollInside={false} />}
    </div>
  );
};
