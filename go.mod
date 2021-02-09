module github.com/gardener/gem

go 1.15

require (
	github.com/Masterminds/semver v1.5.0
	github.com/gardener/gardener v1.16.0
	github.com/google/addlicense v0.0.0-20190510175307-22550fa7c1b0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.0.0
	gopkg.in/src-d/go-git.v4 v4.13.1
	k8s.io/apimachinery v0.19.6
	k8s.io/code-generator v0.19.6
	mvdan.cc/gofumpt v0.0.0-20190729090447-96300e3d49fb
	sigs.k8s.io/controller-tools v0.4.1
)

replace (
	k8s.io/api => k8s.io/api v0.19.6
	k8s.io/client-go => k8s.io/client-go v0.19.6
)
