package runner

/**
  @author: victor2022
  @since: 2025/1/27
*/

// Runner 执行器
type Runner interface {
	MarkStarted()
	Handle()
	IsOn() bool
	Stop()
}
