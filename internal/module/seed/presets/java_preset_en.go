package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var javaPresetJSONEn = map[string]any{
	"title":    "Java Backend Developer Resume",
	"language": "en",
	"Personal": map[string]any{
		"FullName":   "Wei Zhang",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "Java Backend Developer",
		"City":       "Shanghai",
		"Money":      "25k-35k",
		"JoinTime":   "Available immediately",
		"Gender":     "Male",
		"Age":        "28",
		"Degree":     "Master's",
		"CustomInfo": `[{"label":"Political Status","value":"Party Member"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#f97316",
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
					"description": "<ul><li>5 years of Java backend development experience; proficient with Spring Boot / Spring Cloud, with solid performance tuning and reliability practices for high-concurrency systems</li><li>Hands-on with MySQL/Redis/Kafka; experienced in cache consistency, reliable messaging, distributed transactions, and idempotent design</li><li>Engineering and cloud-native practices: Docker/Kubernetes, CI/CD, metrics/logging/tracing</li></ul>",
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
					"title":       "Leading Internet Company",
					"subtitle":    "Java Backend Developer (Trading/Payments)",
					"timeStart":   "2021-07",
					"today":       true,
					"description": "<ul><li>Owned core trading-domain APIs and async flows (Spring Boot + Kafka), raising order creation success rate to 99.95%</li><li>Implemented multi-layer Redis caching and hot-key mitigation, reducing P95 latency from 260ms to 150ms</li><li>Designed idempotency, retry, and degradation strategies, significantly lowering error rates during peak traffic</li><li>Improved observability (Prometheus/Grafana + logging/tracing), cutting incident triage time by 50%</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "Enterprise SaaS Company",
					"subtitle":    "Java Engineer (Platform/Middle Platform)",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-06",
					"description": "<ul><li>Built unified IAM and audit platform (Spring Boot + MySQL) serving 100+ enterprise customers</li><li>Optimized slow SQL and indexing strategies, reducing critical query latency by 60%</li><li>Introduced canary release and rollback mechanisms, improving release success rate and lowering outage risk</li></ul>",
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
					"description": "Focused on distributed systems and highly available architectures; participated in lab platform projects and performance optimization.",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "Zhejiang University",
					"major":       "Software Engineering",
					"degree":      "Bachelor's",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "Studied data structures, operating systems, databases, and computer networks; participated in campus projects and competitions.",
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
					"subtitle":    "Java",
					"description": "Spring Boot / Spring Cloud / MyBatis / REST API",
				},
				map[string]any{
					"id":          "s2",
					"title":       "Storage & Middleware",
					"subtitle":    "Database / Cache / Messaging",
					"description": "MySQL / Redis / Kafka / Elasticsearch",
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
			"id":        "projects",
			"type":      common.ResumeSectionTypeProjects,
			"title":     "Projects",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "p1",
					"title":       "Unified Order and Inventory Platform",
					"subtitle":    "Java Backend Developer",
					"timeStart":   "2022-05",
					"timeEnd":     "2023-02",
					"description": "<ul><li>Tech stack: Spring Boot / MySQL / Redis / Kafka</li><li>Abstracted core order/inventory domain models and APIs to enable reuse across business lines and reduce duplicate development</li><li>Introduced event-driven and eventual-consistency patterns, reducing cross-service coupling and improving evolvability</li><li>Built load-testing and capacity planning processes, doubling peak QPS while keeping the system stable</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "High-Concurrency Coupon System",
					"subtitle":    "Java Backend Developer",
					"timeStart":   "2020-09",
					"timeEnd":     "2021-03",
					"description": "<ul><li>Tech stack: Java / Redis / Kafka</li><li>Designed idempotent inventory deduction and issuance flows to prevent duplicate claims and overselling</li><li>Used Redis pre-warming and Lua atomic scripts, reducing P95 latency by 40% on critical APIs</li><li>Added async peak-shaving and retry mechanisms to ensure stability during major promotions</li></ul>",
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
					"description": `<ul><li>Code repository: <a href="https://code.zhwei.invalid">code.zhwei.invalid</a></li><li>Tech blog: <a href="https://blog.zhwei.invalid">blog.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateJavaPresetEn() []byte {
	javaByte, _ := json.Marshal(javaPresetJSONEn)
	return javaByte
}
