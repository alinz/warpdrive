# Apps API

## Creating a new App

by default whoever creates a new app, automatically becomes agent

```
POST      /apps

{
  name: "My Awesome App"
}
```

## Getting list of all Apps

```
GET       /apps

[
  {
    app_id: 1,
    app_name: "My Awesome App",
    permission: "member|agent|no_access"
  },
  ...
]
```

## Change the app's name

```
PATCH     /apps/:appId

{
  name: "New Awesome App"
}
```
``
