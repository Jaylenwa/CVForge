package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var algoEngineerPresetJSONEn = map[string]any{
	"title":    "Algorithm Engineer Resume",
	"language": "en",
	"Personal": map[string]any{
		"FullName":   "Qiang Wang",
		"Email":      "wangqiang@example.com",
		"Phone":      "13700000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "Algorithm Engineer",
		"City":       "Beijing",
		"Money":      "35k-55k",
		"JoinTime":   "Available immediately",
		"Gender":     "Male",
		"Age":        "27",
		"Degree":     "Master's",
		"CustomInfo": `[{"label":"Political Status","value":"Mass"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#ef6c00",
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
					"description": "<ul><li>4 years in algorithm engineering; strong data structures and algorithms</li><li>Search ranking and recommendation; routing and graph algorithms</li><li>C++/Java and engineering; performance optimization and concurrency</li><li>Modeling and evaluation; metrics and A/B experiments</li></ul>",
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
					"title":       "Search Company",
					"subtitle":    "Algorithm Engineer (Ranking/Recommendation)",
					"timeStart":   "2021-06",
					"today":       true,
					"description": "<ul><li>Ranking and recommendation strategy improvements</li><li>Features and evaluation systems; A/B experiments</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "Logistics Company",
					"subtitle":    "Algorithm Engineer (Routing/Graph)",
					"timeStart":   "2019-08",
					"timeEnd":     "2021-05",
					"description": "<ul><li>Routing and graph algorithms; route and cost optimization</li><li>Engineering and service governance; performance and stability</li></ul>",
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
					"title":       "Search Ranking & Recommendation",
					"subtitle":    "Algorithm Engineer",
					"timeStart":   "2022-03",
					"timeEnd":     "2022-11",
					"description": "<ul><li>Stack: C++ / Java / Feature / AB</li><li>CTR/CVR improvements with ranking and recommendation</li><li>Feature and evaluation system building</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "Routing & Graph Algorithms",
					"subtitle":    "Algorithm Engineer",
					"timeStart":   "2020-06",
					"timeEnd":     "2021-02",
					"description": "<ul><li>Stack: Graph / Shortest Path / Heuristic</li><li>Route and cost optimization; efficiency and stability</li></ul>",
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
					"title":       "Data Structures & Algorithms",
					"subtitle":    "Foundation",
					"description": "Array / LinkedList / Tree / Graph / DP",
				},
				map[string]any{
					"id":          "s2",
					"title":       "Engineering & Optimization",
					"subtitle":    "Performance",
					"description": "C++ / Java / Concurrency / Profiling",
				},
				map[string]any{
					"id":          "s3",
					"title":       "Modeling & Evaluation",
					"subtitle":    "Metrics",
					"description": "Feature / Metrics / AB Test",
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
					"title":       "Beihang University",
					"major":       "Computer Science and Technology",
					"degree":      "Master's",
					"timeStart":   "2016-09",
					"timeEnd":     "2019-06",
					"description": "Algorithms and engineering; competitions and projects.",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "Beihang University",
					"major":       "Software Engineering",
					"degree":      "Bachelor's",
					"timeStart":   "2012-09",
					"timeEnd":     "2016-06",
					"description": "Data structures and algorithms; campus projects.",
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
					"description": `<ul><li>Algorithm notes: <a href="https://algo-notes.wq.invalid">algo-notes.wq.invalid</a></li><li>Ranking demos: <a href="https://ranking-demo.wq.invalid">ranking-demo.wq.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateAlgoEngineerPresetEn() []byte {
	aeByte, _ := json.Marshal(algoEngineerPresetJSONEn)
	return aeByte
}
