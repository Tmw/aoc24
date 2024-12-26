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
- Switching to DFS where we recursively call a check function comparing against the sum and recurse into calling itself three times with the different operations and new sum. This brought the total runtime down to ~1.02 sec for part two and 10ms for part one.
