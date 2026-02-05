package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var harmonyPresetJSONEn = map[string]any{
	"title":    "HarmonyOS Developer Resume",
	"language": "en",
	"Personal": map[string]any{
		"FullName":   "Wei Zhang",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "HarmonyOS Developer",
		"City":       "Shanghai",
		"Money":      "25k-40k",
		"JoinTime":   "Available immediately",
		"Gender":     "Male",
		"Age":        "28",
		"Degree":     "Master's",
		"CustomInfo": `[{"label":"Political Status","value":"Party Member"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#0ea5e9",
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
					"description": "<ul><li>3+ years of HarmonyOS development; familiar with ArkTS/ETS, Ability model and UI framework</li><li>Stage model, distributed capabilities and device collaboration, HAP packaging/signing</li><li>Performance & reliability: startup/rendering, memory & power analysis, crash governance</li><li>Engineering & integration: multi-module projects, capability/permission integration, backend API & messaging</li></ul>",
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
					"title":       "IoT Platform",
					"subtitle":    "HarmonyOS Developer (Device/Control)",
					"timeStart":   "2021-09",
					"today":       true,
					"description": "<ul><li>Developed device management and control app with ArkTS + Ability</li><li>Implemented distributed collaboration and capability sharing</li><li>Optimized performance and stability</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "Lifestyle Services",
					"subtitle":    "HarmonyOS Developer (Services/Operations)",
					"timeStart":   "2019-08",
					"timeEnd":     "2021-08",
					"description": "<ul><li>Built service booking and operations modules with distributed data</li><li>Improved permissions/capabilities integration, tracking and error reporting</li></ul>",
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
					"title":       "Distributed Device Collaboration App",
					"subtitle":    "HarmonyOS Developer",
					"timeStart":   "2022-05",
					"timeEnd":     "2022-12",
					"description": "<ul><li>Stack: ArkTS / Ability / Distributed Data</li><li>Implemented cross-device capability collaboration and state sync</li><li>Optimized performance and reliability</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "UI Framework & Component Library",
					"subtitle":    "HarmonyOS Developer",
					"timeStart":   "2020-08",
					"timeEnd":     "2021-04",
					"description": "<ul><li>Stack: ArkTS / ETS / UI</li><li>Shared components, theming, state management and routing</li><li>Improved consistency and delivery efficiency</li></ul>",
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
					"subtitle":    "HarmonyOS",
					"description": "ArkTS / ETS / Ability / UI",
				},
				map[string]any{
					"id":          "s2",
					"title":       "Engineering & Performance",
					"subtitle":    "Optimization",
					"description": "Multi-module / HAP / Startup & Rendering / Memory & Power / Crash",
				},
				map[string]any{
					"id":          "s3",
					"title":       "Distributed & Integration",
					"subtitle":    "Collaboration",
					"description": "Distributed Capability / Data Sync / Permissions / Backend API",
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
					"title":       "Huazhong University of Science and Technology",
					"major":       "Computer Science and Technology",
					"degree":      "Master's",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "Focus on distributed systems and mobile; HarmonyOS ecosystem projects.",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "Huazhong University of Science and Technology",
					"major":       "Software Engineering",
					"degree":      "Bachelor's",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "System development and engineering; campus projects and competitions.",
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
					"description": `<ul><li>ArkTS kit: <a href="https://arktskit.zhwei.invalid">arktskit.zhwei.invalid</a></li><li>Distributed demo: <a href="https://distributed.zhwei.invalid">distributed.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateHarmonyPresetEn() []byte {
	harmonyByte, _ := json.Marshal(harmonyPresetJSONEn)
	return harmonyByte
}
