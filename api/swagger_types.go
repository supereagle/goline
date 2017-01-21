package api

// A PipelineName parameter model.
//
// This is used for operations that want the name of a pipeline in the path
// swagger:parameters updatePipeline deletePipeline performPipeline
type PipelineName struct {
	// The name of the pipeline
	//
	// in: path
	// required: true
	Name string `json:"name"`
}

// A PipelineParams parameter model.
//
// This is used for operations that want the pipeline config in the body
// swagger:parameters createPipeline
type PipelineParams struct {
	// The config of the pipeline
	//
	// in: body
	// required: true
	Pipeline *Pipeline `json:"pipeline"`
}

// A PipelinePerformParams parameter model.
//
// This is used for operations that want the pipeline perform parameters in the body
// swagger:parameters performPipeline
type PipelinePerformParams struct {
	// The perform parameters of the pipeline
	//
	// in: body
	// required: true
	PerformParams *PerformParams `json:"performParams"`
}

// A PipelineResponse response model
//
// This is used for returning a response with a single pipeline as body
//
// swagger:response pipelineResponse
type PipelineResponse struct {
	// in: body
	Body struct {
		Code       int32  `json:"code"`
		Status     string `json:"status"`
		JsonObject struct {
			Pipeline *Pipeline `json:"pipeline"`
		} `json:"json_object"`
	} `json:"body"`
}

// A NoObjectResponse response model
//
// This is used for returning a response without json object as body
//
// swagger:response noObjectResponse
type NoObjectResponse struct {
	// in: body
	Body struct {
		Code   int32  `json:"code"`
		Status string `json:"status"`
	} `json:"body"`
}

// A GenericErrorResponse is the default error message that is generated.
// For certain status codes there are more appropriate error structures.
//
// swagger:response genericErrorResponse
type GenericErrorResponse struct {
	// in: body
	Body struct {
		Code     int32  `json:"code"`
		Status   string `json:"status"`
		ErrorMsg string `json:"error"`
	} `json:"body"`
}
