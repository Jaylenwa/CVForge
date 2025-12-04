export enum ResumeSectionType {
  Personal = 'personal',
  Summary = 'summary',
  Experience = 'experience',
  Education = 'education',
  Skills = 'skills',
  Projects = 'projects',
  Custom = 'custom'
}

export interface ResumeItem {
  id: string;
  title?: string;
  subtitle?: string; // Company or Degree
  dateRange?: string;
  location?: string;
  description: string; // HTML or Markdown allowed
}

export interface ResumeSection {
  id: string;
  type: ResumeSectionType;
  title: string;
  items: ResumeItem[];
  isVisible: boolean;
}

export interface ThemeConfig {
  color: string;
  fontFamily: string;
  spacing: 'compact' | 'normal' | 'spacious';
}

export interface ResumeData {
  id: string;
  title: string;
  templateId: string;
  themeConfig: ThemeConfig;
  lastModified: number;
  personalInfo: {
    fullName: string;
    email: string;
    phone: string;
    address: string;
    website: string;
    avatarUrl?: string;
    jobTitle: string;
  };
  sections: ResumeSection[];
}

export interface Template {
  id: string;
  name: string;
  thumbnail: string;
  tags: string[]; // 'Professional', 'Creative', 'ATS'
  popularity: number;
  isPremium: boolean;
  category: 'IT' | 'Finance' | 'Creative' | 'General';
  level: 'Intern' | 'Junior' | 'Senior' | 'Executive';
}

export enum AppRoute {
  Home = '/',
  Templates = '/templates',
  Editor = '/editor',
  Dashboard = '/dashboard',
  Login = '/login',
  Register = '/register',
  Pricing = '/pricing',
  Settings = '/settings',
  Print = '/print'
}

export type Language = 'en' | 'zh';
