# Kubernetes Time Travel Debugger (KTTD)

KTTD is a CLI tool that snapshots your Kubernetes cluster's state over time, lets you diff between historical states, and optionally replay any past configuration — like Git for your cluster.

## Features
- 📸 Snapshot Deployments, Pods, Services, ConfigMaps, and more
- 🕵️ Diff any two snapshots and see what changed
- 🔁 Replay past state (coming soon)
- 🔐 Works with RBAC and supports multi-namespace

## Quick Start

```bash
git clone https://github.com/your-username/kttd.git
cd kttd
go mod tidy
go run main.go snapshot --namespace=default
