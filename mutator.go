package exclusion

import (
	"fmt"

	"entgo.io/contrib/entoas"
	"entgo.io/ent/entc/gen"
	"github.com/ogen-go/ogen"
	"golang.org/x/exp/slices"
)

// Mutator searches for skip annotations and removes the requested fields from the operation
func Mutator(graph *gen.Graph, spec *ogen.Spec) error {
	// compile ignore list
	ignores := make(map[string][]string)
	for _, node := range graph.Nodes {
		for _, field := range node.Fields {
			if ann, ok := field.Annotations["OasOperation"]; ok {
				opAnt := &Annotation{}
				if err := opAnt.Decode(ann); err != nil {
					continue
				}

				if opAnt.SkipCreate {
					addIgnore(ignores, entoas.OpCreate.Title(), node.Name, field.Name)
				}
				if opAnt.SkipDelete {
					addIgnore(ignores, entoas.OpDelete.Title(), node.Name, field.Name)
				}
				if opAnt.SkipUpdate {
					addIgnore(ignores, entoas.OpUpdate.Title(), node.Name, field.Name)
				}
				if opAnt.SkipList {
					addIgnore(ignores, entoas.OpList.Title(), node.Name, field.Name)
				}
				if opAnt.SkipRead {
					addIgnore(ignores, entoas.OpRead.Title(), node.Name, field.Name)
				}
			}
		}
	}

	for _, pathItem := range spec.Paths {
		if pathItem.Post != nil {
			parseIgnoreOperation(ignores, pathItem.Post)
		}
		if pathItem.Get != nil {
			parseIgnoreOperation(ignores, pathItem.Get)
		}
		if pathItem.Patch != nil {
			parseIgnoreOperation(ignores, pathItem.Patch)
		}
		if pathItem.Delete != nil {
			parseIgnoreOperation(ignores, pathItem.Delete)
		}
		if pathItem.Head != nil {
			parseIgnoreOperation(ignores, pathItem.Head)
		}
		if pathItem.Put != nil {
			parseIgnoreOperation(ignores, pathItem.Put)
		}
		if pathItem.Trace != nil {
			parseIgnoreOperation(ignores, pathItem.Trace)
		}
	}

	return nil
}

func addIgnore(ignores map[string][]string, opName, nodeName, fieldName string) {
	operationID := fmt.Sprintf("%s%s", entoas.OpCreate, nodeName)
	if ignores[operationID] == nil {
		ignores[operationID] = make([]string, 0)
	}

	ignores[operationID] = append(ignores[operationID], fieldName)
}

func parseIgnoreOperation(ignores map[string][]string, op *ogen.Operation) {
	if val, ok := ignores[op.OperationID]; ok {
		for k, v := range op.RequestBody.Content {
			newProps := make([]ogen.Property, 0)
			newRequired := make([]string, 0)
			for _, prop := range v.Schema.Properties {
				found := false
				for _, field := range val {
					if field == prop.Name {
						found = true
						break
					}
				}
				if !found {
					newProps = append(newProps, prop)
					if slices.Index(op.RequestBody.Content[k].Schema.Required, prop.Name) > -1 {
						newRequired = append(newRequired, prop.Name)
					}
				}
			}

			op.RequestBody.Content[k].Schema.Properties = newProps
			op.RequestBody.Content[k].Schema.Required = newRequired
		}
	}
}
