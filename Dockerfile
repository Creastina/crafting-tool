ARG CI_DEPENDENCY_PROXY_GROUP_IMAGE_PREFIX

FROM ${CI_DEPENDENCY_PROXY_GROUP_IMAGE_PREFIX}/library/alpine:latest

COPY crafting /app/crafting

CMD ["/app/crafting"]
