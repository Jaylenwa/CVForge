package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var nodejsPresetJSON = map[string]any{
	"title":    "Node.js 后端开发工程师简历",
	"language": "zh",
	"Personal": map[string]any{
		"FullName":   "张伟",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "Node.js 后端开发工程师",
		"City":       "上海 / 杭州",
		"Money":      "25k-40k·14薪",
		"JoinTime":   "随时到岗",
		"Gender":     "男",
		"Age":        "28",
		"Degree":     "硕士",
		"CustomInfo": `[{"label":"政治面貌","value":"党员"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#16a34a",
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
					"description": "<ul><li>5 年 Node.js 服务端开发经验，熟悉 TypeScript、NestJS/Express，具备高并发接口优化与稳定性建设经验</li><li>熟悉 MySQL/PostgreSQL、Redis、Kafka/RabbitMQ，掌握缓存一致性、消息可靠投递、幂等与重试设计</li><li>具备工程化与云原生实践：Docker/Kubernetes、CI/CD、指标/日志/链路追踪</li></ul>",
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
					"title":       "某在线教育平台",
					"subtitle":    "Node.js 后端开发工程师（课程/支付）",
					"timeStart":   "2021-06",
					"today":       true,
					"description": "<ul><li>负责课程与支付核心服务（NestJS + TypeScript），推动接口规范化与统一鉴权，提升迭代效率</li><li>落地 Redis 缓存与热点治理，核心接口 P95 从 230ms 降至 140ms</li><li>引入消息队列（Kafka）异步化通知与对账链路，降低主链路耗时并提升可靠性</li><li>建设可观测体系（Prometheus/Grafana + tracing），故障定位时间缩短 50%</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "某内容社区公司",
					"subtitle":    "Node.js 工程师（平台/中台）",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-05",
					"description": "<ul><li>参与账号与权限平台建设（Express + TypeScript + PostgreSQL），支撑多业务线接入</li><li>优化慢查询与索引策略，核心查询耗时下降 60%</li><li>建立 CI/CD 与灰度发布流程，降低发布风险</li></ul>",
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
					"title":       "统一通知与任务调度平台",
					"subtitle":    "Node.js 后端开发工程师",
					"timeStart":   "2022-05",
					"timeEnd":     "2023-02",
					"description": "<ul><li>技术栈：Node.js / TypeScript / NestJS / Redis / Kafka</li><li>抽象任务模板、重试与幂等机制，降低重复开发成本并提升可靠性</li><li>实现任务可观测：状态追踪、失败原因聚合与告警推送</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "高并发活动与优惠券系统",
					"subtitle":    "Node.js 后端开发工程师",
					"timeStart":   "2020-09",
					"timeEnd":     "2021-03",
					"description": "<ul><li>技术栈：Node.js / Redis / RabbitMQ</li><li>设计库存扣减与发放链路幂等方案，避免重复领取与超卖问题</li><li>采用 Redis Lua 原子脚本与热点治理，峰值接口 P95 延迟下降 40%</li></ul>",
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
					"subtitle":    "Node.js",
					"description": "TypeScript / NestJS / Express / REST API / OpenAPI",
				},
				map[string]any{
					"id":          "s2",
					"title":       "存储与中间件",
					"subtitle":    "数据库 / 缓存 / 消息",
					"description": "MySQL / PostgreSQL / Redis / Kafka / RabbitMQ",
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
			"id":        "edu",
			"type":      common.ResumeSectionTypeEducation,
			"title":     "教育背景",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "ed1",
					"title":       "上海交通大学",
					"major":       "计算机科学与技术",
					"degree":      "本科",
					"timeStart":   "2015-09",
					"timeEnd":     "2019-06",
					"description": "主修数据结构、数据库、计算机网络；参与校内项目与竞赛。",
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
					"description": `<ul><li>开源工具：<a href="https://npm.zhwei.invalid">npm.zhwei.invalid</a></li><li>接口示例站：<a href="https://api-demo.zhwei.invalid">api-demo.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateNodejsPreset() []byte {
	nodejsByte, _ := json.Marshal(nodejsPresetJSON)
	return nodejsByte
}
