package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var performanceTestPresetJSONEn = map[string]any{
	"title":    "Performance Testing Engineer Resume",
	"language": "en",
	"Personal": map[string]any{
		"FullName":   "Wei Zhang",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "Performance Testing Engineer",
		"City":       "Shanghai",
		"Money":      "25k-35k",
		"JoinTime":   "Available immediately",
		"Gender":     "Male",
		"Age":        "28",
		"Degree":     "Master's",
		"CustomInfo": `[{"label":"Political Status","value":"Party Member"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#f43f5e",
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
					"description": "<ul><li>5 years of performance testing; pressure/capacity/stability tests and bottleneck analysis</li><li>Tools & frameworks: JMeter/Gatling/Locust; scenario scripting and data preparation</li><li>Metrics & observability: TPS/RT/error rate/resource usage, Prometheus/Grafana/APM</li><li>Performance governance & optimization: caching/concurrency/connection pools/async & degradation</li></ul>",
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
					"subtitle":    "Performance Testing Engineer (Checkout/Orders)",
					"timeStart":   "2021-06",
					"today":       true,
					"description": "<ul><li>Designed pressure tests for checkout/order flows; capacity assessment and bottleneck analysis</li><li>Scenario/data preparation, script reuse and result analysis</li><li>Drove optimizations to improve stability and throughput</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "Content Community",
					"subtitle":    "Performance Testing Engineer (Platform/Middleware)",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-05",
					"description": "<ul><li>Built pressure test plans for growth/messaging; improved metrics and alerting</li><li>Collaborated to optimize resources and concurrency strategies</li></ul>",
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
					"title":       "Unified Pressure Testing Platform & Reporting",
					"subtitle":    "Performance Testing Engineer",
					"timeStart":   "2022-05",
					"timeEnd":     "2023-02",
					"description": "<ul><li>Stack: JMeter / Grafana / Prometheus</li><li>Scenario management & data preparation; reporting and alerting</li><li>APM and tracing integration for bottleneck diagnosis</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "Capacity Assessment & Stability Governance",
					"subtitle":    "Performance Testing Engineer",
					"timeStart":   "2020-09",
					"timeEnd":     "2021-03",
					"description": "<ul><li>Stack: JMeter / Locust / APM</li><li>Capacity evaluation and stability tests; optimization proposals</li><li>Data management and reuse for pressure tests</li></ul>",
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
					"title":       "Tools & Frameworks",
					"subtitle":    "Pressure",
					"description": "JMeter / Gatling / Locust",
				},
				map[string]any{
					"id":          "s2",
					"title":       "Metrics & Observability",
					"subtitle":    "Monitoring",
					"description": "TPS / RT / Error Rate / Resource Usage / Prometheus / Grafana / APM",
				},
				map[string]any{
					"id":          "s3",
					"title":       "Governance & Optimization",
					"subtitle":    "Stability",
					"description": "Caching / Concurrency / Connection Pools / Async / Degradation",
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
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "Performance and stability engineering; pressure testing platform projects.",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "Zhejiang University",
					"major":       "Computer Science and Technology",
					"degree":      "Bachelor's",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "Operating systems and networking; campus projects and competitions.",
				},
			},
		},
	},
}

func GeneratePerformanceTestPresetEn() []byte {
	perfByte, _ := json.Marshal(performanceTestPresetJSONEn)
	return perfByte
}
