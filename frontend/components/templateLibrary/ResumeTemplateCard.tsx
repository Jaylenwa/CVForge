import React, { useLayoutEffect, useRef, useState } from 'react';
import { motion } from 'framer-motion';
import { Star } from 'lucide-react';
import { ResumeArtboard } from '../../pages/editor/ResumePreview';
import { INITIAL_RESUME } from '../../services/mockData';
import { ResumeData } from '../../types';
import { useLanguage } from '../../contexts/LanguageContext';
import { Button } from '../ui/Button';

export const ResumeTemplateCard: React.FC<{
  title: string;
  templateId: string;
  usageCount?: number;
  isPremium?: boolean;
  tag?: string;
  presetData?: Partial<ResumeData> | null;
  onUse: () => void;
  onPreview: () => void;
}> = ({ title, templateId, usageCount = 0, isPremium = false, tag, presetData, onUse, onPreview }) => {
  const { t } = useLanguage();
  const containerRef = useRef<HTMLDivElement | null>(null);
  const rafRef = useRef<number | null>(null);
  const roRef = useRef<ResizeObserver | null>(null);
  const stableTimerRef = useRef<number | null>(null);
  const lastWidthRef = useRef<number>(0);
  const initializedRef = useRef(false);
  const [scale, setScale] = useState<number | null>(null);
  const [ready, setReady] = useState(false);

  useLayoutEffect(() => {
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
          setScale((prev) => (prev === null || Math.abs(prev - s) > 0.002 ? s : prev));
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
      if (stableTimerRef.current) clearTimeout(stableTimerRef.current);
      if (roRef.current) roRef.current.disconnect();
    };
  }, []);

  const mmToPx = 96 / 25.4;
  const a4w = 210 * mmToPx;
  const a4h = 297 * mmToPx;
  const previewData = presetData ? ({ ...INITIAL_RESUME, ...presetData, templateId } as any) : ({ ...INITIAL_RESUME, templateId } as any);

  return (
    <motion.div
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      whileHover={{ y: -6 }}
      transition={{ duration: 0.001, ease: 'easeOut' }}
      className="group relative bg-white rounded-xl border border-slate-200 overflow-hidden hover:shadow-xl transition-all duration-200"
    >
      <div ref={containerRef} className="relative aspect-[210/297] overflow-hidden bg-slate-100">
        {ready && scale !== null ? (
          <div className="absolute inset-0 flex items-center justify-center">
            <div
              style={{ width: a4w * scale, height: a4h * scale }}
              className="relative select-none pointer-events-none bg-white"
            >
              <ResumeArtboard data={previewData} scale={scale} disableShadow={true} showPageHint={false} style={{ margin: 0 }} />
            </div>
          </div>
        ) : (
          <div className="absolute inset-0 bg-white" />
        )}

        <div className="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity duration-200 flex items-center justify-center">
          <div className="flex flex-col items-center space-y-3">
            <Button className="w-40" onClick={onUse}>
              {t('templates.actions.useTemplate')}
            </Button>
            <Button className="w-40" variant="outline" onClick={onPreview}>
              {t('common.preview')}
            </Button>
          </div>
        </div>

        {isPremium ? (
          <div className="absolute top-3 right-3 bg-amber-400 text-amber-900 text-xs font-bold px-2 py-1 rounded-full flex items-center shadow-md">
            <Star size={12} className="mr-1 fill-current" /> {t('templates.badge.premium')}
          </div>
        ) : null}
      </div>

      <div className="p-5">
        <h3 className="text-lg font-bold text-slate-800 mb-1 line-clamp-1">{title}</h3>
        <div className="flex items-center justify-between mt-3">
          <span className="text-sm text-slate-500">{Number(usageCount || 0).toLocaleString()} {t('templates.meta.usageCount')}</span>
          {tag ? (
            <span className="px-2.5 py-1 text-xs font-medium bg-blue-50 text-blue-600 rounded-md border border-blue-100">
              {tag}
            </span>
          ) : null}
        </div>
      </div>
    </motion.div>
  );
};
