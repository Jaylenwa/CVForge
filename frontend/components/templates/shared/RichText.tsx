import React from 'react';
import { sanitizeHtml } from '../../../utils/resume-helpers';

export type RichTextProps = {
  html?: string;
  className?: string;
  fontSize?: number;
  lineHeight?: number;
};

export const RichText: React.FC<RichTextProps> = ({ html, className, fontSize, lineHeight }) => {
  const safe = String(html || '').trim();
  if (!safe) return null;

  const baseClassName = 'resume-rich-content text-sm leading-relaxed';
  const mergedClassName = className ? `${baseClassName} ${className}` : baseClassName;

  return (
    <div
      className={mergedClassName}
      style={{ fontSize, lineHeight }}
      dangerouslySetInnerHTML={{ __html: sanitizeHtml(safe) }}
    />
  );
};

export default RichText;

