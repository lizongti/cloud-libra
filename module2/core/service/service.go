package service

import (
	"github.com/lizongti/libra/router"
)

// type CallOption interface {
// }

// type ServerConn interface {
// 	network.Conn
// }

// // type ClientStream interface {
// // }

// // NewStream(context.Context, StreamHandler, opts ...CallOption) (ClientStream, error)

// type ClientConn interface {
// 	network.Conn
// 	Invoke(context.Context, string, interface{}, interface{}, ...CallOption) error
// }

// NEW ->
// STARTING ->
// RUNNING ->
// STOPPING ->
// TERMINATED

// type Listener interface {
// 		failed(Service.State from, Throwable failure)
// // Called when the service transitions to the FAILED state.
// 	running()
// // Called when the service transitions from STARTING to RUNNING.
// 	starting()
// // Called when the service transitions from NEW to STARTING.
// 	stopping(Service.State from)
// // Called when the service transitions to the STOPPING state.
// 	terminated(Service.State from)
// // Called when the service transitions to the TERMINATED state.
// }

// type Service interface {
// 	// Registers a Service.Listener to be executed on the given executor.
// 	AddListener(Service.Listener listener, Executor executor)
// 	// Waits for the Service to reach the running state.
// 	AwaitRunning()
// 	// Waits for the Service to reach the running state for no more than the given time.
// 	AwaitRunning(long timeout, TimeUnit unit)
// 	// Waits for the Service to reach the terminated state.
// 	AwaitTerminated()
// 	// Waits for the Service to reach a terminal state (either terminated or failed) for no more than the given time.

// 	AwaitTerminated(long timeout, TimeUnit unit)
// 	// Returns true if this service is running.

// 	IsRunning() bool
// 	// If the service state is Service.State.NEW, this initiates service startup and returns immediately.

// 	StartAsync() Service

// 	// Returns the lifecycle state of the service.
// 	State() string

// 	// If the service is starting or running, this initiates service shutdown and returns immediately.
// 	StopAsync() Service

// 	RegisterComponent(c component.Component)
// }

type Service interface {
	router.Router
}

// NEW ->
// STARTING ->
// RUNNING ->
// STOPPING ->
// TERMINATED
