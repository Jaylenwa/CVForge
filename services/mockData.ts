import { ResumeData, ResumeSectionType, Template } from '../types';

export const MOCK_TEMPLATES: Template[] = [
  {
    id: 't1',
    name: 'Classic Professional',
    thumbnail: 'https://picsum.photos/300/400?random=1',
    tags: ['Professional', 'ATS Friendly', 'Finance'],
    popularity: 98,
    isPremium: false
  },
  {
    id: 't2',
    name: 'Modern Creative',
    thumbnail: 'https://picsum.photos/300/400?random=2',
    tags: ['Creative', 'Design', 'Startup'],
    popularity: 85,
    isPremium: true
  },
  {
    id: 't3',
    name: 'Tech Minimalist',
    thumbnail: 'https://picsum.photos/300/400?random=3',
    tags: ['Minimalist', 'Tech', 'Engineering'],
    popularity: 92,
    isPremium: false
  },
    {
    id: 't4',
    name: 'Executive Suite',
    thumbnail: 'https://picsum.photos/300/400?random=4',
    tags: ['Professional', 'Management', 'Senior'],
    popularity: 70,
    isPremium: true
  }
];

export const INITIAL_RESUME: ResumeData = {
  id: 'new',
  title: 'Untitled Resume',
  templateId: 't1',
  themeColor: '#2563eb',
  lastModified: Date.now(),
  personalInfo: {
    fullName: 'Alex Doe',
    jobTitle: 'Software Engineer',
    email: 'alex.doe@example.com',
    phone: '+1 (555) 123-4567',
    address: 'San Francisco, CA',
    website: 'linkedin.com/in/alexdoe',
    avatarUrl: ''
  },
  sections: [
    {
      id: 'summary',
      type: ResumeSectionType.Summary,
      title: 'Professional Summary',
      isVisible: true,
      items: [{
        id: 'sum1',
        description: 'Experienced software engineer with a focus on scalable web applications and React architecture.'
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
