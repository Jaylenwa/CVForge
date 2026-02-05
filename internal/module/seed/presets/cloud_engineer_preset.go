package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var cloudEngineerPresetJSON = map[string]any{
	"title":    "云计算工程师简历",
	"language": "zh",
	"Personal": map[string]any{
		"FullName":   "张伟",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "云计算工程师",
		"City":       "上海",
		"Money":      "30k-45k",
		"JoinTime":   "随时到岗",
		"Gender":     "男",
		"Age":        "29",
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
					"description": "<ul><li>5 年云计算经验，熟悉公有云与云原生：AWS/Azure/阿里云、Kubernetes</li><li>基础设施工程：Terraform/IaC、网络与安全、成本优化</li><li>平台与交付：Helm/Argo CD、多环境治理与策略</li><li>可观测与可靠性：指标/日志/链路，SLO/SLI 与演练</li></ul>",
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
					"title":       "某零售企业",
					"subtitle":    "云计算工程师（平台/IaC）",
					"timeStart":   "2021-05",
					"today":       true,
					"description": "<ul><li>Terraform 管理基础设施与网络安全策略</li><li>Kubernetes 与 GitOps 治理多环境交付</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "某互联网企业",
					"subtitle":    "云计算工程师（可观测/成本）",
					"timeStart":   "2018-07",
					"timeEnd":     "2021-04",
					"description": "<ul><li>建设可观测平台与告警策略</li><li>成本优化与资源评估，提升整体效率</li></ul>",
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
					"title":       "云原生平台与 IaC 治理",
					"subtitle":    "云计算工程师",
					"timeStart":   "2022-03",
					"timeEnd":     "2022-11",
					"description": "<ul><li>技术栈：AWS / Kubernetes / Terraform / Argo CD</li><li>统一资源与交付策略，治理多账户/多环境</li><li>成本优化与可观测建设</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "网络与安全策略建设",
					"subtitle":    "云计算工程师",
					"timeStart":   "2019-07",
					"timeEnd":     "2020-03",
					"description": "<ul><li>技术栈：AWS Security / VPC / WAF</li><li>完善网络与安全策略，接入审计与合规</li></ul>",
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
					"title":       "云与平台",
					"subtitle":    "云原生",
					"description": "AWS / Azure / 阿里云 / Kubernetes",
				},
				map[string]any{
					"id":          "s2",
					"title":       "基础设施工程",
					"subtitle":    "IaC",
					"description": "Terraform / VPC / Security / Cost",
				},
				map[string]any{
					"id":          "s3",
					"title":       "交付与可靠性",
					"subtitle":    "平台",
					"description": "Helm / Argo CD / Observability / SLO/SLI",
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
					"timeStart":   "2015-09",
					"timeEnd":     "2018-06",
					"description": "研究方向为云计算与平台工程；参与平台治理与交付。",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "浙江大学",
					"major":       "计算机科学与技术",
					"degree":      "本科",
					"timeStart":   "2011-09",
					"timeEnd":     "2015-06",
					"description": "主修分布式系统与工程化；参与校内项目。",
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
					"description": `<ul><li>AWS IaC 示例：<a href="https://aws-iac.zhwei.invalid">aws-iac.zhwei.invalid</a></li><li>云原生平台演示：<a href="https://cloud-native.zhwei.invalid">cloud-native.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateCloudEngineerPreset() []byte {
	ceByte, _ := json.Marshal(cloudEngineerPresetJSON)
	return ceByte
}
