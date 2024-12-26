# Day 7

## Run

```console
go run main.go < input.txt
```

## Notes

- Initial implementation was to generate all permutations of operations and apply them in order, checking if we'd reach the target number.
- Adding a cache (`map[int][][]Op`) helped reduce the runtime for part 1 to ~15ms.
- Adding a single operator (part 2) to the permutation raised the runtime to ~3s though.
- Switching to an iterator so it would lazily compute the permutations and quit as soon as the sum was found raised the runtime to ~7.2s, easily explained by no longer having the cache to fall back on.
- Switching to DFS where we recursively call a check function comparing against the sum and recurse into calling itself three times with the different operations and new sum. This brought the total runtime down to ~1.4sec for part two and 10ms for part one.
- Adding an early return if the sum already exceeds the target brought it down to ~1.02sec
- Replacing the `fmt.Sprintf` and `strconv.Atoi` for the concat operation, with a more improved implementation (thanks to ChatGPT) we're now at ~120ms for part two and ~7ms for part one.
