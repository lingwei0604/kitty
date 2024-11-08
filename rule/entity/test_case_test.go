package entity

import (
	"testing"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/lingwei0604/kitty/rule/dto"
	"github.com/stretchr/testify/assert"
)

func TestTestCase_Asserts(t *testing.T) {
	compiledRule := loadRule()

	for _, c := range []struct {
		name    string
		Given   string
		Expect  string
		Asserts func(t *testing.T, err error)
	}{
		{
			"simple pass",
			"http://baidu.com?channel=foo",
			"i == 1",
			func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			"simple fail",
			"http://baidu.com?channel=foo",
			"i == 2",
			func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			"simple error",
			"http://baidu.com?channel=foo",
			"i",
			func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	} {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			testCase := TestCase{
				Given: Given{
					URL: c.Given,
				},
				Expect: c.Expect,
			}
			err := testCase.Asserts(compiledRule, dto.NewDecoder())
			c.Asserts(t, err)
		})
	}
}

func loadRule() *AdvancedRuleCollection {
	ar := NewAdvancedRule()
	k := koanf.New(".")
	_ = k.Load(rawbytes.Provider([]byte(`
style: basic
rule:
  - if: Channel == "foo"
    then:
      i: 1
  - if: true
    then:
      i: 2
`)), yaml.Parser())
	_ = ar.Unmarshal(k)
	_ = ar.Compile()
	return ar
}
