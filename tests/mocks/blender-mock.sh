#!/bin/bash
# Mock Blender script that simulates rendering 8 product photos
# Outputs progress like real Blender and creates PNG files
#
# Usage:
#   blender-mock.sh -- --sku PROD-001 --overlay overlay.png --output ./render
#
# This mock:
# - Parses arguments the same way as real Blender Python script
# - Outputs progress matching real Blender format
# - Creates 8 PNG files (render-1.png through render-8.png)
# - Completes in ~2 seconds instead of 5+ minutes

SKU=""
OVERLAY=""
OUTPUT=""

# Parse arguments (after "--")
while [[ $# -gt 0 ]]; do
    case $1 in
        --sku)
            SKU="$2"
            shift 2
            ;;
        --overlay)
            OVERLAY="$2"
            shift 2
            ;;
        --output)
            OUTPUT="$2"
            shift 2
            ;;
        *)
            shift
            ;;
    esac
done

# Validate required arguments
if [[ -z "$SKU" || -z "$OVERLAY" || -z "$OUTPUT" ]]; then
    echo "Error: Missing required arguments"
    echo "Usage: blender-mock.sh -- --sku SKU --overlay OVERLAY --output OUTPUT"
    exit 1
fi

# Output Blender-like header
echo "Blender 3.6.0"
echo "Rendering product: $SKU"
echo "Overlay: $OVERLAY"
echo "Output: $OUTPUT"
echo ""

# Render 8 angles (simulate 0.2s per frame)
for i in {1..8}; do
    # Output progress like real Blender
    echo "Fra:$i Mem:12.00M (Peak 12.00M) | Rendering $i/8"

    # Create mock PNG file (minimal valid PNG - 1x1 pixel)
    # PNG signature + minimal IHDR chunk + IEND chunk
    printf "\x89\x50\x4e\x47\x0d\x0a\x1a\x0a\x00\x00\x00\x0d\x49\x48\x44\x52\x00\x00\x00\x01\x00\x00\x00\x01\x08\x02\x00\x00\x00\x90\x77\x53\xde\x00\x00\x00\x0c\x49\x44\x41\x54\x08\x99\x63\xf8\xcf\xc0\x00\x00\x00\x03\x00\x01\x00\x18\xdd\x8d\xb4\x00\x00\x00\x00\x49\x45\x4e\x44\xae\x42\x60\x82" > "${OUTPUT}-${i}.png"

    # Simulate render time
    sleep 0.2
done

echo ""
echo "✓ Rendered 8 photos for $SKU"
exit 0
