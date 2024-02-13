#!/bin/bash

# Nome da imagem Docker
IMAGE_NAME="rafaelleoni/rinha-de-backend-2024-q1"

# Versão da imagem
VERSION="latest"

# Construir a imagem Docker
docker buildx build --platform linux/amd64 -t $IMAGE_NAME:$VERSION .