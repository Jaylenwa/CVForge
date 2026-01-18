import React from 'react';

export const TableCard: React.FC<{ children: React.ReactNode; className?: string }> = ({ children, className }) => {
  return (
    <div className={`bg-white border border-gray-100 rounded-2xl overflow-hidden shadow-sm ${className || ''}`}>
      {children}
    </div>
  );
};

