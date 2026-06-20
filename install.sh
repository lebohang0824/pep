#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")" && pwd)"
BINARY="$ROOT/pep"

echo "  • Building pep..."
go build -o "$BINARY" "$ROOT"

SYMLINK_DIR="${HOME}/.local/bin"
if [ ! -d "$SYMLINK_DIR" ]; then
    mkdir -p "$SYMLINK_DIR"
fi

SYMLINK="$SYMLINK_DIR/pep"
if [ -L "$SYMLINK" ] || [ ! -e "$SYMLINK" ]; then
    ln -sf "$BINARY" "$SYMLINK"
    echo "    ✓ Symlinked to $SYMLINK"
else
    echo "    ! $SYMLINK already exists — not a symlink, skipping"
fi

VSCODE_EXT="${HOME}/.vscode/extensions/pep-lang.pep-lang"
if [ -d "$VSCODE_EXT" ]; then
    echo "  • Updating VS Code extension..."
else
    echo "  • Installing VS Code extension..."
fi
rm -rf "$VSCODE_EXT"
cp -r "$ROOT/vscode-pep" "$VSCODE_EXT"
echo "    ✓ Copied to $VSCODE_EXT"

echo ""
echo "  ✓ Pep installed. Reload VS Code to activate the extension."
echo "  ✓ Run 'pep --help' to get started."
echo ""
