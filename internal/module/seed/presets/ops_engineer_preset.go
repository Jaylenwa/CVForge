package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var opsEngineerPresetJSON = map[string]any{
	"title":    "运维工程师简历",
	"language": "zh",
	"Personal": map[string]any{
		"FullName":   "张伟",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "运维工程师",
		"City":       "上海",
		"Money":      "25k-35k",
		"JoinTime":   "随时到岗",
		"Gender":     "男",
		"Age":        "28",
		"Degree":     "硕士",
		"CustomInfo": `[{"label":"政治面貌","value":"党员"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#059669",
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
					"description": "<ul><li>5 年运维经验，熟悉 Linux、网络与服务治理，具备高可用与故障应对能力</li><li>掌握监控与告警体系：Prometheus/Grafana、ELK、日志/指标/链路追踪</li><li>熟悉自动化与基础设施：Ansible/Terraform、CI/CD、配置与变更管理</li><li>具备容量与性能优化经验：缓存/连接池/并发策略，与业务协作提升稳定性</li></ul>",
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
					"title":       "某电商平台",
					"subtitle":    "运维工程师（稳定性/可观测）",
					"timeStart":   "2021-06",
					"today":       true,
					"description": "<ul><li>搭建可观测平台与告警分级，缩短故障定位与恢复时间</li><li>优化 Nginx/Redis/MySQL 的资源与连接策略，提升峰值稳定性</li><li>完善发布与回滚流程，减少上线风险</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "某内容社区",
					"subtitle":    "运维工程师（平台/交付）",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-05",
					"description": "<ul><li>建设 CI/CD 与灰度发布，提升交付效率与质量</li><li>统一日志与指标采集，完善故障演练与预案</li></ul>",
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
					"title":       "可观测平台与告警治理",
					"subtitle":    "运维工程师",
					"timeStart":   "2022-05",
					"timeEnd":     "2023-02",
					"description": "<ul><li>技术栈：Prometheus / Grafana / ELK</li><li>统一采集与告警治理，提升故障响应与定位效率</li><li>建设 SLO/SLI 指标体系，推动稳定性提升</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "发布与变更管理平台",
					"subtitle":    "运维工程师",
					"timeStart":   "2020-09",
					"timeEnd":     "2021-03",
					"description": "<ul><li>技术栈：GitLab CI / Ansible / Terraform</li><li>沉淀发布流程与变更追踪，降低上线与配置风险</li><li>接入质量门禁与回滚策略</li></ul>",
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
					"title":       "系统与网络",
					"subtitle":    "基础",
					"description": "Linux / Shell / TCP/IP / Nginx / HAProxy",
				},
				map[string]any{
					"id":          "s2",
					"title":       "数据库与缓存",
					"subtitle":    "存储",
					"description": "MySQL / PostgreSQL / Redis",
				},
				map[string]any{
					"id":          "s3",
					"title":       "可观测与自动化",
					"subtitle":    "工程",
					"description": "Prometheus / Grafana / ELK / Ansible / Terraform / CI/CD",
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
					"title":       "华南理工大学",
					"major":       "软件工程",
					"degree":      "硕士",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "研究方向为可观测与运维工程；参与平台建设与稳定性治理。",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "华南理工大学",
					"major":       "计算机科学与技术",
					"degree":      "本科",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "主修系统与网络；参与校内项目与竞赛。",
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
					"description": `<ul><li>运维工具集：<a href="https://ops-toolkit.zhwei.invalid">ops-toolkit.zhwei.invalid</a></li><li>技术博客：<a href="https://blog.zhwei.invalid">blog.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateOpsEngineerPreset() []byte {
	opsByte, _ := json.Marshal(opsEngineerPresetJSON)
	return opsByte
}
