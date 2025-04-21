# Hookah - literally passes the hook around

**Hookah** is a lightweight, stateless, and zero-dependency webhook router built in Go. It serves as an intermediary
between webhook sources (like GitLab, GitHub, etc.) and target destinations (such as Discord), forwarding events only
when they match predefined conditions.

Features
------

- **Webhook Receiver:** Accepts incoming webhooks from various platforms.
- **Rule Engine:** Applies filters based on request headers/url query params and body content.
- **Conditional Forwarding:** Sends a message to a target webhook only if the rules match.
- **Reusable Templates:** Define multiple templates and reuse them across different configurations and webhook
  scenarios.
- **Template Support:** Allows dynamic message generation using data from the incoming webhook payload.
- **Lightweight & Extensible:** Simple design with future support for multiple rules, formats, and targets.

Environment Variables
------

The server requires the following environment variables:

- `PORT`: the port on which the server should listen (e.g., `8080`)
- `CONFIG_PATH`: path to the JSON config file defining receivers and rules
- `TEMPALTES_PATH`: path to the templates directory that contains all templates

Endpoint Structure
------

All webhooks should be sent to:

```
POST /webhooks/{receiver}
```

- `{receiver}` must match the `receiver` field defined in your config file.
- If multiple configurations share the same `receiver` name, **all matching receivers will be evaluated**.
- Each matching config can independently authorize, filter, and forward events based on its own rules.

Configuration Format
------

Hookah uses a simple JSON-based configuration to route and filter webhooks. Below is a breakdown of the available fields
and how they work.

### Example Configuration

```json
[
  {
    "receiver": "gitlab",
    "auth": {
      "flow": "gitlab",
      "header_secret_key": "X-Gitlab-Token",
      "secret": "my gitlab webhook secret token"
    },
    "event_type_in": "body",
    "event_type_key": "event_type",
    "events": [
      {
        "event": "merge_request",
        "conditions": [
          "{Header.x-gitlab-label} {in} {Body.object_attributes.labels[].title}"
        ],
        "hooks": [
          {
            "name": "discord",
            "endpoint_key": "x-discord-url",
            "body": "discord.tmpl"
          }
        ]
      }
    ]
  }
]
```

Hookah uses a JSON array of receiver configurations. Each configuration defines:

- **Receiver name:** An identifier for incoming webhook sources, can be any name (e.g., `gitlab`, `cool github repo`,
  etc.).
- **Authentication:** rules for verifying webhook authenticity (`auth` block).
- **Event routing rules:** Defines how to extract event types and which hooks to trigger when conditions are met.

### Multiple Receivers

The configuration supports **multiple receivers** — each with its own auth rules, event types, and hook logic. This
enables you to route webhooks from different sources independently:

```json
[
  {
    "receiver": "gitlab",
    ...
  },
  {
    "receiver": "github",
    ...
  }
]
```

### Event Type Resolution

```json
"event_type_in": "body",
"event_type_key": "event_type"
```

These two keys tell Hookah **where to look** for the event type:

- `event_type_in`: `"body"` or `"header"` — defines the source.
- `event_type_key`: the key name to fetch from the source.

Hookah will match the extracted event type against the entries in the `events` array.

### Authentication (`auth`)

Each receiver must define an `auth` block to control who can send webhooks. The supported flows are:

| Flow           | Description                                                                                                                    |
|----------------|--------------------------------------------------------------------------------------------------------------------------------|
| `none`         | No authentication; accepts all requests.                                                                                       |
| `plain secret` | Matches the value in the request header against the configured `secret`.                                                       |
| `basic auth`   | Verifies username and password in basic auth header matches the `secret`, in the format `username:password`.                   |
| `gitlab`       | Compares the configured `secret` with the GitLab token header using constant-time comparison (SHA-512).                        |
| `github`       | Verifies HMAC SHA-256 signature in a header (e.g. X-Hub-Signature-256 or custom) using the configured secret and request body. |

**Fields:**

- `flow`: One of `gitlab`, `github`, `basic auth`, `plain secret`, or `none`.
- `header_secret_key`: The header to extract the token from (e.g., `X-Gitlab-Token` or `X-Custom-Token`).
- `secret`: The expected secret value (or in `basic auth`, the `username:password` pair).

### Events & Conditional Hooks

```json
{
  "events": [
    {
      "event": "merge_request",
      "conditions": [
        "{Header.x-gitlab-label} {in} {Body.object_attributes.labels[].title}"
      ],
      "hooks": [
        ...
      ]
    }
  ]
}
```

Each event:

- Matches incoming requests based on `event_type`.
- Runs all defined `conditions` — if **all pass**, it triggers the corresponding `hooks`.

You can define **multiple hooks** per event to notify different targets like Discord, Slack, etc. All matching hooks
will be triggered concurrently when conditions are satisfied.

### Condition Syntax

Conditions use a simple templated language:

- `{Header.X-Foo}` refers to a request header/url query param
- `{Body.foo.bar}` refers to a nested body field
- `{Body.foo[].bar}` supports iterating over arrays, should be used with the {in} operator

Example:

```json
"{Header.x-gitlab-label} {in} {Body.object_attributes.labels[].title}"
```

This checks whether the value of `x-gitlab-label` header or url query param exists in any of the `title` fields in the
incoming body
array `object_attributes.labels`.

### Hook Structure & Templating

Each `hook` can look like this:

```json
{
  "name": "discord",
  "endpoint_key": "x-discord-url",
  "body": "template_file_name.some_extension"
}
```

- `endpoint_key`: Specifies the request header key, or the url query param that contains the **target webhook URL**.
  this will be used to make the webhook request for the target hook.
- `body`: The name of the template file to use, from the templates' directory.

> Note: After rendering, the template content must result in a well-formed JSON payload, as it will be used in outgoing
> webhook requests.

### Template Usage

In the template files located in the `templates` directory, you can use Go's native templating language.

You may reference values from the original request body using dot notation like `{{ .some.path }}`.  
For example, if the incoming payload contains a field `user.name`, you can access it in your template as:

```gohtml
{{ .user.name }}
```

### Built-in Template Functions

Your templates also support the following built-in utility functions:

| Function    | Description                                                                                                     |
|-------------|-----------------------------------------------------------------------------------------------------------------|
| `now`       | Returns the current time.                                                                                       |
| `format`    | Formats a `time.Time` object using Go's time layout. Example: `{{ format now "2006-01-02" }}`                   |
| `parseTime` | Parses a string into a `time.Time` using the given layout. Example: `{{ parseTime "2023-01-01" "2006-01-02" }}` |
| `pastTense` | Appends `-ed` or `-d` to a word to form the past tense. Example: `{{ pastTense "open" }}` → `opened`            |
| `lower`     | Converts a string to lowercase. Example: `{{ lower "HELLO" }}` → `hello`                                        |
| `upper`     | Converts a string to uppercase. Example: `{{ upper "hello" }}` → `HELLO`                                        |
| `title`     | Converts a string to title case. Example: `{{ title "hello world" }}` → `HELLO WORLD`                           |
| `trim`      | Trims leading and trailing whitespace. Example: `{{ trim "  hello  " }}` → `hello`                              |
| `contains`  | Checks if a string contains a substring. Example: `{{ contains "hello world" "world" }}` → `true`               |
| `replace`   | Replaces all occurrences of a substring. Example: `{{ replace "hello world" "world" "Go" }}` → `hello Go`       |
| `default`   | Returns a fallback value if the input is empty or nil. Example: `{{ default .user.name "Guest" }}`              |

Running with Docker Compose
------

1. Navigate to the `deploy` directory:

   ```bash
   cd deploy
   ```

2. Run the service:

   ```bash
   docker compose up -d
   ```

---

#### Example `curl` Command (payload copied with some modifications from a gitlab merge_request event)

```bash
curl -X POST http://localhost:3000/webhooks/gitlab?discord-url=your_discord_webhook_url_goes_here \
  -H "Content-Type: application/json" \
  -H "x-gitlab-label: API" \
  -d '{
    "object_kind": "merge_request",
    "event_type": "merge_request",
    "user": {
      "name": "Administrator",
      "username": "root",
      "avatar_url": "http://www.gravatar.com/avatar/e64c7d89f26bd1972efa854d13d7dd61?s=40&d=identicon"
    },
    "project": {
      "web_url": "http://example.com/gitlabhq/gitlab-test",
      "path_with_namespace": "gitlabhq/gitlab-test"
    },
    "object_attributes": {
      "title": "MS-Viewport",
      "updated_at": "2013-12-03T17:23:34Z",
      "labels": [
        {
          "title": "API",
          "description": "API related issues"
        }
      ],
      "action": "open"
    }
  }'
```

Why “Hookah”?
------

Much like a real hookah, this tool filters input before releasing output—except in this case, it's webhooks instead of
smoke.


Contributing
------

Contributions are welcome! If you’d like to help improve **Hookah**, feel free to submit an issue or open a pull
request.

License
------
**Hookah** is released under the [MIT License](LICENSE).



