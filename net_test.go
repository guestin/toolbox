package mod

import (
	"encoding/json"
	"github.com/guestin/mob/merrors"
	"github.com/guestin/mob/mjson"
	"testing"
)

func TestListIf(t *testing.T) {
	netIfs, err := NetworkList(SkipDefault)
	merrors.AssertError(err, "list network")
	jstr, err := json.MarshalIndent(netIfs, mjson.KJsonIndentPrefix, mjson.KJsonIndent)
	merrors.AssertError(err, "marshal indent")
	t.Log(string(jstr))
	for _, net := range netIfs {
		t.Log("Name=", net.Name)
		t.Log("Mac=", net.Mac)
		for _, addr := range net.Address {
			t.Log("IP=", addr.IP)
			t.Log("Mask=", addr.Mask)
		}
		t.Log()
	}
}
