import React, { createContext, useContext, useState, ReactNode } from 'react';
import { Language } from '../types';

interface Translations {
  [key: string]: {
    [key: string]: string;
  };
}

const translations: Translations = {
  en: {
    'nav.templates': 'Templates',
    'nav.dashboard': 'My Resumes',
    'nav.pricing': 'Pricing',
    'nav.signin': 'Sign In',
    'nav.getStarted': 'Get Started',
    'nav.signout': 'Sign Out',
    'nav.settings': 'Settings',
    'nav.profile': 'Profile',
    'hero.title': 'Craft your perfect resume',
    'hero.subtitle': 'with AI precision',
    'hero.desc': 'Build professional, ATS-friendly resumes in minutes. Choose from our curated templates, use AI to polish your text, and land your dream job.',
    'hero.cta': 'Create Resume',
    'hero.cta_secondary': 'My Resumes',
    'home.popular': 'Popular Templates',
    'home.quickAccess': 'Quick Access',
    'common.edit': 'Edit',
    'common.preview': 'Preview',
    'common.duplicate': 'Duplicate',
    'common.delete': 'Delete',
    'common.rename': 'Rename',
    'common.save': 'Save',
    'common.download': 'Download',
    'editor.content': 'Content',
    'editor.design': 'Design',
    'editor.personal': 'Personal Information',
    'editor.theme': 'Theme Settings',
    'editor.font': 'Font Family',
    'editor.spacing': 'Spacing & Density',
    'editor.color': 'Accent Color',
    'editor.ai_polish': 'AI Polish',
    'dashboard.title': 'My Resumes',
    'dashboard.createNew': 'Create New Resume',
    'templates.title': 'Resume Templates',
    'templates.desc': 'Choose a professionally designed template to get started.',
    'templates.filter.industry': 'Industry',
    'templates.filter.level': 'Level',
    'templates.search': 'Search templates...',
    'footer.rights': 'All rights reserved.',
    // Auth
    'auth.welcome': 'Welcome back',
    'auth.welcomeDesc': 'Enter your details to access your account.',
    'auth.createAccount': 'Create an account',
    'auth.createDesc': 'Enter your email to get started.',
    'auth.email': 'Email Address',
    'auth.password': 'Password',
    'auth.confirmPassword': 'Confirm Password',
    'auth.verificationCode': 'Verification Code',
    'auth.sendCode': 'Send Code',
    'auth.sending': 'Sending...',
    'auth.resend': 'Resend in',
    'auth.login': 'Sign In',
    'auth.register': 'Sign Up',
    'auth.noAccount': "Don't have an account?",
    'auth.hasAccount': 'Already have an account?',
    'auth.forgotPassword': 'Forgot password?',
    'auth.terms': 'By clicking continue, you agree to our Terms of Service and Privacy Policy.',
    'auth.error.invalidEmail': 'Please enter a valid email.',
    'auth.error.passwordMismatch': 'Passwords do not match.',
    'auth.error.invalidCode': 'Invalid verification code.',
    'auth.success.codeSent': 'Verification code sent to your email.',
  },
  zh: {
    'nav.templates': '模板库',
    'nav.dashboard': '我的简历',
    'nav.pricing': '订阅方案',
    'nav.signin': '登录',
    'nav.getStarted': '开始制作',
    'nav.signout': '退出登录',
    'nav.settings': '账号设置',
    'nav.profile': '个人资料',
    'hero.title': '打造您的完美简历',
    'hero.subtitle': 'AI 智能辅助',
    'hero.desc': '几分钟内构建专业、通过 ATS 筛选的简历。选择精选模板，使用 AI 润色文本，助您斩获理想 Offer。',
    'hero.cta': '创建简历',
    'hero.cta_secondary': '我的简历',
    'home.popular': '热门模板',
    'home.quickAccess': '快速入口',
    'common.edit': '编辑',
    'common.preview': '预览',
    'common.duplicate': '创建副本',
    'common.delete': '删除',
    'common.rename': '重命名',
    'common.save': '保存',
    'common.download': '下载',
    'editor.content': '内容编辑',
    'editor.design': '外观设计',
    'editor.personal': '个人信息',
    'editor.theme': '主题设置',
    'editor.font': '字体选择',
    'editor.spacing': '间距与密度',
    'editor.color': '主题颜色',
    'editor.ai_polish': 'AI 润色',
    'dashboard.title': '我的简历',
    'dashboard.createNew': '新建简历',
    'templates.title': '简历模板库',
    'templates.desc': '选择一个专业设计的模板，开始您的求职之旅。',
    'templates.filter.industry': '行业分类',
    'templates.filter.level': '职级经验',
    'templates.search': '搜索模板...',
    'footer.rights': '版权所有。',
    // Auth
    'auth.welcome': '欢迎回来',
    'auth.welcomeDesc': '输入您的账号信息以登录。',
    'auth.createAccount': '创建新账号',
    'auth.createDesc': '输入邮箱以开始注册流程。',
    'auth.email': '电子邮箱',
    'auth.password': '密码',
    'auth.confirmPassword': '确认密码',
    'auth.verificationCode': '验证码',
    'auth.sendCode': '获取验证码',
    'auth.sending': '发送中...',
    'auth.resend': '重新发送',
    'auth.login': '登录',
    'auth.register': '注册',
    'auth.noAccount': "还没有账号？",
    'auth.hasAccount': '已有账号？',
    'auth.forgotPassword': '忘记密码？',
    'auth.terms': '点击注册即表示您同意我们的服务条款和隐私政策。',
    'auth.error.invalidEmail': '请输入有效的电子邮箱地址。',
    'auth.error.passwordMismatch': '两次输入的密码不一致。',
    'auth.error.invalidCode': '验证码错误（测试码：123456）。',
    'auth.success.codeSent': '验证码已发送至您的邮箱。',
  }
};

interface LanguageContextType {
  language: Language;
  setLanguage: (lang: Language) => void;
  t: (key: string) => string;
}

const LanguageContext = createContext<LanguageContextType | undefined>(undefined);

export const LanguageProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [language, setLanguage] = useState<Language>('zh'); // Default to Chinese as requested

  const t = (key: string): string => {
    return translations[language][key] || key;
  };

  return (
    <LanguageContext.Provider value={{ language, setLanguage, t }}>
      {children}
    </LanguageContext.Provider>
  );
};

export const useLanguage = () => {
  const context = useContext(LanguageContext);
  if (!context) {
    throw new Error('useLanguage must be used within a LanguageProvider');
  }
  return context;
};