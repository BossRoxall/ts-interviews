# The service has no runtime dependencies, so the image is the standard library
# plus the rate_limiter package.
FROM python:3.13.18-slim

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1

WORKDIR /app
COPY rate_limiter ./rate_limiter

EXPOSE 8080
USER 65532:65532
ENTRYPOINT ["python", "-m", "rate_limiter.server"]
