Read-Me for Gator CLI Blog Aggregator Tool
Thank you for reading!

You will need Postgres and Go installed on your machine in order to run the program.

You can install the program using the command:
    1.) go install github.com/epenick123/gator @latest
    2.) Open your favorite text editor (VSCode, Notepad, etc.)
    3.) create .gatorconfig.json and save it to your home directory. The file should look like:

{
    "db_url": "postgres://postgresUSERNAME:postgresPASSWORD@localhost:5432/gator?sslmode=disable",
    "current_user_name": "CURRENT_USER_NAME"
  }

To execute a command:
In the terminal, type "gator 'cmd'" at the root directory of gator, where 'cmd' is the name of the command:

    "login": login with a new user.
	"register": register a new user.
	"reset": reset the database and delete all data
	"users": display a list of users
	"agg": aggregate the most recent posts from followed feeds within a time frame, arg 1
	"feeds": show all followed feeds for a user
	"follow": follow a new feed with the current user
	"following": display a list of followed feeds for a user
	"addfeed": Add a feed to the database
	"unfollow": remove a feed from a user follow list
	"browse": display recent feeds
    