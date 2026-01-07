import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { AppRoute, ResumeData } from '../../types';
import { API_BASE } from '../../config';
import { ResumePreview } from '../editor/ResumePreview';
import { useLanguage } from '../../contexts/LanguageContext';
import { Button } from '../../components/ui/Button';
import { Lock, Globe } from 'lucide-react';

export const PublicResume: React.FC = () => {
  const { slug } = useParams<{ slug: string }>();
  const [data, setData] = useState<ResumeData | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [notPublic, setNotPublic] = useState<boolean>(false);
  const { t, language, setLanguage } = useLanguage();
  const toggleLanguage = () => setLanguage(language === 'en' ? 'zh' : 'en');

  useEffect(() => {
    (async () => {
      if (!slug) return;
      try {
        const r = await fetch(`${API_BASE}/public/resumes/${slug}`);
        if (!r.ok) {
          if (r.status === 404) {
            setNotPublic(true);
            return;
          }
          let txt = '';
          try { txt = await r.text(); } catch {}
          throw new Error(`HTTP ${r.status} ${r.statusText}${txt ? ' - ' + txt : ''}`);
        }
        const res = await r.json();
        const mapped: ResumeData = {
          id: slug,
          title: res.Title,
          templateId: res.TemplateID,
          Theme: res.Theme,
          lastModified: res.LastModified || Date.now(),
          Personal: res.Personal,
          Job: res.Job,
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
      } catch (err: any) {
        setError(err?.message ? String(err.message) : String(err));
      }
    })();
  }, [slug]);

  if (notPublic) {
    return (
      <div className="min-h-screen bg-slate-50 flex items-center justify-center p-6">
        <div className="fixed top-3 right-4 z-50">
          <button 
            onClick={toggleLanguage}
            className="p-2 text-gray-500 hover:text-gray-900 focus:outline-none"
            title={t('lang.switchTitle')}
          >
            <div className="flex items-center space-x-1">
              <Globe size={18} />
              <span className="text-sm font-medium">{language === 'en' ? t('lang.en_short') : t('lang.zh_short')}</span>
            </div>
          </button>
        </div>
        <div className="bg-white rounded-2xl border border-slate-100 shadow-sm p-8 max-w-lg w-full text-center space-y-4">
          <div className="w-12 h-12 rounded-xl bg-slate-100 text-slate-600 inline-flex items-center justify-center mx-auto">
            <Lock size={20} />
          </div>
          <h2 className="text-xl font-bold text-slate-900">{t('public.unavailable.title')}</h2>
          <p className="text-sm text-slate-600">{t('public.unavailable.desc')}</p>
          <div className="pt-2">
            <a href={`#${AppRoute.Home}`}>
              <Button variant="outline">{t('common.backHome')}</Button>
            </a>
          </div>
        </div>
      </div>
    );
  }
  if (error) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="fixed top-3 right-4 z-50">
          <button 
            onClick={toggleLanguage}
            className="p-2 text-gray-500 hover:text-gray-900 focus:outline-none"
            title={t('lang.switchTitle')}
          >
            <div className="flex items-center space-x-1">
              <Globe size={18} />
              <span className="text-sm font-medium">{language === 'en' ? t('lang.en_short') : t('lang.zh_short')}</span>
            </div>
          </button>
        </div>
        <div className="bg-white shadow rounded-lg p-6 text-gray-700">
          {error}
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white min-h-screen p-4 md:p-8 overflow-visible">
      <div className="fixed top-3 right-4 z-50">
        <button 
          onClick={toggleLanguage}
          className="p-2 text-gray-500 hover:text-gray-900 focus:outline-none"
          title={t('lang.switchTitle')}
        >
          <div className="flex items-center space-x-1">
            <Globe size={18} />
            <span className="text-sm font-medium">{language === 'en' ? t('lang.en_short') : t('lang.zh_short')}</span>
          </div>
        </button>
      </div>
      {data && <ResumePreview data={data} scale={1} scrollInside={false} />}
    </div>
  );
};
