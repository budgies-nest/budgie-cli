# Budgie CLI

A CLI tool for AI-powered conversations that streams responses to the terminal using configurable models and system instructions.

## Usage

```bash
budgie ask --question "your question here"
```

### Commands

- `ask` - Ask a question to the AI agent

### Available Flags for `ask` command

- `-q, --question` - The question to ask the AI (required unless using --prompt)
- `-p, --prompt` - Interactive TUI prompt mode (alternative to --question)
- `-s, --system` (default: ".budgie/budgie.system.md") - Path to system instructions file
- `-c, --config` (default: ".budgie/budgie.config.json") - Path to configuration file
- `-o, --output` (default: ".") - Path where to generate result files
- `-g, --generate` (default: true) - Generate result file

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

The interactive mode will continue asking for new questions after each completion. Type `/bye` to exit the interactive session.

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

Get help:
```bash
budgie --help
budgie ask --help
```

## Requirements

- Docker model runner with `k33g/qwen2.5:0.5b-instruct-q8_0` model
- System instructions file (see `.budgie/budgie.system.md` for example)

## Configuration

The CLI uses a `.budgie/budgie.config.json` file to configure the LLM settings:

```json
{
  "model": "k33g/qwen2.5:0.5b-instruct-q8_0",
  "temperature": 0.8,
  "baseURL": "http://localhost:12434/engines/llama.cpp/v1"
}
```

- `model`: The LLM model to use
- `temperature`: Controls randomness in responses (0.0-1.0)
- `baseURL`: The base URL for the model runner

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
- Reads system instructions from a file
- Command-based interface with `ask` subcommand
- Takes user question as command line argument with `--question` flag
- Streams AI response to terminal in real-time
- Saves response to timestamped markdown file (`result-yyyy-mm-dd-hh-mm-ss.md`)
- Built with budgie agent framework