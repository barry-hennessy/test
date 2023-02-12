VERSION          := snapshot
NAME             := $(shell basename $(CURDIR))

PROFILE_FLAGS    := -v -count=1
TEST_FLAGS       := $(PROFILE_FLAGS) -race

default: test 

# @TODO: As it stands the targets depending on %/.go won't rebuild if packages
# those packages depend on are updated. The PHONY's that depend on these will
# rebuild anyway though, so it's OK for now.
cover.out: *.go
	go test $(PROFILE_FLAGS) -cover -coverprofile cover.out ./...

block.out: *.go
	go test $(PROFILE_FLAGS) -blockprofile block.out ./...

cpu.out: *.go
	go test $(PROFILE_FLAGS) -cpuprofile cpu.out -bench ./...

mem.out: *.go
	go test $(PROFILE_FLAGS) -memprofile mem.out -bench ./...

mutex.out: *.go
	go test $(PROFILE_FLAGS) -mutexprofile mutex.out ./...

# See Datadog's wonderful overview of the profiling tools, go internals and how
# to interpret the results here: https://github.com/DataDog/go-profiler-notes/blob/main/guide/README.md#block
# @TODO: For CI it might be desirable to have a static output, not a web interface
.PHONY: pprof/%
pprof/%: %.out
	go tool pprof -http=localhost:8001 $<

trace.out: *.go
	go test $(PROFILE_FLAGS) -trace trace.out ./...

.PHONY: trace
trace: trace.out
	go tool trace trace.out

.PHONY: test
test:
	go test $(TEST_FLAGS) ./...
	golangci-lint run ./...

.PHONY: clean
clean:
	rm -rvf build coverprofile.txt

.PHONY: tidy
tidy:
	go mod tidy
