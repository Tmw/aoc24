# Day 11

## Run

```console
go run main.go < input.txt
```

## Notes

- For the initial implementation I picked a linked list as it's easy to insert a new node (Stone) in place.
- This worked out fine for part one as iterating over the linked list 25 times to get the answer finished in ~35ms.
- Running the same implementation for part two (blinking 75 times) this never succesfully finished and started to consume over 14GB in memory, it's probably safe to say storing all stones in a linked list isn't viable at this point anymore.
- After taking a hint I'v decided to utilize a map where we map between the number that's on the stone and the number of stones with that number and no longer bother with actually maintaining a list of stones.
- The order of the stones doesn't seem to matter so using a map and applying the blink-rules on each of the keys to form a new map should do just fine.
- After rewriting to a map-based approach, this seems to finish in 620µs for part one and 10ms for part two.
