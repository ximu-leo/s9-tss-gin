package opio

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

var DefaultInterruptSignals = []os.Signal{
	os.Interrupt,
	os.Kill,
	syscall.SIGTERM,
	syscall.SIGQUIT,
}

func BlockOnInterrupts(signals ...os.Signal) {
	if len(signals) == 0 {
		signals = DefaultInterruptSignals
	}
	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, signals...)
	<-interruptChannel
}

func BlockOnInterruptsContext(ctx context.Context, signals ...os.Signal) {
	if len(signals) == 0 {
		signals = DefaultInterruptSignals
	}
	interruptChannel := make(chan os.Signal, 1)
	// signal.Notify 不是“打包”，而是“订阅”
	// 把四个其中之一的信号 放入interruptChannel 通道中
	signal.Notify(interruptChannel, signals...)
	select {
	// 通道中有值，唤醒并退出，不管有没有执行操作，都退出，这里没什么操作
	case <-interruptChannel:
	// 上级通知结束程序，唤醒并退出
	case <-ctx.Done():
		signal.Stop(interruptChannel)
	}
}

type interruptContextKeyType struct{}

var blockerContextKey = interruptContextKeyType{}

type interruptCatcher struct {
	incoming chan os.Signal
}

func (c *interruptCatcher) Block(ctx context.Context) {
	select {
	case <-c.incoming:
	case <-ctx.Done():
	}
}

func WithInterruptBlocker(ctx context.Context) context.Context {
	if ctx.Value(blockerContextKey) != nil { // already has an interrupt handler
		return ctx
	}
	catcher := &interruptCatcher{
		incoming: make(chan os.Signal, 10),
	}
	signal.Notify(catcher.incoming, DefaultInterruptSignals...)

	return context.WithValue(ctx, blockerContextKey, BlockFn(catcher.Block))
}

func WithBlocker(ctx context.Context, fn BlockFn) context.Context {
	return context.WithValue(ctx, blockerContextKey, fn)
}

type BlockFn func(ctx context.Context)

func BlockerFromContext(ctx context.Context) BlockFn {
	v := ctx.Value(blockerContextKey)
	if v == nil {
		return nil
	}
	return v.(BlockFn)
}

func CancelOnInterrupt(ctx context.Context) context.Context {
	inner, cancel := context.WithCancel(ctx)

	blockOnInterrupt := BlockerFromContext(ctx)
	if blockOnInterrupt == nil {
		blockOnInterrupt = func(ctx context.Context) {
			BlockOnInterruptsContext(ctx) // default signals
		}
	}

	go func() {
		blockOnInterrupt(ctx)
		cancel()
	}()

	return inner
}
