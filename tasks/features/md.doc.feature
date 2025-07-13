in the .budgie directory if there is a docs directory (you can specify an alternative path with the --docs option), you can use a new command: budgie generate-embeddings

This command will do this:
- parse all markdown files in the docs directory: you can use the helpers from https://github.com/budgies-nest/budgie/blob/main/helpers/files.go
- for every markdown file, create chunks with ChunkWithMarkdownHierarchy from https://github.com/budgies-nest/budgie/blob/main/rag/chunks.md.advanced.go (there is test code for this: https://github.com/budgies-nest/budgie/blob/main/rag/chunks_md_test.go)
- for every chunk, create an embedding with the CreateAndSaveEmbeddingFromText method of the agent like in https://github.com/budgies-nest/budgie/blob/main/cookbook/13-persist-rag-memory/create_embeddings_file_test.go
- and save the result in the .budgie directory in a file named embeddings.json using the PersistMemoryVectorStore method of the agent.
- you can specifiy the name of the json file with agents.WithMemoryVectorStore("embeddings.json"),
- the model to use is specified in budgie.config.json with the "embedding-model" key.
- you need to create a specific agent for this. Let's call it "budgie-search"