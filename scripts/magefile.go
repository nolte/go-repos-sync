//+build mage

package main

import (
	"context"

	"github.com/magefile/mage/mg"

	"github.com/nolte/plumbing/cmd/golang"
)

// Build configure the Build Targets.
type Build mg.Namespace

func (Build) Lint(ctx context.Context) {
	ctx = context.WithValue(ctx, "basedir", "../")
	mg.CtxDeps(ctx, golang.Golang.Lint)
}

func (Build) Fmt(ctx context.Context) {
	ctx = context.WithValue(ctx, "basedir", "../")
	mg.CtxDeps(ctx, golang.Golang.Fmt)
}
