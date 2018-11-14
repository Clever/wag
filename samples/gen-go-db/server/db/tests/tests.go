package tests

import (
	"context"
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
	t.Run("GetSimpleThing", GetSimpleThing(dbFactory(), t))
	t.Run("SaveSimpleThing", SaveSimpleThing(dbFactory(), t))
	t.Run("DeleteSimpleThing", DeleteSimpleThing(dbFactory(), t))
	t.Run("GetTeacherSharingRule", GetTeacherSharingRule(dbFactory(), t))
	t.Run("GetTeacherSharingRulesByTeacherAndSchoolApp", GetTeacherSharingRulesByTeacherAndSchoolApp(dbFactory(), t))
	t.Run("SaveTeacherSharingRule", SaveTeacherSharingRule(dbFactory(), t))
	t.Run("DeleteTeacherSharingRule", DeleteTeacherSharingRule(dbFactory(), t))
	t.Run("GetTeacherSharingRulesByDistrictAndSchoolTeacherApp", GetTeacherSharingRulesByDistrictAndSchoolTeacherApp(dbFactory(), t))
	t.Run("GetThing", GetThing(dbFactory(), t))
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
	t.Run("GetThingWithDateRange", GetThingWithDateRange(dbFactory(), t))
	t.Run("GetThingWithDateRangesByNameAndDate", GetThingWithDateRangesByNameAndDate(dbFactory(), t))
	t.Run("SaveThingWithDateRange", SaveThingWithDateRange(dbFactory(), t))
	t.Run("DeleteThingWithDateRange", DeleteThingWithDateRange(dbFactory(), t))
	t.Run("GetThingWithUnderscores", GetThingWithUnderscores(dbFactory(), t))
	t.Run("SaveThingWithUnderscores", SaveThingWithUnderscores(dbFactory(), t))
	t.Run("DeleteThingWithUnderscores", DeleteThingWithUnderscores(dbFactory(), t))
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
		require.Equal(t, db.ErrSimpleThingAlreadyExists{
			Name: "string1",
		}, s.SaveSimpleThing(ctx, m))
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
			Teacher: "string1",
			School:  "string1",
			App:     "string1",
			// must specify non-empty string values for attributes
			// in secondary indexes, since dynamodb doesn't support
			// empty strings:
			District: "district",
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
	teacherSharingRules, err := g.d.GetTeacherSharingRulesByTeacherAndSchoolApp(g.input.ctx, g.input.input)
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
						Teacher: "string1",
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
			Teacher: "string1",
			School:  "string1",
			App:     "string1",
			// must specify non-empty string values for attributes
			// in secondary indexes, since dynamodb doesn't support
			// empty strings:
			District: "district",
		}
		require.Nil(t, s.SaveTeacherSharingRule(ctx, m))
	}
}

func DeleteTeacherSharingRule(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.TeacherSharingRule{
			Teacher: "string1",
			School:  "string1",
			App:     "string1",
			// must specify non-empty string values for attributes
			// in secondary indexes, since dynamodb doesn't support
			// empty strings:
			District: "district",
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
	teacherSharingRules, err := g.d.GetTeacherSharingRulesByDistrictAndSchoolTeacherApp(g.input.ctx, g.input.input)
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
						District: "string1",
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
			Name:    "string1",
			Version: 1,
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
	things, err := g.d.GetThingsByNameAndVersion(g.input.ctx, g.input.input)
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
						Name: "string1",
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

func SaveThing(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.Thing{
			Name:    "string1",
			Version: 1,
		}
		require.Nil(t, s.SaveThing(ctx, m))
		require.Equal(t, db.ErrThingAlreadyExists{
			Name:    "string1",
			Version: 1,
		}, s.SaveThing(ctx, m))
	}
}

func DeleteThing(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.Thing{
			Name:    "string1",
			Version: 1,
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
	things, err := g.d.GetThingsByNameAndCreatedAt(g.input.ctx, g.input.input)
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
						Name: "string1",
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

func GetThingWithCompositeAttributes(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithCompositeAttributes{
			Name:   "string1",
			Branch: "string1",
			Date:   mustTime("2018-03-11T15:04:01+07:00"),
		}
		require.Nil(t, s.SaveThingWithCompositeAttributes(ctx, m))
		m2, err := s.GetThingWithCompositeAttributes(ctx, m.Name, m.Branch, m.Date)
		require.Nil(t, err)
		require.Equal(t, m.Name, m2.Name)
		require.Equal(t, m.Branch, m2.Branch)
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
	thingWithCompositeAttributess, err := g.d.GetThingWithCompositeAttributessByNameBranchAndDate(g.input.ctx, g.input.input)
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.thingWithCompositeAttributess, thingWithCompositeAttributess)
}

func GetThingWithCompositeAttributessByNameBranchAndDate(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithCompositeAttributes(ctx, models.ThingWithCompositeAttributes{
			Name:   "string1",
			Branch: "string1",
			Date:   mustTime("2018-03-11T15:04:01+07:00"),
		}))
		require.Nil(t, d.SaveThingWithCompositeAttributes(ctx, models.ThingWithCompositeAttributes{
			Name:   "string1",
			Branch: "string1",
			Date:   mustTime("2018-03-11T15:04:02+07:00"),
		}))
		require.Nil(t, d.SaveThingWithCompositeAttributes(ctx, models.ThingWithCompositeAttributes{
			Name:   "string1",
			Branch: "string1",
			Date:   mustTime("2018-03-11T15:04:03+07:00"),
		}))
		tests := []getThingWithCompositeAttributessByNameBranchAndDateTest{
			{
				testName: "basic",
				d:        d,
				input: getThingWithCompositeAttributessByNameBranchAndDateInput{
					ctx: context.Background(),
					input: db.GetThingWithCompositeAttributessByNameBranchAndDateInput{
						Name:   "string1",
						Branch: "string1",
					},
				},
				output: getThingWithCompositeAttributessByNameBranchAndDateOutput{
					thingWithCompositeAttributess: []models.ThingWithCompositeAttributes{
						models.ThingWithCompositeAttributes{
							Name:   "string1",
							Branch: "string1",
							Date:   mustTime("2018-03-11T15:04:01+07:00"),
						},
						models.ThingWithCompositeAttributes{
							Name:   "string1",
							Branch: "string1",
							Date:   mustTime("2018-03-11T15:04:02+07:00"),
						},
						models.ThingWithCompositeAttributes{
							Name:   "string1",
							Branch: "string1",
							Date:   mustTime("2018-03-11T15:04:03+07:00"),
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
							Name:   "string1",
							Branch: "string1",
							Date:   mustTime("2018-03-11T15:04:03+07:00"),
						},
						models.ThingWithCompositeAttributes{
							Name:   "string1",
							Branch: "string1",
							Date:   mustTime("2018-03-11T15:04:02+07:00"),
						},
						models.ThingWithCompositeAttributes{
							Name:   "string1",
							Branch: "string1",
							Date:   mustTime("2018-03-11T15:04:01+07:00"),
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
						Name:           "string1",
						Branch:         "string1",
						DateStartingAt: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
					},
				},
				output: getThingWithCompositeAttributessByNameBranchAndDateOutput{
					thingWithCompositeAttributess: []models.ThingWithCompositeAttributes{
						models.ThingWithCompositeAttributes{
							Name:   "string1",
							Branch: "string1",
							Date:   mustTime("2018-03-11T15:04:02+07:00"),
						},
						models.ThingWithCompositeAttributes{
							Name:   "string1",
							Branch: "string1",
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

func SaveThingWithCompositeAttributes(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithCompositeAttributes{
			Name:   "string1",
			Branch: "string1",
			Date:   mustTime("2018-03-11T15:04:01+07:00"),
		}
		require.Nil(t, s.SaveThingWithCompositeAttributes(ctx, m))
	}
}

func DeleteThingWithCompositeAttributes(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithCompositeAttributes{
			Name:   "string1",
			Branch: "string1",
			Date:   mustTime("2018-03-11T15:04:01+07:00"),
		}
		require.Nil(t, s.SaveThingWithCompositeAttributes(ctx, m))
		require.Nil(t, s.DeleteThingWithCompositeAttributes(ctx, m.Name, m.Branch, m.Date))
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
	thingWithCompositeAttributess, err := g.d.GetThingWithCompositeAttributessByNameVersionAndDate(g.input.ctx, g.input.input)
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.thingWithCompositeAttributess, thingWithCompositeAttributess)
}

func GetThingWithCompositeAttributessByNameVersionAndDate(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveThingWithCompositeAttributes(ctx, models.ThingWithCompositeAttributes{
			Name:    "string1",
			Version: 1,
			Date:    mustTime("2018-03-11T15:04:01+07:00"),
			Branch:  "string1",
		}))
		require.Nil(t, d.SaveThingWithCompositeAttributes(ctx, models.ThingWithCompositeAttributes{
			Name:    "string1",
			Version: 1,
			Date:    mustTime("2018-03-11T15:04:02+07:00"),
			Branch:  "string3",
		}))
		require.Nil(t, d.SaveThingWithCompositeAttributes(ctx, models.ThingWithCompositeAttributes{
			Name:    "string1",
			Version: 1,
			Date:    mustTime("2018-03-11T15:04:03+07:00"),
			Branch:  "string2",
		}))
		tests := []getThingWithCompositeAttributessByNameVersionAndDateTest{
			{
				testName: "basic",
				d:        d,
				input: getThingWithCompositeAttributessByNameVersionAndDateInput{
					ctx: context.Background(),
					input: db.GetThingWithCompositeAttributessByNameVersionAndDateInput{
						Name:    "string1",
						Version: 1,
					},
				},
				output: getThingWithCompositeAttributessByNameVersionAndDateOutput{
					thingWithCompositeAttributess: []models.ThingWithCompositeAttributes{
						models.ThingWithCompositeAttributes{
							Name:    "string1",
							Version: 1,
							Date:    mustTime("2018-03-11T15:04:01+07:00"),
							Branch:  "string1",
						},
						models.ThingWithCompositeAttributes{
							Name:    "string1",
							Version: 1,
							Date:    mustTime("2018-03-11T15:04:02+07:00"),
							Branch:  "string3",
						},
						models.ThingWithCompositeAttributes{
							Name:    "string1",
							Version: 1,
							Date:    mustTime("2018-03-11T15:04:03+07:00"),
							Branch:  "string2",
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
							Name:    "string1",
							Version: 1,
							Date:    mustTime("2018-03-11T15:04:03+07:00"),
							Branch:  "string2",
						},
						models.ThingWithCompositeAttributes{
							Name:    "string1",
							Version: 1,
							Date:    mustTime("2018-03-11T15:04:02+07:00"),
							Branch:  "string3",
						},
						models.ThingWithCompositeAttributes{
							Name:    "string1",
							Version: 1,
							Date:    mustTime("2018-03-11T15:04:01+07:00"),
							Branch:  "string1",
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
						Name:           "string1",
						Version:        1,
						DateStartingAt: db.DateTime(mustTime("2018-03-11T15:04:02+07:00")),
					},
				},
				output: getThingWithCompositeAttributessByNameVersionAndDateOutput{
					thingWithCompositeAttributess: []models.ThingWithCompositeAttributes{
						models.ThingWithCompositeAttributes{
							Name:    "string1",
							Version: 1,
							Date:    mustTime("2018-03-11T15:04:02+07:00"),
							Branch:  "string3",
						},
						models.ThingWithCompositeAttributes{
							Name:    "string1",
							Version: 1,
							Date:    mustTime("2018-03-11T15:04:03+07:00"),
							Branch:  "string2",
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
			Name: "string1",
			Date: mustTime("2018-03-11T15:04:01+07:00"),
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
	thingWithDateRanges, err := g.d.GetThingWithDateRangesByNameAndDate(g.input.ctx, g.input.input)
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
						Name: "string1",
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
			Name: "string1",
			Date: mustTime("2018-03-11T15:04:01+07:00"),
		}
		require.Nil(t, s.SaveThingWithDateRange(ctx, m))
	}
}

func DeleteThingWithDateRange(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithDateRange{
			Name: "string1",
			Date: mustTime("2018-03-11T15:04:01+07:00"),
		}
		require.Nil(t, s.SaveThingWithDateRange(ctx, m))
		require.Nil(t, s.DeleteThingWithDateRange(ctx, m.Name, m.Date))
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
