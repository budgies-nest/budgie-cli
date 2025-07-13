in the .budgie directory if there is a docs directory (you can specify an alternative path with the --docs option), you can use a new command: budgie generate-embeddings

This command will do this:
- parse all markdown files in the docs directory: you can use the helpers from https://github.com/budgies-nest/budgie/blob/main/helpers/files.go
- for every markdown file, create chunks with ChunkWithMarkdownHierarchy from https://github.com/budgies-nest/budgie/blob/main/rag/chunks.md.advanced.go (there is test code for this: https://github.com/budgies-nest/budgie/blob/main/rag/chunks_md_test.go)
- for every chunk, create an embedding with the CreateAndSaveEmbeddingFromText method of the agent like in https://github.com/budgies-nest/budgie/blob/main/cookbook/13-persist-rag-memory/create_embeddings_file_test.go
- and save the result in the .budgie directory in a file named embeddings.json using the PersistMemoryVectorStore method of the agent.
- you can specifiy the name of the json file with agents.WithMemoryVectorStore("embeddings.json"),
- the model to use is specified in budgie.config.json with the "embedding-model" key.
- you need to create a specific agent for this. Let's call it "budgie-search"

--

in the .budgie directory if there is a docs directory (you can specify an alternative path with the --docs option), the programm will automatically search for similarity in the embeddings.json file.

budgie will use (and load) the embeddings.json file and the model specified in the "embedding-model" key of the budgie.config.json file to search for similarity in the embeddings.json file.

from the user message, first the budgie-search agent will do a similarity search in the embeddings.json file, like in https://github.com/budgies-nest/budgie/blob/main/cookbook/13-persist-rag-memory/main.go using the RAGMemorySearchSimilaritiesWithText method of the agent.
and then it will add the results (if there are results) to the message history of the main agent (budgie agent). before running the chat completion.
The cosine limit (of the RAGMemorySearchSimilaritiesWithText method) is set to 0.7 by default, but you can change it in the budgie.config.json file with the "cosine-limit" key.

the similarity search is triggered when the user message exists (from the interactive mode or with the -q option) and the embeddings.json file exists in the .budgie directory.