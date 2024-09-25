package graph

import "go.opentelemetry.io/otel/trace"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{
	Tracer trace.Tracer
}
