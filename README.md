# Test

Make your code work. Test.

Make your designs clearer. Test.

Make your work faster. Test.

Make your life easier.

TEST

## And how does _this_ help?
This is a collection of modules that help you test in various ways.
It is not a one-true-framework that forces you to do things it's way.
The modules here try keep as close to the go standard way of doing things as
possible.

Keep it lean, use `testing.T`, get it done. And sleep soundly.

## Modules

### [Sweet](./sweet)
A simple test suite 'framework'; if you can even call it a framework.
`t.Run` on steroids.

### [Sweet testcontainers](./sweet/factories/tc)
Run testcontainers with `sweet`. All containers start up as soon as they're
ready to be tested and clean up after themselves.

Roll your own container or use the pre-set up ones:
 - [redis](sweet/factories/tc/redis)
 - [mongodb](sweet/factories/tc/mongodb)
 - [postgres](sweet/factories/tc/postgres)
 - [cockroachdb](sweet/factories/tc/cockroachdb)
