package services_test

import (
	"encoder/application/repositories"
	"encoder/application/services"
	"encoder/domain"
	"encoder/framework/database"
	"log"
	"testing"
	"time"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func init(){
	err := godotenv.Load("../../.env")
	if err != nil{
		log.Fatalf("error loading .env file")
	}
}

func prepare() (*domain.Video, repositories.VideoRepositoryDb){
	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "videos/260310fc-53f1-4d92-8b81-81b25613537f/video.mp4"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{Db: db}
	repo.Insert(video)

	return video, repo
}

func TestVideoServiceDownloadAndFragment(t *testing.T){
	video, repo := prepare()

	videoService := services.NewVideoService()
	videoService.Video = video
	videoService.VideoRepository = repo 

	err := videoService.Download("codeflix-encoder")
	require.Nil(t, err)

	err = videoService.Fragment()
	require.Nil(t, err)

	err = videoService.Encode()
	require.Nil(t, err)

	err = videoService.Finish()
	require.Nil(t, err)

}