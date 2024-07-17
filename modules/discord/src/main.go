package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bwmarrin/discordgo"
	pg "github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	"github.com/olekukonko/tablewriter"
)

var db *pgxpool.Pool
var dg *discordgo.Session
var input Input
var movieChannel = "938973912035901480"
var tvChannel = "938973946575978546"
var goChannel = "1119118378909585408"
var pythonChannel = "1055548690716180532"
var gamesChannel = "938973978263965747"
var githubChannel = "938973612461916169"
var javascriptChannel = "938974163572518943"
var testChannel = "870190331554054194"

type Input struct {
	Test bool `json:"test"`
}

func main() {
	local := os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == ""
	buildLogger(true, false, local)
	if local {
		handle(context.TODO(), Input{
			Test: true,
		})
	} else {
		lambda.Start(handle)
	}
}

func handle(ctx context.Context, i Input) (string, error) {
	var err error
	input = i
	db, err = pgxpool.New(context.Background(), os.Getenv("PG_URI"))
	check(err)
	defer db.Close()
	dg, err = discordgo.New("Bot " + os.Getenv("TOKEN"))
	check(err)
	err = dg.Open()
	check(err)
	games()
	trendingMovies()
	tv()
	upcomingMovies()
	golang()
	python()
	github()
	javascript()
	dg.Close()
	return "", nil
}

func trendingMovies() {
	var movies []*TrendingMovie
	var scrapedAt time.Time
	err := pg.Select(context.Background(), db, &movies, `SELECT * FROM trending_movies LIMIT 25`)
	check(err)

	slog.Info("selected", "rows", len(movies))
	if len(movies) == 0 {
		slog.Warn("no data in trending_movies")
		return
	}

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
	}, movieChannel, "Trending Movies")
}

func tv() {
	var tvs []*TrendingTV
	var scrapedAt time.Time
	err := pg.Select(context.Background(), db, &tvs, `SELECT * FROM trending_tvs LIMIT 25`)
	check(err)

	slog.Info("selected tv", "rows", len(tvs))
	if len(tvs) == 0 {
		slog.Warn("no data in trending_tvs")
		return
	}

	var fields []*discordgo.MessageEmbedField
	for _, v := range tvs {
		scrapedAt = v.UpdatedAt
		name := "#" + strconv.Itoa(v.Rank) + " " + v.Title
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
	}, tvChannel, "TV")
}

func upcomingMovies() {
	var movies []*UpcomingMovie
	var scrapedAt time.Time
	err := pg.Select(context.Background(), db, &movies, `SELECT * FROM upcoming_movies LIMIT 25`)
	check(err)

	slog.Info("selected upcomingMovies", "rows", len(movies))
	if len(movies) == 0 {
		slog.Warn("no data in upcoming_movies")
		return
	}

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
	}, movieChannel, "Upcoming Movies")
}

func golang() {
	var gos []*TrendingGo
	err := pg.Select(context.Background(), db, &gos, `SELECT * FROM trending_gos ORDER BY stars DESC LIMIT 100 OFFSET 19`)
	check(err)

	if len(gos) == 0 {
		slog.Warn("no data in trending_gos")
		return
	}

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
		for j := 0; j < len(reducedArr[i]); j++ {
			if s, ok := reducedArr[i][j].(*TrendingGo); ok {
				stars := strconv.FormatInt(s.Stars/1000, 10) + "k"
				s.Description = strings.ReplaceAll(ShortText(s.Description, 30), "\n", " ")
				tbData = append(tbData, []string{stars, s.Name, s.Description})
				scrapeTime = s.UpdatedAt
			}
		}
		table.AppendBulk(tbData)
		table.Render()
		messages = append(messages, fmt.Sprintf("```md\nTop Go Projects (%d/5), scraped %s\n\n%s```", i+1, scrapeTime.Format("01-02"), tableString))
		table.ClearRows()
		tableString.Reset()
	}
	post(messages, nil, goChannel, "Golang")
}

func python() {
	var pies []*TrendingPY
	err := pg.Select(context.Background(), db, &pies, `SELECT * FROM trending_pies ORDER BY downloads DESC`)
	check(err)

	slog.Info("selected python", "rows", len(pies))

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
		for j := 0; j < len(reducedArr[i]); j++ {
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
	post(messages, nil, pythonChannel, "Python")
}

func games() {
	var gs []*TrendingGame
	err := pg.Select(context.Background(), db, &gs, `SELECT * FROM trending_games LIMIT 16`)
	check(err)

	slog.Info("selected games", "rows", len(gs))

	if len(gs) == 0 {
		slog.Warn("no data in trending_games")
		return
	}

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
		tbData = append(tbData, []string{limitString(game.Title, 25), game.Price})
	}
	table.AppendBulk(tbData)
	table.Render()
	messages := []string{fmt.Sprintf("```md\nTop 16 Selling Games on Steam, scraped %s\n\n%s```", scrapeTime.Format("01-02"), tableString)}
	post(messages, nil, gamesChannel, "Games")
}

func javascript() {
	var jss []*TrendingJS
	err := pg.Select(context.Background(), db, &jss, `SELECT * FROM trending_js ORDER BY subject, rank`)
	check(err)

	slog.Info("selected javascript", "rows", len(jss))

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
		for j := 0; j < len(reducedArr[i]); j++ {
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
	post(messages, nil, javascriptChannel, "JavaScript")
}

func github() {
	var ghs []*TrendingGithub
	err := pg.Select(context.Background(), db, &ghs, `SELECT * FROM trending_githubs ORDER BY stars DESC`)
	check(err)

	slog.Info("selected github", "rows", len(ghs))
	if len(ghs) == 0 {
		slog.Warn("no data in trending_githubs")
		return
	}

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
	post(messages, nil, githubChannel, "github")
}

func post(messages []string, embed *discordgo.MessageEmbed, channelId string, channelName string) {
	slog.Info("Posting to " + channelName)
	if input.Test {
		channelId = "870190331554054194"
	}
	if embed == nil {
		for _, msg := range messages {
			_, err := dg.ChannelMessageSend(channelId, msg)
			check(err)
		}
	} else {
		_, err := dg.ChannelMessageSendEmbed(channelId, embed)
		check(err)
	}
}
