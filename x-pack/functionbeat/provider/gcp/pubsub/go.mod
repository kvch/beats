module github.com/elastic/beats/x-pack/functionbeat/provider/gcp/pubsub

replace github.com/elastic/beats => github.com/kvch/beats v0.0.0-20190809112705-79b88d005ef565ce363b83f6f00c97d3f4f3cc13

replace (
	github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.4.2
	github.com/docker/docker => github.com/docker/engine v0.0.0-20190717161051-705d9623b7c1
)

require (
	cloud.google.com/go v0.43.0 // indirect
	github.com/Shopify/sarama v1.23.1 // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v0.0.0-00010101000000-000000000000 // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/elastic/beats v0.0.0-00010101000000-000000000000 // indirect
	github.com/elastic/ecs v1.0.1 // indirect
	github.com/elastic/go-lumber v0.1.0 // indirect
	github.com/elastic/go-seccomp-bpf v1.1.0 // indirect
	github.com/elastic/go-structform v0.0.6 // indirect
	github.com/elastic/go-sysinfo v1.0.2 // indirect
	github.com/elastic/go-txfile v0.0.6 // indirect
	github.com/elastic/gosigar v0.10.4 // indirect
	github.com/fatih/color v1.7.0 // indirect
	github.com/garyburd/redigo v1.6.0 // indirect
	github.com/gofrs/uuid v3.2.0+incompatible // indirect
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/klauspost/compress v1.7.5 // indirect
	github.com/klauspost/cpuid v1.2.1 // indirect
	github.com/mattn/go-colorable v0.1.2 // indirect
	github.com/miekg/dns v1.1.15 // indirect
	github.com/mitchellh/hashstructure v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20190706150252-9beb055b7962 // indirect
	github.com/sirupsen/logrus v1.4.2 // indirect
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.10.0 // indirect
	k8s.io/api v0.0.0-20190808180749-077ce48e77da // indirect
	k8s.io/apimachinery v0.0.0-20190809020650-423f5d784010 // indirect
	k8s.io/client-go v11.0.0+incompatible // indirect
	k8s.io/utils v0.0.0-20190809000727-6c36bc71fc4a // indirect
)
