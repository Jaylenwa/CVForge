package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var javaPresetJSON = map[string]any{
	"title":    "Java 后端开发工程师简历",
	"language": "zh",
	"Personal": map[string]any{
		"FullName":   "张伟",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "Java 后端开发工程师",
		"City":       "上海 / 杭州",
		"Money":      "25k-35k·14薪",
		"JoinTime":   "随时到岗",
		"Gender":     "男",
		"Age":        "28",
		"Degree":     "硕士",
		"CustomInfo": `[{"label":"政治面貌","value":"党员"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#f97316",
		"Font":     "notosans",
		"Spacing":  "normal",
		"FontSize": "13",
	},
	"sections": []any{
		map[string]any{
			"id":        "summary",
			"type":      common.ResumeSectionTypeSummary,
			"title":     "优势概述",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "sum1",
					"description": "<ul><li>5 年 Java 后端开发经验，熟悉 Spring Boot / Spring Cloud，具备高并发场景性能优化与稳定性建设经验</li><li>熟悉 MySQL/Redis/Kafka，掌握缓存一致性、消息可靠投递、分布式事务与幂等设计</li><li>具备工程化与云原生实践：Docker/Kubernetes、CI/CD、指标/日志/链路追踪</li></ul>",
				},
			},
		},
		map[string]any{
			"id":        "exp",
			"type":      common.ResumeSectionTypeExperience,
			"title":     "工作经历",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "e1",
					"title":       "某知名互联网公司",
					"subtitle":    "Java 后端开发工程师（交易/支付）",
					"timeStart":   "2021-07",
					"today":       true,
					"description": "<ul><li>负责交易域核心接口与异步链路（Spring Boot + Kafka），订单创建成功率提升至 99.95%</li><li>落地 Redis 多级缓存与热点治理，核心接口 P95 从 260ms 降至 150ms</li><li>建设幂等、重试与降级策略，显著降低高峰期异常请求比例</li><li>完善可观测体系（Prometheus/Grafana + 日志/链路追踪），故障定位时间缩短 50%</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "某企业服务 SaaS 公司",
					"subtitle":    "Java 工程师（平台/中台）",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-06",
					"description": "<ul><li>参与统一权限与审计平台建设（Spring Boot + MySQL），支撑 100+ 企业客户</li><li>优化慢 SQL 与索引策略，核心查询耗时下降 60%</li><li>引入灰度发布与回滚机制，发布成功率提升并降低业务中断风险</li></ul>",
				},
			},
		},
		map[string]any{
			"id":        "edu",
			"type":      common.ResumeSectionTypeEducation,
			"title":     "教育背景",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "ed1",
					"title":       "浙江大学",
					"major":       "软件工程",
					"degree":      "硕士",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "研究方向为分布式系统与高可用架构；参与实验室平台项目与性能优化。",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "浙江大学",
					"major":       "软件工程",
					"degree":      "本科",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "主修数据结构、操作系统、数据库、计算机网络；参与校内项目与竞赛。",
				},
			},
		},
		map[string]any{
			"id":        "skills",
			"type":      common.ResumeSectionTypeSkills,
			"title":     "技能清单",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "s1",
					"title":       "语言与框架",
					"subtitle":    "Java",
					"description": "Spring Boot / Spring Cloud / MyBatis / REST API",
				},
				map[string]any{
					"id":          "s2",
					"title":       "存储与中间件",
					"subtitle":    "数据库 / 缓存 / 消息",
					"description": "MySQL / Redis / Kafka / Elasticsearch",
				},
				map[string]any{
					"id":          "s3",
					"title":       "工程化与稳定性",
					"subtitle":    "交付与可观测",
					"description": "Docker / Kubernetes / CI/CD / Prometheus / Grafana",
				},
			},
		},
		map[string]any{
			"id":        "projects",
			"type":      common.ResumeSectionTypeProjects,
			"title":     "项目经历",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "p1",
					"title":       "统一订单与库存中台",
					"subtitle":    "Java 后端开发工程师",
					"timeStart":   "2022-05",
					"timeEnd":     "2023-02",
					"description": "<ul><li>技术栈：Spring Boot / MySQL / Redis / Kafka</li><li>抽象订单/库存核心领域模型与接口，支持多业务线复用并降低重复开发成本</li><li>引入事件驱动与最终一致性方案，降低跨服务耦合并提升可演进性</li><li>建设压测与容量评估体系，峰值 QPS 提升 2 倍且稳定运行</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "高并发优惠券系统",
					"subtitle":    "Java 后端开发工程师",
					"timeStart":   "2020-09",
					"timeEnd":     "2021-03",
					"description": "<ul><li>技术栈：Java / Redis / Kafka</li><li>设计库存扣减与发放链路幂等方案，避免重复领取与超卖问题</li><li>采用 Redis 预热与 Lua 原子脚本，核心接口 P95 延迟下降 40%</li><li>引入异步削峰与失败重试机制，保障大促场景稳定性</li></ul>",
				},
			},
		},
		map[string]any{
			"id":        "portfolio",
			"type":      common.ResumeSectionTypePortfolio,
			"title":     "个人作品",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "pf1",
					"description": `<ul><li>代码仓库：<a href="https://code.zhwei.invalid">code.zhwei.invalid</a></li><li>技术博客：<a href="https://blog.zhwei.invalid">blog.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateJavaPreset() []byte {
	javaBtye, _ := json.Marshal(javaPresetJSON)
	return javaBtye
}
