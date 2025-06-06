package kithttp

import (
	"github.com/fitan/genx/common"
	"github.com/fitan/genx/gen"
	"golang.org/x/exp/slog"
)

type CEPermissionSqlPlug struct {
}

func (p *CEPermissionSqlPlug) Name() string {
	return "@cep"
}

func (p *CEPermissionSqlPlug) Gen(option gen.Option, implGoTypeMetes []gen.InterfaceGoTypeMeta) (res []gen.GenResult, err error) {
	parseImpl := common.NewInterfaceSerialize(option.Pkg)
	for _, v := range implGoTypeMetes {
		meta, err := parseImpl.Parse(v.Obj, v.RawDoc, &v.Doc)
		if err != nil {
			slog.Error("parseImpl.Parse", err, slog.String("name", v.Obj.String()))
			return nil, common.ParseError("failed to parse interface for CEP plugin").
				WithCause(err).
				WithPlugin("@cep").
				WithInterface(v.Name).
				WithDetails("unable to parse interface metadata for CE Permission SQL generation").
				Build()
		}

		gens, err := CEPGen(option, meta)
		if err != nil {
			return nil, common.GenerateError("failed to generate CE Permission SQL code").
				WithCause(err).
				WithPlugin("@cep").
				WithInterface(v.Name).
				WithDetails("error occurred during CE Permission SQL code generation").
				Build()
		}

		res = append(res, gens...)
	}
	return
}

type Plug struct {
}

func (p *Plug) Name() string {
	return "@kit"
}

func (p *Plug) Gen(option gen.Option, implGoTypeMetes []gen.InterfaceGoTypeMeta) (res []gen.GenResult, err error) {

	parseImpl := common.NewInterfaceSerialize(option.Pkg)
	for _, v := range implGoTypeMetes {
		meta, err := parseImpl.Parse(v.Obj, v.RawDoc, &v.Doc)
		if err != nil {
			slog.Error("parseImpl.Parse", err, slog.String("name", v.Obj.String()))
			return nil, common.ParseError("failed to parse interface for kit-http plugin").
				WithCause(err).
				WithPlugin("@kit").
				WithInterface(v.Name).
				WithDetails("unable to parse interface metadata for Go Kit HTTP service generation").
				Build()
		}

		gens, err := Gen(option, meta)
		if err != nil {
			return nil, common.GenerateError("failed to generate Go Kit HTTP code").
				WithCause(err).
				WithPlugin("@kit").
				WithInterface(v.Name).
				WithDetails("error occurred during Go Kit HTTP service code generation").
				Build()
		}

		res = append(res, gens...)
	}
	return
}

type ObserverPlug struct {
}

func (p *ObserverPlug) Name() string {
	return "@observer"
}

func (p *ObserverPlug) Gen(option gen.Option, implGoTypeMetes []gen.InterfaceGoTypeMeta) (res []gen.GenResult, err error) {
	parseImpl := common.NewInterfaceSerialize(option.Pkg)
	for _, v := range implGoTypeMetes {
		meta, err := parseImpl.Parse(v.Obj, v.RawDoc, &v.Doc)
		if err != nil {
			slog.Error("parseImpl.Parse", err, slog.String("name", v.Obj.String()))
			return nil, common.ParseError("failed to parse interface for observer plugin").
				WithCause(err).
				WithPlugin("@observer").
				WithInterface(v.Name).
				WithDetails("unable to parse interface metadata for observer pattern generation").
				Build()
		}

		gens, err := ObserverGen(option, meta)
		if err != nil {
			return nil, common.GenerateError("failed to generate observer code").
				WithCause(err).
				WithPlugin("@observer").
				WithInterface(v.Name).
				WithDetails("error occurred during observer pattern code generation").
				Build()
		}

		res = append(res, gens...)
	}
	return
}
