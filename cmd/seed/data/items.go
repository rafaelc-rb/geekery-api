package data

import (
	"github.com/rafaelc-rb/geekery-api/cmd/seed/seeders"
	"github.com/rafaelc-rb/geekery-api/internal/models"
)

// GetCatalogItems retorna todos os items do catálogo para seed
func GetCatalogItems() []seeders.ItemWithData {
	return []seeders.ItemWithData{
		// ========== ANIME ==========
		{
			Item: models.Item{
				Title:       "Naruto",
				Type:        models.MediaTypeAnime,
				Description: "Naruto Uzumaki, a young ninja with a sealed demon within him, strives to become the strongest ninja and leader of his village.",
				ReleaseDate: seeders.ParseDate("2002-10-03"),
				CoverURL:    "https://cdn.myanimelist.net/images/anime/13/17405.jpg",
			},
			SpecificData: &models.AnimeData{
				Episodes: 220,
				Studio:   "Studio Pierrot",
			},
			TagNames: []string{"Action", "Adventure", "Shounen"},
		},
		{
			Item: models.Item{
				Title:       "Black Clover",
				Type:        models.MediaTypeAnime,
				Description: "Asta and Yuno were abandoned together at the same church and have been inseparable since. As children, they promised that they would compete against each other to see who would become the next Emperor Magus.",
				ReleaseDate: seeders.ParseDate("2017-10-03"),
				CoverURL:    "https://cdn.myanimelist.net/images/anime/1/123486.jpg",
			},
			SpecificData: &models.AnimeData{
				Episodes: 170,
				Studio:   "Studio Pierrot",
			},
			TagNames: []string{"Action", "Fantasy", "Shounen"},
		},
		{
			Item: models.Item{
				Title:       "Death Note",
				Type:        models.MediaTypeAnime,
				Description: "Light Yagami encontra um caderno sobrenatural que permite matar qualquer pessoa cujo nome seja escrito nele.",
				ReleaseDate: seeders.ParseDate("2006-10-04"),
				CoverURL:    "https://cdn.myanimelist.net/images/anime/9/9453.jpg",
			},
			SpecificData: &models.AnimeData{
				Episodes: 37,
				Studio:   "Madhouse",
			},
			TagNames: []string{"Mystery", "Drama", "Shounen"},
		},

		// ========== MOVIES (Star Wars Saga) ==========
		{
			Item: models.Item{
				Title:       "Star Wars: Episode I - The Phantom Menace",
				Type:        models.MediaTypeMovie,
				Description: "Two Jedi escape a hostile blockade to find allies and come across a young boy who may bring balance to the Force.",
				ReleaseDate: seeders.ParseDate("1999-05-19"),
				CoverURL:    "https://image.tmdb.org/t/p/w500/6wkfovpn7Eq8dYNKaG5PY3q2oq6.jpg",
			},
			SpecificData: &models.MovieData{
				Director: "George Lucas",
				Runtime:  136,
			},
			TagNames: []string{"Action", "Adventure", "Sci-Fi"},
		},
		{
			Item: models.Item{
				Title:       "Star Wars: Episode II - Attack of the Clones",
				Type:        models.MediaTypeMovie,
				Description: "Ten years after initially meeting, Anakin Skywalker shares a forbidden romance with Padmé Amidala.",
				ReleaseDate: seeders.ParseDate("2002-05-16"),
				CoverURL:    "https://image.tmdb.org/t/p/w500/oZNPzxqM2s5DyVWab09NTQScDQt.jpg",
			},
			SpecificData: &models.MovieData{
				Director: "George Lucas",
				Runtime:  142,
			},
			TagNames: []string{"Action", "Adventure", "Sci-Fi"},
		},
		{
			Item: models.Item{
				Title:       "Star Wars: Episode III - Revenge of the Sith",
				Type:        models.MediaTypeMovie,
				Description: "Three years into the Clone Wars, the Jedi rescue Palpatine from Count Dooku. As Obi-Wan pursues a new threat, Anakin acts as a double agent.",
				ReleaseDate: seeders.ParseDate("2005-05-19"),
				CoverURL:    "https://image.tmdb.org/t/p/w500/xfSAoBEm9MNBjmlNcDYLvLSMlnq.jpg",
			},
			SpecificData: &models.MovieData{
				Director: "George Lucas",
				Runtime:  140,
			},
			TagNames: []string{"Action", "Adventure", "Sci-Fi"},
		},
		{
			Item: models.Item{
				Title:       "Star Wars: Episode IV - A New Hope",
				Type:        models.MediaTypeMovie,
				Description: "Luke Skywalker joins forces with a Jedi Knight, a cocky pilot, a Wookiee and two droids to save the galaxy from the Empire's world-destroying battle station.",
				ReleaseDate: seeders.ParseDate("1977-05-25"),
				CoverURL:    "https://image.tmdb.org/t/p/w500/6FfCtAuVAW8XJjZ7eWeLibRLWTw.jpg",
			},
			SpecificData: &models.MovieData{
				Director: "George Lucas",
				Runtime:  121,
			},
			TagNames: []string{"Action", "Adventure", "Sci-Fi"},
		},
		{
			Item: models.Item{
				Title:       "Star Wars: Episode V - The Empire Strikes Back",
				Type:        models.MediaTypeMovie,
				Description: "After the Rebels are brutally overpowered by the Empire on the ice planet Hoth, Luke Skywalker begins Jedi training with Yoda.",
				ReleaseDate: seeders.ParseDate("1980-05-21"),
				CoverURL:    "https://image.tmdb.org/t/p/w500/nNAeTmF4CtdSgMDplXTDPOpYzsX.jpg",
			},
			SpecificData: &models.MovieData{
				Director: "Irvin Kershner",
				Runtime:  124,
			},
			TagNames: []string{"Action", "Adventure", "Sci-Fi"},
		},
		{
			Item: models.Item{
				Title:       "Star Wars: Episode VI - Return of the Jedi",
				Type:        models.MediaTypeMovie,
				Description: "After a daring mission to rescue Han Solo from Jabba the Hutt, the Rebels dispatch to Endor to destroy the second Death Star.",
				ReleaseDate: seeders.ParseDate("1983-05-25"),
				CoverURL:    "https://image.tmdb.org/t/p/w500/jx5p0aHlbPXqe3AH9G15NvmWaqQ.jpg",
			},
			SpecificData: &models.MovieData{
				Director: "Richard Marquand",
				Runtime:  131,
			},
			TagNames: []string{"Action", "Adventure", "Sci-Fi"},
		},
		{
			Item: models.Item{
				Title:       "Star Wars: Episode VII - The Force Awakens",
				Type:        models.MediaTypeMovie,
				Description: "As a new threat to the galaxy rises, Rey, a desert scavenger, and Finn, an ex-stormtrooper, must join Han Solo and Chewbacca to search for the one hope of restoring peace.",
				ReleaseDate: seeders.ParseDate("2015-12-18"),
				CoverURL:    "https://image.tmdb.org/t/p/w500/wqnLdwVXoBjKibFRR5U3y0aDUhs.jpg",
			},
			SpecificData: &models.MovieData{
				Director: "J.J. Abrams",
				Runtime:  138,
			},
			TagNames: []string{"Action", "Adventure", "Sci-Fi"},
		},
		{
			Item: models.Item{
				Title:       "Star Wars: Episode VIII - The Last Jedi",
				Type:        models.MediaTypeMovie,
				Description: "Rey develops her newly discovered abilities with the guidance of Luke Skywalker, who is unsettled by the strength of her powers.",
				ReleaseDate: seeders.ParseDate("2017-12-15"),
				CoverURL:    "https://image.tmdb.org/t/p/w500/kOVEVeg59E0wsnXmF9nrh6OmWII.jpg",
			},
			SpecificData: &models.MovieData{
				Director: "Rian Johnson",
				Runtime:  152,
			},
			TagNames: []string{"Action", "Adventure", "Sci-Fi"},
		},
		{
			Item: models.Item{
				Title:       "Star Wars: Episode IX - The Rise of Skywalker",
				Type:        models.MediaTypeMovie,
				Description: "The surviving Resistance faces the First Order once more in the final chapter of the Skywalker saga.",
				ReleaseDate: seeders.ParseDate("2019-12-20"),
				CoverURL:    "https://image.tmdb.org/t/p/w500/db32LaOibwEliAmSL2jjDF6oDdj.jpg",
			},
			SpecificData: &models.MovieData{
				Director: "J.J. Abrams",
				Runtime:  142,
			},
			TagNames: []string{"Action", "Adventure", "Sci-Fi"},
		},

		// ========== SERIES ==========
		{
			Item: models.Item{
				Title:       "Supernatural",
				Type:        models.MediaTypeSeries,
				Description: "Two brothers follow their father's footsteps as hunters, fighting evil supernatural beings of many kinds.",
				ReleaseDate: seeders.ParseDate("2005-09-13"),
				CoverURL:    "https://image.tmdb.org/t/p/w500/KoYWXbnYuS3b0GyQPkbuexlVK9.jpg",
			},
			SpecificData: &models.SeriesData{
				Seasons:  15,
				Episodes: 327,
			},
			TagNames: []string{"Drama", "Fantasy", "Horror"},
		},
		{
			Item: models.Item{
				Title:       "The Big Bang Theory",
				Type:        models.MediaTypeSeries,
				Description: "A woman who moves into an apartment across the hall from two brilliant but socially awkward physicists shows them how little they know about life outside of the laboratory.",
				ReleaseDate: seeders.ParseDate("2007-09-24"),
				CoverURL:    "https://image.tmdb.org/t/p/w500/ooBGRQBdbGzBxAVfExiO8r7kloA.jpg",
			},
			SpecificData: &models.SeriesData{
				Seasons:  12,
				Episodes: 279,
			},
			TagNames: []string{"Comedy", "Romance"},
		},
		{
			Item: models.Item{
				Title:       "How I Met Your Mother",
				Type:        models.MediaTypeSeries,
				Description: "A father recounts to his children - through a series of flashbacks - the journey he and his four best friends took leading up to him meeting their mother.",
				ReleaseDate: seeders.ParseDate("2005-09-19"),
				CoverURL:    "https://image.tmdb.org/t/p/w500/b34jPzmB0wZy7EjUZoleXOl2RRI.jpg",
			},
			SpecificData: &models.SeriesData{
				Seasons:  9,
				Episodes: 208,
			},
			TagNames: []string{"Comedy", "Romance"},
		},

		// ========== GAMES ==========
		{
			Item: models.Item{
				Title:       "The Elder Scrolls V: Skyrim",
				Type:        models.MediaTypeGame,
				Description: "The open-world adventure from Bethesda Game Studios where you can virtually be anyone and do anything.",
				ReleaseDate: seeders.ParseDate("2011-11-11"),
				CoverURL:    "https://cdn.cloudflare.steamstatic.com/steam/apps/489830/header.jpg",
			},
			SpecificData: &models.GameData{
				Platform:        "PC, PS3, PS4, PS5, Xbox 360, Xbox One, Xbox Series X/S, Nintendo Switch",
				Developer:       "Bethesda Game Studios",
				AveragePlaytime: 100,
			},
			TagNames: []string{"Adventure", "Fantasy", "Action"},
		},
		{
			Item: models.Item{
				Title:       "Half-Life: Alyx",
				Type:        models.MediaTypeGame,
				Description: "Half-Life: Alyx is Valve's VR return to the Half-Life series. It's the story of an impossible fight against a vicious alien race known as the Combine.",
				ReleaseDate: seeders.ParseDate("2020-03-23"),
				CoverURL:    "https://cdn.cloudflare.steamstatic.com/steam/apps/546560/header.jpg",
			},
			SpecificData: &models.GameData{
				Platform:        "PC (VR)",
				Developer:       "Valve",
				AveragePlaytime: 15,
			},
			TagNames: []string{"Action", "Sci-Fi", "Adventure"},
		},
		{
			Item: models.Item{
				Title:       "Dota 2",
				Type:        models.MediaTypeGame,
				Description: "Every day, millions of players worldwide enter battle as one of over a hundred Dota heroes in a 5v5 team clash.",
				ReleaseDate: seeders.ParseDate("2013-07-09"),
				CoverURL:    "https://cdn.cloudflare.steamstatic.com/steam/apps/570/header.jpg",
			},
			SpecificData: &models.GameData{
				Platform:        "PC, Mac, Linux",
				Developer:       "Valve",
				AveragePlaytime: 0, // Multiplayer infinito
			},
			TagNames: []string{"Action", "Fantasy"},
		},
		{
			Item: models.Item{
				Title:       "Counter-Strike 2",
				Type:        models.MediaTypeGame,
				Description: "For over two decades, Counter-Strike has offered an elite competitive experience, one shaped by millions of players from across the globe.",
				ReleaseDate: seeders.ParseDate("2023-09-27"),
				CoverURL:    "https://cdn.cloudflare.steamstatic.com/steam/apps/730/header.jpg",
			},
			SpecificData: &models.GameData{
				Platform:        "PC",
				Developer:       "Valve",
				AveragePlaytime: 0, // Multiplayer infinito
			},
			TagNames: []string{"Action"},
		},
		{
			Item: models.Item{
				Title:       "Age of Mythology: Extended Edition",
				Type:        models.MediaTypeGame,
				Description: "The classic real time strategy game that transports players to a time when heroes did battle with monsters of legend and the gods intervened in the affairs of mortals.",
				ReleaseDate: seeders.ParseDate("2014-05-08"),
				CoverURL:    "https://cdn.cloudflare.steamstatic.com/steam/apps/266840/header.jpg",
			},
			SpecificData: &models.GameData{
				Platform:        "PC",
				Developer:       "SkyBox Labs, Ensemble Studios",
				AveragePlaytime: 30,
			},
			TagNames: []string{"Adventure", "Fantasy"},
		},

		// ========== COMICS (Manga, Manhwa) ==========
		{
			Item: models.Item{
				Title:       "Gantz",
				Type:        models.MediaTypeComic,
				Description: "Kei Kurono is killed in a train accident along with his childhood friend, Masaru Kato. They awaken in a room with a black sphere that gives them weapons and orders to kill aliens.",
				ReleaseDate: seeders.ParseDate("2000-07-13"),
				CoverURL:    "https://cdn.myanimelist.net/images/manga/3/54835.jpg",
			},
			SpecificData: &models.BookData{
				Author:   "Hiroya Oku",
				Volumes:  37,
				Chapters: 383,
				Format:   "manga",
			},
			TagNames: []string{"Action", "Sci-Fi", "Horror", "Seinen"},
		},
		{
			Item: models.Item{
				Title:       "Beck: Mongolian Chop Squad",
				Type:        models.MediaTypeComic,
				Description: "Yukio 'Koyuki' Tanaka is a regular 14-year-old Japanese boy who becomes involved with rock music when he saves an odd-looking dog, named Beck, from some kids.",
				ReleaseDate: seeders.ParseDate("1999-09-01"),
				CoverURL:    "https://cdn.myanimelist.net/images/manga/2/54031.jpg",
			},
			SpecificData: &models.BookData{
				Author:   "Harold Sakuishi",
				Volumes:  34,
				Chapters: 103,
				Format:   "manga",
			},
			TagNames: []string{"Drama", "Comedy", "Slice of Life", "Seinen"},
		},
		{
			Item: models.Item{
				Title:       "Solo Leveling",
				Type:        models.MediaTypeComic,
				Description: "Jinwoo Sung, o mais fraco dos hunters, recebe um poder misterioso que o transforma.",
				ReleaseDate: seeders.ParseDate("2018-03-04"),
				CoverURL:    "https://image.tmdb.org/t/p/w500/gUJoJFBnqRdXWjh3TqLjgC5QWCQ.jpg",
			},
			SpecificData: &models.BookData{
				Author:   "Chugong",
				Chapters: 179,
				Format:   "manhwa",
			},
			TagNames: []string{"Action", "Fantasy"},
		},
		{
			Item: models.Item{
				Title:       "Tensei shitara Slime Datta Ken",
				Type:        models.MediaTypeComic,
				Description: "Satoru Mikami is an ordinary 37-year-old who is stabbed to death and reincarnated as a slime in another world.",
				ReleaseDate: seeders.ParseDate("2015-03-26"),
				CoverURL:    "https://cdn.myanimelist.net/images/manga/3/120337.jpg",
			},
			SpecificData: &models.BookData{
				Author:   "Fuse",
				Volumes:  26,
				Chapters: 120,
				Format:   "manga",
			},
			TagNames: []string{"Action", "Adventure", "Fantasy", "Shounen"},
		},
		{
			Item: models.Item{
				Title:       "Blue Lock",
				Type:        models.MediaTypeComic,
				Description: "After a disastrous defeat at the 2018 World Cup, Japan's team struggles to regroup. But what's missing? An absolute Ace Striker.",
				ReleaseDate: seeders.ParseDate("2018-08-01"),
				CoverURL:    "https://cdn.myanimelist.net/images/manga/1/229919.jpg",
			},
			SpecificData: &models.BookData{
				Author:   "Muneyuki Kaneshiro",
				Volumes:  28,
				Chapters: 280,
				Format:   "manga",
			},
			TagNames: []string{"Action", "Drama", "Shounen"},
		},

		// ========== NOVELS (Light Novel, Web Novel) ==========
		{
			Item: models.Item{
				Title:       "The Beginning After The End",
				Type:        models.MediaTypeNovel,
				Description: "Um rei reencarna em um mundo de magia e espada, mantendo memórias de sua vida passada.",
				ReleaseDate: seeders.ParseDate("2016-12-01"),
				CoverURL:    "https://i.pinimg.com/originals/3f/1b/2b/3f1b2b3b9b0f3a2b5b5b3b0f3a2b5b5b.jpg",
			},
			SpecificData: &models.BookData{
				Author:   "TurtleMe",
				Chapters: 500,
				Format:   "web_novel",
			},
			TagNames: []string{"Fantasy", "Adventure"},
		},
		{
			Item: models.Item{
				Title:       "Sword Art Online",
				Type:        models.MediaTypeNovel,
				Description: "Jogadores ficam presos em um MMORPG de realidade virtual onde a morte no jogo significa morte real.",
				ReleaseDate: seeders.ParseDate("2009-04-10"),
				CoverURL:    "https://cdn.myanimelist.net/images/anime/11/39717.jpg",
			},
			SpecificData: &models.BookData{
				Author:   "Reki Kawahara",
				Volumes:  28,
				Format:   "light_novel",
			},
			TagNames: []string{"Action", "Sci-Fi", "Fantasy"},
		},
		{
			Item: models.Item{
				Title:       "Mushoku Tensei: Isekai Ittara Honki Dasu",
				Type:        models.MediaTypeNovel,
				Description: "A 34-year-old NEET gets killed in a traffic accident and finds himself in a world of magic. Rather than waking up as a full-grown mage, he gets reincarnated as a newborn baby.",
				ReleaseDate: seeders.ParseDate("2012-09-22"),
				CoverURL:    "https://cdn.myanimelist.net/images/anime/1530/117776.jpg",
			},
			SpecificData: &models.BookData{
				Author:   "Rifujin na Magonote",
				Volumes:  26,
				Format:   "light_novel",
			},
			TagNames: []string{"Fantasy", "Adventure", "Drama"},
		},
		{
			Item: models.Item{
				Title:       "Hai to Gensou no Grimgar",
				Type:        models.MediaTypeNovel,
				Description: "When Haruhiro awakens, he's in the dark surrounded by people who have no memory of where they came from or how they got there.",
				ReleaseDate: seeders.ParseDate("2013-06-25"),
				CoverURL:    "https://cdn.myanimelist.net/images/anime/5/77809.jpg",
			},
			SpecificData: &models.BookData{
				Author:   "Ao Jyumonji",
				Volumes:  18,
				Format:   "light_novel",
			},
			TagNames: []string{"Fantasy", "Adventure", "Drama"},
		},

		// ========== BOOKS (A Song of Ice and Fire) ==========
		{
			Item: models.Item{
				Title:       "A Game of Thrones",
				Type:        models.MediaTypeBook,
				Description: "The first book in the epic fantasy series A Song of Ice and Fire. Winter is coming to the Seven Kingdoms.",
				ReleaseDate: seeders.ParseDate("1996-08-01"),
				CoverURL:    "https://covers.openlibrary.org/b/id/10521270-L.jpg",
			},
			SpecificData: &models.BookData{
				Author:    "George R. R. Martin",
				Pages:     694,
				Publisher: "Bantam Spectra",
			},
			TagNames: []string{"Fantasy", "Adventure", "Drama"},
		},
		{
			Item: models.Item{
				Title:       "A Clash of Kings",
				Type:        models.MediaTypeBook,
				Description: "The second book in A Song of Ice and Fire. The Seven Kingdoms are torn by strife as five kings claim the Iron Throne.",
				ReleaseDate: seeders.ParseDate("1998-11-16"),
				CoverURL:    "https://covers.openlibrary.org/b/id/10521271-L.jpg",
			},
			SpecificData: &models.BookData{
				Author:    "George R. R. Martin",
				Pages:     768,
				Publisher: "Bantam Spectra",
			},
			TagNames: []string{"Fantasy", "Adventure", "Drama"},
		},
		{
			Item: models.Item{
				Title:       "A Storm of Swords",
				Type:        models.MediaTypeBook,
				Description: "The third book in A Song of Ice and Fire. The War of the Five Kings is in full swing and chaos reigns.",
				ReleaseDate: seeders.ParseDate("2000-08-08"),
				CoverURL:    "https://covers.openlibrary.org/b/id/10521272-L.jpg",
			},
			SpecificData: &models.BookData{
				Author:    "George R. R. Martin",
				Pages:     973,
				Publisher: "Bantam Spectra",
			},
			TagNames: []string{"Fantasy", "Adventure", "Drama"},
		},
		{
			Item: models.Item{
				Title:       "A Feast for Crows",
				Type:        models.MediaTypeBook,
				Description: "The fourth book in A Song of Ice and Fire. In the aftermath of war, the realm struggles to rebuild.",
				ReleaseDate: seeders.ParseDate("2005-10-17"),
				CoverURL:    "https://covers.openlibrary.org/b/id/10521273-L.jpg",
			},
			SpecificData: &models.BookData{
				Author:    "George R. R. Martin",
				Pages:     753,
				Publisher: "Bantam Spectra",
			},
			TagNames: []string{"Fantasy", "Adventure", "Drama"},
		},
		{
			Item: models.Item{
				Title:       "A Dance with Dragons",
				Type:        models.MediaTypeBook,
				Description: "The fifth book in A Song of Ice and Fire. In the east, Daenerys Targaryen, the last heir of House Targaryen, rules with her three dragons.",
				ReleaseDate: seeders.ParseDate("2011-07-12"),
				CoverURL:    "https://covers.openlibrary.org/b/id/10521274-L.jpg",
			},
			SpecificData: &models.BookData{
				Author:    "George R. R. Martin",
				Pages:     1016,
				Publisher: "Bantam Spectra",
			},
			TagNames: []string{"Fantasy", "Adventure", "Drama"},
		},
	}
}
