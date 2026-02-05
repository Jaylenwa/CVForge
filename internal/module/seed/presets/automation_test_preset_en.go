package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var automationTestPresetJSONEn = map[string]any{
	"title":    "Automation Testing Engineer Resume",
	"language": "en",
	"Personal": map[string]any{
		"FullName":   "Wei Zhang",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "Automation Testing Engineer",
		"City":       "Shanghai",
		"Money":      "25k-35k",
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
					"description": "<ul><li>5 years of automation testing; proficient with Web/App automation frameworks and best practices</li><li>UI E2E and API automation: Playwright/Cypress/Selenium, pytest, REST/GraphQL</li><li>Stability & maintainability: layered design, data-driven testing, CI/CD quality gates</li><li>Environment prep and Mock/Stub, containerization/parallelism, reporting and visualization</li></ul>",
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
					"subtitle":    "Automation Testing Engineer (Platform/Checkout)",
					"timeStart":   "2021-06",
					"today":       true,
					"description": "<ul><li>Built Web E2E and API automation (Playwright + pytest) for core flows</li><li>Data factory and environment isolation for stable repeatability</li><li>Integrated CI/CD and quality gates to reduce regression efforts and risks</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "Content Community",
					"subtitle":    "Automation Testing Engineer (Platform/Admin)",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-05",
					"description": "<ul><li>Component/page-level automation framework with robust locator & stability strategies</li><li>Parallel execution and visualized reporting for faster feedback</li></ul>",
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
					"title":       "Unified Automation Framework & Case Governance",
					"subtitle":    "Automation Testing Engineer",
					"timeStart":   "2022-05",
					"timeEnd":     "2023-02",
					"description": "<ul><li>Stack: Playwright / pytest / Docker</li><li>Layered architecture and data-driven testing; Mock/Stub & environment isolation</li><li>CI parallel runs and quality gates</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "API Automation & Contract Testing",
					"subtitle":    "Automation Testing Engineer",
					"timeStart":   "2020-09",
					"timeEnd":     "2021-03",
					"description": "<ul><li>Stack: pytest / Pact / REST / GraphQL</li><li>API automation and contract testing to reduce coupling and regression risks</li><li>Data factory and alerting to speed up incident response</li></ul>",
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
					"title":       "Automation Frameworks",
					"subtitle":    "Web/App",
					"description": "Playwright / Cypress / Selenium / Appium",
				},
				map[string]any{
					"id":          "s2",
					"title":       "API & Contract",
					"subtitle":    "API",
					"description": "pytest / REST / GraphQL / Pact",
				},
				map[string]any{
					"id":          "s3",
					"title":       "Engineering & Execution",
					"subtitle":    "CI",
					"description": "Docker / Parallel / Reporting / Quality Gates",
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
					"title":       "University of Science and Technology of China",
					"major":       "Software Engineering",
					"degree":      "Master's",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "Automation testing and engineering; quality platform projects.",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "University of Science and Technology of China",
					"major":       "Computer Science and Technology",
					"degree":      "Bachelor's",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "Software engineering and testing techniques.",
				},
			},
		},
	},
}

func GenerateAutomationTestPresetEn() []byte {
	autoByte, _ := json.Marshal(automationTestPresetJSONEn)
	return autoByte
}
