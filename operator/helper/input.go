// Copyright 2021, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package helper

import (
	"context"

	"github.com/opentelemetry/opentelemetry-log-collection/entry"
	"github.com/opentelemetry/opentelemetry-log-collection/errors"
	"github.com/opentelemetry/opentelemetry-log-collection/operator"
	"go.uber.org/zap"
)

// NewInputConfig creates a new input config with default values.
func NewInputConfig(operatorID, operatorType string) InputConfig {
	return InputConfig{
		LabelerConfig:    NewLabelerConfig(),
		IdentifierConfig: NewIdentifierConfig(),
		WriterConfig:     NewWriterConfig(operatorID, operatorType),
		WriteTo:          entry.NewRecordField(),
	}
}

// InputConfig provides a basic implementation of an input operator config.
type InputConfig struct {
	LabelerConfig    `yaml:",inline"`
	IdentifierConfig `yaml:",inline"`
	WriterConfig     `yaml:",inline"`
	WriteTo          entry.Field `json:"write_to" yaml:"write_to"`
}

// Build will build a base producer.
func (c InputConfig) Build(context operator.BuildContext) (InputOperator, error) {
	writerOperator, err := c.WriterConfig.Build(context)
	if err != nil {
		return InputOperator{}, errors.WithDetails(err, "operator_id", c.ID())
	}

	labeler, err := c.LabelerConfig.Build()
	if err != nil {
		return InputOperator{}, errors.WithDetails(err, "operator_id", c.ID())
	}

	identifier, err := c.IdentifierConfig.Build()
	if err != nil {
		return InputOperator{}, errors.WithDetails(err, "operator_id", c.ID())
	}

	inputOperator := InputOperator{
		Labeler:        labeler,
		Identifier:     identifier,
		WriterOperator: writerOperator,
		WriteTo:        c.WriteTo,
	}

	return inputOperator, nil
}

// InputOperator provides a basic implementation of an input operator.
type InputOperator struct {
	Labeler
	Identifier
	WriterOperator
	WriteTo entry.Field
}

// NewEntry will create a new entry using the `write_to`, `labels`, and `resource` configuration.
func (i *InputOperator) NewEntry(value interface{}) (*entry.Entry, error) {
	entry := entry.New()
	if err := entry.Set(i.WriteTo, value); err != nil {
		return nil, errors.Wrap(err, "add record to entry")
	}

	if err := i.Label(entry); err != nil {
		return nil, errors.Wrap(err, "add labels to entry")
	}

	if err := i.Identify(entry); err != nil {
		return nil, errors.Wrap(err, "add resource keys to entry")
	}

	return entry, nil
}

// CanProcess will always return false for an input operator.
func (i *InputOperator) CanProcess() bool {
	return false
}

// Process will always return an error if called.
func (i *InputOperator) Process(ctx context.Context, entry *entry.Entry) error {
	i.Errorw("Operator received an entry, but can not process", zap.Any("entry", entry))
	return errors.NewError(
		"Operator can not process logs.",
		"Ensure that operator is not configured to receive logs from other operators",
	)
}
