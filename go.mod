module github.com/elastic/beats

go 1.12

require (
	cloud.google.com/go v0.38.0 // indirect
	github.com/DataDog/zstd v0.0.0-20160706220725-2bf71ec48360 // indirect
	github.com/Microsoft/go-winio v0.4.2 // indirect
	github.com/OneOfOne/xxhash v1.2.3
	github.com/Shopify/sarama v1.20.1
	github.com/Shopify/toxiproxy v2.1.4+incompatible // indirect
	github.com/StackExchange/wmi v0.0.0-20180116203802-5d049714c4a6
	github.com/aerospike/aerospike-client-go v0.0.0-20170612174108-0f3b54da6bdc
	github.com/andrewkroh/sys v0.0.0-20151128191922-287798fe3e43
	github.com/aws/aws-lambda-go v1.6.0
	github.com/aws/aws-sdk-go-v2 v0.5.0
	github.com/awslabs/goformation v0.0.0-20180916202949-d42502ef32a8
	github.com/blakesmith/ar v0.0.0-20190502131153-809d4375e1fb
	github.com/bsm/sarama-cluster v0.0.0-20180625083203-7e67d87a6b3f
	github.com/cavaliercoder/badio v0.0.0-20160213150051-ce5280129e9e // indirect
	github.com/cavaliercoder/go-rpm v0.0.0-20190131055624-7a9c54e3d83e
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/coreos/bbolt v1.3.2
	github.com/coreos/go-systemd v0.0.0-20190618135430-ff7011eec365
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f
	github.com/denisenkom/go-mssqldb v0.0.0-20181014144952-4e0d7dc8888f
	github.com/digitalocean/go-libvirt v0.0.0-20190715144809-7b622097a793
	github.com/dlclark/regexp2 v0.0.0-20171009020623-7632a260cbaf // indirect
	github.com/docker/distribution v0.0.0-20170524205824-1e2f10eb6574 // indirect
	github.com/docker/docker v0.0.0-20170802015333-8af4db6f002a
	github.com/docker/go-connections v0.3.0
	github.com/docker/go-units v0.3.2 // indirect
	github.com/docker/libtrust v0.0.0-20160708172513-aabc10ec26b7 // indirect
	github.com/dop251/goja v0.0.0-00010101000000-000000000000
	github.com/dop251/goja_nodejs v0.0.0-20171011081505-adff31b136e6
	github.com/dustin/go-humanize v1.0.0
	github.com/eapache/go-resiliency v0.0.0-20160104191539-b86b1ec0dd42 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20160609142408-bb955e01b934 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/elastic/ecs v1.0.1
	github.com/elastic/go-libaudit v0.4.0
	github.com/elastic/go-lookslike v0.2.0
	github.com/elastic/go-sysinfo v0.0.0-20190508093345-9a4be54a53be
	github.com/elastic/go-ucfg v0.7.0
	github.com/elastic/gosigar v0.10.3
	github.com/ericchiang/k8s v1.0.0
	github.com/fsnotify/fsevents v0.0.0-00010101000000-000000000000
	github.com/fsnotify/fsnotify v1.4.7
	github.com/garyburd/redigo v0.0.0-20160525165706-b8dc90050f24
	github.com/go-ole/go-ole v1.2.1 // indirect
	github.com/go-sourcemap/sourcemap v2.1.2+incompatible // indirect
	github.com/go-sql-driver/mysql v1.4.1
	github.com/gocarina/gocsv v0.0.0-20170324095351-ffef3ffc77be
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/gogo/protobuf v1.2.1 // indirect
	github.com/golang/protobuf v1.3.1
	github.com/golang/snappy v0.0.0-20170215233205-553a64147049
	github.com/google/flatbuffers v0.0.0-20170925184458-7a6b2bf521e9
	github.com/gorhill/cronexpr v0.0.0-20161205141322-d520615e531a
	github.com/insomniacslk/dhcp v0.0.0-20180716145214-633285ba52b2
	github.com/jmespath/go-jmespath v0.0.0-20180206201540-c2b33e8439af // indirect
	github.com/joeshaw/multierror v0.0.0-20140124173710-69b34d4ec901
	github.com/jstemmer/go-junit-report v0.0.0-20190106144839-af01ea7f8024
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/lib/pq v1.1.1
	github.com/magefile/mage v1.8.0
	github.com/miekg/dns v1.0.8
	github.com/mitchellh/hashstructure v0.0.0-20170116052023-ab25296c0f51
	github.com/mitchellh/mapstructure v1.1.2
	github.com/opencontainers/go-digest v0.0.0-20170510163354-eaa60544f31c // indirect
	github.com/opencontainers/image-spec v0.0.0-20170525204040-4038d4391fe9 // indirect
	github.com/pierrec/lz4 v0.0.0-20170226142621-90290f74b1b4 // indirect
	github.com/pierrec/xxHash v0.0.0-20160112165351-5a004441f897 // indirect
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_model v0.0.0-20190129233127-fd36f4220a90
	github.com/prometheus/common v0.6.0
	github.com/prometheus/procfs v0.0.2
	github.com/rcrowley/go-metrics v0.0.0-20181016184325-3113b8401b8a // indirect
	github.com/rogpeppe/gohack v1.0.2 // indirect
	github.com/samuel/go-parser v0.0.0-20130731160455-ca8abbf65d0e // indirect
	github.com/samuel/go-thrift v0.0.0-20140522043831-2187045faa54
	github.com/shirou/gopsutil v2.18.11+incompatible
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.3
	github.com/stretchr/testify v1.3.0
	github.com/tsg/gopacket v0.0.0-20190320122513-dd3d0e41124a
	github.com/vmware/govmomi v0.20.2
	github.com/yuin/gopher-lua v0.0.0-20170403160031-b402f3114ec7 // indirect
	go.etcd.io/bbolt v1.3.3 // indirect
	go.uber.org/atomic v1.4.0
	go.uber.org/multierr v1.1.0
	go.uber.org/zap v1.10.0
	golang.org/x/crypto v0.0.0-20190618222545-ea8f1a30c443
	golang.org/x/net v0.0.0-20190619014844-b5b0513f8c1b
	golang.org/x/sys v0.0.0-20190616124812-15dcb6c0061f
	golang.org/x/text v0.3.2
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4
	golang.org/x/tools v0.0.0-20190619202714-22e91af008f2
	google.golang.org/appengine v1.5.0 // indirect
	gopkg.in/inf.v0 v0.9.0
	gopkg.in/mgo.v2 v2.0.0-20160818020120-3f83fa500528
	gopkg.in/yaml.v2 v2.2.2
	howett.net/plist v0.0.0-20181124034731-591f970eefbb
	k8s.io/apimachinery v0.0.0-20190719140911-bfcf53abc9f8
)

replace github.com/dop251/goja => github.com/andrewkroh/goja v0.0.0-20190128172624-dd2ac4456e20

replace github.com/fsnotify/fsevents => github.com/elastic/fsevents v0.0.0-20181029231046-e1d381a4d270
