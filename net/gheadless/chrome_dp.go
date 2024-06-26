package gheadless

// This module has been tested successfully.

import (
	"context"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
	"github.com/davidforest123/goutil/basic/gerrors"
	"github.com/davidforest123/goutil/basic/glog"
	"github.com/davidforest123/goutil/container/gstring"
	"io/ioutil"
	"strings"
	"time"
)

type (
	ChromeDP struct {
		userDataDir string
	}
)

var (
	TaskNavigate   = "navigate"
	TaskScreenshot = "screenshot"
	TaskFullHtml   = "full-html"
)

func NewChromeDP() *ChromeDP {
	return &ChromeDP{}
}

func bufToFile(buf []byte, pathToSave string) error {
	return ioutil.WriteFile(pathToSave, buf, 0o644)
}

// SetUserDataDir is the command line option to set the user data dir,
// all the tasks will share the cookies.
//
// Note: set this option to manually set the profile directory used by Chrome.
// When this is not set, then a default path will be created in the /tmp
// directory.
func (ch *ChromeDP) SetUserDataDir(dir string) {
	ch.userDataDir = dir
}

func (ch *ChromeDP) Navigate(urlStr, proxy string, log glog.Interface, timeout time.Duration) ([]byte, error) {
	result, err := ch.DoTask(urlStr, proxy, []string{TaskNavigate}, log, timeout)
	if err != nil {
		return nil, err
	}
	return result[TaskNavigate], nil
}

func (ch *ChromeDP) Screenshot(urlStr, proxy string, log glog.Interface, timeout time.Duration) ([]byte, error) {
	result, err := ch.DoTask(urlStr, proxy, []string{TaskScreenshot}, log, timeout)
	if err != nil {
		return nil, err
	}
	return result[TaskScreenshot], nil
}

func (ch *ChromeDP) GetFullHtml(urlStr, proxy string, log glog.Interface, timeout time.Duration) ([]byte, error) {
	// FIXME
	// when get full html, must exec TaskScreenshot at the same time, otherwise it can't get full html source
	// i don't why, just do it like this.
	result, err := ch.DoTask(urlStr, proxy, []string{TaskScreenshot, TaskFullHtml}, log, timeout)
	if err != nil {
		return nil, err
	}
	return result[TaskFullHtml], nil
}

func (ch *ChromeDP) DoTask(urlStr, proxy string, tasks []string, log glog.Interface, timeout time.Duration) (map[string][]byte, error) {
	// fix urlStr, chromedp doesn't accept URL without "http://" or "https://"
	if !strings.HasPrefix(strings.ToLower(urlStr), "http://") && !strings.HasPrefix(strings.ToLower(urlStr), "https://") {
		urlStr = "http://" + urlStr
	}

	// create context
	opts := chromedp.DefaultExecAllocatorOptions[:]
	if ch.userDataDir != "" {
		opts = append(opts, chromedp.UserDataDir(ch.userDataDir))
	}
	if proxy != "" { // add proxy
		opts = append(opts,
			// 1) specify the proxy server.
			// Note that the username/password is not provided here.
			// Check the link below for the description of the proxy settings:
			// https://www.chromium.org/developers/design-documents/network-settings
			chromedp.ProxyServer(proxy),
			// By default, Chrome will bypass localhost.
			// The test server is bound to localhost, so we should add the
			// following flag to use the proxy for localhost URLs.
			chromedp.Flag("proxy-bypass-list", "<-loopback>"),
		)
	}
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// log the protocol messages to understand how it works.
	debugFunc := func(string, ...interface{}) {}
	if log != nil {
		debugFunc = log.Debgf
	}
	// Note:
	// If remove two lines code below, Screenshot() will report "invalid context",
	// so if input "log" is nil, we give it a fake debug callback,
	// I don't know why.
	ctx, cancel = chromedp.NewContext(ctx, chromedp.WithDebugf(debugFunc))
	defer cancel()

	// capture entire browser viewport, returning png with quality=90
	chDone := make(chan error, 1)
	var resultScreenshot []byte
	resultFullHtml := ""
	var result = map[string][]byte{}
	go func() {
		var actions []chromedp.Action
		if gstring.Contains(tasks, TaskNavigate) {
			actions = append(actions, buildNavigateTask(urlStr))
		}
		if gstring.Contains(tasks, TaskScreenshot) {
			actions = append(actions, buildFullScreenshotTask(urlStr, 90, &resultScreenshot))
		}
		if gstring.Contains(tasks, TaskFullHtml) {
			actions = append(actions, buildFullHtmlSourceTask(&resultFullHtml))
		}
		if err := chromedp.Run(ctx, actions...); err != nil {
			chDone <- err
			return
		}
		close(chDone) // "chan error" returns nil after close() action
	}()

	// wait result
	if timeout > 0 {
		ticker := time.NewTicker(timeout)
		select {
		case <-ticker.C:
			return nil, gerrors.ErrTimeout
		case err := <-chDone:
			result[TaskScreenshot] = resultScreenshot
			result[TaskFullHtml] = []byte(resultFullHtml)
			return result, err
		}
	} else {
		err := <-chDone
		result[TaskScreenshot] = resultScreenshot
		result[TaskFullHtml] = []byte(resultFullHtml)
		return result, err
	}
}

func buildNavigateTask(urlStr string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlStr),
	}
}

// fullScreenshot takes a screenshot of the entire browser viewport.
//
// Note: chromedp.FullScreenshot overrides the device's emulation settings. Use
// device.Reset to reset the emulation and viewport settings.
func buildFullScreenshotTask(urlStr string, quality int, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlStr),
		chromedp.FullScreenshot(res, quality),
	}
}

func buildFullHtmlSourceTask(res *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			node, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}
			*res, err = dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
			return err
		}),
	}
}
