package sweet_test

import (
	"testing"

	"github.com/barry-hennessy/test/sweet"
)

type (
	engine interface {
		Rev()
	}

	truck interface {
		Vroom()
	}

	hose interface {
		IsReeledUp() bool
	}

	fireTruck struct {
		hose   hose
		engine engine
	}
	mockEngine struct{}
	mockHose   struct{}
)

type flammable struct {
	onFire bool
}

func flammableFactory(t *testing.T) *flammable {
	f := &flammable{}

	// Make sure to put the fire out when we're done
	t.Cleanup(func() {
		f.Extinguish()
	})

	return f
}

func (f *flammable) Ignite() {
	f.onFire = true
}

func (f *flammable) Extinguish() {
	f.onFire = false
}

func (_ mockEngine) Rev() {}

func (t fireTruck) Vroom() {
	t.engine.Rev()
}

func (h mockHose) IsReeledUp() bool {
	return true
}

// fireTruckDeps houses everything your tests need to test
// how a fire truck behaves.
type fireTruckDeps struct {
	truck  truck
	hose   hose
	engine engine
}

// A [DepFactory] creates and initialises all the structs and values your
// tests depend on. The software under test itself, but also mocks, databases
// that need spinning up etc. It takes a [testing.T] so it can `Cleanup` once
// the test is over.
//
// [sweet.Run] takes a [DepFactory] and runs a test, just like [testing.T.Run]
// but your test function gets fresh data to test on every time.
func Example() {
	t := &testing.T{}

	// Your dependency factory creates a fresh set of what your tests need and
	// cleanup once the test is done.
	flammableFactory := func(t *testing.T) *flammable {
		f := &flammable{}

		// Make sure to put the fire out when we're done
		t.Cleanup(func() {
			f.Extinguish()
		})

		return f
	}

	// `sweet.Run` creates a new instance and passes it to your test function.
	sweet.Run(t, "it is on fire", flammableFactory, func(t *testing.T, f *flammable) {
		f.Ignite()
		// ...
	})

	// A fresh instance every time.
	sweet.Run(t, "fire spreads", flammableFactory, func(t *testing.T, f *flammable) {
		f.Ignite()
		// ...
	})
}

// This is the most straightforward way to use `sweet`.
//
// If your software under test is wrapped up in one factory function just
// pass it to [sweet.Run] and use your fresh value in every test.
func ExampleRun_direct() {
	t := &testing.T{}

	// flammableFactory returns exactly what we're testing
	sweet.Run(t, "flames are hot", flammableFactory, func(t *testing.T, f *flammable) {
	})

	sweet.Run(t, "flames are orange", flammableFactory, func(t *testing.T, f *flammable) {
	})
}

// A straightforward option for managing multiple dependencies is to create a struct
// to house your related dependencies.
//
// It's more up front work, but everything is typed and it can pay off if you're
// using the struct often.
func ExampleRun_structOfDependencies() {
	t := &testing.T{}

	fireTruckFactory := func(t *testing.T) fireTruckDeps {
		mockHose := mockHose{}
		mockEngine := mockEngine{}

		return fireTruckDeps{
			truck:  fireTruck{mockHose, mockEngine},
			hose:   mockHose,
			engine: mockEngine,
		}
	}

	sweet.Run(t, "when it is on fire", fireTruckFactory, func(t *testing.T, d fireTruckDeps) {
		d.truck.Vroom()

		if !d.hose.IsReeledUp() {
			t.Error("you can't be driving around with a dangling hose")
		}
	})
}

// If your test needs multiple dependencies and you want to avoid the boilerplate
// of creating a struct you can just use a map or slice.
//
// It reduces the up front boilerplate but you need to cast where you use the
// values. Which can be a fine trade off if you only use the dependencies in
// one or two tests.
func ExampleRun_mapOfDependencies() {
	t := &testing.T{}

	fireTruckFactory := func(t *testing.T) map[string]any {
		mockHose := mockHose{}
		mockEngine := mockEngine{}

		return map[string]any{
			"truck":  fireTruck{mockHose, mockEngine},
			"hose":   mockHose,
			"engine": mockEngine,
		}
	}

	sweet.Run(t, "when it is on fire", fireTruckFactory, func(t *testing.T, d map[string]any) {
		d["truck"].(truck).Vroom()

		if !d["hose"].(hose).IsReeledUp() {
			t.Error("you can't be driving around with a dangling hose")
		}
	})
}

// A balanced approach when your test needs multiple dependencies is to use a
// function that returns multiple values. With this you can avoid [sweet.Run]
// altogether.
//
// So why not use the functional approach every time?
//
// You absolutely can.
//
// Remember, the main goal of sweet is to seperate out your test dependencies
// and provide (and share) reusable building blocks.
//
// If your functional factories:
//   - create clean test dependencies
//   - clean up after themselves
//   - can be composed with others
//
// Then go for it!
func ExampleRun_functional() {
	t := &testing.T{}

	fireTruckFactory := func(t *testing.T) (truck, hose, engine) {
		mockHose := mockHose{}
		mockEngine := mockEngine{}

		return fireTruck{mockHose, mockEngine}, mockHose, mockEngine
	}

	t.Run("when it is on fire", func(t *testing.T) {
		truck, hose, _ := fireTruckFactory(t)
		truck.Vroom()

		if !hose.IsReeledUp() {
			t.Error("you can't be driving around with a dangling hose")
		}
	})
}

// Pitfall: Nesting [sweet.Run] calls
//
// If you're used to using a different test suite you might be looking for
// `BeforeSuite` or `AfterSuite` functions that set up some state for all your
// tests and clean up at the end.
//
// You can easily do this with sweet, but it's worth pointing out that this
// undermines the _fresh dependencies_ that sweet tries to provide.
//
// `BeforeSuite` and `AfterSuite` functionality is just another level
// of `sweet.Run` calls. In fact you can nest and organise your test dependencies
// as much or as little as you like.
func ExampleRun_pitfall_nesting() {
	t := &testing.T{}

	flammableFactory := func(t *testing.T) *flammable {
		return &flammable{}
	}

	sweet.Run(t, "when it is on fire", flammableFactory, func(t *testing.T, f *flammable) {
		type fireTruckDeps struct {
			truck  truck
			hose   hose
			engine engine
		}

		fireTruckFactory := func(t *testing.T) fireTruckDeps {
			mockHose := mockHose{}
			mockEngine := mockEngine{}

			return fireTruckDeps{
				truck:  fireTruck{mockHose, mockEngine},
				hose:   mockHose,
				engine: mockEngine,
			}
		}

		// This is set for all tests within this block.
		// If any test calls `f.Extinguish` you have a race condition and flaky tests
		f.Ignite()

		sweet.Run(t, "the alarm goes off", fireTruckFactory, func(t *testing.T, d fireTruckDeps) {
			d.truck.Vroom()
			// ...
		})

		sweet.Run(t, "the fire brigade comes", fireTruckFactory, func(t *testing.T, d fireTruckDeps) {
			d.truck.Vroom()
			// ...
		})
	})
}

// An alternative to `BeforeSuite`/`AfterSuite` that avoids accidental sharing
// of upper level dependencies.
//
// Instead of nesting your `sweet.Run` calls, nesting your dependencies can achieve
// the same effect; just with a fresh top level dependency.
func ExampleRun_pitfall_nesting_alternative() {
	t := &testing.T{}

	fireTruckFactory := func(t *testing.T) fireTruckDeps {
		mockHose := mockHose{}
		mockEngine := mockEngine{}

		flammable := flammableFactory(t)
		flammable.Ignite()

		return fireTruckDeps{
			truck:  fireTruck{mockHose, mockEngine},
			hose:   mockHose,
			engine: mockEngine,
		}
	}

	t.Run("when it is on fire", func(t *testing.T) {
		sweet.Run(t, "the alarm goes off", fireTruckFactory, func(t *testing.T, d fireTruckDeps) {
			d.truck.Vroom()
			// ...
		})

		sweet.Run(t, "the fire brigade comes", fireTruckFactory, func(t *testing.T, d fireTruckDeps) {
			d.truck.Vroom()
			// ...
		})
	})
}
