package sonic

import (
	"context"
	"fmt"

	"github.com/src-d/lookout"
	"gopkg.in/bblfsh/client-go.v2/tools"
	"gopkg.in/bblfsh/sdk.v1/uast"
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
		WantLanguage:    true,
		WantUAST:        true,
		ExcludeVendored: true,
	})

	if err != nil {
		log.Errorf(err, "failed to GetChanges from a DataService")
		return nil, err
	}

	for changes.Next() {
		change := changes.Change()

		log.Infof("got change")
		if change.Head == nil || change.Base == nil {
			continue
		}

		fmt.Println("old uast nodes:")
		for _, n := range toSonicNodes(change.Base.UAST) {
			fmt.Println(n)
		}
		fmt.Println("new uast nodes:")
		for _, n := range toSonicNodes(change.Head.UAST) {
			fmt.Println(n)
		}
	}

	if changes.Err() != nil {
		log.Errorf(changes.Err(), "failed to get a file from DataServer")
	}

	return &lookout.EventResponse{}, nil
}

var uastQuery = "//*[@roleDeclaration and @roleFunction and @startOffset and @endOffset]"

func toSonicNodes(node *uast.Node) []sonicNode {
	nodes, err := tools.Filter(node, uastQuery)
	if err != nil {
		return nil
	}

	var result []sonicNode

	for _, node := range nodes {
		length := node.EndPosition.Offset - node.StartPosition.Offset
		if length == 0 {
			continue
		}

		result = append(result, sonicNode{
			Type:   node.InternalType,
			Token:  getFirstToken(node),
			Lenght: length,
		})
	}

	return result
}

func getFirstToken(node *uast.Node) string {
	var n *uast.Node
	nodesToVisit := []*uast.Node{node}

	for len(nodesToVisit) > 0 {
		n, nodesToVisit = nodesToVisit[0], nodesToVisit[1:]

		for i := 0; i < len(n.Children); i++ {
			nodesToVisit = append(nodesToVisit, n.Children...)
		}

		if hasRole(n.Roles, uast.Name) && n.Token != "" {
			return n.Token
		}
	}

	return ""
}

func hasRole(roles []uast.Role, role uast.Role) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

type sonicNode struct {
	Type   string
	Token  string
	Lenght uint32
}

func (a *Analyzer) NotifyPushEvent(ctx context.Context, e *lookout.PushEvent) (*lookout.EventResponse, error) {
	return &lookout.EventResponse{}, nil
}
