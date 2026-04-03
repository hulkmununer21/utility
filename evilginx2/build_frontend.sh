#!/bin/bash
# Build the frontend React app
cd "$(dirname "$0")/frontend"

echo "Installing dependencies..."
npm install

echo "Building React app..."
npm run build

echo "Build complete! Frontend files are in frontend/dist/"
