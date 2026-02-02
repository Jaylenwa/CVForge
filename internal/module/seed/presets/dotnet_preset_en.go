package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var dotnetPresetJSONEn = map[string]any{
	"title":    ".NET Backend Developer Resume",
	"language": "en",
	"Personal": map[string]any{
		"FullName":   "Wei Zhang",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        ".NET Backend Developer",
		"City":       "Shanghai / Hangzhou",
		"Money":      "28k-45k · 14-month salary",
		"JoinTime":   "Available immediately",
		"Gender":     "Male",
		"Age":        "28",
		"Degree":     "Master's",
		"CustomInfo": `[{"label":"Political Status","value":"Party Member"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#e11d48",
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
					"description": "<ul><li>5 years of .NET backend experience; proficient with ASP.NET Core, microservices, and highly available architectures, with strong performance tuning and reliability practices</li><li>Familiar with gRPC, message queues, and event-driven architecture; experienced in idempotency, retries, degradation, and eventual consistency</li><li>Cloud-native and engineering practices: Docker/Kubernetes, CI/CD, metrics/logging/tracing</li></ul>",
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
					"title":       "Enterprise Digital Transformation Team",
					"subtitle":    ".NET Backend Developer (Platform/Middle Platform)",
					"timeStart":   "2021-08",
					"today":       true,
					"description": "<ul><li>Owned platform core services (ASP.NET Core + gRPC); led service decomposition and API governance to improve iteration speed and stability</li><li>Implemented distributed caching and read/write splitting, reducing P95 latency from 260ms to 160ms</li><li>Built release and rollback pipelines (CI/CD), improving release success rate and lowering outage risk</li><li>Improved observability (Prometheus/Grafana + tracing), cutting incident triage time by 50%</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "Logistics Technology Company",
					"subtitle":    ".NET Engineer (Dispatching/Billing)",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-07",
					"description": "<ul><li>Worked on dispatching and billing systems (.NET + SQL Server + Redis), supporting multi-city business expansion</li><li>Optimized batch processing and reporting pipelines, reducing job duration by 40%</li><li>Improved exception handling and alerting strategies, significantly reducing night-time incidents</li></ul>",
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
					"title":       "Unified Order and Settlement Platform",
					"subtitle":    ".NET Backend Developer",
					"timeStart":   "2022-06",
					"timeEnd":     "2023-03",
					"description": "<ul><li>Tech stack: ASP.NET Core / gRPC / Redis / PostgreSQL / Kafka</li><li>Abstracted core order/settlement domain models and APIs to enable reuse across business lines and reduce duplicate development</li><li>Introduced event-driven and eventual-consistency patterns, reducing cross-service coupling and improving evolvability</li><li>Built load testing and capacity planning processes, doubling peak QPS while keeping the system stable</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "Unified Identity and Authorization Service",
					"subtitle":    ".NET Backend Developer",
					"timeStart":   "2020-09",
					"timeEnd":     "2021-02",
					"description": "<ul><li>Tech stack: .NET / JWT / Redis</li><li>Implemented token strategy, session management, and permission caching to improve auth performance and operability</li><li>Added audit and alerting capabilities to meet compliance requirements</li></ul>",
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
					"title":       "Platforms & Frameworks",
					"subtitle":    ".NET",
					"description": "ASP.NET Core / gRPC / EF Core / Dapper / REST API",
				},
				map[string]any{
					"id":          "s2",
					"title":       "Storage & Middleware",
					"subtitle":    "Database / Cache / Messaging",
					"description": "SQL Server / PostgreSQL / Redis / Kafka",
				},
				map[string]any{
					"id":          "s3",
					"title":       "Engineering & Reliability",
					"subtitle":    "Delivery & Observability",
					"description": "Docker / Kubernetes / CI/CD / Prometheus / Grafana",
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
					"title":       "Harbin Institute of Technology",
					"major":       "Software Engineering",
					"degree":      "Master's",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "Focused on distributed systems and engineering practices; participated in lab platform projects and performance optimization.",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "Harbin Institute of Technology",
					"major":       "Software Engineering",
					"degree":      "Bachelor's",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "Studied software engineering and distributed systems; participated in lab platform projects and performance optimization practice.",
				},
			},
		},
	},
}

func GenerateDotnetPresetEn() []byte {
	dotnetByte, _ := json.Marshal(dotnetPresetJSONEn)
	return dotnetByte
}
