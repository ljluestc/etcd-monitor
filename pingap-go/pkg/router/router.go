package router

import (
	"net/http"
	"regexp"
	"sort"
	"strings"

	"github.com/etcd-monitor/pingap-go/pkg/config"
)

// Router handles request routing
type Router struct {
	locations []RouteEntry
}

// RouteEntry represents a single route
type RouteEntry struct {
	Location     *config.ParsedLocation
	PathMatcher  Matcher
	HostMatcher  Matcher
}

// Matcher interface for matching rules
type Matcher interface {
	Match(value string) bool
}

// ExactMatcher matches exact strings
type ExactMatcher struct {
	value string
}

func (m *ExactMatcher) Match(value string) bool {
	return m.value == value
}

// PrefixMatcher matches string prefixes
type PrefixMatcher struct {
	prefix string
}

func (m *PrefixMatcher) Match(value string) bool {
	return strings.HasPrefix(value, m.prefix)
}

// RegexMatcher matches regular expressions
type RegexMatcher struct {
	regex *regexp.Regexp
}

func (m *RegexMatcher) Match(value string) bool {
	return m.regex.MatchString(value)
}

// NewRouter creates a new router
func NewRouter(locations []config.ParsedLocation) *Router {
	entries := make([]RouteEntry, 0, len(locations))

	for _, loc := range locations {
		entry := RouteEntry{
			Location: &loc,
		}

		// Create path matcher
		pathType := loc.PathType
		if pathType == "" {
			pathType = "prefix"
		}

		switch pathType {
		case "exact":
			entry.PathMatcher = &ExactMatcher{value: loc.Path}
		case "regex":
			if loc.PathPattern != nil {
				entry.PathMatcher = &RegexMatcher{regex: loc.PathPattern.(*regexp.Regexp)}
			}
		default: // prefix
			entry.PathMatcher = &PrefixMatcher{prefix: loc.Path}
		}

		// Create host matcher
		if loc.Host != "" {
			entry.HostMatcher = &ExactMatcher{value: loc.Host}
		}

		entries = append(entries, entry)
	}

	// Sort by priority (higher number = higher priority)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Location.Priority > entries[j].Location.Priority
	})

	return &Router{locations: entries}
}

// Match finds a matching location for the request
func (r *Router) Match(req *http.Request) (*config.ParsedLocation, bool) {
	for _, entry := range r.locations {
		// Check host if specified
		if entry.HostMatcher != nil {
			if !entry.HostMatcher.Match(req.Host) {
				continue
			}
		}

		// Check path
		if entry.PathMatcher != nil {
			if entry.PathMatcher.Match(req.URL.Path) {
				// Check methods if specified
				if len(entry.Location.Methods) > 0 {
					methodMatch := false
					for _, method := range entry.Location.Methods {
						if method == req.Method {
							methodMatch = true
							break
						}
					}
					if !methodMatch {
						continue
					}
				}

				return entry.Location, true
			}
		}
	}

	return nil, false
}
