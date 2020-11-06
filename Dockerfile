FROM golang:1.15-buster

RUN apt-get update -y && apt-get install -y \
    pdf2svg texlive-latex-base \
    texlive-pictures \
    texlive-fonts-recommended xzdec \
    texlive-fonts-extra \
    && rm -rf /var/lib/apt/lists/*

# RUN tlmgr init-usertree \
#  && tlmgr option repository ftp://tug.org/historic/systems/texlive/2018/tlnet-final \
#  && tlmgr install chemformula

# RUN updmap-sys
RUN mkdir -p /usr/src/app/tmpfs

VOLUME /usr/src/app/tmpfs
# RUN mount -t tmpfs swap /usr/src/app/tmpfs

WORKDIR /usr/src/app

COPY . /usr/src/app

EXPOSE 3000

ENV TEXMFHOME /root/texmf

RUN go build

CMD ["go", "run", "main.go"]