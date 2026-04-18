// Package ignore provides support for .driftignore files.
//
// A .driftignore file lists resources that should be excluded from drift
// detection results. Each non-comment line must follow the format:
//
//	resource_type/resource_id
//
// A wildcard (*) may be used as the resource_id to ignore all resources
// of a given type:
//
//	aws_instance/*
//
// Lines starting with '#' are treated as comments and ignored.
package ignore
