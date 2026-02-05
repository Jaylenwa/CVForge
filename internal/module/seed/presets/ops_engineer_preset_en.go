package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var opsEngineerPresetJSONEn = map[string]any{
	"title":    "Operations Engineer Resume",
	"language": "en",
	"Personal": map[string]any{
		"FullName":   "Wei Zhang",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "Operations Engineer",
		"City":       "Shanghai",
		"Money":      "25k-35k",
		"JoinTime":   "Available immediately",
		"Gender":     "Male",
		"Age":        "28",
		"Degree":     "Master's",
		"CustomInfo": `[{"label":"Political Status","value":"Party Member"}]`,
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
			"title":     "Summary",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "sum1",
					"description": "<ul><li>5 years in operations; Linux, networking and service governance; high availability and incident response</li><li>Observability: Prometheus/Grafana, ELK, logs/metrics/tracing</li><li>Automation & infra: Ansible/Terraform, CI/CD, config/change management</li><li>Capacity and performance optimization with teams: caching/connection pools/concurrency</li></ul>",
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
					"title":       "E-commerce Platform",
					"subtitle":    "Operations Engineer (Stability/Observability)",
					"timeStart":   "2021-06",
					"today":       true,
					"description": "<ul><li>Built observability and alerting tiers to shorten MTTR</li><li>Optimized Nginx/Redis/MySQL resource/connection strategies for peak stability</li><li>Improved release and rollback workflows to reduce risks</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "Content Community",
					"subtitle":    "Operations Engineer (Platform/Delivery)",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-05",
					"description": "<ul><li>CI/CD with progressive delivery for efficient releases</li><li>Unified logs and metrics collection; incident drills and playbooks</li></ul>",
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
					"title":       "Observability Platform & Alert Governance",
					"subtitle":    "Operations Engineer",
					"timeStart":   "2022-05",
					"timeEnd":     "2023-02",
					"description": "<ul><li>Stack: Prometheus / Grafana / ELK</li><li>Unified collection and alert governance for faster incident response</li><li>SLO/SLI metrics to drive stability improvements</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "Release & Change Management",
					"subtitle":    "Operations Engineer",
					"timeStart":   "2020-09",
					"timeEnd":     "2021-03",
					"description": "<ul><li>Stack: GitLab CI / Ansible / Terraform</li><li>Release process and change tracking to reduce deployment risks</li><li>Quality gates and rollback strategies</li></ul>",
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
					"title":       "Systems & Networking",
					"subtitle":    "Foundation",
					"description": "Linux / Shell / TCP/IP / Nginx / HAProxy",
				},
				map[string]any{
					"id":          "s2",
					"title":       "Databases & Cache",
					"subtitle":    "Storage",
					"description": "MySQL / PostgreSQL / Redis",
				},
				map[string]any{
					"id":          "s3",
					"title":       "Observability & Automation",
					"subtitle":    "Engineering",
					"description": "Prometheus / Grafana / ELK / Ansible / Terraform / CI/CD",
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
					"title":       "South China University of Technology",
					"major":       "Software Engineering",
					"degree":      "Master's",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "Observability and operations engineering; platform building and stability.",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "South China University of Technology",
					"major":       "Computer Science and Technology",
					"degree":      "Bachelor's",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "Systems and networking; campus projects and competitions.",
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
					"description": `<ul><li>Ops toolkit: <a href="https://ops-toolkit.zhwei.invalid">ops-toolkit.zhwei.invalid</a></li><li>Tech blog: <a href="https://blog.zhwei.invalid">blog.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateOpsEngineerPresetEn() []byte {
	opsByte, _ := json.Marshal(opsEngineerPresetJSONEn)
	return opsByte
}
