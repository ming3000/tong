package common

import "time"

// Job is an interface for job to do.
type Job interface {
	// true to decline the step period, false to stay
	Run() bool
}

/// Cron keeps track of a job,
// invoking associated job's Run() method.
type Cron struct {
	period        time.Duration
	initialPeriod time.Duration
	stepPeriod    time.Duration
	maxPeriod     time.Duration
	prev          time.Time
	next          time.Time
	job           Job
	wait          chan bool
	stop          chan struct{}
	running       bool
}

// SetLogger set the logger field of Cron.
func (c *Cron) SetLogger(logger *Logger) *Cron {
	return c
}

// Every schedules a new Cron and returns it.
func (c *Cron) Every(period time.Duration) *Cron {
	c.schedule(period)
	return c
}

// Set the job.
func (c *Cron) Do(job Job) *Cron {
	c.job = job
	return c
}

// Start the Cron in its own go-routine,
// or do nothing if already started.
func (c *Cron) Start() {
	if c.running {
		return
	}
	c.running = true
	go c.run()
}

// Stop the CronJob if it is running.
func (c *Cron) Stop() {
	if !c.running {
		return
	}
	c.stop <- struct{}{}
	c.running = false
}

func (c *Cron) schedule(period time.Duration) {
	c.period = period
	c.prev = time.Now()
	c.next = c.prev.Add(period)
}

func (c *Cron) run() {
	ticker := time.NewTicker(1 * time.Second)
	var ifWait bool
	go func() {
		for {
			select {
			case <-ticker.C:
				if !ifWait {
					go c.runJob()
				}
				continue
			case <-c.stop:
				ticker.Stop()
				return
			case ifWait = <-c.wait:
				//public.Debug.Println("get msg from the job: ", ifWait, time.Now().Format(time.RFC3339))
				continue
			}
		}
	}()
}

func (c *Cron) runJob() {
	if time.Now().After(c.next) {
		c.wait <- true
		// set-schedule before this job is done
		c.schedule(c.period)
		ifDecline := c.job.Run()
		if ifDecline {
			next := c.period + c.stepPeriod
			if next >= c.maxPeriod {
				next = c.maxPeriod
			}
			// re-schedule after this job is done
			c.schedule(next)
		} else {
			// re-schedule after this job is done
			c.schedule(c.initialPeriod)
		}
		c.wait <- false
	}
}

// returns a new CronJob job runner.
func NewCron(initialPeriod, stepPeriod, maxPeriod time.Duration) *Cron {
	cronObj := &Cron{
		initialPeriod: initialPeriod,
		stepPeriod:    stepPeriod,
		maxPeriod:     maxPeriod,
		job:           nil,
		wait:          make(chan bool),
		stop:          make(chan struct{}),
		running:       false,
	}
	cronObj.schedule(initialPeriod)
	return cronObj
}
