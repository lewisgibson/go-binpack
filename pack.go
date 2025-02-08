package binpack

import (
	"math"
	"sort"
)

// Rectangle represents the dimensions of a rectangle.
type Rectangle struct {
	Width, Height int
}

// Area returns the area of the rectangle.
func (r Rectangle) Area() int {
	return r.Width * r.Height
}

// Packable is the interface for types that support rectangle packing.
type Packable interface {
	Len() int
	Rectangle(n int) Rectangle
	Place(n, x, y int)
}

// placement represents a rectangle placed at a specific position.
type placement struct {
	position, x, y, width, height int
}

// bounds represents the bounding box for a set of rectangles.
type bounds struct {
	minX, minY, maxX, maxY int
}

// Pack arranges rectangles into a compact layout. Larger rectangles are
// placed first to reduce conflicts. The final layout is shifted so that its
// top-left corner is at (0, 0). Returns the overall dimensions.
func Pack(p Packable) (int, int) {
	var count = p.Len()
	if count == 0 {
		return 0, 0
	}

	var positions = make([]int, count)
	for i := 0; i < count; i++ {
		positions[i] = i
	}

	// Sort the positions to prioritize larger rectangles first.
	sort.Slice(positions, func(i, j int) bool {
		return p.Rectangle(positions[i]).Area() > p.Rectangle(positions[j]).Area()
	})

	var placements []placement
	for _, position := range positions {
		var rectangle = p.Rectangle(position)
		if len(placements) == 0 {
			placements = append(placements, placement{
				position: position,
				x:        0,
				y:        0,
				width:    rectangle.Width,
				height:   rectangle.Height,
			})
			continue
		}

		// Derive candidate positions from existing rectangle edges.
		var xCandidates, yCandidates = getCandidatePositions(placements)
		var bounds = computeBounds(placements)

		// Choose the candidate that minimizes the overall bounding box and is as centered as possible.
		var bestX, bestY, candidateFound = findBestPlacement(xCandidates, yCandidates, bounds, rectangle, placements)
		if !candidateFound {
			bestX = bounds.maxX
			bestY = bounds.minY
		}

		placements = append(placements, placement{
			position: position,
			x:        bestX,
			y:        bestY,
			width:    rectangle.Width,
			height:   rectangle.Height,
		})
	}

	// Place all of rectangles at their final positions.
	var bounds = computeBounds(placements)
	for _, placement := range placements {
		p.Place(placement.position, placement.x-bounds.minX, placement.y-bounds.minY)
	}

	// Return the overall dimensions.
	return bounds.maxX - bounds.minX, bounds.maxY - bounds.minY
}

// expandBoundsForPlacement expands b to include rectangle r.
func expandBoundsForPlacement(r placement, b bounds) bounds {
	if r.x < b.minX {
		b.minX = r.x
	}
	if r.y < b.minY {
		b.minY = r.y
	}
	if r.x+r.width > b.maxX {
		b.maxX = r.x + r.width
	}
	if r.y+r.height > b.maxY {
		b.maxY = r.y + r.height
	}
	return b
}

// computeBounds returns the minimal bounding box enclosing all rectangles.
func computeBounds(placements []placement) bounds {
	var b = bounds{
		minX: placements[0].x,
		minY: placements[0].y,
		maxX: placements[0].x + placements[0].width,
		maxY: placements[0].y + placements[0].height,
	}
	// Iterate over all placements, expand the bounding box if necessary.
	for _, r := range placements {
		if r.x < b.minX {
			b.minX = r.x
		}
		if r.y < b.minY {
			b.minY = r.y
		}
		if r.x+r.width > b.maxX {
			b.maxX = r.x + r.width
		}
		if r.y+r.height > b.maxY {
			b.maxY = r.y + r.height
		}
	}
	return b
}

// getCandidatePositions extracts unique x and y coordinates from the edges of placed rectangles.
func getCandidatePositions(rects []placement) ([]int, []int) {
	var x, y = make(map[int]bool), make(map[int]bool)
	for _, r := range rects {
		x[r.x] = true
		x[r.x+r.width] = true
		y[r.y] = true
		y[r.y+r.height] = true
	}

	var xCandidates []int
	for x := range x {
		xCandidates = append(xCandidates, x)
	}

	var yCandidates []int
	for y := range y {
		yCandidates = append(yCandidates, y)
	}

	return xCandidates, yCandidates
}

// doRectanglesIntersect returns true if rectangles a and b intersect.
func doRectanglesIntersect(a, b placement) bool {
	if a.x >= b.x+b.width || b.x >= a.x+a.width {
		return false
	}
	if a.y >= b.y+b.height || b.y >= a.y+a.height {
		return false
	}
	return true
}

// hasIntersection checks if candidate intersects any rectangle in rects.
func hasIntersection(candidate placement, placements []placement) bool {
	for _, p := range placements {
		if doRectanglesIntersect(candidate, p) {
			return true
		}
	}
	return false
}

// findBestPlacement selects the candidate position that minimizes the overall bounding box area,
// favoring positions whose center is closer to the center of the expanded bounding box.
// The area and center are computed inline.
func findBestPlacement(xCandidates, yCandidates []int, b bounds, r Rectangle, placements []placement) (int, int, bool) {
	// Allocate state for the heuristic.
	var bestX, bestY int
	var bestArea = math.MaxInt64
	var bestCenterDistance = math.MaxInt64
	var found = false

	// Evaluate all candidate positions.
	for _, candidateX := range xCandidates {
		for _, candidateY := range yCandidates {
			var candidate = placement{
				x:      candidateX,
				y:      candidateY,
				width:  r.Width,
				height: r.Height,
			}

			// If the candidate intersects any existing rectangle, skip it.
			if hasIntersection(candidate, placements) {
				continue
			}

			candidateBB := expandBoundsForPlacement(candidate, b)
			// Inline area calculation.
			candidateArea := (candidateBB.maxX - candidateBB.minX) * (candidateBB.maxY - candidateBB.minY)
			// Inline center calculation.
			bbCenterX := candidateBB.minX + (candidateBB.maxX-candidateBB.minX)/2
			bbCenterY := candidateBB.minY + (candidateBB.maxY-candidateBB.minY)/2
			candidateCenterX := candidate.x + candidate.width/2
			candidateCenterY := candidate.y + candidate.height/2
			dx := candidateCenterX - bbCenterX
			dy := candidateCenterY - bbCenterY
			centerDistance := dx*dx + dy*dy

			if candidateArea < bestArea || (candidateArea == bestArea && centerDistance < bestCenterDistance) {
				bestArea = candidateArea
				bestCenterDistance = centerDistance
				bestX = candidate.x
				bestY = candidate.y
				found = true
			}
		}
	}

	return bestX, bestY, found
}
