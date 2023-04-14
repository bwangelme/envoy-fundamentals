package main

import (
	"bufio"
	"bytes"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
	"strings"
)

type vmContext struct {
	types.DefaultVMContext
}

func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &pluginContext{
		contextID:         contextID,
		additionalHeaders: map[string]string{},
		// envoy 输出的 prom 指标名称就叫做 hello_header_counter
		helloHeaderCounter: proxywasm.DefineCounterMetric("hello_header_counter"),
	}
}

type pluginContext struct {
	// Embed the default plugin context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultPluginContext
	additionalHeaders  map[string]string
	contextID          uint32
	helloHeaderCounter proxywasm.MetricCounter
}

func (ctx *pluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	config, err := proxywasm.GetPluginConfiguration()
	if err != nil && err != types.ErrorStatusNotFound {
		proxywasm.LogCriticalf("faild to load config %v", err)
		return types.OnPluginStartStatusFailed
	}

	scanner := bufio.NewScanner(bytes.NewReader(config))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}

		tokens := strings.Split(scanner.Text(), "=")
		if len(tokens) == 2 {
			ctx.additionalHeaders[tokens[0]] = tokens[1]
		}
	}

	return types.OnPluginStartStatusOK
}

// 覆盖默认的 types.DefaultPluginContext.
func (ctx *pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &httpHeaders{
		contextID:          contextID,
		additionalHeaders:  ctx.additionalHeaders,
		helloHeaderCounter: ctx.helloHeaderCounter,
	}
}

type httpHeaders struct {
	// Embed the default http context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultHttpContext
	contextID          uint32
	additionalHeaders  map[string]string
	helloHeaderCounter proxywasm.MetricCounter
}

func (ctx *httpHeaders) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	proxywasm.LogInfo("OnHttpRequestHeaders")
	_, err := proxywasm.GetHttpRequestHeader("hello")
	if err != nil {
		// 没有设置 hello 头的花，这忽略它
		return types.ActionContinue
	}

	ctx.helloHeaderCounter.Increment(1)
	proxywasm.LogInfof("hello_header_counter increment")
	return types.ActionContinue
}

func (ctx *httpHeaders) OnHttpResponseHeaders(numHeaders int, endOfStream bool) types.Action {
	proxywasm.LogInfo("OnHttpResponseHeaders")

	for key, value := range ctx.additionalHeaders {
		err := proxywasm.AddHttpResponseHeader(key, value)
		if err != nil {
			proxywasm.LogCriticalf(" failed to add response header: %v", err)
			return types.ActionPause
		}
		proxywasm.LogInfof("header set: %s=%s", key, value)
	}
	return types.ActionContinue
}

func (ctx *httpHeaders) OnHttpStreamDone() {
	proxywasm.LogInfof("%d finished", ctx.contextID)
}

func main() {
	proxywasm.SetVMContext(&vmContext{})
}
