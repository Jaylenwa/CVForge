export enum ResumeSectionType {
  Personal = 'personal',
  Summary = 'summary',
  Experience = 'experience',
  Education = 'education',
  Skills = 'skills',
  Projects = 'projects',
  Internships = 'internships',
  Portfolio = 'portfolio',
  Awards = 'awards',
  Interests = 'interests',
  Exam = 'exam',
  SelfEvaluation = 'selfEvaluation',
  Custom = 'custom'
}

export interface ResumeItem {
  id: string;
  title?: string;
  subtitle?: string; // Company or Degree
  major?: string;
  degree?: string;
  timeStart?: string; // YYYY-MM
  timeEnd?: string;   // YYYY-MM
  today?: boolean;
  description: string; // HTML or Markdown allowed
  orderNum?: number;
}

export interface ResumeSection {
  id: string;
  type: ResumeSectionType;
  title: string;
  items: ResumeItem[];
  isVisible: boolean;
  orderNum?: number;
}

export interface ResumeData {
  id: string;
  title: string;
  templateId: string;
  lastModified: number;
  language?: Language;
  Personal?: {
    FullName?: string;
    Email?: string;
    Phone?: string;
    AvatarURL?: string;
    Gender?: string;
    Age?: string;
    Degree?: string;
    Job?: string;
    City?: string;
    Money?: string;
    JoinTime?: string;
    CustomInfo?: string;
  };
  Theme?: {
    Color?: string;
    Font?: string;
    Spacing?: string;
    FontSize?: string;
  };
  sections: ResumeSection[];
}

export interface Template {
  id: string;
  name: string;
  tags: string[]; // 'Professional', 'Creative', 'ATS'
  usageCount: number;
  isPremium: boolean;
  category: 'IT' | 'Finance' | 'Creative' | 'General';
}

export enum AppRoute {
  Home = '/',
  Templates = '/templates',
  Public = '/public/:slug',
  Editor = '/editor',
  Dashboard = '/dashboard',
  Login = '/login',
  Register = '/register',
  Pricing = '/pricing',
  Settings = '/settings',
  Print = '/print',
  OAuthCallback = '/oauth/callback',
  Admin = '/admin',
  AdminUsers = '/admin/users',
  AdminResumes = '/admin/resumes',
  AdminTemplates = '/admin/templates',
  AdminShares = '/admin/shares',
  AdminConfigs = '/admin/configs'
}

export type Language = 'en' | 'zh';

export interface SystemConfig {
  key: string;
  value: string;
  description: string;
  type: string;
}

export interface AuthConfig {
  enableEmailVerification: boolean;
  enableWeChatLogin: boolean;
  enableGithubLogin: boolean;
  enablePricingPage: boolean;
  weChatAppID: string;
  githubClientID: string;
}
