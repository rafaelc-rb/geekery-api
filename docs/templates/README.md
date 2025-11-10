# CSV Import Templates

Templates para importação em massa de items via CSV.

## 🚀 Quick Start

```bash
# Importar anime
curl -X POST http://localhost:8080/api/items/import/anime -F "file=@anime.csv"

# Importar comics
curl -X POST http://localhost:8080/api/items/import/comic -F "file=@comic.csv"

# Importar novels
curl -X POST http://localhost:8080/api/items/import/novel -F "file=@novel.csv"
```

## 📋 Templates e Endpoints

### 🎌 Anime
**Endpoint:** `POST /api/items/import/anime`

| Campo               | Obrigatório | Formato              | Exemplo                 |
| ------------------- | ----------- | -------------------- | ----------------------- |
| `title`             | ✅ Sim       | String               | Attack on Titan         |
| `description`       | ❌ Não       | String               | Epic story about...     |
| `release_date`      | ❌ Não       | YYYY-MM-DD           | 2013-04-07              |
| `cover_url`         | ❌ Não       | URL                  | https://...             |
| `tags`              | ❌ Não       | tag1\|tag2\|tag3     | action\|shounen\|drama  |
| `episodes`          | ✅ Sim       | Number               | 75                      |
| `studio`            | ✅ Sim       | String               | MAPPA                   |
| `external_metadata` | ❌ Não       | source:id\|source:id | mal:16498\|anilist:1234 |

### 📖 Comic (Manga, Manhwa, Webtoon)
**Endpoint:** `POST /api/items/import/comic`

| Campo               | Obrigatório | Formato              | Exemplo                    |
| ------------------- | ----------- | -------------------- | -------------------------- |
| `title`             | ✅ Sim       | String               | One Piece                  |
| `description`       | ❌ Não       | String               | Pirates search for...      |
| `release_date`      | ❌ Não       | YYYY-MM-DD           | 1997-07-22                 |
| `cover_url`         | ❌ Não       | URL                  | https://...                |
| `tags`              | ❌ Não       | tag1\|tag2           | action\|adventure\|shounen |
| `chapters`          | ✅ Sim       | Number               | 1100                       |
| `volumes`           | ❌ Não       | Number               | 108                        |
| `author`            | ✅ Sim       | String               | Eiichiro Oda               |
| `format`            | ✅ Sim       | String               | manga / manhwa / webtoon   |
| `publisher`         | ❌ Não       | String               | Shueisha                   |
| `external_metadata` | ❌ Não       | source:id\|source:id | mal:13\|anilist:30013      |

### 📚 Novel (Light Novel, Web Novel)
**Endpoint:** `POST /api/items/import/novel`

| Campo               | Obrigatório | Formato              | Exemplo                      |
| ------------------- | ----------- | -------------------- | ---------------------------- |
| `title`             | ✅ Sim       | String               | Sword Art Online             |
| `description`       | ❌ Não       | String               | VRMMORPG adventure...        |
| `release_date`      | ❌ Não       | YYYY-MM-DD           | 2009-04-10                   |
| `cover_url`         | ❌ Não       | URL                  | https://...                  |
| `tags`              | ❌ Não       | tag1\|tag2           | action\|sci-fi\|fantasy      |
| `volumes`           | ❌ Não       | Number               | 28                           |
| `chapters`          | ❌ Não       | Number               | 500                          |
| `author`            | ✅ Sim       | String               | Reki Kawahara                |
| `format`            | ✅ Sim       | String               | light_novel / web_novel      |
| `publisher`         | ❌ Não       | String               | ASCII Media Works            |
| `external_metadata` | ❌ Não       | source:id\|source:id | mal:21479\|anilist:21479     |

### 🎬 Movie
**Endpoint:** `POST /api/items/import/movie`

| Campo               | Obrigatório | Formato              | Exemplo                  |
| ------------------- | ----------- | -------------------- | ------------------------ |
| `title`             | ✅ Sim       | String               | The Matrix               |
| `description`       | ❌ Não       | String               | A hacker discovers...    |
| `release_date`      | ❌ Não       | YYYY-MM-DD           | 1999-03-31               |
| `cover_url`         | ❌ Não       | URL                  | https://...              |
| `tags`              | ❌ Não       | tag1\|tag2           | action\|sci-fi           |
| `runtime`           | ✅ Sim       | Number (minutes)     | 136                      |
| `director`          | ✅ Sim       | String               | Wachowski                |
| `external_metadata` | ❌ Não       | source:id\|source:id | imdb:tt0133093\|tmdb:603 |

### 📺 Series
**Endpoint:** `POST /api/items/import/series`

| Campo               | Obrigatório | Formato              | Exemplo                   |
| ------------------- | ----------- | -------------------- | ------------------------- |
| `title`             | ✅ Sim       | String               | Breaking Bad              |
| `description`       | ❌ Não       | String               | A chemistry teacher...    |
| `release_date`      | ❌ Não       | YYYY-MM-DD           | 2008-01-20                |
| `cover_url`         | ❌ Não       | URL                  | https://...               |
| `tags`              | ❌ Não       | tag1\|tag2           | drama\|crime\|thriller    |
| `seasons`           | ✅ Sim       | Number               | 5                         |
| `episodes`          | ✅ Sim       | Number               | 62                        |
| `network`           | ❌ Não       | String               | AMC                       |
| `external_metadata` | ❌ Não       | source:id\|source:id | imdb:tt0903747\|tmdb:1396 |

### 🎮 Game
**Endpoint:** `POST /api/items/import/game`

| Campo               | Obrigatório | Formato              | Exemplo                    |
| ------------------- | ----------- | -------------------- | -------------------------- |
| `title`             | ✅ Sim       | String               | Elden Ring                 |
| `description`       | ❌ Não       | String               | Action RPG...              |
| `release_date`      | ❌ Não       | YYYY-MM-DD           | 2022-02-25                 |
| `cover_url`         | ❌ Não       | URL                  | https://...                |
| `tags`              | ❌ Não       | tag1\|tag2           | action\|rpg\|souls-like    |
| `platform`          | ✅ Sim       | String               | PC\|PS5\|Xbox Series X     |
| `developer`         | ✅ Sim       | String               | FromSoftware               |
| `publisher`         | ❌ Não       | String               | Bandai Namco               |
| `external_metadata` | ❌ Não       | source:id\|source:id | igdb:119133\|steam:1245620 |

### 📘 Book (Traditional Books)
**Endpoint:** `POST /api/items/import/book`

| Campo               | Obrigatório | Formato              | Exemplo                            |
| ------------------- | ----------- | -------------------- | ---------------------------------- |
| `title`             | ✅ Sim       | String               | The Hobbit                         |
| `description`       | ❌ Não       | String               | Bilbo's adventure...               |
| `release_date`      | ❌ Não       | YYYY-MM-DD           | 1937-09-21                         |
| `cover_url`         | ❌ Não       | URL                  | https://...                        |
| `tags`              | ❌ Não       | tag1\|tag2           | fantasy\|adventure                 |
| `pages`             | ✅ Sim       | Number               | 310                                |
| `author`            | ✅ Sim       | String               | J.R.R. Tolkien                     |
| `publisher`         | ❌ Não       | String               | George Allen & Unwin               |
| `external_metadata` | ❌ Não       | source:id\|source:id | isbn:9780547928227\|goodreads:5907 |

## 📌 Notas Importantes

### External Metadata

Formato: `source:id|source:id` (ex: `mal:123|anilist:456`)

**Sources comuns:**
- Anime/Comic/Novel: `mal`, `anilist`, `kitsu`
- Movies/Series: `imdb`, `tmdb`
- Games: `igdb`, `steam`, `gog`
- Books: `isbn`, `goodreads`

### Tags
Separadas por `|` (ex: `action|adventure|fantasy`). Criadas automaticamente se não existirem.

### Datas
Formato: `YYYY-MM-DD` (ex: `2023-01-15`)

### Validações
- **Duplicatas**: Items com mesmo título e tipo são bloqueados
- **Campos obrigatórios**: Veja tabelas acima
- **Números**: Devem ser inteiros positivos
- **Encoding**: Use UTF-8

## 📊 Response de Exemplo

```json
{
  "success": true,
  "media_type": "anime",
  "total_lines": 10,
  "imported": 8,
  "failed": 2,
  "errors": [
    {
      "line": 5,
      "title": "Duplicate Anime",
      "error": "item 'Duplicate Anime' (type: anime) already exists with ID 3"
    },
    {
      "line": 8,
      "title": "Invalid Item",
      "error": "episodes must be a number"
    }
  ]
}
```

## 🧪 Teste

```bash
# Teste automatizado
./test-import.sh

# Ou manual
curl -X POST http://localhost:8080/api/items/import/anime -F "file=@anime.csv" | jq
```
