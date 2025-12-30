import React from 'react';
import { motion } from 'framer-motion';
import { useLanguage } from '../../contexts/LanguageContext';

export const FaceIDSuccess: React.FC = () => {
  const { t } = useLanguage();
  return (
    <div className="flex flex-col items-center justify-center space-y-4">
      <div className="relative w-24 h-24">
        <motion.svg viewBox="0 0 100 100" className="w-full h-full transform -rotate-90">
          <motion.circle
            cx="50"
            cy="50"
            r="45"
            fill="transparent"
            stroke="#2563eb"
            strokeWidth="6"
            strokeLinecap="round"
            initial={{ pathLength: 0, opacity: 0 }}
            animate={{ pathLength: 1, opacity: 1 }}
            transition={{ duration: 0.6, ease: 'easeOut' }}
          />
        </motion.svg>
        <div className="absolute inset-0 flex items-center justify-center">
          <motion.svg viewBox="0 0 50 50" className="w-12 h-12">
            <motion.path
              d="M10 25 L20 35 L40 15"
              fill="transparent"
              stroke="#2563eb"
              strokeWidth="5"
              strokeLinecap="round"
              strokeLinejoin="round"
              initial={{ pathLength: 0 }}
              animate={{ pathLength: 1 }}
              transition={{ delay: 0.4, duration: 0.4, ease: 'easeInOut' }}
            />
          </motion.svg>
        </div>
      </div>
      <motion.p
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.6 }}
        className="text-primary font-medium text-lg"
      >
        {t('download.success')}
      </motion.p>
    </div>
  );
};
