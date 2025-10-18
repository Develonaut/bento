package bentos

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"bento/pkg/jubako"
	"bento/pkg/omise/emoji"
	"bento/pkg/omise/styles"

	"github.com/charmbracelet/bubbles/list"
)

// bentoItem represents a bento in the list
type bentoItem struct {
	name        string
	path        string
	version     string
	description string
	icon        string
	modified    time.Time
}

// Title returns the item title
func (i bentoItem) Title() string {
	// Use custom icon if set, otherwise generate deterministic emoji from name
	icon := i.icon
	if icon == "" {
		icon = emoji.GetSushi(i.name)
	}

	versionStyled := styles.Subtle.Render(fmt.Sprintf("v%s", i.version))
	return fmt.Sprintf("%s %s %s", icon, i.name, versionStyled)
}

// Description returns the item description
func (i bentoItem) Description() string {
	return i.description
}

// FilterValue returns the value to filter by
func (i bentoItem) FilterValue() string {
	return i.name
}

// loadBentos loads bentos from store
func loadBentos(store *jubako.Store) ([]list.Item, error) {
	infos, err := store.List()
	if err != nil {
		return nil, err
	}

	items := make([]list.Item, 0, len(infos))

	for _, info := range infos {
		if item, ok := loadBentoItem(store, info); ok {
			items = append(items, item)
		}
	}

	return items, nil
}

// loadBentoItem loads a single bento item from store
func loadBentoItem(store *jubako.Store, info jubako.BentoInfo) (bentoItem, bool) {
	def, err := store.Load(extractBentoName(info.Name))
	if err != nil {
		return bentoItem{}, false
	}

	// Use description if available, fallback to type
	description := def.Description
	if description == "" {
		description = def.Type
	}

	return bentoItem{
		name:        extractBentoName(info.Name),
		path:        info.Path,
		version:     def.Version,
		description: description,
		icon:        def.Icon,
		modified:    info.Modified,
	}, true
}

// generateCopyName creates a unique name for a copied bento
func generateCopyName(name string) string {
	base := strings.TrimSuffix(name, ".bento.yaml")
	return fmt.Sprintf("%s-copy", base)
}

// extractBentoName extracts the bento name from a path or filename
func extractBentoName(pathOrName string) string {
	base := filepath.Base(pathOrName)
	return strings.TrimSuffix(base, ".bento.yaml")
}
