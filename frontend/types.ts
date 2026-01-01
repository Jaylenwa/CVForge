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
  JobApplication = 'jobApplication',
  Exam = 'exam',
  SelfEvaluation = 'selfEvaluation',
  Custom = 'custom'
}

export interface ResumeItem {
  id: string;
  title?: string;
  subtitle?: string; // Company or Degree
  timeStart?: string; // YYYY-MM
  timeEnd?: string;   // YYYY-MM
  today?: boolean;
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
    city?: string;
    avatarUrl?: string;
    jobTitle: string;
    gender?: string;
    age?: string;
    maritalStatus?: string;
    politicalStatus?: string;
    birthplace?: string;
    ethnicity?: string;
    height?: string;
    weight?: string;
    customInfo?: Array<{ label: string; value: string }>;
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
  weChatAppID: string;
  githubClientID: string;
}
