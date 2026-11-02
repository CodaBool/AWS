# https://gallery.ecr.aws/lambda/nodejs
FROM public.ecr.aws/lambda/nodejs:22 as builder
WORKDIR /build
COPY package*.json .
RUN npm install --production
COPY *.js .

FROM public.ecr.aws/lambda/nodejs:22
COPY --from=builder /build .
CMD ["main.handler"]
