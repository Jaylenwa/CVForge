export enum ResumeSectionType {
  Personal = 'personal',
  Summary = 'summary',
  Experience = 'experience',
  Education = 'education',
  Skills = 'skills',
  Projects = 'projects'
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

export interface ResumeData {
  id: string;
  title: string;
  templateId: string;
  themeColor: string;
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
}

export enum AppRoute {
  Home = '/',
  Templates = '/templates',
  Editor = '/editor',
  Dashboard = '/dashboard',
  Login = '/login',
  Pricing = '/pricing'
}
