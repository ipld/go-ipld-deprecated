package main

import (
	"io/ioutil"
	"flag"
	"os"

	mc "github.com/jbenet/go-multicodec"
	pb "github.com/ipfs/go-ipld/coding/pb"
	ipld "github.com/ipfs/go-ipld"
	coding "github.com/ipfs/go-ipld/coding"
)

var codecs []mc.Multicodec = []mc.Multicodec{
	coding.CborMulticodec(),
	coding.JsonMulticodec(),
	pb.Multicodec(),
}

func codecByName(name string) mc.Multicodec {
	for _, c := range codecs {
		if name == string(mc.HeaderPath(c.Header())) {
			return c
		}
	}
	return nil
}

func main() {
	infile  := flag.String("i", "", "Input file")
	outfile := flag.String("o", "", "Output file")
	codecid := flag.String("c", "", "Multicodec to use")
	flag.Parse()
	file, err := ioutil.ReadFile(*infile)
	if err != nil {
		panic(err)
	}

	var n ipld.Node
	codec := coding.Multicodec()

	if err := mc.Unmarshal(codec, file, &n); err != nil {
		panic(err)
	}

	codec = codecByName(*codecid)
	if codec == nil {
		panic("Could not find codec " + *codecid)
	}

	delete(n, ipld.CodecKey)

	encoded, err := mc.Marshal(codec, &n)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(*outfile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.Write(encoded);
	if err != nil {
		panic(err)
	}
}


