# 🦜 Budgie CLI

A CLI tool for AI-powered conversations that streams responses to the terminal using configurable models and system instructions.

## Usage

```bash
budgie ask --question "your question here"
```

### Commands

- `init` - Initialize a new Budgie CLI project with default configuration
- `ask` - Ask a question to the AI agent
- `generate-embeddings` - Generate embeddings from markdown files for RAG functionality

### Available Flags for `ask` command

- `-q, --question` - The question to ask the AI (required unless using --prompt or --from)
- `-p, --prompt` - Interactive TUI prompt mode (alternative to --question)
- `-f, --from` - Path to file containing the user question/message (alternative to --question)
- `-s, --system` (default: ".budgie/budgie.system.md") - Path to system instructions file
- `-c, --config` (default: ".budgie/budgie.config.json") - Path to configuration file
- `-o, --output` (default: ".") - Path where to generate result files
- `-g, --generate` (default: true) - Generate result file
- `-u, --use` - Path to file to include as additional system message

### Available Flags for `generate-embeddings` command

- `-c, --config` (default: ".budgie/budgie.config.json") - Path to configuration file
- `-d, --docs` (default: ".budgie/docs") - Path to docs directory containing markdown files

**Chunking Methods** (mutually exclusive):
- `-m, --markdown-hierarchy` - Use markdown hierarchy chunking (default)
- `-s, --markdown-sections` - Use markdown sections chunking
- `-D, --delimiter <string>` - Use delimiter-based chunking with specified delimiter
- `-z, --chunk-size <size>` - Use fixed-size text chunking with specified size
- `-o, --overlap <length>` - Overlap length for fixed-size chunking (requires --chunk-size)
- `-f, --files` - Use whole-file chunking (each file is one chunk)
- `-e, --extension <ext>` - File extension to process (for --delimiter, --chunk-size, and --files methods, default: .md)

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

Generate embeddings from markdown files (default hierarchy chunking):
```bash
budgie generate-embeddings
```

Generate embeddings from custom docs directory:
```bash
budgie generate-embeddings --docs ./my-docs
```

Generate embeddings using different chunking methods:
```bash
# Markdown hierarchy chunking (default) - preserves document structure
budgie generate-embeddings --markdown-hierarchy

# Markdown sections chunking - simpler section-based chunks
budgie generate-embeddings --markdown-sections

# Delimiter-based chunking - split by custom delimiter
budgie generate-embeddings --delimiter "---"
budgie generate-embeddings --delimiter "# "
budgie generate-embeddings --delimiter $'\n\n'  # Split by paragraphs

# Fixed-size text chunking - consistent chunk sizes
budgie generate-embeddings --chunk-size 1024
budgie generate-embeddings --chunk-size 1024 --overlap 256
budgie generate-embeddings --chunk-size 512 --overlap 128

# Whole-file chunking - each file becomes one chunk
budgie generate-embeddings --files
budgie generate-embeddings --files --extension ".py"
budgie generate-embeddings --files --extension "txt"
```

Use additional file as context:
```bash
budgie ask -u ./project-spec.md -q "How should I implement this feature?"
```

Read question from file:
```bash
budgie ask --from ./my-question.txt
# or
budgie ask -f ./my-question.txt
```

Read question from file in interactive mode (triggers immediate completion):
```bash
budgie ask --prompt --from ./my-question.txt
# After completion, continues with interactive mode
```

Initialize new project:
```bash
budgie init
```

Get help:
```bash
budgie --help
budgie ask --help
budgie generate-embeddings --help
budgie init --help
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

1. **Initialize project** (recommended for new projects):
   ```bash
   budgie init
   ```

   Or manually create documentation directory:
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

## Embeddings Generation Methods

Budgie CLI offers multiple chunking strategies to optimize your documentation for different use cases. Choose the method that best fits your content structure and search requirements.

### 1. Markdown Hierarchy Chunking (Default)

**Best for**: Structured documentation with clear hierarchical relationships

```bash
budgie generate-embeddings
# or explicitly:
budgie generate-embeddings --markdown-hierarchy
```

**How it works**:
- Preserves markdown header hierarchy (# ## ### etc.)
- Creates chunks with rich metadata including hierarchical context
- Each chunk includes: TITLE, HIERARCHY path, and CONTENT

**Output format**:
```
TITLE: ## Configuration
HIERARCHY: User Guide > Setup > Configuration
CONTENT: Edit your .budgie/budgie.config.json file...
```

**Use when**:
- Your docs have clear hierarchical structure
- You want to maintain context about where information appears
- You need the most sophisticated chunking with relationship awareness

### 2. Markdown Sections Chunking

**Best for**: Simpler markdown documents where you want section-based chunks without complex hierarchy

```bash
budgie generate-embeddings --markdown-sections
```

**How it works**:
- Splits content by markdown headers (# ## ### etc.)
- Creates clean sections without complex metadata
- Simpler than hierarchy method but still respects markdown structure

**Use when**:
- You have markdown docs but don't need complex hierarchy metadata
- You want faster processing with good structure awareness
- Your documents have flat or simple hierarchical structure

### 3. Delimiter-Based Chunking

**Best for**: Custom content organization or non-markdown documents

```bash
# Split by horizontal rules
budgie generate-embeddings --delimiter "---"

# Split by headers
budgie generate-embeddings --delimiter "# "

# Split by paragraphs (double newlines)
budgie generate-embeddings --delimiter $'\n\n'

# Split by custom markers
budgie generate-embeddings --delimiter "<!-- SECTION -->"
```

**How it works**:
- Splits text at every occurrence of the specified delimiter
- No structural awareness beyond the delimiter
- Flexible for any content type

**Use when**:
- You have custom content markers or separators
- Working with non-markdown documents
- You need precise control over where splits occur
- Your content uses consistent delimiter patterns

**Common delimiters**:
- `"---"` - Horizontal rules in markdown
- `"# "` - Top-level headers
- `$'\n\n'` - Paragraph breaks (use `$''` syntax for newlines)
- `"<!-- SECTION -->"` - HTML comments as markers
- `"\n\n---\n\n"` - Section dividers with spacing

### 4. Fixed-Size Text Chunking

**Best for**: Consistent chunk sizes regardless of content structure, especially for technical content or when optimizing for specific embedding model limits

```bash
# Basic fixed-size chunks
budgie generate-embeddings --chunk-size 1024

# With overlap for better context preservation
budgie generate-embeddings --chunk-size 1024 --overlap 256

# Smaller chunks for detailed content
budgie generate-embeddings --chunk-size 512 --overlap 128

# Larger chunks for overview content
budgie generate-embeddings --chunk-size 2048 --overlap 512
```

**How it works**:
- Splits text into chunks of exactly the specified character count
- Optional overlap preserves context between chunks
- No awareness of document structure - purely size-based

**Parameters**:
- `--chunk-size`: Number of characters per chunk (required)
- `--overlap`: Number of characters to overlap between chunks (optional)

**Use when**:
- You need consistent chunk sizes for embedding model optimization
- Working with very long documents that need systematic splitting
- Your content doesn't have clear structural markers
- You want predictable chunk counts for performance planning

**Overlap benefits**:
- Prevents important context from being lost at chunk boundaries
- Improves search recall for concepts that span chunk boundaries
- Recommended overlap: 20-25% of chunk size (e.g., 256 chars for 1024 chunk size)

### Choosing the Right Method

| Method | Structure Preservation | Processing Speed | Best For |
|--------|----------------------|------------------|----------|
| **Hierarchy** | Excellent | Moderate | Structured docs with clear hierarchy |
| **Sections** | Good | Fast | Simple markdown documents |
| **Delimiter** | Custom | Very Fast | Custom markers, non-markdown content |
| **Fixed-Size** | None | Fastest | Consistent sizes, technical optimization |

### Performance Comparison

Using a typical documentation file:

```bash
# Results may vary based on content structure
budgie generate-embeddings --docs ./example-docs

# Hierarchy: 65 chunks (rich metadata)
# Sections: 58 chunks (clean sections)  
# Delimiter "# ": 66 chunks (header-based)
# Delimiter $'\n\n': 78 chunks (paragraph-based)
# Chunk-size 1024: 16 chunks (consistent size)
# Chunk-size 512: 32 chunks (smaller, more chunks)
```

### Advanced Usage

**Combine with custom docs directory**:
```bash
budgie generate-embeddings --docs ./technical-specs --chunk-size 2048 --overlap 512
```

**Process different content types**:
```bash
# API documentation with clear sections
budgie generate-embeddings --docs ./api-docs --markdown-sections

# Legal documents with custom markers
budgie generate-embeddings --docs ./legal --delimiter "SECTION "

# Mixed content with size constraints
budgie generate-embeddings --docs ./mixed --chunk-size 1500 --overlap 300
```

**Validation and Error Handling**:
- Only one chunking method can be used at a time
- `--overlap` requires `--chunk-size` to be specified
- Overlap must be less than chunk size
- Clear error messages guide correct usage

## Interactive Mode Commands

When using interactive mode (`budgie ask -p`), you have access to special commands:

| Command | Description |
|---------|-------------|
| `/bye` | Exit the interactive session |
| `/clear` | Reset conversation history and reload system instructions from `budgie.system.md` |
| `/use <file-path>` | Load a file and add its content as an additional system message |
| `/from <file-path>` | Load a question from a file and process it immediately |
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

## Reading Questions from Files

Budgie CLI supports reading user questions/messages from files using the `--from` / `-f` flag. This is useful for:

- **Long, complex questions** that are easier to write in a text editor
- **Reusing questions** across multiple sessions
- **Scripting and automation** workflows
- **Template questions** that you use frequently

### Single Question Mode

Read a question from a file and get a response:

```bash
# Create a question file
echo "Explain the differences between Go interfaces and struct embedding" > question.txt

# Ask the question from the file
budgie ask --from question.txt
```

The file content supports all standard features:
- **RAG search**: Start the file content with `#rag ` to trigger similarity search
- **Multi-line questions**: The entire file content becomes the user message

Example question file with RAG:
```
#rag How do I configure the embedding model for RAG functionality?

I'm setting up a new project and want to understand:
1. Which embedding models are recommended
2. How to configure the cosine-limit setting
3. Best practices for organizing documentation
```

### Interactive Mode with File Input

Start interactive mode and immediately process a question from a file:

```bash
budgie ask --prompt --from question.txt
```

This will:
1. **Load and process** the file content immediately
2. **Display the AI response** 
3. **Continue with interactive mode** for follow-up questions

Example workflow:
```bash
# Start with a complex initial question from file
budgie ask -p -f ./project-analysis.txt

# After the AI responds, you can continue interacting:
# > "Can you elaborate on point #3?"
# > "What about performance considerations?"
# > "/bye" (to exit)
```

### Combining with Other Flags

The `--from` flag works seamlessly with other options:

```bash
# Use file question with custom system instructions and context
budgie ask --from question.txt --system custom-system.md --use project-context.md

# Interactive mode with file input and custom output directory
budgie ask -p -f question.txt -o ./results --generate=true
```

### File Format Tips

- **Plain text**: Any text file works (`.txt`, `.md`, etc.)
- **No special formatting required**: The entire file content becomes the user message
- **RAG prefix**: Start with `#rag ` for similarity search
- **Multi-line support**: Complex, formatted questions work perfectly

Example question files:

**simple-question.txt:**
```
What are the best practices for error handling in Go?
```

**complex-analysis.txt:**
```
#rag Please analyze our current codebase structure and provide recommendations.

I need help with:

## Current Situation
- Legacy monolith with 50+ microservices
- Mixed Go versions (1.19, 1.20, 1.21)
- Inconsistent error handling patterns

## Goals
- Standardize error handling
- Improve observability
- Reduce technical debt

## Constraints
- Zero-downtime deployment required
- Budget limitations for refactoring
- Team of 8 developers

Please provide a prioritized roadmap with specific steps.
```

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

### Using `/from`

The `/from` command loads a question from a file and processes it immediately, just like using the `--from` flag:

```
What's your question? > /from ./complex-question.txt
📁 Loaded question from file: ./complex-question.txt
[AI processes the file content and responds]

What's your question? > [continue with follow-up questions]
```

This is particularly useful for:
- **Complex questions** that are easier to write in a text editor
- **Reusing questions** from previous sessions
- **Template questions** that you use frequently

Example workflow:
```
What's your question? > /from ./code-review-checklist.txt
📁 Loaded question from file: ./code-review-checklist.txt
[AI provides detailed code review feedback]

What's your question? > Can you focus more on the security aspects?
[Follow-up conversation continues]

What's your question? > /from ./performance-analysis.txt
📁 Loaded question from file: ./performance-analysis.txt
[AI analyzes performance with fresh context]
```

The `/from` command supports all the same features as the `--from` flag:
- **RAG search**: If the file starts with `#rag `, it will trigger similarity search
- **Multi-line content**: Complex formatted questions work perfectly
- **Any text format**: Works with `.txt`, `.md`, or any plain text file

**Command line equivalent:**
```bash
budgie ask --from ./complex-question.txt
```

## Installation

To make the `budgie` binary available from anywhere:

Install it manually:

**Option 1: Install to `/usr/local/bin` (recommended)**
```bash
sudo cp budgie /usr/local/bin/
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
- Takes user question as command line argument with `--question` flag or from file with `--from` flag
- Interactive TUI prompt mode with conversation history
- Streams AI response to terminal in real-time
- Saves response to timestamped markdown file (`result-yyyy-mm-dd-hh-mm-ss.md`)
- Built with budgie agent framework
- RAG (Retrieval Augmented Generation) support with automatic similarity search
- Intelligent document retrieval from embedded markdown documentation
- Visual feedback showing which documentation influenced responses
- File-based question input for complex queries and automation