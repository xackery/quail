FROM neilotoole/xcgo:latest
RUN gpg --keyserver keyserver.ubuntu.com --recv-keys 16126D3A3E5C1192 && apt update && apt install -y mesa-common-dev libgl1-mesa-dev libglu1-mesa-dev libc6-dev libglu1-mesa-dev libgl1-mesa-dev libxcursor-dev libxi-dev libxinerama-dev libxrandr-dev libxxf86vm-dev libasound2-dev pkg-config

WORKDIR /src