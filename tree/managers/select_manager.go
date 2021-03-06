package managers

import (
  "github.com/chuckpreslar/codex/tree/nodes"
)

// SelectManager manages a tree that compiles to a SQL select statement.
type SelectManager struct {
  Tree    *nodes.SelectStatementNode // The AST for the SQL SELECT statement.
  Context *nodes.SelectCoreNode      // Reference to the Core the manager is curretly operating on.
  engine  interface{}                // The SQL Engine.
}

// Appends a projection to the current Context's Projections slice,
// typically an AttributeNode or string.  If a string is provided, it is
// inserted as a LiteralNode.
func (self *SelectManager) Project(projections ...interface{}) *SelectManager {
  for _, projection := range projections {
    if _, ok := projection.(string); ok {
      projection = nodes.Literal(projection)
    }
    self.Context.Projections = append(self.Context.Projections, projection)
  }
  return self
}

// Appends an expression to the current Context's Wheres slice,
// typically a comparison, i.e. 1 = 1
func (self *SelectManager) Where(expr interface{}) *SelectManager {
  if _, ok := expr.(string); ok {
    expr = nodes.Literal(expr)
  }

  self.Context.Wheres = append(self.Context.Wheres, nodes.Grouping(expr))
  return self
}

// Sets the Tree's Offset to the given integer.
func (self *SelectManager) Offset(skip int) *SelectManager {
  self.Tree.Offset = nodes.Offset(skip)
  return self
}

// Sets the Tree's Limit to the given integer.
func (self *SelectManager) Limit(take int) *SelectManager {
  self.Tree.Limit = nodes.Limit(take)
  return self
}

// Appends a new InnerJoin to the current Context's SourceNode.
func (self *SelectManager) InnerJoin(table interface{}) *SelectManager {
  switch table.(type) {
  case Accessor:
    self.Context.Source.Right = append(self.Context.Source.Right, nodes.InnerJoin(table.(Accessor).Relation(), nil))
  case *nodes.RelationNode:
    self.Context.Source.Right = append(self.Context.Source.Right, nodes.InnerJoin(table.(*nodes.RelationNode), nil))
  }
  return self
}

// Appends a new InnerJoin to the current Context's SourceNode.
func (self *SelectManager) OuterJoin(table interface{}) *SelectManager {
  switch table.(type) {
  case Accessor:
    self.Context.Source.Right = append(self.Context.Source.Right, nodes.OuterJoin(table.(Accessor).Relation(), nil))
  case *nodes.RelationNode:
    self.Context.Source.Right = append(self.Context.Source.Right, nodes.OuterJoin(table.(*nodes.RelationNode), nil))
  }
  return self
}

// Sets the last stored Join's Right leaf to a OnNode containing the
// given expression.
func (self *SelectManager) On(expr interface{}) *SelectManager {
  joins := self.Context.Source.Right

  if 0 == len(joins) {
    return self
  }

  last := joins[len(joins)-1]

  switch last.(type) {
  case *nodes.InnerJoinNode:
    last.(*nodes.InnerJoinNode).Right = nodes.On(expr)
  case *nodes.OuterJoinNode:
    last.(*nodes.OuterJoinNode).Right = nodes.On(expr)
  }

  return self
}

// Appends an expression to the current Context's Orders slice,
// typically an attribute.
func (self *SelectManager) Order(expr interface{}) *SelectManager {
  if _, ok := expr.(string); ok {
    expr = nodes.Literal(expr)
  }

  self.Tree.Orders = append(self.Tree.Orders, expr)
  return self
}

// Appends a node to the current Context's Groups slice,
// typically an attribute or column.
func (self *SelectManager) Group(groupings ...interface{}) *SelectManager {
  for _, group := range groupings {
    if _, ok := group.(string); ok {
      group = nodes.Literal(group)
    }

    self.Context.Groups = append(self.Context.Groups, group)
  }
  return self
}

// Sets the Context's Having member to the given expression.
func (self *SelectManager) Having(expr interface{}) *SelectManager {
  if _, ok := expr.(string); ok {
    expr = nodes.Literal(expr)
  }

  self.Context.Having = nodes.Having(expr)
  return self
}

// Union sets the SelectManager's Tree's Combination member to a
// UnionNode of itself and the parameter `manager`'s Tree.
func (self *SelectManager) Union(manager *SelectManager) *SelectManager {
  self.Tree.Combinator = nodes.Union(self.Tree, manager.Tree)
  return self
}

// Intersect sets the SelectManager's Tree's Combination member to a
// IntersectNode of itself and the parameter `manager`'s Tree.
func (self *SelectManager) Intersect(manager *SelectManager) *SelectManager {
  self.Tree.Combinator = nodes.Intersect(self.Tree, manager.Tree)
  return self
}

// Except sets the SelectManager's Tree's Combination member to a
// ExceptNode of itself and the parameter `manager`'s Tree.
func (self *SelectManager) Except(manager *SelectManager) *SelectManager {
  self.Tree.Combinator = nodes.Except(self.Tree, manager.Tree)
  return self
}

// Sets the SQL Enginge.
func (self *SelectManager) Engine(engine interface{}) *SelectManager {
  if _, ok := VISITORS[engine]; ok {
    self.engine = engine
  }
  return self
}

// Calls a visitor's Accept method based on the manager's SQL Engine.
func (self *SelectManager) ToSql() (string, error) {
  for _, core := range self.Tree.Cores {
    if 0 == len(core.Projections) {
      core.Projections = append(core.Projections, nodes.Attribute(nodes.Star(), core.Relation))
    }
  }

  if nil == self.engine {
    self.engine = "to_sql"
  }

  return VISITORS[self.engine].Accept(self.Tree)
}

// SelectManager factory method.
func Selection(relation *nodes.RelationNode) (selection *SelectManager) {
  selection = new(SelectManager)
  selection.Tree = nodes.SelectStatement(relation)
  selection.Context = selection.Tree.Cores[0]
  return
}
