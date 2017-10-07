package smithwatr_test

import (
	"database/sql"

	. "github.com/dimus/smithwatr"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var b62 Blosum62
var conf Env
var db *sql.DB

func TestSmithwatr(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Smithwatr Suite")
}

var _ = BeforeSuite(func() {
	var err error
	b62 = InitBlosum62()
	conf = EnvVars()
	db, err = Connect(conf)
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	var err error
	err = db.Close()
	Expect(err).NotTo(HaveOccurred())
})
