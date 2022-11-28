package main

import (
	"fmt"
)

type Cache struct {
	Users  map[string]User
	Movies map[int]Movie
	Series map[int]Series
}

func (c *Cache) init() error {

	c.Users = make(map[string]User)
	c.Movies = make(map[int]Movie)
	c.Series = make(map[int]Series)

	ub, err := getAllUsers(connectionPool)
	if err != nil {
		return fmt.Errorf("cache.go: init(): error loading users from database: %s", err.Error())
	}

	mb, err := getAllMoviesBulk(connectionPool)
	if err != nil {
		return fmt.Errorf("cache.go: init(): error loading movies from database: %s", err.Error())
	}

	sb, err := getAllSeriesBulk(connectionPool)
	if err != nil {
		return fmt.Errorf("cache.go: init(): error loading series from database: %s", err.Error())
	}

	c.parseUsersBulk(ub)
	c.parseMoviesBulk(mb)
	c.parseSeriesBulk(sb)

	return nil
}

func (c *Cache) parseUsersBulk(u []UserBulk) {

	permissions := make(map[int]Permission, 0)
	groups := make(map[int]Group, 0)

	// permissions
	for _, ub := range u {
		if _, found := permissions[ub.Permission.Id]; !found {
			permission := Permission{
				Id:   ub.Permission.Id,
				Name: ub.Permission.Name,
			}
			permissions[ub.Permission.Id] = permission
		}
	}

	// groups
	for _, ub := range u {
		if _, found := groups[ub.Group.Id]; !found {
			group := Group{
				Id:          ub.Group.Id,
				Name:        ub.Group.Name,
				Permissions: make(map[string]Permission),
			}
			groups[ub.Group.Id] = group
		}
		groups[ub.Group.Id].Permissions[ub.Permission.Name] = permissions[ub.Permission.Id]
	}

	// users
	for _, ub := range u {
		if _, found := c.Users[ub.User.Email]; !found {
			user := User{
				Id:             ub.User.Id,
				Firstname:      ub.User.Firstname,
				Lastname:       ub.User.Lastname,
				Email:          ub.User.Email,
				HashedPassword: ub.User.HashedPassword,
				ActivatedAt:    ub.User.ActivatedAt,
				CreatedAt:      ub.User.CreatedAt,
				UpdatedAt:      ub.User.UpdatedAt,
				DeletedAt:      ub.User.DeletedAt,
				Groups:         make(map[string]Group),
			}
			c.Users[ub.User.Email] = user
		}
		c.Users[ub.User.Email].Groups[ub.Group.Name] = groups[ub.Group.Id]
	}
}

func (c *Cache) parseMoviesBulk(m []MovieBulk) {

	files := make(map[int]File, 0)
	formats := make(map[int]Format, 0)
	languages := make(map[int]Language, 0)

	// format
	for _, mb := range m {
		if mb.Format.Id.Valid {
			if _, found := formats[int(mb.Format.Id.Int64)]; !found {
				format := Format{
					Id:   int(mb.Format.Id.Int64),
					Name: mb.Format.Name.String,
				}
				formats[int(mb.Format.Id.Int64)] = format
			}
		}
	}

	// type
	for _, mb := range m {
		if mb.Language.Id.Valid {
			if _, found := languages[int(mb.Language.Id.Int64)]; !found {
				language := Language{
					Id:    int(mb.Language.Id.Int64),
					Short: mb.Language.Short.String,
					Name:  mb.Language.Name.String,
				}
				languages[int(mb.Language.Id.Int64)] = language
			}
		}
	}

	// files
	for _, mb := range m {
		if mb.File.Id.Valid {
			if _, found := files[int(mb.File.Id.Int64)]; !found {
				file := File{
					Id:         int(mb.File.Id.Int64),
					FormatId:   mb.File.FormatId,
					LanguageId: mb.File.LanguageId,
					MovieId:    mb.File.MovieId,
					EpisodeId:  mb.File.EpisodeId,
					Hash:       mb.File.Hash,
					CreatedAt:  mb.File.CreatedAt,
					UpdatedAt:  mb.File.UpdatedAt,
					DeletedAt:  mb.File.DeletedAt,
					Format:     formats[int(mb.File.FormatId.Int64)],
					Language:   languages[int(mb.File.LanguageId.Int64)],
				}
				files[int(mb.File.Id.Int64)] = file
			}
		}
	}

	// movies
	for _, mb := range m {
		if _, found := c.Movies[mb.Movie.Id]; !found {
			movie := Movie{
				Id:               mb.Movie.Id,
				Title:            mb.Movie.Title,
				Description:      mb.Movie.Description,
				PrecedingMovieId: mb.Movie.PrecedingMovieId,
				CreatedAt:        mb.Movie.CreatedAt,
				UpdatedAt:        mb.Movie.UpdatedAt,
				DeletedAt:        mb.Movie.DeletedAt,
				Files:            make(map[int]File, 0),
			}
			for _, file := range files {
				if file.MovieId.Int64 == int64(movie.Id) {
					movie.Files[file.Id] = file
				}
			}
			c.Movies[mb.Movie.Id] = movie
		}
	}
}

func (c *Cache) parseSeriesBulk(s []SeriesBulk) {

	seasons := make(map[int]Season, 0)
	episodes := make(map[int]Episode, 0)
	files := make(map[int]File, 0)
	formats := make(map[int]Format, 0)
	languages := make(map[int]Language, 0)

	// format
	for _, sb := range s {
		if sb.Format.Id.Valid {
			if _, found := formats[int(sb.Format.Id.Int64)]; !found {
				format := Format{
					Id:   int(sb.Format.Id.Int64),
					Name: sb.Format.Name.String,
				}
				formats[int(sb.Format.Id.Int64)] = format
			}
		}
	}

	// type
	for _, sb := range s {
		if sb.Language.Id.Valid {
			if _, found := languages[int(sb.Language.Id.Int64)]; !found {
				language := Language{
					Id:    int(sb.Language.Id.Int64),
					Short: sb.Language.Short.String,
					Name:  sb.Language.Name.String,
				}
				languages[int(sb.Language.Id.Int64)] = language
			}
		}
	}

	// files
	for _, sb := range s {
		if sb.File.Id.Valid {
			if _, found := files[int(sb.File.Id.Int64)]; !found {
				file := File{
					Id:         int(sb.File.Id.Int64),
					FormatId:   sb.File.FormatId,
					LanguageId: sb.File.LanguageId,
					MovieId:    sb.File.MovieId,
					EpisodeId:  sb.File.EpisodeId,
					Hash:       sb.File.Hash,
					CreatedAt:  sb.File.CreatedAt,
					UpdatedAt:  sb.File.UpdatedAt,
					DeletedAt:  sb.File.DeletedAt,
					Format:     formats[int(sb.File.FormatId.Int64)],
					Language:   languages[int(sb.File.LanguageId.Int64)],
				}
				files[int(sb.File.Id.Int64)] = file
			}
		}
	}

	for _, sb := range s {
		if sb.Episode.Id.Valid {
			episode := Episode{
				Id:          int(sb.Episode.Id.Int64),
				SeasonId:    sb.Episode.SeasonId,
				Number:      sb.Episode.Number,
				Title:       sb.Episode.Title,
				Description: sb.Episode.Description,
				CreatedAt:   sb.Episode.CreatedAt,
				UpdatedAt:   sb.Episode.UpdatedAt,
				DeletedAt:   sb.Episode.DeletedAt,
				Files:       make(map[int]File, 0),
			}
			for _, file := range files {
				if file.EpisodeId.Int64 == int64(episode.Id) {
					episode.Files[file.Id] = file
				}
			}
			episodes[int(sb.Episode.Id.Int64)] = episode
		}
	}

	for _, sb := range s {
		if sb.Season.Id.Valid {
			season := Season{
				Id:          int(sb.Season.Id.Int64),
				SeriesId:    sb.Season.SeriesId,
				Number:      int(sb.Season.Number.Int64),
				Title:       sb.Season.Title.String,
				Description: sb.Season.Description.String,
				CreatedAt:   sb.Season.CreatedAt,
				UpdatedAt:   sb.Season.UpdatedAt,
				DeletedAt:   sb.Season.DeletedAt,
				Episodes:    make(map[int]Episode, 0),
			}
			for _, episode := range episodes {
				if episode.SeasonId.Int64 == int64(season.Id) {
					season.Episodes[episode.Id] = episode
				}
			}
			seasons[int(sb.Season.Id.Int64)] = season
		}
	}

	for _, sb := range s {
		if _, found := c.Series[sb.Series.Id]; !found {
			s := Series{
				Id:          sb.Series.Id,
				Title:       sb.Series.Title,
				Description: sb.Series.Description,
				CreatedAt:   sb.Series.CreatedAt,
				UpdatedAt:   sb.Series.UpdatedAt,
				DeletedAt:   sb.Series.DeletedAt,
				Seasons:     make(map[int]Season, 0),
			}
			for _, season := range seasons {
				if season.SeriesId.Int64 == int64(s.Id) {
					s.Seasons[season.Id] = season
				}
			}
			c.Series[sb.Series.Id] = s
		}
	}
}
