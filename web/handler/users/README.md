# Users API

## Creating a new user

```
POST      /users

{
  name: "Ali",
  email: "ali@pressly.com",
  password: "1234"
}
```

## Changing Permissions

```
PATCH     /users/:userId/permissions

[
  {
    app_id: 1,
    access_type: "agent|member"
  }
]
```

## Getting list of all users, no pagination yet

```
GET       /users

[
  {
    id: 1,
    name: "Ali",
    email: "ali@pressly.com",
    permissions: [
      {
        app_id: 1,
        app_name: "App1",
        access_type: "agent"
      },
      ...
    ]
  },
  ...
]
```

## Update user

```
PATCH     /users/:userId

{
  name: "Ali",
  email: "ali@pressly.com",
  password: "" <- optional
}
```

## Delete user

```
DELETE    /users/:userId
```
