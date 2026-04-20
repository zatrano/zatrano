package testing

import (
	"database/sql"
	"testing"

	"gorm.io/gorm"
)

// TestSuite provides database transaction rollback for tests.
type TestSuite struct {
	DB       *gorm.DB
	sqlDB    *sql.DB
	original *gorm.DB
	t        *testing.T
}

// NewTestSuite creates a new test suite with transaction rollback.
func NewTestSuite(db *gorm.DB) *TestSuite {
	return &TestSuite{
		original: db,
	}
}

// SetupTest begins a new database transaction for the test.
func (s *TestSuite) SetupTest(t *testing.T) {
	s.t = t
	s.sqlDB, _ = s.original.DB()
	s.DB = s.original.Begin()
}

// TeardownTest rolls back the transaction and cleans up.
func (s *TestSuite) TeardownTest() {
	if s.DB != nil {
		s.DB.Rollback()
	}
	s.DB = nil
	s.t = nil
}

// GetDB returns the transactional database connection.
func (s *TestSuite) GetDB() *gorm.DB {
	return s.DB
}

// AssertNoError fails the test if err is not nil.
func (s *TestSuite) AssertNoError(err error) {
	if err != nil {
		s.t.Fatalf("unexpected error: %v", err)
	}
}

// AssertError fails the test if err is nil.
func (s *TestSuite) AssertError(err error) {
	if err == nil {
		s.t.Fatal("expected error but got none")
	}
}

// AssertTrue fails the test if condition is false.
func (s *TestSuite) AssertTrue(condition bool, msg string) {
	if !condition {
		s.t.Fatal(msg)
	}
}

// AssertEqual fails the test if a != b.
func (s *TestSuite) AssertEqual(a, b interface{}) {
	if a != b {
		s.t.Fatalf("expected %v, got %v", a, b)
	}
}
