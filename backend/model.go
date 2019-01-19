package scores

import (
	"time"
)

// Model is an entity with a primary key `ID` that gets auto
// assigned by the repository when set.
type Model struct {
	ID        int       `json:"id" db:"id"`
}
	
// SetID sets the ID on the model
func (m *Model) SetID(id int) {
	m.ID = id
}

// Tracked adds timestamps `CreatedAt`, `UpdatedAt`,
// `DeletedAt` to the model.
type Tracked struct {
	CreatedAt time.Time  `json:"createdAt" db:"createdAt"`
	UpdatedAt time.Time  `json:"-" db:"updatedAt"`
	DeletedAt *time.Time `json:"-" db:"deletedAt"`
}

// SetCreatedAt sets the `CreatedAt` field.
func (t *Tracked) SetCreatedAt(created time.Time) {
	t.CreatedAt = created
}

// SetUpdatedAt sets the `UpdatedAt` field.
func (t *Tracked) SetUpdatedAt(updated time.Time) {
	t.UpdatedAt = updated
}

// SetDeletedAt sets the `DeletedAt` field.
func (t *Tracked) SetDeletedAt(deleted *time.Time) {
	t.DeletedAt = deleted
}