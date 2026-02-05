package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var miniappPresetJSONEn = map[string]any{
	"title":    "Mini App Developer Resume",
	"language": "en",
	"Personal": map[string]any{
		"FullName":   "Wei Zhang",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "Mini App Developer",
		"City":       "Shanghai",
		"Money":      "25k-35k",
		"JoinTime":   "Available immediately",
		"Gender":     "Male",
		"Age":        "28",
		"Degree":     "Master's",
		"CustomInfo": `[{"label":"Political Status","value":"Party Member"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#10b981",
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
					"description": "<ul><li>5 years of mini app development; familiar with ecosystems of WeChat/Alipay/ByteDance, proficient in native and cross-platform frameworks (Native/uni-app/Taro)</li><li>Componentization and state management; page/router lifecycles, networking/storage, permissions and authentication</li><li>Performance and UX optimization: first screen, asset management, subpackages/preloading, long list optimization</li><li>Cloud development and backend integration: cloud functions, subscription messages, templated notifications, and REST API design</li></ul>",
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
					"title":       "Retail Platform",
					"subtitle":    "Mini App Developer (Checkout/Membership)",
					"timeStart":   "2021-05",
					"today":       true,
					"description": "<ul><li>Led checkout and membership mini apps (WeChat native + TypeScript), building common components and utilities (request/tracking/auth)</li><li>Optimized first screen and interactions: subpackages & preloading, skeleton screens, image lazy loading</li><li>Adopted cloud development (cloud functions/subscription messages) for notification and task processing</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "Lifestyle Platform",
					"subtitle":    "Mini App Developer (Services/Operations)",
					"timeStart":   "2019-08",
					"timeEnd":     "2021-04",
					"description": "<ul><li>Built service booking and operations modules (uni-app/Taro), unified multi-end implementation and component reuse</li><li>Improved auth/permission systems, tracking and error reporting to enhance reliability</li></ul>",
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
					"title":       "Multi-End Component Library & Toolkit",
					"subtitle":    "Mini App Developer",
					"timeStart":   "2022-03",
					"timeEnd":     "2022-11",
					"description": "<ul><li>Stack: Native / uni-app / Taro / TypeScript</li><li>Common utilities for request/tracking/auth/storage and UI components with unified implementations</li><li>Supported subpackage strategies and on-demand loading for performance and maintainability</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "Cloud Development & Subscription Messaging",
					"subtitle":    "Mini App Developer",
					"timeStart":   "2020-09",
					"timeEnd":     "2021-03",
					"description": "<ul><li>Stack: WeChat Cloud Development / Cloud Functions / Cloud Storage</li><li>Migrated orders and notification flows to cloud functions to reduce backend coupling</li><li>Standardized subscription messaging and templates for reliable delivery</li></ul>",
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
					"title":       "Frameworks & Tech",
					"subtitle":    "Miniapp",
					"description": "Native / uni-app / Taro / TypeScript / WXML / WXSS",
				},
				map[string]any{
					"id":          "s2",
					"title":       "Engineering & Performance",
					"subtitle":    "Optimization",
					"description": "Subpackages / Preloading / Skeleton Screens / Image Optimization / Tracking & Error Reporting",
				},
				map[string]any{
					"id":          "s3",
					"title":       "Cloud Dev & Integration",
					"subtitle":    "Backend",
					"description": "Cloud Functions / Subscription Messages / Cloud Storage / REST API / Auth",
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
					"major":       "Computer Science and Technology",
					"degree":      "Master's",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "Focus on mobile and mini app architectures; worked on cloud development and performance optimization.",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "Zhejiang University",
					"major":       "Software Engineering",
					"degree":      "Bachelor's",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "Studied networks and frontend technologies; participated in campus projects and competitions.",
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
					"description": `<ul><li>Miniapp Component Kit: <a href="https://minikit.zhwei.invalid">minikit.zhwei.invalid</a></li><li>Cloud Dev Demo: <a href="https://cloudmini.zhwei.invalid">cloudmini.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateMiniappPresetEn() []byte {
	miniByte, _ := json.Marshal(miniappPresetJSONEn)
	return miniByte
}
