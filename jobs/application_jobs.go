package jobs

import (
	"ai-support-bot-service/config"
	"context"
	"fmt"
	"log"
	"sync"

	chroma "github.com/amikos-tech/chroma-go/pkg/api/v2"
	"github.com/tmc/langchaingo/textsplitter"
	"go.mongodb.org/mongo-driver/bson"
)

type JobStatus string

const (
	JobStatusInitializing   JobStatus = "Initializing"
	JobStatusSplitting      JobStatus = "Splitting"
	JobStatusFailedToSplit  JobStatus = "Failed to split"
	JobStatusEmbedding      JobStatus = "Embedding"
	JobStatusFailedToEmbed  JobStatus = "Failed to embed"
	JobStatusIngesting      JobStatus = "Ingesting"
	JobStatusFailedToIngest JobStatus = "Failed to ingest"
	JobStatusFailed         JobStatus = "Failed"
	JobStatusSuccess        JobStatus = "Success"
)

var (
	jobStatus = make(map[string]string)
	jobMu     = sync.RWMutex{}
)

func setJobStatus(jobID string, status JobStatus) {
	jobMu.Lock()
	defer jobMu.Unlock()
	_, err := config.JobCollection.UpdateOne(context.Background(), bson.D{{Key: "_id", Value: jobID}}, bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: string(status)}}}})
	if err != nil {
		log.Println(err.Error())
	}
	jobStatus[jobID] = string(status)
}

func GetJobStatus(jobID string) string {
	jobMu.RLock()
	defer jobMu.RUnlock()
	return jobStatus[jobID]
}

func IngestJob(jobID string, applicationID string, fileBytes []byte, fileId string) {
	setJobStatus(jobID, JobStatusInitializing)
	log.Println("Ingesting job started")
	splitter := textsplitter.NewRecursiveCharacter()
	splitter.ChunkSize = 1000
	splitter.ChunkOverlap = 200
	setJobStatus(jobID, JobStatusSplitting)
	log.Println("Splitting started")
	chunks, err := splitter.SplitText(string(fileBytes))
	if err != nil {
		setJobStatus(jobID, JobStatusFailedToSplit)
		log.Println("Failed to split")
		return
	}
	setJobStatus(jobID, JobStatusEmbedding)
	log.Println("Embedding started")
	ef := config.ChromaConfig.EmbeddingFunction
	resp, err := ef.EmbedDocuments(context.Background(), chunks)
	if err != nil {
		setJobStatus(jobID, JobStatusFailedToEmbed)
		log.Println("Failed to embed")
		return
	}
	setJobStatus(jobID, JobStatusIngesting)
	log.Println("Ingesting started")
	client := *config.ChromaConfig.Client
	col, err := client.GetOrCreateCollection(context.Background(), "context", chroma.WithEmbeddingFunctionCreate(&ef),
		chroma.WithCollectionMetadataCreate(
			chroma.NewMetadata(
				chroma.NewStringAttribute("application_id", applicationID),
			),
		))
	if err != nil {
		// fmt.Println(err)
		setJobStatus(jobID, JobStatusFailedToIngest)
		log.Println("Failed to ingest")
		return
	}
	ids := make([]chroma.DocumentID, len(chunks))
	docs := make([]string, len(chunks))
	metadatas := make([]chroma.DocumentMetadata, len(chunks))
	for i, chunk := range chunks {
		ids[i] = chroma.DocumentID(fmt.Sprintf("chunk-%d", i))
		docs[i] = chunk
		metadatas[i] = chroma.NewDocumentMetadata(chroma.NewStringAttribute("application_id", applicationID))
	}
	_ = col.Add(context.Background(),
		chroma.WithIDs(ids...),
		chroma.WithTexts(docs...),
		chroma.WithEmbeddings(resp...),
		chroma.WithMetadatas(metadatas...),
	)
	setJobStatus(jobID, JobStatusSuccess)
	_, err = config.FileCollection.UpdateOne(context.Background(), bson.D{{Key: "_id", Value: fileId}}, bson.D{{Key: "$set", Value: bson.D{{Key: "ingested", Value: true}}}})
	if err != nil {
		setJobStatus(jobID, JobStatusFailed)
		log.Println(err.Error())
		return
	}
	log.Println("Ingestion completed")
}
