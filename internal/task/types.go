package task

const (
	TypeContainerCreate = "container:create"
	TypeExperimentExec  = "experiment:exec"
)

var (
	runningMessage string
	pendingMessage string
)

func init() {
	// json 序列化 Running和Pending消息
	//{
	//	"message": "Running",
	//}
	runningMessage = `{"message": "Running"}`
	pendingMessage = `{"message": "Pending"}`
}
