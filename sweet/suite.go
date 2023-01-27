package sweet

import "testing"

// DepFactory creates all dependencies needed for a test run.
//
// It is responsible for cleaning up using `t.Cleanup`. This goes for resources
// it creates, and for state changes made by the test to its dependencies.
type DepFactory[deps any, testingType testing.TB] func(t testingType) deps

type testRun[testRunner any] interface{
	Run(name string, f func(t testRunner)) bool
}

// Run runs a subtest, just like t.Run, except it takes a DepFactory that
// generates a new set of test dependencies for each test. The test is passed
// the dependencies as it's second argument.
//
// Compared to t.Run:
//
//	t.Run("subtest name", func(t *testing.T) {...})
//	sweet.Run(t, "subtest name", func(t *testing.T) deps, func(t *testing.T, d deps) {...})
//
// Run is also compatible with benchmarks; just pass in a *testing.B and a
// *testing.B will be passed to your factory and test functions.
func Run[deps any, ptrDeps *deps, testingType testing.TB, testRunner testRun[testingType]](
	t testRunner,
	testName string,
	factory DepFactory[deps, testingType],
	coreTest func(t testingType, d deps),
) bool {
	return t.Run(testName, func(t testingType) {
		if factory != nil {
			coreTest(t, factory(t))
		} else {
			coreTest(t, *ptrDeps(new(deps)))
		}
	})
}
