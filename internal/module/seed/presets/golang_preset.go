package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var golangPresetJSON = map[string]any{
	"title":    "Go 后端开发工程师简历",
	"language": "zh",
	"Personal": map[string]any{
		"FullName":   "张伟",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "Go 后端开发工程师",
		"City":       "上海",
		"Money":      "25k-35k",
		"JoinTime":   "随时到岗",
		"Gender":     "男",
		"Age":        "28",
		"Degree":     "硕士",
		"CustomInfo": `[{"label":"政治面貌","value":"党员"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#14b8a6",
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
					"description": "<ul><li>4+ 年 Go 后端经验，熟悉微服务、gRPC、分布式一致性与高并发场景优化</li><li>熟悉 MySQL/Redis/Kafka，具备性能调优、故障治理与稳定性建设经验</li><li>熟悉可观测与云原生：Prometheus/Grafana、链路追踪、Docker/Kubernetes</li></ul>",
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
					"title":       "某大型互联网公司",
					"subtitle":    "Go 后端开发工程师（订单/履约）",
					"timeStart":   "2021-08",
					"today":       true,
					"description": "<ul><li>负责订单履约核心服务（Go + gRPC），支持日均千万级请求，接口 P95 降低 35%</li><li>建设幂等、分布式锁与补偿机制，保障高并发下的数据一致性与可恢复</li><li>推动 Redis 热点治理与降级策略，核心链路稳定性显著提升</li><li>完善可观测体系（指标/日志/追踪），重大故障恢复时间从小时级降至 20 分钟内</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "某金融科技公司",
					"subtitle":    "后端工程师（网关/平台）",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-07",
					"description": "<ul><li>参与统一 API 网关建设（鉴权、限流、灰度、熔断），业务接入效率提升</li><li>落地 Kafka 异步链路与延迟队列，解耦核心服务，削峰填谷</li><li>推进数据库治理：慢 SQL 优化、索引与分库分表策略，核心查询耗时下降 60%</li></ul>",
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
					"title":       "统一鉴权与权限平台",
					"subtitle":    "Go 后端开发工程师",
					"timeStart":   "2022-02",
					"timeEnd":     "2022-10",
					"description": "<ul><li>技术栈：Go / gRPC / Redis / MySQL</li><li>设计 RBAC/ABAC 混合模型与权限缓存策略，保障一致性与性能</li><li>实现 token 轮换与多端登录控制，完善审计与告警能力</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "订单履约链路重构",
					"subtitle":    "Go 后端开发工程师",
					"timeStart":   "2023-01",
					"timeEnd":     "2023-12",
					"description": "<ul><li>技术栈：Go / Kafka / Redis</li><li>拆分履约子域与事件驱动流程，降低耦合并提升可演进性</li><li>引入 Saga 补偿与重试策略，显著降低异常订单占比</li></ul>",
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
					"title":       "后端",
					"subtitle":    "Go",
					"description": "Gin / gRPC / 并发与性能优化 / 微服务架构",
				},
				map[string]any{
					"id":          "s2",
					"title":       "数据库与中间件",
					"subtitle":    "MySQL / Redis / Kafka",
					"description": "事务与锁、索引优化、缓存一致性、消息可靠投递",
				},
				map[string]any{
					"id":          "s3",
					"title":       "稳定性与云原生",
					"subtitle":    "可观测",
					"description": "Prometheus / Grafana / 日志与追踪 / Docker / Kubernetes",
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
					"title":       "中山大学",
					"major":       "软件工程",
					"degree":      "硕士",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "研究方向为分布式系统与高并发服务；参与实验室平台项目与性能优化。",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "中山大学",
					"major":       "软件工程",
					"degree":      "本科",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "主修软件工程与分布式系统；参与校内项目实践。",
				},
			},
		},
	},
}

func GenerateGolangPreset() []byte {
	golangByte, _ := json.Marshal(golangPresetJSON)
	return golangByte
}
