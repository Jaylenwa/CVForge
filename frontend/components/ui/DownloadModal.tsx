import React, { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { FileText, Image as ImageIcon, X } from 'lucide-react';
import { LoadingSpinner } from './LoadingSpinner';
import { FaceIDSuccess } from './FaceIDSuccess';
import { useLanguage } from '../../contexts/LanguageContext';

type DownloadType = 'PDF' | 'IMAGE';

enum DownloadState {
  SELECTING = 'SELECTING',
  DOWNLOADING = 'DOWNLOADING',
  SUCCESS = 'SUCCESS',
}

interface DownloadModalProps {
  isOpen: boolean;
  onClose: () => void;
  onExportPDF: () => Promise<void>;
  onExportImage: () => Promise<void>;
  onError?: (err: Error) => void;
}

export const DownloadModal: React.FC<DownloadModalProps> = ({
  isOpen,
  onClose,
  onExportPDF,
  onExportImage,
  onError,
}) => {
  const { t } = useLanguage();
  const [state, setState] = useState<DownloadState>(DownloadState.SELECTING);

  const handleSelect = async (type: DownloadType) => {
    try {
      setState(DownloadState.DOWNLOADING);
      if (type === 'PDF') {
        await onExportPDF();
      } else {
        await onExportImage();
      }
      setState(DownloadState.SUCCESS);
      setTimeout(() => {
        handleReset();
      }, 2000);
    } catch (err: any) {
      if (onError) onError(err instanceof Error ? err : new Error(String(err)));
      handleReset();
    }
  };

  const handleReset = () => {
    onClose();
    setTimeout(() => setState(DownloadState.SELECTING), 300);
  };

  return (
    <AnimatePresence>
      {isOpen && (
        <>
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            onClick={handleReset}
            className="fixed inset-0 bg-black/60 backdrop-blur-sm z-40"
          />

          <div className="fixed inset-0 flex items-center justify-center z-50 pointer-events-none px-4">
            <motion.div
              initial={{ scale: 0.9, opacity: 0, y: 20 }}
              animate={{ scale: 1, opacity: 1, y: 0 }}
              exit={{ scale: 0.9, opacity: 0, y: 20 }}
              transition={{ type: 'spring', damping: 25, stiffness: 300 }}
              className="bg-white border border-gray-200 w-full max-w-sm rounded-2xl overflow-hidden shadow-2xl pointer-events-auto"
            >
              <div className="p-8">
                {state === DownloadState.SELECTING && (
                  <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} exit={{ opacity: 0 }}>
                    <div className="flex justify-between items-start mb-6">
                      <div>
                        <h2 className="text-2xl font-semibold text-gray-900">{t('common.download')}</h2>
                        <p className="text-gray-500 text-sm mt-1">{t('download.selectFormat')}</p>
                      </div>
                      <button onClick={handleReset} aria-label={t('common.close')} className="p-2 hover:bg-gray-100 rounded-full transition-colors">
                        <X className="w-5 h-5 text-gray-500" />
                      </button>
                    </div>

                    <div className="space-y-3">
                      <button
                        onClick={() => handleSelect('PDF')}
                        className="group w-full flex items-center p-4 rounded-2xl bg-gray-50 border border-gray-200 hover:bg-gray-100 hover:border-gray-300 transition-all duration-200 text-left"
                      >
                        <div className="p-3 bg-primary/10 rounded-xl group-hover:bg-primary/20 transition-colors mr-4">
                          <FileText className="w-6 h-6 text-primary" />
                        </div>
                        <div>
                          <div className="text-gray-900 font-medium">{t('editor.export.pdf')}</div>
                          <div className="text-gray-500 text-xs">{t('download.desc.pdf')}</div>
                        </div>
                      </button>
                      <button
                        onClick={() => handleSelect('IMAGE')}
                        className="group w-full flex items-center p-4 rounded-2xl bg-gray-50 border border-gray-200 hover:bg-gray-100 hover:border-gray-300 transition-all duration-200 text-left"
                      >
                        <div className="p-3 bg-primary/10 rounded-xl group-hover:bg-primary/20 transition-colors mr-4">
                          <ImageIcon className="w-6 h-6 text-primary" />
                        </div>
                        <div>
                          <div className="text-gray-900 font-medium">{t('editor.export.png')}</div>
                          <div className="text-gray-500 text-xs">{t('download.desc.png')}</div>
                        </div>
                      </button>
                    </div>
                  </motion.div>
                )}

                {state === DownloadState.DOWNLOADING && (
                  <motion.div key="loading" initial={{ opacity: 0, scale: 0.9 }} animate={{ opacity: 1, scale: 1 }} className="py-12">
                    <LoadingSpinner />
                  </motion.div>
                )}

                {state === DownloadState.SUCCESS && (
                  <motion.div key="success" initial={{ opacity: 0, scale: 0.8 }} animate={{ opacity: 1, scale: 1 }} className="py-10">
                    <FaceIDSuccess />
                  </motion.div>
                )}
              </div>
            </motion.div>
          </div>
        </>
      )}
    </AnimatePresence>
  );
};
