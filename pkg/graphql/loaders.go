package graphql

import (
	_ "github.com/graph-gophers/dataloader/v7" // N+1 batching; see Loaders doc and WithLoaders / LoadersFrom
	"github.com/zatrano/zatrano/pkg/core"
)

// Loaders holds per-request DataLoader instances (create a new Loaders each HTTP request).
//
// Example with graph-gophers/dataloader/v7:
//
//	type Loaders struct {
//		UserByID *dataloader.Loader[uint, *models.User]
//	}
//
//	func NewLoaders(app *core.App) *Loaders {
//		return &Loaders{
//			UserByID: dataloader.NewBatchedLoader(batchUsersByID(app.DB),
//				dataloader.WithWait[uint, *models.User](2*time.Millisecond)),
//		}
//	}
type Loaders struct{}

// NewLoaders builds fresh loaders for one GraphQL request.
func NewLoaders(app *core.App) *Loaders {
	_ = app
	return &Loaders{}
}
