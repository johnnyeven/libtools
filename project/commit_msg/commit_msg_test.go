package commit_msg

import (
	"testing"
)

func Test_CheckCommit(t *testing.T) {
	t.Log(CheckCommit("FIX ZF-1 do what for ..."))
	t.Log(CheckCommit("FIX ZF-1 xxxxx"))
	t.Log(CheckCommit("BUG"))
	t.Log(CheckCommit("TASK "))
}
