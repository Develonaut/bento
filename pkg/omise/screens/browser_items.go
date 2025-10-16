package screens

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"bento/pkg/jubako"

	"github.com/charmbracelet/bubbles/list"
)

// bentoItem represents a bento in the list
type bentoItem struct {
	name      string
	path      string
	version   string
	nodeType  string
	modified  time.Time
	isNewItem bool
}

// Title returns the item title
func (i bentoItem) Title() string {
	if i.isNewItem {
		return "Create New Bento"
	}
	return fmt.Sprintf("%s (v%s)", i.name, i.version)
}

// Description returns the item description
func (i bentoItem) Description() string {
	if i.isNewItem {
		return "Start building a new bento from scratch"
	}
	return fmt.Sprintf("%s • Modified: %s", i.nodeType, i.modified.Format("2006-01-02 15:04"))
}

// FilterValue returns the value to filter by
func (i bentoItem) FilterValue() string {
	if i.isNewItem {
		return "new create"
	}
	return i.name
}

// loadBentos loads bentos from store
func loadBentos(store *jubako.Store) ([]list.Item, error) {
	infos, err := store.List()
	if err != nil {
		return nil, err
	}

	items := make([]list.Item, 0, len(infos)+1)
	items = append(items, createNewBentoItem())

	for _, info := range infos {
		if item, ok := loadBentoItem(store, info); ok {
			items = append(items, item)
		}
	}

	return items, nil
}

// createNewBentoItem creates the special "Create New Bento" item
func createNewBentoItem() bentoItem {
	return bentoItem{isNewItem: true}
}

// loadBentoItem loads a single bento item from store
func loadBentoItem(store *jubako.Store, info jubako.BentoInfo) (bentoItem, bool) {
	def, err := store.Load(extractBentoName(info.Name))
	if err != nil {
		return bentoItem{}, false
	}

	return bentoItem{
		name:     extractBentoName(info.Name),
		path:     info.Path,
		version:  def.Version,
		nodeType: def.Type,
		modified: info.Modified,
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
