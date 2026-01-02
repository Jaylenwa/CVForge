import React, { useEffect, useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import { ResumeData } from '../../types';
import { API_BASE } from '../../config';
import { ResumePreview } from '../editor/ResumePreview';
import { INITIAL_RESUME } from '../../services/mockData';

export const PrintResume: React.FC = () => {
  const [searchParams] = useSearchParams();
  const [data, setData] = useState<ResumeData | null>(null);

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
          themeConfig: { color: res.ThemeColor, fontFamily: res.ThemeFont, spacing: res.ThemeSpacing },
          lastModified: res.LastModified,
          personalInfo: {
            fullName: res.FullName,
            jobTitle: res.JobTitle || '',
            email: res.Email,
            phone: res.Phone,
            avatarUrl: res.AvatarURL,
            gender: res.Gender,
            age: res.Age,
            maritalStatus: res.MaritalStatus,
            politicalStatus: res.PoliticalStatus,
            birthplace: res.Birthplace,
            ethnicity: res.Ethnicity,
            height: res.Height,
            weight: res.Weight,
            customInfo: (() => {
              try {
                if (res.CustomInfo) {
                  const parsed = JSON.parse(res.CustomInfo);
                  if (Array.isArray(parsed)) return parsed;
                }
              } catch {}
              return [];
            })()
          },
          sections: (res.Sections || []).map((s: any) => ({
            id: s.ExternalID || s.ID,
            type: s.Type,
            title: s.Title,
            isVisible: s.IsVisible,
            items: (s.Items || []).map((i: any) => ({
              id: i.ExternalID || i.ID,
              title: i.Title,
              subtitle: i.Subtitle,
              major: i.Major,
              degree: i.Degree,
              timeStart: i.TimeStart,
              timeEnd: i.TimeEnd,
              today: !!i.Today,
              description: i.Description
            }))
          }))
        };
        setData(mapped);
      });
  }, [searchParams]);

  return (
    <div className="bg-white print:bg-white min-h-screen p-0 m-0 overflow-visible">
      {data && <ResumePreview data={data} scale={1} scrollInside={false} />}
    </div>
  );
};
