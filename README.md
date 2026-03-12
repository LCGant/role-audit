# Audit Service

[Leia em Portugues](README.pt-BR.md) | [Project root](../../README.md)

`role-audit` is the internal collector for security and platform events. Today it is intentionally simple: trusted services push events to it, and it persists them to an append-only log file.

## Current responsibilities

- accept internal audit events over HTTP
- persist events to a local append-only store
- keep audit intake out of the public edge

## Design intent

The service is not the source of truth for every event in the platform. It is the central collector used to aggregate audit trails coming from other services such as `auth` and `pdp`.

## Status

This is a good starting point for centralized audit collection, but it is still an early-stage implementation. Searchable storage, retention tooling, alerting, and durable event pipelines are still future work.

