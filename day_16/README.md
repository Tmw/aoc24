# Day 16

## Run

```console
go run . < input.txt
```

## Notes
- First implementation implemented using BFS. After some minor tweaks this worked flawless for the two examples, however it is quite time consuming for the real input. Letting it run until it completes for now, but i'll need to do some profiling.
- It finds the finish three times rather quick, but then seems to hang there for forever.
- let's first try to improve the lookup speed for already visited nodes (don't think that's a major win, but who knows)
- re-implemented this using A*, both examples return the correct path and cost however the real input is failing me again.
