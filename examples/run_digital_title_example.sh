#!/bin/bash

# Digital Title Legal Agreement Verifiable Credential Example
# This script demonstrates how to run the complete end-to-end example

echo "🚀 Digital Title Legal Agreement as Verifiable Credential - Full Example"
echo "======================================================================="
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    exit 1
fi

echo "📋 Running the complete digital title example..."
echo ""

# Run the example
cd "$(dirname "$0")" || exit 1

# Build and run the example
go run digital_title_full_example.go

echo ""
echo "🎉 Digital Title Example Complete!"
echo ""
echo "What this example demonstrated:"
echo "✅ Created DIDs for trustor (Xco) and trustee (Abaxx Technologies)"
echo "✅ Built complete digital title legal agreement structure"
echo "✅ Created full verifiable credential with entire agreement"
echo "✅ Signed credential to generate VC-JWT"
echo "✅ Verified credential and extracted data"
echo "✅ Created reference-based credential for role assertions"
echo "✅ Generated CLI commands for practical usage"
echo "✅ 🔐 Generated cryptographic proofs and verification details"
echo "✅ 📋 Created verifiable presentation (VP) for enhanced proof"
echo "✅ 🔗 Computed integrity hashes for tamper-evidence"
echo "✅ 💾 Exported complete proof package for third-party verification"
echo ""
echo "🔐 CRYPTOGRAPHIC PROOF FEATURES:"
echo "   • Ed25519 digital signature verification"
echo "   • JWT component analysis (header, payload, signature)"
echo "   • Verifiable presentation generation"
echo "   • Content integrity hash computation"
echo "   • Exportable proof package with verification keys"
echo "   • Third-party verification support"
