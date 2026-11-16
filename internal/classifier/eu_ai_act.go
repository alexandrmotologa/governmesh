package classifier

import (
	"strings"

	"github.com/alexandrmotologa/governmesh/internal/config"
	"github.com/alexandrmotologa/governmesh/internal/scanner"
)

// IntentIndicator maps domain keywords to risk classifications.
type IntentIndicator struct {
	Keywords          []string
	Tier              config.RiskLevel
	Article           string
	Justification     string
	ProhibitedReason  string
	RequiresHumanGate bool
	NISTCategories    []string
	ISO42001Controls  []string
}

var (
	// EU AI Act rule indicators covering Article 5 (Prohibited) and Annex III (High Risk).
	Indicators = []IntentIndicator{
		// --- UNACCEPTABLE RISK (Article 5) ---
		{
			Keywords: []string{
				"social_score", "social_scoring", "citizen_trust", "citizen_rating",
				"behavioral_manipulation", "subliminal_exploit", "vulnerability_exploit",
				"predictive_policing", "recidivism_score", "pre_crime",
				"emotion_recognition_employee", "emotion_recognition_workplace", "emotion_recognition_school",
				"biometric_categorization_sensitive", "realtime_biometric_public",
			},
			Tier:              config.RiskUnacceptable,
			Article:           "EU AI Act Article 5 (Prohibited AI Practices)",
			Justification:     "System triggers prohibited capabilities under Article 5 (e.g. social scoring, subliminal manipulation, or prohibited biometric categorization).",
			ProhibitedReason:  "Deployment or distribution of this system is banned within the European Union under Article 5.",
			RequiresHumanGate: true,
			NISTCategories:    []string{"GOVERN-1.1", "MAP-1.1", "MEASURE-2.5"},
			ISO42001Controls:  []string{"A.5.2 (Prohibited Use Verification)", "A.6.2 (Risk Assessment)"},
		},

		// --- HIGH-RISK (Annex III) - Employment & HR ---
		{
			Keywords: []string{
				"resume_screening", "cv_screening", "candidate_evaluation", "candidate_ranking",
				"hiring_decision", "job_applicant", "interview_analysis", "worker_monitoring",
				"employee_performance_score", "termination_decision",
			},
			Tier:              config.RiskHigh,
			Article:           "EU AI Act Annex III Section 4 (Employment, Workers Management)",
			Justification:     "AI used for recruitment, candidate screening, task allocation, or worker evaluation constitutes High-Risk Annex III Section 4.",
			RequiresHumanGate: true,
			NISTCategories:    []string{"GOVERN-2.1", "MAP-1.5", "MEASURE-1.1", "MANAGE-2.4"},
			ISO42001Controls:  []string{"A.6.1 (Risk Assessment)", "A.7.2 (Data Quality)", "A.8.4 (Human Oversight)"},
		},

		// --- HIGH-RISK (Annex III) - Credit, Banking & Essential Services ---
		{
			Keywords: []string{
				"credit_score", "credit_scoring", "loan_underwriting", "mortgage_approval",
				"debt_collection", "insurance_risk", "health_insurance_underwriting",
				"emergency_dispatch", "triage_priority",
			},
			Tier:              config.RiskHigh,
			Article:           "EU AI Act Annex III Section 5 (Access to Essential Private and Public Services)",
			Justification:     "AI used for creditworthiness evaluation, loan eligibility, or emergency service dispatch is classified High-Risk Annex III Section 5.",
			RequiresHumanGate: true,
			NISTCategories:    []string{"GOVERN-1.2", "MAP-2.2", "MEASURE-2.1", "MANAGE-1.3"},
			ISO42001Controls:  []string{"A.6.1 (Risk Management)", "A.7.3 (Fairness)", "A.8.4 (Human Oversight)"},
		},

		// --- HIGH-RISK (Annex III) - Critical Infrastructure & Healthcare ---
		{
			Keywords: []string{
				"critical_infra", "power_grid", "water_supply", "traffic_control",
				"medical_diagnosis", "clinical_decision", "radiation_dose",
			},
			Tier:              config.RiskHigh,
			Article:           "EU AI Act Annex III Section 2 (Critical Infrastructure & Safety Components)",
			Justification:     "AI systems used as safety components in critical infrastructure or medical diagnostics are classified High-Risk.",
			RequiresHumanGate: true,
			NISTCategories:    []string{"GOVERN-1.1", "MAP-1.1", "MEASURE-1.1", "MANAGE-1.1"},
			ISO42001Controls:  []string{"A.6.1 (Risk Management System)", "A.8.1 (Safety Lifecycle)"},
		},

		// --- HIGH-RISK (Annex III) - Education & Admissions ---
		{
			Keywords: []string{
				"student_admission", "exam_proctoring", "exam_grading", "student_evaluation",
			},
			Tier:              config.RiskHigh,
			Article:           "EU AI Act Annex III Section 3 (Education and Vocational Training)",
			Justification:     "AI systems used for admissions, grading, or evaluating learning outcomes are classified High-Risk.",
			RequiresHumanGate: true,
			NISTCategories:    []string{"GOVERN-2.1", "MAP-1.5", "MEASURE-2.1"},
			ISO42001Controls:  []string{"A.6.1 (Risk Assessment)", "A.8.4 (Human Oversight)"},
		},

		// --- SPECIFIC TRANSPARENCY (Article 50) - Chatbots & Synthetic Content ---
		{
			Keywords: []string{
				"chatbot", "customer_agent", "virtual_assistant", "support_bot", "assistant",
				"support_assistant", "conversational_agent", "synthetic_voice", "voice_clone",
				"deepfake", "avatar_generation", "content_generator", "dialogue_agent",
			},
			Tier:              config.RiskTransparency,
			Article:           "EU AI Act Article 50 (Transparency Obligations)",
			Justification:     "Interactive conversational agents or synthetic media generators require explicit end-user transparency disclosures under Article 50.",
			RequiresHumanGate: false,
			NISTCategories:    []string{"GOVERN-1.2", "MAP-1.2", "MANAGE-1.2"},
			ISO42001Controls:  []string{"A.9.1 (Transparency and Disclosure)"},
		},
	}
)

// ClassifyAsset evaluates a discovered AI asset against EU AI Act rules.
func ClassifyAsset(asset *scanner.DiscoveredAsset) scanner.RiskAssessment {
	// Search in IntentKeywords, Function/Variable context, SystemPrompt, and Name
	searchCorpus := strings.ToLower(strings.Join([]string{
		asset.Name,
		asset.ContextSnippet,
		asset.SystemPrompt,
		strings.Join(asset.IntentKeywords, " "),
		strings.Join(asset.PromptVars, " "),
	}, " "))

	// 1. Evaluate indicators in precedence order (Unacceptable -> High -> Transparency)
	for _, ind := range Indicators {
		for _, kw := range ind.Keywords {
			if strings.Contains(searchCorpus, strings.ToLower(kw)) {
				return scanner.RiskAssessment{
					Tier:              ind.Tier,
					ApplicableArticle: ind.Article,
					Justification:     ind.Justification + " Triggered by context identifier: '" + kw + "'.",
					ProhibitedReason:  ind.ProhibitedReason,
					RequiresHumanGate: ind.RequiresHumanGate,
					NISTCategories:    ind.NISTCategories,
					ISO42001Controls:  ind.ISO42001Controls,
				}
			}
		}
	}

	// 2. Check if it's an interactive conversational model (default to transparency if model is conversational)
	lowerModel := strings.ToLower(asset.ModelFamily)
	if strings.Contains(lowerModel, "chat") || strings.Contains(lowerModel, "turbo") ||
		strings.Contains(lowerModel, "gpt") || strings.Contains(lowerModel, "claude") ||
		strings.Contains(lowerModel, "sonnet") || strings.Contains(lowerModel, "opus") ||
		strings.Contains(lowerModel, "gemini") {
		return scanner.RiskAssessment{
			Tier:              config.RiskTransparency,
			ApplicableArticle: "EU AI Act Article 50 (Transparency Obligations)",
			Justification:     "Conversational Foundation Model API in customer-facing or interactive context. Article 50 requires transparency notices.",
			RequiresHumanGate: false,
			NISTCategories:    []string{"GOVERN-1.2", "MANAGE-1.2"},
			ISO42001Controls:  []string{"A.9.1 (Transparency and Disclosure)"},
		}
	}

	// 3. Fallback to Minimal Risk (Article 69 voluntary codes of conduct)
	return scanner.RiskAssessment{
		Tier:              config.RiskMinimal,
		ApplicableArticle: "EU AI Act Article 69 / Recital 68 (Minimal Risk AI)",
		Justification:     "Internal developer utility, embedded search, or backend utility without direct high-risk domain impact.",
		RequiresHumanGate: false,
		NISTCategories:    []string{"GOVERN-1.1", "MAP-1.1"},
		ISO42001Controls:  []string{"A.5.1 (AI Management)", "A.8.1 (AI Lifecycle)"},
	}
}
