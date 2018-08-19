/*

Copyright 2018 This Project Authors.

Author:  seanchann <seanchann@foxmail.com>

See docs/ for more information about the  project.

*/

package install

import (
	"testing"

	samplefuzzer "github.com/seanchann/sample-apiserver/pkg/apis/sample/fuzzer"
	"k8s.io/apimachinery/pkg/api/testing/roundtrip"
)

func TestRoundTripTypes(t *testing.T) {
	roundtrip.RoundTripTestForAPIGroup(t, Install, samplefuzzer.Funcs)
}
