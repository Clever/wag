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

var limit = int64(100)

func mustTime(s string) strfmt.DateTime {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return strfmt.DateTime(t)
}

func RunDBTests(t *testing.T, dbFactory func() db.Interface) {
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
	t.Run("GetThingWithRequiredFields", GetThingWithRequiredFields(dbFactory(), t))
	t.Run("SaveThingWithRequiredFields", SaveThingWithRequiredFields(dbFactory(), t))
	t.Run("DeleteThingWithRequiredFields", DeleteThingWithRequiredFields(dbFactory(), t))
	t.Run("GetThingWithUnderscores", GetThingWithUnderscores(dbFactory(), t))
	t.Run("SaveThingWithUnderscores", SaveThingWithUnderscores(dbFactory(), t))
	t.Run("DeleteThingWithUnderscores", DeleteThingWithUnderscores(dbFactory(), t))
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
		tests := []getNoRangeThingWithCompositeAttributessByNameVersionAndDateTest{
			{
				testName: "basic",
				d:        d,
				input: getNoRangeThingWithCompositeAttributessByNameVersionAndDateInput{
					ctx: context.Background(),
					input: db.GetNoRangeThingWithCompositeAttributessByNameVersionAndDateInput{
						StartingAt: &models.NoRangeThingWithCompositeAttributes{
							Name:    db.String("string1"),
							Version: 1,
							Date:    db.DateTime(mustTime("2018-03-11T15:04:00+07:00")),
							Branch:  db.String("string0"),
						},
						Exclusive: true,
						Limit:     &limit,
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
						StartingAt: &models.NoRangeThingWithCompositeAttributes{
							Name:    db.String("string1"),
							Version: 1,
							Date:    db.DateTime(mustTime("2018-03-11T15:04:04+07:00")),
							Branch:  db.String("string4"),
						},
						Exclusive:  true,
						Descending: true,
						Limit:      &limit,
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
						StartingAt: &models.NoRangeThingWithCompositeAttributes{
							Name:    db.String("string1"),
							Version: 1,
							Date:    db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
							Branch:  db.String("string1"),
						},
						Exclusive: true,
						Limit:     &limit,
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
				testName: "starting at",
				d:        d,
				input: getNoRangeThingWithCompositeAttributessByNameVersionAndDateInput{
					ctx: context.Background(),
					input: db.GetNoRangeThingWithCompositeAttributessByNameVersionAndDateInput{
						StartingAt: &models.NoRangeThingWithCompositeAttributes{
							Name:    db.String("string1"),
							Version: 1,
							Date:    db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
							Branch:  db.String("string1"),
						},
						Limit: &limit,
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
		tests := []getTeacherSharingRulesByTeacherAndSchoolAppTest{
			{
				testName: "basic",
				d:        d,
				input: getTeacherSharingRulesByTeacherAndSchoolAppInput{
					ctx: context.Background(),
					input: db.GetTeacherSharingRulesByTeacherAndSchoolAppInput{
						StartingAt: &models.TeacherSharingRule{
							Teacher:  "string1",
							School:   "string0",
							App:      "string0",
							District: "district",
						},
						Exclusive: true,
						Limit:     &limit,
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
						StartingAt: &models.TeacherSharingRule{
							Teacher:  "string1",
							School:   "string4",
							App:      "string4",
							District: "district",
						},
						Exclusive:  true,
						Descending: true,
						Limit:      &limit,
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
						StartingAt: &models.TeacherSharingRule{
							Teacher:  "string1",
							School:   "string1",
							App:      "string1",
							District: "district",
						},
						Exclusive: true,
						Limit:     &limit,
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
				testName: "starting at",
				d:        d,
				input: getTeacherSharingRulesByTeacherAndSchoolAppInput{
					ctx: context.Background(),
					input: db.GetTeacherSharingRulesByTeacherAndSchoolAppInput{
						StartingAt: &models.TeacherSharingRule{
							Teacher:  "string1",
							School:   "string1",
							App:      "string1",
							District: "district",
						},
						Limit: &limit,
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
		tests := []getTeacherSharingRulesByDistrictAndSchoolTeacherAppTest{
			{
				testName: "basic",
				d:        d,
				input: getTeacherSharingRulesByDistrictAndSchoolTeacherAppInput{
					ctx: context.Background(),
					input: db.GetTeacherSharingRulesByDistrictAndSchoolTeacherAppInput{
						StartingAt: &models.TeacherSharingRule{
							District: "string1",
							School:   "string0",
							Teacher:  "string0",
							App:      "string0",
						},
						Exclusive: true,
						Limit:     &limit,
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
						StartingAt: &models.TeacherSharingRule{
							District: "string1",
							School:   "string4",
							Teacher:  "string4",
							App:      "string4",
						},
						Exclusive:  true,
						Descending: true,
						Limit:      &limit,
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
						StartingAt: &models.TeacherSharingRule{
							District: "string1",
							School:   "string1",
							Teacher:  "string1",
							App:      "string1",
						},
						Exclusive: true,
						Limit:     &limit,
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
				testName: "starting at",
				d:        d,
				input: getTeacherSharingRulesByDistrictAndSchoolTeacherAppInput{
					ctx: context.Background(),
					input: db.GetTeacherSharingRulesByDistrictAndSchoolTeacherAppInput{
						StartingAt: &models.TeacherSharingRule{
							District: "string1",
							School:   "string1",
							Teacher:  "string1",
							App:      "string1",
						},
						Limit: &limit,
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
		tests := []getThingsByNameAndVersionTest{
			{
				testName: "basic",
				d:        d,
				input: getThingsByNameAndVersionInput{
					ctx: context.Background(),
					input: db.GetThingsByNameAndVersionInput{
						StartingAt: &models.Thing{
							Name:    "string1",
							Version: 0,
						},
						Exclusive: true,
						Limit:     &limit,
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
						StartingAt: &models.Thing{
							Name:    "string1",
							Version: 4,
						},
						Exclusive:  true,
						Descending: true,
						Limit:      &limit,
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
						StartingAt: &models.Thing{
							Name:    "string1",
							Version: 1,
						},
						Exclusive: true,
						Limit:     &limit,
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
				testName: "starting at",
				d:        d,
				input: getThingsByNameAndVersionInput{
					ctx: context.Background(),
					input: db.GetThingsByNameAndVersionInput{
						StartingAt: &models.Thing{
							Name:    "string1",
							Version: 1,
						},
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
			Name:    "string1",
			Version: 1,
			ID:      "string1",
		}
		require.Nil(t, s.SaveThing(ctx, m))
		m2, err := s.GetThingByID(ctx, m.ID)
		require.Nil(t, err)
		require.Equal(t, m.Name, m2.Name)
		require.Equal(t, m.Version, m2.Version)
		require.Equal(t, m.ID, m2.ID)

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
		tests := []getThingsByNameAndCreatedAtTest{
			{
				testName: "basic",
				d:        d,
				input: getThingsByNameAndCreatedAtInput{
					ctx: context.Background(),
					input: db.GetThingsByNameAndCreatedAtInput{
						StartingAt: &models.Thing{
							Name:      "string1",
							CreatedAt: mustTime("2018-03-11T15:04:00+07:00"),
							Version:   0,
						},
						Exclusive: true,
						Limit:     &limit,
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
						StartingAt: &models.Thing{
							Name:      "string1",
							CreatedAt: mustTime("2018-03-11T15:04:04+07:00"),
							Version:   4,
						},
						Exclusive:  true,
						Descending: true,
						Limit:      &limit,
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
						StartingAt: &models.Thing{
							Name:      "string1",
							CreatedAt: mustTime("2018-03-11T15:04:01+07:00"),
							Version:   1,
						},
						Exclusive: true,
						Limit:     &limit,
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
				testName: "starting at",
				d:        d,
				input: getThingsByNameAndCreatedAtInput{
					ctx: context.Background(),
					input: db.GetThingsByNameAndCreatedAtInput{
						StartingAt: &models.Thing{
							Name:      "string1",
							CreatedAt: mustTime("2018-03-11T15:04:01+07:00"),
							Version:   1,
						},
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
		tests := []getThingWithCompositeAttributessByNameBranchAndDateTest{
			{
				testName: "basic",
				d:        d,
				input: getThingWithCompositeAttributessByNameBranchAndDateInput{
					ctx: context.Background(),
					input: db.GetThingWithCompositeAttributessByNameBranchAndDateInput{
						StartingAt: &models.ThingWithCompositeAttributes{
							Name:   db.String("string1"),
							Branch: db.String("string1"),
							Date:   db.DateTime(mustTime("2018-03-11T15:04:00+07:00")),
						},
						Exclusive: true,
						Limit:     &limit,
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
						StartingAt: &models.ThingWithCompositeAttributes{
							Name:   db.String("string1"),
							Branch: db.String("string1"),
							Date:   db.DateTime(mustTime("2018-03-11T15:04:04+07:00")),
						},
						Exclusive:  true,
						Descending: true,
						Limit:      &limit,
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
						StartingAt: &models.ThingWithCompositeAttributes{
							Name:   db.String("string1"),
							Branch: db.String("string1"),
							Date:   db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
						},
						Exclusive: true,
						Limit:     &limit,
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
				testName: "starting at",
				d:        d,
				input: getThingWithCompositeAttributessByNameBranchAndDateInput{
					ctx: context.Background(),
					input: db.GetThingWithCompositeAttributessByNameBranchAndDateInput{
						StartingAt: &models.ThingWithCompositeAttributes{
							Name:   db.String("string1"),
							Branch: db.String("string1"),
							Date:   db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
						},
						Limit: &limit,
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
		tests := []getThingWithCompositeAttributessByNameVersionAndDateTest{
			{
				testName: "basic",
				d:        d,
				input: getThingWithCompositeAttributessByNameVersionAndDateInput{
					ctx: context.Background(),
					input: db.GetThingWithCompositeAttributessByNameVersionAndDateInput{
						StartingAt: &models.ThingWithCompositeAttributes{
							Name:    db.String("string1"),
							Version: 1,
							Date:    db.DateTime(mustTime("2018-03-11T15:04:00+07:00")),
							Branch:  db.String("string0"),
						},
						Exclusive: true,
						Limit:     &limit,
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
						StartingAt: &models.ThingWithCompositeAttributes{
							Name:    db.String("string1"),
							Version: 1,
							Date:    db.DateTime(mustTime("2018-03-11T15:04:04+07:00")),
							Branch:  db.String("string4"),
						},
						Exclusive:  true,
						Descending: true,
						Limit:      &limit,
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
						StartingAt: &models.ThingWithCompositeAttributes{
							Name:    db.String("string1"),
							Version: 1,
							Date:    db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
							Branch:  db.String("string1"),
						},
						Exclusive: true,
						Limit:     &limit,
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
				testName: "starting at",
				d:        d,
				input: getThingWithCompositeAttributessByNameVersionAndDateInput{
					ctx: context.Background(),
					input: db.GetThingWithCompositeAttributessByNameVersionAndDateInput{
						StartingAt: &models.ThingWithCompositeAttributes{
							Name:    db.String("string1"),
							Version: 1,
							Date:    db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
							Branch:  db.String("string1"),
						},
						Limit: &limit,
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
		tests := []getThingWithCompositeEnumAttributessByNameBranchAndDateTest{
			{
				testName: "basic",
				d:        d,
				input: getThingWithCompositeEnumAttributessByNameBranchAndDateInput{
					ctx: context.Background(),
					input: db.GetThingWithCompositeEnumAttributessByNameBranchAndDateInput{
						StartingAt: &models.ThingWithCompositeEnumAttributes{
							Name:     db.String("string1"),
							BranchID: models.BranchMaster,
							Date:     db.DateTime(mustTime("2018-03-11T15:04:00+07:00")),
						},
						Exclusive: true,
						Limit:     &limit,
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
						StartingAt: &models.ThingWithCompositeEnumAttributes{
							Name:     db.String("string1"),
							BranchID: models.BranchMaster,
							Date:     db.DateTime(mustTime("2018-03-11T15:04:04+07:00")),
						},
						Exclusive:  true,
						Descending: true,
						Limit:      &limit,
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
						StartingAt: &models.ThingWithCompositeEnumAttributes{
							Name:     db.String("string1"),
							BranchID: models.BranchMaster,
							Date:     db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
						},
						Exclusive: true,
						Limit:     &limit,
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
				testName: "starting at",
				d:        d,
				input: getThingWithCompositeEnumAttributessByNameBranchAndDateInput{
					ctx: context.Background(),
					input: db.GetThingWithCompositeEnumAttributessByNameBranchAndDateInput{
						StartingAt: &models.ThingWithCompositeEnumAttributes{
							Name:     db.String("string1"),
							BranchID: models.BranchMaster,
							Date:     db.DateTime(mustTime("2018-03-11T15:04:01+07:00")),
						},
						Limit: &limit,
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
		tests := []getThingWithDateRangesByNameAndDateTest{
			{
				testName: "basic",
				d:        d,
				input: getThingWithDateRangesByNameAndDateInput{
					ctx: context.Background(),
					input: db.GetThingWithDateRangesByNameAndDateInput{
						StartingAt: &models.ThingWithDateRange{
							Name: "string1",
							Date: mustTime("2018-03-11T15:04:00+07:00"),
						},
						Exclusive: true,
						Limit:     &limit,
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
						StartingAt: &models.ThingWithDateRange{
							Name: "string1",
							Date: mustTime("2018-03-11T15:04:04+07:00"),
						},
						Exclusive:  true,
						Descending: true,
						Limit:      &limit,
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
						StartingAt: &models.ThingWithDateRange{
							Name: "string1",
							Date: mustTime("2018-03-11T15:04:01+07:00"),
						},
						Exclusive: true,
						Limit:     &limit,
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
				testName: "starting at",
				d:        d,
				input: getThingWithDateRangesByNameAndDateInput{
					ctx: context.Background(),
					input: db.GetThingWithDateRangesByNameAndDateInput{
						StartingAt: &models.ThingWithDateRange{
							Name: "string1",
							Date: mustTime("2018-03-11T15:04:01+07:00"),
						},
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
		tests := []getThingWithDateTimeCompositesByTypeIDAndCreatedResourceTest{
			{
				testName: "basic",
				d:        d,
				input: getThingWithDateTimeCompositesByTypeIDAndCreatedResourceInput{
					ctx: context.Background(),
					input: db.GetThingWithDateTimeCompositesByTypeIDAndCreatedResourceInput{
						StartingAt: &models.ThingWithDateTimeComposite{
							Type:     "string1",
							ID:       "string1",
							Created:  mustTime("2018-03-11T15:04:00+07:00"),
							Resource: "string0",
						},
						Exclusive: true,
						Limit:     &limit,
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
						StartingAt: &models.ThingWithDateTimeComposite{
							Type:     "string1",
							ID:       "string1",
							Created:  mustTime("2018-03-11T15:04:04+07:00"),
							Resource: "string4",
						},
						Exclusive:  true,
						Descending: true,
						Limit:      &limit,
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
						StartingAt: &models.ThingWithDateTimeComposite{
							Type:     "string1",
							ID:       "string1",
							Created:  mustTime("2018-03-11T15:04:01+07:00"),
							Resource: "string1",
						},
						Exclusive: true,
						Limit:     &limit,
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
				testName: "starting at",
				d:        d,
				input: getThingWithDateTimeCompositesByTypeIDAndCreatedResourceInput{
					ctx: context.Background(),
					input: db.GetThingWithDateTimeCompositesByTypeIDAndCreatedResourceInput{
						StartingAt: &models.ThingWithDateTimeComposite{
							Type:     "string1",
							ID:       "string1",
							Created:  mustTime("2018-03-11T15:04:01+07:00"),
							Resource: "string1",
						},
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
