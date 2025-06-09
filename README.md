# Gator
RSS Feed Aggregator Buitld in Golang


## PREREQUSITES
- Go (GoLang)
```
curl -sS https://webi.sh/golang | sh
```
- Postgresql
```
sudo apt update
sudo apt install postgresql postgresql-contrib
```

## Installation
to insta gator CLI use the following CLI command:
```
go install github.com/BahryJarbou/Gator
```

## Setting up the config file
the config file must be created at the home directory as follows.
```
~/.gatorconfig.json
```
and include the following:
```
{
    "db_url":"postgres://<psqlusername>:<psqlpassword>@localhost:5432/gator?sslmode=disable",
    "current_user_name":<gatorusername>
}
```

## to run the program use
```
gator <command> <args>
```
commands include:
1. login *username*
2. register *username*
3. reset
4. users
5. agg *time_between_updated*
6. addfeed *feedname* *feedURL*
7. feeds
8. follow *feedurl*
9. following *username*
10. unfollow *feedurl*
11. browse *optional limit: default = 2*