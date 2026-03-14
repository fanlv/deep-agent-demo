/*
 * Copyright 2025 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package tools

import (
	"context"
	"log"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

type ToolRequestPreprocess func(ctx context.Context, baseTool tool.InvokableTool, toolArguments string) (string, error)

type ToolResponsePostprocess func(ctx context.Context, baseTool tool.InvokableTool, toolResponse, toolArguments string) (string, error)

func NewWrapTool(t tool.InvokableTool, preprocess []ToolRequestPreprocess, postprocess []ToolResponsePostprocess) tool.InvokableTool {
	return &wrapTool{
		baseTool:    t,
		preprocess:  preprocess,
		postprocess: postprocess,
	}
}

type wrapTool struct {
	baseTool    tool.InvokableTool
	preprocess  []ToolRequestPreprocess
	postprocess []ToolResponsePostprocess
}

func (w *wrapTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return w.baseTool.Info(ctx)
}

func (w *wrapTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	for _, pre := range w.preprocess {
		var err error
		argumentsInJSON, err = pre(ctx, w.baseTool, argumentsInJSON)
		if err != nil {
			log.Printf("[WrapTool.PreProcess] failed to process response: %v", err)
		}
	}

	resp, err := w.baseTool.InvokableRun(ctx, argumentsInJSON, opts...)
	if err != nil {
		return err.Error(), nil
	}

	for _, post := range w.postprocess {
		resp, err = post(ctx, w.baseTool, resp, argumentsInJSON)
		if err != nil {
			log.Printf("[WrapTool.PostProcess] failed to process response: %v", err)
			return resp, err
		}
	}

	return resp, nil
}
