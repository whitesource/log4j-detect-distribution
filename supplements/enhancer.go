package supplements

import (
	"github.com/go-logr/logr"
	"github.com/whitesource/log4j-detect/records"
)

type Enhancer interface {
	// Enhance assumption: opr not nil
	Enhance(opr records.OperationResult) records.EnhancedResult
	LType() records.LibType
}

func Supplement(logger logr.Logger, opResults []records.OperationResult) []records.EnhancedResult {
	result := make([]records.EnhancedResult, 0, len(opResults))
	enhancers := availableEnhancers(logger)

	for _, opr := range opResults {
		if opr.Err != nil {
			continue
		}

		en := enhancers[opr.LType]
		if en != nil {
			enhanced := en.Enhance(opr)
			result = append(result, enhanced)
		}
	}
	return result
}

func availableEnhancers(logger logr.Logger) map[records.LibType]Enhancer {
	java := NewJavaEnhancer(logger)
	ruby := NewRubyEnhancer(logger)
	fs := NewFsEnhancer()

	available := map[records.LibType]Enhancer{
		java.LType(): java,
		fs.LType():   fs,
		ruby.LType(): ruby,
	}

	return available
}
