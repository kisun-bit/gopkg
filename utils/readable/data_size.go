package readable

import (
	"fmt"
	"strconv"
)

func HumanizeBytes(bytesNum int64) string {
	var size string

	if valPB := bytesNum / (1 << 50); valPB != 0 {
		num1, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(bytesNum)/float64(1<<50)), 64)
		size = fmt.Sprintf("%fPB", num1)
	} else if valTB := bytesNum / (1 << 40); valTB != 0 {
		num2, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(bytesNum)/float64(1<<40)), 64)
		size = fmt.Sprintf("%fTB", num2)
	} else if valGB := bytesNum / (1 << 30); valGB != 0 {
		num3, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(bytesNum)/float64(1<<30)), 64)
		size = fmt.Sprintf("%fGB", num3)
	} else if valMB := bytesNum / (1 << 20); valMB != 0 {
		num4, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(bytesNum)/float64(1<<20)), 64)
		size = fmt.Sprintf("%fMB", num4)
	} else if valKB := bytesNum / (1 << 10); valKB != 0 {
		num5, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(bytesNum)/float64(1<<10)), 64)
		size = fmt.Sprintf("%fKB", num5)
	} else {
		size = fmt.Sprintf("%fB", bytesNum)
	}

	return size
}
