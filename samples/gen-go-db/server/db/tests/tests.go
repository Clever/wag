package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Clever/wag/samples/gen-go-db/models"
	"github.com/Clever/wag/samples/gen-go-db/server/db"
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

func RunDBTests(t *testing.T, dbFactory func() db.Interface) {
	t.Run("GetDeployment", GetDeployment(dbFactory(), t))
	t.Run("GetDeploymentsByEnvAppAndVersion", GetDeploymentsByEnvAppAndVersion(dbFactory(), t))
	t.Run("SaveDeployment", SaveDeployment(dbFactory(), t))
	t.Run("DeleteDeployment", DeleteDeployment(dbFactory(), t))
	t.Run("GetDeploymentsByEnvAppAndDate", GetDeploymentsByEnvAppAndDate(dbFactory(), t))
	t.Run("GetDeploymentsByEnvironmentAndDate", GetDeploymentsByEnvironmentAndDate(dbFactory(), t))
	t.Run("GetDeploymentByVersion", GetDeploymentByVersion(dbFactory(), t))
	t.Run("GetEvent", GetEvent(dbFactory(), t))
	t.Run("GetEventsByPkAndSk", GetEventsByPkAndSk(dbFactory(), t))
	t.Run("SaveEvent", SaveEvent(dbFactory(), t))
	t.Run("DeleteEvent", DeleteEvent(dbFactory(), t))
	t.Run("GetEventsBySkAndData", GetEventsBySkAndData(dbFactory(), t))
	t.Run("GetNoRangeThingWithCompositeAttributes", GetNoRangeThingWithCompositeAttributes(dbFactory(), t))
	t.Run("SaveNoRangeThingWithCompositeAttributes", SaveNoRangeThingWithCompositeAttributes(dbFactory(), t))
	t.Run("DeleteNoRangeThingWithCompositeAttributes", DeleteNoRangeThingWithCompositeAttributes(dbFactory(), t))
	t.Run("GetNoRangeThingWithCompositeAttributessByNameVersionAndDate", GetNoRangeThingWithCompositeAttributessByNameVersionAndDate(dbFactory(), t))
	t.Run("GetSimpleThing", GetSimpleThing(dbFactory(), t))
	t.Run("SaveSimpleThing", SaveSimpleThing(dbFactory(), t))
	t.Run("DeleteSimpleThing", DeleteSimpleThing(dbFactory(), t))
	t.Run("GetTeacherSharingRule", GetTeacherSharingRule(dbFactory(), t))
	t.Run("GetTeacherSharingRulesByTeacherAndSchoolApp", GetTeacherSharingRulesByTeacherAndSchoolApp(dbFactory(), t))
	t.Run("SaveTeacherSharingRule", SaveTeacherSharingRule(dbFactory(), t))
	t.Run("DeleteTeacherSharingRule", DeleteTeacherSharingRule(dbFactory(), t))
	t.Run("GetTeacherSharingRulesByDistrictAndSchoolTeacherApp", GetTeacherSharingRulesByDistrictAndSchoolTeacherApp(dbFactory(), t))
	t.Run("GetThing", GetThing(dbFactory(), t))
	t.Run("ScanThings", ScanThings(dbFactory(), t))
	t.Run("GetThingsByNameAndVersion", GetThingsByNameAndVersion(dbFactory(), t))
	t.Run("SaveThing", SaveThing(dbFactory(), t))
	t.Run("DeleteThing", DeleteThing(dbFactory(), t))
	t.Run("GetThingByID", GetThingByID(dbFactory(), t))
	t.Run("GetThingsByNameAndCreatedAt", GetThingsByNameAndCreatedAt(dbFactory(), t))
	t.Run("GetThingWithCompositeAttributes", GetThingWithCompositeAttributes(dbFactory(), t))
	t.Run("GetThingWithCompositeAttributessByNameBranchAndDate", GetThingWithCompositeAttributessByNameBranchAndDate(dbFactory(), t))
	t.Run("SaveThingWithCompositeAttributes", SaveThingWithCompositeAttributes(dbFactory(), t))
	t.Run("DeleteThingWithCompositeAttributes", DeleteThingWithCompositeAttributes(dbFactory(), t))
	t.Run("GetThingWithCompositeAttributessByNameVersionAndDate", GetThingWithCompositeAttributessByNameVersionAndDate(dbFactory(), t))
	t.Run("GetThingWithCompositeEnumAttributes", GetThingWithCompositeEnumAttributes(dbFactory(), t))
	t.Run("GetThingWithCompositeEnumAttributessByNameBranchAndDate", GetThingWithCompositeEnumAttributessByNameBranchAndDate(dbFactory(), t))
	t.Run("SaveThingWithCompositeEnumAttributes", SaveThingWithCompositeEnumAttributes(dbFactory(), t))
	t.Run("DeleteThingWithCompositeEnumAttributes", DeleteThingWithCompositeEnumAttributes(dbFactory(), t))
	t.Run("GetThingWithDateRange", GetThingWithDateRange(dbFactory(), t))
	t.Run("GetThingWithDateRangesByNameAndDate", GetThingWithDateRangesByNameAndDate(dbFactory(), t))
	t.Run("SaveThingWithDateRange", SaveThingWithDateRange(dbFactory(), t))
	t.Run("DeleteThingWithDateRange", DeleteThingWithDateRange(dbFactory(), t))
	t.Run("GetThingWithDateTimeComposite", GetThingWithDateTimeComposite(dbFactory(), t))
	t.Run("GetThingWithDateTimeCompositesByTypeIDAndCreatedResource", GetThingWithDateTimeCompositesByTypeIDAndCreatedResource(dbFactory(), t))
	t.Run("SaveThingWithDateTimeComposite", SaveThingWithDateTimeComposite(dbFactory(), t))
	t.Run("DeleteThingWithDateTimeComposite", DeleteThingWithDateTimeComposite(dbFactory(), t))
	t.Run("GetThingWithEnumHashKey", GetThingWithEnumHashKey(dbFactory(), t))
	t.Run("GetThingWithEnumHashKeysByBranchAndDate", GetThingWithEnumHashKeysByBranchAndDate(dbFactory(), t))
	t.Run("SaveThingWithEnumHashKey", SaveThingWithEnumHashKey(dbFactory(), t))
	t.Run("DeleteThingWithEnumHashKey", DeleteThingWithEnumHashKey(dbFactory(), t))
	t.Run("GetThingWithEnumHashKeysByBranchAndDate2", GetThingWithEnumHashKeysByBranchAndDate2(dbFactory(), t))
	t.Run("GetThingWithMatchingKeys", GetThingWithMatchingKeys(dbFactory(), t))
	t.Run("GetThingWithMatchingKeyssByBearAndAssocTypeID", GetThingWithMatchingKeyssByBearAndAssocTypeID(dbFactory(), t))
	t.Run("SaveThingWithMatchingKeys", SaveThingWithMatchingKeys(dbFactory(), t))
	t.Run("DeleteThingWithMatchingKeys", DeleteThingWithMatchingKeys(dbFactory(), t))
	t.Run("GetThingWithMatchingKeyssByAssocTypeIDAndCreatedBear", GetThingWithMatchingKeyssByAssocTypeIDAndCreatedBear(dbFactory(), t))
	t.Run("GetThingWithRequiredCompositePropertiesAndKeysOnly", GetThingWithRequiredCompositePropertiesAndKeysOnly(dbFactory(), t))
	t.Run("SaveThingWithRequiredCompositePropertiesAndKeysOnly", SaveThingWithRequiredCompositePropertiesAndKeysOnly(dbFactory(), t))
	t.Run("DeleteThingWithRequiredCompositePropertiesAndKeysOnly", DeleteThingWithRequiredCompositePropertiesAndKeysOnly(dbFactory(), t))
	t.Run("GetThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThree", GetThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThree(dbFactory(), t))
	t.Run("GetThingWithRequiredFields", GetThingWithRequiredFields(dbFactory(), t))
	t.Run("SaveThingWithRequiredFields", SaveThingWithRequiredFields(dbFactory(), t))
	t.Run("DeleteThingWithRequiredFields", DeleteThingWithRequiredFields(dbFactory(), t))
	t.Run("GetThingWithRequiredFields2", GetThingWithRequiredFields2(dbFactory(), t))
	t.Run("GetThingWithRequiredFields2sByNameAndID", GetThingWithRequiredFields2sByNameAndID(dbFactory(), t))
	t.Run("SaveThingWithRequiredFields2", SaveThingWithRequiredFields2(dbFactory(), t))
	t.Run("DeleteThingWithRequiredFields2", DeleteThingWithRequiredFields2(dbFactory(), t))
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
			{
				testName: "starting after",
				d:        d,
				input: getDeploymentsByEnvAppAndVersionInput{
					ctx: context.Background(),
					input: db.GetDeploymentsByEnvAppAndVersionInput{
						Environment: "string1",
						Application: "string1",
						StartingAfter: &models.Deployment{
							Environment: "string1",
							Application: "string1",
							Version:     "string1",
						},
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
			{
				testName: "starting after",
				d:        d,
				input: getDeploymentsByEnvAppAndDateInput{
					ctx: context.Background(),
					input: db.GetDeploymentsByEnvAppAndDateInput{
						Environment: "string1",
						Application: "string1",
						StartingAfter: &models.Deployment{
							Environment: "string1",
							Application: "string1",
							Date:        mustTime("2018-03-11T15:04:01+07:00"),
							Version:     "string1",
						},
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
			{
				testName: "starting after",
				d:        d,
				input: getDeploymentsByEnvironmentAndDateInput{
					ctx: context.Background(),
					input: db.GetDeploymentsByEnvironmentAndDateInput{
						Environment: "string1",
						StartingAfter: &models.Deployment{
							Environment: "string1",
							Date:        mustTime("2018-03-11T15:04:01+07:00"),
							Application: "string1",
							Version:     "string1",
						},
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
			{
				testName: "starting after",
				d:        d,
				input: getEventsByPkAndSkInput{
					ctx: context.Background(),
					input: db.GetEventsByPkAndSkInput{
						Pk: "string1",
						StartingAfter: &models.Event{
							Pk: "string1",
							Sk: "string1",
						},
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
			{
				testName: "starting after",
				d:        d,
				input: getEventsBySkAndDataInput{
					ctx: context.Background(),
					input: db.GetEventsBySkAndDataInput{
						Sk: "string1",
						StartingAfter: &models.Event{
							Sk:   "string1",
							Data: []byte("string1"),
							Pk:   "string1",
						},
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

func GetNoRangeThingWithCompositeAttributes(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.NoRangeThingWithCompositeAttributes{
			Branch:  db.String("string1"),
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

func SaveNoRangeThingWithCompositeAttributes(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.NoRangeThingWithCompositeAttributes{
			Branch:  db.String("string1"),
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
		}))
		require.Nil(t, d.SaveNoRangeThingWithCompositeAttributes(ctx, models.NoRangeThingWithCompositeAttributes{
			Name:    db.String("string1"),
			Version: 1,
			Date:    db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
			Branch:  db.String("string3"),
		}))
		require.Nil(t, d.SaveNoRangeThingWithCompositeAttributes(ctx, models.NoRangeThingWithCompositeAttributes{
			Name:    db.String("string1"),
			Version: 1,
			Date:    db.DateTime(mustTime("2018-03-11T15:04:03+07:00")),
			Branch:  db.String("string2"),
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
						},
						models.NoRangeThingWithCompositeAttributes{
							Name:    db.String("string1"),
							Version: 1,
							Date:    db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
							Branch:  db.String("string3"),
						},
						models.NoRangeThingWithCompositeAttributes{
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
						},
						models.NoRangeThingWithCompositeAttributes{
							Name:    db.String("string1"),
							Version: 1,
							Date:    db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
							Branch:  db.String("string3"),
						},
						models.NoRangeThingWithCompositeAttributes{
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
				testName: "starting after",
				d:        d,
				input: getNoRangeThingWithCompositeAttributessByNameVersionAndDateInput{
					ctx: context.Background(),
					input: db.GetNoRangeThingWithCompositeAttributessByNameVersionAndDateInput{
						Name:    "string1",
						Version: 1,
						StartingAfter: &models.NoRangeThingWithCompositeAttributes{
							Name:    db.String("string1"),
							Version: 1,
							Date:    db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
							Branch:  db.String("string1"),
						},
					},
				},
				output: getNoRangeThingWithCompositeAttributessByNameVersionAndDateOutput{
					noRangeThingWithCompositeAttributess: []models.NoRangeThingWithCompositeAttributes{
						models.NoRangeThingWithCompositeAttributes{
							Name:    db.String("string1"),
							Version: 1,
							Date:    db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
							Branch:  db.String("string3"),
						},
						models.NoRangeThingWithCompositeAttributes{
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
						},
						models.NoRangeThingWithCompositeAttributes{
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
						},
						models.NoRangeThingWithCompositeAttributes{
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
			District: "district",
		}))
		require.Nil(t, d.SaveTeacherSharingRule(ctx, models.TeacherSharingRule{
			Teacher:  "string1",
			School:   "string2",
			App:      "string2",
			District: "district",
		}))
		require.Nil(t, d.SaveTeacherSharingRule(ctx, models.TeacherSharingRule{
			Teacher:  "string1",
			School:   "string3",
			App:      "string3",
			District: "district",
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
							District: "district",
						},
						models.TeacherSharingRule{
							Teacher:  "string1",
							School:   "string2",
							App:      "string2",
							District: "district",
						},
						models.TeacherSharingRule{
							Teacher:  "string1",
							School:   "string3",
							App:      "string3",
							District: "district",
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
							District: "district",
						},
						models.TeacherSharingRule{
							Teacher:  "string1",
							School:   "string2",
							App:      "string2",
							District: "district",
						},
						models.TeacherSharingRule{
							Teacher:  "string1",
							School:   "string1",
							App:      "string1",
							District: "district",
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting after",
				d:        d,
				input: getTeacherSharingRulesByTeacherAndSchoolAppInput{
					ctx: context.Background(),
					input: db.GetTeacherSharingRulesByTeacherAndSchoolAppInput{
						Teacher: "string1",
						StartingAfter: &models.TeacherSharingRule{
							Teacher:  "string1",
							School:   "string1",
							App:      "string1",
							District: "district",
						},
					},
				},
				output: getTeacherSharingRulesByTeacherAndSchoolAppOutput{
					teacherSharingRules: []models.TeacherSharingRule{
						models.TeacherSharingRule{
							Teacher:  "string1",
							School:   "string2",
							App:      "string2",
							District: "district",
						},
						models.TeacherSharingRule{
							Teacher:  "string1",
							School:   "string3",
							App:      "string3",
							District: "district",
						},
					},
					err: nil,
				},
			},
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
							District: "district",
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
							District: "district",
						},
						models.TeacherSharingRule{
							Teacher:  "string1",
							School:   "string1",
							App:      "string1",
							District: "district",
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
							District: "district",
						},
						models.TeacherSharingRule{
							Teacher:  "string1",
							School:   "string3",
							App:      "string3",
							District: "district",
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
			{
				testName: "starting after",
				d:        d,
				input: getTeacherSharingRulesByDistrictAndSchoolTeacherAppInput{
					ctx: context.Background(),
					input: db.GetTeacherSharingRulesByDistrictAndSchoolTeacherAppInput{
						District: "string1",
						StartingAfter: &models.TeacherSharingRule{
							District: "string1",
							School:   "string1",
							Teacher:  "string1",
							App:      "string1",
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

func GetThing(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.Thing{
			CreatedAt: mustTime("2018-03-11T15:04:01+07:00"),
			ID:        "string1",
			Name:      "string1",
			Version:   1,
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
			{
				testName: "starting after",
				d:        d,
				input: getThingsByNameAndVersionInput{
					ctx: context.Background(),
					input: db.GetThingsByNameAndVersionInput{
						Name: "string1",
						StartingAfter: &models.Thing{
							Name:    "string1",
							Version: 1,
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
			},
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

type scanThingsInput struct {
	ctx   context.Context
	input db.ScanThingsInput
}
type scanThingsOutput struct {
	things []models.Thing
	err    error
}
type scanThingsTest struct {
	testName string
	d        db.Interface
	input    scanThingsInput
	output   scanThingsOutput
}

func (g scanThingsTest) run(t *testing.T) {
	things := []models.Thing{}
	err := g.d.ScanThings(g.input.ctx, g.input.input, func(m *models.Thing, last bool) bool {
		things = append(things, *m)
		return true
	})
	var errStr string
	if err != nil {
		errStr = err.Error()
	}
	require.Equal(t, g.output.err, err, errStr)
	require.Equal(t, g.output.things, things)
}

func ScanThings(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThing(ctx, models.Thing{
			Name:    "string1",
			Version: 1,
		}))
		require.Nil(t, d.SaveThing(ctx, models.Thing{
			Name:    "string2",
			Version: 2,
		}))
		require.Nil(t, d.SaveThing(ctx, models.Thing{
			Name:    "string3",
			Version: 3,
		}))
		tests := []scanThingsTest{
			{
				testName: "basic",
				d:        d,
				input: scanThingsInput{
					ctx:   context.Background(),
					input: db.ScanThingsInput{},
				},
				output: scanThingsOutput{
					things: []models.Thing{
						models.Thing{
							Name:    "string3",
							Version: 3,
						},
						models.Thing{
							Name:    "string2",
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
				testName: "starting after",
				d:        d,
				input: scanThingsInput{
					ctx: context.Background(),
					input: db.ScanThingsInput{
						StartingAfter: &models.Thing{
							Name:    "string3",
							Version: 3,
						},
					},
				},
				output: scanThingsOutput{
					things: []models.Thing{
						models.Thing{
							Name:    "string2",
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
		}
		for _, test := range tests {
			t.Run(test.testName, test.run)
		}
	}
}

func SaveThing(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.Thing{
			CreatedAt: mustTime("2018-03-11T15:04:01+07:00"),
			ID:        "string1",
			Name:      "string1",
			Version:   1,
		}
		require.Nil(t, s.SaveThing(ctx, m))
		require.IsType(t, db.ErrThingAlreadyExists{}, s.SaveThing(ctx, m))
	}
}

func DeleteThing(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.Thing{
			CreatedAt: mustTime("2018-03-11T15:04:01+07:00"),
			ID:        "string1",
			Name:      "string1",
			Version:   1,
		}
		require.Nil(t, s.SaveThing(ctx, m))
		require.Nil(t, s.DeleteThing(ctx, m.Name, m.Version))
	}
}

func GetThingByID(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.Thing{
			CreatedAt: mustTime("2018-03-11T15:04:01+07:00"),
			ID:        "string1",
			Name:      "string1",
			Version:   1,
		}
		require.Nil(t, s.SaveThing(ctx, m))
		m2, err := s.GetThingByID(ctx, m.ID)
		require.Nil(t, err)
		require.Equal(t, m.CreatedAt.String(), m2.CreatedAt.String())
		require.Equal(t, m.ID, m2.ID)
		require.Equal(t, m.Name, m2.Name)
		require.Equal(t, m.Version, m2.Version)

		_, err = s.GetThingByID(ctx, "string2")
		require.NotNil(t, err)
		require.IsType(t, err, db.ErrThingByIDNotFound{})
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
			{
				testName: "starting after",
				d:        d,
				input: getThingsByNameAndCreatedAtInput{
					ctx: context.Background(),
					input: db.GetThingsByNameAndCreatedAtInput{
						Name: "string1",
						StartingAfter: &models.Thing{
							Name:      "string1",
							CreatedAt: mustTime("2018-03-11T15:04:01+07:00"),
							Version:   1,
						},
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

type scanThingsByNameAndCreatedAtInput struct {
	ctx   context.Context
	input db.ScanThingsByNameAndCreatedAtInput
}
type scanThingsByNameAndCreatedAtOutput struct {
	things []models.Thing
	err    error
}
type scanThingsByNameAndCreatedAtTest struct {
	testName string
	d        db.Interface
	input    scanThingsByNameAndCreatedAtInput
	output   scanThingsByNameAndCreatedAtOutput
}

func (g scanThingsByNameAndCreatedAtTest) run(t *testing.T) {
	things := []models.Thing{}
	err := g.d.ScanThingsByNameAndCreatedAt(g.input.ctx, g.input.input, func(m *models.Thing, last bool) bool {
		things = append(things, *m)
		return true
	})
	var errStr string
	if err != nil {
		errStr = err.Error()
	}
	require.Equal(t, g.output.err, err, errStr)
	require.Equal(t, g.output.things, things)
}

func ScanThingsByNameAndCreatedAt(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThing(ctx, models.Thing{
			Name:    "string1",
			Version: 1,
		}))
		require.Nil(t, d.SaveThing(ctx, models.Thing{
			Name:    "string2",
			Version: 2,
		}))
		require.Nil(t, d.SaveThing(ctx, models.Thing{
			Name:    "string3",
			Version: 3,
		}))
		tests := []scanThingsByNameAndCreatedAtTest{
			{
				testName: "basic",
				d:        d,
				input: scanThingsByNameAndCreatedAtInput{
					ctx:   context.Background(),
					input: db.ScanThingsByNameAndCreatedAtInput{},
				},
				output: scanThingsByNameAndCreatedAtOutput{
					things: []models.Thing{
						models.Thing{
							Name:    "string3",
							Version: 3,
						},
						models.Thing{
							Name:    "string2",
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
				testName: "starting after",
				d:        d,
				input: scanThingsByNameAndCreatedAtInput{
					ctx: context.Background(),
					input: db.ScanThingsByNameAndCreatedAtInput{
						StartingAfter: &models.Thing{
							Name:    "string3",
							Version: 3,
						},
					},
				},
				output: scanThingsByNameAndCreatedAtOutput{
					things: []models.Thing{
						models.Thing{
							Name:    "string2",
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
			{
				testName: "starting after",
				d:        d,
				input: getThingWithCompositeAttributessByNameBranchAndDateInput{
					ctx: context.Background(),
					input: db.GetThingWithCompositeAttributessByNameBranchAndDateInput{
						Name:   "string1",
						Branch: "string1",
						StartingAfter: &models.ThingWithCompositeAttributes{
							Name:   db.String("string1"),
							Branch: db.String("string1"),
							Date:   db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
						},
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
			{
				testName: "starting after",
				d:        d,
				input: getThingWithCompositeAttributessByNameVersionAndDateInput{
					ctx: context.Background(),
					input: db.GetThingWithCompositeAttributessByNameVersionAndDateInput{
						Name:    "string1",
						Version: 1,
						StartingAfter: &models.ThingWithCompositeAttributes{
							Name:    db.String("string1"),
							Version: 1,
							Date:    db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
							Branch:  db.String("string1"),
						},
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
			{
				testName: "starting after",
				d:        d,
				input: getThingWithCompositeEnumAttributessByNameBranchAndDateInput{
					ctx: context.Background(),
					input: db.GetThingWithCompositeEnumAttributessByNameBranchAndDateInput{
						Name:     "string1",
						BranchID: models.BranchMaster,
						StartingAfter: &models.ThingWithCompositeEnumAttributes{
							Name:     db.String("string1"),
							BranchID: models.BranchMaster,
							Date:     db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
						},
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
			{
				testName: "starting after",
				d:        d,
				input: getThingWithDateRangesByNameAndDateInput{
					ctx: context.Background(),
					input: db.GetThingWithDateRangesByNameAndDateInput{
						Name: "string1",
						StartingAfter: &models.ThingWithDateRange{
							Name: "string1",
							Date: mustTime("2018-03-11T15:04:01+07:00"),
						},
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
			{
				testName: "starting after",
				d:        d,
				input: getThingWithDateTimeCompositesByTypeIDAndCreatedResourceInput{
					ctx: context.Background(),
					input: db.GetThingWithDateTimeCompositesByTypeIDAndCreatedResourceInput{
						Type: "string1",
						ID:   "string1",
						StartingAfter: &models.ThingWithDateTimeComposite{
							Type:     "string1",
							ID:       "string1",
							Created:  mustTime("2018-03-11T15:04:01+07:00"),
							Resource: "string1",
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
			{
				testName: "starting after",
				d:        d,
				input: getThingWithEnumHashKeysByBranchAndDateInput{
					ctx: context.Background(),
					input: db.GetThingWithEnumHashKeysByBranchAndDateInput{
						Branch: models.BranchMaster,
						StartingAfter: &models.ThingWithEnumHashKey{
							Branch: models.BranchMaster,
							Date:   mustTime("2018-03-11T15:04:01+07:00"),
						},
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
			{
				testName: "starting after",
				d:        d,
				input: getThingWithEnumHashKeysByBranchAndDate2Input{
					ctx: context.Background(),
					input: db.GetThingWithEnumHashKeysByBranchAndDate2Input{
						Branch: models.BranchMaster,
						StartingAfter: &models.ThingWithEnumHashKey{
							Branch: models.BranchMaster,
							Date2:  mustTime("2018-03-11T15:04:01+07:00"),
							Date:   mustTime("2018-03-11T15:04:01+07:00"),
						},
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
			{
				testName: "starting after",
				d:        d,
				input: getThingWithMatchingKeyssByBearAndAssocTypeIDInput{
					ctx: context.Background(),
					input: db.GetThingWithMatchingKeyssByBearAndAssocTypeIDInput{
						Bear: "string1",
						StartingAfter: &models.ThingWithMatchingKeys{
							Bear:      "string1",
							AssocType: "string1",
							AssocID:   "string1",
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
			{
				testName: "starting after",
				d:        d,
				input: getThingWithMatchingKeyssByAssocTypeIDAndCreatedBearInput{
					ctx: context.Background(),
					input: db.GetThingWithMatchingKeyssByAssocTypeIDAndCreatedBearInput{
						AssocType: "string1",
						AssocID:   "string1",
						StartingAfter: &models.ThingWithMatchingKeys{
							AssocType: "string1",
							AssocID:   "string1",
							Created:   mustTime("2018-03-11T15:04:01+07:00"),
							Bear:      "string1",
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
			{
				testName: "starting after",
				d:        d,
				input: getThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThreeInput{
					ctx: context.Background(),
					input: db.GetThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThreeInput{
						PropertyOne: "string1",
						PropertyTwo: "string1",
						StartingAfter: &models.ThingWithRequiredCompositePropertiesAndKeysOnly{
							PropertyOne:   db.String("string1"),
							PropertyTwo:   db.String("string1"),
							PropertyThree: db.String("string1"),
						},
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
			{
				testName: "starting after",
				d:        d,
				input: getThingWithRequiredFields2sByNameAndIDInput{
					ctx: context.Background(),
					input: db.GetThingWithRequiredFields2sByNameAndIDInput{
						Name: "string1",
						StartingAfter: &models.ThingWithRequiredFields2{
							Name: db.String("string1"),
							ID:   db.String("string1"),
						},
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
