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

Docs and Support
-----
The documentation for using hookah is available [here](https://adamshannag.github.io/hookah-docs/)

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



