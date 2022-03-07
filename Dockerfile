FROM golang:1.17 as base

FROM base as dev

WORKDIR /queryhw
CMD ["bash"]