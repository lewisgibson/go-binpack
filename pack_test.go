package binpack_test

import (
	"testing"

	"github.com/lewisgibson/go-binpack"
	"github.com/stretchr/testify/require"
)

// testPackable implements binpack.Packable for testing purposes.
// It records the provided rectangles and the placements made.
type testPackable struct {
	rectangles []binpack.Rectangle
	placements []struct{ x, y int }
}

// Ensure that testPackable implements the binpack.Packable interface.
var _ binpack.Packable = (*testPackable)(nil)

// newTestPackable creates a new testPackable with the provided rectangles.
func newTestPackable(rects []binpack.Rectangle) *testPackable {
	return &testPackable{
		rectangles: rects,
		placements: make([]struct{ x, y int }, len(rects)),
	}
}

// Len returns the number of rectangles.
func (tp *testPackable) Len() int {
	return len(tp.rectangles)
}

// Rectangle returns the rectangle at the specified index.
func (tp *testPackable) Rectangle(n int) binpack.Rectangle {
	return tp.rectangles[n]
}

// Place records the placement of the rectangle at the specified index.
func (tp *testPackable) Place(n, x, y int) {
	tp.placements[n].x = x
	tp.placements[n].y = y
}

// rectanglesOverlapTest returns true if the two rectangles intersect.
// The rectangles are defined by their top-left (x,y) and dimensions.
func rectanglesOverlapTest(x1, y1, w1, h1, x2, y2, w2, h2 int) bool {
	if x1 >= x2+w2 || x2 >= x1+w1 {
		return false
	}
	if y1 >= y2+h2 || y2 >= y1+h1 {
		return false
	}
	return true
}

// TestPack_NoRectangles verifies that an empty Packable returns (0,0).
func TestPack_NoRectangles(t *testing.T) {
	t.Parallel()

	// Arrange: create a test packable with no rectangles.
	tp := newTestPackable([]binpack.Rectangle{})

	// Act: pack the rectangles.
	w, h := binpack.Pack(tp)

	// Assert: dimensions should be (0, 0).
	require.Equal(t, 0, w, "expected width 0 for no rectangles")
	require.Equal(t, 0, h, "expected height 0 for no rectangles")
}

// TestPack_SingleRectangle verifies that a single rectangle is placed at (0,0)
// and that the overall dimensions match its size.
func TestPack_SingleRectangle(t *testing.T) {
	t.Parallel()

	// Arrange: create a test packable with one rectangle.
	tp := newTestPackable([]binpack.Rectangle{
		{Width: 100, Height: 200},
	})

	// Act: pack the rectangle.
	w, h := binpack.Pack(tp)

	// Assert: overall dimensions should equal the rectangle's size.
	require.Equal(t, 100, w, "expected width 100")
	require.Equal(t, 200, h, "expected height 200")

	//Assert: the rectangle should be placed at (0, 0).
	require.Equal(t, 0, tp.placements[0].x, "expected x-coordinate 0")
	require.Equal(t, 0, tp.placements[0].y, "expected y-coordinate 0")
}

// TestPack_MultipleRectangles verifies that multiple rectangles are arranged
// into a compact, non-overlapping layout.
func TestPack_MultipleRectangles(t *testing.T) {
	t.Parallel()

	// Arrange: create a test packable with several rectangles.
	rectangles := []binpack.Rectangle{
		{Width: 100, Height: 200},
		{Width: 50, Height: 50},
		{Width: 80, Height: 120},
		{Width: 30, Height: 60},
		{Width: 70, Height: 70},
	}
	tp := newTestPackable(rectangles)

	// Act: pack the rectangles.
	w, h := binpack.Pack(tp)

	// Assert: overall dimensions should be non-zero.
	require.Positive(t, w, "expected positive overall width")
	require.Positive(t, h, "expected positive overall height")

	// Assert: all placements should be non-negative.
	for i, p := range tp.placements {
		require.GreaterOrEqual(t, p.x, 0, "placement x for rectangle %d should be non-negative", i)
		require.GreaterOrEqual(t, p.y, 0, "placement y for rectangle %d should be non-negative", i)
	}

	// Assert: rectangles should not overlap.
	for i := 0; i < len(rectangles); i++ {
		for j := i + 1; j < len(rectangles); j++ {
			require.False(t, rectanglesOverlapTest(
				tp.placements[i].x, tp.placements[i].y,
				rectangles[i].Width, rectangles[i].Height,
				tp.placements[j].x, tp.placements[j].y,
				rectangles[j].Width, rectangles[j].Height,
			), "expected rectangle %d and %d not to overlap", i, j)
		}
	}
}

// TestPack_TenRectangles verifies that a set of ten rectangles is packed
// into a compact, non-overlapping layout.
func TestPack_TenRectangles(t *testing.T) {
	t.Parallel()

	// Arrange: create a test packable with ten rectangles.
	rectangles := []binpack.Rectangle{
		{Width: 100, Height: 200},
		{Width: 150, Height: 150},
		{Width: 80, Height: 120},
		{Width: 50, Height: 70},
		{Width: 60, Height: 90},
		{Width: 120, Height: 80},
		{Width: 200, Height: 100},
		{Width: 40, Height: 40},
		{Width: 90, Height: 110},
		{Width: 70, Height: 130},
	}
	tp := newTestPackable(rectangles)

	// Act: pack the rectangles.
	w, h := binpack.Pack(tp)

	// Assert: overall dimensions should be positive.
	require.Positive(t, w, "expected positive overall width")
	require.Positive(t, h, "expected positive overall height")

	// Assert: all placements should be non-negative.
	for i, p := range tp.placements {
		require.GreaterOrEqual(t, p.x, 0, "placement x for rectangle %d should be non-negative", i)
		require.GreaterOrEqual(t, p.y, 0, "placement y for rectangle %d should be non-negative", i)
	}

	// Assert: rectangles should not overlap.
	for i := 0; i < len(rectangles); i++ {
		for j := i + 1; j < len(rectangles); j++ {
			require.False(t, rectanglesOverlapTest(
				tp.placements[i].x, tp.placements[i].y,
				rectangles[i].Width, rectangles[i].Height,
				tp.placements[j].x, tp.placements[j].y,
				rectangles[j].Width, rectangles[j].Height,
			), "expected rectangle %d and %d not to overlap", i, j)
		}
	}
}
