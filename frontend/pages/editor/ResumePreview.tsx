import React from 'react';
import { useLanguage, LanguageProvider } from '../../contexts/LanguageContext';
import { ResumeData } from '../../types';
import { getThemeStyles } from '../../utils/resume-helpers';
import { AlertCircle } from 'lucide-react';

import { TemplateClassic } from '../../components/templates/TemplateClassic';
import { TemplateMintTimeline } from '../../components/templates/TemplateMintTimeline';
import { TemplateSlate } from '../../components/templates/TemplateSlate';
import { TemplateMonoBar } from '../../components/templates/TemplateMonoBar';
import { TemplateSidebarLabel } from '../../components/templates/TemplateSidebarLabel';
import { TemplateBlueStripe } from '../../components/templates/TemplateBlueStripe';
import { TemplateDarkHeaderIcons } from '../../components/templates/TemplateDarkHeaderIcons';
import { Modal } from '../../components/ui/Modal';
import { Button } from '../../components/ui/Button';

interface PreviewProps {
  data: ResumeData;
  scale?: number;
  disableShadow?: boolean;
  scrollInside?: boolean;
}

export interface ArtboardProps {
  data: ResumeData;
  scale?: number;
  disableShadow?: boolean;
  className?: string;
  style?: React.CSSProperties;
  showPageHint?: boolean;
  transformOrigin?: string;
}

export const ResumeArtboard: React.FC<ArtboardProps> = ({ data, scale = 1, disableShadow = false, className = '', style = {}, showPageHint = true, transformOrigin = 'top left' }) => {
  const { t } = useLanguage();
  const styles = getThemeStyles(data.Theme);
  const rootRef = React.useRef<HTMLDivElement>(null);
  const [pageInfo, setPageInfo] = React.useState<{ pageHeight: number; contentHeight: number; count: number }>({ pageHeight: 0, contentHeight: 0, count: 1 });
  const [tipOpen, setTipOpen] = React.useState(false);
  const [isPrint, setIsPrint] = React.useState<boolean>(() => {
    if (typeof window === 'undefined') return false;
    const mm = window.matchMedia?.('print');
    const hash = String(window.location?.hash || '');
    const byRoute = hash.includes('/print');
    return byRoute || !!mm?.matches;
  });
  const themeColor = data.Theme?.Color || '#2563eb';

  const hexToRgb = React.useCallback((hex: string): { r: number; g: number; b: number } | null => {
    const raw = String(hex || '').trim().replace(/^#/, '');
    const normalized = raw.length === 3 ? raw.split('').map(c => c + c).join('') : raw;
    if (!/^[0-9a-fA-F]{6}$/.test(normalized)) return null;
    const num = parseInt(normalized, 16);
    return { r: (num >> 16) & 255, g: (num >> 8) & 255, b: num & 255 };
  }, []);

  const rgb = hexToRgb(themeColor) || { r: 37, g: 99, b: 235 };
  const rgba = React.useCallback((alpha: number) => `rgba(${rgb.r}, ${rgb.g}, ${rgb.b}, ${alpha})`, [rgb.r, rgb.g, rgb.b]);
  const isLight = ((rgb.r * 0.2126 + rgb.g * 0.7152 + rgb.b * 0.0722) / 255) > 0.72;
  const pillTextColor = isLight ? '#0f172a' : '#ffffff';

  const TEMPLATE_COMPONENTS: Record<string, React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }>> = {
    TemplateClassic: TemplateClassic,
    TemplateMintTimeline: TemplateMintTimeline,
    TemplateSlate: TemplateSlate,
    TemplateMonoBar: TemplateMonoBar,
    TemplateSidebarLabel: TemplateSidebarLabel,
    TemplateBlueStripe: TemplateBlueStripe,
    TemplateDarkHeaderIcons: TemplateDarkHeaderIcons,
  };

  React.useEffect(() => {
    if (typeof window === 'undefined' || !window.matchMedia) return;
    const mm = window.matchMedia('print');
    const byRoute = String(window.location?.hash || '').includes('/print');
    const onChange = () => setIsPrint(mm.matches || byRoute);
    onChange();

    if (typeof mm.addEventListener === 'function') {
      mm.addEventListener('change', onChange);
      return () => mm.removeEventListener('change', onChange);
    }
    mm.addListener(onChange);
    return () => mm.removeListener(onChange);
  }, []);

  const containerStyle: React.CSSProperties = {
    ...(!isPrint
      ? {
          transform: `scale(${scale})`,
          transformOrigin,
          height: pageInfo.pageHeight > 0 ? `${pageInfo.pageHeight * pageInfo.count}px` : '297mm',
        }
      : {}),
    ...style,
  };

  const measure = React.useCallback(() => {
    const el = rootRef.current;
    if (!el) return;
    const widthPx = el.offsetWidth || el.getBoundingClientRect().width || 794; // 210mm ≈ 794px @96dpi
    const pxHeight = (widthPx * 297) / 210;
    const contentHeight = el.scrollHeight;
    const count = Math.max(1, Math.ceil(contentHeight / pxHeight));
    setPageInfo({ pageHeight: pxHeight, contentHeight, count });
  }, []);

  React.useEffect(() => {
    measure();
  }, [measure, data, scale]);

  React.useEffect(() => {
    const onResize = () => measure();
    window.addEventListener('resize', onResize);
    return () => window.removeEventListener('resize', onResize);
  }, [measure]);
  
  React.useEffect(() => {
    let ro: ResizeObserver | null = null;
    const el = rootRef.current;
    if (el) {
      ro = new ResizeObserver(() => measure());
      ro.observe(el);
      const imgs = Array.from(el.querySelectorAll('img')) as HTMLImageElement[];
      const handlers: Array<{ el: HTMLImageElement; handler: () => void }> = [];
      imgs.forEach(img => {
        const handler = () => measure();
        handlers.push({ el: img, handler });
        img.addEventListener('load', handler, { once: true });
      });
      document.fonts?.ready?.then(() => measure());
      return () => {
        ro && ro.disconnect();
        handlers.forEach(({ el, handler }) => el.removeEventListener('load', handler));
      };
    }
  }, [measure, data]);

  const renderTemplate = () => {
    const Comp = TEMPLATE_COMPONENTS[data.templateId] || TemplateClassic;
    return <Comp data={data} styles={styles} disableShadow={true} />;
  };

  const effectiveShowPageHint = showPageHint && !isPrint;

  return (
    <>
      <div 
          id="resume-export-root"
          ref={rootRef}
          className={`relative w-[210mm] min-h-[297mm] print:w-[210mm] print:transform-none bg-white mx-auto box-border ${disableShadow ? 'shadow-none' : 'shadow-md'} print:shadow-none border border-gray-200 print:border-0 ${className}`}
          style={containerStyle}
        >
          {renderTemplate()}

          {effectiveShowPageHint && pageInfo.count > 1 && pageInfo.pageHeight > 0 && (
            <div className="absolute inset-0 pointer-events-none print:hidden">
              {Array.from({ length: pageInfo.count - 1 }).map((_, idx) => {
                const pageIndex = idx + 1;
                const barHeight = 34;
                const top = pageIndex * pageInfo.pageHeight - barHeight / 2;
                return (
                  <div key={`page-break-${pageIndex}`} className="absolute left-0 w-full" style={{ top }}>
                    <div
                      className="h-[34px] w-full backdrop-blur-[1px] shadow-[0_10px_22px_rgba(15,23,42,0.10)] flex items-center justify-between px-3 border-y"
                      style={{
                        borderColor: rgba(0.28),
                        backgroundColor: isLight ? 'rgba(255,255,255,0.24)' : rgba(0.10),
                      }}
                    >
                      <div className="absolute left-0 top-0 h-full w-1" style={{ backgroundColor: rgba(0.55) }} />
                      <button
                        type="button"
                        onClick={() => setTipOpen(true)}
                        className="pointer-events-auto relative inline-flex items-center gap-2 px-3 py-1 rounded-full text-xs font-semibold shadow-sm border hover:brightness-95 active:brightness-90"
                        style={{ backgroundColor: themeColor, color: pillTextColor, borderColor: rgba(0.18) }}
                      >
                        {t('preview.pageBreakTip.badge')}
                        <AlertCircle size={14} />
                      </button>
                      <div
                        className="relative text-xs font-semibold"
                        style={{ color: isLight ? '#0f172a' : themeColor }}
                      >
                        {pageIndex}/{pageInfo.count}
                      </div>
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </div>

        <Modal isOpen={tipOpen} onClose={() => setTipOpen(false)} title={t('preview.pageBreakTip.modalTitle')}>
          <div className="space-y-4 text-sm text-slate-700 leading-relaxed">
            <p>{t('preview.pageBreakTip.modalDesc')}</p>
            <div className="flex justify-end">
              <Button onClick={() => setTipOpen(false)}>{t('preview.pageBreakTip.close')}</Button>
            </div>
          </div>
        </Modal>
      </>
  );
};

export const ResumePreview: React.FC<PreviewProps> = ({ data, scale = 1, disableShadow = false, scrollInside = true }) => {
  return (
    <div className={`w-full flex justify-center bg-white pt-4 md:pt-6 lg:pt-8 xl:pt-10 ${scrollInside ? 'overflow-auto' : 'overflow-visible'} min-h-0 print:pt-0 print:bg-white h-full scrollbar-thin scrollbar-thumb-gray-300 ${disableShadow ? 'shadow-none' : ''} print:shadow-none`}>
      <LanguageProvider languageOverride={data.language}>
        <ResumeArtboard data={data} scale={scale} disableShadow={disableShadow} transformOrigin="top center" />
      </LanguageProvider>
    </div>
  );
};
