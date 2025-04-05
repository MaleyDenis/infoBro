# InfoBro UI Preview

## Dashboard Page
The main dashboard displays the latest tech news from all sources, with filtering capabilities.

```
+-----------------------------------------------------------------------+
|                          InfoBro                 [Refresh News]        |
| [Dashboard] [Sources] [Settings]                                       |
+-----------------------------------------------------------------------+
|                                                                       |
|  Latest Tech News                       |  News by Source             |
|                                         |  +--------------------+     |
|  [Filter News]                          |  |  [PIE CHART]       |     |
|                                         |  |  Reddit: 45%        |     |
|  +---------------------------+          |  |  Telegram: 30%      |     |
|  | Reddit • r/programming    |          |  |  RSS: 25%           |     |
|  | • 2 hours ago             |          |  +--------------------+     |
|  |                           |          |                             |
|  | Go 1.21 Version Released  |          |  Popular Tags               |
|  |                           |          |  +--------------------+     |
|  | The Go development team   |          |  | [go] [programming]  |     |
|  | announced the release...  |          |  | [javascript] [react]|     |
|  |                           |          |  | [mongodb] [python]  |     |
|  | [reddit] [Read original →]|          |  | [ai] [cloud]        |     |
|  +---------------------------+          |  +--------------------+     |
|                                         |                             |
|  +---------------------------+          |                             |
|  | Telegram • Golang News    |          |                             |
|  | • 5 hours ago             |          |                             |
|  |                           |          |                             |
|  | Introducing New Feature   |          |                             |
|  |                           |          |                             |
|  | The latest update brings  |          |                             |
|  | support for generics...   |          |                             |
|  |                           |          |                             |
|  | [telegram] [Read original →]|        |                             |
|  +---------------------------+          |                             |
|                                         |                             |
+-----------------------------------------------------------------------+
|                 © 2025 InfoBro - Personal Technical News Dashboard    |
+-----------------------------------------------------------------------+
```

## News Detail Page
Displays a full article with detailed information.

```
+-----------------------------------------------------------------------+
|                          InfoBro                 [Refresh News]        |
| [Dashboard] [Sources] [Settings]                                       |
+-----------------------------------------------------------------------+
|                                                                       |
|  [← Back to News]                                                     |
|                                                                       |
|  +-------------------------------------------------------------------+|
|  |                                                                   ||
|  |  [reddit] • r/golang • April 2, 2025, 3:30 PM                     ||
|  |                                                                   ||
|  |  Go 1.21 Version Released                                         ||
|  |                                                                   ||
|  |  The Go development team announced the release of a new version   ||
|  |  with significant improvements to the type system and runtime     ||
|  |  performance. This version introduces several new features        ||
|  |  requested by the community:                                      ||
|  |                                                                   ||
|  |  - Improved generics support                                      ||
|  |  - Enhanced error handling                                        ||
|  |  - Better performance for concurrent operations                   ||
|  |  - New standard library additions                                 ||
|  |                                                                   ||
|  |  The language continues to evolve while maintaining its core      ||
|  |  philosophy of simplicity and efficiency...                       ||
|  |                                                                   ||
|  |  +---------------------------------------------------------------+||
|  |  |                                                               |||
|  |  |  [Read Original Article]     View Source: r/golang →          |||
|  |  +---------------------------------------------------------------+||
|  +-------------------------------------------------------------------+|
|                                                                       |
+-----------------------------------------------------------------------+
|                 © 2025 InfoBro - Personal Technical News Dashboard    |
+-----------------------------------------------------------------------+
```

## Sources Page
Shows all available news sources and allows manually refreshing them.

```
+-----------------------------------------------------------------------+
|                          InfoBro                 [Refresh News]        |
| [Dashboard] [Sources] [Settings]                                       |
+-----------------------------------------------------------------------+
|                                                                       |
|  News Sources                                                         |
|                                                                       |
|  +-------------------------+  +-------------------------+             |
|  | [R] Reddit              |  | [T] Telegram            |             |
|  | News from popular       |  | News from tech Telegram |             |
|  | technology subreddits   |  | channels                |             |
|  |                         |  |                         |             |
|  | [reddit]     [Fetch Now]|  | [telegram]   [Fetch Now]|             |
|  |                         |  |                         |             |
|  | Available Sources:      |  | Available Sources:      |             |
|  | • r/golang              |  | • Golang News           |             |
|  | • r/programming         |  | • Rust Language         |             |
|  | • r/rust                |  | • Python Insider        |             |
|  | • r/MachineLearning     |  |                         |             |
|  +-------------------------+  +-------------------------+             |
|                                                                       |
|  +-------------------------+                                          |
|  | [R] RSS                 |                                          |
|  | News from RSS feeds     |                                          |
|  |                         |                                          |
|  | [rss]        [Fetch Now]|                                          |
|  |                         |                                          |
|  | Available Sources:      |                                          |
|  | • Hacker News           |                                          |
|  | • The Verge             |                                          |
|  | • DEV Community         |                                          |
|  |                         |                                          |
|  +-------------------------+                                          |
|                                                                       |
+-----------------------------------------------------------------------+
|                 © 2025 InfoBro - Personal Technical News Dashboard    |
+-----------------------------------------------------------------------+
```

## Settings Page
Allows customizing the application behavior and appearance.

```
+-----------------------------------------------------------------------+
|                          InfoBro                 [Refresh News]        |
| [Dashboard] [Sources] [Settings]                                       |
+-----------------------------------------------------------------------+
|                                                                       |
|  Settings                                                             |
|                                                                       |
|  +-------------------------------------------------------------------+|
|  |                                                                   ||
|  |  Appearance                                                       ||
|  |  Dark Mode                                             [Toggler]  ||
|  |                                                                   ||
|  |  News Preferences                                                 ||
|  |  Auto-refresh Interval (minutes)                      [Dropdown]  ||
|  |  Notifications                                        [Toggler]   ||
|  |                                                                   ||
|  |  News Sources                                                     ||
|  |  Reddit                                               [Toggler]   ||
|  |  Telegram                                             [Toggler]   ||
|  |  RSS                                                  [Toggler]   ||
|  |                                                                   ||
|  |  [                         Save Settings                         ]||
|  |                                                                   ||
|  +-------------------------------------------------------------------+|
|                                                                       |
+-----------------------------------------------------------------------+
|                 © 2025 InfoBro - Personal Technical News Dashboard    |
+-----------------------------------------------------------------------+
```