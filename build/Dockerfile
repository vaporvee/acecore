FROM ubuntu:latest

# Install necessary packages
RUN apt-get update && apt-get install -y \
    postgresql-client \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy the binary executable and other folders
COPY acecore /app/acecore
COPY plugins/ /app/plugins/
COPY web/ /app/web/

# Set executable permissions
RUN chmod +x /app/acecore

# Expose the port
EXPOSE 443

# Command to start the application
CMD ["/app/acecore"]
