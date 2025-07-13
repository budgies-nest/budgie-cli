# Sample Documentation

This is a sample markdown file for testing the embedding generation.

## Introduction

Budgie CLI is a powerful tool for AI-powered conversations.

## Features

### Interactive Mode
- TUI interface with huh library
- Continuous conversation
- Memory retention

### Embedding Generation
- Parse markdown files
- Create semantic chunks
- Generate vector embeddings

## Installation

Run the following commands:

```bash
go build -o budgie .
./budgie ask -p
```

## Configuration

Edit your `.budgie/budgie.config.json` file:

```json
{
  "model": "ai/qwen2.5:latest",
  "embedding-model": "ai/mxbai-embed-large:latest",
  "temperature": 0.8,
  "baseURL": "http://localhost:12434/engines/llama.cpp/v1"
}
```