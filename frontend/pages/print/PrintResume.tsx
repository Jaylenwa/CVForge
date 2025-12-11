import React, { useEffect, useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import { ResumeData } from '../../types';
import { API_BASE } from '../../config';
import { ResumePreview } from '../editor/ResumePreview';

export const PrintResume: React.FC = () => {
  const [searchParams] = useSearchParams();
  const [data, setData] = useState<ResumeData | null>(null);

  useEffect(() => {
    const id = searchParams.get('id');
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
            address: res.Address,
            website: res.Website,
            avatarUrl: res.AvatarURL
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
              location: i.Location,
              description: i.Description
            }))
          }))
        };
        setData(mapped);
      });
  }, [searchParams]);

  return (
    <div className="bg-white print:bg-white min-h-screen p-0 m-0">
      {data && <ResumePreview data={data} scale={1} />}
    </div>
  );
};
