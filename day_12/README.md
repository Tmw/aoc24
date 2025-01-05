# Day 12

## Run

```console
go run main.go < input.txt
```

## Notes

- First approach i thought of is having an array with Tiles, each tile keeps track of its type (byte) and the fences around it using a bitmask. We're starting off with each tile being completely fenced off (0b1111) and then implement an algorithm that iterates the tiles and removes fences to neighbouring cells when the type is the same (removing a fence is done using XOR-ing).
- By iterating the map a final time we can add up all tiles with the same type to get the area, and call `CountBorders()` on each cell to get the perimeter (CountBorders just counts the 1's set in the Border prop).
- Alternative approach I thought of is to find the smallest and greatest X,Y positions of each "region". The delta's would then indicate the length of the borders, however this would still require us to group by type.
- Besides; i have an inkling that for part two we want to save some money and only have a single fence between two neighbouring regions :-)
- Ok, my initial approach worked fine for the simplest example, however it wouldn't correctly calculate the cost for the XOXO example. Rather than treating each X as its own area of one and perimeter of 4, it would calculate over the total area of 4 and perimeter of 16. Rewrote it to use an algorithm to find clusters, and calculate area and perimeter of each cluster now instead. Found the anawer succesfully.
