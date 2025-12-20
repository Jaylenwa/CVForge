import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { ArrowRight, CheckCircle, Cpu, FileCheck, Share2, Briefcase, GraduationCap, PenTool, Code } from 'lucide-react';
import { Button } from '../components/ui/Button';
import { AppRoute } from '../types';
// 后端数据来源
import { useLanguage } from '../contexts/LanguageContext';
import { API_BASE } from '../config';
import { ResumeArtboard } from './editor/ResumePreview';
import { INITIAL_RESUME } from '../services/mockData';
import { Modal } from '../components/ui/Modal';

export const Home: React.FC = () => {
  const navigate = useNavigate();
  const { t } = useLanguage();
  const [popularTemplates, setPopularTemplates] = React.useState<any[]>([]);
  const [previewOpen, setPreviewOpen] = React.useState(false);
  const [previewTemplateId, setPreviewTemplateId] = React.useState<string | null>(null);
  const previewContainerRef = React.useRef<HTMLDivElement | null>(null);
  const [previewScale, setPreviewScale] = React.useState<number | null>(null);
  const previewRafRef = React.useRef<number | null>(null);
  const previewRoRef = React.useRef<ResizeObserver | null>(null);
  React.useEffect(() => {
    (async () => {
      try {
        const res = await fetch(`${API_BASE}/templates`);
        const data = await res.json();
        const items = (data.items || []).slice(0, 4).map((t: any) => ({
          id: t.ExternalID || t.id,
          name: t.Name || t.name,
          category: t.Category || t.category,
          isPremium: t.IsPremium ?? t.isPremium,
        }));
        setPopularTemplates(items);
      } catch {}
    })();
  }, []);

  const quickAccess = [
    { key: 'home.quick.it', icon: <Code size={20} />, query: 'IT' },
    { key: 'home.quick.finance', icon: <Briefcase size={20} />, query: 'Finance' },
    { key: 'home.quick.creative', icon: <PenTool size={20} />, query: 'Creative' },
  ];

  const HomeTemplateCard: React.FC<{ template: any; onUse: () => void; onPreview: () => void }> = ({ template, onUse, onPreview }) => {
    const containerRef = React.useRef<HTMLDivElement | null>(null);
    const rafRef = React.useRef<number | null>(null);
    const roRef = React.useRef<ResizeObserver | null>(null);
    const stableTimerRef = React.useRef<number | null>(null);
    const lastWidthRef = React.useRef<number>(0);
    const initializedRef = React.useRef(false);
    const [scale, setScale] = React.useState<number | null>(null);
    const [ready, setReady] = React.useState(false);
    React.useLayoutEffect(() => {
      const mmToPx = 96 / 25.4;
      const a4w = 210 * mmToPx;
      const scheduleUpdate = () => {
        if (rafRef.current) cancelAnimationFrame(rafRef.current);
        rafRef.current = requestAnimationFrame(() => {
          const el = containerRef.current;
          if (!el) return;
          lastWidthRef.current = el.clientWidth;
          if (stableTimerRef.current) {
            clearTimeout(stableTimerRef.current);
          }
          stableTimerRef.current = window.setTimeout(() => {
            const s = lastWidthRef.current / a4w;
            setScale(prev => (prev === null || Math.abs(prev - s) > 0.002) ? s : prev);
            setReady(true);
          }, 120);
        });
      };
      if (!initializedRef.current) {
        const el = containerRef.current;
        if (el) {
          const s = el.clientWidth / a4w;
          setScale(s);
          setReady(true);
          initializedRef.current = true;
        }
      } else {
        scheduleUpdate();
      }
      const onResize = () => scheduleUpdate();
      window.addEventListener('resize', onResize);
      if (containerRef.current) {
        roRef.current = new ResizeObserver(onResize);
        roRef.current.observe(containerRef.current);
      }
      return () => {
        window.removeEventListener('resize', onResize);
        if (rafRef.current) cancelAnimationFrame(rafRef.current);
        if (stableTimerRef.current) {
          clearTimeout(stableTimerRef.current);
        }
        if (roRef.current) {
          roRef.current.disconnect();
        }
      };
    }, []);
    const mmToPx = 96 / 25.4;
    const a4w = 210 * mmToPx;
    const a4h = 297 * mmToPx;
    return (
      <div className="group relative border border-gray-200 rounded-lg overflow-hidden shadow-sm hover:shadow-lg">
        <div ref={containerRef} className="aspect-[210/297] bg-gray-100 overflow-hidden relative">
          <div className="absolute inset-0 flex items-center justify-center">
            {ready && scale !== null ? (
              <div
                style={{ width: a4w * scale, height: a4h * scale }}
                className="relative select-none pointer-events-none shadow-sm bg-white"
              >
                <ResumeArtboard
                  data={{ ...INITIAL_RESUME, templateId: template.id }}
                  scale={scale}
                  disableShadow={true}
                  style={{ margin: 0 }}
                />
              </div>
            ) : (
              <div className="w-full h-full bg-white" />
            )}
          </div>
          <div className="absolute inset-0 bg-black/0 group-hover:bg-black/20 flex items-center justify-center opacity-0 group-hover:opacity-100">
            <div className="flex flex-col items-center space-y-3">
              <Button className="w-40" onClick={onUse}>{t('home.actions.useTemplate')}</Button>
              <Button className="w-40" variant="outline" onClick={onPreview}>{t('common.preview')}</Button>
            </div>
          </div>
        </div>
        <div className="p-3">
          <h3 className="font-medium text-gray-900 truncate">{template.name}</h3>
          <div className="flex items-center text-xs text-gray-500 mt-1 space-x-2">
            <span className="px-2 py-0.5 bg-gray-100 rounded">{template.category}</span>
            {template.isPremium && <span className="text-yellow-600 font-bold">{t('home.badge.premium')}</span>}
          </div>
        </div>
      </div>
    );
  };

  return (
    <div className="flex flex-col">
      {/* Hero Section */}
      <section className="relative bg-white overflow-hidden">
        <div className="max-w-7xl mx-auto">
          <div className="relative z-10 pb-8 bg-white sm:pb-16 md:pb-20 lg:max-w-2xl lg:w-full lg:pb-28 xl:pb-32 pt-20 px-4 sm:px-6 lg:px-8">
            <main className="mt-10 mx-auto max-w-7xl sm:mt-12 md:mt-16 lg:mt-20 xl:mt-28">
              <div className="sm:text-center lg:text-left">
                <h1 className="text-4xl tracking-tight font-extrabold text-gray-900 sm:text-5xl md:text-6xl">
                  <span className="block xl:inline">{t('hero.title')}</span>{' '}
                  <span className="block text-blue-600 xl:inline">{t('hero.subtitle')}</span>
                </h1>
                <p className="mt-3 text-base text-gray-500 sm:mt-5 sm:text-lg sm:max-w-xl sm:mx-auto md:mt-5 md:text-xl lg:mx-auto">
                  {t('hero.desc')}
                </p>
                <div className="mt-5 sm:mt-8 sm:flex sm:justify-center lg:justify-start">
                  <div className="rounded-md shadow">
                    <Link to={AppRoute.Templates}>
                      <Button size="lg" className="w-full">
                        {t('hero.cta')} <ArrowRight className="ml-2 h-5 w-5" />
                      </Button>
                    </Link>
                  </div>
                  <div className="mt-3 sm:mt-0 sm:ml-3">
                     <Link to={AppRoute.Dashboard}>
                        <Button size="lg" variant="outline" className="w-full">
                            {t('hero.cta_secondary')}
                        </Button>
                    </Link>
                  </div>
                </div>
              </div>
            </main>
          </div>
        </div>
        <div className="lg:absolute lg:inset-y-0 lg:right-0 lg:w-1/2">
          <img
            className="h-56 w-full object-cover sm:h-72 md:h-96 lg:w-full lg:h-full opacity-90"
            src="https://images.unsplash.com/photo-1586281380349-632531db7ed4?ixlib=rb-1.2.1&auto=format&fit=crop&w=1950&q=80"
            alt={t('home.hero.alt')}
          />
        </div>
      </section>

      {/* Quick Access Section */}
      <section className="py-10 bg-gray-50 border-b border-gray-200">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
              <h3 className="text-lg font-semibold text-gray-700 mb-6">{t('home.quickAccess')}</h3>
              <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
                  {quickAccess.map((item) => (
                      <button 
                        key={item.key}
                        onClick={() => navigate(`${AppRoute.Templates}?tag=${item.query}`)}
                        className="flex items-center justify-center p-4 bg-white rounded-lg shadow-sm border border-gray-200 hover:shadow-md hover:border-blue-300 transition-all text-gray-700 hover:text-blue-600"
                      >
                          <span className="mr-3 text-blue-500">{item.icon}</span>
                          <span className="font-medium">{t(item.key)}</span>
                      </button>
                  ))}
              </div>
          </div>
      </section>

      {/* Popular Templates Preview */}
      <section className="py-12 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="flex justify-between items-end mb-6">
                 <h2 className="text-2xl font-bold text-gray-900">{t('home.popular')}</h2>
                 <Link to={AppRoute.Templates} className="text-blue-600 font-medium hover:underline flex items-center">
                    {t('home.actions.viewAll')} <ArrowRight size={16} className="ml-1"/>
                 </Link>
            </div>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
                {popularTemplates.map(template => (
                    <HomeTemplateCard
                      key={template.id}
                      template={template}
                      onUse={() => navigate(`${AppRoute.Editor}?template=${template.id}`)}
                      onPreview={() => { setPreviewTemplateId(template.id); setPreviewOpen(true); }}
                    />
                ))}
            </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-16 bg-gray-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="lg:text-center">
            <h2 className="text-base text-blue-600 font-semibold tracking-wide uppercase">{t('features.title')}</h2>
            <p className="mt-2 text-3xl leading-8 font-extrabold tracking-tight text-gray-900 sm:text-4xl">
              {t('features.subtitle')}
            </p>
          </div>

          <div className="mt-10">
            <dl className="space-y-10 md:space-y-0 md:grid md:grid-cols-2 md:gap-x-8 md:gap-y-10">
              <div className="relative">
                <dt>
                  <div className="absolute flex items-center justify-center h-12 w-12 rounded-md bg-blue-500 text-white">
                    <Cpu size={24} />
                  </div>
                  <p className="ml-16 text-lg leading-6 font-medium text-gray-900">{t('features.ai.title')}</p>
                </dt>
                <dd className="mt-2 ml-16 text-base text-gray-500">
                  {t('features.ai.desc')}
                </dd>
              </div>

              <div className="relative">
                <dt>
                  <div className="absolute flex items-center justify-center h-12 w-12 rounded-md bg-blue-500 text-white">
                    <FileCheck size={24} />
                  </div>
                  <p className="ml-16 text-lg leading-6 font-medium text-gray-900">{t('features.ats.title')}</p>
                </dt>
                <dd className="mt-2 ml-16 text-base text-gray-500">
                  {t('features.ats.desc')}
                </dd>
              </div>

              <div className="relative">
                <dt>
                  <div className="absolute flex items-center justify-center h-12 w-12 rounded-md bg-blue-500 text-white">
                    <Share2 size={24} />
                  </div>
                  <p className="ml-16 text-lg leading-6 font-medium text-gray-900">{t('features.share.title')}</p>
                </dt>
                <dd className="mt-2 ml-16 text-base text-gray-500">
                  {t('features.share.desc')}
                </dd>
              </div>

              <div className="relative">
                <dt>
                  <div className="absolute flex items-center justify-center h-12 w-12 rounded-md bg-blue-500 text-white">
                    <CheckCircle size={24} />
                  </div>
                  <p className="ml-16 text-lg leading-6 font-medium text-gray-900">{t('features.preview.title')}</p>
                </dt>
                <dd className="mt-2 ml-16 text-base text-gray-500">
                  {t('features.preview.desc')}
                </dd>
              </div>
            </dl>
          </div>
        </div>
      </section>
      <Modal isOpen={previewOpen} onClose={() => setPreviewOpen(false)} title={t('common.preview')}>
        <div ref={previewContainerRef} className="aspect-[210/297] bg-gray-100 overflow-hidden relative">
          <div className="absolute inset-0 flex items-center justify-center">
            {previewTemplateId && previewScale !== null ? (
              <div
                style={{ width: (96 / 25.4) * 210 * previewScale, height: (96 / 25.4) * 297 * previewScale }}
                className="relative select-none pointer-events-none shadow-sm bg-white"
              >
                <ResumeArtboard
                  data={{ ...INITIAL_RESUME, templateId: previewTemplateId }}
                  scale={previewScale}
                  disableShadow
                  style={{ margin: 0 }}
                />
              </div>
            ) : (
              <div className="w-full h-full bg-white" />
            )}
          </div>
        </div>
      </Modal>
    </div>
  );
};
