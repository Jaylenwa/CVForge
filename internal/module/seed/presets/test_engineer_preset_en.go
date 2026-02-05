package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var testEngineerPresetJSONEn = map[string]any{
	"title":    "Testing Engineer Resume",
	"language": "en",
	"Personal": map[string]any{
		"FullName":   "Wei Zhang",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "Testing Engineer",
		"City":       "Shanghai",
		"Money":      "20k-30k",
		"JoinTime":   "Available immediately",
		"Gender":     "Male",
		"Age":        "28",
		"Degree":     "Master's",
		"CustomInfo": `[{"label":"Political Status","value":"Party Member"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#ef4444",
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
					"description": "<ul><li>5 years of software testing experience; familiar with testing methodologies, able to design plans and test cases</li><li>Case design (boundary/equivalence/decision table/cause-effect), defect management and quality metrics</li><li>Functional/regression/compatibility/security testing; cross-device/browser adaptation experience</li><li>Basic automation and scripting skills; collaborate with engineering to ensure quality</li></ul>",
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
					"subtitle":    "Testing Engineer (Product/Checkout)",
					"timeStart":   "2021-06",
					"today":       true,
					"description": "<ul><li>Owned testing for product and checkout; improved case library and smoke/regression strategy</li><li>Built defect management and quality dashboards to improve tracking and decisions</li><li>Collaborated to optimize release processes and reduce incidents</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "Content Community",
					"subtitle":    "Testing Engineer (Platform/Admin)",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-05",
					"description": "<ul><li>Tested permissions/messaging platforms; promoted smoke/regression automation</li><li>Established compatibility testing matrix and adaptation solution</li></ul>",
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
					"title":       "Quality Metrics & Defect Governance",
					"subtitle":    "Testing Engineer",
					"timeStart":   "2022-05",
					"timeEnd":     "2023-02",
					"description": "<ul><li>Stack: Jira / Confluence / in-house dashboards</li><li>Established metrics for defects and cases; improved visibility and decisions</li><li>Defect classification and governance to shorten fix cycles</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "Compatibility & Adaptation Matrix",
					"subtitle":    "Testing Engineer",
					"timeStart":   "2020-09",
					"timeEnd":     "2021-03",
					"description": "<ul><li>Stack: BrowserStack / device farms</li><li>Built cross-browser/device matrix; strengthened cross-end testing</li><li>Significantly reduced compatibility incidents</li></ul>",
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
					"title":       "Methods & Design",
					"subtitle":    "Theory",
					"description": "Case design / Boundary / Equivalence / Decision Table / Cause-Effect",
				},
				map[string]any{
					"id":          "s2",
					"title":       "Types & Process",
					"subtitle":    "Practice",
					"description": "Functional / Regression / Compatibility / Security / Exploratory",
				},
				map[string]any{
					"id":          "s3",
					"title":       "Quality & Collaboration",
					"subtitle":    "Governance",
					"description": "Defect management / Metrics / Dashboards / Cross-team collaboration",
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
					"title":       "Nanjing University",
					"major":       "Software Engineering",
					"degree":      "Master's",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "Research on software quality and testing engineering.",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "Nanjing University",
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

func GenerateTestEngineerPresetEn() []byte {
	testByte, _ := json.Marshal(testEngineerPresetJSONEn)
	return testByte
}
