
// Import React to resolve React namespace error for ReactNode type
import React from 'react';

export interface JobSubCategory {
  title: string;
  roles: string[];
}

export interface JobCategory {
  id: string;
  name: string;
  icon: React.ReactNode;
  subCategories?: JobSubCategory[];
}

export interface Template {
  id: number;
  title: string;
  thumbnail: string;
  usageCount: number;
  tag: string;
}