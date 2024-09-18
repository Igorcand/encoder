package utils_test

import (
	"encoder/framework/utils"
	"testing"
	"github.com/stretchr/testify/require"
)

func TestIsJson(t *testing.T){
	json := `{
    "error": "",
    "video": {
        "resource_id": "24846fe1-6218-46bb-96a8-d6d4534e0885.VIDEO",
        "encoded_video_folder": "/path/to/encoded/video"
    },
    "status": "COMPLETED"
	}`

	err := utils.IsJson(json)
	require.Nil(t, err)

	json = "124dc"
	err = utils.IsJson(json)
	require.Error(t, err)
}