package sonic

import (
	"context"
	"fmt"

	"github.com/src-d/lookout"
	log "gopkg.in/src-d/go-log.v1"
)

type Analyzer struct {
	DataClient *lookout.DataClient
}

var _ lookout.AnalyzerServer = &Analyzer{}

func (a *Analyzer) NotifyReviewEvent(ctx context.Context, e *lookout.ReviewEvent) (*lookout.EventResponse, error) {
	changes, err := a.DataClient.GetChanges(ctx, &lookout.ChangesRequest{
		Head:            &e.Head,
		Base:            &e.Base,
		WantContents:    true,
		WantLanguage:    false,
		WantUAST:        false,
		ExcludeVendored: true,
	})

	if err != nil {
		log.Errorf(err, "failed to GetChanges from a DataService")
		return nil, err
	}

	for changes.Next() {
		change := changes.Change()
		if change.Head == nil || change.Base == nil {
			continue
		}

		fmt.Println("change arrived", change.Base.Hash, "->", change.Head.Hash)
	}

	if changes.Err() != nil {
		log.Errorf(changes.Err(), "failed to get a file from DataServer")
	}

	return &lookout.EventResponse{
		Comments: []*lookout.Comment{
			{Text: "it works"},
		},
	}, nil
}

func (a *Analyzer) NotifyPushEvent(ctx context.Context, e *lookout.PushEvent) (*lookout.EventResponse, error) {
	return &lookout.EventResponse{}, nil
}
