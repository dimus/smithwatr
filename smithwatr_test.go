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
			seq1 := []rune("AA")
			seq2 := []rune("AA")
			res := SmithWaterman(seq1, seq2, b62, conf)
			Expect(res.Score).To(Equal(8))
			Expect(res.Identical).To(Equal(2))
			Expect(res.Similar).To(Equal(0))
			i, s := res.IdentitySimilarity()
			Expect(i).To(Equal(float32(100)))
			Expect(s).To(Equal(float32(100)))
		})

		It("calculates 'real' alignment", func() {
			s1 := []rune("MALRGFCSADGSDPLWDWNVTWNTSNPDFTKCF")
			s2 := []rune("MALRGFCSADGAPLWDWDVTWNTSNPDFTKCF")
			res := SmithWaterman(s1, s2, b62, conf)
			identity, similarity := res.IdentitySimilarity()
			log.Println(res)
			log.Println(res.Score, identity, similarity)

		})
	})
})
