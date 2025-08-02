package services

import (
	"ai-support-bot-service/config"
	"ai-support-bot-service/jobs"
	"ai-support-bot-service/models"
	"context"
	"errors"
	"log"
	"net/http"

	chroma "github.com/amikos-tech/chroma-go/pkg/api/v2"

	"github.com/google/uuid"
)

func IngestData(fileId string, applicationID string) ([]map[string]string, int, error) {
	file, status, err := GetFile(fileId)
	if err != nil {
		return nil, status, err
	}
	if file.Ingested {
		return nil, http.StatusBadRequest, errors.New("file already ingested")
	}
	var job models.Job
	job.ID = uuid.New().String()
	job.FileID = fileId
	job.ApplicationID = applicationID
	job.Status = "Initializing"
	job.Log = []string{"Initializing"}
	_, err = config.JobCollection.InsertOne(context.Background(), &job)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	go jobs.IngestJob(job.ID, applicationID, file.FileBytes, fileId)
	return []map[string]string{{"job_id": job.ID}}, http.StatusOK, nil
}

func RetrieveData(query string, applicationID string) ([]chroma.Documents, int, error) {
	client := *config.ChromaConfig.Client
	col, err := client.GetOrCreateCollection(
		context.Background(),
		"context",
		chroma.WithEmbeddingFunctionCreate(&config.ChromaConfig.EmbeddingFunction),
		chroma.WithCollectionMetadataCreate(
			chroma.NewMetadata(
				chroma.NewStringAttribute("application_id", applicationID),
			),
		),
	)
	if err != nil {
		log.Fatalf("Query error: %v", err)
		return nil, http.StatusInternalServerError, errors.New("Query error: " + err.Error())
	}
	log.Println("Querying with query: ", query, " and applicationID: ", applicationID)
	filter := chroma.EqString("application_id", applicationID)
	res, err := col.Query(context.Background(),
		chroma.WithQueryTexts(query),
		chroma.WithNResults(5),
		chroma.WithIncludeQuery(chroma.IncludeDocuments, chroma.IncludeMetadatas),
		chroma.WithWhereQuery(
			filter,
		),
	)
	if err != nil {
		log.Fatalf("Query error: %v", err)
		return nil, http.StatusInternalServerError, errors.New("Query error: " + err.Error())
	}
	log.Println("date retrieved", res.GetMetadatasGroups())
	return res.GetDocumentsGroups(), http.StatusOK, nil
}
