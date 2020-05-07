package shell

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/jerson/deployer/pkg/deployer"

	"github.com/sirupsen/logrus"
)

// WaitIsTrue ...
func WaitIsTrue(ctx context.Context, callback func(ctx context.Context) (ready bool, err error)) (running bool, err error) {

	timeout := GetContextTimeout(ctx, time.Hour)
	expiredTime := time.Now().Add(timeout)
	log := deployer.Log(ctx)
	times := 0
	for {
		diff := expiredTime.Sub(time.Now())
		ready, err := callback(ctx)
		if ready {
			if times > 0 {
				log.Info("is ok but, wait please...")
				time.Sleep(time.Second * 10)
				return WaitIsTrue(
					context.WithValue(ctx, deployer.ContextKeyTimeout, diff),
					callback,
				)
			}
			return ready, err
		}
		time.Sleep(time.Second * 1)
		times++

		if diff < time.Second {
			return false, errors.New("wait timeout")
		}
	}
}

// Run ...
func Run(ctx context.Context, name string, commands ...string) (output string, err error) {

	return RunInDir(ctx, "", name, commands...)
}

// GetContextTimeout ...
func GetContextTimeout(ctx context.Context, defaultValue time.Duration) time.Duration {
	value := ctx.Value(deployer.ContextKeyTimeout)
	if value == nil {
		return defaultValue
	}
	return value.(time.Duration)
}

// RunInDir ...
func RunInDir(ctx context.Context, dir, name string, commands ...string) (output string, err error) {

	log := deployer.Log(ctx)

	disableSTDOUT := false
	if ctx.Value(deployer.ContextKeyDisableSTDOUT) != nil {
		disableSTDOUT = ctx.Value(deployer.ContextKeyDisableSTDOUT).(bool)
	}

	if dir != "" && !disableSTDOUT {
		log.Debug("Running on dir: ", dir)
	}

	if !disableSTDOUT {
		log.Debugf("%s %s", name, strings.Join(commands, " "))
	}
	cmd := exec.Command(name, commands...)
	cmd.Dir = dir
	cmd.Env = os.Environ()

	var stdoutBuf, stderrBuf bytes.Buffer
	if !disableSTDOUT && (log.Logger.Level >= logrus.DebugLevel) {
		cmd.Stdout = io.MultiWriter(log.Logger.Out, &stdoutBuf)
		cmd.Stderr = io.MultiWriter(log.Logger.Out, &stderrBuf)
	} else {
		cmd.Stdout = &stdoutBuf
		cmd.Stderr = &stderrBuf
	}

	errorCode := 0
	err = cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus := exitError.Sys().(syscall.WaitStatus)
			errorCode = waitStatus.ExitStatus()
			if !disableSTDOUT {
				defer log.Error("error code: ", errorCode)
			}
		} else {
			return "", err
		}
	}

	if errorCode != 0 {
		errorString := strings.TrimSpace(string(stderrBuf.Bytes()))
		return "", errors.New(errorString)
	}

	return strings.TrimSpace(string(stdoutBuf.Bytes())), nil
}
