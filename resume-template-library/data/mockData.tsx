
import React from 'react';
import { Briefcase, Code, Palette, Settings, User } from 'lucide-react';
import { JobCategory, Template } from '../types';

export const JOB_CATEGORIES: JobCategory[] = [
  { id: 'all', name: '全部', icon: <User size={18} /> },
  { 
    id: 'tech', 
    name: '技术 / 研发', 
    icon: <Code size={18} />,
    subCategories: [
      { title: '后端开发', roles: ['Java', 'Python', 'Go', 'PHP', 'C++', 'C#', '.NET'] },
      { title: '前端/移动端', roles: ['前端开发', 'iOS', 'Android', 'Flutter', 'React Native'] },
      { title: '数据/人工智能', roles: ['算法', '数据分析', '大数据', '深度学习', '自然语言处理'] },
      { title: '测试/运维', roles: ['自动化测试', '性能测试', '运维工程师', 'SRE', '安全工程师'] }
    ]
  },
  { 
    id: 'product', 
    name: '产品 / 设计', 
    icon: <Palette size={18} />,
    subCategories: [
      { title: '产品', roles: ['产品助理', '产品专员', '产品经理', '产品总监', '其他产品岗'] },
      { title: '设计师', roles: ['设计师助理', '平面设计师', 'UI设计师/交互设计', '室内设计师', '建筑设计师', '服装设计师', '网页设计师', '插画师', '原画师', '3D设计师'] },
      { title: '运营', roles: ['新媒体运营', '产品运营', '电商运营', '内容运营', '用户运营', '游戏运营', '短视频运营', '运营总监'] }
    ]
  },
  { 
    id: 'ops', 
    name: '运维 / 技术支持', 
    icon: <Settings size={18} />,
    subCategories: [
      { title: '技术支持', roles: ['售前技术支持', '售后技术支持', 'IT支持', '系统管理员'] },
      { title: '实施维护', roles: ['实施顾问', 'ERP维护', '网络安全运维'] }
    ]
  },
  { 
    id: 'business', 
    name: '市场 / 商务', 
    icon: <Briefcase size={18} />,
    subCategories: [
      { title: '市场/品牌', roles: ['市场营销/推广', '市场专员', '市场经理', '市场总监', '商务合作/BD', 'SEO/SEM', '品牌公关', '品牌策划'] },
      { title: '销售', roles: ['销售代表', '销售专员', '客户经理', '渠道经理', '销售顾问', '电话销售'] }
    ]
  }
];

export const TEMPLATES: Template[] = [
  {
    id: 1,
    title: '青色时间轴',
    thumbnail: 'https://picsum.photos/seed/resume1/400/560',
    usageCount: 1250,
    tag: 'General'
  },
  {
    id: 2,
    title: '经典专业版',
    thumbnail: 'https://picsum.photos/seed/resume2/400/560',
    usageCount: 3400,
    tag: 'General'
  },
  {
    id: 3,
    title: '极简商务风',
    thumbnail: 'https://picsum.photos/seed/resume3/400/560',
    usageCount: 890,
    tag: 'Pro'
  },
  {
    id: 4,
    title: '创意设计型',
    thumbnail: 'https://picsum.photos/seed/resume4/400/560',
    usageCount: 2100,
    tag: 'Creative'
  },
  {
    id: 5,
    title: '应届生校招款',
    thumbnail: 'https://picsum.photos/seed/resume5/400/560',
    usageCount: 15600,
    tag: 'Student'
  },
  {
    id: 6,
    title: '资深主管模版',
    thumbnail: 'https://picsum.photos/seed/resume6/400/560',
    usageCount: 450,
    tag: 'Senior'
  }
];
