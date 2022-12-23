package validator

import (
	"context"
	"fmt"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	_en "github.com/go-playground/validator/v10/translations/en"
)

func Validate(ctx context.Context, object any, next graphql.Resolver, tags, field string) (any, error) {
	t := en.New()
	universal := ut.New(t, t)
	validate := validator.New()
	translator, _ := universal.GetTranslator("en")

	_ = _en.RegisterDefaultTranslations(validate, translator)
	value := object.(map[string]any)

	if err := validate.Var(value[field], tags); err != nil {
		errs := err.(validator.ValidationErrors)
		for _, vErr := range errs {
			return nil, fmt.Errorf(strings.ToLower(field) + vErr.Translate(translator))
		}
	}

	return next(ctx)
}
