package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var devopsPresetJSON = map[string]any{
	"title":    "DevOps工程师简历",
	"language": "zh",
	"Personal": map[string]any{
		"FullName":   "张伟",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "DevOps工程师",
		"City":       "上海",
		"Money":      "30k-40k",
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
					"description": "<ul><li>5 年 DevOps 经验，熟悉容器与云原生：Docker/Kubernetes/Helm</li><li>擅长 GitOps 与持续交付：Argo CD/Flux、蓝绿/灰度/金丝雀发布</li><li>具备平台与工程化能力：Pipeline 治理、依赖与配置管理、环境与镜像管理</li><li>完善 SRE 体系：可观测/容量评估/故障演练，驱动稳定性与效率提升</li></ul>",
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
					"title":       "某电商中台",
					"subtitle":    "DevOps工程师（交付/平台）",
					"timeStart":   "2021-06",
					"today":       true,
					"description": "<ul><li>建设 GitOps 交付平台（Argo CD），统一环境与发布策略</li><li>沉淀 Helm Chart 与依赖管理，提升复用与一致性</li><li>优化流水线与镜像管理，缩短交付周期</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "某内容平台",
					"subtitle":    "DevOps工程师（云原生/SRE）",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-05",
					"description": "<ul><li>推进可观测与容量评估，完善告警与自动化修复</li><li>落地故障演练与预案，提升 MTTR 与可靠性</li></ul>",
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
					"title":       "GitOps 交付平台",
					"subtitle":    "DevOps工程师",
					"timeStart":   "2022-05",
					"timeEnd":     "2023-02",
					"description": "<ul><li>技术栈：Argo CD / Helm / Kubernetes</li><li>统一环境与发布策略，支持蓝绿/灰度/金丝雀</li><li>治理依赖与配置，提升交付效率与质量</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "可靠性与容量治理",
					"subtitle":    "DevOps工程师",
					"timeStart":   "2020-09",
					"timeEnd":     "2021-03",
					"description": "<ul><li>技术栈：Prometheus / Grafana / KEDA</li><li>完善指标与自动扩缩容方案，提升稳定性与资源效率</li><li>构建故障演练流程与演习</li></ul>",
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
					"title":       "容器与编排",
					"subtitle":    "云原生",
					"description": "Docker / Kubernetes / Helm",
				},
				map[string]any{
					"id":          "s2",
					"title":       "交付与治理",
					"subtitle":    "GitOps",
					"description": "Argo CD / Flux / CI/CD / 蓝绿 / 灰度 / 金丝雀",
				},
				map[string]any{
					"id":          "s3",
					"title":       "可观测与可靠性",
					"subtitle":    "SRE",
					"description": "Prometheus / Grafana / Alerting / Capacity / Chaos",
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
					"description": "研究方向为云原生与交付工程；参与平台建设与治理。",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "同济大学",
					"major":       "计算机科学与技术",
					"degree":      "本科",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "主修分布式系统与工程化；参与校内项目与竞赛。",
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
					"description": `<ul><li>Helm Chart 集合：<a href="https://helm-charts.zhwei.invalid">helm-charts.zhwei.invalid</a></li><li>GitOps 示例：<a href="https://gitops.zhwei.invalid">gitops.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateDevopsPreset() []byte {
	devopsByte, _ := json.Marshal(devopsPresetJSON)
	return devopsByte
}
