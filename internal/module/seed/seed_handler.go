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
	Names            map[string]string
	ParentExternalID string
	OrderNum         int
	IsActive         bool
}

type SeedJobRole struct {
	ExternalID         string
	CategoryExternalID string
	Names              map[string]string
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
				ParentID: parentID,
				OrderNum: sc.OrderNum,
				IsActive: sc.IsActive,
			}
			if err := taxRepo.UpsertJobCategoryWithNames(tx, &m, sc.Names); err != nil {
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
				OrderNum:   sr.OrderNum,
				IsActive:   sr.IsActive,
			}
			if err := taxRepo.UpsertJobRoleWithNames(tx, &m, sr.Names); err != nil {
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
	{Name: "Java 开发", Language: "zh", RoleCode: "Java", DataJSON: string(presets.GenerateJavaPreset()), IsActive: true},
	{Name: "Java Developer", Language: "en", RoleCode: "Java", DataJSON: string(presets.GenerateJavaPresetEn()), IsActive: true},
	{Name: "Python 开发", Language: "zh", RoleCode: "Python", DataJSON: string(presets.GeneratePythonPreset()), IsActive: true},
	{Name: "Python Developer", Language: "en", RoleCode: "Python", DataJSON: string(presets.GeneratePythonPresetEn()), IsActive: true},
	{Name: "Go 开发", Language: "zh", RoleCode: "golang", DataJSON: string(presets.GenerateGolangPreset()), IsActive: true},
	{Name: "Go Developer", Language: "en", RoleCode: "golang", DataJSON: string(presets.GenerateGolangPresetEn()), IsActive: true},
	{Name: "PHP 开发", Language: "zh", RoleCode: "php", DataJSON: string(presets.GeneratePHPPreset()), IsActive: true},
	{Name: "PHP Developer", Language: "en", RoleCode: "php", DataJSON: string(presets.GeneratePHPPresetEn()), IsActive: true},
	{Name: "C/C++ 开发", Language: "zh", RoleCode: "c_cpp", DataJSON: string(presets.GenerateCCppPreset()), IsActive: true},
	{Name: "C/C++ Developer", Language: "en", RoleCode: "c_cpp", DataJSON: string(presets.GenerateCCppPresetEn()), IsActive: true},
	{Name: "C# 开发", Language: "zh", RoleCode: "csharp", DataJSON: string(presets.GenerateCSharpPreset()), IsActive: true},
	{Name: "C# Developer", Language: "en", RoleCode: "csharp", DataJSON: string(presets.GenerateCSharpPresetEn()), IsActive: true},
	{Name: ".NET 开发", Language: "zh", RoleCode: "dotnet", DataJSON: string(presets.GenerateDotnetPreset()), IsActive: true},
	{Name: ".NET Developer", Language: "en", RoleCode: "dotnet", DataJSON: string(presets.GenerateDotnetPresetEn()), IsActive: true},
	{Name: "Node.js 开发", Language: "zh", RoleCode: "nodejs", DataJSON: string(presets.GenerateNodejsPreset()), IsActive: true},
	{Name: "Node.js Developer", Language: "en", RoleCode: "nodejs", DataJSON: string(presets.GenerateNodejsPresetEn()), IsActive: true},
}

var seedRoles = []SeedJobRole{
	{ExternalID: "Java", CategoryExternalID: "it_backend", Names: map[string]string{"zh": "Java", "en": "Java"}, OrderNum: 10, IsActive: true},
	{ExternalID: "Python", CategoryExternalID: "it_backend", Names: map[string]string{"zh": "Python", "en": "Python"}, OrderNum: 20, IsActive: true},
	{ExternalID: "golang", CategoryExternalID: "it_backend", Names: map[string]string{"zh": "Go (Golang)", "en": "Go (Golang)"}, OrderNum: 30, IsActive: true},
	{ExternalID: "php", CategoryExternalID: "it_backend", Names: map[string]string{"zh": "PHP", "en": "PHP"}, OrderNum: 40, IsActive: true},
	{ExternalID: "c_cpp", CategoryExternalID: "it_backend", Names: map[string]string{"zh": "C/C++", "en": "C/C++"}, OrderNum: 50, IsActive: true},
	{ExternalID: "csharp", CategoryExternalID: "it_backend", Names: map[string]string{"zh": "C#", "en": "C#"}, OrderNum: 60, IsActive: true},
	{ExternalID: "dotnet", CategoryExternalID: "it_backend", Names: map[string]string{"zh": ".NET", "en": ".NET"}, OrderNum: 70, IsActive: true},
	{ExternalID: "nodejs", CategoryExternalID: "it_backend", Names: map[string]string{"zh": "Node.js", "en": "Node.js"}, OrderNum: 80, IsActive: true},

	{ExternalID: "web_frontend", CategoryExternalID: "it_frontend", Names: map[string]string{"zh": "Web前端", "en": "Web Frontend"}, OrderNum: 10, IsActive: true},
	{ExternalID: "html5", CategoryExternalID: "it_frontend", Names: map[string]string{"zh": "HTML5", "en": "HTML5"}, OrderNum: 20, IsActive: true},
	{ExternalID: "miniapp", CategoryExternalID: "it_frontend", Names: map[string]string{"zh": "小程序开发工程师", "en": "Mini App Developer"}, OrderNum: 30, IsActive: true},

	{ExternalID: "android", CategoryExternalID: "it_mobile", Names: map[string]string{"zh": "Android开发", "en": "Android Developer"}, OrderNum: 10, IsActive: true},
	{ExternalID: "ios", CategoryExternalID: "it_mobile", Names: map[string]string{"zh": "iOS开发", "en": "iOS Developer"}, OrderNum: 20, IsActive: true},
	{ExternalID: "harmony", CategoryExternalID: "it_mobile", Names: map[string]string{"zh": "鸿蒙开发工程师", "en": "HarmonyOS Developer"}, OrderNum: 30, IsActive: true},

	{ExternalID: "test_engineer", CategoryExternalID: "it_testing", Names: map[string]string{"zh": "测试工程师", "en": "Testing Engineer"}, OrderNum: 10, IsActive: true},
	{ExternalID: "automation_test", CategoryExternalID: "it_testing", Names: map[string]string{"zh": "自动化测试", "en": "Automation Testing"}, OrderNum: 20, IsActive: true},
	{ExternalID: "test_dev", CategoryExternalID: "it_testing", Names: map[string]string{"zh": "测试开发", "en": "Testing Development"}, OrderNum: 30, IsActive: true},
	{ExternalID: "performance_test", CategoryExternalID: "it_testing", Names: map[string]string{"zh": "性能测试", "en": "Performance Testing"}, OrderNum: 40, IsActive: true},
	{ExternalID: "hardware_test", CategoryExternalID: "it_testing", Names: map[string]string{"zh": "硬件测试工程师", "en": "Hardware Testing Engineer"}, OrderNum: 50, IsActive: true},

	{ExternalID: "ops_engineer", CategoryExternalID: "it_ops_sec_dba", Names: map[string]string{"zh": "运维工程师", "en": "Operations Engineer"}, OrderNum: 10, IsActive: true},
	{ExternalID: "devops", CategoryExternalID: "it_ops_sec_dba", Names: map[string]string{"zh": "DevOps工程师", "en": "DevOps Engineer"}, OrderNum: 20, IsActive: true},
	{ExternalID: "sys_admin", CategoryExternalID: "it_ops_sec_dba", Names: map[string]string{"zh": "系统/网络管理员", "en": "System/Network Administrator"}, OrderNum: 30, IsActive: true},
	{ExternalID: "dba", CategoryExternalID: "it_ops_sec_dba", Names: map[string]string{"zh": "数据库管理员 (DBA)", "en": "Database Administrator (DBA)"}, OrderNum: 40, IsActive: true},
	{ExternalID: "security_engineer", CategoryExternalID: "it_ops_sec_dba", Names: map[string]string{"zh": "安全工程师", "en": "Security Engineer"}, OrderNum: 50, IsActive: true},
	{ExternalID: "cloud_engineer", CategoryExternalID: "it_ops_sec_dba", Names: map[string]string{"zh": "云计算工程师", "en": "Cloud Engineer"}, OrderNum: 60, IsActive: true},
	{ExternalID: "ops_manager", CategoryExternalID: "it_ops_sec_dba", Names: map[string]string{"zh": "运维经理/主管", "en": "Operations Manager/Supervisor"}, OrderNum: 70, IsActive: true},

	{ExternalID: "data_mining", CategoryExternalID: "it_ai_bigdata", Names: map[string]string{"zh": "数据挖掘", "en": "Data Mining"}, OrderNum: 10, IsActive: true},
	{ExternalID: "nlp", CategoryExternalID: "it_ai_bigdata", Names: map[string]string{"zh": "自然语言处理", "en": "Natural Language Processing"}, OrderNum: 20, IsActive: true},
	{ExternalID: "ml_ai", CategoryExternalID: "it_ai_bigdata", Names: map[string]string{"zh": "机器学习/AI工程师", "en": "Machine Learning/AI Engineer"}, OrderNum: 30, IsActive: true},
	{ExternalID: "bigdata", CategoryExternalID: "it_ai_bigdata", Names: map[string]string{"zh": "大数据工程师", "en": "Big Data Engineer"}, OrderNum: 40, IsActive: true},
	{ExternalID: "blockchain", CategoryExternalID: "it_ai_bigdata", Names: map[string]string{"zh": "区块链开发工程师", "en": "Blockchain Developer"}, OrderNum: 50, IsActive: true},
	{ExternalID: "algo_engineer", CategoryExternalID: "it_ai_bigdata", Names: map[string]string{"zh": "算法工程师", "en": "Algorithm Engineer"}, OrderNum: 60, IsActive: true},

	{ExternalID: "tech_support", CategoryExternalID: "it_other_tech", Names: map[string]string{"zh": "技术支持工程师", "en": "Technical Support Engineer"}, OrderNum: 10, IsActive: true},
	{ExternalID: "tech_engineer", CategoryExternalID: "it_other_tech", Names: map[string]string{"zh": "技术工程师", "en": "Technical Engineer"}, OrderNum: 20, IsActive: true},
	{ExternalID: "presales_support", CategoryExternalID: "it_other_tech", Names: map[string]string{"zh": "售前售后/技术支持", "en": "Presales Support/Technical Support"}, OrderNum: 30, IsActive: true},
	{ExternalID: "other_tech_roles", CategoryExternalID: "it_other_tech", Names: map[string]string{"zh": "其他技术岗位", "en": "Other Technical Roles"}, OrderNum: 40, IsActive: true},
	{ExternalID: "network_engineer", CategoryExternalID: "it_other_tech", Names: map[string]string{"zh": "网络工程师", "en": "Network Engineer"}, OrderNum: 50, IsActive: true},
	{ExternalID: "hardware_dev", CategoryExternalID: "it_other_tech", Names: map[string]string{"zh": "硬件开发工程师", "en": "Hardware Developer"}, OrderNum: 60, IsActive: true},
	{ExternalID: "system_integrator", CategoryExternalID: "it_other_tech", Names: map[string]string{"zh": "系统集成工程师", "en": "System Integrator"}, OrderNum: 70, IsActive: true},
	{ExternalID: "circuit_engineer", CategoryExternalID: "it_other_tech", Names: map[string]string{"zh": "电路工程师", "en": "Circuit Engineer"}, OrderNum: 80, IsActive: true},

	{ExternalID: "architect", CategoryExternalID: "it_senior", Names: map[string]string{"zh": "架构师", "en": "Architect"}, OrderNum: 10, IsActive: true},
	{ExternalID: "tech_manager", CategoryExternalID: "it_senior", Names: map[string]string{"zh": "技术主管", "en": "Technical Manager"}, OrderNum: 20, IsActive: true},
	{ExternalID: "tech_director", CategoryExternalID: "it_senior", Names: map[string]string{"zh": "技术经理/总监", "en": "Technical Manager/Director"}, OrderNum: 30, IsActive: true},
	{ExternalID: "rd_manager", CategoryExternalID: "it_senior", Names: map[string]string{"zh": "研发经理/总监", "en": "Research and Development Manager/Director"}, OrderNum: 40, IsActive: true},
	{ExternalID: "cto", CategoryExternalID: "it_senior", Names: map[string]string{"zh": "CTO", "en": "CTO"}, OrderNum: 50, IsActive: true},
	{ExternalID: "fullstack", CategoryExternalID: "it_senior", Names: map[string]string{"zh": "全栈工程师", "en": "Fullstack Engineer"}, OrderNum: 60, IsActive: true},

	{ExternalID: "finance_teller", CategoryExternalID: "finance_counter_service", Names: map[string]string{"zh": "银行柜员", "en": "Bank Teller"}, OrderNum: 10, IsActive: true},
	{ExternalID: "finance_general_teller", CategoryExternalID: "finance_counter_service", Names: map[string]string{"zh": "银行综合柜员", "en": "General Bank Teller"}, OrderNum: 20, IsActive: true},
	{ExternalID: "finance_lobby_manager", CategoryExternalID: "finance_counter_service", Names: map[string]string{"zh": "银行大堂经理", "en": "Lobby Manager"}, OrderNum: 30, IsActive: true},
	{ExternalID: "finance_lobby_guide", CategoryExternalID: "finance_counter_service", Names: map[string]string{"zh": "大堂引导员/大堂助理", "en": "Lobby Guide/Lobby Assistant"}, OrderNum: 40, IsActive: true},
	{ExternalID: "finance_bank_cs", CategoryExternalID: "finance_counter_service", Names: map[string]string{"zh": "银行客服/坐席员", "en": "Bank Customer Service/Agent"}, OrderNum: 50, IsActive: true},
	{ExternalID: "finance_bank_frontdesk", CategoryExternalID: "finance_counter_service", Names: map[string]string{"zh": "银行前台", "en": "Bank Front Desk"}, OrderNum: 60, IsActive: true},

	{ExternalID: "finance_account_manager", CategoryExternalID: "finance_personal_wealth", Names: map[string]string{"zh": "客户经理", "en": "Account Manager"}, OrderNum: 10, IsActive: true},
	{ExternalID: "finance_wealth_manager", CategoryExternalID: "finance_personal_wealth", Names: map[string]string{"zh": "理财经理", "en": "Wealth Manager"}, OrderNum: 20, IsActive: true},
	{ExternalID: "finance_wealth_advisor", CategoryExternalID: "finance_personal_wealth", Names: map[string]string{"zh": "理财顾问", "en": "Wealth Advisor"}, OrderNum: 30, IsActive: true},
	{ExternalID: "finance_investment_advisor", CategoryExternalID: "finance_personal_wealth", Names: map[string]string{"zh": "投资顾问", "en": "Investment Advisor"}, OrderNum: 40, IsActive: true},
	{ExternalID: "finance_creditcard_sales", CategoryExternalID: "finance_personal_wealth", Names: map[string]string{"zh": "信用卡销售", "en": "Credit Card Sales"}, OrderNum: 50, IsActive: true},

	{ExternalID: "finance_credit_manager", CategoryExternalID: "finance_credit_approval", Names: map[string]string{"zh": "信贷经理", "en": "Credit Manager"}, OrderNum: 10, IsActive: true},
	{ExternalID: "finance_credit_officer", CategoryExternalID: "finance_credit_approval", Names: map[string]string{"zh": "信贷专员", "en": "Credit Officer"}, OrderNum: 20, IsActive: true},
	{ExternalID: "finance_loan_officer", CategoryExternalID: "finance_credit_approval", Names: map[string]string{"zh": "贷款专员", "en": "Loan Officer"}, OrderNum: 30, IsActive: true},
	{ExternalID: "finance_postloan", CategoryExternalID: "finance_credit_approval", Names: map[string]string{"zh": "贷后管理岗", "en": "Post-Loan Management"}, OrderNum: 40, IsActive: true},
	{ExternalID: "finance_collections", CategoryExternalID: "finance_credit_approval", Names: map[string]string{"zh": "催收岗", "en": "Collections"}, OrderNum: 50, IsActive: true},
	{ExternalID: "finance_mortgage_officer", CategoryExternalID: "finance_credit_approval", Names: map[string]string{"zh": "按揭专员", "en": "Mortgage Officer"}, OrderNum: 60, IsActive: true},
	{ExternalID: "finance_credit_admin", CategoryExternalID: "finance_credit_approval", Names: map[string]string{"zh": "信贷管理", "en": "Credit Administration"}, OrderNum: 70, IsActive: true},

	{ExternalID: "finance_risk_manager", CategoryExternalID: "finance_risk_compliance", Names: map[string]string{"zh": "风险经理", "en": "Risk Manager"}, OrderNum: 10, IsActive: true},
	{ExternalID: "finance_compliance_manager", CategoryExternalID: "finance_risk_compliance", Names: map[string]string{"zh": "合规经理", "en": "Compliance Manager"}, OrderNum: 20, IsActive: true},
	{ExternalID: "finance_risk_control", CategoryExternalID: "finance_risk_compliance", Names: map[string]string{"zh": "风控专员", "en": "Risk Control Officer"}, OrderNum: 30, IsActive: true},
	{ExternalID: "finance_audit", CategoryExternalID: "finance_risk_compliance", Names: map[string]string{"zh": "审计专员", "en": "Audit Officer"}, OrderNum: 40, IsActive: true},
	{ExternalID: "finance_fin_compliance", CategoryExternalID: "finance_risk_compliance", Names: map[string]string{"zh": "金融合规专员", "en": "Financial Compliance Officer"}, OrderNum: 50, IsActive: true},
	{ExternalID: "finance_aml", CategoryExternalID: "finance_risk_compliance", Names: map[string]string{"zh": "反洗钱专员", "en": "Anti-Money Laundering Officer"}, OrderNum: 60, IsActive: true},
	{ExternalID: "finance_legal", CategoryExternalID: "finance_risk_compliance", Names: map[string]string{"zh": "法律事务岗", "en": "Legal Affairs Officer"}, OrderNum: 70, IsActive: true},

	{ExternalID: "finance_securities_broker", CategoryExternalID: "finance_securities_invest", Names: map[string]string{"zh": "证券经纪人", "en": "Securities Broker"}, OrderNum: 10, IsActive: true},
	{ExternalID: "finance_invest_manager", CategoryExternalID: "finance_securities_invest", Names: map[string]string{"zh": "投资经理", "en": "Investment Manager"}, OrderNum: 20, IsActive: true},
	{ExternalID: "finance_fund_manager", CategoryExternalID: "finance_securities_invest", Names: map[string]string{"zh": "基金经理", "en": "Fund Manager"}, OrderNum: 30, IsActive: true},
	{ExternalID: "finance_trader", CategoryExternalID: "finance_securities_invest", Names: map[string]string{"zh": "交易员", "en": "Trader"}, OrderNum: 40, IsActive: true},
	{ExternalID: "finance_fund_accountant", CategoryExternalID: "finance_securities_invest", Names: map[string]string{"zh": "基金会计", "en": "Fund Accountant"}, OrderNum: 50, IsActive: true},
	{ExternalID: "finance_settlement", CategoryExternalID: "finance_securities_invest", Names: map[string]string{"zh": "清算专员", "en": "Settlement Officer"}, OrderNum: 60, IsActive: true},
	{ExternalID: "finance_researcher", CategoryExternalID: "finance_securities_invest", Names: map[string]string{"zh": "研究员", "en": "Researcher"}, OrderNum: 70, IsActive: true},
	{ExternalID: "finance_analyst", CategoryExternalID: "finance_securities_invest", Names: map[string]string{"zh": "金融分析师", "en": "Financial Analyst"}, OrderNum: 80, IsActive: true},
	{ExternalID: "finance_securities_analyst", CategoryExternalID: "finance_securities_invest", Names: map[string]string{"zh": "证券分析师", "en": "Securities Analyst"}, OrderNum: 90, IsActive: true},

	{ExternalID: "finance_insurance_advisor", CategoryExternalID: "finance_insurance_actuary", Names: map[string]string{"zh": "保险顾问", "en": "Insurance Advisor"}, OrderNum: 10, IsActive: true},
	{ExternalID: "finance_insurance_broker", CategoryExternalID: "finance_insurance_actuary", Names: map[string]string{"zh": "保险经纪人", "en": "Insurance Broker"}, OrderNum: 20, IsActive: true},
	{ExternalID: "finance_insurance_agent", CategoryExternalID: "finance_insurance_actuary", Names: map[string]string{"zh": "保险代理专员", "en": "Insurance Agent"}, OrderNum: 30, IsActive: true},
	{ExternalID: "finance_underwriter", CategoryExternalID: "finance_insurance_actuary", Names: map[string]string{"zh": "核保师", "en": "Underwriter"}, OrderNum: 40, IsActive: true},
	{ExternalID: "finance_claims", CategoryExternalID: "finance_insurance_actuary", Names: map[string]string{"zh": "理赔师", "en": "Claims Officer"}, OrderNum: 50, IsActive: true},
	{ExternalID: "finance_actuary", CategoryExternalID: "finance_insurance_actuary", Names: map[string]string{"zh": "精算师", "en": "Actuary"}, OrderNum: 60, IsActive: true},
	{ExternalID: "finance_insurance_trainer", CategoryExternalID: "finance_insurance_actuary", Names: map[string]string{"zh": "保险培训师", "en": "Insurance Trainer"}, OrderNum: 70, IsActive: true},
	{ExternalID: "finance_insurance_ops", CategoryExternalID: "finance_insurance_actuary", Names: map[string]string{"zh": "保险内勤", "en": "Insurance Office Worker"}, OrderNum: 80, IsActive: true},
	{ExternalID: "finance_insurance_coach", CategoryExternalID: "finance_insurance_actuary", Names: map[string]string{"zh": "保险组训", "en": "Insurance Coaching"}, OrderNum: 90, IsActive: true},
	{ExternalID: "finance_surveyor", CategoryExternalID: "finance_insurance_actuary", Names: map[string]string{"zh": "查勘员", "en": "Surveyor"}, OrderNum: 100, IsActive: true},
	{ExternalID: "finance_insurance_sales", CategoryExternalID: "finance_insurance_actuary", Names: map[string]string{"zh": "保险销售", "en": "Insurance Sales"}, OrderNum: 110, IsActive: true},

	{ExternalID: "finance_bank_ops_lead", CategoryExternalID: "finance_banking_support", Names: map[string]string{"zh": "银行运营主管", "en": "Bank Operations Supervisor"}, OrderNum: 10, IsActive: true},
	{ExternalID: "finance_data_entry", CategoryExternalID: "finance_banking_support", Names: map[string]string{"zh": "数据录入员", "en": "Data Entry Operator"}, OrderNum: 20, IsActive: true},
	{ExternalID: "finance_doc_specialist", CategoryExternalID: "finance_banking_support", Names: map[string]string{"zh": "单证处理专员", "en": "Document Specialist"}, OrderNum: 30, IsActive: true},
	{ExternalID: "finance_data_analyst", CategoryExternalID: "finance_banking_support", Names: map[string]string{"zh": "金融数据分析师", "en": "Financial Data Analyst"}, OrderNum: 40, IsActive: true},
	{ExternalID: "finance_funds_settlement", CategoryExternalID: "finance_banking_support", Names: map[string]string{"zh": "资金结算专员", "en": "Funds Settlement Officer"}, OrderNum: 50, IsActive: true},

	{ExternalID: "finance_trust_manager", CategoryExternalID: "finance_trust_futures", Names: map[string]string{"zh": "信托经理", "en": "Trust Manager"}, OrderNum: 10, IsActive: true},
	{ExternalID: "finance_trader_operator", CategoryExternalID: "finance_trust_futures", Names: map[string]string{"zh": "操盘手", "en": "Trader Operator"}, OrderNum: 20, IsActive: true},
	{ExternalID: "finance_futures_analyst", CategoryExternalID: "finance_trust_futures", Names: map[string]string{"zh": "期货分析师", "en": "Futures Analyst"}, OrderNum: 30, IsActive: true},
	{ExternalID: "finance_invest_strategy", CategoryExternalID: "finance_trust_futures", Names: map[string]string{"zh": "投资策略师", "en": "Investment Strategy Consultant"}, OrderNum: 40, IsActive: true},

	{ExternalID: "finance_branch_manager", CategoryExternalID: "finance_bank_management", Names: map[string]string{"zh": "支行行长", "en": "Branch Manager"}, OrderNum: 10, IsActive: true},
	{ExternalID: "finance_branch_vice_manager", CategoryExternalID: "finance_bank_management", Names: map[string]string{"zh": "支行副行长", "en": "Branch Vice Manager"}, OrderNum: 20, IsActive: true},
	{ExternalID: "finance_subbranch_manager", CategoryExternalID: "finance_bank_management", Names: map[string]string{"zh": "分行行长", "en": "Sub-branch Manager"}, OrderNum: 30, IsActive: true},
	{ExternalID: "finance_cfo", CategoryExternalID: "finance_bank_management", Names: map[string]string{"zh": "首席财务官 (CFO)", "en": "Chief Financial Officer (CFO)"}, OrderNum: 40, IsActive: true},
	{ExternalID: "finance_invest_director", CategoryExternalID: "finance_bank_management", Names: map[string]string{"zh": "投资总监", "en": "Investment Director"}, OrderNum: 50, IsActive: true},

	{ExternalID: "finance_mt", CategoryExternalID: "finance_intern", Names: map[string]string{"zh": "银行管理培训生", "en": "Bank Management Training"}, OrderNum: 10, IsActive: true},
	{ExternalID: "finance_securities_intern", CategoryExternalID: "finance_intern", Names: map[string]string{"zh": "证券实习生", "en": "Securities Intern"}, OrderNum: 20, IsActive: true},
	{ExternalID: "finance_insurance_intern", CategoryExternalID: "finance_intern", Names: map[string]string{"zh": "保险实习生", "en": "Insurance Intern"}, OrderNum: 30, IsActive: true},
	{ExternalID: "finance_ib_intern", CategoryExternalID: "finance_intern", Names: map[string]string{"zh": "投资银行实习生", "en": "Investment Bank Intern"}, OrderNum: 40, IsActive: true},
	{ExternalID: "finance_intern_assistant", CategoryExternalID: "finance_intern", Names: map[string]string{"zh": "实习助理", "en": "Intern Assistant"}, OrderNum: 50, IsActive: true},
	{ExternalID: "finance_teller_intern", CategoryExternalID: "finance_intern", Names: map[string]string{"zh": "柜员实习生", "en": "Teller Intern"}, OrderNum: 60, IsActive: true},
	{ExternalID: "finance_bank_intern", CategoryExternalID: "finance_intern", Names: map[string]string{"zh": "银行实习生", "en": "Bank Intern"}, OrderNum: 70, IsActive: true},

	{ExternalID: "edu_chinese_teacher", CategoryExternalID: "education_teacher", Names: map[string]string{"zh": "语文教师", "en": "Chinese Teacher"}, OrderNum: 10, IsActive: true},
	{ExternalID: "edu_english_teacher", CategoryExternalID: "education_teacher", Names: map[string]string{"zh": "英语教师", "en": "English Teacher"}, OrderNum: 20, IsActive: true},
	{ExternalID: "edu_art_teacher", CategoryExternalID: "education_teacher", Names: map[string]string{"zh": "美术老师", "en": "Art Teacher"}, OrderNum: 30, IsActive: true},
	{ExternalID: "edu_math_teacher", CategoryExternalID: "education_teacher", Names: map[string]string{"zh": "数学教师", "en": "Math Teacher"}, OrderNum: 40, IsActive: true},
	{ExternalID: "edu_kindergarten_teacher", CategoryExternalID: "education_teacher", Names: map[string]string{"zh": "幼儿教师", "en": "Kindergarten Teacher"}, OrderNum: 50, IsActive: true},
	{ExternalID: "edu_pe_teacher", CategoryExternalID: "education_teacher", Names: map[string]string{"zh": "体育教师", "en": "Physical Education Teacher"}, OrderNum: 60, IsActive: true},
	{ExternalID: "edu_music_teacher", CategoryExternalID: "education_teacher", Names: map[string]string{"zh": "音乐老师", "en": "Music Teacher"}, OrderNum: 70, IsActive: true},
	{ExternalID: "edu_biology_teacher", CategoryExternalID: "education_teacher", Names: map[string]string{"zh": "生物教师", "en": "Biology Teacher"}, OrderNum: 80, IsActive: true},
	{ExternalID: "edu_dance_teacher", CategoryExternalID: "education_teacher", Names: map[string]string{"zh": "舞蹈老师", "en": "Dance Teacher"}, OrderNum: 90, IsActive: true},
	{ExternalID: "edu_piano_teacher", CategoryExternalID: "education_teacher", Names: map[string]string{"zh": "钢琴教师", "en": "Piano Teacher"}, OrderNum: 100, IsActive: true},
	{ExternalID: "edu_calligraphy_teacher", CategoryExternalID: "education_teacher", Names: map[string]string{"zh": "书法教师", "en": "Calligraphy Teacher"}, OrderNum: 110, IsActive: true},
	{ExternalID: "edu_chemistry_teacher", CategoryExternalID: "education_teacher", Names: map[string]string{"zh": "化学老师", "en": "Chemistry Teacher"}, OrderNum: 120, IsActive: true},
	{ExternalID: "edu_physics_teacher", CategoryExternalID: "education_teacher", Names: map[string]string{"zh": "物理老师", "en": "Physics Teacher"}, OrderNum: 130, IsActive: true},
	{ExternalID: "edu_history_teacher", CategoryExternalID: "education_teacher", Names: map[string]string{"zh": "历史老师", "en": "History Teacher"}, OrderNum: 140, IsActive: true},
	{ExternalID: "edu_politics_teacher", CategoryExternalID: "education_teacher", Names: map[string]string{"zh": "政治老师", "en": "Politics Teacher"}, OrderNum: 150, IsActive: true},
	{ExternalID: "edu_geography_teacher", CategoryExternalID: "education_teacher", Names: map[string]string{"zh": "地理老师", "en": "Geography Teacher"}, OrderNum: 160, IsActive: true},
	{ExternalID: "edu_tcsol_teacher", CategoryExternalID: "education_teacher", Names: map[string]string{"zh": "对外汉语教师", "en": "Chinese Language Teacher (TCFL)"}, OrderNum: 170, IsActive: true},
	{ExternalID: "edu_tutor", CategoryExternalID: "education_teacher", Names: map[string]string{"zh": "家教", "en": "Private Tutor"}, OrderNum: 180, IsActive: true},
	{ExternalID: "edu_university_teacher", CategoryExternalID: "education_teacher", Names: map[string]string{"zh": "大学教师", "en": "University Lecturer"}, OrderNum: 190, IsActive: true},
	{ExternalID: "edu_tutoring_teacher", CategoryExternalID: "education_teacher", Names: map[string]string{"zh": "辅导老师", "en": "Tutoring Teacher"}, OrderNum: 200, IsActive: true},

	{ExternalID: "edu_principal", CategoryExternalID: "education_teaching_admin", Names: map[string]string{"zh": "校长", "en": "Principal"}, OrderNum: 10, IsActive: true},
	{ExternalID: "edu_vice_principal", CategoryExternalID: "education_teaching_admin", Names: map[string]string{"zh": "副校长", "en": "Vice Principal"}, OrderNum: 20, IsActive: true},
	{ExternalID: "edu_kindergarten_principal", CategoryExternalID: "education_teaching_admin", Names: map[string]string{"zh": "园长", "en": "Kindergarten Principal"}, OrderNum: 30, IsActive: true},
	{ExternalID: "edu_academic_director", CategoryExternalID: "education_teaching_admin", Names: map[string]string{"zh": "教务主任", "en": "Academic Director"}, OrderNum: 40, IsActive: true},
	{ExternalID: "edu_research_leader", CategoryExternalID: "education_teaching_admin", Names: map[string]string{"zh": "教研组长", "en": "Teaching Research Lead"}, OrderNum: 50, IsActive: true},
	{ExternalID: "edu_teaching_supervisor", CategoryExternalID: "education_teaching_admin", Names: map[string]string{"zh": "教学主管", "en": "Teaching Supervisor"}, OrderNum: 60, IsActive: true},
	{ExternalID: "edu_academic_specialist", CategoryExternalID: "education_teaching_admin", Names: map[string]string{"zh": "教务专员", "en": "Academic Affairs Specialist"}, OrderNum: 70, IsActive: true},
	{ExternalID: "edu_researcher", CategoryExternalID: "education_teaching_admin", Names: map[string]string{"zh": "教研组长", "en": "Teaching Research Leader"}, OrderNum: 80, IsActive: true},
	{ExternalID: "edu_curriculum_dev", CategoryExternalID: "education_teaching_admin", Names: map[string]string{"zh": "课程开发/设计", "en": "Curriculum Developer"}, OrderNum: 90, IsActive: true},
	{ExternalID: "edu_academic_assistant", CategoryExternalID: "education_teaching_admin", Names: map[string]string{"zh": "教务助理", "en": "Academic Assistant"}, OrderNum: 100, IsActive: true},

	{ExternalID: "edu_dorm_teacher", CategoryExternalID: "education_student_services", Names: map[string]string{"zh": "宿管老师", "en": "Dormitory Supervisor"}, OrderNum: 10, IsActive: true},
	{ExternalID: "edu_admissions", CategoryExternalID: "education_student_services", Names: map[string]string{"zh": "招生顾问", "en": "Admissions Consultant"}, OrderNum: 20, IsActive: true},
	{ExternalID: "edu_counselor", CategoryExternalID: "education_student_services", Names: map[string]string{"zh": "辅导员", "en": "Student Counselor"}, OrderNum: 30, IsActive: true},
	{ExternalID: "edu_ta", CategoryExternalID: "education_student_services", Names: map[string]string{"zh": "助教/教学助理", "en": "Teaching Assistant"}, OrderNum: 40, IsActive: true},
	{ExternalID: "edu_life_teacher", CategoryExternalID: "education_student_services", Names: map[string]string{"zh": "生活老师", "en": "Student Life Teacher"}, OrderNum: 50, IsActive: true},
	{ExternalID: "edu_study_abroad", CategoryExternalID: "education_student_services", Names: map[string]string{"zh": "留学顾问", "en": "Study Abroad Consultant"}, OrderNum: 60, IsActive: true},
	{ExternalID: "edu_homeroom_teacher", CategoryExternalID: "education_student_services", Names: map[string]string{"zh": "班主任", "en": "Homeroom Teacher"}, OrderNum: 70, IsActive: true},
	{ExternalID: "edu_student_manager", CategoryExternalID: "education_student_services", Names: map[string]string{"zh": "学管师", "en": "Student Manager"}, OrderNum: 80, IsActive: true},

	{ExternalID: "edu_trainer", CategoryExternalID: "education_training_lecturer", Names: map[string]string{"zh": "培训师", "en": "Trainer"}, OrderNum: 10, IsActive: true},
	{ExternalID: "edu_lecturer", CategoryExternalID: "education_training_lecturer", Names: map[string]string{"zh": "培训讲师", "en": "Training Lecturer"}, OrderNum: 20, IsActive: true},
	{ExternalID: "edu_english_trainer", CategoryExternalID: "education_training_lecturer", Names: map[string]string{"zh": "英语培训老师", "en": "English Trainer"}, OrderNum: 30, IsActive: true},
	{ExternalID: "edu_corporate_trainer", CategoryExternalID: "education_training_lecturer", Names: map[string]string{"zh": "企业培训师", "en": "Corporate Trainer"}, OrderNum: 40, IsActive: true},
	{ExternalID: "edu_vocational_trainer", CategoryExternalID: "education_training_lecturer", Names: map[string]string{"zh": "职业培训师", "en": "Vocational Trainer"}, OrderNum: 50, IsActive: true},
	{ExternalID: "edu_speaker", CategoryExternalID: "education_training_lecturer", Names: map[string]string{"zh": "讲师", "en": "Speaker"}, OrderNum: 60, IsActive: true},
	{ExternalID: "edu_internal_trainer", CategoryExternalID: "education_training_lecturer", Names: map[string]string{"zh": "内训师", "en": "Internal Trainer"}, OrderNum: 70, IsActive: true},

	{ExternalID: "edu_training_specialist", CategoryExternalID: "education_training_management", Names: map[string]string{"zh": "培训专员", "en": "Training Specialist"}, OrderNum: 10, IsActive: true},
	{ExternalID: "edu_training_supervisor", CategoryExternalID: "education_training_management", Names: map[string]string{"zh": "培训主管", "en": "Training Supervisor"}, OrderNum: 20, IsActive: true},
	{ExternalID: "edu_training_assistant", CategoryExternalID: "education_training_management", Names: map[string]string{"zh": "培训助理", "en": "Training Assistant"}, OrderNum: 30, IsActive: true},
	{ExternalID: "edu_training_director", CategoryExternalID: "education_training_management", Names: map[string]string{"zh": "培训总监", "en": "Training Director"}, OrderNum: 40, IsActive: true},
	{ExternalID: "edu_training_manager", CategoryExternalID: "education_training_management", Names: map[string]string{"zh": "培训经理", "en": "Training Manager"}, OrderNum: 50, IsActive: true},

	{ExternalID: "hc_clinical_doctor", CategoryExternalID: "healthcare_doctor", Names: map[string]string{"zh": "临床医生", "en": "Clinical Doctor"}, OrderNum: 10, IsActive: true},
	{ExternalID: "hc_internal_doctor", CategoryExternalID: "healthcare_doctor", Names: map[string]string{"zh": "内科医生", "en": "Internal Medicine Doctor"}, OrderNum: 20, IsActive: true},
	{ExternalID: "hc_surgery_doctor", CategoryExternalID: "healthcare_doctor", Names: map[string]string{"zh": "外科医生", "en": "Surgeon"}, OrderNum: 30, IsActive: true},
	{ExternalID: "hc_obgyn_doctor", CategoryExternalID: "healthcare_doctor", Names: map[string]string{"zh": "妇产科医生", "en": "Obstetrician-Gynecologist"}, OrderNum: 40, IsActive: true},
	{ExternalID: "hc_pediatrics_doctor", CategoryExternalID: "healthcare_doctor", Names: map[string]string{"zh": "儿科医生", "en": "Pediatrician"}, OrderNum: 50, IsActive: true},
	{ExternalID: "hc_orthopedics_doctor", CategoryExternalID: "healthcare_doctor", Names: map[string]string{"zh": "骨科医生", "en": "Orthopedic Doctor"}, OrderNum: 60, IsActive: true},
	{ExternalID: "hc_anesthesiologist", CategoryExternalID: "healthcare_doctor", Names: map[string]string{"zh": "麻醉医生", "en": "Anesthesiologist"}, OrderNum: 70, IsActive: true},
	{ExternalID: "hc_dentist", CategoryExternalID: "healthcare_doctor", Names: map[string]string{"zh": "口腔医生", "en": "Dentist"}, OrderNum: 80, IsActive: true},
	{ExternalID: "hc_tcm", CategoryExternalID: "healthcare_doctor", Names: map[string]string{"zh": "中医师", "en": "Traditional Chinese Medicine Doctor"}, OrderNum: 90, IsActive: true},
	{ExternalID: "hc_radiologist", CategoryExternalID: "healthcare_doctor", Names: map[string]string{"zh": "放射科医生", "en": "Radiologist"}, OrderNum: 100, IsActive: true},
	{ExternalID: "hc_gp", CategoryExternalID: "healthcare_doctor", Names: map[string]string{"zh": "全科医生", "en": "General Practitioner"}, OrderNum: 110, IsActive: true},
	{ExternalID: "hc_specialist_doctor", CategoryExternalID: "healthcare_doctor", Names: map[string]string{"zh": "专科医生", "en": "Specialist Doctor"}, OrderNum: 120, IsActive: true},
	{ExternalID: "hc_doctor_assistant", CategoryExternalID: "healthcare_doctor", Names: map[string]string{"zh": "医生助理", "en": "Physician Assistant"}, OrderNum: 130, IsActive: true},
	{ExternalID: "hc_resident", CategoryExternalID: "healthcare_doctor", Names: map[string]string{"zh": "住院医师", "en": "Resident Doctor"}, OrderNum: 140, IsActive: true},
	{ExternalID: "hc_dental_doctor", CategoryExternalID: "healthcare_doctor", Names: map[string]string{"zh": "牙科医生", "en": "Dental Doctor"}, OrderNum: 150, IsActive: true},

	{ExternalID: "hc_nurse", CategoryExternalID: "healthcare_nurse", Names: map[string]string{"zh": "护士", "en": "Nurse"}, OrderNum: 10, IsActive: true},
	{ExternalID: "hc_clinical_nurse", CategoryExternalID: "healthcare_nurse", Names: map[string]string{"zh": "临床护士", "en": "Clinical Nurse"}, OrderNum: 20, IsActive: true},
	{ExternalID: "hc_or_nurse", CategoryExternalID: "healthcare_nurse", Names: map[string]string{"zh": "手术室护士", "en": "Operating Room Nurse"}, OrderNum: 30, IsActive: true},
	{ExternalID: "hc_internal_nurse", CategoryExternalID: "healthcare_nurse", Names: map[string]string{"zh": "内科护士", "en": "Internal Medicine Nurse"}, OrderNum: 40, IsActive: true},
	{ExternalID: "hc_obgyn_nurse", CategoryExternalID: "healthcare_nurse", Names: map[string]string{"zh": "妇产科护士", "en": "Obstetrician-Gynecologist Nurse"}, OrderNum: 50, IsActive: true},
	{ExternalID: "hc_midwife", CategoryExternalID: "healthcare_nurse", Names: map[string]string{"zh": "助产士", "en": "Midwife"}, OrderNum: 60, IsActive: true},

	{ExternalID: "hc_lab_tech", CategoryExternalID: "healthcare_medtech", Names: map[string]string{"zh": "医学检验技师", "en": "Laboratory Technician"}, OrderNum: 10, IsActive: true},
	{ExternalID: "hc_imaging_tech", CategoryExternalID: "healthcare_medtech", Names: map[string]string{"zh": "医学影像技师", "en": "Imaging Technician"}, OrderNum: 20, IsActive: true},
	{ExternalID: "hc_rehab_therapist", CategoryExternalID: "healthcare_medtech", Names: map[string]string{"zh": "康复治疗师", "en": "Rehabilitation Therapist"}, OrderNum: 30, IsActive: true},
	{ExternalID: "hc_dietitian", CategoryExternalID: "healthcare_medtech", Names: map[string]string{"zh": "营养师", "en": "Dietitian"}, OrderNum: 40, IsActive: true},
	{ExternalID: "hc_health_manager", CategoryExternalID: "healthcare_medtech", Names: map[string]string{"zh": "健康管理师", "en": "Health Manager"}, OrderNum: 50, IsActive: true},
	{ExternalID: "hc_psych_consultant", CategoryExternalID: "healthcare_medtech", Names: map[string]string{"zh": "心理咨询师", "en": "Psychologist"}, OrderNum: 60, IsActive: true},
	{ExternalID: "hc_acupuncture", CategoryExternalID: "healthcare_medtech", Names: map[string]string{"zh": "针灸推拿", "en": "Acupuncture"}, OrderNum: 70, IsActive: true},
	{ExternalID: "hc_lab_tester", CategoryExternalID: "healthcare_medtech", Names: map[string]string{"zh": "检验师", "en": "Laboratory Tester"}, OrderNum: 80, IsActive: true},
	{ExternalID: "hc_ultrasound_doctor", CategoryExternalID: "healthcare_medtech", Names: map[string]string{"zh": "超声科医师", "en": "Ultrasound Doctor"}, OrderNum: 90, IsActive: true},
	{ExternalID: "hc_pathologist", CategoryExternalID: "healthcare_medtech", Names: map[string]string{"zh": "病理科医师", "en": "Pathologist"}, OrderNum: 100, IsActive: true},

	{ExternalID: "hc_pharma_related", CategoryExternalID: "healthcare_pharma", Names: map[string]string{"zh": "药学相关", "en": "Pharmaceutical Related"}, OrderNum: 10, IsActive: true},
	{ExternalID: "hc_drug_rd", CategoryExternalID: "healthcare_pharma", Names: map[string]string{"zh": "药物研发", "en": "Drug Research and Development"}, OrderNum: 20, IsActive: true},
	{ExternalID: "hc_medicine_rd", CategoryExternalID: "healthcare_pharma", Names: map[string]string{"zh": "药品研发", "en": "Medicine Research and Development"}, OrderNum: 30, IsActive: true},
	{ExternalID: "hc_pharmacist", CategoryExternalID: "healthcare_pharma", Names: map[string]string{"zh": "药剂师/药师", "en": "Pharmacist"}, OrderNum: 40, IsActive: true},
	{ExternalID: "hc_med_qc", CategoryExternalID: "healthcare_pharma", Names: map[string]string{"zh": "医药质检", "en": "Medicine Quality Control"}, OrderNum: 50, IsActive: true},
	{ExternalID: "hc_drug_registration", CategoryExternalID: "healthcare_pharma", Names: map[string]string{"zh": "药品注册", "en": "Drug Registration"}, OrderNum: 60, IsActive: true},
	{ExternalID: "hc_cra", CategoryExternalID: "healthcare_pharma", Names: map[string]string{"zh": "临床监察员 (CRA)", "en": "Clinical Reviewer (CRA)"}, OrderNum: 70, IsActive: true},
	{ExternalID: "hc_crc", CategoryExternalID: "healthcare_pharma", Names: map[string]string{"zh": "临床协调员 (CRC)", "en": "Clinical Coordinator (CRC)"}, OrderNum: 80, IsActive: true},
	{ExternalID: "hc_drug_quality_mgmt", CategoryExternalID: "healthcare_pharma", Names: map[string]string{"zh": "药品质量管理", "en": "Drug Quality Management"}, OrderNum: 90, IsActive: true},

	{ExternalID: "hc_device_sales", CategoryExternalID: "healthcare_devices", Names: map[string]string{"zh": "医疗器械销售", "en": "Device Sales"}, OrderNum: 10, IsActive: true},
	{ExternalID: "hc_after_sales", CategoryExternalID: "healthcare_devices", Names: map[string]string{"zh": "售后工程师", "en": "After-Sales Engineer"}, OrderNum: 20, IsActive: true},
	{ExternalID: "hc_device_inspector", CategoryExternalID: "healthcare_devices", Names: map[string]string{"zh": "检验员", "en": "Device Inspector"}, OrderNum: 30, IsActive: true},
	{ExternalID: "hc_device_qc", CategoryExternalID: "healthcare_devices", Names: map[string]string{"zh": "设备质检", "en": "Device Quality Control"}, OrderNum: 40, IsActive: true},

	{ExternalID: "hc_records_admin", CategoryExternalID: "healthcare_other", Names: map[string]string{"zh": "病案管理员", "en": "Records Administrator"}, OrderNum: 10, IsActive: true},
	{ExternalID: "hc_registration_cashier", CategoryExternalID: "healthcare_other", Names: map[string]string{"zh": "挂号/收费员", "en": "Registration Cashier"}, OrderNum: 20, IsActive: true},
	{ExternalID: "hc_guide", CategoryExternalID: "healthcare_other", Names: map[string]string{"zh": "导医", "en": "Guide"}, OrderNum: 30, IsActive: true},
	{ExternalID: "hc_hospital_frontdesk", CategoryExternalID: "healthcare_other", Names: map[string]string{"zh": "医院前台", "en": "Hospital Front Desk"}, OrderNum: 40, IsActive: true},
	{ExternalID: "hc_inventory", CategoryExternalID: "healthcare_other", Names: map[string]string{"zh": "医疗库管", "en": "Inventory Manager"}, OrderNum: 50, IsActive: true},
	{ExternalID: "hc_security", CategoryExternalID: "healthcare_other", Names: map[string]string{"zh": "医院保安", "en": "Hospital Security"}, OrderNum: 60, IsActive: true},
	{ExternalID: "hc_cleaner", CategoryExternalID: "healthcare_other", Names: map[string]string{"zh": "医院保洁", "en": "Hospital Cleaner"}, OrderNum: 70, IsActive: true},
	{ExternalID: "hc_caregiver", CategoryExternalID: "healthcare_other", Names: map[string]string{"zh": "医院护工", "en": "Hospital Caregiver"}, OrderNum: 80, IsActive: true},
	{ExternalID: "hc_companion", CategoryExternalID: "healthcare_other", Names: map[string]string{"zh": "医院陪护", "en": "Hospital Companion"}, OrderNum: 90, IsActive: true},
	{ExternalID: "hc_med_translator", CategoryExternalID: "healthcare_other", Names: map[string]string{"zh": "医疗翻译", "en": "Medical Translator"}, OrderNum: 100, IsActive: true},
	{ExternalID: "hc_med_legal", CategoryExternalID: "healthcare_other", Names: map[string]string{"zh": "医疗法务", "en": "Medical Legal"}, OrderNum: 110, IsActive: true},
	{ExternalID: "hc_med_consulting", CategoryExternalID: "healthcare_other", Names: map[string]string{"zh": "医疗咨询", "en": "Medical Consulting"}, OrderNum: 120, IsActive: true},
	{ExternalID: "hc_med_training", CategoryExternalID: "healthcare_other", Names: map[string]string{"zh": "医疗培训", "en": "Medical Training"}, OrderNum: 130, IsActive: true},
	{ExternalID: "hc_med_social_worker", CategoryExternalID: "healthcare_other", Names: map[string]string{"zh": "医疗社工", "en": "Medical Social Worker"}, OrderNum: 140, IsActive: true},
	{ExternalID: "hc_med_rep", CategoryExternalID: "healthcare_other", Names: map[string]string{"zh": "医药代表", "en": "Medical Representative"}, OrderNum: 150, IsActive: true},

	{ExternalID: "hc_intern_doctor", CategoryExternalID: "healthcare_intern", Names: map[string]string{"zh": "实习医生", "en": "Intern Doctor"}, OrderNum: 10, IsActive: true},
	{ExternalID: "hc_intern_nurse", CategoryExternalID: "healthcare_intern", Names: map[string]string{"zh": "实习护士", "en": "Intern Nurse"}, OrderNum: 20, IsActive: true},
	{ExternalID: "hc_intern_pharmacist", CategoryExternalID: "healthcare_intern", Names: map[string]string{"zh": "实习药师", "en": "Intern Pharmacist"}, OrderNum: 30, IsActive: true},
	{ExternalID: "hc_intern_tech", CategoryExternalID: "healthcare_intern", Names: map[string]string{"zh": "实习技师", "en": "Intern Technician"}, OrderNum: 40, IsActive: true},
	{ExternalID: "hc_med_intern", CategoryExternalID: "healthcare_intern", Names: map[string]string{"zh": "医学实习生", "en": "Medical Intern"}, OrderNum: 50, IsActive: true},
	{ExternalID: "hc_intern_med_rep", CategoryExternalID: "healthcare_intern", Names: map[string]string{"zh": "医药代表实习生", "en": "Medical Representative Intern"}, OrderNum: 60, IsActive: true},
	{ExternalID: "hc_clinical_research_intern", CategoryExternalID: "healthcare_intern", Names: map[string]string{"zh": "临床研究实习生", "en": "Clinical Research Intern"}, OrderNum: 70, IsActive: true},

	{ExternalID: "hc_hospital_director", CategoryExternalID: "healthcare_management", Names: map[string]string{"zh": "医院院长", "en": "Hospital Director"}, OrderNum: 10, IsActive: true},
	{ExternalID: "hc_department_head", CategoryExternalID: "healthcare_management", Names: map[string]string{"zh": "科室主任", "en": "Department Head"}, OrderNum: 20, IsActive: true},
	{ExternalID: "hc_head_nurse", CategoryExternalID: "healthcare_management", Names: map[string]string{"zh": "护士长", "en": "Head Nurse"}, OrderNum: 30, IsActive: true},

	{ExternalID: "re_structural_engineer", CategoryExternalID: "realestate_design_planning", Names: map[string]string{"zh": "结构工程师", "en": "Structural Engineer"}, OrderNum: 10, IsActive: true},
	{ExternalID: "re_urban_planner", CategoryExternalID: "realestate_design_planning", Names: map[string]string{"zh": "城市规划师", "en": "Urban Planner"}, OrderNum: 20, IsActive: true},
	{ExternalID: "re_landscape_designer", CategoryExternalID: "realestate_design_planning", Names: map[string]string{"zh": "园林设计师", "en": "Landscape Designer"}, OrderNum: 30, IsActive: true},
	{ExternalID: "re_planning_design", CategoryExternalID: "realestate_design_planning", Names: map[string]string{"zh": "规划设计", "en": "Planning Design"}, OrderNum: 40, IsActive: true},
	{ExternalID: "re_bim_engineer", CategoryExternalID: "realestate_design_planning", Names: map[string]string{"zh": "BIM工程师", "en": "BIM Engineer"}, OrderNum: 50, IsActive: true},
	{ExternalID: "re_landscape_constructor", CategoryExternalID: "realestate_design_planning", Names: map[string]string{"zh": "园林施工员", "en": "Landscape Constructor"}, OrderNum: 60, IsActive: true},
	{ExternalID: "re_construction_engineer", CategoryExternalID: "realestate_design_planning", Names: map[string]string{"zh": "建筑工程师", "en": "Construction Engineer"}, OrderNum: 70, IsActive: true},

	{ExternalID: "re_landscape_architect", CategoryExternalID: "realestate_interior_landscape", Names: map[string]string{"zh": "景观设计师", "en": "Landscape Architect"}, OrderNum: 10, IsActive: true},
	{ExternalID: "re_home_designer", CategoryExternalID: "realestate_interior_landscape", Names: map[string]string{"zh": "家装设计师", "en": "Home Designer"}, OrderNum: 20, IsActive: true},

	{ExternalID: "re_budgeter", CategoryExternalID: "realestate_cost_budget", Names: map[string]string{"zh": "预算员", "en": "Budgeter"}, OrderNum: 10, IsActive: true},
	{ExternalID: "re_cost_estimator", CategoryExternalID: "realestate_cost_budget", Names: map[string]string{"zh": "造价员", "en": "Cost Estimator"}, OrderNum: 20, IsActive: true},
	{ExternalID: "re_cost_engineer", CategoryExternalID: "realestate_cost_budget", Names: map[string]string{"zh": "造价工程师", "en": "Cost Engineer"}, OrderNum: 30, IsActive: true},
	{ExternalID: "re_project_cost", CategoryExternalID: "realestate_cost_budget", Names: map[string]string{"zh": "工程造价", "en": "Project Cost"}, OrderNum: 40, IsActive: true},

	{ExternalID: "re_site_worker", CategoryExternalID: "realestate_construction_mgmt", Names: map[string]string{"zh": "施工员", "en": "Site Worker"}, OrderNum: 10, IsActive: true},
	{ExternalID: "re_civil_site_worker", CategoryExternalID: "realestate_construction_mgmt", Names: map[string]string{"zh": "土建施工员", "en": "Civil Site Worker"}, OrderNum: 20, IsActive: true},
	{ExternalID: "re_surveyor", CategoryExternalID: "realestate_construction_mgmt", Names: map[string]string{"zh": "测量员", "en": "Surveyor"}, OrderNum: 30, IsActive: true},
	{ExternalID: "re_mapping_engineer", CategoryExternalID: "realestate_construction_mgmt", Names: map[string]string{"zh": "测绘工程师", "en": "Mapping Engineer"}, OrderNum: 40, IsActive: true},
	{ExternalID: "re_supervisor", CategoryExternalID: "realestate_construction_mgmt", Names: map[string]string{"zh": "工程监理", "en": "Supervisor"}, OrderNum: 50, IsActive: true},
	{ExternalID: "re_supervision_engineer", CategoryExternalID: "realestate_construction_mgmt", Names: map[string]string{"zh": "监理工程师", "en": "Supervision Engineer"}, OrderNum: 60, IsActive: true},
	{ExternalID: "re_project_admin", CategoryExternalID: "realestate_construction_mgmt", Names: map[string]string{"zh": "工程管理员", "en": "Project Administrator"}, OrderNum: 70, IsActive: true},
	{ExternalID: "re_document_controller", CategoryExternalID: "realestate_construction_mgmt", Names: map[string]string{"zh": "资料员", "en": "Document Controller"}, OrderNum: 80, IsActive: true},
	{ExternalID: "re_archive_admin", CategoryExternalID: "realestate_construction_mgmt", Names: map[string]string{"zh": "档案管理员", "en": "Archive Administrator"}, OrderNum: 90, IsActive: true},
	{ExternalID: "re_engineering_manager", CategoryExternalID: "realestate_construction_mgmt", Names: map[string]string{"zh": "工程经理", "en": "Engineering Manager"}, OrderNum: 100, IsActive: true},
	{ExternalID: "re_civil_engineer", CategoryExternalID: "realestate_construction_mgmt", Names: map[string]string{"zh": "土木工程师", "en": "Civil Engineer"}, OrderNum: 110, IsActive: true},
	{ExternalID: "re_safety_officer", CategoryExternalID: "realestate_construction_mgmt", Names: map[string]string{"zh": "安全员", "en": "Safety Officer"}, OrderNum: 120, IsActive: true},
	{ExternalID: "re_quality_inspector", CategoryExternalID: "realestate_construction_mgmt", Names: map[string]string{"zh": "工程质检员", "en": "Quality Inspector"}, OrderNum: 130, IsActive: true},
	{ExternalID: "re_project_engineer", CategoryExternalID: "realestate_construction_mgmt", Names: map[string]string{"zh": "项目工程师", "en": "Project Engineer"}, OrderNum: 140, IsActive: true},
	{ExternalID: "re_civil_project_engineer", CategoryExternalID: "realestate_construction_mgmt", Names: map[string]string{"zh": "土建工程师", "en": "Civil Project Engineer"}, OrderNum: 150, IsActive: true},

	{ExternalID: "re_pm", CategoryExternalID: "realestate_project_mgmt", Names: map[string]string{"zh": "项目经理", "en": "Project Manager"}, OrderNum: 10, IsActive: true},
	{ExternalID: "re_pm_assistant", CategoryExternalID: "realestate_project_mgmt", Names: map[string]string{"zh": "项目助理", "en": "Project Assistant"}, OrderNum: 20, IsActive: true},
	{ExternalID: "re_pm_specialist", CategoryExternalID: "realestate_project_mgmt", Names: map[string]string{"zh": "项目专员", "en": "Project Specialist"}, OrderNum: 30, IsActive: true},
	{ExternalID: "re_pm_supervisor", CategoryExternalID: "realestate_project_mgmt", Names: map[string]string{"zh": "项目主管", "en": "Project Supervisor"}, OrderNum: 40, IsActive: true},
	{ExternalID: "re_pm_director", CategoryExternalID: "realestate_project_mgmt", Names: map[string]string{"zh": "项目总监", "en": "Project Director"}, OrderNum: 50, IsActive: true},
	{ExternalID: "re_bid_specialist", CategoryExternalID: "realestate_project_mgmt", Names: map[string]string{"zh": "投标专员", "en": "Bid Specialist"}, OrderNum: 60, IsActive: true},

	{ExternalID: "re_sales", CategoryExternalID: "realestate_sales_planning", Names: map[string]string{"zh": "房地产销售", "en": "Real Estate Sales"}, OrderNum: 10, IsActive: true},
	{ExternalID: "re_property_consultant", CategoryExternalID: "realestate_sales_planning", Names: map[string]string{"zh": "置业顾问", "en": "Property Consultant"}, OrderNum: 20, IsActive: true},
	{ExternalID: "re_agent", CategoryExternalID: "realestate_sales_planning", Names: map[string]string{"zh": "房产经纪人", "en": "Property Agent"}, OrderNum: 30, IsActive: true},
	{ExternalID: "re_marketing_planner", CategoryExternalID: "realestate_sales_planning", Names: map[string]string{"zh": "房地产策划", "en": "Real Estate Marketing Planner"}, OrderNum: 40, IsActive: true},
	{ExternalID: "re_leasing_manager", CategoryExternalID: "realestate_sales_planning", Names: map[string]string{"zh": "招商经理", "en": "Leasing Manager"}, OrderNum: 50, IsActive: true},
	{ExternalID: "re_channel_manager", CategoryExternalID: "realestate_sales_planning", Names: map[string]string{"zh": "渠道经理", "en": "Channel Manager"}, OrderNum: 60, IsActive: true},
	{ExternalID: "re_other_roles", CategoryExternalID: "realestate_sales_planning", Names: map[string]string{"zh": "房产其他岗位", "en": "Other Roles"}, OrderNum: 70, IsActive: true},

	{ExternalID: "re_property_mgmt", CategoryExternalID: "realestate_property_mgmt", Names: map[string]string{"zh": "物业管理", "en": "Property Management"}, OrderNum: 10, IsActive: true},
	{ExternalID: "re_property_manager", CategoryExternalID: "realestate_property_mgmt", Names: map[string]string{"zh": "物业经理", "en": "Property Manager"}, OrderNum: 20, IsActive: true},
	{ExternalID: "re_property_cs", CategoryExternalID: "realestate_property_mgmt", Names: map[string]string{"zh": "物业客服", "en": "Property Customer Service"}, OrderNum: 30, IsActive: true},
	{ExternalID: "re_property_steward", CategoryExternalID: "realestate_property_mgmt", Names: map[string]string{"zh": "物业管家", "en": "Property Steward"}, OrderNum: 40, IsActive: true},

	{ExternalID: "mfg_mech_engineer", CategoryExternalID: "mfg_mechanical", Names: map[string]string{"zh": "机械工程师", "en": "Mechanical Engineer"}, OrderNum: 10, IsActive: true},
	{ExternalID: "mfg_mech_design", CategoryExternalID: "mfg_mechanical", Names: map[string]string{"zh": "机械设计工程师", "en": "Mechanical Design Engineer"}, OrderNum: 20, IsActive: true},
	{ExternalID: "mfg_mechatronics", CategoryExternalID: "mfg_mechanical", Names: map[string]string{"zh": "机电工程师", "en": "Mechatronics Engineer"}, OrderNum: 30, IsActive: true},
	{ExternalID: "mfg_equipment_engineer", CategoryExternalID: "mfg_mechanical", Names: map[string]string{"zh": "设备工程师", "en": "Equipment Engineer"}, OrderNum: 40, IsActive: true},
	{ExternalID: "mfg_mechanical_manufacturing", CategoryExternalID: "mfg_mechanical", Names: map[string]string{"zh": "机械制造", "en": "Mechanical Manufacturing"}, OrderNum: 50, IsActive: true},
	{ExternalID: "mfg_mechanical_maintenance", CategoryExternalID: "mfg_mechanical", Names: map[string]string{"zh": "机械维修", "en": "Mechanical Maintenance"}, OrderNum: 60, IsActive: true},
	{ExternalID: "mfg_hydraulics", CategoryExternalID: "mfg_mechanical", Names: map[string]string{"zh": "液压工程师", "en": "Hydraulics Engineer"}, OrderNum: 70, IsActive: true},
	{ExternalID: "mfg_nc_programmer", CategoryExternalID: "mfg_mechanical", Names: map[string]string{"zh": "数控编程", "en": "NC Programmer"}, OrderNum: 80, IsActive: true},
	{ExternalID: "mfg_other_engineers", CategoryExternalID: "mfg_mechanical", Names: map[string]string{"zh": "其他工程师岗位", "en": "Other Engineers"}, OrderNum: 90, IsActive: true},

	{ExternalID: "mfg_elec_engineer", CategoryExternalID: "mfg_electrical", Names: map[string]string{"zh": "电子工程师", "en": "Electrical Engineer"}, OrderNum: 10, IsActive: true},
	{ExternalID: "mfg_electrical_engineer", CategoryExternalID: "mfg_electrical", Names: map[string]string{"zh": "电气工程师", "en": "Electrical Engineer"}, OrderNum: 20, IsActive: true},
	{ExternalID: "mfg_hw_engineer", CategoryExternalID: "mfg_electrical", Names: map[string]string{"zh": "硬件工程师", "en": "Hardware Engineer"}, OrderNum: 30, IsActive: true},
	{ExternalID: "mfg_industrial_automation", CategoryExternalID: "mfg_electrical", Names: map[string]string{"zh": "电气自动化工程师", "en": "Industrial Automation Engineer"}, OrderNum: 40, IsActive: true},
	{ExternalID: "mfg_embedded", CategoryExternalID: "mfg_electrical", Names: map[string]string{"zh": "嵌入式工程师", "en": "Embedded Engineer"}, OrderNum: 50, IsActive: true},
	{ExternalID: "mfg_automation", CategoryExternalID: "mfg_electrical", Names: map[string]string{"zh": "自动化工程师", "en": "Automation Engineer"}, OrderNum: 60, IsActive: true},
	{ExternalID: "mfg_semiconductor_tech", CategoryExternalID: "mfg_electrical", Names: map[string]string{"zh": "半导体技术员", "en": "Semiconductor Technician"}, OrderNum: 70, IsActive: true},
	{ExternalID: "mfg_circuit_design", CategoryExternalID: "mfg_electrical", Names: map[string]string{"zh": "电路设计", "en": "Circuit Design"}, OrderNum: 80, IsActive: true},

	{ExternalID: "mfg_auto_engineer", CategoryExternalID: "mfg_auto_transport", Names: map[string]string{"zh": "汽车工程师", "en": "Automobile Engineer"}, OrderNum: 10, IsActive: true},
	{ExternalID: "mfg_vehicle_engineer", CategoryExternalID: "mfg_auto_transport", Names: map[string]string{"zh": "车辆工程师", "en": "Vehicle Engineer"}, OrderNum: 20, IsActive: true},
	{ExternalID: "mfg_vehicle_design", CategoryExternalID: "mfg_auto_transport", Names: map[string]string{"zh": "汽车设计", "en": "Vehicle Design"}, OrderNum: 30, IsActive: true},
	{ExternalID: "mfg_powertrain", CategoryExternalID: "mfg_auto_transport", Names: map[string]string{"zh": "动力总成工程师", "en": "Powertrain Engineer"}, OrderNum: 40, IsActive: true},

	{ExternalID: "mfg_process_engineer", CategoryExternalID: "mfg_process_mold", Names: map[string]string{"zh": "工艺工程师", "en": "Process Engineer"}, OrderNum: 10, IsActive: true},
	{ExternalID: "mfg_mould_engineer", CategoryExternalID: "mfg_process_mold", Names: map[string]string{"zh": "模具工程师", "en": "Mould Engineer"}, OrderNum: 20, IsActive: true},
	{ExternalID: "mfg_welding_engineer", CategoryExternalID: "mfg_process_mold", Names: map[string]string{"zh": "焊接工程师", "en": "Welding Engineer"}, OrderNum: 30, IsActive: true},
	{ExternalID: "mfg_mould_design", CategoryExternalID: "mfg_process_mold", Names: map[string]string{"zh": "模具设计师", "en": "Mould Designer"}, OrderNum: 40, IsActive: true},
	{ExternalID: "mfg_stamping_process", CategoryExternalID: "mfg_process_mold", Names: map[string]string{"zh": "冲压工艺师/模具设计师", "en": "Stamping Process/ Mould Designer"}, OrderNum: 50, IsActive: true},

	{ExternalID: "mfg_prod_management", CategoryExternalID: "mfg_prod_equip", Names: map[string]string{"zh": "生产管理", "en": "Production Management"}, OrderNum: 10, IsActive: true},
	{ExternalID: "mfg_prod_supervisor", CategoryExternalID: "mfg_prod_equip", Names: map[string]string{"zh": "生产主管", "en": "Production Supervisor"}, OrderNum: 20, IsActive: true},
	{ExternalID: "mfg_production_manager", CategoryExternalID: "mfg_prod_equip", Names: map[string]string{"zh": "生产经理", "en": "Production Manager"}, OrderNum: 30, IsActive: true},
	{ExternalID: "mfg_line_leader", CategoryExternalID: "mfg_prod_equip", Names: map[string]string{"zh": "车间主任", "en": "Line Leader"}, OrderNum: 40, IsActive: true},
	{ExternalID: "mfg_equipment_maintenance", CategoryExternalID: "mfg_prod_equip", Names: map[string]string{"zh": "设备维护", "en": "Equipment Maintenance"}, OrderNum: 50, IsActive: true},
	{ExternalID: "mfg_equipment_manager", CategoryExternalID: "mfg_prod_equip", Names: map[string]string{"zh": "设备管理", "en": "Equipment Manager"}, OrderNum: 60, IsActive: true},
	{ExternalID: "mfg_shift_leader", CategoryExternalID: "mfg_prod_equip", Names: map[string]string{"zh": "生产班长", "en": "Shift Leader"}, OrderNum: 70, IsActive: true},
	{ExternalID: "mfg_team_leader", CategoryExternalID: "mfg_prod_equip", Names: map[string]string{"zh": "工段长", "en": "Team Leader"}, OrderNum: 80, IsActive: true},
	{ExternalID: "mfg_group_leader", CategoryExternalID: "mfg_prod_equip", Names: map[string]string{"zh": "班组长", "en": "Group Leader"}, OrderNum: 90, IsActive: true},
	{ExternalID: "mfg_factory_manager", CategoryExternalID: "mfg_prod_equip", Names: map[string]string{"zh": "厂长", "en": "Factory Manager"}, OrderNum: 100, IsActive: true},
	{ExternalID: "mfg_production_planner", CategoryExternalID: "mfg_prod_equip", Names: map[string]string{"zh": "生产计划员", "en": "Production Planner"}, OrderNum: 110, IsActive: true},

	{ExternalID: "mfg_quality_engineer", CategoryExternalID: "mfg_quality", Names: map[string]string{"zh": "质量工程师", "en": "Quality Engineer"}, OrderNum: 10, IsActive: true},
	{ExternalID: "mfg_quality_specialist", CategoryExternalID: "mfg_quality", Names: map[string]string{"zh": "品质工程师", "en": "Quality Specialist"}, OrderNum: 20, IsActive: true},
	{ExternalID: "mfg_qc", CategoryExternalID: "mfg_quality", Names: map[string]string{"zh": "QC", "en": "QC"}, OrderNum: 30, IsActive: true},
	{ExternalID: "mfg_quality_manager", CategoryExternalID: "mfg_quality", Names: map[string]string{"zh": "质量管理工程师", "en": "Quality Management Engineer"}, OrderNum: 40, IsActive: true},
	{ExternalID: "mfg_quality_inspector", CategoryExternalID: "mfg_quality", Names: map[string]string{"zh": "质量检测员", "en": "Quality Inspector"}, OrderNum: 50, IsActive: true},

	{ExternalID: "mfg_rd_engineer", CategoryExternalID: "mfg_rd_design", Names: map[string]string{"zh": "研发工程师", "en": "Research and Development Engineer"}, OrderNum: 10, IsActive: true},
	{ExternalID: "mfg_tech_engineer", CategoryExternalID: "mfg_rd_design", Names: map[string]string{"zh": "技术员", "en": "Technician"}, OrderNum: 20, IsActive: true},
	{ExternalID: "mfg_design_engineer", CategoryExternalID: "mfg_rd_design", Names: map[string]string{"zh": "设计工程师", "en": "Design Engineer"}, OrderNum: 30, IsActive: true},
	{ExternalID: "mfg_product_designer", CategoryExternalID: "mfg_rd_design", Names: map[string]string{"zh": "产品设计师", "en": "Product Designer"}, OrderNum: 40, IsActive: true},

	{ExternalID: "mfg_operator", CategoryExternalID: "mfg_worker", Names: map[string]string{"zh": "操作工", "en": "Operator"}, OrderNum: 10, IsActive: true},
	{ExternalID: "mfg_general_worker", CategoryExternalID: "mfg_worker", Names: map[string]string{"zh": "普工", "en": "General Worker"}, OrderNum: 20, IsActive: true},
	{ExternalID: "mfg_assembler", CategoryExternalID: "mfg_worker", Names: map[string]string{"zh": "装配工", "en": "Assembler"}, OrderNum: 30, IsActive: true},
	{ExternalID: "mfg_technician", CategoryExternalID: "mfg_worker", Names: map[string]string{"zh": "技术工人", "en": "Technician"}, OrderNum: 40, IsActive: true},
	{ExternalID: "mfg_fitter", CategoryExternalID: "mfg_worker", Names: map[string]string{"zh": "钳工", "en": "Fitter"}, OrderNum: 50, IsActive: true},
	{ExternalID: "mfg_welder", CategoryExternalID: "mfg_worker", Names: map[string]string{"zh": "焊工", "en": "Welder"}, OrderNum: 60, IsActive: true},
	{ExternalID: "mfg_driller", CategoryExternalID: "mfg_worker", Names: map[string]string{"zh": "钻工", "en": "Driller"}, OrderNum: 70, IsActive: true},
	{ExternalID: "mfg_turner", CategoryExternalID: "mfg_worker", Names: map[string]string{"zh": "车工", "en": "Turner"}, OrderNum: 80, IsActive: true},
	{ExternalID: "mfg_miller", CategoryExternalID: "mfg_worker", Names: map[string]string{"zh": "铣工", "en": "Millwright"}, OrderNum: 90, IsActive: true},
	{ExternalID: "mfg_cnc_operator", CategoryExternalID: "mfg_worker", Names: map[string]string{"zh": "CNC操作工", "en": "CNC Operator"}, OrderNum: 100, IsActive: true},
	{ExternalID: "mfg_machine_operator", CategoryExternalID: "mfg_worker", Names: map[string]string{"zh": "机床操作工", "en": "Machine Operator"}, OrderNum: 110, IsActive: true},
	{ExternalID: "mfg_packer", CategoryExternalID: "mfg_worker", Names: map[string]string{"zh": "包装工", "en": "Packer"}, OrderNum: 120, IsActive: true},

	{ExternalID: "logi_captain", CategoryExternalID: "logi_transport_service", Names: map[string]string{"zh": "机长", "en": "Captain"}, OrderNum: 10, IsActive: true},
	{ExternalID: "logi_pilot", CategoryExternalID: "logi_transport_service", Names: map[string]string{"zh": "飞行员", "en": "Pilot"}, OrderNum: 20, IsActive: true},
	{ExternalID: "logi_air_security", CategoryExternalID: "logi_transport_service", Names: map[string]string{"zh": "空中安全员", "en": "Air Security"}, OrderNum: 30, IsActive: true},
	{ExternalID: "logi_flight_attendant", CategoryExternalID: "logi_transport_service", Names: map[string]string{"zh": "空姐", "en": "Flight Attendant"}, OrderNum: 40, IsActive: true},
	{ExternalID: "logi_male_attendant", CategoryExternalID: "logi_transport_service", Names: map[string]string{"zh": "空少", "en": "Male Attendant"}, OrderNum: 50, IsActive: true},
	{ExternalID: "logi_ground", CategoryExternalID: "logi_transport_service", Names: map[string]string{"zh": "地勤", "en": "Ground Staff"}, OrderNum: 60, IsActive: true},
	{ExternalID: "logi_tickets", CategoryExternalID: "logi_transport_service", Names: map[string]string{"zh": "票务员", "en": "Ticketer"}, OrderNum: 70, IsActive: true},
	{ExternalID: "logi_security", CategoryExternalID: "logi_transport_service", Names: map[string]string{"zh": "安检员", "en": "Security"}, OrderNum: 80, IsActive: true},
	{ExternalID: "logi_metro_staff", CategoryExternalID: "logi_transport_service", Names: map[string]string{"zh": "地铁站务员", "en": "Metro Staff"}, OrderNum: 90, IsActive: true},
	{ExternalID: "logi_metro_security", CategoryExternalID: "logi_transport_service", Names: map[string]string{"zh": "地铁安检员", "en": "Metro Security"}, OrderNum: 100, IsActive: true},
	{ExternalID: "logi_metro_driver", CategoryExternalID: "logi_transport_service", Names: map[string]string{"zh": "地铁驾驶员", "en": "Metro Driver"}, OrderNum: 110, IsActive: true},
	{ExternalID: "logi_highspeed_attendant", CategoryExternalID: "logi_transport_service", Names: map[string]string{"zh": "高铁乘务", "en": "Highspeed Attendant"}, OrderNum: 120, IsActive: true},
	{ExternalID: "logi_call_center", CategoryExternalID: "logi_transport_service", Names: map[string]string{"zh": "话务员", "en": "Call Center"}, OrderNum: 130, IsActive: true},
	{ExternalID: "logi_flight_dispatcher", CategoryExternalID: "logi_transport_service", Names: map[string]string{"zh": "签派员", "en": "Flight Dispatcher"}, OrderNum: 140, IsActive: true},

	{ExternalID: "logi_driver", CategoryExternalID: "logi_delivery", Names: map[string]string{"zh": "司机", "en": "Driver"}, OrderNum: 10, IsActive: true},
	{ExternalID: "logi_dispatch_manager", CategoryExternalID: "logi_delivery", Names: map[string]string{"zh": "配送经理", "en": "Dispatch Manager"}, OrderNum: 20, IsActive: true},
	{ExternalID: "logi_transport_admin", CategoryExternalID: "logi_delivery", Names: map[string]string{"zh": "运输主管", "en": "Transport Admin"}, OrderNum: 30, IsActive: true},
	{ExternalID: "logi_fleet_manager", CategoryExternalID: "logi_delivery", Names: map[string]string{"zh": "车队管理", "en": "Fleet Manager"}, OrderNum: 40, IsActive: true},
	{ExternalID: "logi_scheduler", CategoryExternalID: "logi_delivery", Names: map[string]string{"zh": "调度员", "en": "Scheduler"}, OrderNum: 50, IsActive: true},

	{ExternalID: "logi_warehouse_admin", CategoryExternalID: "logi_warehouse", Names: map[string]string{"zh": "仓库管理员/库管", "en": "Warehouse Admin"}, OrderNum: 10, IsActive: true},
	{ExternalID: "logi_storekeeper", CategoryExternalID: "logi_warehouse", Names: map[string]string{"zh": "仓管员", "en": "Storekeeper"}, OrderNum: 20, IsActive: true},
	{ExternalID: "logi_warehouse_manager", CategoryExternalID: "logi_warehouse", Names: map[string]string{"zh": "仓储管理", "en": "Warehouse Manager"}, OrderNum: 30, IsActive: true},
	{ExternalID: "logi_loader", CategoryExternalID: "logi_warehouse", Names: map[string]string{"zh": "装卸工", "en": "Loader"}, OrderNum: 40, IsActive: true},
	{ExternalID: "logi_packer", CategoryExternalID: "logi_warehouse", Names: map[string]string{"zh": "包装员", "en": "Packer"}, OrderNum: 50, IsActive: true},
	{ExternalID: "logi_warehouse_specialist", CategoryExternalID: "logi_warehouse", Names: map[string]string{"zh": "仓储专员", "en": "Warehouse Specialist"}, OrderNum: 60, IsActive: true},
	{ExternalID: "logi_warehouse_supervisor", CategoryExternalID: "logi_warehouse", Names: map[string]string{"zh": "仓储主管", "en": "Warehouse Supervisor"}, OrderNum: 70, IsActive: true},
	{ExternalID: "logi_warehouse_manager_role", CategoryExternalID: "logi_warehouse", Names: map[string]string{"zh": "仓储经理", "en": "Warehouse Manager"}, OrderNum: 80, IsActive: true},
	{ExternalID: "logi_warehouse_reserve", CategoryExternalID: "logi_warehouse", Names: map[string]string{"zh": "储备干部", "en": "Warehouse Reserve"}, OrderNum: 90, IsActive: true},
	{ExternalID: "logi_warehouse_superintendent", CategoryExternalID: "logi_warehouse", Names: map[string]string{"zh": "储备主管", "en": "Warehouse Superintendent"}, OrderNum: 100, IsActive: true},

	{ExternalID: "logi_ops_specialist", CategoryExternalID: "logi_ops", Names: map[string]string{"zh": "物流专员/助理", "en": "Ops Specialist"}, OrderNum: 10, IsActive: true},
	{ExternalID: "logi_ops_supervisor", CategoryExternalID: "logi_ops", Names: map[string]string{"zh": "物流主管", "en": "Ops Supervisor"}, OrderNum: 20, IsActive: true},
	{ExternalID: "logi_ops_manager", CategoryExternalID: "logi_ops", Names: map[string]string{"zh": "物流经理", "en": "Ops Manager"}, OrderNum: 30, IsActive: true},
	{ExternalID: "logi_order_follower", CategoryExternalID: "logi_ops", Names: map[string]string{"zh": "跟单员", "en": "Order Follower"}, OrderNum: 40, IsActive: true},
	{ExternalID: "logi_courier", CategoryExternalID: "logi_ops", Names: map[string]string{"zh": "物流员", "en": "Courier"}, OrderNum: 50, IsActive: true},
	{ExternalID: "logi_express", CategoryExternalID: "logi_ops", Names: map[string]string{"zh": "快递员", "en": "Express"}, OrderNum: 60, IsActive: true},
	{ExternalID: "logi_delivery_staff", CategoryExternalID: "logi_ops", Names: map[string]string{"zh": "配送员", "en": "Delivery Staff"}, OrderNum: 70, IsActive: true},

	{ExternalID: "logi_supply_specialist", CategoryExternalID: "logi_supply", Names: map[string]string{"zh": "供应链专员", "en": "Supply Specialist"}, OrderNum: 10, IsActive: true},
	{ExternalID: "logi_supply_manager", CategoryExternalID: "logi_supply", Names: map[string]string{"zh": "供应链管理", "en": "Supply Manager"}, OrderNum: 20, IsActive: true},
	{ExternalID: "logi_supply_director", CategoryExternalID: "logi_supply", Names: map[string]string{"zh": "供应链经理", "en": "Supply Director"}, OrderNum: 30, IsActive: true},
	{ExternalID: "logi_supply_supervisor", CategoryExternalID: "logi_supply", Names: map[string]string{"zh": "供应链总监", "en": "Supply Supervisor"}, OrderNum: 40, IsActive: true},
	{ExternalID: "logi_supply_intern", CategoryExternalID: "logi_supply", Names: map[string]string{"zh": "供应链实习生", "en": "Supply Intern"}, OrderNum: 50, IsActive: true},

	{ExternalID: "svc_supermarket_manager", CategoryExternalID: "svc_food", Names: map[string]string{"zh": "超市店长", "en": "Supermarket Manager"}, OrderNum: 10, IsActive: true},
	{ExternalID: "svc_cashier_staff", CategoryExternalID: "svc_food", Names: map[string]string{"zh": "收银员", "en": "Cashier Staff"}, OrderNum: 20, IsActive: true},
	{ExternalID: "svc_waiter", CategoryExternalID: "svc_food", Names: map[string]string{"zh": "服务员", "en": "Waiter"}, OrderNum: 30, IsActive: true},
	{ExternalID: "svc_barista", CategoryExternalID: "svc_food", Names: map[string]string{"zh": "咖啡师", "en": "Barista"}, OrderNum: 40, IsActive: true},
	{ExternalID: "svc_chef", CategoryExternalID: "svc_food", Names: map[string]string{"zh": "厨师", "en": "Chef"}, OrderNum: 50, IsActive: true},
	{ExternalID: "svc_baker", CategoryExternalID: "svc_food", Names: map[string]string{"zh": "面包师", "en": "Baker"}, OrderNum: 60, IsActive: true},
	{ExternalID: "svc_western_chef", CategoryExternalID: "svc_food", Names: map[string]string{"zh": "西点师", "en": "Western Chef"}, OrderNum: 70, IsActive: true},
	{ExternalID: "svc_food_manager", CategoryExternalID: "svc_food", Names: map[string]string{"zh": "餐饮管理", "en": "Food Manager"}, OrderNum: 80, IsActive: true},
	{ExternalID: "svc_food_store_manager", CategoryExternalID: "svc_food", Names: map[string]string{"zh": "餐饮店长", "en": "Food Store Manager"}, OrderNum: 90, IsActive: true},
	{ExternalID: "svc_back_kitchen", CategoryExternalID: "svc_food", Names: map[string]string{"zh": "后厨", "en": "Back Kitchen"}, OrderNum: 100, IsActive: true},
	{ExternalID: "svc_restaurant_manager", CategoryExternalID: "svc_food", Names: map[string]string{"zh": "餐厅经理", "en": "Restaurant Manager"}, OrderNum: 110, IsActive: true},
	{ExternalID: "svc_pastry_chef", CategoryExternalID: "svc_food", Names: map[string]string{"zh": "面点师", "en": "Pastry Chef"}, OrderNum: 120, IsActive: true},
	{ExternalID: "svc_head_chef", CategoryExternalID: "svc_food", Names: map[string]string{"zh": "厨师长", "en": "Head Chef"}, OrderNum: 130, IsActive: true},
	{ExternalID: "svc_baking", CategoryExternalID: "svc_food", Names: map[string]string{"zh": "烘焙师", "en": "Baking"}, OrderNum: 140, IsActive: true},
	{ExternalID: "svc_bartender", CategoryExternalID: "svc_food", Names: map[string]string{"zh": "调酒师", "en": "Bartender"}, OrderNum: 150, IsActive: true},
	{ExternalID: "svc_tea_artist", CategoryExternalID: "svc_food", Names: map[string]string{"zh": "茶艺师", "en": "Tea Artist"}, OrderNum: 160, IsActive: true},
	{ExternalID: "svc_food_director", CategoryExternalID: "svc_food", Names: map[string]string{"zh": "餐饮总监", "en": "Food Director"}, OrderNum: 170, IsActive: true},

	{ExternalID: "svc_tour_guide", CategoryExternalID: "svc_hotel_travel", Names: map[string]string{"zh": "导游", "en": "Tour Guide"}, OrderNum: 10, IsActive: true},
	{ExternalID: "svc_travel_consultant", CategoryExternalID: "svc_hotel_travel", Names: map[string]string{"zh": "旅游顾问", "en": "Travel Consultant"}, OrderNum: 20, IsActive: true},
	{ExternalID: "svc_planner", CategoryExternalID: "svc_hotel_travel", Names: map[string]string{"zh": "计调", "en": "Planner"}, OrderNum: 30, IsActive: true},
	{ExternalID: "svc_hotel_manager", CategoryExternalID: "svc_hotel_travel", Names: map[string]string{"zh": "酒店管理", "en": "Hotel Manager"}, OrderNum: 40, IsActive: true},
	{ExternalID: "svc_frontdesk", CategoryExternalID: "svc_hotel_travel", Names: map[string]string{"zh": "酒店前台", "en": "Frontdesk"}, OrderNum: 50, IsActive: true},
	{ExternalID: "svc_guest_service", CategoryExternalID: "svc_hotel_travel", Names: map[string]string{"zh": "客房服务", "en": "Guest Service"}, OrderNum: 60, IsActive: true},
	{ExternalID: "svc_porter", CategoryExternalID: "svc_hotel_travel", Names: map[string]string{"zh": "行李员", "en": "Porter"}, OrderNum: 70, IsActive: true},
	{ExternalID: "svc_travel_interpreter", CategoryExternalID: "svc_hotel_travel", Names: map[string]string{"zh": "景区讲解员", "en": "Travel Interpreter"}, OrderNum: 80, IsActive: true},
	{ExternalID: "svc_travel_custom", CategoryExternalID: "svc_hotel_travel", Names: map[string]string{"zh": "旅游定制师", "en": "Travel Custom"}, OrderNum: 90, IsActive: true},
	{ExternalID: "svc_greeter", CategoryExternalID: "svc_hotel_travel", Names: map[string]string{"zh": "前厅经理", "en": "Greeter"}, OrderNum: 100, IsActive: true},
	{ExternalID: "svc_concierge", CategoryExternalID: "svc_hotel_travel", Names: map[string]string{"zh": "礼宾相关岗位", "en": "Concierge"}, OrderNum: 110, IsActive: true},
	{ExternalID: "svc_hotel_sales", CategoryExternalID: "svc_hotel_travel", Names: map[string]string{"zh": "酒店销售", "en": "Hotel Sales"}, OrderNum: 120, IsActive: true},
	{ExternalID: "svc_banquets", CategoryExternalID: "svc_hotel_travel", Names: map[string]string{"zh": "宴会服务", "en": "Banquets"}, OrderNum: 130, IsActive: true},

	{ExternalID: "svc_cleaner", CategoryExternalID: "svc_personal", Names: map[string]string{"zh": "保洁", "en": "Cleaner"}, OrderNum: 10, IsActive: true},
	{ExternalID: "svc_guard", CategoryExternalID: "svc_personal", Names: map[string]string{"zh": "保安", "en": "Guard"}, OrderNum: 20, IsActive: true},
	{ExternalID: "svc_housekeeping", CategoryExternalID: "svc_personal", Names: map[string]string{"zh": "家政", "en": "Housekeeping"}, OrderNum: 30, IsActive: true},
	{ExternalID: "svc_babysitter", CategoryExternalID: "svc_personal", Names: map[string]string{"zh": "保姆", "en": "Babysitter"}, OrderNum: 40, IsActive: true},
	{ExternalID: "svc_maternity_matron", CategoryExternalID: "svc_personal", Names: map[string]string{"zh": "月嫂", "en": "Maternity Matron"}, OrderNum: 50, IsActive: true},
	{ExternalID: "svc_housekeeper", CategoryExternalID: "svc_personal", Names: map[string]string{"zh": "管家", "en": "Housekeeper"}, OrderNum: 60, IsActive: true},
	{ExternalID: "svc_pet_beautician", CategoryExternalID: "svc_personal", Names: map[string]string{"zh": "宠物美容师", "en": "Pet Beautician"}, OrderNum: 70, IsActive: true},
	{ExternalID: "svc_pet_doctor", CategoryExternalID: "svc_personal", Names: map[string]string{"zh": "宠物医生", "en": "Pet Doctor"}, OrderNum: 80, IsActive: true},
	{ExternalID: "svc_fashion_buyer", CategoryExternalID: "svc_personal", Names: map[string]string{"zh": "服装买手", "en": "Fashion Buyer"}, OrderNum: 90, IsActive: true},
	{ExternalID: "svc_stylist", CategoryExternalID: "svc_personal", Names: map[string]string{"zh": "服装搭配师", "en": "Stylist"}, OrderNum: 100, IsActive: true},

	{ExternalID: "svc_fitness_coach", CategoryExternalID: "svc_fitness", Names: map[string]string{"zh": "健身教练", "en": "Fitness Coach"}, OrderNum: 10, IsActive: true},
	{ExternalID: "svc_swimming_coach", CategoryExternalID: "svc_fitness", Names: map[string]string{"zh": "游泳教练", "en": "Swimming Coach"}, OrderNum: 20, IsActive: true},
	{ExternalID: "svc_yoga_coach", CategoryExternalID: "svc_fitness", Names: map[string]string{"zh": "瑜伽教练", "en": "Yoga Coach"}, OrderNum: 30, IsActive: true},
	{ExternalID: "svc_dance_teacher", CategoryExternalID: "svc_fitness", Names: map[string]string{"zh": "舞蹈老师", "en": "Dance Teacher"}, OrderNum: 40, IsActive: true},
	{ExternalID: "svc_basketball_coach", CategoryExternalID: "svc_fitness", Names: map[string]string{"zh": "篮球教练", "en": "Basketball Coach"}, OrderNum: 50, IsActive: true},
	{ExternalID: "svc_badminton_coach", CategoryExternalID: "svc_fitness", Names: map[string]string{"zh": "羽毛球教练", "en": "Badminton Coach"}, OrderNum: 60, IsActive: true},
	{ExternalID: "svc_taekwondo_coach", CategoryExternalID: "svc_fitness", Names: map[string]string{"zh": "跆拳道教练", "en": "Taekwondo Coach"}, OrderNum: 70, IsActive: true},
	{ExternalID: "svc_martial_teacher", CategoryExternalID: "svc_fitness", Names: map[string]string{"zh": "武术教练", "en": "Martial Teacher"}, OrderNum: 80, IsActive: true},
	{ExternalID: "svc_street_dance_teacher", CategoryExternalID: "svc_fitness", Names: map[string]string{"zh": "街舞老师", "en": "Street Dance Teacher"}, OrderNum: 90, IsActive: true},

	{ExternalID: "svc_beauty_consultant", CategoryExternalID: "svc_beauty", Names: map[string]string{"zh": "美容顾问", "en": "Beauty Consultant"}, OrderNum: 10, IsActive: true},
	{ExternalID: "svc_store_manager", CategoryExternalID: "svc_beauty", Names: map[string]string{"zh": "美容店长", "en": "Store Manager"}, OrderNum: 20, IsActive: true},
	{ExternalID: "svc_beautician", CategoryExternalID: "svc_beauty", Names: map[string]string{"zh": "美容师", "en": "Beautician"}, OrderNum: 30, IsActive: true},
	{ExternalID: "svc_makeup_artist", CategoryExternalID: "svc_beauty", Names: map[string]string{"zh": "化妆师", "en": "Makeup Artist"}, OrderNum: 40, IsActive: true},
	{ExternalID: "svc_manicurist", CategoryExternalID: "svc_beauty", Names: map[string]string{"zh": "美甲师", "en": "Manicurist"}, OrderNum: 50, IsActive: true},
	{ExternalID: "svc_hairdresser", CategoryExternalID: "svc_beauty", Names: map[string]string{"zh": "美发师", "en": "Hairdresser"}, OrderNum: 60, IsActive: true},
	{ExternalID: "svc_masseur", CategoryExternalID: "svc_beauty", Names: map[string]string{"zh": "按摩师", "en": "Masseur"}, OrderNum: 70, IsActive: true},
	{ExternalID: "svc_physiotherapist", CategoryExternalID: "svc_beauty", Names: map[string]string{"zh": "理疗师", "en": "Physiotherapist"}, OrderNum: 80, IsActive: true},
	{ExternalID: "svc_tattoo_artist", CategoryExternalID: "svc_beauty", Names: map[string]string{"zh": "纹绣师", "en": "Tattoo Artist"}, OrderNum: 90, IsActive: true},

	{ExternalID: "svc_course_consultant", CategoryExternalID: "svc_consult_translate", Names: map[string]string{"zh": "课程顾问", "en": "Course Consultant"}, OrderNum: 10, IsActive: true},
	{ExternalID: "svc_headhunter", CategoryExternalID: "svc_consult_translate", Names: map[string]string{"zh": "猎头顾问", "en": "Headhunter Consultant"}, OrderNum: 20, IsActive: true},
	{ExternalID: "svc_consultant", CategoryExternalID: "svc_consult_translate", Names: map[string]string{"zh": "咨询顾问", "en": "Consultant"}, OrderNum: 30, IsActive: true},
	{ExternalID: "svc_legal_consultant", CategoryExternalID: "svc_consult_translate", Names: map[string]string{"zh": "法律顾问", "en": "Legal Consultant"}, OrderNum: 40, IsActive: true},
	{ExternalID: "svc_translator", CategoryExternalID: "svc_consult_translate", Names: map[string]string{"zh": "翻译", "en": "Translator"}, OrderNum: 50, IsActive: true},

	{ExternalID: "svc_auto_mechanic", CategoryExternalID: "svc_repair", Names: map[string]string{"zh": "汽车维修", "en": "Auto Mechanic"}, OrderNum: 10, IsActive: true},
	{ExternalID: "svc_mobile_repair", CategoryExternalID: "svc_repair", Names: map[string]string{"zh": "机务维修", "en": "Mobile Repair"}, OrderNum: 20, IsActive: true},
	{ExternalID: "svc_electric_repair", CategoryExternalID: "svc_repair", Names: map[string]string{"zh": "维修电工", "en": "Electric Repair"}, OrderNum: 30, IsActive: true},
	{ExternalID: "svc_mould_maintenance", CategoryExternalID: "svc_repair", Names: map[string]string{"zh": "模具维修", "en": "Mould Maintenance"}, OrderNum: 40, IsActive: true},
	{ExternalID: "svc_device_repair", CategoryExternalID: "svc_repair", Names: map[string]string{"zh": "器械维修", "en": "Device Repair"}, OrderNum: 50, IsActive: true},
	{ExternalID: "svc_other_repair", CategoryExternalID: "svc_repair", Names: map[string]string{"zh": "其他维修服务岗", "en": "Other Repair"}, OrderNum: 60, IsActive: true},

	{ExternalID: "svc_store_clerk", CategoryExternalID: "svc_other", Names: map[string]string{"zh": "店员", "en": "Store Clerk"}, OrderNum: 10, IsActive: true},
	{ExternalID: "svc_bid_staff", CategoryExternalID: "svc_other", Names: map[string]string{"zh": "项目招投标", "en": "Bid Staff"}, OrderNum: 20, IsActive: true},
	{ExternalID: "svc_tickets_staff", CategoryExternalID: "svc_other", Names: map[string]string{"zh": "票务", "en": "Tickets Staff"}, OrderNum: 30, IsActive: true},
	{ExternalID: "svc_trustee", CategoryExternalID: "svc_other", Names: map[string]string{"zh": "托管", "en": "Trustee"}, OrderNum: 40, IsActive: true},
	{ExternalID: "svc_fee_collector", CategoryExternalID: "svc_other", Names: map[string]string{"zh": "收费员", "en": "Fee Collector"}, OrderNum: 50, IsActive: true},
	{ExternalID: "svc_transport_agent", CategoryExternalID: "svc_other", Names: map[string]string{"zh": "客运员", "en": "Transport Agent"}, OrderNum: 60, IsActive: true},
	{ExternalID: "svc_uav_pilot", CategoryExternalID: "svc_other", Names: map[string]string{"zh": "无人机飞手", "en": "UAV Pilot"}, OrderNum: 70, IsActive: true},

	{ExternalID: "media_journalist", CategoryExternalID: "media_news_publishing", Names: map[string]string{"zh": "记者", "en": "Journalist"}, OrderNum: 10, IsActive: true},
	{ExternalID: "media_editor", CategoryExternalID: "media_news_publishing", Names: map[string]string{"zh": "编辑", "en": "Editor"}, OrderNum: 20, IsActive: true},
	{ExternalID: "media_reporter", CategoryExternalID: "media_news_publishing", Names: map[string]string{"zh": "采编", "en": "Reporter"}, OrderNum: 30, IsActive: true},
	{ExternalID: "media_chief_editor", CategoryExternalID: "media_news_publishing", Names: map[string]string{"zh": "主编/副主编", "en": "Chief Editor"}, OrderNum: 40, IsActive: true},
	{ExternalID: "media_proofreader", CategoryExternalID: "media_news_publishing", Names: map[string]string{"zh": "校对", "en": "Proofreader"}, OrderNum: 50, IsActive: true},
	{ExternalID: "media_writer", CategoryExternalID: "media_news_publishing", Names: map[string]string{"zh": "撰稿", "en": "Writer"}, OrderNum: 60, IsActive: true},
	{ExternalID: "media_reviewer", CategoryExternalID: "media_news_publishing", Names: map[string]string{"zh": "审核", "en": "Reviewer"}, OrderNum: 70, IsActive: true},

	{ExternalID: "media_host", CategoryExternalID: "media_broadcast_tv", Names: map[string]string{"zh": "主持人", "en": "Host"}, OrderNum: 10, IsActive: true},
	{ExternalID: "media_announcer", CategoryExternalID: "media_broadcast_tv", Names: map[string]string{"zh": "播音员", "en": "Announcer"}, OrderNum: 20, IsActive: true},
	{ExternalID: "media_director_tv", CategoryExternalID: "media_broadcast_tv", Names: map[string]string{"zh": "导演", "en": "Director"}, OrderNum: 30, IsActive: true},
	{ExternalID: "media_cameraman", CategoryExternalID: "media_broadcast_tv", Names: map[string]string{"zh": "摄像师", "en": "Cameraman"}, OrderNum: 40, IsActive: true},
	{ExternalID: "media_camera_assistant", CategoryExternalID: "media_broadcast_tv", Names: map[string]string{"zh": "摄影助理", "en": "Camera Assistant"}, OrderNum: 50, IsActive: true},
	{ExternalID: "media_postprod", CategoryExternalID: "media_broadcast_tv", Names: map[string]string{"zh": "后期制作", "en": "Post Production"}, OrderNum: 60, IsActive: true},
	{ExternalID: "media_vfx", CategoryExternalID: "media_broadcast_tv", Names: map[string]string{"zh": "特效师", "en": "VFX"}, OrderNum: 70, IsActive: true},
	{ExternalID: "media_sound_engineer", CategoryExternalID: "media_broadcast_tv", Names: map[string]string{"zh": "音效师", "en": "Sound Engineer"}, OrderNum: 80, IsActive: true},
	{ExternalID: "media_switcher", CategoryExternalID: "media_broadcast_tv", Names: map[string]string{"zh": "导播", "en": "Switcher"}, OrderNum: 90, IsActive: true},
	{ExternalID: "media_program_planner", CategoryExternalID: "media_broadcast_tv", Names: map[string]string{"zh": "节目策划", "en": "Program Planner"}, OrderNum: 100, IsActive: true},
	{ExternalID: "media_producer", CategoryExternalID: "media_broadcast_tv", Names: map[string]string{"zh": "制片人", "en": "Producer"}, OrderNum: 110, IsActive: true},
	{ExternalID: "media_film_making", CategoryExternalID: "media_broadcast_tv", Names: map[string]string{"zh": "影视制作", "en": "Film Making"}, OrderNum: 120, IsActive: true},
	{ExternalID: "media_channel_specialist", CategoryExternalID: "media_broadcast_tv", Names: map[string]string{"zh": "渠道专员", "en": "Channel Specialist"}, OrderNum: 130, IsActive: true},

	{ExternalID: "media_director", CategoryExternalID: "media_film_performance", Names: map[string]string{"zh": "导演", "en": "Director"}, OrderNum: 10, IsActive: true},
	{ExternalID: "media_screenwriter", CategoryExternalID: "media_film_performance", Names: map[string]string{"zh": "编剧", "en": "Screenwriter"}, OrderNum: 20, IsActive: true},
	{ExternalID: "media_actor", CategoryExternalID: "media_film_performance", Names: map[string]string{"zh": "演员", "en": "Actor"}, OrderNum: 30, IsActive: true},
	{ExternalID: "media_script_supervisor", CategoryExternalID: "media_film_performance", Names: map[string]string{"zh": "场记", "en": "Script Supervisor"}, OrderNum: 40, IsActive: true},
	{ExternalID: "media_artist_assistant", CategoryExternalID: "media_film_performance", Names: map[string]string{"zh": "艺人助理", "en": "Artist Assistant"}, OrderNum: 50, IsActive: true},
	{ExternalID: "media_agent", CategoryExternalID: "media_film_performance", Names: map[string]string{"zh": "经纪人", "en": "Agent"}, OrderNum: 60, IsActive: true},
	{ExternalID: "media_model", CategoryExternalID: "media_film_performance", Names: map[string]string{"zh": "模特", "en": "Model"}, OrderNum: 70, IsActive: true},
	{ExternalID: "media_stage_designer", CategoryExternalID: "media_film_performance", Names: map[string]string{"zh": "舞美设计", "en": "Stage Designer"}, OrderNum: 80, IsActive: true},
	{ExternalID: "media_star_mapper", CategoryExternalID: "media_film_performance", Names: map[string]string{"zh": "星探", "en": "Star Mapper"}, OrderNum: 90, IsActive: true},
	{ExternalID: "media_intern", CategoryExternalID: "media_film_performance", Names: map[string]string{"zh": "练习生", "en": "Intern"}, OrderNum: 100, IsActive: true},

	{ExternalID: "trade_business_staff", CategoryExternalID: "trade_foreign_trade", Names: map[string]string{"zh": "外贸业务员", "en": "Business Staff"}, OrderNum: 10, IsActive: true},
	{ExternalID: "trade_specialist", CategoryExternalID: "trade_foreign_trade", Names: map[string]string{"zh": "外贸专员", "en": "Specialist"}, OrderNum: 20, IsActive: true},
	{ExternalID: "trade_doc_tracker", CategoryExternalID: "trade_foreign_trade", Names: map[string]string{"zh": "外贸跟单员", "en": "Doc Tracker"}, OrderNum: 30, IsActive: true},
	{ExternalID: "trade_assistant", CategoryExternalID: "trade_foreign_trade", Names: map[string]string{"zh": "外贸助理", "en": "Assistant"}, OrderNum: 40, IsActive: true},
	{ExternalID: "trade_manager", CategoryExternalID: "trade_foreign_trade", Names: map[string]string{"zh": "外贸经理", "en": "Manager"}, OrderNum: 50, IsActive: true},
	{ExternalID: "trade_doc_staff", CategoryExternalID: "trade_foreign_trade", Names: map[string]string{"zh": "外贸单证员", "en": "Doc Staff"}, OrderNum: 60, IsActive: true},
	{ExternalID: "trade_cs", CategoryExternalID: "trade_foreign_trade", Names: map[string]string{"zh": "外贸客服", "en": "Customer Service"}, OrderNum: 70, IsActive: true},
	{ExternalID: "trade_director", CategoryExternalID: "trade_foreign_trade", Names: map[string]string{"zh": "外贸总监", "en": "Director"}, OrderNum: 80, IsActive: true},

	{ExternalID: "trade_customs_broker", CategoryExternalID: "trade_trade_support", Names: map[string]string{"zh": "报关员", "en": "Customs Broker"}, OrderNum: 10, IsActive: true},
	{ExternalID: "trade_customs_supervisor", CategoryExternalID: "trade_trade_support", Names: map[string]string{"zh": "报关主管", "en": "Customs Supervisor"}, OrderNum: 20, IsActive: true},
	{ExternalID: "trade_docs", CategoryExternalID: "trade_trade_support", Names: map[string]string{"zh": "单证员", "en": "Docs"}, OrderNum: 30, IsActive: true},

	{ExternalID: "trade_cbe", CategoryExternalID: "trade_crossborder_ecom", Names: map[string]string{"zh": "跨境电商", "en": "Cross-border E-commerce"}, OrderNum: 10, IsActive: true},
	{ExternalID: "trade_cbe_specialist", CategoryExternalID: "trade_crossborder_ecom", Names: map[string]string{"zh": "跨境电商专员", "en": "Cross-border E-commerce Specialist"}, OrderNum: 20, IsActive: true},
	{ExternalID: "trade_cbe_ops", CategoryExternalID: "trade_crossborder_ecom", Names: map[string]string{"zh": "跨境电商业务员", "en": "Cross-border E-commerce Ops"}, OrderNum: 30, IsActive: true},
	{ExternalID: "trade_cbe_assistant", CategoryExternalID: "trade_crossborder_ecom", Names: map[string]string{"zh": "跨境电商运营助理", "en": "Cross-border E-commerce Assistant"}, OrderNum: 40, IsActive: true},
	{ExternalID: "trade_amazon_ops", CategoryExternalID: "trade_crossborder_ecom", Names: map[string]string{"zh": "亚马逊运营", "en": "Amazon Ops"}, OrderNum: 50, IsActive: true},
	{ExternalID: "trade_amazon_sales", CategoryExternalID: "trade_crossborder_ecom", Names: map[string]string{"zh": "亚马逊销售", "en": "Amazon Sales"}, OrderNum: 60, IsActive: true},

	{ExternalID: "trade_translator", CategoryExternalID: "trade_translation_support", Names: map[string]string{"zh": "外贸翻译", "en": "Translator"}, OrderNum: 10, IsActive: true},
	{ExternalID: "trade_en_translator", CategoryExternalID: "trade_translation_support", Names: map[string]string{"zh": "英语翻译", "en": "English Translator"}, OrderNum: 20, IsActive: true},
	{ExternalID: "trade_jp_translator", CategoryExternalID: "trade_translation_support", Names: map[string]string{"zh": "日语翻译", "en": "Japanese Translator"}, OrderNum: 30, IsActive: true},
	{ExternalID: "trade_kr_translator", CategoryExternalID: "trade_translation_support", Names: map[string]string{"zh": "韩语翻译", "en": "Korean Translator"}, OrderNum: 40, IsActive: true},

	{ExternalID: "energy_power_engineer", CategoryExternalID: "energy_traditional", Names: map[string]string{"zh": "电力工程师", "en": "Power Engineer"}, OrderNum: 10, IsActive: true},
	{ExternalID: "energy_new_energy_engineer", CategoryExternalID: "energy_traditional", Names: map[string]string{"zh": "新能源工程师", "en": "New Energy Engineer"}, OrderNum: 20, IsActive: true},
	{ExternalID: "energy_thermal_engineer", CategoryExternalID: "energy_traditional", Names: map[string]string{"zh": "热能工程师", "en": "Thermal Engineer"}, OrderNum: 30, IsActive: true},
	{ExternalID: "energy_oil_engineer", CategoryExternalID: "energy_traditional", Names: map[string]string{"zh": "石油工程师", "en": "Oil Engineer"}, OrderNum: 40, IsActive: true},
	{ExternalID: "energy_gas_engineer", CategoryExternalID: "energy_traditional", Names: map[string]string{"zh": "燃气工程师", "en": "Gas Engineer"}, OrderNum: 50, IsActive: true},
	{ExternalID: "energy_hvac_engineer", CategoryExternalID: "energy_traditional", Names: map[string]string{"zh": "暖通工程师", "en": "HVAC Engineer"}, OrderNum: 60, IsActive: true},
	{ExternalID: "energy_engineer", CategoryExternalID: "energy_traditional", Names: map[string]string{"zh": "能源工程师", "en": "Energy Engineer"}, OrderNum: 70, IsActive: true},
	{ExternalID: "energy_pv_engineer", CategoryExternalID: "energy_traditional", Names: map[string]string{"zh": "光伏系统工程师", "en": "PV System Engineer"}, OrderNum: 80, IsActive: true},
	{ExternalID: "energy_wind_engineer", CategoryExternalID: "energy_traditional", Names: map[string]string{"zh": "风电工程师", "en": "Wind Engineer"}, OrderNum: 90, IsActive: true},

	{ExternalID: "env_engineer", CategoryExternalID: "energy_environment", Names: map[string]string{"zh": "环保工程师", "en": "Environmental Engineer"}, OrderNum: 10, IsActive: true},
	{ExternalID: "env_environment_engineer", CategoryExternalID: "energy_environment", Names: map[string]string{"zh": "环境工程师", "en": "Environment Engineer"}, OrderNum: 20, IsActive: true},
	{ExternalID: "env_ehs_engineer", CategoryExternalID: "energy_environment", Names: map[string]string{"zh": "EHS工程师", "en": "EHS Engineer"}, OrderNum: 30, IsActive: true},
	{ExternalID: "env_water_treatment", CategoryExternalID: "energy_environment", Names: map[string]string{"zh": "水处理工程师", "en": "Water Treatment Engineer"}, OrderNum: 40, IsActive: true},
	{ExternalID: "env_mep_water", CategoryExternalID: "energy_environment", Names: map[string]string{"zh": "给排水工程师", "en": "Water Engineering Engineer"}, OrderNum: 50, IsActive: true},
	{ExternalID: "env_eia_engineer", CategoryExternalID: "energy_environment", Names: map[string]string{"zh": "环评工程师", "en": "Environmental Impact Assessment Engineer"}, OrderNum: 60, IsActive: true},
	{ExternalID: "env_tech", CategoryExternalID: "energy_environment", Names: map[string]string{"zh": "环保技术员", "en": "Environmental Technician"}, OrderNum: 70, IsActive: true},
	{ExternalID: "env_inspection", CategoryExternalID: "energy_environment", Names: map[string]string{"zh": "环保检测", "en": "Environmental Inspection"}, OrderNum: 80, IsActive: true},
	{ExternalID: "env_specialist", CategoryExternalID: "energy_environment", Names: map[string]string{"zh": "环保专员", "en": "Environmental Specialist"}, OrderNum: 90, IsActive: true},
	{ExternalID: "env_supervisor", CategoryExternalID: "energy_environment", Names: map[string]string{"zh": "环保主管", "en": "Environmental Supervisor"}, OrderNum: 100, IsActive: true},

	{ExternalID: "agri_tech", CategoryExternalID: "agri_planting", Names: map[string]string{"zh": "农业技术员", "en": "Agricultural Technician"}, OrderNum: 10, IsActive: true},
	{ExternalID: "agri_agronomist", CategoryExternalID: "agri_planting", Names: map[string]string{"zh": "农艺师", "en": "Agronomist"}, OrderNum: 20, IsActive: true},
	{ExternalID: "agri_horticulturist", CategoryExternalID: "agri_planting", Names: map[string]string{"zh": "园艺师", "en": "Horticulturist"}, OrderNum: 30, IsActive: true},
	{ExternalID: "agri_florist", CategoryExternalID: "agri_planting", Names: map[string]string{"zh": "花艺师", "en": "Florist"}, OrderNum: 40, IsActive: true},
	{ExternalID: "agri_farm_machine", CategoryExternalID: "agri_planting", Names: map[string]string{"zh": "农机操作修理", "en": "Farm Machine Operator"}, OrderNum: 50, IsActive: true},

	{ExternalID: "forestry_engineer", CategoryExternalID: "agri_forestry", Names: map[string]string{"zh": "林业工程师", "en": "Forestry Engineer"}, OrderNum: 10, IsActive: true},
	{ExternalID: "forestry_tech", CategoryExternalID: "agri_forestry", Names: map[string]string{"zh": "林业技术员", "en": "Forestry Technician"}, OrderNum: 20, IsActive: true},
	{ExternalID: "forestry_garden_engineer", CategoryExternalID: "agri_forestry", Names: map[string]string{"zh": "园林工程师", "en": "Garden Engineer"}, OrderNum: 30, IsActive: true},
	{ExternalID: "forestry_ranger", CategoryExternalID: "agri_forestry", Names: map[string]string{"zh": "护林员", "en": "Ranger"}, OrderNum: 40, IsActive: true},

	{ExternalID: "livestock_specialist", CategoryExternalID: "agri_livestock", Names: map[string]string{"zh": "畜牧师", "en": "Livestock Specialist"}, OrderNum: 10, IsActive: true},
	{ExternalID: "livestock_vet", CategoryExternalID: "agri_livestock", Names: map[string]string{"zh": "兽医", "en": "Vet"}, OrderNum: 20, IsActive: true},
	{ExternalID: "livestock_breeding_tech", CategoryExternalID: "agri_livestock", Names: map[string]string{"zh": "养殖技术员", "en": "Breeding Technician"}, OrderNum: 30, IsActive: true},
	{ExternalID: "livestock_farm_mgmt", CategoryExternalID: "agri_livestock", Names: map[string]string{"zh": "牧场管理", "en": "Farm Management"}, OrderNum: 40, IsActive: true},
	{ExternalID: "livestock_feeder", CategoryExternalID: "agri_livestock", Names: map[string]string{"zh": "饲养员", "en": "Feeder"}, OrderNum: 50, IsActive: true},
	{ExternalID: "livestock_quarantine", CategoryExternalID: "agri_livestock", Names: map[string]string{"zh": "动物检疫员", "en": "Quarantine"}, OrderNum: 60, IsActive: true},
	{ExternalID: "livestock_trainer", CategoryExternalID: "agri_livestock", Names: map[string]string{"zh": "动物驯养师", "en": "Trainer"}, OrderNum: 70, IsActive: true},
	{ExternalID: "livestock_guide", CategoryExternalID: "agri_livestock", Names: map[string]string{"zh": "动物讲解员", "en": "Guide"}, OrderNum: 80, IsActive: true},

	{ExternalID: "fishery_farmer", CategoryExternalID: "agri_fishery", Names: map[string]string{"zh": "水产养殖员", "en": "Fishery Farmer"}, OrderNum: 10, IsActive: true},
	{ExternalID: "fishery_tech", CategoryExternalID: "agri_fishery", Names: map[string]string{"zh": "水产技术员", "en": "Fishery Technician"}, OrderNum: 20, IsActive: true},
	{ExternalID: "fishery_aquaculture", CategoryExternalID: "agri_fishery", Names: map[string]string{"zh": "渔业养殖员", "en": "Aquaculture"}, OrderNum: 30, IsActive: true},
	{ExternalID: "fishery_fisher", CategoryExternalID: "agri_fishery", Names: map[string]string{"zh": "捕捞员", "en": "Fisher"}, OrderNum: 40, IsActive: true},

	{ExternalID: "public_civil_servant", CategoryExternalID: "public_services", Names: map[string]string{"zh": "公务员", "en": "Civil Servant"}, OrderNum: 10, IsActive: true},
	{ExternalID: "public_ngo_specialist", CategoryExternalID: "public_services", Names: map[string]string{"zh": "非营利组织专员", "en": "NGO Specialist"}, OrderNum: 20, IsActive: true},
	{ExternalID: "public_urban_mgmt", CategoryExternalID: "public_services", Names: map[string]string{"zh": "城管", "en": "Urban Management"}, OrderNum: 30, IsActive: true},
	{ExternalID: "public_military", CategoryExternalID: "public_services", Names: map[string]string{"zh": "军人", "en": "Military"}, OrderNum: 40, IsActive: true},
	{ExternalID: "public_firefighter", CategoryExternalID: "public_services", Names: map[string]string{"zh": "消防员", "en": "Firefighter"}, OrderNum: 50, IsActive: true},
	{ExternalID: "public_community_worker", CategoryExternalID: "public_services", Names: map[string]string{"zh": "社区工作者", "en": "Community Worker"}, OrderNum: 60, IsActive: true},
	{ExternalID: "public_police", CategoryExternalID: "public_services", Names: map[string]string{"zh": "警察/辅警", "en": "Police/Paramilitary"}, OrderNum: 70, IsActive: true},
	{ExternalID: "public_party_affairs", CategoryExternalID: "public_services", Names: map[string]string{"zh": "党建管理岗", "en": "Party Affairs"}, OrderNum: 80, IsActive: true},

	{ExternalID: "research_assistant", CategoryExternalID: "public_research", Names: map[string]string{"zh": "科研助理", "en": "Research Assistant"}, OrderNum: 10, IsActive: true},
	{ExternalID: "research_staff", CategoryExternalID: "public_research", Names: map[string]string{"zh": "科研人员", "en": "Research Staff"}, OrderNum: 20, IsActive: true},
	{ExternalID: "research_management", CategoryExternalID: "public_research", Names: map[string]string{"zh": "科研管理", "en": "Research Management"}, OrderNum: 30, IsActive: true},
	{ExternalID: "research_academic_promo", CategoryExternalID: "public_research", Names: map[string]string{"zh": "学术推广", "en": "Academic Promotion"}, OrderNum: 40, IsActive: true},
	{ExternalID: "research_chem_analysis", CategoryExternalID: "public_research", Names: map[string]string{"zh": "化学分析", "en": "Chemical Analysis"}, OrderNum: 50, IsActive: true},

	{ExternalID: "social_worker", CategoryExternalID: "public_social", Names: map[string]string{"zh": "社工", "en": "Social Worker"}, OrderNum: 10, IsActive: true},
	{ExternalID: "social_worker_assistant", CategoryExternalID: "public_social", Names: map[string]string{"zh": "社工助理", "en": "Social Worker Assistant"}, OrderNum: 20, IsActive: true},
	{ExternalID: "social_worker_intern", CategoryExternalID: "public_social", Names: map[string]string{"zh": "社工实习生", "en": "Social Worker Intern"}, OrderNum: 30, IsActive: true},
	{ExternalID: "volunteer", CategoryExternalID: "public_social", Names: map[string]string{"zh": "志愿者", "en": "Volunteer"}, OrderNum: 40, IsActive: true},
	{ExternalID: "unpaid_volunteer", CategoryExternalID: "public_social", Names: map[string]string{"zh": "义工", "en": "Unpaid Volunteer"}, OrderNum: 50, IsActive: true},
	{ExternalID: "caregiver", CategoryExternalID: "public_social", Names: map[string]string{"zh": "护理员", "en": "Caregiver"}, OrderNum: 60, IsActive: true},
	{ExternalID: "rehab_therapist", CategoryExternalID: "public_social", Names: map[string]string{"zh": "康复师", "en": "Rehabilitation Therapist"}, OrderNum: 70, IsActive: true},
}

var seedCategories = []SeedJobCategory{
	{ExternalID: "it", Names: map[string]string{"zh": "IT | 互联网", "en": "Internet Industry"}, ParentExternalID: "", OrderNum: 10, IsActive: true},
	{ExternalID: "finance", Names: map[string]string{"zh": "金融 | 银行", "en": "Banking Finance"}, ParentExternalID: "", OrderNum: 20, IsActive: true},
	{ExternalID: "education", Names: map[string]string{"zh": "教育 | 培训", "en": "Education"}, ParentExternalID: "", OrderNum: 30, IsActive: true},
	{ExternalID: "healthcare", Names: map[string]string{"zh": "医疗 | 健康", "en": "Healthcare"}, ParentExternalID: "", OrderNum: 40, IsActive: true},
	{ExternalID: "realestate", Names: map[string]string{"zh": "建筑 | 房地产", "en": "Real Estate"}, ParentExternalID: "", OrderNum: 50, IsActive: true},
	{ExternalID: "manufacturing", Names: map[string]string{"zh": "制造 | 生产", "en": "Manufacturing"}, ParentExternalID: "", OrderNum: 60, IsActive: true},
	{ExternalID: "logistics", Names: map[string]string{"zh": "交通 | 物流", "en": "Logistics"}, ParentExternalID: "", OrderNum: 70, IsActive: true},
	{ExternalID: "services", Names: map[string]string{"zh": "服务 | 消费", "en": "Service Industry"}, ParentExternalID: "", OrderNum: 80, IsActive: true},
	{ExternalID: "media", Names: map[string]string{"zh": "文化 | 传媒", "en": "Media"}, ParentExternalID: "", OrderNum: 90, IsActive: true},
	{ExternalID: "trade", Names: map[string]string{"zh": "贸易 | 进出口", "en": "Trade"}, ParentExternalID: "", OrderNum: 100, IsActive: true},
	{ExternalID: "energy", Names: map[string]string{"zh": "能源 | 环保", "en": "Eco-energy"}, ParentExternalID: "", OrderNum: 110, IsActive: true},
	{ExternalID: "agriculture", Names: map[string]string{"zh": "农林 | 牧渔", "en": "Agriculture"}, ParentExternalID: "", OrderNum: 120, IsActive: true},
	{ExternalID: "public", Names: map[string]string{"zh": "公共 | 非盈利", "en": "Public Service"}, ParentExternalID: "", OrderNum: 130, IsActive: true},
	{ExternalID: "others", Names: map[string]string{"zh": "其他 | 岗位", "en": "Others"}, ParentExternalID: "", OrderNum: 140, IsActive: true},

	{ExternalID: "it_backend", Names: map[string]string{"zh": "后端开发/程序员", "en": "Backend Developer/Programmer"}, ParentExternalID: "it", OrderNum: 10, IsActive: true},
	{ExternalID: "it_frontend", Names: map[string]string{"zh": "前端开发", "en": "Frontend Developer"}, ParentExternalID: "it", OrderNum: 20, IsActive: true},
	{ExternalID: "it_mobile", Names: map[string]string{"zh": "移动开发", "en": "Mobile Developer"}, ParentExternalID: "it", OrderNum: 30, IsActive: true},
	{ExternalID: "it_testing", Names: map[string]string{"zh": "软件测试", "en": "Software Tester"}, ParentExternalID: "it", OrderNum: 40, IsActive: true},
	{ExternalID: "it_ops_sec_dba", Names: map[string]string{"zh": "运维 / 安全 / DBA", "en": "Operations / Security / DBA"}, ParentExternalID: "it", OrderNum: 50, IsActive: true},
	{ExternalID: "it_ai_bigdata", Names: map[string]string{"zh": "新兴技术", "en": "Emerging Technologies"}, ParentExternalID: "it", OrderNum: 60, IsActive: true},
	{ExternalID: "it_other_tech", Names: map[string]string{"zh": "其他技术岗", "en": "Other Technical Positions"}, ParentExternalID: "it", OrderNum: 70, IsActive: true},
	{ExternalID: "it_senior", Names: map[string]string{"zh": "高端技术职位", "en": "Senior Technical Positions"}, ParentExternalID: "it", OrderNum: 80, IsActive: true},

	{ExternalID: "finance_counter_service", Names: map[string]string{"zh": "银行柜台/服务", "en": "Bank Counter/Service"}, ParentExternalID: "finance", OrderNum: 10, IsActive: true},
	{ExternalID: "finance_personal_wealth", Names: map[string]string{"zh": "个人金融与理财", "en": "Personal Finance and Wealth Management"}, ParentExternalID: "finance", OrderNum: 20, IsActive: true},
	{ExternalID: "finance_credit_approval", Names: map[string]string{"zh": "信贷/审批", "en": "Credit Approval"}, ParentExternalID: "finance", OrderNum: 30, IsActive: true},
	{ExternalID: "finance_risk_compliance", Names: map[string]string{"zh": "风险管理/合规", "en": "Risk Management/Compliance"}, ParentExternalID: "finance", OrderNum: 40, IsActive: true},
	{ExternalID: "finance_securities_invest", Names: map[string]string{"zh": "证券与投资", "en": "Securities and Investment"}, ParentExternalID: "finance", OrderNum: 50, IsActive: true},
	{ExternalID: "finance_insurance_actuary", Names: map[string]string{"zh": "保险/精算", "en": "Insurance/Actuary"}, ParentExternalID: "finance", OrderNum: 60, IsActive: true},
	{ExternalID: "finance_banking_support", Names: map[string]string{"zh": "银行业务支持", "en": "Banking Support"}, ParentExternalID: "finance", OrderNum: 70, IsActive: true},
	{ExternalID: "finance_trust_futures", Names: map[string]string{"zh": "信托/期货类", "en": "Trust/Futures"}, ParentExternalID: "finance", OrderNum: 80, IsActive: true},
	{ExternalID: "finance_bank_management", Names: map[string]string{"zh": "银行管理类岗位", "en": "Bank Management"}, ParentExternalID: "finance", OrderNum: 90, IsActive: true},
	{ExternalID: "finance_intern", Names: map[string]string{"zh": "银行新人/实习生", "en": "Bank Intern/Trainee"}, ParentExternalID: "finance", OrderNum: 100, IsActive: true},

	{ExternalID: "education_teacher", Names: map[string]string{"zh": "教师", "en": "Teacher"}, ParentExternalID: "education", OrderNum: 10, IsActive: true},
	{ExternalID: "education_teaching_admin", Names: map[string]string{"zh": "教学管理", "en": "Teaching Administration"}, ParentExternalID: "education", OrderNum: 20, IsActive: true},
	{ExternalID: "education_student_services", Names: map[string]string{"zh": "学生服务", "en": "Student Services"}, ParentExternalID: "education", OrderNum: 30, IsActive: true},
	{ExternalID: "education_training_lecturer", Names: map[string]string{"zh": "培训/讲师", "en": "Training/Lecturer"}, ParentExternalID: "education", OrderNum: 40, IsActive: true},
	{ExternalID: "education_training_management", Names: map[string]string{"zh": "培训管理", "en": "Training Management"}, ParentExternalID: "education", OrderNum: 50, IsActive: true},

	{ExternalID: "healthcare_doctor", Names: map[string]string{"zh": "医生", "en": "Doctor"}, ParentExternalID: "healthcare", OrderNum: 10, IsActive: true},
	{ExternalID: "healthcare_nurse", Names: map[string]string{"zh": "护士", "en": "Nurse"}, ParentExternalID: "healthcare", OrderNum: 20, IsActive: true},
	{ExternalID: "healthcare_medtech", Names: map[string]string{"zh": "医学技术岗", "en": "Medical Technology"}, ParentExternalID: "healthcare", OrderNum: 30, IsActive: true},
	{ExternalID: "healthcare_pharma", Names: map[string]string{"zh": "药学岗", "en": "Pharmacy"}, ParentExternalID: "healthcare", OrderNum: 40, IsActive: true},
	{ExternalID: "healthcare_devices", Names: map[string]string{"zh": "医疗器械岗", "en": "Medical Devices"}, ParentExternalID: "healthcare", OrderNum: 50, IsActive: true},
	{ExternalID: "healthcare_other", Names: map[string]string{"zh": "其他医疗岗", "en": "Other Healthcare Positions"}, ParentExternalID: "healthcare", OrderNum: 60, IsActive: true},
	{ExternalID: "healthcare_intern", Names: map[string]string{"zh": "医药实习", "en": "Healthcare Intern"}, ParentExternalID: "healthcare", OrderNum: 70, IsActive: true},
	{ExternalID: "healthcare_management", Names: map[string]string{"zh": "医疗管理", "en": "Healthcare Management"}, ParentExternalID: "healthcare", OrderNum: 80, IsActive: true},

	{ExternalID: "realestate_design_planning", Names: map[string]string{"zh": "建筑设计/规划", "en": "Architecture Design/Planning"}, ParentExternalID: "realestate", OrderNum: 10, IsActive: true},
	{ExternalID: "realestate_interior_landscape", Names: map[string]string{"zh": "室内/景观设计", "en": "Interior/Landscape Design"}, ParentExternalID: "realestate", OrderNum: 20, IsActive: true},
	{ExternalID: "realestate_cost_budget", Names: map[string]string{"zh": "工程造价/预算", "en": "Cost/Budget"}, ParentExternalID: "realestate", OrderNum: 30, IsActive: true},
	{ExternalID: "realestate_construction_mgmt", Names: map[string]string{"zh": "工程施工/管理", "en": "Construction Management"}, ParentExternalID: "realestate", OrderNum: 40, IsActive: true},
	{ExternalID: "realestate_project_mgmt", Names: map[string]string{"zh": "项目管理", "en": "Project Management"}, ParentExternalID: "realestate", OrderNum: 50, IsActive: true},
	{ExternalID: "realestate_sales_planning", Names: map[string]string{"zh": "房地产销售/策划", "en": "Real Estate Sales/Planning"}, ParentExternalID: "realestate", OrderNum: 60, IsActive: true},
	{ExternalID: "realestate_property_mgmt", Names: map[string]string{"zh": "物业管理", "en": "Property Management"}, ParentExternalID: "realestate", OrderNum: 70, IsActive: true},

	{ExternalID: "mfg_mechanical", Names: map[string]string{"zh": "机械制造", "en": "Mechanical Engineering"}, ParentExternalID: "manufacturing", OrderNum: 10, IsActive: true},
	{ExternalID: "mfg_electrical", Names: map[string]string{"zh": "电子/电气制造", "en": "Electrical/Electronic Engineering"}, ParentExternalID: "manufacturing", OrderNum: 20, IsActive: true},
	{ExternalID: "mfg_auto_transport", Names: map[string]string{"zh": "汽车/交通制造", "en": "Automotive/Transportation Engineering"}, ParentExternalID: "manufacturing", OrderNum: 30, IsActive: true},
	{ExternalID: "mfg_process_mold", Names: map[string]string{"zh": "工艺/模具工程", "en": "Process/Mold Engineering"}, ParentExternalID: "manufacturing", OrderNum: 40, IsActive: true},
	{ExternalID: "mfg_prod_equip", Names: map[string]string{"zh": "生产/设备管理", "en": "Production/Equipment Management"}, ParentExternalID: "manufacturing", OrderNum: 50, IsActive: true},
	{ExternalID: "mfg_quality", Names: map[string]string{"zh": "质量管理", "en": "Quality Management"}, ParentExternalID: "manufacturing", OrderNum: 60, IsActive: true},
	{ExternalID: "mfg_rd_design", Names: map[string]string{"zh": "研发设计", "en": "Research and Development Design"}, ParentExternalID: "manufacturing", OrderNum: 70, IsActive: true},
	{ExternalID: "mfg_worker", Names: map[string]string{"zh": "生产技工", "en": "Production Worker"}, ParentExternalID: "manufacturing", OrderNum: 80, IsActive: true},

	{ExternalID: "logi_transport_service", Names: map[string]string{"zh": "运输服务", "en": "Transportation Service"}, ParentExternalID: "logistics", OrderNum: 10, IsActive: true},
	{ExternalID: "logi_delivery", Names: map[string]string{"zh": "运输配送", "en": "Delivery"}, ParentExternalID: "logistics", OrderNum: 20, IsActive: true},
	{ExternalID: "logi_warehouse", Names: map[string]string{"zh": "仓储管理", "en": "Warehouse Management"}, ParentExternalID: "logistics", OrderNum: 30, IsActive: true},
	{ExternalID: "logi_ops", Names: map[string]string{"zh": "物流运营", "en": "Logistics Operations"}, ParentExternalID: "logistics", OrderNum: 40, IsActive: true},
	{ExternalID: "logi_supply", Names: map[string]string{"zh": "供应链", "en": "Supply Chain"}, ParentExternalID: "logistics", OrderNum: 50, IsActive: true},

	{ExternalID: "svc_food", Names: map[string]string{"zh": "餐饮服务", "en": "Food Services"}, ParentExternalID: "services", OrderNum: 10, IsActive: true},
	{ExternalID: "svc_hotel_travel", Names: map[string]string{"zh": "酒店旅游", "en": "Hotel Travel"}, ParentExternalID: "services", OrderNum: 20, IsActive: true},
	{ExternalID: "svc_personal", Names: map[string]string{"zh": "个人服务", "en": "Personal Services"}, ParentExternalID: "services", OrderNum: 30, IsActive: true},
	{ExternalID: "svc_fitness", Names: map[string]string{"zh": "运动健身", "en": "Fitness"}, ParentExternalID: "services", OrderNum: 40, IsActive: true},
	{ExternalID: "svc_beauty", Names: map[string]string{"zh": "美容保健", "en": "Beauty"}, ParentExternalID: "services", OrderNum: 50, IsActive: true},
	{ExternalID: "svc_consult_translate", Names: map[string]string{"zh": "咨询翻译", "en": "Consultation/Translation"}, ParentExternalID: "services", OrderNum: 60, IsActive: true},
	{ExternalID: "svc_repair", Names: map[string]string{"zh": "维修岗位", "en": "Repair"}, ParentExternalID: "services", OrderNum: 70, IsActive: true},
	{ExternalID: "svc_other", Names: map[string]string{"zh": "其他服务岗", "en": "Other Services"}, ParentExternalID: "services", OrderNum: 80, IsActive: true},

	{ExternalID: "media_news_publishing", Names: map[string]string{"zh": "新闻/出版", "en": "News/Publishing"}, ParentExternalID: "media", OrderNum: 10, IsActive: true},
	{ExternalID: "media_broadcast_tv", Names: map[string]string{"zh": "广播电视", "en": "Broadcast/TV"}, ParentExternalID: "media", OrderNum: 20, IsActive: true},
	{ExternalID: "media_film_performance", Names: map[string]string{"zh": "影视演艺", "en": "Film/Performance"}, ParentExternalID: "media", OrderNum: 30, IsActive: true},

	{ExternalID: "trade_foreign_trade", Names: map[string]string{"zh": "外贸业务", "en": "Foreign Trade"}, ParentExternalID: "trade", OrderNum: 10, IsActive: true},
	{ExternalID: "trade_trade_support", Names: map[string]string{"zh": "国际贸易支持", "en": "Trade Support"}, ParentExternalID: "trade", OrderNum: 20, IsActive: true},
	{ExternalID: "trade_crossborder_ecom", Names: map[string]string{"zh": "跨境电商", "en": "Cross-Border E-Commerce"}, ParentExternalID: "trade", OrderNum: 30, IsActive: true},
	{ExternalID: "trade_translation_support", Names: map[string]string{"zh": "翻译支持", "en": "Translation Support"}, ParentExternalID: "trade", OrderNum: 40, IsActive: true},

	{ExternalID: "energy_traditional", Names: map[string]string{"zh": "传统能源", "en": "Traditional Energy"}, ParentExternalID: "energy", OrderNum: 10, IsActive: true},
	{ExternalID: "energy_environment", Names: map[string]string{"zh": "环境保护", "en": "Environmental Protection"}, ParentExternalID: "energy", OrderNum: 20, IsActive: true},

	{ExternalID: "agri_planting", Names: map[string]string{"zh": "农业种植", "en": "Agricultural Planting"}, ParentExternalID: "agriculture", OrderNum: 10, IsActive: true},
	{ExternalID: "agri_forestry", Names: map[string]string{"zh": "林业", "en": "Forestry"}, ParentExternalID: "agriculture", OrderNum: 20, IsActive: true},
	{ExternalID: "agri_livestock", Names: map[string]string{"zh": "畜牧业", "en": "Livestock"}, ParentExternalID: "agriculture", OrderNum: 30, IsActive: true},
	{ExternalID: "agri_fishery", Names: map[string]string{"zh": "渔业水产", "en": "Fishery"}, ParentExternalID: "agriculture", OrderNum: 40, IsActive: true},

	{ExternalID: "public_services", Names: map[string]string{"zh": "公共服务", "en": "Public Services"}, ParentExternalID: "public", OrderNum: 10, IsActive: true},
	{ExternalID: "public_research", Names: map[string]string{"zh": "科研相关", "en": "Research"}, ParentExternalID: "public", OrderNum: 20, IsActive: true},
	{ExternalID: "public_social", Names: map[string]string{"zh": "社会服务", "en": "Social Services"}, ParentExternalID: "public", OrderNum: 30, IsActive: true},
}
