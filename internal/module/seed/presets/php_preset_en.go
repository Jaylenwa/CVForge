package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var phpPresetJSONEn = map[string]any{
	"title":    "PHP Backend Developer Resume",
	"language": "en",
	"Personal": map[string]any{
		"FullName":   "Wei Zhang",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "PHP Backend Developer",
		"City":       "Shanghai / Hangzhou",
		"Money":      "25k-35k · 14-month salary",
		"JoinTime":   "Available immediately",
		"Gender":     "Male",
		"Age":        "28",
		"Degree":     "Master's",
		"CustomInfo": `[{"label":"Political Status","value":"Party Member"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#7c3aed",
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
					"description": "<ul><li>5 years of PHP backend experience; familiar with Laravel/Symfony and domain modeling, able to independently deliver features and optimize performance</li><li>Hands-on with MySQL/Redis/Elasticsearch and Kafka/RabbitMQ; experienced in cache consistency, reliable messaging, idempotency, and retry design</li><li>Engineering and cloud-native practices: Docker/Kubernetes, CI/CD, monitoring/alerting, and tracing</li></ul>",
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
					"title":       "Leading Local Services Platform",
					"subtitle":    "PHP Backend Developer (Orders/Fulfillment)",
					"timeStart":   "2021-06",
					"today":       true,
					"description": "<ul><li>Owned core order and fulfillment APIs (Laravel + MySQL + Redis); implemented rate limiting/circuit breaking/degradation, significantly improving peak stability</li><li>Built multi-layer Redis caching and hot-key mitigation, reducing P95 latency from 220ms to 140ms</li><li>Introduced Kafka-based async payment callbacks and notification flows, reducing critical-path latency and improving success rate</li><li>Improved monitoring/alerting and logging standards, cutting incident triage time by 50%</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "Enterprise SaaS Company",
					"subtitle":    "PHP Engineer (Platform/Middle Platform)",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-05",
					"description": "<ul><li>Built unified IAM and audit platform (Symfony + MySQL) serving 100+ enterprise customers</li><li>Optimized slow SQL and indexing strategies, reducing critical query latency by 60%</li><li>Established API standards and automated testing (OpenAPI + PHPUnit), improving regression efficiency</li></ul>",
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
					"title":       "High-Concurrency Flash Sale and Inventory System",
					"subtitle":    "PHP Backend Developer",
					"timeStart":   "2022-04",
					"timeEnd":     "2022-12",
					"description": "<ul><li>Tech stack: PHP / Laravel / Redis / MySQL / Kafka</li><li>Designed idempotent inventory reservation and order creation to prevent duplicate orders and overselling</li><li>Used Redis Lua atomic scripts and hot-key governance, reducing P95 latency by 40% during peaks</li><li>Implemented message retries and compensation jobs to ensure eventual consistency and recoverability</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "Unified API Gateway and Authorization Component",
					"subtitle":    "PHP Backend Developer",
					"timeStart":   "2020-10",
					"timeEnd":     "2021-03",
					"description": "<ul><li>Tech stack: PHP / Nginx / Redis / JWT</li><li>Abstracted auth, signature validation, rate limiting, and canary release capabilities to speed up service onboarding and reduce duplicated work</li><li>Implemented permission caching and audit logs to meet compliance and troubleshooting needs</li></ul>",
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
					"subtitle":    "PHP",
					"description": "Laravel / Symfony / Composer / REST API",
				},
				map[string]any{
					"id":          "s2",
					"title":       "Storage & Middleware",
					"subtitle":    "Database / Cache / Messaging",
					"description": "MySQL / Redis / Elasticsearch / Kafka / RabbitMQ",
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
					"title":       "East China University of Science and Technology",
					"major":       "Computer Science and Technology",
					"degree":      "Master's",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "Focused on distributed systems and high-concurrency services; participated in lab platform projects and performance optimization.",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "East China University of Science and Technology",
					"major":       "Computer Science and Technology",
					"degree":      "Bachelor's",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "Studied data structures, operating systems, and databases; participated in campus projects and competitions.",
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
					"description": `<ul><li>Open-source packages: <a href="https://packagist.zhwei.invalid">packagist.zhwei.invalid</a></li><li>Tech blog: <a href="https://blog.zhwei.invalid">blog.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GeneratePHPPresetEn() []byte {
	phpByte, _ := json.Marshal(phpPresetJSONEn)
	return phpByte
}
