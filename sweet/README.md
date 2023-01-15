# Sweet; a testing suite

`sweet` is a simple test suite 'framework'; if you can even call it
a framework.

Simply put `sweet` supplies the most 'test suite' functionality possible,
while staying close to standard golang testing idioms.

`Before`, `After` etc. are handled by a factory function that returns
new test dependencies (`Before`) and optionally calls `t.Cleanup` on
them (`After`).

```go
flammibleFactory := func(t *testing.T) *flammable {
    f := &flammable{}
    t.Cleanup(func() {
        putOutFlames(f)
    })
    return f
}

sweet.Run(t, "it is not on fire", flammibleFactory, func(t *testing.T, f *flammable) {
   // ... 
})
```

`BeforeSuite`/`AfterSuite` kind of functionality is just another level
of `sweet.Run` calls.

```go
sweet.Run(t, "when it is on fire", flammibleFactory, func(t *testing.T, f *flammable) {
    depFactory := func(t *testing.T) *flammable {
        f := &deps{
            alarm: alarm{},
            fireBrigade: fireBrigade{},
        }

        t.Cleanup(func() {
            f.alarm.reset()
            f.fireBrigade.reset()
        })

        return f
    }

    setFire(f)

    sweet.Run(t, "the alarm goes off", depFactory, func(t *testing.T, d *deps) {
       // ... 
    })

    sweet.Run(t, "the fire brigade comes", depFactory, func(t *testing.T, d *deps) {
       // ... 
    })
})
```

Sets of dependencies can be passed in as custom structs.

The definition of factory functions can be centralised. The goal is that
all database/other external deps having a set of test containers that
can spin themselves up, and clean up, for every test. Defined once and
reused throughout every project.

If you want to take some action without returning a dep you're free to.
E.g: outputting some text, asserting some value (as usual with `testing.T`),
waiting for a dependency to spin up and be ready.
