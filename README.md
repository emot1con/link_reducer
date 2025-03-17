# Link Reducer - Go

A URL shortener service built with Go to create and manage shortened links efficiently.

## Features
- Shorten long URLs into concise links.
- Redirect users to the original URLs.
- Store and manage shortened links in a database.
- Track basic link usage statistics.

## Getting Started
1. **Clone the repository:**
   ```sh
   git clone https://github.com/emot1con/link_reducer.git
   ```
2. **Navigate to the project folder:**
   ```sh
   cd link_reducer
   ```
3. **Install dependencies:**
   ```sh
   go mod tidy
   ```
4. **Set up environment variables:**
   - Create a `.env` file in the root directory.
   - Add the following:
     ```env
     DATABASE_URL=your_database_connection_string
     SERVER_PORT=8080
     ```
5. **Run the application:**
   ```sh
   go run main.go
   ```

## API Endpoints
- `POST /shorten` - Generate a short URL.
- `GET /{shortcode}` - Redirect to the original URL.
