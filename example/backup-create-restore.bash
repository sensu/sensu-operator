#!/bin/bash

SENSU_OPERATOR_ROOT=${SENSU_OPERATOR_ROOT:-$(dirname "${BASH_SOURCE[0]}")/..}

print_usage() {
  echo "$(basename "$0") - Create SensuCluster backups
Usage: $(basename "$0") [options...]
Options:
  --cluster-name=STRING    Name of SensuCluster to backup
                             (default=\"example-sensu-cluster\", environment variable: CLUSTER_NAME)
  --backup-name=STRING     Name for new SensuClusterBackup
                             (default=\"CLUSTER_NAME-backup-EPOCH\", environment variable: BACKUP_NAME)
  --aws-secret-name=STRING Name of the AWS secret to use
                             (default=\"sensu-backups-aws-secret\", environment variable: AWS_SECRET_NAME)
  --aws-bucket-name=STRING Name of the AWS bucket to use for backups
                             (default=\"sensu-backups\", environment variable: AWS_BUCKET_NAME)
" >&2
}

CLUSTER_NAME="${CLUSTER_NAME:-example-sensu-cluster}"
BACKUP_NAME="${BACKUP_NAME:-${CLUSTER_NAME}}"
AWS_SECRET_NAME="${AWS_SECRET_NAME:-sensu-backups-aws-secret}"
AWS_BUCKET_NAME="${AWS_BUCKET_NAME:-sensu-backups}"

parse_options() {
  for i in "$@"; do
    case $i in
    --cluster-name=*)
      CLUSTER_NAME="${i#*=}"
      ;;
    --backup-name=*)
      BACKUP_NAME="${i#*=}"
      ;;
    --aws-secret-name=*)
      AWS_SECRET_NAME="${i#*=}"
      ;;
    --aws-bucket-name=*)
      AWS_BUCKET_NAME="${i#*=}"
      ;;
    -h | --help)
      print_usage
      exit 0
      ;;
    *)
      print_usage
      exit 1
      ;;
    esac
  done
}

_do() {
  local action="${1}"
  echo "${action^} of cluster '${CLUSTER_NAME}' with backup named '${BACKUP_NAME}'"
  sed \
    -e "s/<CLUSTER_NAME>/${CLUSTER_NAME}/g" \
    -e "s/<BACKUP_NAME>/${BACKUP_NAME}/g" \
    -e "s/<AWS_SECRET_NAME>/${AWS_SECRET_NAME}/g" \
    -e "s/<AWS_BUCKET_NAME>/${AWS_BUCKET_NAME}/g" \
    "${SENSU_OPERATOR_ROOT}/example/${action}-operator/${action}-template.yaml" |
    kubectl create -f -
}

backup() {
  _do backup
}

restore() {
  _do restore
}
