package brtool

// InstringSlice 判断一个元素是否在string类型的slice里面
func InstringSlice(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
