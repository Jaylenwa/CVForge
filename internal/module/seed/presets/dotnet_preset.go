package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var dotnetPresetJSON = map[string]any{
	"title":    ".NET 后端开发工程师简历",
	"language": "zh",
	"Personal": map[string]any{
		"FullName":   "张伟",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        ".NET 后端开发工程师",
		"City":       "上海",
		"Money":      "25k-40k",
		"JoinTime":   "随时到岗",
		"Gender":     "男",
		"Age":        "28",
		"Degree":     "硕士",
		"CustomInfo": `[{"label":"政治面貌","value":"党员"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#e11d48",
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
					"description": "<ul><li>5 年 .NET 服务端开发经验，熟悉 ASP.NET Core、微服务与高可用架构，具备性能优化与稳定性建设经验</li><li>熟悉 gRPC、消息队列与事件驱动架构，掌握幂等、重试、降级与最终一致性实践</li><li>熟悉云原生与工程化：Docker/Kubernetes、CI/CD、指标/日志/链路追踪</li></ul>",
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
					"title":       "某大型企业数字化团队",
					"subtitle":    ".NET 后端开发工程师（平台/中台）",
					"timeStart":   "2021-08",
					"today":       true,
					"description": "<ul><li>负责平台核心服务（ASP.NET Core + gRPC），推进服务拆分与接口治理，提升迭代效率与稳定性</li><li>落地分布式缓存与读写分离策略，核心接口 P95 从 260ms 降至 160ms</li><li>建设发布流程与回滚方案（CI/CD），发布成功率提升并降低业务中断风险</li><li>完善可观测体系（Prometheus/Grafana + tracing），故障定位时间缩短 50%</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "某物流科技公司",
					"subtitle":    ".NET 工程师（调度/计费）",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-07",
					"description": "<ul><li>参与调度与计费系统建设（.NET + SQL Server + Redis），支撑多城市业务扩张</li><li>优化批处理与报表链路，任务耗时下降 40%</li><li>完善异常处理与告警策略，显著降低夜间故障次数</li></ul>",
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
					"title":       "统一订单与结算中台",
					"subtitle":    ".NET 后端开发工程师",
					"timeStart":   "2022-06",
					"timeEnd":     "2023-03",
					"description": "<ul><li>技术栈：ASP.NET Core / gRPC / Redis / PostgreSQL / Kafka</li><li>抽象订单/结算核心领域模型与接口，支撑多业务线复用并降低重复开发</li><li>引入事件驱动与最终一致性方案，降低跨服务耦合并提升可演进性</li><li>建设压测与容量评估体系，峰值 QPS 提升 2 倍且稳定运行</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "统一身份与权限服务",
					"subtitle":    ".NET 后端开发工程师",
					"timeStart":   "2020-09",
					"timeEnd":     "2021-02",
					"description": "<ul><li>技术栈：.NET / JWT / Redis</li><li>实现 token 策略、登录态管理与权限缓存，提升鉴权性能与可运维性</li><li>完善审计与告警能力，满足合规需求</li></ul>",
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
					"title":       "平台与框架",
					"subtitle":    ".NET",
					"description": "ASP.NET Core / gRPC / EF Core / Dapper / REST API",
				},
				map[string]any{
					"id":          "s2",
					"title":       "存储与中间件",
					"subtitle":    "数据库 / 缓存 / 消息",
					"description": "SQL Server / PostgreSQL / Redis / Kafka",
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
					"title":       "哈尔滨工业大学",
					"major":       "软件工程",
					"degree":      "硕士",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "研究方向为分布式系统与工程化实践；参与实验室平台项目与性能优化。",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "哈尔滨工业大学",
					"major":       "软件工程",
					"degree":      "本科",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "主修软件工程与分布式系统；参与校内项目与性能优化实践。",
				},
			},
		},
	},
}

func GenerateDotnetPreset() []byte {
	dotnetByte, _ := json.Marshal(dotnetPresetJSON)
	return dotnetByte
}
