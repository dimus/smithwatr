package smithwatr_test

import (
	. "github.com/dimus/smithwatr"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var b62 Blosum62
var conf Env

func TestSmithwatr(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Smithwatr Suite")
}

var _ = BeforeSuite(func() {
	b62 = InitBlosum62()
	conf = EnvVars()
})
