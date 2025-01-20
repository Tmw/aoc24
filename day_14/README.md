# Day 14

## Run

```console
go run main.go < input.txt
```

## Notes

- Fairly quickly I figured this should be solvable using modulo due to the wrapping and it wouldn't require actually looping 100 times. Fully expecting this assignment to add a few trillion to it as well ;)
- Counting the robots per quadrant and multiplying them together was easy.
- This got me the answer in less than a milisecond.
- I do not understand part two just yet.. Likely we'll need to animate it in the terminal to see if a christmas tree pattern emerges.. but that'll be for another day =)
- ok, ended up same day.. Assumed this would be in ascii art and implemented a very scrummy loop that would break if it found 10 robots in sequence (arbitrarily chosen) and let it run for 20_000 seconds to start with.. And it found it relatively fast (1.3sec)! Side note: this was seriously awesome, didn't see this coming at all. Kudos :)

