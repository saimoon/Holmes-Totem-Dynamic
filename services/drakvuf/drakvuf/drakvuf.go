package drakvuf

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type Drakvuf struct {
	IncomingDir   string
	ProcessingDir string
	FinishedDir   string
	PendingTasks  map[string]string
}

type StatusResp struct {
	PendingTaskNum int `json:"pendingtasknum"`
}

type TasksReport struct {
	TaskID string `json;"taskid"`
	Log    []byte `json:"log"`
}

// Drakvuf constructor
func New(incomingDir string, processingDir string, finishedDir string) (*Drakvuf, error) {
	r := &Drakvuf{
		IncomingDir:   incomingDir,
		ProcessingDir: processingDir,
		FinishedDir:   finishedDir,
	}
	r.PendingTasks = make(map[string]string)
	return r, nil
}

// GetStatus: return Drakfuv service status (pending tasks, disk free space)
func (c *Drakvuf) GetStatus() (*StatusResp, error) {
	r := &StatusResp{
		PendingTaskNum: len(c.PendingTasks),
	}

	/* TODO: calculate free disk space and return error if low
	freeDiskSpace, err := CalculateDiskSpace()
	if err != nil {
		err = errors.New(fmt.Sprintf("Alert: Free Disk Space %d", freeDiskSpace))
		return nil, err
	}
	*/

	return r, nil
}

// submitTask submits a new task to the drakvuf api.
func (c *Drakvuf) NewTask(fileBytes []byte, fileName string) (string, error) {
	// calculate MD5 sum on sample content
	hash := md5.Sum(fileBytes)
	md5hash := hex.EncodeToString(hash[:])

	// check if task is already pending
	_, ok := c.PendingTasks[md5hash]
	if ok {
		err := errors.New(fmt.Sprintf("Task already pending: ", md5hash))
		return "", err
	}

	// check if file exists already
	_, err := os.Stat(c.IncomingDir + "/" + fileName)
	if err == nil {
		err = errors.New(fmt.Sprintf("Sample already exists in incoming"))
		return "", err
	}

	// save file to Drakvuf incoming dir
	err = ioutil.WriteFile(c.IncomingDir+"/"+fileName, fileBytes, 0644)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error writing sample in incoming"))
		return "", err
	}

	// update pending task map
	c.PendingTasks[md5hash] = fileName

	// md5 hash of the sample is the taskID
	return md5hash, nil
}

func (c *Drakvuf) TaskStatus(taskID string) (int, error) {
	// check if task exists
	fileName, ok := c.PendingTasks[taskID]
	if !ok {
		err := errors.New(fmt.Sprintf("Task not pending: ", taskID))
		return 0, err
	}

	// check if taskID dir exists in Drakvuf finished dir
	_, err := os.Stat(c.ProcessingDir + "/" + fileName)
	if err == nil {
		// Found, task already pending
		return 0, nil
	}

	// Task finished
	return 1, nil
}

func (c *Drakvuf) TaskReport(taskID string) (*TasksReport, error) {
	r := &TasksReport{
		TaskID: taskID,
	}

	// check if task exists
	_, ok := c.PendingTasks[taskID]
	if !ok {
		err := errors.New(fmt.Sprintf("Task not pending: ", taskID))
		return nil, err
	}

	// get Drakvuf result
	resultBytes, err := ioutil.ReadFile(c.FinishedDir + "/" + taskID + "/drakvuf.log")
	if err != nil {
		err = errors.New(fmt.Sprintf("Drakvuf result not found: ", taskID))
		return nil, err
	}

	// remove from pending task
	delete(c.PendingTasks, taskID)

	r.Log = resultBytes

	return r, nil
}

func (c *Drakvuf) DeleteTask(taskID string) error {
	return nil
}
