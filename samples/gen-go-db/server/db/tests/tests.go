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
	t.Run("GetThing", GetThing(dbFactory(), t))
	t.Run("GetThingsByNameAndVersion", GetThingsByNameAndVersion(dbFactory(), t))
	t.Run("SaveThing", SaveThing(dbFactory(), t))
	t.Run("DeleteThing", DeleteThing(dbFactory(), t))
	t.Run("GetThingByID", GetThingByID(dbFactory(), t))
	t.Run("GetThingsByNameAndCreatedAt", GetThingsByNameAndCreatedAt(dbFactory(), t))
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
