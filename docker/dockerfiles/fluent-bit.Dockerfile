FROM fluent/fluent-bit:2.2
COPY --from=tarampampam/curl /bin/curl /usr/bin/curl
