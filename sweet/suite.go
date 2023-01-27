package sweet

import "testing"

// DepFactory creates all dependencies needed for a test run.
//
// It is responsible for cleaning up using `t.Cleanup`. This goes for resources
// it creates, and for state changes made by the test to its dependencies.
type DepFactory[deps any] func(t *testing.T) deps

// Run runs a subtest, just like t.Run, except it takes a DepFactory that
// generates a new set of test dependencies for each test. The test is passed
// the dependencies as it's second argument.
//
// Compared to t.Run:
//
//	t.Run("subtest name", func(t *testing.T) {...})
//	sweet.Run(t, "subtest name", func(t *testing.T) deps, func(t *testing.T, d deps) {...})
func Run[deps any, ptrDeps *deps](
	t *testing.T,
	testName string,
	factory DepFactory[deps],
	coreTest func(t *testing.T, d deps),
) bool {
	return t.Run(testName, func(t *testing.T) {
		if factory != nil {
			coreTest(t, factory(t))
		} else {
			coreTest(t, *ptrDeps(new(deps)))
		}
	})
}

// DepFactories is a map of dependency factory functions.
//
// It can be used to instantiate a set of dependencies before each test run by
// `Run`.
//
// Note: the concrete set of dependencies will have to be mapped back from
// `interface{}` to their concrete type. See `FactoryForMap`.
// To avoid this consider using a custom struct to return your dependencies.
type DepFactories[deps any] map[string]DepFactory[any]

func (d DepFactories[deps]) DepFactory() func(t *testing.T) DepsMapped[deps] {
	return func(t *testing.T) DepsMapped[deps] {
		ds := make(DepsMapped[deps], len(d))
		for i := range d {
			ds[i] = FactoryForMap(d[i])(t)
		}
		return ds
	}
}

// FactoryForMap is a helper function to convert dependency factories into a form
// compatible with `DepFactories`.
func FactoryForMap[deps any](factory DepFactory[deps]) DepFactory[interface{}] {
	return func(t *testing.T) interface{} {
		return factory(t)
	}
}

// DepsMapped is a map of a concrete set of dependencies.
//
// If returned by an instance of DepMapFactory it will be a mirror of the
// `DepFactories` you passed in, only with concrete implementations where you left
// factory functions.
type DepsMapped[deps any] map[string]deps


// DepMapFactory is a factory function that returns a map of new instances of
// its dependencies. For use by Run.
func DepMapFactory[deps any](depFactories DepFactories[any]) func(t *testing.T) DepsMapped {
	return func(t *testing.T) DepsMapped {
		ds := make(DepsMapped, len(depFactories))
		for i := range depFactories {
			ds[i] = depFactories[i](t)
		}
		return ds
	}
}
