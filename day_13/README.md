# Day 13

## Run

```console
go run main.go < input.txt
```

## Notes

- Initial thought: this sounds like a problem for DFS but this was not the right tool for the job. It found the answer for the first machine quite quickly, however finding out that the second machine was not solvable took forever.
- Solving part two with a simple and naively nested for-loop got me the answer in about four milliseconds.
- Obviously for part two this didn't scale.. at all! So after trying some sort of binary search approach and some other brute force mechanisms I decided to take a hint and treat it as a math problem. Readin up on linear algebra taught me Cramers rule again which ultimately gave me the answer to both parts in under a millisecond.

### Cramers Rule

# Learning about Cramer's rule

first thing first is we need to grab the determinent.
In a 2x2 grid this is pretty straight forward:

### equation
in the following linear equation system equation:

a * 94 + b * 22 = 8400
a * 34 + b * 67 = 5400

### matrix

[ 94 22 ]
[ 34 67 ]

### determinant
determinant is calculated by multiplying top-left and bottom right and subtract top-right times bottom-left from that, so:
d = 94 * 67 - 22 * 34
d = 5550

### now for each of the unknowns
replace the first column in the matrix with the constants from after the eq sign.

For finding X:

dx = [ 8400 22 ]
     [ 5400 67 ]

and repeat the same process:
dx = 8400 * 67 - 22 * 5400
dx = 444000
x = dx / d
x = 80

And now to find Y:

dy = [ 94 8400 ]
     [ 34 5400 ]
dy = 94 * 5400 - 8400 * 34
dy = 222000
y = dy / d
y = 40

### Next steps
let's try to apply this on a 3x3 matrix too.

__Equation:__
3x + 3y +  5z = 1
3x + 5y +  9z = 0
5x + 9y + 17z = 0

__matrix:__

[ 3 3  5 ]
[ 3 5  9 ]
[ 5 9 17 ]

now, to take the determenant, we'll need to add the first two columns to the end again so we always have the same number of elements to work with.

Added the first two element divided by semicolon:
[ 3 3  5 ;3 3 ]
[ 3 5  9 ;3 5 ]
[ 5 9 17 ;5 9 ]

then taking from the first row the diagonals down and to the right, multiplying them and adding those together,
then taking the diagonals from the last element of the first row (5) down and to the left and subtracting them

det =  (3 * 5 * 17) + (3 * 9 * 5) + (5 * 3 * 9) - 
        (5 * 5 * 5) - (3 * 9 * 9) - (3 * 3 * 17)
det = 4

now that we have the determinant of the whole matrix,
let's solve for each X, Y and Z determinant.


dx = [ 1 3  5; 1  3 ]
     [ 0 5  9; 0  5 ]
     [ 0 9 17; 0 17 ]
dx = 1 * 5 * 17 + 3 * 9 * 0 + 5 * 0 * 17 - 5 * 5 * 0 - 1 * 9 * 9 - 3 * 0 * 17
x = dx / d
x = 1


dy = [ 3 1  5; 3 1 ]
     [ 3 0  9; 3 0 ]
     [ 5 0 17; 5 0 ]
dy =  3 * 0 * 17 + 1 * 9 * 5 + 5 * 3 * 0 - 5 * 0 * 5 - 3 * 9 * 0 - 1 * 3 * 17
y = dy / d
y = -1.5


dz = [ 3 3 1 ;3 3 ]
     [ 3 5 0 ;3 5 ]
     [ 5 9 0 ;5 9 ]

dz = 3 * 5 * 0 + 3 * 0 * 5 + 1 * 3 * 9 - 1 * 5 * 5 - 3 * 0 * 9 - 3 * 3 * 0
z = dz / d
z = 0.5


so, proof:
3 * 1 + 3 * -1.5 +  5 * 0.5 = 1
3 * 1 + 5 * -1.5 +  9 * 0.5 = 0
5 * 1 + 9 * -1.5 + 17 * 0.5 = 0

and it matches.

### Gauss Jordan Elimination
From the input:
---------------
Button A: X+94, Y+34
Button B: X+22, Y+67
Prize: X=8400, Y=5400

As a linear equation system:
----------------------------
a * 94 + b * 22 = 8400
a * 34 + b * 67 = 5400

As an augmented matrix
----------------------
[ 94 22 | 8400 ] -> [ 1 0 | ?? ]
[ 34 67 | 5400 | -> [ 0 1 | ?? ]

Gauss-jordan elimination:
-------------------------
[ 1  22/94 | 8400/94 ] R1 * 1/94 -> R1
[ 34 67    | 5400 ]

[ 1  22/94 | 8400/94 ]
[ 0 33 | 5366] R2 -> R2 - (34 * R1)

by applying elementary operations such as swap, pivot and scale we should try to get the matrix down to something that looks like this:

[1 0 | X]
[0 1 | Y]

where X and Y are the actual X and Y values of the system.

This works for N dimensions, here's one with three
[1 0 0 | X]
[0 1 0 | Y]
[0 0 1 | Z]

this matrix is in "reduced echelon form"

again - we can employ swapping rows, scaling rows by positive numbers or fractions or pivoting where we do something similar but based on another row.

As long as we do the same operation to the entire row, meaning also the augmented bit on the right hand side of the pipe.

and once we got it all the way down to a reduced form, it should directly give us the answer.

Not entirely sure how we'd do this programmatically tho, but its should be possible.

we don't have to go all the way, we can also go up-to echelon form, this is where we have the following structure:
[ 1 34 55 | X ]
[ 0 1  55 | Y ]
[ 0 0  1  | Z ]

as you can see we still have numbers that are non-zero and non-one - but this would already give us the answer to Z which we then can use to find the answer of Y which we then can use to find the answer in X.
