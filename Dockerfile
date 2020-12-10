FROM gcr.io/distroless/base-debian10
COPY aws-rds-logs-s3-copier /
CMD ["/aws-rds-logs-s3-copier"]

ARG BUILD_DATE
ARG VCS_REF

LABEL maintainer="aws-rds-logs-s3-copier <https://github.com/hipagesgroup/aws-rds-logs-s3-copier/issues>" \
      org.opencontainers.image.title="aws-rds-logs-s3-copier" \
      org.opencontainers.image.description="Store AWS RDS logs in S3 without CloudWatch" \
      org.opencontainers.image.url="https://github.com/hipagesgroup/aws-rds-logs-s3-copier" \
      org.opencontainers.image.source="git@github.com:hipagesgroup/aws-rds-logs-s3-copier.git" \
      org.opencontainers.image.vendor="hipages" \
      org.label-schema.schema-version="1.0" \
      org.label-schema.name="aws-rds-logs-s3-copier" \
      org.label-schema.description="Store AWS RDS logs in S3 without CloudWatch" \
      org.label-schema.url="https://github.com/hipagesgroup/aws-rds-logs-s3-copier" \
      org.label-schema.vcs-url="git@github.com:hipagesgroup/aws-rds-logs-s3-copier.git" \
      org.label-schema.vendor="hipages" \
      org.opencontainers.image.revision="$VCS_REF" \
      org.opencontainers.image.created="$BUILD_DATE" \
      org.label-schema.vcs-ref="$VCS_REF" \
      org.label-schema.build-date="$BUILD_DATE"
