package http

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cloudlibraries/libra/internal/boost/coroutine"
	"github.com/cloudlibraries/libra/internal/boost/magic"
	"github.com/cloudlibraries/libra/internal/core/device"
	"github.com/cloudlibraries/libra/internal/core/encoding"
	"github.com/cloudlibraries/libra/internal/core/message"
	"github.com/cloudlibraries/libra/internal/core/route"
	"github.com/cloudlibraries/libra/internal/core/scheduler"
)

var (
	ErrRequestNotFound               = errors.New("request not found by name")
	ErrCoroutineYieldOutputEmpty     = errors.New("coroutine yield output is empty ")
	ErrCoroutineYieldOutputTypeError = errors.New("coroutine yield output type error")
)

type Commander struct {
	*device.Client
	opts       commanderOptions
	controller *scheduler.TPSController
	reqIndex   int
	reqChan    chan *ServiceRequest
	respChan   chan *ServiceResponse
	errorChan  chan error
	dieChan    chan struct{}
	exitChan   chan struct{}
}

func NewCommander(opt ...ApplyCommanderOption) *Commander {
	opts := defaultCommanderOptions
	for _, o := range opt {
		o.apply(&opts)
	}

	return &Commander{
		Client:    device.NewClient(""),
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
	c.reqChan = make(chan *ServiceRequest, c.opts.reqBacklog)
	c.respChan = make(chan *ServiceResponse, c.opts.respBacklog)

	c.controller = scheduler.NewTPSController(
		scheduler.WithTPSSafety(),
		scheduler.WithTPSBackground(),
		scheduler.WithTPSErrorChan(c.errorChan),
		scheduler.WithTPSParallel(c.opts.parallel),
		scheduler.WithTPSTaskBacklog(c.opts.taskBacklog),
		scheduler.WithTPSReportBacklog(c.opts.reportBacklog),
		scheduler.WithTPSParallelBacklog(c.opts.parallelBacklog),
		scheduler.WithTPSParallelTick(c.opts.parallelTick),
		scheduler.WithTPSParallelIncrease(c.opts.parallelIncrease),
		scheduler.WithTPSLimit(c.opts.tpsLimit),
	)

	if err := c.controller.Serve(); err != nil {
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

func Invoke(c *Commander, req *ServiceRequest) *ServiceResponse {
	return c.Invoke(req)
}

func (c *Commander) Invoke(req *ServiceRequest) *ServiceResponse {
	var resp *ServiceResponse
	if err := coroutine.Start(c.opts.ctx, func(co *coroutine.Coroutine) error {
		req.CoroutineID = co.ID()
		c.RequestChan() <- req
		output, err := co.Yield()
		if err != nil {
			return err
		}
		if len(output) == 0 {
			return ErrCoroutineYieldOutputEmpty
		}
		var ok bool
		resp, ok = output[0].(*ServiceResponse)
		if !ok {
			return ErrCoroutineYieldOutputTypeError
		}
		return nil
	}); err != nil {
		return &ServiceResponse{
			Err: err,
		}
	}

	return resp
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

	r := route.NewChainRoute(device.Addr(c), magic.GoogleChain("/http"))
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
		scheduler.WithTaskName(fmt.Sprintf("%s[%d]", c.String(), c.reqIndex)),
		scheduler.WithTaskStage(func(task *scheduler.Task) error {
			if req.CoroutineID != "" {
				coroutine.TryResume(req.CoroutineID, c.invoke(c.opts.ctx, c.String(), req))
			} else {
				c.respChan <- c.invoke(c.opts.ctx, c.String(), req)
			}

			return nil
		}),
	)
}

func (c *Commander) pushTask(task *scheduler.Task) {
	task.Publish(c.controller.Scheduler())
}

type commanderOptions struct {
	name             string
	ctx              context.Context
	safety           bool
	background       bool
	errorChan        chan<- error
	parallel         int
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
	ctx:             context.Background(),
	safety:          false,
	background:      false,
	errorChan:       nil,
	parallel:        1,
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

func (f funcCommanderOption) apply(opt *commanderOptions) {
	f(opt)
}

func WithCommanderName(name string) funcCommanderOption {
	return func(c *commanderOptions) {
		c.name = name
	}
}

func WithCommanderContext(ctx context.Context) funcCommanderOption {
	return func(c *commanderOptions) {
		c.ctx = ctx
	}
}

func WithCommanderSafety() funcCommanderOption {
	return func(c *commanderOptions) {
		c.safety = true
	}
}

func WithCommanderBackground() funcCommanderOption {
	return func(c *commanderOptions) {
		c.background = true
	}
}

func WithCommanderErrorChan(errorChan chan<- error) funcCommanderOption {
	return func(c *commanderOptions) {
		c.errorChan = errorChan
	}
}

func WithCommanderParallel(parallel int) funcCommanderOption {
	return func(c *commanderOptions) {
		if parallel > 0 {
			c.parallel = parallel
		}
	}
}

func WithCommanderParallelTick(parallelTick time.Duration) funcCommanderOption {
	return func(c *commanderOptions) {
		c.parallelTick = parallelTick
	}
}

func WithCommanderParallelIncrease(parallelIncrease int) funcCommanderOption {
	return func(c *commanderOptions) {
		if parallelIncrease >= 0 {
			c.parallelIncrease = parallelIncrease
		}
	}
}

func WithCommanderTPSLimit(tpsLimit int) funcCommanderOption {
	return func(c *commanderOptions) {
		if tpsLimit > 0 {
			c.tpsLimit = tpsLimit
		}
	}
}

func WithCommanderRequestBacklog(reqChanBacklog int) funcCommanderOption {
	return func(c *commanderOptions) {
		if reqChanBacklog > 1 {
			c.reqBacklog = reqChanBacklog
		}
	}
}

func WithCommanderResponseBacklog(respBacklog int) funcCommanderOption {
	return func(c *commanderOptions) {
		if respBacklog > 1 {
			c.respBacklog = respBacklog
		}
	}
}

func WithCommanderReportBacklog(reportBacklog int) funcCommanderOption {
	return func(c *commanderOptions) {
		if reportBacklog > 1 {
			c.reportBacklog = reportBacklog
		}
	}
}

func WithCommanderTaskBacklog(taskBacklog int) funcCommanderOption {
	return func(c *commanderOptions) {
		if taskBacklog > 1 {
			c.taskBacklog = taskBacklog
		}
	}
}

func WithCommanderParallelBacklog(parallelBacklog int) funcCommanderOption {
	return func(c *commanderOptions) {
		if parallelBacklog > 1 {
			c.parallelBacklog = parallelBacklog
		}
	}
}
