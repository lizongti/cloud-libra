package http

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aceaura/libra/core/device"
	"github.com/aceaura/libra/core/encoding"
	"github.com/aceaura/libra/core/magic"
	"github.com/aceaura/libra/core/message"
	"github.com/aceaura/libra/core/route"
	"github.com/aceaura/libra/core/scheduler"
)

var (
	ErrRequestNotFound = errors.New("request not found by name")
)

type Commander struct {
	*device.Client
	opts         commanderOptions
	scheduler    *scheduler.Scheduler
	controller   *scheduler.TPSController
	reqIndex     int
	reqChan      chan *ServiceRequest
	respChan     chan *ServiceResponse
	errorChan    chan error
	reportChan   chan *scheduler.Report
	parallelChan chan int
	dieChan      chan struct{}
	exitChan     chan struct{}
}

func NewCommander(opt ...ApplyCommanderOption) *Commander {
	opts := defaultCommanderOptions
	for _, o := range opt {
		o.apply(&opts)
	}

	return &Commander{
		Client:    device.NewClient(),
		opts:      opts,
		reqIndex:  0,
		errorChan: make(chan error),
		dieChan:   make(chan struct{}),
		exitChan:  make(chan struct{}),
	}
}

func (c *Commander) String() string {
	return c.opts.name
}

func (c *Commander) Close() {
	close(c.dieChan)
	<-c.exitChan
}

func (c *Commander) Serve() error {
	c.reportChan = make(chan *scheduler.Report, c.opts.reportBacklog)
	c.parallelChan = make(chan int, c.opts.parallelBacklog)
	c.reqChan = make(chan *ServiceRequest, c.opts.reqBacklog)
	c.respChan = make(chan *ServiceResponse, c.opts.respBacklog)

	c.scheduler = scheduler.NewScheduler(
		scheduler.SchedulerOption.Safety(),
		scheduler.SchedulerOption.Background(),
		scheduler.SchedulerOption.ErrorChan(c.errorChan),
		scheduler.SchedulerOption.TaskBacklog(c.opts.taskBacklog),
		scheduler.SchedulerOption.Parallel(c.opts.parallelInit),
		scheduler.SchedulerOption.ReportChan(c.reportChan),
		scheduler.SchedulerOption.ParallelChan(c.parallelChan),
	)

	c.controller = scheduler.NewTPSController(
		scheduler.ControllerOption.Safety(),
		scheduler.ControllerOption.Background(),
		scheduler.ControllerOption.ErrorChan(c.errorChan),
		scheduler.ControllerOption.ParallelTick(c.opts.parallelTick),
		scheduler.ControllerOption.ParallelIncrease(c.opts.parallelIncrease),
		scheduler.ControllerOption.TPSLimit(c.opts.tpsLimit),
	)

	if err := c.scheduler.Serve(); err != nil {
		return err
	}
	if err := c.controller.Serve(c.reportChan, c.parallelChan); err != nil {
		return err
	}
	if c.opts.background {
		go c.serve()
		return nil
	}
	return c.serve()
}

func (c *Commander) RequestChan() chan<- *ServiceRequest {
	return c.reqChan
}

func (c *Commander) ResponseChan() <-chan *ServiceResponse {
	return c.respChan
}

func (c *Commander) serve() (err error) {
	defer func() {
		if v := recover(); v != nil {
			err = fmt.Errorf("%v", v)
			if c.opts.errorChan != nil {
				c.opts.errorChan <- err
			}
		}
	}()

	defer close(c.exitChan)
	defer c.scheduler.Close()
	defer c.controller.Close()

	for {
		select {
		case req := <-c.reqChan:
			task := c.createTask(req)
			c.pushTask(task)
		case <-c.dieChan:
			return
		}
	}
}

func (c *Commander) invoke(ctx context.Context, deviceName string, req *ServiceRequest) (resp *ServiceResponse) {
	resp = new(ServiceResponse)
	defer func() {
		resp.Request = req
	}()

	defer func() {
		if v := recover(); v != nil {
			err := fmt.Errorf("%v", v)
			resp.Err = err
		}
	}()

	e := encoding.NewJSON()
	data, err := e.Marshal(req)
	if err != nil {
		resp.Err = err
		return resp
	}

	routeStyle := magic.NewChainStyle(magic.SeparatorSlash, magic.SeparatorUnderscore)
	r := route.NewChainRoute(device.Addr(c), routeStyle.Chain("/http"))
	msg := &message.Message{
		Route:    r,
		Encoding: e,
		Data:     data,
	}
	processor := device.NewFuncProcessor(func(ctx context.Context, msg *message.Message) error {
		return msg.Encoding.Unmarshal(msg.Data, resp)
	})
	if err := c.Client.Invoke(ctx, msg, processor); err != nil {
		resp.Err = err
		return resp
	}
	return resp
}

func (c *Commander) createTask(req *ServiceRequest) *scheduler.Task {
	return scheduler.NewTask(
		scheduler.TaskOption.Name(fmt.Sprintf("%s[%d]", c.String(), c.reqIndex)),
		scheduler.TaskOption.Stage(func(task *scheduler.Task) error {
			c.respChan <- c.invoke(c.opts.context, c.String(), req)
			return nil
		}),
	)
}

func (c *Commander) pushTask(task *scheduler.Task) {
	task.Publish(c.scheduler)
}

type commanderOptions struct {
	name             string
	context          context.Context
	safety           bool
	background       bool
	errorChan        chan<- error
	parallelInit     int
	parallelTick     time.Duration
	parallelIncrease int
	tpsLimit         int
	reqBacklog       int
	respBacklog      int
	reportBacklog    int
	taskBacklog      int
	parallelBacklog  int
}

var defaultCommanderOptions = commanderOptions{
	name:            "",
	context:         context.Background(),
	safety:          false,
	background:      false,
	errorChan:       nil,
	parallelInit:    1,
	tpsLimit:        -1,
	reqBacklog:      0,
	respBacklog:     0,
	reportBacklog:   0,
	taskBacklog:     0,
	parallelBacklog: 0,
}

type ApplyCommanderOption interface {
	apply(*commanderOptions)
}

type funcCommanderOption func(*commanderOptions)

func (fco funcCommanderOption) apply(co *commanderOptions) {
	fco(co)
}

type commanderOption int

var CommanderOption commanderOption

func (commanderOption) Name(name string) funcCommanderOption {
	return func(c *commanderOptions) {
		c.name = name
	}
}

func (c *Commander) WithName(name string) *Commander {
	CommanderOption.Name(name).apply(&c.opts)
	return c
}

func (commanderOption) Context(context context.Context) funcCommanderOption {
	return func(c *commanderOptions) {
		c.context = context
	}
}

func (c *Commander) WithContext(context context.Context) *Commander {
	CommanderOption.Context(context).apply(&c.opts)
	return c
}

func (commanderOption) Safety() funcCommanderOption {
	return func(c *commanderOptions) {
		c.safety = true
	}
}

func (c *Commander) WithSafety() *Commander {
	CommanderOption.Safety().apply(&c.opts)
	return c
}

func (commanderOption) Background() funcCommanderOption {
	return func(c *commanderOptions) {
		c.background = true
	}
}

func (c *Commander) WithBackground(background bool) *Commander {
	CommanderOption.Background().apply(&c.opts)
	return c
}

func (commanderOption) ErrorChan(errorChan chan<- error) funcCommanderOption {
	return func(c *commanderOptions) {
		c.errorChan = errorChan
	}
}

func (c *Commander) WithErrorChan(errorChan chan<- error) *Commander {
	CommanderOption.ErrorChan(errorChan).apply(&c.opts)
	return c
}

func (commanderOption) ParallelInit(parallelInit int) funcCommanderOption {
	return func(c *commanderOptions) {
		if parallelInit > 0 {
			c.parallelInit = parallelInit
		}
	}
}

func (c *Commander) WithParallelInit(parallelInit int) *Commander {
	CommanderOption.ParallelInit(parallelInit).apply(&c.opts)
	return c
}

func (commanderOption) ParallelTick(parallelTick time.Duration) funcCommanderOption {
	return func(c *commanderOptions) {
		c.parallelTick = parallelTick
	}
}

func (c *Commander) WithParallelTick(parallelTick time.Duration) *Commander {
	CommanderOption.ParallelTick(parallelTick).apply(&c.opts)
	return c
}

func (commanderOption) ParallelIncrease(parallelIncrease int) funcCommanderOption {
	return func(c *commanderOptions) {
		if parallelIncrease >= 0 {
			c.parallelIncrease = parallelIncrease
		}
	}
}

func (c *Commander) WithParallelIncrease(parallelIncrease int) *Commander {
	CommanderOption.ParallelIncrease(parallelIncrease).apply(&c.opts)
	return c
}

func (commanderOption) TPSLimit(tpsLimit int) funcCommanderOption {
	return func(c *commanderOptions) {
		if tpsLimit > 0 {
			c.tpsLimit = tpsLimit
		}
	}
}

func (c *Commander) WithTPSLimit(tpsLimit int) *Commander {
	CommanderOption.TPSLimit(tpsLimit).apply(&c.opts)
	return c
}

func (commanderOption) RequestBacklog(reqChanBacklog int) funcCommanderOption {
	return func(c *commanderOptions) {
		if reqChanBacklog > 1 {
			c.reqBacklog = reqChanBacklog
		}
	}
}

func (c *Commander) WithRequestBacklog(reqBacklog int) *Commander {
	CommanderOption.RequestBacklog(reqBacklog).apply(&c.opts)
	return c
}

func (commanderOption) ResponseBacklog(respBacklog int) funcCommanderOption {
	return func(c *commanderOptions) {
		if respBacklog > 1 {
			c.respBacklog = respBacklog
		}
	}
}

func (c *Commander) WithResponseBacklog(respBacklog int) *Commander {
	CommanderOption.ResponseBacklog(respBacklog).apply(&c.opts)
	return c
}

func (commanderOption) ReportBacklog(reportBacklog int) funcCommanderOption {
	return func(c *commanderOptions) {
		if reportBacklog > 1 {
			c.reportBacklog = reportBacklog
		}
	}
}

func (c *Commander) WithReportBacklog(reportBacklog int) *Commander {
	CommanderOption.ReportBacklog(reportBacklog).apply(&c.opts)
	return c
}

func (commanderOption) TaskBacklog(taskBacklog int) funcCommanderOption {
	return func(c *commanderOptions) {
		if taskBacklog > 1 {
			c.taskBacklog = taskBacklog
		}
	}
}

func (c *Commander) WithTaskBacklog(taskBacklog int) *Commander {
	CommanderOption.TaskBacklog(taskBacklog).apply(&c.opts)
	return c
}

func (commanderOption) ParallelBacklog(parallelBacklog int) funcCommanderOption {
	return func(c *commanderOptions) {
		if parallelBacklog > 1 {
			c.parallelBacklog = parallelBacklog
		}
	}
}

func (c *Commander) WithParallelBacklog(parallelBacklog int) *Commander {
	CommanderOption.ParallelBacklog(parallelBacklog).apply(&c.opts)
	return c
}
