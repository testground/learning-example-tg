package tgsync

import "github.com/testground/sdk-go/runtime"

const (
	ProducerRole = "producer"
	ConsumerRole = "consumer"
)

// Returns instance's role, based on the passed run params or instance's sequence ID,
// if the run param is unset
func getInstanceRole(runenv *runtime.RunEnv, seqId int64) string {

	if runenv.IsParamSet("role") {
		var roleParam = runenv.StringParam("role")
		if roleParam == ProducerRole {
			return ProducerRole
		} else if roleParam == ConsumerRole {
			return ConsumerRole
		}
	}

	if seqId == 1 {
		return ConsumerRole
	} else {
		return ProducerRole
	}
}
