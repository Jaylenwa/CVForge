import { ResumeData, ResumeSectionType, Template } from '../types';

export const MOCK_TEMPLATES: Template[] = [
  {
    id: 't1',
    name: 'Classic Professional',
    thumbnail: 'https://images.unsplash.com/photo-1586281380349-632531db7ed4?w=500&q=80',
    tags: ['Professional', 'Simple', 'ATS Friendly'],
    popularity: 98,
    isPremium: false,
    category: 'General',
    level: 'Senior'
  },
  {
    id: 't2',
    name: 'Modern Dark',
    thumbnail: 'https://images.unsplash.com/photo-1616628188859-7a11abb6fcc9?w=500&q=80',
    tags: ['Creative', 'Design', 'Startup'],
    popularity: 85,
    isPremium: true,
    category: 'Creative',
    level: 'Junior'
  },
  {
    id: 't3',
    name: 'Tech Minimalist',
    thumbnail: 'https://images.unsplash.com/photo-1512486130939-2c4f79935e4f?w=500&q=80',
    tags: ['Minimalist', 'Tech', 'Clean'],
    popularity: 92,
    isPremium: false,
    category: 'IT',
    level: 'Junior'
  },
    {
    id: 't4',
    name: 'Executive Serif',
    thumbnail: 'https://images.unsplash.com/photo-1507679799987-c73779587ccf?w=500&q=80',
    tags: ['Professional', 'Management', 'Senior'],
    popularity: 70,
    isPremium: true,
    category: 'Finance',
    level: 'Executive'
  },
  {
    id: 't5',
    name: 'Creative Bold',
    thumbnail: 'https://images.unsplash.com/photo-1513542789411-b6a5d4f31634?w=500&q=80',
    tags: ['Creative', 'Marketing', 'Colorful'],
    popularity: 65,
    isPremium: true,
    category: 'Creative',
    level: 'Senior'
  },
  {
    id: 't6',
    name: 'Elegant Teal',
    thumbnail: 'https://images.unsplash.com/photo-1515378791036-0648a3ef77b2?w=500&q=80',
    tags: ['Modern', 'Fresh', 'Entry Level'],
    popularity: 88,
    isPremium: false,
    category: 'General',
    level: 'Intern'
  }
];

export const INITIAL_RESUME: ResumeData = {
  id: 'new',
  title: 'Untitled Resume',
  templateId: 't1',
  themeConfig: {
    color: '#2563eb',
    fontFamily: 'inter',
    spacing: 'normal'
  },
  lastModified: Date.now(),
  personalInfo: {
    fullName: 'Alex Doe',
    jobTitle: 'Software Engineer',
    email: 'alex.doe@example.com',
    phone: '+1 (555) 123-4567',
    address: 'San Francisco, CA',
    website: 'linkedin.com/in/alexdoe',
    avatarUrl: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?ixlib=rb-1.2.1&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80'
  },
  sections: [
    {
      id: 'summary',
      type: ResumeSectionType.Summary,
      title: 'Professional Summary',
      isVisible: true,
      items: [{
        id: 'sum1',
        description: 'Experienced software engineer with a focus on scalable web applications and React architecture. Proven track record of delivering high-quality code and leading agile teams.'
      }]
    },
    {
      id: 'exp',
      type: ResumeSectionType.Experience,
      title: 'Work Experience',
      isVisible: true,
      items: [
        {
          id: 'exp1',
          title: 'Senior Frontend Developer',
          subtitle: 'Tech Corp Inc.',
          dateRange: '2020 - Present',
          location: 'Remote',
          description: 'Led a team of 5 developers to rebuild the main dashboard using React and TypeScript. Improved load times by 40%.'
        }
      ]
    },
    {
      id: 'edu',
      type: ResumeSectionType.Education,
      title: 'Education',
      isVisible: true,
      items: [
        {
          id: 'edu1',
          title: 'B.S. Computer Science',
          subtitle: 'University of Technology',
          dateRange: '2016 - 2020',
          location: 'New York, NY',
          description: 'Graduated with Honors. Member of the ACM chapter.'
        }
      ]
    },
    {
      id: 'skills',
      type: ResumeSectionType.Skills,
      title: 'Skills',
      isVisible: true,
      items: [
        { id: 's1', description: 'React, TypeScript, Tailwind CSS, Node.js, GraphQL' }
      ]
    }
  ]
};

// Mock user's existing resumes
export const MOCK_USER_RESUMES: ResumeData[] = [
  { 
    ...INITIAL_RESUME, 
    id: '1', 
    title: 'Software Engineer Resume', 
    lastModified: Date.now() - 1000000 
  },
  { 
    ...INITIAL_RESUME, 
    id: '2', 
    title: 'Executive Director CV', 
    templateId: 't4',
    themeConfig: {
        color: '#111827',
        fontFamily: 'merriweather',
        spacing: 'spacious'
    },
    lastModified: Date.now() - 86400000 * 5,
    personalInfo: {
        ...INITIAL_RESUME.personalInfo,
        fullName: 'Sarah Jenkins',
        jobTitle: 'Executive Director'
    }
  }
];