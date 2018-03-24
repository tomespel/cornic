package cornic

import (
	"fmt"

	"../cornic/fin"
)

// TestStrategy tests that the strategy file is active
func TestStrategy() bool {
	a := fin.NewAccount("test", "GBP", 20.0, 10.0, "myid")
	fmt.Println(a.ID)
	fmt.Println(a.LastUpdated)
	return fin.TestIncentive()
}
