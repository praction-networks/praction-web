package logger

import (
	"context"

	"github.com/olivere/elastic/v7"
	"go.uber.org/zap/zapcore"
)

type elkWriter struct {
	client      *elastic.Client
	searchIndex string
}

func getElkWriter(client *elastic.Client, index string) zapcore.WriteSyncer {
	return &elkWriter{
		client:      client,
		searchIndex: index,
	}
}

func (w *elkWriter) Write(p []byte) (n int, err error) {
	// Send the log message to Elasticsearch
	ctx := context.Background()
	_, err = w.client.Index().
		Index(w.searchIndex).
		BodyJson(string(p)).
		Do(ctx)

	if err != nil {
		return 0, err
	}

	return len(p), nil
}

func (w *elkWriter) Sync() error {
	// Optional sync method if needed
	return nil
}
