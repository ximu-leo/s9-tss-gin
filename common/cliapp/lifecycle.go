package cliapp

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/ximu-leo/s9-tss-gin/common/opio"
)

type Lifecycle interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Stopped() bool
}

// 定义一个“函数模板”
type LifecycleAction func(ctx *cli.Context, close context.CancelCauseFunc) (Lifecycle, error)

func LifecycleCmd(fn LifecycleAction) cli.ActionFunc {
	// 乌拉：opio.BlockOnInterruptsContext 这里只是  引用函数，不是调用函数，所以无参数
	return lifecycleCmd(fn, opio.BlockOnInterruptsContext)
}

type waitSignalFn func(ctx context.Context, signals ...os.Signal)

var interruptErr = errors.New("interrupt signal")

func lifecycleCmd(fn LifecycleAction, blockOnInterrupt waitSignalFn) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		hostCtx := ctx.Context
		appCtx, appCancel := context.WithCancelCause(hostCtx)
		ctx.Context = appCtx

		// 启动 1个协程  G1 协程
		go func() {
			// 乌拉：真正调用函数
			blockOnInterrupt(appCtx)
			// 调用 appCancel(interruptErr) 时，Go 内部会做两件事：
			// 关闭 appCtx.Done() 这个通道（给警报器通电）。
			// 记录错误原因。
			appCancel(interruptErr)
		}()
		//实际上调用就是 runGinHttp(ctx, appCancel)
		appLifecycle, err := fn(ctx, appCancel)
		if err != nil {
			return errors.Join(
				fmt.Errorf("failed to setup: %w", err),
				context.Cause(appCtx),
			)
		}

		if err := appLifecycle.Start(appCtx); err != nil {
			return errors.Join(
				fmt.Errorf("failed to start: %w", err),
				context.Cause(appCtx),
			)
		}

		// appCtx.Done() 是 Context 自带的“关门通道”，它被 appCancel 函数控制
		// 觉察到 关闭 appCtx.Done()，继续往下
		<-appCtx.Done()

		stopCtx, stopCancel := context.WithCancelCause(hostCtx)
		go func() {
			blockOnInterrupt(stopCtx)
			stopCancel(interruptErr)
		}()

		stopErr := appLifecycle.Stop(stopCtx)
		stopCancel(nil)
		if stopErr != nil {
			return errors.Join(
				fmt.Errorf("failed to stop: %w", stopErr),
				context.Cause(stopCtx),
			)
		}
		return nil
	}
}
