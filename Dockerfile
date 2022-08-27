#Build Stage
FROM golang:1.16-alpine As build-env

WORKDIR /clovemind

COPY . .

RUN go build -o main main.go

# dev stage
FROM surnet/alpine-wkhtmltopdf:3.9-0.12.5-full as wkhtmltopdf
FROM openjdk:8-jdk-alpine3.9

WORKDIR /clovemind

EXPOSE 5000

COPY --from=build-env /clovemind/main .

COPY htmlfile/ ./htmlfile/

COPY .test-env .

# RUN apk add --no-cache wkhtmltopdf

RUN apk add --no-cache \
  libstdc++ \
  libx11 \
  libxrender \
  libxext \
  libssl1.1 \
  ca-certificates \
  fontconfig \
  freetype \
  ttf-dejavu \
  ttf-droid \
  ttf-freefont \
  ttf-liberation \
  ttf-ubuntu-font-family \
&& apk add --no-cache --virtual .build-deps \
  msttcorefonts-installer \
\
# Install microsoft fonts
&& update-ms-fonts \
&& fc-cache -f \
\
# Clean up when done
&& rm -rf /tmp/* \
&& apk del .build-deps

# Copy wkhtmltopdf files from docker-wkhtmltopdf image
COPY --from=wkhtmltopdf /bin/wkhtmltopdf /bin/wkhtmltopdf



CMD [ "/clovemind/main" ]