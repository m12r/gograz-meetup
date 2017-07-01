package meetupcom

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

type MemberPhoto struct {
	ID          int64
	HighResLink string `json:"highres_link"`
	ThumbLink   string `json:"thumb_link"`
	PhotoLink   string `json:"photo_link"`
}

type RSVPsResponseItemMember struct {
	ID    int64
	Name  string
	Photo MemberPhoto
}

type RSVPsResponseItem struct {
	Respone string                  `json:"response"`
	Member  RSVPsResponseItemMember `json:"member"`
	Guests  int64                   `json:"guests"`
}

type RSVPsResponse []RSVPsResponseItem

func (c *Client) GetRSVPs(ctx context.Context, eventID string, urlName string) (*RSVPsResponse, error) {
	var opts url.Values = make(url.Values)
	r, err := c.executeGet(ctx, fmt.Sprintf("/%s/events/%s/rsvps", url.PathEscape(urlName), url.PathEscape(eventID)), opts)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	var resp RSVPsResponse
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
