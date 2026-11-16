FROM vllm/vllm-openai:latest

ENV MODEL_NAME=mistralai/Mistral-7B-Instruct-v0.3
ENV HOST=0.0.0.0
ENV PORT=8000

CMD ["python3", "-m", "vllm.entrypoints.openai.api_server"]
