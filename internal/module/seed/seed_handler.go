package seed

import (
	"net/http"

	"openresume/internal/infra/database"
	"openresume/internal/module/preset"
	"openresume/internal/module/seed/presets"
	"openresume/internal/module/taxonomy"
	"openresume/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SeedData struct {
	Categories []SeedJobCategory
	Roles      []SeedJobRole
	Presets    []SeedContentPreset
}

type SeedJobCategory struct {
	ExternalID       string
	Name             string
	ParentExternalID string
	OrderNum         int
	IsActive         bool
}

type SeedJobRole struct {
	ExternalID         string
	CategoryExternalID string
	Name               string
	OrderNum           int
	IsActive           bool
}

type SeedContentPreset struct {
	Name     string
	Language string
	RoleCode string
	DataJSON string
	IsActive bool
}

type ImportCounts struct {
	JobCategories  int `json:"jobCategories"`
	JobRoles       int `json:"jobRoles"`
	ContentPresets int `json:"contentPresets"`
}

type AdminHandler struct{}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

func (h *AdminHandler) AdminImportDefault(c *gin.Context) {
	seed := DefaultSeed()

	var counts ImportCounts
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		taxRepo := taxonomy.NewRepo(tx)
		presetRepo := preset.NewRepo(tx)

		categoryIDByExternal := make(map[string]uint, len(seed.Categories))
		for _, sc := range seed.Categories {
			var parentID *uint
			if sc.ParentExternalID != "" {
				pid, ok := categoryIDByExternal[sc.ParentExternalID]
				if !ok || pid == 0 {
					return gorm.ErrInvalidData
				}
				parentID = &pid
			}
			m := taxonomy.JobCategory{
				Name:     sc.Name,
				ParentID: parentID,
				OrderNum: sc.OrderNum,
				IsActive: sc.IsActive,
			}
			if err := taxRepo.UpsertJobCategory(tx, &m); err != nil {
				return err
			}
			categoryIDByExternal[sc.ExternalID] = m.ID
			counts.JobCategories++
		}

		roleIDByExternal := make(map[string]uint, len(seed.Roles))
		for _, sr := range seed.Roles {
			cid, ok := categoryIDByExternal[sr.CategoryExternalID]
			if !ok || cid == 0 {
				return gorm.ErrInvalidData
			}
			m := taxonomy.JobRole{
				CategoryID: cid,
				Name:       sr.Name,
				OrderNum:   sr.OrderNum,
				IsActive:   sr.IsActive,
			}
			if err := taxRepo.UpsertJobRole(tx, &m); err != nil {
				return err
			}
			roleIDByExternal[sr.ExternalID] = m.ID
			counts.JobRoles++
		}

		for _, sp := range seed.Presets {
			rid, ok := roleIDByExternal[sp.RoleCode]
			if !ok || rid == 0 {
				return gorm.ErrInvalidData
			}
			m := preset.ContentPreset{
				Name:     sp.Name,
				Language: sp.Language,
				RoleID:   rid,
				DataJSON: sp.DataJSON,
				IsActive: sp.IsActive,
			}
			if err := presetRepo.UpsertContentPreset(tx, &m); err != nil {
				return err
			}
			counts.ContentPresets++
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

func DefaultSeed() SeedData {
	return SeedData{
		Categories: seedCategories,
		Roles:      seedRoles,
		Presets:    seedPresets,
	}
}

var seedPresets = []SeedContentPreset{
	{Name: "Java 开发（中文示例）", Language: "zh", RoleCode: "Java", DataJSON: string(presets.GenerateJavaPreset()), IsActive: true},
	{Name: "Python 开发（中文示例）", Language: "zh", RoleCode: "Python", DataJSON: string(presets.GeneratePythonPreset()), IsActive: true},
	{Name: "Go 开发（中文示例）", Language: "zh", RoleCode: "golang", DataJSON: string(presets.GenerateGolangPreset()), IsActive: true},
	{Name: "PHP 开发（中文示例）", Language: "zh", RoleCode: "php", DataJSON: string(presets.GeneratePHPPreset()), IsActive: true},
	{Name: "C/C++ 开发（中文示例）", Language: "zh", RoleCode: "c_cpp", DataJSON: string(presets.GenerateCCppPreset()), IsActive: true},
	{Name: "C# 开发（中文示例）", Language: "zh", RoleCode: "csharp", DataJSON: string(presets.GenerateCSharpPreset()), IsActive: true},
	{Name: ".NET 开发（中文示例）", Language: "zh", RoleCode: "dotnet", DataJSON: string(presets.GenerateDotnetPreset()), IsActive: true},
	{Name: "Node.js 开发（中文示例）", Language: "zh", RoleCode: "nodejs", DataJSON: string(presets.GenerateNodejsPreset()), IsActive: true},
}

var seedRoles = []SeedJobRole{
	{ExternalID: "Java", CategoryExternalID: "it_backend", Name: "Java", OrderNum: 10, IsActive: true},
	{ExternalID: "Python", CategoryExternalID: "it_backend", Name: "Python", OrderNum: 20, IsActive: true},
	{ExternalID: "golang", CategoryExternalID: "it_backend", Name: "Go (Golang)", OrderNum: 30, IsActive: true},
	{ExternalID: "php", CategoryExternalID: "it_backend", Name: "PHP", OrderNum: 40, IsActive: true},
	{ExternalID: "c_cpp", CategoryExternalID: "it_backend", Name: "C/C++", OrderNum: 50, IsActive: true},
	{ExternalID: "csharp", CategoryExternalID: "it_backend", Name: "C#", OrderNum: 60, IsActive: true},
	{ExternalID: "dotnet", CategoryExternalID: "it_backend", Name: ".NET", OrderNum: 70, IsActive: true},
	{ExternalID: "nodejs", CategoryExternalID: "it_backend", Name: "node.js", OrderNum: 80, IsActive: true},

	{ExternalID: "web_frontend", CategoryExternalID: "it_frontend", Name: "Web前端", OrderNum: 10, IsActive: true},
	{ExternalID: "html5", CategoryExternalID: "it_frontend", Name: "HTML5", OrderNum: 20, IsActive: true},
	{ExternalID: "miniapp", CategoryExternalID: "it_frontend", Name: "小程序开发工程师", OrderNum: 30, IsActive: true},

	{ExternalID: "android", CategoryExternalID: "it_mobile", Name: "Android开发", OrderNum: 10, IsActive: true},
	{ExternalID: "ios", CategoryExternalID: "it_mobile", Name: "iOS开发", OrderNum: 20, IsActive: true},
	{ExternalID: "harmony", CategoryExternalID: "it_mobile", Name: "鸿蒙开发工程师", OrderNum: 30, IsActive: true},

	{ExternalID: "test_engineer", CategoryExternalID: "it_testing", Name: "测试工程师", OrderNum: 10, IsActive: true},
	{ExternalID: "automation_test", CategoryExternalID: "it_testing", Name: "自动化测试", OrderNum: 20, IsActive: true},
	{ExternalID: "test_dev", CategoryExternalID: "it_testing", Name: "测试开发", OrderNum: 30, IsActive: true},
	{ExternalID: "performance_test", CategoryExternalID: "it_testing", Name: "性能测试", OrderNum: 40, IsActive: true},
	{ExternalID: "hardware_test", CategoryExternalID: "it_testing", Name: "硬件测试工程师", OrderNum: 50, IsActive: true},

	{ExternalID: "ops_engineer", CategoryExternalID: "it_ops_sec_dba", Name: "运维工程师", OrderNum: 10, IsActive: true},
	{ExternalID: "devops", CategoryExternalID: "it_ops_sec_dba", Name: "DevOps工程师", OrderNum: 20, IsActive: true},
	{ExternalID: "sys_admin", CategoryExternalID: "it_ops_sec_dba", Name: "系统/网络管理员", OrderNum: 30, IsActive: true},
	{ExternalID: "dba", CategoryExternalID: "it_ops_sec_dba", Name: "数据库管理员 (DBA)", OrderNum: 40, IsActive: true},
	{ExternalID: "security_engineer", CategoryExternalID: "it_ops_sec_dba", Name: "安全工程师", OrderNum: 50, IsActive: true},
	{ExternalID: "cloud_engineer", CategoryExternalID: "it_ops_sec_dba", Name: "云计算工程师", OrderNum: 60, IsActive: true},
	{ExternalID: "ops_manager", CategoryExternalID: "it_ops_sec_dba", Name: "运维经理/主管", OrderNum: 70, IsActive: true},

	{ExternalID: "data_mining", CategoryExternalID: "it_ai_bigdata", Name: "数据挖掘", OrderNum: 10, IsActive: true},
	{ExternalID: "nlp", CategoryExternalID: "it_ai_bigdata", Name: "自然语言处理", OrderNum: 20, IsActive: true},
	{ExternalID: "ml_ai", CategoryExternalID: "it_ai_bigdata", Name: "机器学习/AI工程师", OrderNum: 30, IsActive: true},
	{ExternalID: "bigdata", CategoryExternalID: "it_ai_bigdata", Name: "大数据工程师", OrderNum: 40, IsActive: true},
	{ExternalID: "blockchain", CategoryExternalID: "it_ai_bigdata", Name: "区块链开发", OrderNum: 50, IsActive: true},
	{ExternalID: "algo_engineer", CategoryExternalID: "it_ai_bigdata", Name: "算法工程师", OrderNum: 60, IsActive: true},

	{ExternalID: "tech_support", CategoryExternalID: "it_other_tech", Name: "技术支持工程师", OrderNum: 10, IsActive: true},
	{ExternalID: "tech_engineer", CategoryExternalID: "it_other_tech", Name: "技术工程师", OrderNum: 20, IsActive: true},
	{ExternalID: "presales_support", CategoryExternalID: "it_other_tech", Name: "售前售后/技术支持", OrderNum: 30, IsActive: true},
	{ExternalID: "other_tech_roles", CategoryExternalID: "it_other_tech", Name: "其他技术岗位", OrderNum: 40, IsActive: true},
	{ExternalID: "network_engineer", CategoryExternalID: "it_other_tech", Name: "网络工程师", OrderNum: 50, IsActive: true},
	{ExternalID: "hardware_dev", CategoryExternalID: "it_other_tech", Name: "硬件开发工程师", OrderNum: 60, IsActive: true},
	{ExternalID: "system_integrator", CategoryExternalID: "it_other_tech", Name: "系统集成工程师", OrderNum: 70, IsActive: true},
	{ExternalID: "circuit_engineer", CategoryExternalID: "it_other_tech", Name: "电路工程师", OrderNum: 80, IsActive: true},

	{ExternalID: "architect", CategoryExternalID: "it_senior", Name: "架构师", OrderNum: 10, IsActive: true},
	{ExternalID: "tech_manager", CategoryExternalID: "it_senior", Name: "技术主管", OrderNum: 20, IsActive: true},
	{ExternalID: "tech_director", CategoryExternalID: "it_senior", Name: "技术经理/总监", OrderNum: 30, IsActive: true},
	{ExternalID: "rd_manager", CategoryExternalID: "it_senior", Name: "研发经理/总监", OrderNum: 40, IsActive: true},
	{ExternalID: "cto", CategoryExternalID: "it_senior", Name: "CTO", OrderNum: 50, IsActive: true},
	{ExternalID: "fullstack", CategoryExternalID: "it_senior", Name: "全栈工程师", OrderNum: 60, IsActive: true},

	{ExternalID: "finance_teller", CategoryExternalID: "finance_counter_service", Name: "银行柜员", OrderNum: 10, IsActive: true},
	{ExternalID: "finance_general_teller", CategoryExternalID: "finance_counter_service", Name: "银行综合柜员", OrderNum: 20, IsActive: true},
	{ExternalID: "finance_lobby_manager", CategoryExternalID: "finance_counter_service", Name: "银行大堂经理", OrderNum: 30, IsActive: true},
	{ExternalID: "finance_lobby_guide", CategoryExternalID: "finance_counter_service", Name: "大堂引导员/大堂助理", OrderNum: 40, IsActive: true},
	{ExternalID: "finance_bank_cs", CategoryExternalID: "finance_counter_service", Name: "银行客服/坐席员", OrderNum: 50, IsActive: true},
	{ExternalID: "finance_bank_frontdesk", CategoryExternalID: "finance_counter_service", Name: "银行前台", OrderNum: 60, IsActive: true},

	{ExternalID: "finance_account_manager", CategoryExternalID: "finance_personal_wealth", Name: "客户经理", OrderNum: 10, IsActive: true},
	{ExternalID: "finance_wealth_manager", CategoryExternalID: "finance_personal_wealth", Name: "理财经理", OrderNum: 20, IsActive: true},
	{ExternalID: "finance_wealth_advisor", CategoryExternalID: "finance_personal_wealth", Name: "理财顾问", OrderNum: 30, IsActive: true},
	{ExternalID: "finance_investment_advisor", CategoryExternalID: "finance_personal_wealth", Name: "投资顾问", OrderNum: 40, IsActive: true},
	{ExternalID: "finance_creditcard_sales", CategoryExternalID: "finance_personal_wealth", Name: "信用卡销售", OrderNum: 50, IsActive: true},

	{ExternalID: "finance_credit_manager", CategoryExternalID: "finance_credit_approval", Name: "信贷经理", OrderNum: 10, IsActive: true},
	{ExternalID: "finance_credit_officer", CategoryExternalID: "finance_credit_approval", Name: "信贷专员", OrderNum: 20, IsActive: true},
	{ExternalID: "finance_loan_officer", CategoryExternalID: "finance_credit_approval", Name: "贷款专员", OrderNum: 30, IsActive: true},
	{ExternalID: "finance_postloan", CategoryExternalID: "finance_credit_approval", Name: "贷后管理岗", OrderNum: 40, IsActive: true},
	{ExternalID: "finance_collections", CategoryExternalID: "finance_credit_approval", Name: "催收岗", OrderNum: 50, IsActive: true},
	{ExternalID: "finance_mortgage_officer", CategoryExternalID: "finance_credit_approval", Name: "按揭专员", OrderNum: 60, IsActive: true},
	{ExternalID: "finance_credit_admin", CategoryExternalID: "finance_credit_approval", Name: "信贷管理", OrderNum: 70, IsActive: true},

	{ExternalID: "finance_risk_manager", CategoryExternalID: "finance_risk_compliance", Name: "风险经理", OrderNum: 10, IsActive: true},
	{ExternalID: "finance_compliance_manager", CategoryExternalID: "finance_risk_compliance", Name: "合规经理", OrderNum: 20, IsActive: true},
	{ExternalID: "finance_risk_control", CategoryExternalID: "finance_risk_compliance", Name: "风控专员", OrderNum: 30, IsActive: true},
	{ExternalID: "finance_audit", CategoryExternalID: "finance_risk_compliance", Name: "审计专员", OrderNum: 40, IsActive: true},
	{ExternalID: "finance_fin_compliance", CategoryExternalID: "finance_risk_compliance", Name: "金融合规专员", OrderNum: 50, IsActive: true},
	{ExternalID: "finance_aml", CategoryExternalID: "finance_risk_compliance", Name: "反洗钱专员", OrderNum: 60, IsActive: true},
	{ExternalID: "finance_legal", CategoryExternalID: "finance_risk_compliance", Name: "法律事务岗", OrderNum: 70, IsActive: true},

	{ExternalID: "finance_securities_broker", CategoryExternalID: "finance_securities_invest", Name: "证券经纪人", OrderNum: 10, IsActive: true},
	{ExternalID: "finance_invest_manager", CategoryExternalID: "finance_securities_invest", Name: "投资经理", OrderNum: 20, IsActive: true},
	{ExternalID: "finance_fund_manager", CategoryExternalID: "finance_securities_invest", Name: "基金经理", OrderNum: 30, IsActive: true},
	{ExternalID: "finance_trader", CategoryExternalID: "finance_securities_invest", Name: "交易员", OrderNum: 40, IsActive: true},
	{ExternalID: "finance_fund_accountant", CategoryExternalID: "finance_securities_invest", Name: "基金会计", OrderNum: 50, IsActive: true},
	{ExternalID: "finance_settlement", CategoryExternalID: "finance_securities_invest", Name: "清算专员", OrderNum: 60, IsActive: true},
	{ExternalID: "finance_researcher", CategoryExternalID: "finance_securities_invest", Name: "研究员", OrderNum: 70, IsActive: true},
	{ExternalID: "finance_analyst", CategoryExternalID: "finance_securities_invest", Name: "金融分析师", OrderNum: 80, IsActive: true},
	{ExternalID: "finance_securities_analyst", CategoryExternalID: "finance_securities_invest", Name: "证券分析师", OrderNum: 90, IsActive: true},

	{ExternalID: "finance_insurance_advisor", CategoryExternalID: "finance_insurance_actuary", Name: "保险顾问", OrderNum: 10, IsActive: true},
	{ExternalID: "finance_insurance_broker", CategoryExternalID: "finance_insurance_actuary", Name: "保险经纪人", OrderNum: 20, IsActive: true},
	{ExternalID: "finance_insurance_agent", CategoryExternalID: "finance_insurance_actuary", Name: "保险代理专员", OrderNum: 30, IsActive: true},
	{ExternalID: "finance_underwriter", CategoryExternalID: "finance_insurance_actuary", Name: "核保师", OrderNum: 40, IsActive: true},
	{ExternalID: "finance_claims", CategoryExternalID: "finance_insurance_actuary", Name: "理赔师", OrderNum: 50, IsActive: true},
	{ExternalID: "finance_actuary", CategoryExternalID: "finance_insurance_actuary", Name: "精算师", OrderNum: 60, IsActive: true},
	{ExternalID: "finance_insurance_trainer", CategoryExternalID: "finance_insurance_actuary", Name: "保险培训师", OrderNum: 70, IsActive: true},
	{ExternalID: "finance_insurance_ops", CategoryExternalID: "finance_insurance_actuary", Name: "保险内勤", OrderNum: 80, IsActive: true},
	{ExternalID: "finance_insurance_coach", CategoryExternalID: "finance_insurance_actuary", Name: "保险组训", OrderNum: 90, IsActive: true},
	{ExternalID: "finance_surveyor", CategoryExternalID: "finance_insurance_actuary", Name: "查勘员", OrderNum: 100, IsActive: true},
	{ExternalID: "finance_insurance_sales", CategoryExternalID: "finance_insurance_actuary", Name: "保险销售", OrderNum: 110, IsActive: true},

	{ExternalID: "finance_bank_ops_lead", CategoryExternalID: "finance_banking_support", Name: "银行运营主管", OrderNum: 10, IsActive: true},
	{ExternalID: "finance_data_entry", CategoryExternalID: "finance_banking_support", Name: "数据录入员", OrderNum: 20, IsActive: true},
	{ExternalID: "finance_doc_specialist", CategoryExternalID: "finance_banking_support", Name: "单证处理专员", OrderNum: 30, IsActive: true},
	{ExternalID: "finance_data_analyst", CategoryExternalID: "finance_banking_support", Name: "金融数据分析师", OrderNum: 40, IsActive: true},
	{ExternalID: "finance_funds_settlement", CategoryExternalID: "finance_banking_support", Name: "资金结算专员", OrderNum: 50, IsActive: true},

	{ExternalID: "finance_trust_manager", CategoryExternalID: "finance_trust_futures", Name: "信托经理", OrderNum: 10, IsActive: true},
	{ExternalID: "finance_trader_operator", CategoryExternalID: "finance_trust_futures", Name: "操盘手", OrderNum: 20, IsActive: true},
	{ExternalID: "finance_futures_analyst", CategoryExternalID: "finance_trust_futures", Name: "期货分析师", OrderNum: 30, IsActive: true},
	{ExternalID: "finance_invest_strategy", CategoryExternalID: "finance_trust_futures", Name: "投资策略师", OrderNum: 40, IsActive: true},

	{ExternalID: "finance_branch_manager", CategoryExternalID: "finance_bank_management", Name: "支行行长", OrderNum: 10, IsActive: true},
	{ExternalID: "finance_branch_vice_manager", CategoryExternalID: "finance_bank_management", Name: "支行副行长", OrderNum: 20, IsActive: true},
	{ExternalID: "finance_subbranch_manager", CategoryExternalID: "finance_bank_management", Name: "分行行长", OrderNum: 30, IsActive: true},
	{ExternalID: "finance_cfo", CategoryExternalID: "finance_bank_management", Name: "首席财务官 (CFO)", OrderNum: 40, IsActive: true},
	{ExternalID: "finance_invest_director", CategoryExternalID: "finance_bank_management", Name: "投资总监", OrderNum: 50, IsActive: true},

	{ExternalID: "finance_mt", CategoryExternalID: "finance_intern", Name: "银行管理培训生", OrderNum: 10, IsActive: true},
	{ExternalID: "finance_securities_intern", CategoryExternalID: "finance_intern", Name: "证券实习生", OrderNum: 20, IsActive: true},
	{ExternalID: "finance_insurance_intern", CategoryExternalID: "finance_intern", Name: "保险实习生", OrderNum: 30, IsActive: true},
	{ExternalID: "finance_ib_intern", CategoryExternalID: "finance_intern", Name: "投资银行实习生", OrderNum: 40, IsActive: true},
	{ExternalID: "finance_intern_assistant", CategoryExternalID: "finance_intern", Name: "实习助理", OrderNum: 50, IsActive: true},
	{ExternalID: "finance_teller_intern", CategoryExternalID: "finance_intern", Name: "柜员实习生", OrderNum: 60, IsActive: true},
	{ExternalID: "finance_bank_intern", CategoryExternalID: "finance_intern", Name: "银行实习生", OrderNum: 70, IsActive: true},

	{ExternalID: "edu_chinese_teacher", CategoryExternalID: "education_teacher", Name: "语文教师", OrderNum: 10, IsActive: true},
	{ExternalID: "edu_english_teacher", CategoryExternalID: "education_teacher", Name: "英语教师", OrderNum: 20, IsActive: true},
	{ExternalID: "edu_art_teacher", CategoryExternalID: "education_teacher", Name: "美术老师", OrderNum: 30, IsActive: true},
	{ExternalID: "edu_math_teacher", CategoryExternalID: "education_teacher", Name: "数学教师", OrderNum: 40, IsActive: true},
	{ExternalID: "edu_kindergarten_teacher", CategoryExternalID: "education_teacher", Name: "幼儿教师", OrderNum: 50, IsActive: true},
	{ExternalID: "edu_pe_teacher", CategoryExternalID: "education_teacher", Name: "体育教师", OrderNum: 60, IsActive: true},
	{ExternalID: "edu_music_teacher", CategoryExternalID: "education_teacher", Name: "音乐老师", OrderNum: 70, IsActive: true},
	{ExternalID: "edu_biology_teacher", CategoryExternalID: "education_teacher", Name: "生物教师", OrderNum: 80, IsActive: true},
	{ExternalID: "edu_dance_teacher", CategoryExternalID: "education_teacher", Name: "舞蹈老师", OrderNum: 90, IsActive: true},
	{ExternalID: "edu_piano_teacher", CategoryExternalID: "education_teacher", Name: "钢琴教师", OrderNum: 100, IsActive: true},
	{ExternalID: "edu_calligraphy_teacher", CategoryExternalID: "education_teacher", Name: "书法教师", OrderNum: 110, IsActive: true},
	{ExternalID: "edu_chemistry_teacher", CategoryExternalID: "education_teacher", Name: "化学老师", OrderNum: 120, IsActive: true},
	{ExternalID: "edu_physics_teacher", CategoryExternalID: "education_teacher", Name: "物理老师", OrderNum: 130, IsActive: true},
	{ExternalID: "edu_history_teacher", CategoryExternalID: "education_teacher", Name: "历史老师", OrderNum: 140, IsActive: true},
	{ExternalID: "edu_politics_teacher", CategoryExternalID: "education_teacher", Name: "政治老师", OrderNum: 150, IsActive: true},
	{ExternalID: "edu_geography_teacher", CategoryExternalID: "education_teacher", Name: "地理老师", OrderNum: 160, IsActive: true},
	{ExternalID: "edu_tcsol_teacher", CategoryExternalID: "education_teacher", Name: "对外汉语教师", OrderNum: 170, IsActive: true},
	{ExternalID: "edu_tutor", CategoryExternalID: "education_teacher", Name: "家教", OrderNum: 180, IsActive: true},
	{ExternalID: "edu_university_teacher", CategoryExternalID: "education_teacher", Name: "大学教师", OrderNum: 190, IsActive: true},
	{ExternalID: "edu_tutoring_teacher", CategoryExternalID: "education_teacher", Name: "辅导老师", OrderNum: 200, IsActive: true},

	{ExternalID: "edu_principal", CategoryExternalID: "education_teaching_admin", Name: "校长", OrderNum: 10, IsActive: true},
	{ExternalID: "edu_vice_principal", CategoryExternalID: "education_teaching_admin", Name: "副校长", OrderNum: 20, IsActive: true},
	{ExternalID: "edu_kindergarten_principal", CategoryExternalID: "education_teaching_admin", Name: "园长", OrderNum: 30, IsActive: true},
	{ExternalID: "edu_academic_director", CategoryExternalID: "education_teaching_admin", Name: "教务主任", OrderNum: 40, IsActive: true},
	{ExternalID: "edu_research_leader", CategoryExternalID: "education_teaching_admin", Name: "教研组长", OrderNum: 50, IsActive: true},
	{ExternalID: "edu_teaching_supervisor", CategoryExternalID: "education_teaching_admin", Name: "教学主管", OrderNum: 60, IsActive: true},
	{ExternalID: "edu_academic_specialist", CategoryExternalID: "education_teaching_admin", Name: "教务专员", OrderNum: 70, IsActive: true},
	{ExternalID: "edu_researcher", CategoryExternalID: "education_teaching_admin", Name: "教研组长", OrderNum: 80, IsActive: true},
	{ExternalID: "edu_curriculum_dev", CategoryExternalID: "education_teaching_admin", Name: "课程开发/设计", OrderNum: 90, IsActive: true},
	{ExternalID: "edu_academic_assistant", CategoryExternalID: "education_teaching_admin", Name: "教务助理", OrderNum: 100, IsActive: true},

	{ExternalID: "edu_dorm_teacher", CategoryExternalID: "education_student_services", Name: "宿管老师", OrderNum: 10, IsActive: true},
	{ExternalID: "edu_admissions", CategoryExternalID: "education_student_services", Name: "招生顾问", OrderNum: 20, IsActive: true},
	{ExternalID: "edu_counselor", CategoryExternalID: "education_student_services", Name: "辅导员", OrderNum: 30, IsActive: true},
	{ExternalID: "edu_ta", CategoryExternalID: "education_student_services", Name: "助教/教学助理", OrderNum: 40, IsActive: true},
	{ExternalID: "edu_life_teacher", CategoryExternalID: "education_student_services", Name: "生活老师", OrderNum: 50, IsActive: true},
	{ExternalID: "edu_study_abroad", CategoryExternalID: "education_student_services", Name: "留学顾问", OrderNum: 60, IsActive: true},
	{ExternalID: "edu_homeroom_teacher", CategoryExternalID: "education_student_services", Name: "班主任", OrderNum: 70, IsActive: true},
	{ExternalID: "edu_student_manager", CategoryExternalID: "education_student_services", Name: "学管师", OrderNum: 80, IsActive: true},

	{ExternalID: "edu_trainer", CategoryExternalID: "education_training_lecturer", Name: "培训师", OrderNum: 10, IsActive: true},
	{ExternalID: "edu_lecturer", CategoryExternalID: "education_training_lecturer", Name: "培训讲师", OrderNum: 20, IsActive: true},
	{ExternalID: "edu_english_trainer", CategoryExternalID: "education_training_lecturer", Name: "英语培训老师", OrderNum: 30, IsActive: true},
	{ExternalID: "edu_corporate_trainer", CategoryExternalID: "education_training_lecturer", Name: "企业培训师", OrderNum: 40, IsActive: true},
	{ExternalID: "edu_vocational_trainer", CategoryExternalID: "education_training_lecturer", Name: "职业培训师", OrderNum: 50, IsActive: true},
	{ExternalID: "edu_speaker", CategoryExternalID: "education_training_lecturer", Name: "讲师", OrderNum: 60, IsActive: true},
	{ExternalID: "edu_internal_trainer", CategoryExternalID: "education_training_lecturer", Name: "内训师", OrderNum: 70, IsActive: true},

	{ExternalID: "edu_training_specialist", CategoryExternalID: "education_training_management", Name: "培训专员", OrderNum: 10, IsActive: true},
	{ExternalID: "edu_training_supervisor", CategoryExternalID: "education_training_management", Name: "培训主管", OrderNum: 20, IsActive: true},
	{ExternalID: "edu_training_assistant", CategoryExternalID: "education_training_management", Name: "培训助理", OrderNum: 30, IsActive: true},
	{ExternalID: "edu_training_director", CategoryExternalID: "education_training_management", Name: "培训总监", OrderNum: 40, IsActive: true},
	{ExternalID: "edu_training_manager", CategoryExternalID: "education_training_management", Name: "培训经理", OrderNum: 50, IsActive: true},

	{ExternalID: "hc_clinical_doctor", CategoryExternalID: "healthcare_doctor", Name: "临床医生", OrderNum: 10, IsActive: true},
	{ExternalID: "hc_internal_doctor", CategoryExternalID: "healthcare_doctor", Name: "内科医生", OrderNum: 20, IsActive: true},
	{ExternalID: "hc_surgery_doctor", CategoryExternalID: "healthcare_doctor", Name: "外科医生", OrderNum: 30, IsActive: true},
	{ExternalID: "hc_obgyn_doctor", CategoryExternalID: "healthcare_doctor", Name: "妇产科医生", OrderNum: 40, IsActive: true},
	{ExternalID: "hc_pediatrics_doctor", CategoryExternalID: "healthcare_doctor", Name: "儿科医生", OrderNum: 50, IsActive: true},
	{ExternalID: "hc_orthopedics_doctor", CategoryExternalID: "healthcare_doctor", Name: "骨科医生", OrderNum: 60, IsActive: true},
	{ExternalID: "hc_anesthesiologist", CategoryExternalID: "healthcare_doctor", Name: "麻醉医生", OrderNum: 70, IsActive: true},
	{ExternalID: "hc_dentist", CategoryExternalID: "healthcare_doctor", Name: "口腔医生", OrderNum: 80, IsActive: true},
	{ExternalID: "hc_tcm", CategoryExternalID: "healthcare_doctor", Name: "中医师", OrderNum: 90, IsActive: true},
	{ExternalID: "hc_radiologist", CategoryExternalID: "healthcare_doctor", Name: "放射科医生", OrderNum: 100, IsActive: true},
	{ExternalID: "hc_gp", CategoryExternalID: "healthcare_doctor", Name: "全科医生", OrderNum: 110, IsActive: true},
	{ExternalID: "hc_specialist_doctor", CategoryExternalID: "healthcare_doctor", Name: "专科医生", OrderNum: 120, IsActive: true},
	{ExternalID: "hc_doctor_assistant", CategoryExternalID: "healthcare_doctor", Name: "医生助理", OrderNum: 130, IsActive: true},
	{ExternalID: "hc_resident", CategoryExternalID: "healthcare_doctor", Name: "住院医师", OrderNum: 140, IsActive: true},
	{ExternalID: "hc_dental_doctor", CategoryExternalID: "healthcare_doctor", Name: "牙科医生", OrderNum: 150, IsActive: true},

	{ExternalID: "hc_nurse", CategoryExternalID: "healthcare_nurse", Name: "护士", OrderNum: 10, IsActive: true},
	{ExternalID: "hc_clinical_nurse", CategoryExternalID: "healthcare_nurse", Name: "临床护士", OrderNum: 20, IsActive: true},
	{ExternalID: "hc_or_nurse", CategoryExternalID: "healthcare_nurse", Name: "手术室护士", OrderNum: 30, IsActive: true},
	{ExternalID: "hc_internal_nurse", CategoryExternalID: "healthcare_nurse", Name: "内科护士", OrderNum: 40, IsActive: true},
	{ExternalID: "hc_obgyn_nurse", CategoryExternalID: "healthcare_nurse", Name: "妇产科护士", OrderNum: 50, IsActive: true},
	{ExternalID: "hc_midwife", CategoryExternalID: "healthcare_nurse", Name: "助产士", OrderNum: 60, IsActive: true},

	{ExternalID: "hc_lab_tech", CategoryExternalID: "healthcare_medtech", Name: "医学检验技师", OrderNum: 10, IsActive: true},
	{ExternalID: "hc_imaging_tech", CategoryExternalID: "healthcare_medtech", Name: "医学影像技师", OrderNum: 20, IsActive: true},
	{ExternalID: "hc_rehab_therapist", CategoryExternalID: "healthcare_medtech", Name: "康复治疗师", OrderNum: 30, IsActive: true},
	{ExternalID: "hc_dietitian", CategoryExternalID: "healthcare_medtech", Name: "营养师", OrderNum: 40, IsActive: true},
	{ExternalID: "hc_health_manager", CategoryExternalID: "healthcare_medtech", Name: "健康管理师", OrderNum: 50, IsActive: true},
	{ExternalID: "hc_psych_consultant", CategoryExternalID: "healthcare_medtech", Name: "心理咨询师", OrderNum: 60, IsActive: true},
	{ExternalID: "hc_acupuncture", CategoryExternalID: "healthcare_medtech", Name: "针灸推拿", OrderNum: 70, IsActive: true},
	{ExternalID: "hc_lab_tester", CategoryExternalID: "healthcare_medtech", Name: "检验师", OrderNum: 80, IsActive: true},
	{ExternalID: "hc_ultrasound_doctor", CategoryExternalID: "healthcare_medtech", Name: "超声科医师", OrderNum: 90, IsActive: true},
	{ExternalID: "hc_pathologist", CategoryExternalID: "healthcare_medtech", Name: "病理科医师", OrderNum: 100, IsActive: true},

	{ExternalID: "hc_pharma_related", CategoryExternalID: "healthcare_pharma", Name: "药学相关", OrderNum: 10, IsActive: true},
	{ExternalID: "hc_drug_rd", CategoryExternalID: "healthcare_pharma", Name: "药物研发", OrderNum: 20, IsActive: true},
	{ExternalID: "hc_medicine_rd", CategoryExternalID: "healthcare_pharma", Name: "药品研发", OrderNum: 30, IsActive: true},
	{ExternalID: "hc_pharmacist", CategoryExternalID: "healthcare_pharma", Name: "药剂师/药师", OrderNum: 40, IsActive: true},
	{ExternalID: "hc_med_qc", CategoryExternalID: "healthcare_pharma", Name: "医药质检", OrderNum: 50, IsActive: true},
	{ExternalID: "hc_drug_registration", CategoryExternalID: "healthcare_pharma", Name: "药品注册", OrderNum: 60, IsActive: true},
	{ExternalID: "hc_cra", CategoryExternalID: "healthcare_pharma", Name: "临床监察员 (CRA)", OrderNum: 70, IsActive: true},
	{ExternalID: "hc_crc", CategoryExternalID: "healthcare_pharma", Name: "临床协调员 (CRC)", OrderNum: 80, IsActive: true},
	{ExternalID: "hc_drug_quality_mgmt", CategoryExternalID: "healthcare_pharma", Name: "药品质量管理", OrderNum: 90, IsActive: true},

	{ExternalID: "hc_device_sales", CategoryExternalID: "healthcare_devices", Name: "医疗器械销售", OrderNum: 10, IsActive: true},
	{ExternalID: "hc_after_sales", CategoryExternalID: "healthcare_devices", Name: "售后工程师", OrderNum: 20, IsActive: true},
	{ExternalID: "hc_device_inspector", CategoryExternalID: "healthcare_devices", Name: "检验员", OrderNum: 30, IsActive: true},
	{ExternalID: "hc_device_qc", CategoryExternalID: "healthcare_devices", Name: "质检员", OrderNum: 40, IsActive: true},

	{ExternalID: "hc_records_admin", CategoryExternalID: "healthcare_other", Name: "病案管理员", OrderNum: 10, IsActive: true},
	{ExternalID: "hc_registration_cashier", CategoryExternalID: "healthcare_other", Name: "挂号/收费员", OrderNum: 20, IsActive: true},
	{ExternalID: "hc_guide", CategoryExternalID: "healthcare_other", Name: "导医", OrderNum: 30, IsActive: true},
	{ExternalID: "hc_hospital_frontdesk", CategoryExternalID: "healthcare_other", Name: "医院前台", OrderNum: 40, IsActive: true},
	{ExternalID: "hc_inventory", CategoryExternalID: "healthcare_other", Name: "医疗库管", OrderNum: 50, IsActive: true},
	{ExternalID: "hc_security", CategoryExternalID: "healthcare_other", Name: "医院保安", OrderNum: 60, IsActive: true},
	{ExternalID: "hc_cleaner", CategoryExternalID: "healthcare_other", Name: "医院保洁", OrderNum: 70, IsActive: true},
	{ExternalID: "hc_caregiver", CategoryExternalID: "healthcare_other", Name: "医院护工", OrderNum: 80, IsActive: true},
	{ExternalID: "hc_companion", CategoryExternalID: "healthcare_other", Name: "医院陪护", OrderNum: 90, IsActive: true},
	{ExternalID: "hc_med_translator", CategoryExternalID: "healthcare_other", Name: "医疗翻译", OrderNum: 100, IsActive: true},
	{ExternalID: "hc_med_legal", CategoryExternalID: "healthcare_other", Name: "医疗法务", OrderNum: 110, IsActive: true},
	{ExternalID: "hc_med_consulting", CategoryExternalID: "healthcare_other", Name: "医疗咨询", OrderNum: 120, IsActive: true},
	{ExternalID: "hc_med_training", CategoryExternalID: "healthcare_other", Name: "医疗培训", OrderNum: 130, IsActive: true},
	{ExternalID: "hc_med_social_worker", CategoryExternalID: "healthcare_other", Name: "医疗社工", OrderNum: 140, IsActive: true},
	{ExternalID: "hc_med_rep", CategoryExternalID: "healthcare_other", Name: "医药代表", OrderNum: 150, IsActive: true},

	{ExternalID: "hc_intern_doctor", CategoryExternalID: "healthcare_intern", Name: "实习医生", OrderNum: 10, IsActive: true},
	{ExternalID: "hc_intern_nurse", CategoryExternalID: "healthcare_intern", Name: "实习护士", OrderNum: 20, IsActive: true},
	{ExternalID: "hc_intern_pharmacist", CategoryExternalID: "healthcare_intern", Name: "实习药师", OrderNum: 30, IsActive: true},
	{ExternalID: "hc_intern_tech", CategoryExternalID: "healthcare_intern", Name: "实习技师", OrderNum: 40, IsActive: true},
	{ExternalID: "hc_med_intern", CategoryExternalID: "healthcare_intern", Name: "医学实习生", OrderNum: 50, IsActive: true},
	{ExternalID: "hc_intern_med_rep", CategoryExternalID: "healthcare_intern", Name: "医药代表实习生", OrderNum: 60, IsActive: true},
	{ExternalID: "hc_clinical_research_intern", CategoryExternalID: "healthcare_intern", Name: "临床研究实习生", OrderNum: 70, IsActive: true},

	{ExternalID: "hc_hospital_director", CategoryExternalID: "healthcare_management", Name: "医院院长", OrderNum: 10, IsActive: true},
	{ExternalID: "hc_department_head", CategoryExternalID: "healthcare_management", Name: "科室主任", OrderNum: 20, IsActive: true},
	{ExternalID: "hc_head_nurse", CategoryExternalID: "healthcare_management", Name: "护士长", OrderNum: 30, IsActive: true},

	{ExternalID: "re_structural_engineer", CategoryExternalID: "realestate_design_planning", Name: "结构工程师", OrderNum: 10, IsActive: true},
	{ExternalID: "re_urban_planner", CategoryExternalID: "realestate_design_planning", Name: "城市规划师", OrderNum: 20, IsActive: true},
	{ExternalID: "re_landscape_designer", CategoryExternalID: "realestate_design_planning", Name: "园林设计师", OrderNum: 30, IsActive: true},
	{ExternalID: "re_planning_design", CategoryExternalID: "realestate_design_planning", Name: "规划设计", OrderNum: 40, IsActive: true},
	{ExternalID: "re_bim_engineer", CategoryExternalID: "realestate_design_planning", Name: "BIM工程师", OrderNum: 50, IsActive: true},
	{ExternalID: "re_landscape_constructor", CategoryExternalID: "realestate_design_planning", Name: "园林施工员", OrderNum: 60, IsActive: true},
	{ExternalID: "re_construction_engineer", CategoryExternalID: "realestate_design_planning", Name: "建筑工程师", OrderNum: 70, IsActive: true},

	{ExternalID: "re_landscape_architect", CategoryExternalID: "realestate_interior_landscape", Name: "景观设计师", OrderNum: 10, IsActive: true},
	{ExternalID: "re_home_designer", CategoryExternalID: "realestate_interior_landscape", Name: "家装设计师", OrderNum: 20, IsActive: true},

	{ExternalID: "re_budgeter", CategoryExternalID: "realestate_cost_budget", Name: "预算员", OrderNum: 10, IsActive: true},
	{ExternalID: "re_cost_estimator", CategoryExternalID: "realestate_cost_budget", Name: "造价员", OrderNum: 20, IsActive: true},
	{ExternalID: "re_cost_engineer", CategoryExternalID: "realestate_cost_budget", Name: "造价工程师", OrderNum: 30, IsActive: true},
	{ExternalID: "re_project_cost", CategoryExternalID: "realestate_cost_budget", Name: "工程造价", OrderNum: 40, IsActive: true},

	{ExternalID: "re_site_worker", CategoryExternalID: "realestate_construction_mgmt", Name: "施工员", OrderNum: 10, IsActive: true},
	{ExternalID: "re_civil_site_worker", CategoryExternalID: "realestate_construction_mgmt", Name: "土建施工员", OrderNum: 20, IsActive: true},
	{ExternalID: "re_surveyor", CategoryExternalID: "realestate_construction_mgmt", Name: "测量员", OrderNum: 30, IsActive: true},
	{ExternalID: "re_mapping_engineer", CategoryExternalID: "realestate_construction_mgmt", Name: "测绘工程师", OrderNum: 40, IsActive: true},
	{ExternalID: "re_supervisor", CategoryExternalID: "realestate_construction_mgmt", Name: "工程监理", OrderNum: 50, IsActive: true},
	{ExternalID: "re_supervision_engineer", CategoryExternalID: "realestate_construction_mgmt", Name: "监理工程师", OrderNum: 60, IsActive: true},
	{ExternalID: "re_project_admin", CategoryExternalID: "realestate_construction_mgmt", Name: "工程管理员", OrderNum: 70, IsActive: true},
	{ExternalID: "re_document_controller", CategoryExternalID: "realestate_construction_mgmt", Name: "资料员", OrderNum: 80, IsActive: true},
	{ExternalID: "re_archive_admin", CategoryExternalID: "realestate_construction_mgmt", Name: "档案管理员", OrderNum: 90, IsActive: true},
	{ExternalID: "re_engineering_manager", CategoryExternalID: "realestate_construction_mgmt", Name: "工程经理", OrderNum: 100, IsActive: true},
	{ExternalID: "re_civil_engineer", CategoryExternalID: "realestate_construction_mgmt", Name: "土木工程师", OrderNum: 110, IsActive: true},
	{ExternalID: "re_safety_officer", CategoryExternalID: "realestate_construction_mgmt", Name: "安全员", OrderNum: 120, IsActive: true},
	{ExternalID: "re_quality_inspector", CategoryExternalID: "realestate_construction_mgmt", Name: "工程质检员", OrderNum: 130, IsActive: true},
	{ExternalID: "re_project_engineer", CategoryExternalID: "realestate_construction_mgmt", Name: "项目工程师", OrderNum: 140, IsActive: true},
	{ExternalID: "re_civil_project_engineer", CategoryExternalID: "realestate_construction_mgmt", Name: "土建工程师", OrderNum: 150, IsActive: true},

	{ExternalID: "re_pm", CategoryExternalID: "realestate_project_mgmt", Name: "项目经理", OrderNum: 10, IsActive: true},
	{ExternalID: "re_pm_assistant", CategoryExternalID: "realestate_project_mgmt", Name: "项目助理", OrderNum: 20, IsActive: true},
	{ExternalID: "re_pm_specialist", CategoryExternalID: "realestate_project_mgmt", Name: "项目专员", OrderNum: 30, IsActive: true},
	{ExternalID: "re_pm_supervisor", CategoryExternalID: "realestate_project_mgmt", Name: "项目主管", OrderNum: 40, IsActive: true},
	{ExternalID: "re_pm_director", CategoryExternalID: "realestate_project_mgmt", Name: "项目总监", OrderNum: 50, IsActive: true},
	{ExternalID: "re_bid_specialist", CategoryExternalID: "realestate_project_mgmt", Name: "投标专员", OrderNum: 60, IsActive: true},

	{ExternalID: "re_sales", CategoryExternalID: "realestate_sales_planning", Name: "房地产销售", OrderNum: 10, IsActive: true},
	{ExternalID: "re_property_consultant", CategoryExternalID: "realestate_sales_planning", Name: "置业顾问", OrderNum: 20, IsActive: true},
	{ExternalID: "re_agent", CategoryExternalID: "realestate_sales_planning", Name: "房产经纪人", OrderNum: 30, IsActive: true},
	{ExternalID: "re_marketing_planner", CategoryExternalID: "realestate_sales_planning", Name: "房地产策划", OrderNum: 40, IsActive: true},
	{ExternalID: "re_leasing_manager", CategoryExternalID: "realestate_sales_planning", Name: "招商经理", OrderNum: 50, IsActive: true},
	{ExternalID: "re_channel_manager", CategoryExternalID: "realestate_sales_planning", Name: "渠道经理", OrderNum: 60, IsActive: true},
	{ExternalID: "re_other_roles", CategoryExternalID: "realestate_sales_planning", Name: "房产其他岗位", OrderNum: 70, IsActive: true},

	{ExternalID: "re_property_mgmt", CategoryExternalID: "realestate_property_mgmt", Name: "物业管理", OrderNum: 10, IsActive: true},
	{ExternalID: "re_property_manager", CategoryExternalID: "realestate_property_mgmt", Name: "物业经理", OrderNum: 20, IsActive: true},
	{ExternalID: "re_property_cs", CategoryExternalID: "realestate_property_mgmt", Name: "物业客服", OrderNum: 30, IsActive: true},
	{ExternalID: "re_property_steward", CategoryExternalID: "realestate_property_mgmt", Name: "物业管家", OrderNum: 40, IsActive: true},

	{ExternalID: "mfg_mech_engineer", CategoryExternalID: "mfg_mechanical", Name: "机械工程师", OrderNum: 10, IsActive: true},
	{ExternalID: "mfg_mech_design", CategoryExternalID: "mfg_mechanical", Name: "机械设计工程师", OrderNum: 20, IsActive: true},
	{ExternalID: "mfg_mechatronics", CategoryExternalID: "mfg_mechanical", Name: "机电工程师", OrderNum: 30, IsActive: true},
	{ExternalID: "mfg_equipment_engineer", CategoryExternalID: "mfg_mechanical", Name: "设备工程师", OrderNum: 40, IsActive: true},
	{ExternalID: "mfg_mechanical_manufacturing", CategoryExternalID: "mfg_mechanical", Name: "机械制造", OrderNum: 50, IsActive: true},
	{ExternalID: "mfg_mechanical_maintenance", CategoryExternalID: "mfg_mechanical", Name: "机械维修", OrderNum: 60, IsActive: true},
	{ExternalID: "mfg_hydraulics", CategoryExternalID: "mfg_mechanical", Name: "液压工程师", OrderNum: 70, IsActive: true},
	{ExternalID: "mfg_nc_programmer", CategoryExternalID: "mfg_mechanical", Name: "数控编程", OrderNum: 80, IsActive: true},
	{ExternalID: "mfg_other_engineers", CategoryExternalID: "mfg_mechanical", Name: "其他工程师岗位", OrderNum: 90, IsActive: true},

	{ExternalID: "mfg_elec_engineer", CategoryExternalID: "mfg_electrical", Name: "电子工程师", OrderNum: 10, IsActive: true},
	{ExternalID: "mfg_electrical_engineer", CategoryExternalID: "mfg_electrical", Name: "电气工程师", OrderNum: 20, IsActive: true},
	{ExternalID: "mfg_hw_engineer", CategoryExternalID: "mfg_electrical", Name: "硬件工程师", OrderNum: 30, IsActive: true},
	{ExternalID: "mfg_industrial_automation", CategoryExternalID: "mfg_electrical", Name: "电气自动化工程师", OrderNum: 40, IsActive: true},
	{ExternalID: "mfg_embedded", CategoryExternalID: "mfg_electrical", Name: "嵌入式工程师", OrderNum: 50, IsActive: true},
	{ExternalID: "mfg_automation", CategoryExternalID: "mfg_electrical", Name: "自动化工程师", OrderNum: 60, IsActive: true},
	{ExternalID: "mfg_semiconductor_tech", CategoryExternalID: "mfg_electrical", Name: "半导体技术员", OrderNum: 70, IsActive: true},
	{ExternalID: "mfg_circuit_design", CategoryExternalID: "mfg_electrical", Name: "电路设计", OrderNum: 80, IsActive: true},

	{ExternalID: "mfg_auto_engineer", CategoryExternalID: "mfg_auto_transport", Name: "汽车工程师", OrderNum: 10, IsActive: true},
	{ExternalID: "mfg_vehicle_engineer", CategoryExternalID: "mfg_auto_transport", Name: "车辆工程师", OrderNum: 20, IsActive: true},
	{ExternalID: "mfg_vehicle_design", CategoryExternalID: "mfg_auto_transport", Name: "汽车设计", OrderNum: 30, IsActive: true},
	{ExternalID: "mfg_powertrain", CategoryExternalID: "mfg_auto_transport", Name: "动力总成工程师", OrderNum: 40, IsActive: true},

	{ExternalID: "mfg_process_engineer", CategoryExternalID: "mfg_process_mold", Name: "工艺工程师", OrderNum: 10, IsActive: true},
	{ExternalID: "mfg_mould_engineer", CategoryExternalID: "mfg_process_mold", Name: "模具工程师", OrderNum: 20, IsActive: true},
	{ExternalID: "mfg_welding_engineer", CategoryExternalID: "mfg_process_mold", Name: "焊接工程师", OrderNum: 30, IsActive: true},
	{ExternalID: "mfg_mould_design", CategoryExternalID: "mfg_process_mold", Name: "模具设计师", OrderNum: 40, IsActive: true},
	{ExternalID: "mfg_stamping_process", CategoryExternalID: "mfg_process_mold", Name: "冲压工艺师/模具设计师", OrderNum: 50, IsActive: true},

	{ExternalID: "mfg_prod_management", CategoryExternalID: "mfg_prod_equip", Name: "生产管理", OrderNum: 10, IsActive: true},
	{ExternalID: "mfg_prod_supervisor", CategoryExternalID: "mfg_prod_equip", Name: "生产主管", OrderNum: 20, IsActive: true},
	{ExternalID: "mfg_production_manager", CategoryExternalID: "mfg_prod_equip", Name: "生产经理", OrderNum: 30, IsActive: true},
	{ExternalID: "mfg_line_leader", CategoryExternalID: "mfg_prod_equip", Name: "车间主任", OrderNum: 40, IsActive: true},
	{ExternalID: "mfg_equipment_maintenance", CategoryExternalID: "mfg_prod_equip", Name: "设备维护", OrderNum: 50, IsActive: true},
	{ExternalID: "mfg_equipment_manager", CategoryExternalID: "mfg_prod_equip", Name: "设备管理", OrderNum: 60, IsActive: true},
	{ExternalID: "mfg_shift_leader", CategoryExternalID: "mfg_prod_equip", Name: "生产班长", OrderNum: 70, IsActive: true},
	{ExternalID: "mfg_team_leader", CategoryExternalID: "mfg_prod_equip", Name: "工段长", OrderNum: 80, IsActive: true},
	{ExternalID: "mfg_group_leader", CategoryExternalID: "mfg_prod_equip", Name: "班组长", OrderNum: 90, IsActive: true},
	{ExternalID: "mfg_factory_manager", CategoryExternalID: "mfg_prod_equip", Name: "厂长", OrderNum: 100, IsActive: true},
	{ExternalID: "mfg_production_planner", CategoryExternalID: "mfg_prod_equip", Name: "生产计划员", OrderNum: 110, IsActive: true},

	{ExternalID: "mfg_quality_engineer", CategoryExternalID: "mfg_quality", Name: "质量工程师", OrderNum: 10, IsActive: true},
	{ExternalID: "mfg_quality_specialist", CategoryExternalID: "mfg_quality", Name: "品质工程师", OrderNum: 20, IsActive: true},
	{ExternalID: "mfg_qc", CategoryExternalID: "mfg_quality", Name: "QC", OrderNum: 30, IsActive: true},
	{ExternalID: "mfg_quality_manager", CategoryExternalID: "mfg_quality", Name: "质量管理工程师", OrderNum: 40, IsActive: true},
	{ExternalID: "mfg_quality_inspector", CategoryExternalID: "mfg_quality", Name: "质量检测员", OrderNum: 50, IsActive: true},

	{ExternalID: "mfg_rd_engineer", CategoryExternalID: "mfg_rd_design", Name: "研发工程师", OrderNum: 10, IsActive: true},
	{ExternalID: "mfg_tech_engineer", CategoryExternalID: "mfg_rd_design", Name: "技术员", OrderNum: 20, IsActive: true},
	{ExternalID: "mfg_design_engineer", CategoryExternalID: "mfg_rd_design", Name: "设计工程师", OrderNum: 30, IsActive: true},
	{ExternalID: "mfg_product_designer", CategoryExternalID: "mfg_rd_design", Name: "产品设计师", OrderNum: 40, IsActive: true},

	{ExternalID: "mfg_operator", CategoryExternalID: "mfg_worker", Name: "操作工", OrderNum: 10, IsActive: true},
	{ExternalID: "mfg_general_worker", CategoryExternalID: "mfg_worker", Name: "普工", OrderNum: 20, IsActive: true},
	{ExternalID: "mfg_assembler", CategoryExternalID: "mfg_worker", Name: "装配工", OrderNum: 30, IsActive: true},
	{ExternalID: "mfg_technician", CategoryExternalID: "mfg_worker", Name: "技术工人", OrderNum: 40, IsActive: true},
	{ExternalID: "mfg_fitter", CategoryExternalID: "mfg_worker", Name: "钳工", OrderNum: 50, IsActive: true},
	{ExternalID: "mfg_welder", CategoryExternalID: "mfg_worker", Name: "焊工", OrderNum: 60, IsActive: true},
	{ExternalID: "mfg_driller", CategoryExternalID: "mfg_worker", Name: "钻工", OrderNum: 70, IsActive: true},
	{ExternalID: "mfg_turner", CategoryExternalID: "mfg_worker", Name: "车工", OrderNum: 80, IsActive: true},
	{ExternalID: "mfg_miller", CategoryExternalID: "mfg_worker", Name: "铣工", OrderNum: 90, IsActive: true},
	{ExternalID: "mfg_cnc_operator", CategoryExternalID: "mfg_worker", Name: "CNC操作工", OrderNum: 100, IsActive: true},
	{ExternalID: "mfg_machine_operator", CategoryExternalID: "mfg_worker", Name: "机床操作工", OrderNum: 110, IsActive: true},
	{ExternalID: "mfg_packer", CategoryExternalID: "mfg_worker", Name: "包装工", OrderNum: 120, IsActive: true},

	{ExternalID: "logi_captain", CategoryExternalID: "logi_transport_service", Name: "机长", OrderNum: 10, IsActive: true},
	{ExternalID: "logi_pilot", CategoryExternalID: "logi_transport_service", Name: "飞行员", OrderNum: 20, IsActive: true},
	{ExternalID: "logi_air_security", CategoryExternalID: "logi_transport_service", Name: "空中安全员", OrderNum: 30, IsActive: true},
	{ExternalID: "logi_flight_attendant", CategoryExternalID: "logi_transport_service", Name: "空姐", OrderNum: 40, IsActive: true},
	{ExternalID: "logi_male_attendant", CategoryExternalID: "logi_transport_service", Name: "空少", OrderNum: 50, IsActive: true},
	{ExternalID: "logi_ground", CategoryExternalID: "logi_transport_service", Name: "地勤", OrderNum: 60, IsActive: true},
	{ExternalID: "logi_tickets", CategoryExternalID: "logi_transport_service", Name: "票务员", OrderNum: 70, IsActive: true},
	{ExternalID: "logi_security", CategoryExternalID: "logi_transport_service", Name: "安检员", OrderNum: 80, IsActive: true},
	{ExternalID: "logi_metro_staff", CategoryExternalID: "logi_transport_service", Name: "地铁站务员", OrderNum: 90, IsActive: true},
	{ExternalID: "logi_metro_security", CategoryExternalID: "logi_transport_service", Name: "地铁安检员", OrderNum: 100, IsActive: true},
	{ExternalID: "logi_metro_driver", CategoryExternalID: "logi_transport_service", Name: "地铁驾驶员", OrderNum: 110, IsActive: true},
	{ExternalID: "logi_highspeed_attendant", CategoryExternalID: "logi_transport_service", Name: "高铁乘务", OrderNum: 120, IsActive: true},
	{ExternalID: "logi_call_center", CategoryExternalID: "logi_transport_service", Name: "话务员", OrderNum: 130, IsActive: true},
	{ExternalID: "logi_flight_dispatcher", CategoryExternalID: "logi_transport_service", Name: "签派员", OrderNum: 140, IsActive: true},

	{ExternalID: "logi_driver", CategoryExternalID: "logi_delivery", Name: "司机", OrderNum: 10, IsActive: true},
	{ExternalID: "logi_dispatch_manager", CategoryExternalID: "logi_delivery", Name: "配送经理", OrderNum: 20, IsActive: true},
	{ExternalID: "logi_transport_admin", CategoryExternalID: "logi_delivery", Name: "运输主管", OrderNum: 30, IsActive: true},
	{ExternalID: "logi_fleet_manager", CategoryExternalID: "logi_delivery", Name: "车队管理", OrderNum: 40, IsActive: true},
	{ExternalID: "logi_scheduler", CategoryExternalID: "logi_delivery", Name: "调度员", OrderNum: 50, IsActive: true},

	{ExternalID: "logi_warehouse_admin", CategoryExternalID: "logi_warehouse", Name: "仓库管理员/库管", OrderNum: 10, IsActive: true},
	{ExternalID: "logi_storekeeper", CategoryExternalID: "logi_warehouse", Name: "仓管员", OrderNum: 20, IsActive: true},
	{ExternalID: "logi_warehouse_manager", CategoryExternalID: "logi_warehouse", Name: "仓储管理", OrderNum: 30, IsActive: true},
	{ExternalID: "logi_loader", CategoryExternalID: "logi_warehouse", Name: "装卸工", OrderNum: 40, IsActive: true},
	{ExternalID: "logi_packer", CategoryExternalID: "logi_warehouse", Name: "包装员", OrderNum: 50, IsActive: true},
	{ExternalID: "logi_warehouse_specialist", CategoryExternalID: "logi_warehouse", Name: "仓储专员", OrderNum: 60, IsActive: true},
	{ExternalID: "logi_warehouse_supervisor", CategoryExternalID: "logi_warehouse", Name: "仓储主管", OrderNum: 70, IsActive: true},
	{ExternalID: "logi_warehouse_manager_role", CategoryExternalID: "logi_warehouse", Name: "仓储经理", OrderNum: 80, IsActive: true},
	{ExternalID: "logi_warehouse_reserve", CategoryExternalID: "logi_warehouse", Name: "储备干部", OrderNum: 90, IsActive: true},
	{ExternalID: "logi_warehouse_superintendent", CategoryExternalID: "logi_warehouse", Name: "储备主管", OrderNum: 100, IsActive: true},

	{ExternalID: "logi_ops_specialist", CategoryExternalID: "logi_ops", Name: "物流专员/助理", OrderNum: 10, IsActive: true},
	{ExternalID: "logi_ops_supervisor", CategoryExternalID: "logi_ops", Name: "物流主管", OrderNum: 20, IsActive: true},
	{ExternalID: "logi_ops_manager", CategoryExternalID: "logi_ops", Name: "物流经理", OrderNum: 30, IsActive: true},
	{ExternalID: "logi_order_follower", CategoryExternalID: "logi_ops", Name: "跟单员", OrderNum: 40, IsActive: true},
	{ExternalID: "logi_courier", CategoryExternalID: "logi_ops", Name: "物流员", OrderNum: 50, IsActive: true},
	{ExternalID: "logi_express", CategoryExternalID: "logi_ops", Name: "快递员", OrderNum: 60, IsActive: true},
	{ExternalID: "logi_delivery_staff", CategoryExternalID: "logi_ops", Name: "配送员", OrderNum: 70, IsActive: true},

	{ExternalID: "logi_supply_specialist", CategoryExternalID: "logi_supply", Name: "供应链专员", OrderNum: 10, IsActive: true},
	{ExternalID: "logi_supply_manager", CategoryExternalID: "logi_supply", Name: "供应链管理", OrderNum: 20, IsActive: true},
	{ExternalID: "logi_supply_director", CategoryExternalID: "logi_supply", Name: "供应链经理", OrderNum: 30, IsActive: true},
	{ExternalID: "logi_supply_supervisor", CategoryExternalID: "logi_supply", Name: "供应链总监", OrderNum: 40, IsActive: true},
	{ExternalID: "logi_supply_intern", CategoryExternalID: "logi_supply", Name: "供应链实习生", OrderNum: 50, IsActive: true},

	{ExternalID: "svc_supermarket_manager", CategoryExternalID: "svc_food", Name: "超市店长", OrderNum: 10, IsActive: true},
	{ExternalID: "svc_cashier_staff", CategoryExternalID: "svc_food", Name: "收银员", OrderNum: 20, IsActive: true},
	{ExternalID: "svc_waiter", CategoryExternalID: "svc_food", Name: "服务员", OrderNum: 30, IsActive: true},
	{ExternalID: "svc_barista", CategoryExternalID: "svc_food", Name: "咖啡师", OrderNum: 40, IsActive: true},
	{ExternalID: "svc_chef", CategoryExternalID: "svc_food", Name: "厨师", OrderNum: 50, IsActive: true},
	{ExternalID: "svc_baker", CategoryExternalID: "svc_food", Name: "面包师", OrderNum: 60, IsActive: true},
	{ExternalID: "svc_western_chef", CategoryExternalID: "svc_food", Name: "西点师", OrderNum: 70, IsActive: true},
	{ExternalID: "svc_food_manager", CategoryExternalID: "svc_food", Name: "餐饮管理", OrderNum: 80, IsActive: true},
	{ExternalID: "svc_food_store_manager", CategoryExternalID: "svc_food", Name: "餐饮店长", OrderNum: 90, IsActive: true},
	{ExternalID: "svc_back_kitchen", CategoryExternalID: "svc_food", Name: "后厨", OrderNum: 100, IsActive: true},
	{ExternalID: "svc_restaurant_manager", CategoryExternalID: "svc_food", Name: "餐厅经理", OrderNum: 110, IsActive: true},
	{ExternalID: "svc_pastry_chef", CategoryExternalID: "svc_food", Name: "面点师", OrderNum: 120, IsActive: true},
	{ExternalID: "svc_head_chef", CategoryExternalID: "svc_food", Name: "厨师长", OrderNum: 130, IsActive: true},
	{ExternalID: "svc_baking", CategoryExternalID: "svc_food", Name: "烘焙师", OrderNum: 140, IsActive: true},
	{ExternalID: "svc_bartender", CategoryExternalID: "svc_food", Name: "调酒师", OrderNum: 150, IsActive: true},
	{ExternalID: "svc_tea_artist", CategoryExternalID: "svc_food", Name: "茶艺师", OrderNum: 160, IsActive: true},
	{ExternalID: "svc_food_director", CategoryExternalID: "svc_food", Name: "餐饮总监", OrderNum: 170, IsActive: true},

	{ExternalID: "svc_tour_guide", CategoryExternalID: "svc_hotel_travel", Name: "导游", OrderNum: 10, IsActive: true},
	{ExternalID: "svc_travel_consultant", CategoryExternalID: "svc_hotel_travel", Name: "旅游顾问", OrderNum: 20, IsActive: true},
	{ExternalID: "svc_planner", CategoryExternalID: "svc_hotel_travel", Name: "计调", OrderNum: 30, IsActive: true},
	{ExternalID: "svc_hotel_manager", CategoryExternalID: "svc_hotel_travel", Name: "酒店管理", OrderNum: 40, IsActive: true},
	{ExternalID: "svc_frontdesk", CategoryExternalID: "svc_hotel_travel", Name: "酒店前台", OrderNum: 50, IsActive: true},
	{ExternalID: "svc_guest_service", CategoryExternalID: "svc_hotel_travel", Name: "客房服务", OrderNum: 60, IsActive: true},
	{ExternalID: "svc_porter", CategoryExternalID: "svc_hotel_travel", Name: "行李员", OrderNum: 70, IsActive: true},
	{ExternalID: "svc_travel_interpreter", CategoryExternalID: "svc_hotel_travel", Name: "景区讲解员", OrderNum: 80, IsActive: true},
	{ExternalID: "svc_travel_custom", CategoryExternalID: "svc_hotel_travel", Name: "旅游定制师", OrderNum: 90, IsActive: true},
	{ExternalID: "svc_greeter", CategoryExternalID: "svc_hotel_travel", Name: "前厅经理", OrderNum: 100, IsActive: true},
	{ExternalID: "svc_concierge", CategoryExternalID: "svc_hotel_travel", Name: "礼宾相关岗位", OrderNum: 110, IsActive: true},
	{ExternalID: "svc_hotel_sales", CategoryExternalID: "svc_hotel_travel", Name: "酒店销售", OrderNum: 120, IsActive: true},
	{ExternalID: "svc_banquets", CategoryExternalID: "svc_hotel_travel", Name: "宴会服务", OrderNum: 130, IsActive: true},

	{ExternalID: "svc_cleaner", CategoryExternalID: "svc_personal", Name: "保洁", OrderNum: 10, IsActive: true},
	{ExternalID: "svc_guard", CategoryExternalID: "svc_personal", Name: "保安", OrderNum: 20, IsActive: true},
	{ExternalID: "svc_housekeeping", CategoryExternalID: "svc_personal", Name: "家政", OrderNum: 30, IsActive: true},
	{ExternalID: "svc_babysitter", CategoryExternalID: "svc_personal", Name: "保姆", OrderNum: 40, IsActive: true},
	{ExternalID: "svc_maternity_matron", CategoryExternalID: "svc_personal", Name: "月嫂", OrderNum: 50, IsActive: true},
	{ExternalID: "svc_housekeeper", CategoryExternalID: "svc_personal", Name: "管家", OrderNum: 60, IsActive: true},
	{ExternalID: "svc_pet_beautician", CategoryExternalID: "svc_personal", Name: "宠物美容师", OrderNum: 70, IsActive: true},
	{ExternalID: "svc_pet_doctor", CategoryExternalID: "svc_personal", Name: "宠物医生", OrderNum: 80, IsActive: true},
	{ExternalID: "svc_fashion_buyer", CategoryExternalID: "svc_personal", Name: "服装买手", OrderNum: 90, IsActive: true},
	{ExternalID: "svc_stylist", CategoryExternalID: "svc_personal", Name: "服装搭配师", OrderNum: 100, IsActive: true},

	{ExternalID: "svc_fitness_coach", CategoryExternalID: "svc_fitness", Name: "健身教练", OrderNum: 10, IsActive: true},
	{ExternalID: "svc_swimming_coach", CategoryExternalID: "svc_fitness", Name: "游泳教练", OrderNum: 20, IsActive: true},
	{ExternalID: "svc_yoga_coach", CategoryExternalID: "svc_fitness", Name: "瑜伽教练", OrderNum: 30, IsActive: true},
	{ExternalID: "svc_dance_teacher", CategoryExternalID: "svc_fitness", Name: "舞蹈老师", OrderNum: 40, IsActive: true},
	{ExternalID: "svc_basketball_coach", CategoryExternalID: "svc_fitness", Name: "篮球教练", OrderNum: 50, IsActive: true},
	{ExternalID: "svc_badminton_coach", CategoryExternalID: "svc_fitness", Name: "羽毛球教练", OrderNum: 60, IsActive: true},
	{ExternalID: "svc_taekwondo_coach", CategoryExternalID: "svc_fitness", Name: "跆拳道教练", OrderNum: 70, IsActive: true},
	{ExternalID: "svc_martial_teacher", CategoryExternalID: "svc_fitness", Name: "武术教练", OrderNum: 80, IsActive: true},
	{ExternalID: "svc_street_dance_teacher", CategoryExternalID: "svc_fitness", Name: "街舞老师", OrderNum: 90, IsActive: true},

	{ExternalID: "svc_beauty_consultant", CategoryExternalID: "svc_beauty", Name: "美容顾问", OrderNum: 10, IsActive: true},
	{ExternalID: "svc_store_manager", CategoryExternalID: "svc_beauty", Name: "美容店长", OrderNum: 20, IsActive: true},
	{ExternalID: "svc_beautician", CategoryExternalID: "svc_beauty", Name: "美容师", OrderNum: 30, IsActive: true},
	{ExternalID: "svc_makeup_artist", CategoryExternalID: "svc_beauty", Name: "化妆师", OrderNum: 40, IsActive: true},
	{ExternalID: "svc_manicurist", CategoryExternalID: "svc_beauty", Name: "美甲师", OrderNum: 50, IsActive: true},
	{ExternalID: "svc_hairdresser", CategoryExternalID: "svc_beauty", Name: "美发师", OrderNum: 60, IsActive: true},
	{ExternalID: "svc_masseur", CategoryExternalID: "svc_beauty", Name: "按摩师", OrderNum: 70, IsActive: true},
	{ExternalID: "svc_physiotherapist", CategoryExternalID: "svc_beauty", Name: "理疗师", OrderNum: 80, IsActive: true},
	{ExternalID: "svc_tattoo_artist", CategoryExternalID: "svc_beauty", Name: "纹绣师", OrderNum: 90, IsActive: true},

	{ExternalID: "svc_course_consultant", CategoryExternalID: "svc_consult_translate", Name: "课程顾问", OrderNum: 10, IsActive: true},
	{ExternalID: "svc_headhunter", CategoryExternalID: "svc_consult_translate", Name: "猎头顾问", OrderNum: 20, IsActive: true},
	{ExternalID: "svc_consultant", CategoryExternalID: "svc_consult_translate", Name: "咨询顾问", OrderNum: 30, IsActive: true},
	{ExternalID: "svc_legal_consultant", CategoryExternalID: "svc_consult_translate", Name: "法律顾问", OrderNum: 40, IsActive: true},
	{ExternalID: "svc_translator", CategoryExternalID: "svc_consult_translate", Name: "翻译", OrderNum: 50, IsActive: true},

	{ExternalID: "svc_auto_mechanic", CategoryExternalID: "svc_repair", Name: "汽车维修", OrderNum: 10, IsActive: true},
	{ExternalID: "svc_mobile_repair", CategoryExternalID: "svc_repair", Name: "机务维修", OrderNum: 20, IsActive: true},
	{ExternalID: "svc_electric_repair", CategoryExternalID: "svc_repair", Name: "维修电工", OrderNum: 30, IsActive: true},
	{ExternalID: "svc_mould_maintenance", CategoryExternalID: "svc_repair", Name: "模具维修", OrderNum: 40, IsActive: true},
	{ExternalID: "svc_device_repair", CategoryExternalID: "svc_repair", Name: "器械维修", OrderNum: 50, IsActive: true},
	{ExternalID: "svc_other_repair", CategoryExternalID: "svc_repair", Name: "其他维修服务岗", OrderNum: 60, IsActive: true},

	{ExternalID: "svc_store_clerk", CategoryExternalID: "svc_other", Name: "店长", OrderNum: 10, IsActive: true},
	{ExternalID: "svc_bid_staff", CategoryExternalID: "svc_other", Name: "项目招投标", OrderNum: 20, IsActive: true},
	{ExternalID: "svc_tickets_staff", CategoryExternalID: "svc_other", Name: "票务", OrderNum: 30, IsActive: true},
	{ExternalID: "svc_trustee", CategoryExternalID: "svc_other", Name: "托管", OrderNum: 40, IsActive: true},
	{ExternalID: "svc_fee_collector", CategoryExternalID: "svc_other", Name: "收费员", OrderNum: 50, IsActive: true},
	{ExternalID: "svc_transport_agent", CategoryExternalID: "svc_other", Name: "客运员", OrderNum: 60, IsActive: true},
	{ExternalID: "svc_uav_pilot", CategoryExternalID: "svc_other", Name: "无人机飞手", OrderNum: 70, IsActive: true},

	{ExternalID: "media_journalist", CategoryExternalID: "media_news_publishing", Name: "记者", OrderNum: 10, IsActive: true},
	{ExternalID: "media_editor", CategoryExternalID: "media_news_publishing", Name: "编辑", OrderNum: 20, IsActive: true},
	{ExternalID: "media_reporter", CategoryExternalID: "media_news_publishing", Name: "采编", OrderNum: 30, IsActive: true},
	{ExternalID: "media_chief_editor", CategoryExternalID: "media_news_publishing", Name: "主编/副主编", OrderNum: 40, IsActive: true},
	{ExternalID: "media_proofreader", CategoryExternalID: "media_news_publishing", Name: "校对", OrderNum: 50, IsActive: true},
	{ExternalID: "media_writer", CategoryExternalID: "media_news_publishing", Name: "撰稿", OrderNum: 60, IsActive: true},
	{ExternalID: "media_reviewer", CategoryExternalID: "media_news_publishing", Name: "审核", OrderNum: 70, IsActive: true},

	{ExternalID: "media_host", CategoryExternalID: "media_broadcast_tv", Name: "主持人", OrderNum: 10, IsActive: true},
	{ExternalID: "media_announcer", CategoryExternalID: "media_broadcast_tv", Name: "播音员", OrderNum: 20, IsActive: true},
	{ExternalID: "media_director_tv", CategoryExternalID: "media_broadcast_tv", Name: "编导", OrderNum: 30, IsActive: true},
	{ExternalID: "media_cameraman", CategoryExternalID: "media_broadcast_tv", Name: "摄像师", OrderNum: 40, IsActive: true},
	{ExternalID: "media_camera_assistant", CategoryExternalID: "media_broadcast_tv", Name: "摄影助理", OrderNum: 50, IsActive: true},
	{ExternalID: "media_postprod", CategoryExternalID: "media_broadcast_tv", Name: "后期制作", OrderNum: 60, IsActive: true},
	{ExternalID: "media_vfx", CategoryExternalID: "media_broadcast_tv", Name: "特效师", OrderNum: 70, IsActive: true},
	{ExternalID: "media_sound_engineer", CategoryExternalID: "media_broadcast_tv", Name: "音效师", OrderNum: 80, IsActive: true},
	{ExternalID: "media_switcher", CategoryExternalID: "media_broadcast_tv", Name: "导播", OrderNum: 90, IsActive: true},
	{ExternalID: "media_program_planner", CategoryExternalID: "media_broadcast_tv", Name: "节目策划", OrderNum: 100, IsActive: true},
	{ExternalID: "media_producer", CategoryExternalID: "media_broadcast_tv", Name: "制片人", OrderNum: 110, IsActive: true},
	{ExternalID: "media_film_making", CategoryExternalID: "media_broadcast_tv", Name: "影视制作", OrderNum: 120, IsActive: true},
	{ExternalID: "media_channel_specialist", CategoryExternalID: "media_broadcast_tv", Name: "渠道专员", OrderNum: 130, IsActive: true},

	{ExternalID: "media_director", CategoryExternalID: "media_film_performance", Name: "导演", OrderNum: 10, IsActive: true},
	{ExternalID: "media_screenwriter", CategoryExternalID: "media_film_performance", Name: "编剧", OrderNum: 20, IsActive: true},
	{ExternalID: "media_actor", CategoryExternalID: "media_film_performance", Name: "演员", OrderNum: 30, IsActive: true},
	{ExternalID: "media_script_supervisor", CategoryExternalID: "media_film_performance", Name: "场记", OrderNum: 40, IsActive: true},
	{ExternalID: "media_artist_assistant", CategoryExternalID: "media_film_performance", Name: "艺人助理", OrderNum: 50, IsActive: true},
	{ExternalID: "media_agent", CategoryExternalID: "media_film_performance", Name: "经纪人", OrderNum: 60, IsActive: true},
	{ExternalID: "media_model", CategoryExternalID: "media_film_performance", Name: "模特", OrderNum: 70, IsActive: true},
	{ExternalID: "media_stage_designer", CategoryExternalID: "media_film_performance", Name: "舞美设计", OrderNum: 80, IsActive: true},
	{ExternalID: "media_star_mapper", CategoryExternalID: "media_film_performance", Name: "星探", OrderNum: 90, IsActive: true},
	{ExternalID: "media_intern", CategoryExternalID: "media_film_performance", Name: "练习生", OrderNum: 100, IsActive: true},

	{ExternalID: "trade_business_staff", CategoryExternalID: "trade_foreign_trade", Name: "外贸业务员", OrderNum: 10, IsActive: true},
	{ExternalID: "trade_specialist", CategoryExternalID: "trade_foreign_trade", Name: "外贸专员", OrderNum: 20, IsActive: true},
	{ExternalID: "trade_doc_tracker", CategoryExternalID: "trade_foreign_trade", Name: "外贸跟单员", OrderNum: 30, IsActive: true},
	{ExternalID: "trade_assistant", CategoryExternalID: "trade_foreign_trade", Name: "外贸助理", OrderNum: 40, IsActive: true},
	{ExternalID: "trade_manager", CategoryExternalID: "trade_foreign_trade", Name: "外贸经理", OrderNum: 50, IsActive: true},
	{ExternalID: "trade_doc_staff", CategoryExternalID: "trade_foreign_trade", Name: "外贸单证员", OrderNum: 60, IsActive: true},
	{ExternalID: "trade_cs", CategoryExternalID: "trade_foreign_trade", Name: "外贸客服", OrderNum: 70, IsActive: true},
	{ExternalID: "trade_director", CategoryExternalID: "trade_foreign_trade", Name: "外贸总监", OrderNum: 80, IsActive: true},

	{ExternalID: "trade_customs_broker", CategoryExternalID: "trade_trade_support", Name: "报关员", OrderNum: 10, IsActive: true},
	{ExternalID: "trade_customs_supervisor", CategoryExternalID: "trade_trade_support", Name: "报关主管", OrderNum: 20, IsActive: true},
	{ExternalID: "trade_docs", CategoryExternalID: "trade_trade_support", Name: "单证员", OrderNum: 30, IsActive: true},

	{ExternalID: "trade_cbe", CategoryExternalID: "trade_crossborder_ecom", Name: "跨境电商", OrderNum: 10, IsActive: true},
	{ExternalID: "trade_cbe_specialist", CategoryExternalID: "trade_crossborder_ecom", Name: "跨境电商专员", OrderNum: 20, IsActive: true},
	{ExternalID: "trade_cbe_ops", CategoryExternalID: "trade_crossborder_ecom", Name: "跨境电商业务员", OrderNum: 30, IsActive: true},
	{ExternalID: "trade_cbe_assistant", CategoryExternalID: "trade_crossborder_ecom", Name: "跨境电商运营助理", OrderNum: 40, IsActive: true},
	{ExternalID: "trade_amazon_ops", CategoryExternalID: "trade_crossborder_ecom", Name: "Amazon运营", OrderNum: 50, IsActive: true},
	{ExternalID: "trade_amazon_sales", CategoryExternalID: "trade_crossborder_ecom", Name: "亚马逊销售", OrderNum: 60, IsActive: true},

	{ExternalID: "trade_translator", CategoryExternalID: "trade_translation_support", Name: "外贸翻译", OrderNum: 10, IsActive: true},
	{ExternalID: "trade_en_translator", CategoryExternalID: "trade_translation_support", Name: "英语翻译", OrderNum: 20, IsActive: true},
	{ExternalID: "trade_jp_translator", CategoryExternalID: "trade_translation_support", Name: "日语翻译", OrderNum: 30, IsActive: true},
	{ExternalID: "trade_kr_translator", CategoryExternalID: "trade_translation_support", Name: "韩语翻译", OrderNum: 40, IsActive: true},

	{ExternalID: "energy_power_engineer", CategoryExternalID: "energy_traditional", Name: "电力工程师", OrderNum: 10, IsActive: true},
	{ExternalID: "energy_new_energy_engineer", CategoryExternalID: "energy_traditional", Name: "新能源工程师", OrderNum: 20, IsActive: true},
	{ExternalID: "energy_thermal_engineer", CategoryExternalID: "energy_traditional", Name: "热能工程师", OrderNum: 30, IsActive: true},
	{ExternalID: "energy_oil_engineer", CategoryExternalID: "energy_traditional", Name: "石油工程师", OrderNum: 40, IsActive: true},
	{ExternalID: "energy_gas_engineer", CategoryExternalID: "energy_traditional", Name: "燃气工程师", OrderNum: 50, IsActive: true},
	{ExternalID: "energy_hvac_engineer", CategoryExternalID: "energy_traditional", Name: "暖通工程师", OrderNum: 60, IsActive: true},
	{ExternalID: "energy_engineer", CategoryExternalID: "energy_traditional", Name: "能源工程师", OrderNum: 70, IsActive: true},
	{ExternalID: "energy_pv_engineer", CategoryExternalID: "energy_traditional", Name: "光伏系统工程师", OrderNum: 80, IsActive: true},
	{ExternalID: "energy_wind_engineer", CategoryExternalID: "energy_traditional", Name: "风电工程师", OrderNum: 90, IsActive: true},

	{ExternalID: "env_engineer", CategoryExternalID: "energy_environment", Name: "环保工程师", OrderNum: 10, IsActive: true},
	{ExternalID: "env_environment_engineer", CategoryExternalID: "energy_environment", Name: "环境工程师", OrderNum: 20, IsActive: true},
	{ExternalID: "env_ehs_engineer", CategoryExternalID: "energy_environment", Name: "EHS工程师", OrderNum: 30, IsActive: true},
	{ExternalID: "env_water_treatment", CategoryExternalID: "energy_environment", Name: "水处理工程师", OrderNum: 40, IsActive: true},
	{ExternalID: "env_mep_water", CategoryExternalID: "energy_environment", Name: "给排水工程师", OrderNum: 50, IsActive: true},
	{ExternalID: "env_eia_engineer", CategoryExternalID: "energy_environment", Name: "环评工程师", OrderNum: 60, IsActive: true},
	{ExternalID: "env_tech", CategoryExternalID: "energy_environment", Name: "环保技术员", OrderNum: 70, IsActive: true},
	{ExternalID: "env_inspection", CategoryExternalID: "energy_environment", Name: "环保检测", OrderNum: 80, IsActive: true},
	{ExternalID: "env_specialist", CategoryExternalID: "energy_environment", Name: "环保专员", OrderNum: 90, IsActive: true},
	{ExternalID: "env_supervisor", CategoryExternalID: "energy_environment", Name: "环保主管", OrderNum: 100, IsActive: true},

	{ExternalID: "agri_tech", CategoryExternalID: "agri_planting", Name: "农业技术员", OrderNum: 10, IsActive: true},
	{ExternalID: "agri_agronomist", CategoryExternalID: "agri_planting", Name: "农艺师", OrderNum: 20, IsActive: true},
	{ExternalID: "agri_horticulturist", CategoryExternalID: "agri_planting", Name: "园艺师", OrderNum: 30, IsActive: true},
	{ExternalID: "agri_florist", CategoryExternalID: "agri_planting", Name: "花艺师", OrderNum: 40, IsActive: true},
	{ExternalID: "agri_farm_machine", CategoryExternalID: "agri_planting", Name: "农机操作修理", OrderNum: 50, IsActive: true},

	{ExternalID: "forestry_engineer", CategoryExternalID: "agri_forestry", Name: "林业工程师", OrderNum: 10, IsActive: true},
	{ExternalID: "forestry_tech", CategoryExternalID: "agri_forestry", Name: "林业技术员", OrderNum: 20, IsActive: true},
	{ExternalID: "forestry_garden_engineer", CategoryExternalID: "agri_forestry", Name: "园林工程师", OrderNum: 30, IsActive: true},
	{ExternalID: "forestry_ranger", CategoryExternalID: "agri_forestry", Name: "护林员", OrderNum: 40, IsActive: true},

	{ExternalID: "livestock_specialist", CategoryExternalID: "agri_livestock", Name: "畜牧师", OrderNum: 10, IsActive: true},
	{ExternalID: "livestock_vet", CategoryExternalID: "agri_livestock", Name: "兽医", OrderNum: 20, IsActive: true},
	{ExternalID: "livestock_breeding_tech", CategoryExternalID: "agri_livestock", Name: "养殖技术员", OrderNum: 30, IsActive: true},
	{ExternalID: "livestock_farm_mgmt", CategoryExternalID: "agri_livestock", Name: "牧场管理", OrderNum: 40, IsActive: true},
	{ExternalID: "livestock_feeder", CategoryExternalID: "agri_livestock", Name: "饲养员", OrderNum: 50, IsActive: true},
	{ExternalID: "livestock_quarantine", CategoryExternalID: "agri_livestock", Name: "动物检疫员", OrderNum: 60, IsActive: true},
	{ExternalID: "livestock_trainer", CategoryExternalID: "agri_livestock", Name: "动物驯养师", OrderNum: 70, IsActive: true},
	{ExternalID: "livestock_guide", CategoryExternalID: "agri_livestock", Name: "动物讲解员", OrderNum: 80, IsActive: true},

	{ExternalID: "fishery_farmer", CategoryExternalID: "agri_fishery", Name: "水产养殖员", OrderNum: 10, IsActive: true},
	{ExternalID: "fishery_tech", CategoryExternalID: "agri_fishery", Name: "水产技术员", OrderNum: 20, IsActive: true},
	{ExternalID: "fishery_aquaculture", CategoryExternalID: "agri_fishery", Name: "渔业养殖员", OrderNum: 30, IsActive: true},
	{ExternalID: "fishery_fisher", CategoryExternalID: "agri_fishery", Name: "捕捞员", OrderNum: 40, IsActive: true},

	{ExternalID: "public_civil_servant", CategoryExternalID: "public_services", Name: "公务员", OrderNum: 10, IsActive: true},
	{ExternalID: "public_ngo_specialist", CategoryExternalID: "public_services", Name: "非营利组织专员", OrderNum: 20, IsActive: true},
	{ExternalID: "public_urban_mgmt", CategoryExternalID: "public_services", Name: "城管", OrderNum: 30, IsActive: true},
	{ExternalID: "public_military", CategoryExternalID: "public_services", Name: "军人", OrderNum: 40, IsActive: true},
	{ExternalID: "public_firefighter", CategoryExternalID: "public_services", Name: "消防员", OrderNum: 50, IsActive: true},
	{ExternalID: "public_community_worker", CategoryExternalID: "public_services", Name: "社区工作者", OrderNum: 60, IsActive: true},
	{ExternalID: "public_police", CategoryExternalID: "public_services", Name: "警察/辅警", OrderNum: 70, IsActive: true},
	{ExternalID: "public_party_affairs", CategoryExternalID: "public_services", Name: "党建管理岗", OrderNum: 80, IsActive: true},

	{ExternalID: "research_assistant", CategoryExternalID: "public_research", Name: "科研助理", OrderNum: 10, IsActive: true},
	{ExternalID: "research_staff", CategoryExternalID: "public_research", Name: "科研人员", OrderNum: 20, IsActive: true},
	{ExternalID: "research_management", CategoryExternalID: "public_research", Name: "科研管理", OrderNum: 30, IsActive: true},
	{ExternalID: "research_academic_promo", CategoryExternalID: "public_research", Name: "学术推广", OrderNum: 40, IsActive: true},
	{ExternalID: "research_chem_analysis", CategoryExternalID: "public_research", Name: "化学分析", OrderNum: 50, IsActive: true},

	{ExternalID: "social_worker", CategoryExternalID: "public_social", Name: "社工", OrderNum: 10, IsActive: true},
	{ExternalID: "social_worker_assistant", CategoryExternalID: "public_social", Name: "社工助理", OrderNum: 20, IsActive: true},
	{ExternalID: "social_worker_intern", CategoryExternalID: "public_social", Name: "社工实习生", OrderNum: 30, IsActive: true},
	{ExternalID: "volunteer", CategoryExternalID: "public_social", Name: "志愿者", OrderNum: 40, IsActive: true},
	{ExternalID: "unpaid_volunteer", CategoryExternalID: "public_social", Name: "义工", OrderNum: 50, IsActive: true},
	{ExternalID: "caregiver", CategoryExternalID: "public_social", Name: "护理员", OrderNum: 60, IsActive: true},
	{ExternalID: "rehab_therapist", CategoryExternalID: "public_social", Name: "康复师", OrderNum: 70, IsActive: true},
}

var seedCategories = []SeedJobCategory{
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
}
