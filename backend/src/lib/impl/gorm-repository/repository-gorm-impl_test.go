package gorm_impl

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/cmo7/folly4/src/lib/generics/common"
	"github.com/cmo7/folly4/src/lib/generics/filter"
	"github.com/cmo7/folly4/src/lib/generics/order"
	"github.com/cmo7/folly4/src/lib/generics/pagination"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ParentEntity struct {
	ID        common.ID `gorm:"type:char(36);primary_key"`
	Name      string
	BirthDate time.Time     `gorm:"type:date"`
	Children  []ChildEntity `gorm:"foreignKey:ParentID"`
}

type ChildEntity struct {
	ID       common.ID `gorm:"type:char(36);primary_key"`
	Name     string
	ParentID common.ID      `gorm:"type:char(36)"`
	Siblings []*ChildEntity `gorm:"many2many:child_entity_siblings"`
}

func (e *ParentEntity) BeforeCreate(tx *gorm.DB) error {
	e.ID = common.NewID()
	return nil
}

func (e *ChildEntity) BeforeCreate(tx *gorm.DB) error {
	e.ID = common.NewID()
	return nil
}

func (e ParentEntity) GetID() common.ID {
	return e.ID
}

func (e ParentEntity) GetName() string {
	return e.Name
}

func (e ChildEntity) GetID() common.ID {
	return e.ID
}

func (e ChildEntity) GetName() string {
	return e.Name
}

var parentRepository *GormGenericRepository[ParentEntity]
var childRepository *GormGenericRepository[ChildEntity]

func setupDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&ParentEntity{}, &ChildEntity{})
	if err != nil {
		panic("failed to migrate database schema")
	}

	return db
}

func truncateDB(db *gorm.DB) {
	db.Migrator().DropTable(&ChildEntity{}, &ParentEntity{})
	db.AutoMigrate(&ParentEntity{}, &ChildEntity{})
}

func setupTest(t *testing.T) {
	truncateDB(parentRepository.db)
	t.Cleanup(func() {
		truncateDB(parentRepository.db)
	})
}

func createContext(timeout time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	return ctx
}

func TestMain(m *testing.M) {
	db := setupDB()
	parentRepository = NewGormGenericRepository[ParentEntity](db)
	childRepository = NewGormGenericRepository[ChildEntity](db)

	code := m.Run()
	os.Exit(code)
}

func TestGormGenericRepositoryCreateParent(t *testing.T) {
	setupTest(t)

	ctx := createContext(5 * time.Second)
	parent := ParentEntity{Name: "Parent"}
	result, err := parentRepository.Create(ctx, parent)

	assert.Nil(t, err)
	assert.Equal(t, "Parent", result.Name)
}

func TestGormGenericRepositoryCreateChild(t *testing.T) {
	setupTest(t)

	ctx := createContext(5 * time.Second)

	parent := ParentEntity{Name: "Parent"}
	createdParent, err := parentRepository.Create(ctx, parent)
	assert.Nil(t, err)

	child := ChildEntity{Name: "Child", ParentID: createdParent.ID}
	result, err := childRepository.Create(ctx, child)

	assert.Nil(t, err)
	assert.Equal(t, "Child", result.Name)
	assert.Equal(t, createdParent.ID, result.ParentID)
}

func TestGormGenericRepositoryFindParentWithChildren(t *testing.T) {
	setupTest(t)

	ctx := createContext(5 * time.Second)

	parent := ParentEntity{Name: "Parent"}
	child := ChildEntity{Name: "Child"}

	createdParent, err := parentRepository.Create(ctx, parent)
	assert.Nil(t, err)

	child.ParentID = createdParent.ID
	_, err = childRepository.Create(ctx, child)
	assert.Nil(t, err)

	var loadedParent ParentEntity
	err = parentRepository.db.WithContext(ctx).Preload("Children").First(&loadedParent, createdParent.ID).Error

	assert.Nil(t, err)
	assert.Equal(t, 1, len(loadedParent.Children))
	assert.Equal(t, "Child", loadedParent.Children[0].Name)
}

func TestGormGenericRepositoryPagination(t *testing.T) {
	setupTest(t)

	ctx := createContext(5 * time.Second)

	parents := []ParentEntity{
		{Name: "Parent1"},
		{Name: "Parent2"},
		{Name: "Parent3"},
	}

	for _, p := range parents {
		_, err := parentRepository.Create(ctx, p)
		assert.Nil(t, err)
	}

	pageable := pagination.Pageable{Page: 0, Size: 2}
	page, err := parentRepository.FindAll(ctx, pageable, filter.Composite{Operator: filter.And, Filters: []filter.Filter{}}, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(page.Content))
}

func TestGormGenericRepositorySortByNameCases(t *testing.T) {
	setupTest(t)

	ctx := createContext(5 * time.Second)

	parents := []ParentEntity{
		{Name: "Charlie"},
		{Name: "Alice"},
		{Name: "Bob"},
		{Name: "Aaron"},
		{Name: "Amanda"},
	}

	for _, p := range parents {
		_, err := parentRepository.Create(ctx, p)
		assert.Nil(t, err)
	}

	testCases := []struct {
		name           string
		orderBy        []order.OrderBy
		expectedOrder  []string
		expectedLength int
	}{
		{"Sort by Name Ascending", []order.OrderBy{{Field: "Name", Direction: order.Asc}}, []string{"Aaron", "Alice", "Amanda", "Bob", "Charlie"}, 5},
		{"Sort by Name Descending", []order.OrderBy{{Field: "Name", Direction: order.Desc}}, []string{"Charlie", "Bob", "Amanda", "Alice", "Aaron"}, 5},
		{"Paginated Sort Ascending", []order.OrderBy{{Field: "Name", Direction: order.Asc}}, []string{"Aaron", "Alice"}, 2},
		{"Paginated Sort Descending", []order.OrderBy{{Field: "Name", Direction: order.Desc}}, []string{"Charlie", "Bob"}, 2},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var pageable pagination.Pageable
			if tc.name == "Paginated Sort Ascending" || tc.name == "Paginated Sort Descending" {
				pageable = pagination.Pageable{Page: 0, Size: 2}
			} else {
				pageable = pagination.Pageable{Page: 0, Size: 10}
			}

			page, err := parentRepository.FindAll(ctx, pageable, filter.Composite{Operator: filter.And, Filters: []filter.Filter{}}, nil, tc.orderBy)

			assert.Nil(t, err)
			assert.Equal(t, tc.expectedLength, len(page.Content))

			actualOrder := make([]string, len(page.Content))
			for i, p := range page.Content {
				actualOrder[i] = p.Name
			}

			assert.Equal(t, tc.expectedOrder, actualOrder)
		})
	}
}

func TestGormGenericRepositoryTimeout(t *testing.T) {
	setupTest(t)

	ctx := createContext(0 * time.Millisecond) // Timeout immediately
	parent := ParentEntity{Name: "Parent"}
	_, err := parentRepository.Create(ctx, parent)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}
