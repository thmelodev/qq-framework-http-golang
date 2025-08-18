# qq-framework-http-golang

Framework HTTP para aplicações Go, integrando health checks, logging, tratamento de erros e modularização via Uber FX. Projetado para uso corporativo, com foco em observabilidade, resiliência e integração com outros módulos do ecossistema qq-framework.

## Visão Geral dos Componentes

### 1. Configuração HTTP (`config.go`)

Define a interface `IHttpProvider`, que abstrai configurações essenciais do servidor HTTP:

- Porta do servidor
- Limite de paginação
- Limite de rotinas concorrentes

### 2. Tratamento de Erros (`errorHandler.go`)

Middleware para Echo que padroniza respostas de erro, mapeando exceções do framework para códigos HTTP adequados:

- 400 para erros de validação/negócio
- 404 para dados não encontrados
- 500 para erros internos

### 3. Health Checks (`health.go`, `healthHandler.go`, `health.module.go`)

Infraestrutura robusta para health checks de serviços críticos:

- **Kafka**: Verifica conectividade e status dos brokers
- **Postgres**: Testa conexão e ping ao banco
- **Redis**: Testa conexão ao cache

O handler retorna status detalhado de cada serviço, facilitando integração com sistemas de monitoramento.

### 4. Módulo Health (`health.module.go`)

Exporta o health check como módulo FX, facilitando injeção e composição em aplicações modulares.

### 5. Módulo HTTP (`module.go`)

Exporta o servidor HTTP como módulo FX, permitindo inicialização e shutdown automáticos via ciclo de vida do FX.

### 6. Servidor HTTP (`server.go`)

- Inicializa o Echo
- Aplica middleware de logging
- Cria grupo de rotas baseado no nome da aplicação
- Expõe endpoints `/alive` (simples) e `/health` (detalhado, se health checks estiverem configurados)
- Gerencia ciclo de vida do servidor via FX

## Integração com Outros Módulos

O framework depende de outros pacotes do ecossistema qq-framework:

- `qq-framework-basic-golang`: Exceções e utilitários
- `qq-framework-basic-kafka`: Configuração e integração Kafka
- `qq-framework-db-golang`: Integração com banco de dados
- `qq-framework-log-golang`: Logging estruturado

## Exemplo de Uso

```go
import (
    "github.com/qq-mercantil/qq-framework-http-golang/http"
    "go.uber.org/fx"
)

func main() {
    app := fx.New(
        http.ServerModule(),
        http.HealthModule(),
        // outros módulos...
    )
    app.Run()
}
```

## Endpoints Padrão

- `GET /<APP_NAME>/alive`: Verifica se o serviço está rodando
- `GET /<APP_NAME>/health`: Health check detalhado (Kafka, Postgres, Redis)

## Licença

MIT

---

Este README foi gerado automaticamente a partir do código-fonte para garantir precisão e detalhamento técnico.
