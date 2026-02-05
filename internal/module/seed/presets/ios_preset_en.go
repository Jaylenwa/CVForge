package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var iosPresetJSONEn = map[string]any{
	"title":    "iOS Developer Resume",
	"language": "en",
	"Personal": map[string]any{
		"FullName":   "Wei Zhang",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "iOS Developer",
		"City":       "Shanghai",
		"Money":      "25k-40k",
		"JoinTime":   "Available immediately",
		"Gender":     "Male",
		"Age":        "28",
		"Degree":     "Master's",
		"CustomInfo": `[{"label":"Political Status","value":"Party Member"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#1e90ff",
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
					"description": "<ul><li>5 years of iOS development; proficient in Swift/Objective-C, UIKit/SwiftUI, Combine</li><li>Networking & concurrency (URLSession/Alamofire/GCD/Async-Await), persistence (Core Data)</li><li>Performance & reliability: startup optimization, memory leak fixes, crash analytics</li><li>Engineering & delivery: SPM/CocoaPods, unit/UI tests, TestFlight/enterprise release</li></ul>",
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
					"title":       "E-commerce App",
					"subtitle":    "iOS Developer (Product/Orders)",
					"timeStart":   "2021-06",
					"today":       true,
					"description": "<ul><li>Owned product and order modules with MVVM + Combine</li><li>Optimized rendering and interactions to reduce frame drops</li><li>Crash and error reporting (Crashlytics/Sentry) for faster incident response</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "Content Community App",
					"subtitle":    "iOS Developer (Account/Messaging)",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-05",
					"description": "<ul><li>Refactored modules; introduced SPM and graphics optimizations</li><li>Built CI/CD and automated testing workflows for reliable releases</li></ul>",
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
					"title":       "SwiftUI Components & Theming",
					"subtitle":    "iOS Developer",
					"timeStart":   "2022-05",
					"timeEnd":     "2023-02",
					"description": "<ul><li>Stack: Swift / SwiftUI / Combine</li><li>Abstracted UI components and theming for reuse and consistency</li><li>Optimized rendering and state management</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "Offline Capabilities & Data Sync",
					"subtitle":    "iOS Developer",
					"timeStart":   "2020-09",
					"timeEnd":     "2021-03",
					"description": "<ul><li>Stack: Core Data / Background Tasks / URLSession</li><li>Implemented offline caching, background sync, retries and conflict handling</li><li>Measured UX improvements via metrics</li></ul>",
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
					"subtitle":    "iOS",
					"description": "Swift / Objective-C / UIKit / SwiftUI / Combine",
				},
				map[string]any{
					"id":          "s2",
					"title":       "Networking & Concurrency",
					"subtitle":    "Data",
					"description": "URLSession / Alamofire / GCD / Async-Await / Core Data",
				},
				map[string]any{
					"id":          "s3",
					"title":       "Engineering & Delivery",
					"subtitle":    "Release",
					"description": "SPM / CocoaPods / Unit Tests / UI Tests / TestFlight",
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
					"title":       "Tsinghua University",
					"major":       "Computer Science and Technology",
					"degree":      "Master's",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "Focus on iOS architecture and performance optimization.",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "Tsinghua University",
					"major":       "Software Engineering",
					"degree":      "Bachelor's",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "Studied mobile development, compilation/linking.",
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
					"description": `<ul><li>iOS toolkit: <a href="https://ioskit.zhwei.invalid">ioskit.zhwei.invalid</a></li><li>Tech blog: <a href="https://blog.zhwei.invalid">blog.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateIOSPresetEn() []byte {
	iosByte, _ := json.Marshal(iosPresetJSONEn)
	return iosByte
}
