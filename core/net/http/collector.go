package http

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/aceaura/libra/core/device"
	"github.com/aceaura/libra/core/encoding"
	"github.com/aceaura/libra/core/magic"
	"github.com/aceaura/libra/core/message"
	"github.com/aceaura/libra/core/route"
	"github.com/aceaura/libra/core/scheduler"
	"github.com/mohae/deepcopy"
)

var (
	ErrRequestNotFound = errors.New("request not found by name")
)

type Collector struct {
	*device.Client
	opts         collectorOptions
	scheduler    *scheduler.Scheduler
	reqMap       map[string]*ServiceRequest
	stateMap     map[scheduler.TaskStateType]int
	reqChan      chan *ServiceRequest
	respChan     chan *ServiceResponse
	reportChan   chan *scheduler.Report
	parallelChan chan int
	dieChan      chan struct{}
	exitChan     chan struct{}
	tpsMax       int
	tpsFinished  int
	reqIndex     int
}

func NewCollector(opt ...ApplyCollectorOption) *Collector {
	opts := defaultCollectorOptions
	for _, o := range opt {
		o.apply(&opts)
	}

	reqChan := make(chan *ServiceRequest, opts.reqBacklog)
	respChan := make(chan *ServiceResponse, opts.respBacklog)
	reportChan := make(chan *scheduler.Report, opts.reportBacklog)
	parallelChan := make(chan int)
	s := scheduler.New(
		scheduler.SchedulerOption.Backlog(opts.reqBacklog),
		scheduler.SchedulerOption.Parallel(opts.parallelInit),
		scheduler.SchedulerOption.ReportChan(reportChan),
		scheduler.SchedulerOption.ParallelChan(parallelChan),
		scheduler.SchedulerOption.Background(),
		scheduler.SchedulerOption.Safety(),
	)

	return &Collector{
		Client:       device.NewClient(),
		opts:         opts,
		scheduler:    s,
		reqMap:       make(map[string]*ServiceRequest),
		reqChan:      reqChan,
		respChan:     respChan,
		stateMap:     make(map[scheduler.TaskStateType]int),
		reportChan:   reportChan,
		parallelChan: parallelChan,
		dieChan:      make(chan struct{}),
		exitChan:     make(chan struct{}),
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
	if err := c.scheduler.Serve(); err != nil {
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
	if c.opts.safety {
		defer func() {
			if e := recover(); e != nil {
				err = fmt.Errorf("%v", e)
			}
		}()
	}

	defer func() {
		c.scheduler.Close()
		close(c.exitChan)
	}()

	var (
		dying      = false
		tickerChan = time.NewTicker(c.opts.parallelTick).C
	)

	for {
		select {
		case req := <-c.reqChan:
			name := c.addRequest(req)
			c.scheduleRequest(name)

		case r := <-c.reportChan:
			c.updateStateMap(r)
			c.updateTPSFinished(r)

			switch r.State {
			case scheduler.TaskStateDone:
				c.removeRequest(r.Name)
			case scheduler.TaskStateFailed:
				c.scheduleRequest(r.Name)
			}

			if dying && len(c.reqMap) == 0 {
				return
			}

		case <-tickerChan:
			c.updateTPSMax()
			c.updateParallel()

		case <-c.dieChan:
			dying = true
		}
	}
}

func (c *Collector) updateStateMap(r *scheduler.Report) {
	if r.Progress == 0 {
		c.stateMap[r.State]++
	}
}

func (c *Collector) updateTPSFinished(r *scheduler.Report) {
	if r.State == scheduler.TaskStateDone || r.State == scheduler.TaskStateFailed {
		c.tpsFinished++
	}
}

func (c *Collector) updateTPSMax() {
	tps := int(float64(c.tpsFinished) * float64(c.opts.parallelTick) / float64(time.Second))
	c.tpsMax = int(math.Max(float64(tps), float64(c.tpsMax)))
	c.tpsFinished = 0
}

func (c *Collector) updateParallel() {
	if c.stateMap[scheduler.TaskStatePending]-c.stateMap[scheduler.TaskStateRunning] > 0 {
		if c.tpsMax < c.opts.tpsLimit || c.opts.tpsLimit < 0 {
			c.parallelChan <- c.opts.parallelIncrease
		}
	}
}

func (c *Collector) invoke(ctx context.Context, req *ServiceRequest, deviceName string) (*ServiceResponse, error) {
	e := encoding.Empty()
	data, err := e.Marshal(req)
	if err != nil {
		return nil, err
	}

	routeStyle := magic.NewChainStyle(magic.SeparatorSlash, magic.SeparatorUnderscore)

	r := route.NewChainRoute(device.Addr(c), routeStyle.Chain("/http"))

	msg := &message.Message{
		Route:    r,
		Encoding: e,
		Data:     data,
	}
	resp := new(ServiceResponse)
	processor := device.NewFuncProcessor(func(ctx context.Context, msg *message.Message) error {
		return msg.Encoding.Unmarshal(msg.Data, resp)
	})
	if err := c.Client.Invoke(ctx, msg, processor); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Collector) addRequest(req *ServiceRequest) string {
	name := fmt.Sprintf("%s[%d]", c.String(), c.reqIndex)
	c.reqIndex++
	c.reqMap[name] = req
	return name
}

func (c *Collector) removeRequest(name string) {
	delete(c.reqMap, name)
}

func (c *Collector) copyRequest(name string) *ServiceRequest {
	if req, ok := c.reqMap[name]; ok {
		return deepcopy.Copy(req).(*ServiceRequest)
	}
	return nil
}

func (c *Collector) scheduleRequest(name string) {
	req := c.copyRequest(name)
	stage := func(task *scheduler.Task) error {
		resp, err := c.invoke(c.opts.context, req, c.String())
		if err != nil {
			return err
		}
		c.respChan <- resp
		return nil
	}

	scheduler.NewTask(
		scheduler.TaskOption.Name(name),
		scheduler.TaskOption.Stage(stage),
	).Publish(c.scheduler)
}

type collectorOptions struct {
	name             string
	context          context.Context
	background       bool
	safety           bool
	parallelInit     int
	parallelTick     time.Duration
	parallelIncrease int
	tpsLimit         int
	reqBacklog       int
	respBacklog      int
	reportBacklog    int
}

var defaultCollectorOptions = collectorOptions{
	name:             "",
	context:          context.Background(),
	background:       false,
	safety:           false,
	parallelInit:     1,
	parallelTick:     time.Second,
	parallelIncrease: 1,
	tpsLimit:         -1,
	reqBacklog:       0,
	respBacklog:      0,
	reportBacklog:    0,
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

func (collectorOption) Background() funcCollectorOption {
	return func(c *collectorOptions) {
		c.background = true
	}
}

func (c *Collector) WithBackground(background bool) *Collector {
	CollectorOption.Background().apply(&c.opts)
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

func (collectorOption) ParallelInit(parallelInit int) funcCollectorOption {
	return func(c *collectorOptions) {
		c.parallelInit = parallelInit
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
		c.parallelIncrease = parallelIncrease
	}
}

func (c *Collector) WithParallelIncrease(parallelIncrease int) *Collector {
	CollectorOption.ParallelIncrease(parallelIncrease).apply(&c.opts)
	return c
}

func (collectorOption) TPSLimit(tpsLimit int) funcCollectorOption {
	return func(c *collectorOptions) {
		c.tpsLimit = tpsLimit
	}
}

func (c *Collector) WithTPSLimit(tpsLimit int) *Collector {
	CollectorOption.TPSLimit(tpsLimit).apply(&c.opts)
	return c
}

func (collectorOption) RequestBacklog(reqChanBacklog int) funcCollectorOption {
	return func(c *collectorOptions) {
		c.reqBacklog = reqChanBacklog
	}
}

func (c *Collector) WithRequestBacklog(reqBacklog int) *Collector {
	CollectorOption.RequestBacklog(reqBacklog).apply(&c.opts)
	return c
}

func (collectorOption) ResponseBacklog(respBacklog int) funcCollectorOption {
	return func(c *collectorOptions) {
		c.respBacklog = respBacklog
	}
}

func (c *Collector) WithResponseBacklog(respBacklog int) *Collector {
	CollectorOption.ResponseBacklog(respBacklog).apply(&c.opts)
	return c
}

func (collectorOption) ReportBacklog(reportBacklog int) funcCollectorOption {
	return func(c *collectorOptions) {
		c.reportBacklog = reportBacklog
	}
}

func (c *Collector) WithReportBacklog(reportBacklog int) *Collector {
	CollectorOption.ReportBacklog(reportBacklog).apply(&c.opts)
	return c
}
