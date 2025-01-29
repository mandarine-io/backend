package security

import (
	"github.com/mandarine-io/backend/internal/util/security"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"testing"
)

type SecurityUtilSuite struct {
	suite.Suite
}

func TestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(SecurityUtilSuite))
}

func (s *SecurityUtilSuite) Test_HashPassword(t provider.T) {
	t.Title("Test HashPassword function")
	t.Severity(allure.NORMAL)
	t.Tags("security-utils", "positive")
	t.Parallel()

	testCases := []struct {
		title        string
		password     string
		expectsError bool
	}{
		{"hash password successfully", "strongpassword123", false},
	}

	for _, tc := range testCases {
		t.Run(
			tc.title, func(t provider.T) {
				hash, err := security.HashPassword(tc.password)
				if tc.expectsError {
					t.Require().Error(err)
					t.Require().Empty(hash)
				} else {
					t.Require().NoError(err)
					t.Require().NotEmpty(hash)
				}
			},
		)
	}
}

func (s *SecurityUtilSuite) Test_CheckPasswordHash_Positive(t provider.T) {
	t.Title("Test CheckPasswordHash functions")
	t.Severity(allure.NORMAL)
	t.Tags("security-utils", "positive")
	t.Parallel()

	password := "strongpassword123"
	hash, err := security.HashPassword(password)
	t.Require().NoError(err)

	match := security.CheckPasswordHash(password, hash)
	t.Require().True(match)
}

func (s *SecurityUtilSuite) Test_CheckPasswordHash_Negative(t provider.T) {
	t.Title("Test HashPassword and CheckPasswordHash functions")
	t.Severity(allure.NORMAL)
	t.Tags("security-utils", "negative")
	t.Parallel()

	password := "strongpassword123"
	hash, err := security.HashPassword(password)
	t.Require().NoError(err)

	wrongPassword := "wrongpassword"
	match := security.CheckPasswordHash(wrongPassword, hash)
	t.Require().False(match)
}
