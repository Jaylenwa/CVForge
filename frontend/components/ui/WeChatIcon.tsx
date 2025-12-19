import React from 'react';

export const WeChatIcon: React.FC<{ className?: string; size?: number }> = ({ className, size = 20 }) => {
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 64 64"
      xmlns="http://www.w3.org/2000/svg"
      className={className}
    >
      <rect x="0" y="0" width="64" height="64" rx="12" fill="#1AAD19" />
      <circle cx="28" cy="26" r="16" fill="#FFFFFF" />
      <path d="M18 40 L12 52 L25 43 Z" fill="#FFFFFF" />
      <circle cx="46" cy="38" r="12" fill="#FFFFFF" />
      <path d="M50 46 L62 50 L54 54 Z" fill="#FFFFFF" />
      <circle cx="22" cy="24" r="2.2" fill="#1AAD19" />
      <circle cx="34" cy="24" r="2.2" fill="#1AAD19" />
      <circle cx="42" cy="36" r="2" fill="#1AAD19" />
      <circle cx="50" cy="36" r="2" fill="#1AAD19" />
    </svg>
  );
};
