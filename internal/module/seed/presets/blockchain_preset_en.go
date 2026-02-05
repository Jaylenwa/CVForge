package presets

import (
	"encoding/json"
	"openresume/internal/common"
)

var blockchainPresetJSONEn = map[string]any{
	"title":    "Blockchain Developer Resume",
	"language": "en",
	"Personal": map[string]any{
		"FullName":   "Na Li",
		"Email":      "lina@example.com",
		"Phone":      "13900000000",
		"AvatarURL":  "/avatar.avif",
		"Job":        "Blockchain Developer",
		"City":       "Beijing",
		"Money":      "35k-55k",
		"JoinTime":   "Available immediately",
		"Gender":     "Female",
		"Age":        "27",
		"Degree":     "Master's",
		"CustomInfo": `[{"label":"Political Status","value":"Mass"}]`,
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
			"title":     "Summary",
			"isVisible": true,
			"items": []any{
				map[string]any{
					"id":          "sum1",
					"description": "<ul><li>4 years in blockchain development; smart contracts and dApps</li><li>Solidity/Hardhat/EVM; cross-chain and security audits</li><li>Web3 with wallet/signature/transaction; frontend/backend collaboration</li><li>Performance and gas optimization; contract versioning and risk governance</li></ul>",
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
					"title":       "Technology Company",
					"subtitle":    "Blockchain Developer (Contracts/Platform)",
					"timeStart":   "2021-06",
					"today":       true,
					"description": "<ul><li>Designed and implemented smart contracts and dApps</li><li>Testing and audit processes for security and compliance</li></ul>",
				},
				map[string]any{
					"id":          "e2",
					"title":       "Internet Company",
					"subtitle":    "Blockchain Developer (Web3/Interaction)",
					"timeStart":   "2019-08",
					"timeEnd":     "2021-05",
					"description": "<ul><li>Web3 and wallet interaction flows</li><li>Transaction/events governance; cost and UX optimization</li></ul>",
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
					"title":       "Smart Contracts & Platform Governance",
					"subtitle":    "Blockchain Developer",
					"timeStart":   "2022-03",
					"timeEnd":     "2022-11",
					"description": "<ul><li>Stack: Solidity / Hardhat / Ethers.js</li><li>Contracts and events; frontend interactions and governance</li><li>Testing and audits for security and compliance</li></ul>",
				},
				map[string]any{
					"id":          "p2",
					"title":       "Web3 & Wallet Interaction",
					"subtitle":    "Blockchain Developer",
					"timeStart":   "2020-06",
					"timeEnd":     "2021-02",
					"description": "<ul><li>Stack: Web3 / Metamask / Ethers.js</li><li>Signatures and transactions; cost and UX improvements</li></ul>",
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
					"title":       "Contracts & Platforms",
					"subtitle":    "Blockchain",
					"description": "Solidity / Hardhat / EVM",
				},
				map[string]any{
					"id":          "s2",
					"title":       "Interaction & Services",
					"subtitle":    "Web3",
					"description": "Web3 / Ethers.js / Wallet",
				},
				map[string]any{
					"id":          "s3",
					"title":       "Performance & Security",
					"subtitle":    "Governance",
					"description": "Audit / Gas Optimize / Versioning",
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
					"title":       "Harbin Institute of Technology",
					"major":       "Software Engineering",
					"degree":      "Master's",
					"timeStart":   "2016-09",
					"timeEnd":     "2019-06",
					"description": "Blockchain and distributed systems; platforms and contracts.",
				},
				map[string]any{
					"id":          "ed2",
					"title":       "Harbin Institute of Technology",
					"major":       "Computer Science and Technology",
					"degree":      "Bachelor's",
					"timeStart":   "2012-09",
					"timeEnd":     "2016-06",
					"description": "Distributed systems and computing; campus projects.",
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
					"description": `<ul><li>Contract examples: <a href="https://contracts.zhwei.invalid">contracts.zhwei.invalid</a></li><li>Web3 demos: <a href="https://web3-demo.zhwei.invalid">web3-demo.zhwei.invalid</a></li></ul>`,
				},
			},
		},
	},
}

func GenerateBlockchainPresetEn() []byte {
	bcByte, _ := json.Marshal(blockchainPresetJSONEn)
	return bcByte
}
