/*

Copyright 2018 This Project Authors.

Author:  seanchann <seanchann@foxmail.com>

See docs/ for more information about the  project.

*/

package fuzzer

import (
	fuzz "github.com/google/gofuzz"
	"github.com/seanchann/sample-apiserver/pkg/apis/sample"

	runtimeserializer "k8s.io/apimachinery/pkg/runtime/serializer"
)

// Funcs returns the fuzzer functions for the apps api group.
var Funcs = func(codecs runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		func(s *sample.TestSpec, c fuzz.Continue) {
			c.FuzzNoCustom(s) // fuzz self without calling this function again

		},
	}
}
