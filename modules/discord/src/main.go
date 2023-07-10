package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bwmarrin/discordgo"
	pg "github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	"github.com/olekukonko/tablewriter"
)

var dg *discordgo.Session

func main() {
	buildLogger()
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == "" {
		handle(context.TODO(), nil)
	} else {
		lambda.Start(handle)
	}
}

func handle(ctx context.Context, _ any) (string, error) {
	db, err := pgxpool.New(context.Background(), os.Getenv("PG_URI"))
	check(err)
	defer db.Close()
	check(err)
	dg, err = discordgo.New("Bot " + os.Getenv("TOKEN"))
	check(err)
	err = dg.Open()
	check(err)
	trendingMovies(db, "938973912035901480")
	tv(db, "938973946575978546")
	upcomingMovies(db, "938973912035901480")
	golang(db, "1119118378909585408")
	python(db, "1055548690716180532")
	games(db, "938973978263965747")
	github(db, "938973612461916169")
	javascript(db, "938974163572518943")
	dg.Close()
	return "", nil
}

func trendingMovies(db *pgxpool.Pool, channelId string) {
	log := logger.With().Str("func", "trendingMovies").Logger()
	var movies []*TrendingMovie
	var scrapedAt time.Time
	err := pg.Select(context.Background(), db, &movies, `SELECT * FROM trending_movies LIMIT 25`)
	check(err, log)

	log.Print("selected rows ", len(movies))

	var fields []*discordgo.MessageEmbedField
	for _, m := range movies {
		scrapedAt = m.UpdatedAt
		name := "#" + strconv.Itoa(m.Rank) + " " + m.Title
		value := "★" + m.Rating + " (" + m.Velocity + ")"
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  name,
			Value: value,
		})
	}
	post(nil, &discordgo.MessageEmbed{
		Fields: fields,
		Title:  "Trending Movies",
		Color:  16776960,
		URL:    "https://www.imdb.com/chart/moviemeter",
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "http://icons.iconarchive.com/icons/danleech/simple/1024/imdb-icon.png",
		},
		Author: &discordgo.MessageEmbedAuthor{
			URL:     "https://codabool.com",
			Name:    "CodaBot",
			IconURL: "https://avatars.githubusercontent.com/u/61724833?v=4",
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Scraped " + scrapedAt.Format("January 2"),
			IconURL: "http://icons.iconarchive.com/icons/danleech/simple/1024/imdb-icon.png",
		},
	}, channelId, "Trending Movies")
}

func tv(db *pgxpool.Pool, channelId string) {
	log := logger.With().Str("func", "tv").Logger()
	var tvs []*TrendingTV
	var scrapedAt time.Time
	err := pg.Select(context.Background(), db, &tvs, `SELECT * FROM trending_tvs LIMIT 25`)
	check(err, log)

	log.Print("selected rows ", len(tvs))

	var fields []*discordgo.MessageEmbedField
	for _, v := range tvs {
		scrapedAt = v.UpdatedAt
		name := "#" + strconv.Itoa(v.Rank+1) + " " + v.Title
		value := "★" + v.Rating + " (" + v.Velocity + ")"
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  name,
			Value: value,
		})
	}
	post(nil, &discordgo.MessageEmbed{
		Fields: fields,
		Title:  "Trending TV",
		Color:  16776960,
		URL:    "https://www.imdb.com/chart/tvmeter",
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "http://icons.iconarchive.com/icons/danleech/simple/1024/imdb-icon.png",
		},
		Author: &discordgo.MessageEmbedAuthor{
			URL:     "https://codabool.com",
			Name:    "CodaBot",
			IconURL: "https://avatars.githubusercontent.com/u/61724833?v=4",
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Scraped " + scrapedAt.Format("January 2"),
			IconURL: "http://icons.iconarchive.com/icons/danleech/simple/1024/imdb-icon.png",
		},
	}, channelId, "TV")
}

func GroupBy(xs []map[string]interface{}, key string) map[interface{}][]map[string]interface{} {
	rv := make(map[interface{}][]map[string]interface{})
	for _, x := range xs {
		k := x[key]
		rv[k] = append(rv[k], x)
	}
	return rv
}

func upcomingMovies(db *pgxpool.Pool, channelId string) {
	log := logger.With().Str("func", "upcomingMovies").Logger()
	var movies []*UpcomingMovie
	var scrapedAt time.Time
	err := pg.Select(context.Background(), db, &movies, `SELECT * FROM upcoming_movies LIMIT 25`)
	check(err, log)

	log.Print("selected rows ", len(movies))

	var newMovies []map[string]interface{}
	for _, movie := range movies {
		newMovies = append(newMovies, map[string]interface{}{
			"Title":   movie.Title,
			"Release": movie.Release,
		})
		scrapedAt = movie.UpdatedAt
	}
	var releases []time.Time
	grouped := GroupBy(newMovies, "Release")
	for key, _ := range grouped {
		if t, ok := key.(time.Time); ok {
			releases = append(releases, t)
		}
	}
	sort.Slice(releases, func(i, j int) bool {
		return releases[i].Before(releases[j])
	})
	var fields []*discordgo.MessageEmbedField
	for _, releaseDate := range releases {
		for key, val := range grouped {
			if release, ok := key.(time.Time); ok {
				if release == releaseDate {
					var value string
					for _, v := range val {
						value += fmt.Sprintf("%v", v["Title"]) + "\n"
					}
					fields = append(fields, &discordgo.MessageEmbedField{
						Name:  release.Format("January 2"),
						Value: value,
					})
				}
			}
		}
	}

	post(nil, &discordgo.MessageEmbed{
		Fields: fields,
		Title:  "Upcoming Movies",
		Color:  16776960,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "http://icons.iconarchive.com/icons/danleech/simple/1024/imdb-icon.png",
		},
		Author: &discordgo.MessageEmbedAuthor{
			URL:     "https://codabool.com",
			Name:    "CodaBot",
			IconURL: "https://avatars.githubusercontent.com/u/61724833?v=4",
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Scraped " + scrapedAt.Format("January 2"),
			IconURL: "http://icons.iconarchive.com/icons/danleech/simple/1024/imdb-icon.png",
		},
	}, channelId, "Upcoming Movies")
}

func golang(db *pgxpool.Pool, channelId string) {
	log := logger.With().Str("func", "golang").Logger()
	var gos []*TrendingGo
	err := pg.Select(context.Background(), db, &gos, `SELECT * FROM trending_gos ORDER BY stars DESC LIMIT 100 OFFSET 19`)
	check(err, log)

	log.Print("selected rows ", len(gos))

	interfaces := make([]interface{}, len(gos))
	for i, v := range gos {
		interfaces[i] = v
	}
	reducedArr := reduce(interfaces, 20)

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"Stars", "Repo", "Description"})
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	// table.SetRowSeparator("")
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	var messages []string
	for i := 0; i < 5; i++ {
		var scrapeTime time.Time
		var tbData [][]string
		for j := 0; j < 20; j++ {
			if s, ok := reducedArr[i][j].(*TrendingGo); ok {
				stars := strconv.FormatInt(s.Stars/1000, 10) + "k"
				tbData = append(tbData, []string{stars, s.Name, ShortText(s.Description, 30)})
				scrapeTime = s.UpdatedAt
			}
		}
		table.AppendBulk(tbData)
		table.Render()
		messages = append(messages, fmt.Sprintf("```md\nTop Go Projects (%d/5), scraped %s\n\n%s```", i+1, scrapeTime.Format("01-02"), tableString))
		table.ClearRows()
		tableString.Reset()
	}
	post(messages, nil, channelId, "Golang")
}

func python(db *pgxpool.Pool, channelId string) {
	log := logger.With().Str("func", "python").Logger()
	var pies []*TrendingPY
	err := pg.Select(context.Background(), db, &pies, `SELECT * FROM trending_pies ORDER BY downloads DESC`)
	check(err, log)

	log.Print("selected rows ", len(pies))

	interfaces := make([]interface{}, len(pies))
	for i, v := range pies {
		interfaces[i] = v
	}
	reducedArr := reduce(interfaces, 20)

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"Downloads", "Package", "Description"})
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	var messages []string
	for i := 0; i < 5; i++ {
		var tbData [][]string
		var scrapeTime time.Time
		for j := 0; j < 20; j++ {
			if s, ok := reducedArr[i][j].(*TrendingPY); ok {
				scrapeTime = s.UpdatedAt
				downloads := strconv.FormatInt(s.Downloads/1000/1000, 10) + "m"
				tbData = append(tbData, []string{downloads, s.Name, ShortText(s.Description, 30)})
			}
		}
		table.AppendBulk(tbData)
		table.Render()
		messages = append(messages, fmt.Sprintf("```md\nTop Python packages (%d/5), scraped %s\n\n%s```", i+1, scrapeTime.Format("01-02"), tableString))
		table.ClearRows()
		tableString.Reset()
	}
	post(messages, nil, channelId, "Python")
}

func games(db *pgxpool.Pool, channelId string) {
	log := logger.With().Str("func", "games").Logger()
	var gs []*TrendingGame
	err := pg.Select(context.Background(), db, &gs, `SELECT * FROM trending_games LIMIT 16`)
	check(err, log)

	log.Print("selected rows ", len(gs))
	log.Trace().Int("rows", len(gs)).Msg("")

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"Title", "Price"})
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	var tbData [][]string
	var scrapeTime time.Time
	for _, game := range gs {
		scrapeTime = game.UpdatedAt
		tbData = append(tbData, []string{game.Title, game.Price})
	}
	table.AppendBulk(tbData)
	table.Render()
	messages := []string{fmt.Sprintf("```md\nTop 16 Selling Games on Steam, scraped %s\n\n%s```", scrapeTime.Format("01-02"), tableString)}
	post(messages, nil, channelId, "Games")
}

func javascript(db *pgxpool.Pool, channelId string) {
	log := logger.With().Str("func", "javascript").Logger()
	var jss []*TrendingJS
	err := pg.Select(context.Background(), db, &jss, `SELECT * FROM trending_js ORDER BY subject, rank`)
	check(err, log)

	log.Print("selected rows ", len(jss))

	interfaces := make([]interface{}, len(jss))
	for i, v := range jss {
		interfaces[i] = v
	}
	reducedArr := reduce(interfaces, 20)

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"Rank", "Package", "Description"})
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	var messages []string
	for i := 0; i < 4; i++ {
		var subject string
		var sTime time.Time
		var tbData [][]string
		for j := 0; j < 20; j++ {
			if s, ok := reducedArr[i][j].(*TrendingJS); ok {
				tbData = append(tbData, []string{strconv.Itoa(s.Rank), s.Title, ShortText(s.Description, 30)})
				subject = s.Subject
				sTime = s.UpdatedAt
			}
		}
		table.AppendBulk(tbData)
		table.Render()
		messages = append(messages, fmt.Sprintf("```md\nTop 20 %s JavaScript packages, scraped %s\n\n%s```", subject, sTime.Format("01-02"), tableString))
		table.ClearRows()
		tableString.Reset()
	}
	post(messages, nil, channelId, "JavaScript")
}

func ShortText(s string, i int) string {
	if len(s) < i {
		return s
	}
	if utf8.ValidString(s[:i]) {
		return s[:i]
	}
	return s[:i+1]
}

func github(db *pgxpool.Pool, channelId string) {
	log := logger.With().Str("func", "github").Logger()
	var ghs []*TrendingGithub
	err := pg.Select(context.Background(), db, &ghs, `SELECT * FROM trending_githubs ORDER BY stars DESC`)
	check(err, log)

	log.Print("selected rows ", len(ghs))

	interfaces := make([]interface{}, len(ghs))
	for i, v := range ghs {
		interfaces[i] = v
	}
	reducedArr := reduce(interfaces, 20)

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"Stars", "Repo", "Description"})
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	var messages []string
	for i := 0; i < 5; i++ {
		var scrapeTime time.Time
		var tbData [][]string
		for j := 0; j < 20; j++ {
			if s, ok := reducedArr[i][j].(*TrendingGithub); ok {
				stars := strconv.FormatInt(s.Stars/1000, 10) + "k"
				tbData = append(tbData, []string{stars, s.Name, ShortText(s.Description, 30)})
				scrapeTime = s.UpdatedAt
			}
		}
		table.AppendBulk(tbData)
		table.Render()
		messages = append(messages, fmt.Sprintf("```md\nTop GitHub Repos (%d/5), scraped %s\n\n%s```", i+1, scrapeTime.Format("01-02"), tableString))
		table.ClearRows()
		tableString.Reset()
	}
	post(messages, nil, channelId, "github")
}

func reduce(arr []interface{}, chunkSize int) [][]interface{} {
	chunks := make([][]interface{}, 0)
	chunk := make([]interface{}, 0)
	for i, item := range arr {
		chunkIndex := i / chunkSize
		if chunkIndex >= len(chunks) {
			chunks = append(chunks, chunk)
			chunk = make([]interface{}, 0)
		}
		chunk = append(chunk, item)
		chunks[chunkIndex] = chunk
	}
	if len(chunk) > 0 {
		chunks = append(chunks, chunk)
	}
	return chunks
}

func post(messages []string, embed *discordgo.MessageEmbed, channelId string, channelName string) {
	log := logger.With().Str("func", "post").Logger()
	log.Info().Msg("Posting to " + channelName)
	if os.Getenv("POST_TO_TEST") != "" {
		channelId = "870190331554054194"
	}
	if embed == nil {
		for _, msg := range messages {
			_, err := dg.ChannelMessageSend(channelId, msg)
			check(err, log)
		}
	} else {
		_, err := dg.ChannelMessageSendEmbed(channelId, embed)
		check(err, log)
	}
}
