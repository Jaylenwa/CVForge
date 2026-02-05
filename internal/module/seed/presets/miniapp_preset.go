package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var miniappPresetJSON = map[string]any{
	"title":    "小程序开发工程师简历",
	"language": "zh",
	"Personal": map[string]any{
		"FullName":   "张伟",
		"Email":      "zhangwei@example.com",
		"Phone":      "13800000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "小程序开发工程师",
		"City":       "上海",
		"Money":      "25k-35k",
		"JoinTime":   "随时到岗",
		"Gender":     "男",
		"Age":        "28",
		"Degree":     "硕士",
		"CustomInfo": `[{"label":"政治面貌","value":"党员"}]`,
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
			"title":     "优势概述",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "sum1",
					"description": "<ul><li>5 年小程序开发经验，熟悉微信/支付宝/字节等多端生态，掌握原生与跨端框架（原生/uni-app/Taro）</li><li>熟悉组件化与状态管理，掌握页面与路由生命周期、网络与存储、权限与登录体系</li><li>具备性能优化与体验提升：首屏优化、图片与资源管理、分包与预加载、长列表优化</li><li>熟悉云开发与后端对接：云函数、消息订阅、模板消息/消息推送与服务端接口设计</li></ul>",
				},
			},
		},
		map[string]any{
			"id":        "exp",
			"type":      common.ResumeSectionTypeExperience,
			"title":     "工作经历",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "e1",
					"title":       "某零售平台",
					"subtitle":    "小程序开发工程师（交易/会员）",
					"timeStart":   "2021-05",
					"today":       true,
					"description": "<ul><li>负责交易与会员小程序（微信原生 + TypeScript），沉淀组件与工具（请求/埋点/权限）</li><li>优化首屏与交互性能：分包与预加载、骨架屏与图片懒加载，首页渲染加速显著</li><li>引入云开发能力（云函数/订阅消息），完善通知与任务处理链路</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "某生活服务平台",
					"subtitle":    "小程序开发工程师（服务/运营）",
					"timeStart":   "2019-08",
					"timeEnd":     "2021-04",
					"description": "<ul><li>参与服务预约与运营活动模块（uni-app/Taro），实现多端统一与组件复用</li><li>完善登录与权限体系、埋点与异常上报，提升稳定性与问题定位效率</li></ul>",
				},
			},
		},
		map[string]any{
			"id":        "projects",
			"type":      common.ResumeSectionTypeProjects,
			"title":     "项目经历",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "p1",
					"title":       "多端统一组件库与工具集",
					"subtitle":    "小程序开发工程师",
					"timeStart":   "2022-03",
					"timeEnd":     "2022-11",
					"description": "<ul><li>技术栈：原生 / uni-app / Taro / TypeScript</li><li>沉淀请求/埋点/权限/存储等工具与 UI 组件，统一多端实现</li><li>支持分包策略与按需加载，提升性能与可维护性</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "云开发与订阅消息方案",
					"subtitle":    "小程序开发工程师",
					"timeStart":   "2020-09",
					"timeEnd":     "2021-03",
					"description": "<ul><li>技术栈：微信云开发 / 云函数 / 云存储</li><li>实现订单与通知链路云端化，降低服务端耦合与运维成本</li><li>沉淀订阅消息与模板管理，提升触达与稳定性</li></ul>",
				},
			},
		},
		map[string]any{
			"id":        "skills",
			"type":      common.ResumeSectionTypeSkills,
			"title":     "技能清单",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "s1",
					"title":       "框架与技术",
					"subtitle":    "小程序",
					"description": "原生 / uni-app / Taro / TypeScript / WXML / WXSS",
				},
				map[string]any{
					"id":          "s2",
					"title":       "工程与性能",
					"subtitle":    "优化",
					"description": "分包 / 预加载 / 骨架屏 / 图片优化 / 埋点与异常上报",
				},
				map[string]any{
					"id":          "s3",
					"title":       "云开发与对接",
					"subtitle":    "后端",
					"description": "云函数 / 订阅消息 / 云存储 / REST API / 鉴权",
				},
			},
		},
		map[string]any{
			"id":        "edu",
			"type":      common.ResumeSectionTypeEducation,
			"title":     "教育背景",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "ed1",
					"title":       "浙江大学",
					"major":       "计算机科学与技术",
					"degree":      "硕士",
					"timeStart":   "2017-09",
					"timeEnd":     "2019-06",
					"description": "研究方向为移动端与小程序架构；参与云开发与性能优化项目。",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "浙江大学",
					"major":       "软件工程",
					"degree":      "本科",
					"timeStart":   "2013-09",
					"timeEnd":     "2017-06",
					"description": "主修计算机网络与前端技术；参与校内项目与竞赛。",
				},
			},
		},
		map[string]any{
			"id":        "portfolio",
			"type":      common.ResumeSectionTypePortfolio,
			"title":     "个人作品",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "pf1",
					"description": `<ul><li>小程序组件库：<a href="https://minikit.zhwei.invalid">minikit.zhwei.invalid</a></li><li>云开发示例：<a href="https://cloudmini.zhwei.invalid">cloudmini.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateMiniappPreset() []byte {
	miniByte, _ := json.Marshal(miniappPresetJSON)
	return miniByte
}
