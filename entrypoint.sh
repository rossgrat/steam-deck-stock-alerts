#!/bin/sh
chown -R appuser:appuser /var/log/steam-deck-alerts /data
exec su-exec appuser "$@"
