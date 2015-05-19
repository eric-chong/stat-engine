package main

import (
  "fmt"
  "os"
  "strings"
  "strconv"
  "regexp"
  "github.com/PuerkitoBio/goquery"
)

var (
  games = []game{}
  teams = []team{}
  test = []string{}
  statPageBaseUrl = "http://www.nhl.com/scores/htmlreports/"
  season = "20142015"
  gameSummaryPrefix = "GS"
)

type game struct {
  id string
  date string
  time string
  homeTeam string
  awayTeam string
  gameSeq gameSequence
  info gameInfo
}

// game summary: http://www.nhl.com/scores/htmlreports/20142015/GS020003.HTM
func (g *game) pullGameSummary() {
  gameSummaryUrl := statPageBaseUrl + season + "/" + gameSummaryPrefix + g.id[4:10] + ".HTM"
  doc, err := goquery.NewDocument(gameSummaryUrl) 
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  } else {
    doc.Find("table#MainTable tbody tr").Each(func(i int, s *goquery.Selection) {      // fmt.Println(s.Text())
      if (i==0) {
        visitorStringR, _ := regexp.Compile("([A-Z, ' ']+)Game ([0-9]+) Away Game ([0-9]+)")
        homeStringR, _ := regexp.Compile("([A-Z, ' ']+)Game ([0-9]+) Home Game ([0-9]+)")
        visitorContainer := s.Find("table#Visitor > tbody > tr")
        homeContainer := s.Find("table#Home > tbody > tr")
        gameInfoContainer := s.Find("table#GameInfo tbody")

        _, g.gameSeq.awayTeamSeq, g.gameSeq.awayTeamAwaySeq = getGameSeq(visitorContainer, visitorStringR)
        _, g.gameSeq.homeTeamSeq, g.gameSeq.homeTeamHomeSeq = getGameSeq(homeContainer, homeStringR)
        g.date,g.info.attendance,g.info.arena,g.info.startTime,g.info.startTimeZone,g.info.endTime,g.info.endTimeZone,g.gameSeq.globalSeq = getGameInfo(gameInfoContainer)
        fmt.Println(getGameInfo(gameInfoContainer))
      }
    })
  }
}

func getGameSeq(s *goquery.Selection, r *regexp.Regexp) (string, int, int) {
  var (
    teamName string
    teamSeq int
    teamAwaySeq int
  )
  s.Each(func(i int, s *goquery.Selection) {
    if (i==2) {
      for i, v := range r.FindStringSubmatch(s.Text()) {
        switch i {
        case 1:
          teamName = v
        case 2:
          teamSeq = stringToInt(v)
        case 3:
          teamAwaySeq = stringToInt(v)
        }
      }
    }  
  })

  return teamName, teamSeq, teamAwaySeq
}

func getGameInfo(s *goquery.Selection) (string,string,string,string,string,string,string,int) {
  var (
    date string
    attendance string
    arena string
    startTime string
    startTimeZone string
    endTime string
    endTimeZone string
    gameSeq int
  )
  attendanceRegex, _ := regexp.Compile("Attendance ([0-9,',']+) at ([a-z,A-Z,' ']+)")
  timeRegex, _ := regexp.Compile("Start ([0-9,':']+) ([A-Z]+); End ([0-9,':']+) ([A-Z]+)")
  gameSeqRegex, _ := regexp.Compile("Game ([0-9]+)")
  s.Find("tr > td").Each(func(i int, s *goquery.Selection) {
    switch i {
    case 3:
      date = s.Text()
    case 4:
      for i, v := range attendanceRegex.FindStringSubmatch(s.Text()) {
        switch i {
        case 1:
          attendance = v
        case 2:
          arena = v
        }
      }
    case 5:
      for i, v := range timeRegex.FindStringSubmatch(strings.Replace(s.Text(), "&nbsp;", " ", -1)) {
        switch i {
        case 1:
          startTime = v
        case 2:
          startTimeZone = v
        case 3:
          endTime = v
        case 4:
          endTimeZone = v
        }
      }
    case 6:
      for i, v := range gameSeqRegex.FindStringSubmatch(strings.Replace(s.Text(), "&nbsp;", " ", -1)) {
        switch i {
        case 1:
          gameSeq = stringToInt(v)
        }
      }
    }
  })

  return date, attendance, arena, startTime, startTimeZone, endTime, endTimeZone, gameSeq
}

// event summary: http://www.nhl.com/scores/htmlreports/20142015/ES020003.HTM
// faceoff comparison: http://www.nhl.com/scores/htmlreports/20142015/FC020003.HTM
// faceoff summary: http://www.nhl.com/scores/htmlreports/20142015/FS020003.HTM
// Shift Report (Away): http://www.nhl.com/scores/htmlreports/20142015/TV020003.HTM
// Shift Report (Home): http://www.nhl.com/scores/htmlreports/20142015/TH020003.HTM
// Roster Report: http://www.nhl.com/scores/htmlreports/20142015/RO020003.HTM
// Shot Report: http://www.nhl.com/scores/htmlreports/20142015/SS020003.HTM

type gameSequence struct {
  globalSeq int
  homeTeamSeq int
  awayTeamSeq int
  homeTeamHomeSeq int
  awayTeamAwaySeq int
}

type gameInfo struct {
  startTime string
  startTimeZone string
  endTime string
  endTimeZone string
  attendance string
  arena string
}

type team struct {
  key string
  conference string
  division string
  city string
  name string
  site string
}


func main() {
  // fmt.Println("NHL Stat App")
  // generateTeams()
  // fmt.Println("Number of teams: ", len(teams))
  // pullGames()
  // games[len(games)-1].pullGameSummary()
  
  test_game := game{"2014021230", "Sat Apr 11, 2015", "10:00 PM ET", "Edmonton", "Vancouver", gameSequence{}, gameInfo{}}
  test_game.pullGameSummary()
  fmt.Println(test_game)

  // test_string := "EDMONTON OILERSGame 82 Away Game 41"
  // r, _ := regexp.Compile("([A-Z, ' ']+)Game ([0-9]+) Away Game ([0-9]+)")
  // for _, v := range r.FindStringSubmatch(test_string) {
  //   fmt.Println(v)
  // }
}

func generateTeams() {
  fmt.Println("Get Teams :-")
  teamPageUrl := "http://www.nhl.com/ice/teams.htm?navid=nav-tms-main"
  doc, err := goquery.NewDocument(teamPageUrl) 
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  } else {
    doc.Find("div.teamContainer > div.container").Each(func(i int, s *goquery.Selection) {
      conference := strings.Split(s.AttrOr("class", ""), " ")[0]
      s.Children().Each(func(i int, s * goquery.Selection) {
        if (!strings.Contains(s.AttrOr("class",""), "divhead")) {
          division := s.AttrOr("class","")
          s.Find("div.teamCard").Each(func(i int, s *goquery.Selection) {
            var key string
            for i, k := range strings.Split(s.AttrOr("class",""), " ") {
              if (i==1) {
                key = strings.ToUpper(k)
              }
            }
            city := s.Find("div.teamName span.teamPlace").Text()
            name := s.Find("div.teamName span.teamCommon").Text()
            site := s.Find("div.teamLogo a").AttrOr("href","")
            teams = append(teams,team{key,conference,division,city,name,site,})
          })
        }
      })
    })
    for i, t := range teams {
      fmt.Println(i, t)
    }
  }
}

func pullGames() {
  fmt.Println("Get Games :-")
  regularSeasonUrl := "http://www.nhl.com/ice/schedulebyseason.htm?season=20142015&gameType=2&team=&network=&venue="
  doc, err := goquery.NewDocument(regularSeasonUrl) 
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  } else {
    doc.Find("table.data.schedTbl tbody tr").Each(func(i int, s *goquery.Selection) {
      var (
        id string
        awayTeam string
        homeTeam string
      )
      date := s.Find("td.date div.skedStartDateSite").Text()
      time := s.Find("td.time div.skedStartTimeEST").Text()
      if (date != "") {      
        s.Find("td.skedLinks a.btn").Each(func(i int, s *goquery.Selection) {
          if (s.Text() == "RECAP›") {
            recapLink := s.AttrOr("href", "")
            for i, value := range strings.Split(recapLink, "?id=") {
              if (i==1) {
                id = value
              }
            }
          }
        })
        s.Find("td.team div.teamName").Each(func(i int, s *goquery.Selection) {
          switch i {
          case 0:
            awayTeam = s.Text()
          case 1:
            homeTeam = s.Text()
          }
        })
        if (validTeam(awayTeam) && validTeam(homeTeam)) {
          g := game{id, date, time, awayTeam, homeTeam, gameSequence{}, gameInfo{} }
          games = append(games, g)
        } else {
          fmt.Println(awayTeam)
          fmt.Println(homeTeam)
        }
      }
    })
    for _, g := range games {
      fmt.Println(g)
    }
    fmt.Printf("Number of games: %d \n", len(games))
  }
}

func stringToInt(s string) int {
  i, err := strconv.Atoi(s)
  if (err != nil) {
    fmt.Println(err)
    return 0
  } else {
    return i
  }
}

func validTeam(a string) bool {
  for _, b := range teams {
    switch {
    case strings.ToLower(b.city) == strings.ToLower(a):
      return true
    case strings.Replace(a,"NY", "New York", 1) == b.city + " " + b.name:
      return true
    }
  }
  return false
}
