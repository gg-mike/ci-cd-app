package scheduler

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/gg-mike/ci-cd-app/backend/internal/model"
)

func TestSelectWorker(t *testing.T) {
	type fields struct {
		cfg     model.PipelineConfig
		workers []model.Worker
	}
	tests := []struct {
		name    string
		args    fields
		want    model.Worker
		wantErr string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SelectWorker(tt.args.cfg, tt.args.workers)
			if fmt.Sprint(err) != tt.wantErr {
				t.Errorf("error\nactual: %v\nexpect: %v", err, tt.wantErr)
				return
			}
			if got.Name != tt.want.Name {
				t.Errorf("worker\nactual: %v\nexpect: %v", got.Name, tt.want.Name)
			}
		})
	}
}

func Test_filterWorkers(t *testing.T) {
	type fields struct {
		workers []model.Worker
		system  string
		image   string
	}
	tests := []struct {
		name string
		args fields
		want []model.Worker
	}{
		{
			"No requirements",
			fields{
				[]model.Worker{
					{Name: "worker-1", Capacity: 3},
					{Name: "worker-2", Capacity: 3},
				},
				"", "",
			},
			[]model.Worker{
				{Name: "worker-1", Capacity: 3},
				{Name: "worker-2", Capacity: 3},
			},
		},
		{
			"System requirement",
			fields{
				[]model.Worker{
					{Name: "worker-1", Capacity: 3, System: "foo"},
					{Name: "worker-2", Capacity: 3, System: "bar"},
					{Name: "worker-3", Capacity: 3},
				},
				"foo", "",
			},
			[]model.Worker{
				{Name: "worker-1", Capacity: 3, System: "foo"},
			},
		},
		{
			"Image requirement",
			fields{
				[]model.Worker{
					{Name: "worker-1", Capacity: 3, Type: model.WorkerDockerHost},
					{Name: "worker-2", Capacity: 3, Type: model.WorkerStatic},
					{Name: "worker-3", Capacity: 3},
				},
				"", "foo",
			},
			[]model.Worker{
				{Name: "worker-1", Capacity: 3, Type: model.WorkerDockerHost},
			},
		},
		{
			"Capacity requirements",
			fields{
				[]model.Worker{
					{Name: "worker-1", Capacity: 3, ActiveBuilds: 3},
					{Name: "worker-2", Capacity: 3, ActiveBuilds: 2},
					{Name: "worker-3", Capacity: 3, ActiveBuilds: 0},
				},
				"", "",
			},
			[]model.Worker{
				{Name: "worker-2", Capacity: 3, ActiveBuilds: 2},
				{Name: "worker-3", Capacity: 3, ActiveBuilds: 0},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterWorkers(tt.args.workers, tt.args.system, tt.args.image); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("workers\nactual: %v\nexpect: %v", workersToNames(got), workersToNames(tt.want))
			}
		})
	}
}

func Test_sortWorkers(t *testing.T) {
	tests := []struct {
		name string
		arg  []model.Worker
		want []model.Worker
	}{
		{
			"Strategy based difference",
			[]model.Worker{
				{Name: "worker-1", Strategy: model.WorkerBalanced},
				{Name: "worker-2", Strategy: model.WorkerUseMin},
				{Name: "worker-3", Strategy: model.WorkerUseMax},
			},
			[]model.Worker{
				{Name: "worker-3", Strategy: model.WorkerUseMax},
				{Name: "worker-1", Strategy: model.WorkerBalanced},
				{Name: "worker-2", Strategy: model.WorkerUseMin},
			},
		},
		{
			"Active builds based difference",
			[]model.Worker{
				{Name: "worker-1", ActiveBuilds: 6},
				{Name: "worker-2", ActiveBuilds: 3},
				{Name: "worker-3", ActiveBuilds: 2},
			},
			[]model.Worker{
				{Name: "worker-3", ActiveBuilds: 2},
				{Name: "worker-2", ActiveBuilds: 3},
				{Name: "worker-1", ActiveBuilds: 6},
			},
		},
		{
			"Name based difference",
			[]model.Worker{
				{Name: "worker-2"},
				{Name: "worker-3"},
				{Name: "worker-1"},
			},
			[]model.Worker{
				{Name: "worker-1"},
				{Name: "worker-2"},
				{Name: "worker-3"},
			},
		},
		{
			"Combined",
			[]model.Worker{
				{Name: "worker-1", ActiveBuilds: 3, Strategy: model.WorkerBalanced},
				{Name: "worker-2", ActiveBuilds: 1, Strategy: model.WorkerUseMin},
				{Name: "worker-3", ActiveBuilds: 2, Strategy: model.WorkerUseMax},
				{Name: "worker-4", ActiveBuilds: 5, Strategy: model.WorkerBalanced},
				{Name: "worker-5", ActiveBuilds: 2, Strategy: model.WorkerUseMax},
				{Name: "worker-6", ActiveBuilds: 0, Strategy: model.WorkerUseMin},
				{Name: "worker-7", ActiveBuilds: 1, Strategy: model.WorkerBalanced},
			},
			[]model.Worker{
				{Name: "worker-3", ActiveBuilds: 2, Strategy: model.WorkerUseMax},
				{Name: "worker-5", ActiveBuilds: 2, Strategy: model.WorkerUseMax},
				{Name: "worker-7", ActiveBuilds: 1, Strategy: model.WorkerBalanced},
				{Name: "worker-1", ActiveBuilds: 3, Strategy: model.WorkerBalanced},
				{Name: "worker-4", ActiveBuilds: 5, Strategy: model.WorkerBalanced},
				{Name: "worker-6", ActiveBuilds: 0, Strategy: model.WorkerUseMin},
				{Name: "worker-2", ActiveBuilds: 1, Strategy: model.WorkerUseMin},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sortWorkers(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("workers\nactual: %v\nexpect: %v", workersToNames(got), workersToNames(tt.want))
			}
		})
	}
}

func workersToNames(workers []model.Worker) []string {
	names := []string{}
	for _, worker := range workers {
		names = append(names, worker.Name)
	}
	return names
}
