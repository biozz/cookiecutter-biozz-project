#!/bin/bash

task deps-frontend
task build-frontend
task build

cp .env.example .env

echo ""
echo "You are all set!"
echo "The server is going to run in mock mode. To enable real storage edit cmd/server.go."
echo ""
echo " Next steps:"
echo " - cd {{ cookiecutter.slug }}"
echo " - git init"
echo " - task dev"
