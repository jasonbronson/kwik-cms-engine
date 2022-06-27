package repositories

import (
	"strconv"
	"strings"
	"time"

	"github.com/jasonbronson/kwik-cms-engine/library/helpers"
	"github.com/jasonbronson/kwik-cms-engine/request/response"
	"gorm.io/gorm"
)

const timeFormat = "2006-01-02 15:04:05"
const dateFormat = "2006-01-02"

func metaBuild(data interface{}, params helpers.DefaultParameters) *response.Response {
	p := response.Pagination{
		PageOffset:  params.PageOffset,
		PageSize:    params.PageSize,
		ResultTotal: params.ResultTotal,
		Total:       params.Total,
	}
	s := response.Sort{
		Sort: params.SortOrder,
	}
	return &response.Response{
		Data: data,
		Meta: response.MetaData{
			p,
			s,
			response.ResponseAction{
				Message: "",
			},
		},
	}
}

func Count(db *gorm.DB, model interface{}) int {
	var count int64
	db.Model(model).Count(&count)
	return int(count)
}
func Published() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("published_at < NOW()")
	}
}
func CategoryPostJoinByID(value string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Joins("JOIN categories_post_links cal on cal.post_id=post.id").Where("cal.category_id = ?", value)
	}
}
func AuthorPostJoinByID(value string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", value)
	}
}
func ByID(value string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", value)
	}
}
func ByStatus(value string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", value)
	}
}

func OrderChipsDesc() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order("chips desc")
	}
}

func FilterShowOnlyPublished(value string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if value == "true" {
			return db.Where("status = 'publish' and published_at < NOW()")
		} else {
			return nil
		}
	}
}
func FilterTitle(value string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("LOWER(title) LIKE ?", "%"+strings.ToLower(value)+"%")
	}
}
func FilterStatus(value string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch value {
		case "published":
			db.Where("publish_date <= now() and status = 'publish'")
			break
		case "scheduled":
			db.Where("publish_date > now() and status = 'publish'")
			break
		case "draft":
			db.Where("status = 'draft'")
			break
		}
		return db
	}
}
func FilterEventName(value string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("LOWER(event_name) LIKE ?", "%"+strings.ToLower(value)+"%")
	}
}
func FilterEventTime(value string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if value == "upcoming" {
			return db.Where("event_start_time > now() and status = 'publish'")
		} else {
			return db.Where("event_start_time < now() and status = 'publish' and has_results = true")
		}
	}
}
func FilterEndedAt(value string) func(db *gorm.DB) *gorm.DB {
	//This is used in series call and could be used in other places
	return func(db *gorm.DB) *gorm.DB {
		if value == "past" {
			return db.Where("ended_at < now() and status = 'publish'")
		} else if value == "current" {
			return db.Where("started_at > now() and status = 'publish'")
		} else if value == "latest" {
			return db.Where("started_at - INTERVAL '7 day' < now() and ended_at > now() and status = 'publish'")
		} else {

			return FilterEndedAtOnDay(value)(db)
		}
	}
}
func FilterEndedAtOnDay(value string) func(db *gorm.DB) *gorm.DB {
	endedAt, _ := time.Parse(dateFormat, value)
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("event_end_date > ? and event_end_date < ? and status = 'publish'", endedAt.Truncate(24*time.Hour), endedAt.AddDate(0, 0, 1).Truncate(24*time.Hour))
	}
}
func FilterStartedAtOnDay(value string) func(db *gorm.DB) *gorm.DB {
	startedAt, _ := time.Parse(dateFormat, value)
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("event_start_date > ? and event_start_date < ? and status = 'publish'", startedAt.Truncate(24*time.Hour), startedAt.AddDate(0, 0, 1).Truncate(24*time.Hour))
	}
}
func FilterScheduledEventDates(dateRange string) func(db *gorm.DB) *gorm.DB {
	started := strings.Split(dateRange, "|")[0]
	ended := strings.Split(dateRange, "|")[1]
	startedAt, _ := time.Parse(dateFormat, started)
	endedAt, _ := time.Parse(dateFormat, ended)
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("date(event_start_time) between ? and ? and status = 'publish'", startedAt.Truncate(24*time.Hour), endedAt.Truncate(24*time.Hour))
	}
}
func FilterArticleDateRange(dateRange string) func(db *gorm.DB) *gorm.DB {
	started := strings.Split(dateRange, "|")[0]
	ended := strings.Split(dateRange, "|")[1]
	startedAt, _ := time.Parse(dateFormat, started)
	endedAt, _ := time.Parse(dateFormat, ended)
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("date(published_at) between ? and ?", startedAt.Truncate(24*time.Hour), endedAt.Truncate(24*time.Hour))
	}
}
func FilterLocation(value string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("event_location ILIKE ?", "%"+value+"%")
	}
}
func FilterGameType(value string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("game_type = ?", value)
	}
}
func FilterMinBuyin(value string) func(db *gorm.DB) *gorm.DB {
	val, _ := strconv.Atoi(value)
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("buy_in >= ?", val)
	}
}
func FilterMaxBuyin(value string) func(db *gorm.DB) *gorm.DB {
	val, _ := strconv.Atoi(value)
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("buy_in <= ?", val)
	}
}
func FilterTrash(value string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status <> ?", "closed")
	}
}
func FilterSeriesByTag(value string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Joins("JOIN tags_lr_series_links l on l.lr_series_id = id").Joins("JOIN tag t on t.id = l.tag_id").Where("t.name = ?", value)
	}
}
func FilterDescription(value string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("LOWER(description) LIKE ?", "%"+strings.ToLower(value)+"%")
	}
}
func FilterName(value string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(value)+"%")
	}
}
func FilterLastNameStartsWith(value string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("LOWER(name) ~ ?", "\\y"+strings.ToLower(value)+"\\S+$").Order("substring(name, '([^[:space:]]+)(?:,|$)')").Group("id")
	}
}
func FilterContent(value string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("LOWER(content) LIKE ?", "%"+strings.ToLower(value)+"%")
	}
}
func FilterUsername(value string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("LOWER(username) LIKE ?", "%"+strings.ToLower(value)+"%")
	}
}
func FilterEmail(value string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("LOWER(email) LIKE ?", "%"+strings.ToLower(value)+"%")
	}
}
func FilterYear(value string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("date_part('year', created_at) = ?", value)
	}
}
func FilterScheduleID(value string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("schedule_id = ?", value)
	}
}
func FilterHasMedia(value string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if value == "true" {
			return db.Where("media_id <> 0")
		} else {
			return nil
		}
	}
}
