module github.com/seanchann/sample-apiserver

require (
	cloud.google.com/go v0.34.0 // indirect
	contrib.go.opencensus.io/exporter/ocagent v0.3.0 // indirect
	github.com/Azure/go-autorest v10.6.2+incompatible // indirect
	github.com/CyCoreSystems/ari v4.8.3+incompatible
	github.com/DataDog/datadog-go v0.0.0-20180822151419-281ae9f2d895 // indirect
	github.com/NYTimes/gziphandler v1.0.1 // indirect
	github.com/armon/go-metrics v0.0.0-20180917152333-f0300d1749da // indirect`
	github.com/circonus-labs/circonus-gometrics v2.2.5+incompatible // indirect
	github.com/circonus-labs/circonusllhist v0.1.3 // indirect
	github.com/cockroachdb/cmux v0.0.0-20170110192607-30d10be49292 // indirect
	github.com/coreos/go-oidc v2.0.0+incompatible // indirect
	github.com/coreos/go-semver v0.2.0 // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/go-openapi/spec v0.17.2
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/groupcache v0.0.0-20181024230925-c65c006176ff // indirect
	github.com/google/gofuzz v0.0.0-20170612174753-24818f796faf
	github.com/gophercloud/gophercloud v0.0.0-20181207024449-3038305ba4ed // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.5.1 // indirect
	github.com/hashicorp/consul v1.4.0
	github.com/hashicorp/go-immutable-radix v1.0.0 // indirect
	github.com/hashicorp/go-msgpack v0.0.0-20150518234257-fa3f63826f7c // indirect
	github.com/hashicorp/go-multierror v1.0.0 // indirect
	github.com/hashicorp/go-retryablehttp v0.5.0 // indirect
	github.com/hashicorp/go-rootcerts v0.0.0-20160503143440-6bb64b370b90 // indirect
	github.com/hashicorp/go-sockaddr v0.0.0-20180320115054-6d291a969b86 // indirect
	github.com/hashicorp/memberlist v0.1.0 // indirect
	github.com/hashicorp/serf v0.8.1 // indirect
	github.com/hashicorp/yamux v0.0.0-20181012175058-2f1d1f20f75d // indirect
	github.com/inconshreveable/log15 v0.0.0-20180818164646-67afb5ed74ec
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jonboulle/clockwork v0.1.0 // indirect
	github.com/mattn/go-colorable v0.0.9 // indirect
	github.com/mattn/go-isatty v0.0.4 // indirect
	github.com/miekg/dns v1.1.1 // indirect
	github.com/mitchellh/go-homedir v1.0.0 // indirect
	github.com/mitchellh/go-testing-interface v1.0.0 // indirect
	github.com/mitchellh/mapstructure v1.1.2 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/pascaldekloe/goe v0.0.0-20180627143212-57f6aae5913c // indirect
	github.com/pkg/errors v0.8.0 // indirect
	github.com/pquerna/cachecontrol v0.0.0-20180517163645-1555304b9b35 // indirect
	github.com/sean-/seed v0.0.0-20170313163322-e2103e2c3529 // indirect
	github.com/seanchann/apimaster v0.0.0-20181125155804-5b2322e9cd5f
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.3
	github.com/tv42/httpunix v0.0.0-20150427012821-b75d8614f926 // indirect
	go.opencensus.io v0.18.0 // indirect
	golang.org/x/crypto v0.0.0-20181203042331-505ab145d0a9
	golang.org/x/time v0.0.0-20181108054448-85acf8d2951c // indirect
	golang.org/x/tools v0.0.0-20181207222222-4c874b978acb // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/square/go-jose.v2 v2.2.1 // indirect
	gopkg.in/vmihailenco/msgpack.v2 v2.9.1 // indirect
	k8s.io/apimachinery v0.0.0-20181108192626-90473842928c
	k8s.io/apiserver v0.0.0-20181206230536-f9c2597c8687
	k8s.io/client-go v9.0.0+incompatible
	k8s.io/code-generator v0.0.0-20181206115026-3a2206dd6a78 // indirect
	k8s.io/gengo v0.0.0-20181113154421-fd15ee9cc2f7 // indirect
	k8s.io/klog v0.1.0 // indirect
	k8s.io/kube-openapi v0.0.0-20181114233023-0317810137be
)

replace (
	//if you need to compiler on local configure your custom gopath
	github.com/Sirupsen/logrus => ../../../../../pkg/mod/github.com/sirupsen/logrus@v1.2.0
	k8s.io/apiserver => github.com/seanchann/apiserver v0.0.0-20181125155507-2d612e0f9460
)
