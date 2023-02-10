package sweet_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/barry-hennessy/test/sweet"
)

func TestUseCases(t *testing.T) {
	sweet.Run(t, "set a global", nil, func(t *testing.T, d any) {
		var global bool

		before := func(t *testing.T) any {
			global = true
			t.Cleanup(func() {
				global = false
			})
			return nil
		}

		sweet.Run(t, "before testing and reset it after testing", before, func(t *testing.T, deps any) {
			if !global {
				t.Error("global was not set")
			}
		})

		if global {
			t.Error("global was not reset")
		}
	})

	sweet.Run(t, "can start a dep that", nil, func(t *testing.T, deps any) {
		sweet.Run(t, "needs to be set before and reset after each test", nil, func(t *testing.T, deps any) {
			type resettableDep struct {
				state bool
			}

			// Say this models a connection to an external data source, we define
			// it like this so we can inspect if the 'connection' is open or closed
			// after our tests
			var someGlobalConnection bool

			before := func(t *testing.T) *resettableDep {
				someGlobalConnection = true
				deps := &resettableDep{someGlobalConnection}

				t.Cleanup(func() {
					someGlobalConnection = false
					deps.state = someGlobalConnection
				})

				return deps
			}

			for i := 0; i < 3; i++ {
				sweet.Run(t, fmt.Sprintf("test %d", i), before, func(t *testing.T, deps *resettableDep) {
					if !deps.state {
						t.Error("a connection was not made before the test")
					}
				})

				if someGlobalConnection {
					t.Error("the connection was not cleaned up after the test")
				}
			}
		})

		sweet.Run(t, "needs to be set before and reset after a set of tests", nil, func(t *testing.T, deps any) {
			type dataDep struct {
				db map[string]bool
			}

			// Say this models a connection to an external data source, we define
			// it like this so we can inspect if the 'connection' is open or closed
			// after our tests
			var someGlobalDatabase map[string]bool

			before := func(t *testing.T) *dataDep {
				someGlobalDatabase := map[string]bool{}
				deps := &dataDep{someGlobalDatabase}

				t.Cleanup(func() {
					someGlobalDatabase = nil
					deps.db = someGlobalDatabase
				})

				return deps
			}

			sweet.Run(t, "outer test", before, func(t *testing.T, deps *dataDep) {
				if deps.db == nil {
					t.Error("a database was not made before the test")
				}

				for i := 0; i < 3; i++ {
					runName := fmt.Sprintf("inner test %d", i)
					// Be aware that you could easily shadow the outer `deps` here
					sweet.Run(t, runName, nil, func(t *testing.T, innerDeps *dataDep) {
						deps.db[runName] = true
					})
				}

				if len(deps.db) != 3 {
					t.Errorf("the database was not seen by all tests. Size was %d", len(deps.db))
				}
			})

			if someGlobalDatabase != nil {
				t.Error("the database was not cleaned up after the test")
			}
		})
	})

	t.Run("can run with multiple dependencies", func(t *testing.T) {
		// It's worth pointing out that these deps, and their factory functions
		// can be used by multiple different kinds of tests. I.e. instances of
		// this kind could be centralised in a repo or foreign project.
		// Think a struct like this per database/foreign system. Or even one for
		// each kind of docker based testcontainer.
		type dep struct {
			state string
			rand  float64
		}

		type otherDep struct {
			state         string
			rand          float64
			somethingElse bool
		}

		depAFactory := func(t *testing.T) *dep {
			d := &dep{"depA", rand.Float64()}
			t.Cleanup(func() {
				d.state = d.state + " cleaned up"
			})
			return d
		}

		depBFactory := func(t *testing.T) *otherDep {
			d := &otherDep{"depB", rand.Float64(), true}
			t.Cleanup(func() {
				d.state = d.state + " cleaned up"
			})
			return d
		}

		t.Run("struct of dependencies", func(t *testing.T) {
			type depStruct struct {
				A *dep
				B *otherDep
			}

			depFactory := func(t *testing.T) depStruct {
				return depStruct{
					A: depAFactory(t),
					B: depBFactory(t),
				}
			}

			sweet.Run(t, "instantiate correctly", depFactory, func(t *testing.T, deps depStruct) {
				if deps.A.state != "depA" {
					t.Error("dep map was not instantiated correctly")
				}

				if deps.B.state != "depB" {
					t.Error("dep map was not instantiated correctly")
				}
			})

			t.Run("instiantiate new deps every time", func(t *testing.T) {
				deps := []depStruct{}
				for i := 0; i < 2; i++ {
					sweet.Run(t, fmt.Sprintf("run %d", i), depFactory, func(t *testing.T, d depStruct) {
						deps = append(deps, d)
					})
				}

				if deps[0].A.rand == deps[1].A.rand {
					t.Error("deps were not instantiated correctly")
				}

				if deps[0].B.rand == deps[1].B.rand {
					t.Error("deps were not instantiated correctly")
				}
			})
		})
	})
}
