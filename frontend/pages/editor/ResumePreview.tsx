import React from 'react';
import { useLanguage, LanguageProvider } from '../../contexts/LanguageContext';
import { ResumeData } from '../../types';
import { getThemeStyles } from '../../utils/resume-helpers';
import { AlertCircle } from 'lucide-react';

import { TemplateClassic } from '../../components/templates/TemplateClassic';
import { TemplateMintTimeline } from '../../components/templates/TemplateMintTimeline';
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

  const TEMPLATE_COMPONENTS: Record<string, React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }>> = {
    TemplateClassic: TemplateClassic,
    TemplateMintTimeline: TemplateMintTimeline,
  };

  const containerStyle: React.CSSProperties = {
    transform: `scale(${scale})`,
    transformOrigin,
    height: pageInfo.pageHeight > 0 ? `${pageInfo.pageHeight * pageInfo.count}px` : '297mm',
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

  return (
    <>
      <div 
          id="resume-export-root"
          ref={rootRef}
          className={`relative w-[210mm] min-h-[297mm] print:w-[210mm] print:transform-none bg-white mx-auto box-border ${disableShadow ? 'shadow-none' : 'shadow-md'} print:shadow-none border border-gray-200 print:border-0 ${className}`}
          style={containerStyle}
        >
          {renderTemplate()}

          {showPageHint && pageInfo.count > 1 && pageInfo.pageHeight > 0 && (
            <div className="absolute inset-0 pointer-events-none print:hidden">
              {Array.from({ length: pageInfo.count - 1 }).map((_, idx) => {
                const pageIndex = idx + 1;
                const barHeight = 36;
                const top = pageIndex * pageInfo.pageHeight - barHeight / 2;
                return (
                  <div key={`page-break-${pageIndex}`} className="absolute left-0 w-full" style={{ top }}>
                    <div className="h-9 w-full bg-gradient-to-b from-slate-800/90 to-slate-700/70 shadow-[0_0_24px_rgba(0,0,0,0.35)] flex items-center justify-between px-3">
                      <button
                        type="button"
                        onClick={() => setTipOpen(true)}
                        className="pointer-events-auto inline-flex items-center gap-2 px-3 py-1 rounded-full bg-rose-500 text-white text-xs font-semibold shadow-sm hover:bg-rose-600"
                      >
                        {t('preview.pageBreakTip.badge')}
                        <AlertCircle size={14} className="text-white" />
                      </button>
                      <div className="text-white/90 text-xs font-semibold">
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
          <div className="space-y-4 text-sm text-slate-700">
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
