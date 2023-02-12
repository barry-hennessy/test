# Testcontainer factories

Test container factories are `sweet.DepFactory` implementations that bootstrap
starting containers for testing.

The containers
 - clean themselves up after tests
 - start up ready to serve traffic
   - client tests should never have to call `time.Sleep` before testing
   - client tests should never fail due to a container dependency not being
     ready
 - start tests as soon as they're ready to serve traffic
   - client tests should not be waiting because a factory is calling `time.Sleep`
