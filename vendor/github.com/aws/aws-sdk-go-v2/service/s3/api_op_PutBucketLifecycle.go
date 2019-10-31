// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awsutil"
	"github.com/aws/aws-sdk-go-v2/private/protocol"
	"github.com/aws/aws-sdk-go-v2/private/protocol/restxml"
)

// Please also see https://docs.aws.amazon.com/goto/WebAPI/s3-2006-03-01/PutBucketLifecycleRequest
type PutBucketLifecycleInput struct {
	_ struct{} `type:"structure" payload:"LifecycleConfiguration"`

	// Bucket is a required field
	Bucket *string `location:"uri" locationName:"Bucket" type:"string" required:"true"`

	LifecycleConfiguration *LifecycleConfiguration `locationName:"LifecycleConfiguration" type:"structure" xmlURI:"http://s3.amazonaws.com/doc/2006-03-01/"`
}

// String returns the string representation
func (s PutBucketLifecycleInput) String() string {
	return awsutil.Prettify(s)
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *PutBucketLifecycleInput) Validate() error {
	invalidParams := aws.ErrInvalidParams{Context: "PutBucketLifecycleInput"}

	if s.Bucket == nil {
		invalidParams.Add(aws.NewErrParamRequired("Bucket"))
	}
	if s.LifecycleConfiguration != nil {
		if err := s.LifecycleConfiguration.Validate(); err != nil {
			invalidParams.AddNested("LifecycleConfiguration", err.(aws.ErrInvalidParams))
		}
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

func (s *PutBucketLifecycleInput) getBucket() (v string) {
	if s.Bucket == nil {
		return v
	}
	return *s.Bucket
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s PutBucketLifecycleInput) MarshalFields(e protocol.FieldEncoder) error {

	if s.Bucket != nil {
		v := *s.Bucket

		metadata := protocol.Metadata{}
		e.SetValue(protocol.PathTarget, "Bucket", protocol.StringValue(v), metadata)
	}
	if s.LifecycleConfiguration != nil {
		v := s.LifecycleConfiguration

		metadata := protocol.Metadata{XMLNamespaceURI: "http://s3.amazonaws.com/doc/2006-03-01/"}
		e.SetFields(protocol.PayloadTarget, "LifecycleConfiguration", v, metadata)
	}
	return nil
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/s3-2006-03-01/PutBucketLifecycleOutput
type PutBucketLifecycleOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s PutBucketLifecycleOutput) String() string {
	return awsutil.Prettify(s)
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s PutBucketLifecycleOutput) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

const opPutBucketLifecycle = "PutBucketLifecycle"

// PutBucketLifecycleRequest returns a request value for making API operation for
// Amazon Simple Storage Service.
//
// No longer used, see the PutBucketLifecycleConfiguration operation.
//
//    // Example sending a request using PutBucketLifecycleRequest.
//    req := client.PutBucketLifecycleRequest(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
//
// Please also see https://docs.aws.amazon.com/goto/WebAPI/s3-2006-03-01/PutBucketLifecycle
func (c *Client) PutBucketLifecycleRequest(input *PutBucketLifecycleInput) PutBucketLifecycleRequest {
	if c.Client.Config.Logger != nil {
		c.Client.Config.Logger.Log("This operation, PutBucketLifecycle, has been deprecated")
	}
	op := &aws.Operation{
		Name:       opPutBucketLifecycle,
		HTTPMethod: "PUT",
		HTTPPath:   "/{Bucket}?lifecycle",
	}

	if input == nil {
		input = &PutBucketLifecycleInput{}
	}

	req := c.newRequest(op, input, &PutBucketLifecycleOutput{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	return PutBucketLifecycleRequest{Request: req, Input: input, Copy: c.PutBucketLifecycleRequest}
}

// PutBucketLifecycleRequest is the request type for the
// PutBucketLifecycle API operation.
type PutBucketLifecycleRequest struct {
	*aws.Request
	Input *PutBucketLifecycleInput
	Copy  func(*PutBucketLifecycleInput) PutBucketLifecycleRequest
}

// Send marshals and sends the PutBucketLifecycle API request.
func (r PutBucketLifecycleRequest) Send(ctx context.Context) (*PutBucketLifecycleResponse, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &PutBucketLifecycleResponse{
		PutBucketLifecycleOutput: r.Request.Data.(*PutBucketLifecycleOutput),
		response:                 &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// PutBucketLifecycleResponse is the response type for the
// PutBucketLifecycle API operation.
type PutBucketLifecycleResponse struct {
	*PutBucketLifecycleOutput

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// PutBucketLifecycle request.
func (r *PutBucketLifecycleResponse) SDKResponseMetdata() *aws.Response {
	return r.response
}
