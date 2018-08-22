/*

Copyright 2018 This Project Authors.

Author:  seanchann <seanchann@foxmail.com>

See docs/ for more information about the  project.

*/

package main

import (
	"flag"

	"github.com/seanchann/apimaster/pkg/swaggerdoc"
	swaggerdocopt "github.com/seanchann/apimaster/pkg/swaggerdoc/options"
	"github.com/spf13/pflag"
)

func main() {
	opt := swaggerdocopt.NewSwaggerDocOptions()

	opt.AddFlags(pflag.CommandLine)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	swaggerdoc.GenerateDoc(opt)
}
