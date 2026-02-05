package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var opsManagerPresetJSONEn = map[string]any{
	"title":    "Operations Manager/Supervisor Resume",
	"language": "en",
	"Personal": map[string]any{
		"FullName":   "Wei Zhang",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "Operations Manager/Supervisor",
		"City":       "Shanghai",
		"Money":      "40k-60k",
		"JoinTime":   "Available immediately",
		"Gender":     "Male",
		"Age":        "32",
		"Degree":     "Master's",
		"CustomInfo": `[{"label":"Political Status","value":"Party Member"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#f59e0b",
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
					"description": "<ul><li>8 years in operations/platform; led teams to build observability, delivery and SRE</li><li>Stability and cost governance: SLA/SLO/SLI; capacity and cost optimization</li><li>Processes and collaboration: change/release/incident response and drills; cross-team collaboration</li><li>People development: ladder building; policies/reviews; metrics-driven improvements</li></ul>",
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
					"title":       "Retail Platform",
					"subtitle":    "Operations Manager (Platform/SRE)",
					"timeStart":   "2020-03",
					"today":       true,
					"description": "<ul><li>Led GitOps and observability platform building</li><li>Stability metrics and drill mechanisms</li><li>Cross-team collaboration and process improvements</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "Internet Company",
					"subtitle":    "Operations Supervisor (Delivery/Governance)",
					"timeStart":   "2016-07",
					"timeEnd":     "2020-02",
					"description": "<ul><li>Delivery platform and process standards</li><li>Incident response and postmortems</li></ul>",
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
					"title":       "Platform & Stability System",
					"subtitle":    "Operations Manager/Supervisor",
					"timeStart":   "2021-06",
					"timeEnd":     "2022-12",
					"description": "<ul><li>Stack: Argo CD / Helm / Prometheus / Grafana</li><li>Unified delivery and observability governance; stability metrics</li><li>Drills and continuous improvements</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "Process & Organization Capability",
					"subtitle":    "Operations Manager/Supervisor",
					"timeStart":   "2018-08",
					"timeEnd":     "2019-06",
					"description": "<ul><li>Change/release/incident response processes</li><li>Ladder/policy building; metrics-driven improvements</li></ul>",
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
					"title":       "Platform & Delivery",
					"subtitle":    "Engineering",
					"description": "GitOps / Helm / CI/CD",
				},
				map[string]any{
					"id":          "s2",
					"title":       "Stability & Reliability",
					"subtitle":    "SRE",
					"description": "Observability / SLO / SLI / Incident",
				},
				map[string]any{
					"id":          "s3",
					"title":       "Process & Organization",
					"subtitle":    "Governance",
					"description": "Change / Release / Drill / Training",
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
					"title":       "Fudan University",
					"major":       "Software Engineering",
					"degree":      "Master's",
					"timeStart":   "2013-09",
					"timeEnd":     "2016-06",
					"description": "Platform engineering and organizational governance.",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "Fudan University",
					"major":       "Computer Science and Technology",
					"degree":      "Bachelor's",
					"timeStart":   "2009-09",
					"timeEnd":     "2013-06",
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
					"description": `<ul><li>Platform governance handbook: <a href="https://platform-governance.zhwei.invalid">platform-governance.zhwei.invalid</a></li><li>Stability metrics plan: <a href="https://sre-metrics.zhwei.invalid">sre-metrics.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateOpsManagerPresetEn() []byte {
	omByte, _ := json.Marshal(opsManagerPresetJSONEn)
	return omByte
}
