package task

import "fmt"

const (
	TypeContainerCreate = "container:create"
	TypeContainerExec   = "container:exec"
)

var (
	ContainerCreateChannelName = func(userID, templateID uint) string {
		return fmt.Sprintf("%d:%d:create", userID, templateID)
	}
	ContainerExecChannelName = func(userID, experimentID uint) string {
		return fmt.Sprintf("%d:%d:exec", userID, experimentID)
	}
)

var (
	runningMessage string
	pendingMessage string
	passMessage    string
	failMessage    string
)

func init() {
	// json 序列化 Running和Pending消息
	//{
	//	"message": "Running",
	//}
	runningMessage = `{"status": "Running"}`
	pendingMessage = `{"status": "Pending"}`
	passMessage = `{"status": "Pass"}`
	failMessage = `{"status": "Fail"}`

}
