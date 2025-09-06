package search

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/query"
	"github.com/mohan2020coder/mSpace/internal/models"
)

type BleveEvent struct {
	Date  string `json:"Date"`
	Event string `json:"Event"`
}

type BleveDoc struct {
	ID           uint         `json:"ID"`
	Title        string       `json:"Title"`
	Author       string       `json:"Author"`
	Abstract     string       `json:"Abstract"`
	FullText     string       `json:"FullText"`
	CollectionID uint         `json:"CollectionID"`
	Visibility   string       `json:"Visibility"`
	Petitioners  []string     `json:"Petitioners"`
	Respondents  []string     `json:"Respondents"`
	Events       []BleveEvent `json:"Events"`
	Synopsis     string       `json:"Synopsis"`
}

type SearchIndex struct {
	Index bleve.Index
}

func NewIndex(path string) (*SearchIndex, error) {
	var idx bleve.Index
	var err error

	if _, err = os.Stat(path); os.IsNotExist(err) {
		mapping := bleve.NewIndexMapping()
		idx, err = bleve.New(path, mapping)
		if err != nil {
			return nil, err
		}
	} else {
		idx, err = bleve.Open(path)
		if err != nil {
			return nil, err
		}
	}

	return &SearchIndex{Index: idx}, nil
}

func (s *SearchIndex) IndexItem(item *models.Item) error {
	bleveDoc := BleveDoc{
		ID:           item.ID,
		Title:        item.Title,
		Author:       item.Author,
		Abstract:     item.Abstract,
		FullText:     item.FullText,
		CollectionID: item.CollectionID,
		Visibility:   item.Visibility,
	}

	if item.LegalJSON != "" {
		var ld LegalDocument
		if err := json.Unmarshal([]byte(item.LegalJSON), &ld); err == nil {
			bleveDoc.Petitioners = ld.Petitioners
			bleveDoc.Respondents = ld.Respondents
			for _, e := range ld.Events {
				bleveDoc.Events = append(bleveDoc.Events, BleveEvent(e))
			}
			bleveDoc.Synopsis = ld.Synopsis
		}
	}

	return s.Index.Index(fmt.Sprintf("%d", item.ID), bleveDoc)
}

func (s *SearchIndex) Search(queryStr string, collectionID uint, author string) ([]uint, error) {
	var queries []query.Query

	if queryStr == "*" || queryStr == "" {
		queries = append(queries, bleve.NewMatchAllQuery())
	} else {
		fields := []string{"Title", "Abstract", "FullText", "Petitioners", "Respondents", "Events.Event", "Synopsis"}
		for _, f := range fields {
			q := bleve.NewMatchQuery(queryStr)
			q.SetField(f)
			queries = append(queries, q)
		}
	}

	if collectionID > 0 {
		val := float64(collectionID)
		numQuery := bleve.NewNumericRangeQuery(&val, &val)
		numQuery.SetField("CollectionID")
		queries = append(queries, numQuery)
	}

	if author != "" {
		authorQuery := bleve.NewMatchQuery(author)
		authorQuery.SetField("Author")
		queries = append(queries, authorQuery)
	}

	var finalQuery query.Query
	if len(queries) == 1 {
		finalQuery = queries[0]
	} else {
		finalQuery = bleve.NewDisjunctionQuery(queries...)
	}

	searchRequest := bleve.NewSearchRequestOptions(finalQuery, 100, 0, false)
	searchResult, err := s.Index.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	var ids []uint
	for _, hit := range searchResult.Hits {
		var id uint
		fmt.Sscanf(hit.ID, "%d", &id)
		ids = append(ids, id)
	}

	return ids, nil
}
