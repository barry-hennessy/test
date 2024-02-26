package sweet_test

import (
	"testing"

	suite "github.com/barry-hennessy/test/sweet"
)

type depsF struct {
	b bool
}

func TestRun(t *testing.T) {
	t.Run("DepFactory gets called", func(t *testing.T) {
		innerCalled := false

		depFactory := func(t *testing.T) *depsF {
			return &depsF{
				b: true,
			}
		}

		suite.Run(t, "everything gets called", depFactory, func(t *testing.T, deps *depsF) {
			innerCalled = true
		})

		if !innerCalled {
			t.Error("the inner test function was never called")
			t.Fail()
		}
	})

	t.Run("DepFactory gets called - pass by value", func(t *testing.T) {
		innerCalled := false

		depFactory := func(t *testing.T) depsF {
			return depsF{
				b: true,
			}
		}

		suite.Run(t, "everything gets called", depFactory, func(t *testing.T, deps depsF) {
			innerCalled = true
		})

		if !innerCalled {
			t.Error("the inner test function was never called")
			t.Fail()
		}
	})

	t.Run("succeeds if nil deps created", func(t *testing.T) {
		innerCalled := false

		depFactory := func(t *testing.T) *depsF {
			return nil
		}

		suite.Run(t, "everything gets called", depFactory, func(t *testing.T, deps *depsF) {
			innerCalled = true
		})

		if !innerCalled {
			t.Error("the inner test function was never called")
			t.Fail()
		}
	})

	t.Run("DepsFactory is optional", func(t *testing.T) {
		innerCalled := false

		suite.Run(t, "everything gets called", nil, func(t *testing.T, deps depsF) {
			innerCalled = true
		})

		if !innerCalled {
			t.Error("the inner test function was never called")
			t.Fail()
		}
	})
}
