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

The list starts with a special **`<未選択 / unset current-context>`** entry, which is the default cursor position. Selecting it clears the current context (equivalent to `kubectl config unset current-context`), leaving no context selected. The entry matching the active context is annotated with `(current)`.

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

### Unset Context
Run interactive mode and choose the `<未選択 / unset current-context>` entry (selected by default):
```bash
kubec
```
This removes the `current-context` key from your kubeconfig so that no context is selected. All other kubeconfig fields are preserved.

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
