## Similarity search with RAG


### With Delimiter-based chunking

Generate embeddings using delimiter-based chunking:
```bash
budgie generate-embeddings --delimiter "-----"
```

Test the embeddings with a query:
```bash
budgie ask --question "#rag give me the description of the hawaiian pizza"
budgie ask --question "#rag what is Mapo Tofu?"
budgie ask --question "#rag invent a new recipe of the Margherita inspired by the Sweet and Sour Pork"

budgie ask --prompt --rag # activate the rag mode
```