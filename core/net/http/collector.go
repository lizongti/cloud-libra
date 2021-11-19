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

type Collector struct {
	*device.Client
	opts         collectorOptions
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

func NewCollector(opt ...ApplyCollectorOption) *Collector {
	opts := defaultCollectorOptions
	for _, o := range opt {
		o.apply(&opts)
	}

	return &Collector{
		Client:    device.NewClient(),
		opts:      opts,
		reqIndex:  0,
		errorChan: make(chan error),
		dieChan:   make(chan struct{}),
		exitChan:  make(chan struct{}),
	}
}

func (c *Collector) String() string {
	return c.opts.name
}

func (c *Collector) Close() {
	close(c.dieChan)
	<-c.exitChan
}

func (c *Collector) Serve() error {
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

func (c *Collector) RequestChan() chan<- *ServiceRequest {
	return c.reqChan
}

func (c *Collector) ResponseChan() <-chan *ServiceResponse {
	return c.respChan
}

func (c *Collector) serve() (err error) {
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

func (c *Collector) invoke(ctx context.Context, deviceName string, req *ServiceRequest) (resp *ServiceResponse) {
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

func (c *Collector) createTask(req *ServiceRequest) *scheduler.Task {
	return scheduler.NewTask(
		scheduler.TaskOption.Name(fmt.Sprintf("%s[%d]", c.String(), c.reqIndex)),
		scheduler.TaskOption.Stage(func(task *scheduler.Task) error {
			c.respChan <- c.invoke(c.opts.context, c.String(), req)
			return nil
		}),
	)
}

func (c *Collector) pushTask(task *scheduler.Task) {
	task.Publish(c.scheduler)
}

type collectorOptions struct {
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

var defaultCollectorOptions = collectorOptions{
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

type ApplyCollectorOption interface {
	apply(*collectorOptions)
}

type funcCollectorOption func(*collectorOptions)

func (fco funcCollectorOption) apply(co *collectorOptions) {
	fco(co)
}

type collectorOption int

var CollectorOption collectorOption

func (collectorOption) Name(name string) funcCollectorOption {
	return func(c *collectorOptions) {
		c.name = name
	}
}

func (c *Collector) WithName(name string) *Collector {
	CollectorOption.Name(name).apply(&c.opts)
	return c
}

func (collectorOption) Context(context context.Context) funcCollectorOption {
	return func(c *collectorOptions) {
		c.context = context
	}
}

func (c *Collector) WithContext(context context.Context) *Collector {
	CollectorOption.Context(context).apply(&c.opts)
	return c
}

func (collectorOption) Safety() funcCollectorOption {
	return func(c *collectorOptions) {
		c.safety = true
	}
}

func (c *Collector) WithSafety() *Collector {
	CollectorOption.Safety().apply(&c.opts)
	return c
}

func (collectorOption) Background() funcCollectorOption {
	return func(c *collectorOptions) {
		c.background = true
	}
}

func (c *Collector) WithBackground(background bool) *Collector {
	CollectorOption.Background().apply(&c.opts)
	return c
}

func (collectorOption) ErrorChan(errorChan chan<- error) funcCollectorOption {
	return func(c *collectorOptions) {
		c.errorChan = errorChan
	}
}

func (c *Collector) WithErrorChan(errorChan chan<- error) *Collector {
	CollectorOption.ErrorChan(errorChan).apply(&c.opts)
	return c
}

func (collectorOption) ParallelInit(parallelInit int) funcCollectorOption {
	return func(c *collectorOptions) {
		if parallelInit > 0 {
			c.parallelInit = parallelInit
		}
	}
}

func (c *Collector) WithParallelInit(parallelInit int) *Collector {
	CollectorOption.ParallelInit(parallelInit).apply(&c.opts)
	return c
}

func (collectorOption) ParallelTick(parallelTick time.Duration) funcCollectorOption {
	return func(c *collectorOptions) {
		c.parallelTick = parallelTick
	}
}

func (c *Collector) WithParallelTick(parallelTick time.Duration) *Collector {
	CollectorOption.ParallelTick(parallelTick).apply(&c.opts)
	return c
}

func (collectorOption) ParallelIncrease(parallelIncrease int) funcCollectorOption {
	return func(c *collectorOptions) {
		if parallelIncrease >= 0 {
			c.parallelIncrease = parallelIncrease
		}
	}
}

func (c *Collector) WithParallelIncrease(parallelIncrease int) *Collector {
	CollectorOption.ParallelIncrease(parallelIncrease).apply(&c.opts)
	return c
}

func (collectorOption) TPSLimit(tpsLimit int) funcCollectorOption {
	return func(c *collectorOptions) {
		if tpsLimit > 0 {
			c.tpsLimit = tpsLimit
		}
	}
}

func (c *Collector) WithTPSLimit(tpsLimit int) *Collector {
	CollectorOption.TPSLimit(tpsLimit).apply(&c.opts)
	return c
}

func (collectorOption) RequestBacklog(reqChanBacklog int) funcCollectorOption {
	return func(c *collectorOptions) {
		if reqChanBacklog > 1 {
			c.reqBacklog = reqChanBacklog
		}
	}
}

func (c *Collector) WithRequestBacklog(reqBacklog int) *Collector {
	CollectorOption.RequestBacklog(reqBacklog).apply(&c.opts)
	return c
}

func (collectorOption) ResponseBacklog(respBacklog int) funcCollectorOption {
	return func(c *collectorOptions) {
		if respBacklog > 1 {
			c.respBacklog = respBacklog
		}
	}
}

func (c *Collector) WithResponseBacklog(respBacklog int) *Collector {
	CollectorOption.ResponseBacklog(respBacklog).apply(&c.opts)
	return c
}

func (collectorOption) ReportBacklog(reportBacklog int) funcCollectorOption {
	return func(c *collectorOptions) {
		if reportBacklog > 1 {
			c.reportBacklog = reportBacklog
		}
	}
}

func (c *Collector) WithReportBacklog(reportBacklog int) *Collector {
	CollectorOption.ReportBacklog(reportBacklog).apply(&c.opts)
	return c
}

func (collectorOption) TaskBacklog(taskBacklog int) funcCollectorOption {
	return func(c *collectorOptions) {
		if taskBacklog > 1 {
			c.taskBacklog = taskBacklog
		}
	}
}

func (c *Collector) WithTaskBacklog(taskBacklog int) *Collector {
	CollectorOption.TaskBacklog(taskBacklog).apply(&c.opts)
	return c
}

func (collectorOption) ParallelBacklog(parallelBacklog int) funcCollectorOption {
	return func(c *collectorOptions) {
		if parallelBacklog > 1 {
			c.parallelBacklog = parallelBacklog
		}
	}
}

func (c *Collector) WithParallelBacklog(parallelBacklog int) *Collector {
	CollectorOption.ParallelBacklog(parallelBacklog).apply(&c.opts)
	return c
}
