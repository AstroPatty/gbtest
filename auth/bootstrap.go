package auth

import (
    "fmt"
    "github.com/aws/aws-cdk-go/awscdk/v2/cxapi"
)


func CheckBoostrap(assembly cxapi.CloudAssembly) error {
    // var assets *cloudassemblyschema.ArtifactManifest

    for _, artifact := range *assembly.Artifacts() {
        fmt.Println(artifact.Manifest().Type)
    }

    //fmt.Println(assets.Properties)
    return nil
}

