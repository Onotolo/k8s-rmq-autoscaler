module github.com/medal-labs/k8s-rmq-autoscaler

go 1.13

require (
	4d63.com/gochecknoglobals v0.0.0-20190118042838-abbdf6ec0afb // indirect
	4d63.com/gochecknoinits v0.0.0-20180528051558-14d5915061e5 // indirect
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/alecthomas/gocyclo v0.0.0-20150208221726-aa8f8b160214 // indirect
	github.com/alecthomas/gometalinter v3.0.0+incompatible // indirect
	github.com/alecthomas/units v0.0.0-20151022065526-2efee857e7cf // indirect
	github.com/alexflint/go-arg v1.0.0 // indirect
	github.com/alexkohler/nakedret v0.0.0-20171106223215-c0e305a4f690 // indirect
	github.com/client9/misspell v0.3.4 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/lint v0.0.0-20181217174547-8f45f776aaf1 // indirect
	github.com/google/shlex v0.0.0-20181106134648-c34317bd91bf // indirect
	github.com/gordonklaus/ineffassign v0.0.0-20180909121442-1003c8bd00dc // indirect
	github.com/gregjones/httpcache v0.0.0-20190212212710-3befbb6ad0cc // indirect
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/jgautheron/goconst v0.0.0-20170703170152-9740945f5dcb // indirect
	github.com/kisielk/errcheck v1.2.0 // indirect
	github.com/kisielk/gotool v1.0.0 // indirect
	github.com/mdempsky/maligned v0.0.0-20180708014732-6e39bd26a8c8 // indirect
	github.com/mdempsky/unconvert v0.0.0-20190117010209-2db5a8ead8e7 // indirect
	github.com/mibk/dupl v1.0.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/modocache/gover v0.0.0-20171022184752-b58185e213c5 // indirect
	github.com/namsral/flag v1.7.4-pre
	github.com/nicksnyder/go-i18n v1.10.0 // indirect
	github.com/onsi/ginkgo v1.12.3 // indirect
	github.com/opennota/check v0.0.0-20180911053232-0c771f5545ff // indirect
	github.com/pelletier/go-toml v1.2.0 // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/securego/gosec v0.0.0-20190213104759-9cdfec40ca54 // indirect
	github.com/stripe/safesql v0.0.0-20171221195208-cddf355596fe // indirect
	github.com/tsenart/deadcode v0.0.0-20160724212837-210d2dc333e9 // indirect
	github.com/walle/lll v0.0.0-20160702150637-8b13b3fbf731 // indirect
	golang.org/x/sys v0.0.0-20200602225109-6fdc65e7d980 // indirect
	gopkg.in/alecthomas/kingpin.v3-unstable v3.0.0-20180810215634-df19058c872c // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	k8s.io/api v0.0.0-20200603011159-afb0842feaf5
	k8s.io/apimachinery v0.0.0-20200601184421-76330795f827
	k8s.io/client-go v0.0.0-20200603035352-be97aaa976ad
	k8s.io/klog v0.2.0
	mvdan.cc/interfacer v0.0.0-20180901003855-c20040233aed // indirect
	mvdan.cc/lint v0.0.0-20170908181259-adc824a0674b // indirect
	mvdan.cc/unparam v0.0.0-20190213212834-da01123e7b4f // indirect
)

replace (
	golang.org/x/sys => golang.org/x/sys v0.0.0-20190813064441-fde4db37ae7a // pinned to release-branch.go1.13
	golang.org/x/tools => golang.org/x/tools v0.0.0-20190821162956-65e3620a7ae7 // pinned to release-branch.go1.13
)
