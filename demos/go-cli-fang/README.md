# Golang expert

## Prerequisites

### Model
```bash
docker model pull hf.co/unsloth/Qwen2.5-Coder-3B-Instruct-128K-GGUF:Q4_K_M
```

### Settings

`budgie.config.json`:
```json
{
  "model": "hf.co/unsloth/qwen2.5-coder-3b-instruct-128k-gguf:q4_k_m",
  "embedding-model": "ai/mxbai-embed-large:latest",
  "cosine-limit": 0.4,
  "temperature": 0.8,
  "baseURL": "http://localhost:12434/engines/llama.cpp/v1"
}
```

`budgie.system.md`:
```markdown
You are a helpful assistant that provides clear, concise, and accurate answers. 
You are a Golang expert. 

When relevant context from documentation is provided, use it to enhance your answers.
Be specific and reference the documentation when appropriate.
```

## Query

```bash
budgie ask --question "how to write a string to a file in Go?"
```

## Generate Embeddings
To generate embeddings for Go files, you can use the following command:
```bash
budgie generate-embeddings --chunk-size 1024 --overlap 512 --extension ".go" 
```

### Query with RAG
To query with RAG (Retrieval-Augmented Generation) using the generated embeddings, you can use:
```bash
budgie ask --question "how to write a string to a file in Go?" --rag
```