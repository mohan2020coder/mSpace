📌 Golang Alternatives / Building Blocks for Digital Repositories

Since Go excels at APIs, microservices, and scalability, people often use it to build custom repositories or DAM (Digital Asset Management) systems.

1. Dataverse (Harvard) – not Go, but API-friendly

Written in Java, but has strong APIs.

You can build a Go-based frontend or services around it.

2. Go-based Digital Storage/Repository Components

Not direct DSpace replacements, but usable to build your own DSpace-like system:

MinIO (Go) – S3-compatible object storage (great for storing files, datasets, preservation).

SeaweedFS (Go) – distributed file system (for large collections).

Bleve (Go) – full-text search library.

Gorse (Go) – recommendation system (for suggesting related documents).

Cayley (Go) – graph database for relationships between authors, works, etc.

These can be stitched together into a repository system.

3. Community Projects

Some smaller open-source DAMS and CMS exist in Go, but not as institutionalized as DSpace.

Pachyderm (Go + data pipelines, but mainly ML/data science workflows).

Photoprism (Go-based photo management, can be adapted for digital collections).

Perkeep (formerly Camlistore, Go-based personal archival system).

4. Build Your Own Repository with Go

A typical Go replacement for DSpace would look like:

Storage: MinIO / SeaweedFS

Database: PostgreSQL / SQLite (metadata)

Indexing: Bleve / Elasticsearch

API: Go + REST/GraphQL

Frontend: React/Next.js

Protocols: Implement OAI-PMH, Dublin Core metadata export in Go

This way you get a DSpace-like system but fully in Go + modern stack.