package resource

import (
   "fmt"
   "github.com/newrelic/newrelic-cloudformation-resource-providers-common/cferror"
   "github.com/newrelic/newrelic-cloudformation-resource-providers-common/client/nerdgraph"
   "github.com/newrelic/newrelic-cloudformation-resource-providers-common/model"
   log "github.com/sirupsen/logrus"
   "strings"
)

// ErrorHandler at a minimum provides access to the default error processing.
// If required we can provide custom processing here via composition overrides
type ErrorHandler struct {
   // Use Go composition to access the default implementation
   model.ErrorHandler
   M model.Model
}

// NewErrorHandler This is all pretty magical. We return the interface so common is insulated from an implementation. Payload implements model.Model so all is good
func NewErrorHandler(p *Payload) (h model.ErrorHandler) {
   defer func() {
      log.Debugf("errorHandler.NewErrorHandler: exit %p", h)
   }()
   // Initialize ourself with the common core
   h = &ErrorHandler{ErrorHandler: nerdgraph.NewCommonErrorHandler(p), M: p}
   return
}

func (e *ErrorHandler) TypeSpecificError(data *[]byte, s string) (err error) {
   log.Debugf("TypeSpecificError: %p enter", e)
   defer func() {
      log.Debugf("TypeSpecificError: %p returning %v", e, err)
   }()
   v, err := nerdgraph.FindKeyValue(*data, "errors")
   log.Debugf("TypeSpecificError: found: %v %T", v, v)
   if err != nil {
      return
   }
   if v == nil {
      return
   }

   errorMap := make(map[string]interface{})
   e.GetErrorMap(v, errorMap)

   if errorMap == nil {
      log.Warnf("Empty errors array: %v+ %T", e, e)
      return
   }
   _type := fmt.Sprintf("%v", errorMap[e.M.GetErrorKey()])
   if strings.Contains(strings.ToLower(_type), "not_found") || strings.Contains(strings.ToLower(_type), "not found") {
      err = fmt.Errorf("%w Not found", &cferror.NotFound{})
      return
   }
   return
}
