package worker
import (
    "encoding/json"
    "log"
    "time"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/port/service"
)

type EmailNotificationWorker struct {
    redisService       service.RedisService
    bookingEmailService service.BookingEmailService
    interval           time.Duration
    running            bool
    stopChan           chan struct{}
}

func NewEmailNotificationWorker(
    redisService service.RedisService,
    bookingEmailService service.BookingEmailService,
    interval time.Duration,
) *EmailNotificationWorker {
    return &EmailNotificationWorker{
        redisService:       redisService,
        bookingEmailService: bookingEmailService,
        interval:           interval,
        stopChan:           make(chan struct{}),
    }
}


func (w *EmailNotificationWorker) Start() {
    if w.running {
        return
    }

    w.running = true
    log.Println("Email notification worker started")
    go w.run()
}

func (w *EmailNotificationWorker) Stop() {
    if !w.running {
        return
    }

    w.running = false
    w.stopChan <- struct{}{}
    log.Println("Email notification worker stopped")
}

func (w *EmailNotificationWorker) run() {
    ticker := time.NewTicker(w.interval)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            w.processFlightCancellationNotifications()
            w.processFlightDelayNotifications()
        case <-w.stopChan:
            return
        }
    }
}
func (w *EmailNotificationWorker) processFlightCancellationNotifications() {
    keys, err := w.redisService.Keys("notification:cancellation:*")
    if err != nil {
        log.Printf("Error fetching cancellation notification keys: %v", err)
        return
    }

    if len(keys) == 0 {
        return
    }

    log.Printf("Found %d flight cancellation notifications to process", len(keys))

    for _, key := range keys {
        jsonData, err := w.redisService.Get(key)
        if err != nil {
            log.Printf("Error fetching notification data for key %s: %v", key, err)
            continue
        }

        var notification struct {
            BookingID int64  `json:"booking_id"`
            Reason    string `json:"reason"`
        }

        if err := json.Unmarshal([]byte(jsonData), &notification); err != nil {
            log.Printf("Error parsing notification data for key %s: %v", key, err)
            continue
        }

        response := w.bookingEmailService.SendFlightCancelEmail(notification.BookingID, notification.Reason)
        if response.Status {
            if err := w.redisService.Del(key); err != nil {
                log.Printf("Error deleting processed notification key %s: %v", key, err)
            }
        } else {
            log.Printf("Failed to send cancellation email for booking %d: %v", notification.BookingID, response.ErrorMessage)
        }

        time.Sleep(100 * time.Millisecond)
    }
}
func (w *EmailNotificationWorker) processFlightDelayNotifications() {
    keys, err := w.redisService.Keys("notification:delay:*")
    if err != nil {
        log.Printf("Error fetching delay notification keys: %v", err)
        return
    }

    if len(keys) == 0 {
        return
    }

    log.Printf("Found %d flight delay notifications to process", len(keys))

    for _, key := range keys {
        jsonData, err := w.redisService.Get(key)
        if err != nil {
            log.Printf("Error fetching notification data for key %s: %v", key, err)
            continue
        }

        var notification struct {
            BookingID        int64  `json:"booking_id"`
            NewDepartureTime int64  `json:"new_departure_time"`
            Reason           string `json:"reason"`
        }

        if err := json.Unmarshal([]byte(jsonData), &notification); err != nil {
            log.Printf("Error parsing notification data for key %s: %v", key, err)
            continue
        }

        newDepartureTime := time.Unix(notification.NewDepartureTime, 0)
        response := w.bookingEmailService.SendFlightDelayEmail(notification.BookingID, newDepartureTime, notification.Reason)
        if response.Status {
            if err := w.redisService.Del(key); err != nil {
                log.Printf("Error deleting processed notification key %s: %v", key, err)
            }
        } else {
            log.Printf("Failed to send delay email for booking %d: %v", notification.BookingID, response.ErrorMessage)
        }

        time.Sleep(100 * time.Millisecond)
    }
}