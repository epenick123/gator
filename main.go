package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/epenick123/blogagg/internal/config"
	"github.com/epenick123/blogagg/internal/database"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// Define your structs
type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	// name and args fields
	name string
	args []string
}

type commands struct {
	// your map of handlers
	cmd_names map[string]func(*state, command) error
}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// Define your methods and functions
func handlerLogin(s *state, cmd command) error {
	// implementation
	if len(cmd.args) == 0 {
		return fmt.Errorf("The login handler expects a single argument, the username.")
	}
	user, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err == sql.ErrNoRows {
		fmt.Println("user not found")
		os.Exit(1)
	} else if err != nil {
		return err
	}
	if err := s.cfg.SetUser(user.Name); err != nil {
		return err
	}
	fmt.Println("The user has been set.")
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	// implementation
	handler, exists := c.cmd_names[cmd.name]
	if !exists {
		return fmt.Errorf("unknown command: %s", cmd.name)
	}
	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	// implementation
	c.cmd_names[name] = f

}
func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		fmt.Println("username required")
		os.Exit(1)
	}
	name := cmd.args[0]
	// Try to get the user:
	user, err := s.db.GetUser(context.Background(), name)
	if err == nil {
		fmt.Println("user already exists")
		os.Exit(1)
	}
	if err != sql.ErrNoRows {
		// Some other database error
		return err
	}

	// Construct input for CreateUser (see your generated code!)
	user, err = s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Set user in config and persist if needed
	if err := s.cfg.SetUser(user.Name); err != nil {
		log.Fatal(err)
	}
	fmt.Println("User created:", user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.Reset(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("Database reset successfully") // Add this line
	return nil
}

func handlerUsers(s *state, cmd command) error {
	currentUser := s.cfg.CurrentUserName
	all_users, err := s.db.GetUsers(context.Background())

	if err != nil {
		return err
	}
	for _, user := range all_users {
		if user.Name == currentUser {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	rssf, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Printf("%+v", rssf)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("Missing name and/or url")
	}
	name := cmd.args[0]
	url := cmd.args[1]

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}
	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", feed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.Feeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("Feed Name: %s | Feed URL: %s | User: %s\n", feed.FeedName, feed.Url, feed.UserName)
	}
	return nil
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := &http.Client{}
	method := http.MethodGet

	req, err := http.NewRequestWithContext(ctx, method, feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	rssf := RSSFeed{}
	err = xml.Unmarshal(body, &rssf)
	if err != nil {
		return nil, err
	}

	rssf.Channel.Description = html.UnescapeString(rssf.Channel.Description)
	rssf.Channel.Title = html.UnescapeString(rssf.Channel.Title)
	for i, _ := range rssf.Channel.Item {
		rssf.Channel.Item[i].Description = html.UnescapeString(rssf.Channel.Item[i].Description)
		rssf.Channel.Item[i].Title = html.UnescapeString(rssf.Channel.Item[i].Title)
	}

	return &rssf, err
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	s2 := state{
		db:  dbQueries,
		cfg: &cfg,
	}

	cmds := commands{
		cmd_names: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerFeeds)

	if len(os.Args) < 2 {
		fmt.Println("not enough arguments provided")
		os.Exit(1)
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	cmd := command{
		name: cmdName,
		args: cmdArgs,
	}

	if err := cmds.run(&s2, cmd); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
