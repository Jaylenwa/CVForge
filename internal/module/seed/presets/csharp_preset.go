package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var csharpPresetJSON = map[string]any{
	"title":    "C# 后端开发工程师简历",
	"language": "zh",
	"Personal": map[string]any{
		"FullName":   "张伟",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "C# 后端开发工程师",
		"City":       "上海",
		"Money":      "25k-40k",
		"JoinTime":   "随时到岗",
		"Gender":     "男",
		"Age":        "28",
		"Degree":     "硕士",
		"CustomInfo": `[{"label":"政治面貌","value":"党员"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#22c55e",
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
					"description": "<ul><li>5 年 C# 服务端开发经验，熟悉 ASP.NET Core、微服务与高并发场景优化，具备稳定性与可观测建设经验</li><li>熟悉 SQL Server/PostgreSQL、Redis、Kafka，掌握缓存一致性、消息可靠投递、分布式事务与幂等设计</li><li>具备工程化能力：CI/CD、自动化测试、容器化部署（Docker/Kubernetes）</li></ul>",
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
					"title":       "某企业服务 SaaS 公司",
					"subtitle":    "C# 后端开发工程师（多租户平台）",
					"timeStart":   "2021-05",
					"today":       true,
					"description": "<ul><li>负责多租户核心服务（ASP.NET Core + EF Core），设计租户隔离与权限模型，支撑 100+ 企业客户</li><li>落地 Redis 缓存与热点治理，核心接口 P95 从 240ms 降至 150ms</li><li>引入 Kafka 异步化审计与通知链路，降低主链路耗时并提升可靠性</li><li>完善监控告警与链路追踪（Prometheus/Grafana），故障定位时间缩短 50%</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "某零售科技公司",
					"subtitle":    "C# 工程师（订单/支付）",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-04",
					"description": "<ul><li>参与订单与支付系统改造（ASP.NET Core + SQL Server），引入幂等与重试策略，异常请求比例显著下降</li><li>优化慢 SQL 与索引策略，核心查询耗时下降 60%</li><li>建立接口规范与自动化测试流程（xUnit + Postman Collection），回归效率提升</li></ul>",
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
					"title":       "统一权限与审计平台",
					"subtitle":    "C# 后端开发工程师",
					"timeStart":   "2022-03",
					"timeEnd":     "2022-11",
					"description": "<ul><li>技术栈：ASP.NET Core / EF Core / Redis / PostgreSQL</li><li>设计 RBAC/ABAC 混合模型与权限缓存策略，兼顾一致性与性能</li><li>实现审计事件异步采集与查询聚合，满足合规与排障需求</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "网关与服务治理组件",
					"subtitle":    "C# 后端开发工程师",
					"timeStart":   "2020-08",
					"timeEnd":     "2021-02",
					"description": "<ul><li>技术栈：.NET / Ocelot / Redis</li><li>落地鉴权、限流、灰度、熔断与统一日志，提升多团队协作效率</li><li>完善配置化与回滚机制，降低发布风险</li></ul>",
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
					"subtitle":    "C#",
					"description": "ASP.NET Core / EF Core / Dapper / REST API / gRPC",
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
					"title":       "同济大学",
					"major":       "软件工程",
					"degree":      "硕士",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "研究方向为企业级应用架构与分布式系统；参与实验室平台项目与性能优化。",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "同济大学",
					"major":       "软件工程",
					"degree":      "本科",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "主修面向对象设计、数据库、分布式系统；参与校内项目实践。",
				},
			},
		},
	},
}

func GenerateCSharpPreset() []byte {
	csharpByte, _ := json.Marshal(csharpPresetJSON)
	return csharpByte
}
