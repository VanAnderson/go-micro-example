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



### userapi_test.go
```
package api_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/sksmith/go-micro-example/api"
	"github.com/sksmith/go-micro-example/core"
	"github.com/sksmith/go-micro-example/core/user"
	"github.com/sksmith/go-micro-example/testutil"
)

func TestUserCreate(t *testing.T) {
	ts, mockSvc := setupUserTestServer()
	defer ts.Close()

	tests := []struct {
		name           string
		loginFunc      func(ctx context.Context, username, password string) (user.User, error)
		createFunc     func(ctx context.Context, user user.CreateUserRequest) (user.User, error)
		url            string
		request        interface{}
		wantResponse   interface{}
		wantStatusCode int
	}{
		{
			name: "admin users can create valid user",
			loginFunc: func(ctx context.Context, username, password string) (user.User, error) {
				return createUser("someadmin", "", true), nil
			},
			createFunc: func(ctx context.Context, usr user.CreateUserRequest) (user.User, error) {
				return createUser(usr.Username, "somepasswordhash", usr.IsAdmin), nil
			},
			url:            ts.URL,
			request:        createUserReq("someuser", "somepass", false),
			wantResponse:   nil,
			wantStatusCode: http.StatusOK,
		},
		{
			name: "non-admin users are unable to create users",
			loginFunc: func(ctx context.Context, username, password string) (user.User, error) {
				return createUser("someadmin", "", false), nil
			},
			createFunc: func(ctx context.Context, usr user.CreateUserRequest) (user.User, error) {
				return createUser(usr.Username, "somepasswordhash", usr.IsAdmin), nil
			},
			url:            ts.URL,
			request:        createUserReq("someuser", "somepass", false),
			wantResponse:   nil,
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name: "when the creating user is not found, server returns unauthorized",
			loginFunc: func(ctx context.Context, username, password string) (user.User, error) {
				return user.User{}, core.ErrNotFound
			},
			createFunc: func(ctx context.Context, usr user.CreateUserRequest) (user.User, error) {
				return createUser(usr.Username, "somepasswordhash", usr.IsAdmin), nil
			},
			url:            ts.URL,
			request:        createUserReq("someuser", "somepass", false),
			wantResponse:   nil,
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name: "when an unexpected error occurs logging in, an internal server error is returned",
			loginFunc: func(ctx context.Context, username, password string) (user.User, error) {
				return user.User{}, errors.New("some unexpected error")
			},
			createFunc:     nil,
			url:            ts.URL,
			request:        createUserReq("someuser", "somepass", false),
			wantResponse:   api.ErrInternalServer,
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name: "when an error occurs creating the user, an internal server error is returned",
			loginFunc: func(ctx context.Context, username, password string) (user.User, error) {
				return createUser("someadmin", "", true), nil
			},
			createFunc: func(ctx context.Context, usr user.CreateUserRequest) (user.User, error) {
				return user.User{}, errors.New("some unexpected error")
			},
			url:            ts.URL,
			request:        createUserReq("someuser", "somepass", false),
			wantResponse:   api.ErrInternalServer,
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockSvc.LoginFunc = test.loginFunc
			mockSvc.CreateFunc = test.createFunc

			res := testutil.Post(test.url, test.request, t, testutil.RequestOptions{Username: "someuser", Password: "somepass"})

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
			}
		})
	}
}

func createUser(username, password string, isAdmin bool) user.User {
	return user.User{Username: username, HashedPassword: password, IsAdmin: isAdmin}
}

func createUserReq(username, password string, isAdmin bool) api.CreateUserRequestDto {
	return api.CreateUserRequestDto{CreateUserRequest: &user.CreateUserRequest{Username: username, IsAdmin: isAdmin}, Password: password}
}

func setupUserTestServer() (*httptest.Server, *user.MockUserService) {
	svc := user.NewMockUserService()
	usrApi := api.NewUserApi(svc)
	r := chi.NewRouter()
	r.With(api.Authenticate(svc)).Route("/", func(r chi.Router) {
		usrApi.ConfigureRouter(r)
	})
	ts := httptest.NewServer(r)

	return ts, svc
}

```
