import logging
import sys

from rate_limiter.limiter import InvalidConfigError
from rate_limiter.server.app import Server
from rate_limiter.server.config import ServerConfig

logging.basicConfig(level=logging.INFO, format='{"level":"%(levelname)s","msg":"%(message)s"}')
log = logging.getLogger("rate_limiter")


def main() -> None:
    config = ServerConfig.from_env()
    try:
        server = Server(config)
    except InvalidConfigError as exc:
        log.error("invalid rate limit config: %s", exc)
        sys.exit(1)
    server.run()


if __name__ == "__main__":
    main()
