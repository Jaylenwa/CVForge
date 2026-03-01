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
  const [expired, setExpired] = useState<boolean>(false);
  const [requiresPassword, setRequiresPassword] = useState<boolean>(false);
  const [password, setPassword] = useState<string>('');
  const [passwordError, setPasswordError] = useState<string | null>(null);
  const { t, language, setLanguage } = useLanguage();
  const toggleLanguage = () => setLanguage(language === 'en' ? 'zh' : 'en');

  const load = async () => {
    if (!slug) return;
    setError(null);
    setNotPublic(false);
    setExpired(false);
    setRequiresPassword(false);
    const toOrderNum = (v: any): number | undefined => {
      const n = typeof v === 'number' ? v : Number(v);
      return Number.isFinite(n) ? n : undefined;
    };
    try {
      const token = sessionStorage.getItem(`share_token:${slug}`) || '';
      const r = await fetch(`${API_BASE}/public/resumes/${slug}`, {
        headers: token ? { 'X-Share-Token': token } : undefined
      });
      if (!r.ok) {
        if (r.status === 404) {
          setNotPublic(true);
          return;
        }
        if (r.status === 410) {
          setExpired(true);
          return;
        }
        if (r.status === 401) {
          setRequiresPassword(true);
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
        language: (res.Language || '') === 'en' ? 'en' : 'zh',
        Theme: res.Theme,
        lastModified: res.LastModified || Date.now(),
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
    } catch (err: any) {
      setError(err?.message ? String(err.message) : String(err));
    }
  };

  useEffect(() => {
    (async () => {
      await load();
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
  if (expired) {
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
          <h2 className="text-xl font-bold text-slate-900">{t('public.expired.title')}</h2>
          <p className="text-sm text-slate-600">{t('public.expired.desc')}</p>
          <div className="pt-2">
            <a href={`#${AppRoute.Home}`}>
              <Button variant="outline">{t('common.backHome')}</Button>
            </a>
          </div>
        </div>
      </div>
    );
  }
  if (requiresPassword) {
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
          <h2 className="text-xl font-bold text-slate-900">{t('public.password.title')}</h2>
          <p className="text-sm text-slate-600">{t('public.password.desc')}</p>
          <div className="pt-2 space-y-3 text-left" onKeyDown={(e) => {
            if (e.key === 'Enter') {
              (async () => {
                if (!slug) return;
                setPasswordError(null);
                const r = await fetch(`${API_BASE}/public/resumes/${slug}/auth`, {
                  method: 'POST',
                  headers: { 'Content-Type': 'application/json' },
                  body: JSON.stringify({ password })
                });
                if (!r.ok) {
                  setPasswordError(t('public.password.invalid'));
                  return;
                }
                const res = await r.json();
                if (res?.token) {
                  sessionStorage.setItem(`share_token:${slug}`, String(res.token));
                }
                await load();
              })();
            }
          }}>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder={t('public.password.placeholder')}
              className="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all"
            />
            {passwordError && <div className="text-sm text-red-600">{passwordError}</div>}
            <div className="flex items-center gap-2">
              <Button
                onClick={async () => {
                  if (!slug) return;
                  setPasswordError(null);
                  const r = await fetch(`${API_BASE}/public/resumes/${slug}/auth`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ password })
                  });
                  if (!r.ok) {
                    setPasswordError(t('public.password.invalid'));
                    return;
                  }
                  const res = await r.json();
                  if (res?.token) {
                    sessionStorage.setItem(`share_token:${slug}`, String(res.token));
                  }
                  await load();
                }}
              >
                {t('public.password.submit')}
              </Button>
              <a href={`#${AppRoute.Home}`}>
                <Button variant="outline">{t('common.backHome')}</Button>
              </a>
            </div>
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
