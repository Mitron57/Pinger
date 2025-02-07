package http

import (
    "bytes"
    "encoding/json"
    "errors"
    "go.uber.org/zap"
    "io"
    "net/http"
    "net/url"
    "pinger/internal/domain/interfaces/handlers"
    "pinger/internal/domain/models"
)

type errorMessage struct {
    message string
}

type Streamer struct {
    client *http.Client
    url    *url.URL
    logger *zap.Logger
}

func NewHttpStreamer(url *url.URL, logger *zap.Logger) handlers.Streamer {
    return Streamer{
        client: http.DefaultClient,
        url:    url,
        logger: logger,
    }
}

func (h Streamer) Send(machine *models.Machine) error {
    body, err := json.Marshal(machine)
    if err != nil {
        return err
    }
    request := h.buildRequest(body)
    response, err := h.client.Do(&request)
    if err != nil {
        return err
    }
    defer response.Body.Close()
    if response.StatusCode != http.StatusOK {
        var msg errorMessage
        err = json.NewDecoder(response.Body).Decode(&msg)
        if err != nil {
            return err
        }
        return errors.New(msg.message)
    }
    h.logger.Info("Sent info about machine", zap.String("machine", machine.IP))
    return nil
}

func (h Streamer) buildRequest(body []byte) http.Request {
    return http.Request{
        Header: http.Header{
            "Content-Type": []string{"application/json"},
        },
        Method: http.MethodPut,
        URL:    h.url,
        Body:   io.NopCloser(bytes.NewBuffer(body)),
    }
}
