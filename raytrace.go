package voxelraytrace

import (
	"errors"
	"github.com/go-gl/mathgl/mgl64"
	"math"
)

// InDirection performs a ray trace from the start position in the given direction, for a distance of the maxDistance.
// This returns a Generator which yields Vector3s containing the coordinates of voxels it passes through.
func InDirection(start, directionVector mgl64.Vec3, maxDistance float64) (vectors []mgl64.Vec3, err error) {
	return BetweenPoints(start, start.Add(directionVector.Mul(maxDistance)))
}

// BetweenPoints performs a ray trace between the start and end coordinates.
// This returns an array of vectors containing the coordinates of voxels it passes through.
// http://www.cse.yorku.ca/~amana/research/grid.pdf
func BetweenPoints(start, end mgl64.Vec3) (vectors []mgl64.Vec3, err error) {
	currentPoint := mgl64.Vec3{math.Floor(start.X()), math.Floor(start.Y()), math.Floor(start.Z())}

	directionVector := end.Sub(start).Normalize()
	if directionVector.LenSqr() <= 0 {
		return nil, errors.New("start and end points are the same, giving a zero direction vector")
	}

	radius := distance(start, end)

	stepX := compareTo(directionVector.X(), 0)
	stepY := compareTo(directionVector.Y(), 0)
	stepZ := compareTo(directionVector.Z(), 0)

	tMaxX := rayTraceDistanceToBoundary(start.X(), directionVector.X())
	tMaxY := rayTraceDistanceToBoundary(start.Y(), directionVector.Y())
	tMaxZ := rayTraceDistanceToBoundary(start.Z(), directionVector.Z())

	tDeltaX := findDelta(directionVector.X(), stepX)
	tDeltaY := findDelta(directionVector.Y(), stepY)
	tDeltaZ := findDelta(directionVector.Z(), stepZ)

	for {
		vectors = append(vectors, currentPoint)

		if tMaxX < tMaxY && tMaxX < tMaxZ {
			if tMaxX > radius {
				break
			}
			currentPoint = currentPoint.Add(mgl64.Vec3{stepX})
			tMaxX += tDeltaX
		} else if tMaxY < tMaxZ {
			if tMaxY > radius {
				break
			}
			currentPoint = currentPoint.Add(mgl64.Vec3{0, stepY})
			tMaxY += tDeltaY
		} else {
			if tMaxZ > radius {
				break
			}
			currentPoint = currentPoint.Add(mgl64.Vec3{0, 0, stepZ})
			tMaxZ += tDeltaZ
		}
	}

	return
}

// findDelta finds the change in t on an axis when taking a step on that axis (always positive).
func findDelta(first, second float64) float64 {
	if first == 0 {
		return 0
	}
	return second / first
}

// rayTraceDistanceToBoundary returns the distance that must be travelled on an axis
// from the start point with the direction vector component to cross a block boundary.
func rayTraceDistanceToBoundary(first, second float64) float64 {
	if second == 0 {
		return math.Inf(0)
	}

	if second < 0 {
		first = -first
		second = -second

		if math.Floor(first) == first {
			return 0
		}
	}

	return (1 - (first - math.Floor(first))) / second
}

// compareTo compares the first and second float. It returns 0 if they are both the same,
// -1 if the second float is bigger than the first, and 1 if the first is bigger than the second.
// It is similar to the spaceship operator in PHP, and the Java comparable class.
func compareTo(first, second float64) float64 {
	if first == second {
		return 0
	} else if first < second {
		return -1
	} else {
		return 1
	}
}

// distance measures the distance between two vectors.
func distance(a, b mgl64.Vec3) float64 {
	xDiff, yDiff, zDiff := b[0]-a[0], b[1]-a[1], b[2]-a[2]
	return math.Sqrt(xDiff*xDiff + yDiff*yDiff + zDiff*zDiff)
}
