package plugin

import (
	"fmt"
	"sort"
)

// Chain manages plugin execution
type Chain struct {
	plugins []Plugin
}

// NewChain creates a new plugin chain
func NewChain() *Chain {
	return &Chain{
		plugins: make([]Plugin, 0),
	}
}

// Add adds a plugin to the chain
func (c *Chain) Add(plugin Plugin) {
	c.plugins = append(c.plugins, plugin)
	// Sort by priority (lower number = higher priority)
	sort.Slice(c.plugins, func(i, j int) bool {
		return c.plugins[i].Priority() < c.plugins[j].Priority()
	})
}

// ExecuteRequest runs all plugins' OnRequest methods
func (c *Chain) ExecuteRequest(ctx *Context) error {
	for _, plugin := range c.plugins {
		if ctx.Skip {
			break
		}

		if err := plugin.OnRequest(ctx); err != nil {
			return fmt.Errorf("plugin %s: %w", plugin.Name(), err)
		}

		if ctx.Error != nil {
			return ctx.Error
		}
	}
	return nil
}

// ExecuteResponse runs all plugins' OnResponse methods
func (c *Chain) ExecuteResponse(ctx *Context) error {
	// Execute in reverse order for response phase
	for i := len(c.plugins) - 1; i >= 0; i-- {
		if ctx.Skip {
			break
		}

		plugin := c.plugins[i]
		if err := plugin.OnResponse(ctx); err != nil {
			return fmt.Errorf("plugin %s: %w", plugin.Name(), err)
		}

		if ctx.Error != nil {
			return ctx.Error
		}
	}
	return nil
}

// List returns all plugins in the chain
func (c *Chain) List() []Plugin {
	return c.plugins
}
