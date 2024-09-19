package services

import (
	"encoder/domain"
	"encoder/framework/utils"
	"encoding/json"
	"os"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"
)

type JobWorkerResult struct {
	Job 		domain.Job
	Message 	*amqp.Delivery
	Error 		error
}

func JobWorker(messageChannel chan amqp.Delivery, returnChan chan JobWorkerResult, jobService JobService, job domain.Job, workerID int){
	// {   
	// 	"resource_id": "260310fc-53f1-4d92-8b81-81b25613537f.VIDEO", 
	// 	"file_path": "videos/260310fc-53f1-4d92-8b81-81b25613537f/video.mp4"
	// }

	for message := range messageChannel{
		err := utils.IsJson(string(message.Body))
		if err != nil{
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		err = json.Unmarshal(message.Body, &jobService.VideoService.Video)
		jobService.VideoService.Video.ID = uuid.NewV4().String()
		if err != nil{
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		err = jobService.VideoService.Video.Validate()
		if err != nil{
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		err = jobService.VideoService.InsertVideo()
		if err != nil{
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		job.Video = jobService.VideoService.Video
		job.OutputBucketPath = os.Getenv("outputBucketName")
		job.ID = uuid.NewV4().String()
		job.Status = "STARTING"
		job.CreatedAt = time.Now()

		_, err = jobService.JobRepository.Insert(&job)
		if err != nil{
			returnChan <- returnJobResult(job, message, err)
			continue
		}

		jobService.Job = &job
		err = jobService.Start()
		if err != nil{
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		returnChan <- returnJobResult(job, message, nil)
		
		

	}
}


func returnJobResult(job domain.Job, message amqp.Delivery, err error) JobWorkerResult{
	result := JobWorkerResult{
		Job: job,
		Message: &message,
		Error: err,
	}
	return result
}
