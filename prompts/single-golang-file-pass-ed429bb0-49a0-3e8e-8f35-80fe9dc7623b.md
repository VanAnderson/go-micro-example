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



### service_test.go
```
package user_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/sksmith/go-micro-example/core"
	"github.com/sksmith/go-micro-example/core/user"
	"github.com/sksmith/go-micro-example/db/usrrepo"
)

func TestGet(t *testing.T) {
	usr := user.User{Username: "someuser", HashedPassword: "somehashedpassword", IsAdmin: false, Created: time.Now()}
	tests := []struct {
		name     string
		username string

		getFunc func(ctx context.Context, username string, options ...core.QueryOptions) (user.User, error)

		wantUser user.User
		wantErr  bool
	}{
		{
			name:     "user is returned",
			username: "someuser",

			getFunc: func(ctx context.Context, username string, options ...core.QueryOptions) (user.User, error) {
				return usr, nil
			},

			wantUser: usr,
		},
		{
			name:     "error is returned",
			username: "someuser",

			getFunc: func(ctx context.Context, username string, options ...core.QueryOptions) (user.User, error) {
				return user.User{}, errors.New("some unexpected error")
			},

			wantErr:  true,
			wantUser: user.User{},
		},
	}

	for _, test := range tests {
		mockRepo := usrrepo.NewMockRepo()
		if test.getFunc != nil {
			mockRepo.GetFunc = test.getFunc
		}

		service := user.NewService(mockRepo)

		t.Run(test.name, func(t *testing.T) {
			got, err := service.Get(context.Background(), test.username)
			if test.wantErr && err == nil {
				t.Errorf("expected error, got none")
			} else if !test.wantErr && err != nil {
				t.Errorf("did not want error, got=%v", err)
			}

			if !reflect.DeepEqual(got, test.wantUser) {
				t.Errorf("unexpected user\n got=%+v\nwant=%+v", got, test.wantUser)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	tests := []struct {
		name    string
		request user.CreateUserRequest

		createFunc func(ctx context.Context, user *user.User, tx ...core.UpdateOptions) error

		wantUsername    string
		wantRepoCallCnt map[string]int
		wantErr         bool
	}{
		{
			name:    "user is returned",
			request: user.CreateUserRequest{Username: "someuser", IsAdmin: false, PlainTextPassword: "plaintextpw"},

			wantRepoCallCnt: map[string]int{"Create": 1},
			wantUsername:    "someuser",
		},
	}

	for _, test := range tests {
		mockRepo := usrrepo.NewMockRepo()
		if test.createFunc != nil {
			mockRepo.CreateFunc = test.createFunc
		}

		service := user.NewService(mockRepo)

		t.Run(test.name, func(t *testing.T) {
			got, err := service.Create(context.Background(), test.request)
			if test.wantErr && err == nil {
				t.Errorf("expected error, got none")
			} else if !test.wantErr && err != nil {
				t.Errorf("did not want error, got=%v", err)
			}

			if got.Username != test.wantUsername {
				t.Errorf("unexpected username got=%+v want=%+v", got.Username, test.wantUsername)
			}

			for f, c := range test.wantRepoCallCnt {
				mockRepo.VerifyCount(f, c, t)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name     string
		username string

		deleteFunc func(ctx context.Context, username string, tx ...core.UpdateOptions) error

		wantRepoCallCnt map[string]int
		wantErr         bool
	}{
		{
			name:            "user is deleted",
			username:        "someuser",
			wantRepoCallCnt: map[string]int{"Delete": 1},
		},
		{
			name:     "error is returned",
			username: "someuser",

			deleteFunc: func(ctx context.Context, username string, tx ...core.UpdateOptions) error {
				return errors.New("some unexpected error")
			},

			wantRepoCallCnt: map[string]int{"Delete": 1},
			wantErr:         true,
		},
	}

	for _, test := range tests {
		mockRepo := usrrepo.NewMockRepo()
		if test.deleteFunc != nil {
			mockRepo.DeleteFunc = test.deleteFunc
		}

		service := user.NewService(mockRepo)

		t.Run(test.name, func(t *testing.T) {
			err := service.Delete(context.Background(), test.username)
			if test.wantErr && err == nil {
				t.Errorf("expected error, got none")
			} else if !test.wantErr && err != nil {
				t.Errorf("did not want error, got=%v", err)
			}

			for f, c := range test.wantRepoCallCnt {
				mockRepo.VerifyCount(f, c, t)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	usr := user.User{Username: "someuser", HashedPassword: "$2a$10$t67eB.bOkZGovKD8wqqppO7q.SqWwTS8FUrUx3GAW57GMhkD2Zcwy", IsAdmin: false, Created: time.Now()}
	tests := []struct {
		name     string
		username string
		password string

		getFunc func(ctx context.Context, username string, options ...core.QueryOptions) (user.User, error)

		wantUsername string
		wantErr      bool
	}{
		{
			name:     "correct password",
			username: "someuser",
			password: "plaintextpw",

			getFunc: func(ctx context.Context, username string, options ...core.QueryOptions) (user.User, error) {
				return usr, nil
			},

			wantUsername: "someuser",
		},
		{
			name:     "wrong password",
			username: "someuser",
			password: "wrongpw",

			getFunc: func(ctx context.Context, username string, options ...core.QueryOptions) (user.User, error) {
				return usr, nil
			},

			wantErr:      true,
			wantUsername: "",
		},
		{
			name:     "unexpected error getting user",
			username: "someuser",
			password: "wrongpw",

			getFunc: func(ctx context.Context, username string, options ...core.QueryOptions) (user.User, error) {
				return user.User{}, errors.New("some unexpected error")
			},

			wantErr:      true,
			wantUsername: "",
		},
	}

	for _, test := range tests {
		mockRepo := usrrepo.NewMockRepo()
		if test.getFunc != nil {
			mockRepo.GetFunc = test.getFunc
		}

		service := user.NewService(mockRepo)

		t.Run(test.name, func(t *testing.T) {
			got, err := service.Login(context.Background(), test.username, test.password)
			if test.wantErr && err == nil {
				t.Errorf("expected error, got none")
			} else if !test.wantErr && err != nil {
				t.Errorf("did not want error, got=%v", err)
			}

			if got.Username != test.wantUsername {
				t.Errorf("unexpected username got=%v want=%v", got.Username, test.wantUsername)
			}
		})
	}
}

```
