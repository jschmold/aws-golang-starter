package http

// TestContext is a basic, empty, no-effect HttpContext for use in unit testing controllers
type TestContext struct {
	req IRequest
	res IResponse
}

// NewTestContext creates an http context object where you can define the request and response objects
func NewTestContext(req IRequest, res IResponse) *TestContext {
	return &TestContext{req: req, res: res}
}

// GetRequest gets the defined IRequest object initialized with NewTestContext
func (ctx *TestContext) GetRequest() IRequest {
	return ctx.req
}

// GetResponse gets the defined IResponse object initialized with NewTestContext
func (ctx *TestContext) GetResponse() IResponse {
	return ctx.res
}
