package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/bytebase/bytebase/backend/common"
	api "github.com/bytebase/bytebase/backend/legacyapi"
)

// GetPipelineByID gets an instance of Pipeline.
func (s *Store) GetPipelineByID(ctx context.Context, id int) (*api.Pipeline, error) {
	pipeline, err := s.GetPipelineV2ByID(ctx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get pipeline with ID %d", id)
	}
	if pipeline == nil {
		return nil, nil
	}
	composedPipeline, err := s.composePipeline(ctx, pipeline)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to compose pipeline")
	}
	return composedPipeline, nil
}

func (s *Store) composePipeline(ctx context.Context, pipeline *PipelineMessage) (*api.Pipeline, error) {
	composedPipeline := &api.Pipeline{
		ID:   pipeline.ID,
		Name: pipeline.Name,
	}

	tasks, err := s.FindTask(ctx, &api.TaskFind{PipelineID: &pipeline.ID})
	if err != nil {
		return nil, err
	}

	stages, err := s.ListStageV2(ctx, pipeline.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find stages for pipeline %d", pipeline.ID)
	}
	var composedStages []*api.Stage
	for _, stage := range stages {
		environment, err := s.GetEnvironmentByID(ctx, stage.EnvironmentID)
		if err != nil {
			return nil, err
		}
		composedStage := &api.Stage{
			ID:            stage.ID,
			PipelineID:    stage.PipelineID,
			EnvironmentID: stage.EnvironmentID,
			Environment:   environment,
			Name:          stage.Name,
		}
		for _, task := range tasks {
			if task.StageID == stage.ID {
				composedStage.TaskList = append(composedStage.TaskList, task)
			}
		}

		composedStages = append(composedStages, composedStage)
	}
	composedPipeline.StageList = composedStages

	return composedPipeline, nil
}

func (s *Store) composeSimplePipeline(ctx context.Context, pipeline *PipelineMessage) (*api.Pipeline, error) {
	tasks, err := s.ListTasks(ctx, &api.TaskFind{PipelineID: &pipeline.ID})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find tasks for pipeline %d", pipeline.ID)
	}

	stages, err := s.ListStageV2(ctx, pipeline.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find stage list")
	}
	var composedStages []*api.Stage
	for _, stage := range stages {
		environment, err := s.GetEnvironmentByID(ctx, stage.EnvironmentID)
		if err != nil {
			return nil, err
		}
		composedStage := &api.Stage{
			ID:            stage.ID,
			Name:          stage.Name,
			EnvironmentID: stage.EnvironmentID,
			Environment:   environment,
			PipelineID:    stage.PipelineID,
		}

		for _, task := range tasks {
			if task.StageID == stage.ID {
				composedStage.TaskList = append(composedStage.TaskList, task.toTask())
			}
		}
		composedStages = append(composedStages, composedStage)
	}

	return &api.Pipeline{
		ID:        pipeline.ID,
		Name:      pipeline.Name,
		StageList: composedStages,
	}, nil
}

// PipelineMessage is the message for pipelines.
type PipelineMessage struct {
	Name string
	// Output only.
	ID int
}

// PipelineFind is the API message for finding pipelines.
type PipelineFind struct {
	ID *int
}

// CreatePipelineV2 creates a pipeline.
func (s *Store) CreatePipelineV2(ctx context.Context, create *PipelineMessage, creatorID int) (*PipelineMessage, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, FormatError(err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO pipeline (
			creator_id,
			updater_id,
			name
		)
		VALUES ($1, $2, $3)
		RETURNING id, name
	`
	pipeline := &PipelineMessage{}
	if err := tx.QueryRowContext(ctx, query,
		creatorID,
		creatorID,
		create.Name,
	).Scan(
		&pipeline.ID,
		&pipeline.Name,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, common.FormatDBErrorEmptyRowWithQuery(query)
		}
		return nil, FormatError(err)
	}

	if err := tx.Commit(); err != nil {
		return nil, FormatError(err)
	}

	s.pipelineCache.Store(pipeline.ID, pipeline)
	return pipeline, nil
}

// GetPipelineV2ByID gets the pipeline by ID.
func (s *Store) GetPipelineV2ByID(ctx context.Context, id int) (*PipelineMessage, error) {
	if pipeline, ok := s.pipelineCache.Load(id); ok {
		return pipeline.(*PipelineMessage), nil
	}
	pipelines, err := s.ListPipelineV2(ctx, &PipelineFind{ID: &id})
	if err != nil {
		return nil, err
	}

	if len(pipelines) == 0 {
		return nil, nil
	} else if len(pipelines) > 1 {
		return nil, &common.Error{Code: common.Conflict, Err: errors.Errorf("found %d pipelines, expect 1", len(pipelines))}
	}
	pipeline := pipelines[0]
	return pipeline, nil
}

// ListPipelineV2 lists pipelines.
func (s *Store) ListPipelineV2(ctx context.Context, find *PipelineFind) ([]*PipelineMessage, error) {
	where, args := []string{"TRUE"}, []interface{}{}
	if v := find.ID; v != nil {
		where, args = append(where, fmt.Sprintf("pipeline.id = $%d", len(args)+1)), append(args, *v)
	}
	query := fmt.Sprintf(`
		SELECT
			pipeline.id,
			pipeline.name
		FROM pipeline
		WHERE %s`, strings.Join(where, " AND "))

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, FormatError(err)
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, FormatError(err)
	}
	defer rows.Close()

	var pipelines []*PipelineMessage
	for rows.Next() {
		var pipeline PipelineMessage
		if err := rows.Scan(
			&pipeline.ID,
			&pipeline.Name,
		); err != nil {
			return nil, FormatError(err)
		}
		pipelines = append(pipelines, &pipeline)
	}
	if err := rows.Err(); err != nil {
		return nil, FormatError(err)
	}

	if err := tx.Commit(); err != nil {
		return nil, FormatError(err)
	}

	for _, pipeline := range pipelines {
		s.pipelineCache.Store(pipeline.ID, pipeline)
	}
	return pipelines, nil
}
