# Day 11

## Run

```console
go run main.go < input.txt
```

## Notes

- For the initial implementation I picked a linked list as it's easy to insert a new node (Stone) in place.
- This worked out fine for part one as iterating over the linked list 25 times to get the answer finished in ~35ms.
- Running the same implementation for part two (blinking 75 times) this never succesfully finished.
