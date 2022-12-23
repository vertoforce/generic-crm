module github.com/vertoforce/generic-crm

go 1.18

require (
	github.com/araddon/dateparse v0.0.0-20210429162001-6b43995a97de
	github.com/go-sql-driver/mysql v1.5.0
	github.com/jmoiron/sqlx v1.2.0
	github.com/juju/ratelimit v1.0.1
	github.com/lib/pq v1.7.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/segmentio/agecache v0.0.2
	github.com/stretchr/testify v1.7.1
	github.com/vertoforce/airtable-go v0.0.0-20200608173945-23baffada355
	github.com/vertoforce/regexgrouphelp v0.1.1
	go.elastic.co/apm v1.12.0
	go.elastic.co/apm/module/apmot v1.12.0
	golang.org/x/net v0.0.0-20200226121028-0de0cce0169b
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	gopkg.in/Iwark/spreadsheet.v2 v2.0.0-20191122095212-08231195c43b
)

require (
	cloud.google.com/go v0.34.0 // indirect
	github.com/armon/go-radix v1.0.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/elastic/go-licenser v0.3.1 // indirect
	github.com/elastic/go-sysinfo v1.1.1 // indirect
	github.com/elastic/go-windows v1.0.0 // indirect
	github.com/fabioberger/airtable-go v3.1.0+incompatible // indirect
	github.com/golang/protobuf v1.2.0 // indirect
	github.com/joeshaw/multierror v0.0.0-20140124173710-69b34d4ec901 // indirect
	github.com/makeitraina/airtable-go v3.1.0+incompatible // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pkg/errors v0.8.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/procfs v0.0.0-20190425082905-87a4384529e0 // indirect
	github.com/santhosh-tekuri/jsonschema v1.2.4 // indirect
	go.elastic.co/apm/module/apmhttp v1.12.0 // indirect
	go.elastic.co/fastjson v1.1.0 // indirect
	golang.org/x/lint v0.0.0-20201208152925-83fdc39ff7b5 // indirect
	golang.org/x/mod v0.2.0 // indirect
	golang.org/x/sys v0.0.0-20191204072324-ce4227a45e2e // indirect
	golang.org/x/tools v0.0.0-20200509030707-2212a7e161a5 // indirect
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543 // indirect
	google.golang.org/appengine v1.4.0 // indirect
	gopkg.in/check.v1 v1.0.0-20200902074654-038fdea0a05b // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
	howett.net/plist v0.0.0-20181124034731-591f970eefbb // indirect
)

replace github.com/segmentio/agecache v0.0.2 => github.com/deankarn/agecache v0.0.0-20220327155826-92fabfdefbd2
