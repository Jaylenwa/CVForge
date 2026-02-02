package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var csharpPresetJSONEn = map[string]any{
	"title":    "C# Backend Developer Resume",
	"language": "en",
	"Personal": map[string]any{
		"FullName":   "Wei Zhang",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "C# Backend Developer",
		"City":       "Shanghai / Hangzhou",
		"Money":      "25k-40k · 14-month salary",
		"JoinTime":   "Available immediately",
		"Gender":     "Male",
		"Age":        "28",
		"Degree":     "Master's",
		"CustomInfo": `[{"label":"Political Status","value":"Party Member"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#22c55e",
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
					"description": "<ul><li>5 years of C# backend experience; proficient with ASP.NET Core, microservices, and high-concurrency optimization, with strong reliability and observability practices</li><li>Hands-on with SQL Server/PostgreSQL, Redis, and Kafka; experienced in cache consistency, reliable messaging, distributed transactions, and idempotent design</li><li>Strong engineering skills: CI/CD, automated testing, and containerized deployment (Docker/Kubernetes)</li></ul>",
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
					"title":       "Enterprise SaaS Company",
					"subtitle":    "C# Backend Developer (Multi-tenant Platform)",
					"timeStart":   "2021-05",
					"today":       true,
					"description": "<ul><li>Owned multi-tenant core services (ASP.NET Core + EF Core); designed tenant isolation and permission model to support 100+ enterprise customers</li><li>Implemented Redis caching and hot-key mitigation, reducing P95 latency from 240ms to 150ms</li><li>Introduced Kafka-based async audit and notification flows, reducing critical-path latency and improving reliability</li><li>Improved monitoring/alerting and tracing (Prometheus/Grafana), cutting incident triage time by 50%</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "Retail Tech Company",
					"subtitle":    "C# Engineer (Orders/Payments)",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-04",
					"description": "<ul><li>Worked on order and payment system refactoring (ASP.NET Core + SQL Server); introduced idempotency and retry strategies, significantly reducing error rates</li><li>Optimized slow SQL and indexing strategies, reducing critical query time by 60%</li><li>Established API standards and automated testing (xUnit + Postman Collections), improving regression efficiency</li></ul>",
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
					"title":       "Unified Authorization and Audit Platform",
					"subtitle":    "C# Backend Developer",
					"timeStart":   "2022-03",
					"timeEnd":     "2022-11",
					"description": "<ul><li>Tech stack: ASP.NET Core / EF Core / Redis / PostgreSQL</li><li>Designed a hybrid RBAC/ABAC model and permission cache strategy to balance consistency and performance</li><li>Implemented async ingestion and query aggregation for audit events to meet compliance and troubleshooting needs</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "Gateway and Service Governance Components",
					"subtitle":    "C# Backend Developer",
					"timeStart":   "2020-08",
					"timeEnd":     "2021-02",
					"description": "<ul><li>Tech stack: .NET / Ocelot / Redis</li><li>Implemented auth, rate limiting, canary release, circuit breaking, and unified logging to improve cross-team collaboration</li><li>Built configuration-driven governance and rollback mechanisms to reduce release risk</li></ul>",
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
					"subtitle":    "C#",
					"description": "ASP.NET Core / EF Core / Dapper / REST API / gRPC",
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
					"title":       "Tongji University",
					"major":       "Software Engineering",
					"degree":      "Master's",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "Focused on enterprise application architecture and distributed systems; participated in lab platform projects and performance optimization.",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "Tongji University",
					"major":       "Software Engineering",
					"degree":      "Bachelor's",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "Studied object-oriented design, databases, and distributed systems; participated in campus project practice.",
				},
			},
		},
	},
}

func GenerateCSharpPresetEn() []byte {
	csharpByte, _ := json.Marshal(csharpPresetJSONEn)
	return csharpByte
}
