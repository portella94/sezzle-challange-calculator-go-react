# --- Build stage ---------------------------------------------------------
FROM node:24.19-alpine AS build

WORKDIR /app

# Install dependencies from the lockfile for reproducible builds.
# --legacy-peer-deps matches how the lockfile was produced and avoids an
# optional-peer resolution edge case on some npm versions.
COPY package.json package-lock.json ./
RUN npm ci --legacy-peer-deps

COPY . .
RUN npm run build

# --- Runtime stage -------------------------------------------------------
FROM nginx:alpine

# Static assets and the reverse-proxy config that forwards /api to the backend.
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
