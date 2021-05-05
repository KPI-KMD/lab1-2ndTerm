package gomodule

import (
	"fmt"
	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
	"path"
	"runtime"
)
// Package context used to define Ninja build rules.
var (
	pctx = blueprint.NewPackageContext("github.com/KPI-KMD/lab1-2ndTerm/bood/gomodule")
	goBuild, goVendor, goTest blueprint.Rule
)

func init() {

	var build, vendor, test string
	switch os := runtime.GOOS; os {
	case "windows":
		build = "cmd /c cd $workDir && go build -o $outputPath $pkg"
		vendor = "cmd /c cd $workDir && go mod vendor"
		test = "cmd /c cd $workDir && go test -v $testPkg > $outputPath"
	default:
		build = "cd $workDir && go build -o $outputPath $pkg"
		vendor = "cd $workDir && go mod vendor"
		test = "cd $workDir && go test $testPkg > $output"
	}

	goBuild = pctx.StaticRule("binaryBuild", blueprint.RuleParams{
		Command:    build,
		Description: "build go command $pkg",
	}, "workDir", "outputPath", "pkg")

	goVendor = pctx.StaticRule("vendor", blueprint.RuleParams{
		Command:     vendor,
		Description: "vendor dependencies of $name",
	}, "workDir", "name")

	goTest = pctx.StaticRule("test", blueprint.RuleParams{
		Command: test,
		Description: "test package $testPkg",
	}, "workDir", "testPkg", "output")
}



// goBinaryModuleType implements the simplest Go binary build without running tests for the target Go package.
type goBinaryModuleType struct {
	blueprint.SimpleName

	properties struct {
		// Go package name to build as a command with "go build".
		Pkg string
		// Go package name to test
		TestPkg string
		// Go test excludes
		TestSrcs []string
		// List of source files.
		Srcs []string
		// Exclude patterns.
		SrcsExclude []string
		// If to call vendor command.
		VendorFirst bool

		// Example of how to specify dependencies.
		Deps []string
	}
}

func (gb *goBinaryModuleType) DynamicDependencies(blueprint.DynamicDependerModuleContext) []string {
	return gb.properties.Deps
}

func (gb *goBinaryModuleType) GenerateBuildActions(ctx blueprint.ModuleContext) {
	name := ctx.ModuleName()
	config := bood.ExtractConfig(ctx)
	config.Debug.Printf("Adding build actions for go binary module '%s'", name)

	outputPath := path.Join(config.BaseOutputDir, "bin", name)
	output := path.Join(config.BaseOutputDir, name, "bood_test")

	var inputs []string
	var testInputs []string
	inputErors := false
	for _, src := range gb.properties.Srcs {
		if matches, err := ctx.GlobWithDeps(src, gb.properties.SrcsExclude); err == nil {
			inputs = append(inputs, matches...)
		} else {
			ctx.PropertyErrorf("srcs", "Cannot resolve files that match pattern %s", src)
			inputErors = true
		}
	}
	for _, src := range gb.properties.TestSrcs {
		if matches, err := ctx.GlobWithDeps(src, nil); err == nil {
			testInputs = append(testInputs, matches...)
		} else {
			ctx.PropertyErrorf("srcs", "Cannot resolve files that match pattern %s", src)
			inputErors = true
		}
	}
	if inputErors {
		return
	}

	if gb.properties.VendorFirst {
		vendorDirPath := path.Join(ctx.ModuleDir(), "vendor")
		ctx.Build(pctx, blueprint.BuildParams{
			Description: fmt.Sprintf("Vendor dependencies of %s", name),
			Rule:        goVendor,
			Outputs:     []string{vendorDirPath},
			Implicits:   []string{path.Join(ctx.ModuleDir(), "go.mod")},
			Optional:    true,
			Args: map[string]string{
				"workDir": ctx.ModuleDir(),
				"name":    name,
			},
		})
		inputs = append(inputs, vendorDirPath)
	}

	ctx.Build(pctx, blueprint.BuildParams{
		Description: fmt.Sprintf("Build %s as Go binary", name),
		Rule:        goBuild,
		Outputs:     []string{outputPath},
		Implicits:   inputs,
		Args: map[string]string{
			"outputPath": outputPath,
			"workDir":    ctx.ModuleDir(),
			"pkg":        gb.properties.Pkg,
		},
	})

	testInputs = append(testInputs, inputs...)
	ctx.Build(pctx, blueprint.BuildParams{
		Description: fmt.Sprintf("Test %s as Go binary", name),
		Rule:        goTest,
		Outputs:     []string{output},
		Implicits:   testInputs,
		Args: map[string]string{
			"output":     output,
			"workDir":    ctx.ModuleDir(),
			"testPkg":    gb.properties.TestPkg,
		},
	})
}

// TestBinFactory is a factory for go binary module type which supports Go command packages with running tests.
func TestBinFactory() (blueprint.Module, []interface{}) {
	mType := &goBinaryModuleType{}
	return mType, []interface{}{&mType.SimpleName.Properties, &mType.properties}
}
