package crawler

import (
	"fmt"
	"time"
)

func GetCoordinates() {
	i := 0
	for range time.Tick(time.Second * 1) {
		fmt.Println(i)
		i++
	}
}
