package utils

func MapKey2Slice(mapData map[uint]uint) []uint {
	outList := make([]uint, 0)
	for key, _ := range mapData {
		outList = append(outList, key)
	}
	return outList
}

func MapVal2Slice(mapData map[uint]uint) []uint {
	outList := make([]uint, 0)
	mapValData := make(map[uint]uint)
	for _, Val := range mapData {
		mapValData[Val] = Val
	}
	for key, _ := range mapValData {
		outList = append(outList, key)
	}
	return outList
}
