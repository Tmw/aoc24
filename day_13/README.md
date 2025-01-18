# Day 13

## Run

```console
go run main.go < input.txt
```

## Notes

- Initial thought: this sounds like a problem for DFS but this was not the right tool for the job. It found the answer for the first machine quite quickly, however finding out that the second machine was not solvable took forever.
- Solving part two with a simple and naively nested for-loop got me the answer in about four milliseconds.
- Obviously for part two this didn't scale.. at all! So after trying some sort of binary search approach and some other brute force mechanisms I decided to take a hint and treat it as a math problem. Readin up on linear algebra taught me Cramers rule again which ultimately gave me the answer to both parts in under a millisecond.
