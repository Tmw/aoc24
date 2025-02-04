# Day 15

## Run

```console
go run main.go < input.txt
```

## Notes

- Nothing really challenging, but fun assignment to work out and worth letting the animation play out in the terminal.
- Detecting movable boxes in a chain was easy enough for part one by using recursion
- Found the answer to part one in about 307µs
- Part two is giving some issues => 1488820 is too high
    - probably let it run again, writing out every frame.. See if we have odd behavior?
    - perhaps we can run the example input again, see if that number still tracks? It doesn't..
    - perhaps rewrite it so we can set the state like we want it (widened) and write test against certain scenario's in terms of moving,
    - i think it now moves too few boxes around?
    - moving away from simple grid system, getting the algorithm right when doing recursion
    and trying to match up [ with ] is too much of a hassle. Will try out a free-form system instead?
