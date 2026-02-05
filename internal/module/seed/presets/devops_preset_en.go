package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var devopsPresetJSONEn = map[string]any{
	"title":    "DevOps Engineer Resume",
	"language": "en",
	"Personal": map[string]any{
		"FullName":   "Wei Zhang",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "DevOps Engineer",
		"City":       "Shanghai",
		"Money":      "30k-40k",
		"JoinTime":   "Available immediately",
		"Gender":     "Male",
		"Age":        "28",
		"Degree":     "Master's",
		"CustomInfo": `[{"label":"Political Status","value":"Party Member"}]`,
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
			"title":     "Summary",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "sum1",
					"description": "<ul><li>5 years in DevOps; cloud-native stack: Docker/Kubernetes/Helm</li><li>GitOps and continuous delivery: Argo CD/Flux; blue/green/canary releases</li><li>Platform engineering: pipeline governance, dependency/config management, env/image management</li><li>SRE practices: observability/capacity drills/incident playbooks to drive reliability</li></ul>",
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
					"title":       "E-commerce Middleware",
					"subtitle":    "DevOps Engineer (Delivery/Platform)",
					"timeStart":   "2021-06",
					"today":       true,
					"description": "<ul><li>GitOps delivery with Argo CD; unified environments and release strategies</li><li>Helm charts and dependency management for reuse and consistency</li><li>Pipeline and image management to shorten delivery cycles</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "Content Platform",
					"subtitle":    "DevOps Engineer (Cloud-native/SRE)",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-05",
					"description": "<ul><li>Observability and capacity; improved alerting and automated remediation</li><li>Incident drills and response plans to improve MTTR and reliability</li></ul>",
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
					"title":       "GitOps Delivery Platform",
					"subtitle":    "DevOps Engineer",
					"timeStart":   "2022-05",
					"timeEnd":     "2023-02",
					"description": "<ul><li>Stack: Argo CD / Helm / Kubernetes</li><li>Unified environment and release strategies (blue/green/canary)</li><li>Dependency/config governance for efficient delivery</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "Reliability & Capacity Governance",
					"subtitle":    "DevOps Engineer",
					"timeStart":   "2020-09",
					"timeEnd":     "2021-03",
					"description": "<ul><li>Stack: Prometheus / Grafana / KEDA</li><li>Metrics and auto-scaling for stability and resource efficiency</li><li>Incident drill processes and exercises</li></ul>",
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
					"title":       "Containers & Orchestration",
					"subtitle":    "Cloud-native",
					"description": "Docker / Kubernetes / Helm",
				},
				map[string]any{
					"id":          "s2",
					"title":       "Delivery & Governance",
					"subtitle":    "GitOps",
					"description": "Argo CD / Flux / CI/CD / Blue-Green / Canary",
				},
				map[string]any{
					"id":          "s3",
					"title":       "Observability & Reliability",
					"subtitle":    "SRE",
					"description": "Prometheus / Grafana / Alerting / Capacity / Chaos",
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
					"title":       "Tongji University",
					"major":       "Software Engineering",
					"degree":      "Master's",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "Cloud-native and delivery engineering; platform building and governance.",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "Tongji University",
					"major":       "Computer Science and Technology",
					"degree":      "Bachelor's",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "Distributed systems and engineering; campus projects and competitions.",
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
					"description": `<ul><li>Helm chart collection: <a href="https://helm-charts.zhwei.invalid">helm-charts.zhwei.invalid</a></li><li>GitOps examples: <a href="https://gitops.zhwei.invalid">gitops.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateDevopsPresetEn() []byte {
	devopsByte, _ := json.Marshal(devopsPresetJSONEn)
	return devopsByte
}
