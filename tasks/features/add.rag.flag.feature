for the ask command, add a new --rag flag (false by default) to enable RAG (Retrieval-Augmented Generation) mode. 
When this flag is set to true, the command will use the RAG system to retrieve relevant context from the embeddings before generating a response.
if rag is activated you do not the need to specify `#rag` in the query.