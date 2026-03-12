# Servico de Audit

[Read in English](README.md) | [Raiz do projeto](../../README.pt-BR.md)

`role-audit` e o coletor interno de eventos de seguranca e plataforma. Hoje ele e simples de forma intencional: servicos confiaveis enviam eventos para ele, e o servico persiste esses eventos em um log append-only.

## Responsabilidades atuais

- aceitar eventos internos de auditoria via HTTP
- persistir eventos em armazenamento local append-only
- manter a entrada de auditoria fora da borda publica

## Intencao de desenho

Este servico nao e a fonte de verdade de todos os eventos da plataforma. Ele e o coletor central usado para agregar trilhas de auditoria vindas de outros servicos, como `auth` e `pdp`.

## Estado atual

Este modulo ja e um bom ponto de partida para auditoria centralizada, mas ainda esta em fase inicial. Storage consultavel, retencao, alertas e pipeline duravel de eventos continuam como trabalho futuro.

