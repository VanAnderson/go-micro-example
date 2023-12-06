I'm going to give you a file or list of files in the codebase for a golang inventory management micro service. Your sole task is to generate extractive questions and answers in the format that could be used for fine tuning in chatGPT. Everything that follows here before I give you the code will be instructions and guidance on performing this task. I want you to give me the output in pure JSONL with no further commentary. Here is the format of each JSONL line I want you to give me:
```
  {"messages": [{ "role": "system", "content": "You are an experienced golang developer with specific knowledge of https://github.com/VanAnderson/go-micro-example. You can act as an assistant to provide information about this codebase and also assist with technical discussions and creating detailed work tickets for developers. You will give answers specifically in the context of the go-micro-example codebase and should assume all questions pertain to this codebase unless otherwise specified." },{ "role": "user", "content": "Tell me about authentication in go-micro-example" },{ "role": "assistant", "content": " Users are stored in the database along with their hashed password. Users are locally cached using [golang-lru](https://github.com/hashicorp/golang-lru)." }] }
```
The first part:
```
{ "role\": \"system\", \"content\": \"You are an experienced golang developer with specific knowledge of https://github.com/VanAnderson/go-micro-example (or simply go-micro-example). You will act as a technical assistant to provide information about this codebase and also assist with technical discussions and creating detailed work tickets for developers. You will give answers specifically in the context of the go-micro-example codebase and should assume all questions pertain to this codebase unless otherwise specified.\" }
```
Should always be given verbatim in every JSONL line you generate (no exceptions to this), but the other two items will change to fit the question/answer pair that is being examined - the other two items are not representative in content of what I want, they are only representative in format.
*Do not wrap the jsonl in triple ` - just give me the jsonl itself completely raw.*
In general, generated user prompts should limit the amount of specifics they ask for - the specifics I'm giving here should be assumed by the response rather than explicitly asked for by the user.

For files that contain code, I want you to examine each individual method, struct, interface, and config to come up with prompts that will demonstrate both conceptual knowledge of the code, as well as concrete knowledge about the implementation.
Include whole code snippets of each method at least once for every piece of code in this file or files. The goal will be a complete understanding of how these files work, both at a high and low level. Every line of code should be accounted for. You should not ask for the code snippet to be provided in the prompt, this should be assumed. 
Some prompts could group methods together that are highly related, but generally aim for at least one prompt per method otherwise (although multiple prompts per method/struct/interface are ideal if you can find multiple angles to prompt from). Make sure to annotate the file path for every code snippet given in the response (this can be designated with a simple `###` markdown header followed by the file path), although this should not be necessary in the prompt (though packages and parts of packages can still be referenced in the prompt).

Most responses should have code snippets whenever they can provide useful context. Forego code snippets only when dealing with purely conceptual topics. Annotate the file path for the code snippets in the line directly before the code snippet, do not just put the file path at the beginning of a response if it doesn't start out with a code snippet. So every time you give a code snippet it should be:

[descriptive text about whatever]

### path/to/file.go
```golang
[codesnippet]
```

Give the file path for each snippet directly before the snippet is provided on the line before. You should not ask for file paths to be given in the prompt, it should be assumed that the will be provided when they can provide helpful context.

**Only give me the raw JSONL in the requested format, please do not include any commentary. If explicitly instructed to, you may return a blank response. **

# here are the files I want you to analyze:



### reservationapi_test.go
```
package api_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/gobwas/ws"
	"github.com/sksmith/go-micro-example/api"
	"github.com/sksmith/go-micro-example/core"
	"github.com/sksmith/go-micro-example/core/inventory"
	"github.com/sksmith/go-micro-example/testutil"
)

func TestReservationSubscribe(t *testing.T) {
	mockSvc := inventory.NewMockReservationService()

	subscribeCalled := false
	expectedSubId := inventory.ReservationsSubID("subid1")
	unsubscribeCalled := false

	mockSvc.SubscribeReservationsFunc = func(ch chan<- inventory.Reservation) (id inventory.ReservationsSubID) {
		subscribeCalled = true
		go func() {
			res := getTestReservations()
			for i := 0; i < 3; i++ {
				ch <- res[i]
			}
			close(ch)
		}()

		return expectedSubId
	}

	mockSvc.UnsubscribeReservationsFunc = func(id inventory.ReservationsSubID) {
		unsubscribeCalled = true
	}

	resApi := api.NewReservationApi(mockSvc)
	r := chi.NewRouter()
	resApi.ConfigureRouter(r)
	ts := httptest.NewServer(r)
	defer ts.Close()

	url := strings.Replace(ts.URL, "http", "ws", 1) + "/subscribe"

	conn, _, _, err := ws.DefaultDialer.Dial(context.Background(), url)
	if err != nil {
		t.Fatal(err)
	}

	curRes := getTestReservations()
	for i := 0; i < 3; i++ {
		got := &inventory.Reservation{}
		testutil.ReadWs(conn, got, t)

		reflect.DeepEqual(got, curRes[i])
	}

	if !subscribeCalled {
		t.Errorf("subscribe never called")
	}

	if !unsubscribeCalled {
		t.Errorf("unsubscribe never called")
	}
}

func TestReservationGet(t *testing.T) {
	ts, mockResSvc := setupReservationTestServer()
	defer ts.Close()

	tests := []struct {
		getReservationFunc func(ctx context.Context, ID uint64) (inventory.Reservation, error)
		ID                 string
		wantResponse       *api.ReservationResponse
		wantErr            *api.ErrResponse
		wantStatusCode     int
	}{
		{
			getReservationFunc: func(ctx context.Context, ID uint64) (inventory.Reservation, error) {
				return getTestReservations()[0], nil
			},
			ID:             "1",
			wantResponse:   &api.ReservationResponse{Reservation: getTestReservations()[0]},
			wantErr:        nil,
			wantStatusCode: http.StatusOK,
		},
		{
			getReservationFunc: func(ctx context.Context, ID uint64) (inventory.Reservation, error) {
				return inventory.Reservation{}, core.ErrNotFound
			},
			ID:             "1",
			wantResponse:   nil,
			wantErr:        api.ErrNotFound,
			wantStatusCode: http.StatusNotFound,
		},
		{
			getReservationFunc: func(ctx context.Context, ID uint64) (inventory.Reservation, error) {
				return inventory.Reservation{}, errors.New("some unexpected error")
			},
			ID:             "1",
			wantResponse:   nil,
			wantErr:        api.ErrInternalServer,
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		mockResSvc.GetReservationFunc = test.getReservationFunc

		url := ts.URL + "/" + test.ID
		res, err := http.Get(url)
		if err != nil {
			t.Fatal(err)
		}

		if res.StatusCode != test.wantStatusCode {
			t.Errorf("status code got=%d want=%d", res.StatusCode, test.wantStatusCode)
		}

		if test.wantErr == nil {
			got := api.ReservationResponse{}
			testutil.Unmarshal(res, &got, t)

			if !reflect.DeepEqual(got, *test.wantResponse) {
				t.Errorf("reservation\n got=%+v\nwant=%+v", got, *test.wantResponse)
			}
		} else {
			got := &api.ErrResponse{}
			testutil.Unmarshal(res, got, t)

			if got.StatusText != test.wantErr.StatusText {
				t.Errorf("status text got=%s want=%s", got.StatusText, test.wantErr.StatusText)
			}
			if got.ErrorText != test.wantErr.ErrorText {
				t.Errorf("error text got=%s want=%s", got.ErrorText, test.wantErr.ErrorText)
			}
		}
	}
}

func TestReservationCreate(t *testing.T) {
	ts, mockResSvc := setupReservationTestServer()
	defer ts.Close()

	tests := []struct {
		reserveFunc    func(ctx context.Context, rr inventory.ReservationRequest) (inventory.Reservation, error)
		request        *api.ReservationRequest
		wantResponse   *api.ReservationResponse
		wantErr        *api.ErrResponse
		wantStatusCode int
	}{
		{
			reserveFunc: func(ctx context.Context, rr inventory.ReservationRequest) (inventory.Reservation, error) {
				return getTestReservations()[0], nil
			},
			request:        createReservationRequest("requestid1", "requester1", "sku1", 1),
			wantResponse:   &api.ReservationResponse{Reservation: getTestReservations()[0]},
			wantErr:        nil,
			wantStatusCode: http.StatusCreated,
		},
		{
			reserveFunc: func(ctx context.Context, rr inventory.ReservationRequest) (inventory.Reservation, error) {
				return inventory.Reservation{}, core.ErrNotFound
			},
			request:        createReservationRequest("requestid1", "requester1", "sku1", 1),
			wantResponse:   nil,
			wantErr:        api.ErrNotFound,
			wantStatusCode: http.StatusNotFound,
		},
		{
			reserveFunc: func(ctx context.Context, rr inventory.ReservationRequest) (inventory.Reservation, error) {
				return inventory.Reservation{}, errors.New("some unexpected error")
			},
			request:        createReservationRequest("requestid1", "requester1", "sku1", 1),
			wantResponse:   nil,
			wantErr:        api.ErrInternalServer,
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		mockResSvc.ReserveFunc = test.reserveFunc

		url := ts.URL
		res := testutil.Put(url, test.request, t)

		if res.StatusCode != test.wantStatusCode {
			t.Errorf("status code got=%d want=%d", res.StatusCode, test.wantStatusCode)
		}

		if test.wantErr == nil {
			got := api.ReservationResponse{}
			testutil.Unmarshal(res, &got, t)

			if !reflect.DeepEqual(got, *test.wantResponse) {
				t.Errorf("reservation\n got=%+v\nwant=%+v", got, *test.wantResponse)
			}
		} else {
			got := &api.ErrResponse{}
			testutil.Unmarshal(res, got, t)

			if got.StatusText != test.wantErr.StatusText {
				t.Errorf("status text got=%s want=%s", got.StatusText, test.wantErr.StatusText)
			}
			if got.ErrorText != test.wantErr.ErrorText {
				t.Errorf("error text got=%s want=%s", got.ErrorText, test.wantErr.ErrorText)
			}
		}
	}
}

func TestReservationList(t *testing.T) {
	ts, mockResSvc := setupReservationTestServer()
	defer ts.Close()

	tests := []struct {
		getReservationsFunc func(ctx context.Context, options inventory.GetReservationsOptions, limit int, offset int) ([]inventory.Reservation, error)
		url                 string
		wantResponse        interface{}
		wantStatusCode      int
	}{
		{
			getReservationsFunc: func(ctx context.Context, options inventory.GetReservationsOptions, limit int, offset int) ([]inventory.Reservation, error) {
				if options.Sku != "" {
					t.Errorf("sku got=%s want=%s", options.Sku, "")
				}
				if options.State != inventory.None {
					t.Errorf("state got=%s want=%s", options.State, inventory.None)
				}
				if limit != 50 {
					t.Errorf("limit got=%d want=%d", limit, 50)
				}
				if offset != 0 {
					t.Errorf("offset got=%d want=%d", offset, 0)
				}
				return getTestReservations(), nil
			},
			url:            ts.URL,
			wantResponse:   getTestReservationResponses(),
			wantStatusCode: http.StatusOK,
		},
		{
			getReservationsFunc: func(ctx context.Context, options inventory.GetReservationsOptions, limit int, offset int) ([]inventory.Reservation, error) {
				if options.Sku != "somesku" {
					t.Errorf("sku got=%s want=%s", options.Sku, "somesku")
				}
				if options.State != inventory.Open {
					t.Errorf("state got=%s want=%s", options.State, inventory.Open)
				}
				if limit != 111 {
					t.Errorf("limit got=%d want=%d", limit, 111)
				}
				if offset != 222 {
					t.Errorf("offset got=%d want=%d", offset, 0)
				}
				return getTestReservations(), nil
			},
			url:            ts.URL + "?sku=somesku&state=Open&limit=111&offset=222",
			wantResponse:   getTestReservationResponses(),
			wantStatusCode: http.StatusOK,
		},
		{
			getReservationsFunc: func(ctx context.Context, options inventory.GetReservationsOptions, limit int, offset int) ([]inventory.Reservation, error) {
				if options.State != inventory.Closed {
					t.Errorf("state got=%s want=%s", options.State, inventory.Closed)
				}
				return getTestReservations(), nil
			},
			url:            ts.URL + "?state=Closed",
			wantResponse:   getTestReservationResponses(),
			wantStatusCode: http.StatusOK,
		},
		{
			getReservationsFunc: nil,
			url:                 ts.URL + "?state=SomeInvalidState",
			wantResponse:        api.ErrInvalidRequest(errors.New("invalid state")),
			wantStatusCode:      http.StatusBadRequest,
		},
		{
			getReservationsFunc: func(ctx context.Context, options inventory.GetReservationsOptions, limit int, offset int) ([]inventory.Reservation, error) {
				return []inventory.Reservation{}, core.ErrNotFound
			},
			url:            ts.URL,
			wantResponse:   api.ErrNotFound,
			wantStatusCode: http.StatusNotFound,
		},
		{
			getReservationsFunc: func(ctx context.Context, options inventory.GetReservationsOptions, limit int, offset int) ([]inventory.Reservation, error) {
				return []inventory.Reservation{}, nil
			},
			url:            ts.URL + "?sku=someunknownsku",
			wantResponse:   convertReservationsToResponse([]inventory.Reservation{}),
			wantStatusCode: http.StatusOK,
		},
		{
			getReservationsFunc: func(ctx context.Context, options inventory.GetReservationsOptions, limit int, offset int) ([]inventory.Reservation, error) {
				return []inventory.Reservation{}, errors.New("some unexpected error")
			},
			url:            ts.URL,
			wantResponse:   api.ErrInternalServer,
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		mockResSvc.GetReservationsFunc = test.getReservationsFunc

		res, err := http.Get(test.url)
		if err != nil {
			t.Fatal(err)
		}

		if res.StatusCode != test.wantStatusCode {
			t.Errorf("status code got=%d want=%d", res.StatusCode, test.wantStatusCode)
		}

		if test.wantStatusCode == http.StatusBadRequest ||
			test.wantStatusCode == http.StatusInternalServerError ||
			test.wantStatusCode == http.StatusNotFound {

			want := test.wantResponse.(*api.ErrResponse)
			got := &api.ErrResponse{}
			testutil.Unmarshal(res, got, t)

			if got.StatusText != want.StatusText {
				t.Errorf("status text got=%s want=%s", got.StatusText, want.StatusText)
			}
			if got.ErrorText != want.ErrorText {
				t.Errorf("error text got=%s want=%s", got.ErrorText, want.ErrorText)
			}
		} else {
			want := test.wantResponse.([]api.ReservationResponse)
			got := []api.ReservationResponse{}
			testutil.Unmarshal(res, &got, t)

			if !reflect.DeepEqual(got, want) {
				t.Errorf("reservation\n got=%+v\nwant=%+v", got, want)
			}
		}
	}
}

func createReservationRequest(requestID, requester, sku string, quantity int64) *api.ReservationRequest {
	return &api.ReservationRequest{ReservationRequest: &inventory.ReservationRequest{
		Sku: sku, RequestID: requestID, Requester: requester, Quantity: quantity},
	}
}

func setupReservationTestServer() (*httptest.Server, *inventory.MockReservationService) {
	mockSvc := inventory.NewMockReservationService()
	invApi := api.NewReservationApi(mockSvc)
	r := chi.NewRouter()
	invApi.ConfigureRouter(r)
	ts := httptest.NewServer(r)

	return ts, mockSvc
}

var testReservations = []inventory.Reservation{
	{ID: 1, RequestID: "requestID1", Requester: "requester1", Sku: "sku1", State: inventory.Closed, ReservedQuantity: 1, RequestedQuantity: 1, Created: getTime("2020-01-01T01:01:01Z")},
	{ID: 2, RequestID: "requestID2", Requester: "requester2", Sku: "sku2", State: inventory.Open, ReservedQuantity: 1, RequestedQuantity: 2, Created: getTime("2020-01-01T01:01:01Z")},
	{ID: 3, RequestID: "requestID3", Requester: "requester3", Sku: "sku3", State: inventory.None, ReservedQuantity: 0, RequestedQuantity: 3, Created: getTime("2020-01-01T01:01:01Z")},
}

func getTestReservations() []inventory.Reservation {
	return testReservations
}

func getTestReservationResponses() []api.ReservationResponse {
	responses := []api.ReservationResponse{}

	for _, res := range testReservations {
		responses = append(responses, api.ReservationResponse{Reservation: res})
	}

	return responses
}

func convertReservationsToResponse(reservations []inventory.Reservation) []api.ReservationResponse {
	responses := []api.ReservationResponse{}

	for _, res := range reservations {
		responses = append(responses, api.ReservationResponse{Reservation: res})
	}

	return responses
}

func getTime(t string) time.Time {
	tm, err := time.Parse(time.RFC3339, t)
	if err != nil {
		panic(err)
	}
	return tm
}

```
