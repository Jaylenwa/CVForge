import { ResumeData, ResumeSectionType, Template } from '../types';

export const MOCK_TEMPLATES: Template[] = [
  { id: 't1', name: '经典专业版', tags: ['专业', '简洁', 'ATS 友好'], usageCount: 0, isPremium: false, category: 'General' },
  { id: 't2', name: '创意大胆', tags: ['创意', '营销', '色彩'], usageCount: 0, isPremium: true, category: 'Creative' },
  { id: 't3', name: '蓝色简约', tags: ['通用', '中文', 'ATS 友好'], usageCount: 0, isPremium: false, category: 'General' },
  { id: 't4', name: '现代中文', tags: ['现代', '中文', '清爽'], usageCount: 0, isPremium: false, category: 'General' },
  { id: 't5', name: '商务专业', tags: ['商务', '中文', '专业'], usageCount: 0, isPremium: false, category: 'Finance' },
  { id: 't6', name: '中文创意', tags: ['创意', '中文', '设计'], usageCount: 0, isPremium: true, category: 'Creative' },
  { id: 't7', name: '清新绿色', tags: ['清新', '中文', '绿色'], usageCount: 0, isPremium: false, category: 'General' },
  { id: 't8', name: '时间轴信息图', tags: ['时间轴', '中文', '信息图'], usageCount: 0, isPremium: false, category: 'General' },
  { id: 't9', name: '毛笔艺术', tags: ['毛笔', '中文', '艺术'], usageCount: 0, isPremium: true, category: 'Creative' },
  { id: 't10', name: '云朵柔和', tags: ['云端', '中文', '柔和'], usageCount: 0, isPremium: false, category: 'IT' },
  { id: 't11', name: '金融数据', tags: ['金融', '中文', '数据'], usageCount: 0, isPremium: false, category: 'Finance' },
  { id: 't12', name: '游戏趣味', tags: ['游戏', '中文', '趣味'], usageCount: 0, isPremium: false, category: 'IT' },
  { id: 't13', name: '几何图形', tags: ['几何', '中文', '造型'], usageCount: 0, isPremium: false, category: 'General' },
  { id: 't14', name: '医疗洁净', tags: ['医疗', '中文', '洁净'], usageCount: 0, isPremium: false, category: 'General' },
  { id: 't15', name: '折纸艺术', tags: ['折纸', '中文', '纸感'], usageCount: 0, isPremium: true, category: 'Creative' },
  { id: 't16', name: '像素复古', tags: ['像素', '中文', '复古'], usageCount: 0, isPremium: true, category: 'IT' },
  { id: 't17', name: '流线波纹', tags: ['波纹', '中文', '流线'], usageCount: 0, isPremium: false, category: 'Creative' },
  { id: 't18', name: '米色宣纸', tags: ['中国风', '单列', '温润'], usageCount: 0, isPremium: false, category: 'General' },
  { id: 't19', name: '竹韵清雅', tags: ['中国风', '单列', '清雅'], usageCount: 0, isPremium: false, category: 'General' },
  { id: 't20', name: '梅韵简雅', tags: ['中国风', '单列', '简雅'], usageCount: 0, isPremium: false, category: 'General' },
];

export const INITIAL_RESUME: ResumeData = {
  id: 'new',
  title: '我的简历',
  templateId: 't1',
  themeConfig: {
    color: '#000000',
    fontFamily: 'yahei',
    spacing: 'normal'
  },
  lastModified: Date.now(),
  personalInfo: {
    fullName: '张伟',
    email: 'zhangwei@example.com',
    phone: '13800000000',
    avatarUrl: '/avator.avif',
    jobTitle: '高级软件工程师',
    gender: '男',
    age: '28',
    maritalStatus: '未婚',
    politicalStatus: '中共党员',
    birthplace: '江苏苏州',
    ethnicity: '汉族',
    height: '180cm',
    weight: '70kg',
    city: '上海',
    customInfo: [
      { label: '期望城市', value: '上海 / 杭州' }
    ]
  },
  sections: [
    {
      id: 'summary',
      type: ResumeSectionType.Summary,
      title: '个人简介',
      isVisible: true,
      items: [
        {
          id: 's1',
          description: '拥有 5 年以上全栈开发经验，专注于 React/TypeScript/Node.js 与 Go，擅长前后端协同与系统性能优化，具备从 0 到 1 搭建企业级项目的实战能力。熟悉微服务架构、CI/CD 与云原生部署，能有效推动跨团队协作与交付。'
        },
        {
          id: 's2',
          description: '<ul><li>核心页面首屏加载降低 45%，转化率提升 18%</li><li>主导组件库与工程化改造，构建速度提升 35%</li><li>设计并落地日志与监控体系，故障定位时间缩短 60%</li></ul>'
        }
      ]
    },
    {
      id: 'exp',
      type: ResumeSectionType.Experience,
      title: '工作经历',
      isVisible: true,
      items: [
        {
          id: 'e1',
          title: '高级软件工程师',
          subtitle: '某知名互联网公司',
          dateRange: '2021.07 - 至今',
          description: '<ul><li>负责核心前端架构重构，搭建 React 18 + TypeScript + Vite 技术栈与组件库</li><li>引入服务端渲染与按需加载，页面首屏从 3.2s 降至 1.8s</li><li>协同后端落地基于 Go 的微服务，拆分单体应用，稳定性提升</li><li>建设自动化测试与 CI/CD 流程，发布效率提升，故障率下降</li><li>推动数据可观测方案（埋点/日志/监控/告警），显著缩短问题定位时间</li></ul>'
        },
        {
          id: 'e2',
          title: '全栈工程师',
          subtitle: '某 SaaS 初创公司',
          dateRange: '2019.07 - 2021.06',
          description: '<ul><li>独立负责三个核心模块的端到端开发（Web + API + 数据库）</li><li>基于 Node.js/NestJS 实现多租户权限与审计日志，提高安全合规性</li><li>设计并优化数据库索引与缓存策略，接口 P95 延迟下降 40%</li><li>与产品团队合作完善需求与交互，推动版本迭代与用户增长</li></ul>'
        }
      ]
    },
    {
      id: 'edu',
      type: ResumeSectionType.Education,
      title: '教育背景',
      isVisible: true,
      items: [
        {
          id: 'ed1',
          title: '计算机科学与技术（本科）',
          subtitle: '上海交通大学',
          dateRange: '2015 - 2019',
          description: '主修数据结构与算法、操作系统、计算机网络；ACM 校队成员，省级竞赛获奖。'
        },
        {
          id: 'ed2',
          title: '软件工程（硕士）',
          subtitle: '浙江大学',
          dateRange: '2019 - 2021',
          description: '研究方向为 Web 性能优化与前端工程化；以第一作者发表相关论文。'
        }
      ]
    },
    {
      id: 'skills',
      type: ResumeSectionType.Skills,
      title: '技能特长',
      isVisible: true,
      items: [
        { id: 'sk1', description: 'Go / Gin / 微服务' }
      ]
    },
    {
      id: 'projects',
      type: ResumeSectionType.Projects,
      title: '项目经历',
      isVisible: true,
      items: [
        {
          id: 'p1',
          title: 'OpenResume 开源简历生成器',
          subtitle: '开源项目',
          dateRange: '2023',
          description: '<ul><li>设计并实现拖拽编辑、模板切换、PDF 导出等核心功能</li><li>支持多语言与主题配置，适配多行业模板</li><li>开源社区维护，累计 Star 及 Issue 处理</li></ul>'
        },
        {
          id: 'p2',
          title: '智能报表平台',
          subtitle: '企业内部系统',
          dateRange: '2022',
          description: '<ul><li>前端采用可视化搭建与低代码思路，提升报表配置效率</li><li>后端提供统一数据网关与权限体系，保障数据安全</li><li>上线后减少运营成本 30%，报告生成时间从小时级降至分钟级</li></ul>'
        },
        {
          id: 'p3',
          title: '跨境电商多语言官网',
          subtitle: '自由职业',
          dateRange: '2024',
          description: '<ul><li>基于 Next.js 与国际化方案实现中/英/西多语言站点</li><li>接入 Stripe 支付与内容管理系统（CMS），支持活动快速上线</li><li>实现 SEO 友好与性能优化，Google PageSpeed 评分提升至 95+</li></ul>'
        }
      ]
    }
  ]
};
