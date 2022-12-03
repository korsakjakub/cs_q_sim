#!/bin/bash
printf "source with: . config/export_envs.sh\n"
export $(grep -v ^# config/.env | xargs);