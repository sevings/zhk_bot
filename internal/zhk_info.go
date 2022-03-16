package internal

const (
	minFloor      = 2
	maxFloor      = 23
	floorCount    = maxFloor - minFloor + 1
	buildingCount = 4
	minFlat       = 1
	liftCount     = 3
)

var (
	flatPerFloor    []int
	flatPerBuilding []int
	maxFlat         int
)

func init() {
	flatPerFloor = []int{10, 9, 8, 9}

	for _, k := range flatPerFloor {
		flatPerBuilding = append(flatPerBuilding, k*floorCount)
	}

	for _, k := range flatPerBuilding {
		maxFlat += k
	}
}

func getBuilding(flat int) int {
	for i := 0; i < buildingCount; i++ {
		if flat <= flatPerBuilding[i] {
			return i + 1
		}

		flat -= flatPerBuilding[i]
	}

	return 0
}

func getFloor(flat int) int {
	var i int
	for ; i < buildingCount; i++ {
		if flat <= flatPerBuilding[i] {
			break
		}

		flat -= flatPerBuilding[i]
	}

	return (flat-1)/flatPerFloor[i] + minFloor
}

func getMinBuildingFlat(building int) int {
	flat := 1

	for i := 0; i < building-1; i++ {
		flat += flatPerBuilding[i]
	}

	return flat
}

func getMaxBuildingFlat(building int) int {
	flat := 0

	for i := 0; i < building; i++ {
		flat += flatPerBuilding[i]
	}

	return flat
}

func getMinFloorFlat(building, floor int) int {
	flat := getMinBuildingFlat(building)
	flat += (floor - minFloor) * flatPerFloor[building-1]

	return flat
}

func getMaxFloorFlat(building, floor int) int {
	flat := getMinBuildingFlat(building) - 1
	flat += (floor - minFloor + 1) * flatPerFloor[building-1]

	return flat
}
