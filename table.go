package codex

import (
  "github.com/chuckpreslar/codex/tree/managers"
  "github.com/chuckpreslar/codex/tree/nodes"
)

// Table returns an Accessor
func Table(name string) managers.Accessor {
  relation := nodes.Relation(name)
  return func(name interface{}) *nodes.AttributeNode {
    if _, ok := name.(string); ok {
      return nodes.Attribute(nodes.Column(name), relation)
    }

    return nodes.Attribute(name, relation)
  }
}
