import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { ResumeData } from '../../types';
import { API_BASE } from '../../config';
import { ResumePreview } from '../editor/ResumePreview';

export const PublicResume: React.FC = () => {
  const { slug } = useParams<{ slug: string }>();
  const [data, setData] = useState<ResumeData | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!slug) return;
    fetch(`${API_BASE}/public/resumes/${slug}`)
      .then(async r => {
        if (!r.ok) {
          let txt = '';
          try { txt = await r.text(); } catch {}
          throw new Error(`HTTP ${r.status} ${r.statusText}${txt ? ' - ' + txt : ''}`);
        }
        return r.json();
      })
      .then((res: any) => {
        const mapped: ResumeData = {
          id: slug,
          title: res.Title,
          templateId: res.TemplateID,
          themeConfig: { color: res.ThemeColor, fontFamily: res.ThemeFont, spacing: res.ThemeSpacing },
          lastModified: Date.now(),
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
              dateRange: i.DateRange,
              description: i.Description
            }))
          }))
        };
        setData(mapped);
      })
      .catch((err) => setError(err?.message ? String(err.message) : String(err)));
  }, [slug]);

  if (error) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="bg-white shadow rounded-lg p-6 text-gray-700">
          {error}
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white min-h-screen p-4 md:p-8">
      {data && <ResumePreview data={data} scale={1} />}
    </div>
  );
};

