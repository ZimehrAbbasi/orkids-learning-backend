package utils

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func WrapWithTracer(ctx context.Context, domain, name, string1 string, funcToExecute func(string) error) error {
	tracer := otel.Tracer(domain)
	ctx, span := tracer.Start(ctx, name)
	err := funcToExecute(string1)
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
	}
	span.End()
	return err
}

func WrapWithTracer2(ctx context.Context, domain, name, string1, string2 string, funcToExecute func(string, string) error) error {
	tracer := otel.Tracer(domain)
	ctx, span := tracer.Start(ctx, name)
	err := funcToExecute(string1, string2)
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
	}
	span.End()
	return err
}
