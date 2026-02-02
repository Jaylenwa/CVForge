package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var pythonPresetJSON = map[string]any{
	"title":    "Python 后端开发工程师简历",
	"language": "zh",
	"Personal": map[string]any{
		"FullName":   "张伟",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "Python 后端开发工程师",
		"City":       "上海 / 杭州",
		"Money":      "25k-35k·14薪",
		"JoinTime":   "随时到岗",
		"Gender":     "男",
		"Age":        "28",
		"Degree":     "硕士",
		"CustomInfo": `[{"label":"政治面貌","value":"党员"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#2563eb",
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
					"description": "<ul><li>5 年 Python 后端开发经验，熟悉 FastAPI/Django，具备从 0 到 1 业务落地与性能优化能力</li><li>熟悉 PostgreSQL/MySQL、Redis、Kafka、Celery，能构建高并发与高可用服务</li><li>具备工程化与云原生实践：Docker、Kubernetes、CI/CD、日志/监控/告警</li></ul>",
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
					"title":       "某头部电商平台",
					"subtitle":    "Python 后端开发工程师（交易/风控）",
					"timeStart":   "2021-05",
					"today":       true,
					"description": "<ul><li>负责交易域核心接口与异步链路（FastAPI + Kafka + Celery），订单创建成功率提升至 99.95%</li><li>落地 Redis 多级缓存与热点治理，核心接口 P95 从 240ms 降至 130ms</li><li>推动风控规则引擎重构，规则配置化与灰度发布，策略迭代效率提升 3 倍</li><li>完善监控与告警（Prometheus/Grafana），故障定位时间缩短 50%</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "某 SaaS 创业公司",
					"subtitle":    "Python 后端工程师（数据与平台）",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-04",
					"description": "<ul><li>基于 Django 构建多租户权限、审计日志与配置中心，支撑 100+ 企业客户</li><li>实现 ETL 数据处理与任务编排（Celery + PostgreSQL），报表生成耗时下降 40%</li><li>参与 API 规范与自动化测试建设（OpenAPI + pytest），回归效率显著提升</li></ul>",
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
					"title":       "实时风控规则引擎",
					"subtitle":    "Python 后端开发工程师",
					"timeStart":   "2022-06",
					"timeEnd":     "2023-03",
					"description": "<ul><li>技术栈：FastAPI / Kafka / Redis / PostgreSQL</li><li>设计规则 DSL 与执行器，支持按商家/人群维度动态下发与快速回滚</li><li>接入 Kafka 流式事件与 Redis 画像缓存，实现毫秒级判定</li><li>完善压测与容量评估，峰值 QPS 提升 2 倍且稳定运行</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "任务编排与异步处理平台",
					"subtitle":    "Python 后端开发工程师",
					"timeStart":   "2020-03",
					"timeEnd":     "2020-11",
					"description": "<ul><li>技术栈：Django / Celery / Redis</li><li>抽象任务模板与重试/幂等机制，降低重复开发成本</li><li>实现任务可观测：状态追踪、失败原因聚合与告警推送</li></ul>",
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
					"subtitle":    "Python",
					"description": "FastAPI / Django / REST API / OpenAPI",
				},
				map[string]any{
					"id":          "s2",
					"title":       "存储与中间件",
					"subtitle":    "数据库 / 缓存 / 消息",
					"description": "PostgreSQL / MySQL / Redis / Kafka / Celery",
				},
				map[string]any{
					"id":          "s3",
					"title":       "工程化",
					"subtitle":    "可观测与交付",
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
					"title":       "华中科技大学",
					"major":       "计算机科学与技术",
					"degree":      "硕士",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "研究方向为分布式系统与高可用架构；参与实验室平台项目与性能优化。",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "华中科技大学",
					"major":       "计算机科学与技术",
					"degree":      "本科",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "主修数据结构、数据库、计算机网络；参与校内项目与竞赛。",
				},
			},
		},
	},
}

func GeneratePythonPreset() []byte {
	pythonByte, _ := json.Marshal(pythonPresetJSON)
	return pythonByte
}
