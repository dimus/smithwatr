package smithwatr_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSmithwatr(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Smithwatr Suite")
}
