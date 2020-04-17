// Code generated by piper's step-generator. DO NOT EDIT.

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/SAP/jenkins-library/pkg/config"
	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/SAP/jenkins-library/pkg/piperenv"
	"github.com/SAP/jenkins-library/pkg/telemetry"
	"github.com/spf13/cobra"
)

type artifactPrepareVersionOptions struct {
	BuildTool            string `json:"buildTool,omitempty"`
	CommitUserName       string `json:"commitUserName,omitempty"`
	CustomversionField   string `json:"customversionField,omitempty"`
	CustomVersionSection string `json:"customVersionSection,omitempty"`
	DockerVersionSource  string `json:"dockerVersionSource,omitempty"`
	FilePath             string `json:"filePath,omitempty"`
	GlobalSettingsFile   string `json:"globalSettingsFile,omitempty"`
	IncludeCommitID      bool   `json:"includeCommitId,omitempty"`
	M2Path               string `json:"m2Path,omitempty"`
	Password             string `json:"password,omitempty"`
	ProjectSettingsFile  string `json:"projectSettingsFile,omitempty"`
	TagPrefix            string `json:"tagPrefix,omitempty"`
	Username             string `json:"username,omitempty"`
	VersioningTemplate   string `json:"versioningTemplate,omitempty"`
	VersioningType       string `json:"versioningType,omitempty"`
}

type artifactPrepareVersionCommonPipelineEnvironment struct {
	artifactVersion string
	git             struct {
		commitID string
	}
}

func (p *artifactPrepareVersionCommonPipelineEnvironment) persist(path, resourceName string) {
	content := []struct {
		category string
		name     string
		value    string
	}{
		{category: "", name: "artifactVersion", value: p.artifactVersion},
		{category: "git", name: "commitId", value: p.git.commitID},
	}

	errCount := 0
	for _, param := range content {
		err := piperenv.SetResourceParameter(path, resourceName, filepath.Join(param.category, param.name), param.value)
		if err != nil {
			log.Entry().WithError(err).Error("Error persisting piper environment.")
			errCount++
		}
	}
	if errCount > 0 {
		os.Exit(1)
	}
}

// ArtifactPrepareVersionCommand Prepares and potentially updates the artifact's version before building the artifact.
func ArtifactPrepareVersionCommand() *cobra.Command {
	metadata := artifactPrepareVersionMetadata()
	var stepConfig artifactPrepareVersionOptions
	var startTime time.Time
	var commonPipelineEnvironment artifactPrepareVersionCommonPipelineEnvironment

	var createArtifactPrepareVersionCmd = &cobra.Command{
		Use:   "artifactPrepareVersion",
		Short: "Prepares and potentially updates the artifact's version before building the artifact.",
		Long: `Prepares and potentially updates the artifact's version before building the artifact.

The continuous delivery process requires that each build is done with a unique version number.
There are two common patterns found:

### 1. Continuous Deployment pattern with automatic versioning

The team has full authority on ` + "`" + `<major>.<minor>.<patch>` + "`" + ` and can increase any part whenever required.
Nonetheless, the automatic versioning makes sure that every build will create a unique version by appending ` + "`" + `<major>.<minor>.<patch>` + "`" + ` with a buildversion (we use a timestamp) and optinally the commitId.

In order to represent this version also in the version control system the new unique version will be pushed with a dedicated tag (` + "`" + `<tagPrefix><major>.<minor>.<patch><unique extension>` + "`" + `).

Depending on the build tool used and thus the allowed versioning format the ` + "`" + `<unique extension>` + "`" + ` varies.

**Remarks:**

* There is no commit to master since this would create a perpetuum mobile and just trigger the next automatic build with automatic versioning, and so on ...
* Not creating a tag would lead to a loss of the final artifact version in scm which often is not acceptable
* You need to ensure that your CI/CD system can push back to your SCM (via providing ssh or HTTP(s) credentials)

**This pattern is the default** behavior (` + "`" + `versioningType: cloud` + "`" + `) since this is suitable for for most cloud deliveries.

It is possible to use ` + "`" + `versioningType: cloud_noTag` + "`" + ` which has a slighly different behavior than described above:

* The new version will NOT be written as tag into the SCM but it is only available in the corresponding CI/CD workspace
* IMPORTANT NOTICE: Using the option ` + "`" + `cloud_noTag` + "`" + ` should not be picked in case you need to ensure a fully traceable path from SCM commit to your build artifact.

### 2. Pure version ` + "`" + `<major>.<minor>.<patch>` + "`" + `

This pattern is often used by teams that have cloud deliveries with no fully automated procedure, e.g. delivery after each takt.
Another typical use-case is development of a library with regular releases where the versioning pattern should be consumable and thus ideally complies to a ` + "`" + `<major>.<minor>.<patch>` + "`" + ` pattern.

The version is then either manually set by the team in the course of the development process or automatically pushed to master after a successful release.

Unlike for the _Continuous Deloyment_ pattern descibed above, in this case there is no dedicated tagging required for the build process since the version is already available in the repository.

Configuration of this pattern is done via ` + "`" + `versioningType: library` + "`" + `.

### Support of additional build tools

Besides the ` + "`" + `buildTools` + "`" + ` provided out of the box (like ` + "`" + `maven` + "`" + `, ` + "`" + `mta` + "`" + `, ` + "`" + `npm` + "`" + `, ...) it is possible to set ` + "`" + `buildTool: custom` + "`" + `.

This allows you to provide automatic versioning for tools using a:

#### file with the version as only content:

Define ` + "`" + `buildTool: custom` + "`" + ` as well as ` + "`" + `filePath: <path to your file>` + "`" + `

**Please note:** ` + "`" + `<path to your file>` + "`" + ` need to point either to a ` + "`" + `*.txt` + "`" + ` file or to a file without extension.

#### ` + "`" + `ini` + "`" + ` file containing the version:

Define ` + "`" + `buildTool: custom` + "`" + `, ` + "`" + `filePath: <path to your ini-file>` + "`" + ` as well as parameters ` + "`" + `versionSection` + "`" + ` and ` + "`" + `versionSource` + "`" + ` to point to the version location (section & parameter name) within the file.

**Please note:** ` + "`" + `<path to your file>` + "`" + ` need to point either to a ` + "`" + `*.cfg` + "`" + ` or a ` + "`" + `*.ini` + "`" + ` file.

#### ` + "`" + `json` + "`" + ` file containing the version:

Define ` + "`" + `buildTool: custom` + "`" + `, ` + "`" + `filePath: <path to your *.json file` + "`" + ` as well as parameter ` + "`" + `versionSource` + "`" + ` to point to the parameter containing the version.

#### ` + "`" + `yaml` + "`" + ` file containing the version

Define ` + "`" + `buildTool: custom` + "`" + `, ` + "`" + `filePath: <path to your *.yml/*.yaml file` + "`" + ` as well as parameter ` + "`" + `versionSource` + "`" + ` to point to the parameter containing the version.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			startTime = time.Now()
			log.SetStepName("artifactPrepareVersion")
			log.SetVerbose(GeneralConfig.Verbose)
			err := PrepareConfig(cmd, &metadata, "artifactPrepareVersion", &stepConfig, config.OpenPiperFile)
			if err != nil {
				return err
			}
			log.RegisterSecret(stepConfig.Password)
			log.RegisterSecret(stepConfig.Username)
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			telemetryData := telemetry.CustomData{}
			telemetryData.ErrorCode = "1"
			handler := func() {
				commonPipelineEnvironment.persist(GeneralConfig.EnvRootPath, "commonPipelineEnvironment")
				telemetryData.Duration = fmt.Sprintf("%v", time.Since(startTime).Milliseconds())
				telemetry.Send(&telemetryData)
			}
			log.DeferExitHandler(handler)
			defer handler()
			telemetry.Initialize(GeneralConfig.NoTelemetry, "artifactPrepareVersion")
			artifactPrepareVersion(stepConfig, &telemetryData, &commonPipelineEnvironment)
			telemetryData.ErrorCode = "0"
		},
	}

	addArtifactPrepareVersionFlags(createArtifactPrepareVersionCmd, &stepConfig)
	return createArtifactPrepareVersionCmd
}

func addArtifactPrepareVersionFlags(cmd *cobra.Command, stepConfig *artifactPrepareVersionOptions) {
	cmd.Flags().StringVar(&stepConfig.BuildTool, "buildTool", os.Getenv("PIPER_buildTool"), "Defines the tool which is used for building the artifact. Supports `custom`, `dub`, `golang`, `maven`, `mta`, `npm`, `pip`, `sbt`.")
	cmd.Flags().StringVar(&stepConfig.CommitUserName, "commitUserName", "Project Piper", "Defines the user name which appears in version control for the versioning update (in case `versioningType: cloud`).")
	cmd.Flags().StringVar(&stepConfig.CustomversionField, "customversionField", os.Getenv("PIPER_customversionField"), "For `buildTool: custom`: Defines the field which contains the version in the descriptor file.")
	cmd.Flags().StringVar(&stepConfig.CustomVersionSection, "customVersionSection", os.Getenv("PIPER_customVersionSection"), "For `buildTool: custom`: Defines the section for version retrieval in vase a *.ini/*.cfg file is used.")
	cmd.Flags().StringVar(&stepConfig.DockerVersionSource, "dockerVersionSource", os.Getenv("PIPER_dockerVersionSource"), "For `buildTool: docker`: Defines the source of the version. Can be `FROM`, any supported _buildTool_ or an environment variable name.")
	cmd.Flags().StringVar(&stepConfig.FilePath, "filePath", os.Getenv("PIPER_filePath"), "Defines a custom path to the descriptor file. Build tool specific defaults are used (e.g. `maven: pom.xml`, `npm: package.json`, `mta: mta.yaml`).")
	cmd.Flags().StringVar(&stepConfig.GlobalSettingsFile, "globalSettingsFile", os.Getenv("PIPER_globalSettingsFile"), "Maven only - Path to the mvn settings file that should be used as global settings file.")
	cmd.Flags().BoolVar(&stepConfig.IncludeCommitID, "includeCommitId", true, "Defines if the automatically generated version (`versioningType: cloud`) should include the commit id hash.")
	cmd.Flags().StringVar(&stepConfig.M2Path, "m2Path", os.Getenv("PIPER_m2Path"), "Maven only - Path to the location of the local repository that should be used.")
	cmd.Flags().StringVar(&stepConfig.Password, "password", os.Getenv("PIPER_password"), "Password/token for git authentication.")
	cmd.Flags().StringVar(&stepConfig.ProjectSettingsFile, "projectSettingsFile", os.Getenv("PIPER_projectSettingsFile"), "Maven only - Path to the mvn settings file that should be used as project settings file.")
	cmd.Flags().StringVar(&stepConfig.TagPrefix, "tagPrefix", "build_", "Defines the prefix which is used for the git tag which is written during the versioning run (only `versioningType: cloud`).")
	cmd.Flags().StringVar(&stepConfig.Username, "username", os.Getenv("PIPER_username"), "User name for git authentication")
	cmd.Flags().StringVar(&stepConfig.VersioningTemplate, "versioningTemplate", os.Getenv("PIPER_versioningTemplate"), "DEPRECATED: Defines the template for the automatic version which will be created")
	cmd.Flags().StringVar(&stepConfig.VersioningType, "versioningType", "cloud", "Defines the type of versioning (`cloud`: fully automatic, `cloud_noTag`: automatic but no tag created, `library`: manual)")

	cmd.MarkFlagRequired("buildTool")
}

// retrieve step metadata
func artifactPrepareVersionMetadata() config.StepData {
	var theMetaData = config.StepData{
		Metadata: config.StepMetadata{
			Name:    "artifactPrepareVersion",
			Aliases: []config.Alias{{Name: "artifactSetVersion", Deprecated: false}, {Name: "setVersion", Deprecated: true}},
		},
		Spec: config.StepSpec{
			Inputs: config.StepInputs{
				Parameters: []config.StepParameters{
					{
						Name:        "buildTool",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"GENERAL", "PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   true,
						Aliases:     []config.Alias{},
					},
					{
						Name:        "commitUserName",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{{Name: "gitUserName"}},
					},
					{
						Name:        "customversionField",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
					},
					{
						Name:        "customVersionSection",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
					},
					{
						Name:        "dockerVersionSource",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
					},
					{
						Name:        "filePath",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
					},
					{
						Name:        "globalSettingsFile",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"GENERAL", "STEPS", "STAGES", "PARAMETERS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{{Name: "maven/globalSettingsFile"}},
					},
					{
						Name:        "includeCommitId",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "bool",
						Mandatory:   false,
						Aliases:     []config.Alias{},
					},
					{
						Name:        "m2Path",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"GENERAL", "STEPS", "STAGES", "PARAMETERS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{{Name: "maven/m2Path"}},
					},
					{
						Name:        "password",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
					},
					{
						Name:        "projectSettingsFile",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"GENERAL", "STEPS", "STAGES", "PARAMETERS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{{Name: "maven/projectSettingsFile"}},
					},
					{
						Name:        "tagPrefix",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
					},
					{
						Name:        "username",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
					},
					{
						Name:        "versioningTemplate",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
					},
					{
						Name:        "versioningType",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
					},
				},
			},
		},
	}
	return theMetaData
}