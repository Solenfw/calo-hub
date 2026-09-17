// Package service contains outbound integrations and other service-layer helpers
// used by the application.
//
// This package is intentionally kept separate from the repo and HTTP layers so
// the app can call external APIs without leaking transport or HTTP concerns into
// the domain logic.
package service

// The project currently uses this package as the integration boundary for
// outbound HTTP calls and future service-layer orchestration.
