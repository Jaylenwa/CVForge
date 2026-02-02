package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var phpPresetJSON = map[string]any{
	"title":    "PHP 后端开发工程师简历",
	"language": "zh",
	"Personal": map[string]any{
		"FullName":   "张伟",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "PHP 后端开发工程师",
		"City":       "上海 / 杭州",
		"Money":      "25k-35k·14薪",
		"JoinTime":   "随时到岗",
		"Gender":     "男",
		"Age":        "28",
		"Degree":     "硕士",
		"CustomInfo": `[{"label":"政治面貌","value":"党员"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#7c3aed",
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
					"description": "<ul><li>5 年 PHP 后端开发经验，熟悉 Laravel/Symfony 与领域建模，能独立推进需求落地与性能优化</li><li>熟悉 MySQL/Redis/Elasticsearch、Kafka/RabbitMQ，掌握缓存一致性、消息可靠投递、幂等与重试设计</li><li>具备工程化与云原生实践：Docker/Kubernetes、CI/CD、监控告警与链路追踪</li></ul>",
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
					"title":       "某头部本地生活平台",
					"subtitle":    "PHP 后端开发工程师（订单/履约）",
					"timeStart":   "2021-06",
					"today":       true,
					"description": "<ul><li>负责订单与履约核心接口（Laravel + MySQL + Redis），落地限流/熔断/降级，峰值期稳定性显著提升</li><li>建设 Redis 多级缓存与热点 Key 治理，核心接口 P95 从 220ms 降至 140ms</li><li>引入消息队列（Kafka）异步化支付回调与通知链路，降低主链路耗时并提升成功率</li><li>完善监控告警与日志规范，线上故障定位时间缩短 50%</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "某企业服务 SaaS 公司",
					"subtitle":    "PHP 工程师（平台/中台）",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-05",
					"description": "<ul><li>参与统一权限与审计平台建设（Symfony + MySQL），支撑 100+ 企业客户</li><li>优化慢 SQL 与索引策略，核心查询耗时下降 60%</li><li>建立 API 规范与自动化测试流程（OpenAPI + PHPUnit），回归效率提升</li></ul>",
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
					"title":       "高并发秒杀与库存系统",
					"subtitle":    "PHP 后端开发工程师",
					"timeStart":   "2022-04",
					"timeEnd":     "2022-12",
					"description": "<ul><li>技术栈：PHP / Laravel / Redis / MySQL / Kafka</li><li>设计库存预扣与订单创建幂等方案，避免重复下单与超卖</li><li>采用 Redis Lua 原子脚本与热点治理，峰值接口 P95 延迟下降 40%</li><li>落地消息重试与补偿任务，保障最终一致性与可恢复</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "统一 API 网关与权限组件",
					"subtitle":    "PHP 后端开发工程师",
					"timeStart":   "2020-10",
					"timeEnd":     "2021-03",
					"description": "<ul><li>技术栈：PHP / Nginx / Redis / JWT</li><li>抽象鉴权、签名、限流与灰度能力，提升业务接入效率并降低重复开发</li><li>实现权限缓存与审计日志，满足合规与排障需求</li></ul>",
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
					"subtitle":    "PHP",
					"description": "Laravel / Symfony / Composer / REST API",
				},
				map[string]any{
					"id":          "s2",
					"title":       "存储与中间件",
					"subtitle":    "数据库 / 缓存 / 消息",
					"description": "MySQL / Redis / Elasticsearch / Kafka / RabbitMQ",
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
					"title":       "华东理工大学",
					"major":       "计算机科学与技术",
					"degree":      "硕士",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "研究方向为分布式系统与高并发服务；参与实验室平台项目与性能优化。",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "华东理工大学",
					"major":       "计算机科学与技术",
					"degree":      "本科",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "主修数据结构、操作系统、数据库；参与校内项目与竞赛。",
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
					"description": `<ul><li>开源组件：<a href="https://packagist.zhwei.invalid">packagist.zhwei.invalid</a></li><li>技术博客：<a href="https://blog.zhwei.invalid">blog.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GeneratePHPPreset() []byte {
	phpByte, _ := json.Marshal(phpPresetJSON)
	return phpByte
}
