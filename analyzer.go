package sonic

import (
	"context"
	"crypto/sha1"
	"fmt"

	"github.com/MLonCode/sonic/src/sound"
	"github.com/rakyll/portmidi"

	"github.com/src-d/lookout"
	"gopkg.in/bblfsh/client-go.v2/tools"
	"gopkg.in/bblfsh/sdk.v1/uast"
	log "gopkg.in/src-d/go-log.v1"
)

type Analyzer struct {
	DataClient *lookout.DataClient
	DeviceID   portmidi.DeviceID
}

var _ lookout.AnalyzerServer = &Analyzer{}

var m1 = sound.NewMarkov("song1.midi")
var m2 = sound.NewMarkov("song2.midi")

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

	var total int

	for changes.Next() {
		log.Infof("got change")
		change := changes.Change()

		var baseNodes []sonicNode
		var headNodes []sonicNode

		if change.Base != nil && change.Base.UAST != nil {
			baseNodes = toSonicNodes(change.Base)
		}

		if change.Head != nil && change.Head.UAST != nil {
			headNodes = toSonicNodes(change.Head)
		}

		deleted, added, changed := diffNodes(baseNodes, headNodes)
		printNodes("deleted:", deleted)
		printNodes("added:", added)
		printNodes("changed:", changed)

		deletedSeq := sound.NewSequence("prophet", ConvertMarkov(m2, deleted))
		deletedSeq.Play(a.DeviceID)

		addedSeq := sound.NewSequence("prophet", ConvertMarkov(m2, added))
		addedSeq.Play(a.DeviceID)

		total += len(deleted) + len(added) + len(changed)
	}

	fmt.Println("total nodes:", total)

	if changes.Err() != nil {
		log.Errorf(changes.Err(), "failed to get a file from DataServer")
	}

	return &lookout.EventResponse{}, nil
}

type sonicNode struct {
	Type   string
	Token  string
	Lenght uint32
	Hash   [20]byte
}

func (n *sonicNode) Key() string {
	return fmt.Sprintf("%s%s", n.Type, n.Token)
}

func printNodes(header string, nodes []sonicNode) {
	fmt.Println(header)
	for _, n := range nodes {
		fmt.Println(n.Type, n.Token, n.Lenght)
	}
}

var uastQuery = "//*[(@roleDeclaration or @roleIdentifier or @roleLiteral) and @startOffset and @endOffset]"

func toSonicNodes(file *lookout.File) []sonicNode {
	nodes, err := tools.Filter(file.UAST, uastQuery)
	if err != nil {
		return nil
	}

	var result []sonicNode

	for _, node := range nodes {
		length := node.EndPosition.Offset - node.StartPosition.Offset
		if length == 0 {
			continue
		}

		token := getFirstToken(node)
		if token == "" {
			continue
		}

		content := file.Content[node.StartPosition.Offset:node.EndPosition.Offset]

		result = append(result, sonicNode{
			Type:   node.InternalType,
			Token:  token,
			Lenght: length,
			Hash:   sha1.Sum(content),
		})
	}

	return result
}

func getFirstToken(node *uast.Node) string {
	if node.Token != "" {
		return node.Token
	}

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
		} else if oldN.Hash != n.Hash {
			modifiedNodes = append(modifiedNodes, n)
		}
	}

	return deletedNodes, addNodes, modifiedNodes
}

// we don't need code below but need to satisfy interface

func (a *Analyzer) NotifyPushEvent(ctx context.Context, e *lookout.PushEvent) (*lookout.EventResponse, error) {
	return &lookout.EventResponse{}, nil
}
