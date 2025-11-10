# 🌱 Database Seeding

Este diretório contém os scripts para popular o banco de dados com dados de exemplo.

## 📁 Estrutura

```
cmd/seed/
├── main.go              # Arquivo principal - orquestra o processo de seed
├── README.md            # Este arquivo
├── data/                # Dados para seed organizados por tipo
│   ├── tags.go         # Tags do sistema
│   ├── items.go        # Items do catálogo global
│   └── user_items.go   # Items na lista pessoal dos usuários
└── seeders/            # Funções auxiliares para seed
    └── base.go         # Helpers reutilizáveis
```

## 🚀 Como Usar

### Executar Seed Completo

```bash
# Via Makefile (recomendado)
make seed

# Ou diretamente com Go
go run cmd/seed/main.go
```

### O que é Criado

1. **Usuário Demo**
   - Email: `demo@geekery.com`
   - Password: `password123`
   - Username: `demo`

2. **Tags** (12 tags)
   - Action, Adventure, Drama, Comedy, Fantasy, Sci-Fi, Shounen, Seinen, Romance, Mystery, Horror, Slice of Life

3. **Items do Catálogo** (34 items)
   - **Anime (3)**: Naruto, Black Clover, Death Note
   - **Movies (9)**: Star Wars Saga completa (Episodes I-IX)
   - **Series (3)**: Supernatural, The Big Bang Theory, How I Met Your Mother
   - **Games (5)**: Skyrim, Half-Life: Alyx, Dota 2, Counter-Strike 2, Age of Mythology
   - **Comics (5)**: Gantz (manga), Beck (manga), Solo Leveling (manhwa), Tensei shitara Slime Datta Ken (manga), Blue Lock (manga)
   - **Novels (4)**: The Beginning After The End (web_novel), Sword Art Online (light_novel), Mushoku Tensei (light_novel), Hai to Gensou no Grimgar (light_novel)
   - **Books (5)**: A Song of Ice and Fire series - A Game of Thrones, A Clash of Kings, A Storm of Swords, A Feast for Crows, A Dance with Dragons

4. **Lista Pessoal do Usuário**
   - Inicialmente vazia - usuário pode adicionar items do catálogo
   - Sistema preparado para tracking de progresso de todos os tipos de mídia

## 📝 Adicionar Novos Dados

### Adicionar Nova Tag

Edite `data/tags.go`:

```go
func GetTags() []models.Tag {
    return []models.Tag{
        // ... tags existentes
        {Name: "Nova Tag"},
    }
}
```

### Adicionar Novo Item ao Catálogo

Edite `data/items.go`:

```go
{
    Item: models.Item{
        Title:       "Novo Item",
        Type:        models.MediaTypeAnime,
        Description: "Descrição do item",
        ReleaseDate: seeders.ParseDate("2024-01-01"),
        CoverURL:    "https://...",
    },
    SpecificData: &models.AnimeData{
        Episodes: 24,
        Studio:   "Studio Name",
    },
    TagNames: []string{"Action", "Adventure"},
},
```

### Adicionar Item à Lista de Usuário

Edite `data/user_items.go`:

```go
{
    UserID:          userID,
    ItemID:          itemIDs["Nome do Item"],
    Status:          models.StatusInProgress,
    Rating:          8.5,
    Favorite:        true,
    Notes:           "Minhas notas",
    ProgressType:    models.ProgressTypeEpisodic,
    CompletionCount: 0,
    ProgressData: models.JSONB{
        "season":  1,
        "episode": 12,
        "history": []interface{}{
            map[string]interface{}{
                "started_at":  "2024-01-01T00:00:00Z",
                "finished_at": nil,
            },
        },
    },
},
```

## 🔧 Funções Auxiliares

### `seeders.ParseDate(dateStr string)`
Converte string no formato `YYYY-MM-DD` para `*time.Time`.

```go
releaseDate: seeders.ParseDate("2024-01-15")
```

### `seeders.CreateOrSkip(db, model, where)`
Cria um registro ou pula se já existir. Útil para evitar duplicatas.

```go
created, err := seeders.CreateOrSkip(db, &tag, models.Tag{Name: tag.Name})
```

### `seeders.CreateItemWithSpecificData(db, itemData, allTags)`
Cria um item completo com dados específicos e tags associadas.

```go
err := seeders.CreateItemWithSpecificData(db, itemData, tags)
```

### `seeders.CreateUserItemIfNotExists(db, userItem)`
Cria um user item se não existir para aquele usuário/item.

```go
err := seeders.CreateUserItemIfNotExists(db, &userItem)
```

## 🎯 Vantagens da Nova Estrutura

✅ **Modular**: Dados separados por tipo em arquivos diferentes
✅ **Manutenível**: Fácil adicionar/editar dados sem mexer na lógica
✅ **Reutilizável**: Helpers podem ser usados em outros scripts
✅ **Limpo**: main.go com apenas 170 linhas vs 383 linhas antigas
✅ **Testável**: Cada componente pode ser testado isoladamente
✅ **Idiomático**: Segue padrões de organização Go

## 📚 Exemplos de ProgressData por Tipo

### Episodic (Anime/Series)
```go
ProgressData: models.JSONB{
    "season":  2,
    "episode": 15,
    "history": []interface{}{...},
}
```

### Reading (Comics/Novels/Books)
```go
// Para Comics (manga, manhwa)
ProgressData: models.JSONB{
    "chapter": 450,
    "volume":  45,
    "history": []interface{}{...},
}

// Para Novels (light_novel, web_novel)
ProgressData: models.JSONB{
    "chapter": 200,
    "volume":  20,
    "history": []interface{}{...},
}

// Para Books tradicionais
ProgressData: models.JSONB{
    "page":    250,
    "history": []interface{}{...},
}
```

### Time (Movies)
```go
ProgressData: models.JSONB{
    "minutes_watched": 90,
    "last_position":   90,
    "history": []interface{}{...},
}
```

### Percent (Games)
```go
ProgressData: models.JSONB{
    "percent": 75,
    "hours":   120,
    "extras": map[string]interface{}{
        "achievements": 85,
    },
    "history": []interface{}{...},
}
```

## ⚠️ Notas Importantes

- Os dados são criados com **idempotência** - executar múltiplas vezes não cria duplicatas
- IDs são determinados dinamicamente - não hardcode IDs
- Use `itemIDs` map para referenciar items por título
- O seed é seguro para rodar em qualquer ambiente (dev/staging)
- Password do usuário demo é hashada com bcrypt

## 🔄 Reset Database

Para limpar e recriar tudo:

```bash
make db-reset
make seed
```
