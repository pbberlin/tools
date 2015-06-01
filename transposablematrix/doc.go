/*

Transposable Matrix is a two dimensional slice of Slots.
Change the type of Slot field A, to accommodate yourType.

You can now View, Get and Set your Slots from *four* perspectives.

0 - North
1 - West
2 - South
3 - East
(opposite of classic watch direction)

The coordinates are called mapped ("perspectivized") xm, ym.

The coordinates of the constituting slices are called *base* coordinates xb, yb.

Base coordinates are checked for boundaries.
Requested Slots outside physical limits are reported as empty slots.
Todo: return an error - similar to map

You can further change the center of mapped coordinates, independently of the perspective.



For example see:
	visual_test.go - go test

*/
package transposablematrix
