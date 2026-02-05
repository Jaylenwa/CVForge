package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var blockchainPresetJSON = map[string]any{
	"title":    "区块链开发工程师简历",
	"language": "zh",
	"Personal": map[string]any{
		"FullName":   "李娜",
		"Email":      "lina@example.com",
		"Phone":      "13900000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "区块链开发工程师",
		"City":       "北京",
		"Money":      "35k-55k",
		"JoinTime":   "随时到岗",
		"Gender":     "女",
		"Age":        "27",
		"Degree":     "硕士",
		"CustomInfo": `[{"label":"政治面貌","value":"群众"}]`,
	},
	"Theme": map[string]any{
		"Color":    "#8b5cf6",
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
					"description": "<ul><li>4 年区块链开发经验，熟悉智能合约与链上应用</li><li>掌握 Solidity/Hardhat/EVM，了解跨链与安全审计</li><li>熟悉 Web3 与钱包/签名/交易流程，前后端协作</li><li>具备性能与成本优化经验，治理合约版本与风险</li></ul>",
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
					"title":       "某技术公司",
					"subtitle":    "区块链开发工程师（合约/平台）",
					"timeStart":   "2021-06",
					"today":       true,
					"description": "<ul><li>设计与实现智能合约与链上应用</li><li>完善测试与审计流程，保障安全与合规</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "某互联网公司",
					"subtitle":    "区块链开发工程师（Web3/交互）",
					"timeStart":   "2019-08",
					"timeEnd":     "2021-05",
					"description": "<ul><li>建设 Web3 与钱包交互流程</li><li>治理交易与事件，优化成本与体验</li></ul>",
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
					"title":       "智能合约与平台治理",
					"subtitle":    "区块链开发工程师",
					"timeStart":   "2022-03",
					"timeEnd":     "2022-11",
					"description": "<ul><li>技术栈：Solidity / Hardhat / Ethers.js</li><li>设计合约与事件，与前端交互与治理</li><li>测试与审计，保障安全与合规</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "Web3 与钱包交互",
					"subtitle":    "区块链开发工程师",
					"timeStart":   "2020-06",
					"timeEnd":     "2021-02",
					"description": "<ul><li>技术栈：Web3 / Metamask / Ethers.js</li><li>完善签名与交易流程，优化成本与体验</li></ul>",
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
					"title":       "合约与平台",
					"subtitle":    "区块链",
					"description": "Solidity / Hardhat / EVM",
				},
				map[string]any{
					"id":          "s2",
					"title":       "交互与服务",
					"subtitle":    "Web3",
					"description": "Web3 / Ethers.js / Wallet",
				},
				map[string]any{
					"id":          "s3",
					"title":       "性能与安全",
					"subtitle":    "治理",
					"description": "Audit / Gas Optimize / Versioning",
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
					"title":       "哈尔滨工业大学",
					"major":       "软件工程",
					"degree":      "硕士",
					"timeStart":   "2016-09",
					"timeEnd":     "2019-06",
					"description": "研究方向为区块链与分布式系统；参与平台与合约。",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "哈尔滨工业大学",
					"major":       "计算机科学与技术",
					"degree":      "本科",
					"timeStart":   "2012-09",
					"timeEnd":     "2016-06",
					"description": "主修分布式与系统；参与校内项目与竞赛。",
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
					"description": `<ul><li>合约示例：<a href="https://contracts.zhwei.invalid">contracts.zhwei.invalid</a></li><li>Web3 演示：<a href="https://web3-demo.zhwei.invalid">web3-demo.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateBlockchainPreset() []byte {
	bcByte, _ := json.Marshal(blockchainPresetJSON)
	return bcByte
}
