package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var androidPresetJSONEn = map[string]any{
	"title":    "Android Developer Resume",
	"language": "en",
	"Personal": map[string]any{
		"FullName":   "Wei Zhang",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "Android Developer",
		"City":       "Shanghai",
		"Money":      "25k-40k",
		"JoinTime":   "Available immediately",
		"Gender":     "Male",
		"Age":        "28",
		"Degree":     "Master's",
		"CustomInfo": `[{"label":"Political Status","value":"Party Member"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#3ddc84",
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
					"description": "<ul><li>5 years of Android development; proficient in Kotlin/Java, Jetpack, Compose/MVVM</li><li>Concurrency & data flow (Coroutines/Flow), networking & storage (Retrofit/Room/OKHttp)</li><li>Performance & reliability: cold start, memory & battery, crash analytics</li><li>Engineering & DI: Hilt/Dagger, Gradle, modularization, multi-channel release</li></ul>",
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
					"title":       "Online Education App",
					"subtitle":    "Android Developer (Course/Payment)",
					"timeStart":   "2021-06",
					"today":       true,
					"description": "<ul><li>Owned course & payment modules with MVVM + Jetpack</li><li>Optimized cold start and rendering; reduced homepage launch time by 30%</li><li>Established crash analytics and tracking (Crashlytics + in-house SDK)</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "Content Community App",
					"subtitle":    "Android Developer (Account/Messaging)",
					"timeStart":   "2019-07",
					"timeEnd":     "2021-05",
					"description": "<ul><li>Refactored account and messaging modules with Hilt and coroutines</li><li>Improved networking & caching; reduced P95 latency by 40%</li><li>Built multi-module and multi-channel release pipelines</li></ul>",
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
					"title":       "Compose Component Migration",
					"subtitle":    "Android Developer",
					"timeStart":   "2022-05",
					"timeEnd":     "2023-02",
					"description": "<ul><li>Stack: Kotlin / Jetpack Compose / Hilt / Room</li><li>Common UI components and state management for reuse and consistency</li><li>Snapshot & diff rendering for better UX and performance</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "Offline Caching & Weak Network Optimization",
					"subtitle":    "Android Developer",
					"timeStart":   "2020-09",
					"timeEnd":     "2021-03",
					"description": "<ul><li>Stack: OKHttp / Retrofit / Room / WorkManager</li><li>Implemented resource caching, resume downloads, retries</li><li>Measured improvements via metrics & tracking</li></ul>",
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
					"subtitle":    "Android",
					"description": "Kotlin / Java / Jetpack / Compose / MVVM",
				},
				map[string]any{
					"id":          "s2",
					"title":       "Networking & Storage",
					"subtitle":    "Data",
					"description": "Retrofit / OKHttp / Room / SQLite",
				},
				map[string]any{
					"id":          "s3",
					"title":       "Engineering & Release",
					"subtitle":    "Delivery",
					"description": "Gradle / Modularization / Multi-channel / CI/CD / Crashlytics",
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
					"title":       "Peking University",
					"major":       "Computer Science and Technology",
					"degree":      "Master's",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "Focus on mobile architecture and performance optimization.",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "Peking University",
					"major":       "Software Engineering",
					"degree":      "Bachelor's",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "Studied mobile development, networks, operating systems.",
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
					"description": `<ul><li>Android toolkit: <a href="https://androidkit.zhwei.invalid">androidkit.zhwei.invalid</a></li><li>Tech blog: <a href="https://blog.zhwei.invalid">blog.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateAndroidPresetEn() []byte {
	androidByte, _ := json.Marshal(androidPresetJSONEn)
	return androidByte
}
