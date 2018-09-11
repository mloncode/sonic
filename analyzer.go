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
		log.Infof("got change")
		change := changes.Change()

		var baseNodes []sonicNode
		var headNodes []sonicNode

		if change.Base != nil && change.Base.UAST != nil {
			baseNodes = toSonicNodes(change.Base.UAST)
		}

		if change.Head != nil && change.Head.UAST != nil {
			headNodes = toSonicNodes(change.Head.UAST)
		}

		deleted, added, changed := diffNodes(baseNodes, headNodes)
		printNodes("deleted:", deleted)
		printNodes("added:", added)
		printNodes("changed:", changed)
	}

	if changes.Err() != nil {
		log.Errorf(changes.Err(), "failed to get a file from DataServer")
	}

	return &lookout.EventResponse{}, nil
}

func printNodes(header string, nodes []sonicNode) {
	fmt.Println(header)
	for _, n := range nodes {
		fmt.Println(n)
	}
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

func (n *sonicNode) Key() string {
	return fmt.Sprintf("%s%s", n.Type, n.Token)
}

func diffNodes(oldNodes, newNodes []sonicNode) ([]sonicNode, []sonicNode, []sonicNode) {
	oldMap := make(map[string]sonicNode, len(oldNodes))
	for _, n := range oldNodes {
		oldMap[n.Key()] = n
	}

	newMap := make(map[string]sonicNode, len(newNodes))
	for _, n := range newNodes {
		newMap[n.Key()] = n
	}

	var deletedNodes []sonicNode
	for key, n := range oldMap {
		if _, ok := newMap[key]; !ok {
			deletedNodes = append(deletedNodes, n)
		}
	}

	var addNodes []sonicNode
	var modifiedNodes []sonicNode
	for key, n := range newMap {
		if oldN, ok := oldMap[key]; !ok {
			addNodes = append(addNodes, n)
		} else if oldN.Lenght != n.Lenght {
			modifiedNodes = append(modifiedNodes, n)
		}
	}

	return deletedNodes, addNodes, modifiedNodes
}

// we don't need code below but need to satisfy interface

func (a *Analyzer) NotifyPushEvent(ctx context.Context, e *lookout.PushEvent) (*lookout.EventResponse, error) {
	return &lookout.EventResponse{}, nil
}
