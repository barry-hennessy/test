package sweet_test

import (
	"testing"

	suite "github.com/barry-hennessy/test/sweet"
)

func TestDogfoodRunF(t *testing.T) {
	type depsRunF struct {
		depFactory, after bool
	}

	depFactory := func(t *testing.T) *depsRunF {
		return &depsRunF{true, false}
	}

	type dogfood struct {
		innerCalled bool
	}

	dogfoodFactory := func(t *testing.T) *dogfood {
		d := &dogfood{}
		t.Cleanup(func() {
			// It's not the best idea to put your tests in the cleanup function;
			// better to make your tests obvious. This just shows that it's
			// possible, and that's the point of _these_ tests.
			if !d.innerCalled {
				t.Error("the inner test function was never called")
				t.Fail()
			}
		})
		return d
	}

	suite.Run(t, "DepFactory gets called", dogfoodFactory, func(t *testing.T, d *dogfood) {
		suite.Run(t, "everything gets called", depFactory, func(t *testing.T, deps *depsRunF) {
			d.innerCalled = true
		})
	})

	suite.Run(t, "modifications to the dependency are visible within the test, if passed by reference", dogfoodFactory, func(t *testing.T, d *dogfood) {
		suite.Run(t, "everything gets called", depFactory, func(t *testing.T, deps *depsRunF) {
			d.innerCalled = true
		})

		if !d.innerCalled {
			t.Error("the inner test function was never called")
			t.Fail()
		}
	})

	suite.Run(t, "DepFactory gets called - pass by value", dogfoodFactory, func(t *testing.T, d *dogfood) {
		depFactoryByValue := func(t *testing.T) depsRunF {
			return *depFactory(t)
		}

		suite.Run(t, "everything gets called", depFactoryByValue, func(t *testing.T, deps depsRunF) {
			d.innerCalled = true
		})
	})

	suite.Run(t, "succeeds if nil deps created", dogfoodFactory, func(t *testing.T, d *dogfood) {
		nilDepFactory := func(t *testing.T) *depsRunF {
			return nil
		}

		suite.Run(t, "everything gets called", nilDepFactory, func(t *testing.T, deps *depsRunF) {
			d.innerCalled = true
		})
	})

	suite.Run(t, "DepFactory is optional", dogfoodFactory, func(t *testing.T, d *dogfood) {
		suite.Run(t, "everything gets called", nil, func(t *testing.T, deps *depsRunF) {
			d.innerCalled = true
		})
	})
}
