FROM ubuntu:latest

# Set dependenences
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Create work directory
WORKDIR /app

# Copy executable application file and frontend's directory
COPY todo-app.exe /app/todo-app.exe
COPY web /app/web

# Set environment variables
ENV TODO_PORT=7540
ENV TODO_DBFILE=scheduler.db
ENV TODO_PASSWORD=password

# Open application port
EXPOSE 7540

# Run application
CMD ["./todo-app.exe"]
