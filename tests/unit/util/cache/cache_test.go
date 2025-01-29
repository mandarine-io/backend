package cache

import (
	"github.com/mandarine-io/backend/internal/util/cache"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"testing"
)

type CacheUtilSuite struct {
	suite.Suite
}

func TestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(CacheUtilSuite))
}

func (s *CacheUtilSuite) Test_CreateCacheKey(t provider.T) {
	t.Title("Create cache key")
	t.Severity(allure.NORMAL)
	t.Tag("positive")
	t.Parallel()

	type testCase struct {
		title    string
		parts    []string
		expected string
	}

	testData := []testCase{
		{"empty parts", []string{}, ""},
		{"single part", []string{"part1"}, "part1"},
		{"multiple parts", []string{"part1", "part2", "part3"}, "part1.part2.part3"},
		{"parts with empty strings", []string{"part1", "", "part3"}, "part1..part3"},
		{"parts with dots", []string{"part1.", ".part2.", ".part3"}, "part1...part2...part3"},
	}

	for _, data := range testData {
		t.Run(
			data.title, func(t provider.T) {
				t.Severity(allure.NORMAL)
				t.Tag("positive")
				t.Parallel()

				actual := cache.CreateCacheKey(data.parts...)
				t.Require().Equal(data.expected, actual)
			},
		)
	}
}
