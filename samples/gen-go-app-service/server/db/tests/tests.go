package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/Clever/wag/samples/gen-go-app-service/models"
	"github.com/Clever/wag/samples/gen-go-app-service/server/db"
	"github.com/stretchr/testify/require"
)

func RunDBTests(t *testing.T, dbFactory func() db.Interface) {
	t.Run("GetSetupStep", GetSetupStep(dbFactory(), t))
	t.Run("GetSetupStepsByAppIDAndID", GetSetupStepsByAppIDAndID(dbFactory(), t))
	t.Run("SaveSetupStep", SaveSetupStep(dbFactory(), t))
	t.Run("DeleteSetupStep", DeleteSetupStep(dbFactory(), t))
}

func GetSetupStep(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.SetupStep{
			AppID: "string1",
			ID:    "string1",
		}
		require.Nil(t, s.SaveSetupStep(ctx, m))
		m2, err := s.GetSetupStep(ctx, m.AppID, m.ID)
		require.Nil(t, err)
		require.Equal(t, m.AppID, m2.AppID)
		require.Equal(t, m.ID, m2.ID)

		_, err = s.GetSetupStep(ctx, "string2", "string2")
		require.NotNil(t, err)
		require.IsType(t, err, db.ErrSetupStepNotFound{})
	}
}

type getSetupStepsByAppIDAndIDInput struct {
	ctx   context.Context
	input db.GetSetupStepsByAppIDAndIDInput
}
type getSetupStepsByAppIDAndIDOutput struct {
	setupSteps []models.SetupStep
	err        error
}
type getSetupStepsByAppIDAndIDTest struct {
	testName string
	d        db.Interface
	input    getSetupStepsByAppIDAndIDInput
	output   getSetupStepsByAppIDAndIDOutput
}

func (g getSetupStepsByAppIDAndIDTest) run(t *testing.T) {
	setupSteps := []models.SetupStep{}
	fn := func(m *models.SetupStep, lastSetupStep bool) bool {
		setupSteps = append(setupSteps, *m)
		if lastSetupStep {
			return false
		}
		return true
	}
	err := g.d.GetSetupStepsByAppIDAndID(g.input.ctx, g.input.input, fn)
	if err != nil {
		fmt.Println(err.Error())
	}
	require.Equal(t, g.output.err, err)
	require.Equal(t, g.output.setupSteps, setupSteps)
}

func GetSetupStepsByAppIDAndID(d db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		require.Nil(t, d.SaveSetupStep(ctx, models.SetupStep{
			AppID: "string1",
			ID:    "string1",
		}))
		require.Nil(t, d.SaveSetupStep(ctx, models.SetupStep{
			AppID: "string1",
			ID:    "string2",
		}))
		require.Nil(t, d.SaveSetupStep(ctx, models.SetupStep{
			AppID: "string1",
			ID:    "string3",
		}))
		limit := int64(3)
		tests := []getSetupStepsByAppIDAndIDTest{
			{
				testName: "basic",
				d:        d,
				input: getSetupStepsByAppIDAndIDInput{
					ctx: context.Background(),
					input: db.GetSetupStepsByAppIDAndIDInput{
						AppID: "string1",
						Limit: &limit,
					},
				},
				output: getSetupStepsByAppIDAndIDOutput{
					setupSteps: []models.SetupStep{
						models.SetupStep{
							AppID: "string1",
							ID:    "string1",
						},
						models.SetupStep{
							AppID: "string1",
							ID:    "string2",
						},
						models.SetupStep{
							AppID: "string1",
							ID:    "string3",
						},
					},
					err: nil,
				},
			},
			{
				testName: "descending",
				d:        d,
				input: getSetupStepsByAppIDAndIDInput{
					ctx: context.Background(),
					input: db.GetSetupStepsByAppIDAndIDInput{
						AppID:      "string1",
						Descending: true,
					},
				},
				output: getSetupStepsByAppIDAndIDOutput{
					setupSteps: []models.SetupStep{
						models.SetupStep{
							AppID: "string1",
							ID:    "string3",
						},
						models.SetupStep{
							AppID: "string1",
							ID:    "string2",
						},
						models.SetupStep{
							AppID: "string1",
							ID:    "string1",
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting after",
				d:        d,
				input: getSetupStepsByAppIDAndIDInput{
					ctx: context.Background(),
					input: db.GetSetupStepsByAppIDAndIDInput{
						AppID: "string1",
						StartingAfter: &models.SetupStep{
							AppID: "string1",
							ID:    "string1",
						},
					},
				},
				output: getSetupStepsByAppIDAndIDOutput{
					setupSteps: []models.SetupStep{
						models.SetupStep{
							AppID: "string1",
							ID:    "string2",
						},
						models.SetupStep{
							AppID: "string1",
							ID:    "string3",
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting after descending",
				d:        d,
				input: getSetupStepsByAppIDAndIDInput{
					ctx: context.Background(),
					input: db.GetSetupStepsByAppIDAndIDInput{
						AppID: "string1",
						StartingAfter: &models.SetupStep{
							AppID: "string1",
							ID:    "string3",
						},
						Descending: true,
					},
				},
				output: getSetupStepsByAppIDAndIDOutput{
					setupSteps: []models.SetupStep{
						models.SetupStep{
							AppID: "string1",
							ID:    "string2",
						},
						models.SetupStep{
							AppID: "string1",
							ID:    "string1",
						},
					},
					err: nil,
				},
			},
			{
				testName: "starting at",
				d:        d,
				input: getSetupStepsByAppIDAndIDInput{
					ctx: context.Background(),
					input: db.GetSetupStepsByAppIDAndIDInput{
						AppID:        "string1",
						IDStartingAt: db.String("string2"),
					},
				},
				output: getSetupStepsByAppIDAndIDOutput{
					setupSteps: []models.SetupStep{
						models.SetupStep{
							AppID: "string1",
							ID:    "string2",
						},
						models.SetupStep{
							AppID: "string1",
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

func SaveSetupStep(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.SetupStep{
			AppID: "string1",
			ID:    "string1",
		}
		require.Nil(t, s.SaveSetupStep(ctx, m))
	}
}

func DeleteSetupStep(s db.Interface, t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		m := models.SetupStep{
			AppID: "string1",
			ID:    "string1",
		}
		require.Nil(t, s.SaveSetupStep(ctx, m))
		require.Nil(t, s.DeleteSetupStep(ctx, m.AppID, m.ID))
	}
}
