module github.com/elastic/beats/libbeat

go 1.12

replace (
	github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.4.2
	github.com/Sirupsen/logrus v1.0.5 => github.com/sirupsen/logrus v1.0.5
	github.com/Sirupsen/logrus v1.3.0 => github.com/Sirupsen/logrus v1.0.6
	github.com/Sirupsen/logrus v1.4.0 => github.com/sirupsen/logrus v1.0.6
	github.com/docker/docker => github.com/docker/engine v0.0.0-20190717161051-705d9623b7c1
	github.com/dop251/goja => github.com/andrewkroh/goja v0.0.0-20190128172624-dd2ac4456e20
	github.com/fsnotify/fsevents => github.com/elastic/fsevents v0.0.0-20181029231046-e1d381a4d270
)

require (
	github.com/Shopify/sarama v1.23.1
	github.com/armon/go-socks5 v0.0.0-20160902184237-e75332964ef5
	github.com/davecgh/go-spew v1.1.1
	github.com/dlclark/regexp2 v1.2.0 // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v0.0.0-00010101000000-000000000000
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-units v0.4.0 // indirect
	github.com/dop251/goja v0.0.0-00010101000000-000000000000
	github.com/dop251/goja_nodejs v0.0.0-20171011081505-adff31b136e6
	github.com/dustin/go-humanize v1.0.0
	github.com/elastic/ecs v1.0.1
	github.com/elastic/go-lumber v0.1.0
	github.com/elastic/go-seccomp-bpf v1.1.0
	github.com/elastic/go-structform v0.0.6
	github.com/elastic/go-sysinfo v1.0.2
	github.com/elastic/go-txfile v0.0.6
	github.com/elastic/go-ucfg v0.7.0
	github.com/elastic/gosigar v0.10.4
	github.com/ericchiang/k8s v1.2.0
	github.com/fatih/color v1.7.0
	github.com/garyburd/redigo v1.6.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-sourcemap/sourcemap v2.1.2+incompatible // indirect
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/gogo/protobuf v1.2.1 // indirect
	github.com/golang/protobuf v1.3.2 // indirect
	github.com/joeshaw/multierror v0.0.0-20140124173710-69b34d4ec901
	github.com/klauspost/compress v1.7.5 // indirect
	github.com/klauspost/cpuid v1.2.1 // indirect
	github.com/mattn/go-colorable v0.1.2
	github.com/miekg/dns v1.1.15
	github.com/mitchellh/hashstructure v1.0.0
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/pkg/errors v0.8.1
	github.com/rcrowley/go-metrics v0.0.0-20190706150252-9beb055b7962
	github.com/sirupsen/logrus v1.4.2 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.3
	github.com/stretchr/testify v1.3.0
	github.com/tsg/gopacket v0.0.0-20190320122513-dd3d0e41124a
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.10.0
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4
	golang.org/x/net v0.0.0-20190724013045-ca1201d0de80
	golang.org/x/sys v0.0.0-20190804053845-51ab0e2deafa
	golang.org/x/text v0.3.2
	google.golang.org/grpc v1.22.1 // indirect
	gopkg.in/yaml.v2 v2.2.2
)
