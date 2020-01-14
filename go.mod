module github.com/IanCassTwo/akamai-property-to-terraform

go 1.12

replace github.com/h2non/gock => gopkg.in/h2non/gock.v1 v1.0.14

replace github.com/akamai/cli-common-golang => github.com/IanCassTwo/cli-common-golang v0.0.0-20200114191904-0fe43bdd18ff

replace github.com/janeczku/go-spinner => github.com/dmolesUC/go-spinner v0.0.0-20190903171623-0c332afb0926

require (
	github.com/akamai/AkamaiOPEN-edgegrid-golang v0.9.5
	github.com/akamai/cli-common-golang v0.0.0-20190304161655-ff6f30c42c80
	github.com/fatih/color v1.9.0
	github.com/janeczku/go-spinner v0.0.0-20150530144529-cf8ef1d64394
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/stretchr/testify v1.4.0 // indirect
	github.com/urfave/cli v1.22.2
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	golang.org/x/sys v0.0.0-20200107162124-548cf772de50 // indirect
	gopkg.in/ini.v1 v1.51.1 // indirect
)
