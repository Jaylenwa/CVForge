import React, { useState, useEffect } from 'react';

interface AvatarProps {
  src?: string;
  name?: string;
  className?: string;
  alt?: string;
}

export const Avatar: React.FC<AvatarProps> = ({ src, name, className = '', alt = 'Avatar' }) => {
  const [imgError, setImgError] = useState(false);

  // Reset error state when src changes
  useEffect(() => {
    setImgError(false);
  }, [src]);

  // If src exists and no error, show image
  if (src && !imgError) {
    return (
      <img 
        src={src} 
        alt={alt} 
        className={`object-cover ${className}`}
        onError={() => setImgError(true)}
      />
    );
  }

  // Fallback to initials
  const initial = name ? name.charAt(0).toUpperCase() : '?';
  
  return (
    <div className={`flex items-center justify-center bg-blue-600 text-white font-medium ${className}`}>
      {initial}
    </div>
  );
};
