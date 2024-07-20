module scrapeops

go 1.22

replace github.com/Radicalius/scrapeops/shared => ../shared

require (
	github.com/PuerkitoBio/goquery v1.9.2
	github.com/Radicalius/scrapeops/shared v0.0.0-20240629173140-10e98d66c3fd
	github.com/mattn/go-sqlite3 v1.14.22
	github.com/montanaflynn/stats v0.7.1
	github.com/robfig/cron v1.2.0
	gorm.io/driver/sqlite v1.5.6
	gorm.io/gorm v1.25.11
)

require (
	github.com/andybalholm/cascadia v1.3.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)
