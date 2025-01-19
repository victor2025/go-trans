package consts

/**
  @author: victor2022
  @since: 2025/1/19
*/
type TaskStatus int

const (
	Waiting    TaskStatus = 0
	Processing TaskStatus = 1
	Finished   TaskStatus = 2
	Fail       TaskStatus = -1
)
