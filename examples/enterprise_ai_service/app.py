import os
from openai import OpenAI
from langchain_openai import ChatOpenAI
from langchain_anthropic import ChatAnthropic

client = OpenAI(api_key=os.environ.get("OPENAI_API_KEY"))

# 1. Annex III High-Risk: Employment & Resume Screening
def evaluate_candidate_cv(resume_text: str, candidate_ssn: str, salary_expectation: float):
    prompt = f"""
    You are an automated employment evaluator. Evaluate applicant suitability.
    Candidate National ID: {candidate_ssn}
    Expected Salary: {salary_expectation}
    Resume Content: {resume_text}
    """
    response = client.chat.completions.create(
        model="gpt-4o",
        temperature=0.2,
        messages=[
            {"role": "system", "content": "You are an automated recruiting and hiring decision engine."},
            {"role": "user", "content": prompt}
        ]
    )
    return response.choices[0].message.content

# 2. Annex III High-Risk: Credit Underwriting
def calculate_credit_score_eligibility(user_income: float, credit_score: int, bank_account: str):
    llm = ChatOpenAI(model="gpt-4o-mini", temperature=0.0)
    system_prompt = "You evaluate loan risk, mortgage eligibility, and creditworthiness."
    prompt = f"Analyze creditworthiness for applicant with score {credit_score} and account {bank_account}."
    return llm.invoke(prompt)

# 3. Article 50: Transparency Chatbot
def customer_support_chatbot(user_message: str):
    llm = ChatAnthropic(model="claude-3-5-sonnet", temperature=0.7)
    system_prompt = "You are a customer assistant answering billing queries."
    return llm.invoke(user_message)
