package protoc_gen_base

import (
	"io/ioutil"
	"log"
	"os"

	"google.golang.org/protobuf/proto"
)

func (g *Generator) toolName() string {
	return "protoc-gen-" + g.name
}
func (g *Generator) pluginName() string {
	return plugins[0].Name()
}
func Main(name string, ps ...Plugin) {
	for _, p := range ps {
		RegisterPlugin(p)
	}
	// Begin by allocating a generator. The request and response structures are stored there
	// so we can do error handling easily - the response structure contains the field to
	// report failure.
	var (
		data []byte
		err  error
	)
	// just for debug INPUT=xx protoc ....
	if in, ok := os.LookupEnv("INPUT"); ok {
		f, err := os.Open(in)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		data, err = ioutil.ReadAll(f)
	} else {
		data, err = ioutil.ReadAll(os.Stdin)
	}

	g := New(name)

	if err != nil {
		g.Error(err, "reading input")
	}

	if err := proto.Unmarshal(data, g.Request); err != nil {
		g.Error(err, "parsing input proto")
	}

	if len(g.Request.FileToGenerate) == 0 {
		g.Fail("no files to generate")
	}

	//g.CommandLineParameters(AddPluginToParams(g.Request.GetParameter()))
	g.CommandLineParameters(g.Request.GetParameter())

	// Create a wrapped version of the Descriptors and EnumDescriptors that
	// point to the file that defines them.
	g.WrapTypes()

	g.SetPackageNames()
	g.BuildTypeNameMap()

	g.GenerateAllFiles()

	// Send back the results.
	data, err = proto.Marshal(g.Response)
	if err != nil {
		g.Error(err, "failed to marshal output proto")
	}
	_, err = os.Stdout.Write(data)
	if err != nil {
		g.Error(err, "failed to write output proto")
	}
}
