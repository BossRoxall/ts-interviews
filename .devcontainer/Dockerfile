FROM mcr.microsoft.com/devcontainers/base:3.0.12-ubuntu24.04

# Fetch Ubuntu packages over HTTPS. apt verifies package signatures either way,
# but the plain-HTTP archive endpoints are not always reachable from build hosts.
# The docker-in-docker feature installs after this and uses the same sources.
RUN sed -i -E 's#http://(archive|security|ports)\.ubuntu\.com#https://\1.ubuntu.com#g' \
      /etc/apt/sources.list.d/ubuntu.sources

COPY install-tools.sh /tmp/install-tools.sh
RUN bash /tmp/install-tools.sh && rm /tmp/install-tools.sh
