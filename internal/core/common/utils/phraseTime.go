package utils

import (
    "time"
    "fmt"
)

type TimeUtils interface {
    Now() time.Time
    NowUnix() int64
    ParseTime(dateStr string) (time.Time, error)
    FormatTime(t time.Time) string
    ParseDate(dateStr string) (time.Time, error)
}

func ParseTime(dateStr string) (time.Time, error) {

    formats := []string{
        "02/01/2006 15:04",       // dd/mm/yyyy HH:mm
        "02/01/2006",             // dd/mm/yyyy
    }
    
    for _, format := range formats {
        if parsedTime, err := time.Parse(format, dateStr); err == nil {
            return parsedTime, nil
        }
    }
    
    return time.Time{}, fmt.Errorf("invalid datetime format: %s. Supported formats: dd/mm/yyyy HH:mm, dd/mm/yyyy, yyyy-mm-dd HH:mm:ss", dateStr)
}

func FormatTime(t time.Time) string {
    return t.Format("2006-01-02 15:04:05")
}

func ParseDate(dateStr string) (time.Time, error) {
    formats := []string{
        "02/01/2006",    // dd/mm/yyyy
        "02-01-2006",    // dd-mm-yyyy
        "2006-01-02",    // yyyy-mm-dd
    }
    
    for _, format := range formats {
        if parsedTime, err := time.Parse(format, dateStr); err == nil {
            return parsedTime, nil
        }
    }
    
    return time.Time{}, fmt.Errorf("invalid date format: %s. Supported formats: dd/mm/yyyy, yyyy-mm-dd", dateStr)
}

func FormatTimeVN(t time.Time) string {
    return t.Format("02/01/2006 15:04") // dd/mm/yyyy HH:mm
}

func FormatDateVN(t time.Time) string {
    return t.Format("02/01/2006") // dd/mm/yyyy
}