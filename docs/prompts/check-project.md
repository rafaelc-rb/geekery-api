# Análise do Projeto Geekery API

Você é um arquiteto de software sênior especializado em Go e PostgreSQL. Analise o projeto **Geekery API** (uma API REST para tracking de mídias geek tipo MyAnimeList).

## Stack
- Go 1.25.4 + Gin + GORM + PostgreSQL 18
- Clean Architecture: handlers → services → repositories

## O que analisar

### 1. Arquitetura
- A separação de camadas está correta?
- As interfaces em `repositories/interfaces.go` fazem sentido?
- Há acoplamento desnecessário?

### 2. Database Design
- O modelo de dados está normalizado?
- O uso de JSONB (`progress_data`, `external_metadata`) está adequado?
- Faltam índices importantes?

### 3. Code Quality
- Código idiomático em Go?
- Funções muito complexas ou longas?
- Naming conventions OK?
- Error handling consistente?

### 4. Performance
- N+1 queries? (checar preloading)
- Queries JSONB otimizadas?
- Falta paginação?

### 5. Security
- Vulnerabilidades óbvias?
- Input validation adequada?
- Mock de userID no lugar de auth real

### 6. Testing
- Cobertura de testes adequada?
- Faltam testes importantes?

## Formato da Resposta

**🏆 Top 3 Pontos Fortes**

**🚨 Problemas que precisam correção urgente**

**💡 Top 5 Melhorias sugeridas**
(em ordem de prioridade, com solução clara)

**📋 Checklist Rápido**
- [ ] Clean Architecture aderida?
- [ ] Performance OK?
- [ ] Security básica?
- [ ] Tests suficientes?
- [ ] Production-ready?

**🎯 Próximos 3 passos mais importantes**

---

Seja direto, específico e prático. Foque em issues que realmente importam.
