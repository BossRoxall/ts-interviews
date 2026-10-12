# Node runs the TypeScript sources directly (type stripping), and the service
# has no runtime dependencies, so there is no build or install step.
FROM node:26.13.0-slim

ENV NODE_ENV=production

WORKDIR /app
COPY package.json ./
COPY src ./src

EXPOSE 8080
USER node
ENTRYPOINT ["node", "src/main.ts"]
