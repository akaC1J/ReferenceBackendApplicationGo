#!/bin/bash

REPO_URL=$1
BRANCH=$2
SPARSE_PATH=$3
TARGET_DIR=$4

git clone -b $BRANCH --single-branch -n --depth=1 --filter=tree:0 $REPO_URL $TARGET_DIR && \
cd $TARGET_DIR && \
git sparse-checkout set --no-cone $SPARSE_PATH && \
git checkout

# Перемещение файлов в нужную директорию
mkdir -p $(dirname $SPARSE_PATH)
mv $TARGET_DIR/$SPARSE_PATH $(dirname $SPARSE_PATH)

# Удаление временных файлов
rm -rf $TARGET_DIR
