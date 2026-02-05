package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var cloudEngineerPresetJSONEn = map[string]any{
	"title":    "Cloud Engineer Resume",
	"language": "en",
	"Personal": map[string]any{
		"FullName":   "Wei Zhang",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "Cloud Engineer",
		"City":       "Shanghai",
		"Money":      "30k-45k",
		"JoinTime":   "Available immediately",
		"Gender":     "Male",
		"Age":        "29",
		"Degree":     "Master's",
		"CustomInfo": `[{"label":"Political Status","value":"Party Member"}]`,
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
			"title":     "Summary",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "sum1",
					"description": "<ul><li>5 years in cloud engineering; public cloud and cloud-native: AWS/Azure/Alibaba Cloud, Kubernetes</li><li>Infrastructure engineering: Terraform/IaC, networking/security, cost optimization</li><li>Platform and delivery: Helm/Argo CD, multi-environment governance</li><li>Observability & reliability: logs/metrics/tracing; SLO/SLI and drills</li></ul>",
				},
			},
		},
		map[string]any{
			"id":        "exp",
			"type":      common.ResumeSectionTypeExperience,
			"title":     "Work Experience",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "e1",
					"title":       "Retail Company",
					"subtitle":    "Cloud Engineer (Platform/IaC)",
					"timeStart":   "2021-05",
					"today":       true,
					"description": "<ul><li>Terraform for infra and network security policies</li><li>Kubernetes and GitOps for multi-environment delivery</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "Internet Company",
					"subtitle":    "Cloud Engineer (Observability/Cost)",
					"timeStart":   "2018-07",
					"timeEnd":     "2021-04",
					"description": "<ul><li>Built observability platform and alert strategies</li><li>Cost optimization and resource evaluation for efficiency</li></ul>",
				},
			},
		},
		map[string]any{
			"id":        "projects",
			"type":      common.ResumeSectionTypeProjects,
			"title":     "Projects",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "p1",
					"title":       "Cloud-native Platform & IaC Governance",
					"subtitle":    "Cloud Engineer",
					"timeStart":   "2022-03",
					"timeEnd":     "2022-11",
					"description": "<ul><li>Stack: AWS / Kubernetes / Terraform / Argo CD</li><li>Unified resources and delivery strategies across accounts/environments</li><li>Cost optimization and observability</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "Networking & Security Policies",
					"subtitle":    "Cloud Engineer",
					"timeStart":   "2019-07",
					"timeEnd":     "2020-03",
					"description": "<ul><li>Stack: AWS Security / VPC / WAF</li><li>Networking and security policies; audits and compliance</li></ul>",
				},
			},
		},
		map[string]any{
			"id":        "skills",
			"type":      common.ResumeSectionTypeSkills,
			"title":     "Skills",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "s1",
					"title":       "Cloud & Platform",
					"subtitle":    "Cloud-native",
					"description": "AWS / Azure / Alibaba Cloud / Kubernetes",
				},
				map[string]any{
					"id":          "s2",
					"title":       "Infrastructure Engineering",
					"subtitle":    "IaC",
					"description": "Terraform / VPC / Security / Cost",
				},
				map[string]any{
					"id":          "s3",
					"title":       "Delivery & Reliability",
					"subtitle":    "Platform",
					"description": "Helm / Argo CD / Observability / SLO/SLI",
				},
			},
		},
		map[string]any{
			"id":        "edu",
			"type":      common.ResumeSectionTypeEducation,
			"title":     "Education",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "ed1",
					"title":       "Zhejiang University",
					"major":       "Software Engineering",
					"degree":      "Master's",
					"timeStart":   "2015-09",
					"timeEnd":     "2018-06",
					"description": "Cloud computing and platform engineering; governance and delivery.",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "Zhejiang University",
					"major":       "Computer Science and Technology",
					"degree":      "Bachelor's",
					"timeStart":   "2011-09",
					"timeEnd":     "2015-06",
					"description": "Distributed systems and engineering; campus projects.",
				},
			},
		},
		map[string]any{
			"id":        "portfolio",
			"type":      common.ResumeSectionTypePortfolio,
			"title":     "Portfolio",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "pf1",
					"description": `<ul><li>AWS IaC examples: <a href="https://aws-iac.zhwei.invalid">aws-iac.zhwei.invalid</a></li><li>Cloud-native demos: <a href="https://cloud-native.zhwei.invalid">cloud-native.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateCloudEngineerPresetEn() []byte {
	ceByte, _ := json.Marshal(cloudEngineerPresetJSONEn)
	return ceByte
}
