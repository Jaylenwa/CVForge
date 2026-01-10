import React from 'react';
import { useLanguage, LanguageProvider } from '../../contexts/LanguageContext';
import { ResumeData } from '../../types';
import { getThemeStyles } from '../../utils/resume-helpers';

import { TemplateClassic } from '../../components/templates/TemplateClassic';
import { TemplateMintTimeline } from '../../components/templates/TemplateMintTimeline';

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
  const sourceRef = React.useRef<HTMLDivElement>(null);
  const [pageInfo, setPageInfo] = React.useState<{ pageHeight: number; contentHeight: number; count: number }>({ pageHeight: 0, contentHeight: 0, count: 1 });
  const [pageOffsets, setPageOffsets] = React.useState<number[]>([0]);

  const TEMPLATE_COMPONENTS: Record<string, React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }>> = {
    TemplateClassic: TemplateClassic,
    TemplateMintTimeline: TemplateMintTimeline,
  };

  const containerStyle: React.CSSProperties = {
    transform: `scale(${scale})`,
    transformOrigin,
    height: pageInfo.pageHeight > 0 ? `${pageInfo.count * pageInfo.pageHeight}px` : '297mm',
    ...style,
  };

  const measure = React.useCallback(() => {
    const el = sourceRef.current || rootRef.current;
    if (!el) return;
    const widthPx = el.offsetWidth || el.getBoundingClientRect().width || 794; // 210mm ≈ 794px @96dpi
    const pxHeight = (widthPx * 297) / 210;
    const contentHeight = el.scrollHeight;
    const count = Math.max(1, Math.ceil(contentHeight / pxHeight));
    setPageInfo({ pageHeight: pxHeight, contentHeight, count });

    const offsets: number[] = [];
    for (let i = 0; i < count; i++) {
      offsets.push(i * pxHeight);
    }
    setPageOffsets(offsets);
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
    const el = sourceRef.current || rootRef.current;
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
    return <Comp data={data} styles={styles} disableShadow={disableShadow} />;
  };

  return (
      <div 
        id="resume-export-root"
        ref={rootRef}
        className={`relative w-[210mm] min-h-[297mm] print:w-[210mm] print:transform-none bg-white mx-auto box-border ${disableShadow ? 'shadow-none' : 'shadow-md'} print:shadow-none border border-gray-200 print:border-0 ${className}`}
        style={containerStyle}
      >
        {showPageHint && (
          <div className="absolute left-2 top-2 px-2 py-1 text-xs rounded bg-amber-100 text-amber-800 shadow-sm print:hidden">
            {t('preview.pagesHint').replace('{count}', String(pageInfo.count))}
          </div>
        )}

        <div className="print:hidden">
          {Array.from({ length: pageInfo.count }).map((_, i) => {
            const start = pageOffsets[i] ?? i * pageInfo.pageHeight;
            return (
              <div
                key={`page-${i}`}
                className="w-[210mm] h-[297mm] overflow-hidden bg-white mx-auto border border-gray-200 shadow-sm mb-4"
              >
                <div style={{ transform: `translateY(-${start}px)` }}>
                  {renderTemplate()}
                </div>
              </div>
            );
          })}
        </div>

        <div
          ref={sourceRef}
          className="absolute left-0 top-0 w-[210mm] min-h-[297mm] invisible -z-10 pointer-events-none print:visible print:static print:z-auto print:pointer-events-auto"
        >
          <div>
            {renderTemplate()}
          </div>
        </div>
      </div>
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
