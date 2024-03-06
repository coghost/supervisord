package supervisord

import (
	"errors"
	"fmt"
	"os"
	"time"
)

var ErrTailLog = errors.New("tail log error")

func TailLogError(op string) error {
	return fmt.Errorf("TailLogError %w: %s", ErrTailLog, op)
}

func (c *Client) ReadProcessStdoutLog(name string, offset, length int) (string, error) {
	return c.CallAsStr(readProcessStdoutLog, name, offset, length)
}

func (c *Client) ReadProcessStderrLog(name string, offset, length int) (string, error) {
	return c.CallAsStr(readProcessStderrLog, name, offset, length)
}

func (c *Client) TailProcessStdoutLog(name string, offset int64, length int) (string, int64, bool, error) {
	replies, err := c.CallAsInterfaceArray(tailProcessStdoutLog, name, offset, length)
	if err != nil {
		return "", 0, false, fmt.Errorf("cannot get stdoutlog: %w", err)
	}

	return c.fmtTail(replies)
}

func (c *Client) fmtTail(replies []interface{}) (string, int64, bool, error) {
	raw, success := replies[0].(string)
	if !success {
		return "", 0, false, nil
	}

	offsetNew, success := replies[1].(int64)
	if !success {
		return "", 0, false, TailLogError("cannot get offset")
	}

	b, success := replies[2].(bool)
	if !success {
		return "", 0, false, TailLogError("cannot get overflow")
	}

	return raw, offsetNew, b, nil
}

func (c *Client) TailProcessStderrLog(name string, offset int64, length int) (string, int64, bool, error) {
	replies, err := c.CallAsInterfaceArray(tailProcessStderrLog, name, offset, length)
	if err != nil {
		return "", 0, false, fmt.Errorf("cannot get process Stderrlog: %w", err)
	}

	return c.fmtTail(replies)
}

func (c *Client) ClearProcessLogs(name string) error {
	return c.CallAsBool(clearProcessLogs, name)
}

func (c *Client) ClearAllProcessLogs() ([]interface{}, error) {
	return c.CallAsInterfaceArray(clearAllProcessLogs)
}

type tailFn func(string, int64, int) (string, int64, bool, error)

// TailFProcessLog
//
//	@param bufferSize: if 0 is passed in, will use default 5120
//	@param processFn: if nil is passed in, will use default TailProcessStdoutLog
func (c *Client) TailFProcessLog(name string, bufferSize int, processFn tailFn) {
	var err error

	const readIntervalMs = 100

	if bufferSize == 0 {
		bufferSize = 5120
	}

	if processFn == nil {
		processFn = c.TailProcessStdoutLog
	}

	got, offset, overflow := "", int64(0), false
	lastOffset := offset

	for {
		got, offset, overflow, err = processFn(name, lastOffset, bufferSize)
		if err != nil {
			panic(err)
		}

		if offset > lastOffset {
			newChars := 0
			if lastOffset != 0 {
				newChars = int(offset) - int(lastOffset)
			}

			lastOffset = offset

			// when overflow is true, means the total log size is greate than got, so print all got logs.
			if overflow {
				fmt.Fprint(os.Stdout, got)
				continue
			}

			if newChars > 0 {
				buff := got[len(got)-newChars:]
				if buff != "" {
					fmt.Fprint(os.Stdout, buff)
				}
			}
		} else {
			time.Sleep(time.Millisecond * time.Duration(readIntervalMs))
		}
	}
}
