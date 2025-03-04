# Marina - Game Launcher

## Description

Marina is an **unofficial** launcher & version manager for the PC ports made by Harbour Masters.

Currently supported games include:

- Ship of Harkinian
- 2 Ship 2 Harkinian
- Starship

## Limitations

Marina fetches the version lists using the public GitHub api. This api is rate limited by IP address, so run the chance of being temporarily banned. To mitigate this, Marina only fetches new versions once every 1hr on launch.

You can manually refresh using the button in the top right, but if you refresh too often you will likely hit the rate limit and be blocked.

## Testing

If I was a good developer there would be tests.

## Support

I created Marina as a pet project to try out golang. There will surely be many bugs. I will try to fix anything major and within my ability, but I don't plan to maintain this project indefinitely. Don't expect major breakages to be resolved with any urgency.

> [!CAUTION]
> This alpha software and part of its feature-set includes deleting files. Be VERY careful when changing the install location.
