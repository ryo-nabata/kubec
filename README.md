# kubec

A command-line tool for easily switching Kubernetes current-context.

## Features

- **Interactive mode**: Select from available contexts using a menu
- **Direct specification**: Quickly switch by specifying context name directly
- **Current context display**: Check the currently active context

## Installation

```bash
git clone <repository-url>
cd kubec
go build -o kubec
```

## Usage

### Interactive Mode
```bash
kubec
```
A list of available contexts will be displayed and you can select using arrow keys.

### Direct Specification
```bash
kubec my-cluster
```
Switch directly to the specified context.

### Show Current Context
```bash
kubec --current
# or
kubec -c
```
Display the currently active context.

## Prerequisites

- Access to a Kubernetes cluster environment
- `~/.kube/config` file must exist
- Multiple contexts must be configured

## Configuration File

kubec uses standard Kubernetes configuration files:
- Default: `~/.kube/config`
- If `KUBECONFIG` environment variable is set, it takes priority

## Reference

This tool is inspired by the implementation of [awsd](https://github.com/radiusmethod/awsd).
