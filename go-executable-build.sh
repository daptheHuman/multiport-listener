#!/usr/bin/env bash

tag_name=$1
platforms=("windows/amd64" "linux/amd64" "darwin/arm64")
current_arch=$(uname -m)

echo "Current architecture is $current_arch"

# Build
for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}

    if [ $GOARCH = "amd64" ] && [ $GOOS = "linux" ]; then
        if [ $current_arch = "arm64" ]; then
            echo "Skipping $GOOS/$GOARCH build on $current_arch"
            continue
        fi
    fi
  
    output_name='multiport-listener-'$GOOS'-'$GOARCH
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi
    echo "Building to release/$output_name..."
    env GOOS=$GOOS GOARCH=$GOARCH go build -o release/$output_name
    done

