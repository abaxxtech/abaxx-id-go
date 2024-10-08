package main

import (
	"context"

	"github.com/alecthomas/kong"
)

// CLI is the main command line interface
// more information about this struct can be found in the [kong documentation]
//
// [kong documentation]: https://github.com/alecthomas/kong
type CLI struct {
	JWT struct {
		Sign   jwtSignCMD   `cmd:"" help:"Sign a JWT."`
		Decode jwtDecodeCMD `cmd:"" help:"Decode a JWT."`
		Verify jwtVerifyCMD `cmd:"" help:"Verify a JWT."`
	} `cmd:"" help:"Interface with JWT's."`
	DID struct {
		Resolve didResolveCMD `cmd:"" help:"Resolve a DID."`
		Create  didCreateCMD  `cmd:"" help:"Create a DID."`
	} `cmd:"" help:"Interface with DID's."`
	VC struct {
		Create vcCreateCMD `cmd:"" help:"Create a VC."`
		Sign   vcSignCMD   `cmd:"" help:"Sign a VC."`
	} `cmd:"" help:"Interface with VC's."`
	VCJWT struct {
		Verify vcjwtVerifyCMD `cmd:"" help:"Verify a VC-JWT."`
		Decode vcjwtDecodeCMD `cmd:"" help:"Decode a VC-JWT."`
	} `cmd:"" help:"Interface with VC-JWT's."`
}

func main() {
	kctx := kong.Parse(&CLI{},
		kong.Description("Abaxx ID++ - A decentralized identity platform."),
	)

	ctx := context.Background()
	kctx.BindTo(ctx, (*context.Context)(nil))
	err := kctx.Run(ctx)
	kctx.FatalIfErrorf(err)
}
