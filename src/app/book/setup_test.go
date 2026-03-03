package book_test

import (
	"context"
	"testing"

	"github.com/ariesmaulana/api-kit/src/app/book"
	testsuite "github.com/ariesmaulana/api-kit/testing"
)

// BookApp holds the initialized book application components
type BookApp struct {
	*testsuite.AppContext
	Helper  *TestHelper
	Storage book.Storage
	Service book.Service
}

// TestSuite wraps testsuite.Suite for book tests
type TestSuite struct {
	*testsuite.Suite
}

// Run executes a test scenario with initialized BookApp and context
func (ts *TestSuite) Run(t *testing.T, scenario string, fn func(t *testing.T, ctx context.Context, app *BookApp)) {
	ts.Runs(t, scenario, func(t *testing.T, appCtx *testsuite.AppContext) {
		ctx := context.Background()
		app := initBookApp(appCtx)
		fn(t, ctx, app)
	})
}

// Setup registers a function to run before each test scenario with initialized BookApp and context
func (ts *TestSuite) Setup(fn func(ctx context.Context, app *BookApp)) {
	ts.Before(func(appCtx *testsuite.AppContext) {
		ctx := context.Background()
		app := initBookApp(appCtx)
		fn(ctx, app)
	})
}

// initBookApp initializes book app components from the app context
func initBookApp(app *testsuite.AppContext) *BookApp {
	helper := NewTestHelper(app.Pool)
	storage := book.NewStorage(app.Pool)
	service := book.NewService(storage)

	return &BookApp{
		AppContext: app,
		Helper:     helper,
		Storage:    storage,
		Service:    service,
	}
}

// RunTest is a wrapper that automatically sets up and tears down the test suite
func RunTest(t *testing.T, testFunc func(t *testing.T, suite *TestSuite)) {
	t.Parallel()
	cfg := testsuite.InitTestConfig()

	baseSuite, err := testsuite.NewSuite(cfg, "books.sql")
	if err != nil {
		t.Fatalf("Failed to create test suite: %v", err)
	}

	t.Cleanup(func() {
		baseSuite.Close()
	})

	suite := &TestSuite{Suite: baseSuite}
	testFunc(t, suite)
}
