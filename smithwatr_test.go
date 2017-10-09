package smithwatr_test

import (
	"errors"
	"log"

	. "github.com/dimus/smithwatr"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Smithwatr", func() {
	Describe("Check()", func() {
		It("ignores `nil` errors", func() {
			err := error(nil)
			a := "one"
			Check(err)
			Expect(a).To(Equal("one"))
		})

		It("panics if err is not `nil`", func() {
			defer func() {
				if r := recover(); r != nil {
					e := r.(error)
					Expect(e).To(Equal(errors.New("My error")))
				}
			}()
			err := errors.New("My error")
			Check(err)
		})
	})

	Describe("InitBlosum62()", func() {
		It("returns BLOSUM62 weights", func() {
			b := InitBlosum62()
			Expect(b['A']['N']).To(Equal(-2))
			Expect(b['*']['N']).To(Equal(-4))
			Expect(b['W']['W']).To(Equal(11))
		})
	})

	Describe("EnvVars()", func() {
		It("reads data from environment", func() {
			env := EnvVars()
			Expect(env.DbHost).To(Equal("pg"))
			Expect(env.DbUser).To(Equal("postgres"))
			Expect(env.Db).To(Equal("smithwatr"))
			Expect(env.GapOpens).To(Equal(10))
			Expect(env.GapExtends).To(Equal(1))
		})
	})

	Describe("SmithWaterman()", func() {
		It("calculates identical alignment", func() {
			g1 := Gene{Seq: []rune("AA"), SeqLen: 2, Gene: "gene1"}
			g2 := Gene{Seq: []rune("AA"), SeqLen: 2, Gene: "gene2"}

			res := SmithWaterman(g1, g2, b62, conf)
			Expect(res.Score).To(Equal(8))
			Expect(res.Identical).To(Equal(2))
			Expect(res.Similar).To(Equal(0))
			i, s := res.IdentitySimilarity()
			Expect(i).To(Equal(float32(100)))
			Expect(s).To(Equal(float32(100)))
			log.Println(res.Show(50))
		})

		It("calculates 'real' alignment", func() {
			s1 := []rune("MADRGFCSADGSDPLWDWNVTWNTSNPDFTKCF")
			g1 := Gene{Seq: s1, SeqLen: len(s1), Gene: "gene1"}
			s2 := []rune("MANRGFCSADGWPLWDWDVTWNTSNPDFTKCF")
			g2 := Gene{Seq: s2, SeqLen: len(s2), Gene: "gene2"}
			res := SmithWaterman(g1, g2, b62, conf)
			identity, similarity := res.IdentitySimilarity()
			Expect(res.Score).To(Equal(177))
			Expect(identity).To(BeNumerically("~", 87.8, 0.1))
			Expect(similarity).To(BeNumerically("~", 93.9, 0.1))
			log.Println(res.Show(50))
		})
	})

	Describe("ImportData()", func() {
		It("imports data to the database", func() {
			ImportData(db, conf)
			Expect(NotEmpty(db, "genes")).To(Equal(true))
		})
	})

	Describe("ImportJobs()", func() {
		It("imports jobs to the database", func() {
			ImportData(db, conf)
			ImportJobs(db, 1)
			Expect(NotEmpty(db, "jobs")).To(Equal(true))
		})
	})

	Describe("Align()", func() {
		It("Aligns genes and saves data", func() {
			ImportData(db, conf)
			ImportJobs(db, 1)
			Align(db, 2, -1, b62, conf)
		})
	})
})
