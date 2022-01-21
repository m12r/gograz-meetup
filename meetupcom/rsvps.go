package meetupcom

import (
	"context"
	"fmt"
	"strings"
)

type MemberPhoto struct {
	ID        int64
	ThumbLink string `json:"thumb_link"`
	PhotoLink string `json:"photo_link"`
}

type RSVPsResponseItemMember struct {
	ID    string
	Name  string
	Photo MemberPhoto
}

type RSVPsResponseItem struct {
	Response string                  `json:"response"`
	Member   RSVPsResponseItemMember `json:"member"`
	Guests   int64                   `json:"guests"`
}

type RSVPsResponse []RSVPsResponseItem

type graphQLResponse struct {
	Data struct {
		Event struct {
			Tickets struct {
				Count int64 `json:"count"`
				Edges []struct {
					Node struct {
						User struct {
							ID          string `json:"id"`
							MemberURL   string `json:"memberUrl"`
							Name        string `json:"name"`
							MemberPhoto struct {
								ID      string `json:"id"`
								BaseURL string `json:"baseUrl"`
							} `json:"memberPhoto"`
						} `json:"user"`
						GuestsCount int64  `json:"guestsCount"`
						Status      string `json:"status"`
					} `json:"node"`
				} `json:"edges`
			} `json:"tickets"`
		} `json:"event"`
	} `json:"data"`
}

func (c *Client) GetRSVPs(ctx context.Context, eventID string, urlName string) (*RSVPsResponse, error) {
	query := `
  query($eventID: ID) {
    event(id: $eventID) {
      tickets {
        count
        edges {
          node {
            user {
							id
							memberUrl
              name
							memberPhoto {
								baseUrl
								id
							}
            }
            status
          }
        }
      }
    }
  }
  `

	gr := graphQLResponse{}
	if err := c.executeGraphQLQuery(ctx, query, map[string]string{"eventID": eventID}, &gr); err != nil {
		return nil, err
	}
	var resp RSVPsResponse
	for _, ticket := range gr.Data.Event.Tickets.Edges {
		thumbLink := ""
		photoLink := ""
		if ticket.Node.User.MemberPhoto.ID != "0" {
			thumbLink = fmt.Sprintf("%s%s/40x40.jpg", ticket.Node.User.MemberPhoto.BaseURL, ticket.Node.User.MemberPhoto.ID)
			photoLink = fmt.Sprintf("%s%s/100x100.jpg", ticket.Node.User.MemberPhoto.BaseURL, ticket.Node.User.MemberPhoto.ID)
		}
		resp = append(resp, RSVPsResponseItem{
			Response: ticket.Node.Status,
			Member: RSVPsResponseItemMember{
				Name: ticket.Node.User.Name,
				ID:   strings.TrimSuffix(ticket.Node.User.ID, "!chp"),
				Photo: MemberPhoto{
					ThumbLink: thumbLink,
					PhotoLink: photoLink,
				},
			},
			Guests: ticket.Node.GuestsCount,
		})
	}

	return &resp, nil
}
