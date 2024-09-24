package repositories_test

import (
	"encoder/application/repositories"
	"encoder/domain"
	"encoder/adapters/database"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestNewJobRepositoryDb(t *testing.T){
	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	repoVideo := repositories.VideoRepositoryDb{Db: db}
	repoVideo.Insert(video)

	job, err := domain.NewJob("output_path", "Pending", video)
	require.Nil(t, err)

	repoJob := repositories.JobRepositoryDb{Db: db}
	repoJob.Insert(job)

	j, err := repoJob.Find(job.ID)
	require.NotEmpty(t, j.ID)
	require.Nil(t, err)
	require.Equal(t, j.ID, job.ID)
	require.Equal(t, j.VideoID, video.ID)

	job.Status =  "Complete"
	repoJob.Update(job)

	j, err = repoJob.Find(job.ID)
	require.NotEmpty(t, j.ID)
	require.Nil(t, err)
	require.Equal(t, j.Status, job.Status)

}