package gomodule

import (
	"fmt"
	"path"

	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
)

var (
	archiveBin = pctx.StaticRule("archiveBin", blueprint.RuleParams{
		Command:     "cd $workDir && zip $outputPath -j $toBinary",
		Description: "Archive $toBinary binary",
	}, "workDir", "toBinary", "outputPath")
)

type archiveModuleType struct {
	blueprint.SimpleName

	properties struct {
		ToBinary string
	}
}

func (am *archiveModuleType)GenerateBuildActions(ctx blueprint.ModuleContext) {
	name := ctx.ModuleName()
	config := bood.ExtractConfig(ctx)
	config.Debug.Printf("Adding build actions for go binary module '%s'", name)
	outputPath := path.Join(config.BaseOutputDir, "archive", name)

	ctx.Build(pctx, blueprint.BuildParams{
		Description: fmt.Sprintf("go binary archivation by module %s", name),
		Rule:        archiveBin,
		Outputs:     []string{outputPath},
		Inputs:      []string{path.Join(ctx.ModuleDir(), "out", "bin", am.properties.ToBinary)},
		Args: map[string]string{
			"outputPath": outputPath,
			"workDir":    ctx.ModuleDir(),
			"toBinary":  path.Join("out", "bin", am.properties.ToBinary),
		},
	})
}

func ArchiveBinFactory() (blueprint.Module, []interface{}) {
	mType := &archiveModuleType{}
	return mType, []interface{}{&mType.SimpleName.Properties, &mType.properties}
}