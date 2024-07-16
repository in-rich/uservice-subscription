# Contributions

## Project structure

This project uses a variant of [hexagonal architecture](https://en.wikipedia.org/wiki/Hexagonal_architecture_(software))
(or whatever you like to call it). The goal is to write small, testable, and maintainable components.

### PKG

The main logic of the project is located under the `/pkg` directory. Here are located the main components of the
hexagonal architecture:

- `dao`: Data Access Object. This package provides I/O methods to interact with data, such as databases or
  file systems.
- `service`: Business logic. This package contains the main logic of the application.
- `client`: External dependencies. This package contains the interfaces and adapters for external dependencies, such as
  databases or file systems, if needed. This does not include external dependencies that already provide properly
  interfaced libraries, or if you don't need to mock them for tests.
- `handlers`: An external API to connect to your services.

#### DAO

This package should only contain basic I/O operations. Data should only be validated against storage rules (such as 
unique constraints in a database).

```go
// ✅ Good
func (d *DAO) Create(ctx context.Context, data *Data) error {
    err := d.ORM.Create(data)
    if errors.Is(err, orm.ErrUniqueConstraint) {
        return ErrAlreadyExists
    }
    
    return nil
}

// ❌ Bad
func (d *DAO) Create(ctx context.Context, data *Data) error {
    if err := validateData(data); err != nil {
        return err
    }
    
    // Complex computation
    for _, item := range data.Items {
        // ....
    }
    
    err := d.ORM.Create(data)
    if errors.Is(err, orm.ErrUniqueConstraint) {
        return ErrAlreadyExists
    }
    
    return nil
}
```

Each DAO should only export one method, with straightforward logic. Ideally, it should never be more than 50 lines long.

```go
// ✅ Good
type MyDAO interface {
    DoSomething(ctx context.Context, data *Data) error
}

// ❌ Bad
type MyDAO interface {
    DoSomething(ctx context.Context, data *Data) error
    DoSomethingElse(ctx context.Context, data *Data) error
    DoAnotherThing(ctx context.Context, data *Data) error
}
```

The only dependency of a DAO should be an ORM, or any library used to interact with the specific file system. If a
DAO can work under multiple storage systems (for example, one local and one remote), then it should use a
[client](#clients), with a generic interface for all storage systems.
In general, if available, prefer using interfaces as dependencies, as it makes mocking for tests easier.

```go
// ✅ Good
type MyDAO interface {
    DoSomething(ctx context.Context, data *Data) error
}

type myDAOImpl struct {
    ORM clients.ORM
}

func NewDAOWithThisORM(thisORM this.ORM) MyDAO {
    adapter := clients.NewThisORMAdapter(thisORM)
    return &myDAOImpl{ORM: adapter}
}

func NewDAOWithThatORM(thatORM that.ORM) MyDAO {
    adapter := clients.NewThatORMAdapter(thatORM)
    return &myDAOImpl{ORM: adapter}
}

// ❌ Bad
type MyDAO interface {
    DoSomething(ctx context.Context, data *Data) error
}

type myDAOThisORMImpl struct {
    ORM this.ORM
}

type myDAOThatORMImpl struct {
    ORM that.ORM
}

func NewDAOWithThisORM(thisORM this.ORM) MyDAO {
    return &myDAOThisORMImpl{ORM: adapter}
}

func NewDAOWithThatORM(thatORM that.ORM) MyDAO {
    return &myDAOThatORMImpl{ORM: adapter}
}
```

#### Services

This package should contain the main logic of the application. If you ever need some logic done, it should be there.

A service must be environment-agnostic: it should only use interfaced dependencies. This way, you can easily test it
using mocks.

```go
type MyService interface {
    Exec(ctx context.Context, data *Data) error
}

// ✅ Good

type myServiceImpl struct {
    DAO dao.MyDAOInterface
    Logger log.LoggerInterface
}

func NewService(d dao.MyDAOInterface, l log.LoggerInterface) MyService {
    return &myServiceImpl{DAO: d, Logger: l}
}

// ❌ Bad

type myServiceImpl struct {
    DAO *dao.MyDAOImpl
    Logger *log.LoggerImpl
}

func NewService() MyService {
    return &myServiceImpl{
        DAO: &dao.MyDAOImpl{},
        Logger: &log.LoggerImpl{},
    }
}
```

Each service should only contain one method named `Exec`.

```go
// ✅ Good
type MyService interface {
    Exec(ctx context.Context, data *Data) error
}

// ❌ Bad
type MyService interface {
    DoSomething(ctx context.Context, data *Data) error
    DoSomethingElse(ctx context.Context, data *Data) error
}
```

#### Clients

A client is a simple adapter between one or more external dependencies and your application. It covers:

- Poorly designed / cumbersome libraries that require a simplified interface for usage in the DAO.
- Generic interfaces that are agnostic to multiple dependencies (for example, different storage methods depending on the
  environment).
- Internal libraries with generic logic that can be reused in multiple services.

A client interface should always be fine-tuned to the requirements of the current project. You are developing an
internal tool, not a library for the world.

```go
// ✅ Good
type SecretFileHandler interface {
    GetSecret(ctx context.Context, key string) (string, error)
    RotateSecret(ctx context.Context, key string) error
}

type secretFileHandlerLocalImpl struct {
    fileSystem embed.FS
}

type secretFileHandlerDBImpl struct {
    orm orm.ORM
}

func NewSecretFileHandler(fileSystem embed.FS) SecretFileHandler {
    return &secretFileHandlerLocalImpl{fileSystem: fileSystem}
}

func NewSecretFileHandler(orm orm.ORM) SecretFileHandler {
    return &secretFileHandlerDBImpl{orm: orm}
}
```

#### Handlers

This package is merely an adapter for your services.

A handler creates an interface between your service and the external world. It only serves as a passthrough:
 - Prepare input for the service.
 - Convert the output of the service to a format that the external world can understand.
 - Handle errors

Anything else should be delegated to the service.

```go
type MyHandler interface {
    HandleRequest(ctx context.Context, req *Request) (*Response, error)
}

type myHandlerImpl struct {
    Service service.MyService
}

// ✅ Good
func (h *myHandlerImpl) HandleRequest(ctx context.Context, req *Request) (*Response, error) {
    const input := requestToData(req)

    output, err := h.Service.Exec(ctx, input)
    if err != nil {
        return nil, handleError(err)
    }
    
    return outputToResponse(output), nil
}

// ❌ Bad
func (h *myHandlerImpl) HandleRequest(ctx context.Context, req *Request) (*Response, error) {
    const input := requestToData(req)

    // Complex computation
    for _, item := range input.Items {
        // ....
    }

    output, err := h.Service.Exec(ctx, input)
    if err != nil {
        return nil, handleError(err)
    }
    
    return outputToResponse(output), nil
}
```

The only dependency of a handler should be one (avoid more) service.
