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

Here are three templates for interface, struct, and function documentation. If and only if this file defines a struct, interface, methods on a struct/interface, or individual function, I want you to generate extractive questions and answers where a user asks for documentation of a particular entity in the context of a package, then reply with the proper documentation with one of these templates.
You should only use the function documentation if the function is not a method on a struct or interface.
**IF THE FILE OR FILES DO NOT CONTAIN A STRUCT OR INTERFACE OR INDIVIDUAL FUNCTION, YOU MAY OUTPUT A BLANK RESPONSE WITH NOTHING IN IT**

```
## Struct Documentation

### `TypeName`

**Description:**

`TypeName` represents a ... It is used to ...

**Fields:**

- `FieldName`: The ... It is used for ...

**Constructors:**

#### `NewTypeName(param1 ParamType1, param2 ParamType2) *TypeName`

Creates a new instance of `TypeName`. It initializes the struct with the given parameters.

**Parameters:**

- `param1` (ParamType1): The first parameter used for ...
- `param2` (ParamType2): The second parameter used for ...

**Returns:**

- A pointer to a newly allocated `TypeName` with fields initialized.

**Methods:**

#### `MethodName(param1 ParamType1, param2 ParamType2) (returnType, error)`

Performs a specific operation on `TypeName`.

**Parameters:**

- `param1` (ParamType1): The first parameter used for ...
- `param2` (ParamType2): The second parameter used for ...

**Returns:**

- `returnType`: Description of the return value.
- `error`: Any error that occurred during the operation, or `nil` if no error occurred.

## Interface Documentation

### `InterfaceName`

**Description:**

`InterfaceName` defines the behavior for ... It is implemented by types that need to ...

**Methods:**

#### `MethodName(param1 ParamType1, param2 ParamType2) (returnType, error)`

Performs an action. It is used to ...

**Parameters:**

- `param1` (ParamType1): The first parameter used for ...
- `param2` (ParamType2): The second parameter used for ...

**Returns:**

- `returnType`: Description of the return value.
- `error`: Any error that occurred during the operation, or `nil` if no error occurred.

### Implementations

#### `ImplementingTypeName`

**Description:**

`ImplementingTypeName` is a struct that implements `InterfaceName`. It has the following properties and behaviors...

**Methods:**

#### `MethodName(param1 ParamType1, param2 ParamType2) (returnType, error)`

Implements the `InterfaceName.MethodName` for `ImplementingTypeName`.

**Parameters:**

- `param1` (ParamType1): The first parameter used for ...
- `param2` (ParamType2): The second parameter used for ...

**Returns:**

- `returnType`: Description of the return value.
- `error`: Any error that occurred during the operation, or `nil` if no error occurred.
## Function Documentation

### `FunctionName`

**Description:**

`FunctionName` performs a specific operation or calculation. It is used to ...

**Parameters:**

- `param1` (ParamType1): Description of the first parameter and its purpose.
- `param2` (ParamType2): Description of the second parameter and its purpose.
- ... (additional parameters as needed)

**Returns:**

- `returnType`: Description of the return value and what it represents.
- `error`: Description of the error returned, if applicable, or `nil` if no error occurred.

**Example:**

```go
result, err := FunctionName(param1, param2)
if err != nil {
    // Handle error
}
// Use result
```

**Notes:**

- Additional notes or considerations regarding the function's behavior, side effects, or usage in certain contexts.
```

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
