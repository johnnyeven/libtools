package gen_method

import (
	"testing"
)

func TestScanner(t *testing.T) {
	scanner := NewScanner("github.com/johnnyeven/libtools/gen_method/examples")
	scanner.Output("CustomerG7", false)
	scanner.Output("User", true)
	scanner.Output("PhysicsDeleteByUniquustomerG7", false)
}
