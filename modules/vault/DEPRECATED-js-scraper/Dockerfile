FROM public.ecr.aws/lambda/nodejs:18 as builder
WORKDIR /build
COPY package*.json .
RUN npm install --production
COPY main.js .

FROM public.ecr.aws/lambda/nodejs:18
COPY --from=builder /build .
CMD ["main.handler"]