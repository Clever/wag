package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Clever/wag/samples/gen-go-db-custom-path/models/v9"
	"github.com/Clever/wag/samples/v9/gen-go-db-custom-path/db"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/require"
)

func mustTime(s string) strfmt.DateTime {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return strfmt.DateTime(t)
}

func mustDate(s string) strfmt.Date {
	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		panic(err)
	}
	return strfmt.Date(t)
}

func pointerToString(str string) *string {
	return &str
}

func RunDBTests(t *testing.T, dbFactory func() db.Interface) {
	t.Run("GetDeployment", GetDeployment(dbFactory(), t))
	t.Run("ScanDeployments", ScanDeployments(dbFactory(), t))
	t.Run("GetDeploymentsByEnvAppAndVersion", GetDeploymentsByEnvAppAndVersion(dbFactory(), t))
	t.Run("SaveDeployment", SaveDeployment(dbFactory(), t))
	t.Run("DeleteDeployment", DeleteDeployment(dbFactory(), t))
	t.Run("GetDeploymentsByEnvAppAndDate", GetDeploymentsByEnvAppAndDate(dbFactory(), t))
	t.Run("ScanDeploymentsByEnvAppAndDate", ScanDeploymentsByEnvAppAndDate(dbFactory(), t))
	t.Run("GetDeploymentsByEnvironmentAndDate", GetDeploymentsByEnvironmentAndDate(dbFactory(), t))
	t.Run("GetDeploymentByVersion", GetDeploymentByVersion(dbFactory(), t))
	t.Run("ScanDeploymentsByVersion", ScanDeploymentsByVersion(dbFactory(), t))
	t.Run("GetEvent", GetEvent(dbFactory(), t))
	t.Run("ScanEvents", ScanEvents(dbFactory(), t))
	t.Run("GetEventsByPkAndSk", GetEventsByPkAndSk(dbFactory(), t))
	t.Run("SaveEvent", SaveEvent(dbFactory(), t))
	t.Run("DeleteEvent", DeleteEvent(dbFactory(), t))
	t.Run("GetEventsBySkAndData", GetEventsBySkAndData(dbFactory(), t))
	t.Run("ScanEventsBySkAndData", ScanEventsBySkAndData(dbFactory(), t))
	t.Run("GetNoRangeThingWithCompositeAttributes", GetNoRangeThingWithCompositeAttributes(dbFactory(), t))
	t.Run("ScanNoRangeThingWithCompositeAttributess", ScanNoRangeThingWithCompositeAttributess(dbFactory(), t))
	t.Run("SaveNoRangeThingWithCompositeAttributes", SaveNoRangeThingWithCompositeAttributes(dbFactory(), t))
	t.Run("DeleteNoRangeThingWithCompositeAttributes", DeleteNoRangeThingWithCompositeAttributes(dbFactory(), t))
	t.Run("GetNoRangeThingWithCompositeAttributessByNameVersionAndDate", GetNoRangeThingWithCompositeAttributessByNameVersionAndDate(dbFactory(), t))
	t.Run("ScanNoRangeThingWithCompositeAttributessByNameVersionAndDate", ScanNoRangeThingWithCompositeAttributessByNameVersionAndDate(dbFactory(), t))
	t.Run("GetNoRangeThingWithCompositeAttributesByNameBranchCommit", GetNoRangeThingWithCompositeAttributesByNameBranchCommit(dbFactory(), t))
	t.Run("GetSimpleThing", GetSimpleThing(dbFactory(), t))
	t.Run("ScanSimpleThings", ScanSimpleThings(dbFactory(), t))
	t.Run("SaveSimpleThing", SaveSimpleThing(dbFactory(), t))
	t.Run("DeleteSimpleThing", DeleteSimpleThing(dbFactory(), t))
	t.Run("GetTeacherSharingRule", GetTeacherSharingRule(dbFactory(), t))
	t.Run("ScanTeacherSharingRules", ScanTeacherSharingRules(dbFactory(), t))
	t.Run("GetTeacherSharingRulesByTeacherAndSchoolApp", GetTeacherSharingRulesByTeacherAndSchoolApp(dbFactory(), t))
	t.Run("SaveTeacherSharingRule", SaveTeacherSharingRule(dbFactory(), t))
	t.Run("DeleteTeacherSharingRule", DeleteTeacherSharingRule(dbFactory(), t))
	t.Run("GetTeacherSharingRulesByDistrictAndSchoolTeacherApp", GetTeacherSharingRulesByDistrictAndSchoolTeacherApp(dbFactory(), t))
	t.Run("ScanTeacherSharingRulesByDistrictAndSchoolTeacherApp", ScanTeacherSharingRulesByDistrictAndSchoolTeacherApp(dbFactory(), t))
	t.Run("GetThing", GetThing(dbFactory(), t))
	t.Run("ScanThings", ScanThings(dbFactory(), t))
	t.Run("GetThingsByNameAndVersion", GetThingsByNameAndVersion(dbFactory(), t))
	t.Run("SaveThing", SaveThing(dbFactory(), t))
	t.Run("DeleteThing", DeleteThing(dbFactory(), t))
	t.Run("GetThingByID", GetThingByID(dbFactory(), t))
	t.Run("ScanThingsByID", ScanThingsByID(dbFactory(), t))
	t.Run("GetThingsByNameAndCreatedAt", GetThingsByNameAndCreatedAt(dbFactory(), t))
	t.Run("ScanThingsByNameAndCreatedAt", ScanThingsByNameAndCreatedAt(dbFactory(), t))
	t.Run("GetThingsByNameAndRangeNullable", GetThingsByNameAndRangeNullable(dbFactory(), t))
	t.Run("ScanThingsByNameAndRangeNullable", ScanThingsByNameAndRangeNullable(dbFactory(), t))
	t.Run("GetThingsByHashNullableAndName", GetThingsByHashNullableAndName(dbFactory(), t))
	t.Run("GetThingAllowingBatchWrites", GetThingAllowingBatchWrites(dbFactory(), t))
	t.Run("ScanThingAllowingBatchWritess", ScanThingAllowingBatchWritess(dbFactory(), t))
	t.Run("GetThingAllowingBatchWritessByNameAndVersion", GetThingAllowingBatchWritessByNameAndVersion(dbFactory(), t))
	t.Run("SaveThingAllowingBatchWrites", SaveThingAllowingBatchWrites(dbFactory(), t))
	t.Run("DeleteThingAllowingBatchWrites", DeleteThingAllowingBatchWrites(dbFactory(), t))
	t.Run("GetThingAllowingBatchWritesWithCompositeAttributes", GetThingAllowingBatchWritesWithCompositeAttributes(dbFactory(), t))
	t.Run("ScanThingAllowingBatchWritesWithCompositeAttributess", ScanThingAllowingBatchWritesWithCompositeAttributess(dbFactory(), t))
	t.Run("GetThingAllowingBatchWritesWithCompositeAttributessByNameIDAndDate", GetThingAllowingBatchWritesWithCompositeAttributessByNameIDAndDate(dbFactory(), t))
	t.Run("SaveThingAllowingBatchWritesWithCompositeAttributes", SaveThingAllowingBatchWritesWithCompositeAttributes(dbFactory(), t))
	t.Run("DeleteThingAllowingBatchWritesWithCompositeAttributes", DeleteThingAllowingBatchWritesWithCompositeAttributes(dbFactory(), t))
	t.Run("GetThingWithAdditionalAttributes", GetThingWithAdditionalAttributes(dbFactory(), t))
	t.Run("ScanThingWithAdditionalAttributess", ScanThingWithAdditionalAttributess(dbFactory(), t))
	t.Run("GetThingWithAdditionalAttributessByNameAndVersion", GetThingWithAdditionalAttributessByNameAndVersion(dbFactory(), t))
	t.Run("SaveThingWithAdditionalAttributes", SaveThingWithAdditionalAttributes(dbFactory(), t))
	t.Run("DeleteThingWithAdditionalAttributes", DeleteThingWithAdditionalAttributes(dbFactory(), t))
	t.Run("GetThingWithAdditionalAttributesByID", GetThingWithAdditionalAttributesByID(dbFactory(), t))
	t.Run("ScanThingWithAdditionalAttributessByID", ScanThingWithAdditionalAttributessByID(dbFactory(), t))
	t.Run("GetThingWithAdditionalAttributessByNameAndCreatedAt", GetThingWithAdditionalAttributessByNameAndCreatedAt(dbFactory(), t))
	t.Run("ScanThingWithAdditionalAttributessByNameAndCreatedAt", ScanThingWithAdditionalAttributessByNameAndCreatedAt(dbFactory(), t))
	t.Run("GetThingWithAdditionalAttributessByNameAndRangeNullable", GetThingWithAdditionalAttributessByNameAndRangeNullable(dbFactory(), t))
	t.Run("ScanThingWithAdditionalAttributessByNameAndRangeNullable", ScanThingWithAdditionalAttributessByNameAndRangeNullable(dbFactory(), t))
	t.Run("GetThingWithAdditionalAttributessByHashNullableAndName", GetThingWithAdditionalAttributessByHashNullableAndName(dbFactory(), t))
	t.Run("GetThingWithCompositeAttributes", GetThingWithCompositeAttributes(dbFactory(), t))
	t.Run("ScanThingWithCompositeAttributess", ScanThingWithCompositeAttributess(dbFactory(), t))
	t.Run("GetThingWithCompositeAttributessByNameBranchAndDate", GetThingWithCompositeAttributessByNameBranchAndDate(dbFactory(), t))
	t.Run("SaveThingWithCompositeAttributes", SaveThingWithCompositeAttributes(dbFactory(), t))
	t.Run("DeleteThingWithCompositeAttributes", DeleteThingWithCompositeAttributes(dbFactory(), t))
	t.Run("GetThingWithCompositeAttributessByNameVersionAndDate", GetThingWithCompositeAttributessByNameVersionAndDate(dbFactory(), t))
	t.Run("ScanThingWithCompositeAttributessByNameVersionAndDate", ScanThingWithCompositeAttributessByNameVersionAndDate(dbFactory(), t))
	t.Run("GetThingWithCompositeEnumAttributes", GetThingWithCompositeEnumAttributes(dbFactory(), t))
	t.Run("ScanThingWithCompositeEnumAttributess", ScanThingWithCompositeEnumAttributess(dbFactory(), t))
	t.Run("GetThingWithCompositeEnumAttributessByNameBranchAndDate", GetThingWithCompositeEnumAttributessByNameBranchAndDate(dbFactory(), t))
	t.Run("SaveThingWithCompositeEnumAttributes", SaveThingWithCompositeEnumAttributes(dbFactory(), t))
	t.Run("DeleteThingWithCompositeEnumAttributes", DeleteThingWithCompositeEnumAttributes(dbFactory(), t))
	t.Run("GetThingWithDateGSI", GetThingWithDateGSI(dbFactory(), t))
	t.Run("ScanThingWithDateGSIs", ScanThingWithDateGSIs(dbFactory(), t))
	t.Run("SaveThingWithDateGSI", SaveThingWithDateGSI(dbFactory(), t))
	t.Run("DeleteThingWithDateGSI", DeleteThingWithDateGSI(dbFactory(), t))
	t.Run("GetThingWithDateGSIsByIDAndDateR", GetThingWithDateGSIsByIDAndDateR(dbFactory(), t))
	t.Run("GetThingWithDateGSIsByDateHAndID", GetThingWithDateGSIsByDateHAndID(dbFactory(), t))
	t.Run("GetThingWithDateRange", GetThingWithDateRange(dbFactory(), t))
	t.Run("ScanThingWithDateRanges", ScanThingWithDateRanges(dbFactory(), t))
	t.Run("GetThingWithDateRangesByNameAndDate", GetThingWithDateRangesByNameAndDate(dbFactory(), t))
	t.Run("SaveThingWithDateRange", SaveThingWithDateRange(dbFactory(), t))
	t.Run("DeleteThingWithDateRange", DeleteThingWithDateRange(dbFactory(), t))
	t.Run("GetThingWithDateRangeKey", GetThingWithDateRangeKey(dbFactory(), t))
	t.Run("ScanThingWithDateRangeKeys", ScanThingWithDateRangeKeys(dbFactory(), t))
	t.Run("GetThingWithDateRangeKeysByIDAndDate", GetThingWithDateRangeKeysByIDAndDate(dbFactory(), t))
	t.Run("SaveThingWithDateRangeKey", SaveThingWithDateRangeKey(dbFactory(), t))
	t.Run("DeleteThingWithDateRangeKey", DeleteThingWithDateRangeKey(dbFactory(), t))
	t.Run("GetThingWithDateTimeComposite", GetThingWithDateTimeComposite(dbFactory(), t))
	t.Run("ScanThingWithDateTimeComposites", ScanThingWithDateTimeComposites(dbFactory(), t))
	t.Run("GetThingWithDateTimeCompositesByTypeIDAndCreatedResource", GetThingWithDateTimeCompositesByTypeIDAndCreatedResource(dbFactory(), t))
	t.Run("SaveThingWithDateTimeComposite", SaveThingWithDateTimeComposite(dbFactory(), t))
	t.Run("DeleteThingWithDateTimeComposite", DeleteThingWithDateTimeComposite(dbFactory(), t))
	t.Run("GetThingWithDatetimeGSI", GetThingWithDatetimeGSI(dbFactory(), t))
	t.Run("ScanThingWithDatetimeGSIs", ScanThingWithDatetimeGSIs(dbFactory(), t))
	t.Run("SaveThingWithDatetimeGSI", SaveThingWithDatetimeGSI(dbFactory(), t))
	t.Run("DeleteThingWithDatetimeGSI", DeleteThingWithDatetimeGSI(dbFactory(), t))
	t.Run("GetThingWithDatetimeGSIsByDatetimeAndID", GetThingWithDatetimeGSIsByDatetimeAndID(dbFactory(), t))
	t.Run("ScanThingWithDatetimeGSIsByDatetimeAndID", ScanThingWithDatetimeGSIsByDatetimeAndID(dbFactory(), t))
	t.Run("GetThingWithEnumHashKey", GetThingWithEnumHashKey(dbFactory(), t))
	t.Run("ScanThingWithEnumHashKeys", ScanThingWithEnumHashKeys(dbFactory(), t))
	t.Run("GetThingWithEnumHashKeysByBranchAndDate", GetThingWithEnumHashKeysByBranchAndDate(dbFactory(), t))
	t.Run("SaveThingWithEnumHashKey", SaveThingWithEnumHashKey(dbFactory(), t))
	t.Run("DeleteThingWithEnumHashKey", DeleteThingWithEnumHashKey(dbFactory(), t))
	t.Run("GetThingWithEnumHashKeysByBranchAndDate2", GetThingWithEnumHashKeysByBranchAndDate2(dbFactory(), t))
	t.Run("ScanThingWithEnumHashKeysByBranchAndDate2", ScanThingWithEnumHashKeysByBranchAndDate2(dbFactory(), t))
	t.Run("GetThingWithMatchingKeys", GetThingWithMatchingKeys(dbFactory(), t))
	t.Run("ScanThingWithMatchingKeyss", ScanThingWithMatchingKeyss(dbFactory(), t))
	t.Run("GetThingWithMatchingKeyssByBearAndAssocTypeID", GetThingWithMatchingKeyssByBearAndAssocTypeID(dbFactory(), t))
	t.Run("SaveThingWithMatchingKeys", SaveThingWithMatchingKeys(dbFactory(), t))
	t.Run("DeleteThingWithMatchingKeys", DeleteThingWithMatchingKeys(dbFactory(), t))
	t.Run("GetThingWithMatchingKeyssByAssocTypeIDAndCreatedBear", GetThingWithMatchingKeyssByAssocTypeIDAndCreatedBear(dbFactory(), t))
	t.Run("ScanThingWithMatchingKeyssByAssocTypeIDAndCreatedBear", ScanThingWithMatchingKeyssByAssocTypeIDAndCreatedBear(dbFactory(), t))
	t.Run("GetThingWithMultiUseCompositeAttribute", GetThingWithMultiUseCompositeAttribute(dbFactory(), t))
	t.Run("ScanThingWithMultiUseCompositeAttributes", ScanThingWithMultiUseCompositeAttributes(dbFactory(), t))
	t.Run("SaveThingWithMultiUseCompositeAttribute", SaveThingWithMultiUseCompositeAttribute(dbFactory(), t))
	t.Run("DeleteThingWithMultiUseCompositeAttribute", DeleteThingWithMultiUseCompositeAttribute(dbFactory(), t))
	t.Run("GetThingWithMultiUseCompositeAttributesByThreeAndOneTwo", GetThingWithMultiUseCompositeAttributesByThreeAndOneTwo(dbFactory(), t))
	t.Run("ScanThingWithMultiUseCompositeAttributesByThreeAndOneTwo", ScanThingWithMultiUseCompositeAttributesByThreeAndOneTwo(dbFactory(), t))
	t.Run("GetThingWithMultiUseCompositeAttributesByFourAndOneTwo", GetThingWithMultiUseCompositeAttributesByFourAndOneTwo(dbFactory(), t))
	t.Run("ScanThingWithMultiUseCompositeAttributesByFourAndOneTwo", ScanThingWithMultiUseCompositeAttributesByFourAndOneTwo(dbFactory(), t))
	t.Run("GetThingWithRequiredCompositePropertiesAndKeysOnly", GetThingWithRequiredCompositePropertiesAndKeysOnly(dbFactory(), t))
	t.Run("ScanThingWithRequiredCompositePropertiesAndKeysOnlys", ScanThingWithRequiredCompositePropertiesAndKeysOnlys(dbFactory(), t))
	t.Run("SaveThingWithRequiredCompositePropertiesAndKeysOnly", SaveThingWithRequiredCompositePropertiesAndKeysOnly(dbFactory(), t))
	t.Run("DeleteThingWithRequiredCompositePropertiesAndKeysOnly", DeleteThingWithRequiredCompositePropertiesAndKeysOnly(dbFactory(), t))
	t.Run("GetThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThree", GetThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThree(dbFactory(), t))
	t.Run("ScanThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThree", ScanThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThree(dbFactory(), t))
	t.Run("GetThingWithRequiredFields", GetThingWithRequiredFields(dbFactory(), t))
	t.Run("ScanThingWithRequiredFieldss", ScanThingWithRequiredFieldss(dbFactory(), t))
	t.Run("SaveThingWithRequiredFields", SaveThingWithRequiredFields(dbFactory(), t))
	t.Run("DeleteThingWithRequiredFields", DeleteThingWithRequiredFields(dbFactory(), t))
	t.Run("GetThingWithRequiredFields2", GetThingWithRequiredFields2(dbFactory(), t))
	t.Run("ScanThingWithRequiredFields2s", ScanThingWithRequiredFields2s(dbFactory(), t))
	t.Run("GetThingWithRequiredFields2sByNameAndID", GetThingWithRequiredFields2sByNameAndID(dbFactory(), t))
	t.Run("SaveThingWithRequiredFields2", SaveThingWithRequiredFields2(dbFactory(), t))
	t.Run("DeleteThingWithRequiredFields2", DeleteThingWithRequiredFields2(dbFactory(), t))
	t.Run("GetThingWithTransactMultipleGSI", GetThingWithTransactMultipleGSI(dbFactory(), t))
	t.Run("ScanThingWithTransactMultipleGSIs", ScanThingWithTransactMultipleGSIs(dbFactory(), t))
	t.Run("SaveThingWithTransactMultipleGSI", SaveThingWithTransactMultipleGSI(dbFactory(), t))
	t.Run("DeleteThingWithTransactMultipleGSI", DeleteThingWithTransactMultipleGSI(dbFactory(), t))
	t.Run("GetThingWithTransactMultipleGSIsByIDAndDateR", GetThingWithTransactMultipleGSIsByIDAndDateR(dbFactory(), t))
	t.Run("TransactSaveThingWithTransactMultipleGSIAndThing", TransactSaveThingWithTransactMultipleGSIAndThing(dbFactory(), t))
	t.Run("GetThingWithTransactMultipleGSIsByDateHAndID", GetThingWithTransactMultipleGSIsByDateHAndID(dbFactory(), t))
	t.Run("TransactSaveThingWithTransactMultipleGSIAndThing", TransactSaveThingWithTransactMultipleGSIAndThing(dbFactory(), t))
	t.Run("GetThingWithTransaction", GetThingWithTransaction(dbFactory(), t))
	t.Run("ScanThingWithTransactions", ScanThingWithTransactions(dbFactory(), t))
	t.Run("SaveThingWithTransaction", SaveThingWithTransaction(dbFactory(), t))
	t.Run("DeleteThingWithTransaction", DeleteThingWithTransaction(dbFactory(), t))
	t.Run("GetThingWithTransactionWithSimpleThing", GetThingWithTransactionWithSimpleThing(dbFactory(), t))
	t.Run("ScanThingWithTransactionWithSimpleThings", ScanThingWithTransactionWithSimpleThings(dbFactory(), t))
	t.Run("SaveThingWithTransactionWithSimpleThing", SaveThingWithTransactionWithSimpleThing(dbFactory(), t))
	t.Run("DeleteThingWithTransactionWithSimpleThing", DeleteThingWithTransactionWithSimpleThing(dbFactory(), t))
	t.Run("GetThingWithUnderscores", GetThingWithUnderscores(dbFactory(), t))
	t.Run("SaveThingWithUnderscores", SaveThingWithUnderscores(dbFactory(), t))
	t.Run("DeleteThingWithUnderscores", DeleteThingWithUnderscores(dbFactory(), t))
}

func GetDeployment(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.Deployment{
			Application: "string1",
			Date:        mustTime("2018-03-11T15:04:01+07:00"),
			Environment: "string1",
			Version:     "string1",
		}
		require.Nil(t, s.SaveDeployment(ctx, m))
		m2, err := s.GetDeployment(ctx, m.Environment, m.Application, m.Version)
		require.Nil(t, err)
		require.Equal(t, m.Environment, m2.Environment)
		require.Equal(t, m.Application, m2.Application)
		require.Equal(t, m.Version, m2.Version)

		_, err = s.GetDeployment(ctx, "string2", "string2", "string2")
		require.NotNil(t, err)
		require.IsType(t, err, db.ErrDeploymentNotFound{})
	}
}

type getDeploymentsByEnvAppAndVersionInput struct {
	ctx   context.Context
	input db.GetDeploymentsByEnvAppAndVersionInput
}
type getDeploymentsByEnvAppAndVersionOutput struct {
	deployments []models.Deployment
	err         error
}
type getDeploymentsByEnvAppAndVersionTest struct {
	testName string
	d        db.Interface
	input    getDeploymentsByEnvAppAndVersionInput
	output   getDeploymentsByEnvAppAndVersionOutput
}

func (g getDeploymentsByEnvAppAndVersionTest) run(t *testing.T) {
	deployments := []models.Deployment{}
	fn := func(m *models.Deployment, lastDeployment bool) bool {
		deployments = append(deployments, *m)
		if lastDeployment {
			return false
		}
		return true
	}
	err := g.d.GetDeploymentsByEnvAppAndVersion(g.input.ctx, g.input.input, fn)
	if err != nil {
		fmt.Println(err.Error())
	}
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.deployments, deployments)
}

func GetDeploymentsByEnvAppAndVersion(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveDeployment(ctx, models.Deployment{
			Environment: "string1",
			Application: "string1",
			Version:     "string1",
		}))
		require.Nil(t, d.SaveDeployment(ctx, models.Deployment{
			Environment: "string1",
			Application: "string1",
			Version:     "string2",
		}))
		require.Nil(t, d.SaveDeployment(ctx, models.Deployment{
			Environment: "string1",
			Application: "string1",
			Version:     "string3",
		}))
		limit := int64(3)
		tests := []getDeploymentsByEnvAppAndVersionTest{
			{
				testName: "basic",
				d:        d,
				input: getDeploymentsByEnvAppAndVersionInput{
					ctx: context.Background(),
					input: db.GetDeploymentsByEnvAppAndVersionInput{
						Environment: "string1",
						Application: "string1",
						Limit:       &limit,
					},
				},
				output: getDeploymentsByEnvAppAndVersionOutput{
					deployments: []models.Deployment{
						models.Deployment{
							Environment: "string1",
							Application: "string1",
							Version:     "string1",
						},
						models.Deployment{
							Environment: "string1",
							Application: "string1",
							Version:     "string2",
						},
						models.Deployment{
							Environment: "string1",
							Application: "string1",
							Version:     "string3",
						},
					},
					err: nil,
				},
			},
			{
				testName: "descending",
				d:        d,
				input: getDeploymentsByEnvAppAndVersionInput{
					ctx: context.Background(),
					input: db.GetDeploymentsByEnvAppAndVersionInput{
						Environment: "string1",
						Application: "string1",
						Descending:  true,
					},
				},
				output: getDeploymentsByEnvAppAndVersionOutput{
					deployments: []models.Deployment{
						models.Deployment{
							Environment: "string1",
							Application: "string1",
							Version:     "string3",
						},
						models.Deployment{
							Environment: "string1",
							Application: "string1",
							Version:     "string2",
						},
						models.Deployment{
							Environment: "string1",
							Application: "string1",
							Version:     "string1",
						},
					},
					err: nil,
				},
			},
			/* FAILING_TEST */
			/* {
			   testName: "starting after",
			   d:    d,
			   input: getDeploymentsByEnvAppAndVersionInput{
			     ctx: context.Background(),
			     input: db.GetDeploymentsByEnvAppAndVersionInput{
			         Environment: "string1",
			         Application: "string1",
			       StartingAfter: &models.Deployment{
			           Environment:    "string1",
			           Application:    "string1",
			           Version:    "string1",
			       },
			     },
			   },
			   output: getDeploymentsByEnvAppAndVersionOutput{
			     deployments: []models.Deployment{
			       models.Deployment{
			           Environment:    "string1",
			           Application:    "string1",
			           Version: "string2",
			       },
			       models.Deployment{
			           Environment:    "string1",
			           Application:    "string1",
			           Version: "string3",
			       },
			     },
			     err: nil,
			   },
			 }, */
			{
				testName: "starting after descending",
				d:        d,
				input: getDeploymentsByEnvAppAndVersionInput{
					ctx: context.Background(),
					input: db.GetDeploymentsByEnvAppAndVersionInput{
						Environment: "string1",
						Application: "string1",
						StartingAfter: &models.Deployment{
							Environment: "string1",
							Application: "string1",
							Version:     "string3",
						},
						Descending: true,
					},
				},
				output: getDeploymentsByEnvAppAndVersionOutput{
					deployments: []models.Deployment{
						models.Deployment{
							Environment: "string1",
							Application: "string1",
							Version:     "string2",
						},
						models.Deployment{
							Environment: "string1",
							Application: "string1",
							Version:     "string1",
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting at",
				d:        d,
				input: getDeploymentsByEnvAppAndVersionInput{
					ctx: context.Background(),
					input: db.GetDeploymentsByEnvAppAndVersionInput{
						Environment:       "string1",
						Application:       "string1",
						VersionStartingAt: db.String("string2"),
					},
				},
				output: getDeploymentsByEnvAppAndVersionOutput{
					deployments: []models.Deployment{
						models.Deployment{
							Environment: "string1",
							Application: "string1",
							Version:     "string2",
						},
						models.Deployment{
							Environment: "string1",
							Application: "string1",
							Version:     "string3",
						},
					},
					err: nil,
				},
			},
		}
		for _, test := range tests {
			t.Run(test.testName, test.run)
		}
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanDeployments(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveDeployment(ctx, models.Deployment{
			Application: "string1",
			Date:        mustTime("2018-03-11T15:04:01+07:00"),
			Environment: "string1",
			Version:     "string1",
		}))
		require.Nil(t, d.SaveDeployment(ctx, models.Deployment{
			Application: "string2",
			Date:        mustTime("2018-03-11T15:04:02+07:00"),
			Environment: "string2",
			Version:     "string2",
		}))
		require.Nil(t, d.SaveDeployment(ctx, models.Deployment{
			Application: "string3",
			Date:        mustTime("2018-03-11T15:04:03+07:00"),
			Environment: "string3",
			Version:     "string3",
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.Deployment{
				models.Deployment{
					Application: "string1",
					Date:        mustTime("2018-03-11T15:04:01+07:00"),
					Environment: "string1",
					Version:     "string1",
				},
				models.Deployment{
					Application: "string2",
					Date:        mustTime("2018-03-11T15:04:02+07:00"),
					Environment: "string2",
					Version:     "string2",
				},
				models.Deployment{
					Application: "string3",
					Date:        mustTime("2018-03-11T15:04:03+07:00"),
					Environment: "string3",
					Version:     "string3",
				},
			}
			actual := []models.Deployment{}
			err := d.ScanDeployments(ctx, db.ScanDeploymentsInput{}, func(m *models.Deployment, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		// FAILING_TEST
		t.Run("starting after", func(t *testing.T) {
			t.Skip()
			// Scan for everything.
			allItems := []models.Deployment{}
			err := d.ScanDeployments(ctx, db.ScanDeploymentsInput{}, func(m *models.Deployment, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput := db.ScanDeploymentsInput{
				StartingAfter: &models.Deployment{
					Environment: firstItem.Environment,
					Application: firstItem.Application,
					Version:     firstItem.Version,
				},
			}
			actual := []models.Deployment{}
			err = d.ScanDeployments(ctx, scanInput, func(m *models.Deployment, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanDeploymentsInput{
				Limit: &limit,
			}
			actual := []models.Deployment{}
			err := d.ScanDeployments(ctx, scanInput, func(m *models.Deployment, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

func SaveDeployment(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.Deployment{
			Application: "string1",
			Date:        mustTime("2018-03-11T15:04:01+07:00"),
			Environment: "string1",
			Version:     "string1",
		}
		require.Nil(t, s.SaveDeployment(ctx, m))
	}
}

func DeleteDeployment(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.Deployment{
			Application: "string1",
			Date:        mustTime("2018-03-11T15:04:01+07:00"),
			Environment: "string1",
			Version:     "string1",
		}
		require.Nil(t, s.SaveDeployment(ctx, m))
		require.Nil(t, s.DeleteDeployment(ctx, m.Environment, m.Application, m.Version))
	}
}

type getDeploymentsByEnvAppAndDateInput struct {
	ctx   context.Context
	input db.GetDeploymentsByEnvAppAndDateInput
}
type getDeploymentsByEnvAppAndDateOutput struct {
	deployments []models.Deployment
	err         error
}
type getDeploymentsByEnvAppAndDateTest struct {
	testName string
	d        db.Interface
	input    getDeploymentsByEnvAppAndDateInput
	output   getDeploymentsByEnvAppAndDateOutput
}

func (g getDeploymentsByEnvAppAndDateTest) run(t *testing.T) {
	deployments := []models.Deployment{}
	fn := func(m *models.Deployment, lastDeployment bool) bool {
		deployments = append(deployments, *m)
		if lastDeployment {
			return false
		}
		return true
	}
	err := g.d.GetDeploymentsByEnvAppAndDate(g.input.ctx, g.input.input, fn)
	if err != nil {
		fmt.Println(err.Error())
	}
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.deployments, deployments)
}

func GetDeploymentsByEnvAppAndDate(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveDeployment(ctx, models.Deployment{
			Environment: "string1",
			Application: "string1",
			Date:        mustTime("2018-03-11T15:04:01+07:00"),
			Version:     "string1",
		}))
		require.Nil(t, d.SaveDeployment(ctx, models.Deployment{
			Environment: "string1",
			Application: "string1",
			Date:        mustTime("2018-03-11T15:04:02+07:00"),
			Version:     "string3",
		}))
		require.Nil(t, d.SaveDeployment(ctx, models.Deployment{
			Environment: "string1",
			Application: "string1",
			Date:        mustTime("2018-03-11T15:04:03+07:00"),
			Version:     "string2",
		}))
		limit := int64(3)
		tests := []getDeploymentsByEnvAppAndDateTest{
			{
				testName: "basic",
				d:        d,
				input: getDeploymentsByEnvAppAndDateInput{
					ctx: context.Background(),
					input: db.GetDeploymentsByEnvAppAndDateInput{
						Environment: "string1",
						Application: "string1",
						Limit:       &limit,
					},
				},
				output: getDeploymentsByEnvAppAndDateOutput{
					deployments: []models.Deployment{
						models.Deployment{
							Environment: "string1",
							Application: "string1",
							Date:        mustTime("2018-03-11T15:04:01+07:00"),
							Version:     "string1",
						},
						models.Deployment{
							Environment: "string1",
							Application: "string1",
							Date:        mustTime("2018-03-11T15:04:02+07:00"),
							Version:     "string3",
						},
						models.Deployment{
							Environment: "string1",
							Application: "string1",
							Date:        mustTime("2018-03-11T15:04:03+07:00"),
							Version:     "string2",
						},
					},
					err: nil,
				},
			},
			{
				testName: "descending",
				d:        d,
				input: getDeploymentsByEnvAppAndDateInput{
					ctx: context.Background(),
					input: db.GetDeploymentsByEnvAppAndDateInput{
						Environment: "string1",
						Application: "string1",
						Descending:  true,
					},
				},
				output: getDeploymentsByEnvAppAndDateOutput{
					deployments: []models.Deployment{
						models.Deployment{
							Environment: "string1",
							Application: "string1",
							Date:        mustTime("2018-03-11T15:04:03+07:00"),
							Version:     "string2",
						},
						models.Deployment{
							Environment: "string1",
							Application: "string1",
							Date:        mustTime("2018-03-11T15:04:02+07:00"),
							Version:     "string3",
						},
						models.Deployment{
							Environment: "string1",
							Application: "string1",
							Date:        mustTime("2018-03-11T15:04:01+07:00"),
							Version:     "string1",
						},
					},
					err: nil,
				},
			},
			/* FAILING_TEST */
			/* {
			   testName: "starting after",
			   d:    d,
			   input: getDeploymentsByEnvAppAndDateInput{
			     ctx: context.Background(),
			     input: db.GetDeploymentsByEnvAppAndDateInput{
			         Environment: "string1",
			         Application: "string1",
			       StartingAfter: &models.Deployment{
			         Environment:    "string1",
			         Application:    "string1",
			         Date: mustTime("2018-03-11T15:04:01+07:00"),
			         Version:    "string1",
			       },
			     },
			   },
			   output: getDeploymentsByEnvAppAndDateOutput{
			     deployments: []models.Deployment{
			       models.Deployment{
			         Environment:    "string1",
			         Application:    "string1",
			         Date: mustTime("2018-03-11T15:04:02+07:00"),
			         Version:    "string3",
			       },
			       models.Deployment{
			         Environment:    "string1",
			         Application:    "string1",
			         Date: mustTime("2018-03-11T15:04:03+07:00"),
			         Version:    "string2",
			       },
			     },
			     err: nil,
			   },
			 }, */
			{
				testName: "starting after descending",
				d:        d,
				input: getDeploymentsByEnvAppAndDateInput{
					ctx: context.Background(),
					input: db.GetDeploymentsByEnvAppAndDateInput{
						Environment: "string1",
						Application: "string1",
						StartingAfter: &models.Deployment{
							Environment: "string1",
							Application: "string1",
							Date:        mustTime("2018-03-11T15:04:03+07:00"),
							Version:     "string2",
						},
						Descending: true,
					},
				},
				output: getDeploymentsByEnvAppAndDateOutput{
					deployments: []models.Deployment{
						models.Deployment{
							Environment: "string1",
							Application: "string1",
							Date:        mustTime("2018-03-11T15:04:02+07:00"),
							Version:     "string3",
						},
						models.Deployment{
							Environment: "string1",
							Application: "string1",
							Date:        mustTime("2018-03-11T15:04:01+07:00"),
							Version:     "string1",
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting at",
				d:        d,
				input: getDeploymentsByEnvAppAndDateInput{
					ctx: context.Background(),
					input: db.GetDeploymentsByEnvAppAndDateInput{
						Environment:    "string1",
						Application:    "string1",
						DateStartingAt: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
					},
				},
				output: getDeploymentsByEnvAppAndDateOutput{
					deployments: []models.Deployment{
						models.Deployment{
							Environment: "string1",
							Application: "string1",
							Date:        mustTime("2018-03-11T15:04:02+07:00"),
							Version:     "string3",
						},
						models.Deployment{
							Environment: "string1",
							Application: "string1",
							Date:        mustTime("2018-03-11T15:04:03+07:00"),
							Version:     "string2",
						},
					},
					err: nil,
				},
			},
		}
		for _, test := range tests {
			t.Run(test.testName, test.run)
		}
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanDeploymentsByEnvAppAndDate(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveDeployment(ctx, models.Deployment{
			Environment: "string1",
			Application: "string1",
			Date:        mustTime("2018-03-11T15:04:01+07:00"),
			Version:     "string1",
		}))
		require.Nil(t, d.SaveDeployment(ctx, models.Deployment{
			Environment: "string2",
			Application: "string2",
			Date:        mustTime("2018-03-11T15:04:02+07:00"),
			Version:     "string2",
		}))
		require.Nil(t, d.SaveDeployment(ctx, models.Deployment{
			Environment: "string3",
			Application: "string3",
			Date:        mustTime("2018-03-11T15:04:03+07:00"),
			Version:     "string3",
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.Deployment{
				models.Deployment{
					Environment: "string1",
					Application: "string1",
					Date:        mustTime("2018-03-11T15:04:01+07:00"),
					Version:     "string1",
				},
				models.Deployment{
					Environment: "string2",
					Application: "string2",
					Date:        mustTime("2018-03-11T15:04:02+07:00"),
					Version:     "string2",
				},
				models.Deployment{
					Environment: "string3",
					Application: "string3",
					Date:        mustTime("2018-03-11T15:04:03+07:00"),
					Version:     "string3",
				},
			}
			// Consistent read must be disabled when scaning a GSI.
			scanInput := db.ScanDeploymentsByEnvAppAndDateInput{DisableConsistentRead: true}
			actual := []models.Deployment{}
			err := d.ScanDeploymentsByEnvAppAndDate(ctx, scanInput, func(m *models.Deployment, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		t.Run("starting after", func(t *testing.T) {
			// Scan for everything.
			allItems := []models.Deployment{}
			// Consistent read must be disabled when scaning a GSI.
			scanInput := db.ScanDeploymentsByEnvAppAndDateInput{DisableConsistentRead: true}
			err := d.ScanDeploymentsByEnvAppAndDate(ctx, scanInput, func(m *models.Deployment, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput = db.ScanDeploymentsByEnvAppAndDateInput{
				DisableConsistentRead: true,
				StartingAfter: &models.Deployment{
					Environment: firstItem.Environment,
					Application: firstItem.Application,
					Date:        firstItem.Date,
					Version:     firstItem.Version,
				},
			}
			actual := []models.Deployment{}
			err = d.ScanDeploymentsByEnvAppAndDate(ctx, scanInput, func(m *models.Deployment, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanDeploymentsInput{
				Limit: &limit,
			}
			actual := []models.Deployment{}
			err := d.ScanDeployments(ctx, scanInput, func(m *models.Deployment, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

type getDeploymentsByEnvironmentAndDateInput struct {
	ctx   context.Context
	input db.GetDeploymentsByEnvironmentAndDateInput
}
type getDeploymentsByEnvironmentAndDateOutput struct {
	deployments []models.Deployment
	err         error
}
type getDeploymentsByEnvironmentAndDateTest struct {
	testName string
	d        db.Interface
	input    getDeploymentsByEnvironmentAndDateInput
	output   getDeploymentsByEnvironmentAndDateOutput
}

func (g getDeploymentsByEnvironmentAndDateTest) run(t *testing.T) {
	deployments := []models.Deployment{}
	fn := func(m *models.Deployment, lastDeployment bool) bool {
		deployments = append(deployments, *m)
		if lastDeployment {
			return false
		}
		return true
	}
	err := g.d.GetDeploymentsByEnvironmentAndDate(g.input.ctx, g.input.input, fn)
	if err != nil {
		fmt.Println(err.Error())
	}
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.deployments, deployments)
}

func GetDeploymentsByEnvironmentAndDate(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveDeployment(ctx, models.Deployment{
			Environment: "string1",
			Date:        mustTime("2018-03-11T15:04:01+07:00"),
			Application: "string1",
			Version:     "string1",
		}))
		require.Nil(t, d.SaveDeployment(ctx, models.Deployment{
			Environment: "string1",
			Date:        mustTime("2018-03-11T15:04:02+07:00"),
			Application: "string3",
			Version:     "string3",
		}))
		require.Nil(t, d.SaveDeployment(ctx, models.Deployment{
			Environment: "string1",
			Date:        mustTime("2018-03-11T15:04:03+07:00"),
			Application: "string2",
			Version:     "string2",
		}))
		limit := int64(3)
		tests := []getDeploymentsByEnvironmentAndDateTest{
			{
				testName: "basic",
				d:        d,
				input: getDeploymentsByEnvironmentAndDateInput{
					ctx: context.Background(),
					input: db.GetDeploymentsByEnvironmentAndDateInput{
						Environment: "string1",
						Limit:       &limit,
					},
				},
				output: getDeploymentsByEnvironmentAndDateOutput{
					deployments: []models.Deployment{
						models.Deployment{
							Environment: "string1",
							Date:        mustTime("2018-03-11T15:04:01+07:00"),
							Application: "string1",
							Version:     "string1",
						},
						models.Deployment{
							Environment: "string1",
							Date:        mustTime("2018-03-11T15:04:02+07:00"),
							Application: "string3",
							Version:     "string3",
						},
						models.Deployment{
							Environment: "string1",
							Date:        mustTime("2018-03-11T15:04:03+07:00"),
							Application: "string2",
							Version:     "string2",
						},
					},
					err: nil,
				},
			},
			{
				testName: "descending",
				d:        d,
				input: getDeploymentsByEnvironmentAndDateInput{
					ctx: context.Background(),
					input: db.GetDeploymentsByEnvironmentAndDateInput{
						Environment: "string1",
						Descending:  true,
					},
				},
				output: getDeploymentsByEnvironmentAndDateOutput{
					deployments: []models.Deployment{
						models.Deployment{
							Environment: "string1",
							Date:        mustTime("2018-03-11T15:04:03+07:00"),
							Application: "string2",
							Version:     "string2",
						},
						models.Deployment{
							Environment: "string1",
							Date:        mustTime("2018-03-11T15:04:02+07:00"),
							Application: "string3",
							Version:     "string3",
						},
						models.Deployment{
							Environment: "string1",
							Date:        mustTime("2018-03-11T15:04:01+07:00"),
							Application: "string1",
							Version:     "string1",
						},
					},
					err: nil,
				},
			},
			/* FAILING_TEST */
			/* {
			   testName: "starting after",
			   d:    d,
			   input: getDeploymentsByEnvironmentAndDateInput{
			     ctx: context.Background(),
			     input: db.GetDeploymentsByEnvironmentAndDateInput{
			         Environment: "string1",
			       StartingAfter: &models.Deployment{
			         Environment:    "string1",
			         Date: mustTime("2018-03-11T15:04:01+07:00"),
			         Application:    "string1",
			         Version:    "string1",
			       },
			     },
			   },
			   output: getDeploymentsByEnvironmentAndDateOutput{
			     deployments: []models.Deployment{
			       models.Deployment{
			         Environment:    "string1",
			         Date: mustTime("2018-03-11T15:04:02+07:00"),
			         Application:    "string3",
			         Version:    "string3",
			       },
			       models.Deployment{
			         Environment:    "string1",
			         Date: mustTime("2018-03-11T15:04:03+07:00"),
			         Application:    "string2",
			         Version:    "string2",
			       },
			     },
			     err: nil,
			   },
			 }, */
			{
				testName: "starting after descending",
				d:        d,
				input: getDeploymentsByEnvironmentAndDateInput{
					ctx: context.Background(),
					input: db.GetDeploymentsByEnvironmentAndDateInput{
						Environment: "string1",
						StartingAfter: &models.Deployment{
							Environment: "string1",
							Date:        mustTime("2018-03-11T15:04:03+07:00"),
							Application: "string2",
							Version:     "string2",
						},
						Descending: true,
					},
				},
				output: getDeploymentsByEnvironmentAndDateOutput{
					deployments: []models.Deployment{
						models.Deployment{
							Environment: "string1",
							Date:        mustTime("2018-03-11T15:04:02+07:00"),
							Application: "string3",
							Version:     "string3",
						},
						models.Deployment{
							Environment: "string1",
							Date:        mustTime("2018-03-11T15:04:01+07:00"),
							Application: "string1",
							Version:     "string1",
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting at",
				d:        d,
				input: getDeploymentsByEnvironmentAndDateInput{
					ctx: context.Background(),
					input: db.GetDeploymentsByEnvironmentAndDateInput{
						Environment:    "string1",
						DateStartingAt: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
					},
				},
				output: getDeploymentsByEnvironmentAndDateOutput{
					deployments: []models.Deployment{
						models.Deployment{
							Environment: "string1",
							Date:        mustTime("2018-03-11T15:04:02+07:00"),
							Application: "string3",
							Version:     "string3",
						},
						models.Deployment{
							Environment: "string1",
							Date:        mustTime("2018-03-11T15:04:03+07:00"),
							Application: "string2",
							Version:     "string2",
						},
					},
					err: nil,
				},
			},
		}
		for _, test := range tests {
			t.Run(test.testName, test.run)
		}
	}
}

func GetDeploymentByVersion(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.Deployment{
			Application: "string1",
			Date:        mustTime("2018-03-11T15:04:01+07:00"),
			Environment: "string1",
			Version:     "string1",
		}
		require.Nil(t, s.SaveDeployment(ctx, m))
		m2, err := s.GetDeploymentByVersion(ctx, m.Version)
		require.Nil(t, err)
		require.Equal(t, m.Application, m2.Application)
		require.Equal(t, m.Date.String(), m2.Date.String())
		require.Equal(t, m.Environment, m2.Environment)
		require.Equal(t, m.Version, m2.Version)

		_, err = s.GetDeploymentByVersion(ctx, "string2")
		require.NotNil(t, err)
		require.IsType(t, err, db.ErrDeploymentByVersionNotFound{})
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanDeploymentsByVersion(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveDeployment(ctx, models.Deployment{
			Version:     "string1",
			Application: "string1",
			Date:        mustTime("2018-03-11T15:04:01+07:00"),
			Environment: "string1",
		}))
		require.Nil(t, d.SaveDeployment(ctx, models.Deployment{
			Version:     "string2",
			Application: "string2",
			Date:        mustTime("2018-03-11T15:04:02+07:00"),
			Environment: "string2",
		}))
		require.Nil(t, d.SaveDeployment(ctx, models.Deployment{
			Version:     "string3",
			Application: "string3",
			Date:        mustTime("2018-03-11T15:04:03+07:00"),
			Environment: "string3",
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.Deployment{
				models.Deployment{
					Version:     "string1",
					Application: "string1",
					Date:        mustTime("2018-03-11T15:04:01+07:00"),
					Environment: "string1",
				},
				models.Deployment{
					Version:     "string2",
					Application: "string2",
					Date:        mustTime("2018-03-11T15:04:02+07:00"),
					Environment: "string2",
				},
				models.Deployment{
					Version:     "string3",
					Application: "string3",
					Date:        mustTime("2018-03-11T15:04:03+07:00"),
					Environment: "string3",
				},
			}
			// Consistent read must be disabled when scaning a GSI.
			scanInput := db.ScanDeploymentsByVersionInput{DisableConsistentRead: true}
			actual := []models.Deployment{}
			err := d.ScanDeploymentsByVersion(ctx, scanInput, func(m *models.Deployment, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		t.Run("starting after", func(t *testing.T) {
			// Scan for everything.
			allItems := []models.Deployment{}
			// Consistent read must be disabled when scaning a GSI.
			scanInput := db.ScanDeploymentsByVersionInput{DisableConsistentRead: true}
			err := d.ScanDeploymentsByVersion(ctx, scanInput, func(m *models.Deployment, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput = db.ScanDeploymentsByVersionInput{
				DisableConsistentRead: true,
				StartingAfter: &models.Deployment{
					Version:     firstItem.Version,
					Environment: firstItem.Environment,
					Application: firstItem.Application,
				},
			}
			actual := []models.Deployment{}
			err = d.ScanDeploymentsByVersion(ctx, scanInput, func(m *models.Deployment, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanDeploymentsInput{
				Limit: &limit,
			}
			actual := []models.Deployment{}
			err := d.ScanDeployments(ctx, scanInput, func(m *models.Deployment, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

func GetEvent(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.Event{
			Data: []byte("string1"),
			Pk:   "string1",
			Sk:   "string1",
		}
		require.Nil(t, s.SaveEvent(ctx, m))
		m2, err := s.GetEvent(ctx, m.Pk, m.Sk)
		require.Nil(t, err)
		require.Equal(t, m.Pk, m2.Pk)
		require.Equal(t, m.Sk, m2.Sk)

		_, err = s.GetEvent(ctx, "string2", "string2")
		require.NotNil(t, err)
		require.IsType(t, err, db.ErrEventNotFound{})
	}
}

type getEventsByPkAndSkInput struct {
	ctx   context.Context
	input db.GetEventsByPkAndSkInput
}
type getEventsByPkAndSkOutput struct {
	events []models.Event
	err    error
}
type getEventsByPkAndSkTest struct {
	testName string
	d        db.Interface
	input    getEventsByPkAndSkInput
	output   getEventsByPkAndSkOutput
}

func (g getEventsByPkAndSkTest) run(t *testing.T) {
	events := []models.Event{}
	fn := func(m *models.Event, lastEvent bool) bool {
		events = append(events, *m)
		if lastEvent {
			return false
		}
		return true
	}
	err := g.d.GetEventsByPkAndSk(g.input.ctx, g.input.input, fn)
	if err != nil {
		fmt.Println(err.Error())
	}
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.events, events)
}

func GetEventsByPkAndSk(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveEvent(ctx, models.Event{
			Pk: "string1",
			Sk: "string1",
		}))
		require.Nil(t, d.SaveEvent(ctx, models.Event{
			Pk: "string1",
			Sk: "string2",
		}))
		require.Nil(t, d.SaveEvent(ctx, models.Event{
			Pk: "string1",
			Sk: "string3",
		}))
		limit := int64(3)
		tests := []getEventsByPkAndSkTest{
			{
				testName: "basic",
				d:        d,
				input: getEventsByPkAndSkInput{
					ctx: context.Background(),
					input: db.GetEventsByPkAndSkInput{
						Pk:    "string1",
						Limit: &limit,
					},
				},
				output: getEventsByPkAndSkOutput{
					events: []models.Event{
						models.Event{
							Pk: "string1",
							Sk: "string1",
						},
						models.Event{
							Pk: "string1",
							Sk: "string2",
						},
						models.Event{
							Pk: "string1",
							Sk: "string3",
						},
					},
					err: nil,
				},
			},
			{
				testName: "descending",
				d:        d,
				input: getEventsByPkAndSkInput{
					ctx: context.Background(),
					input: db.GetEventsByPkAndSkInput{
						Pk:         "string1",
						Descending: true,
					},
				},
				output: getEventsByPkAndSkOutput{
					events: []models.Event{
						models.Event{
							Pk: "string1",
							Sk: "string3",
						},
						models.Event{
							Pk: "string1",
							Sk: "string2",
						},
						models.Event{
							Pk: "string1",
							Sk: "string1",
						},
					},
					err: nil,
				},
			},
			/* FAILING_TEST */
			/* {
			   testName: "starting after",
			   d:    d,
			   input: getEventsByPkAndSkInput{
			     ctx: context.Background(),
			     input: db.GetEventsByPkAndSkInput{
			         Pk: "string1",
			       StartingAfter: &models.Event{
			           Pk:    "string1",
			           Sk:    "string1",
			       },
			     },
			   },
			   output: getEventsByPkAndSkOutput{
			     events: []models.Event{
			       models.Event{
			           Pk:    "string1",
			           Sk: "string2",
			       },
			       models.Event{
			           Pk:    "string1",
			           Sk: "string3",
			       },
			     },
			     err: nil,
			   },
			 }, */
			{
				testName: "starting after descending",
				d:        d,
				input: getEventsByPkAndSkInput{
					ctx: context.Background(),
					input: db.GetEventsByPkAndSkInput{
						Pk: "string1",
						StartingAfter: &models.Event{
							Pk: "string1",
							Sk: "string3",
						},
						Descending: true,
					},
				},
				output: getEventsByPkAndSkOutput{
					events: []models.Event{
						models.Event{
							Pk: "string1",
							Sk: "string2",
						},
						models.Event{
							Pk: "string1",
							Sk: "string1",
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting at",
				d:        d,
				input: getEventsByPkAndSkInput{
					ctx: context.Background(),
					input: db.GetEventsByPkAndSkInput{
						Pk:           "string1",
						SkStartingAt: db.String("string2"),
					},
				},
				output: getEventsByPkAndSkOutput{
					events: []models.Event{
						models.Event{
							Pk: "string1",
							Sk: "string2",
						},
						models.Event{
							Pk: "string1",
							Sk: "string3",
						},
					},
					err: nil,
				},
			},
		}
		for _, test := range tests {
			t.Run(test.testName, test.run)
		}
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanEvents(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveEvent(ctx, models.Event{
			Data: []byte("string1"),
			Pk:   "string1",
			Sk:   "string1",
		}))
		require.Nil(t, d.SaveEvent(ctx, models.Event{
			Data: []byte("string2"),
			Pk:   "string2",
			Sk:   "string2",
		}))
		require.Nil(t, d.SaveEvent(ctx, models.Event{
			Data: []byte("string3"),
			Pk:   "string3",
			Sk:   "string3",
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.Event{
				models.Event{
					Data: []byte("string1"),
					Pk:   "string1",
					Sk:   "string1",
				},
				models.Event{
					Data: []byte("string2"),
					Pk:   "string2",
					Sk:   "string2",
				},
				models.Event{
					Data: []byte("string3"),
					Pk:   "string3",
					Sk:   "string3",
				},
			}
			actual := []models.Event{}
			err := d.ScanEvents(ctx, db.ScanEventsInput{}, func(m *models.Event, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		// FAILING_TEST
		t.Run("starting after", func(t *testing.T) {
			t.Skip()
			// Scan for everything.
			allItems := []models.Event{}
			err := d.ScanEvents(ctx, db.ScanEventsInput{}, func(m *models.Event, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput := db.ScanEventsInput{
				StartingAfter: &models.Event{
					Pk: firstItem.Pk,
					Sk: firstItem.Sk,
				},
			}
			actual := []models.Event{}
			err = d.ScanEvents(ctx, scanInput, func(m *models.Event, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanEventsInput{
				Limit: &limit,
			}
			actual := []models.Event{}
			err := d.ScanEvents(ctx, scanInput, func(m *models.Event, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

func SaveEvent(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.Event{
			Data: []byte("string1"),
			Pk:   "string1",
			Sk:   "string1",
		}
		require.Nil(t, s.SaveEvent(ctx, m))
	}
}

func DeleteEvent(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.Event{
			Data: []byte("string1"),
			Pk:   "string1",
			Sk:   "string1",
		}
		require.Nil(t, s.SaveEvent(ctx, m))
		require.Nil(t, s.DeleteEvent(ctx, m.Pk, m.Sk))
	}
}

type getEventsBySkAndDataInput struct {
	ctx   context.Context
	input db.GetEventsBySkAndDataInput
}
type getEventsBySkAndDataOutput struct {
	events []models.Event
	err    error
}
type getEventsBySkAndDataTest struct {
	testName string
	d        db.Interface
	input    getEventsBySkAndDataInput
	output   getEventsBySkAndDataOutput
}

func (g getEventsBySkAndDataTest) run(t *testing.T) {
	events := []models.Event{}
	fn := func(m *models.Event, lastEvent bool) bool {
		events = append(events, *m)
		if lastEvent {
			return false
		}
		return true
	}
	err := g.d.GetEventsBySkAndData(g.input.ctx, g.input.input, fn)
	if err != nil {
		fmt.Println(err.Error())
	}
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.events, events)
}

func GetEventsBySkAndData(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveEvent(ctx, models.Event{
			Sk:   "string1",
			Data: []byte("string1"),
			Pk:   "string1",
		}))
		require.Nil(t, d.SaveEvent(ctx, models.Event{
			Sk:   "string1",
			Data: []byte("string2"),
			Pk:   "string3",
		}))
		require.Nil(t, d.SaveEvent(ctx, models.Event{
			Sk:   "string1",
			Data: []byte("string3"),
			Pk:   "string2",
		}))
		limit := int64(3)
		tests := []getEventsBySkAndDataTest{
			{
				testName: "basic",
				d:        d,
				input: getEventsBySkAndDataInput{
					ctx: context.Background(),
					input: db.GetEventsBySkAndDataInput{
						Sk:    "string1",
						Limit: &limit,
					},
				},
				output: getEventsBySkAndDataOutput{
					events: []models.Event{
						models.Event{
							Sk:   "string1",
							Data: []byte("string1"),
							Pk:   "string1",
						},
						models.Event{
							Sk:   "string1",
							Data: []byte("string2"),
							Pk:   "string3",
						},
						models.Event{
							Sk:   "string1",
							Data: []byte("string3"),
							Pk:   "string2",
						},
					},
					err: nil,
				},
			},
			{
				testName: "descending",
				d:        d,
				input: getEventsBySkAndDataInput{
					ctx: context.Background(),
					input: db.GetEventsBySkAndDataInput{
						Sk:         "string1",
						Descending: true,
					},
				},
				output: getEventsBySkAndDataOutput{
					events: []models.Event{
						models.Event{
							Sk:   "string1",
							Data: []byte("string3"),
							Pk:   "string2",
						},
						models.Event{
							Sk:   "string1",
							Data: []byte("string2"),
							Pk:   "string3",
						},
						models.Event{
							Sk:   "string1",
							Data: []byte("string1"),
							Pk:   "string1",
						},
					},
					err: nil,
				},
			},
			/* FAILING_TEST */
			/* {
			   testName: "starting after",
			   d:    d,
			   input: getEventsBySkAndDataInput{
			     ctx: context.Background(),
			     input: db.GetEventsBySkAndDataInput{
			         Sk: "string1",
			       StartingAfter: &models.Event{
			         Sk:    "string1",
			         Data: []byte("string1"),
			         Pk:    "string1",
			       },
			     },
			   },
			   output: getEventsBySkAndDataOutput{
			     events: []models.Event{
			       models.Event{
			         Sk:    "string1",
			         Data: []byte("string2"),
			         Pk:    "string3",
			       },
			       models.Event{
			         Sk:    "string1",
			         Data: []byte("string3"),
			         Pk:    "string2",
			       },
			     },
			     err: nil,
			   },
			 }, */
			{
				testName: "starting after descending",
				d:        d,
				input: getEventsBySkAndDataInput{
					ctx: context.Background(),
					input: db.GetEventsBySkAndDataInput{
						Sk: "string1",
						StartingAfter: &models.Event{
							Sk:   "string1",
							Data: []byte("string3"),
							Pk:   "string2",
						},
						Descending: true,
					},
				},
				output: getEventsBySkAndDataOutput{
					events: []models.Event{
						models.Event{
							Sk:   "string1",
							Data: []byte("string2"),
							Pk:   "string3",
						},
						models.Event{
							Sk:   "string1",
							Data: []byte("string1"),
							Pk:   "string1",
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting at",
				d:        d,
				input: getEventsBySkAndDataInput{
					ctx: context.Background(),
					input: db.GetEventsBySkAndDataInput{
						Sk:             "string1",
						DataStartingAt: []byte("string2"),
					},
				},
				output: getEventsBySkAndDataOutput{
					events: []models.Event{
						models.Event{
							Sk:   "string1",
							Data: []byte("string2"),
							Pk:   "string3",
						},
						models.Event{
							Sk:   "string1",
							Data: []byte("string3"),
							Pk:   "string2",
						},
					},
					err: nil,
				},
			},
		}
		for _, test := range tests {
			t.Run(test.testName, test.run)
		}
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanEventsBySkAndData(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveEvent(ctx, models.Event{
			Sk:   "string1",
			Data: []byte("string1"),
			Pk:   "string1",
		}))
		require.Nil(t, d.SaveEvent(ctx, models.Event{
			Sk:   "string2",
			Data: []byte("string2"),
			Pk:   "string2",
		}))
		require.Nil(t, d.SaveEvent(ctx, models.Event{
			Sk:   "string3",
			Data: []byte("string3"),
			Pk:   "string3",
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.Event{
				models.Event{
					Sk:   "string1",
					Data: []byte("string1"),
					Pk:   "string1",
				},
				models.Event{
					Sk:   "string2",
					Data: []byte("string2"),
					Pk:   "string2",
				},
				models.Event{
					Sk:   "string3",
					Data: []byte("string3"),
					Pk:   "string3",
				},
			}
			// Consistent read must be disabled when scaning a GSI.
			scanInput := db.ScanEventsBySkAndDataInput{DisableConsistentRead: true}
			actual := []models.Event{}
			err := d.ScanEventsBySkAndData(ctx, scanInput, func(m *models.Event, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		t.Run("starting after", func(t *testing.T) {
			// Scan for everything.
			allItems := []models.Event{}
			// Consistent read must be disabled when scaning a GSI.
			scanInput := db.ScanEventsBySkAndDataInput{DisableConsistentRead: true}
			err := d.ScanEventsBySkAndData(ctx, scanInput, func(m *models.Event, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput = db.ScanEventsBySkAndDataInput{
				DisableConsistentRead: true,
				StartingAfter: &models.Event{
					Sk:   firstItem.Sk,
					Data: firstItem.Data,
					Pk:   firstItem.Pk,
				},
			}
			actual := []models.Event{}
			err = d.ScanEventsBySkAndData(ctx, scanInput, func(m *models.Event, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanEventsInput{
				Limit: &limit,
			}
			actual := []models.Event{}
			err := d.ScanEvents(ctx, scanInput, func(m *models.Event, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

func GetNoRangeThingWithCompositeAttributes(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.NoRangeThingWithCompositeAttributes{
			Branch:  db.String("string1"),
			Commit:  db.String("string1"),
			Date:    db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			Name:    db.String("string1"),
			Version: 1,
		}
		require.Nil(t, s.SaveNoRangeThingWithCompositeAttributes(ctx, m))
		m2, err := s.GetNoRangeThingWithCompositeAttributes(ctx, *m.Name, *m.Branch)
		require.Nil(t, err)
		require.Equal(t, *m.Name, *m2.Name)
		require.Equal(t, *m.Branch, *m2.Branch)

		_, err = s.GetNoRangeThingWithCompositeAttributes(ctx, "string2", "string2")
		require.NotNil(t, err)
		require.IsType(t, err, db.ErrNoRangeThingWithCompositeAttributesNotFound{})
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanNoRangeThingWithCompositeAttributess(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveNoRangeThingWithCompositeAttributes(ctx, models.NoRangeThingWithCompositeAttributes{
			Branch:  db.String("string1"),
			Commit:  db.String("string1"),
			Date:    db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			Name:    db.String("string1"),
			Version: 1,
		}))
		require.Nil(t, d.SaveNoRangeThingWithCompositeAttributes(ctx, models.NoRangeThingWithCompositeAttributes{
			Branch:  db.String("string2"),
			Commit:  db.String("string2"),
			Date:    db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
			Name:    db.String("string2"),
			Version: 2,
		}))
		require.Nil(t, d.SaveNoRangeThingWithCompositeAttributes(ctx, models.NoRangeThingWithCompositeAttributes{
			Branch:  db.String("string3"),
			Commit:  db.String("string3"),
			Date:    db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
			Name:    db.String("string3"),
			Version: 3,
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.NoRangeThingWithCompositeAttributes{
				models.NoRangeThingWithCompositeAttributes{
					Branch:  db.String("string1"),
					Commit:  db.String("string1"),
					Date:    db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
					Name:    db.String("string1"),
					Version: 1,
				},
				models.NoRangeThingWithCompositeAttributes{
					Branch:  db.String("string2"),
					Commit:  db.String("string2"),
					Date:    db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
					Name:    db.String("string2"),
					Version: 2,
				},
				models.NoRangeThingWithCompositeAttributes{
					Branch:  db.String("string3"),
					Commit:  db.String("string3"),
					Date:    db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
					Name:    db.String("string3"),
					Version: 3,
				},
			}
			actual := []models.NoRangeThingWithCompositeAttributes{}
			err := d.ScanNoRangeThingWithCompositeAttributess(ctx, db.ScanNoRangeThingWithCompositeAttributessInput{}, func(m *models.NoRangeThingWithCompositeAttributes, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		// FAILING_TEST
		t.Run("starting after", func(t *testing.T) {
			t.Skip()
			// Scan for everything.
			allItems := []models.NoRangeThingWithCompositeAttributes{}
			err := d.ScanNoRangeThingWithCompositeAttributess(ctx, db.ScanNoRangeThingWithCompositeAttributessInput{}, func(m *models.NoRangeThingWithCompositeAttributes, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput := db.ScanNoRangeThingWithCompositeAttributessInput{
				StartingAfter: &models.NoRangeThingWithCompositeAttributes{
					Name:   firstItem.Name,
					Branch: firstItem.Branch,
				},
			}
			actual := []models.NoRangeThingWithCompositeAttributes{}
			err = d.ScanNoRangeThingWithCompositeAttributess(ctx, scanInput, func(m *models.NoRangeThingWithCompositeAttributes, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanNoRangeThingWithCompositeAttributessInput{
				Limit: &limit,
			}
			actual := []models.NoRangeThingWithCompositeAttributes{}
			err := d.ScanNoRangeThingWithCompositeAttributess(ctx, scanInput, func(m *models.NoRangeThingWithCompositeAttributes, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

func SaveNoRangeThingWithCompositeAttributes(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.NoRangeThingWithCompositeAttributes{
			Branch:  db.String("string1"),
			Commit:  db.String("string1"),
			Date:    db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			Name:    db.String("string1"),
			Version: 1,
		}
		require.Nil(t, s.SaveNoRangeThingWithCompositeAttributes(ctx, m))
		require.IsType(t, db.ErrNoRangeThingWithCompositeAttributesAlreadyExists{}, s.SaveNoRangeThingWithCompositeAttributes(ctx, m))
	}
}

func DeleteNoRangeThingWithCompositeAttributes(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.NoRangeThingWithCompositeAttributes{
			Branch:  db.String("string1"),
			Commit:  db.String("string1"),
			Date:    db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			Name:    db.String("string1"),
			Version: 1,
		}
		require.Nil(t, s.SaveNoRangeThingWithCompositeAttributes(ctx, m))
		require.Nil(t, s.DeleteNoRangeThingWithCompositeAttributes(ctx, *m.Name, *m.Branch))
	}
}

type getNoRangeThingWithCompositeAttributessByNameVersionAndDateInput struct {
	ctx   context.Context
	input db.GetNoRangeThingWithCompositeAttributessByNameVersionAndDateInput
}
type getNoRangeThingWithCompositeAttributessByNameVersionAndDateOutput struct {
	noRangeThingWithCompositeAttributess []models.NoRangeThingWithCompositeAttributes
	err                                  error
}
type getNoRangeThingWithCompositeAttributessByNameVersionAndDateTest struct {
	testName string
	d        db.Interface
	input    getNoRangeThingWithCompositeAttributessByNameVersionAndDateInput
	output   getNoRangeThingWithCompositeAttributessByNameVersionAndDateOutput
}

func (g getNoRangeThingWithCompositeAttributessByNameVersionAndDateTest) run(t *testing.T) {
	noRangeThingWithCompositeAttributess := []models.NoRangeThingWithCompositeAttributes{}
	fn := func(m *models.NoRangeThingWithCompositeAttributes, lastNoRangeThingWithCompositeAttributes bool) bool {
		noRangeThingWithCompositeAttributess = append(noRangeThingWithCompositeAttributess, *m)
		if lastNoRangeThingWithCompositeAttributes {
			return false
		}
		return true
	}
	err := g.d.GetNoRangeThingWithCompositeAttributessByNameVersionAndDate(g.input.ctx, g.input.input, fn)
	if err != nil {
		fmt.Println(err.Error())
	}
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.noRangeThingWithCompositeAttributess, noRangeThingWithCompositeAttributess)
}

func GetNoRangeThingWithCompositeAttributessByNameVersionAndDate(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveNoRangeThingWithCompositeAttributes(ctx, models.NoRangeThingWithCompositeAttributes{
			Name:    db.String("string1"),
			Version: 1,
			Date:    db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			Branch:  db.String("string1"),
			Commit:  db.String("string1"),
		}))
		require.Nil(t, d.SaveNoRangeThingWithCompositeAttributes(ctx, models.NoRangeThingWithCompositeAttributes{
			Name:    db.String("string1"),
			Version: 1,
			Date:    db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
			Branch:  db.String("string3"),
			Commit:  db.String("string3"),
		}))
		require.Nil(t, d.SaveNoRangeThingWithCompositeAttributes(ctx, models.NoRangeThingWithCompositeAttributes{
			Name:    db.String("string1"),
			Version: 1,
			Date:    db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
			Branch:  db.String("string2"),
			Commit:  db.String("string2"),
		}))
		limit := int64(3)
		tests := []getNoRangeThingWithCompositeAttributessByNameVersionAndDateTest{
			{
				testName: "basic",
				d:        d,
				input: getNoRangeThingWithCompositeAttributessByNameVersionAndDateInput{
					ctx: context.Background(),
					input: db.GetNoRangeThingWithCompositeAttributessByNameVersionAndDateInput{
						Name:    "string1",
						Version: 1,
						Limit:   &limit,
					},
				},
				output: getNoRangeThingWithCompositeAttributessByNameVersionAndDateOutput{
					noRangeThingWithCompositeAttributess: []models.NoRangeThingWithCompositeAttributes{
						models.NoRangeThingWithCompositeAttributes{
							Name:    db.String("string1"),
							Version: 1,
							Date:    db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
							Branch:  db.String("string1"),
							Commit:  db.String("string1"),
						},
						models.NoRangeThingWithCompositeAttributes{
							Name:    db.String("string1"),
							Version: 1,
							Date:    db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
							Branch:  db.String("string3"),
							Commit:  db.String("string3"),
						},
						models.NoRangeThingWithCompositeAttributes{
							Name:    db.String("string1"),
							Version: 1,
							Date:    db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
							Branch:  db.String("string2"),
							Commit:  db.String("string2"),
						},
					},
					err: nil,
				},
			},
			{
				testName: "descending",
				d:        d,
				input: getNoRangeThingWithCompositeAttributessByNameVersionAndDateInput{
					ctx: context.Background(),
					input: db.GetNoRangeThingWithCompositeAttributessByNameVersionAndDateInput{
						Name:       "string1",
						Version:    1,
						Descending: true,
					},
				},
				output: getNoRangeThingWithCompositeAttributessByNameVersionAndDateOutput{
					noRangeThingWithCompositeAttributess: []models.NoRangeThingWithCompositeAttributes{
						models.NoRangeThingWithCompositeAttributes{
							Name:    db.String("string1"),
							Version: 1,
							Date:    db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
							Branch:  db.String("string2"),
							Commit:  db.String("string2"),
						},
						models.NoRangeThingWithCompositeAttributes{
							Name:    db.String("string1"),
							Version: 1,
							Date:    db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
							Branch:  db.String("string3"),
							Commit:  db.String("string3"),
						},
						models.NoRangeThingWithCompositeAttributes{
							Name:    db.String("string1"),
							Version: 1,
							Date:    db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
							Branch:  db.String("string1"),
							Commit:  db.String("string1"),
						},
					},
					err: nil,
				},
			},
			/* FAILING_TEST */
			/* {
			   testName: "starting after",
			   d:    d,
			   input: getNoRangeThingWithCompositeAttributessByNameVersionAndDateInput{
			     ctx: context.Background(),
			     input: db.GetNoRangeThingWithCompositeAttributessByNameVersionAndDateInput{
			         Name: "string1",
			         Version: 1,
			       StartingAfter: &models.NoRangeThingWithCompositeAttributes{
			         Name:    db.String("string1"),
			         Version:    1,
			         Date: db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			         Branch:    db.String("string1"),
			         Commit:    db.String("string1"),
			       },
			     },
			   },
			   output: getNoRangeThingWithCompositeAttributessByNameVersionAndDateOutput{
			     noRangeThingWithCompositeAttributess: []models.NoRangeThingWithCompositeAttributes{
			       models.NoRangeThingWithCompositeAttributes{
			         Name:    db.String("string1"),
			         Version:    1,
			         Date: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
			         Branch:    db.String("string3"),
			         Commit:    db.String("string3"),
			       },
			       models.NoRangeThingWithCompositeAttributes{
			         Name:    db.String("string1"),
			         Version:    1,
			         Date: db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
			         Branch:    db.String("string2"),
			         Commit:    db.String("string2"),
			       },
			     },
			     err: nil,
			   },
			 }, */
			{
				testName: "starting after descending",
				d:        d,
				input: getNoRangeThingWithCompositeAttributessByNameVersionAndDateInput{
					ctx: context.Background(),
					input: db.GetNoRangeThingWithCompositeAttributessByNameVersionAndDateInput{
						Name:    "string1",
						Version: 1,
						StartingAfter: &models.NoRangeThingWithCompositeAttributes{
							Name:    db.String("string1"),
							Version: 1,
							Date:    db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
							Branch:  db.String("string2"),
							Commit:  db.String("string2"),
						},
						Descending: true,
					},
				},
				output: getNoRangeThingWithCompositeAttributessByNameVersionAndDateOutput{
					noRangeThingWithCompositeAttributess: []models.NoRangeThingWithCompositeAttributes{
						models.NoRangeThingWithCompositeAttributes{
							Name:    db.String("string1"),
							Version: 1,
							Date:    db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
							Branch:  db.String("string3"),
							Commit:  db.String("string3"),
						},
						models.NoRangeThingWithCompositeAttributes{
							Name:    db.String("string1"),
							Version: 1,
							Date:    db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
							Branch:  db.String("string1"),
							Commit:  db.String("string1"),
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting at",
				d:        d,
				input: getNoRangeThingWithCompositeAttributessByNameVersionAndDateInput{
					ctx: context.Background(),
					input: db.GetNoRangeThingWithCompositeAttributessByNameVersionAndDateInput{
						Name:           "string1",
						Version:        1,
						DateStartingAt: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
					},
				},
				output: getNoRangeThingWithCompositeAttributessByNameVersionAndDateOutput{
					noRangeThingWithCompositeAttributess: []models.NoRangeThingWithCompositeAttributes{
						models.NoRangeThingWithCompositeAttributes{
							Name:    db.String("string1"),
							Version: 1,
							Date:    db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
							Branch:  db.String("string3"),
							Commit:  db.String("string3"),
						},
						models.NoRangeThingWithCompositeAttributes{
							Name:    db.String("string1"),
							Version: 1,
							Date:    db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
							Branch:  db.String("string2"),
							Commit:  db.String("string2"),
						},
					},
					err: nil,
				},
			},
		}
		for _, test := range tests {
			t.Run(test.testName, test.run)
		}
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanNoRangeThingWithCompositeAttributessByNameVersionAndDate(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveNoRangeThingWithCompositeAttributes(ctx, models.NoRangeThingWithCompositeAttributes{
			Name:    db.String("string1"),
			Version: 1,
			Date:    db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			Branch:  db.String("string1"),
			Commit:  db.String("string1"),
		}))
		require.Nil(t, d.SaveNoRangeThingWithCompositeAttributes(ctx, models.NoRangeThingWithCompositeAttributes{
			Name:    db.String("string2"),
			Version: 2,
			Date:    db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
			Branch:  db.String("string2"),
			Commit:  db.String("string2"),
		}))
		require.Nil(t, d.SaveNoRangeThingWithCompositeAttributes(ctx, models.NoRangeThingWithCompositeAttributes{
			Name:    db.String("string3"),
			Version: 3,
			Date:    db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
			Branch:  db.String("string3"),
			Commit:  db.String("string3"),
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.NoRangeThingWithCompositeAttributes{
				models.NoRangeThingWithCompositeAttributes{
					Name:    db.String("string1"),
					Version: 1,
					Date:    db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
					Branch:  db.String("string1"),
					Commit:  db.String("string1"),
				},
				models.NoRangeThingWithCompositeAttributes{
					Name:    db.String("string2"),
					Version: 2,
					Date:    db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
					Branch:  db.String("string2"),
					Commit:  db.String("string2"),
				},
				models.NoRangeThingWithCompositeAttributes{
					Name:    db.String("string3"),
					Version: 3,
					Date:    db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
					Branch:  db.String("string3"),
					Commit:  db.String("string3"),
				},
			}
			// Consistent read must be disabled when scaning a GSI.
			scanInput := db.ScanNoRangeThingWithCompositeAttributessByNameVersionAndDateInput{DisableConsistentRead: true}
			actual := []models.NoRangeThingWithCompositeAttributes{}
			err := d.ScanNoRangeThingWithCompositeAttributessByNameVersionAndDate(ctx, scanInput, func(m *models.NoRangeThingWithCompositeAttributes, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		t.Run("starting after", func(t *testing.T) {
			// Scan for everything.
			allItems := []models.NoRangeThingWithCompositeAttributes{}
			// Consistent read must be disabled when scaning a GSI.
			scanInput := db.ScanNoRangeThingWithCompositeAttributessByNameVersionAndDateInput{DisableConsistentRead: true}
			err := d.ScanNoRangeThingWithCompositeAttributessByNameVersionAndDate(ctx, scanInput, func(m *models.NoRangeThingWithCompositeAttributes, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput = db.ScanNoRangeThingWithCompositeAttributessByNameVersionAndDateInput{
				DisableConsistentRead: true,
				StartingAfter: &models.NoRangeThingWithCompositeAttributes{
					Name:    firstItem.Name,
					Version: firstItem.Version,
					Date:    firstItem.Date,
					Branch:  firstItem.Branch,
				},
			}
			actual := []models.NoRangeThingWithCompositeAttributes{}
			err = d.ScanNoRangeThingWithCompositeAttributessByNameVersionAndDate(ctx, scanInput, func(m *models.NoRangeThingWithCompositeAttributes, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanNoRangeThingWithCompositeAttributessInput{
				Limit: &limit,
			}
			actual := []models.NoRangeThingWithCompositeAttributes{}
			err := d.ScanNoRangeThingWithCompositeAttributess(ctx, scanInput, func(m *models.NoRangeThingWithCompositeAttributes, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

func GetNoRangeThingWithCompositeAttributesByNameBranchCommit(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.NoRangeThingWithCompositeAttributes{
			Branch:  db.String("string1"),
			Commit:  db.String("string1"),
			Date:    db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			Name:    db.String("string1"),
			Version: 1,
		}
		require.Nil(t, s.SaveNoRangeThingWithCompositeAttributes(ctx, m))
		m2, err := s.GetNoRangeThingWithCompositeAttributesByNameBranchCommit(ctx, *m.Name, *m.Branch, *m.Commit)
		require.Nil(t, err)
		require.Equal(t, m.Branch, m2.Branch)
		require.Equal(t, m.Commit, m2.Commit)
		require.Equal(t, m.Date.String(), m2.Date.String())
		require.Equal(t, m.Name, m2.Name)
		require.Equal(t, m.Version, m2.Version)

		_, err = s.GetNoRangeThingWithCompositeAttributesByNameBranchCommit(ctx, "string2", "string2", "string2")
		require.NotNil(t, err)
		require.IsType(t, err, db.ErrNoRangeThingWithCompositeAttributesByNameBranchCommitNotFound{})
	}
}

func GetSimpleThing(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.SimpleThing{
			Name: "string1",
		}
		require.Nil(t, s.SaveSimpleThing(ctx, m))
		m2, err := s.GetSimpleThing(ctx, m.Name)
		require.Nil(t, err)
		require.Equal(t, m.Name, m2.Name)

		_, err = s.GetSimpleThing(ctx, "string2")
		require.NotNil(t, err)
		require.IsType(t, err, db.ErrSimpleThingNotFound{})
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanSimpleThings(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveSimpleThing(ctx, models.SimpleThing{
			Name: "string1",
		}))
		require.Nil(t, d.SaveSimpleThing(ctx, models.SimpleThing{
			Name: "string2",
		}))
		require.Nil(t, d.SaveSimpleThing(ctx, models.SimpleThing{
			Name: "string3",
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.SimpleThing{
				models.SimpleThing{
					Name: "string1",
				},
				models.SimpleThing{
					Name: "string2",
				},
				models.SimpleThing{
					Name: "string3",
				},
			}
			actual := []models.SimpleThing{}
			err := d.ScanSimpleThings(ctx, db.ScanSimpleThingsInput{}, func(m *models.SimpleThing, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		// FAILING_TEST
		t.Run("starting after", func(t *testing.T) {
			t.Skip()
			// Scan for everything.
			allItems := []models.SimpleThing{}
			err := d.ScanSimpleThings(ctx, db.ScanSimpleThingsInput{}, func(m *models.SimpleThing, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput := db.ScanSimpleThingsInput{
				StartingAfter: &models.SimpleThing{
					Name: firstItem.Name,
				},
			}
			actual := []models.SimpleThing{}
			err = d.ScanSimpleThings(ctx, scanInput, func(m *models.SimpleThing, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanSimpleThingsInput{
				Limit: &limit,
			}
			actual := []models.SimpleThing{}
			err := d.ScanSimpleThings(ctx, scanInput, func(m *models.SimpleThing, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

func SaveSimpleThing(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.SimpleThing{
			Name: "string1",
		}
		require.Nil(t, s.SaveSimpleThing(ctx, m))
		require.IsType(t, db.ErrSimpleThingAlreadyExists{}, s.SaveSimpleThing(ctx, m))
	}
}

func DeleteSimpleThing(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.SimpleThing{
			Name: "string1",
		}
		require.Nil(t, s.SaveSimpleThing(ctx, m))
		require.Nil(t, s.DeleteSimpleThing(ctx, m.Name))
	}
}

func GetTeacherSharingRule(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.TeacherSharingRule{
			App:      "string1",
			District: "string1",
			School:   "string1",
			Teacher:  "string1",
		}
		require.Nil(t, s.SaveTeacherSharingRule(ctx, m))
		m2, err := s.GetTeacherSharingRule(ctx, m.Teacher, m.School, m.App)
		require.Nil(t, err)
		require.Equal(t, m.Teacher, m2.Teacher)
		require.Equal(t, m.School, m2.School)
		require.Equal(t, m.App, m2.App)

		_, err = s.GetTeacherSharingRule(ctx, "string2", "string2", "string2")
		require.NotNil(t, err)
		require.IsType(t, err, db.ErrTeacherSharingRuleNotFound{})
	}
}

type getTeacherSharingRulesByTeacherAndSchoolAppInput struct {
	ctx   context.Context
	input db.GetTeacherSharingRulesByTeacherAndSchoolAppInput
}
type getTeacherSharingRulesByTeacherAndSchoolAppOutput struct {
	teacherSharingRules []models.TeacherSharingRule
	err                 error
}
type getTeacherSharingRulesByTeacherAndSchoolAppTest struct {
	testName string
	d        db.Interface
	input    getTeacherSharingRulesByTeacherAndSchoolAppInput
	output   getTeacherSharingRulesByTeacherAndSchoolAppOutput
}

func (g getTeacherSharingRulesByTeacherAndSchoolAppTest) run(t *testing.T) {
	teacherSharingRules := []models.TeacherSharingRule{}
	fn := func(m *models.TeacherSharingRule, lastTeacherSharingRule bool) bool {
		teacherSharingRules = append(teacherSharingRules, *m)
		if lastTeacherSharingRule {
			return false
		}
		return true
	}
	err := g.d.GetTeacherSharingRulesByTeacherAndSchoolApp(g.input.ctx, g.input.input, fn)
	if err != nil {
		fmt.Println(err.Error())
	}
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.teacherSharingRules, teacherSharingRules)
}

func GetTeacherSharingRulesByTeacherAndSchoolApp(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveTeacherSharingRule(ctx, models.TeacherSharingRule{
			Teacher:  "string1",
			School:   "string1",
			App:      "string1",
			District: "district0",
		}))
		require.Nil(t, d.SaveTeacherSharingRule(ctx, models.TeacherSharingRule{
			Teacher:  "string1",
			School:   "string2",
			App:      "string2",
			District: "district1",
		}))
		require.Nil(t, d.SaveTeacherSharingRule(ctx, models.TeacherSharingRule{
			Teacher:  "string1",
			School:   "string3",
			App:      "string3",
			District: "district2",
		}))
		limit := int64(3)
		tests := []getTeacherSharingRulesByTeacherAndSchoolAppTest{
			{
				testName: "basic",
				d:        d,
				input: getTeacherSharingRulesByTeacherAndSchoolAppInput{
					ctx: context.Background(),
					input: db.GetTeacherSharingRulesByTeacherAndSchoolAppInput{
						Teacher: "string1",
						Limit:   &limit,
					},
				},
				output: getTeacherSharingRulesByTeacherAndSchoolAppOutput{
					teacherSharingRules: []models.TeacherSharingRule{
						models.TeacherSharingRule{
							Teacher:  "string1",
							School:   "string1",
							App:      "string1",
							District: "district0",
						},
						models.TeacherSharingRule{
							Teacher:  "string1",
							School:   "string2",
							App:      "string2",
							District: "district1",
						},
						models.TeacherSharingRule{
							Teacher:  "string1",
							School:   "string3",
							App:      "string3",
							District: "district2",
						},
					},
					err: nil,
				},
			},
			{
				testName: "descending",
				d:        d,
				input: getTeacherSharingRulesByTeacherAndSchoolAppInput{
					ctx: context.Background(),
					input: db.GetTeacherSharingRulesByTeacherAndSchoolAppInput{
						Teacher:    "string1",
						Descending: true,
					},
				},
				output: getTeacherSharingRulesByTeacherAndSchoolAppOutput{
					teacherSharingRules: []models.TeacherSharingRule{
						models.TeacherSharingRule{
							Teacher:  "string1",
							School:   "string3",
							App:      "string3",
							District: "district2",
						},
						models.TeacherSharingRule{
							Teacher:  "string1",
							School:   "string2",
							App:      "string2",
							District: "district1",
						},
						models.TeacherSharingRule{
							Teacher:  "string1",
							School:   "string1",
							App:      "string1",
							District: "district0",
						},
					},
					err: nil,
				},
			},
			/* FAILING_TEST */
			/* {
			   testName: "starting after",
			   d:    d,
			   input: getTeacherSharingRulesByTeacherAndSchoolAppInput{
			     ctx: context.Background(),
			     input: db.GetTeacherSharingRulesByTeacherAndSchoolAppInput{
			         Teacher: "string1",
			       StartingAfter: &models.TeacherSharingRule{
			           Teacher:    "string1",
			           School:    "string1",
			           App:    "string1",
			               District: "district0",
			       },
			     },
			   },
			   output: getTeacherSharingRulesByTeacherAndSchoolAppOutput{
			     teacherSharingRules: []models.TeacherSharingRule{
			       models.TeacherSharingRule{
			           Teacher:    "string1",
			           School: "string2",
			           App: "string2",
			               District: "district1",
			       },
			       models.TeacherSharingRule{
			           Teacher:    "string1",
			           School: "string3",
			           App: "string3",
			               District: "district2",
			       },
			     },
			     err: nil,
			   },
			 }, */
			{
				testName: "starting after descending",
				d:        d,
				input: getTeacherSharingRulesByTeacherAndSchoolAppInput{
					ctx: context.Background(),
					input: db.GetTeacherSharingRulesByTeacherAndSchoolAppInput{
						Teacher: "string1",
						StartingAfter: &models.TeacherSharingRule{
							Teacher:  "string1",
							School:   "string3",
							App:      "string3",
							District: "district2",
						},
						Descending: true,
					},
				},
				output: getTeacherSharingRulesByTeacherAndSchoolAppOutput{
					teacherSharingRules: []models.TeacherSharingRule{
						models.TeacherSharingRule{
							Teacher:  "string1",
							School:   "string2",
							App:      "string2",
							District: "district1",
						},
						models.TeacherSharingRule{
							Teacher:  "string1",
							School:   "string1",
							App:      "string1",
							District: "district0",
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting at",
				d:        d,
				input: getTeacherSharingRulesByTeacherAndSchoolAppInput{
					ctx: context.Background(),
					input: db.GetTeacherSharingRulesByTeacherAndSchoolAppInput{
						Teacher: "string1",
						StartingAt: &db.SchoolApp{
							School: "string2",
							App:    "string2",
						},
					},
				},
				output: getTeacherSharingRulesByTeacherAndSchoolAppOutput{
					teacherSharingRules: []models.TeacherSharingRule{
						models.TeacherSharingRule{
							Teacher:  "string1",
							School:   "string2",
							App:      "string2",
							District: "district1",
						},
						models.TeacherSharingRule{
							Teacher:  "string1",
							School:   "string3",
							App:      "string3",
							District: "district2",
						},
					},
					err: nil,
				},
			},
			{
				testName: "filtering",
				d:        d,
				input: getTeacherSharingRulesByTeacherAndSchoolAppInput{
					ctx: context.Background(),
					input: db.GetTeacherSharingRulesByTeacherAndSchoolAppInput{
						Teacher: "string1",
						FilterValues: []db.TeacherSharingRuleByTeacherAndSchoolAppFilterValues{
							db.TeacherSharingRuleByTeacherAndSchoolAppFilterValues{
								AttributeName:   db.TeacherSharingRuleDistrict,
								AttributeValues: []interface{}{"district0"},
							},
						},
						FilterExpression: "#DISTRICT = :district_value0",
						Limit:            &limit,
					},
				},
				output: getTeacherSharingRulesByTeacherAndSchoolAppOutput{
					teacherSharingRules: []models.TeacherSharingRule{
						models.TeacherSharingRule{
							Teacher:  "string1",
							School:   "string1",
							App:      "string1",
							District: "district0",
						},
					},
					err: nil,
				},
			},
		}
		for _, test := range tests {
			t.Run(test.testName, test.run)
		}
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanTeacherSharingRules(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveTeacherSharingRule(ctx, models.TeacherSharingRule{
			App:      "string1",
			District: "string1",
			School:   "string1",
			Teacher:  "string1",
		}))
		require.Nil(t, d.SaveTeacherSharingRule(ctx, models.TeacherSharingRule{
			App:      "string2",
			District: "string2",
			School:   "string2",
			Teacher:  "string2",
		}))
		require.Nil(t, d.SaveTeacherSharingRule(ctx, models.TeacherSharingRule{
			App:      "string3",
			District: "string3",
			School:   "string3",
			Teacher:  "string3",
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.TeacherSharingRule{
				models.TeacherSharingRule{
					App:      "string1",
					District: "string1",
					School:   "string1",
					Teacher:  "string1",
				},
				models.TeacherSharingRule{
					App:      "string2",
					District: "string2",
					School:   "string2",
					Teacher:  "string2",
				},
				models.TeacherSharingRule{
					App:      "string3",
					District: "string3",
					School:   "string3",
					Teacher:  "string3",
				},
			}
			actual := []models.TeacherSharingRule{}
			err := d.ScanTeacherSharingRules(ctx, db.ScanTeacherSharingRulesInput{}, func(m *models.TeacherSharingRule, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		// FAILING_TEST
		t.Run("starting after", func(t *testing.T) {
			t.Skip()
			// Scan for everything.
			allItems := []models.TeacherSharingRule{}
			err := d.ScanTeacherSharingRules(ctx, db.ScanTeacherSharingRulesInput{}, func(m *models.TeacherSharingRule, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput := db.ScanTeacherSharingRulesInput{
				StartingAfter: &models.TeacherSharingRule{
					Teacher: firstItem.Teacher,
					School:  firstItem.School,
					App:     firstItem.App,
					// must specify non-empty string values for attributes
					// in secondary indexes, since dynamodb doesn't support
					// empty strings:
					District: "district",
				},
			}
			actual := []models.TeacherSharingRule{}
			err = d.ScanTeacherSharingRules(ctx, scanInput, func(m *models.TeacherSharingRule, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanTeacherSharingRulesInput{
				Limit: &limit,
			}
			actual := []models.TeacherSharingRule{}
			err := d.ScanTeacherSharingRules(ctx, scanInput, func(m *models.TeacherSharingRule, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

func SaveTeacherSharingRule(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.TeacherSharingRule{
			App:      "string1",
			District: "string1",
			School:   "string1",
			Teacher:  "string1",
		}
		require.Nil(t, s.SaveTeacherSharingRule(ctx, m))
	}
}

func DeleteTeacherSharingRule(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.TeacherSharingRule{
			App:      "string1",
			District: "string1",
			School:   "string1",
			Teacher:  "string1",
		}
		require.Nil(t, s.SaveTeacherSharingRule(ctx, m))
		require.Nil(t, s.DeleteTeacherSharingRule(ctx, m.Teacher, m.School, m.App))
	}
}

type getTeacherSharingRulesByDistrictAndSchoolTeacherAppInput struct {
	ctx   context.Context
	input db.GetTeacherSharingRulesByDistrictAndSchoolTeacherAppInput
}
type getTeacherSharingRulesByDistrictAndSchoolTeacherAppOutput struct {
	teacherSharingRules []models.TeacherSharingRule
	err                 error
}
type getTeacherSharingRulesByDistrictAndSchoolTeacherAppTest struct {
	testName string
	d        db.Interface
	input    getTeacherSharingRulesByDistrictAndSchoolTeacherAppInput
	output   getTeacherSharingRulesByDistrictAndSchoolTeacherAppOutput
}

func (g getTeacherSharingRulesByDistrictAndSchoolTeacherAppTest) run(t *testing.T) {
	teacherSharingRules := []models.TeacherSharingRule{}
	fn := func(m *models.TeacherSharingRule, lastTeacherSharingRule bool) bool {
		teacherSharingRules = append(teacherSharingRules, *m)
		if lastTeacherSharingRule {
			return false
		}
		return true
	}
	err := g.d.GetTeacherSharingRulesByDistrictAndSchoolTeacherApp(g.input.ctx, g.input.input, fn)
	if err != nil {
		fmt.Println(err.Error())
	}
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.teacherSharingRules, teacherSharingRules)
}

func GetTeacherSharingRulesByDistrictAndSchoolTeacherApp(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveTeacherSharingRule(ctx, models.TeacherSharingRule{
			District: "string1",
			School:   "string1",
			Teacher:  "string1",
			App:      "string1",
		}))
		require.Nil(t, d.SaveTeacherSharingRule(ctx, models.TeacherSharingRule{
			District: "string1",
			School:   "string2",
			Teacher:  "string2",
			App:      "string2",
		}))
		require.Nil(t, d.SaveTeacherSharingRule(ctx, models.TeacherSharingRule{
			District: "string1",
			School:   "string3",
			Teacher:  "string3",
			App:      "string3",
		}))
		limit := int64(3)
		tests := []getTeacherSharingRulesByDistrictAndSchoolTeacherAppTest{
			{
				testName: "basic",
				d:        d,
				input: getTeacherSharingRulesByDistrictAndSchoolTeacherAppInput{
					ctx: context.Background(),
					input: db.GetTeacherSharingRulesByDistrictAndSchoolTeacherAppInput{
						District: "string1",
						Limit:    &limit,
					},
				},
				output: getTeacherSharingRulesByDistrictAndSchoolTeacherAppOutput{
					teacherSharingRules: []models.TeacherSharingRule{
						models.TeacherSharingRule{
							District: "string1",
							School:   "string1",
							Teacher:  "string1",
							App:      "string1",
						},
						models.TeacherSharingRule{
							District: "string1",
							School:   "string2",
							Teacher:  "string2",
							App:      "string2",
						},
						models.TeacherSharingRule{
							District: "string1",
							School:   "string3",
							Teacher:  "string3",
							App:      "string3",
						},
					},
					err: nil,
				},
			},
			{
				testName: "descending",
				d:        d,
				input: getTeacherSharingRulesByDistrictAndSchoolTeacherAppInput{
					ctx: context.Background(),
					input: db.GetTeacherSharingRulesByDistrictAndSchoolTeacherAppInput{
						District:   "string1",
						Descending: true,
					},
				},
				output: getTeacherSharingRulesByDistrictAndSchoolTeacherAppOutput{
					teacherSharingRules: []models.TeacherSharingRule{
						models.TeacherSharingRule{
							District: "string1",
							School:   "string3",
							Teacher:  "string3",
							App:      "string3",
						},
						models.TeacherSharingRule{
							District: "string1",
							School:   "string2",
							Teacher:  "string2",
							App:      "string2",
						},
						models.TeacherSharingRule{
							District: "string1",
							School:   "string1",
							Teacher:  "string1",
							App:      "string1",
						},
					},
					err: nil,
				},
			},
			/* FAILING_TEST */
			/* {
			   testName: "starting after",
			   d:    d,
			   input: getTeacherSharingRulesByDistrictAndSchoolTeacherAppInput{
			     ctx: context.Background(),
			     input: db.GetTeacherSharingRulesByDistrictAndSchoolTeacherAppInput{
			         District: "string1",
			       StartingAfter: &models.TeacherSharingRule{
			         District:    "string1",
			         School: "string1",
			         Teacher: "string1",
			         App: "string1",
			       },
			     },
			   },
			   output: getTeacherSharingRulesByDistrictAndSchoolTeacherAppOutput{
			     teacherSharingRules: []models.TeacherSharingRule{
			       models.TeacherSharingRule{
			         District:    "string1",
			         School: "string2",
			         Teacher: "string2",
			         App: "string2",
			       },
			       models.TeacherSharingRule{
			         District:    "string1",
			         School: "string3",
			         Teacher: "string3",
			         App: "string3",
			       },
			     },
			     err: nil,
			   },
			 }, */
			{
				testName: "starting after descending",
				d:        d,
				input: getTeacherSharingRulesByDistrictAndSchoolTeacherAppInput{
					ctx: context.Background(),
					input: db.GetTeacherSharingRulesByDistrictAndSchoolTeacherAppInput{
						District: "string1",
						StartingAfter: &models.TeacherSharingRule{
							District: "string1",
							School:   "string3",
							Teacher:  "string3",
							App:      "string3",
						},
						Descending: true,
					},
				},
				output: getTeacherSharingRulesByDistrictAndSchoolTeacherAppOutput{
					teacherSharingRules: []models.TeacherSharingRule{
						models.TeacherSharingRule{
							District: "string1",
							School:   "string2",
							Teacher:  "string2",
							App:      "string2",
						},
						models.TeacherSharingRule{
							District: "string1",
							School:   "string1",
							Teacher:  "string1",
							App:      "string1",
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting at",
				d:        d,
				input: getTeacherSharingRulesByDistrictAndSchoolTeacherAppInput{
					ctx: context.Background(),
					input: db.GetTeacherSharingRulesByDistrictAndSchoolTeacherAppInput{
						District: "string1",
						StartingAt: &db.SchoolTeacherApp{
							School:  "string2",
							Teacher: "string2",
							App:     "string2",
						},
					},
				},
				output: getTeacherSharingRulesByDistrictAndSchoolTeacherAppOutput{
					teacherSharingRules: []models.TeacherSharingRule{
						models.TeacherSharingRule{
							District: "string1",
							School:   "string2",
							Teacher:  "string2",
							App:      "string2",
						},
						models.TeacherSharingRule{
							District: "string1",
							School:   "string3",
							Teacher:  "string3",
							App:      "string3",
						},
					},
					err: nil,
				},
			},
		}
		for _, test := range tests {
			t.Run(test.testName, test.run)
		}
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanTeacherSharingRulesByDistrictAndSchoolTeacherApp(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveTeacherSharingRule(ctx, models.TeacherSharingRule{
			District: "string1",
			School:   "string1",
			Teacher:  "string1",
			App:      "string1",
		}))
		require.Nil(t, d.SaveTeacherSharingRule(ctx, models.TeacherSharingRule{
			District: "string2",
			School:   "string2",
			Teacher:  "string2",
			App:      "string2",
		}))
		require.Nil(t, d.SaveTeacherSharingRule(ctx, models.TeacherSharingRule{
			District: "string3",
			School:   "string3",
			Teacher:  "string3",
			App:      "string3",
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.TeacherSharingRule{
				models.TeacherSharingRule{
					District: "string1",
					School:   "string1",
					Teacher:  "string1",
					App:      "string1",
				},
				models.TeacherSharingRule{
					District: "string2",
					School:   "string2",
					Teacher:  "string2",
					App:      "string2",
				},
				models.TeacherSharingRule{
					District: "string3",
					School:   "string3",
					Teacher:  "string3",
					App:      "string3",
				},
			}
			// Consistent read must be disabled when scaning a GSI.
			scanInput := db.ScanTeacherSharingRulesByDistrictAndSchoolTeacherAppInput{DisableConsistentRead: true}
			actual := []models.TeacherSharingRule{}
			err := d.ScanTeacherSharingRulesByDistrictAndSchoolTeacherApp(ctx, scanInput, func(m *models.TeacherSharingRule, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		t.Run("starting after", func(t *testing.T) {
			// Scan for everything.
			allItems := []models.TeacherSharingRule{}
			// Consistent read must be disabled when scaning a GSI.
			scanInput := db.ScanTeacherSharingRulesByDistrictAndSchoolTeacherAppInput{DisableConsistentRead: true}
			err := d.ScanTeacherSharingRulesByDistrictAndSchoolTeacherApp(ctx, scanInput, func(m *models.TeacherSharingRule, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput = db.ScanTeacherSharingRulesByDistrictAndSchoolTeacherAppInput{
				DisableConsistentRead: true,
				StartingAfter: &models.TeacherSharingRule{
					District: firstItem.District,
					School:   firstItem.School,
					Teacher:  firstItem.Teacher,
					App:      firstItem.App,
				},
			}
			actual := []models.TeacherSharingRule{}
			err = d.ScanTeacherSharingRulesByDistrictAndSchoolTeacherApp(ctx, scanInput, func(m *models.TeacherSharingRule, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanTeacherSharingRulesInput{
				Limit: &limit,
			}
			actual := []models.TeacherSharingRule{}
			err := d.ScanTeacherSharingRules(ctx, scanInput, func(m *models.TeacherSharingRule, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

func GetThing(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.Thing{
			CreatedAt:     mustTime("2018-03-11T15:04:01+07:00"),
			HashNullable:  db.String("string1"),
			ID:            "string1",
			Name:          "string1",
			RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			Version:       1,
		}
		require.Nil(t, s.SaveThing(ctx, m))
		m2, err := s.GetThing(ctx, m.Name, m.Version)
		require.Nil(t, err)
		require.Equal(t, m.Name, m2.Name)
		require.Equal(t, m.Version, m2.Version)

		_, err = s.GetThing(ctx, "string2", 2)
		require.NotNil(t, err)
		require.IsType(t, err, db.ErrThingNotFound{})
	}
}

type getThingsByNameAndVersionInput struct {
	ctx   context.Context
	input db.GetThingsByNameAndVersionInput
}
type getThingsByNameAndVersionOutput struct {
	things []models.Thing
	err    error
}
type getThingsByNameAndVersionTest struct {
	testName string
	d        db.Interface
	input    getThingsByNameAndVersionInput
	output   getThingsByNameAndVersionOutput
}

func (g getThingsByNameAndVersionTest) run(t *testing.T) {
	things := []models.Thing{}
	fn := func(m *models.Thing, lastThing bool) bool {
		things = append(things, *m)
		if lastThing {
			return false
		}
		return true
	}
	err := g.d.GetThingsByNameAndVersion(g.input.ctx, g.input.input, fn)
	if err != nil {
		fmt.Println(err.Error())
	}
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.things, things)
}

func GetThingsByNameAndVersion(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThing(ctx, models.Thing{
			Name:    "string1",
			Version: 1,
		}))
		require.Nil(t, d.SaveThing(ctx, models.Thing{
			Name:    "string1",
			Version: 2,
		}))
		require.Nil(t, d.SaveThing(ctx, models.Thing{
			Name:    "string1",
			Version: 3,
		}))
		limit := int64(3)
		tests := []getThingsByNameAndVersionTest{
			{
				testName: "basic",
				d:        d,
				input: getThingsByNameAndVersionInput{
					ctx: context.Background(),
					input: db.GetThingsByNameAndVersionInput{
						Name:  "string1",
						Limit: &limit,
					},
				},
				output: getThingsByNameAndVersionOutput{
					things: []models.Thing{
						models.Thing{
							Name:    "string1",
							Version: 1,
						},
						models.Thing{
							Name:    "string1",
							Version: 2,
						},
						models.Thing{
							Name:    "string1",
							Version: 3,
						},
					},
					err: nil,
				},
			},
			{
				testName: "descending",
				d:        d,
				input: getThingsByNameAndVersionInput{
					ctx: context.Background(),
					input: db.GetThingsByNameAndVersionInput{
						Name:       "string1",
						Descending: true,
					},
				},
				output: getThingsByNameAndVersionOutput{
					things: []models.Thing{
						models.Thing{
							Name:    "string1",
							Version: 3,
						},
						models.Thing{
							Name:    "string1",
							Version: 2,
						},
						models.Thing{
							Name:    "string1",
							Version: 1,
						},
					},
					err: nil,
				},
			},
			/* FAILING_TEST */
			/* {
			   testName: "starting after",
			   d:    d,
			   input: getThingsByNameAndVersionInput{
			     ctx: context.Background(),
			     input: db.GetThingsByNameAndVersionInput{
			         Name: "string1",
			       StartingAfter: &models.Thing{
			           Name:    "string1",
			           Version:    1,
			       },
			     },
			   },
			   output: getThingsByNameAndVersionOutput{
			     things: []models.Thing{
			       models.Thing{
			           Name:    "string1",
			           Version: 2,
			       },
			       models.Thing{
			           Name:    "string1",
			           Version: 3,
			       },
			     },
			     err: nil,
			   },
			 }, */
			{
				testName: "starting after descending",
				d:        d,
				input: getThingsByNameAndVersionInput{
					ctx: context.Background(),
					input: db.GetThingsByNameAndVersionInput{
						Name: "string1",
						StartingAfter: &models.Thing{
							Name:    "string1",
							Version: 3,
						},
						Descending: true,
					},
				},
				output: getThingsByNameAndVersionOutput{
					things: []models.Thing{
						models.Thing{
							Name:    "string1",
							Version: 2,
						},
						models.Thing{
							Name:    "string1",
							Version: 1,
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting at",
				d:        d,
				input: getThingsByNameAndVersionInput{
					ctx: context.Background(),
					input: db.GetThingsByNameAndVersionInput{
						Name:              "string1",
						VersionStartingAt: db.Int64(2),
					},
				},
				output: getThingsByNameAndVersionOutput{
					things: []models.Thing{
						models.Thing{
							Name:    "string1",
							Version: 2,
						},
						models.Thing{
							Name:    "string1",
							Version: 3,
						},
					},
					err: nil,
				},
			},
		}
		for _, test := range tests {
			t.Run(test.testName, test.run)
		}
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanThings(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThing(ctx, models.Thing{
			CreatedAt:     mustTime("2018-03-11T15:04:01+07:00"),
			HashNullable:  db.String("string1"),
			ID:            "string1",
			Name:          "string1",
			RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			Version:       1,
		}))
		require.Nil(t, d.SaveThing(ctx, models.Thing{
			CreatedAt:     mustTime("2018-03-11T15:04:02+07:00"),
			HashNullable:  db.String("string2"),
			ID:            "string2",
			Name:          "string2",
			RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
			Version:       2,
		}))
		require.Nil(t, d.SaveThing(ctx, models.Thing{
			CreatedAt:     mustTime("2018-03-11T15:04:03+07:00"),
			HashNullable:  db.String("string3"),
			ID:            "string3",
			Name:          "string3",
			RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
			Version:       3,
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.Thing{
				models.Thing{
					CreatedAt:     mustTime("2018-03-11T15:04:01+07:00"),
					HashNullable:  db.String("string1"),
					ID:            "string1",
					Name:          "string1",
					RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
					Version:       1,
				},
				models.Thing{
					CreatedAt:     mustTime("2018-03-11T15:04:02+07:00"),
					HashNullable:  db.String("string2"),
					ID:            "string2",
					Name:          "string2",
					RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
					Version:       2,
				},
				models.Thing{
					CreatedAt:     mustTime("2018-03-11T15:04:03+07:00"),
					HashNullable:  db.String("string3"),
					ID:            "string3",
					Name:          "string3",
					RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
					Version:       3,
				},
			}
			actual := []models.Thing{}
			err := d.ScanThings(ctx, db.ScanThingsInput{}, func(m *models.Thing, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		// FAILING_TEST
		t.Run("starting after", func(t *testing.T) {
			t.Skip()
			// Scan for everything.
			allItems := []models.Thing{}
			err := d.ScanThings(ctx, db.ScanThingsInput{}, func(m *models.Thing, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput := db.ScanThingsInput{
				StartingAfter: &models.Thing{
					Name:    firstItem.Name,
					Version: firstItem.Version,
				},
			}
			actual := []models.Thing{}
			err = d.ScanThings(ctx, scanInput, func(m *models.Thing, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanThingsInput{
				Limit: &limit,
			}
			actual := []models.Thing{}
			err := d.ScanThings(ctx, scanInput, func(m *models.Thing, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

func SaveThing(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.Thing{
			CreatedAt:     mustTime("2018-03-11T15:04:01+07:00"),
			HashNullable:  db.String("string1"),
			ID:            "string1",
			Name:          "string1",
			RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			Version:       1,
		}
		require.Nil(t, s.SaveThing(ctx, m))
		require.IsType(t, db.ErrThingAlreadyExists{}, s.SaveThing(ctx, m))
	}
}

func DeleteThing(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.Thing{
			CreatedAt:     mustTime("2018-03-11T15:04:01+07:00"),
			HashNullable:  db.String("string1"),
			ID:            "string1",
			Name:          "string1",
			RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			Version:       1,
		}
		require.Nil(t, s.SaveThing(ctx, m))
		require.Nil(t, s.DeleteThing(ctx, m.Name, m.Version))
	}
}

func GetThingByID(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.Thing{
			CreatedAt:     mustTime("2018-03-11T15:04:01+07:00"),
			HashNullable:  db.String("string1"),
			ID:            "string1",
			Name:          "string1",
			RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			Version:       1,
		}
		require.Nil(t, s.SaveThing(ctx, m))
		m2, err := s.GetThingByID(ctx, m.ID)
		require.Nil(t, err)
		require.Equal(t, m.CreatedAt.String(), m2.CreatedAt.String())
		require.Equal(t, m.HashNullable, m2.HashNullable)
		require.Equal(t, m.ID, m2.ID)
		require.Equal(t, m.Name, m2.Name)
		require.Equal(t, m.RangeNullable.String(), m2.RangeNullable.String())
		require.Equal(t, m.Version, m2.Version)

		_, err = s.GetThingByID(ctx, "string2")
		require.NotNil(t, err)
		require.IsType(t, err, db.ErrThingByIDNotFound{})
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanThingsByID(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThing(ctx, models.Thing{
			ID:      "string1",
			Name:    "string1",
			Version: 1,
		}))
		require.Nil(t, d.SaveThing(ctx, models.Thing{
			ID:      "string2",
			Name:    "string2",
			Version: 2,
		}))
		require.Nil(t, d.SaveThing(ctx, models.Thing{
			ID:      "string3",
			Name:    "string3",
			Version: 3,
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.Thing{
				models.Thing{
					ID:      "string1",
					Name:    "string1",
					Version: 1,
				},
				models.Thing{
					ID:      "string2",
					Name:    "string2",
					Version: 2,
				},
				models.Thing{
					ID:      "string3",
					Name:    "string3",
					Version: 3,
				},
			}
			// Consistent read must be disabled when scaning a GSI.
			scanInput := db.ScanThingsByIDInput{DisableConsistentRead: true}
			actual := []models.Thing{}
			err := d.ScanThingsByID(ctx, scanInput, func(m *models.Thing, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		t.Run("starting after", func(t *testing.T) {
			// Scan for everything.
			allItems := []models.Thing{}
			// Consistent read must be disabled when scaning a GSI.
			scanInput := db.ScanThingsByIDInput{DisableConsistentRead: true}
			err := d.ScanThingsByID(ctx, scanInput, func(m *models.Thing, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput = db.ScanThingsByIDInput{
				DisableConsistentRead: true,
				StartingAfter: &models.Thing{
					ID:      firstItem.ID,
					Name:    firstItem.Name,
					Version: firstItem.Version,
				},
			}
			actual := []models.Thing{}
			err = d.ScanThingsByID(ctx, scanInput, func(m *models.Thing, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanThingsInput{
				Limit: &limit,
			}
			actual := []models.Thing{}
			err := d.ScanThings(ctx, scanInput, func(m *models.Thing, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

type getThingsByNameAndCreatedAtInput struct {
	ctx   context.Context
	input db.GetThingsByNameAndCreatedAtInput
}
type getThingsByNameAndCreatedAtOutput struct {
	things []models.Thing
	err    error
}
type getThingsByNameAndCreatedAtTest struct {
	testName string
	d        db.Interface
	input    getThingsByNameAndCreatedAtInput
	output   getThingsByNameAndCreatedAtOutput
}

func (g getThingsByNameAndCreatedAtTest) run(t *testing.T) {
	things := []models.Thing{}
	fn := func(m *models.Thing, lastThing bool) bool {
		things = append(things, *m)
		if lastThing {
			return false
		}
		return true
	}
	err := g.d.GetThingsByNameAndCreatedAt(g.input.ctx, g.input.input, fn)
	if err != nil {
		fmt.Println(err.Error())
	}
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.things, things)
}

func GetThingsByNameAndCreatedAt(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThing(ctx, models.Thing{
			Name:      "string1",
			CreatedAt: mustTime("2018-03-11T15:04:01+07:00"),
			Version:   1,
		}))
		require.Nil(t, d.SaveThing(ctx, models.Thing{
			Name:      "string1",
			CreatedAt: mustTime("2018-03-11T15:04:02+07:00"),
			Version:   3,
		}))
		require.Nil(t, d.SaveThing(ctx, models.Thing{
			Name:      "string1",
			CreatedAt: mustTime("2018-03-11T15:04:03+07:00"),
			Version:   2,
		}))
		limit := int64(3)
		tests := []getThingsByNameAndCreatedAtTest{
			{
				testName: "basic",
				d:        d,
				input: getThingsByNameAndCreatedAtInput{
					ctx: context.Background(),
					input: db.GetThingsByNameAndCreatedAtInput{
						Name:  "string1",
						Limit: &limit,
					},
				},
				output: getThingsByNameAndCreatedAtOutput{
					things: []models.Thing{
						models.Thing{
							Name:      "string1",
							CreatedAt: mustTime("2018-03-11T15:04:01+07:00"),
							Version:   1,
						},
						models.Thing{
							Name:      "string1",
							CreatedAt: mustTime("2018-03-11T15:04:02+07:00"),
							Version:   3,
						},
						models.Thing{
							Name:      "string1",
							CreatedAt: mustTime("2018-03-11T15:04:03+07:00"),
							Version:   2,
						},
					},
					err: nil,
				},
			},
			{
				testName: "descending",
				d:        d,
				input: getThingsByNameAndCreatedAtInput{
					ctx: context.Background(),
					input: db.GetThingsByNameAndCreatedAtInput{
						Name:       "string1",
						Descending: true,
					},
				},
				output: getThingsByNameAndCreatedAtOutput{
					things: []models.Thing{
						models.Thing{
							Name:      "string1",
							CreatedAt: mustTime("2018-03-11T15:04:03+07:00"),
							Version:   2,
						},
						models.Thing{
							Name:      "string1",
							CreatedAt: mustTime("2018-03-11T15:04:02+07:00"),
							Version:   3,
						},
						models.Thing{
							Name:      "string1",
							CreatedAt: mustTime("2018-03-11T15:04:01+07:00"),
							Version:   1,
						},
					},
					err: nil,
				},
			},
			/* FAILING_TEST */
			/* {
			   testName: "starting after",
			   d:    d,
			   input: getThingsByNameAndCreatedAtInput{
			     ctx: context.Background(),
			     input: db.GetThingsByNameAndCreatedAtInput{
			         Name: "string1",
			       StartingAfter: &models.Thing{
			         Name:    "string1",
			         CreatedAt: mustTime("2018-03-11T15:04:01+07:00"),
			         Version:    1,
			       },
			     },
			   },
			   output: getThingsByNameAndCreatedAtOutput{
			     things: []models.Thing{
			       models.Thing{
			         Name:    "string1",
			         CreatedAt: mustTime("2018-03-11T15:04:02+07:00"),
			         Version:    3,
			       },
			       models.Thing{
			         Name:    "string1",
			         CreatedAt: mustTime("2018-03-11T15:04:03+07:00"),
			         Version:    2,
			       },
			     },
			     err: nil,
			   },
			 }, */
			{
				testName: "starting after descending",
				d:        d,
				input: getThingsByNameAndCreatedAtInput{
					ctx: context.Background(),
					input: db.GetThingsByNameAndCreatedAtInput{
						Name: "string1",
						StartingAfter: &models.Thing{
							Name:      "string1",
							CreatedAt: mustTime("2018-03-11T15:04:03+07:00"),
							Version:   2,
						},
						Descending: true,
					},
				},
				output: getThingsByNameAndCreatedAtOutput{
					things: []models.Thing{
						models.Thing{
							Name:      "string1",
							CreatedAt: mustTime("2018-03-11T15:04:02+07:00"),
							Version:   3,
						},
						models.Thing{
							Name:      "string1",
							CreatedAt: mustTime("2018-03-11T15:04:01+07:00"),
							Version:   1,
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting at",
				d:        d,
				input: getThingsByNameAndCreatedAtInput{
					ctx: context.Background(),
					input: db.GetThingsByNameAndCreatedAtInput{
						Name:                "string1",
						CreatedAtStartingAt: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
					},
				},
				output: getThingsByNameAndCreatedAtOutput{
					things: []models.Thing{
						models.Thing{
							Name:      "string1",
							CreatedAt: mustTime("2018-03-11T15:04:02+07:00"),
							Version:   3,
						},
						models.Thing{
							Name:      "string1",
							CreatedAt: mustTime("2018-03-11T15:04:03+07:00"),
							Version:   2,
						},
					},
					err: nil,
				},
			},
		}
		for _, test := range tests {
			t.Run(test.testName, test.run)
		}
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanThingsByNameAndCreatedAt(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThing(ctx, models.Thing{
			Name:      "string1",
			CreatedAt: mustTime("2018-03-11T15:04:01+07:00"),
			Version:   1,
		}))
		require.Nil(t, d.SaveThing(ctx, models.Thing{
			Name:      "string2",
			CreatedAt: mustTime("2018-03-11T15:04:02+07:00"),
			Version:   2,
		}))
		require.Nil(t, d.SaveThing(ctx, models.Thing{
			Name:      "string3",
			CreatedAt: mustTime("2018-03-11T15:04:03+07:00"),
			Version:   3,
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.Thing{
				models.Thing{
					Name:      "string1",
					CreatedAt: mustTime("2018-03-11T15:04:01+07:00"),
					Version:   1,
				},
				models.Thing{
					Name:      "string2",
					CreatedAt: mustTime("2018-03-11T15:04:02+07:00"),
					Version:   2,
				},
				models.Thing{
					Name:      "string3",
					CreatedAt: mustTime("2018-03-11T15:04:03+07:00"),
					Version:   3,
				},
			}
			// Consistent read must be disabled when scaning a GSI.
			scanInput := db.ScanThingsByNameAndCreatedAtInput{DisableConsistentRead: true}
			actual := []models.Thing{}
			err := d.ScanThingsByNameAndCreatedAt(ctx, scanInput, func(m *models.Thing, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		t.Run("starting after", func(t *testing.T) {
			// Scan for everything.
			allItems := []models.Thing{}
			// Consistent read must be disabled when scaning a GSI.
			scanInput := db.ScanThingsByNameAndCreatedAtInput{DisableConsistentRead: true}
			err := d.ScanThingsByNameAndCreatedAt(ctx, scanInput, func(m *models.Thing, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput = db.ScanThingsByNameAndCreatedAtInput{
				DisableConsistentRead: true,
				StartingAfter: &models.Thing{
					Name:      firstItem.Name,
					CreatedAt: firstItem.CreatedAt,
					Version:   firstItem.Version,
				},
			}
			actual := []models.Thing{}
			err = d.ScanThingsByNameAndCreatedAt(ctx, scanInput, func(m *models.Thing, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanThingsInput{
				Limit: &limit,
			}
			actual := []models.Thing{}
			err := d.ScanThings(ctx, scanInput, func(m *models.Thing, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

type getThingsByNameAndRangeNullableInput struct {
	ctx   context.Context
	input db.GetThingsByNameAndRangeNullableInput
}
type getThingsByNameAndRangeNullableOutput struct {
	things []models.Thing
	err    error
}
type getThingsByNameAndRangeNullableTest struct {
	testName string
	d        db.Interface
	input    getThingsByNameAndRangeNullableInput
	output   getThingsByNameAndRangeNullableOutput
}

func (g getThingsByNameAndRangeNullableTest) run(t *testing.T) {
	things := []models.Thing{}
	fn := func(m *models.Thing, lastThing bool) bool {
		things = append(things, *m)
		if lastThing {
			return false
		}
		return true
	}
	err := g.d.GetThingsByNameAndRangeNullable(g.input.ctx, g.input.input, fn)
	if err != nil {
		fmt.Println(err.Error())
	}
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.things, things)
}

func GetThingsByNameAndRangeNullable(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThing(ctx, models.Thing{
			Name:          "string1",
			RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			Version:       1,
		}))
		require.Nil(t, d.SaveThing(ctx, models.Thing{
			Name:          "string1",
			RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
			Version:       3,
		}))
		require.Nil(t, d.SaveThing(ctx, models.Thing{
			Name:          "string1",
			RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
			Version:       2,
		}))
		limit := int64(3)
		tests := []getThingsByNameAndRangeNullableTest{
			{
				testName: "basic",
				d:        d,
				input: getThingsByNameAndRangeNullableInput{
					ctx: context.Background(),
					input: db.GetThingsByNameAndRangeNullableInput{
						Name:  "string1",
						Limit: &limit,
					},
				},
				output: getThingsByNameAndRangeNullableOutput{
					things: []models.Thing{
						models.Thing{
							Name:          "string1",
							RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
							Version:       1,
						},
						models.Thing{
							Name:          "string1",
							RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
							Version:       3,
						},
						models.Thing{
							Name:          "string1",
							RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
							Version:       2,
						},
					},
					err: nil,
				},
			},
			{
				testName: "descending",
				d:        d,
				input: getThingsByNameAndRangeNullableInput{
					ctx: context.Background(),
					input: db.GetThingsByNameAndRangeNullableInput{
						Name:       "string1",
						Descending: true,
					},
				},
				output: getThingsByNameAndRangeNullableOutput{
					things: []models.Thing{
						models.Thing{
							Name:          "string1",
							RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
							Version:       2,
						},
						models.Thing{
							Name:          "string1",
							RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
							Version:       3,
						},
						models.Thing{
							Name:          "string1",
							RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
							Version:       1,
						},
					},
					err: nil,
				},
			},
			/* FAILING_TEST */
			/* {
			   testName: "starting after",
			   d:    d,
			   input: getThingsByNameAndRangeNullableInput{
			     ctx: context.Background(),
			     input: db.GetThingsByNameAndRangeNullableInput{
			         Name: "string1",
			       StartingAfter: &models.Thing{
			         Name:    "string1",
			         RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			         Version:    1,
			       },
			     },
			   },
			   output: getThingsByNameAndRangeNullableOutput{
			     things: []models.Thing{
			       models.Thing{
			         Name:    "string1",
			         RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
			         Version:    3,
			       },
			       models.Thing{
			         Name:    "string1",
			         RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
			         Version:    2,
			       },
			     },
			     err: nil,
			   },
			 }, */
			{
				testName: "starting after descending",
				d:        d,
				input: getThingsByNameAndRangeNullableInput{
					ctx: context.Background(),
					input: db.GetThingsByNameAndRangeNullableInput{
						Name: "string1",
						StartingAfter: &models.Thing{
							Name:          "string1",
							RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
							Version:       2,
						},
						Descending: true,
					},
				},
				output: getThingsByNameAndRangeNullableOutput{
					things: []models.Thing{
						models.Thing{
							Name:          "string1",
							RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
							Version:       3,
						},
						models.Thing{
							Name:          "string1",
							RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
							Version:       1,
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting at",
				d:        d,
				input: getThingsByNameAndRangeNullableInput{
					ctx: context.Background(),
					input: db.GetThingsByNameAndRangeNullableInput{
						Name:                    "string1",
						RangeNullableStartingAt: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
					},
				},
				output: getThingsByNameAndRangeNullableOutput{
					things: []models.Thing{
						models.Thing{
							Name:          "string1",
							RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
							Version:       3,
						},
						models.Thing{
							Name:          "string1",
							RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
							Version:       2,
						},
					},
					err: nil,
				},
			},
		}
		for _, test := range tests {
			t.Run(test.testName, test.run)
		}
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanThingsByNameAndRangeNullable(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThing(ctx, models.Thing{
			Name:          "string1",
			RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			Version:       1,
		}))
		require.Nil(t, d.SaveThing(ctx, models.Thing{
			Name:          "string2",
			RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
			Version:       2,
		}))
		require.Nil(t, d.SaveThing(ctx, models.Thing{
			Name:          "string3",
			RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
			Version:       3,
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.Thing{
				models.Thing{
					Name:          "string1",
					RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
					Version:       1,
				},
				models.Thing{
					Name:          "string2",
					RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
					Version:       2,
				},
				models.Thing{
					Name:          "string3",
					RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
					Version:       3,
				},
			}
			// Consistent read must be disabled when scaning a GSI.
			scanInput := db.ScanThingsByNameAndRangeNullableInput{DisableConsistentRead: true}
			actual := []models.Thing{}
			err := d.ScanThingsByNameAndRangeNullable(ctx, scanInput, func(m *models.Thing, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		t.Run("starting after", func(t *testing.T) {
			// Scan for everything.
			allItems := []models.Thing{}
			// Consistent read must be disabled when scaning a GSI.
			scanInput := db.ScanThingsByNameAndRangeNullableInput{DisableConsistentRead: true}
			err := d.ScanThingsByNameAndRangeNullable(ctx, scanInput, func(m *models.Thing, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput = db.ScanThingsByNameAndRangeNullableInput{
				DisableConsistentRead: true,
				StartingAfter: &models.Thing{
					Name:          firstItem.Name,
					RangeNullable: firstItem.RangeNullable,
					Version:       firstItem.Version,
				},
			}
			actual := []models.Thing{}
			err = d.ScanThingsByNameAndRangeNullable(ctx, scanInput, func(m *models.Thing, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanThingsInput{
				Limit: &limit,
			}
			actual := []models.Thing{}
			err := d.ScanThings(ctx, scanInput, func(m *models.Thing, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

type getThingsByHashNullableAndNameInput struct {
	ctx   context.Context
	input db.GetThingsByHashNullableAndNameInput
}
type getThingsByHashNullableAndNameOutput struct {
	things []models.Thing
	err    error
}
type getThingsByHashNullableAndNameTest struct {
	testName string
	d        db.Interface
	input    getThingsByHashNullableAndNameInput
	output   getThingsByHashNullableAndNameOutput
}

func (g getThingsByHashNullableAndNameTest) run(t *testing.T) {
	things := []models.Thing{}
	fn := func(m *models.Thing, lastThing bool) bool {
		things = append(things, *m)
		if lastThing {
			return false
		}
		return true
	}
	err := g.d.GetThingsByHashNullableAndName(g.input.ctx, g.input.input, fn)
	if err != nil {
		fmt.Println(err.Error())
	}
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.things, things)
}

func GetThingsByHashNullableAndName(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThing(ctx, models.Thing{
			HashNullable: db.String("string1"),
			Name:         "string1",
			Version:      1,
		}))
		require.Nil(t, d.SaveThing(ctx, models.Thing{
			HashNullable: db.String("string1"),
			Name:         "string2",
			Version:      3,
		}))
		require.Nil(t, d.SaveThing(ctx, models.Thing{
			HashNullable: db.String("string1"),
			Name:         "string3",
			Version:      2,
		}))
		limit := int64(3)
		tests := []getThingsByHashNullableAndNameTest{
			{
				testName: "basic",
				d:        d,
				input: getThingsByHashNullableAndNameInput{
					ctx: context.Background(),
					input: db.GetThingsByHashNullableAndNameInput{
						HashNullable: "string1",
						Limit:        &limit,
					},
				},
				output: getThingsByHashNullableAndNameOutput{
					things: []models.Thing{
						models.Thing{
							HashNullable: db.String("string1"),
							Name:         "string1",
							Version:      1,
						},
						models.Thing{
							HashNullable: db.String("string1"),
							Name:         "string2",
							Version:      3,
						},
						models.Thing{
							HashNullable: db.String("string1"),
							Name:         "string3",
							Version:      2,
						},
					},
					err: nil,
				},
			},
			{
				testName: "descending",
				d:        d,
				input: getThingsByHashNullableAndNameInput{
					ctx: context.Background(),
					input: db.GetThingsByHashNullableAndNameInput{
						HashNullable: "string1",
						Descending:   true,
					},
				},
				output: getThingsByHashNullableAndNameOutput{
					things: []models.Thing{
						models.Thing{
							HashNullable: db.String("string1"),
							Name:         "string3",
							Version:      2,
						},
						models.Thing{
							HashNullable: db.String("string1"),
							Name:         "string2",
							Version:      3,
						},
						models.Thing{
							HashNullable: db.String("string1"),
							Name:         "string1",
							Version:      1,
						},
					},
					err: nil,
				},
			},
			/* FAILING_TEST */
			/* {
			   testName: "starting after",
			   d:    d,
			   input: getThingsByHashNullableAndNameInput{
			     ctx: context.Background(),
			     input: db.GetThingsByHashNullableAndNameInput{
			         HashNullable: "string1",
			       StartingAfter: &models.Thing{
			         HashNullable:    db.String("string1"),
			         Name: "string1",
			         Version:    1,
			       },
			     },
			   },
			   output: getThingsByHashNullableAndNameOutput{
			     things: []models.Thing{
			       models.Thing{
			         HashNullable:    db.String("string1"),
			         Name: "string2",
			         Version:    3,
			       },
			       models.Thing{
			         HashNullable:    db.String("string1"),
			         Name: "string3",
			         Version:    2,
			       },
			     },
			     err: nil,
			   },
			 }, */
			{
				testName: "starting after descending",
				d:        d,
				input: getThingsByHashNullableAndNameInput{
					ctx: context.Background(),
					input: db.GetThingsByHashNullableAndNameInput{
						HashNullable: "string1",
						StartingAfter: &models.Thing{
							HashNullable: db.String("string1"),
							Name:         "string3",
							Version:      2,
						},
						Descending: true,
					},
				},
				output: getThingsByHashNullableAndNameOutput{
					things: []models.Thing{
						models.Thing{
							HashNullable: db.String("string1"),
							Name:         "string2",
							Version:      3,
						},
						models.Thing{
							HashNullable: db.String("string1"),
							Name:         "string1",
							Version:      1,
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting at",
				d:        d,
				input: getThingsByHashNullableAndNameInput{
					ctx: context.Background(),
					input: db.GetThingsByHashNullableAndNameInput{
						HashNullable:   "string1",
						NameStartingAt: db.String("string2"),
					},
				},
				output: getThingsByHashNullableAndNameOutput{
					things: []models.Thing{
						models.Thing{
							HashNullable: db.String("string1"),
							Name:         "string2",
							Version:      3,
						},
						models.Thing{
							HashNullable: db.String("string1"),
							Name:         "string3",
							Version:      2,
						},
					},
					err: nil,
				},
			},
		}
		for _, test := range tests {
			t.Run(test.testName, test.run)
		}
	}
}

func GetThingAllowingBatchWrites(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingAllowingBatchWrites{
			Name:    "string1",
			Version: 1,
		}
		require.Nil(t, s.SaveThingAllowingBatchWrites(ctx, m))
		m2, err := s.GetThingAllowingBatchWrites(ctx, m.Name, m.Version)
		require.Nil(t, err)
		require.Equal(t, m.Name, m2.Name)
		require.Equal(t, m.Version, m2.Version)

		_, err = s.GetThingAllowingBatchWrites(ctx, "string2", 2)
		require.NotNil(t, err)
		require.IsType(t, err, db.ErrThingAllowingBatchWritesNotFound{})
	}
}

type getThingAllowingBatchWritessByNameAndVersionInput struct {
	ctx   context.Context
	input db.GetThingAllowingBatchWritessByNameAndVersionInput
}
type getThingAllowingBatchWritessByNameAndVersionOutput struct {
	thingAllowingBatchWritess []models.ThingAllowingBatchWrites
	err                       error
}
type getThingAllowingBatchWritessByNameAndVersionTest struct {
	testName string
	d        db.Interface
	input    getThingAllowingBatchWritessByNameAndVersionInput
	output   getThingAllowingBatchWritessByNameAndVersionOutput
}

func (g getThingAllowingBatchWritessByNameAndVersionTest) run(t *testing.T) {
	thingAllowingBatchWritess := []models.ThingAllowingBatchWrites{}
	fn := func(m *models.ThingAllowingBatchWrites, lastThingAllowingBatchWrites bool) bool {
		thingAllowingBatchWritess = append(thingAllowingBatchWritess, *m)
		if lastThingAllowingBatchWrites {
			return false
		}
		return true
	}
	err := g.d.GetThingAllowingBatchWritessByNameAndVersion(g.input.ctx, g.input.input, fn)
	if err != nil {
		fmt.Println(err.Error())
	}
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.thingAllowingBatchWritess, thingAllowingBatchWritess)
}

func GetThingAllowingBatchWritessByNameAndVersion(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingAllowingBatchWrites(ctx, models.ThingAllowingBatchWrites{
			Name:    "string1",
			Version: 1,
		}))
		require.Nil(t, d.SaveThingAllowingBatchWrites(ctx, models.ThingAllowingBatchWrites{
			Name:    "string1",
			Version: 2,
		}))
		require.Nil(t, d.SaveThingAllowingBatchWrites(ctx, models.ThingAllowingBatchWrites{
			Name:    "string1",
			Version: 3,
		}))
		limit := int64(3)
		tests := []getThingAllowingBatchWritessByNameAndVersionTest{
			{
				testName: "basic",
				d:        d,
				input: getThingAllowingBatchWritessByNameAndVersionInput{
					ctx: context.Background(),
					input: db.GetThingAllowingBatchWritessByNameAndVersionInput{
						Name:  "string1",
						Limit: &limit,
					},
				},
				output: getThingAllowingBatchWritessByNameAndVersionOutput{
					thingAllowingBatchWritess: []models.ThingAllowingBatchWrites{
						models.ThingAllowingBatchWrites{
							Name:    "string1",
							Version: 1,
						},
						models.ThingAllowingBatchWrites{
							Name:    "string1",
							Version: 2,
						},
						models.ThingAllowingBatchWrites{
							Name:    "string1",
							Version: 3,
						},
					},
					err: nil,
				},
			},
			{
				testName: "descending",
				d:        d,
				input: getThingAllowingBatchWritessByNameAndVersionInput{
					ctx: context.Background(),
					input: db.GetThingAllowingBatchWritessByNameAndVersionInput{
						Name:       "string1",
						Descending: true,
					},
				},
				output: getThingAllowingBatchWritessByNameAndVersionOutput{
					thingAllowingBatchWritess: []models.ThingAllowingBatchWrites{
						models.ThingAllowingBatchWrites{
							Name:    "string1",
							Version: 3,
						},
						models.ThingAllowingBatchWrites{
							Name:    "string1",
							Version: 2,
						},
						models.ThingAllowingBatchWrites{
							Name:    "string1",
							Version: 1,
						},
					},
					err: nil,
				},
			},
			/* FAILING_TEST */
			/* {
			   testName: "starting after",
			   d:    d,
			   input: getThingAllowingBatchWritessByNameAndVersionInput{
			     ctx: context.Background(),
			     input: db.GetThingAllowingBatchWritessByNameAndVersionInput{
			         Name: "string1",
			       StartingAfter: &models.ThingAllowingBatchWrites{
			           Name:    "string1",
			           Version:    1,
			       },
			     },
			   },
			   output: getThingAllowingBatchWritessByNameAndVersionOutput{
			     thingAllowingBatchWritess: []models.ThingAllowingBatchWrites{
			       models.ThingAllowingBatchWrites{
			           Name:    "string1",
			           Version: 2,
			       },
			       models.ThingAllowingBatchWrites{
			           Name:    "string1",
			           Version: 3,
			       },
			     },
			     err: nil,
			   },
			 }, */
			{
				testName: "starting after descending",
				d:        d,
				input: getThingAllowingBatchWritessByNameAndVersionInput{
					ctx: context.Background(),
					input: db.GetThingAllowingBatchWritessByNameAndVersionInput{
						Name: "string1",
						StartingAfter: &models.ThingAllowingBatchWrites{
							Name:    "string1",
							Version: 3,
						},
						Descending: true,
					},
				},
				output: getThingAllowingBatchWritessByNameAndVersionOutput{
					thingAllowingBatchWritess: []models.ThingAllowingBatchWrites{
						models.ThingAllowingBatchWrites{
							Name:    "string1",
							Version: 2,
						},
						models.ThingAllowingBatchWrites{
							Name:    "string1",
							Version: 1,
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting at",
				d:        d,
				input: getThingAllowingBatchWritessByNameAndVersionInput{
					ctx: context.Background(),
					input: db.GetThingAllowingBatchWritessByNameAndVersionInput{
						Name:              "string1",
						VersionStartingAt: db.Int64(2),
					},
				},
				output: getThingAllowingBatchWritessByNameAndVersionOutput{
					thingAllowingBatchWritess: []models.ThingAllowingBatchWrites{
						models.ThingAllowingBatchWrites{
							Name:    "string1",
							Version: 2,
						},
						models.ThingAllowingBatchWrites{
							Name:    "string1",
							Version: 3,
						},
					},
					err: nil,
				},
			},
		}
		for _, test := range tests {
			t.Run(test.testName, test.run)
		}
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanThingAllowingBatchWritess(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingAllowingBatchWrites(ctx, models.ThingAllowingBatchWrites{
			Name:    "string1",
			Version: 1,
		}))
		require.Nil(t, d.SaveThingAllowingBatchWrites(ctx, models.ThingAllowingBatchWrites{
			Name:    "string2",
			Version: 2,
		}))
		require.Nil(t, d.SaveThingAllowingBatchWrites(ctx, models.ThingAllowingBatchWrites{
			Name:    "string3",
			Version: 3,
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.ThingAllowingBatchWrites{
				models.ThingAllowingBatchWrites{
					Name:    "string1",
					Version: 1,
				},
				models.ThingAllowingBatchWrites{
					Name:    "string2",
					Version: 2,
				},
				models.ThingAllowingBatchWrites{
					Name:    "string3",
					Version: 3,
				},
			}
			actual := []models.ThingAllowingBatchWrites{}
			err := d.ScanThingAllowingBatchWritess(ctx, db.ScanThingAllowingBatchWritessInput{}, func(m *models.ThingAllowingBatchWrites, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		// FAILING_TEST
		t.Run("starting after", func(t *testing.T) {
			t.Skip()
			// Scan for everything.
			allItems := []models.ThingAllowingBatchWrites{}
			err := d.ScanThingAllowingBatchWritess(ctx, db.ScanThingAllowingBatchWritessInput{}, func(m *models.ThingAllowingBatchWrites, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput := db.ScanThingAllowingBatchWritessInput{
				StartingAfter: &models.ThingAllowingBatchWrites{
					Name:    firstItem.Name,
					Version: firstItem.Version,
				},
			}
			actual := []models.ThingAllowingBatchWrites{}
			err = d.ScanThingAllowingBatchWritess(ctx, scanInput, func(m *models.ThingAllowingBatchWrites, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanThingAllowingBatchWritessInput{
				Limit: &limit,
			}
			actual := []models.ThingAllowingBatchWrites{}
			err := d.ScanThingAllowingBatchWritess(ctx, scanInput, func(m *models.ThingAllowingBatchWrites, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

func SaveThingAllowingBatchWrites(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingAllowingBatchWrites{
			Name:    "string1",
			Version: 1,
		}
		require.Nil(t, s.SaveThingAllowingBatchWrites(ctx, m))
		require.IsType(t, db.ErrThingAllowingBatchWritesAlreadyExists{}, s.SaveThingAllowingBatchWrites(ctx, m))
	}
}

func DeleteThingAllowingBatchWrites(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingAllowingBatchWrites{
			Name:    "string1",
			Version: 1,
		}
		require.Nil(t, s.SaveThingAllowingBatchWrites(ctx, m))
		require.Nil(t, s.DeleteThingAllowingBatchWrites(ctx, m.Name, m.Version))
	}
}

func GetThingAllowingBatchWritesWithCompositeAttributes(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingAllowingBatchWritesWithCompositeAttributes{
			Date: db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			ID:   db.String("string1"),
			Name: db.String("string1"),
		}
		require.Nil(t, s.SaveThingAllowingBatchWritesWithCompositeAttributes(ctx, m))
		m2, err := s.GetThingAllowingBatchWritesWithCompositeAttributes(ctx, *m.Name, *m.ID, *m.Date)
		require.Nil(t, err)
		require.Equal(t, *m.Name, *m2.Name)
		require.Equal(t, *m.ID, *m2.ID)
		require.Equal(t, m.Date.String(), m2.Date.String())

		_, err = s.GetThingAllowingBatchWritesWithCompositeAttributes(ctx, "string2", "string2", mustTime("2018-03-11T15:04:02+07:00"))
		require.NotNil(t, err)
		require.IsType(t, err, db.ErrThingAllowingBatchWritesWithCompositeAttributesNotFound{})
	}
}

type getThingAllowingBatchWritesWithCompositeAttributessByNameIDAndDateInput struct {
	ctx   context.Context
	input db.GetThingAllowingBatchWritesWithCompositeAttributessByNameIDAndDateInput
}
type getThingAllowingBatchWritesWithCompositeAttributessByNameIDAndDateOutput struct {
	thingAllowingBatchWritesWithCompositeAttributess []models.ThingAllowingBatchWritesWithCompositeAttributes
	err                                              error
}
type getThingAllowingBatchWritesWithCompositeAttributessByNameIDAndDateTest struct {
	testName string
	d        db.Interface
	input    getThingAllowingBatchWritesWithCompositeAttributessByNameIDAndDateInput
	output   getThingAllowingBatchWritesWithCompositeAttributessByNameIDAndDateOutput
}

func (g getThingAllowingBatchWritesWithCompositeAttributessByNameIDAndDateTest) run(t *testing.T) {
	thingAllowingBatchWritesWithCompositeAttributess := []models.ThingAllowingBatchWritesWithCompositeAttributes{}
	fn := func(m *models.ThingAllowingBatchWritesWithCompositeAttributes, lastThingAllowingBatchWritesWithCompositeAttributes bool) bool {
		thingAllowingBatchWritesWithCompositeAttributess = append(thingAllowingBatchWritesWithCompositeAttributess, *m)
		if lastThingAllowingBatchWritesWithCompositeAttributes {
			return false
		}
		return true
	}
	err := g.d.GetThingAllowingBatchWritesWithCompositeAttributessByNameIDAndDate(g.input.ctx, g.input.input, fn)
	if err != nil {
		fmt.Println(err.Error())
	}
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.thingAllowingBatchWritesWithCompositeAttributess, thingAllowingBatchWritesWithCompositeAttributess)
}

func GetThingAllowingBatchWritesWithCompositeAttributessByNameIDAndDate(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingAllowingBatchWritesWithCompositeAttributes(ctx, models.ThingAllowingBatchWritesWithCompositeAttributes{
			Name: db.String("string1"),
			ID:   db.String("string1"),
			Date: db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
		}))
		require.Nil(t, d.SaveThingAllowingBatchWritesWithCompositeAttributes(ctx, models.ThingAllowingBatchWritesWithCompositeAttributes{
			Name: db.String("string1"),
			ID:   db.String("string1"),
			Date: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
		}))
		require.Nil(t, d.SaveThingAllowingBatchWritesWithCompositeAttributes(ctx, models.ThingAllowingBatchWritesWithCompositeAttributes{
			Name: db.String("string1"),
			ID:   db.String("string1"),
			Date: db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
		}))
		limit := int64(3)
		tests := []getThingAllowingBatchWritesWithCompositeAttributessByNameIDAndDateTest{
			{
				testName: "basic",
				d:        d,
				input: getThingAllowingBatchWritesWithCompositeAttributessByNameIDAndDateInput{
					ctx: context.Background(),
					input: db.GetThingAllowingBatchWritesWithCompositeAttributessByNameIDAndDateInput{
						Name:  "string1",
						ID:    "string1",
						Limit: &limit,
					},
				},
				output: getThingAllowingBatchWritesWithCompositeAttributessByNameIDAndDateOutput{
					thingAllowingBatchWritesWithCompositeAttributess: []models.ThingAllowingBatchWritesWithCompositeAttributes{
						models.ThingAllowingBatchWritesWithCompositeAttributes{
							Name: db.String("string1"),
							ID:   db.String("string1"),
							Date: db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
						},
						models.ThingAllowingBatchWritesWithCompositeAttributes{
							Name: db.String("string1"),
							ID:   db.String("string1"),
							Date: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
						},
						models.ThingAllowingBatchWritesWithCompositeAttributes{
							Name: db.String("string1"),
							ID:   db.String("string1"),
							Date: db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
						},
					},
					err: nil,
				},
			},
			{
				testName: "descending",
				d:        d,
				input: getThingAllowingBatchWritesWithCompositeAttributessByNameIDAndDateInput{
					ctx: context.Background(),
					input: db.GetThingAllowingBatchWritesWithCompositeAttributessByNameIDAndDateInput{
						Name:       "string1",
						ID:         "string1",
						Descending: true,
					},
				},
				output: getThingAllowingBatchWritesWithCompositeAttributessByNameIDAndDateOutput{
					thingAllowingBatchWritesWithCompositeAttributess: []models.ThingAllowingBatchWritesWithCompositeAttributes{
						models.ThingAllowingBatchWritesWithCompositeAttributes{
							Name: db.String("string1"),
							ID:   db.String("string1"),
							Date: db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
						},
						models.ThingAllowingBatchWritesWithCompositeAttributes{
							Name: db.String("string1"),
							ID:   db.String("string1"),
							Date: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
						},
						models.ThingAllowingBatchWritesWithCompositeAttributes{
							Name: db.String("string1"),
							ID:   db.String("string1"),
							Date: db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
						},
					},
					err: nil,
				},
			},
			/* FAILING_TEST */
			/* {
			   testName: "starting after",
			   d:    d,
			   input: getThingAllowingBatchWritesWithCompositeAttributessByNameIDAndDateInput{
			     ctx: context.Background(),
			     input: db.GetThingAllowingBatchWritesWithCompositeAttributessByNameIDAndDateInput{
			         Name: "string1",
			         ID: "string1",
			       StartingAfter: &models.ThingAllowingBatchWritesWithCompositeAttributes{
			           Name:    db.String("string1"),
			           ID:    db.String("string1"),
			           Date:    db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			       },
			     },
			   },
			   output: getThingAllowingBatchWritesWithCompositeAttributessByNameIDAndDateOutput{
			     thingAllowingBatchWritesWithCompositeAttributess: []models.ThingAllowingBatchWritesWithCompositeAttributes{
			       models.ThingAllowingBatchWritesWithCompositeAttributes{
			           Name:    db.String("string1"),
			           ID:    db.String("string1"),
			           Date: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
			       },
			       models.ThingAllowingBatchWritesWithCompositeAttributes{
			           Name:    db.String("string1"),
			           ID:    db.String("string1"),
			           Date: db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
			       },
			     },
			     err: nil,
			   },
			 }, */
			{
				testName: "starting after descending",
				d:        d,
				input: getThingAllowingBatchWritesWithCompositeAttributessByNameIDAndDateInput{
					ctx: context.Background(),
					input: db.GetThingAllowingBatchWritesWithCompositeAttributessByNameIDAndDateInput{
						Name: "string1",
						ID:   "string1",
						StartingAfter: &models.ThingAllowingBatchWritesWithCompositeAttributes{
							Name: db.String("string1"),
							ID:   db.String("string1"),
							Date: db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
						},
						Descending: true,
					},
				},
				output: getThingAllowingBatchWritesWithCompositeAttributessByNameIDAndDateOutput{
					thingAllowingBatchWritesWithCompositeAttributess: []models.ThingAllowingBatchWritesWithCompositeAttributes{
						models.ThingAllowingBatchWritesWithCompositeAttributes{
							Name: db.String("string1"),
							ID:   db.String("string1"),
							Date: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
						},
						models.ThingAllowingBatchWritesWithCompositeAttributes{
							Name: db.String("string1"),
							ID:   db.String("string1"),
							Date: db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting at",
				d:        d,
				input: getThingAllowingBatchWritesWithCompositeAttributessByNameIDAndDateInput{
					ctx: context.Background(),
					input: db.GetThingAllowingBatchWritesWithCompositeAttributessByNameIDAndDateInput{
						Name:           "string1",
						ID:             "string1",
						DateStartingAt: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
					},
				},
				output: getThingAllowingBatchWritesWithCompositeAttributessByNameIDAndDateOutput{
					thingAllowingBatchWritesWithCompositeAttributess: []models.ThingAllowingBatchWritesWithCompositeAttributes{
						models.ThingAllowingBatchWritesWithCompositeAttributes{
							Name: db.String("string1"),
							ID:   db.String("string1"),
							Date: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
						},
						models.ThingAllowingBatchWritesWithCompositeAttributes{
							Name: db.String("string1"),
							ID:   db.String("string1"),
							Date: db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
						},
					},
					err: nil,
				},
			},
		}
		for _, test := range tests {
			t.Run(test.testName, test.run)
		}
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanThingAllowingBatchWritesWithCompositeAttributess(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingAllowingBatchWritesWithCompositeAttributes(ctx, models.ThingAllowingBatchWritesWithCompositeAttributes{
			Date: db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			ID:   db.String("string1"),
			Name: db.String("string1"),
		}))
		require.Nil(t, d.SaveThingAllowingBatchWritesWithCompositeAttributes(ctx, models.ThingAllowingBatchWritesWithCompositeAttributes{
			Date: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
			ID:   db.String("string2"),
			Name: db.String("string2"),
		}))
		require.Nil(t, d.SaveThingAllowingBatchWritesWithCompositeAttributes(ctx, models.ThingAllowingBatchWritesWithCompositeAttributes{
			Date: db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
			ID:   db.String("string3"),
			Name: db.String("string3"),
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.ThingAllowingBatchWritesWithCompositeAttributes{
				models.ThingAllowingBatchWritesWithCompositeAttributes{
					Date: db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
					ID:   db.String("string1"),
					Name: db.String("string1"),
				},
				models.ThingAllowingBatchWritesWithCompositeAttributes{
					Date: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
					ID:   db.String("string2"),
					Name: db.String("string2"),
				},
				models.ThingAllowingBatchWritesWithCompositeAttributes{
					Date: db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
					ID:   db.String("string3"),
					Name: db.String("string3"),
				},
			}
			actual := []models.ThingAllowingBatchWritesWithCompositeAttributes{}
			err := d.ScanThingAllowingBatchWritesWithCompositeAttributess(ctx, db.ScanThingAllowingBatchWritesWithCompositeAttributessInput{}, func(m *models.ThingAllowingBatchWritesWithCompositeAttributes, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		// FAILING_TEST
		t.Run("starting after", func(t *testing.T) {
			t.Skip()
			// Scan for everything.
			allItems := []models.ThingAllowingBatchWritesWithCompositeAttributes{}
			err := d.ScanThingAllowingBatchWritesWithCompositeAttributess(ctx, db.ScanThingAllowingBatchWritesWithCompositeAttributessInput{}, func(m *models.ThingAllowingBatchWritesWithCompositeAttributes, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput := db.ScanThingAllowingBatchWritesWithCompositeAttributessInput{
				StartingAfter: &models.ThingAllowingBatchWritesWithCompositeAttributes{
					Name: firstItem.Name,
					ID:   firstItem.ID,
					Date: firstItem.Date,
				},
			}
			actual := []models.ThingAllowingBatchWritesWithCompositeAttributes{}
			err = d.ScanThingAllowingBatchWritesWithCompositeAttributess(ctx, scanInput, func(m *models.ThingAllowingBatchWritesWithCompositeAttributes, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanThingAllowingBatchWritesWithCompositeAttributessInput{
				Limit: &limit,
			}
			actual := []models.ThingAllowingBatchWritesWithCompositeAttributes{}
			err := d.ScanThingAllowingBatchWritesWithCompositeAttributess(ctx, scanInput, func(m *models.ThingAllowingBatchWritesWithCompositeAttributes, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

func SaveThingAllowingBatchWritesWithCompositeAttributes(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingAllowingBatchWritesWithCompositeAttributes{
			Date: db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			ID:   db.String("string1"),
			Name: db.String("string1"),
		}
		require.Nil(t, s.SaveThingAllowingBatchWritesWithCompositeAttributes(ctx, m))
		require.IsType(t, db.ErrThingAllowingBatchWritesWithCompositeAttributesAlreadyExists{}, s.SaveThingAllowingBatchWritesWithCompositeAttributes(ctx, m))
	}
}

func DeleteThingAllowingBatchWritesWithCompositeAttributes(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingAllowingBatchWritesWithCompositeAttributes{
			Date: db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			ID:   db.String("string1"),
			Name: db.String("string1"),
		}
		require.Nil(t, s.SaveThingAllowingBatchWritesWithCompositeAttributes(ctx, m))
		require.Nil(t, s.DeleteThingAllowingBatchWritesWithCompositeAttributes(ctx, *m.Name, *m.ID, *m.Date))
	}
}

func GetThingWithAdditionalAttributes(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithAdditionalAttributes{
			AdditionalBAttribute: []byte("string1"),
			AdditionalNAttribute: db.Int64(1),
			AdditionalSAttribute: db.String("string1"),
			CreatedAt:            mustTime("2018-03-11T15:04:01+07:00"),
			HashNullable:         db.String("string1"),
			ID:                   "string1",
			Name:                 "string1",
			RangeNullable:        db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			Version:              1,
		}
		require.Nil(t, s.SaveThingWithAdditionalAttributes(ctx, m))
		m2, err := s.GetThingWithAdditionalAttributes(ctx, m.Name, m.Version)
		require.Nil(t, err)
		require.Equal(t, m.Name, m2.Name)
		require.Equal(t, m.Version, m2.Version)

		_, err = s.GetThingWithAdditionalAttributes(ctx, "string2", 2)
		require.NotNil(t, err)
		require.IsType(t, err, db.ErrThingWithAdditionalAttributesNotFound{})
	}
}

type getThingWithAdditionalAttributessByNameAndVersionInput struct {
	ctx   context.Context
	input db.GetThingWithAdditionalAttributessByNameAndVersionInput
}
type getThingWithAdditionalAttributessByNameAndVersionOutput struct {
	thingWithAdditionalAttributess []models.ThingWithAdditionalAttributes
	err                            error
}
type getThingWithAdditionalAttributessByNameAndVersionTest struct {
	testName string
	d        db.Interface
	input    getThingWithAdditionalAttributessByNameAndVersionInput
	output   getThingWithAdditionalAttributessByNameAndVersionOutput
}

func (g getThingWithAdditionalAttributessByNameAndVersionTest) run(t *testing.T) {
	thingWithAdditionalAttributess := []models.ThingWithAdditionalAttributes{}
	fn := func(m *models.ThingWithAdditionalAttributes, lastThingWithAdditionalAttributes bool) bool {
		thingWithAdditionalAttributess = append(thingWithAdditionalAttributess, *m)
		if lastThingWithAdditionalAttributes {
			return false
		}
		return true
	}
	err := g.d.GetThingWithAdditionalAttributessByNameAndVersion(g.input.ctx, g.input.input, fn)
	if err != nil {
		fmt.Println(err.Error())
	}
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.thingWithAdditionalAttributess, thingWithAdditionalAttributess)
}

func GetThingWithAdditionalAttributessByNameAndVersion(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithAdditionalAttributes(ctx, models.ThingWithAdditionalAttributes{
			Name:                 "string1",
			Version:              1,
			AdditionalSAttribute: pointerToString("additionalSAttribute0"),
		}))
		require.Nil(t, d.SaveThingWithAdditionalAttributes(ctx, models.ThingWithAdditionalAttributes{
			Name:                 "string1",
			Version:              2,
			AdditionalSAttribute: pointerToString("additionalSAttribute1"),
		}))
		require.Nil(t, d.SaveThingWithAdditionalAttributes(ctx, models.ThingWithAdditionalAttributes{
			Name:                 "string1",
			Version:              3,
			AdditionalSAttribute: pointerToString("additionalSAttribute2"),
		}))
		limit := int64(3)
		tests := []getThingWithAdditionalAttributessByNameAndVersionTest{
			{
				testName: "basic",
				d:        d,
				input: getThingWithAdditionalAttributessByNameAndVersionInput{
					ctx: context.Background(),
					input: db.GetThingWithAdditionalAttributessByNameAndVersionInput{
						Name:  "string1",
						Limit: &limit,
					},
				},
				output: getThingWithAdditionalAttributessByNameAndVersionOutput{
					thingWithAdditionalAttributess: []models.ThingWithAdditionalAttributes{
						models.ThingWithAdditionalAttributes{
							Name:                 "string1",
							Version:              1,
							AdditionalSAttribute: pointerToString("additionalSAttribute0"),
						},
						models.ThingWithAdditionalAttributes{
							Name:                 "string1",
							Version:              2,
							AdditionalSAttribute: pointerToString("additionalSAttribute1"),
						},
						models.ThingWithAdditionalAttributes{
							Name:                 "string1",
							Version:              3,
							AdditionalSAttribute: pointerToString("additionalSAttribute2"),
						},
					},
					err: nil,
				},
			},
			{
				testName: "descending",
				d:        d,
				input: getThingWithAdditionalAttributessByNameAndVersionInput{
					ctx: context.Background(),
					input: db.GetThingWithAdditionalAttributessByNameAndVersionInput{
						Name:       "string1",
						Descending: true,
					},
				},
				output: getThingWithAdditionalAttributessByNameAndVersionOutput{
					thingWithAdditionalAttributess: []models.ThingWithAdditionalAttributes{
						models.ThingWithAdditionalAttributes{
							Name:                 "string1",
							Version:              3,
							AdditionalSAttribute: pointerToString("additionalSAttribute2"),
						},
						models.ThingWithAdditionalAttributes{
							Name:                 "string1",
							Version:              2,
							AdditionalSAttribute: pointerToString("additionalSAttribute1"),
						},
						models.ThingWithAdditionalAttributes{
							Name:                 "string1",
							Version:              1,
							AdditionalSAttribute: pointerToString("additionalSAttribute0"),
						},
					},
					err: nil,
				},
			},
			/* FAILING_TEST */
			/* {
			   testName: "starting after",
			   d:    d,
			   input: getThingWithAdditionalAttributessByNameAndVersionInput{
			     ctx: context.Background(),
			     input: db.GetThingWithAdditionalAttributessByNameAndVersionInput{
			         Name: "string1",
			       StartingAfter: &models.ThingWithAdditionalAttributes{
			           Name:    "string1",
			           Version:    1,
			               AdditionalSAttribute: pointerToString("additionalSAttribute0"),
			       },
			     },
			   },
			   output: getThingWithAdditionalAttributessByNameAndVersionOutput{
			     thingWithAdditionalAttributess: []models.ThingWithAdditionalAttributes{
			       models.ThingWithAdditionalAttributes{
			           Name:    "string1",
			           Version: 2,
			               AdditionalSAttribute:pointerToString("additionalSAttribute1"),
			       },
			       models.ThingWithAdditionalAttributes{
			           Name:    "string1",
			           Version: 3,
			               AdditionalSAttribute: pointerToString("additionalSAttribute2"),
			       },
			     },
			     err: nil,
			   },
			 }, */
			{
				testName: "starting after descending",
				d:        d,
				input: getThingWithAdditionalAttributessByNameAndVersionInput{
					ctx: context.Background(),
					input: db.GetThingWithAdditionalAttributessByNameAndVersionInput{
						Name: "string1",
						StartingAfter: &models.ThingWithAdditionalAttributes{
							Name:                 "string1",
							Version:              3,
							AdditionalSAttribute: pointerToString("additionalSAttribute2"),
						},
						Descending: true,
					},
				},
				output: getThingWithAdditionalAttributessByNameAndVersionOutput{
					thingWithAdditionalAttributess: []models.ThingWithAdditionalAttributes{
						models.ThingWithAdditionalAttributes{
							Name:                 "string1",
							Version:              2,
							AdditionalSAttribute: pointerToString("additionalSAttribute1"),
						},
						models.ThingWithAdditionalAttributes{
							Name:                 "string1",
							Version:              1,
							AdditionalSAttribute: pointerToString("additionalSAttribute0"),
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting at",
				d:        d,
				input: getThingWithAdditionalAttributessByNameAndVersionInput{
					ctx: context.Background(),
					input: db.GetThingWithAdditionalAttributessByNameAndVersionInput{
						Name:              "string1",
						VersionStartingAt: db.Int64(2),
					},
				},
				output: getThingWithAdditionalAttributessByNameAndVersionOutput{
					thingWithAdditionalAttributess: []models.ThingWithAdditionalAttributes{
						models.ThingWithAdditionalAttributes{
							Name:                 "string1",
							Version:              2,
							AdditionalSAttribute: pointerToString("additionalSAttribute1"),
						},
						models.ThingWithAdditionalAttributes{
							Name:                 "string1",
							Version:              3,
							AdditionalSAttribute: pointerToString("additionalSAttribute2"),
						},
					},
					err: nil,
				},
			},
			{
				testName: "filtering",
				d:        d,
				input: getThingWithAdditionalAttributessByNameAndVersionInput{
					ctx: context.Background(),
					input: db.GetThingWithAdditionalAttributessByNameAndVersionInput{
						Name: "string1",
						FilterValues: []db.ThingWithAdditionalAttributesByNameAndVersionFilterValues{
							db.ThingWithAdditionalAttributesByNameAndVersionFilterValues{
								AttributeName:   db.ThingWithAdditionalAttributesAdditionalSAttribute,
								AttributeValues: []interface{}{"additionalSAttribute0"},
							},
						},
						FilterExpression: "#ADDITIONALSATTRIBUTE = :additionalSAttribute_value0",
						Limit:            &limit,
					},
				},
				output: getThingWithAdditionalAttributessByNameAndVersionOutput{
					thingWithAdditionalAttributess: []models.ThingWithAdditionalAttributes{
						models.ThingWithAdditionalAttributes{
							Name:                 "string1",
							Version:              1,
							AdditionalSAttribute: pointerToString("additionalSAttribute0"),
						},
					},
					err: nil,
				},
			},
		}
		for _, test := range tests {
			t.Run(test.testName, test.run)
		}
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanThingWithAdditionalAttributess(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithAdditionalAttributes(ctx, models.ThingWithAdditionalAttributes{
			AdditionalBAttribute: []byte("string1"),
			AdditionalNAttribute: db.Int64(1),
			AdditionalSAttribute: db.String("string1"),
			CreatedAt:            mustTime("2018-03-11T15:04:01+07:00"),
			HashNullable:         db.String("string1"),
			ID:                   "string1",
			Name:                 "string1",
			RangeNullable:        db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			Version:              1,
		}))
		require.Nil(t, d.SaveThingWithAdditionalAttributes(ctx, models.ThingWithAdditionalAttributes{
			AdditionalBAttribute: []byte("string2"),
			AdditionalNAttribute: db.Int64(2),
			AdditionalSAttribute: db.String("string2"),
			CreatedAt:            mustTime("2018-03-11T15:04:02+07:00"),
			HashNullable:         db.String("string2"),
			ID:                   "string2",
			Name:                 "string2",
			RangeNullable:        db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
			Version:              2,
		}))
		require.Nil(t, d.SaveThingWithAdditionalAttributes(ctx, models.ThingWithAdditionalAttributes{
			AdditionalBAttribute: []byte("string3"),
			AdditionalNAttribute: db.Int64(3),
			AdditionalSAttribute: db.String("string3"),
			CreatedAt:            mustTime("2018-03-11T15:04:03+07:00"),
			HashNullable:         db.String("string3"),
			ID:                   "string3",
			Name:                 "string3",
			RangeNullable:        db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
			Version:              3,
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.ThingWithAdditionalAttributes{
				models.ThingWithAdditionalAttributes{
					AdditionalBAttribute: []byte("string1"),
					AdditionalNAttribute: db.Int64(1),
					AdditionalSAttribute: db.String("string1"),
					CreatedAt:            mustTime("2018-03-11T15:04:01+07:00"),
					HashNullable:         db.String("string1"),
					ID:                   "string1",
					Name:                 "string1",
					RangeNullable:        db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
					Version:              1,
				},
				models.ThingWithAdditionalAttributes{
					AdditionalBAttribute: []byte("string2"),
					AdditionalNAttribute: db.Int64(2),
					AdditionalSAttribute: db.String("string2"),
					CreatedAt:            mustTime("2018-03-11T15:04:02+07:00"),
					HashNullable:         db.String("string2"),
					ID:                   "string2",
					Name:                 "string2",
					RangeNullable:        db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
					Version:              2,
				},
				models.ThingWithAdditionalAttributes{
					AdditionalBAttribute: []byte("string3"),
					AdditionalNAttribute: db.Int64(3),
					AdditionalSAttribute: db.String("string3"),
					CreatedAt:            mustTime("2018-03-11T15:04:03+07:00"),
					HashNullable:         db.String("string3"),
					ID:                   "string3",
					Name:                 "string3",
					RangeNullable:        db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
					Version:              3,
				},
			}
			actual := []models.ThingWithAdditionalAttributes{}
			err := d.ScanThingWithAdditionalAttributess(ctx, db.ScanThingWithAdditionalAttributessInput{}, func(m *models.ThingWithAdditionalAttributes, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		// FAILING_TEST
		t.Run("starting after", func(t *testing.T) {
			t.Skip()
			// Scan for everything.
			allItems := []models.ThingWithAdditionalAttributes{}
			err := d.ScanThingWithAdditionalAttributess(ctx, db.ScanThingWithAdditionalAttributessInput{}, func(m *models.ThingWithAdditionalAttributes, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput := db.ScanThingWithAdditionalAttributessInput{
				StartingAfter: &models.ThingWithAdditionalAttributes{
					Name:    firstItem.Name,
					Version: firstItem.Version,
					// must specify non-empty string values for attributes
					// in secondary indexes, since dynamodb doesn't support
					// empty strings:
					AdditionalSAttribute: pointerToString("additionalSAttribute"),
				},
			}
			actual := []models.ThingWithAdditionalAttributes{}
			err = d.ScanThingWithAdditionalAttributess(ctx, scanInput, func(m *models.ThingWithAdditionalAttributes, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanThingWithAdditionalAttributessInput{
				Limit: &limit,
			}
			actual := []models.ThingWithAdditionalAttributes{}
			err := d.ScanThingWithAdditionalAttributess(ctx, scanInput, func(m *models.ThingWithAdditionalAttributes, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

func SaveThingWithAdditionalAttributes(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithAdditionalAttributes{
			AdditionalBAttribute: []byte("string1"),
			AdditionalNAttribute: db.Int64(1),
			AdditionalSAttribute: db.String("string1"),
			CreatedAt:            mustTime("2018-03-11T15:04:01+07:00"),
			HashNullable:         db.String("string1"),
			ID:                   "string1",
			Name:                 "string1",
			RangeNullable:        db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			Version:              1,
		}
		require.Nil(t, s.SaveThingWithAdditionalAttributes(ctx, m))
		require.IsType(t, db.ErrThingWithAdditionalAttributesAlreadyExists{}, s.SaveThingWithAdditionalAttributes(ctx, m))
	}
}

func DeleteThingWithAdditionalAttributes(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithAdditionalAttributes{
			AdditionalBAttribute: []byte("string1"),
			AdditionalNAttribute: db.Int64(1),
			AdditionalSAttribute: db.String("string1"),
			CreatedAt:            mustTime("2018-03-11T15:04:01+07:00"),
			HashNullable:         db.String("string1"),
			ID:                   "string1",
			Name:                 "string1",
			RangeNullable:        db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			Version:              1,
		}
		require.Nil(t, s.SaveThingWithAdditionalAttributes(ctx, m))
		require.Nil(t, s.DeleteThingWithAdditionalAttributes(ctx, m.Name, m.Version))
	}
}

func GetThingWithAdditionalAttributesByID(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithAdditionalAttributes{
			AdditionalBAttribute: []byte("string1"),
			AdditionalNAttribute: db.Int64(1),
			AdditionalSAttribute: db.String("string1"),
			CreatedAt:            mustTime("2018-03-11T15:04:01+07:00"),
			HashNullable:         db.String("string1"),
			ID:                   "string1",
			Name:                 "string1",
			RangeNullable:        db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			Version:              1,
		}
		require.Nil(t, s.SaveThingWithAdditionalAttributes(ctx, m))
		m2, err := s.GetThingWithAdditionalAttributesByID(ctx, m.ID)
		require.Nil(t, err)
		require.Equal(t, m.AdditionalBAttribute, m2.AdditionalBAttribute)
		require.Equal(t, m.AdditionalNAttribute, m2.AdditionalNAttribute)
		require.Equal(t, m.AdditionalSAttribute, m2.AdditionalSAttribute)
		require.Equal(t, m.CreatedAt.String(), m2.CreatedAt.String())
		require.Equal(t, m.HashNullable, m2.HashNullable)
		require.Equal(t, m.ID, m2.ID)
		require.Equal(t, m.Name, m2.Name)
		require.Equal(t, m.RangeNullable.String(), m2.RangeNullable.String())
		require.Equal(t, m.Version, m2.Version)

		_, err = s.GetThingWithAdditionalAttributesByID(ctx, "string2")
		require.NotNil(t, err)
		require.IsType(t, err, db.ErrThingWithAdditionalAttributesByIDNotFound{})
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanThingWithAdditionalAttributessByID(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithAdditionalAttributes(ctx, models.ThingWithAdditionalAttributes{
			ID:      "string1",
			Name:    "string1",
			Version: 1,
		}))
		require.Nil(t, d.SaveThingWithAdditionalAttributes(ctx, models.ThingWithAdditionalAttributes{
			ID:      "string2",
			Name:    "string2",
			Version: 2,
		}))
		require.Nil(t, d.SaveThingWithAdditionalAttributes(ctx, models.ThingWithAdditionalAttributes{
			ID:      "string3",
			Name:    "string3",
			Version: 3,
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.ThingWithAdditionalAttributes{
				models.ThingWithAdditionalAttributes{
					ID:      "string1",
					Name:    "string1",
					Version: 1,
				},
				models.ThingWithAdditionalAttributes{
					ID:      "string2",
					Name:    "string2",
					Version: 2,
				},
				models.ThingWithAdditionalAttributes{
					ID:      "string3",
					Name:    "string3",
					Version: 3,
				},
			}
			// Consistent read must be disabled when scaning a GSI.
			scanInput := db.ScanThingWithAdditionalAttributessByIDInput{DisableConsistentRead: true}
			actual := []models.ThingWithAdditionalAttributes{}
			err := d.ScanThingWithAdditionalAttributessByID(ctx, scanInput, func(m *models.ThingWithAdditionalAttributes, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		t.Run("starting after", func(t *testing.T) {
			// Scan for everything.
			allItems := []models.ThingWithAdditionalAttributes{}
			// Consistent read must be disabled when scaning a GSI.
			scanInput := db.ScanThingWithAdditionalAttributessByIDInput{DisableConsistentRead: true}
			err := d.ScanThingWithAdditionalAttributessByID(ctx, scanInput, func(m *models.ThingWithAdditionalAttributes, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput = db.ScanThingWithAdditionalAttributessByIDInput{
				DisableConsistentRead: true,
				StartingAfter: &models.ThingWithAdditionalAttributes{
					ID:      firstItem.ID,
					Name:    firstItem.Name,
					Version: firstItem.Version,
				},
			}
			actual := []models.ThingWithAdditionalAttributes{}
			err = d.ScanThingWithAdditionalAttributessByID(ctx, scanInput, func(m *models.ThingWithAdditionalAttributes, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanThingWithAdditionalAttributessInput{
				Limit: &limit,
			}
			actual := []models.ThingWithAdditionalAttributes{}
			err := d.ScanThingWithAdditionalAttributess(ctx, scanInput, func(m *models.ThingWithAdditionalAttributes, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

type getThingWithAdditionalAttributessByNameAndCreatedAtInput struct {
	ctx   context.Context
	input db.GetThingWithAdditionalAttributessByNameAndCreatedAtInput
}
type getThingWithAdditionalAttributessByNameAndCreatedAtOutput struct {
	thingWithAdditionalAttributess []models.ThingWithAdditionalAttributes
	err                            error
}
type getThingWithAdditionalAttributessByNameAndCreatedAtTest struct {
	testName string
	d        db.Interface
	input    getThingWithAdditionalAttributessByNameAndCreatedAtInput
	output   getThingWithAdditionalAttributessByNameAndCreatedAtOutput
}

func (g getThingWithAdditionalAttributessByNameAndCreatedAtTest) run(t *testing.T) {
	thingWithAdditionalAttributess := []models.ThingWithAdditionalAttributes{}
	fn := func(m *models.ThingWithAdditionalAttributes, lastThingWithAdditionalAttributes bool) bool {
		thingWithAdditionalAttributess = append(thingWithAdditionalAttributess, *m)
		if lastThingWithAdditionalAttributes {
			return false
		}
		return true
	}
	err := g.d.GetThingWithAdditionalAttributessByNameAndCreatedAt(g.input.ctx, g.input.input, fn)
	if err != nil {
		fmt.Println(err.Error())
	}
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.thingWithAdditionalAttributess, thingWithAdditionalAttributess)
}

func GetThingWithAdditionalAttributessByNameAndCreatedAt(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithAdditionalAttributes(ctx, models.ThingWithAdditionalAttributes{
			Name:      "string1",
			CreatedAt: mustTime("2018-03-11T15:04:01+07:00"),
			Version:   1,
		}))
		require.Nil(t, d.SaveThingWithAdditionalAttributes(ctx, models.ThingWithAdditionalAttributes{
			Name:      "string1",
			CreatedAt: mustTime("2018-03-11T15:04:02+07:00"),
			Version:   3,
		}))
		require.Nil(t, d.SaveThingWithAdditionalAttributes(ctx, models.ThingWithAdditionalAttributes{
			Name:      "string1",
			CreatedAt: mustTime("2018-03-11T15:04:03+07:00"),
			Version:   2,
		}))
		limit := int64(3)
		tests := []getThingWithAdditionalAttributessByNameAndCreatedAtTest{
			{
				testName: "basic",
				d:        d,
				input: getThingWithAdditionalAttributessByNameAndCreatedAtInput{
					ctx: context.Background(),
					input: db.GetThingWithAdditionalAttributessByNameAndCreatedAtInput{
						Name:  "string1",
						Limit: &limit,
					},
				},
				output: getThingWithAdditionalAttributessByNameAndCreatedAtOutput{
					thingWithAdditionalAttributess: []models.ThingWithAdditionalAttributes{
						models.ThingWithAdditionalAttributes{
							Name:      "string1",
							CreatedAt: mustTime("2018-03-11T15:04:01+07:00"),
							Version:   1,
						},
						models.ThingWithAdditionalAttributes{
							Name:      "string1",
							CreatedAt: mustTime("2018-03-11T15:04:02+07:00"),
							Version:   3,
						},
						models.ThingWithAdditionalAttributes{
							Name:      "string1",
							CreatedAt: mustTime("2018-03-11T15:04:03+07:00"),
							Version:   2,
						},
					},
					err: nil,
				},
			},
			{
				testName: "descending",
				d:        d,
				input: getThingWithAdditionalAttributessByNameAndCreatedAtInput{
					ctx: context.Background(),
					input: db.GetThingWithAdditionalAttributessByNameAndCreatedAtInput{
						Name:       "string1",
						Descending: true,
					},
				},
				output: getThingWithAdditionalAttributessByNameAndCreatedAtOutput{
					thingWithAdditionalAttributess: []models.ThingWithAdditionalAttributes{
						models.ThingWithAdditionalAttributes{
							Name:      "string1",
							CreatedAt: mustTime("2018-03-11T15:04:03+07:00"),
							Version:   2,
						},
						models.ThingWithAdditionalAttributes{
							Name:      "string1",
							CreatedAt: mustTime("2018-03-11T15:04:02+07:00"),
							Version:   3,
						},
						models.ThingWithAdditionalAttributes{
							Name:      "string1",
							CreatedAt: mustTime("2018-03-11T15:04:01+07:00"),
							Version:   1,
						},
					},
					err: nil,
				},
			},
			/* FAILING_TEST */
			/* {
			   testName: "starting after",
			   d:    d,
			   input: getThingWithAdditionalAttributessByNameAndCreatedAtInput{
			     ctx: context.Background(),
			     input: db.GetThingWithAdditionalAttributessByNameAndCreatedAtInput{
			         Name: "string1",
			       StartingAfter: &models.ThingWithAdditionalAttributes{
			         Name:    "string1",
			         CreatedAt: mustTime("2018-03-11T15:04:01+07:00"),
			         Version:    1,
			       },
			     },
			   },
			   output: getThingWithAdditionalAttributessByNameAndCreatedAtOutput{
			     thingWithAdditionalAttributess: []models.ThingWithAdditionalAttributes{
			       models.ThingWithAdditionalAttributes{
			         Name:    "string1",
			         CreatedAt: mustTime("2018-03-11T15:04:02+07:00"),
			         Version:    3,
			       },
			       models.ThingWithAdditionalAttributes{
			         Name:    "string1",
			         CreatedAt: mustTime("2018-03-11T15:04:03+07:00"),
			         Version:    2,
			       },
			     },
			     err: nil,
			   },
			 }, */
			{
				testName: "starting after descending",
				d:        d,
				input: getThingWithAdditionalAttributessByNameAndCreatedAtInput{
					ctx: context.Background(),
					input: db.GetThingWithAdditionalAttributessByNameAndCreatedAtInput{
						Name: "string1",
						StartingAfter: &models.ThingWithAdditionalAttributes{
							Name:      "string1",
							CreatedAt: mustTime("2018-03-11T15:04:03+07:00"),
							Version:   2,
						},
						Descending: true,
					},
				},
				output: getThingWithAdditionalAttributessByNameAndCreatedAtOutput{
					thingWithAdditionalAttributess: []models.ThingWithAdditionalAttributes{
						models.ThingWithAdditionalAttributes{
							Name:      "string1",
							CreatedAt: mustTime("2018-03-11T15:04:02+07:00"),
							Version:   3,
						},
						models.ThingWithAdditionalAttributes{
							Name:      "string1",
							CreatedAt: mustTime("2018-03-11T15:04:01+07:00"),
							Version:   1,
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting at",
				d:        d,
				input: getThingWithAdditionalAttributessByNameAndCreatedAtInput{
					ctx: context.Background(),
					input: db.GetThingWithAdditionalAttributessByNameAndCreatedAtInput{
						Name:                "string1",
						CreatedAtStartingAt: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
					},
				},
				output: getThingWithAdditionalAttributessByNameAndCreatedAtOutput{
					thingWithAdditionalAttributess: []models.ThingWithAdditionalAttributes{
						models.ThingWithAdditionalAttributes{
							Name:      "string1",
							CreatedAt: mustTime("2018-03-11T15:04:02+07:00"),
							Version:   3,
						},
						models.ThingWithAdditionalAttributes{
							Name:      "string1",
							CreatedAt: mustTime("2018-03-11T15:04:03+07:00"),
							Version:   2,
						},
					},
					err: nil,
				},
			},
		}
		for _, test := range tests {
			t.Run(test.testName, test.run)
		}
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanThingWithAdditionalAttributessByNameAndCreatedAt(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithAdditionalAttributes(ctx, models.ThingWithAdditionalAttributes{
			Name:      "string1",
			CreatedAt: mustTime("2018-03-11T15:04:01+07:00"),
			Version:   1,
		}))
		require.Nil(t, d.SaveThingWithAdditionalAttributes(ctx, models.ThingWithAdditionalAttributes{
			Name:      "string2",
			CreatedAt: mustTime("2018-03-11T15:04:02+07:00"),
			Version:   2,
		}))
		require.Nil(t, d.SaveThingWithAdditionalAttributes(ctx, models.ThingWithAdditionalAttributes{
			Name:      "string3",
			CreatedAt: mustTime("2018-03-11T15:04:03+07:00"),
			Version:   3,
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.ThingWithAdditionalAttributes{
				models.ThingWithAdditionalAttributes{
					Name:      "string1",
					CreatedAt: mustTime("2018-03-11T15:04:01+07:00"),
					Version:   1,
				},
				models.ThingWithAdditionalAttributes{
					Name:      "string2",
					CreatedAt: mustTime("2018-03-11T15:04:02+07:00"),
					Version:   2,
				},
				models.ThingWithAdditionalAttributes{
					Name:      "string3",
					CreatedAt: mustTime("2018-03-11T15:04:03+07:00"),
					Version:   3,
				},
			}
			// Consistent read must be disabled when scaning a GSI.
			scanInput := db.ScanThingWithAdditionalAttributessByNameAndCreatedAtInput{DisableConsistentRead: true}
			actual := []models.ThingWithAdditionalAttributes{}
			err := d.ScanThingWithAdditionalAttributessByNameAndCreatedAt(ctx, scanInput, func(m *models.ThingWithAdditionalAttributes, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		t.Run("starting after", func(t *testing.T) {
			// Scan for everything.
			allItems := []models.ThingWithAdditionalAttributes{}
			// Consistent read must be disabled when scaning a GSI.
			scanInput := db.ScanThingWithAdditionalAttributessByNameAndCreatedAtInput{DisableConsistentRead: true}
			err := d.ScanThingWithAdditionalAttributessByNameAndCreatedAt(ctx, scanInput, func(m *models.ThingWithAdditionalAttributes, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput = db.ScanThingWithAdditionalAttributessByNameAndCreatedAtInput{
				DisableConsistentRead: true,
				StartingAfter: &models.ThingWithAdditionalAttributes{
					Name:      firstItem.Name,
					CreatedAt: firstItem.CreatedAt,
					Version:   firstItem.Version,
				},
			}
			actual := []models.ThingWithAdditionalAttributes{}
			err = d.ScanThingWithAdditionalAttributessByNameAndCreatedAt(ctx, scanInput, func(m *models.ThingWithAdditionalAttributes, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanThingWithAdditionalAttributessInput{
				Limit: &limit,
			}
			actual := []models.ThingWithAdditionalAttributes{}
			err := d.ScanThingWithAdditionalAttributess(ctx, scanInput, func(m *models.ThingWithAdditionalAttributes, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

type getThingWithAdditionalAttributessByNameAndRangeNullableInput struct {
	ctx   context.Context
	input db.GetThingWithAdditionalAttributessByNameAndRangeNullableInput
}
type getThingWithAdditionalAttributessByNameAndRangeNullableOutput struct {
	thingWithAdditionalAttributess []models.ThingWithAdditionalAttributes
	err                            error
}
type getThingWithAdditionalAttributessByNameAndRangeNullableTest struct {
	testName string
	d        db.Interface
	input    getThingWithAdditionalAttributessByNameAndRangeNullableInput
	output   getThingWithAdditionalAttributessByNameAndRangeNullableOutput
}

func (g getThingWithAdditionalAttributessByNameAndRangeNullableTest) run(t *testing.T) {
	thingWithAdditionalAttributess := []models.ThingWithAdditionalAttributes{}
	fn := func(m *models.ThingWithAdditionalAttributes, lastThingWithAdditionalAttributes bool) bool {
		thingWithAdditionalAttributess = append(thingWithAdditionalAttributess, *m)
		if lastThingWithAdditionalAttributes {
			return false
		}
		return true
	}
	err := g.d.GetThingWithAdditionalAttributessByNameAndRangeNullable(g.input.ctx, g.input.input, fn)
	if err != nil {
		fmt.Println(err.Error())
	}
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.thingWithAdditionalAttributess, thingWithAdditionalAttributess)
}

func GetThingWithAdditionalAttributessByNameAndRangeNullable(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithAdditionalAttributes(ctx, models.ThingWithAdditionalAttributes{
			Name:          "string1",
			RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			Version:       1,
		}))
		require.Nil(t, d.SaveThingWithAdditionalAttributes(ctx, models.ThingWithAdditionalAttributes{
			Name:          "string1",
			RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
			Version:       3,
		}))
		require.Nil(t, d.SaveThingWithAdditionalAttributes(ctx, models.ThingWithAdditionalAttributes{
			Name:          "string1",
			RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
			Version:       2,
		}))
		limit := int64(3)
		tests := []getThingWithAdditionalAttributessByNameAndRangeNullableTest{
			{
				testName: "basic",
				d:        d,
				input: getThingWithAdditionalAttributessByNameAndRangeNullableInput{
					ctx: context.Background(),
					input: db.GetThingWithAdditionalAttributessByNameAndRangeNullableInput{
						Name:  "string1",
						Limit: &limit,
					},
				},
				output: getThingWithAdditionalAttributessByNameAndRangeNullableOutput{
					thingWithAdditionalAttributess: []models.ThingWithAdditionalAttributes{
						models.ThingWithAdditionalAttributes{
							Name:          "string1",
							RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
							Version:       1,
						},
						models.ThingWithAdditionalAttributes{
							Name:          "string1",
							RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
							Version:       3,
						},
						models.ThingWithAdditionalAttributes{
							Name:          "string1",
							RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
							Version:       2,
						},
					},
					err: nil,
				},
			},
			{
				testName: "descending",
				d:        d,
				input: getThingWithAdditionalAttributessByNameAndRangeNullableInput{
					ctx: context.Background(),
					input: db.GetThingWithAdditionalAttributessByNameAndRangeNullableInput{
						Name:       "string1",
						Descending: true,
					},
				},
				output: getThingWithAdditionalAttributessByNameAndRangeNullableOutput{
					thingWithAdditionalAttributess: []models.ThingWithAdditionalAttributes{
						models.ThingWithAdditionalAttributes{
							Name:          "string1",
							RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
							Version:       2,
						},
						models.ThingWithAdditionalAttributes{
							Name:          "string1",
							RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
							Version:       3,
						},
						models.ThingWithAdditionalAttributes{
							Name:          "string1",
							RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
							Version:       1,
						},
					},
					err: nil,
				},
			},
			/* FAILING_TEST */
			/* {
			   testName: "starting after",
			   d:    d,
			   input: getThingWithAdditionalAttributessByNameAndRangeNullableInput{
			     ctx: context.Background(),
			     input: db.GetThingWithAdditionalAttributessByNameAndRangeNullableInput{
			         Name: "string1",
			       StartingAfter: &models.ThingWithAdditionalAttributes{
			         Name:    "string1",
			         RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			         Version:    1,
			       },
			     },
			   },
			   output: getThingWithAdditionalAttributessByNameAndRangeNullableOutput{
			     thingWithAdditionalAttributess: []models.ThingWithAdditionalAttributes{
			       models.ThingWithAdditionalAttributes{
			         Name:    "string1",
			         RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
			         Version:    3,
			       },
			       models.ThingWithAdditionalAttributes{
			         Name:    "string1",
			         RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
			         Version:    2,
			       },
			     },
			     err: nil,
			   },
			 }, */
			{
				testName: "starting after descending",
				d:        d,
				input: getThingWithAdditionalAttributessByNameAndRangeNullableInput{
					ctx: context.Background(),
					input: db.GetThingWithAdditionalAttributessByNameAndRangeNullableInput{
						Name: "string1",
						StartingAfter: &models.ThingWithAdditionalAttributes{
							Name:          "string1",
							RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
							Version:       2,
						},
						Descending: true,
					},
				},
				output: getThingWithAdditionalAttributessByNameAndRangeNullableOutput{
					thingWithAdditionalAttributess: []models.ThingWithAdditionalAttributes{
						models.ThingWithAdditionalAttributes{
							Name:          "string1",
							RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
							Version:       3,
						},
						models.ThingWithAdditionalAttributes{
							Name:          "string1",
							RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
							Version:       1,
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting at",
				d:        d,
				input: getThingWithAdditionalAttributessByNameAndRangeNullableInput{
					ctx: context.Background(),
					input: db.GetThingWithAdditionalAttributessByNameAndRangeNullableInput{
						Name:                    "string1",
						RangeNullableStartingAt: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
					},
				},
				output: getThingWithAdditionalAttributessByNameAndRangeNullableOutput{
					thingWithAdditionalAttributess: []models.ThingWithAdditionalAttributes{
						models.ThingWithAdditionalAttributes{
							Name:          "string1",
							RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
							Version:       3,
						},
						models.ThingWithAdditionalAttributes{
							Name:          "string1",
							RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
							Version:       2,
						},
					},
					err: nil,
				},
			},
		}
		for _, test := range tests {
			t.Run(test.testName, test.run)
		}
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanThingWithAdditionalAttributessByNameAndRangeNullable(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithAdditionalAttributes(ctx, models.ThingWithAdditionalAttributes{
			Name:          "string1",
			RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			Version:       1,
		}))
		require.Nil(t, d.SaveThingWithAdditionalAttributes(ctx, models.ThingWithAdditionalAttributes{
			Name:          "string2",
			RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
			Version:       2,
		}))
		require.Nil(t, d.SaveThingWithAdditionalAttributes(ctx, models.ThingWithAdditionalAttributes{
			Name:          "string3",
			RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
			Version:       3,
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.ThingWithAdditionalAttributes{
				models.ThingWithAdditionalAttributes{
					Name:          "string1",
					RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
					Version:       1,
				},
				models.ThingWithAdditionalAttributes{
					Name:          "string2",
					RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
					Version:       2,
				},
				models.ThingWithAdditionalAttributes{
					Name:          "string3",
					RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
					Version:       3,
				},
			}
			// Consistent read must be disabled when scaning a GSI.
			scanInput := db.ScanThingWithAdditionalAttributessByNameAndRangeNullableInput{DisableConsistentRead: true}
			actual := []models.ThingWithAdditionalAttributes{}
			err := d.ScanThingWithAdditionalAttributessByNameAndRangeNullable(ctx, scanInput, func(m *models.ThingWithAdditionalAttributes, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		t.Run("starting after", func(t *testing.T) {
			// Scan for everything.
			allItems := []models.ThingWithAdditionalAttributes{}
			// Consistent read must be disabled when scaning a GSI.
			scanInput := db.ScanThingWithAdditionalAttributessByNameAndRangeNullableInput{DisableConsistentRead: true}
			err := d.ScanThingWithAdditionalAttributessByNameAndRangeNullable(ctx, scanInput, func(m *models.ThingWithAdditionalAttributes, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput = db.ScanThingWithAdditionalAttributessByNameAndRangeNullableInput{
				DisableConsistentRead: true,
				StartingAfter: &models.ThingWithAdditionalAttributes{
					Name:          firstItem.Name,
					RangeNullable: firstItem.RangeNullable,
					Version:       firstItem.Version,
				},
			}
			actual := []models.ThingWithAdditionalAttributes{}
			err = d.ScanThingWithAdditionalAttributessByNameAndRangeNullable(ctx, scanInput, func(m *models.ThingWithAdditionalAttributes, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanThingWithAdditionalAttributessInput{
				Limit: &limit,
			}
			actual := []models.ThingWithAdditionalAttributes{}
			err := d.ScanThingWithAdditionalAttributess(ctx, scanInput, func(m *models.ThingWithAdditionalAttributes, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

type getThingWithAdditionalAttributessByHashNullableAndNameInput struct {
	ctx   context.Context
	input db.GetThingWithAdditionalAttributessByHashNullableAndNameInput
}
type getThingWithAdditionalAttributessByHashNullableAndNameOutput struct {
	thingWithAdditionalAttributess []models.ThingWithAdditionalAttributes
	err                            error
}
type getThingWithAdditionalAttributessByHashNullableAndNameTest struct {
	testName string
	d        db.Interface
	input    getThingWithAdditionalAttributessByHashNullableAndNameInput
	output   getThingWithAdditionalAttributessByHashNullableAndNameOutput
}

func (g getThingWithAdditionalAttributessByHashNullableAndNameTest) run(t *testing.T) {
	thingWithAdditionalAttributess := []models.ThingWithAdditionalAttributes{}
	fn := func(m *models.ThingWithAdditionalAttributes, lastThingWithAdditionalAttributes bool) bool {
		thingWithAdditionalAttributess = append(thingWithAdditionalAttributess, *m)
		if lastThingWithAdditionalAttributes {
			return false
		}
		return true
	}
	err := g.d.GetThingWithAdditionalAttributessByHashNullableAndName(g.input.ctx, g.input.input, fn)
	if err != nil {
		fmt.Println(err.Error())
	}
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.thingWithAdditionalAttributess, thingWithAdditionalAttributess)
}

func GetThingWithAdditionalAttributessByHashNullableAndName(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithAdditionalAttributes(ctx, models.ThingWithAdditionalAttributes{
			HashNullable: db.String("string1"),
			Name:         "string1",
			Version:      1,
		}))
		require.Nil(t, d.SaveThingWithAdditionalAttributes(ctx, models.ThingWithAdditionalAttributes{
			HashNullable: db.String("string1"),
			Name:         "string2",
			Version:      3,
		}))
		require.Nil(t, d.SaveThingWithAdditionalAttributes(ctx, models.ThingWithAdditionalAttributes{
			HashNullable: db.String("string1"),
			Name:         "string3",
			Version:      2,
		}))
		limit := int64(3)
		tests := []getThingWithAdditionalAttributessByHashNullableAndNameTest{
			{
				testName: "basic",
				d:        d,
				input: getThingWithAdditionalAttributessByHashNullableAndNameInput{
					ctx: context.Background(),
					input: db.GetThingWithAdditionalAttributessByHashNullableAndNameInput{
						HashNullable: "string1",
						Limit:        &limit,
					},
				},
				output: getThingWithAdditionalAttributessByHashNullableAndNameOutput{
					thingWithAdditionalAttributess: []models.ThingWithAdditionalAttributes{
						models.ThingWithAdditionalAttributes{
							HashNullable: db.String("string1"),
							Name:         "string1",
							Version:      1,
						},
						models.ThingWithAdditionalAttributes{
							HashNullable: db.String("string1"),
							Name:         "string2",
							Version:      3,
						},
						models.ThingWithAdditionalAttributes{
							HashNullable: db.String("string1"),
							Name:         "string3",
							Version:      2,
						},
					},
					err: nil,
				},
			},
			{
				testName: "descending",
				d:        d,
				input: getThingWithAdditionalAttributessByHashNullableAndNameInput{
					ctx: context.Background(),
					input: db.GetThingWithAdditionalAttributessByHashNullableAndNameInput{
						HashNullable: "string1",
						Descending:   true,
					},
				},
				output: getThingWithAdditionalAttributessByHashNullableAndNameOutput{
					thingWithAdditionalAttributess: []models.ThingWithAdditionalAttributes{
						models.ThingWithAdditionalAttributes{
							HashNullable: db.String("string1"),
							Name:         "string3",
							Version:      2,
						},
						models.ThingWithAdditionalAttributes{
							HashNullable: db.String("string1"),
							Name:         "string2",
							Version:      3,
						},
						models.ThingWithAdditionalAttributes{
							HashNullable: db.String("string1"),
							Name:         "string1",
							Version:      1,
						},
					},
					err: nil,
				},
			},
			/* FAILING_TEST */
			/* {
			   testName: "starting after",
			   d:    d,
			   input: getThingWithAdditionalAttributessByHashNullableAndNameInput{
			     ctx: context.Background(),
			     input: db.GetThingWithAdditionalAttributessByHashNullableAndNameInput{
			         HashNullable: "string1",
			       StartingAfter: &models.ThingWithAdditionalAttributes{
			         HashNullable:    db.String("string1"),
			         Name: "string1",
			         Version:    1,
			       },
			     },
			   },
			   output: getThingWithAdditionalAttributessByHashNullableAndNameOutput{
			     thingWithAdditionalAttributess: []models.ThingWithAdditionalAttributes{
			       models.ThingWithAdditionalAttributes{
			         HashNullable:    db.String("string1"),
			         Name: "string2",
			         Version:    3,
			       },
			       models.ThingWithAdditionalAttributes{
			         HashNullable:    db.String("string1"),
			         Name: "string3",
			         Version:    2,
			       },
			     },
			     err: nil,
			   },
			 }, */
			{
				testName: "starting after descending",
				d:        d,
				input: getThingWithAdditionalAttributessByHashNullableAndNameInput{
					ctx: context.Background(),
					input: db.GetThingWithAdditionalAttributessByHashNullableAndNameInput{
						HashNullable: "string1",
						StartingAfter: &models.ThingWithAdditionalAttributes{
							HashNullable: db.String("string1"),
							Name:         "string3",
							Version:      2,
						},
						Descending: true,
					},
				},
				output: getThingWithAdditionalAttributessByHashNullableAndNameOutput{
					thingWithAdditionalAttributess: []models.ThingWithAdditionalAttributes{
						models.ThingWithAdditionalAttributes{
							HashNullable: db.String("string1"),
							Name:         "string2",
							Version:      3,
						},
						models.ThingWithAdditionalAttributes{
							HashNullable: db.String("string1"),
							Name:         "string1",
							Version:      1,
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting at",
				d:        d,
				input: getThingWithAdditionalAttributessByHashNullableAndNameInput{
					ctx: context.Background(),
					input: db.GetThingWithAdditionalAttributessByHashNullableAndNameInput{
						HashNullable:   "string1",
						NameStartingAt: db.String("string2"),
					},
				},
				output: getThingWithAdditionalAttributessByHashNullableAndNameOutput{
					thingWithAdditionalAttributess: []models.ThingWithAdditionalAttributes{
						models.ThingWithAdditionalAttributes{
							HashNullable: db.String("string1"),
							Name:         "string2",
							Version:      3,
						},
						models.ThingWithAdditionalAttributes{
							HashNullable: db.String("string1"),
							Name:         "string3",
							Version:      2,
						},
					},
					err: nil,
				},
			},
		}
		for _, test := range tests {
			t.Run(test.testName, test.run)
		}
	}
}

func GetThingWithCompositeAttributes(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithCompositeAttributes{
			Branch:  db.String("string1"),
			Date:    db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			Name:    db.String("string1"),
			Version: 1,
		}
		require.Nil(t, s.SaveThingWithCompositeAttributes(ctx, m))
		m2, err := s.GetThingWithCompositeAttributes(ctx, *m.Name, *m.Branch, *m.Date)
		require.Nil(t, err)
		require.Equal(t, *m.Name, *m2.Name)
		require.Equal(t, *m.Branch, *m2.Branch)
		require.Equal(t, m.Date.String(), m2.Date.String())

		_, err = s.GetThingWithCompositeAttributes(ctx, "string2", "string2", mustTime("2018-03-11T15:04:02+07:00"))
		require.NotNil(t, err)
		require.IsType(t, err, db.ErrThingWithCompositeAttributesNotFound{})
	}
}

type getThingWithCompositeAttributessByNameBranchAndDateInput struct {
	ctx   context.Context
	input db.GetThingWithCompositeAttributessByNameBranchAndDateInput
}
type getThingWithCompositeAttributessByNameBranchAndDateOutput struct {
	thingWithCompositeAttributess []models.ThingWithCompositeAttributes
	err                           error
}
type getThingWithCompositeAttributessByNameBranchAndDateTest struct {
	testName string
	d        db.Interface
	input    getThingWithCompositeAttributessByNameBranchAndDateInput
	output   getThingWithCompositeAttributessByNameBranchAndDateOutput
}

func (g getThingWithCompositeAttributessByNameBranchAndDateTest) run(t *testing.T) {
	thingWithCompositeAttributess := []models.ThingWithCompositeAttributes{}
	fn := func(m *models.ThingWithCompositeAttributes, lastThingWithCompositeAttributes bool) bool {
		thingWithCompositeAttributess = append(thingWithCompositeAttributess, *m)
		if lastThingWithCompositeAttributes {
			return false
		}
		return true
	}
	err := g.d.GetThingWithCompositeAttributessByNameBranchAndDate(g.input.ctx, g.input.input, fn)
	if err != nil {
		fmt.Println(err.Error())
	}
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.thingWithCompositeAttributess, thingWithCompositeAttributess)
}

func GetThingWithCompositeAttributessByNameBranchAndDate(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithCompositeAttributes(ctx, models.ThingWithCompositeAttributes{
			Name:   db.String("string1"),
			Branch: db.String("string1"),
			Date:   db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
		}))
		require.Nil(t, d.SaveThingWithCompositeAttributes(ctx, models.ThingWithCompositeAttributes{
			Name:   db.String("string1"),
			Branch: db.String("string1"),
			Date:   db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
		}))
		require.Nil(t, d.SaveThingWithCompositeAttributes(ctx, models.ThingWithCompositeAttributes{
			Name:   db.String("string1"),
			Branch: db.String("string1"),
			Date:   db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
		}))
		limit := int64(3)
		tests := []getThingWithCompositeAttributessByNameBranchAndDateTest{
			{
				testName: "basic",
				d:        d,
				input: getThingWithCompositeAttributessByNameBranchAndDateInput{
					ctx: context.Background(),
					input: db.GetThingWithCompositeAttributessByNameBranchAndDateInput{
						Name:   "string1",
						Branch: "string1",
						Limit:  &limit,
					},
				},
				output: getThingWithCompositeAttributessByNameBranchAndDateOutput{
					thingWithCompositeAttributess: []models.ThingWithCompositeAttributes{
						models.ThingWithCompositeAttributes{
							Name:   db.String("string1"),
							Branch: db.String("string1"),
							Date:   db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
						},
						models.ThingWithCompositeAttributes{
							Name:   db.String("string1"),
							Branch: db.String("string1"),
							Date:   db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
						},
						models.ThingWithCompositeAttributes{
							Name:   db.String("string1"),
							Branch: db.String("string1"),
							Date:   db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
						},
					},
					err: nil,
				},
			},
			{
				testName: "descending",
				d:        d,
				input: getThingWithCompositeAttributessByNameBranchAndDateInput{
					ctx: context.Background(),
					input: db.GetThingWithCompositeAttributessByNameBranchAndDateInput{
						Name:       "string1",
						Branch:     "string1",
						Descending: true,
					},
				},
				output: getThingWithCompositeAttributessByNameBranchAndDateOutput{
					thingWithCompositeAttributess: []models.ThingWithCompositeAttributes{
						models.ThingWithCompositeAttributes{
							Name:   db.String("string1"),
							Branch: db.String("string1"),
							Date:   db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
						},
						models.ThingWithCompositeAttributes{
							Name:   db.String("string1"),
							Branch: db.String("string1"),
							Date:   db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
						},
						models.ThingWithCompositeAttributes{
							Name:   db.String("string1"),
							Branch: db.String("string1"),
							Date:   db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
						},
					},
					err: nil,
				},
			},
			/* FAILING_TEST */
			/* {
			   testName: "starting after",
			   d:    d,
			   input: getThingWithCompositeAttributessByNameBranchAndDateInput{
			     ctx: context.Background(),
			     input: db.GetThingWithCompositeAttributessByNameBranchAndDateInput{
			         Name: "string1",
			         Branch: "string1",
			       StartingAfter: &models.ThingWithCompositeAttributes{
			           Name:    db.String("string1"),
			           Branch:    db.String("string1"),
			           Date:    db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			       },
			     },
			   },
			   output: getThingWithCompositeAttributessByNameBranchAndDateOutput{
			     thingWithCompositeAttributess: []models.ThingWithCompositeAttributes{
			       models.ThingWithCompositeAttributes{
			           Name:    db.String("string1"),
			           Branch:    db.String("string1"),
			           Date: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
			       },
			       models.ThingWithCompositeAttributes{
			           Name:    db.String("string1"),
			           Branch:    db.String("string1"),
			           Date: db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
			       },
			     },
			     err: nil,
			   },
			 }, */
			{
				testName: "starting after descending",
				d:        d,
				input: getThingWithCompositeAttributessByNameBranchAndDateInput{
					ctx: context.Background(),
					input: db.GetThingWithCompositeAttributessByNameBranchAndDateInput{
						Name:   "string1",
						Branch: "string1",
						StartingAfter: &models.ThingWithCompositeAttributes{
							Name:   db.String("string1"),
							Branch: db.String("string1"),
							Date:   db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
						},
						Descending: true,
					},
				},
				output: getThingWithCompositeAttributessByNameBranchAndDateOutput{
					thingWithCompositeAttributess: []models.ThingWithCompositeAttributes{
						models.ThingWithCompositeAttributes{
							Name:   db.String("string1"),
							Branch: db.String("string1"),
							Date:   db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
						},
						models.ThingWithCompositeAttributes{
							Name:   db.String("string1"),
							Branch: db.String("string1"),
							Date:   db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting at",
				d:        d,
				input: getThingWithCompositeAttributessByNameBranchAndDateInput{
					ctx: context.Background(),
					input: db.GetThingWithCompositeAttributessByNameBranchAndDateInput{
						Name:           "string1",
						Branch:         "string1",
						DateStartingAt: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
					},
				},
				output: getThingWithCompositeAttributessByNameBranchAndDateOutput{
					thingWithCompositeAttributess: []models.ThingWithCompositeAttributes{
						models.ThingWithCompositeAttributes{
							Name:   db.String("string1"),
							Branch: db.String("string1"),
							Date:   db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
						},
						models.ThingWithCompositeAttributes{
							Name:   db.String("string1"),
							Branch: db.String("string1"),
							Date:   db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
						},
					},
					err: nil,
				},
			},
		}
		for _, test := range tests {
			t.Run(test.testName, test.run)
		}
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanThingWithCompositeAttributess(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithCompositeAttributes(ctx, models.ThingWithCompositeAttributes{
			Branch:  db.String("string1"),
			Date:    db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			Name:    db.String("string1"),
			Version: 1,
		}))
		require.Nil(t, d.SaveThingWithCompositeAttributes(ctx, models.ThingWithCompositeAttributes{
			Branch:  db.String("string2"),
			Date:    db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
			Name:    db.String("string2"),
			Version: 2,
		}))
		require.Nil(t, d.SaveThingWithCompositeAttributes(ctx, models.ThingWithCompositeAttributes{
			Branch:  db.String("string3"),
			Date:    db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
			Name:    db.String("string3"),
			Version: 3,
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.ThingWithCompositeAttributes{
				models.ThingWithCompositeAttributes{
					Branch:  db.String("string1"),
					Date:    db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
					Name:    db.String("string1"),
					Version: 1,
				},
				models.ThingWithCompositeAttributes{
					Branch:  db.String("string2"),
					Date:    db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
					Name:    db.String("string2"),
					Version: 2,
				},
				models.ThingWithCompositeAttributes{
					Branch:  db.String("string3"),
					Date:    db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
					Name:    db.String("string3"),
					Version: 3,
				},
			}
			actual := []models.ThingWithCompositeAttributes{}
			err := d.ScanThingWithCompositeAttributess(ctx, db.ScanThingWithCompositeAttributessInput{}, func(m *models.ThingWithCompositeAttributes, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		// FAILING_TEST
		t.Run("starting after", func(t *testing.T) {
			t.Skip()
			// Scan for everything.
			allItems := []models.ThingWithCompositeAttributes{}
			err := d.ScanThingWithCompositeAttributess(ctx, db.ScanThingWithCompositeAttributessInput{}, func(m *models.ThingWithCompositeAttributes, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput := db.ScanThingWithCompositeAttributessInput{
				StartingAfter: &models.ThingWithCompositeAttributes{
					Name:   firstItem.Name,
					Branch: firstItem.Branch,
					Date:   firstItem.Date,
				},
			}
			actual := []models.ThingWithCompositeAttributes{}
			err = d.ScanThingWithCompositeAttributess(ctx, scanInput, func(m *models.ThingWithCompositeAttributes, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanThingWithCompositeAttributessInput{
				Limit: &limit,
			}
			actual := []models.ThingWithCompositeAttributes{}
			err := d.ScanThingWithCompositeAttributess(ctx, scanInput, func(m *models.ThingWithCompositeAttributes, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

func SaveThingWithCompositeAttributes(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithCompositeAttributes{
			Branch:  db.String("string1"),
			Date:    db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			Name:    db.String("string1"),
			Version: 1,
		}
		require.Nil(t, s.SaveThingWithCompositeAttributes(ctx, m))
		require.IsType(t, db.ErrThingWithCompositeAttributesAlreadyExists{}, s.SaveThingWithCompositeAttributes(ctx, m))
	}
}

func DeleteThingWithCompositeAttributes(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithCompositeAttributes{
			Branch:  db.String("string1"),
			Date:    db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			Name:    db.String("string1"),
			Version: 1,
		}
		require.Nil(t, s.SaveThingWithCompositeAttributes(ctx, m))
		require.Nil(t, s.DeleteThingWithCompositeAttributes(ctx, *m.Name, *m.Branch, *m.Date))
	}
}

type getThingWithCompositeAttributessByNameVersionAndDateInput struct {
	ctx   context.Context
	input db.GetThingWithCompositeAttributessByNameVersionAndDateInput
}
type getThingWithCompositeAttributessByNameVersionAndDateOutput struct {
	thingWithCompositeAttributess []models.ThingWithCompositeAttributes
	err                           error
}
type getThingWithCompositeAttributessByNameVersionAndDateTest struct {
	testName string
	d        db.Interface
	input    getThingWithCompositeAttributessByNameVersionAndDateInput
	output   getThingWithCompositeAttributessByNameVersionAndDateOutput
}

func (g getThingWithCompositeAttributessByNameVersionAndDateTest) run(t *testing.T) {
	thingWithCompositeAttributess := []models.ThingWithCompositeAttributes{}
	fn := func(m *models.ThingWithCompositeAttributes, lastThingWithCompositeAttributes bool) bool {
		thingWithCompositeAttributess = append(thingWithCompositeAttributess, *m)
		if lastThingWithCompositeAttributes {
			return false
		}
		return true
	}
	err := g.d.GetThingWithCompositeAttributessByNameVersionAndDate(g.input.ctx, g.input.input, fn)
	if err != nil {
		fmt.Println(err.Error())
	}
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.thingWithCompositeAttributess, thingWithCompositeAttributess)
}

func GetThingWithCompositeAttributessByNameVersionAndDate(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithCompositeAttributes(ctx, models.ThingWithCompositeAttributes{
			Name:    db.String("string1"),
			Version: 1,
			Date:    db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			Branch:  db.String("string1"),
		}))
		require.Nil(t, d.SaveThingWithCompositeAttributes(ctx, models.ThingWithCompositeAttributes{
			Name:    db.String("string1"),
			Version: 1,
			Date:    db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
			Branch:  db.String("string3"),
		}))
		require.Nil(t, d.SaveThingWithCompositeAttributes(ctx, models.ThingWithCompositeAttributes{
			Name:    db.String("string1"),
			Version: 1,
			Date:    db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
			Branch:  db.String("string2"),
		}))
		limit := int64(3)
		tests := []getThingWithCompositeAttributessByNameVersionAndDateTest{
			{
				testName: "basic",
				d:        d,
				input: getThingWithCompositeAttributessByNameVersionAndDateInput{
					ctx: context.Background(),
					input: db.GetThingWithCompositeAttributessByNameVersionAndDateInput{
						Name:    "string1",
						Version: 1,
						Limit:   &limit,
					},
				},
				output: getThingWithCompositeAttributessByNameVersionAndDateOutput{
					thingWithCompositeAttributess: []models.ThingWithCompositeAttributes{
						models.ThingWithCompositeAttributes{
							Name:    db.String("string1"),
							Version: 1,
							Date:    db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
							Branch:  db.String("string1"),
						},
						models.ThingWithCompositeAttributes{
							Name:    db.String("string1"),
							Version: 1,
							Date:    db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
							Branch:  db.String("string3"),
						},
						models.ThingWithCompositeAttributes{
							Name:    db.String("string1"),
							Version: 1,
							Date:    db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
							Branch:  db.String("string2"),
						},
					},
					err: nil,
				},
			},
			{
				testName: "descending",
				d:        d,
				input: getThingWithCompositeAttributessByNameVersionAndDateInput{
					ctx: context.Background(),
					input: db.GetThingWithCompositeAttributessByNameVersionAndDateInput{
						Name:       "string1",
						Version:    1,
						Descending: true,
					},
				},
				output: getThingWithCompositeAttributessByNameVersionAndDateOutput{
					thingWithCompositeAttributess: []models.ThingWithCompositeAttributes{
						models.ThingWithCompositeAttributes{
							Name:    db.String("string1"),
							Version: 1,
							Date:    db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
							Branch:  db.String("string2"),
						},
						models.ThingWithCompositeAttributes{
							Name:    db.String("string1"),
							Version: 1,
							Date:    db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
							Branch:  db.String("string3"),
						},
						models.ThingWithCompositeAttributes{
							Name:    db.String("string1"),
							Version: 1,
							Date:    db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
							Branch:  db.String("string1"),
						},
					},
					err: nil,
				},
			},
			/* FAILING_TEST */
			/* {
			   testName: "starting after",
			   d:    d,
			   input: getThingWithCompositeAttributessByNameVersionAndDateInput{
			     ctx: context.Background(),
			     input: db.GetThingWithCompositeAttributessByNameVersionAndDateInput{
			         Name: "string1",
			         Version: 1,
			       StartingAfter: &models.ThingWithCompositeAttributes{
			         Name:    db.String("string1"),
			         Version:    1,
			         Date: db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			         Branch:    db.String("string1"),
			       },
			     },
			   },
			   output: getThingWithCompositeAttributessByNameVersionAndDateOutput{
			     thingWithCompositeAttributess: []models.ThingWithCompositeAttributes{
			       models.ThingWithCompositeAttributes{
			         Name:    db.String("string1"),
			         Version:    1,
			         Date: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
			         Branch:    db.String("string3"),
			       },
			       models.ThingWithCompositeAttributes{
			         Name:    db.String("string1"),
			         Version:    1,
			         Date: db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
			         Branch:    db.String("string2"),
			       },
			     },
			     err: nil,
			   },
			 }, */
			{
				testName: "starting after descending",
				d:        d,
				input: getThingWithCompositeAttributessByNameVersionAndDateInput{
					ctx: context.Background(),
					input: db.GetThingWithCompositeAttributessByNameVersionAndDateInput{
						Name:    "string1",
						Version: 1,
						StartingAfter: &models.ThingWithCompositeAttributes{
							Name:    db.String("string1"),
							Version: 1,
							Date:    db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
							Branch:  db.String("string2"),
						},
						Descending: true,
					},
				},
				output: getThingWithCompositeAttributessByNameVersionAndDateOutput{
					thingWithCompositeAttributess: []models.ThingWithCompositeAttributes{
						models.ThingWithCompositeAttributes{
							Name:    db.String("string1"),
							Version: 1,
							Date:    db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
							Branch:  db.String("string3"),
						},
						models.ThingWithCompositeAttributes{
							Name:    db.String("string1"),
							Version: 1,
							Date:    db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
							Branch:  db.String("string1"),
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting at",
				d:        d,
				input: getThingWithCompositeAttributessByNameVersionAndDateInput{
					ctx: context.Background(),
					input: db.GetThingWithCompositeAttributessByNameVersionAndDateInput{
						Name:           "string1",
						Version:        1,
						DateStartingAt: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
					},
				},
				output: getThingWithCompositeAttributessByNameVersionAndDateOutput{
					thingWithCompositeAttributess: []models.ThingWithCompositeAttributes{
						models.ThingWithCompositeAttributes{
							Name:    db.String("string1"),
							Version: 1,
							Date:    db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
							Branch:  db.String("string3"),
						},
						models.ThingWithCompositeAttributes{
							Name:    db.String("string1"),
							Version: 1,
							Date:    db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
							Branch:  db.String("string2"),
						},
					},
					err: nil,
				},
			},
		}
		for _, test := range tests {
			t.Run(test.testName, test.run)
		}
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanThingWithCompositeAttributessByNameVersionAndDate(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithCompositeAttributes(ctx, models.ThingWithCompositeAttributes{
			Name:    db.String("string1"),
			Version: 1,
			Date:    db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			Branch:  db.String("string1"),
		}))
		require.Nil(t, d.SaveThingWithCompositeAttributes(ctx, models.ThingWithCompositeAttributes{
			Name:    db.String("string2"),
			Version: 2,
			Date:    db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
			Branch:  db.String("string2"),
		}))
		require.Nil(t, d.SaveThingWithCompositeAttributes(ctx, models.ThingWithCompositeAttributes{
			Name:    db.String("string3"),
			Version: 3,
			Date:    db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
			Branch:  db.String("string3"),
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.ThingWithCompositeAttributes{
				models.ThingWithCompositeAttributes{
					Name:    db.String("string1"),
					Version: 1,
					Date:    db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
					Branch:  db.String("string1"),
				},
				models.ThingWithCompositeAttributes{
					Name:    db.String("string2"),
					Version: 2,
					Date:    db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
					Branch:  db.String("string2"),
				},
				models.ThingWithCompositeAttributes{
					Name:    db.String("string3"),
					Version: 3,
					Date:    db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
					Branch:  db.String("string3"),
				},
			}
			// Consistent read must be disabled when scaning a GSI.
			scanInput := db.ScanThingWithCompositeAttributessByNameVersionAndDateInput{DisableConsistentRead: true}
			actual := []models.ThingWithCompositeAttributes{}
			err := d.ScanThingWithCompositeAttributessByNameVersionAndDate(ctx, scanInput, func(m *models.ThingWithCompositeAttributes, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		t.Run("starting after", func(t *testing.T) {
			// Scan for everything.
			allItems := []models.ThingWithCompositeAttributes{}
			// Consistent read must be disabled when scaning a GSI.
			scanInput := db.ScanThingWithCompositeAttributessByNameVersionAndDateInput{DisableConsistentRead: true}
			err := d.ScanThingWithCompositeAttributessByNameVersionAndDate(ctx, scanInput, func(m *models.ThingWithCompositeAttributes, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput = db.ScanThingWithCompositeAttributessByNameVersionAndDateInput{
				DisableConsistentRead: true,
				StartingAfter: &models.ThingWithCompositeAttributes{
					Name:    firstItem.Name,
					Version: firstItem.Version,
					Date:    firstItem.Date,
					Branch:  firstItem.Branch,
				},
			}
			actual := []models.ThingWithCompositeAttributes{}
			err = d.ScanThingWithCompositeAttributessByNameVersionAndDate(ctx, scanInput, func(m *models.ThingWithCompositeAttributes, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanThingWithCompositeAttributessInput{
				Limit: &limit,
			}
			actual := []models.ThingWithCompositeAttributes{}
			err := d.ScanThingWithCompositeAttributess(ctx, scanInput, func(m *models.ThingWithCompositeAttributes, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

func GetThingWithCompositeEnumAttributes(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithCompositeEnumAttributes{
			BranchID: models.BranchMaster,
			Date:     db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			Name:     db.String("string1"),
		}
		require.Nil(t, s.SaveThingWithCompositeEnumAttributes(ctx, m))
		m2, err := s.GetThingWithCompositeEnumAttributes(ctx, *m.Name, m.BranchID, *m.Date)
		require.Nil(t, err)
		require.Equal(t, *m.Name, *m2.Name)
		require.Equal(t, m.BranchID, m2.BranchID)
		require.Equal(t, m.Date.String(), m2.Date.String())

		_, err = s.GetThingWithCompositeEnumAttributes(ctx, "string2", models.BranchDEVBRANCH, mustTime("2018-03-11T15:04:02+07:00"))
		require.NotNil(t, err)
		require.IsType(t, err, db.ErrThingWithCompositeEnumAttributesNotFound{})
	}
}

type getThingWithCompositeEnumAttributessByNameBranchAndDateInput struct {
	ctx   context.Context
	input db.GetThingWithCompositeEnumAttributessByNameBranchAndDateInput
}
type getThingWithCompositeEnumAttributessByNameBranchAndDateOutput struct {
	thingWithCompositeEnumAttributess []models.ThingWithCompositeEnumAttributes
	err                               error
}
type getThingWithCompositeEnumAttributessByNameBranchAndDateTest struct {
	testName string
	d        db.Interface
	input    getThingWithCompositeEnumAttributessByNameBranchAndDateInput
	output   getThingWithCompositeEnumAttributessByNameBranchAndDateOutput
}

func (g getThingWithCompositeEnumAttributessByNameBranchAndDateTest) run(t *testing.T) {
	thingWithCompositeEnumAttributess := []models.ThingWithCompositeEnumAttributes{}
	fn := func(m *models.ThingWithCompositeEnumAttributes, lastThingWithCompositeEnumAttributes bool) bool {
		thingWithCompositeEnumAttributess = append(thingWithCompositeEnumAttributess, *m)
		if lastThingWithCompositeEnumAttributes {
			return false
		}
		return true
	}
	err := g.d.GetThingWithCompositeEnumAttributessByNameBranchAndDate(g.input.ctx, g.input.input, fn)
	if err != nil {
		fmt.Println(err.Error())
	}
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.thingWithCompositeEnumAttributess, thingWithCompositeEnumAttributess)
}

func GetThingWithCompositeEnumAttributessByNameBranchAndDate(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithCompositeEnumAttributes(ctx, models.ThingWithCompositeEnumAttributes{
			Name:     db.String("string1"),
			BranchID: models.BranchMaster,
			Date:     db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
		}))
		require.Nil(t, d.SaveThingWithCompositeEnumAttributes(ctx, models.ThingWithCompositeEnumAttributes{
			Name:     db.String("string1"),
			BranchID: models.BranchMaster,
			Date:     db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
		}))
		require.Nil(t, d.SaveThingWithCompositeEnumAttributes(ctx, models.ThingWithCompositeEnumAttributes{
			Name:     db.String("string1"),
			BranchID: models.BranchMaster,
			Date:     db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
		}))
		limit := int64(3)
		tests := []getThingWithCompositeEnumAttributessByNameBranchAndDateTest{
			{
				testName: "basic",
				d:        d,
				input: getThingWithCompositeEnumAttributessByNameBranchAndDateInput{
					ctx: context.Background(),
					input: db.GetThingWithCompositeEnumAttributessByNameBranchAndDateInput{
						Name:     "string1",
						BranchID: models.BranchMaster,
						Limit:    &limit,
					},
				},
				output: getThingWithCompositeEnumAttributessByNameBranchAndDateOutput{
					thingWithCompositeEnumAttributess: []models.ThingWithCompositeEnumAttributes{
						models.ThingWithCompositeEnumAttributes{
							Name:     db.String("string1"),
							BranchID: models.BranchMaster,
							Date:     db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
						},
						models.ThingWithCompositeEnumAttributes{
							Name:     db.String("string1"),
							BranchID: models.BranchMaster,
							Date:     db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
						},
						models.ThingWithCompositeEnumAttributes{
							Name:     db.String("string1"),
							BranchID: models.BranchMaster,
							Date:     db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
						},
					},
					err: nil,
				},
			},
			{
				testName: "descending",
				d:        d,
				input: getThingWithCompositeEnumAttributessByNameBranchAndDateInput{
					ctx: context.Background(),
					input: db.GetThingWithCompositeEnumAttributessByNameBranchAndDateInput{
						Name:       "string1",
						BranchID:   models.BranchMaster,
						Descending: true,
					},
				},
				output: getThingWithCompositeEnumAttributessByNameBranchAndDateOutput{
					thingWithCompositeEnumAttributess: []models.ThingWithCompositeEnumAttributes{
						models.ThingWithCompositeEnumAttributes{
							Name:     db.String("string1"),
							BranchID: models.BranchMaster,
							Date:     db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
						},
						models.ThingWithCompositeEnumAttributes{
							Name:     db.String("string1"),
							BranchID: models.BranchMaster,
							Date:     db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
						},
						models.ThingWithCompositeEnumAttributes{
							Name:     db.String("string1"),
							BranchID: models.BranchMaster,
							Date:     db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
						},
					},
					err: nil,
				},
			},
			/* FAILING_TEST */
			/* {
			   testName: "starting after",
			   d:    d,
			   input: getThingWithCompositeEnumAttributessByNameBranchAndDateInput{
			     ctx: context.Background(),
			     input: db.GetThingWithCompositeEnumAttributessByNameBranchAndDateInput{
			         Name: "string1",
			         BranchID: models.BranchMaster,
			       StartingAfter: &models.ThingWithCompositeEnumAttributes{
			           Name:    db.String("string1"),
			           BranchID:    models.BranchMaster,
			           Date:    db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			       },
			     },
			   },
			   output: getThingWithCompositeEnumAttributessByNameBranchAndDateOutput{
			     thingWithCompositeEnumAttributess: []models.ThingWithCompositeEnumAttributes{
			       models.ThingWithCompositeEnumAttributes{
			           Name:    db.String("string1"),
			           BranchID:    models.BranchMaster,
			           Date: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
			       },
			       models.ThingWithCompositeEnumAttributes{
			           Name:    db.String("string1"),
			           BranchID:    models.BranchMaster,
			           Date: db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
			       },
			     },
			     err: nil,
			   },
			 }, */
			{
				testName: "starting after descending",
				d:        d,
				input: getThingWithCompositeEnumAttributessByNameBranchAndDateInput{
					ctx: context.Background(),
					input: db.GetThingWithCompositeEnumAttributessByNameBranchAndDateInput{
						Name:     "string1",
						BranchID: models.BranchMaster,
						StartingAfter: &models.ThingWithCompositeEnumAttributes{
							Name:     db.String("string1"),
							BranchID: models.BranchMaster,
							Date:     db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
						},
						Descending: true,
					},
				},
				output: getThingWithCompositeEnumAttributessByNameBranchAndDateOutput{
					thingWithCompositeEnumAttributess: []models.ThingWithCompositeEnumAttributes{
						models.ThingWithCompositeEnumAttributes{
							Name:     db.String("string1"),
							BranchID: models.BranchMaster,
							Date:     db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
						},
						models.ThingWithCompositeEnumAttributes{
							Name:     db.String("string1"),
							BranchID: models.BranchMaster,
							Date:     db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting at",
				d:        d,
				input: getThingWithCompositeEnumAttributessByNameBranchAndDateInput{
					ctx: context.Background(),
					input: db.GetThingWithCompositeEnumAttributessByNameBranchAndDateInput{
						Name:           "string1",
						BranchID:       models.BranchMaster,
						DateStartingAt: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
					},
				},
				output: getThingWithCompositeEnumAttributessByNameBranchAndDateOutput{
					thingWithCompositeEnumAttributess: []models.ThingWithCompositeEnumAttributes{
						models.ThingWithCompositeEnumAttributes{
							Name:     db.String("string1"),
							BranchID: models.BranchMaster,
							Date:     db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
						},
						models.ThingWithCompositeEnumAttributes{
							Name:     db.String("string1"),
							BranchID: models.BranchMaster,
							Date:     db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
						},
					},
					err: nil,
				},
			},
		}
		for _, test := range tests {
			t.Run(test.testName, test.run)
		}
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanThingWithCompositeEnumAttributess(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithCompositeEnumAttributes(ctx, models.ThingWithCompositeEnumAttributes{
			BranchID: models.BranchMaster,
			Date:     db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			Name:     db.String("string1"),
		}))
		require.Nil(t, d.SaveThingWithCompositeEnumAttributes(ctx, models.ThingWithCompositeEnumAttributes{
			BranchID: models.BranchDEVBRANCH,
			Date:     db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
			Name:     db.String("string2"),
		}))
		require.Nil(t, d.SaveThingWithCompositeEnumAttributes(ctx, models.ThingWithCompositeEnumAttributes{
			BranchID: models.BranchTest,
			Date:     db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
			Name:     db.String("string3"),
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.ThingWithCompositeEnumAttributes{
				models.ThingWithCompositeEnumAttributes{
					BranchID: models.BranchMaster,
					Date:     db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
					Name:     db.String("string1"),
				},
				models.ThingWithCompositeEnumAttributes{
					BranchID: models.BranchDEVBRANCH,
					Date:     db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
					Name:     db.String("string2"),
				},
				models.ThingWithCompositeEnumAttributes{
					BranchID: models.BranchTest,
					Date:     db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
					Name:     db.String("string3"),
				},
			}
			actual := []models.ThingWithCompositeEnumAttributes{}
			err := d.ScanThingWithCompositeEnumAttributess(ctx, db.ScanThingWithCompositeEnumAttributessInput{}, func(m *models.ThingWithCompositeEnumAttributes, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		// FAILING_TEST
		t.Run("starting after", func(t *testing.T) {
			t.Skip()
			// Scan for everything.
			allItems := []models.ThingWithCompositeEnumAttributes{}
			err := d.ScanThingWithCompositeEnumAttributess(ctx, db.ScanThingWithCompositeEnumAttributessInput{}, func(m *models.ThingWithCompositeEnumAttributes, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput := db.ScanThingWithCompositeEnumAttributessInput{
				StartingAfter: &models.ThingWithCompositeEnumAttributes{
					Name:     firstItem.Name,
					BranchID: firstItem.BranchID,
					Date:     firstItem.Date,
				},
			}
			actual := []models.ThingWithCompositeEnumAttributes{}
			err = d.ScanThingWithCompositeEnumAttributess(ctx, scanInput, func(m *models.ThingWithCompositeEnumAttributes, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanThingWithCompositeEnumAttributessInput{
				Limit: &limit,
			}
			actual := []models.ThingWithCompositeEnumAttributes{}
			err := d.ScanThingWithCompositeEnumAttributess(ctx, scanInput, func(m *models.ThingWithCompositeEnumAttributes, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

func SaveThingWithCompositeEnumAttributes(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithCompositeEnumAttributes{
			BranchID: models.BranchMaster,
			Date:     db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			Name:     db.String("string1"),
		}
		require.Nil(t, s.SaveThingWithCompositeEnumAttributes(ctx, m))
		require.IsType(t, db.ErrThingWithCompositeEnumAttributesAlreadyExists{}, s.SaveThingWithCompositeEnumAttributes(ctx, m))
	}
}

func DeleteThingWithCompositeEnumAttributes(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithCompositeEnumAttributes{
			BranchID: models.BranchMaster,
			Date:     db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			Name:     db.String("string1"),
		}
		require.Nil(t, s.SaveThingWithCompositeEnumAttributes(ctx, m))
		require.Nil(t, s.DeleteThingWithCompositeEnumAttributes(ctx, *m.Name, m.BranchID, *m.Date))
	}
}

func GetThingWithDateGSI(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithDateGSI{
			DateH: mustDate("2018-03-11"),
			DateR: mustDate("2018-03-11"),
			ID:    "string1",
		}
		require.Nil(t, s.SaveThingWithDateGSI(ctx, m))
		m2, err := s.GetThingWithDateGSI(ctx, m.DateH)
		require.Nil(t, err)
		require.Equal(t, m.DateH, m2.DateH)

		_, err = s.GetThingWithDateGSI(ctx, mustDate("2018-03-12"))
		require.NotNil(t, err)
		require.IsType(t, err, db.ErrThingWithDateGSINotFound{})
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanThingWithDateGSIs(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithDateGSI(ctx, models.ThingWithDateGSI{
			DateH: mustDate("2018-03-11"),
			DateR: mustDate("2018-03-11"),
			ID:    "string1",
		}))
		require.Nil(t, d.SaveThingWithDateGSI(ctx, models.ThingWithDateGSI{
			DateH: mustDate("2018-03-12"),
			DateR: mustDate("2018-03-12"),
			ID:    "string2",
		}))
		require.Nil(t, d.SaveThingWithDateGSI(ctx, models.ThingWithDateGSI{
			DateH: mustDate("2018-03-13"),
			DateR: mustDate("2018-03-13"),
			ID:    "string3",
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.ThingWithDateGSI{
				models.ThingWithDateGSI{
					DateH: mustDate("2018-03-11"),
					DateR: mustDate("2018-03-11"),
					ID:    "string1",
				},
				models.ThingWithDateGSI{
					DateH: mustDate("2018-03-12"),
					DateR: mustDate("2018-03-12"),
					ID:    "string2",
				},
				models.ThingWithDateGSI{
					DateH: mustDate("2018-03-13"),
					DateR: mustDate("2018-03-13"),
					ID:    "string3",
				},
			}
			actual := []models.ThingWithDateGSI{}
			err := d.ScanThingWithDateGSIs(ctx, db.ScanThingWithDateGSIsInput{}, func(m *models.ThingWithDateGSI, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		// FAILING_TEST
		t.Run("starting after", func(t *testing.T) {
			t.Skip()
			// Scan for everything.
			allItems := []models.ThingWithDateGSI{}
			err := d.ScanThingWithDateGSIs(ctx, db.ScanThingWithDateGSIsInput{}, func(m *models.ThingWithDateGSI, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput := db.ScanThingWithDateGSIsInput{
				StartingAfter: &models.ThingWithDateGSI{
					DateH: firstItem.DateH,
				},
			}
			actual := []models.ThingWithDateGSI{}
			err = d.ScanThingWithDateGSIs(ctx, scanInput, func(m *models.ThingWithDateGSI, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanThingWithDateGSIsInput{
				Limit: &limit,
			}
			actual := []models.ThingWithDateGSI{}
			err := d.ScanThingWithDateGSIs(ctx, scanInput, func(m *models.ThingWithDateGSI, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

func SaveThingWithDateGSI(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithDateGSI{
			DateH: mustDate("2018-03-11"),
			DateR: mustDate("2018-03-11"),
			ID:    "string1",
		}
		require.Nil(t, s.SaveThingWithDateGSI(ctx, m))
		require.IsType(t, db.ErrThingWithDateGSIAlreadyExists{}, s.SaveThingWithDateGSI(ctx, m))
	}
}

func DeleteThingWithDateGSI(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithDateGSI{
			DateH: mustDate("2018-03-11"),
			DateR: mustDate("2018-03-11"),
			ID:    "string1",
		}
		require.Nil(t, s.SaveThingWithDateGSI(ctx, m))
		require.Nil(t, s.DeleteThingWithDateGSI(ctx, m.DateH))
	}
}

type getThingWithDateGSIsByIDAndDateRInput struct {
	ctx   context.Context
	input db.GetThingWithDateGSIsByIDAndDateRInput
}
type getThingWithDateGSIsByIDAndDateROutput struct {
	thingWithDateGSIs []models.ThingWithDateGSI
	err               error
}
type getThingWithDateGSIsByIDAndDateRTest struct {
	testName string
	d        db.Interface
	input    getThingWithDateGSIsByIDAndDateRInput
	output   getThingWithDateGSIsByIDAndDateROutput
}

func (g getThingWithDateGSIsByIDAndDateRTest) run(t *testing.T) {
	thingWithDateGSIs := []models.ThingWithDateGSI{}
	fn := func(m *models.ThingWithDateGSI, lastThingWithDateGSI bool) bool {
		thingWithDateGSIs = append(thingWithDateGSIs, *m)
		if lastThingWithDateGSI {
			return false
		}
		return true
	}
	err := g.d.GetThingWithDateGSIsByIDAndDateR(g.input.ctx, g.input.input, fn)
	if err != nil {
		fmt.Println(err.Error())
	}
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.thingWithDateGSIs, thingWithDateGSIs)
}

func GetThingWithDateGSIsByIDAndDateR(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithDateGSI(ctx, models.ThingWithDateGSI{
			ID:    "string1",
			DateR: mustDate("2018-03-11"),
			DateH: mustDate("2018-03-11"),
		}))
		require.Nil(t, d.SaveThingWithDateGSI(ctx, models.ThingWithDateGSI{
			ID:    "string1",
			DateR: mustDate("2018-03-12"),
			DateH: mustDate("2018-03-13"),
		}))
		require.Nil(t, d.SaveThingWithDateGSI(ctx, models.ThingWithDateGSI{
			ID:    "string1",
			DateR: mustDate("2018-03-13"),
			DateH: mustDate("2018-03-12"),
		}))
		limit := int64(3)
		tests := []getThingWithDateGSIsByIDAndDateRTest{
			{
				testName: "basic",
				d:        d,
				input: getThingWithDateGSIsByIDAndDateRInput{
					ctx: context.Background(),
					input: db.GetThingWithDateGSIsByIDAndDateRInput{
						ID:    "string1",
						Limit: &limit,
					},
				},
				output: getThingWithDateGSIsByIDAndDateROutput{
					thingWithDateGSIs: []models.ThingWithDateGSI{
						models.ThingWithDateGSI{
							ID:    "string1",
							DateR: mustDate("2018-03-11"),
							DateH: mustDate("2018-03-11"),
						},
						models.ThingWithDateGSI{
							ID:    "string1",
							DateR: mustDate("2018-03-12"),
							DateH: mustDate("2018-03-13"),
						},
						models.ThingWithDateGSI{
							ID:    "string1",
							DateR: mustDate("2018-03-13"),
							DateH: mustDate("2018-03-12"),
						},
					},
					err: nil,
				},
			},
			{
				testName: "descending",
				d:        d,
				input: getThingWithDateGSIsByIDAndDateRInput{
					ctx: context.Background(),
					input: db.GetThingWithDateGSIsByIDAndDateRInput{
						ID:         "string1",
						Descending: true,
					},
				},
				output: getThingWithDateGSIsByIDAndDateROutput{
					thingWithDateGSIs: []models.ThingWithDateGSI{
						models.ThingWithDateGSI{
							ID:    "string1",
							DateR: mustDate("2018-03-13"),
							DateH: mustDate("2018-03-12"),
						},
						models.ThingWithDateGSI{
							ID:    "string1",
							DateR: mustDate("2018-03-12"),
							DateH: mustDate("2018-03-13"),
						},
						models.ThingWithDateGSI{
							ID:    "string1",
							DateR: mustDate("2018-03-11"),
							DateH: mustDate("2018-03-11"),
						},
					},
					err: nil,
				},
			},
			/* FAILING_TEST */
			/* {
			   testName: "starting after",
			   d:    d,
			   input: getThingWithDateGSIsByIDAndDateRInput{
			     ctx: context.Background(),
			     input: db.GetThingWithDateGSIsByIDAndDateRInput{
			         ID: "string1",
			       StartingAfter: &models.ThingWithDateGSI{
			         ID:    "string1",
			         DateR: mustDate("2018-03-11"),
			         DateH:    mustDate("2018-03-11"),
			       },
			     },
			   },
			   output: getThingWithDateGSIsByIDAndDateROutput{
			     thingWithDateGSIs: []models.ThingWithDateGSI{
			       models.ThingWithDateGSI{
			         ID:    "string1",
			         DateR: mustDate("2018-03-12"),
			         DateH:    mustDate("2018-03-13"),
			       },
			       models.ThingWithDateGSI{
			         ID:    "string1",
			         DateR: mustDate("2018-03-13"),
			         DateH:    mustDate("2018-03-12"),
			       },
			     },
			     err: nil,
			   },
			 }, */
			{
				testName: "starting after descending",
				d:        d,
				input: getThingWithDateGSIsByIDAndDateRInput{
					ctx: context.Background(),
					input: db.GetThingWithDateGSIsByIDAndDateRInput{
						ID: "string1",
						StartingAfter: &models.ThingWithDateGSI{
							ID:    "string1",
							DateR: mustDate("2018-03-13"),
							DateH: mustDate("2018-03-12"),
						},
						Descending: true,
					},
				},
				output: getThingWithDateGSIsByIDAndDateROutput{
					thingWithDateGSIs: []models.ThingWithDateGSI{
						models.ThingWithDateGSI{
							ID:    "string1",
							DateR: mustDate("2018-03-12"),
							DateH: mustDate("2018-03-13"),
						},
						models.ThingWithDateGSI{
							ID:    "string1",
							DateR: mustDate("2018-03-11"),
							DateH: mustDate("2018-03-11"),
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting at",
				d:        d,
				input: getThingWithDateGSIsByIDAndDateRInput{
					ctx: context.Background(),
					input: db.GetThingWithDateGSIsByIDAndDateRInput{
						ID:              "string1",
						DateRStartingAt: db.Date(mustDate("2018-03-12")),
					},
				},
				output: getThingWithDateGSIsByIDAndDateROutput{
					thingWithDateGSIs: []models.ThingWithDateGSI{
						models.ThingWithDateGSI{
							ID:    "string1",
							DateR: mustDate("2018-03-12"),
							DateH: mustDate("2018-03-13"),
						},
						models.ThingWithDateGSI{
							ID:    "string1",
							DateR: mustDate("2018-03-13"),
							DateH: mustDate("2018-03-12"),
						},
					},
					err: nil,
				},
			},
		}
		for _, test := range tests {
			t.Run(test.testName, test.run)
		}
	}
}

type getThingWithDateGSIsByDateHAndIDInput struct {
	ctx   context.Context
	input db.GetThingWithDateGSIsByDateHAndIDInput
}
type getThingWithDateGSIsByDateHAndIDOutput struct {
	thingWithDateGSIs []models.ThingWithDateGSI
	err               error
}
type getThingWithDateGSIsByDateHAndIDTest struct {
	testName string
	d        db.Interface
	input    getThingWithDateGSIsByDateHAndIDInput
	output   getThingWithDateGSIsByDateHAndIDOutput
}

func (g getThingWithDateGSIsByDateHAndIDTest) run(t *testing.T) {
	thingWithDateGSIs := []models.ThingWithDateGSI{}
	fn := func(m *models.ThingWithDateGSI, lastThingWithDateGSI bool) bool {
		thingWithDateGSIs = append(thingWithDateGSIs, *m)
		if lastThingWithDateGSI {
			return false
		}
		return true
	}
	err := g.d.GetThingWithDateGSIsByDateHAndID(g.input.ctx, g.input.input, fn)
	if err != nil {
		fmt.Println(err.Error())
	}
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.thingWithDateGSIs, thingWithDateGSIs)
}

func GetThingWithDateGSIsByDateHAndID(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithDateGSI(ctx, models.ThingWithDateGSI{
			DateH: mustDate("2018-03-11"),
			ID:    "string1",
		}))
		require.Nil(t, d.SaveThingWithDateGSI(ctx, models.ThingWithDateGSI{
			DateH: mustDate("2018-03-11"),
			ID:    "string2",
		}))
		require.Nil(t, d.SaveThingWithDateGSI(ctx, models.ThingWithDateGSI{
			DateH: mustDate("2018-03-11"),
			ID:    "string3",
		}))
		limit := int64(3)
		tests := []getThingWithDateGSIsByDateHAndIDTest{
			{
				testName: "basic",
				d:        d,
				input: getThingWithDateGSIsByDateHAndIDInput{
					ctx: context.Background(),
					input: db.GetThingWithDateGSIsByDateHAndIDInput{
						DateH: mustDate("2018-03-11"),
						Limit: &limit,
					},
				},
				output: getThingWithDateGSIsByDateHAndIDOutput{
					thingWithDateGSIs: []models.ThingWithDateGSI{
						models.ThingWithDateGSI{
							DateH: mustDate("2018-03-11"),
							ID:    "string1",
						},
						models.ThingWithDateGSI{
							DateH: mustDate("2018-03-11"),
							ID:    "string2",
						},
						models.ThingWithDateGSI{
							DateH: mustDate("2018-03-11"),
							ID:    "string3",
						},
					},
					err: nil,
				},
			},
			{
				testName: "descending",
				d:        d,
				input: getThingWithDateGSIsByDateHAndIDInput{
					ctx: context.Background(),
					input: db.GetThingWithDateGSIsByDateHAndIDInput{
						DateH:      mustDate("2018-03-11"),
						Descending: true,
					},
				},
				output: getThingWithDateGSIsByDateHAndIDOutput{
					thingWithDateGSIs: []models.ThingWithDateGSI{
						models.ThingWithDateGSI{
							DateH: mustDate("2018-03-11"),
							ID:    "string3",
						},
						models.ThingWithDateGSI{
							DateH: mustDate("2018-03-11"),
							ID:    "string2",
						},
						models.ThingWithDateGSI{
							DateH: mustDate("2018-03-11"),
							ID:    "string1",
						},
					},
					err: nil,
				},
			},
			/* FAILING_TEST */
			/* {
			   testName: "starting after",
			   d:    d,
			   input: getThingWithDateGSIsByDateHAndIDInput{
			     ctx: context.Background(),
			     input: db.GetThingWithDateGSIsByDateHAndIDInput{
			         DateH: mustDate("2018-03-11"),
			       StartingAfter: &models.ThingWithDateGSI{
			         DateH:    mustDate("2018-03-11"),
			         ID: "string1",
			       },
			     },
			   },
			   output: getThingWithDateGSIsByDateHAndIDOutput{
			     thingWithDateGSIs: []models.ThingWithDateGSI{
			       models.ThingWithDateGSI{
			         DateH:    mustDate("2018-03-11"),
			         ID: "string2",
			       },
			       models.ThingWithDateGSI{
			         DateH:    mustDate("2018-03-11"),
			         ID: "string3",
			       },
			     },
			     err: nil,
			   },
			 }, */
			{
				testName: "starting after descending",
				d:        d,
				input: getThingWithDateGSIsByDateHAndIDInput{
					ctx: context.Background(),
					input: db.GetThingWithDateGSIsByDateHAndIDInput{
						DateH: mustDate("2018-03-11"),
						StartingAfter: &models.ThingWithDateGSI{
							DateH: mustDate("2018-03-11"),
							ID:    "string3",
						},
						Descending: true,
					},
				},
				output: getThingWithDateGSIsByDateHAndIDOutput{
					thingWithDateGSIs: []models.ThingWithDateGSI{
						models.ThingWithDateGSI{
							DateH: mustDate("2018-03-11"),
							ID:    "string2",
						},
						models.ThingWithDateGSI{
							DateH: mustDate("2018-03-11"),
							ID:    "string1",
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting at",
				d:        d,
				input: getThingWithDateGSIsByDateHAndIDInput{
					ctx: context.Background(),
					input: db.GetThingWithDateGSIsByDateHAndIDInput{
						DateH:        mustDate("2018-03-11"),
						IDStartingAt: db.String("string2"),
					},
				},
				output: getThingWithDateGSIsByDateHAndIDOutput{
					thingWithDateGSIs: []models.ThingWithDateGSI{
						models.ThingWithDateGSI{
							DateH: mustDate("2018-03-11"),
							ID:    "string2",
						},
						models.ThingWithDateGSI{
							DateH: mustDate("2018-03-11"),
							ID:    "string3",
						},
					},
					err: nil,
				},
			},
		}
		for _, test := range tests {
			t.Run(test.testName, test.run)
		}
	}
}

func GetThingWithDateRange(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithDateRange{
			Date: mustTime("2018-03-11T15:04:01+07:00"),
			Name: "string1",
		}
		require.Nil(t, s.SaveThingWithDateRange(ctx, m))
		m2, err := s.GetThingWithDateRange(ctx, m.Name, m.Date)
		require.Nil(t, err)
		require.Equal(t, m.Name, m2.Name)
		require.Equal(t, m.Date.String(), m2.Date.String())

		_, err = s.GetThingWithDateRange(ctx, "string2", mustTime("2018-03-11T15:04:02+07:00"))
		require.NotNil(t, err)
		require.IsType(t, err, db.ErrThingWithDateRangeNotFound{})
	}
}

type getThingWithDateRangesByNameAndDateInput struct {
	ctx   context.Context
	input db.GetThingWithDateRangesByNameAndDateInput
}
type getThingWithDateRangesByNameAndDateOutput struct {
	thingWithDateRanges []models.ThingWithDateRange
	err                 error
}
type getThingWithDateRangesByNameAndDateTest struct {
	testName string
	d        db.Interface
	input    getThingWithDateRangesByNameAndDateInput
	output   getThingWithDateRangesByNameAndDateOutput
}

func (g getThingWithDateRangesByNameAndDateTest) run(t *testing.T) {
	thingWithDateRanges := []models.ThingWithDateRange{}
	fn := func(m *models.ThingWithDateRange, lastThingWithDateRange bool) bool {
		thingWithDateRanges = append(thingWithDateRanges, *m)
		if lastThingWithDateRange {
			return false
		}
		return true
	}
	err := g.d.GetThingWithDateRangesByNameAndDate(g.input.ctx, g.input.input, fn)
	if err != nil {
		fmt.Println(err.Error())
	}
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.thingWithDateRanges, thingWithDateRanges)
}

func GetThingWithDateRangesByNameAndDate(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithDateRange(ctx, models.ThingWithDateRange{
			Name: "string1",
			Date: mustTime("2018-03-11T15:04:01+07:00"),
		}))
		require.Nil(t, d.SaveThingWithDateRange(ctx, models.ThingWithDateRange{
			Name: "string1",
			Date: mustTime("2018-03-11T15:04:02+07:00"),
		}))
		require.Nil(t, d.SaveThingWithDateRange(ctx, models.ThingWithDateRange{
			Name: "string1",
			Date: mustTime("2018-03-11T15:04:03+07:00"),
		}))
		limit := int64(3)
		tests := []getThingWithDateRangesByNameAndDateTest{
			{
				testName: "basic",
				d:        d,
				input: getThingWithDateRangesByNameAndDateInput{
					ctx: context.Background(),
					input: db.GetThingWithDateRangesByNameAndDateInput{
						Name:  "string1",
						Limit: &limit,
					},
				},
				output: getThingWithDateRangesByNameAndDateOutput{
					thingWithDateRanges: []models.ThingWithDateRange{
						models.ThingWithDateRange{
							Name: "string1",
							Date: mustTime("2018-03-11T15:04:01+07:00"),
						},
						models.ThingWithDateRange{
							Name: "string1",
							Date: mustTime("2018-03-11T15:04:02+07:00"),
						},
						models.ThingWithDateRange{
							Name: "string1",
							Date: mustTime("2018-03-11T15:04:03+07:00"),
						},
					},
					err: nil,
				},
			},
			{
				testName: "descending",
				d:        d,
				input: getThingWithDateRangesByNameAndDateInput{
					ctx: context.Background(),
					input: db.GetThingWithDateRangesByNameAndDateInput{
						Name:       "string1",
						Descending: true,
					},
				},
				output: getThingWithDateRangesByNameAndDateOutput{
					thingWithDateRanges: []models.ThingWithDateRange{
						models.ThingWithDateRange{
							Name: "string1",
							Date: mustTime("2018-03-11T15:04:03+07:00"),
						},
						models.ThingWithDateRange{
							Name: "string1",
							Date: mustTime("2018-03-11T15:04:02+07:00"),
						},
						models.ThingWithDateRange{
							Name: "string1",
							Date: mustTime("2018-03-11T15:04:01+07:00"),
						},
					},
					err: nil,
				},
			},
			/* FAILING_TEST */
			/* {
			   testName: "starting after",
			   d:    d,
			   input: getThingWithDateRangesByNameAndDateInput{
			     ctx: context.Background(),
			     input: db.GetThingWithDateRangesByNameAndDateInput{
			         Name: "string1",
			       StartingAfter: &models.ThingWithDateRange{
			           Name:    "string1",
			           Date:    mustTime("2018-03-11T15:04:01+07:00"),
			       },
			     },
			   },
			   output: getThingWithDateRangesByNameAndDateOutput{
			     thingWithDateRanges: []models.ThingWithDateRange{
			       models.ThingWithDateRange{
			           Name:    "string1",
			           Date: mustTime("2018-03-11T15:04:02+07:00"),
			       },
			       models.ThingWithDateRange{
			           Name:    "string1",
			           Date: mustTime("2018-03-11T15:04:03+07:00"),
			       },
			     },
			     err: nil,
			   },
			 }, */
			{
				testName: "starting after descending",
				d:        d,
				input: getThingWithDateRangesByNameAndDateInput{
					ctx: context.Background(),
					input: db.GetThingWithDateRangesByNameAndDateInput{
						Name: "string1",
						StartingAfter: &models.ThingWithDateRange{
							Name: "string1",
							Date: mustTime("2018-03-11T15:04:03+07:00"),
						},
						Descending: true,
					},
				},
				output: getThingWithDateRangesByNameAndDateOutput{
					thingWithDateRanges: []models.ThingWithDateRange{
						models.ThingWithDateRange{
							Name: "string1",
							Date: mustTime("2018-03-11T15:04:02+07:00"),
						},
						models.ThingWithDateRange{
							Name: "string1",
							Date: mustTime("2018-03-11T15:04:01+07:00"),
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting at",
				d:        d,
				input: getThingWithDateRangesByNameAndDateInput{
					ctx: context.Background(),
					input: db.GetThingWithDateRangesByNameAndDateInput{
						Name:           "string1",
						DateStartingAt: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
					},
				},
				output: getThingWithDateRangesByNameAndDateOutput{
					thingWithDateRanges: []models.ThingWithDateRange{
						models.ThingWithDateRange{
							Name: "string1",
							Date: mustTime("2018-03-11T15:04:02+07:00"),
						},
						models.ThingWithDateRange{
							Name: "string1",
							Date: mustTime("2018-03-11T15:04:03+07:00"),
						},
					},
					err: nil,
				},
			},
		}
		for _, test := range tests {
			t.Run(test.testName, test.run)
		}
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanThingWithDateRanges(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithDateRange(ctx, models.ThingWithDateRange{
			Date: mustTime("2018-03-11T15:04:01+07:00"),
			Name: "string1",
		}))
		require.Nil(t, d.SaveThingWithDateRange(ctx, models.ThingWithDateRange{
			Date: mustTime("2018-03-11T15:04:02+07:00"),
			Name: "string2",
		}))
		require.Nil(t, d.SaveThingWithDateRange(ctx, models.ThingWithDateRange{
			Date: mustTime("2018-03-11T15:04:03+07:00"),
			Name: "string3",
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.ThingWithDateRange{
				models.ThingWithDateRange{
					Date: mustTime("2018-03-11T15:04:01+07:00"),
					Name: "string1",
				},
				models.ThingWithDateRange{
					Date: mustTime("2018-03-11T15:04:02+07:00"),
					Name: "string2",
				},
				models.ThingWithDateRange{
					Date: mustTime("2018-03-11T15:04:03+07:00"),
					Name: "string3",
				},
			}
			actual := []models.ThingWithDateRange{}
			err := d.ScanThingWithDateRanges(ctx, db.ScanThingWithDateRangesInput{}, func(m *models.ThingWithDateRange, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		// FAILING_TEST
		t.Run("starting after", func(t *testing.T) {
			t.Skip()
			// Scan for everything.
			allItems := []models.ThingWithDateRange{}
			err := d.ScanThingWithDateRanges(ctx, db.ScanThingWithDateRangesInput{}, func(m *models.ThingWithDateRange, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput := db.ScanThingWithDateRangesInput{
				StartingAfter: &models.ThingWithDateRange{
					Name: firstItem.Name,
					Date: firstItem.Date,
				},
			}
			actual := []models.ThingWithDateRange{}
			err = d.ScanThingWithDateRanges(ctx, scanInput, func(m *models.ThingWithDateRange, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanThingWithDateRangesInput{
				Limit: &limit,
			}
			actual := []models.ThingWithDateRange{}
			err := d.ScanThingWithDateRanges(ctx, scanInput, func(m *models.ThingWithDateRange, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

func SaveThingWithDateRange(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithDateRange{
			Date: mustTime("2018-03-11T15:04:01+07:00"),
			Name: "string1",
		}
		require.Nil(t, s.SaveThingWithDateRange(ctx, m))
	}
}

func DeleteThingWithDateRange(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithDateRange{
			Date: mustTime("2018-03-11T15:04:01+07:00"),
			Name: "string1",
		}
		require.Nil(t, s.SaveThingWithDateRange(ctx, m))
		require.Nil(t, s.DeleteThingWithDateRange(ctx, m.Name, m.Date))
	}
}

func GetThingWithDateRangeKey(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithDateRangeKey{
			Date: mustDate("2018-03-11"),
			ID:   "string1",
		}
		require.Nil(t, s.SaveThingWithDateRangeKey(ctx, m))
		m2, err := s.GetThingWithDateRangeKey(ctx, m.ID, m.Date)
		require.Nil(t, err)
		require.Equal(t, m.ID, m2.ID)
		require.Equal(t, m.Date, m2.Date)

		_, err = s.GetThingWithDateRangeKey(ctx, "string2", mustDate("2018-03-12"))
		require.NotNil(t, err)
		require.IsType(t, err, db.ErrThingWithDateRangeKeyNotFound{})
	}
}

type getThingWithDateRangeKeysByIDAndDateInput struct {
	ctx   context.Context
	input db.GetThingWithDateRangeKeysByIDAndDateInput
}
type getThingWithDateRangeKeysByIDAndDateOutput struct {
	thingWithDateRangeKeys []models.ThingWithDateRangeKey
	err                    error
}
type getThingWithDateRangeKeysByIDAndDateTest struct {
	testName string
	d        db.Interface
	input    getThingWithDateRangeKeysByIDAndDateInput
	output   getThingWithDateRangeKeysByIDAndDateOutput
}

func (g getThingWithDateRangeKeysByIDAndDateTest) run(t *testing.T) {
	thingWithDateRangeKeys := []models.ThingWithDateRangeKey{}
	fn := func(m *models.ThingWithDateRangeKey, lastThingWithDateRangeKey bool) bool {
		thingWithDateRangeKeys = append(thingWithDateRangeKeys, *m)
		if lastThingWithDateRangeKey {
			return false
		}
		return true
	}
	err := g.d.GetThingWithDateRangeKeysByIDAndDate(g.input.ctx, g.input.input, fn)
	if err != nil {
		fmt.Println(err.Error())
	}
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.thingWithDateRangeKeys, thingWithDateRangeKeys)
}

func GetThingWithDateRangeKeysByIDAndDate(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithDateRangeKey(ctx, models.ThingWithDateRangeKey{
			ID:   "string1",
			Date: mustDate("2018-03-11"),
		}))
		require.Nil(t, d.SaveThingWithDateRangeKey(ctx, models.ThingWithDateRangeKey{
			ID:   "string1",
			Date: mustDate("2018-03-12"),
		}))
		require.Nil(t, d.SaveThingWithDateRangeKey(ctx, models.ThingWithDateRangeKey{
			ID:   "string1",
			Date: mustDate("2018-03-13"),
		}))
		limit := int64(3)
		tests := []getThingWithDateRangeKeysByIDAndDateTest{
			{
				testName: "basic",
				d:        d,
				input: getThingWithDateRangeKeysByIDAndDateInput{
					ctx: context.Background(),
					input: db.GetThingWithDateRangeKeysByIDAndDateInput{
						ID:    "string1",
						Limit: &limit,
					},
				},
				output: getThingWithDateRangeKeysByIDAndDateOutput{
					thingWithDateRangeKeys: []models.ThingWithDateRangeKey{
						models.ThingWithDateRangeKey{
							ID:   "string1",
							Date: mustDate("2018-03-11"),
						},
						models.ThingWithDateRangeKey{
							ID:   "string1",
							Date: mustDate("2018-03-12"),
						},
						models.ThingWithDateRangeKey{
							ID:   "string1",
							Date: mustDate("2018-03-13"),
						},
					},
					err: nil,
				},
			},
			{
				testName: "descending",
				d:        d,
				input: getThingWithDateRangeKeysByIDAndDateInput{
					ctx: context.Background(),
					input: db.GetThingWithDateRangeKeysByIDAndDateInput{
						ID:         "string1",
						Descending: true,
					},
				},
				output: getThingWithDateRangeKeysByIDAndDateOutput{
					thingWithDateRangeKeys: []models.ThingWithDateRangeKey{
						models.ThingWithDateRangeKey{
							ID:   "string1",
							Date: mustDate("2018-03-13"),
						},
						models.ThingWithDateRangeKey{
							ID:   "string1",
							Date: mustDate("2018-03-12"),
						},
						models.ThingWithDateRangeKey{
							ID:   "string1",
							Date: mustDate("2018-03-11"),
						},
					},
					err: nil,
				},
			},
			/* FAILING_TEST */
			/* {
			   testName: "starting after",
			   d:    d,
			   input: getThingWithDateRangeKeysByIDAndDateInput{
			     ctx: context.Background(),
			     input: db.GetThingWithDateRangeKeysByIDAndDateInput{
			         ID: "string1",
			       StartingAfter: &models.ThingWithDateRangeKey{
			           ID:    "string1",
			           Date:    mustDate("2018-03-11"),
			       },
			     },
			   },
			   output: getThingWithDateRangeKeysByIDAndDateOutput{
			     thingWithDateRangeKeys: []models.ThingWithDateRangeKey{
			       models.ThingWithDateRangeKey{
			           ID:    "string1",
			           Date: mustDate("2018-03-12"),
			       },
			       models.ThingWithDateRangeKey{
			           ID:    "string1",
			           Date: mustDate("2018-03-13"),
			       },
			     },
			     err: nil,
			   },
			 }, */
			{
				testName: "starting after descending",
				d:        d,
				input: getThingWithDateRangeKeysByIDAndDateInput{
					ctx: context.Background(),
					input: db.GetThingWithDateRangeKeysByIDAndDateInput{
						ID: "string1",
						StartingAfter: &models.ThingWithDateRangeKey{
							ID:   "string1",
							Date: mustDate("2018-03-13"),
						},
						Descending: true,
					},
				},
				output: getThingWithDateRangeKeysByIDAndDateOutput{
					thingWithDateRangeKeys: []models.ThingWithDateRangeKey{
						models.ThingWithDateRangeKey{
							ID:   "string1",
							Date: mustDate("2018-03-12"),
						},
						models.ThingWithDateRangeKey{
							ID:   "string1",
							Date: mustDate("2018-03-11"),
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting at",
				d:        d,
				input: getThingWithDateRangeKeysByIDAndDateInput{
					ctx: context.Background(),
					input: db.GetThingWithDateRangeKeysByIDAndDateInput{
						ID:             "string1",
						DateStartingAt: db.Date(mustDate("2018-03-12")),
					},
				},
				output: getThingWithDateRangeKeysByIDAndDateOutput{
					thingWithDateRangeKeys: []models.ThingWithDateRangeKey{
						models.ThingWithDateRangeKey{
							ID:   "string1",
							Date: mustDate("2018-03-12"),
						},
						models.ThingWithDateRangeKey{
							ID:   "string1",
							Date: mustDate("2018-03-13"),
						},
					},
					err: nil,
				},
			},
		}
		for _, test := range tests {
			t.Run(test.testName, test.run)
		}
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanThingWithDateRangeKeys(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithDateRangeKey(ctx, models.ThingWithDateRangeKey{
			Date: mustDate("2018-03-11"),
			ID:   "string1",
		}))
		require.Nil(t, d.SaveThingWithDateRangeKey(ctx, models.ThingWithDateRangeKey{
			Date: mustDate("2018-03-12"),
			ID:   "string2",
		}))
		require.Nil(t, d.SaveThingWithDateRangeKey(ctx, models.ThingWithDateRangeKey{
			Date: mustDate("2018-03-13"),
			ID:   "string3",
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.ThingWithDateRangeKey{
				models.ThingWithDateRangeKey{
					Date: mustDate("2018-03-11"),
					ID:   "string1",
				},
				models.ThingWithDateRangeKey{
					Date: mustDate("2018-03-12"),
					ID:   "string2",
				},
				models.ThingWithDateRangeKey{
					Date: mustDate("2018-03-13"),
					ID:   "string3",
				},
			}
			actual := []models.ThingWithDateRangeKey{}
			err := d.ScanThingWithDateRangeKeys(ctx, db.ScanThingWithDateRangeKeysInput{}, func(m *models.ThingWithDateRangeKey, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		// FAILING_TEST
		t.Run("starting after", func(t *testing.T) {
			t.Skip()
			// Scan for everything.
			allItems := []models.ThingWithDateRangeKey{}
			err := d.ScanThingWithDateRangeKeys(ctx, db.ScanThingWithDateRangeKeysInput{}, func(m *models.ThingWithDateRangeKey, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput := db.ScanThingWithDateRangeKeysInput{
				StartingAfter: &models.ThingWithDateRangeKey{
					ID:   firstItem.ID,
					Date: firstItem.Date,
				},
			}
			actual := []models.ThingWithDateRangeKey{}
			err = d.ScanThingWithDateRangeKeys(ctx, scanInput, func(m *models.ThingWithDateRangeKey, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanThingWithDateRangeKeysInput{
				Limit: &limit,
			}
			actual := []models.ThingWithDateRangeKey{}
			err := d.ScanThingWithDateRangeKeys(ctx, scanInput, func(m *models.ThingWithDateRangeKey, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

func SaveThingWithDateRangeKey(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithDateRangeKey{
			Date: mustDate("2018-03-11"),
			ID:   "string1",
		}
		require.Nil(t, s.SaveThingWithDateRangeKey(ctx, m))
		require.IsType(t, db.ErrThingWithDateRangeKeyAlreadyExists{}, s.SaveThingWithDateRangeKey(ctx, m))
	}
}

func DeleteThingWithDateRangeKey(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithDateRangeKey{
			Date: mustDate("2018-03-11"),
			ID:   "string1",
		}
		require.Nil(t, s.SaveThingWithDateRangeKey(ctx, m))
		require.Nil(t, s.DeleteThingWithDateRangeKey(ctx, m.ID, m.Date))
	}
}

func GetThingWithDateTimeComposite(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithDateTimeComposite{
			Created:  mustTime("2018-03-11T15:04:01+07:00"),
			ID:       "string1",
			Resource: "string1",
			Type:     "string1",
		}
		require.Nil(t, s.SaveThingWithDateTimeComposite(ctx, m))
		m2, err := s.GetThingWithDateTimeComposite(ctx, m.Type, m.ID, m.Created, m.Resource)
		require.Nil(t, err)
		require.Equal(t, m.Type, m2.Type)
		require.Equal(t, m.ID, m2.ID)
		require.Equal(t, m.Created.String(), m2.Created.String())
		require.Equal(t, m.Resource, m2.Resource)

		_, err = s.GetThingWithDateTimeComposite(ctx, "string2", "string2", mustTime("2018-03-11T15:04:02+07:00"), "string2")
		require.NotNil(t, err)
		require.IsType(t, err, db.ErrThingWithDateTimeCompositeNotFound{})
	}
}

type getThingWithDateTimeCompositesByTypeIDAndCreatedResourceInput struct {
	ctx   context.Context
	input db.GetThingWithDateTimeCompositesByTypeIDAndCreatedResourceInput
}
type getThingWithDateTimeCompositesByTypeIDAndCreatedResourceOutput struct {
	thingWithDateTimeComposites []models.ThingWithDateTimeComposite
	err                         error
}
type getThingWithDateTimeCompositesByTypeIDAndCreatedResourceTest struct {
	testName string
	d        db.Interface
	input    getThingWithDateTimeCompositesByTypeIDAndCreatedResourceInput
	output   getThingWithDateTimeCompositesByTypeIDAndCreatedResourceOutput
}

func (g getThingWithDateTimeCompositesByTypeIDAndCreatedResourceTest) run(t *testing.T) {
	thingWithDateTimeComposites := []models.ThingWithDateTimeComposite{}
	fn := func(m *models.ThingWithDateTimeComposite, lastThingWithDateTimeComposite bool) bool {
		thingWithDateTimeComposites = append(thingWithDateTimeComposites, *m)
		if lastThingWithDateTimeComposite {
			return false
		}
		return true
	}
	err := g.d.GetThingWithDateTimeCompositesByTypeIDAndCreatedResource(g.input.ctx, g.input.input, fn)
	if err != nil {
		fmt.Println(err.Error())
	}
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.thingWithDateTimeComposites, thingWithDateTimeComposites)
}

func GetThingWithDateTimeCompositesByTypeIDAndCreatedResource(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithDateTimeComposite(ctx, models.ThingWithDateTimeComposite{
			Type:     "string1",
			ID:       "string1",
			Created:  mustTime("2018-03-11T15:04:01+07:00"),
			Resource: "string1",
		}))
		require.Nil(t, d.SaveThingWithDateTimeComposite(ctx, models.ThingWithDateTimeComposite{
			Type:     "string1",
			ID:       "string1",
			Created:  mustTime("2018-03-11T15:04:02+07:00"),
			Resource: "string2",
		}))
		require.Nil(t, d.SaveThingWithDateTimeComposite(ctx, models.ThingWithDateTimeComposite{
			Type:     "string1",
			ID:       "string1",
			Created:  mustTime("2018-03-11T15:04:03+07:00"),
			Resource: "string3",
		}))
		limit := int64(3)
		tests := []getThingWithDateTimeCompositesByTypeIDAndCreatedResourceTest{
			{
				testName: "basic",
				d:        d,
				input: getThingWithDateTimeCompositesByTypeIDAndCreatedResourceInput{
					ctx: context.Background(),
					input: db.GetThingWithDateTimeCompositesByTypeIDAndCreatedResourceInput{
						Type:  "string1",
						ID:    "string1",
						Limit: &limit,
					},
				},
				output: getThingWithDateTimeCompositesByTypeIDAndCreatedResourceOutput{
					thingWithDateTimeComposites: []models.ThingWithDateTimeComposite{
						models.ThingWithDateTimeComposite{
							Type:     "string1",
							ID:       "string1",
							Created:  mustTime("2018-03-11T15:04:01+07:00"),
							Resource: "string1",
						},
						models.ThingWithDateTimeComposite{
							Type:     "string1",
							ID:       "string1",
							Created:  mustTime("2018-03-11T15:04:02+07:00"),
							Resource: "string2",
						},
						models.ThingWithDateTimeComposite{
							Type:     "string1",
							ID:       "string1",
							Created:  mustTime("2018-03-11T15:04:03+07:00"),
							Resource: "string3",
						},
					},
					err: nil,
				},
			},
			{
				testName: "descending",
				d:        d,
				input: getThingWithDateTimeCompositesByTypeIDAndCreatedResourceInput{
					ctx: context.Background(),
					input: db.GetThingWithDateTimeCompositesByTypeIDAndCreatedResourceInput{
						Type:       "string1",
						ID:         "string1",
						Descending: true,
					},
				},
				output: getThingWithDateTimeCompositesByTypeIDAndCreatedResourceOutput{
					thingWithDateTimeComposites: []models.ThingWithDateTimeComposite{
						models.ThingWithDateTimeComposite{
							Type:     "string1",
							ID:       "string1",
							Created:  mustTime("2018-03-11T15:04:03+07:00"),
							Resource: "string3",
						},
						models.ThingWithDateTimeComposite{
							Type:     "string1",
							ID:       "string1",
							Created:  mustTime("2018-03-11T15:04:02+07:00"),
							Resource: "string2",
						},
						models.ThingWithDateTimeComposite{
							Type:     "string1",
							ID:       "string1",
							Created:  mustTime("2018-03-11T15:04:01+07:00"),
							Resource: "string1",
						},
					},
					err: nil,
				},
			},
			/* FAILING_TEST */
			/* {
			   testName: "starting after",
			   d:    d,
			   input: getThingWithDateTimeCompositesByTypeIDAndCreatedResourceInput{
			     ctx: context.Background(),
			     input: db.GetThingWithDateTimeCompositesByTypeIDAndCreatedResourceInput{
			         Type: "string1",
			         ID: "string1",
			       StartingAfter: &models.ThingWithDateTimeComposite{
			           Type:    "string1",
			           ID:    "string1",
			           Created:    mustTime("2018-03-11T15:04:01+07:00"),
			           Resource:    "string1",
			       },
			     },
			   },
			   output: getThingWithDateTimeCompositesByTypeIDAndCreatedResourceOutput{
			     thingWithDateTimeComposites: []models.ThingWithDateTimeComposite{
			       models.ThingWithDateTimeComposite{
			           Type:    "string1",
			           ID:    "string1",
			           Created: mustTime("2018-03-11T15:04:02+07:00"),
			           Resource: "string2",
			       },
			       models.ThingWithDateTimeComposite{
			           Type:    "string1",
			           ID:    "string1",
			           Created: mustTime("2018-03-11T15:04:03+07:00"),
			           Resource: "string3",
			       },
			     },
			     err: nil,
			   },
			 }, */
			{
				testName: "starting after descending",
				d:        d,
				input: getThingWithDateTimeCompositesByTypeIDAndCreatedResourceInput{
					ctx: context.Background(),
					input: db.GetThingWithDateTimeCompositesByTypeIDAndCreatedResourceInput{
						Type: "string1",
						ID:   "string1",
						StartingAfter: &models.ThingWithDateTimeComposite{
							Type:     "string1",
							ID:       "string1",
							Created:  mustTime("2018-03-11T15:04:03+07:00"),
							Resource: "string3",
						},
						Descending: true,
					},
				},
				output: getThingWithDateTimeCompositesByTypeIDAndCreatedResourceOutput{
					thingWithDateTimeComposites: []models.ThingWithDateTimeComposite{
						models.ThingWithDateTimeComposite{
							Type:     "string1",
							ID:       "string1",
							Created:  mustTime("2018-03-11T15:04:02+07:00"),
							Resource: "string2",
						},
						models.ThingWithDateTimeComposite{
							Type:     "string1",
							ID:       "string1",
							Created:  mustTime("2018-03-11T15:04:01+07:00"),
							Resource: "string1",
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting at",
				d:        d,
				input: getThingWithDateTimeCompositesByTypeIDAndCreatedResourceInput{
					ctx: context.Background(),
					input: db.GetThingWithDateTimeCompositesByTypeIDAndCreatedResourceInput{
						Type: "string1",
						ID:   "string1",
						StartingAt: &db.CreatedResource{
							Created:  mustTime("2018-03-11T15:04:02+07:00"),
							Resource: "string2",
						},
					},
				},
				output: getThingWithDateTimeCompositesByTypeIDAndCreatedResourceOutput{
					thingWithDateTimeComposites: []models.ThingWithDateTimeComposite{
						models.ThingWithDateTimeComposite{
							Type:     "string1",
							ID:       "string1",
							Created:  mustTime("2018-03-11T15:04:02+07:00"),
							Resource: "string2",
						},
						models.ThingWithDateTimeComposite{
							Type:     "string1",
							ID:       "string1",
							Created:  mustTime("2018-03-11T15:04:03+07:00"),
							Resource: "string3",
						},
					},
					err: nil,
				},
			},
		}
		for _, test := range tests {
			t.Run(test.testName, test.run)
		}
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanThingWithDateTimeComposites(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithDateTimeComposite(ctx, models.ThingWithDateTimeComposite{
			Created:  mustTime("2018-03-11T15:04:01+07:00"),
			ID:       "string1",
			Resource: "string1",
			Type:     "string1",
		}))
		require.Nil(t, d.SaveThingWithDateTimeComposite(ctx, models.ThingWithDateTimeComposite{
			Created:  mustTime("2018-03-11T15:04:02+07:00"),
			ID:       "string2",
			Resource: "string2",
			Type:     "string2",
		}))
		require.Nil(t, d.SaveThingWithDateTimeComposite(ctx, models.ThingWithDateTimeComposite{
			Created:  mustTime("2018-03-11T15:04:03+07:00"),
			ID:       "string3",
			Resource: "string3",
			Type:     "string3",
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.ThingWithDateTimeComposite{
				models.ThingWithDateTimeComposite{
					Created:  mustTime("2018-03-11T15:04:01+07:00"),
					ID:       "string1",
					Resource: "string1",
					Type:     "string1",
				},
				models.ThingWithDateTimeComposite{
					Created:  mustTime("2018-03-11T15:04:02+07:00"),
					ID:       "string2",
					Resource: "string2",
					Type:     "string2",
				},
				models.ThingWithDateTimeComposite{
					Created:  mustTime("2018-03-11T15:04:03+07:00"),
					ID:       "string3",
					Resource: "string3",
					Type:     "string3",
				},
			}
			actual := []models.ThingWithDateTimeComposite{}
			err := d.ScanThingWithDateTimeComposites(ctx, db.ScanThingWithDateTimeCompositesInput{}, func(m *models.ThingWithDateTimeComposite, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		// FAILING_TEST
		t.Run("starting after", func(t *testing.T) {
			t.Skip()
			// Scan for everything.
			allItems := []models.ThingWithDateTimeComposite{}
			err := d.ScanThingWithDateTimeComposites(ctx, db.ScanThingWithDateTimeCompositesInput{}, func(m *models.ThingWithDateTimeComposite, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput := db.ScanThingWithDateTimeCompositesInput{
				StartingAfter: &models.ThingWithDateTimeComposite{
					Type:     firstItem.Type,
					ID:       firstItem.ID,
					Created:  firstItem.Created,
					Resource: firstItem.Resource,
				},
			}
			actual := []models.ThingWithDateTimeComposite{}
			err = d.ScanThingWithDateTimeComposites(ctx, scanInput, func(m *models.ThingWithDateTimeComposite, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanThingWithDateTimeCompositesInput{
				Limit: &limit,
			}
			actual := []models.ThingWithDateTimeComposite{}
			err := d.ScanThingWithDateTimeComposites(ctx, scanInput, func(m *models.ThingWithDateTimeComposite, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

func SaveThingWithDateTimeComposite(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithDateTimeComposite{
			Created:  mustTime("2018-03-11T15:04:01+07:00"),
			ID:       "string1",
			Resource: "string1",
			Type:     "string1",
		}
		require.Nil(t, s.SaveThingWithDateTimeComposite(ctx, m))
	}
}

func DeleteThingWithDateTimeComposite(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithDateTimeComposite{
			Created:  mustTime("2018-03-11T15:04:01+07:00"),
			ID:       "string1",
			Resource: "string1",
			Type:     "string1",
		}
		require.Nil(t, s.SaveThingWithDateTimeComposite(ctx, m))
		require.Nil(t, s.DeleteThingWithDateTimeComposite(ctx, m.Type, m.ID, m.Created, m.Resource))
	}
}

func GetThingWithDatetimeGSI(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithDatetimeGSI{
			Datetime: mustTime("2018-03-11T15:04:01+07:00"),
			ID:       "string1",
		}
		require.Nil(t, s.SaveThingWithDatetimeGSI(ctx, m))
		m2, err := s.GetThingWithDatetimeGSI(ctx, m.ID)
		require.Nil(t, err)
		require.Equal(t, m.ID, m2.ID)

		_, err = s.GetThingWithDatetimeGSI(ctx, "string2")
		require.NotNil(t, err)
		require.IsType(t, err, db.ErrThingWithDatetimeGSINotFound{})
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanThingWithDatetimeGSIs(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithDatetimeGSI(ctx, models.ThingWithDatetimeGSI{
			Datetime: mustTime("2018-03-11T15:04:01+07:00"),
			ID:       "string1",
		}))
		require.Nil(t, d.SaveThingWithDatetimeGSI(ctx, models.ThingWithDatetimeGSI{
			Datetime: mustTime("2018-03-11T15:04:02+07:00"),
			ID:       "string2",
		}))
		require.Nil(t, d.SaveThingWithDatetimeGSI(ctx, models.ThingWithDatetimeGSI{
			Datetime: mustTime("2018-03-11T15:04:03+07:00"),
			ID:       "string3",
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.ThingWithDatetimeGSI{
				models.ThingWithDatetimeGSI{
					Datetime: mustTime("2018-03-11T15:04:01+07:00"),
					ID:       "string1",
				},
				models.ThingWithDatetimeGSI{
					Datetime: mustTime("2018-03-11T15:04:02+07:00"),
					ID:       "string2",
				},
				models.ThingWithDatetimeGSI{
					Datetime: mustTime("2018-03-11T15:04:03+07:00"),
					ID:       "string3",
				},
			}
			actual := []models.ThingWithDatetimeGSI{}
			err := d.ScanThingWithDatetimeGSIs(ctx, db.ScanThingWithDatetimeGSIsInput{}, func(m *models.ThingWithDatetimeGSI, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		// FAILING_TEST
		t.Run("starting after", func(t *testing.T) {
			t.Skip()
			// Scan for everything.
			allItems := []models.ThingWithDatetimeGSI{}
			err := d.ScanThingWithDatetimeGSIs(ctx, db.ScanThingWithDatetimeGSIsInput{}, func(m *models.ThingWithDatetimeGSI, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput := db.ScanThingWithDatetimeGSIsInput{
				StartingAfter: &models.ThingWithDatetimeGSI{
					ID: firstItem.ID,
				},
			}
			actual := []models.ThingWithDatetimeGSI{}
			err = d.ScanThingWithDatetimeGSIs(ctx, scanInput, func(m *models.ThingWithDatetimeGSI, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanThingWithDatetimeGSIsInput{
				Limit: &limit,
			}
			actual := []models.ThingWithDatetimeGSI{}
			err := d.ScanThingWithDatetimeGSIs(ctx, scanInput, func(m *models.ThingWithDatetimeGSI, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

func SaveThingWithDatetimeGSI(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithDatetimeGSI{
			Datetime: mustTime("2018-03-11T15:04:01+07:00"),
			ID:       "string1",
		}
		require.Nil(t, s.SaveThingWithDatetimeGSI(ctx, m))
		require.IsType(t, db.ErrThingWithDatetimeGSIAlreadyExists{}, s.SaveThingWithDatetimeGSI(ctx, m))
	}
}

func DeleteThingWithDatetimeGSI(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithDatetimeGSI{
			Datetime: mustTime("2018-03-11T15:04:01+07:00"),
			ID:       "string1",
		}
		require.Nil(t, s.SaveThingWithDatetimeGSI(ctx, m))
		require.Nil(t, s.DeleteThingWithDatetimeGSI(ctx, m.ID))
	}
}

type getThingWithDatetimeGSIsByDatetimeAndIDInput struct {
	ctx   context.Context
	input db.GetThingWithDatetimeGSIsByDatetimeAndIDInput
}
type getThingWithDatetimeGSIsByDatetimeAndIDOutput struct {
	thingWithDatetimeGSIs []models.ThingWithDatetimeGSI
	err                   error
}
type getThingWithDatetimeGSIsByDatetimeAndIDTest struct {
	testName string
	d        db.Interface
	input    getThingWithDatetimeGSIsByDatetimeAndIDInput
	output   getThingWithDatetimeGSIsByDatetimeAndIDOutput
}

func (g getThingWithDatetimeGSIsByDatetimeAndIDTest) run(t *testing.T) {
	thingWithDatetimeGSIs := []models.ThingWithDatetimeGSI{}
	fn := func(m *models.ThingWithDatetimeGSI, lastThingWithDatetimeGSI bool) bool {
		thingWithDatetimeGSIs = append(thingWithDatetimeGSIs, *m)
		if lastThingWithDatetimeGSI {
			return false
		}
		return true
	}
	err := g.d.GetThingWithDatetimeGSIsByDatetimeAndID(g.input.ctx, g.input.input, fn)
	if err != nil {
		fmt.Println(err.Error())
	}
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.thingWithDatetimeGSIs, thingWithDatetimeGSIs)
}

func GetThingWithDatetimeGSIsByDatetimeAndID(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithDatetimeGSI(ctx, models.ThingWithDatetimeGSI{
			Datetime: mustTime("2018-03-11T15:04:01+07:00"),
			ID:       "string1",
		}))
		require.Nil(t, d.SaveThingWithDatetimeGSI(ctx, models.ThingWithDatetimeGSI{
			Datetime: mustTime("2018-03-11T15:04:01+07:00"),
			ID:       "string2",
		}))
		require.Nil(t, d.SaveThingWithDatetimeGSI(ctx, models.ThingWithDatetimeGSI{
			Datetime: mustTime("2018-03-11T15:04:01+07:00"),
			ID:       "string3",
		}))
		limit := int64(3)
		tests := []getThingWithDatetimeGSIsByDatetimeAndIDTest{
			{
				testName: "basic",
				d:        d,
				input: getThingWithDatetimeGSIsByDatetimeAndIDInput{
					ctx: context.Background(),
					input: db.GetThingWithDatetimeGSIsByDatetimeAndIDInput{
						Datetime: mustTime("2018-03-11T15:04:01+07:00"),
						Limit:    &limit,
					},
				},
				output: getThingWithDatetimeGSIsByDatetimeAndIDOutput{
					thingWithDatetimeGSIs: []models.ThingWithDatetimeGSI{
						models.ThingWithDatetimeGSI{
							Datetime: mustTime("2018-03-11T15:04:01+07:00"),
							ID:       "string1",
						},
						models.ThingWithDatetimeGSI{
							Datetime: mustTime("2018-03-11T15:04:01+07:00"),
							ID:       "string2",
						},
						models.ThingWithDatetimeGSI{
							Datetime: mustTime("2018-03-11T15:04:01+07:00"),
							ID:       "string3",
						},
					},
					err: nil,
				},
			},
			{
				testName: "descending",
				d:        d,
				input: getThingWithDatetimeGSIsByDatetimeAndIDInput{
					ctx: context.Background(),
					input: db.GetThingWithDatetimeGSIsByDatetimeAndIDInput{
						Datetime:   mustTime("2018-03-11T15:04:01+07:00"),
						Descending: true,
					},
				},
				output: getThingWithDatetimeGSIsByDatetimeAndIDOutput{
					thingWithDatetimeGSIs: []models.ThingWithDatetimeGSI{
						models.ThingWithDatetimeGSI{
							Datetime: mustTime("2018-03-11T15:04:01+07:00"),
							ID:       "string3",
						},
						models.ThingWithDatetimeGSI{
							Datetime: mustTime("2018-03-11T15:04:01+07:00"),
							ID:       "string2",
						},
						models.ThingWithDatetimeGSI{
							Datetime: mustTime("2018-03-11T15:04:01+07:00"),
							ID:       "string1",
						},
					},
					err: nil,
				},
			},
			/* FAILING_TEST */
			/* {
			   testName: "starting after",
			   d:    d,
			   input: getThingWithDatetimeGSIsByDatetimeAndIDInput{
			     ctx: context.Background(),
			     input: db.GetThingWithDatetimeGSIsByDatetimeAndIDInput{
			         Datetime: mustTime("2018-03-11T15:04:01+07:00"),
			       StartingAfter: &models.ThingWithDatetimeGSI{
			         Datetime:    mustTime("2018-03-11T15:04:01+07:00"),
			         ID: "string1",
			       },
			     },
			   },
			   output: getThingWithDatetimeGSIsByDatetimeAndIDOutput{
			     thingWithDatetimeGSIs: []models.ThingWithDatetimeGSI{
			       models.ThingWithDatetimeGSI{
			         Datetime:    mustTime("2018-03-11T15:04:01+07:00"),
			         ID: "string2",
			       },
			       models.ThingWithDatetimeGSI{
			         Datetime:    mustTime("2018-03-11T15:04:01+07:00"),
			         ID: "string3",
			       },
			     },
			     err: nil,
			   },
			 }, */
			{
				testName: "starting after descending",
				d:        d,
				input: getThingWithDatetimeGSIsByDatetimeAndIDInput{
					ctx: context.Background(),
					input: db.GetThingWithDatetimeGSIsByDatetimeAndIDInput{
						Datetime: mustTime("2018-03-11T15:04:01+07:00"),
						StartingAfter: &models.ThingWithDatetimeGSI{
							Datetime: mustTime("2018-03-11T15:04:01+07:00"),
							ID:       "string3",
						},
						Descending: true,
					},
				},
				output: getThingWithDatetimeGSIsByDatetimeAndIDOutput{
					thingWithDatetimeGSIs: []models.ThingWithDatetimeGSI{
						models.ThingWithDatetimeGSI{
							Datetime: mustTime("2018-03-11T15:04:01+07:00"),
							ID:       "string2",
						},
						models.ThingWithDatetimeGSI{
							Datetime: mustTime("2018-03-11T15:04:01+07:00"),
							ID:       "string1",
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting at",
				d:        d,
				input: getThingWithDatetimeGSIsByDatetimeAndIDInput{
					ctx: context.Background(),
					input: db.GetThingWithDatetimeGSIsByDatetimeAndIDInput{
						Datetime:     mustTime("2018-03-11T15:04:01+07:00"),
						IDStartingAt: db.String("string2"),
					},
				},
				output: getThingWithDatetimeGSIsByDatetimeAndIDOutput{
					thingWithDatetimeGSIs: []models.ThingWithDatetimeGSI{
						models.ThingWithDatetimeGSI{
							Datetime: mustTime("2018-03-11T15:04:01+07:00"),
							ID:       "string2",
						},
						models.ThingWithDatetimeGSI{
							Datetime: mustTime("2018-03-11T15:04:01+07:00"),
							ID:       "string3",
						},
					},
					err: nil,
				},
			},
		}
		for _, test := range tests {
			t.Run(test.testName, test.run)
		}
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanThingWithDatetimeGSIsByDatetimeAndID(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithDatetimeGSI(ctx, models.ThingWithDatetimeGSI{
			Datetime: mustTime("2018-03-11T15:04:01+07:00"),
			ID:       "string1",
		}))
		require.Nil(t, d.SaveThingWithDatetimeGSI(ctx, models.ThingWithDatetimeGSI{
			Datetime: mustTime("2018-03-11T15:04:02+07:00"),
			ID:       "string2",
		}))
		require.Nil(t, d.SaveThingWithDatetimeGSI(ctx, models.ThingWithDatetimeGSI{
			Datetime: mustTime("2018-03-11T15:04:03+07:00"),
			ID:       "string3",
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.ThingWithDatetimeGSI{
				models.ThingWithDatetimeGSI{
					Datetime: mustTime("2018-03-11T15:04:01+07:00"),
					ID:       "string1",
				},
				models.ThingWithDatetimeGSI{
					Datetime: mustTime("2018-03-11T15:04:02+07:00"),
					ID:       "string2",
				},
				models.ThingWithDatetimeGSI{
					Datetime: mustTime("2018-03-11T15:04:03+07:00"),
					ID:       "string3",
				},
			}
			// Consistent read must be disabled when scaning a GSI.
			scanInput := db.ScanThingWithDatetimeGSIsByDatetimeAndIDInput{DisableConsistentRead: true}
			actual := []models.ThingWithDatetimeGSI{}
			err := d.ScanThingWithDatetimeGSIsByDatetimeAndID(ctx, scanInput, func(m *models.ThingWithDatetimeGSI, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		t.Run("starting after", func(t *testing.T) {
			// Scan for everything.
			allItems := []models.ThingWithDatetimeGSI{}
			// Consistent read must be disabled when scaning a GSI.
			scanInput := db.ScanThingWithDatetimeGSIsByDatetimeAndIDInput{DisableConsistentRead: true}
			err := d.ScanThingWithDatetimeGSIsByDatetimeAndID(ctx, scanInput, func(m *models.ThingWithDatetimeGSI, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput = db.ScanThingWithDatetimeGSIsByDatetimeAndIDInput{
				DisableConsistentRead: true,
				StartingAfter: &models.ThingWithDatetimeGSI{
					Datetime: firstItem.Datetime,
					ID:       firstItem.ID,
				},
			}
			actual := []models.ThingWithDatetimeGSI{}
			err = d.ScanThingWithDatetimeGSIsByDatetimeAndID(ctx, scanInput, func(m *models.ThingWithDatetimeGSI, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanThingWithDatetimeGSIsInput{
				Limit: &limit,
			}
			actual := []models.ThingWithDatetimeGSI{}
			err := d.ScanThingWithDatetimeGSIs(ctx, scanInput, func(m *models.ThingWithDatetimeGSI, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

func GetThingWithEnumHashKey(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithEnumHashKey{
			Branch: models.BranchMaster,
			Date:   mustTime("2018-03-11T15:04:01+07:00"),
			Date2:  mustTime("2018-03-11T15:04:01+07:00"),
		}
		require.Nil(t, s.SaveThingWithEnumHashKey(ctx, m))
		m2, err := s.GetThingWithEnumHashKey(ctx, m.Branch, m.Date)
		require.Nil(t, err)
		require.Equal(t, m.Branch, m2.Branch)
		require.Equal(t, m.Date.String(), m2.Date.String())

		_, err = s.GetThingWithEnumHashKey(ctx, models.BranchDEVBRANCH, mustTime("2018-03-11T15:04:02+07:00"))
		require.NotNil(t, err)
		require.IsType(t, err, db.ErrThingWithEnumHashKeyNotFound{})
	}
}

type getThingWithEnumHashKeysByBranchAndDateInput struct {
	ctx   context.Context
	input db.GetThingWithEnumHashKeysByBranchAndDateInput
}
type getThingWithEnumHashKeysByBranchAndDateOutput struct {
	thingWithEnumHashKeys []models.ThingWithEnumHashKey
	err                   error
}
type getThingWithEnumHashKeysByBranchAndDateTest struct {
	testName string
	d        db.Interface
	input    getThingWithEnumHashKeysByBranchAndDateInput
	output   getThingWithEnumHashKeysByBranchAndDateOutput
}

func (g getThingWithEnumHashKeysByBranchAndDateTest) run(t *testing.T) {
	thingWithEnumHashKeys := []models.ThingWithEnumHashKey{}
	fn := func(m *models.ThingWithEnumHashKey, lastThingWithEnumHashKey bool) bool {
		thingWithEnumHashKeys = append(thingWithEnumHashKeys, *m)
		if lastThingWithEnumHashKey {
			return false
		}
		return true
	}
	err := g.d.GetThingWithEnumHashKeysByBranchAndDate(g.input.ctx, g.input.input, fn)
	if err != nil {
		fmt.Println(err.Error())
	}
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.thingWithEnumHashKeys, thingWithEnumHashKeys)
}

func GetThingWithEnumHashKeysByBranchAndDate(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithEnumHashKey(ctx, models.ThingWithEnumHashKey{
			Branch: models.BranchMaster,
			Date:   mustTime("2018-03-11T15:04:01+07:00"),
		}))
		require.Nil(t, d.SaveThingWithEnumHashKey(ctx, models.ThingWithEnumHashKey{
			Branch: models.BranchMaster,
			Date:   mustTime("2018-03-11T15:04:02+07:00"),
		}))
		require.Nil(t, d.SaveThingWithEnumHashKey(ctx, models.ThingWithEnumHashKey{
			Branch: models.BranchMaster,
			Date:   mustTime("2018-03-11T15:04:03+07:00"),
		}))
		limit := int64(3)
		tests := []getThingWithEnumHashKeysByBranchAndDateTest{
			{
				testName: "basic",
				d:        d,
				input: getThingWithEnumHashKeysByBranchAndDateInput{
					ctx: context.Background(),
					input: db.GetThingWithEnumHashKeysByBranchAndDateInput{
						Branch: models.BranchMaster,
						Limit:  &limit,
					},
				},
				output: getThingWithEnumHashKeysByBranchAndDateOutput{
					thingWithEnumHashKeys: []models.ThingWithEnumHashKey{
						models.ThingWithEnumHashKey{
							Branch: models.BranchMaster,
							Date:   mustTime("2018-03-11T15:04:01+07:00"),
						},
						models.ThingWithEnumHashKey{
							Branch: models.BranchMaster,
							Date:   mustTime("2018-03-11T15:04:02+07:00"),
						},
						models.ThingWithEnumHashKey{
							Branch: models.BranchMaster,
							Date:   mustTime("2018-03-11T15:04:03+07:00"),
						},
					},
					err: nil,
				},
			},
			{
				testName: "descending",
				d:        d,
				input: getThingWithEnumHashKeysByBranchAndDateInput{
					ctx: context.Background(),
					input: db.GetThingWithEnumHashKeysByBranchAndDateInput{
						Branch:     models.BranchMaster,
						Descending: true,
					},
				},
				output: getThingWithEnumHashKeysByBranchAndDateOutput{
					thingWithEnumHashKeys: []models.ThingWithEnumHashKey{
						models.ThingWithEnumHashKey{
							Branch: models.BranchMaster,
							Date:   mustTime("2018-03-11T15:04:03+07:00"),
						},
						models.ThingWithEnumHashKey{
							Branch: models.BranchMaster,
							Date:   mustTime("2018-03-11T15:04:02+07:00"),
						},
						models.ThingWithEnumHashKey{
							Branch: models.BranchMaster,
							Date:   mustTime("2018-03-11T15:04:01+07:00"),
						},
					},
					err: nil,
				},
			},
			/* FAILING_TEST */
			/* {
			   testName: "starting after",
			   d:    d,
			   input: getThingWithEnumHashKeysByBranchAndDateInput{
			     ctx: context.Background(),
			     input: db.GetThingWithEnumHashKeysByBranchAndDateInput{
			         Branch: models.BranchMaster,
			       StartingAfter: &models.ThingWithEnumHashKey{
			           Branch:    models.BranchMaster,
			           Date:    mustTime("2018-03-11T15:04:01+07:00"),
			       },
			     },
			   },
			   output: getThingWithEnumHashKeysByBranchAndDateOutput{
			     thingWithEnumHashKeys: []models.ThingWithEnumHashKey{
			       models.ThingWithEnumHashKey{
			           Branch:    models.BranchMaster,
			           Date: mustTime("2018-03-11T15:04:02+07:00"),
			       },
			       models.ThingWithEnumHashKey{
			           Branch:    models.BranchMaster,
			           Date: mustTime("2018-03-11T15:04:03+07:00"),
			       },
			     },
			     err: nil,
			   },
			 }, */
			{
				testName: "starting after descending",
				d:        d,
				input: getThingWithEnumHashKeysByBranchAndDateInput{
					ctx: context.Background(),
					input: db.GetThingWithEnumHashKeysByBranchAndDateInput{
						Branch: models.BranchMaster,
						StartingAfter: &models.ThingWithEnumHashKey{
							Branch: models.BranchMaster,
							Date:   mustTime("2018-03-11T15:04:03+07:00"),
						},
						Descending: true,
					},
				},
				output: getThingWithEnumHashKeysByBranchAndDateOutput{
					thingWithEnumHashKeys: []models.ThingWithEnumHashKey{
						models.ThingWithEnumHashKey{
							Branch: models.BranchMaster,
							Date:   mustTime("2018-03-11T15:04:02+07:00"),
						},
						models.ThingWithEnumHashKey{
							Branch: models.BranchMaster,
							Date:   mustTime("2018-03-11T15:04:01+07:00"),
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting at",
				d:        d,
				input: getThingWithEnumHashKeysByBranchAndDateInput{
					ctx: context.Background(),
					input: db.GetThingWithEnumHashKeysByBranchAndDateInput{
						Branch:         models.BranchMaster,
						DateStartingAt: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
					},
				},
				output: getThingWithEnumHashKeysByBranchAndDateOutput{
					thingWithEnumHashKeys: []models.ThingWithEnumHashKey{
						models.ThingWithEnumHashKey{
							Branch: models.BranchMaster,
							Date:   mustTime("2018-03-11T15:04:02+07:00"),
						},
						models.ThingWithEnumHashKey{
							Branch: models.BranchMaster,
							Date:   mustTime("2018-03-11T15:04:03+07:00"),
						},
					},
					err: nil,
				},
			},
		}
		for _, test := range tests {
			t.Run(test.testName, test.run)
		}
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanThingWithEnumHashKeys(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithEnumHashKey(ctx, models.ThingWithEnumHashKey{
			Branch: models.BranchMaster,
			Date:   mustTime("2018-03-11T15:04:01+07:00"),
			Date2:  mustTime("2018-03-11T15:04:01+07:00"),
		}))
		require.Nil(t, d.SaveThingWithEnumHashKey(ctx, models.ThingWithEnumHashKey{
			Branch: models.BranchDEVBRANCH,
			Date:   mustTime("2018-03-11T15:04:02+07:00"),
			Date2:  mustTime("2018-03-11T15:04:02+07:00"),
		}))
		require.Nil(t, d.SaveThingWithEnumHashKey(ctx, models.ThingWithEnumHashKey{
			Branch: models.BranchTest,
			Date:   mustTime("2018-03-11T15:04:03+07:00"),
			Date2:  mustTime("2018-03-11T15:04:03+07:00"),
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.ThingWithEnumHashKey{
				models.ThingWithEnumHashKey{
					Branch: models.BranchMaster,
					Date:   mustTime("2018-03-11T15:04:01+07:00"),
					Date2:  mustTime("2018-03-11T15:04:01+07:00"),
				},
				models.ThingWithEnumHashKey{
					Branch: models.BranchDEVBRANCH,
					Date:   mustTime("2018-03-11T15:04:02+07:00"),
					Date2:  mustTime("2018-03-11T15:04:02+07:00"),
				},
				models.ThingWithEnumHashKey{
					Branch: models.BranchTest,
					Date:   mustTime("2018-03-11T15:04:03+07:00"),
					Date2:  mustTime("2018-03-11T15:04:03+07:00"),
				},
			}
			actual := []models.ThingWithEnumHashKey{}
			err := d.ScanThingWithEnumHashKeys(ctx, db.ScanThingWithEnumHashKeysInput{}, func(m *models.ThingWithEnumHashKey, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		// FAILING_TEST
		t.Run("starting after", func(t *testing.T) {
			t.Skip()
			// Scan for everything.
			allItems := []models.ThingWithEnumHashKey{}
			err := d.ScanThingWithEnumHashKeys(ctx, db.ScanThingWithEnumHashKeysInput{}, func(m *models.ThingWithEnumHashKey, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput := db.ScanThingWithEnumHashKeysInput{
				StartingAfter: &models.ThingWithEnumHashKey{
					Branch: firstItem.Branch,
					Date:   firstItem.Date,
				},
			}
			actual := []models.ThingWithEnumHashKey{}
			err = d.ScanThingWithEnumHashKeys(ctx, scanInput, func(m *models.ThingWithEnumHashKey, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanThingWithEnumHashKeysInput{
				Limit: &limit,
			}
			actual := []models.ThingWithEnumHashKey{}
			err := d.ScanThingWithEnumHashKeys(ctx, scanInput, func(m *models.ThingWithEnumHashKey, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

func SaveThingWithEnumHashKey(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithEnumHashKey{
			Branch: models.BranchMaster,
			Date:   mustTime("2018-03-11T15:04:01+07:00"),
			Date2:  mustTime("2018-03-11T15:04:01+07:00"),
		}
		require.Nil(t, s.SaveThingWithEnumHashKey(ctx, m))
		require.IsType(t, db.ErrThingWithEnumHashKeyAlreadyExists{}, s.SaveThingWithEnumHashKey(ctx, m))
	}
}

func DeleteThingWithEnumHashKey(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithEnumHashKey{
			Branch: models.BranchMaster,
			Date:   mustTime("2018-03-11T15:04:01+07:00"),
			Date2:  mustTime("2018-03-11T15:04:01+07:00"),
		}
		require.Nil(t, s.SaveThingWithEnumHashKey(ctx, m))
		require.Nil(t, s.DeleteThingWithEnumHashKey(ctx, m.Branch, m.Date))
	}
}

type getThingWithEnumHashKeysByBranchAndDate2Input struct {
	ctx   context.Context
	input db.GetThingWithEnumHashKeysByBranchAndDate2Input
}
type getThingWithEnumHashKeysByBranchAndDate2Output struct {
	thingWithEnumHashKeys []models.ThingWithEnumHashKey
	err                   error
}
type getThingWithEnumHashKeysByBranchAndDate2Test struct {
	testName string
	d        db.Interface
	input    getThingWithEnumHashKeysByBranchAndDate2Input
	output   getThingWithEnumHashKeysByBranchAndDate2Output
}

func (g getThingWithEnumHashKeysByBranchAndDate2Test) run(t *testing.T) {
	thingWithEnumHashKeys := []models.ThingWithEnumHashKey{}
	fn := func(m *models.ThingWithEnumHashKey, lastThingWithEnumHashKey bool) bool {
		thingWithEnumHashKeys = append(thingWithEnumHashKeys, *m)
		if lastThingWithEnumHashKey {
			return false
		}
		return true
	}
	err := g.d.GetThingWithEnumHashKeysByBranchAndDate2(g.input.ctx, g.input.input, fn)
	if err != nil {
		fmt.Println(err.Error())
	}
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.thingWithEnumHashKeys, thingWithEnumHashKeys)
}

func GetThingWithEnumHashKeysByBranchAndDate2(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithEnumHashKey(ctx, models.ThingWithEnumHashKey{
			Branch: models.BranchMaster,
			Date2:  mustTime("2018-03-11T15:04:01+07:00"),
			Date:   mustTime("2018-03-11T15:04:01+07:00"),
		}))
		require.Nil(t, d.SaveThingWithEnumHashKey(ctx, models.ThingWithEnumHashKey{
			Branch: models.BranchMaster,
			Date2:  mustTime("2018-03-11T15:04:02+07:00"),
			Date:   mustTime("2018-03-11T15:04:03+07:00"),
		}))
		require.Nil(t, d.SaveThingWithEnumHashKey(ctx, models.ThingWithEnumHashKey{
			Branch: models.BranchMaster,
			Date2:  mustTime("2018-03-11T15:04:03+07:00"),
			Date:   mustTime("2018-03-11T15:04:02+07:00"),
		}))
		limit := int64(3)
		tests := []getThingWithEnumHashKeysByBranchAndDate2Test{
			{
				testName: "basic",
				d:        d,
				input: getThingWithEnumHashKeysByBranchAndDate2Input{
					ctx: context.Background(),
					input: db.GetThingWithEnumHashKeysByBranchAndDate2Input{
						Branch: models.BranchMaster,
						Limit:  &limit,
					},
				},
				output: getThingWithEnumHashKeysByBranchAndDate2Output{
					thingWithEnumHashKeys: []models.ThingWithEnumHashKey{
						models.ThingWithEnumHashKey{
							Branch: models.BranchMaster,
							Date2:  mustTime("2018-03-11T15:04:01+07:00"),
							Date:   mustTime("2018-03-11T15:04:01+07:00"),
						},
						models.ThingWithEnumHashKey{
							Branch: models.BranchMaster,
							Date2:  mustTime("2018-03-11T15:04:02+07:00"),
							Date:   mustTime("2018-03-11T15:04:03+07:00"),
						},
						models.ThingWithEnumHashKey{
							Branch: models.BranchMaster,
							Date2:  mustTime("2018-03-11T15:04:03+07:00"),
							Date:   mustTime("2018-03-11T15:04:02+07:00"),
						},
					},
					err: nil,
				},
			},
			{
				testName: "descending",
				d:        d,
				input: getThingWithEnumHashKeysByBranchAndDate2Input{
					ctx: context.Background(),
					input: db.GetThingWithEnumHashKeysByBranchAndDate2Input{
						Branch:     models.BranchMaster,
						Descending: true,
					},
				},
				output: getThingWithEnumHashKeysByBranchAndDate2Output{
					thingWithEnumHashKeys: []models.ThingWithEnumHashKey{
						models.ThingWithEnumHashKey{
							Branch: models.BranchMaster,
							Date2:  mustTime("2018-03-11T15:04:03+07:00"),
							Date:   mustTime("2018-03-11T15:04:02+07:00"),
						},
						models.ThingWithEnumHashKey{
							Branch: models.BranchMaster,
							Date2:  mustTime("2018-03-11T15:04:02+07:00"),
							Date:   mustTime("2018-03-11T15:04:03+07:00"),
						},
						models.ThingWithEnumHashKey{
							Branch: models.BranchMaster,
							Date2:  mustTime("2018-03-11T15:04:01+07:00"),
							Date:   mustTime("2018-03-11T15:04:01+07:00"),
						},
					},
					err: nil,
				},
			},
			/* FAILING_TEST */
			/* {
			   testName: "starting after",
			   d:    d,
			   input: getThingWithEnumHashKeysByBranchAndDate2Input{
			     ctx: context.Background(),
			     input: db.GetThingWithEnumHashKeysByBranchAndDate2Input{
			         Branch: models.BranchMaster,
			       StartingAfter: &models.ThingWithEnumHashKey{
			         Branch:    models.BranchMaster,
			         Date2: mustTime("2018-03-11T15:04:01+07:00"),
			         Date:    mustTime("2018-03-11T15:04:01+07:00"),
			       },
			     },
			   },
			   output: getThingWithEnumHashKeysByBranchAndDate2Output{
			     thingWithEnumHashKeys: []models.ThingWithEnumHashKey{
			       models.ThingWithEnumHashKey{
			         Branch:    models.BranchMaster,
			         Date2: mustTime("2018-03-11T15:04:02+07:00"),
			         Date:    mustTime("2018-03-11T15:04:03+07:00"),
			       },
			       models.ThingWithEnumHashKey{
			         Branch:    models.BranchMaster,
			         Date2: mustTime("2018-03-11T15:04:03+07:00"),
			         Date:    mustTime("2018-03-11T15:04:02+07:00"),
			       },
			     },
			     err: nil,
			   },
			 }, */
			{
				testName: "starting after descending",
				d:        d,
				input: getThingWithEnumHashKeysByBranchAndDate2Input{
					ctx: context.Background(),
					input: db.GetThingWithEnumHashKeysByBranchAndDate2Input{
						Branch: models.BranchMaster,
						StartingAfter: &models.ThingWithEnumHashKey{
							Branch: models.BranchMaster,
							Date2:  mustTime("2018-03-11T15:04:03+07:00"),
							Date:   mustTime("2018-03-11T15:04:02+07:00"),
						},
						Descending: true,
					},
				},
				output: getThingWithEnumHashKeysByBranchAndDate2Output{
					thingWithEnumHashKeys: []models.ThingWithEnumHashKey{
						models.ThingWithEnumHashKey{
							Branch: models.BranchMaster,
							Date2:  mustTime("2018-03-11T15:04:02+07:00"),
							Date:   mustTime("2018-03-11T15:04:03+07:00"),
						},
						models.ThingWithEnumHashKey{
							Branch: models.BranchMaster,
							Date2:  mustTime("2018-03-11T15:04:01+07:00"),
							Date:   mustTime("2018-03-11T15:04:01+07:00"),
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting at",
				d:        d,
				input: getThingWithEnumHashKeysByBranchAndDate2Input{
					ctx: context.Background(),
					input: db.GetThingWithEnumHashKeysByBranchAndDate2Input{
						Branch:          models.BranchMaster,
						Date2StartingAt: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
					},
				},
				output: getThingWithEnumHashKeysByBranchAndDate2Output{
					thingWithEnumHashKeys: []models.ThingWithEnumHashKey{
						models.ThingWithEnumHashKey{
							Branch: models.BranchMaster,
							Date2:  mustTime("2018-03-11T15:04:02+07:00"),
							Date:   mustTime("2018-03-11T15:04:03+07:00"),
						},
						models.ThingWithEnumHashKey{
							Branch: models.BranchMaster,
							Date2:  mustTime("2018-03-11T15:04:03+07:00"),
							Date:   mustTime("2018-03-11T15:04:02+07:00"),
						},
					},
					err: nil,
				},
			},
		}
		for _, test := range tests {
			t.Run(test.testName, test.run)
		}
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanThingWithEnumHashKeysByBranchAndDate2(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithEnumHashKey(ctx, models.ThingWithEnumHashKey{
			Branch: models.BranchMaster,
			Date2:  mustTime("2018-03-11T15:04:01+07:00"),
			Date:   mustTime("2018-03-11T15:04:01+07:00"),
		}))
		require.Nil(t, d.SaveThingWithEnumHashKey(ctx, models.ThingWithEnumHashKey{
			Branch: models.BranchDEVBRANCH,
			Date2:  mustTime("2018-03-11T15:04:02+07:00"),
			Date:   mustTime("2018-03-11T15:04:02+07:00"),
		}))
		require.Nil(t, d.SaveThingWithEnumHashKey(ctx, models.ThingWithEnumHashKey{
			Branch: models.BranchTest,
			Date2:  mustTime("2018-03-11T15:04:03+07:00"),
			Date:   mustTime("2018-03-11T15:04:03+07:00"),
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.ThingWithEnumHashKey{
				models.ThingWithEnumHashKey{
					Branch: models.BranchMaster,
					Date2:  mustTime("2018-03-11T15:04:01+07:00"),
					Date:   mustTime("2018-03-11T15:04:01+07:00"),
				},
				models.ThingWithEnumHashKey{
					Branch: models.BranchDEVBRANCH,
					Date2:  mustTime("2018-03-11T15:04:02+07:00"),
					Date:   mustTime("2018-03-11T15:04:02+07:00"),
				},
				models.ThingWithEnumHashKey{
					Branch: models.BranchTest,
					Date2:  mustTime("2018-03-11T15:04:03+07:00"),
					Date:   mustTime("2018-03-11T15:04:03+07:00"),
				},
			}
			// Consistent read must be disabled when scaning a GSI.
			scanInput := db.ScanThingWithEnumHashKeysByBranchAndDate2Input{DisableConsistentRead: true}
			actual := []models.ThingWithEnumHashKey{}
			err := d.ScanThingWithEnumHashKeysByBranchAndDate2(ctx, scanInput, func(m *models.ThingWithEnumHashKey, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		t.Run("starting after", func(t *testing.T) {
			// Scan for everything.
			allItems := []models.ThingWithEnumHashKey{}
			// Consistent read must be disabled when scaning a GSI.
			scanInput := db.ScanThingWithEnumHashKeysByBranchAndDate2Input{DisableConsistentRead: true}
			err := d.ScanThingWithEnumHashKeysByBranchAndDate2(ctx, scanInput, func(m *models.ThingWithEnumHashKey, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput = db.ScanThingWithEnumHashKeysByBranchAndDate2Input{
				DisableConsistentRead: true,
				StartingAfter: &models.ThingWithEnumHashKey{
					Branch: firstItem.Branch,
					Date2:  firstItem.Date2,
					Date:   firstItem.Date,
				},
			}
			actual := []models.ThingWithEnumHashKey{}
			err = d.ScanThingWithEnumHashKeysByBranchAndDate2(ctx, scanInput, func(m *models.ThingWithEnumHashKey, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanThingWithEnumHashKeysInput{
				Limit: &limit,
			}
			actual := []models.ThingWithEnumHashKey{}
			err := d.ScanThingWithEnumHashKeys(ctx, scanInput, func(m *models.ThingWithEnumHashKey, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

func GetThingWithMatchingKeys(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithMatchingKeys{
			AssocID:   "string1",
			AssocType: "string1",
			Bear:      "string1",
			Created:   mustTime("2018-03-11T15:04:01+07:00"),
		}
		require.Nil(t, s.SaveThingWithMatchingKeys(ctx, m))
		m2, err := s.GetThingWithMatchingKeys(ctx, m.Bear, m.AssocType, m.AssocID)
		require.Nil(t, err)
		require.Equal(t, m.Bear, m2.Bear)
		require.Equal(t, m.AssocType, m2.AssocType)
		require.Equal(t, m.AssocID, m2.AssocID)

		_, err = s.GetThingWithMatchingKeys(ctx, "string2", "string2", "string2")
		require.NotNil(t, err)
		require.IsType(t, err, db.ErrThingWithMatchingKeysNotFound{})
	}
}

type getThingWithMatchingKeyssByBearAndAssocTypeIDInput struct {
	ctx   context.Context
	input db.GetThingWithMatchingKeyssByBearAndAssocTypeIDInput
}
type getThingWithMatchingKeyssByBearAndAssocTypeIDOutput struct {
	thingWithMatchingKeyss []models.ThingWithMatchingKeys
	err                    error
}
type getThingWithMatchingKeyssByBearAndAssocTypeIDTest struct {
	testName string
	d        db.Interface
	input    getThingWithMatchingKeyssByBearAndAssocTypeIDInput
	output   getThingWithMatchingKeyssByBearAndAssocTypeIDOutput
}

func (g getThingWithMatchingKeyssByBearAndAssocTypeIDTest) run(t *testing.T) {
	thingWithMatchingKeyss := []models.ThingWithMatchingKeys{}
	fn := func(m *models.ThingWithMatchingKeys, lastThingWithMatchingKeys bool) bool {
		thingWithMatchingKeyss = append(thingWithMatchingKeyss, *m)
		if lastThingWithMatchingKeys {
			return false
		}
		return true
	}
	err := g.d.GetThingWithMatchingKeyssByBearAndAssocTypeID(g.input.ctx, g.input.input, fn)
	if err != nil {
		fmt.Println(err.Error())
	}
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.thingWithMatchingKeyss, thingWithMatchingKeyss)
}

func GetThingWithMatchingKeyssByBearAndAssocTypeID(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithMatchingKeys(ctx, models.ThingWithMatchingKeys{
			Bear:      "string1",
			AssocType: "string1",
			AssocID:   "string1",
		}))
		require.Nil(t, d.SaveThingWithMatchingKeys(ctx, models.ThingWithMatchingKeys{
			Bear:      "string1",
			AssocType: "string2",
			AssocID:   "string2",
		}))
		require.Nil(t, d.SaveThingWithMatchingKeys(ctx, models.ThingWithMatchingKeys{
			Bear:      "string1",
			AssocType: "string3",
			AssocID:   "string3",
		}))
		limit := int64(3)
		tests := []getThingWithMatchingKeyssByBearAndAssocTypeIDTest{
			{
				testName: "basic",
				d:        d,
				input: getThingWithMatchingKeyssByBearAndAssocTypeIDInput{
					ctx: context.Background(),
					input: db.GetThingWithMatchingKeyssByBearAndAssocTypeIDInput{
						Bear:  "string1",
						Limit: &limit,
					},
				},
				output: getThingWithMatchingKeyssByBearAndAssocTypeIDOutput{
					thingWithMatchingKeyss: []models.ThingWithMatchingKeys{
						models.ThingWithMatchingKeys{
							Bear:      "string1",
							AssocType: "string1",
							AssocID:   "string1",
						},
						models.ThingWithMatchingKeys{
							Bear:      "string1",
							AssocType: "string2",
							AssocID:   "string2",
						},
						models.ThingWithMatchingKeys{
							Bear:      "string1",
							AssocType: "string3",
							AssocID:   "string3",
						},
					},
					err: nil,
				},
			},
			{
				testName: "descending",
				d:        d,
				input: getThingWithMatchingKeyssByBearAndAssocTypeIDInput{
					ctx: context.Background(),
					input: db.GetThingWithMatchingKeyssByBearAndAssocTypeIDInput{
						Bear:       "string1",
						Descending: true,
					},
				},
				output: getThingWithMatchingKeyssByBearAndAssocTypeIDOutput{
					thingWithMatchingKeyss: []models.ThingWithMatchingKeys{
						models.ThingWithMatchingKeys{
							Bear:      "string1",
							AssocType: "string3",
							AssocID:   "string3",
						},
						models.ThingWithMatchingKeys{
							Bear:      "string1",
							AssocType: "string2",
							AssocID:   "string2",
						},
						models.ThingWithMatchingKeys{
							Bear:      "string1",
							AssocType: "string1",
							AssocID:   "string1",
						},
					},
					err: nil,
				},
			},
			/* FAILING_TEST */
			/* {
			   testName: "starting after",
			   d:    d,
			   input: getThingWithMatchingKeyssByBearAndAssocTypeIDInput{
			     ctx: context.Background(),
			     input: db.GetThingWithMatchingKeyssByBearAndAssocTypeIDInput{
			         Bear: "string1",
			       StartingAfter: &models.ThingWithMatchingKeys{
			           Bear:    "string1",
			           AssocType:    "string1",
			           AssocID:    "string1",
			       },
			     },
			   },
			   output: getThingWithMatchingKeyssByBearAndAssocTypeIDOutput{
			     thingWithMatchingKeyss: []models.ThingWithMatchingKeys{
			       models.ThingWithMatchingKeys{
			           Bear:    "string1",
			           AssocType: "string2",
			           AssocID: "string2",
			       },
			       models.ThingWithMatchingKeys{
			           Bear:    "string1",
			           AssocType: "string3",
			           AssocID: "string3",
			       },
			     },
			     err: nil,
			   },
			 }, */
			{
				testName: "starting after descending",
				d:        d,
				input: getThingWithMatchingKeyssByBearAndAssocTypeIDInput{
					ctx: context.Background(),
					input: db.GetThingWithMatchingKeyssByBearAndAssocTypeIDInput{
						Bear: "string1",
						StartingAfter: &models.ThingWithMatchingKeys{
							Bear:      "string1",
							AssocType: "string3",
							AssocID:   "string3",
						},
						Descending: true,
					},
				},
				output: getThingWithMatchingKeyssByBearAndAssocTypeIDOutput{
					thingWithMatchingKeyss: []models.ThingWithMatchingKeys{
						models.ThingWithMatchingKeys{
							Bear:      "string1",
							AssocType: "string2",
							AssocID:   "string2",
						},
						models.ThingWithMatchingKeys{
							Bear:      "string1",
							AssocType: "string1",
							AssocID:   "string1",
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting at",
				d:        d,
				input: getThingWithMatchingKeyssByBearAndAssocTypeIDInput{
					ctx: context.Background(),
					input: db.GetThingWithMatchingKeyssByBearAndAssocTypeIDInput{
						Bear: "string1",
						StartingAt: &db.AssocTypeAssocID{
							AssocType: "string2",
							AssocID:   "string2",
						},
					},
				},
				output: getThingWithMatchingKeyssByBearAndAssocTypeIDOutput{
					thingWithMatchingKeyss: []models.ThingWithMatchingKeys{
						models.ThingWithMatchingKeys{
							Bear:      "string1",
							AssocType: "string2",
							AssocID:   "string2",
						},
						models.ThingWithMatchingKeys{
							Bear:      "string1",
							AssocType: "string3",
							AssocID:   "string3",
						},
					},
					err: nil,
				},
			},
		}
		for _, test := range tests {
			t.Run(test.testName, test.run)
		}
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanThingWithMatchingKeyss(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithMatchingKeys(ctx, models.ThingWithMatchingKeys{
			AssocID:   "string1",
			AssocType: "string1",
			Bear:      "string1",
			Created:   mustTime("2018-03-11T15:04:01+07:00"),
		}))
		require.Nil(t, d.SaveThingWithMatchingKeys(ctx, models.ThingWithMatchingKeys{
			AssocID:   "string2",
			AssocType: "string2",
			Bear:      "string2",
			Created:   mustTime("2018-03-11T15:04:02+07:00"),
		}))
		require.Nil(t, d.SaveThingWithMatchingKeys(ctx, models.ThingWithMatchingKeys{
			AssocID:   "string3",
			AssocType: "string3",
			Bear:      "string3",
			Created:   mustTime("2018-03-11T15:04:03+07:00"),
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.ThingWithMatchingKeys{
				models.ThingWithMatchingKeys{
					AssocID:   "string1",
					AssocType: "string1",
					Bear:      "string1",
					Created:   mustTime("2018-03-11T15:04:01+07:00"),
				},
				models.ThingWithMatchingKeys{
					AssocID:   "string2",
					AssocType: "string2",
					Bear:      "string2",
					Created:   mustTime("2018-03-11T15:04:02+07:00"),
				},
				models.ThingWithMatchingKeys{
					AssocID:   "string3",
					AssocType: "string3",
					Bear:      "string3",
					Created:   mustTime("2018-03-11T15:04:03+07:00"),
				},
			}
			actual := []models.ThingWithMatchingKeys{}
			err := d.ScanThingWithMatchingKeyss(ctx, db.ScanThingWithMatchingKeyssInput{}, func(m *models.ThingWithMatchingKeys, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		// FAILING_TEST
		t.Run("starting after", func(t *testing.T) {
			t.Skip()
			// Scan for everything.
			allItems := []models.ThingWithMatchingKeys{}
			err := d.ScanThingWithMatchingKeyss(ctx, db.ScanThingWithMatchingKeyssInput{}, func(m *models.ThingWithMatchingKeys, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput := db.ScanThingWithMatchingKeyssInput{
				StartingAfter: &models.ThingWithMatchingKeys{
					Bear:      firstItem.Bear,
					AssocType: firstItem.AssocType,
					AssocID:   firstItem.AssocID,
				},
			}
			actual := []models.ThingWithMatchingKeys{}
			err = d.ScanThingWithMatchingKeyss(ctx, scanInput, func(m *models.ThingWithMatchingKeys, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanThingWithMatchingKeyssInput{
				Limit: &limit,
			}
			actual := []models.ThingWithMatchingKeys{}
			err := d.ScanThingWithMatchingKeyss(ctx, scanInput, func(m *models.ThingWithMatchingKeys, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

func SaveThingWithMatchingKeys(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithMatchingKeys{
			AssocID:   "string1",
			AssocType: "string1",
			Bear:      "string1",
			Created:   mustTime("2018-03-11T15:04:01+07:00"),
		}
		require.Nil(t, s.SaveThingWithMatchingKeys(ctx, m))
	}
}

func DeleteThingWithMatchingKeys(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithMatchingKeys{
			AssocID:   "string1",
			AssocType: "string1",
			Bear:      "string1",
			Created:   mustTime("2018-03-11T15:04:01+07:00"),
		}
		require.Nil(t, s.SaveThingWithMatchingKeys(ctx, m))
		require.Nil(t, s.DeleteThingWithMatchingKeys(ctx, m.Bear, m.AssocType, m.AssocID))
	}
}

type getThingWithMatchingKeyssByAssocTypeIDAndCreatedBearInput struct {
	ctx   context.Context
	input db.GetThingWithMatchingKeyssByAssocTypeIDAndCreatedBearInput
}
type getThingWithMatchingKeyssByAssocTypeIDAndCreatedBearOutput struct {
	thingWithMatchingKeyss []models.ThingWithMatchingKeys
	err                    error
}
type getThingWithMatchingKeyssByAssocTypeIDAndCreatedBearTest struct {
	testName string
	d        db.Interface
	input    getThingWithMatchingKeyssByAssocTypeIDAndCreatedBearInput
	output   getThingWithMatchingKeyssByAssocTypeIDAndCreatedBearOutput
}

func (g getThingWithMatchingKeyssByAssocTypeIDAndCreatedBearTest) run(t *testing.T) {
	thingWithMatchingKeyss := []models.ThingWithMatchingKeys{}
	fn := func(m *models.ThingWithMatchingKeys, lastThingWithMatchingKeys bool) bool {
		thingWithMatchingKeyss = append(thingWithMatchingKeyss, *m)
		if lastThingWithMatchingKeys {
			return false
		}
		return true
	}
	err := g.d.GetThingWithMatchingKeyssByAssocTypeIDAndCreatedBear(g.input.ctx, g.input.input, fn)
	if err != nil {
		fmt.Println(err.Error())
	}
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.thingWithMatchingKeyss, thingWithMatchingKeyss)
}

func GetThingWithMatchingKeyssByAssocTypeIDAndCreatedBear(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithMatchingKeys(ctx, models.ThingWithMatchingKeys{
			AssocType: "string1",
			AssocID:   "string1",
			Created:   mustTime("2018-03-11T15:04:01+07:00"),
			Bear:      "string1",
		}))
		require.Nil(t, d.SaveThingWithMatchingKeys(ctx, models.ThingWithMatchingKeys{
			AssocType: "string1",
			AssocID:   "string1",
			Created:   mustTime("2018-03-11T15:04:02+07:00"),
			Bear:      "string2",
		}))
		require.Nil(t, d.SaveThingWithMatchingKeys(ctx, models.ThingWithMatchingKeys{
			AssocType: "string1",
			AssocID:   "string1",
			Created:   mustTime("2018-03-11T15:04:03+07:00"),
			Bear:      "string3",
		}))
		limit := int64(3)
		tests := []getThingWithMatchingKeyssByAssocTypeIDAndCreatedBearTest{
			{
				testName: "basic",
				d:        d,
				input: getThingWithMatchingKeyssByAssocTypeIDAndCreatedBearInput{
					ctx: context.Background(),
					input: db.GetThingWithMatchingKeyssByAssocTypeIDAndCreatedBearInput{
						AssocType: "string1",
						AssocID:   "string1",
						Limit:     &limit,
					},
				},
				output: getThingWithMatchingKeyssByAssocTypeIDAndCreatedBearOutput{
					thingWithMatchingKeyss: []models.ThingWithMatchingKeys{
						models.ThingWithMatchingKeys{
							AssocType: "string1",
							AssocID:   "string1",
							Created:   mustTime("2018-03-11T15:04:01+07:00"),
							Bear:      "string1",
						},
						models.ThingWithMatchingKeys{
							AssocType: "string1",
							AssocID:   "string1",
							Created:   mustTime("2018-03-11T15:04:02+07:00"),
							Bear:      "string2",
						},
						models.ThingWithMatchingKeys{
							AssocType: "string1",
							AssocID:   "string1",
							Created:   mustTime("2018-03-11T15:04:03+07:00"),
							Bear:      "string3",
						},
					},
					err: nil,
				},
			},
			{
				testName: "descending",
				d:        d,
				input: getThingWithMatchingKeyssByAssocTypeIDAndCreatedBearInput{
					ctx: context.Background(),
					input: db.GetThingWithMatchingKeyssByAssocTypeIDAndCreatedBearInput{
						AssocType:  "string1",
						AssocID:    "string1",
						Descending: true,
					},
				},
				output: getThingWithMatchingKeyssByAssocTypeIDAndCreatedBearOutput{
					thingWithMatchingKeyss: []models.ThingWithMatchingKeys{
						models.ThingWithMatchingKeys{
							AssocType: "string1",
							AssocID:   "string1",
							Created:   mustTime("2018-03-11T15:04:03+07:00"),
							Bear:      "string3",
						},
						models.ThingWithMatchingKeys{
							AssocType: "string1",
							AssocID:   "string1",
							Created:   mustTime("2018-03-11T15:04:02+07:00"),
							Bear:      "string2",
						},
						models.ThingWithMatchingKeys{
							AssocType: "string1",
							AssocID:   "string1",
							Created:   mustTime("2018-03-11T15:04:01+07:00"),
							Bear:      "string1",
						},
					},
					err: nil,
				},
			},
			/* FAILING_TEST */
			/* {
			   testName: "starting after",
			   d:    d,
			   input: getThingWithMatchingKeyssByAssocTypeIDAndCreatedBearInput{
			     ctx: context.Background(),
			     input: db.GetThingWithMatchingKeyssByAssocTypeIDAndCreatedBearInput{
			         AssocType: "string1",
			         AssocID: "string1",
			       StartingAfter: &models.ThingWithMatchingKeys{
			         AssocType:    "string1",
			         AssocID:    "string1",
			         Created: mustTime("2018-03-11T15:04:01+07:00"),
			         Bear: "string1",
			       },
			     },
			   },
			   output: getThingWithMatchingKeyssByAssocTypeIDAndCreatedBearOutput{
			     thingWithMatchingKeyss: []models.ThingWithMatchingKeys{
			       models.ThingWithMatchingKeys{
			         AssocType:    "string1",
			         AssocID:    "string1",
			         Created: mustTime("2018-03-11T15:04:02+07:00"),
			         Bear: "string2",
			       },
			       models.ThingWithMatchingKeys{
			         AssocType:    "string1",
			         AssocID:    "string1",
			         Created: mustTime("2018-03-11T15:04:03+07:00"),
			         Bear: "string3",
			       },
			     },
			     err: nil,
			   },
			 }, */
			{
				testName: "starting after descending",
				d:        d,
				input: getThingWithMatchingKeyssByAssocTypeIDAndCreatedBearInput{
					ctx: context.Background(),
					input: db.GetThingWithMatchingKeyssByAssocTypeIDAndCreatedBearInput{
						AssocType: "string1",
						AssocID:   "string1",
						StartingAfter: &models.ThingWithMatchingKeys{
							AssocType: "string1",
							AssocID:   "string1",
							Created:   mustTime("2018-03-11T15:04:03+07:00"),
							Bear:      "string3",
						},
						Descending: true,
					},
				},
				output: getThingWithMatchingKeyssByAssocTypeIDAndCreatedBearOutput{
					thingWithMatchingKeyss: []models.ThingWithMatchingKeys{
						models.ThingWithMatchingKeys{
							AssocType: "string1",
							AssocID:   "string1",
							Created:   mustTime("2018-03-11T15:04:02+07:00"),
							Bear:      "string2",
						},
						models.ThingWithMatchingKeys{
							AssocType: "string1",
							AssocID:   "string1",
							Created:   mustTime("2018-03-11T15:04:01+07:00"),
							Bear:      "string1",
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting at",
				d:        d,
				input: getThingWithMatchingKeyssByAssocTypeIDAndCreatedBearInput{
					ctx: context.Background(),
					input: db.GetThingWithMatchingKeyssByAssocTypeIDAndCreatedBearInput{
						AssocType: "string1",
						AssocID:   "string1",
						StartingAt: &db.CreatedBear{
							Created: mustTime("2018-03-11T15:04:02+07:00"),
							Bear:    "string2",
						},
					},
				},
				output: getThingWithMatchingKeyssByAssocTypeIDAndCreatedBearOutput{
					thingWithMatchingKeyss: []models.ThingWithMatchingKeys{
						models.ThingWithMatchingKeys{
							AssocType: "string1",
							AssocID:   "string1",
							Created:   mustTime("2018-03-11T15:04:02+07:00"),
							Bear:      "string2",
						},
						models.ThingWithMatchingKeys{
							AssocType: "string1",
							AssocID:   "string1",
							Created:   mustTime("2018-03-11T15:04:03+07:00"),
							Bear:      "string3",
						},
					},
					err: nil,
				},
			},
		}
		for _, test := range tests {
			t.Run(test.testName, test.run)
		}
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanThingWithMatchingKeyssByAssocTypeIDAndCreatedBear(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithMatchingKeys(ctx, models.ThingWithMatchingKeys{
			AssocType: "string1",
			AssocID:   "string1",
			Created:   mustTime("2018-03-11T15:04:01+07:00"),
			Bear:      "string1",
		}))
		require.Nil(t, d.SaveThingWithMatchingKeys(ctx, models.ThingWithMatchingKeys{
			AssocType: "string2",
			AssocID:   "string2",
			Created:   mustTime("2018-03-11T15:04:02+07:00"),
			Bear:      "string2",
		}))
		require.Nil(t, d.SaveThingWithMatchingKeys(ctx, models.ThingWithMatchingKeys{
			AssocType: "string3",
			AssocID:   "string3",
			Created:   mustTime("2018-03-11T15:04:03+07:00"),
			Bear:      "string3",
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.ThingWithMatchingKeys{
				models.ThingWithMatchingKeys{
					AssocType: "string1",
					AssocID:   "string1",
					Created:   mustTime("2018-03-11T15:04:01+07:00"),
					Bear:      "string1",
				},
				models.ThingWithMatchingKeys{
					AssocType: "string2",
					AssocID:   "string2",
					Created:   mustTime("2018-03-11T15:04:02+07:00"),
					Bear:      "string2",
				},
				models.ThingWithMatchingKeys{
					AssocType: "string3",
					AssocID:   "string3",
					Created:   mustTime("2018-03-11T15:04:03+07:00"),
					Bear:      "string3",
				},
			}
			// Consistent read must be disabled when scaning a GSI.
			scanInput := db.ScanThingWithMatchingKeyssByAssocTypeIDAndCreatedBearInput{DisableConsistentRead: true}
			actual := []models.ThingWithMatchingKeys{}
			err := d.ScanThingWithMatchingKeyssByAssocTypeIDAndCreatedBear(ctx, scanInput, func(m *models.ThingWithMatchingKeys, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		t.Run("starting after", func(t *testing.T) {
			// Scan for everything.
			allItems := []models.ThingWithMatchingKeys{}
			// Consistent read must be disabled when scaning a GSI.
			scanInput := db.ScanThingWithMatchingKeyssByAssocTypeIDAndCreatedBearInput{DisableConsistentRead: true}
			err := d.ScanThingWithMatchingKeyssByAssocTypeIDAndCreatedBear(ctx, scanInput, func(m *models.ThingWithMatchingKeys, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput = db.ScanThingWithMatchingKeyssByAssocTypeIDAndCreatedBearInput{
				DisableConsistentRead: true,
				StartingAfter: &models.ThingWithMatchingKeys{
					AssocType: firstItem.AssocType,
					AssocID:   firstItem.AssocID,
					Created:   firstItem.Created,
					Bear:      firstItem.Bear,
				},
			}
			actual := []models.ThingWithMatchingKeys{}
			err = d.ScanThingWithMatchingKeyssByAssocTypeIDAndCreatedBear(ctx, scanInput, func(m *models.ThingWithMatchingKeys, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanThingWithMatchingKeyssInput{
				Limit: &limit,
			}
			actual := []models.ThingWithMatchingKeys{}
			err := d.ScanThingWithMatchingKeyss(ctx, scanInput, func(m *models.ThingWithMatchingKeys, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

func GetThingWithMultiUseCompositeAttribute(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithMultiUseCompositeAttribute{
			Four:  db.String("string1"),
			One:   db.String("string1"),
			Three: db.String("string1"),
			Two:   db.String("string1"),
		}
		require.Nil(t, s.SaveThingWithMultiUseCompositeAttribute(ctx, m))
		m2, err := s.GetThingWithMultiUseCompositeAttribute(ctx, *m.One)
		require.Nil(t, err)
		require.Equal(t, *m.One, *m2.One)

		_, err = s.GetThingWithMultiUseCompositeAttribute(ctx, "string2")
		require.NotNil(t, err)
		require.IsType(t, err, db.ErrThingWithMultiUseCompositeAttributeNotFound{})
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanThingWithMultiUseCompositeAttributes(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithMultiUseCompositeAttribute(ctx, models.ThingWithMultiUseCompositeAttribute{
			Four:  db.String("string1"),
			One:   db.String("string1"),
			Three: db.String("string1"),
			Two:   db.String("string1"),
		}))
		require.Nil(t, d.SaveThingWithMultiUseCompositeAttribute(ctx, models.ThingWithMultiUseCompositeAttribute{
			Four:  db.String("string2"),
			One:   db.String("string2"),
			Three: db.String("string2"),
			Two:   db.String("string2"),
		}))
		require.Nil(t, d.SaveThingWithMultiUseCompositeAttribute(ctx, models.ThingWithMultiUseCompositeAttribute{
			Four:  db.String("string3"),
			One:   db.String("string3"),
			Three: db.String("string3"),
			Two:   db.String("string3"),
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.ThingWithMultiUseCompositeAttribute{
				models.ThingWithMultiUseCompositeAttribute{
					Four:  db.String("string1"),
					One:   db.String("string1"),
					Three: db.String("string1"),
					Two:   db.String("string1"),
				},
				models.ThingWithMultiUseCompositeAttribute{
					Four:  db.String("string2"),
					One:   db.String("string2"),
					Three: db.String("string2"),
					Two:   db.String("string2"),
				},
				models.ThingWithMultiUseCompositeAttribute{
					Four:  db.String("string3"),
					One:   db.String("string3"),
					Three: db.String("string3"),
					Two:   db.String("string3"),
				},
			}
			actual := []models.ThingWithMultiUseCompositeAttribute{}
			err := d.ScanThingWithMultiUseCompositeAttributes(ctx, db.ScanThingWithMultiUseCompositeAttributesInput{}, func(m *models.ThingWithMultiUseCompositeAttribute, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		// FAILING_TEST
		t.Run("starting after", func(t *testing.T) {
			t.Skip()
			// Scan for everything.
			allItems := []models.ThingWithMultiUseCompositeAttribute{}
			err := d.ScanThingWithMultiUseCompositeAttributes(ctx, db.ScanThingWithMultiUseCompositeAttributesInput{}, func(m *models.ThingWithMultiUseCompositeAttribute, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput := db.ScanThingWithMultiUseCompositeAttributesInput{
				StartingAfter: &models.ThingWithMultiUseCompositeAttribute{
					One: firstItem.One,
				},
			}
			actual := []models.ThingWithMultiUseCompositeAttribute{}
			err = d.ScanThingWithMultiUseCompositeAttributes(ctx, scanInput, func(m *models.ThingWithMultiUseCompositeAttribute, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanThingWithMultiUseCompositeAttributesInput{
				Limit: &limit,
			}
			actual := []models.ThingWithMultiUseCompositeAttribute{}
			err := d.ScanThingWithMultiUseCompositeAttributes(ctx, scanInput, func(m *models.ThingWithMultiUseCompositeAttribute, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

func SaveThingWithMultiUseCompositeAttribute(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithMultiUseCompositeAttribute{
			Four:  db.String("string1"),
			One:   db.String("string1"),
			Three: db.String("string1"),
			Two:   db.String("string1"),
		}
		require.Nil(t, s.SaveThingWithMultiUseCompositeAttribute(ctx, m))
	}
}

func DeleteThingWithMultiUseCompositeAttribute(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithMultiUseCompositeAttribute{
			Four:  db.String("string1"),
			One:   db.String("string1"),
			Three: db.String("string1"),
			Two:   db.String("string1"),
		}
		require.Nil(t, s.SaveThingWithMultiUseCompositeAttribute(ctx, m))
		require.Nil(t, s.DeleteThingWithMultiUseCompositeAttribute(ctx, *m.One))
	}
}

type getThingWithMultiUseCompositeAttributesByThreeAndOneTwoInput struct {
	ctx   context.Context
	input db.GetThingWithMultiUseCompositeAttributesByThreeAndOneTwoInput
}
type getThingWithMultiUseCompositeAttributesByThreeAndOneTwoOutput struct {
	thingWithMultiUseCompositeAttributes []models.ThingWithMultiUseCompositeAttribute
	err                                  error
}
type getThingWithMultiUseCompositeAttributesByThreeAndOneTwoTest struct {
	testName string
	d        db.Interface
	input    getThingWithMultiUseCompositeAttributesByThreeAndOneTwoInput
	output   getThingWithMultiUseCompositeAttributesByThreeAndOneTwoOutput
}

func (g getThingWithMultiUseCompositeAttributesByThreeAndOneTwoTest) run(t *testing.T) {
	thingWithMultiUseCompositeAttributes := []models.ThingWithMultiUseCompositeAttribute{}
	fn := func(m *models.ThingWithMultiUseCompositeAttribute, lastThingWithMultiUseCompositeAttribute bool) bool {
		thingWithMultiUseCompositeAttributes = append(thingWithMultiUseCompositeAttributes, *m)
		if lastThingWithMultiUseCompositeAttribute {
			return false
		}
		return true
	}
	err := g.d.GetThingWithMultiUseCompositeAttributesByThreeAndOneTwo(g.input.ctx, g.input.input, fn)
	if err != nil {
		fmt.Println(err.Error())
	}
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.thingWithMultiUseCompositeAttributes, thingWithMultiUseCompositeAttributes)
}

func GetThingWithMultiUseCompositeAttributesByThreeAndOneTwo(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithMultiUseCompositeAttribute(ctx, models.ThingWithMultiUseCompositeAttribute{
			Three: db.String("string1"),
			One:   db.String("string1"),
			Two:   db.String("string1"),
			Four:  db.String("string1"),
		}))
		require.Nil(t, d.SaveThingWithMultiUseCompositeAttribute(ctx, models.ThingWithMultiUseCompositeAttribute{
			Three: db.String("string1"),
			One:   db.String("string2"),
			Two:   db.String("string2"),
			Four:  db.String("string3"),
		}))
		require.Nil(t, d.SaveThingWithMultiUseCompositeAttribute(ctx, models.ThingWithMultiUseCompositeAttribute{
			Three: db.String("string1"),
			One:   db.String("string3"),
			Two:   db.String("string3"),
			Four:  db.String("string2"),
		}))
		limit := int64(3)
		tests := []getThingWithMultiUseCompositeAttributesByThreeAndOneTwoTest{
			{
				testName: "basic",
				d:        d,
				input: getThingWithMultiUseCompositeAttributesByThreeAndOneTwoInput{
					ctx: context.Background(),
					input: db.GetThingWithMultiUseCompositeAttributesByThreeAndOneTwoInput{
						Three: "string1",
						Limit: &limit,
					},
				},
				output: getThingWithMultiUseCompositeAttributesByThreeAndOneTwoOutput{
					thingWithMultiUseCompositeAttributes: []models.ThingWithMultiUseCompositeAttribute{
						models.ThingWithMultiUseCompositeAttribute{
							Three: db.String("string1"),
							One:   db.String("string1"),
							Two:   db.String("string1"),
							Four:  db.String("string1"),
						},
						models.ThingWithMultiUseCompositeAttribute{
							Three: db.String("string1"),
							One:   db.String("string2"),
							Two:   db.String("string2"),
							Four:  db.String("string3"),
						},
						models.ThingWithMultiUseCompositeAttribute{
							Three: db.String("string1"),
							One:   db.String("string3"),
							Two:   db.String("string3"),
							Four:  db.String("string2"),
						},
					},
					err: nil,
				},
			},
			{
				testName: "descending",
				d:        d,
				input: getThingWithMultiUseCompositeAttributesByThreeAndOneTwoInput{
					ctx: context.Background(),
					input: db.GetThingWithMultiUseCompositeAttributesByThreeAndOneTwoInput{
						Three:      "string1",
						Descending: true,
					},
				},
				output: getThingWithMultiUseCompositeAttributesByThreeAndOneTwoOutput{
					thingWithMultiUseCompositeAttributes: []models.ThingWithMultiUseCompositeAttribute{
						models.ThingWithMultiUseCompositeAttribute{
							Three: db.String("string1"),
							One:   db.String("string3"),
							Two:   db.String("string3"),
							Four:  db.String("string2"),
						},
						models.ThingWithMultiUseCompositeAttribute{
							Three: db.String("string1"),
							One:   db.String("string2"),
							Two:   db.String("string2"),
							Four:  db.String("string3"),
						},
						models.ThingWithMultiUseCompositeAttribute{
							Three: db.String("string1"),
							One:   db.String("string1"),
							Two:   db.String("string1"),
							Four:  db.String("string1"),
						},
					},
					err: nil,
				},
			},
			/* FAILING_TEST */
			/* {
			   testName: "starting after",
			   d:    d,
			   input: getThingWithMultiUseCompositeAttributesByThreeAndOneTwoInput{
			     ctx: context.Background(),
			     input: db.GetThingWithMultiUseCompositeAttributesByThreeAndOneTwoInput{
			         Three: "string1",
			       StartingAfter: &models.ThingWithMultiUseCompositeAttribute{
			         Three:    db.String("string1"),
			         One: db.String("string1"),
			         Two: db.String("string1"),
			         Four:    db.String("string1"),
			       },
			     },
			   },
			   output: getThingWithMultiUseCompositeAttributesByThreeAndOneTwoOutput{
			     thingWithMultiUseCompositeAttributes: []models.ThingWithMultiUseCompositeAttribute{
			       models.ThingWithMultiUseCompositeAttribute{
			         Three:    db.String("string1"),
			         One: db.String("string2"),
			         Two: db.String("string2"),
			         Four:    db.String("string3"),
			       },
			       models.ThingWithMultiUseCompositeAttribute{
			         Three:    db.String("string1"),
			         One: db.String("string3"),
			         Two: db.String("string3"),
			         Four:    db.String("string2"),
			       },
			     },
			     err: nil,
			   },
			 }, */
			{
				testName: "starting after descending",
				d:        d,
				input: getThingWithMultiUseCompositeAttributesByThreeAndOneTwoInput{
					ctx: context.Background(),
					input: db.GetThingWithMultiUseCompositeAttributesByThreeAndOneTwoInput{
						Three: "string1",
						StartingAfter: &models.ThingWithMultiUseCompositeAttribute{
							Three: db.String("string1"),
							One:   db.String("string3"),
							Two:   db.String("string3"),
							Four:  db.String("string2"),
						},
						Descending: true,
					},
				},
				output: getThingWithMultiUseCompositeAttributesByThreeAndOneTwoOutput{
					thingWithMultiUseCompositeAttributes: []models.ThingWithMultiUseCompositeAttribute{
						models.ThingWithMultiUseCompositeAttribute{
							Three: db.String("string1"),
							One:   db.String("string2"),
							Two:   db.String("string2"),
							Four:  db.String("string3"),
						},
						models.ThingWithMultiUseCompositeAttribute{
							Three: db.String("string1"),
							One:   db.String("string1"),
							Two:   db.String("string1"),
							Four:  db.String("string1"),
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting at",
				d:        d,
				input: getThingWithMultiUseCompositeAttributesByThreeAndOneTwoInput{
					ctx: context.Background(),
					input: db.GetThingWithMultiUseCompositeAttributesByThreeAndOneTwoInput{
						Three: "string1",
						StartingAt: &db.OneTwo{
							One: "string2",
							Two: "string2",
						},
					},
				},
				output: getThingWithMultiUseCompositeAttributesByThreeAndOneTwoOutput{
					thingWithMultiUseCompositeAttributes: []models.ThingWithMultiUseCompositeAttribute{
						models.ThingWithMultiUseCompositeAttribute{
							Three: db.String("string1"),
							One:   db.String("string2"),
							Two:   db.String("string2"),
							Four:  db.String("string3"),
						},
						models.ThingWithMultiUseCompositeAttribute{
							Three: db.String("string1"),
							One:   db.String("string3"),
							Two:   db.String("string3"),
							Four:  db.String("string2"),
						},
					},
					err: nil,
				},
			},
		}
		for _, test := range tests {
			t.Run(test.testName, test.run)
		}
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanThingWithMultiUseCompositeAttributesByThreeAndOneTwo(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithMultiUseCompositeAttribute(ctx, models.ThingWithMultiUseCompositeAttribute{
			Three: db.String("string1"),
			One:   db.String("string1"),
			Two:   db.String("string1"),
			Four:  db.String("string1"),
		}))
		require.Nil(t, d.SaveThingWithMultiUseCompositeAttribute(ctx, models.ThingWithMultiUseCompositeAttribute{
			Three: db.String("string2"),
			One:   db.String("string2"),
			Two:   db.String("string2"),
			Four:  db.String("string2"),
		}))
		require.Nil(t, d.SaveThingWithMultiUseCompositeAttribute(ctx, models.ThingWithMultiUseCompositeAttribute{
			Three: db.String("string3"),
			One:   db.String("string3"),
			Two:   db.String("string3"),
			Four:  db.String("string3"),
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.ThingWithMultiUseCompositeAttribute{
				models.ThingWithMultiUseCompositeAttribute{
					Three: db.String("string1"),
					One:   db.String("string1"),
					Two:   db.String("string1"),
					Four:  db.String("string1"),
				},
				models.ThingWithMultiUseCompositeAttribute{
					Three: db.String("string2"),
					One:   db.String("string2"),
					Two:   db.String("string2"),
					Four:  db.String("string2"),
				},
				models.ThingWithMultiUseCompositeAttribute{
					Three: db.String("string3"),
					One:   db.String("string3"),
					Two:   db.String("string3"),
					Four:  db.String("string3"),
				},
			}
			// Consistent read must be disabled when scaning a GSI.
			scanInput := db.ScanThingWithMultiUseCompositeAttributesByThreeAndOneTwoInput{DisableConsistentRead: true}
			actual := []models.ThingWithMultiUseCompositeAttribute{}
			err := d.ScanThingWithMultiUseCompositeAttributesByThreeAndOneTwo(ctx, scanInput, func(m *models.ThingWithMultiUseCompositeAttribute, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		t.Run("starting after", func(t *testing.T) {
			// Scan for everything.
			allItems := []models.ThingWithMultiUseCompositeAttribute{}
			// Consistent read must be disabled when scaning a GSI.
			scanInput := db.ScanThingWithMultiUseCompositeAttributesByThreeAndOneTwoInput{DisableConsistentRead: true}
			err := d.ScanThingWithMultiUseCompositeAttributesByThreeAndOneTwo(ctx, scanInput, func(m *models.ThingWithMultiUseCompositeAttribute, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput = db.ScanThingWithMultiUseCompositeAttributesByThreeAndOneTwoInput{
				DisableConsistentRead: true,
				StartingAfter: &models.ThingWithMultiUseCompositeAttribute{
					Three: firstItem.Three,
					One:   firstItem.One,
					Two:   firstItem.Two,
				},
			}
			actual := []models.ThingWithMultiUseCompositeAttribute{}
			err = d.ScanThingWithMultiUseCompositeAttributesByThreeAndOneTwo(ctx, scanInput, func(m *models.ThingWithMultiUseCompositeAttribute, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanThingWithMultiUseCompositeAttributesInput{
				Limit: &limit,
			}
			actual := []models.ThingWithMultiUseCompositeAttribute{}
			err := d.ScanThingWithMultiUseCompositeAttributes(ctx, scanInput, func(m *models.ThingWithMultiUseCompositeAttribute, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

type getThingWithMultiUseCompositeAttributesByFourAndOneTwoInput struct {
	ctx   context.Context
	input db.GetThingWithMultiUseCompositeAttributesByFourAndOneTwoInput
}
type getThingWithMultiUseCompositeAttributesByFourAndOneTwoOutput struct {
	thingWithMultiUseCompositeAttributes []models.ThingWithMultiUseCompositeAttribute
	err                                  error
}
type getThingWithMultiUseCompositeAttributesByFourAndOneTwoTest struct {
	testName string
	d        db.Interface
	input    getThingWithMultiUseCompositeAttributesByFourAndOneTwoInput
	output   getThingWithMultiUseCompositeAttributesByFourAndOneTwoOutput
}

func (g getThingWithMultiUseCompositeAttributesByFourAndOneTwoTest) run(t *testing.T) {
	thingWithMultiUseCompositeAttributes := []models.ThingWithMultiUseCompositeAttribute{}
	fn := func(m *models.ThingWithMultiUseCompositeAttribute, lastThingWithMultiUseCompositeAttribute bool) bool {
		thingWithMultiUseCompositeAttributes = append(thingWithMultiUseCompositeAttributes, *m)
		if lastThingWithMultiUseCompositeAttribute {
			return false
		}
		return true
	}
	err := g.d.GetThingWithMultiUseCompositeAttributesByFourAndOneTwo(g.input.ctx, g.input.input, fn)
	if err != nil {
		fmt.Println(err.Error())
	}
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.thingWithMultiUseCompositeAttributes, thingWithMultiUseCompositeAttributes)
}

func GetThingWithMultiUseCompositeAttributesByFourAndOneTwo(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithMultiUseCompositeAttribute(ctx, models.ThingWithMultiUseCompositeAttribute{
			Four:  db.String("string1"),
			One:   db.String("string1"),
			Two:   db.String("string1"),
			Three: db.String("string1"),
		}))
		require.Nil(t, d.SaveThingWithMultiUseCompositeAttribute(ctx, models.ThingWithMultiUseCompositeAttribute{
			Four:  db.String("string1"),
			One:   db.String("string2"),
			Two:   db.String("string2"),
			Three: db.String("string3"),
		}))
		require.Nil(t, d.SaveThingWithMultiUseCompositeAttribute(ctx, models.ThingWithMultiUseCompositeAttribute{
			Four:  db.String("string1"),
			One:   db.String("string3"),
			Two:   db.String("string3"),
			Three: db.String("string2"),
		}))
		limit := int64(3)
		tests := []getThingWithMultiUseCompositeAttributesByFourAndOneTwoTest{
			{
				testName: "basic",
				d:        d,
				input: getThingWithMultiUseCompositeAttributesByFourAndOneTwoInput{
					ctx: context.Background(),
					input: db.GetThingWithMultiUseCompositeAttributesByFourAndOneTwoInput{
						Four:  "string1",
						Limit: &limit,
					},
				},
				output: getThingWithMultiUseCompositeAttributesByFourAndOneTwoOutput{
					thingWithMultiUseCompositeAttributes: []models.ThingWithMultiUseCompositeAttribute{
						models.ThingWithMultiUseCompositeAttribute{
							Four:  db.String("string1"),
							One:   db.String("string1"),
							Two:   db.String("string1"),
							Three: db.String("string1"),
						},
						models.ThingWithMultiUseCompositeAttribute{
							Four:  db.String("string1"),
							One:   db.String("string2"),
							Two:   db.String("string2"),
							Three: db.String("string3"),
						},
						models.ThingWithMultiUseCompositeAttribute{
							Four:  db.String("string1"),
							One:   db.String("string3"),
							Two:   db.String("string3"),
							Three: db.String("string2"),
						},
					},
					err: nil,
				},
			},
			{
				testName: "descending",
				d:        d,
				input: getThingWithMultiUseCompositeAttributesByFourAndOneTwoInput{
					ctx: context.Background(),
					input: db.GetThingWithMultiUseCompositeAttributesByFourAndOneTwoInput{
						Four:       "string1",
						Descending: true,
					},
				},
				output: getThingWithMultiUseCompositeAttributesByFourAndOneTwoOutput{
					thingWithMultiUseCompositeAttributes: []models.ThingWithMultiUseCompositeAttribute{
						models.ThingWithMultiUseCompositeAttribute{
							Four:  db.String("string1"),
							One:   db.String("string3"),
							Two:   db.String("string3"),
							Three: db.String("string2"),
						},
						models.ThingWithMultiUseCompositeAttribute{
							Four:  db.String("string1"),
							One:   db.String("string2"),
							Two:   db.String("string2"),
							Three: db.String("string3"),
						},
						models.ThingWithMultiUseCompositeAttribute{
							Four:  db.String("string1"),
							One:   db.String("string1"),
							Two:   db.String("string1"),
							Three: db.String("string1"),
						},
					},
					err: nil,
				},
			},
			/* FAILING_TEST */
			/* {
			   testName: "starting after",
			   d:    d,
			   input: getThingWithMultiUseCompositeAttributesByFourAndOneTwoInput{
			     ctx: context.Background(),
			     input: db.GetThingWithMultiUseCompositeAttributesByFourAndOneTwoInput{
			         Four: "string1",
			       StartingAfter: &models.ThingWithMultiUseCompositeAttribute{
			         Four:    db.String("string1"),
			         One: db.String("string1"),
			         Two: db.String("string1"),
			         Three:    db.String("string1"),
			       },
			     },
			   },
			   output: getThingWithMultiUseCompositeAttributesByFourAndOneTwoOutput{
			     thingWithMultiUseCompositeAttributes: []models.ThingWithMultiUseCompositeAttribute{
			       models.ThingWithMultiUseCompositeAttribute{
			         Four:    db.String("string1"),
			         One: db.String("string2"),
			         Two: db.String("string2"),
			         Three:    db.String("string3"),
			       },
			       models.ThingWithMultiUseCompositeAttribute{
			         Four:    db.String("string1"),
			         One: db.String("string3"),
			         Two: db.String("string3"),
			         Three:    db.String("string2"),
			       },
			     },
			     err: nil,
			   },
			 }, */
			{
				testName: "starting after descending",
				d:        d,
				input: getThingWithMultiUseCompositeAttributesByFourAndOneTwoInput{
					ctx: context.Background(),
					input: db.GetThingWithMultiUseCompositeAttributesByFourAndOneTwoInput{
						Four: "string1",
						StartingAfter: &models.ThingWithMultiUseCompositeAttribute{
							Four:  db.String("string1"),
							One:   db.String("string3"),
							Two:   db.String("string3"),
							Three: db.String("string2"),
						},
						Descending: true,
					},
				},
				output: getThingWithMultiUseCompositeAttributesByFourAndOneTwoOutput{
					thingWithMultiUseCompositeAttributes: []models.ThingWithMultiUseCompositeAttribute{
						models.ThingWithMultiUseCompositeAttribute{
							Four:  db.String("string1"),
							One:   db.String("string2"),
							Two:   db.String("string2"),
							Three: db.String("string3"),
						},
						models.ThingWithMultiUseCompositeAttribute{
							Four:  db.String("string1"),
							One:   db.String("string1"),
							Two:   db.String("string1"),
							Three: db.String("string1"),
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting at",
				d:        d,
				input: getThingWithMultiUseCompositeAttributesByFourAndOneTwoInput{
					ctx: context.Background(),
					input: db.GetThingWithMultiUseCompositeAttributesByFourAndOneTwoInput{
						Four: "string1",
						StartingAt: &db.OneTwo{
							One: "string2",
							Two: "string2",
						},
					},
				},
				output: getThingWithMultiUseCompositeAttributesByFourAndOneTwoOutput{
					thingWithMultiUseCompositeAttributes: []models.ThingWithMultiUseCompositeAttribute{
						models.ThingWithMultiUseCompositeAttribute{
							Four:  db.String("string1"),
							One:   db.String("string2"),
							Two:   db.String("string2"),
							Three: db.String("string3"),
						},
						models.ThingWithMultiUseCompositeAttribute{
							Four:  db.String("string1"),
							One:   db.String("string3"),
							Two:   db.String("string3"),
							Three: db.String("string2"),
						},
					},
					err: nil,
				},
			},
		}
		for _, test := range tests {
			t.Run(test.testName, test.run)
		}
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanThingWithMultiUseCompositeAttributesByFourAndOneTwo(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithMultiUseCompositeAttribute(ctx, models.ThingWithMultiUseCompositeAttribute{
			Four:  db.String("string1"),
			One:   db.String("string1"),
			Two:   db.String("string1"),
			Three: db.String("string1"),
		}))
		require.Nil(t, d.SaveThingWithMultiUseCompositeAttribute(ctx, models.ThingWithMultiUseCompositeAttribute{
			Four:  db.String("string2"),
			One:   db.String("string2"),
			Two:   db.String("string2"),
			Three: db.String("string2"),
		}))
		require.Nil(t, d.SaveThingWithMultiUseCompositeAttribute(ctx, models.ThingWithMultiUseCompositeAttribute{
			Four:  db.String("string3"),
			One:   db.String("string3"),
			Two:   db.String("string3"),
			Three: db.String("string3"),
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.ThingWithMultiUseCompositeAttribute{
				models.ThingWithMultiUseCompositeAttribute{
					Four:  db.String("string1"),
					One:   db.String("string1"),
					Two:   db.String("string1"),
					Three: db.String("string1"),
				},
				models.ThingWithMultiUseCompositeAttribute{
					Four:  db.String("string2"),
					One:   db.String("string2"),
					Two:   db.String("string2"),
					Three: db.String("string2"),
				},
				models.ThingWithMultiUseCompositeAttribute{
					Four:  db.String("string3"),
					One:   db.String("string3"),
					Two:   db.String("string3"),
					Three: db.String("string3"),
				},
			}
			// Consistent read must be disabled when scaning a GSI.
			scanInput := db.ScanThingWithMultiUseCompositeAttributesByFourAndOneTwoInput{DisableConsistentRead: true}
			actual := []models.ThingWithMultiUseCompositeAttribute{}
			err := d.ScanThingWithMultiUseCompositeAttributesByFourAndOneTwo(ctx, scanInput, func(m *models.ThingWithMultiUseCompositeAttribute, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		t.Run("starting after", func(t *testing.T) {
			// Scan for everything.
			allItems := []models.ThingWithMultiUseCompositeAttribute{}
			// Consistent read must be disabled when scaning a GSI.
			scanInput := db.ScanThingWithMultiUseCompositeAttributesByFourAndOneTwoInput{DisableConsistentRead: true}
			err := d.ScanThingWithMultiUseCompositeAttributesByFourAndOneTwo(ctx, scanInput, func(m *models.ThingWithMultiUseCompositeAttribute, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput = db.ScanThingWithMultiUseCompositeAttributesByFourAndOneTwoInput{
				DisableConsistentRead: true,
				StartingAfter: &models.ThingWithMultiUseCompositeAttribute{
					Four: firstItem.Four,
					One:  firstItem.One,
					Two:  firstItem.Two,
				},
			}
			actual := []models.ThingWithMultiUseCompositeAttribute{}
			err = d.ScanThingWithMultiUseCompositeAttributesByFourAndOneTwo(ctx, scanInput, func(m *models.ThingWithMultiUseCompositeAttribute, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanThingWithMultiUseCompositeAttributesInput{
				Limit: &limit,
			}
			actual := []models.ThingWithMultiUseCompositeAttribute{}
			err := d.ScanThingWithMultiUseCompositeAttributes(ctx, scanInput, func(m *models.ThingWithMultiUseCompositeAttribute, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

func GetThingWithRequiredCompositePropertiesAndKeysOnly(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithRequiredCompositePropertiesAndKeysOnly{
			PropertyOne:   db.String("string1"),
			PropertyThree: db.String("string1"),
			PropertyTwo:   db.String("string1"),
		}
		require.Nil(t, s.SaveThingWithRequiredCompositePropertiesAndKeysOnly(ctx, m))
		m2, err := s.GetThingWithRequiredCompositePropertiesAndKeysOnly(ctx, *m.PropertyThree)
		require.Nil(t, err)
		require.Equal(t, *m.PropertyThree, *m2.PropertyThree)

		_, err = s.GetThingWithRequiredCompositePropertiesAndKeysOnly(ctx, "string2")
		require.NotNil(t, err)
		require.IsType(t, err, db.ErrThingWithRequiredCompositePropertiesAndKeysOnlyNotFound{})
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanThingWithRequiredCompositePropertiesAndKeysOnlys(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithRequiredCompositePropertiesAndKeysOnly(ctx, models.ThingWithRequiredCompositePropertiesAndKeysOnly{
			PropertyOne:   db.String("string1"),
			PropertyThree: db.String("string1"),
			PropertyTwo:   db.String("string1"),
		}))
		require.Nil(t, d.SaveThingWithRequiredCompositePropertiesAndKeysOnly(ctx, models.ThingWithRequiredCompositePropertiesAndKeysOnly{
			PropertyOne:   db.String("string2"),
			PropertyThree: db.String("string2"),
			PropertyTwo:   db.String("string2"),
		}))
		require.Nil(t, d.SaveThingWithRequiredCompositePropertiesAndKeysOnly(ctx, models.ThingWithRequiredCompositePropertiesAndKeysOnly{
			PropertyOne:   db.String("string3"),
			PropertyThree: db.String("string3"),
			PropertyTwo:   db.String("string3"),
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.ThingWithRequiredCompositePropertiesAndKeysOnly{
				models.ThingWithRequiredCompositePropertiesAndKeysOnly{
					PropertyOne:   db.String("string1"),
					PropertyThree: db.String("string1"),
					PropertyTwo:   db.String("string1"),
				},
				models.ThingWithRequiredCompositePropertiesAndKeysOnly{
					PropertyOne:   db.String("string2"),
					PropertyThree: db.String("string2"),
					PropertyTwo:   db.String("string2"),
				},
				models.ThingWithRequiredCompositePropertiesAndKeysOnly{
					PropertyOne:   db.String("string3"),
					PropertyThree: db.String("string3"),
					PropertyTwo:   db.String("string3"),
				},
			}
			actual := []models.ThingWithRequiredCompositePropertiesAndKeysOnly{}
			err := d.ScanThingWithRequiredCompositePropertiesAndKeysOnlys(ctx, db.ScanThingWithRequiredCompositePropertiesAndKeysOnlysInput{}, func(m *models.ThingWithRequiredCompositePropertiesAndKeysOnly, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		// FAILING_TEST
		t.Run("starting after", func(t *testing.T) {
			t.Skip()
			// Scan for everything.
			allItems := []models.ThingWithRequiredCompositePropertiesAndKeysOnly{}
			err := d.ScanThingWithRequiredCompositePropertiesAndKeysOnlys(ctx, db.ScanThingWithRequiredCompositePropertiesAndKeysOnlysInput{}, func(m *models.ThingWithRequiredCompositePropertiesAndKeysOnly, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput := db.ScanThingWithRequiredCompositePropertiesAndKeysOnlysInput{
				StartingAfter: &models.ThingWithRequiredCompositePropertiesAndKeysOnly{
					PropertyThree: firstItem.PropertyThree,
					// must specify non-empty string values for attributes
					// in secondary indexes, since dynamodb doesn't support
					// empty strings:
					PropertyOne: pointerToString("propertyOne"),
					PropertyTwo: pointerToString("propertyTwo"),
				},
			}
			actual := []models.ThingWithRequiredCompositePropertiesAndKeysOnly{}
			err = d.ScanThingWithRequiredCompositePropertiesAndKeysOnlys(ctx, scanInput, func(m *models.ThingWithRequiredCompositePropertiesAndKeysOnly, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanThingWithRequiredCompositePropertiesAndKeysOnlysInput{
				Limit: &limit,
			}
			actual := []models.ThingWithRequiredCompositePropertiesAndKeysOnly{}
			err := d.ScanThingWithRequiredCompositePropertiesAndKeysOnlys(ctx, scanInput, func(m *models.ThingWithRequiredCompositePropertiesAndKeysOnly, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

func SaveThingWithRequiredCompositePropertiesAndKeysOnly(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithRequiredCompositePropertiesAndKeysOnly{
			PropertyOne:   db.String("string1"),
			PropertyThree: db.String("string1"),
			PropertyTwo:   db.String("string1"),
		}
		require.Nil(t, s.SaveThingWithRequiredCompositePropertiesAndKeysOnly(ctx, m))
	}
}

func DeleteThingWithRequiredCompositePropertiesAndKeysOnly(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithRequiredCompositePropertiesAndKeysOnly{
			PropertyOne:   db.String("string1"),
			PropertyThree: db.String("string1"),
			PropertyTwo:   db.String("string1"),
		}
		require.Nil(t, s.SaveThingWithRequiredCompositePropertiesAndKeysOnly(ctx, m))
		require.Nil(t, s.DeleteThingWithRequiredCompositePropertiesAndKeysOnly(ctx, *m.PropertyThree))
	}
}

type getThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThreeInput struct {
	ctx   context.Context
	input db.GetThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThreeInput
}
type getThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThreeOutput struct {
	thingWithRequiredCompositePropertiesAndKeysOnlys []models.ThingWithRequiredCompositePropertiesAndKeysOnly
	err                                              error
}
type getThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThreeTest struct {
	testName string
	d        db.Interface
	input    getThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThreeInput
	output   getThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThreeOutput
}

func (g getThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThreeTest) run(t *testing.T) {
	thingWithRequiredCompositePropertiesAndKeysOnlys := []models.ThingWithRequiredCompositePropertiesAndKeysOnly{}
	fn := func(m *models.ThingWithRequiredCompositePropertiesAndKeysOnly, lastThingWithRequiredCompositePropertiesAndKeysOnly bool) bool {
		thingWithRequiredCompositePropertiesAndKeysOnlys = append(thingWithRequiredCompositePropertiesAndKeysOnlys, *m)
		if lastThingWithRequiredCompositePropertiesAndKeysOnly {
			return false
		}
		return true
	}
	err := g.d.GetThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThree(g.input.ctx, g.input.input, fn)
	if err != nil {
		fmt.Println(err.Error())
	}
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.thingWithRequiredCompositePropertiesAndKeysOnlys, thingWithRequiredCompositePropertiesAndKeysOnlys)
}

func GetThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThree(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithRequiredCompositePropertiesAndKeysOnly(ctx, models.ThingWithRequiredCompositePropertiesAndKeysOnly{
			PropertyOne:   db.String("string1"),
			PropertyTwo:   db.String("string1"),
			PropertyThree: db.String("string1"),
		}))
		require.Nil(t, d.SaveThingWithRequiredCompositePropertiesAndKeysOnly(ctx, models.ThingWithRequiredCompositePropertiesAndKeysOnly{
			PropertyOne:   db.String("string1"),
			PropertyTwo:   db.String("string1"),
			PropertyThree: db.String("string2"),
		}))
		require.Nil(t, d.SaveThingWithRequiredCompositePropertiesAndKeysOnly(ctx, models.ThingWithRequiredCompositePropertiesAndKeysOnly{
			PropertyOne:   db.String("string1"),
			PropertyTwo:   db.String("string1"),
			PropertyThree: db.String("string3"),
		}))
		limit := int64(3)
		tests := []getThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThreeTest{
			{
				testName: "basic",
				d:        d,
				input: getThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThreeInput{
					ctx: context.Background(),
					input: db.GetThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThreeInput{
						PropertyOne: "string1",
						PropertyTwo: "string1",
						Limit:       &limit,
					},
				},
				output: getThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThreeOutput{
					thingWithRequiredCompositePropertiesAndKeysOnlys: []models.ThingWithRequiredCompositePropertiesAndKeysOnly{
						models.ThingWithRequiredCompositePropertiesAndKeysOnly{
							PropertyOne:   db.String("string1"),
							PropertyTwo:   db.String("string1"),
							PropertyThree: db.String("string1"),
						},
						models.ThingWithRequiredCompositePropertiesAndKeysOnly{
							PropertyOne:   db.String("string1"),
							PropertyTwo:   db.String("string1"),
							PropertyThree: db.String("string2"),
						},
						models.ThingWithRequiredCompositePropertiesAndKeysOnly{
							PropertyOne:   db.String("string1"),
							PropertyTwo:   db.String("string1"),
							PropertyThree: db.String("string3"),
						},
					},
					err: nil,
				},
			},
			{
				testName: "descending",
				d:        d,
				input: getThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThreeInput{
					ctx: context.Background(),
					input: db.GetThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThreeInput{
						PropertyOne: "string1",
						PropertyTwo: "string1",
						Descending:  true,
					},
				},
				output: getThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThreeOutput{
					thingWithRequiredCompositePropertiesAndKeysOnlys: []models.ThingWithRequiredCompositePropertiesAndKeysOnly{
						models.ThingWithRequiredCompositePropertiesAndKeysOnly{
							PropertyOne:   db.String("string1"),
							PropertyTwo:   db.String("string1"),
							PropertyThree: db.String("string3"),
						},
						models.ThingWithRequiredCompositePropertiesAndKeysOnly{
							PropertyOne:   db.String("string1"),
							PropertyTwo:   db.String("string1"),
							PropertyThree: db.String("string2"),
						},
						models.ThingWithRequiredCompositePropertiesAndKeysOnly{
							PropertyOne:   db.String("string1"),
							PropertyTwo:   db.String("string1"),
							PropertyThree: db.String("string1"),
						},
					},
					err: nil,
				},
			},
			/* FAILING_TEST */
			/* {
			   testName: "starting after",
			   d:    d,
			   input: getThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThreeInput{
			     ctx: context.Background(),
			     input: db.GetThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThreeInput{
			         PropertyOne: "string1",
			         PropertyTwo: "string1",
			       StartingAfter: &models.ThingWithRequiredCompositePropertiesAndKeysOnly{
			         PropertyOne:    db.String("string1"),
			         PropertyTwo:    db.String("string1"),
			         PropertyThree: db.String("string1"),
			       },
			     },
			   },
			   output: getThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThreeOutput{
			     thingWithRequiredCompositePropertiesAndKeysOnlys: []models.ThingWithRequiredCompositePropertiesAndKeysOnly{
			       models.ThingWithRequiredCompositePropertiesAndKeysOnly{
			         PropertyOne:    db.String("string1"),
			         PropertyTwo:    db.String("string1"),
			         PropertyThree: db.String("string2"),
			       },
			       models.ThingWithRequiredCompositePropertiesAndKeysOnly{
			         PropertyOne:    db.String("string1"),
			         PropertyTwo:    db.String("string1"),
			         PropertyThree: db.String("string3"),
			       },
			     },
			     err: nil,
			   },
			 }, */
			{
				testName: "starting after descending",
				d:        d,
				input: getThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThreeInput{
					ctx: context.Background(),
					input: db.GetThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThreeInput{
						PropertyOne: "string1",
						PropertyTwo: "string1",
						StartingAfter: &models.ThingWithRequiredCompositePropertiesAndKeysOnly{
							PropertyOne:   db.String("string1"),
							PropertyTwo:   db.String("string1"),
							PropertyThree: db.String("string3"),
						},
						Descending: true,
					},
				},
				output: getThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThreeOutput{
					thingWithRequiredCompositePropertiesAndKeysOnlys: []models.ThingWithRequiredCompositePropertiesAndKeysOnly{
						models.ThingWithRequiredCompositePropertiesAndKeysOnly{
							PropertyOne:   db.String("string1"),
							PropertyTwo:   db.String("string1"),
							PropertyThree: db.String("string2"),
						},
						models.ThingWithRequiredCompositePropertiesAndKeysOnly{
							PropertyOne:   db.String("string1"),
							PropertyTwo:   db.String("string1"),
							PropertyThree: db.String("string1"),
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting at",
				d:        d,
				input: getThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThreeInput{
					ctx: context.Background(),
					input: db.GetThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThreeInput{
						PropertyOne:             "string1",
						PropertyTwo:             "string1",
						PropertyThreeStartingAt: db.String("string2"),
					},
				},
				output: getThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThreeOutput{
					thingWithRequiredCompositePropertiesAndKeysOnlys: []models.ThingWithRequiredCompositePropertiesAndKeysOnly{
						models.ThingWithRequiredCompositePropertiesAndKeysOnly{
							PropertyOne:   db.String("string1"),
							PropertyTwo:   db.String("string1"),
							PropertyThree: db.String("string2"),
						},
						models.ThingWithRequiredCompositePropertiesAndKeysOnly{
							PropertyOne:   db.String("string1"),
							PropertyTwo:   db.String("string1"),
							PropertyThree: db.String("string3"),
						},
					},
					err: nil,
				},
			},
		}
		for _, test := range tests {
			t.Run(test.testName, test.run)
		}
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThree(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithRequiredCompositePropertiesAndKeysOnly(ctx, models.ThingWithRequiredCompositePropertiesAndKeysOnly{
			PropertyOne:   db.String("string1"),
			PropertyTwo:   db.String("string1"),
			PropertyThree: db.String("string1"),
		}))
		require.Nil(t, d.SaveThingWithRequiredCompositePropertiesAndKeysOnly(ctx, models.ThingWithRequiredCompositePropertiesAndKeysOnly{
			PropertyOne:   db.String("string2"),
			PropertyTwo:   db.String("string2"),
			PropertyThree: db.String("string2"),
		}))
		require.Nil(t, d.SaveThingWithRequiredCompositePropertiesAndKeysOnly(ctx, models.ThingWithRequiredCompositePropertiesAndKeysOnly{
			PropertyOne:   db.String("string3"),
			PropertyTwo:   db.String("string3"),
			PropertyThree: db.String("string3"),
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.ThingWithRequiredCompositePropertiesAndKeysOnly{
				models.ThingWithRequiredCompositePropertiesAndKeysOnly{
					PropertyOne:   db.String("string1"),
					PropertyTwo:   db.String("string1"),
					PropertyThree: db.String("string1"),
				},
				models.ThingWithRequiredCompositePropertiesAndKeysOnly{
					PropertyOne:   db.String("string2"),
					PropertyTwo:   db.String("string2"),
					PropertyThree: db.String("string2"),
				},
				models.ThingWithRequiredCompositePropertiesAndKeysOnly{
					PropertyOne:   db.String("string3"),
					PropertyTwo:   db.String("string3"),
					PropertyThree: db.String("string3"),
				},
			}
			// Consistent read must be disabled when scaning a GSI.
			scanInput := db.ScanThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThreeInput{DisableConsistentRead: true}
			actual := []models.ThingWithRequiredCompositePropertiesAndKeysOnly{}
			err := d.ScanThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThree(ctx, scanInput, func(m *models.ThingWithRequiredCompositePropertiesAndKeysOnly, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		t.Run("starting after", func(t *testing.T) {
			// Scan for everything.
			allItems := []models.ThingWithRequiredCompositePropertiesAndKeysOnly{}
			// Consistent read must be disabled when scaning a GSI.
			scanInput := db.ScanThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThreeInput{DisableConsistentRead: true}
			err := d.ScanThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThree(ctx, scanInput, func(m *models.ThingWithRequiredCompositePropertiesAndKeysOnly, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput = db.ScanThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThreeInput{
				DisableConsistentRead: true,
				StartingAfter: &models.ThingWithRequiredCompositePropertiesAndKeysOnly{
					PropertyOne:   firstItem.PropertyOne,
					PropertyTwo:   firstItem.PropertyTwo,
					PropertyThree: firstItem.PropertyThree,
				},
			}
			actual := []models.ThingWithRequiredCompositePropertiesAndKeysOnly{}
			err = d.ScanThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThree(ctx, scanInput, func(m *models.ThingWithRequiredCompositePropertiesAndKeysOnly, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanThingWithRequiredCompositePropertiesAndKeysOnlysInput{
				Limit: &limit,
			}
			actual := []models.ThingWithRequiredCompositePropertiesAndKeysOnly{}
			err := d.ScanThingWithRequiredCompositePropertiesAndKeysOnlys(ctx, scanInput, func(m *models.ThingWithRequiredCompositePropertiesAndKeysOnly, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

func GetThingWithRequiredFields(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithRequiredFields{
			Name: db.String("string1"),
		}
		require.Nil(t, s.SaveThingWithRequiredFields(ctx, m))
		m2, err := s.GetThingWithRequiredFields(ctx, *m.Name)
		require.Nil(t, err)
		require.Equal(t, *m.Name, *m2.Name)

		_, err = s.GetThingWithRequiredFields(ctx, "string2")
		require.NotNil(t, err)
		require.IsType(t, err, db.ErrThingWithRequiredFieldsNotFound{})
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanThingWithRequiredFieldss(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithRequiredFields(ctx, models.ThingWithRequiredFields{
			Name: db.String("string1"),
		}))
		require.Nil(t, d.SaveThingWithRequiredFields(ctx, models.ThingWithRequiredFields{
			Name: db.String("string2"),
		}))
		require.Nil(t, d.SaveThingWithRequiredFields(ctx, models.ThingWithRequiredFields{
			Name: db.String("string3"),
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.ThingWithRequiredFields{
				models.ThingWithRequiredFields{
					Name: db.String("string1"),
				},
				models.ThingWithRequiredFields{
					Name: db.String("string2"),
				},
				models.ThingWithRequiredFields{
					Name: db.String("string3"),
				},
			}
			actual := []models.ThingWithRequiredFields{}
			err := d.ScanThingWithRequiredFieldss(ctx, db.ScanThingWithRequiredFieldssInput{}, func(m *models.ThingWithRequiredFields, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		// FAILING_TEST
		t.Run("starting after", func(t *testing.T) {
			t.Skip()
			// Scan for everything.
			allItems := []models.ThingWithRequiredFields{}
			err := d.ScanThingWithRequiredFieldss(ctx, db.ScanThingWithRequiredFieldssInput{}, func(m *models.ThingWithRequiredFields, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput := db.ScanThingWithRequiredFieldssInput{
				StartingAfter: &models.ThingWithRequiredFields{
					Name: firstItem.Name,
				},
			}
			actual := []models.ThingWithRequiredFields{}
			err = d.ScanThingWithRequiredFieldss(ctx, scanInput, func(m *models.ThingWithRequiredFields, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanThingWithRequiredFieldssInput{
				Limit: &limit,
			}
			actual := []models.ThingWithRequiredFields{}
			err := d.ScanThingWithRequiredFieldss(ctx, scanInput, func(m *models.ThingWithRequiredFields, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

func SaveThingWithRequiredFields(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithRequiredFields{
			Name: db.String("string1"),
		}
		require.Nil(t, s.SaveThingWithRequiredFields(ctx, m))
		require.IsType(t, db.ErrThingWithRequiredFieldsAlreadyExists{}, s.SaveThingWithRequiredFields(ctx, m))
	}
}

func DeleteThingWithRequiredFields(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithRequiredFields{
			Name: db.String("string1"),
		}
		require.Nil(t, s.SaveThingWithRequiredFields(ctx, m))
		require.Nil(t, s.DeleteThingWithRequiredFields(ctx, *m.Name))
	}
}

func GetThingWithRequiredFields2(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithRequiredFields2{
			ID:   db.String("string1"),
			Name: db.String("string1"),
		}
		require.Nil(t, s.SaveThingWithRequiredFields2(ctx, m))
		m2, err := s.GetThingWithRequiredFields2(ctx, *m.Name, *m.ID)
		require.Nil(t, err)
		require.Equal(t, *m.Name, *m2.Name)
		require.Equal(t, *m.ID, *m2.ID)

		_, err = s.GetThingWithRequiredFields2(ctx, "string2", "string2")
		require.NotNil(t, err)
		require.IsType(t, err, db.ErrThingWithRequiredFields2NotFound{})
	}
}

type getThingWithRequiredFields2sByNameAndIDInput struct {
	ctx   context.Context
	input db.GetThingWithRequiredFields2sByNameAndIDInput
}
type getThingWithRequiredFields2sByNameAndIDOutput struct {
	thingWithRequiredFields2s []models.ThingWithRequiredFields2
	err                       error
}
type getThingWithRequiredFields2sByNameAndIDTest struct {
	testName string
	d        db.Interface
	input    getThingWithRequiredFields2sByNameAndIDInput
	output   getThingWithRequiredFields2sByNameAndIDOutput
}

func (g getThingWithRequiredFields2sByNameAndIDTest) run(t *testing.T) {
	thingWithRequiredFields2s := []models.ThingWithRequiredFields2{}
	fn := func(m *models.ThingWithRequiredFields2, lastThingWithRequiredFields2 bool) bool {
		thingWithRequiredFields2s = append(thingWithRequiredFields2s, *m)
		if lastThingWithRequiredFields2 {
			return false
		}
		return true
	}
	err := g.d.GetThingWithRequiredFields2sByNameAndID(g.input.ctx, g.input.input, fn)
	if err != nil {
		fmt.Println(err.Error())
	}
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.thingWithRequiredFields2s, thingWithRequiredFields2s)
}

func GetThingWithRequiredFields2sByNameAndID(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithRequiredFields2(ctx, models.ThingWithRequiredFields2{
			Name: db.String("string1"),
			ID:   db.String("string1"),
		}))
		require.Nil(t, d.SaveThingWithRequiredFields2(ctx, models.ThingWithRequiredFields2{
			Name: db.String("string1"),
			ID:   db.String("string2"),
		}))
		require.Nil(t, d.SaveThingWithRequiredFields2(ctx, models.ThingWithRequiredFields2{
			Name: db.String("string1"),
			ID:   db.String("string3"),
		}))
		limit := int64(3)
		tests := []getThingWithRequiredFields2sByNameAndIDTest{
			{
				testName: "basic",
				d:        d,
				input: getThingWithRequiredFields2sByNameAndIDInput{
					ctx: context.Background(),
					input: db.GetThingWithRequiredFields2sByNameAndIDInput{
						Name:  "string1",
						Limit: &limit,
					},
				},
				output: getThingWithRequiredFields2sByNameAndIDOutput{
					thingWithRequiredFields2s: []models.ThingWithRequiredFields2{
						models.ThingWithRequiredFields2{
							Name: db.String("string1"),
							ID:   db.String("string1"),
						},
						models.ThingWithRequiredFields2{
							Name: db.String("string1"),
							ID:   db.String("string2"),
						},
						models.ThingWithRequiredFields2{
							Name: db.String("string1"),
							ID:   db.String("string3"),
						},
					},
					err: nil,
				},
			},
			{
				testName: "descending",
				d:        d,
				input: getThingWithRequiredFields2sByNameAndIDInput{
					ctx: context.Background(),
					input: db.GetThingWithRequiredFields2sByNameAndIDInput{
						Name:       "string1",
						Descending: true,
					},
				},
				output: getThingWithRequiredFields2sByNameAndIDOutput{
					thingWithRequiredFields2s: []models.ThingWithRequiredFields2{
						models.ThingWithRequiredFields2{
							Name: db.String("string1"),
							ID:   db.String("string3"),
						},
						models.ThingWithRequiredFields2{
							Name: db.String("string1"),
							ID:   db.String("string2"),
						},
						models.ThingWithRequiredFields2{
							Name: db.String("string1"),
							ID:   db.String("string1"),
						},
					},
					err: nil,
				},
			},
			/* FAILING_TEST */
			/* {
			   testName: "starting after",
			   d:    d,
			   input: getThingWithRequiredFields2sByNameAndIDInput{
			     ctx: context.Background(),
			     input: db.GetThingWithRequiredFields2sByNameAndIDInput{
			         Name: "string1",
			       StartingAfter: &models.ThingWithRequiredFields2{
			           Name:    db.String("string1"),
			           ID:    db.String("string1"),
			       },
			     },
			   },
			   output: getThingWithRequiredFields2sByNameAndIDOutput{
			     thingWithRequiredFields2s: []models.ThingWithRequiredFields2{
			       models.ThingWithRequiredFields2{
			           Name:    db.String("string1"),
			           ID: db.String("string2"),
			       },
			       models.ThingWithRequiredFields2{
			           Name:    db.String("string1"),
			           ID: db.String("string3"),
			       },
			     },
			     err: nil,
			   },
			 }, */
			{
				testName: "starting after descending",
				d:        d,
				input: getThingWithRequiredFields2sByNameAndIDInput{
					ctx: context.Background(),
					input: db.GetThingWithRequiredFields2sByNameAndIDInput{
						Name: "string1",
						StartingAfter: &models.ThingWithRequiredFields2{
							Name: db.String("string1"),
							ID:   db.String("string3"),
						},
						Descending: true,
					},
				},
				output: getThingWithRequiredFields2sByNameAndIDOutput{
					thingWithRequiredFields2s: []models.ThingWithRequiredFields2{
						models.ThingWithRequiredFields2{
							Name: db.String("string1"),
							ID:   db.String("string2"),
						},
						models.ThingWithRequiredFields2{
							Name: db.String("string1"),
							ID:   db.String("string1"),
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting at",
				d:        d,
				input: getThingWithRequiredFields2sByNameAndIDInput{
					ctx: context.Background(),
					input: db.GetThingWithRequiredFields2sByNameAndIDInput{
						Name:         "string1",
						IDStartingAt: db.String("string2"),
					},
				},
				output: getThingWithRequiredFields2sByNameAndIDOutput{
					thingWithRequiredFields2s: []models.ThingWithRequiredFields2{
						models.ThingWithRequiredFields2{
							Name: db.String("string1"),
							ID:   db.String("string2"),
						},
						models.ThingWithRequiredFields2{
							Name: db.String("string1"),
							ID:   db.String("string3"),
						},
					},
					err: nil,
				},
			},
		}
		for _, test := range tests {
			t.Run(test.testName, test.run)
		}
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanThingWithRequiredFields2s(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithRequiredFields2(ctx, models.ThingWithRequiredFields2{
			ID:   db.String("string1"),
			Name: db.String("string1"),
		}))
		require.Nil(t, d.SaveThingWithRequiredFields2(ctx, models.ThingWithRequiredFields2{
			ID:   db.String("string2"),
			Name: db.String("string2"),
		}))
		require.Nil(t, d.SaveThingWithRequiredFields2(ctx, models.ThingWithRequiredFields2{
			ID:   db.String("string3"),
			Name: db.String("string3"),
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.ThingWithRequiredFields2{
				models.ThingWithRequiredFields2{
					ID:   db.String("string1"),
					Name: db.String("string1"),
				},
				models.ThingWithRequiredFields2{
					ID:   db.String("string2"),
					Name: db.String("string2"),
				},
				models.ThingWithRequiredFields2{
					ID:   db.String("string3"),
					Name: db.String("string3"),
				},
			}
			actual := []models.ThingWithRequiredFields2{}
			err := d.ScanThingWithRequiredFields2s(ctx, db.ScanThingWithRequiredFields2sInput{}, func(m *models.ThingWithRequiredFields2, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		// FAILING_TEST
		t.Run("starting after", func(t *testing.T) {
			t.Skip()
			// Scan for everything.
			allItems := []models.ThingWithRequiredFields2{}
			err := d.ScanThingWithRequiredFields2s(ctx, db.ScanThingWithRequiredFields2sInput{}, func(m *models.ThingWithRequiredFields2, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput := db.ScanThingWithRequiredFields2sInput{
				StartingAfter: &models.ThingWithRequiredFields2{
					Name: firstItem.Name,
					ID:   firstItem.ID,
				},
			}
			actual := []models.ThingWithRequiredFields2{}
			err = d.ScanThingWithRequiredFields2s(ctx, scanInput, func(m *models.ThingWithRequiredFields2, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanThingWithRequiredFields2sInput{
				Limit: &limit,
			}
			actual := []models.ThingWithRequiredFields2{}
			err := d.ScanThingWithRequiredFields2s(ctx, scanInput, func(m *models.ThingWithRequiredFields2, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

func SaveThingWithRequiredFields2(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithRequiredFields2{
			ID:   db.String("string1"),
			Name: db.String("string1"),
		}
		require.Nil(t, s.SaveThingWithRequiredFields2(ctx, m))
		require.IsType(t, db.ErrThingWithRequiredFields2AlreadyExists{}, s.SaveThingWithRequiredFields2(ctx, m))
	}
}

func DeleteThingWithRequiredFields2(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithRequiredFields2{
			ID:   db.String("string1"),
			Name: db.String("string1"),
		}
		require.Nil(t, s.SaveThingWithRequiredFields2(ctx, m))
		require.Nil(t, s.DeleteThingWithRequiredFields2(ctx, *m.Name, *m.ID))
	}
}

func GetThingWithTransactMultipleGSI(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithTransactMultipleGSI{
			DateH: mustDate("2018-03-11"),
			DateR: mustDate("2018-03-11"),
			ID:    "string1",
		}
		require.Nil(t, s.SaveThingWithTransactMultipleGSI(ctx, m))
		m2, err := s.GetThingWithTransactMultipleGSI(ctx, m.DateH)
		require.Nil(t, err)
		require.Equal(t, m.DateH, m2.DateH)

		_, err = s.GetThingWithTransactMultipleGSI(ctx, mustDate("2018-03-12"))
		require.NotNil(t, err)
		require.IsType(t, err, db.ErrThingWithTransactMultipleGSINotFound{})
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanThingWithTransactMultipleGSIs(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithTransactMultipleGSI(ctx, models.ThingWithTransactMultipleGSI{
			DateH: mustDate("2018-03-11"),
			DateR: mustDate("2018-03-11"),
			ID:    "string1",
		}))
		require.Nil(t, d.SaveThingWithTransactMultipleGSI(ctx, models.ThingWithTransactMultipleGSI{
			DateH: mustDate("2018-03-12"),
			DateR: mustDate("2018-03-12"),
			ID:    "string2",
		}))
		require.Nil(t, d.SaveThingWithTransactMultipleGSI(ctx, models.ThingWithTransactMultipleGSI{
			DateH: mustDate("2018-03-13"),
			DateR: mustDate("2018-03-13"),
			ID:    "string3",
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.ThingWithTransactMultipleGSI{
				models.ThingWithTransactMultipleGSI{
					DateH: mustDate("2018-03-11"),
					DateR: mustDate("2018-03-11"),
					ID:    "string1",
				},
				models.ThingWithTransactMultipleGSI{
					DateH: mustDate("2018-03-12"),
					DateR: mustDate("2018-03-12"),
					ID:    "string2",
				},
				models.ThingWithTransactMultipleGSI{
					DateH: mustDate("2018-03-13"),
					DateR: mustDate("2018-03-13"),
					ID:    "string3",
				},
			}
			actual := []models.ThingWithTransactMultipleGSI{}
			err := d.ScanThingWithTransactMultipleGSIs(ctx, db.ScanThingWithTransactMultipleGSIsInput{}, func(m *models.ThingWithTransactMultipleGSI, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		// FAILING_TEST
		t.Run("starting after", func(t *testing.T) {
			t.Skip()
			// Scan for everything.
			allItems := []models.ThingWithTransactMultipleGSI{}
			err := d.ScanThingWithTransactMultipleGSIs(ctx, db.ScanThingWithTransactMultipleGSIsInput{}, func(m *models.ThingWithTransactMultipleGSI, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput := db.ScanThingWithTransactMultipleGSIsInput{
				StartingAfter: &models.ThingWithTransactMultipleGSI{
					DateH: firstItem.DateH,
				},
			}
			actual := []models.ThingWithTransactMultipleGSI{}
			err = d.ScanThingWithTransactMultipleGSIs(ctx, scanInput, func(m *models.ThingWithTransactMultipleGSI, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanThingWithTransactMultipleGSIsInput{
				Limit: &limit,
			}
			actual := []models.ThingWithTransactMultipleGSI{}
			err := d.ScanThingWithTransactMultipleGSIs(ctx, scanInput, func(m *models.ThingWithTransactMultipleGSI, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

func SaveThingWithTransactMultipleGSI(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithTransactMultipleGSI{
			DateH: mustDate("2018-03-11"),
			DateR: mustDate("2018-03-11"),
			ID:    "string1",
		}
		require.Nil(t, s.SaveThingWithTransactMultipleGSI(ctx, m))
		require.IsType(t, db.ErrThingWithTransactMultipleGSIAlreadyExists{}, s.SaveThingWithTransactMultipleGSI(ctx, m))
	}
}

func DeleteThingWithTransactMultipleGSI(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithTransactMultipleGSI{
			DateH: mustDate("2018-03-11"),
			DateR: mustDate("2018-03-11"),
			ID:    "string1",
		}
		require.Nil(t, s.SaveThingWithTransactMultipleGSI(ctx, m))
		require.Nil(t, s.DeleteThingWithTransactMultipleGSI(ctx, m.DateH))
	}
}

type getThingWithTransactMultipleGSIsByIDAndDateRInput struct {
	ctx   context.Context
	input db.GetThingWithTransactMultipleGSIsByIDAndDateRInput
}
type getThingWithTransactMultipleGSIsByIDAndDateROutput struct {
	thingWithTransactMultipleGSIs []models.ThingWithTransactMultipleGSI
	err                           error
}
type getThingWithTransactMultipleGSIsByIDAndDateRTest struct {
	testName string
	d        db.Interface
	input    getThingWithTransactMultipleGSIsByIDAndDateRInput
	output   getThingWithTransactMultipleGSIsByIDAndDateROutput
}

func (g getThingWithTransactMultipleGSIsByIDAndDateRTest) run(t *testing.T) {
	thingWithTransactMultipleGSIs := []models.ThingWithTransactMultipleGSI{}
	fn := func(m *models.ThingWithTransactMultipleGSI, lastThingWithTransactMultipleGSI bool) bool {
		thingWithTransactMultipleGSIs = append(thingWithTransactMultipleGSIs, *m)
		if lastThingWithTransactMultipleGSI {
			return false
		}
		return true
	}
	err := g.d.GetThingWithTransactMultipleGSIsByIDAndDateR(g.input.ctx, g.input.input, fn)
	if err != nil {
		fmt.Println(err.Error())
	}
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.thingWithTransactMultipleGSIs, thingWithTransactMultipleGSIs)
}

func GetThingWithTransactMultipleGSIsByIDAndDateR(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithTransactMultipleGSI(ctx, models.ThingWithTransactMultipleGSI{
			ID:    "string1",
			DateR: mustDate("2018-03-11"),
			DateH: mustDate("2018-03-11"),
		}))
		require.Nil(t, d.SaveThingWithTransactMultipleGSI(ctx, models.ThingWithTransactMultipleGSI{
			ID:    "string1",
			DateR: mustDate("2018-03-12"),
			DateH: mustDate("2018-03-13"),
		}))
		require.Nil(t, d.SaveThingWithTransactMultipleGSI(ctx, models.ThingWithTransactMultipleGSI{
			ID:    "string1",
			DateR: mustDate("2018-03-13"),
			DateH: mustDate("2018-03-12"),
		}))
		limit := int64(3)
		tests := []getThingWithTransactMultipleGSIsByIDAndDateRTest{
			{
				testName: "basic",
				d:        d,
				input: getThingWithTransactMultipleGSIsByIDAndDateRInput{
					ctx: context.Background(),
					input: db.GetThingWithTransactMultipleGSIsByIDAndDateRInput{
						ID:    "string1",
						Limit: &limit,
					},
				},
				output: getThingWithTransactMultipleGSIsByIDAndDateROutput{
					thingWithTransactMultipleGSIs: []models.ThingWithTransactMultipleGSI{
						models.ThingWithTransactMultipleGSI{
							ID:    "string1",
							DateR: mustDate("2018-03-11"),
							DateH: mustDate("2018-03-11"),
						},
						models.ThingWithTransactMultipleGSI{
							ID:    "string1",
							DateR: mustDate("2018-03-12"),
							DateH: mustDate("2018-03-13"),
						},
						models.ThingWithTransactMultipleGSI{
							ID:    "string1",
							DateR: mustDate("2018-03-13"),
							DateH: mustDate("2018-03-12"),
						},
					},
					err: nil,
				},
			},
			{
				testName: "descending",
				d:        d,
				input: getThingWithTransactMultipleGSIsByIDAndDateRInput{
					ctx: context.Background(),
					input: db.GetThingWithTransactMultipleGSIsByIDAndDateRInput{
						ID:         "string1",
						Descending: true,
					},
				},
				output: getThingWithTransactMultipleGSIsByIDAndDateROutput{
					thingWithTransactMultipleGSIs: []models.ThingWithTransactMultipleGSI{
						models.ThingWithTransactMultipleGSI{
							ID:    "string1",
							DateR: mustDate("2018-03-13"),
							DateH: mustDate("2018-03-12"),
						},
						models.ThingWithTransactMultipleGSI{
							ID:    "string1",
							DateR: mustDate("2018-03-12"),
							DateH: mustDate("2018-03-13"),
						},
						models.ThingWithTransactMultipleGSI{
							ID:    "string1",
							DateR: mustDate("2018-03-11"),
							DateH: mustDate("2018-03-11"),
						},
					},
					err: nil,
				},
			},
			/* FAILING_TEST */
			/* {
			   testName: "starting after",
			   d:    d,
			   input: getThingWithTransactMultipleGSIsByIDAndDateRInput{
			     ctx: context.Background(),
			     input: db.GetThingWithTransactMultipleGSIsByIDAndDateRInput{
			         ID: "string1",
			       StartingAfter: &models.ThingWithTransactMultipleGSI{
			         ID:    "string1",
			         DateR: mustDate("2018-03-11"),
			         DateH:    mustDate("2018-03-11"),
			       },
			     },
			   },
			   output: getThingWithTransactMultipleGSIsByIDAndDateROutput{
			     thingWithTransactMultipleGSIs: []models.ThingWithTransactMultipleGSI{
			       models.ThingWithTransactMultipleGSI{
			         ID:    "string1",
			         DateR: mustDate("2018-03-12"),
			         DateH:    mustDate("2018-03-13"),
			       },
			       models.ThingWithTransactMultipleGSI{
			         ID:    "string1",
			         DateR: mustDate("2018-03-13"),
			         DateH:    mustDate("2018-03-12"),
			       },
			     },
			     err: nil,
			   },
			 }, */
			{
				testName: "starting after descending",
				d:        d,
				input: getThingWithTransactMultipleGSIsByIDAndDateRInput{
					ctx: context.Background(),
					input: db.GetThingWithTransactMultipleGSIsByIDAndDateRInput{
						ID: "string1",
						StartingAfter: &models.ThingWithTransactMultipleGSI{
							ID:    "string1",
							DateR: mustDate("2018-03-13"),
							DateH: mustDate("2018-03-12"),
						},
						Descending: true,
					},
				},
				output: getThingWithTransactMultipleGSIsByIDAndDateROutput{
					thingWithTransactMultipleGSIs: []models.ThingWithTransactMultipleGSI{
						models.ThingWithTransactMultipleGSI{
							ID:    "string1",
							DateR: mustDate("2018-03-12"),
							DateH: mustDate("2018-03-13"),
						},
						models.ThingWithTransactMultipleGSI{
							ID:    "string1",
							DateR: mustDate("2018-03-11"),
							DateH: mustDate("2018-03-11"),
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting at",
				d:        d,
				input: getThingWithTransactMultipleGSIsByIDAndDateRInput{
					ctx: context.Background(),
					input: db.GetThingWithTransactMultipleGSIsByIDAndDateRInput{
						ID:              "string1",
						DateRStartingAt: db.Date(mustDate("2018-03-12")),
					},
				},
				output: getThingWithTransactMultipleGSIsByIDAndDateROutput{
					thingWithTransactMultipleGSIs: []models.ThingWithTransactMultipleGSI{
						models.ThingWithTransactMultipleGSI{
							ID:    "string1",
							DateR: mustDate("2018-03-12"),
							DateH: mustDate("2018-03-13"),
						},
						models.ThingWithTransactMultipleGSI{
							ID:    "string1",
							DateR: mustDate("2018-03-13"),
							DateH: mustDate("2018-03-12"),
						},
					},
					err: nil,
				},
			},
		}
		for _, test := range tests {
			t.Run(test.testName, test.run)
		}
	}
}

type getThingWithTransactMultipleGSIsByDateHAndIDInput struct {
	ctx   context.Context
	input db.GetThingWithTransactMultipleGSIsByDateHAndIDInput
}
type getThingWithTransactMultipleGSIsByDateHAndIDOutput struct {
	thingWithTransactMultipleGSIs []models.ThingWithTransactMultipleGSI
	err                           error
}
type getThingWithTransactMultipleGSIsByDateHAndIDTest struct {
	testName string
	d        db.Interface
	input    getThingWithTransactMultipleGSIsByDateHAndIDInput
	output   getThingWithTransactMultipleGSIsByDateHAndIDOutput
}

func (g getThingWithTransactMultipleGSIsByDateHAndIDTest) run(t *testing.T) {
	thingWithTransactMultipleGSIs := []models.ThingWithTransactMultipleGSI{}
	fn := func(m *models.ThingWithTransactMultipleGSI, lastThingWithTransactMultipleGSI bool) bool {
		thingWithTransactMultipleGSIs = append(thingWithTransactMultipleGSIs, *m)
		if lastThingWithTransactMultipleGSI {
			return false
		}
		return true
	}
	err := g.d.GetThingWithTransactMultipleGSIsByDateHAndID(g.input.ctx, g.input.input, fn)
	if err != nil {
		fmt.Println(err.Error())
	}
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.thingWithTransactMultipleGSIs, thingWithTransactMultipleGSIs)
}

func GetThingWithTransactMultipleGSIsByDateHAndID(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithTransactMultipleGSI(ctx, models.ThingWithTransactMultipleGSI{
			DateH: mustDate("2018-03-11"),
			ID:    "string1",
		}))
		require.Nil(t, d.SaveThingWithTransactMultipleGSI(ctx, models.ThingWithTransactMultipleGSI{
			DateH: mustDate("2018-03-11"),
			ID:    "string2",
		}))
		require.Nil(t, d.SaveThingWithTransactMultipleGSI(ctx, models.ThingWithTransactMultipleGSI{
			DateH: mustDate("2018-03-11"),
			ID:    "string3",
		}))
		limit := int64(3)
		tests := []getThingWithTransactMultipleGSIsByDateHAndIDTest{
			{
				testName: "basic",
				d:        d,
				input: getThingWithTransactMultipleGSIsByDateHAndIDInput{
					ctx: context.Background(),
					input: db.GetThingWithTransactMultipleGSIsByDateHAndIDInput{
						DateH: mustDate("2018-03-11"),
						Limit: &limit,
					},
				},
				output: getThingWithTransactMultipleGSIsByDateHAndIDOutput{
					thingWithTransactMultipleGSIs: []models.ThingWithTransactMultipleGSI{
						models.ThingWithTransactMultipleGSI{
							DateH: mustDate("2018-03-11"),
							ID:    "string1",
						},
						models.ThingWithTransactMultipleGSI{
							DateH: mustDate("2018-03-11"),
							ID:    "string2",
						},
						models.ThingWithTransactMultipleGSI{
							DateH: mustDate("2018-03-11"),
							ID:    "string3",
						},
					},
					err: nil,
				},
			},
			{
				testName: "descending",
				d:        d,
				input: getThingWithTransactMultipleGSIsByDateHAndIDInput{
					ctx: context.Background(),
					input: db.GetThingWithTransactMultipleGSIsByDateHAndIDInput{
						DateH:      mustDate("2018-03-11"),
						Descending: true,
					},
				},
				output: getThingWithTransactMultipleGSIsByDateHAndIDOutput{
					thingWithTransactMultipleGSIs: []models.ThingWithTransactMultipleGSI{
						models.ThingWithTransactMultipleGSI{
							DateH: mustDate("2018-03-11"),
							ID:    "string3",
						},
						models.ThingWithTransactMultipleGSI{
							DateH: mustDate("2018-03-11"),
							ID:    "string2",
						},
						models.ThingWithTransactMultipleGSI{
							DateH: mustDate("2018-03-11"),
							ID:    "string1",
						},
					},
					err: nil,
				},
			},
			/* FAILING_TEST */
			/* {
			   testName: "starting after",
			   d:    d,
			   input: getThingWithTransactMultipleGSIsByDateHAndIDInput{
			     ctx: context.Background(),
			     input: db.GetThingWithTransactMultipleGSIsByDateHAndIDInput{
			         DateH: mustDate("2018-03-11"),
			       StartingAfter: &models.ThingWithTransactMultipleGSI{
			         DateH:    mustDate("2018-03-11"),
			         ID: "string1",
			       },
			     },
			   },
			   output: getThingWithTransactMultipleGSIsByDateHAndIDOutput{
			     thingWithTransactMultipleGSIs: []models.ThingWithTransactMultipleGSI{
			       models.ThingWithTransactMultipleGSI{
			         DateH:    mustDate("2018-03-11"),
			         ID: "string2",
			       },
			       models.ThingWithTransactMultipleGSI{
			         DateH:    mustDate("2018-03-11"),
			         ID: "string3",
			       },
			     },
			     err: nil,
			   },
			 }, */
			{
				testName: "starting after descending",
				d:        d,
				input: getThingWithTransactMultipleGSIsByDateHAndIDInput{
					ctx: context.Background(),
					input: db.GetThingWithTransactMultipleGSIsByDateHAndIDInput{
						DateH: mustDate("2018-03-11"),
						StartingAfter: &models.ThingWithTransactMultipleGSI{
							DateH: mustDate("2018-03-11"),
							ID:    "string3",
						},
						Descending: true,
					},
				},
				output: getThingWithTransactMultipleGSIsByDateHAndIDOutput{
					thingWithTransactMultipleGSIs: []models.ThingWithTransactMultipleGSI{
						models.ThingWithTransactMultipleGSI{
							DateH: mustDate("2018-03-11"),
							ID:    "string2",
						},
						models.ThingWithTransactMultipleGSI{
							DateH: mustDate("2018-03-11"),
							ID:    "string1",
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting at",
				d:        d,
				input: getThingWithTransactMultipleGSIsByDateHAndIDInput{
					ctx: context.Background(),
					input: db.GetThingWithTransactMultipleGSIsByDateHAndIDInput{
						DateH:        mustDate("2018-03-11"),
						IDStartingAt: db.String("string2"),
					},
				},
				output: getThingWithTransactMultipleGSIsByDateHAndIDOutput{
					thingWithTransactMultipleGSIs: []models.ThingWithTransactMultipleGSI{
						models.ThingWithTransactMultipleGSI{
							DateH: mustDate("2018-03-11"),
							ID:    "string2",
						},
						models.ThingWithTransactMultipleGSI{
							DateH: mustDate("2018-03-11"),
							ID:    "string3",
						},
					},
					err: nil,
				},
			},
		}
		for _, test := range tests {
			t.Run(test.testName, test.run)
		}
	}
}

func TransactSaveThingWithTransactMultipleGSIAndThing(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m1 := models.ThingWithTransactMultipleGSI{
			DateH: mustDate("2018-03-11"),
			DateR: mustDate("2018-03-11"),
			ID:    "string1",
		}
		m2 := models.Thing{
			CreatedAt:     mustTime("2018-03-11T15:04:01+07:00"),
			HashNullable:  db.String("string1"),
			ID:            "string1",
			Name:          "string1",
			RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			Version:       1,
		}
		require.Nil(t, s.TransactSaveThingWithTransactMultipleGSIAndThing(ctx, m1, nil, m2, nil))
	}
}

func GetThingWithTransaction(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithTransaction{
			Name: "string1",
		}
		require.Nil(t, s.SaveThingWithTransaction(ctx, m))
		m2, err := s.GetThingWithTransaction(ctx, m.Name)
		require.Nil(t, err)
		require.Equal(t, m.Name, m2.Name)

		_, err = s.GetThingWithTransaction(ctx, "string2")
		require.NotNil(t, err)
		require.IsType(t, err, db.ErrThingWithTransactionNotFound{})
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanThingWithTransactions(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithTransaction(ctx, models.ThingWithTransaction{
			Name: "string1",
		}))
		require.Nil(t, d.SaveThingWithTransaction(ctx, models.ThingWithTransaction{
			Name: "string2",
		}))
		require.Nil(t, d.SaveThingWithTransaction(ctx, models.ThingWithTransaction{
			Name: "string3",
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.ThingWithTransaction{
				models.ThingWithTransaction{
					Name: "string1",
				},
				models.ThingWithTransaction{
					Name: "string2",
				},
				models.ThingWithTransaction{
					Name: "string3",
				},
			}
			actual := []models.ThingWithTransaction{}
			err := d.ScanThingWithTransactions(ctx, db.ScanThingWithTransactionsInput{}, func(m *models.ThingWithTransaction, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		// FAILING_TEST
		t.Run("starting after", func(t *testing.T) {
			t.Skip()
			// Scan for everything.
			allItems := []models.ThingWithTransaction{}
			err := d.ScanThingWithTransactions(ctx, db.ScanThingWithTransactionsInput{}, func(m *models.ThingWithTransaction, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput := db.ScanThingWithTransactionsInput{
				StartingAfter: &models.ThingWithTransaction{
					Name: firstItem.Name,
				},
			}
			actual := []models.ThingWithTransaction{}
			err = d.ScanThingWithTransactions(ctx, scanInput, func(m *models.ThingWithTransaction, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanThingWithTransactionsInput{
				Limit: &limit,
			}
			actual := []models.ThingWithTransaction{}
			err := d.ScanThingWithTransactions(ctx, scanInput, func(m *models.ThingWithTransaction, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

func SaveThingWithTransaction(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithTransaction{
			Name: "string1",
		}
		require.Nil(t, s.SaveThingWithTransaction(ctx, m))
		require.IsType(t, db.ErrThingWithTransactionAlreadyExists{}, s.SaveThingWithTransaction(ctx, m))
	}
}

func DeleteThingWithTransaction(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithTransaction{
			Name: "string1",
		}
		require.Nil(t, s.SaveThingWithTransaction(ctx, m))
		require.Nil(t, s.DeleteThingWithTransaction(ctx, m.Name))
	}
}

func TransactSaveThingWithTransactionAndThing(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m1 := models.ThingWithTransaction{
			Name: "string1",
		}
		m2 := models.Thing{
			CreatedAt:     mustTime("2018-03-11T15:04:01+07:00"),
			HashNullable:  db.String("string1"),
			ID:            "string1",
			Name:          "string1",
			RangeNullable: db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
			Version:       1,
		}
		require.Nil(t, s.TransactSaveThingWithTransactionAndThing(ctx, m1, nil, m2, nil))
	}
}

func GetThingWithTransactionWithSimpleThing(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithTransactionWithSimpleThing{
			Name: "string1",
		}
		require.Nil(t, s.SaveThingWithTransactionWithSimpleThing(ctx, m))
		m2, err := s.GetThingWithTransactionWithSimpleThing(ctx, m.Name)
		require.Nil(t, err)
		require.Equal(t, m.Name, m2.Name)

		_, err = s.GetThingWithTransactionWithSimpleThing(ctx, "string2")
		require.NotNil(t, err)
		require.IsType(t, err, db.ErrThingWithTransactionWithSimpleThingNotFound{})
	}
}

// The scan tests are structured differently compared to other tests in because items returned by scans
// are not returned in any particular order, so we can't simply declare what the expected arrays of items are.
func ScanThingWithTransactionWithSimpleThings(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithTransactionWithSimpleThing(ctx, models.ThingWithTransactionWithSimpleThing{
			Name: "string1",
		}))
		require.Nil(t, d.SaveThingWithTransactionWithSimpleThing(ctx, models.ThingWithTransactionWithSimpleThing{
			Name: "string2",
		}))
		require.Nil(t, d.SaveThingWithTransactionWithSimpleThing(ctx, models.ThingWithTransactionWithSimpleThing{
			Name: "string3",
		}))

		t.Run("basic", func(t *testing.T) {
			expected := []models.ThingWithTransactionWithSimpleThing{
				models.ThingWithTransactionWithSimpleThing{
					Name: "string1",
				},
				models.ThingWithTransactionWithSimpleThing{
					Name: "string2",
				},
				models.ThingWithTransactionWithSimpleThing{
					Name: "string3",
				},
			}
			actual := []models.ThingWithTransactionWithSimpleThing{}
			err := d.ScanThingWithTransactionWithSimpleThings(ctx, db.ScanThingWithTransactionWithSimpleThingsInput{}, func(m *models.ThingWithTransactionWithSimpleThing, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)
			// We can't use Equal here because Scan doesn't return items in any specific order.
			require.ElementsMatch(t, expected, actual)
		})

		// FAILING_TEST
		t.Run("starting after", func(t *testing.T) {
			t.Skip()
			// Scan for everything.
			allItems := []models.ThingWithTransactionWithSimpleThing{}
			err := d.ScanThingWithTransactionWithSimpleThings(ctx, db.ScanThingWithTransactionWithSimpleThingsInput{}, func(m *models.ThingWithTransactionWithSimpleThing, last bool) bool {
				allItems = append(allItems, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			firstItem := allItems[0]

			// Scan for everything after the first item.
			scanInput := db.ScanThingWithTransactionWithSimpleThingsInput{
				StartingAfter: &models.ThingWithTransactionWithSimpleThing{
					Name: firstItem.Name,
				},
			}
			actual := []models.ThingWithTransactionWithSimpleThing{}
			err = d.ScanThingWithTransactionWithSimpleThings(ctx, scanInput, func(m *models.ThingWithTransactionWithSimpleThing, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			expected := allItems[1:]
			require.Equal(t, expected, actual)
		})

		t.Run("limit", func(t *testing.T) {
			limit := int64(1)
			// Scan for just the first item.
			scanInput := db.ScanThingWithTransactionWithSimpleThingsInput{
				Limit: &limit,
			}
			actual := []models.ThingWithTransactionWithSimpleThing{}
			err := d.ScanThingWithTransactionWithSimpleThings(ctx, scanInput, func(m *models.ThingWithTransactionWithSimpleThing, last bool) bool {
				actual = append(actual, *m)
				return true
			})
			var errStr string
			if err != nil {
				errStr = err.Error()
			}
			require.NoError(t, err, errStr)

			require.Len(t, actual, 1)
		})
	}
}

func SaveThingWithTransactionWithSimpleThing(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithTransactionWithSimpleThing{
			Name: "string1",
		}
		require.Nil(t, s.SaveThingWithTransactionWithSimpleThing(ctx, m))
		require.IsType(t, db.ErrThingWithTransactionWithSimpleThingAlreadyExists{}, s.SaveThingWithTransactionWithSimpleThing(ctx, m))
	}
}

func DeleteThingWithTransactionWithSimpleThing(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithTransactionWithSimpleThing{
			Name: "string1",
		}
		require.Nil(t, s.SaveThingWithTransactionWithSimpleThing(ctx, m))
		require.Nil(t, s.DeleteThingWithTransactionWithSimpleThing(ctx, m.Name))
	}
}

func TransactSaveThingWithTransactionWithSimpleThingAndSimpleThing(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m1 := models.ThingWithTransactionWithSimpleThing{
			Name: "string1",
		}
		m2 := models.SimpleThing{
			Name: "string1",
		}
		require.Nil(t, s.TransactSaveThingWithTransactionWithSimpleThingAndSimpleThing(ctx, m1, nil, m2, nil))
	}
}

func GetThingWithUnderscores(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithUnderscores{
			IDApp: "string1",
		}
		require.Nil(t, s.SaveThingWithUnderscores(ctx, m))
		m2, err := s.GetThingWithUnderscores(ctx, m.IDApp)
		require.Nil(t, err)
		require.Equal(t, m.IDApp, m2.IDApp)

		_, err = s.GetThingWithUnderscores(ctx, "string2")
		require.NotNil(t, err)
		require.IsType(t, err, db.ErrThingWithUnderscoresNotFound{})
	}
}

func SaveThingWithUnderscores(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithUnderscores{
			IDApp: "string1",
		}
		require.Nil(t, s.SaveThingWithUnderscores(ctx, m))
	}
}

func DeleteThingWithUnderscores(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithUnderscores{
			IDApp: "string1",
		}
		require.Nil(t, s.SaveThingWithUnderscores(ctx, m))
		require.Nil(t, s.DeleteThingWithUnderscores(ctx, m.IDApp))
	}
}
