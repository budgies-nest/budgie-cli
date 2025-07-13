# Budgie CLI

A CLI tool for AI-powered conversations that streams responses to the terminal using configurable models and system instructions.

## Usage

```bash
budgie ask --question "your question here"
```

### Commands

- `ask` - Ask a question to the AI agent
- `generate-embeddings` - Generate embeddings from markdown files for RAG functionality

### Available Flags for `ask` command

- `-q, --question` - The question to ask the AI (required unless using --prompt)
- `-p, --prompt` - Interactive TUI prompt mode (alternative to --question)
- `-s, --system` (default: ".budgie/budgie.system.md") - Path to system instructions file
- `-c, --config` (default: ".budgie/budgie.config.json") - Path to configuration file
- `-o, --output` (default: ".") - Path where to generate result files
- `-g, --generate` (default: true) - Generate result file
- `-u, --use` - Path to file to include as additional system message

### Available Flags for `generate-embeddings` command

- `-c, --config` (default: ".budgie/budgie.config.json") - Path to configuration file
- `-d, --docs` (default: ".budgie/docs") - Path to docs directory containing markdown files

### Examples

Basic usage:
```bash
budgie ask --question "What is Go programming language?"
```

Using short flags:
```bash
budgie ask -q "What is Go programming language?"
```

Interactive TUI mode:
```bash
budgie ask --prompt
# or
budgie ask -p
```

The interactive mode will continue asking for new questions after each completion. Available commands:
- Type `/bye` to exit the interactive session
- Type `/clear` to reset the conversation and reload system instructions

Custom system file and config:
```bash
budgie ask --system custom-system.txt --config prod-budgie.config.json --question "Explain machine learning"
```

Save results to specific directory:
```bash
budgie ask --output ./results --question "What are design patterns?"
```

Disable result file generation:
```bash
budgie ask --generate=false --question "Quick question about Go"
```

Generate embeddings from markdown files:
```bash
budgie generate-embeddings
```

Generate embeddings from custom docs directory:
```bash
budgie generate-embeddings --docs ./my-docs
```

Use additional file as context:
```bash
budgie ask -u ./project-spec.md -q "How should I implement this feature?"
```

Get help:
```bash
budgie --help
budgie ask --help
budgie generate-embeddings --help
```

## Requirements

- Docker model runner with `k33g/qwen2.5:0.5b-instruct-q8_0` model
- System instructions file (see `.budgie/budgie.system.md` for example)

## Configuration

The CLI uses a `.budgie/budgie.config.json` file to configure the LLM settings:

```json
{
  "model": "ai/qwen2.5:latest",
  "embedding-model": "ai/mxbai-embed-large:latest",
  "cosine-limit": 0.6,
  "temperature": 0.8,
  "baseURL": "http://localhost:12434/engines/llama.cpp/v1"
}
```

- `model`: The LLM model to use for chat completions
- `embedding-model`: The model to use for generating embeddings (required for `generate-embeddings` command)
- `cosine-limit`: Similarity threshold for RAG search (default: 0.7, lower values = more results)
- `temperature`: Controls randomness in responses (0.0-1.0)
- `baseURL`: The base URL for the model runner

## RAG (Retrieval Augmented Generation) with Similarity Search

Budgie CLI includes intelligent document search capabilities that automatically enhance your conversations with relevant context from your documentation.

### How It Works

When you ask a question with the `#rag` prefix, Budgie:

1. **🔍 Searches** your embedded documentation for relevant content
2. **📚 Displays** the found documentation chunks in green
3. **🤖 Enhances** the AI response with this context

### Setting Up RAG

1. **Create documentation directory**:
   ```bash
   mkdir -p .budgie/docs
   ```

2. **Add your markdown files** to `.budgie/docs/`

3. **Generate embeddings**:
   ```bash
   budgie generate-embeddings
   ```

4. **Ask questions with RAG** - Use the `#rag` prefix to search for relevant context:
   ```bash
   budgie ask -q "#rag How do I configure the system?"
   # or use interactive mode
   budgie ask -p
   # then type: #rag How do I configure the system?
   ```

### What You'll See

When you use the `#rag` prefix and relevant documentation is found, you'll see:

```
What's your question? > #rag How do I configure budgie?

🔍 Searching... ✓
📚 Found 3 relevant documentation chunks:

   1. TITLE: ## Configuration
      HIERARCHY: User Guide > Configuration
      CONTENT: Edit your .budgie/budgie.config.json file...

   2. TITLE: ## Installation
      HIERARCHY: User Guide > Installation  
      CONTENT: Run the following commands...

   3. [Additional relevant chunks...]

[AI response using this context]
```

### Configuring Similarity Search

The `cosine-limit` setting in your config controls how strict the similarity matching is:

- **0.9**: Very strict - only highly relevant matches
- **0.7**: Balanced - good mix of relevance and coverage (default)
- **0.5**: Permissive - more results, some may be less relevant
- **0.3**: Very permissive - many results, varying relevance

Lower values return more documentation chunks but may include less relevant content.

### Benefits

- **Contextual Responses**: AI answers are enhanced with your specific documentation
- **Automatic Discovery**: No need to manually specify which docs to reference
- **Visual Transparency**: See exactly what documentation influenced the response
- **Memory Efficient**: Embeddings are loaded once per session for fast subsequent searches

## Interactive Mode Commands

When using interactive mode (`budgie ask -p`), you have access to special commands:

| Command | Description |
|---------|-------------|
| `/bye` | Exit the interactive session |
| `/clear` | Reset conversation history and reload system instructions from `budgie.system.md` |
| `/use <file-path>` | Load a file and add its content as an additional system message |
| `#rag <question>` | Search documentation and enhance response with relevant context |

### Using `/clear`

The `/clear` command is useful when you want to:
- Start a fresh conversation without exiting the session
- Apply changes made to your `budgie.system.md` file
- Reset the conversational context while keeping the same session

Example:
```
What's your question? > /clear
✅ Conversation cleared and system instructions reloaded

What's your question? > [fresh conversation starts here]
```

### Using `#rag`

The `#rag` prefix activates similarity search for that specific question:

```
What's your question? > Hello, how are you?
[Normal AI response without documentation search]

What's your question? > #rag How do I install budgie?
🔍 Searching... ✓
📚 Found 2 relevant documentation chunks:
[Documentation-enhanced response]
```

This gives you control over when to use RAG search versus having normal conversations.

### Using `/use`

The `/use` command loads a file and adds its content as a system message:

```
What's your question? > /use ./project-context.md
✅ File ./project-context.md loaded as system message

What's your question? > Now help me with my project
[AI response enhanced with the loaded file context]
```

You can load multiple files during a session:
```
What's your question? > /use ./api-docs.md
✅ File ./api-docs.md loaded as system message

What's your question? > /use ./coding-standards.md  
✅ File ./coding-standards.md loaded as system message

What's your question? > Write a function that follows our standards
[AI response using both loaded contexts]
```

**Command line equivalent:**
```bash
budgie ask -u ./project-context.md -q "Help me with my project"
```

## Installation

To make the `budgie` binary available from anywhere:

Use the provided installation script:
```bash
./install.sh
```

Or install manually:

**Option 1: Install to `/usr/local/bin` (recommended)**
```bash
sudo cp budgie /usr/local/bin/
sudo cp -r .budgie /usr/local/bin/
```

**Option 2: Install to `~/bin`**
```bash
mkdir -p ~/bin
cp budgie ~/bin/
cp -r .budgie ~/bin/
echo 'export PATH="$HOME/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

**Option 3: Create symlink**
```bash
sudo ln -s $(pwd)/budgie /usr/local/bin/budgie
```

After installation, you can run `budgie` from any directory.

## Features

- Configurable LLM model, temperature, and base URL
- Reads system instructions from a file (hot-reloadable in interactive mode)
- Command-based interface with `ask` subcommand
- Takes user question as command line argument with `--question` flag
- Streams AI response to terminal in real-time
- Saves response to timestamped markdown file (`result-yyyy-mm-dd-hh-mm-ss.md`)
- Built with budgie agent framework
- RAG (Retrieval Augmented Generation) support with automatic similarity search
- Intelligent document retrieval from embedded markdown documentation
- Visual feedback showing which documentation influenced responses