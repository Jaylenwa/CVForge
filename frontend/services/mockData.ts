import { ResumeData, ResumeSectionType, Template } from '../types';

export const MOCK_TEMPLATES: Template[] = [
  { id: 'TemplateMintTimeline', name: '青色时间轴', tags: ['美观', '中文', 'ATS 友好'] },
  { id: 'TemplateWaveTimeline', name: '波浪时间轴', tags: ['美观', '中文', '时间轴'] },
  { id: 'TemplateBluePhotoColumns', name: '蓝头头像三栏版', tags: ['中文', '三栏', '图标'] },
  { id: 'TemplateClassic', name: '经典专业版', tags: ['专业', '简洁', 'ATS 友好'] },
  { id: 'TemplateSlate', name: '石板简洁版', tags: ['单栏', '简洁', 'ATS 友好'] },
  { id: 'TemplateBlueStripe', name: '蓝条灰底版', tags: ['中文', '清晰', 'ATS 友好'] },
  { id: 'TemplateDarkHeaderIcons', name: '深色头图标版', tags: ['中文', '图标', 'ATS 友好'] },
  { id: 'TemplateMonoBar', name: '黑白竖线版', tags: ['极简', '中文', 'ATS 友好'] },
  { id: 'TemplateSidebarLabel', name: '左栏标签版', tags: ['边框', '中文', '图标'] },
];

export const INITIAL_RESUME: ResumeData = {
  id: 'new',
  title: '我的简历',
  templateId: 'TemplateMintTimeline',
  language: 'zh',
  Theme: {
    Color: '#14b8a6',
    Font: 'notosans',
    Spacing: 'normal',
    FontSize: '13'
  },
  lastModified: Date.now(),
  Personal: {
    FullName: '张伟',
    Email: 'zhangwei@example.com',
    Phone: '13800000000',
    AvatarURL: '/avator.avif',
    Job: '前端/全栈工程师',
    City: '上海 / 杭州',
    Money: '25k-35k·14薪',
    JoinTime: '随时到岗',
    Gender: '男',
    Age: '28',
    Degree: '硕士',
    CustomInfo: JSON.stringify([{ label: '政治面貌', value: '党员' }])
  },
  sections: [
    {
      id: 'exam',
      type: ResumeSectionType.Exam,
      title: '报考信息',
      isVisible: true,
      items: [
        { id: 'ex_meta', title: '浙江大学', subtitle: '计算机科学与技术', description: '考研成绩' },
        { id: 'ex_s1', title: '政治', subtitle: '70', description: '' },
        { id: 'ex_s2', title: '英语', subtitle: '78', description: '' },
        { id: 'ex_s3', title: '数学', subtitle: '120', description: '' },
        { id: 'ex_s4', title: '专业课', subtitle: '130', description: '' }
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
          title: '某知名互联网公司',
          subtitle: '高级软件工程师',
          timeStart: '2021-07',
          today: true,
          description: '<ul><li>负责核心前端架构重构，搭建 React 18 + TypeScript + Vite 技术栈与组件库</li><li>引入服务端渲染与按需加载，页面首屏从 3.2s 降至 1.8s</li><li>协同后端落地基于 Go 的微服务，拆分单体应用，稳定性提升</li><li>建设自动化测试与 CI/CD 流程，发布效率提升，故障率下降</li><li>推动数据可观测方案（埋点/日志/监控/告警），显著缩短问题定位时间</li></ul>'
        },
        {
          id: 'e2',
          title: '某 SaaS 初创公司',
          subtitle: '全栈工程师',
          timeStart: '2019-07',
          timeEnd: '2021-06',
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
          title: '上海交通大学',
          major: '计算机科学与技术',
          degree: '本科',
          timeStart: '2015-09',
          timeEnd: '2019-07',
          description: '主修数据结构与算法、操作系统、计算机网络；ACM 校队成员，省级竞赛获奖。'
        },
        {
          id: 'ed2',
          title: '浙江大学',
          major: '软件工程',
          degree: '硕士',
          timeStart: '2019-09',
          timeEnd: '2021-07',
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
          subtitle: '核心开发者',
          timeStart: '2023-01',
          timeEnd: '2023-12',
          description: '<ul><li>设计并实现拖拽编辑、模板切换、PDF 导出等核心功能</li><li>支持多语言与主题配置，适配多行业模板</li><li>开源社区维护，累计 Star 及 Issue 处理</li></ul>'
        },
        {
          id: 'p2',
          title: '智能报表平台',
          subtitle: '前端负责人',
          timeStart: '2022-01',
          timeEnd: '2022-12',
          description: '<ul><li>前端采用可视化搭建与低代码思路，提升报表配置效率</li><li>后端提供统一数据网关与权限体系，保障数据安全</li><li>上线后减少运营成本 30%，报告生成时间从小时级降至分钟级</li></ul>'
        },
        {
          id: 'p3',
          title: '跨境电商多语言官网',
          subtitle: '独立开发者',
          timeStart: '2024-01',
          timeEnd: '2024-12',
          description: '<ul><li>基于 Next.js 与国际化方案实现中/英/西多语言站点</li><li>接入 Stripe 支付与内容管理系统（CMS），支持活动快速上线</li><li>实现 SEO 友好与性能优化，Google PageSpeed 评分提升至 95+</li></ul>'
        }
      ]
    },
    {
      id: 'intern',
      type: ResumeSectionType.Internships,
      title: '实习经验',
      isVisible: true,
      items: [
        {
          id: 'in1',
          title: '某大型互联网公司',
          subtitle: '前端实习生',
          timeStart: '2020-07',
          timeEnd: '2020-09',
          description: '<ul><li>参与业务组件开发与单元测试，修复若干 UI 缺陷</li><li>协助搭建内部文档站点与脚手架使用指引</li></ul>'
        }
      ]
    },
    {
      id: 'portfolio',
      type: ResumeSectionType.Portfolio,
      title: '个人作品',
      isVisible: true,
      items: [
        {
          id: 'pf1',
          description: '<ul><li>GitHub：<a href="https://github.com/zhangwei">github.com/zhangwei</a></li><li>技术博客：<a href="https://blog.example.com">blog.example.com</a></li></ul>'
        }
      ]
    },
    {
      id: 'awards',
      type: ResumeSectionType.Awards,
      title: '荣誉证书',
      isVisible: true,
      items: [
        {
          id: 'a1',
          description: '<ul><li>国家奖学金（2018）</li><li>ACM 省赛二等奖（2017）</li></ul>'
        }
      ]
    },
    {
      id: 'self',
      type: ResumeSectionType.SelfEvaluation,
      title: '自我评价',
      isVisible: true,
      items: [
        { 
          id: 'self1', 
          description: '责任心强，沟通协作能力好，持续学习新技术，具备独立解决问题与推动项目落地的能力。' 
        }
      ]
    },
    {
      id: 'interests',
      type: ResumeSectionType.Interests,
      title: '兴趣爱好',
      isVisible: true,
      items: [
        { id: 'int1', description: '阅读 / 摄影 / 跑步 / 桌游' }
      ]
    }
  ]
};
