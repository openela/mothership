FROM node:20-slim AS base
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable
WORKDIR /app
ADD package.json /app/package.json
ADD pnpm-lock.yaml /app/pnpm-lock.yaml

FROM base AS install
RUN mkdir -p /temp/dev
COPY package.json pnpm-lock.yaml /temp/dev/
RUN --mount=type=cache,id=pnpm,target=/pnpm/store cd /temp/dev && pnpm install --frozen-lockfile

RUN mkdir -p /temp/prod
COPY package.json pnpm-lock.yaml /temp/prod/
RUN --mount=type=cache,id=pnpm,target=/pnpm/store cd /temp/prod && pnpm install --frozen-lockfile --production

FROM base AS build
ADD . .
COPY --from=install /temp/dev/node_modules node_modules
ENV NODE_ENV=development
RUN pnpm run build

FROM cgr.dev/chainguard/bun
WORKDIR /app
COPY --from=install /temp/prod/node_modules node_modules
COPY --from=build /app/dist dist
COPY server.ts server.ts
COPY favicon.png favicon.png
ENV NODE_ENV=production
CMD ["run", "server.ts"]