// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

// Package servicediscovery provides the client and types for making API
// requests to Amazon Route 53 Auto Naming.
//
// Amazon Route 53 auto naming lets you configure public or private namespaces
// that your microservice applications run in. When instances of the service
// become available, you can call the auto naming API to register the instance,
// and Route 53 automatically creates up to five DNS records and an optional
// health check. Clients that submit DNS queries for the service receive an
// answer that contains up to eight healthy records.
//
// See https://docs.aws.amazon.com/goto/WebAPI/servicediscovery-2017-03-14 for more information on this service.
//
// See servicediscovery package documentation for more information.
// https://docs.aws.amazon.com/sdk-for-go/api/service/servicediscovery/
//
// Using the Client
//
// To contact Amazon Route 53 Auto Naming with the SDK use the New function to create
// a new service client. With that client you can make API requests to the service.
// These clients are safe to use concurrently.
//
// See the SDK's documentation for more information on how to use the SDK.
// https://docs.aws.amazon.com/sdk-for-go/api/
//
// See aws.Config documentation for more information on configuring SDK clients.
// https://docs.aws.amazon.com/sdk-for-go/api/aws/#Config
//
// See the Amazon Route 53 Auto Naming client ServiceDiscovery for more
// information on creating client for this service.
// https://docs.aws.amazon.com/sdk-for-go/api/service/servicediscovery/#New
package servicediscovery
