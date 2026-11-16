# EU AI Act Regulatory Classification Taxonomy

governmesh implements a deterministic decision engine for classifying AI workflows according to the four risk tiers established in the European Union Artificial Intelligence Act.

## 1. Unacceptable Risk (Article 5: Prohibited AI Practices)

AI systems in this tier are banned under Article 5. governmesh treats any detected practice in this category as an immediate CI/CD hard block (Exit Code 1).

### Key Covered Practices
- **Subliminal and Behavioral Manipulation**: Systems deploying subliminal techniques beyond personal consciousness to distort human behavior in a manner that causes physical or psychological harm.
- **Vulnerability Exploitation**: Exploitation of age, disability, or social/economic vulnerability.
- **Social Scoring**: Trustworthiness classification of individuals based on social behavior or personality traits leading to detrimental treatment.
- **Biometric Categorization of Sensitive Attributes**: Systems that deduce race, political opinions, trade union membership, religious beliefs, sex life, or sexual orientation.
- **Individual Predictive Policing**: Risk assessments based solely on profiling or past criminal behavior.
- **Real-time Remote Biometric Identification**: In publicly accessible spaces for law enforcement (subject to statutory judicial exceptions).
- **Workplace and School Emotion Recognition**: Unless deployed for safety or medical purposes.

## 2. High-Risk (Article 6 and Annex III)

Permitted under strict conformity obligations, including risk management systems, data governance, technical documentation, record-keeping, transparency, and human oversight.

### Key Covered Domains (Annex III)
- **Biometrics (Annex III.1)**: Remote biometric identification, emotion recognition, biometric categorization.
- **Critical Infrastructure (Annex III.2)**: Safety components in digital infrastructure, road traffic, water, gas, heating, or electricity.
- **Education and Vocational Training (Annex III.3)**: Admissions, allocation, grading, and monitoring exam behavior.
- **Employment and Worker Management (Annex III.4)**: Recruitment, resume and CV screening, job ad targeting, candidate evaluation, task allocation, worker monitoring, and performance evaluation.
- **Access to Essential Private and Public Services (Annex III.5)**: Creditworthiness evaluation, loan approval, risk evaluation for life and health insurance, and emergency service dispatch prioritization.
- **Law Enforcement (Annex III.6)**: Individual risk assessments, polygraphs, reliability evaluation of evidence.
- **Migration and Border Control (Annex III.7)**: Polygraphs, verification of travel documents, visa and asylum application examination.
- **Administration of Justice and Democratic Processes (Annex III.8)**: Assisting judicial authorities in researching facts and interpreting law.

## 3. Specific Transparency Risk (Article 50)

Permitted with clear end-user disclosure obligations.

### Key Covered Systems
- **Direct Interaction with Humans**: Conversational chatbots, virtual customer agents, support bots.
- **Synthetic Content Generation**: Generation or manipulation of image, audio, or video content (deepfakes).
- **Text Generation for Public Dissemination**: Automated text generators published without verified human review.

## 4. Minimal and Low Risk (Article 69)

Permitted freely without mandatory statutory conformity. Encouraged to adopt voluntary codes of conduct.

### Key Covered Applications
- Internal developer tooling and code completion
- AI-enabled search indexing
- Spam filtering and email sorting
- Automated translation and spell-checking
- Internal document summarization without high-risk domain impact
