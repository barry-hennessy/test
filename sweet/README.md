# Sweet; a simple test suite

`sweet` is a simple test suite 'framework'; if you can even call it
a framework.

`sweet` supplies the most 'test suite' functionality possible,
while staying close to standard golang testing idioms.

The focus is on setting up, sharing and reusing test dependencies.
With the goal of minimising the risk of accidentally mutating
shared state between runs. And promoting reuse of test boilerplate to
speed up testing, and help get started with trickier testing situations.

There's no big up front work to implement it, and when you don't need it you
can stick to `testing.T` as usual.

# Use

See the [package example](https://pkg.go.dev/github.com/barry-hennessy/test/sweet#example-package)
for a quick intro. And other examples for details.

# Reuse & skipping boilerplate

`DepFactory` functions are the interface that can be centralised, reused
and shared.

For example common database or other external dependencies can have a set of
stock test containers that can spin themselves up for every test and clean
themselves after. All defined once and reused throughout every project.
The same goes for fixture loaders, http Servers etc.

This is probably sweet's biggest power: a shared pattern for seperating test
setup from the test itself. That pattern being a building block you can stack up
and go higher with.

See:
  - [Factories for testcontainers](https://github.com/barry-hennessy/test/tree/main/sweet/factories/tc)

## Links
 - [Project overview and background](https://barryhennessy.com/projects/test-sweet/)
