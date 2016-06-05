FROM golang:1.6

RUN apt-get update && apt-get install -y --no-install-recommends \
        pkg-config \
        cmake \
    && rm -rf /var/lib/apt/lists/*

RUN git config --global user.email "you@example.com" \
    && git config --global user.name "Your Name"

RUN git clone --branch v0.24.0 https://github.com/libgit2/libgit2
RUN cd libgit2 && \
    mkdir build && cd build && \
    cmake .. && \
    cmake --build . --target install

# Without this it cannot find the shared library
ENV LD_LIBRARY_PATH=/usr/local/lib:$LD_LIBRARY_PATH

RUN mkdir -p /go/src/app
WORKDIR /go/src/app
COPY . /go/src/app

RUN go-wrapper download
RUN go-wrapper install

CMD ["go", "run", "main.go"]