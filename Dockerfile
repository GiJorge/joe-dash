# --- STAGE 1: Build Frontend ---
FROM node:20-slim AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ .
RUN npm run build

# --- STAGE 2: Build Backend ---
FROM golang:1.26-alpine AS backend-builder
WORKDIR /app/backend
# Install build dependencies for SQLite
RUN apk add --no-cache gcc musl-dev sqlite-dev
#RUN apk add --no-cache gcc musl-dev
COPY backend/go.mod ./
# If you have a go.sum, copy it too: COPY backend/go.sum ./ 
RUN go mod download
COPY backend/ .

# CRITICAL: CGO_ENABLED=1 is required for SQLite
RUN CGO_ENABLED=1 GOOS=linux go build -o main .

#RUN go build -o main .

# --- STAGE 3: Final Production Image ---
FROM alpine:latest
# libc6-compat is often needed for CGO binaries on Alpine
RUN apk add --no-cache ca-certificates libc6-compat
#RUN apk add --no-cache ca-certificates
#WORKDIR /root/
WORKDIR /app/backend

# Copy Go binary from Stage 2
COPY --from=backend-builder /app/backend/main .

# Copy Vue static files from Stage 1
# This puts the "dist" folder right next to the "main" binary
COPY --from=frontend-builder /app/frontend/dist ./dist

# Create data directory for SQLite
RUN mkdir ./data

EXPOSE 8080
CMD ["./main"]








