package inline

import "sync"

// DryFruit 干果接口.
type DryFruit interface {
	Name() string
	Price() uint64
	Family() string
	Distribution() string
	Increase()
}

// Chestnut 定义一个栗子🌰对象，它实现了干果接口.
type Chestnut struct {
	name  string
	count uint64
}

var _ DryFruit = (*Chestnut)(nil)

// Name 名称.
func (c Chestnut) Name() string {
	return c.name
}

// Price 价格.
func (c Chestnut) Price() uint64 {
	return 10
}

// Family Family name.
func (c Chestnut) Family() string {
	return "Fagaceae"
}

// Distribution 分布.
func (c Chestnut) Distribution() string {
	return "East Asia"
}

// Increase 数量加一
func (c *Chestnut) Increase() {
	c.count++
}

// OriginGift 未优化之前的礼物对象.
type OriginGift struct {
	mu       sync.Mutex
	dryFruit DryFruit
}

// Access 访问接口对象.
func (g *OriginGift) Access() {
	g.dryFruit.Name()
	g.dryFruit.Price()
	g.dryFruit.Family()
	g.dryFruit.Distribution()
	g.dryFruit.Increase()
}

// ImprovedGift 优化后的礼物对象.
type ImprovedGift struct {
	mu       sync.Mutex
	dryFruit *Chestnut
}

// Access 访问具体的字段对象.
func (g *ImprovedGift) Access() {
	g.dryFruit.Name()
	g.dryFruit.Price()
	g.dryFruit.Family()
	g.dryFruit.Distribution()
	g.dryFruit.Increase()
}
