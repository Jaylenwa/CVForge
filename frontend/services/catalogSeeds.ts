import { ResumeData, ResumeSectionType } from '../types';

export type JobCategorySeed = {
  id: string;
  name: string;
  parentId?: string;
  orderNum?: number;
};

export type JobRoleSeed = {
  id: string;
  categoryId: string;
  name: string;
  tags?: string[];
  orderNum?: number;
};

export type ContentPresetSeed = {
  id: string;
  name: string;
  language: 'zh' | 'en';
  roleId: string;
  data: Partial<ResumeData>;
  tags?: string[];
};

export type TemplateVariantSeed = {
  id: string;
  name: string;
  layoutTemplateId: string;
  presetId: string;
  roleId: string;
  tags?: string[];
  usageCount?: number;
  isPremium?: boolean;
};

export const JOB_CATEGORIES_SEED: JobCategorySeed[] = [
  { id: 'it', name: 'IT | 互联网', orderNum: 10 },
  { id: 'it_backend', name: '后端开发/程序员', parentId: 'it', orderNum: 10 },
];

export const JOB_ROLES_SEED: JobRoleSeed[] = [
  { id: 'java', categoryId: 'it_backend', name: 'Java', tags: ['Java', '后端'], orderNum: 10 },
];

export const CONTENT_PRESETS_SEED: ContentPresetSeed[] = [
  {
    id: 'it_backend_java_zh',
    name: 'Java 开发（中文示例）',
    language: 'zh',
    roleId: 'java',
    tags: ['Java', '后端', 'Spring'],
    data: {
      title: 'Java 开发工程师简历',
      Personal: {
        FullName: '李雷',
        Email: 'lilei@example.com',
        Phone: '13800000001',
        Job: 'Java 开发工程师',
        City: '北京 / 上海',
        Money: '25k-35k·14薪',
        JoinTime: '1 个月内',
        Degree: '本科',
        CustomInfo: JSON.stringify([{ label: '技术栈', value: 'Java / Spring / MySQL / Redis' }]),
      } as any,
      sections: [
        {
          id: 'exp',
          type: ResumeSectionType.Experience,
          title: '工作经历',
          isVisible: true,
          items: [
            {
              id: 'e1',
              title: '某大型互联网公司',
              subtitle: 'Java 开发工程师（支付/交易）',
              timeStart: '2021-06',
              today: true,
              description:
                '<ul><li>负责交易链路核心接口开发与重构，推动单体拆分为可复用服务，接口 P95 延迟降低 30%</li><li>基于 Spring Boot + MyBatis 构建高并发 API，完善限流/熔断/降级策略，显著提升稳定性</li><li>落地 Redis 缓存与异步消息（MQ）削峰填谷，日均处理千万级请求</li><li>建设 SQL 慢查询治理与索引优化，核心查询耗时从 600ms 降至 120ms</li></ul>',
            },
          ],
        },
        {
          id: 'proj',
          type: ResumeSectionType.Projects,
          title: '项目经历',
          isVisible: true,
          items: [
            {
              id: 'p1',
              title: '订单中心重构',
              subtitle: 'Java / Spring Cloud / MySQL / Redis',
              timeStart: '2022-03',
              timeEnd: '2022-12',
              description:
                '<ul><li>设计领域模型与接口分层，拆分订单、库存、优惠等子域，实现职责清晰与可演进</li><li>引入幂等、分布式锁、补偿机制，保障高并发下的数据一致性</li><li>完善链路追踪与监控告警，问题定位效率提升</li></ul>',
            },
          ],
        },
        {
          id: 'skills',
          type: ResumeSectionType.Skills,
          title: '技能清单',
          isVisible: true,
          items: [
            {
              id: 's1',
              title: '后端',
              subtitle: 'Java',
              description: 'Spring Boot / Spring Cloud / MyBatis / JPA',
            },
            {
              id: 's2',
              title: '数据库与缓存',
              subtitle: 'MySQL / Redis',
              description: '索引优化、事务与锁、缓存一致性、热点治理',
            },
            {
              id: 's3',
              title: '工程化',
              subtitle: 'CI/CD',
              description: 'Docker、日志与监控、性能压测与容量评估',
            },
          ],
        },
      ],
    },
  },
];

export const TEMPLATE_VARIANTS_SEED: TemplateVariantSeed[] = [
  {
    id: 'mint_it_backend_java_zh',
    name: '青色时间轴 - Java 开发',
    layoutTemplateId: 'TemplateMintTimeline',
    presetId: 'it_backend_java_zh',
    roleId: 'java',
    tags: ['Java', '后端', '中文'],
    usageCount: 0,
    isPremium: false,
  },
  {
    id: 'classic_it_backend_java_zh',
    name: '经典专业版 - Java 开发',
    layoutTemplateId: 'TemplateClassic',
    presetId: 'it_backend_java_zh',
    roleId: 'java',
    tags: ['Java', '后端', '中文'],
    usageCount: 0,
    isPremium: false,
  },
];
