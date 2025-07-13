# Documentation

This directory contains documentation files that can be used with Budgie CLI's RAG (Retrieval Augmented Generation) functionality.

## Getting Started

1. Add your markdown documentation files to this directory
2. Generate embeddings: `budgie generate-embeddings`
3. Use RAG search in your questions: `budgie ask -q "#rag your question here"`

## Organization

You can organize your documentation using subdirectories. The embedding generation will recursively find all `.md` files.

Example structure:
```
docs/
├── README.md
├── api/
│   ├── endpoints.md
│   └── authentication.md
└── guides/
    ├── getting-started.md
    └── deployment.md
```

## Tips

- Use clear headings and sections in your markdown files
- Include relevant keywords that users might search for
- Keep content focused and well-structured
- Update embeddings after adding new documentation: `budgie generate-embeddings`