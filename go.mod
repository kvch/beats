module github.com/elastic/beats

go 1.13

replace (
	github.com/Shopify/sarama => github.com/elastic/sarama v0.0.0-20191122160421-355d120d0970
	github.com/docker/docker => github.com/docker/engine v0.0.0-20191113042239-ea84732a7725
	github.com/dop251/goja => github.com/andrewkroh/goja v0.0.0-20190128172624-dd2ac4456e20
	github.com/fsnotify/fsevents => github.com/elastic/fsevents v0.0.0-20181029231046-e1d381a4d270
	github.com/fsnotify/fsnotify => github.com/adriansr/fsnotify v0.0.0-20180417234312-c9bbe1f46f1d
	github.com/tonistiigi/fifo => github.com/containerd/fifo v0.0.0-20190816180239-bda0ff6ed73c
)

require (
	4d63.com/tz v1.1.0
	cloud.google.com/go v0.50.0
	cloud.google.com/go/pubsub v1.1.0
	cloud.google.com/go/storage v1.5.0
	github.com/Azure/azure-event-hubs-go/v3 v3.1.0
	github.com/Azure/azure-sdk-for-go v37.1.0+incompatible
	github.com/Azure/azure-storage-blob-go v0.8.0
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/Azure/go-autorest/autorest v0.9.4
	github.com/Azure/go-autorest/autorest/azure/auth v0.4.2
	github.com/Azure/go-autorest/autorest/date v0.2.0
	github.com/Microsoft/go-winio v0.4.15-0.20190919025122-fc70bd9a86b5
	github.com/Microsoft/hcsshim v0.8.6 // indirect
	github.com/Shopify/sarama v0.0.0-00010101000000-000000000000
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d
	github.com/aerospike/aerospike-client-go v0.0.0-20170612174108-0f3b54da6bdc
	github.com/andrewkroh/sys v0.0.0-20151128191922-287798fe3e43
	github.com/armon/go-socks5 v0.0.0-20160902184237-e75332964ef5
	github.com/aws/aws-lambda-go v1.12.0
	github.com/aws/aws-sdk-go-v2 v0.9.0
	github.com/awslabs/goformation/v4 v4.1.0
	github.com/blakesmith/ar v0.0.0-20190502131153-809d4375e1fb
	github.com/bsm/sarama-cluster v2.1.15+incompatible
	github.com/cavaliercoder/badio v0.0.0-20160213150051-ce5280129e9e // indirect
	github.com/cavaliercoder/go-rpm v0.0.0-20190131055624-7a9c54e3d83e
	github.com/cespare/xxhash/v2 v2.0.0-20191031200418-2372543dd2bb
	github.com/containerd/containerd v1.3.2 // indirect
	github.com/containerd/continuity v0.0.0-20200107194136-26c1120b8d41 // indirect
	github.com/containerd/fifo v0.0.0-20191213151349-ff969a566b00
	github.com/coreos/bbolt v1.3.3
	github.com/coreos/go-systemd v0.0.0-20190321100706-95778dfbb74e // indirect
	github.com/coreos/go-systemd/v22 v22.0.0
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f
	github.com/denisenkom/go-mssqldb v0.0.0-20191128021309-1d7a30a10f73
	github.com/digitalocean/go-libvirt v0.0.0-20190715144809-7b622097a793
	github.com/dlclark/regexp2 v1.2.0 // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v0.0.0-00010101000000-000000000000
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-metrics v0.0.1 // indirect
	github.com/docker/go-plugins-helpers v0.0.0-20200102110956-c9a8a2d92ccc
	github.com/docker/go-units v0.4.0
	github.com/dop251/goja v0.0.0-00010101000000-000000000000
	github.com/dop251/goja_nodejs v0.0.0-20171011081505-adff31b136e6
	github.com/dustin/go-humanize v1.0.0
	github.com/elastic/ecs v1.4.0
	github.com/elastic/go-libaudit v0.4.0
	github.com/elastic/go-lookslike v0.3.0
	github.com/elastic/go-lumber v0.1.0
	github.com/elastic/go-perf v0.0.0-20191212140718-9c656876f595
	github.com/elastic/go-seccomp-bpf v1.1.0
	github.com/elastic/go-structform v0.0.6
	github.com/elastic/go-sysinfo v1.3.0
	github.com/elastic/go-txfile v0.0.7
	github.com/elastic/go-ucfg v0.7.0
	github.com/elastic/gosigar v0.10.5
	github.com/fatih/color v1.9.0
	github.com/fsnotify/fsevents v0.0.0-00010101000000-000000000000
	github.com/fsnotify/fsnotify v1.4.7
	github.com/garyburd/redigo v1.6.0
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/go-sourcemap/sourcemap v2.1.2+incompatible // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gocarina/gocsv v0.0.0-20191214001331-e6697589f2e0
	github.com/godror/godror v0.10.4
	github.com/gofrs/flock v0.7.1
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.3.2
	github.com/golang/snappy v0.0.1
	github.com/google/flatbuffers v1.11.0
	github.com/google/gopacket v1.1.17
	github.com/gorhill/cronexpr v0.0.0-20180427100037-88b0669f7d75
	github.com/gorilla/mux v1.7.3 // indirect
	github.com/insomniacslk/dhcp v0.0.0-20180716145214-633285ba52b2
	github.com/jmoiron/sqlx v1.2.0
	github.com/joeshaw/multierror v0.0.0-20140124173710-69b34d4ec901
	github.com/jstemmer/go-junit-report v0.9.1
	github.com/klauspost/cpuid v1.2.2 // indirect
	github.com/lib/pq v1.3.0
	github.com/magefile/mage v1.9.0
	github.com/mattn/go-colorable v0.1.4
	github.com/miekg/dns v1.1.27
	github.com/mitchellh/gox v1.0.1
	github.com/mitchellh/hashstructure v1.0.0
	github.com/mitchellh/mapstructure v1.1.2
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/opencontainers/runc v0.1.1 // indirect
	github.com/pierrre/gotestcover v0.0.0-20160517101806-924dca7d15f0
	github.com/pkg/errors v0.9.1
	github.com/pmezard/go-difflib v1.0.0
	github.com/prometheus/client_model v0.2.0
	github.com/prometheus/common v0.6.0
	github.com/prometheus/procfs v0.0.9-0.20191208103036-42f6e295b56f
	github.com/rcrowley/go-metrics v0.0.0-20190826022208-cac0b30c2563
	github.com/samuel/go-parser v0.0.0-20170131185712-99744db8e45c // indirect
	github.com/samuel/go-thrift v0.0.0-20140522043831-2187045faa54
	github.com/shirou/gopsutil v2.19.12+incompatible
	github.com/sirupsen/logrus v1.4.2 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.4.0
	github.com/tsg/go-daemon v0.0.0-20200123164349-4b60efc26d5f
	github.com/tsg/gopacket v0.0.0-20190320122513-dd3d0e41124a
	github.com/vmware/govmomi v0.0.0-20170802214208-2cad15190b41
	github.com/yuin/gopher-lua v0.0.0-20191220021717-ab39c6098bdb // indirect
	go.etcd.io/bbolt v1.3.3 // indirect
	go.uber.org/atomic v1.5.1
	go.uber.org/multierr v1.4.0
	go.uber.org/zap v1.13.0
	golang.org/x/crypto v0.0.0-20200117160349-530e935923ad
	golang.org/x/net v0.0.0-20200114155413-6afb5195e5aa
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	golang.org/x/sys v0.0.0-20200117145432-59e60aa80a0c
	golang.org/x/text v0.3.2
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	golang.org/x/tools v0.0.0-20200119215504-eb0d8dd85bcc
	google.golang.org/api v0.15.0
	google.golang.org/genproto v0.0.0-20200117163144-32f20d992d24
	gopkg.in/inf.v0 v0.9.1
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
	gopkg.in/yaml.v2 v2.2.7
	gotest.tools v2.2.0+incompatible // indirect
	howett.net/plist v0.0.0-20181124034731-591f970eefbb
	k8s.io/api v0.17.1
	k8s.io/apimachinery v0.17.1
	k8s.io/client-go v0.17.1
)
