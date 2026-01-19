package seed

import (
	"encoding/json"
	"net/http"

	"openresume/internal/infra/database"
	"openresume/internal/module/library"
	"openresume/internal/module/preset"
	"openresume/internal/module/taxonomy"
	"openresume/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SeedData struct {
	Categories []taxonomy.JobCategory
	Roles      []taxonomy.JobRole
	Presets    []preset.ContentPreset
	Variants   []library.TemplateVariant
}

type ImportCounts struct {
	JobCategories    int `json:"jobCategories"`
	JobRoles         int `json:"jobRoles"`
	ContentPresets   int `json:"contentPresets"`
	TemplateVariants int `json:"templateVariants"`
}

type AdminHandler struct{}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

func (h *AdminHandler) AdminImportDefault(c *gin.Context) {
	seed, err := DefaultSeed()
	if err != nil {
		logger.WithCtx(c).Error("seed.build failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "seed error"})
		return
	}

	var counts ImportCounts
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		taxRepo := taxonomy.NewRepo(tx)
		presetRepo := preset.NewRepo(tx)
		libraryRepo := library.NewRepo(tx)

		for i := range seed.Categories {
			if err := taxRepo.UpsertJobCategory(tx, &seed.Categories[i]); err != nil {
				return err
			}
			counts.JobCategories++
		}
		for i := range seed.Roles {
			if err := taxRepo.UpsertJobRole(tx, &seed.Roles[i]); err != nil {
				return err
			}
			counts.JobRoles++
		}
		for i := range seed.Presets {
			if err := presetRepo.UpsertContentPreset(tx, &seed.Presets[i]); err != nil {
				return err
			}
			counts.ContentPresets++
		}
		for i := range seed.Variants {
			if err := libraryRepo.UpsertTemplateVariant(tx, &seed.Variants[i]); err != nil {
				return err
			}
			counts.TemplateVariants++
		}
		return nil
	})
	if err != nil {
		logger.WithCtx(c).Error("seed.import failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "counts": counts})
}

func DefaultSeed() (SeedData, error) {
	presetJSON := map[string]any{
		"title":    "Java 开发工程师简历",
		"language": "zh",
		"Personal": map[string]any{
			"FullName":   "李雷",
			"Email":      "lilei@example.com",
			"Phone":      "13800000001",
			"Job":        "Java 开发工程师",
			"City":       "北京 / 上海",
			"Money":      "25k-35k·14薪",
			"JoinTime":   "1 个月内",
			"Degree":     "本科",
			"CustomInfo": `[{"label":"技术栈","value":"Java / Spring / MySQL / Redis"}]`,
		},
		"Theme": map[string]any{
			"Color":    "#14b8a6",
			"Font":     "notosans",
			"Spacing":  "normal",
			"FontSize": "13",
		},
		"sections": []any{
			map[string]any{
				"id":        "exp",
				"type":      "Experience",
				"title":     "工作经历",
				"isVisible": true,
				"items": []any{
					map[string]any{
						"id":          "e1",
						"title":       "某大型互联网公司",
						"subtitle":    "Java 开发工程师（支付/交易）",
						"timeStart":   "2021-06",
						"today":       true,
						"description": "<ul><li>负责交易链路核心接口开发与重构，推动服务化拆分，接口 P95 延迟降低</li><li>基于 Spring Boot 构建高并发 API，完善限流/熔断/降级</li><li>引入 Redis 缓存与异步消息削峰填谷，提升稳定性</li></ul>",
					},
				},
			},
			map[string]any{
				"id":        "skills",
				"type":      "Skills",
				"title":     "技能清单",
				"isVisible": true,
				"items": []any{
					map[string]any{
						"id":          "s1",
						"title":       "后端",
						"subtitle":    "Java",
						"description": "Spring Boot / Spring Cloud / MyBatis",
					},
				},
			},
		},
	}
	b, err := json.Marshal(presetJSON)
	if err != nil {
		return SeedData{}, err
	}

	return SeedData{
		Categories: []taxonomy.JobCategory{
			{ExternalID: "it", Name: "IT | 互联网", ParentExternalID: "", OrderNum: 10, IsActive: true},
			{ExternalID: "finance", Name: "金融 | 银行", ParentExternalID: "", OrderNum: 20, IsActive: true},
			{ExternalID: "education", Name: "教育 | 培训", ParentExternalID: "", OrderNum: 30, IsActive: true},
			{ExternalID: "healthcare", Name: "医疗 | 健康", ParentExternalID: "", OrderNum: 40, IsActive: true},
			{ExternalID: "realestate", Name: "建筑 | 房地产", ParentExternalID: "", OrderNum: 50, IsActive: true},
			{ExternalID: "manufacturing", Name: "制造 | 生产", ParentExternalID: "", OrderNum: 60, IsActive: true},
			{ExternalID: "logistics", Name: "交通 | 物流", ParentExternalID: "", OrderNum: 70, IsActive: true},
			{ExternalID: "services", Name: "服务 | 消费", ParentExternalID: "", OrderNum: 80, IsActive: true},
			{ExternalID: "media", Name: "文化 | 传媒", ParentExternalID: "", OrderNum: 90, IsActive: true},
			{ExternalID: "trade", Name: "贸易 | 进出口", ParentExternalID: "", OrderNum: 100, IsActive: true},
			{ExternalID: "energy", Name: "能源 | 环保", ParentExternalID: "", OrderNum: 110, IsActive: true},
			{ExternalID: "agriculture", Name: "农林 | 牧渔", ParentExternalID: "", OrderNum: 120, IsActive: true},
			{ExternalID: "public", Name: "公共 | 非盈利", ParentExternalID: "", OrderNum: 130, IsActive: true},
			{ExternalID: "others", Name: "其他 | 岗位", ParentExternalID: "", OrderNum: 140, IsActive: true},

			{ExternalID: "it_backend", Name: "后端开发/程序员", ParentExternalID: "it", OrderNum: 10, IsActive: true},
			{ExternalID: "it_frontend", Name: "前端开发", ParentExternalID: "it", OrderNum: 20, IsActive: true},
			{ExternalID: "it_mobile", Name: "移动开发", ParentExternalID: "it", OrderNum: 30, IsActive: true},
			{ExternalID: "it_testing", Name: "软件测试", ParentExternalID: "it", OrderNum: 40, IsActive: true},
			{ExternalID: "it_ops_sec_dba", Name: "运维 / 安全 / DBA", ParentExternalID: "it", OrderNum: 50, IsActive: true},
			{ExternalID: "it_ai_bigdata", Name: "新兴技术", ParentExternalID: "it", OrderNum: 60, IsActive: true},
			{ExternalID: "it_other_tech", Name: "其他技术岗", ParentExternalID: "it", OrderNum: 70, IsActive: true},
			{ExternalID: "it_senior", Name: "高端技术职位", ParentExternalID: "it", OrderNum: 80, IsActive: true},

			{ExternalID: "finance_counter_service", Name: "银行柜台/服务", ParentExternalID: "finance", OrderNum: 10, IsActive: true},
			{ExternalID: "finance_personal_wealth", Name: "个人金融与理财", ParentExternalID: "finance", OrderNum: 20, IsActive: true},
			{ExternalID: "finance_credit_approval", Name: "信贷/审批", ParentExternalID: "finance", OrderNum: 30, IsActive: true},
			{ExternalID: "finance_risk_compliance", Name: "风险管理/合规", ParentExternalID: "finance", OrderNum: 40, IsActive: true},
			{ExternalID: "finance_securities_invest", Name: "证券与投资", ParentExternalID: "finance", OrderNum: 50, IsActive: true},
			{ExternalID: "finance_insurance_actuary", Name: "保险/精算", ParentExternalID: "finance", OrderNum: 60, IsActive: true},
			{ExternalID: "finance_banking_support", Name: "银行业务支持", ParentExternalID: "finance", OrderNum: 70, IsActive: true},
			{ExternalID: "finance_trust_futures", Name: "信托/期货类", ParentExternalID: "finance", OrderNum: 80, IsActive: true},
			{ExternalID: "finance_bank_management", Name: "银行管理类岗位", ParentExternalID: "finance", OrderNum: 90, IsActive: true},
			{ExternalID: "finance_intern", Name: "银行新人/实习生", ParentExternalID: "finance", OrderNum: 100, IsActive: true},

			{ExternalID: "education_teacher", Name: "教师", ParentExternalID: "education", OrderNum: 10, IsActive: true},
			{ExternalID: "education_teaching_admin", Name: "教学管理", ParentExternalID: "education", OrderNum: 20, IsActive: true},
			{ExternalID: "education_student_services", Name: "学生服务", ParentExternalID: "education", OrderNum: 30, IsActive: true},
			{ExternalID: "education_training_lecturer", Name: "培训/讲师", ParentExternalID: "education", OrderNum: 40, IsActive: true},
			{ExternalID: "education_training_management", Name: "培训管理", ParentExternalID: "education", OrderNum: 50, IsActive: true},

			{ExternalID: "healthcare_doctor", Name: "医生", ParentExternalID: "healthcare", OrderNum: 10, IsActive: true},
			{ExternalID: "healthcare_nurse", Name: "护士", ParentExternalID: "healthcare", OrderNum: 20, IsActive: true},
			{ExternalID: "healthcare_medtech", Name: "医学技术岗", ParentExternalID: "healthcare", OrderNum: 30, IsActive: true},
			{ExternalID: "healthcare_pharma", Name: "药学岗", ParentExternalID: "healthcare", OrderNum: 40, IsActive: true},
			{ExternalID: "healthcare_devices", Name: "医疗器械岗", ParentExternalID: "healthcare", OrderNum: 50, IsActive: true},
			{ExternalID: "healthcare_other", Name: "其他医疗岗", ParentExternalID: "healthcare", OrderNum: 60, IsActive: true},
			{ExternalID: "healthcare_intern", Name: "医药实习", ParentExternalID: "healthcare", OrderNum: 70, IsActive: true},
			{ExternalID: "healthcare_management", Name: "医疗管理", ParentExternalID: "healthcare", OrderNum: 80, IsActive: true},

			{ExternalID: "realestate_design_planning", Name: "建筑设计/规划", ParentExternalID: "realestate", OrderNum: 10, IsActive: true},
			{ExternalID: "realestate_interior_landscape", Name: "室内/景观设计", ParentExternalID: "realestate", OrderNum: 20, IsActive: true},
			{ExternalID: "realestate_cost_budget", Name: "工程造价/预算", ParentExternalID: "realestate", OrderNum: 30, IsActive: true},
			{ExternalID: "realestate_construction_mgmt", Name: "工程施工/管理", ParentExternalID: "realestate", OrderNum: 40, IsActive: true},
			{ExternalID: "realestate_project_mgmt", Name: "项目管理", ParentExternalID: "realestate", OrderNum: 50, IsActive: true},
			{ExternalID: "realestate_sales_planning", Name: "房地产销售/策划", ParentExternalID: "realestate", OrderNum: 60, IsActive: true},
			{ExternalID: "realestate_property_mgmt", Name: "物业管理", ParentExternalID: "realestate", OrderNum: 70, IsActive: true},

			{ExternalID: "mfg_mechanical", Name: "机械制造", ParentExternalID: "manufacturing", OrderNum: 10, IsActive: true},
			{ExternalID: "mfg_electrical", Name: "电子/电气制造", ParentExternalID: "manufacturing", OrderNum: 20, IsActive: true},
			{ExternalID: "mfg_auto_transport", Name: "汽车/交通制造", ParentExternalID: "manufacturing", OrderNum: 30, IsActive: true},
			{ExternalID: "mfg_process_mold", Name: "工艺/模具工程", ParentExternalID: "manufacturing", OrderNum: 40, IsActive: true},
			{ExternalID: "mfg_prod_equip", Name: "生产/设备管理", ParentExternalID: "manufacturing", OrderNum: 50, IsActive: true},
			{ExternalID: "mfg_quality", Name: "质量管理", ParentExternalID: "manufacturing", OrderNum: 60, IsActive: true},
			{ExternalID: "mfg_rd_design", Name: "研发设计", ParentExternalID: "manufacturing", OrderNum: 70, IsActive: true},
			{ExternalID: "mfg_worker", Name: "生产技工", ParentExternalID: "manufacturing", OrderNum: 80, IsActive: true},

			{ExternalID: "logi_transport_service", Name: "运输服务", ParentExternalID: "logistics", OrderNum: 10, IsActive: true},
			{ExternalID: "logi_delivery", Name: "运输配送", ParentExternalID: "logistics", OrderNum: 20, IsActive: true},
			{ExternalID: "logi_warehouse", Name: "仓储管理", ParentExternalID: "logistics", OrderNum: 30, IsActive: true},
			{ExternalID: "logi_ops", Name: "物流运营", ParentExternalID: "logistics", OrderNum: 40, IsActive: true},
			{ExternalID: "logi_supply", Name: "供应链", ParentExternalID: "logistics", OrderNum: 50, IsActive: true},

			{ExternalID: "svc_food", Name: "餐饮服务", ParentExternalID: "services", OrderNum: 10, IsActive: true},
			{ExternalID: "svc_hotel_travel", Name: "酒店旅游", ParentExternalID: "services", OrderNum: 20, IsActive: true},
			{ExternalID: "svc_personal", Name: "个人服务", ParentExternalID: "services", OrderNum: 30, IsActive: true},
			{ExternalID: "svc_fitness", Name: "运动健身", ParentExternalID: "services", OrderNum: 40, IsActive: true},
			{ExternalID: "svc_beauty", Name: "美容保健", ParentExternalID: "services", OrderNum: 50, IsActive: true},
			{ExternalID: "svc_consult_translate", Name: "咨询翻译", ParentExternalID: "services", OrderNum: 60, IsActive: true},
			{ExternalID: "svc_repair", Name: "维修岗位", ParentExternalID: "services", OrderNum: 70, IsActive: true},
			{ExternalID: "svc_other", Name: "其他服务岗", ParentExternalID: "services", OrderNum: 80, IsActive: true},

			{ExternalID: "media_news_publishing", Name: "新闻/出版", ParentExternalID: "media", OrderNum: 10, IsActive: true},
			{ExternalID: "media_broadcast_tv", Name: "广播电视", ParentExternalID: "media", OrderNum: 20, IsActive: true},
			{ExternalID: "media_film_performance", Name: "影视演艺", ParentExternalID: "media", OrderNum: 30, IsActive: true},

			{ExternalID: "trade_foreign_trade", Name: "外贸业务", ParentExternalID: "trade", OrderNum: 10, IsActive: true},
			{ExternalID: "trade_trade_support", Name: "国际贸易支持", ParentExternalID: "trade", OrderNum: 20, IsActive: true},
			{ExternalID: "trade_crossborder_ecom", Name: "跨境电商", ParentExternalID: "trade", OrderNum: 30, IsActive: true},
			{ExternalID: "trade_translation_support", Name: "翻译支持", ParentExternalID: "trade", OrderNum: 40, IsActive: true},

			{ExternalID: "energy_traditional", Name: "传统能源", ParentExternalID: "energy", OrderNum: 10, IsActive: true},
			{ExternalID: "energy_environment", Name: "环境保护", ParentExternalID: "energy", OrderNum: 20, IsActive: true},

			{ExternalID: "agri_planting", Name: "农业种植", ParentExternalID: "agriculture", OrderNum: 10, IsActive: true},
			{ExternalID: "agri_forestry", Name: "林业", ParentExternalID: "agriculture", OrderNum: 20, IsActive: true},
			{ExternalID: "agri_livestock", Name: "畜牧业", ParentExternalID: "agriculture", OrderNum: 30, IsActive: true},
			{ExternalID: "agri_fishery", Name: "渔业水产", ParentExternalID: "agriculture", OrderNum: 40, IsActive: true},

			{ExternalID: "public_services", Name: "公共服务", ParentExternalID: "public", OrderNum: 10, IsActive: true},
			{ExternalID: "public_research", Name: "科研相关", ParentExternalID: "public", OrderNum: 20, IsActive: true},
			{ExternalID: "public_social", Name: "社会服务", ParentExternalID: "public", OrderNum: 30, IsActive: true},
		},
		Roles: []taxonomy.JobRole{
			{ExternalID: "java", CategoryExternalID: "it_backend", Name: "java", Tags: "Java,后端", OrderNum: 10, IsActive: true},
			{ExternalID: "python", CategoryExternalID: "it_backend", Name: "python", Tags: "Python,后端", OrderNum: 20, IsActive: true},
			{ExternalID: "golang", CategoryExternalID: "it_backend", Name: "Go (Golang)", Tags: "Go,后端", OrderNum: 30, IsActive: true},
			{ExternalID: "php", CategoryExternalID: "it_backend", Name: "php", Tags: "PHP,后端", OrderNum: 40, IsActive: true},
			{ExternalID: "c_cpp", CategoryExternalID: "it_backend", Name: "C/C++", Tags: "C,C++,后端", OrderNum: 50, IsActive: true},
			{ExternalID: "csharp", CategoryExternalID: "it_backend", Name: "C#", Tags: "C#,后端", OrderNum: 60, IsActive: true},
			{ExternalID: "dotnet", CategoryExternalID: "it_backend", Name: ".NET", Tags: ".NET,后端", OrderNum: 70, IsActive: true},
			{ExternalID: "nodejs", CategoryExternalID: "it_backend", Name: "node.js", Tags: "Node.js,后端", OrderNum: 80, IsActive: true},

			{ExternalID: "web_frontend", CategoryExternalID: "it_frontend", Name: "Web前端", Tags: "前端", OrderNum: 10, IsActive: true},
			{ExternalID: "html5", CategoryExternalID: "it_frontend", Name: "HTML5", Tags: "前端", OrderNum: 20, IsActive: true},
			{ExternalID: "miniapp", CategoryExternalID: "it_frontend", Name: "小程序开发工程师", Tags: "前端,小程序", OrderNum: 30, IsActive: true},

			{ExternalID: "android", CategoryExternalID: "it_mobile", Name: "Android开发", Tags: "移动端", OrderNum: 10, IsActive: true},
			{ExternalID: "ios", CategoryExternalID: "it_mobile", Name: "iOS开发", Tags: "移动端", OrderNum: 20, IsActive: true},
			{ExternalID: "harmony", CategoryExternalID: "it_mobile", Name: "鸿蒙开发工程师", Tags: "移动端,鸿蒙", OrderNum: 30, IsActive: true},

			{ExternalID: "test_engineer", CategoryExternalID: "it_testing", Name: "测试工程师", Tags: "测试", OrderNum: 10, IsActive: true},
			{ExternalID: "automation_test", CategoryExternalID: "it_testing", Name: "自动化测试", Tags: "测试,自动化", OrderNum: 20, IsActive: true},
			{ExternalID: "test_dev", CategoryExternalID: "it_testing", Name: "测试开发", Tags: "测试,开发", OrderNum: 30, IsActive: true},
			{ExternalID: "performance_test", CategoryExternalID: "it_testing", Name: "性能测试", Tags: "测试,性能", OrderNum: 40, IsActive: true},
			{ExternalID: "hardware_test", CategoryExternalID: "it_testing", Name: "硬件测试工程师", Tags: "测试,硬件", OrderNum: 50, IsActive: true},

			{ExternalID: "ops_engineer", CategoryExternalID: "it_ops_sec_dba", Name: "运维工程师", Tags: "运维", OrderNum: 10, IsActive: true},
			{ExternalID: "devops", CategoryExternalID: "it_ops_sec_dba", Name: "DevOps工程师", Tags: "运维,DevOps", OrderNum: 20, IsActive: true},
			{ExternalID: "sys_admin", CategoryExternalID: "it_ops_sec_dba", Name: "系统/网络管理员", Tags: "运维,系统", OrderNum: 30, IsActive: true},
			{ExternalID: "dba", CategoryExternalID: "it_ops_sec_dba", Name: "数据库管理员 (DBA)", Tags: "数据库,DBA", OrderNum: 40, IsActive: true},
			{ExternalID: "security_engineer", CategoryExternalID: "it_ops_sec_dba", Name: "安全工程师", Tags: "安全", OrderNum: 50, IsActive: true},
			{ExternalID: "cloud_engineer", CategoryExternalID: "it_ops_sec_dba", Name: "云计算工程师", Tags: "云计算", OrderNum: 60, IsActive: true},
			{ExternalID: "ops_manager", CategoryExternalID: "it_ops_sec_dba", Name: "运维经理/主管", Tags: "运维,管理", OrderNum: 70, IsActive: true},

			{ExternalID: "data_mining", CategoryExternalID: "it_ai_bigdata", Name: "数据挖掘", Tags: "数据挖掘", OrderNum: 10, IsActive: true},
			{ExternalID: "nlp", CategoryExternalID: "it_ai_bigdata", Name: "自然语言处理", Tags: "NLP,AI", OrderNum: 20, IsActive: true},
			{ExternalID: "ml_ai", CategoryExternalID: "it_ai_bigdata", Name: "机器学习/AI工程师", Tags: "机器学习,AI", OrderNum: 30, IsActive: true},
			{ExternalID: "bigdata", CategoryExternalID: "it_ai_bigdata", Name: "大数据工程师", Tags: "大数据", OrderNum: 40, IsActive: true},
			{ExternalID: "blockchain", CategoryExternalID: "it_ai_bigdata", Name: "区块链开发", Tags: "区块链", OrderNum: 50, IsActive: true},
			{ExternalID: "algo_engineer", CategoryExternalID: "it_ai_bigdata", Name: "算法工程师", Tags: "算法", OrderNum: 60, IsActive: true},

			{ExternalID: "tech_support", CategoryExternalID: "it_other_tech", Name: "技术支持工程师", Tags: "技术支持", OrderNum: 10, IsActive: true},
			{ExternalID: "tech_engineer", CategoryExternalID: "it_other_tech", Name: "技术工程师", Tags: "技术", OrderNum: 20, IsActive: true},
			{ExternalID: "presales_support", CategoryExternalID: "it_other_tech", Name: "售前售后/技术支持", Tags: "售前,售后,技术支持", OrderNum: 30, IsActive: true},
			{ExternalID: "other_tech_roles", CategoryExternalID: "it_other_tech", Name: "其他技术岗位", Tags: "技术", OrderNum: 40, IsActive: true},
			{ExternalID: "network_engineer", CategoryExternalID: "it_other_tech", Name: "网络工程师", Tags: "网络", OrderNum: 50, IsActive: true},
			{ExternalID: "hardware_dev", CategoryExternalID: "it_other_tech", Name: "硬件开发工程师", Tags: "硬件", OrderNum: 60, IsActive: true},
			{ExternalID: "system_integrator", CategoryExternalID: "it_other_tech", Name: "系统集成工程师", Tags: "系统集成", OrderNum: 70, IsActive: true},
			{ExternalID: "circuit_engineer", CategoryExternalID: "it_other_tech", Name: "电路工程师", Tags: "电路", OrderNum: 80, IsActive: true},

			{ExternalID: "architect", CategoryExternalID: "it_senior", Name: "架构师", Tags: "架构", OrderNum: 10, IsActive: true},
			{ExternalID: "tech_manager", CategoryExternalID: "it_senior", Name: "技术主管", Tags: "管理", OrderNum: 20, IsActive: true},
			{ExternalID: "tech_director", CategoryExternalID: "it_senior", Name: "技术经理/总监", Tags: "管理", OrderNum: 30, IsActive: true},
			{ExternalID: "rd_manager", CategoryExternalID: "it_senior", Name: "研发经理/总监", Tags: "研发,管理", OrderNum: 40, IsActive: true},
			{ExternalID: "cto", CategoryExternalID: "it_senior", Name: "CTO", Tags: "管理", OrderNum: 50, IsActive: true},
			{ExternalID: "fullstack", CategoryExternalID: "it_senior", Name: "全栈工程师", Tags: "全栈", OrderNum: 60, IsActive: true},

			{ExternalID: "finance_teller", CategoryExternalID: "finance_counter_service", Name: "银行柜员", Tags: "金融,银行,柜台", OrderNum: 10, IsActive: true},
			{ExternalID: "finance_general_teller", CategoryExternalID: "finance_counter_service", Name: "银行综合柜员", Tags: "金融,银行,柜台", OrderNum: 20, IsActive: true},
			{ExternalID: "finance_lobby_manager", CategoryExternalID: "finance_counter_service", Name: "银行大堂经理", Tags: "金融,银行,服务", OrderNum: 30, IsActive: true},
			{ExternalID: "finance_lobby_guide", CategoryExternalID: "finance_counter_service", Name: "大堂引导员/大堂助理", Tags: "金融,银行,服务", OrderNum: 40, IsActive: true},
			{ExternalID: "finance_bank_cs", CategoryExternalID: "finance_counter_service", Name: "银行客服/坐席员", Tags: "金融,银行,客服", OrderNum: 50, IsActive: true},
			{ExternalID: "finance_bank_frontdesk", CategoryExternalID: "finance_counter_service", Name: "银行前台", Tags: "金融,银行,前台", OrderNum: 60, IsActive: true},

			{ExternalID: "finance_account_manager", CategoryExternalID: "finance_personal_wealth", Name: "客户经理", Tags: "金融,理财", OrderNum: 10, IsActive: true},
			{ExternalID: "finance_wealth_manager", CategoryExternalID: "finance_personal_wealth", Name: "理财经理", Tags: "金融,理财", OrderNum: 20, IsActive: true},
			{ExternalID: "finance_wealth_advisor", CategoryExternalID: "finance_personal_wealth", Name: "理财顾问", Tags: "金融,理财", OrderNum: 30, IsActive: true},
			{ExternalID: "finance_investment_advisor", CategoryExternalID: "finance_personal_wealth", Name: "投资顾问", Tags: "金融,投资", OrderNum: 40, IsActive: true},
			{ExternalID: "finance_creditcard_sales", CategoryExternalID: "finance_personal_wealth", Name: "信用卡销售", Tags: "金融,销售", OrderNum: 50, IsActive: true},

			{ExternalID: "finance_credit_manager", CategoryExternalID: "finance_credit_approval", Name: "信贷经理", Tags: "金融,信贷", OrderNum: 10, IsActive: true},
			{ExternalID: "finance_credit_officer", CategoryExternalID: "finance_credit_approval", Name: "信贷专员", Tags: "金融,信贷", OrderNum: 20, IsActive: true},
			{ExternalID: "finance_loan_officer", CategoryExternalID: "finance_credit_approval", Name: "贷款专员", Tags: "金融,贷款", OrderNum: 30, IsActive: true},
			{ExternalID: "finance_postloan", CategoryExternalID: "finance_credit_approval", Name: "贷后管理岗", Tags: "金融,贷后", OrderNum: 40, IsActive: true},
			{ExternalID: "finance_collections", CategoryExternalID: "finance_credit_approval", Name: "催收岗", Tags: "金融,催收", OrderNum: 50, IsActive: true},
			{ExternalID: "finance_mortgage_officer", CategoryExternalID: "finance_credit_approval", Name: "按揭专员", Tags: "金融,按揭", OrderNum: 60, IsActive: true},
			{ExternalID: "finance_credit_admin", CategoryExternalID: "finance_credit_approval", Name: "信贷管理", Tags: "金融,信贷", OrderNum: 70, IsActive: true},

			{ExternalID: "finance_risk_manager", CategoryExternalID: "finance_risk_compliance", Name: "风险经理", Tags: "金融,风控", OrderNum: 10, IsActive: true},
			{ExternalID: "finance_compliance_manager", CategoryExternalID: "finance_risk_compliance", Name: "合规经理", Tags: "金融,合规", OrderNum: 20, IsActive: true},
			{ExternalID: "finance_risk_control", CategoryExternalID: "finance_risk_compliance", Name: "风控专员", Tags: "金融,风控", OrderNum: 30, IsActive: true},
			{ExternalID: "finance_audit", CategoryExternalID: "finance_risk_compliance", Name: "审计专员", Tags: "金融,审计", OrderNum: 40, IsActive: true},
			{ExternalID: "finance_fin_compliance", CategoryExternalID: "finance_risk_compliance", Name: "金融合规专员", Tags: "金融,合规", OrderNum: 50, IsActive: true},
			{ExternalID: "finance_aml", CategoryExternalID: "finance_risk_compliance", Name: "反洗钱专员", Tags: "金融,合规", OrderNum: 60, IsActive: true},
			{ExternalID: "finance_legal", CategoryExternalID: "finance_risk_compliance", Name: "法律事务岗", Tags: "金融,法务", OrderNum: 70, IsActive: true},

			{ExternalID: "finance_securities_broker", CategoryExternalID: "finance_securities_invest", Name: "证券经纪人", Tags: "金融,证券", OrderNum: 10, IsActive: true},
			{ExternalID: "finance_invest_manager", CategoryExternalID: "finance_securities_invest", Name: "投资经理", Tags: "金融,投资", OrderNum: 20, IsActive: true},
			{ExternalID: "finance_fund_manager", CategoryExternalID: "finance_securities_invest", Name: "基金经理", Tags: "金融,基金", OrderNum: 30, IsActive: true},
			{ExternalID: "finance_trader", CategoryExternalID: "finance_securities_invest", Name: "交易员", Tags: "金融,交易", OrderNum: 40, IsActive: true},
			{ExternalID: "finance_fund_accountant", CategoryExternalID: "finance_securities_invest", Name: "基金会计", Tags: "金融,基金", OrderNum: 50, IsActive: true},
			{ExternalID: "finance_settlement", CategoryExternalID: "finance_securities_invest", Name: "清算专员", Tags: "金融,清算", OrderNum: 60, IsActive: true},
			{ExternalID: "finance_researcher", CategoryExternalID: "finance_securities_invest", Name: "研究员", Tags: "金融,研究", OrderNum: 70, IsActive: true},
			{ExternalID: "finance_analyst", CategoryExternalID: "finance_securities_invest", Name: "金融分析师", Tags: "金融,分析", OrderNum: 80, IsActive: true},
			{ExternalID: "finance_securities_analyst", CategoryExternalID: "finance_securities_invest", Name: "证券分析师", Tags: "金融,证券,分析", OrderNum: 90, IsActive: true},

			{ExternalID: "finance_insurance_advisor", CategoryExternalID: "finance_insurance_actuary", Name: "保险顾问", Tags: "金融,保险", OrderNum: 10, IsActive: true},
			{ExternalID: "finance_insurance_broker", CategoryExternalID: "finance_insurance_actuary", Name: "保险经纪人", Tags: "金融,保险", OrderNum: 20, IsActive: true},
			{ExternalID: "finance_insurance_agent", CategoryExternalID: "finance_insurance_actuary", Name: "保险代理专员", Tags: "金融,保险", OrderNum: 30, IsActive: true},
			{ExternalID: "finance_underwriter", CategoryExternalID: "finance_insurance_actuary", Name: "核保师", Tags: "金融,保险,核保", OrderNum: 40, IsActive: true},
			{ExternalID: "finance_claims", CategoryExternalID: "finance_insurance_actuary", Name: "理赔师", Tags: "金融,保险,理赔", OrderNum: 50, IsActive: true},
			{ExternalID: "finance_actuary", CategoryExternalID: "finance_insurance_actuary", Name: "精算师", Tags: "金融,保险,精算", OrderNum: 60, IsActive: true},
			{ExternalID: "finance_insurance_trainer", CategoryExternalID: "finance_insurance_actuary", Name: "保险培训师", Tags: "金融,保险,培训", OrderNum: 70, IsActive: true},
			{ExternalID: "finance_insurance_ops", CategoryExternalID: "finance_insurance_actuary", Name: "保险内勤", Tags: "金融,保险", OrderNum: 80, IsActive: true},
			{ExternalID: "finance_insurance_coach", CategoryExternalID: "finance_insurance_actuary", Name: "保险组训", Tags: "金融,保险,培训", OrderNum: 90, IsActive: true},
			{ExternalID: "finance_surveyor", CategoryExternalID: "finance_insurance_actuary", Name: "查勘员", Tags: "金融,保险,查勘", OrderNum: 100, IsActive: true},
			{ExternalID: "finance_insurance_sales", CategoryExternalID: "finance_insurance_actuary", Name: "保险销售", Tags: "金融,保险,销售", OrderNum: 110, IsActive: true},

			{ExternalID: "finance_bank_ops_lead", CategoryExternalID: "finance_banking_support", Name: "银行运营主管", Tags: "金融,银行,运营", OrderNum: 10, IsActive: true},
			{ExternalID: "finance_data_entry", CategoryExternalID: "finance_banking_support", Name: "数据录入员", Tags: "金融,数据", OrderNum: 20, IsActive: true},
			{ExternalID: "finance_doc_specialist", CategoryExternalID: "finance_banking_support", Name: "单证处理专员", Tags: "金融,单证", OrderNum: 30, IsActive: true},
			{ExternalID: "finance_data_analyst", CategoryExternalID: "finance_banking_support", Name: "金融数据分析师", Tags: "金融,数据,分析", OrderNum: 40, IsActive: true},
			{ExternalID: "finance_funds_settlement", CategoryExternalID: "finance_banking_support", Name: "资金结算专员", Tags: "金融,结算", OrderNum: 50, IsActive: true},

			{ExternalID: "finance_trust_manager", CategoryExternalID: "finance_trust_futures", Name: "信托经理", Tags: "金融,信托", OrderNum: 10, IsActive: true},
			{ExternalID: "finance_trader_operator", CategoryExternalID: "finance_trust_futures", Name: "操盘手", Tags: "金融,交易", OrderNum: 20, IsActive: true},
			{ExternalID: "finance_futures_analyst", CategoryExternalID: "finance_trust_futures", Name: "期货分析师", Tags: "金融,期货,分析", OrderNum: 30, IsActive: true},
			{ExternalID: "finance_invest_strategy", CategoryExternalID: "finance_trust_futures", Name: "投资策略师", Tags: "金融,策略", OrderNum: 40, IsActive: true},

			{ExternalID: "finance_branch_manager", CategoryExternalID: "finance_bank_management", Name: "支行行长", Tags: "金融,银行,管理", OrderNum: 10, IsActive: true},
			{ExternalID: "finance_branch_vice_manager", CategoryExternalID: "finance_bank_management", Name: "支行副行长", Tags: "金融,银行,管理", OrderNum: 20, IsActive: true},
			{ExternalID: "finance_subbranch_manager", CategoryExternalID: "finance_bank_management", Name: "分行行长", Tags: "金融,银行,管理", OrderNum: 30, IsActive: true},
			{ExternalID: "finance_cfo", CategoryExternalID: "finance_bank_management", Name: "首席财务官 (CFO)", Tags: "金融,银行,管理", OrderNum: 40, IsActive: true},
			{ExternalID: "finance_invest_director", CategoryExternalID: "finance_bank_management", Name: "投资总监", Tags: "金融,投资,管理", OrderNum: 50, IsActive: true},

			{ExternalID: "finance_mt", CategoryExternalID: "finance_intern", Name: "银行管理培训生", Tags: "金融,银行,校招", OrderNum: 10, IsActive: true},
			{ExternalID: "finance_securities_intern", CategoryExternalID: "finance_intern", Name: "证券实习生", Tags: "金融,证券,实习", OrderNum: 20, IsActive: true},
			{ExternalID: "finance_insurance_intern", CategoryExternalID: "finance_intern", Name: "保险实习生", Tags: "金融,保险,实习", OrderNum: 30, IsActive: true},
			{ExternalID: "finance_ib_intern", CategoryExternalID: "finance_intern", Name: "投资银行实习生", Tags: "金融,投行,实习", OrderNum: 40, IsActive: true},
			{ExternalID: "finance_intern_assistant", CategoryExternalID: "finance_intern", Name: "实习助理", Tags: "金融,实习", OrderNum: 50, IsActive: true},
			{ExternalID: "finance_teller_intern", CategoryExternalID: "finance_intern", Name: "柜员实习生", Tags: "金融,银行,实习", OrderNum: 60, IsActive: true},
			{ExternalID: "finance_bank_intern", CategoryExternalID: "finance_intern", Name: "银行实习生", Tags: "金融,银行,实习", OrderNum: 70, IsActive: true},

			{ExternalID: "edu_chinese_teacher", CategoryExternalID: "education_teacher", Name: "语文教师", Tags: "教育,教师", OrderNum: 10, IsActive: true},
			{ExternalID: "edu_english_teacher", CategoryExternalID: "education_teacher", Name: "英语教师", Tags: "教育,教师", OrderNum: 20, IsActive: true},
			{ExternalID: "edu_art_teacher", CategoryExternalID: "education_teacher", Name: "美术老师", Tags: "教育,教师", OrderNum: 30, IsActive: true},
			{ExternalID: "edu_math_teacher", CategoryExternalID: "education_teacher", Name: "数学教师", Tags: "教育,教师", OrderNum: 40, IsActive: true},
			{ExternalID: "edu_kindergarten_teacher", CategoryExternalID: "education_teacher", Name: "幼儿教师", Tags: "教育,教师", OrderNum: 50, IsActive: true},
			{ExternalID: "edu_pe_teacher", CategoryExternalID: "education_teacher", Name: "体育教师", Tags: "教育,教师", OrderNum: 60, IsActive: true},
			{ExternalID: "edu_music_teacher", CategoryExternalID: "education_teacher", Name: "音乐老师", Tags: "教育,教师", OrderNum: 70, IsActive: true},
			{ExternalID: "edu_biology_teacher", CategoryExternalID: "education_teacher", Name: "生物教师", Tags: "教育,教师", OrderNum: 80, IsActive: true},
			{ExternalID: "edu_dance_teacher", CategoryExternalID: "education_teacher", Name: "舞蹈老师", Tags: "教育,教师", OrderNum: 90, IsActive: true},
			{ExternalID: "edu_piano_teacher", CategoryExternalID: "education_teacher", Name: "钢琴教师", Tags: "教育,教师", OrderNum: 100, IsActive: true},
			{ExternalID: "edu_calligraphy_teacher", CategoryExternalID: "education_teacher", Name: "书法教师", Tags: "教育,教师", OrderNum: 110, IsActive: true},
			{ExternalID: "edu_chemistry_teacher", CategoryExternalID: "education_teacher", Name: "化学老师", Tags: "教育,教师", OrderNum: 120, IsActive: true},
			{ExternalID: "edu_physics_teacher", CategoryExternalID: "education_teacher", Name: "物理老师", Tags: "教育,教师", OrderNum: 130, IsActive: true},
			{ExternalID: "edu_history_teacher", CategoryExternalID: "education_teacher", Name: "历史老师", Tags: "教育,教师", OrderNum: 140, IsActive: true},
			{ExternalID: "edu_politics_teacher", CategoryExternalID: "education_teacher", Name: "政治老师", Tags: "教育,教师", OrderNum: 150, IsActive: true},
			{ExternalID: "edu_geography_teacher", CategoryExternalID: "education_teacher", Name: "地理老师", Tags: "教育,教师", OrderNum: 160, IsActive: true},
			{ExternalID: "edu_tcsol_teacher", CategoryExternalID: "education_teacher", Name: "对外汉语教师", Tags: "教育,教师", OrderNum: 170, IsActive: true},
			{ExternalID: "edu_tutor", CategoryExternalID: "education_teacher", Name: "家教", Tags: "教育,教师", OrderNum: 180, IsActive: true},
			{ExternalID: "edu_university_teacher", CategoryExternalID: "education_teacher", Name: "大学教师", Tags: "教育,教师", OrderNum: 190, IsActive: true},
			{ExternalID: "edu_tutoring_teacher", CategoryExternalID: "education_teacher", Name: "辅导老师", Tags: "教育,教师", OrderNum: 200, IsActive: true},

			{ExternalID: "edu_principal", CategoryExternalID: "education_teaching_admin", Name: "校长", Tags: "教育,管理", OrderNum: 10, IsActive: true},
			{ExternalID: "edu_vice_principal", CategoryExternalID: "education_teaching_admin", Name: "副校长", Tags: "教育,管理", OrderNum: 20, IsActive: true},
			{ExternalID: "edu_kindergarten_principal", CategoryExternalID: "education_teaching_admin", Name: "园长", Tags: "教育,管理", OrderNum: 30, IsActive: true},
			{ExternalID: "edu_academic_director", CategoryExternalID: "education_teaching_admin", Name: "教务主任", Tags: "教育,教务", OrderNum: 40, IsActive: true},
			{ExternalID: "edu_research_leader", CategoryExternalID: "education_teaching_admin", Name: "教研组长", Tags: "教育,教研", OrderNum: 50, IsActive: true},
			{ExternalID: "edu_teaching_supervisor", CategoryExternalID: "education_teaching_admin", Name: "教学主管", Tags: "教育,管理", OrderNum: 60, IsActive: true},
			{ExternalID: "edu_academic_specialist", CategoryExternalID: "education_teaching_admin", Name: "教务专员", Tags: "教育,教务", OrderNum: 70, IsActive: true},
			{ExternalID: "edu_researcher", CategoryExternalID: "education_teaching_admin", Name: "教研员", Tags: "教育,教研", OrderNum: 80, IsActive: true},
			{ExternalID: "edu_curriculum_dev", CategoryExternalID: "education_teaching_admin", Name: "课程开发/设计", Tags: "教育,课程", OrderNum: 90, IsActive: true},
			{ExternalID: "edu_academic_assistant", CategoryExternalID: "education_teaching_admin", Name: "教务助理", Tags: "教育,教务", OrderNum: 100, IsActive: true},

			{ExternalID: "edu_dorm_teacher", CategoryExternalID: "education_student_services", Name: "宿管老师", Tags: "教育,学生", OrderNum: 10, IsActive: true},
			{ExternalID: "edu_admissions", CategoryExternalID: "education_student_services", Name: "招生顾问", Tags: "教育,招生", OrderNum: 20, IsActive: true},
			{ExternalID: "edu_counselor", CategoryExternalID: "education_student_services", Name: "辅导员", Tags: "教育,学生", OrderNum: 30, IsActive: true},
			{ExternalID: "edu_ta", CategoryExternalID: "education_student_services", Name: "助教/教学助理", Tags: "教育,助教", OrderNum: 40, IsActive: true},
			{ExternalID: "edu_life_teacher", CategoryExternalID: "education_student_services", Name: "生活老师", Tags: "教育,学生", OrderNum: 50, IsActive: true},
			{ExternalID: "edu_study_abroad", CategoryExternalID: "education_student_services", Name: "留学顾问", Tags: "教育,留学", OrderNum: 60, IsActive: true},
			{ExternalID: "edu_homeroom_teacher", CategoryExternalID: "education_student_services", Name: "班主任", Tags: "教育,学生", OrderNum: 70, IsActive: true},
			{ExternalID: "edu_student_manager", CategoryExternalID: "education_student_services", Name: "学管师", Tags: "教育,学生", OrderNum: 80, IsActive: true},

			{ExternalID: "edu_trainer", CategoryExternalID: "education_training_lecturer", Name: "培训师", Tags: "教育,培训", OrderNum: 10, IsActive: true},
			{ExternalID: "edu_lecturer", CategoryExternalID: "education_training_lecturer", Name: "培训讲师", Tags: "教育,培训", OrderNum: 20, IsActive: true},
			{ExternalID: "edu_english_trainer", CategoryExternalID: "education_training_lecturer", Name: "英语培训老师", Tags: "教育,培训", OrderNum: 30, IsActive: true},
			{ExternalID: "edu_corporate_trainer", CategoryExternalID: "education_training_lecturer", Name: "企业培训师", Tags: "教育,培训", OrderNum: 40, IsActive: true},
			{ExternalID: "edu_vocational_trainer", CategoryExternalID: "education_training_lecturer", Name: "职业培训师", Tags: "教育,培训", OrderNum: 50, IsActive: true},
			{ExternalID: "edu_speaker", CategoryExternalID: "education_training_lecturer", Name: "讲师", Tags: "教育,培训", OrderNum: 60, IsActive: true},
			{ExternalID: "edu_internal_trainer", CategoryExternalID: "education_training_lecturer", Name: "内训师", Tags: "教育,培训", OrderNum: 70, IsActive: true},

			{ExternalID: "edu_training_specialist", CategoryExternalID: "education_training_management", Name: "培训专员", Tags: "教育,培训,管理", OrderNum: 10, IsActive: true},
			{ExternalID: "edu_training_supervisor", CategoryExternalID: "education_training_management", Name: "培训主管", Tags: "教育,培训,管理", OrderNum: 20, IsActive: true},
			{ExternalID: "edu_training_assistant", CategoryExternalID: "education_training_management", Name: "培训助理", Tags: "教育,培训,管理", OrderNum: 30, IsActive: true},
			{ExternalID: "edu_training_director", CategoryExternalID: "education_training_management", Name: "培训总监", Tags: "教育,培训,管理", OrderNum: 40, IsActive: true},
			{ExternalID: "edu_training_manager", CategoryExternalID: "education_training_management", Name: "培训经理", Tags: "教育,培训,管理", OrderNum: 50, IsActive: true},

			{ExternalID: "hc_clinical_doctor", CategoryExternalID: "healthcare_doctor", Name: "临床医生", Tags: "医疗,医生", OrderNum: 10, IsActive: true},
			{ExternalID: "hc_internal_doctor", CategoryExternalID: "healthcare_doctor", Name: "内科医生", Tags: "医疗,医生", OrderNum: 20, IsActive: true},
			{ExternalID: "hc_surgery_doctor", CategoryExternalID: "healthcare_doctor", Name: "外科医生", Tags: "医疗,医生", OrderNum: 30, IsActive: true},
			{ExternalID: "hc_obgyn_doctor", CategoryExternalID: "healthcare_doctor", Name: "妇产科医生", Tags: "医疗,医生", OrderNum: 40, IsActive: true},
			{ExternalID: "hc_pediatrics_doctor", CategoryExternalID: "healthcare_doctor", Name: "儿科医生", Tags: "医疗,医生", OrderNum: 50, IsActive: true},
			{ExternalID: "hc_orthopedics_doctor", CategoryExternalID: "healthcare_doctor", Name: "骨科医生", Tags: "医疗,医生", OrderNum: 60, IsActive: true},
			{ExternalID: "hc_anesthesiologist", CategoryExternalID: "healthcare_doctor", Name: "麻醉医生", Tags: "医疗,医生", OrderNum: 70, IsActive: true},
			{ExternalID: "hc_dentist", CategoryExternalID: "healthcare_doctor", Name: "口腔医生", Tags: "医疗,医生,口腔", OrderNum: 80, IsActive: true},
			{ExternalID: "hc_tcm", CategoryExternalID: "healthcare_doctor", Name: "中医师", Tags: "医疗,医生,中医", OrderNum: 90, IsActive: true},
			{ExternalID: "hc_radiologist", CategoryExternalID: "healthcare_doctor", Name: "放射科医生", Tags: "医疗,医生", OrderNum: 100, IsActive: true},
			{ExternalID: "hc_gp", CategoryExternalID: "healthcare_doctor", Name: "全科医生", Tags: "医疗,医生", OrderNum: 110, IsActive: true},
			{ExternalID: "hc_specialist_doctor", CategoryExternalID: "healthcare_doctor", Name: "专科医生", Tags: "医疗,医生", OrderNum: 120, IsActive: true},
			{ExternalID: "hc_doctor_assistant", CategoryExternalID: "healthcare_doctor", Name: "医生助理", Tags: "医疗,医生", OrderNum: 130, IsActive: true},
			{ExternalID: "hc_resident", CategoryExternalID: "healthcare_doctor", Name: "住院医师", Tags: "医疗,医生", OrderNum: 140, IsActive: true},
			{ExternalID: "hc_dental_doctor", CategoryExternalID: "healthcare_doctor", Name: "牙科医生", Tags: "医疗,医生,口腔", OrderNum: 150, IsActive: true},

			{ExternalID: "hc_nurse", CategoryExternalID: "healthcare_nurse", Name: "护士", Tags: "医疗,护士", OrderNum: 10, IsActive: true},
			{ExternalID: "hc_clinical_nurse", CategoryExternalID: "healthcare_nurse", Name: "临床护士", Tags: "医疗,护士", OrderNum: 20, IsActive: true},
			{ExternalID: "hc_or_nurse", CategoryExternalID: "healthcare_nurse", Name: "手术室护士", Tags: "医疗,护士", OrderNum: 30, IsActive: true},
			{ExternalID: "hc_internal_nurse", CategoryExternalID: "healthcare_nurse", Name: "内科护士", Tags: "医疗,护士", OrderNum: 40, IsActive: true},
			{ExternalID: "hc_obgyn_nurse", CategoryExternalID: "healthcare_nurse", Name: "妇产科护士", Tags: "医疗,护士", OrderNum: 50, IsActive: true},
			{ExternalID: "hc_midwife", CategoryExternalID: "healthcare_nurse", Name: "助产士", Tags: "医疗,护士", OrderNum: 60, IsActive: true},

			{ExternalID: "hc_lab_tech", CategoryExternalID: "healthcare_medtech", Name: "医学检验技师", Tags: "医疗,技师", OrderNum: 10, IsActive: true},
			{ExternalID: "hc_imaging_tech", CategoryExternalID: "healthcare_medtech", Name: "医学影像技师", Tags: "医疗,技师", OrderNum: 20, IsActive: true},
			{ExternalID: "hc_rehab_therapist", CategoryExternalID: "healthcare_medtech", Name: "康复治疗师", Tags: "医疗,康复", OrderNum: 30, IsActive: true},
			{ExternalID: "hc_dietitian", CategoryExternalID: "healthcare_medtech", Name: "营养师", Tags: "医疗,营养", OrderNum: 40, IsActive: true},
			{ExternalID: "hc_health_manager", CategoryExternalID: "healthcare_medtech", Name: "健康管理师", Tags: "医疗,健康管理", OrderNum: 50, IsActive: true},
			{ExternalID: "hc_psych_consultant", CategoryExternalID: "healthcare_medtech", Name: "心理咨询师", Tags: "医疗,心理", OrderNum: 60, IsActive: true},
			{ExternalID: "hc_acupuncture", CategoryExternalID: "healthcare_medtech", Name: "针灸推拿", Tags: "医疗,中医", OrderNum: 70, IsActive: true},
			{ExternalID: "hc_lab_tester", CategoryExternalID: "healthcare_medtech", Name: "检验师", Tags: "医疗,检验", OrderNum: 80, IsActive: true},
			{ExternalID: "hc_ultrasound_doctor", CategoryExternalID: "healthcare_medtech", Name: "超声科医师", Tags: "医疗,医生", OrderNum: 90, IsActive: true},
			{ExternalID: "hc_pathologist", CategoryExternalID: "healthcare_medtech", Name: "病理科医师", Tags: "医疗,医生", OrderNum: 100, IsActive: true},

			{ExternalID: "hc_pharma_related", CategoryExternalID: "healthcare_pharma", Name: "药学相关", Tags: "医疗,药学", OrderNum: 10, IsActive: true},
			{ExternalID: "hc_drug_rd", CategoryExternalID: "healthcare_pharma", Name: "药物研发", Tags: "医疗,药学,研发", OrderNum: 20, IsActive: true},
			{ExternalID: "hc_medicine_rd", CategoryExternalID: "healthcare_pharma", Name: "药品研发", Tags: "医疗,药学,研发", OrderNum: 30, IsActive: true},
			{ExternalID: "hc_pharmacist", CategoryExternalID: "healthcare_pharma", Name: "药剂师/药师", Tags: "医疗,药学", OrderNum: 40, IsActive: true},
			{ExternalID: "hc_med_qc", CategoryExternalID: "healthcare_pharma", Name: "医药质检", Tags: "医疗,药学,质检", OrderNum: 50, IsActive: true},
			{ExternalID: "hc_drug_registration", CategoryExternalID: "healthcare_pharma", Name: "药品注册", Tags: "医疗,药学,注册", OrderNum: 60, IsActive: true},
			{ExternalID: "hc_cra", CategoryExternalID: "healthcare_pharma", Name: "临床监察员 (CRA)", Tags: "医疗,药学,临床", OrderNum: 70, IsActive: true},
			{ExternalID: "hc_crc", CategoryExternalID: "healthcare_pharma", Name: "临床协调员 (CRC)", Tags: "医疗,药学,临床", OrderNum: 80, IsActive: true},
			{ExternalID: "hc_drug_quality_mgmt", CategoryExternalID: "healthcare_pharma", Name: "药品质量管理", Tags: "医疗,药学,质量", OrderNum: 90, IsActive: true},

			{ExternalID: "hc_device_sales", CategoryExternalID: "healthcare_devices", Name: "医疗器械销售", Tags: "医疗,器械,销售", OrderNum: 10, IsActive: true},
			{ExternalID: "hc_after_sales", CategoryExternalID: "healthcare_devices", Name: "售后工程师", Tags: "医疗,器械,售后", OrderNum: 20, IsActive: true},
			{ExternalID: "hc_device_inspector", CategoryExternalID: "healthcare_devices", Name: "检验员", Tags: "医疗,器械,检验", OrderNum: 30, IsActive: true},
			{ExternalID: "hc_device_qc", CategoryExternalID: "healthcare_devices", Name: "质检员", Tags: "医疗,器械,质检", OrderNum: 40, IsActive: true},

			{ExternalID: "hc_records_admin", CategoryExternalID: "healthcare_other", Name: "病案管理员", Tags: "医疗,行政", OrderNum: 10, IsActive: true},
			{ExternalID: "hc_registration_cashier", CategoryExternalID: "healthcare_other", Name: "挂号/收费员", Tags: "医疗,行政", OrderNum: 20, IsActive: true},
			{ExternalID: "hc_guide", CategoryExternalID: "healthcare_other", Name: "导医", Tags: "医疗,服务", OrderNum: 30, IsActive: true},
			{ExternalID: "hc_hospital_frontdesk", CategoryExternalID: "healthcare_other", Name: "医院前台", Tags: "医疗,前台", OrderNum: 40, IsActive: true},
			{ExternalID: "hc_inventory", CategoryExternalID: "healthcare_other", Name: "医疗库管", Tags: "医疗,仓储", OrderNum: 50, IsActive: true},
			{ExternalID: "hc_security", CategoryExternalID: "healthcare_other", Name: "医院保安", Tags: "医疗,安保", OrderNum: 60, IsActive: true},
			{ExternalID: "hc_cleaner", CategoryExternalID: "healthcare_other", Name: "医院保洁", Tags: "医疗,后勤", OrderNum: 70, IsActive: true},
			{ExternalID: "hc_caregiver", CategoryExternalID: "healthcare_other", Name: "医院护工", Tags: "医疗,后勤", OrderNum: 80, IsActive: true},
			{ExternalID: "hc_companion", CategoryExternalID: "healthcare_other", Name: "医院陪护", Tags: "医疗,后勤", OrderNum: 90, IsActive: true},
			{ExternalID: "hc_med_translator", CategoryExternalID: "healthcare_other", Name: "医疗翻译", Tags: "医疗,翻译", OrderNum: 100, IsActive: true},
			{ExternalID: "hc_med_legal", CategoryExternalID: "healthcare_other", Name: "医疗法务", Tags: "医疗,法务", OrderNum: 110, IsActive: true},
			{ExternalID: "hc_med_consulting", CategoryExternalID: "healthcare_other", Name: "医疗咨询", Tags: "医疗,咨询", OrderNum: 120, IsActive: true},
			{ExternalID: "hc_med_training", CategoryExternalID: "healthcare_other", Name: "医疗培训", Tags: "医疗,培训", OrderNum: 130, IsActive: true},
			{ExternalID: "hc_med_social_worker", CategoryExternalID: "healthcare_other", Name: "医疗社工", Tags: "医疗,社工", OrderNum: 140, IsActive: true},
			{ExternalID: "hc_med_rep", CategoryExternalID: "healthcare_other", Name: "医药代表", Tags: "医疗,销售", OrderNum: 150, IsActive: true},

			{ExternalID: "hc_intern_doctor", CategoryExternalID: "healthcare_intern", Name: "实习医生", Tags: "医疗,实习", OrderNum: 10, IsActive: true},
			{ExternalID: "hc_intern_nurse", CategoryExternalID: "healthcare_intern", Name: "实习护士", Tags: "医疗,实习", OrderNum: 20, IsActive: true},
			{ExternalID: "hc_intern_pharmacist", CategoryExternalID: "healthcare_intern", Name: "实习药师", Tags: "医疗,实习", OrderNum: 30, IsActive: true},
			{ExternalID: "hc_intern_tech", CategoryExternalID: "healthcare_intern", Name: "实习技师", Tags: "医疗,实习", OrderNum: 40, IsActive: true},
			{ExternalID: "hc_med_intern", CategoryExternalID: "healthcare_intern", Name: "医学实习生", Tags: "医疗,实习", OrderNum: 50, IsActive: true},
			{ExternalID: "hc_intern_med_rep", CategoryExternalID: "healthcare_intern", Name: "医药代表实习生", Tags: "医疗,实习,销售", OrderNum: 60, IsActive: true},
			{ExternalID: "hc_clinical_research_intern", CategoryExternalID: "healthcare_intern", Name: "临床研究实习生", Tags: "医疗,实习,临床", OrderNum: 70, IsActive: true},

			{ExternalID: "hc_hospital_director", CategoryExternalID: "healthcare_management", Name: "医院院长", Tags: "医疗,管理", OrderNum: 10, IsActive: true},
			{ExternalID: "hc_department_head", CategoryExternalID: "healthcare_management", Name: "科室主任", Tags: "医疗,管理", OrderNum: 20, IsActive: true},
			{ExternalID: "hc_head_nurse", CategoryExternalID: "healthcare_management", Name: "护士长", Tags: "医疗,管理", OrderNum: 30, IsActive: true},

			{ExternalID: "re_structural_engineer", CategoryExternalID: "realestate_design_planning", Name: "结构工程师", Tags: "建筑,设计", OrderNum: 10, IsActive: true},
			{ExternalID: "re_urban_planner", CategoryExternalID: "realestate_design_planning", Name: "城市规划师", Tags: "建筑,规划", OrderNum: 20, IsActive: true},
			{ExternalID: "re_landscape_designer", CategoryExternalID: "realestate_design_planning", Name: "园林设计师", Tags: "建筑,园林", OrderNum: 30, IsActive: true},
			{ExternalID: "re_planning_design", CategoryExternalID: "realestate_design_planning", Name: "规划设计", Tags: "建筑,规划", OrderNum: 40, IsActive: true},
			{ExternalID: "re_bim_engineer", CategoryExternalID: "realestate_design_planning", Name: "BIM工程师", Tags: "建筑,BIM", OrderNum: 50, IsActive: true},
			{ExternalID: "re_landscape_constructor", CategoryExternalID: "realestate_design_planning", Name: "园林施工员", Tags: "建筑,园林,施工", OrderNum: 60, IsActive: true},
			{ExternalID: "re_construction_engineer", CategoryExternalID: "realestate_design_planning", Name: "建筑工程师", Tags: "建筑,工程", OrderNum: 70, IsActive: true},

			{ExternalID: "re_landscape_architect", CategoryExternalID: "realestate_interior_landscape", Name: "景观设计师", Tags: "建筑,设计", OrderNum: 10, IsActive: true},
			{ExternalID: "re_home_designer", CategoryExternalID: "realestate_interior_landscape", Name: "家装设计师", Tags: "建筑,设计", OrderNum: 20, IsActive: true},

			{ExternalID: "re_budgeter", CategoryExternalID: "realestate_cost_budget", Name: "预算员", Tags: "建筑,造价", OrderNum: 10, IsActive: true},
			{ExternalID: "re_cost_estimator", CategoryExternalID: "realestate_cost_budget", Name: "造价员", Tags: "建筑,造价", OrderNum: 20, IsActive: true},
			{ExternalID: "re_cost_engineer", CategoryExternalID: "realestate_cost_budget", Name: "造价工程师", Tags: "建筑,造价", OrderNum: 30, IsActive: true},
			{ExternalID: "re_project_cost", CategoryExternalID: "realestate_cost_budget", Name: "工程造价", Tags: "建筑,造价", OrderNum: 40, IsActive: true},

			{ExternalID: "re_site_worker", CategoryExternalID: "realestate_construction_mgmt", Name: "施工员", Tags: "建筑,施工", OrderNum: 10, IsActive: true},
			{ExternalID: "re_civil_site_worker", CategoryExternalID: "realestate_construction_mgmt", Name: "土建施工员", Tags: "建筑,施工,土建", OrderNum: 20, IsActive: true},
			{ExternalID: "re_surveyor", CategoryExternalID: "realestate_construction_mgmt", Name: "测量员", Tags: "建筑,测量", OrderNum: 30, IsActive: true},
			{ExternalID: "re_mapping_engineer", CategoryExternalID: "realestate_construction_mgmt", Name: "测绘工程师", Tags: "建筑,测绘", OrderNum: 40, IsActive: true},
			{ExternalID: "re_supervisor", CategoryExternalID: "realestate_construction_mgmt", Name: "工程监理", Tags: "建筑,监理", OrderNum: 50, IsActive: true},
			{ExternalID: "re_supervision_engineer", CategoryExternalID: "realestate_construction_mgmt", Name: "监理工程师", Tags: "建筑,监理", OrderNum: 60, IsActive: true},
			{ExternalID: "re_project_admin", CategoryExternalID: "realestate_construction_mgmt", Name: "工程管理员", Tags: "建筑,工程", OrderNum: 70, IsActive: true},
			{ExternalID: "re_document_controller", CategoryExternalID: "realestate_construction_mgmt", Name: "资料员", Tags: "建筑,资料", OrderNum: 80, IsActive: true},
			{ExternalID: "re_archive_admin", CategoryExternalID: "realestate_construction_mgmt", Name: "档案管理员", Tags: "建筑,档案", OrderNum: 90, IsActive: true},
			{ExternalID: "re_engineering_manager", CategoryExternalID: "realestate_construction_mgmt", Name: "工程经理", Tags: "建筑,工程,管理", OrderNum: 100, IsActive: true},
			{ExternalID: "re_civil_engineer", CategoryExternalID: "realestate_construction_mgmt", Name: "土木工程师", Tags: "建筑,土木", OrderNum: 110, IsActive: true},
			{ExternalID: "re_safety_officer", CategoryExternalID: "realestate_construction_mgmt", Name: "安全员", Tags: "建筑,安全", OrderNum: 120, IsActive: true},
			{ExternalID: "re_quality_inspector", CategoryExternalID: "realestate_construction_mgmt", Name: "工程质检员", Tags: "建筑,质检", OrderNum: 130, IsActive: true},
			{ExternalID: "re_project_engineer", CategoryExternalID: "realestate_construction_mgmt", Name: "项目工程师", Tags: "建筑,项目", OrderNum: 140, IsActive: true},
			{ExternalID: "re_civil_project_engineer", CategoryExternalID: "realestate_construction_mgmt", Name: "土建工程师", Tags: "建筑,土建", OrderNum: 150, IsActive: true},

			{ExternalID: "re_pm", CategoryExternalID: "realestate_project_mgmt", Name: "项目经理", Tags: "建筑,项目,管理", OrderNum: 10, IsActive: true},
			{ExternalID: "re_pm_assistant", CategoryExternalID: "realestate_project_mgmt", Name: "项目助理", Tags: "建筑,项目", OrderNum: 20, IsActive: true},
			{ExternalID: "re_pm_specialist", CategoryExternalID: "realestate_project_mgmt", Name: "项目专员", Tags: "建筑,项目", OrderNum: 30, IsActive: true},
			{ExternalID: "re_pm_supervisor", CategoryExternalID: "realestate_project_mgmt", Name: "项目主管", Tags: "建筑,项目,管理", OrderNum: 40, IsActive: true},
			{ExternalID: "re_pm_director", CategoryExternalID: "realestate_project_mgmt", Name: "项目总监", Tags: "建筑,项目,管理", OrderNum: 50, IsActive: true},
			{ExternalID: "re_bid_specialist", CategoryExternalID: "realestate_project_mgmt", Name: "投标专员", Tags: "建筑,投标", OrderNum: 60, IsActive: true},

			{ExternalID: "re_sales", CategoryExternalID: "realestate_sales_planning", Name: "房地产销售", Tags: "房产,销售", OrderNum: 10, IsActive: true},
			{ExternalID: "re_property_consultant", CategoryExternalID: "realestate_sales_planning", Name: "置业顾问", Tags: "房产,销售", OrderNum: 20, IsActive: true},
			{ExternalID: "re_agent", CategoryExternalID: "realestate_sales_planning", Name: "房产经纪人", Tags: "房产,销售", OrderNum: 30, IsActive: true},
			{ExternalID: "re_marketing_planner", CategoryExternalID: "realestate_sales_planning", Name: "房地产策划", Tags: "房产,策划", OrderNum: 40, IsActive: true},
			{ExternalID: "re_leasing_manager", CategoryExternalID: "realestate_sales_planning", Name: "招商经理", Tags: "房产,招商", OrderNum: 50, IsActive: true},
			{ExternalID: "re_channel_manager", CategoryExternalID: "realestate_sales_planning", Name: "渠道经理", Tags: "房产,渠道", OrderNum: 60, IsActive: true},
			{ExternalID: "re_other_roles", CategoryExternalID: "realestate_sales_planning", Name: "房产其他岗位", Tags: "房产", OrderNum: 70, IsActive: true},

			{ExternalID: "re_property_mgmt", CategoryExternalID: "realestate_property_mgmt", Name: "物业管理", Tags: "物业,管理", OrderNum: 10, IsActive: true},
			{ExternalID: "re_property_manager", CategoryExternalID: "realestate_property_mgmt", Name: "物业经理", Tags: "物业,管理", OrderNum: 20, IsActive: true},
			{ExternalID: "re_property_cs", CategoryExternalID: "realestate_property_mgmt", Name: "物业客服", Tags: "物业,客服", OrderNum: 30, IsActive: true},
			{ExternalID: "re_property_steward", CategoryExternalID: "realestate_property_mgmt", Name: "物业管家", Tags: "物业,服务", OrderNum: 40, IsActive: true},

			{ExternalID: "mfg_mech_engineer", CategoryExternalID: "mfg_mechanical", Name: "机械工程师", Tags: "制造,机械", OrderNum: 10, IsActive: true},
			{ExternalID: "mfg_mech_design", CategoryExternalID: "mfg_mechanical", Name: "机械设计工程师", Tags: "制造,机械,设计", OrderNum: 20, IsActive: true},
			{ExternalID: "mfg_mechatronics", CategoryExternalID: "mfg_mechanical", Name: "机电工程师", Tags: "制造,机电", OrderNum: 30, IsActive: true},
			{ExternalID: "mfg_equipment_engineer", CategoryExternalID: "mfg_mechanical", Name: "设备工程师", Tags: "制造,设备", OrderNum: 40, IsActive: true},
			{ExternalID: "mfg_mechanical_manufacturing", CategoryExternalID: "mfg_mechanical", Name: "机械制造", Tags: "制造,机械", OrderNum: 50, IsActive: true},
			{ExternalID: "mfg_mechanical_maintenance", CategoryExternalID: "mfg_mechanical", Name: "机械维修", Tags: "制造,维修", OrderNum: 60, IsActive: true},
			{ExternalID: "mfg_hydraulics", CategoryExternalID: "mfg_mechanical", Name: "液压工程师", Tags: "制造,液压", OrderNum: 70, IsActive: true},
			{ExternalID: "mfg_nc_programmer", CategoryExternalID: "mfg_mechanical", Name: "数控编程", Tags: "制造,数控", OrderNum: 80, IsActive: true},
			{ExternalID: "mfg_other_engineers", CategoryExternalID: "mfg_mechanical", Name: "其他工程师岗位", Tags: "制造,工程师", OrderNum: 90, IsActive: true},

			{ExternalID: "mfg_elec_engineer", CategoryExternalID: "mfg_electrical", Name: "电子工程师", Tags: "制造,电子", OrderNum: 10, IsActive: true},
			{ExternalID: "mfg_electrical_engineer", CategoryExternalID: "mfg_electrical", Name: "电气工程师", Tags: "制造,电气", OrderNum: 20, IsActive: true},
			{ExternalID: "mfg_hw_engineer", CategoryExternalID: "mfg_electrical", Name: "硬件工程师", Tags: "制造,硬件", OrderNum: 30, IsActive: true},
			{ExternalID: "mfg_industrial_automation", CategoryExternalID: "mfg_electrical", Name: "电气自动化工程师", Tags: "制造,自动化", OrderNum: 40, IsActive: true},
			{ExternalID: "mfg_embedded", CategoryExternalID: "mfg_electrical", Name: "嵌入式工程师", Tags: "制造,嵌入式", OrderNum: 50, IsActive: true},
			{ExternalID: "mfg_automation", CategoryExternalID: "mfg_electrical", Name: "自动化工程师", Tags: "制造,自动化", OrderNum: 60, IsActive: true},
			{ExternalID: "mfg_semiconductor_tech", CategoryExternalID: "mfg_electrical", Name: "半导体技术员", Tags: "制造,半导体", OrderNum: 70, IsActive: true},
			{ExternalID: "mfg_circuit_design", CategoryExternalID: "mfg_electrical", Name: "电路设计", Tags: "制造,电路", OrderNum: 80, IsActive: true},

			{ExternalID: "mfg_auto_engineer", CategoryExternalID: "mfg_auto_transport", Name: "汽车工程师", Tags: "制造,汽车", OrderNum: 10, IsActive: true},
			{ExternalID: "mfg_vehicle_engineer", CategoryExternalID: "mfg_auto_transport", Name: "车辆工程师", Tags: "制造,汽车", OrderNum: 20, IsActive: true},
			{ExternalID: "mfg_vehicle_design", CategoryExternalID: "mfg_auto_transport", Name: "汽车设计", Tags: "制造,汽车,设计", OrderNum: 30, IsActive: true},
			{ExternalID: "mfg_powertrain", CategoryExternalID: "mfg_auto_transport", Name: "动力总成工程师", Tags: "制造,汽车", OrderNum: 40, IsActive: true},

			{ExternalID: "mfg_process_engineer", CategoryExternalID: "mfg_process_mold", Name: "工艺工程师", Tags: "制造,工艺", OrderNum: 10, IsActive: true},
			{ExternalID: "mfg_mould_engineer", CategoryExternalID: "mfg_process_mold", Name: "模具工程师", Tags: "制造,模具", OrderNum: 20, IsActive: true},
			{ExternalID: "mfg_welding_engineer", CategoryExternalID: "mfg_process_mold", Name: "焊接工程师", Tags: "制造,焊接", OrderNum: 30, IsActive: true},
			{ExternalID: "mfg_mould_design", CategoryExternalID: "mfg_process_mold", Name: "模具设计师", Tags: "制造,模具,设计", OrderNum: 40, IsActive: true},
			{ExternalID: "mfg_stamping_process", CategoryExternalID: "mfg_process_mold", Name: "冲压工艺师/模具设计师", Tags: "制造,冲压,模具", OrderNum: 50, IsActive: true},

			{ExternalID: "mfg_prod_management", CategoryExternalID: "mfg_prod_equip", Name: "生产管理", Tags: "制造,生产", OrderNum: 10, IsActive: true},
			{ExternalID: "mfg_prod_supervisor", CategoryExternalID: "mfg_prod_equip", Name: "生产主管", Tags: "制造,生产", OrderNum: 20, IsActive: true},
			{ExternalID: "mfg_production_manager", CategoryExternalID: "mfg_prod_equip", Name: "生产经理", Tags: "制造,生产,管理", OrderNum: 30, IsActive: true},
			{ExternalID: "mfg_line_leader", CategoryExternalID: "mfg_prod_equip", Name: "车间主任", Tags: "制造,车间,管理", OrderNum: 40, IsActive: true},
			{ExternalID: "mfg_equipment_maintenance", CategoryExternalID: "mfg_prod_equip", Name: "设备维护", Tags: "制造,设备,维护", OrderNum: 50, IsActive: true},
			{ExternalID: "mfg_equipment_manager", CategoryExternalID: "mfg_prod_equip", Name: "设备管理", Tags: "制造,设备,管理", OrderNum: 60, IsActive: true},
			{ExternalID: "mfg_shift_leader", CategoryExternalID: "mfg_prod_equip", Name: "生产班长", Tags: "制造,生产,班长", OrderNum: 70, IsActive: true},
			{ExternalID: "mfg_team_leader", CategoryExternalID: "mfg_prod_equip", Name: "工段长", Tags: "制造,工段", OrderNum: 80, IsActive: true},
			{ExternalID: "mfg_group_leader", CategoryExternalID: "mfg_prod_equip", Name: "班组长", Tags: "制造,班组", OrderNum: 90, IsActive: true},
			{ExternalID: "mfg_factory_manager", CategoryExternalID: "mfg_prod_equip", Name: "厂长", Tags: "制造,管理", OrderNum: 100, IsActive: true},
			{ExternalID: "mfg_production_planner", CategoryExternalID: "mfg_prod_equip", Name: "生产计划员", Tags: "制造,计划", OrderNum: 110, IsActive: true},

			{ExternalID: "mfg_quality_engineer", CategoryExternalID: "mfg_quality", Name: "质量工程师", Tags: "制造,质量", OrderNum: 10, IsActive: true},
			{ExternalID: "mfg_quality_specialist", CategoryExternalID: "mfg_quality", Name: "品质工程师", Tags: "制造,品质", OrderNum: 20, IsActive: true},
			{ExternalID: "mfg_qc", CategoryExternalID: "mfg_quality", Name: "QC", Tags: "制造,质量", OrderNum: 30, IsActive: true},
			{ExternalID: "mfg_quality_manager", CategoryExternalID: "mfg_quality", Name: "质量管理工程师", Tags: "制造,质量,管理", OrderNum: 40, IsActive: true},
			{ExternalID: "mfg_quality_inspector", CategoryExternalID: "mfg_quality", Name: "质量检测员", Tags: "制造,质量,检测", OrderNum: 50, IsActive: true},

			{ExternalID: "mfg_rd_engineer", CategoryExternalID: "mfg_rd_design", Name: "研发工程师", Tags: "制造,研发", OrderNum: 10, IsActive: true},
			{ExternalID: "mfg_tech_engineer", CategoryExternalID: "mfg_rd_design", Name: "技术员", Tags: "制造,技术", OrderNum: 20, IsActive: true},
			{ExternalID: "mfg_design_engineer", CategoryExternalID: "mfg_rd_design", Name: "设计工程师", Tags: "制造,设计", OrderNum: 30, IsActive: true},
			{ExternalID: "mfg_product_designer", CategoryExternalID: "mfg_rd_design", Name: "产品设计师", Tags: "制造,产品,设计", OrderNum: 40, IsActive: true},

			{ExternalID: "mfg_operator", CategoryExternalID: "mfg_worker", Name: "操作工", Tags: "制造,技工", OrderNum: 10, IsActive: true},
			{ExternalID: "mfg_general_worker", CategoryExternalID: "mfg_worker", Name: "普工", Tags: "制造,技工", OrderNum: 20, IsActive: true},
			{ExternalID: "mfg_assembler", CategoryExternalID: "mfg_worker", Name: "装配工", Tags: "制造,装配", OrderNum: 30, IsActive: true},
			{ExternalID: "mfg_technician", CategoryExternalID: "mfg_worker", Name: "技术工人", Tags: "制造,技工", OrderNum: 40, IsActive: true},
			{ExternalID: "mfg_fitter", CategoryExternalID: "mfg_worker", Name: "钳工", Tags: "制造,加工", OrderNum: 50, IsActive: true},
			{ExternalID: "mfg_welder", CategoryExternalID: "mfg_worker", Name: "焊工", Tags: "制造,焊接", OrderNum: 60, IsActive: true},
			{ExternalID: "mfg_driller", CategoryExternalID: "mfg_worker", Name: "钻工", Tags: "制造,加工", OrderNum: 70, IsActive: true},
			{ExternalID: "mfg_turner", CategoryExternalID: "mfg_worker", Name: "车工", Tags: "制造,加工", OrderNum: 80, IsActive: true},
			{ExternalID: "mfg_miller", CategoryExternalID: "mfg_worker", Name: "铣工", Tags: "制造,加工", OrderNum: 90, IsActive: true},
			{ExternalID: "mfg_cnc_operator", CategoryExternalID: "mfg_worker", Name: "CNC操作工", Tags: "制造,CNC", OrderNum: 100, IsActive: true},
			{ExternalID: "mfg_machine_operator", CategoryExternalID: "mfg_worker", Name: "机床操作工", Tags: "制造,机床", OrderNum: 110, IsActive: true},
			{ExternalID: "mfg_packer", CategoryExternalID: "mfg_worker", Name: "包装工", Tags: "制造,包装", OrderNum: 120, IsActive: true},

			{ExternalID: "logi_captain", CategoryExternalID: "logi_transport_service", Name: "机长", Tags: "物流,运输,航空", OrderNum: 10, IsActive: true},
			{ExternalID: "logi_pilot", CategoryExternalID: "logi_transport_service", Name: "飞行员", Tags: "物流,运输,航空", OrderNum: 20, IsActive: true},
			{ExternalID: "logi_air_security", CategoryExternalID: "logi_transport_service", Name: "空中安全员", Tags: "物流,运输,航空", OrderNum: 30, IsActive: true},
			{ExternalID: "logi_flight_attendant", CategoryExternalID: "logi_transport_service", Name: "空姐", Tags: "物流,运输,航空", OrderNum: 40, IsActive: true},
			{ExternalID: "logi_male_attendant", CategoryExternalID: "logi_transport_service", Name: "空少", Tags: "物流,运输,航空", OrderNum: 50, IsActive: true},
			{ExternalID: "logi_ground", CategoryExternalID: "logi_transport_service", Name: "地勤", Tags: "物流,运输,航空", OrderNum: 60, IsActive: true},
			{ExternalID: "logi_tickets", CategoryExternalID: "logi_transport_service", Name: "票务员", Tags: "物流,票务", OrderNum: 70, IsActive: true},
			{ExternalID: "logi_security", CategoryExternalID: "logi_transport_service", Name: "安检员", Tags: "物流,安检", OrderNum: 80, IsActive: true},
			{ExternalID: "logi_metro_staff", CategoryExternalID: "logi_transport_service", Name: "地铁站务员", Tags: "物流,地铁", OrderNum: 90, IsActive: true},
			{ExternalID: "logi_metro_security", CategoryExternalID: "logi_transport_service", Name: "地铁安检员", Tags: "物流,地铁,安检", OrderNum: 100, IsActive: true},
			{ExternalID: "logi_metro_driver", CategoryExternalID: "logi_transport_service", Name: "地铁驾驶员", Tags: "物流,地铁,驾驶", OrderNum: 110, IsActive: true},
			{ExternalID: "logi_highspeed_attendant", CategoryExternalID: "logi_transport_service", Name: "高铁乘务", Tags: "物流,铁路", OrderNum: 120, IsActive: true},
			{ExternalID: "logi_call_center", CategoryExternalID: "logi_transport_service", Name: "话务员", Tags: "物流,运输", OrderNum: 130, IsActive: true},
			{ExternalID: "logi_flight_dispatcher", CategoryExternalID: "logi_transport_service", Name: "签派员", Tags: "物流,航空", OrderNum: 140, IsActive: true},

			{ExternalID: "logi_driver", CategoryExternalID: "logi_delivery", Name: "司机", Tags: "物流,司机", OrderNum: 10, IsActive: true},
			{ExternalID: "logi_dispatch_manager", CategoryExternalID: "logi_delivery", Name: "配送经理", Tags: "物流,配送,管理", OrderNum: 20, IsActive: true},
			{ExternalID: "logi_transport_admin", CategoryExternalID: "logi_delivery", Name: "运输主管", Tags: "物流,运输,管理", OrderNum: 30, IsActive: true},
			{ExternalID: "logi_fleet_manager", CategoryExternalID: "logi_delivery", Name: "车队管理", Tags: "物流,车队", OrderNum: 40, IsActive: true},
			{ExternalID: "logi_scheduler", CategoryExternalID: "logi_delivery", Name: "调度员", Tags: "物流,调度", OrderNum: 50, IsActive: true},

			{ExternalID: "logi_warehouse_admin", CategoryExternalID: "logi_warehouse", Name: "仓库管理员/库管", Tags: "物流,仓储", OrderNum: 10, IsActive: true},
			{ExternalID: "logi_storekeeper", CategoryExternalID: "logi_warehouse", Name: "仓管员", Tags: "物流,仓储", OrderNum: 20, IsActive: true},
			{ExternalID: "logi_warehouse_manager", CategoryExternalID: "logi_warehouse", Name: "仓储管理", Tags: "物流,仓储", OrderNum: 30, IsActive: true},
			{ExternalID: "logi_loader", CategoryExternalID: "logi_warehouse", Name: "装卸工", Tags: "物流,仓储", OrderNum: 40, IsActive: true},
			{ExternalID: "logi_packer", CategoryExternalID: "logi_warehouse", Name: "包装员", Tags: "物流,仓储", OrderNum: 50, IsActive: true},
			{ExternalID: "logi_warehouse_specialist", CategoryExternalID: "logi_warehouse", Name: "仓储专员", Tags: "物流,仓储", OrderNum: 60, IsActive: true},
			{ExternalID: "logi_warehouse_supervisor", CategoryExternalID: "logi_warehouse", Name: "仓储主管", Tags: "物流,仓储", OrderNum: 70, IsActive: true},
			{ExternalID: "logi_warehouse_manager_role", CategoryExternalID: "logi_warehouse", Name: "仓储经理", Tags: "物流,仓储,管理", OrderNum: 80, IsActive: true},
			{ExternalID: "logi_warehouse_reserve", CategoryExternalID: "logi_warehouse", Name: "储备干部", Tags: "物流,仓储", OrderNum: 90, IsActive: true},
			{ExternalID: "logi_warehouse_superintendent", CategoryExternalID: "logi_warehouse", Name: "储备主管", Tags: "物流,仓储", OrderNum: 100, IsActive: true},

			{ExternalID: "logi_ops_specialist", CategoryExternalID: "logi_ops", Name: "物流专员/助理", Tags: "物流,运营", OrderNum: 10, IsActive: true},
			{ExternalID: "logi_ops_supervisor", CategoryExternalID: "logi_ops", Name: "物流主管", Tags: "物流,运营", OrderNum: 20, IsActive: true},
			{ExternalID: "logi_ops_manager", CategoryExternalID: "logi_ops", Name: "物流经理", Tags: "物流,运营,管理", OrderNum: 30, IsActive: true},
			{ExternalID: "logi_order_follower", CategoryExternalID: "logi_ops", Name: "跟单员", Tags: "物流,跟单", OrderNum: 40, IsActive: true},
			{ExternalID: "logi_courier", CategoryExternalID: "logi_ops", Name: "物流员", Tags: "物流", OrderNum: 50, IsActive: true},
			{ExternalID: "logi_express", CategoryExternalID: "logi_ops", Name: "快递员", Tags: "物流,快递", OrderNum: 60, IsActive: true},
			{ExternalID: "logi_delivery_staff", CategoryExternalID: "logi_ops", Name: "配送员", Tags: "物流,配送", OrderNum: 70, IsActive: true},

			{ExternalID: "logi_supply_specialist", CategoryExternalID: "logi_supply", Name: "供应链专员", Tags: "供应链", OrderNum: 10, IsActive: true},
			{ExternalID: "logi_supply_manager", CategoryExternalID: "logi_supply", Name: "供应链管理", Tags: "供应链,管理", OrderNum: 20, IsActive: true},
			{ExternalID: "logi_supply_director", CategoryExternalID: "logi_supply", Name: "供应链经理", Tags: "供应链,管理", OrderNum: 30, IsActive: true},
			{ExternalID: "logi_supply_supervisor", CategoryExternalID: "logi_supply", Name: "供应链总监", Tags: "供应链,管理", OrderNum: 40, IsActive: true},
			{ExternalID: "logi_supply_intern", CategoryExternalID: "logi_supply", Name: "供应链实习生", Tags: "供应链,实习", OrderNum: 50, IsActive: true},

			{ExternalID: "svc_supermarket_manager", CategoryExternalID: "svc_food", Name: "超市店长", Tags: "服务,餐饮", OrderNum: 10, IsActive: true},
			{ExternalID: "svc_cashier_staff", CategoryExternalID: "svc_food", Name: "收银员", Tags: "服务,餐饮", OrderNum: 20, IsActive: true},
			{ExternalID: "svc_waiter", CategoryExternalID: "svc_food", Name: "服务员", Tags: "服务,餐饮", OrderNum: 30, IsActive: true},
			{ExternalID: "svc_barista", CategoryExternalID: "svc_food", Name: "咖啡师", Tags: "服务,餐饮", OrderNum: 40, IsActive: true},
			{ExternalID: "svc_chef", CategoryExternalID: "svc_food", Name: "厨师", Tags: "服务,餐饮", OrderNum: 50, IsActive: true},
			{ExternalID: "svc_baker", CategoryExternalID: "svc_food", Name: "面包师", Tags: "服务,餐饮", OrderNum: 60, IsActive: true},
			{ExternalID: "svc_western_chef", CategoryExternalID: "svc_food", Name: "西点师", Tags: "服务,餐饮", OrderNum: 70, IsActive: true},
			{ExternalID: "svc_food_manager", CategoryExternalID: "svc_food", Name: "餐饮管理", Tags: "服务,餐饮,管理", OrderNum: 80, IsActive: true},
			{ExternalID: "svc_food_store_manager", CategoryExternalID: "svc_food", Name: "餐饮店长", Tags: "服务,餐饮,管理", OrderNum: 90, IsActive: true},
			{ExternalID: "svc_back_kitchen", CategoryExternalID: "svc_food", Name: "后厨", Tags: "服务,餐饮", OrderNum: 100, IsActive: true},
			{ExternalID: "svc_restaurant_manager", CategoryExternalID: "svc_food", Name: "餐厅经理", Tags: "服务,餐饮,管理", OrderNum: 110, IsActive: true},
			{ExternalID: "svc_pastry_chef", CategoryExternalID: "svc_food", Name: "面点师", Tags: "服务,餐饮", OrderNum: 120, IsActive: true},
			{ExternalID: "svc_head_chef", CategoryExternalID: "svc_food", Name: "厨师长", Tags: "服务,餐饮,管理", OrderNum: 130, IsActive: true},
			{ExternalID: "svc_baking", CategoryExternalID: "svc_food", Name: "烘焙师", Tags: "服务,餐饮", OrderNum: 140, IsActive: true},
			{ExternalID: "svc_bartender", CategoryExternalID: "svc_food", Name: "调酒师", Tags: "服务,餐饮", OrderNum: 150, IsActive: true},
			{ExternalID: "svc_tea_artist", CategoryExternalID: "svc_food", Name: "茶艺师", Tags: "服务,餐饮", OrderNum: 160, IsActive: true},
			{ExternalID: "svc_food_director", CategoryExternalID: "svc_food", Name: "餐饮总监", Tags: "服务,餐饮,管理", OrderNum: 170, IsActive: true},

			{ExternalID: "svc_tour_guide", CategoryExternalID: "svc_hotel_travel", Name: "导游", Tags: "服务,旅游", OrderNum: 10, IsActive: true},
			{ExternalID: "svc_travel_consultant", CategoryExternalID: "svc_hotel_travel", Name: "旅游顾问", Tags: "服务,旅游", OrderNum: 20, IsActive: true},
			{ExternalID: "svc_planner", CategoryExternalID: "svc_hotel_travel", Name: "计调", Tags: "服务,旅游", OrderNum: 30, IsActive: true},
			{ExternalID: "svc_hotel_manager", CategoryExternalID: "svc_hotel_travel", Name: "酒店管理", Tags: "服务,酒店", OrderNum: 40, IsActive: true},
			{ExternalID: "svc_frontdesk", CategoryExternalID: "svc_hotel_travel", Name: "酒店前台", Tags: "服务,酒店", OrderNum: 50, IsActive: true},
			{ExternalID: "svc_guest_service", CategoryExternalID: "svc_hotel_travel", Name: "客房服务", Tags: "服务,酒店", OrderNum: 60, IsActive: true},
			{ExternalID: "svc_porter", CategoryExternalID: "svc_hotel_travel", Name: "行李员", Tags: "服务,酒店", OrderNum: 70, IsActive: true},
			{ExternalID: "svc_travel_interpreter", CategoryExternalID: "svc_hotel_travel", Name: "景区讲解员", Tags: "服务,旅游", OrderNum: 80, IsActive: true},
			{ExternalID: "svc_travel_custom", CategoryExternalID: "svc_hotel_travel", Name: "旅游定制师", Tags: "服务,旅游", OrderNum: 90, IsActive: true},
			{ExternalID: "svc_greeter", CategoryExternalID: "svc_hotel_travel", Name: "前厅经理", Tags: "服务,酒店", OrderNum: 100, IsActive: true},
			{ExternalID: "svc_concierge", CategoryExternalID: "svc_hotel_travel", Name: "礼宾相关岗位", Tags: "服务,酒店", OrderNum: 110, IsActive: true},
			{ExternalID: "svc_hotel_sales", CategoryExternalID: "svc_hotel_travel", Name: "酒店销售", Tags: "服务,酒店,销售", OrderNum: 120, IsActive: true},
			{ExternalID: "svc_banquets", CategoryExternalID: "svc_hotel_travel", Name: "宴会服务", Tags: "服务,酒店", OrderNum: 130, IsActive: true},

			{ExternalID: "svc_cleaner", CategoryExternalID: "svc_personal", Name: "保洁", Tags: "服务,保洁", OrderNum: 10, IsActive: true},
			{ExternalID: "svc_guard", CategoryExternalID: "svc_personal", Name: "保安", Tags: "服务,保安", OrderNum: 20, IsActive: true},
			{ExternalID: "svc_housekeeping", CategoryExternalID: "svc_personal", Name: "家政", Tags: "服务,家政", OrderNum: 30, IsActive: true},
			{ExternalID: "svc_babysitter", CategoryExternalID: "svc_personal", Name: "保姆", Tags: "服务,家政", OrderNum: 40, IsActive: true},
			{ExternalID: "svc_maternity_matron", CategoryExternalID: "svc_personal", Name: "月嫂", Tags: "服务,家政", OrderNum: 50, IsActive: true},
			{ExternalID: "svc_housekeeper", CategoryExternalID: "svc_personal", Name: "管家", Tags: "服务,家政", OrderNum: 60, IsActive: true},
			{ExternalID: "svc_pet_beautician", CategoryExternalID: "svc_personal", Name: "宠物美容师", Tags: "服务,宠物", OrderNum: 70, IsActive: true},
			{ExternalID: "svc_pet_doctor", CategoryExternalID: "svc_personal", Name: "宠物医生", Tags: "服务,宠物,医生", OrderNum: 80, IsActive: true},
			{ExternalID: "svc_fashion_buyer", CategoryExternalID: "svc_personal", Name: "服装买手", Tags: "服务,服装", OrderNum: 90, IsActive: true},
			{ExternalID: "svc_stylist", CategoryExternalID: "svc_personal", Name: "服装搭配师", Tags: "服务,服装", OrderNum: 100, IsActive: true},

			{ExternalID: "svc_fitness_coach", CategoryExternalID: "svc_fitness", Name: "健身教练", Tags: "服务,健身", OrderNum: 10, IsActive: true},
			{ExternalID: "svc_swimming_coach", CategoryExternalID: "svc_fitness", Name: "游泳教练", Tags: "服务,健身", OrderNum: 20, IsActive: true},
			{ExternalID: "svc_yoga_coach", CategoryExternalID: "svc_fitness", Name: "瑜伽教练", Tags: "服务,健身", OrderNum: 30, IsActive: true},
			{ExternalID: "svc_dance_teacher", CategoryExternalID: "svc_fitness", Name: "舞蹈老师", Tags: "服务,健身", OrderNum: 40, IsActive: true},
			{ExternalID: "svc_basketball_coach", CategoryExternalID: "svc_fitness", Name: "篮球教练", Tags: "服务,健身", OrderNum: 50, IsActive: true},
			{ExternalID: "svc_badminton_coach", CategoryExternalID: "svc_fitness", Name: "羽毛球教练", Tags: "服务,健身", OrderNum: 60, IsActive: true},
			{ExternalID: "svc_taekwondo_coach", CategoryExternalID: "svc_fitness", Name: "跆拳道教练", Tags: "服务,健身", OrderNum: 70, IsActive: true},
			{ExternalID: "svc_martial_teacher", CategoryExternalID: "svc_fitness", Name: "武术教练", Tags: "服务,健身", OrderNum: 80, IsActive: true},
			{ExternalID: "svc_street_dance_teacher", CategoryExternalID: "svc_fitness", Name: "街舞老师", Tags: "服务,健身", OrderNum: 90, IsActive: true},

			{ExternalID: "svc_beauty_consultant", CategoryExternalID: "svc_beauty", Name: "美容顾问", Tags: "服务,美容", OrderNum: 10, IsActive: true},
			{ExternalID: "svc_store_manager", CategoryExternalID: "svc_beauty", Name: "美容店长", Tags: "服务,美容,管理", OrderNum: 20, IsActive: true},
			{ExternalID: "svc_beautician", CategoryExternalID: "svc_beauty", Name: "美容师", Tags: "服务,美容", OrderNum: 30, IsActive: true},
			{ExternalID: "svc_makeup_artist", CategoryExternalID: "svc_beauty", Name: "化妆师", Tags: "服务,美容", OrderNum: 40, IsActive: true},
			{ExternalID: "svc_manicurist", CategoryExternalID: "svc_beauty", Name: "美甲师", Tags: "服务,美容", OrderNum: 50, IsActive: true},
			{ExternalID: "svc_hairdresser", CategoryExternalID: "svc_beauty", Name: "美发师", Tags: "服务,美容", OrderNum: 60, IsActive: true},
			{ExternalID: "svc_masseur", CategoryExternalID: "svc_beauty", Name: "按摩师", Tags: "服务,美容", OrderNum: 70, IsActive: true},
			{ExternalID: "svc_physiotherapist", CategoryExternalID: "svc_beauty", Name: "理疗师", Tags: "服务,美容", OrderNum: 80, IsActive: true},
			{ExternalID: "svc_tattoo_artist", CategoryExternalID: "svc_beauty", Name: "纹绣师", Tags: "服务,美容", OrderNum: 90, IsActive: true},

			{ExternalID: "svc_course_consultant", CategoryExternalID: "svc_consult_translate", Name: "课程顾问", Tags: "服务,咨询", OrderNum: 10, IsActive: true},
			{ExternalID: "svc_headhunter", CategoryExternalID: "svc_consult_translate", Name: "猎头顾问", Tags: "服务,咨询", OrderNum: 20, IsActive: true},
			{ExternalID: "svc_consultant", CategoryExternalID: "svc_consult_translate", Name: "咨询顾问", Tags: "服务,咨询", OrderNum: 30, IsActive: true},
			{ExternalID: "svc_legal_consultant", CategoryExternalID: "svc_consult_translate", Name: "法律顾问", Tags: "服务,咨询", OrderNum: 40, IsActive: true},
			{ExternalID: "svc_translator", CategoryExternalID: "svc_consult_translate", Name: "翻译", Tags: "服务,翻译", OrderNum: 50, IsActive: true},

			{ExternalID: "svc_auto_mechanic", CategoryExternalID: "svc_repair", Name: "汽车维修", Tags: "服务,维修", OrderNum: 10, IsActive: true},
			{ExternalID: "svc_mobile_repair", CategoryExternalID: "svc_repair", Name: "机务维修", Tags: "服务,维修", OrderNum: 20, IsActive: true},
			{ExternalID: "svc_electric_repair", CategoryExternalID: "svc_repair", Name: "维修电工", Tags: "服务,维修", OrderNum: 30, IsActive: true},
			{ExternalID: "svc_mould_maintenance", CategoryExternalID: "svc_repair", Name: "模具维修", Tags: "服务,维修", OrderNum: 40, IsActive: true},
			{ExternalID: "svc_device_repair", CategoryExternalID: "svc_repair", Name: "器械维修", Tags: "服务,维修", OrderNum: 50, IsActive: true},
			{ExternalID: "svc_other_repair", CategoryExternalID: "svc_repair", Name: "其他维修服务岗", Tags: "服务,维修", OrderNum: 60, IsActive: true},

			{ExternalID: "svc_store_clerk", CategoryExternalID: "svc_other", Name: "店长", Tags: "服务,其他", OrderNum: 10, IsActive: true},
			{ExternalID: "svc_bid_staff", CategoryExternalID: "svc_other", Name: "项目招投标", Tags: "服务,其他", OrderNum: 20, IsActive: true},
			{ExternalID: "svc_tickets_staff", CategoryExternalID: "svc_other", Name: "票务", Tags: "服务,其他", OrderNum: 30, IsActive: true},
			{ExternalID: "svc_trustee", CategoryExternalID: "svc_other", Name: "托管", Tags: "服务,其他", OrderNum: 40, IsActive: true},
			{ExternalID: "svc_fee_collector", CategoryExternalID: "svc_other", Name: "收费员", Tags: "服务,其他", OrderNum: 50, IsActive: true},
			{ExternalID: "svc_transport_agent", CategoryExternalID: "svc_other", Name: "客运员", Tags: "服务,其他", OrderNum: 60, IsActive: true},
			{ExternalID: "svc_uav_pilot", CategoryExternalID: "svc_other", Name: "无人机飞手", Tags: "服务,其他,无人机", OrderNum: 70, IsActive: true},

			{ExternalID: "media_journalist", CategoryExternalID: "media_news_publishing", Name: "记者", Tags: "传媒,新闻", OrderNum: 10, IsActive: true},
			{ExternalID: "media_editor", CategoryExternalID: "media_news_publishing", Name: "编辑", Tags: "传媒,新闻", OrderNum: 20, IsActive: true},
			{ExternalID: "media_reporter", CategoryExternalID: "media_news_publishing", Name: "采编", Tags: "传媒,新闻", OrderNum: 30, IsActive: true},
			{ExternalID: "media_chief_editor", CategoryExternalID: "media_news_publishing", Name: "主编/副主编", Tags: "传媒,新闻", OrderNum: 40, IsActive: true},
			{ExternalID: "media_proofreader", CategoryExternalID: "media_news_publishing", Name: "校对", Tags: "传媒,新闻", OrderNum: 50, IsActive: true},
			{ExternalID: "media_writer", CategoryExternalID: "media_news_publishing", Name: "撰稿", Tags: "传媒,新闻", OrderNum: 60, IsActive: true},
			{ExternalID: "media_reviewer", CategoryExternalID: "media_news_publishing", Name: "审核", Tags: "传媒,新闻", OrderNum: 70, IsActive: true},

			{ExternalID: "media_host", CategoryExternalID: "media_broadcast_tv", Name: "主持人", Tags: "传媒,广电", OrderNum: 10, IsActive: true},
			{ExternalID: "media_announcer", CategoryExternalID: "media_broadcast_tv", Name: "播音员", Tags: "传媒,广电", OrderNum: 20, IsActive: true},
			{ExternalID: "media_director_tv", CategoryExternalID: "media_broadcast_tv", Name: "编导", Tags: "传媒,广电", OrderNum: 30, IsActive: true},
			{ExternalID: "media_cameraman", CategoryExternalID: "media_broadcast_tv", Name: "摄像师", Tags: "传媒,广电", OrderNum: 40, IsActive: true},
			{ExternalID: "media_camera_assistant", CategoryExternalID: "media_broadcast_tv", Name: "摄影助理", Tags: "传媒,广电", OrderNum: 50, IsActive: true},
			{ExternalID: "media_postprod", CategoryExternalID: "media_broadcast_tv", Name: "后期制作", Tags: "传媒,后期", OrderNum: 60, IsActive: true},
			{ExternalID: "media_vfx", CategoryExternalID: "media_broadcast_tv", Name: "特效师", Tags: "传媒,后期", OrderNum: 70, IsActive: true},
			{ExternalID: "media_sound_engineer", CategoryExternalID: "media_broadcast_tv", Name: "音效师", Tags: "传媒,后期", OrderNum: 80, IsActive: true},
			{ExternalID: "media_switcher", CategoryExternalID: "media_broadcast_tv", Name: "导播", Tags: "传媒,广电", OrderNum: 90, IsActive: true},
			{ExternalID: "media_program_planner", CategoryExternalID: "media_broadcast_tv", Name: "节目策划", Tags: "传媒,广电", OrderNum: 100, IsActive: true},
			{ExternalID: "media_producer", CategoryExternalID: "media_broadcast_tv", Name: "制片人", Tags: "传媒,广电", OrderNum: 110, IsActive: true},
			{ExternalID: "media_film_making", CategoryExternalID: "media_broadcast_tv", Name: "影视制作", Tags: "传媒,广电", OrderNum: 120, IsActive: true},
			{ExternalID: "media_channel_specialist", CategoryExternalID: "media_broadcast_tv", Name: "渠道专员", Tags: "传媒,广电", OrderNum: 130, IsActive: true},

			{ExternalID: "media_director", CategoryExternalID: "media_film_performance", Name: "导演", Tags: "传媒,影视", OrderNum: 10, IsActive: true},
			{ExternalID: "media_screenwriter", CategoryExternalID: "media_film_performance", Name: "编剧", Tags: "传媒,影视", OrderNum: 20, IsActive: true},
			{ExternalID: "media_actor", CategoryExternalID: "media_film_performance", Name: "演员", Tags: "传媒,影视", OrderNum: 30, IsActive: true},
			{ExternalID: "media_script_supervisor", CategoryExternalID: "media_film_performance", Name: "场记", Tags: "传媒,影视", OrderNum: 40, IsActive: true},
			{ExternalID: "media_artist_assistant", CategoryExternalID: "media_film_performance", Name: "艺人助理", Tags: "传媒,影视", OrderNum: 50, IsActive: true},
			{ExternalID: "media_agent", CategoryExternalID: "media_film_performance", Name: "经纪人", Tags: "传媒,影视", OrderNum: 60, IsActive: true},
			{ExternalID: "media_model", CategoryExternalID: "media_film_performance", Name: "模特", Tags: "传媒,影视", OrderNum: 70, IsActive: true},
			{ExternalID: "media_stage_designer", CategoryExternalID: "media_film_performance", Name: "舞美设计", Tags: "传媒,影视", OrderNum: 80, IsActive: true},
			{ExternalID: "media_star_mapper", CategoryExternalID: "media_film_performance", Name: "星探", Tags: "传媒,影视", OrderNum: 90, IsActive: true},
			{ExternalID: "media_intern", CategoryExternalID: "media_film_performance", Name: "练习生", Tags: "传媒,影视", OrderNum: 100, IsActive: true},

			{ExternalID: "trade_business_staff", CategoryExternalID: "trade_foreign_trade", Name: "外贸业务员", Tags: "外贸,业务", OrderNum: 10, IsActive: true},
			{ExternalID: "trade_specialist", CategoryExternalID: "trade_foreign_trade", Name: "外贸专员", Tags: "外贸,业务", OrderNum: 20, IsActive: true},
			{ExternalID: "trade_doc_tracker", CategoryExternalID: "trade_foreign_trade", Name: "外贸跟单员", Tags: "外贸,跟单", OrderNum: 30, IsActive: true},
			{ExternalID: "trade_assistant", CategoryExternalID: "trade_foreign_trade", Name: "外贸助理", Tags: "外贸,助理", OrderNum: 40, IsActive: true},
			{ExternalID: "trade_manager", CategoryExternalID: "trade_foreign_trade", Name: "外贸经理", Tags: "外贸,管理", OrderNum: 50, IsActive: true},
			{ExternalID: "trade_doc_staff", CategoryExternalID: "trade_foreign_trade", Name: "外贸单证员", Tags: "外贸,单证", OrderNum: 60, IsActive: true},
			{ExternalID: "trade_cs", CategoryExternalID: "trade_foreign_trade", Name: "外贸客服", Tags: "外贸,客服", OrderNum: 70, IsActive: true},
			{ExternalID: "trade_director", CategoryExternalID: "trade_foreign_trade", Name: "外贸总监", Tags: "外贸,管理", OrderNum: 80, IsActive: true},

			{ExternalID: "trade_customs_broker", CategoryExternalID: "trade_trade_support", Name: "报关员", Tags: "外贸,报关", OrderNum: 10, IsActive: true},
			{ExternalID: "trade_customs_supervisor", CategoryExternalID: "trade_trade_support", Name: "报关主管", Tags: "外贸,报关", OrderNum: 20, IsActive: true},
			{ExternalID: "trade_docs", CategoryExternalID: "trade_trade_support", Name: "单证员", Tags: "外贸,单证", OrderNum: 30, IsActive: true},

			{ExternalID: "trade_cbe", CategoryExternalID: "trade_crossborder_ecom", Name: "跨境电商", Tags: "外贸,电商", OrderNum: 10, IsActive: true},
			{ExternalID: "trade_cbe_specialist", CategoryExternalID: "trade_crossborder_ecom", Name: "跨境电商专员", Tags: "外贸,电商", OrderNum: 20, IsActive: true},
			{ExternalID: "trade_cbe_ops", CategoryExternalID: "trade_crossborder_ecom", Name: "跨境电商业务员", Tags: "外贸,电商", OrderNum: 30, IsActive: true},
			{ExternalID: "trade_cbe_assistant", CategoryExternalID: "trade_crossborder_ecom", Name: "跨境电商运营助理", Tags: "外贸,电商", OrderNum: 40, IsActive: true},
			{ExternalID: "trade_amazon_ops", CategoryExternalID: "trade_crossborder_ecom", Name: "Amazon运营", Tags: "外贸,电商,Amazon", OrderNum: 50, IsActive: true},
			{ExternalID: "trade_amazon_sales", CategoryExternalID: "trade_crossborder_ecom", Name: "亚马逊销售", Tags: "外贸,电商,Amazon", OrderNum: 60, IsActive: true},

			{ExternalID: "trade_translator", CategoryExternalID: "trade_translation_support", Name: "外贸翻译", Tags: "外贸,翻译", OrderNum: 10, IsActive: true},
			{ExternalID: "trade_en_translator", CategoryExternalID: "trade_translation_support", Name: "英语翻译", Tags: "翻译,英语", OrderNum: 20, IsActive: true},
			{ExternalID: "trade_jp_translator", CategoryExternalID: "trade_translation_support", Name: "日语翻译", Tags: "翻译,日语", OrderNum: 30, IsActive: true},
			{ExternalID: "trade_kr_translator", CategoryExternalID: "trade_translation_support", Name: "韩语翻译", Tags: "翻译,韩语", OrderNum: 40, IsActive: true},

			{ExternalID: "energy_power_engineer", CategoryExternalID: "energy_traditional", Name: "电力工程师", Tags: "能源,电力", OrderNum: 10, IsActive: true},
			{ExternalID: "energy_new_energy_engineer", CategoryExternalID: "energy_traditional", Name: "新能源工程师", Tags: "能源,新能源", OrderNum: 20, IsActive: true},
			{ExternalID: "energy_thermal_engineer", CategoryExternalID: "energy_traditional", Name: "热能工程师", Tags: "能源,热能", OrderNum: 30, IsActive: true},
			{ExternalID: "energy_oil_engineer", CategoryExternalID: "energy_traditional", Name: "石油工程师", Tags: "能源,石油", OrderNum: 40, IsActive: true},
			{ExternalID: "energy_gas_engineer", CategoryExternalID: "energy_traditional", Name: "燃气工程师", Tags: "能源,燃气", OrderNum: 50, IsActive: true},
			{ExternalID: "energy_hvac_engineer", CategoryExternalID: "energy_traditional", Name: "暖通工程师", Tags: "能源,暖通", OrderNum: 60, IsActive: true},
			{ExternalID: "energy_engineer", CategoryExternalID: "energy_traditional", Name: "能源工程师", Tags: "能源", OrderNum: 70, IsActive: true},
			{ExternalID: "energy_pv_engineer", CategoryExternalID: "energy_traditional", Name: "光伏系统工程师", Tags: "能源,光伏", OrderNum: 80, IsActive: true},
			{ExternalID: "energy_wind_engineer", CategoryExternalID: "energy_traditional", Name: "风电工程师", Tags: "能源,风电", OrderNum: 90, IsActive: true},

			{ExternalID: "env_engineer", CategoryExternalID: "energy_environment", Name: "环保工程师", Tags: "环保", OrderNum: 10, IsActive: true},
			{ExternalID: "env_environment_engineer", CategoryExternalID: "energy_environment", Name: "环境工程师", Tags: "环保,环境", OrderNum: 20, IsActive: true},
			{ExternalID: "env_ehs_engineer", CategoryExternalID: "energy_environment", Name: "EHS工程师", Tags: "环保,EHS", OrderNum: 30, IsActive: true},
			{ExternalID: "env_water_treatment", CategoryExternalID: "energy_environment", Name: "水处理工程师", Tags: "环保,水处理", OrderNum: 40, IsActive: true},
			{ExternalID: "env_mep_water", CategoryExternalID: "energy_environment", Name: "给排水工程师", Tags: "环保,给排水", OrderNum: 50, IsActive: true},
			{ExternalID: "env_eia_engineer", CategoryExternalID: "energy_environment", Name: "环评工程师", Tags: "环保,环评", OrderNum: 60, IsActive: true},
			{ExternalID: "env_tech", CategoryExternalID: "energy_environment", Name: "环保技术员", Tags: "环保,技术", OrderNum: 70, IsActive: true},
			{ExternalID: "env_inspection", CategoryExternalID: "energy_environment", Name: "环保检测", Tags: "环保,检测", OrderNum: 80, IsActive: true},
			{ExternalID: "env_specialist", CategoryExternalID: "energy_environment", Name: "环保专员", Tags: "环保", OrderNum: 90, IsActive: true},
			{ExternalID: "env_supervisor", CategoryExternalID: "energy_environment", Name: "环保主管", Tags: "环保,管理", OrderNum: 100, IsActive: true},

			{ExternalID: "agri_tech", CategoryExternalID: "agri_planting", Name: "农业技术员", Tags: "农业", OrderNum: 10, IsActive: true},
			{ExternalID: "agri_agronomist", CategoryExternalID: "agri_planting", Name: "农艺师", Tags: "农业", OrderNum: 20, IsActive: true},
			{ExternalID: "agri_horticulturist", CategoryExternalID: "agri_planting", Name: "园艺师", Tags: "农业,园艺", OrderNum: 30, IsActive: true},
			{ExternalID: "agri_florist", CategoryExternalID: "agri_planting", Name: "花艺师", Tags: "农业,花艺", OrderNum: 40, IsActive: true},
			{ExternalID: "agri_farm_machine", CategoryExternalID: "agri_planting", Name: "农机操作修理", Tags: "农业,农机", OrderNum: 50, IsActive: true},

			{ExternalID: "forestry_engineer", CategoryExternalID: "agri_forestry", Name: "林业工程师", Tags: "林业", OrderNum: 10, IsActive: true},
			{ExternalID: "forestry_tech", CategoryExternalID: "agri_forestry", Name: "林业技术员", Tags: "林业", OrderNum: 20, IsActive: true},
			{ExternalID: "forestry_garden_engineer", CategoryExternalID: "agri_forestry", Name: "园林工程师", Tags: "林业,园林", OrderNum: 30, IsActive: true},
			{ExternalID: "forestry_ranger", CategoryExternalID: "agri_forestry", Name: "护林员", Tags: "林业", OrderNum: 40, IsActive: true},

			{ExternalID: "livestock_specialist", CategoryExternalID: "agri_livestock", Name: "畜牧师", Tags: "畜牧", OrderNum: 10, IsActive: true},
			{ExternalID: "livestock_vet", CategoryExternalID: "agri_livestock", Name: "兽医", Tags: "畜牧,兽医", OrderNum: 20, IsActive: true},
			{ExternalID: "livestock_breeding_tech", CategoryExternalID: "agri_livestock", Name: "养殖技术员", Tags: "畜牧,养殖", OrderNum: 30, IsActive: true},
			{ExternalID: "livestock_farm_mgmt", CategoryExternalID: "agri_livestock", Name: "牧场管理", Tags: "畜牧,管理", OrderNum: 40, IsActive: true},
			{ExternalID: "livestock_feeder", CategoryExternalID: "agri_livestock", Name: "饲养员", Tags: "畜牧,饲养", OrderNum: 50, IsActive: true},
			{ExternalID: "livestock_quarantine", CategoryExternalID: "agri_livestock", Name: "动物检疫员", Tags: "畜牧,检疫", OrderNum: 60, IsActive: true},
			{ExternalID: "livestock_trainer", CategoryExternalID: "agri_livestock", Name: "动物驯养师", Tags: "畜牧,驯养", OrderNum: 70, IsActive: true},
			{ExternalID: "livestock_guide", CategoryExternalID: "agri_livestock", Name: "动物讲解员", Tags: "畜牧,讲解", OrderNum: 80, IsActive: true},

			{ExternalID: "fishery_farmer", CategoryExternalID: "agri_fishery", Name: "水产养殖员", Tags: "水产,养殖", OrderNum: 10, IsActive: true},
			{ExternalID: "fishery_tech", CategoryExternalID: "agri_fishery", Name: "水产技术员", Tags: "水产", OrderNum: 20, IsActive: true},
			{ExternalID: "fishery_aquaculture", CategoryExternalID: "agri_fishery", Name: "渔业养殖员", Tags: "水产,养殖", OrderNum: 30, IsActive: true},
			{ExternalID: "fishery_fisher", CategoryExternalID: "agri_fishery", Name: "捕捞员", Tags: "水产,捕捞", OrderNum: 40, IsActive: true},

			{ExternalID: "public_civil_servant", CategoryExternalID: "public_services", Name: "公务员", Tags: "公共服务", OrderNum: 10, IsActive: true},
			{ExternalID: "public_ngo_specialist", CategoryExternalID: "public_services", Name: "非营利组织专员", Tags: "公共服务,NGO", OrderNum: 20, IsActive: true},
			{ExternalID: "public_urban_mgmt", CategoryExternalID: "public_services", Name: "城管", Tags: "公共服务", OrderNum: 30, IsActive: true},
			{ExternalID: "public_military", CategoryExternalID: "public_services", Name: "军人", Tags: "公共服务", OrderNum: 40, IsActive: true},
			{ExternalID: "public_firefighter", CategoryExternalID: "public_services", Name: "消防员", Tags: "公共服务", OrderNum: 50, IsActive: true},
			{ExternalID: "public_community_worker", CategoryExternalID: "public_services", Name: "社区工作者", Tags: "公共服务", OrderNum: 60, IsActive: true},
			{ExternalID: "public_police", CategoryExternalID: "public_services", Name: "警察/辅警", Tags: "公共服务", OrderNum: 70, IsActive: true},
			{ExternalID: "public_party_affairs", CategoryExternalID: "public_services", Name: "党建管理岗", Tags: "公共服务,党建", OrderNum: 80, IsActive: true},

			{ExternalID: "research_assistant", CategoryExternalID: "public_research", Name: "科研助理", Tags: "科研", OrderNum: 10, IsActive: true},
			{ExternalID: "research_staff", CategoryExternalID: "public_research", Name: "科研人员", Tags: "科研", OrderNum: 20, IsActive: true},
			{ExternalID: "research_management", CategoryExternalID: "public_research", Name: "科研管理", Tags: "科研,管理", OrderNum: 30, IsActive: true},
			{ExternalID: "research_academic_promo", CategoryExternalID: "public_research", Name: "学术推广", Tags: "科研,推广", OrderNum: 40, IsActive: true},
			{ExternalID: "research_chem_analysis", CategoryExternalID: "public_research", Name: "化学分析", Tags: "科研,化学", OrderNum: 50, IsActive: true},

			{ExternalID: "social_worker", CategoryExternalID: "public_social", Name: "社工", Tags: "社会服务", OrderNum: 10, IsActive: true},
			{ExternalID: "social_worker_assistant", CategoryExternalID: "public_social", Name: "社工助理", Tags: "社会服务", OrderNum: 20, IsActive: true},
			{ExternalID: "social_worker_intern", CategoryExternalID: "public_social", Name: "社工实习生", Tags: "社会服务,实习", OrderNum: 30, IsActive: true},
			{ExternalID: "volunteer", CategoryExternalID: "public_social", Name: "志愿者", Tags: "社会服务", OrderNum: 40, IsActive: true},
			{ExternalID: "unpaid_volunteer", CategoryExternalID: "public_social", Name: "义工", Tags: "社会服务", OrderNum: 50, IsActive: true},
			{ExternalID: "caregiver", CategoryExternalID: "public_social", Name: "护理员", Tags: "社会服务,护理", OrderNum: 60, IsActive: true},
			{ExternalID: "rehab_therapist", CategoryExternalID: "public_social", Name: "康复师", Tags: "社会服务,康复", OrderNum: 70, IsActive: true},
		},
		Presets: []preset.ContentPreset{
			{ExternalID: "it_backend_java_zh", Name: "Java 开发（中文示例）", Language: "zh", RoleExternalID: "java", Tags: "Java,后端,中文", DataJSON: string(b), IsActive: true},
		},
		Variants: []library.TemplateVariant{
			{ExternalID: "mint_it_backend_java_zh", Name: "青色时间轴 - Java 开发", LayoutTemplateExternalID: "TemplateMintTimeline", PresetExternalID: "it_backend_java_zh", RoleExternalID: "java", Tags: "Java,后端,中文", UsageCount: 0, IsPremium: false, IsActive: true},
			{ExternalID: "classic_it_backend_java_zh", Name: "经典专业版 - Java 开发", LayoutTemplateExternalID: "TemplateClassic", PresetExternalID: "it_backend_java_zh", RoleExternalID: "java", Tags: "Java,后端,中文", UsageCount: 0, IsPremium: false, IsActive: true},
		},
	}, nil
}
