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

func RunDBTests(t *testing.T, dbFactory func() db.Interface) {
	t.Run("GetSimpleThing", GetSimpleThing(dbFactory(), t))
	t.Run("SaveSimpleThing", SaveSimpleThing(dbFactory(), t))
	t.Run("DeleteSimpleThing", DeleteSimpleThing(dbFactory(), t))
	t.Run("GetThing", GetThing(dbFactory(), t))
	t.Run("SaveThing", SaveThing(dbFactory(), t))
	t.Run("DeleteThing", DeleteThing(dbFactory(), t))
	t.Run("GetThingWithDateRange", GetThingWithDateRange(dbFactory(), t))
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
func GetThingWithDateRange(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithDateRange{
			Name: "string1",
			Date: strfmt.DateTime(time.Unix(1522279646, 0)),
		}
		require.Nil(t, s.SaveThingWithDateRange(ctx, m))
		m2, err := s.GetThingWithDateRange(ctx, m.Name, m.Date)
		require.Nil(t, err)
		require.Equal(t, m.Name, m2.Name)
		require.Equal(t, m.Date.String(), m2.Date.String())

		_, err = s.GetThingWithDateRange(ctx, "string2", strfmt.DateTime(time.Unix(2522279646, 0)))
		require.NotNil(t, err)
		require.IsType(t, err, db.ErrThingWithDateRangeNotFound{})
	}
}

func SaveThingWithDateRange(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithDateRange{
			Name: "string1",
			Date: strfmt.DateTime(time.Unix(1522279646, 0)),
		}
		require.Nil(t, s.SaveThingWithDateRange(ctx, m))
	}
}

func DeleteThingWithDateRange(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.ThingWithDateRange{
			Name: "string1",
			Date: strfmt.DateTime(time.Unix(1522279646, 0)),
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
