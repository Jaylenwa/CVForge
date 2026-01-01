import React from 'react';
import { useLanguage } from '../../contexts/LanguageContext';
import { ResumeData, ThemeConfig } from '../../types';
import { getThemeStyles } from '../../utils/resume-helpers';

import { TemplateClassic } from '../../components/templates/base/TemplateClassic';
import { TemplateBold } from '../../components/templates/base/TemplateBold';
import { TemplateCNBlue } from '../../components/templates/cn/TemplateCNBlue';
import { TemplateCNModern } from '../../components/templates/cn/TemplateCNModern';
import { TemplateCNBusiness } from '../../components/templates/cn/TemplateCNBusiness';
import { TemplateCNCreative } from '../../components/templates/cn/TemplateCNCreative';
import { TemplateCNFresh } from '../../components/templates/cn/TemplateCNFresh';
import { TemplateCNTimeline } from '../../components/templates/cn/TemplateCNTimeline';
import { TemplateCNBrush } from '../../components/templates/cn/TemplateCNBrush';
import { TemplateCNCloud } from '../../components/templates/cn/TemplateCNCloud';
import { TemplateCNFinance } from '../../components/templates/cn/TemplateCNFinance';
import { TemplateCNGame } from '../../components/templates/cn/TemplateCNGame';
import { TemplateCNGeometric } from '../../components/templates/cn/TemplateCNGeometric';
import { TemplateCNMedical } from '../../components/templates/cn/TemplateCNMedical';
import { TemplateCNOrigami } from '../../components/templates/cn/TemplateCNOrigami';
import { TemplateCNPixel } from '../../components/templates/cn/TemplateCNPixel';
import { TemplateCNWave } from '../../components/templates/cn/TemplateCNWave';

interface PreviewProps {
  data: ResumeData;
  scale?: number;
  disableShadow?: boolean;
}

export interface ArtboardProps {
  data: ResumeData;
  scale?: number;
  disableShadow?: boolean;
  className?: string;
  style?: React.CSSProperties;
  showPageHint?: boolean;
}

export const ResumeArtboard: React.FC<ArtboardProps> = ({ data, scale = 1, disableShadow = false, className = '', style = {}, showPageHint = true }) => {
  const { t } = useLanguage();
  const styles = getThemeStyles(data.themeConfig);
  const rootRef = React.useRef<HTMLDivElement>(null);
  const [pageInfo, setPageInfo] = React.useState<{ pageHeight: number; contentHeight: number; count: number }>({ pageHeight: 0, contentHeight: 0, count: 1 });

  const TEMPLATE_COMPONENTS: Record<string, React.FC<{ data: ResumeData; styles: any; disableShadow?: boolean }>> = {
    TemplateClassic: TemplateClassic,
    TemplateBold: TemplateBold,
    TemplateCNBlue: TemplateCNBlue,
    TemplateCNModern: TemplateCNModern,
    TemplateCNBusiness: TemplateCNBusiness,
    TemplateCNCreative: TemplateCNCreative,
    TemplateCNFresh: TemplateCNFresh,
    TemplateCNTimeline: TemplateCNTimeline,
    TemplateCNBrush: TemplateCNBrush,
    TemplateCNCloud: TemplateCNCloud,
    TemplateCNFinance: TemplateCNFinance,
    TemplateCNGame: TemplateCNGame,
    TemplateCNGeometric: TemplateCNGeometric,
    TemplateCNMedical: TemplateCNMedical,
    TemplateCNOrigami: TemplateCNOrigami,
    TemplateCNPixel: TemplateCNPixel,
    TemplateCNWave: TemplateCNWave,
  };

  const containerStyle: React.CSSProperties = {
    transform: `scale(${scale})`,
    transformOrigin: 'top left',
    height: pageInfo.pageHeight > 0 ? `${pageInfo.count * pageInfo.pageHeight}px` : '297mm',
    ...style,
  };

  const measure = React.useCallback(() => {
    const el = rootRef.current;
    if (!el) return;
    const mmProbe = document.createElement('div');
    mmProbe.style.position = 'absolute';
    mmProbe.style.visibility = 'hidden';
    mmProbe.style.height = '297mm';
    mmProbe.style.width = '0';
    document.body.appendChild(mmProbe);
    const pxHeight = mmProbe.offsetHeight || 1123;
    document.body.removeChild(mmProbe);
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
    return <Comp data={data} styles={styles} disableShadow={disableShadow} />;
  };

  return (
      <div 
        id="resume-export-root"
        ref={rootRef}
        className={`relative w-[210mm] min-h-[297mm] print:w-[210mm] print:transform-none bg-white mx-auto box-border px-[1.5mm] py-[2mm] print:px-[1.5mm] print:py-[2mm] ${disableShadow ? 'shadow-none' : 'shadow-md'} print:shadow-none border border-gray-200 print:border-0 ${className}`}
        style={containerStyle}
      >
        {showPageHint && (
          <div className="absolute left-2 top-2 px-2 py-1 text-xs rounded bg-amber-100 text-amber-800 shadow-sm print:hidden">
            {t('preview.pagesHint').replace('{count}', String(pageInfo.count))}
          </div>
        )}
        {Array.from({ length: Math.max(0, pageInfo.count - 1) }).map((_, i) => {
          const top = (i + 1) * pageInfo.pageHeight;
          return (
            <div
              key={i}
              className="absolute left-0 right-0 border-b border-dashed border-red-300 print:hidden z-50 flex items-center justify-end pr-2"
              style={{ top: `${top}px`, height: '1px' }}
            >
              <span className="text-xs text-red-400 bg-white px-1">Page {i + 2}</span>
            </div>
          );
        })}
        {renderTemplate()}
      </div>
  );
};

export const ResumePreview: React.FC<PreviewProps> = ({ data, scale = 1, disableShadow = false }) => {
  return (
    <div className={`w-full flex justify-center bg-white p-0 overflow-auto min-h-0 print:p-0 print:bg-white h-full scrollbar-thin scrollbar-thumb-gray-300 ${disableShadow ? 'shadow-none' : ''} print:shadow-none`}>
      <ResumeArtboard data={data} scale={scale} disableShadow={disableShadow} />
    </div>
  );
};
