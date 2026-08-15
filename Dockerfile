# ローカル開発用のイメージ。Lambda 向けの bootstrap をビルドする本番用マルチステージ
# ビルドはここでは扱わない（設計書 §9 第2段階で追加する）。
#
# alpine を採用しないのは musl libc であるため。本番の Lambda ランタイムは
# provided.al2023（glibc）であり、ローカルと libc を揃えたほうが切り分けが容易になる。
FROM golang:1.26.6-bookworm

WORKDIR /app

# 開発中はホストのソースをバインドマウントするため、ここでは COPY しない。
# .git をマウントしない構成でも go build が通るよう VCS スタンプを無効化する。
ENV GOFLAGS=-buildvcs=false

# compose から常駐させて exec で使う。
CMD ["sleep", "infinity"]
