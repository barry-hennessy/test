# Sweet; a testing suite

`sweet` is a simple test suite 'framework'; if you can even call it
a framework.

Simply put `sweet` supplies the most 'test suite' functionality possible,
while staying close to standard golang testing idioms.

Sweet helps you set up and reuse test dependencies with minimal risk of mutating
shared state between runs. It stays as close to go's standard `testing` package
as possible, there's no big up front work to slow you down.
Just drop it in and get going.

## Manage state; your tests dependencies
Many of the mistakes made writing tests come from accidentally sharing state
between runs. Either by aliasing a variable in a for loop, sharing one struct
for all tests or testing against the same mock database.

Sweet minimises the risk of one test sharing state with another, because the factories
create new objects for each test. You could go out of your way to use pointers
to globals, but... don't!

Managing dependencies `Before` and `After` each test is handled by a factory
function that returns new test dependencies (`Before`) and optionally calls
`t.Cleanup` on them (`After`).

The test call itself is typed. No need to cast.

```go
flammibleFactory := func(t *testing.T) *flammable {
    f := &flammable{}
    t.Cleanup(func() {
        putOutFlames(f)
    })
    return f
}

sweet.Run(t, "it is on fire", flammibleFactory, func(t *testing.T, f *flammable) {
    f.Ignite()
    // ... 
})

sweet.Run(t, "the fire has spread", flammibleFactory, func(t *testing.T, f *flammable) {
    f.Ignite()
    // ... 
})
// ...
```

`BeforeSuite` and `AfterSuite` functionality is just another level
of `sweet.Run` calls. In fact you can nest and organise your test dependencies
as much or as little as you like.

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


    sweet.Run(t, "the alarm goes off", depFactory, func(t *testing.T, d *deps) {
        d.alarm.Trigger()
        // ...
    })

    sweet.Run(t, "the fire brigade comes", depFactory, func(t *testing.T, d *deps) {
        d.alarm.Trigger()
        d.fireBrigade.Go()
        // ...
    })
})
```

Sets of dependencies can be passed in as custom structs, slices or maps. 

The definition of factory functions can be centralised. For example common database
or other external dependency can have a set of stock test containers that
can spin themselves up for every test and clean themselves after. ALl defined 
once and reused throughout every project.

### Links
 - [Project overview and background](https://barryhennessy.com/projects/test-sweet/)
