package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var testDevPresetJSONEn = map[string]any{
	"title":    "Testing Development Engineer Resume",
	"language": "en",
	"Personal": map[string]any{
		"FullName":   "Wei Zhang",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "Testing Development Engineer",
		"City":       "Shanghai",
		"Money":      "30k-40k",
		"JoinTime":   "Available immediately",
		"Gender":     "Male",
		"Age":        "28",
		"Degree":     "Master's",
		"CustomInfo": `[{"label":"Political Status","value":"Party Member"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#8b5cf6",
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
					"description": "<ul><li>5 years of test development; experienced in building testing frameworks and platforms</li><li>Frontend/backend automation, API contracts, data factory and environment management</li><li>Quality platforms, case governance, coverage and metrics dashboards</li><li>CI/CD with quality gates, stability and observability</li></ul>",
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
					"title":       "Platform Middleware",
					"subtitle":    "Testing Development Engineer (Platform/Quality)",
					"timeStart":   "2021-06",
					"today":       true,
					"description": "<ul><li>Built unified testing platform/framework; layered architecture and case governance</li><li>Data factory, environment isolation, parallel execution for efficiency and stability</li><li>CI/CD quality gates with metrics and dashboards</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "E-commerce Middleware",
					"subtitle":    "Testing Development Engineer (API/Contract)",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-05",
					"description": "<ul><li>API automation and contract testing platform to reduce coupling and regression risks</li><li>Data generation and simulation to speed up debugging and fixes</li></ul>",
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
					"title":       "Unified Testing Platform & Quality Dashboard",
					"subtitle":    "Testing Development Engineer",
					"timeStart":   "2022-05",
					"timeEnd":     "2023-02",
					"description": "<ul><li>Stack: Go / Python / Playwright / Grafana</li><li>Unified platform with case governance; coverage and quality metrics</li><li>Alerting and quality gates to ensure delivery quality</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "Contract Testing Platform & Data Factory",
					"subtitle":    "Testing Development Engineer",
					"timeStart":   "2020-09",
					"timeEnd":     "2021-03",
					"description": "<ul><li>Stack: Pact / pytest / Postgres / Redis</li><li>Contract tests, data generation/isolation to reduce coupling and risks</li><li>Reporting and diagnostic tools for faster incident response</li></ul>",
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
					"title":       "Languages & Frameworks",
					"subtitle":    "Test Dev",
					"description": "Go / Python / Playwright / pytest / Pact",
				},
				map[string]any{
					"id":          "s2",
					"title":       "Platforms & Engineering",
					"subtitle":    "Quality",
					"description": "Testing Platform / Data Factory / Parallel Execution / Metrics",
				},
				map[string]any{
					"id":          "s3",
					"title":       "CI/CD & Governance",
					"subtitle":    "Delivery",
					"description": "CI/CD / Quality Gates / Coverage / Alerting & Observability",
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
					"title":       "Shanghai Jiao Tong University",
					"major":       "Software Engineering",
					"degree":      "Master's",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "Test engineering and platform building; quality systems and metrics.",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "Shanghai Jiao Tong University",
					"major":       "Computer Science and Technology",
					"degree":      "Bachelor's",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "Software engineering and systems development.",
				},
			},
		},
	},
}

func GenerateTestDevPresetEn() []byte {
	tdByte, _ := json.Marshal(testDevPresetJSONEn)
	return tdByte
}
