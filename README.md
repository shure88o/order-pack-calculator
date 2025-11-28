# Order Pack Calculator

Calculate the optimal number of packs to ship for any order quantity.

## The Problem

Products are sold in fixed pack sizes (e.g. 250, 500, 1000). When a customer orders any quantity, we need to figure out which combination of packs to send.

Rules:
1. Only send whole packs
2. Minimize extra items shipped
3. Use fewer packs when possible

For example, if someone orders 251 items and we have packs of 250 and 500:
- Sending 2×250 = 500 items (249 extra)
- Sending 1×500 = 500 items (249 extra, but only 1 pack) ✓ Better

## How to Run

You need Go 1.22+ installed.

```bash
go run ./cmd/server
```

Open http://localhost:8080

Or with Docker:
```bash
# Build the image
docker build -t pack-calc .

# Run the container
docker run -p 8080:8080 pack-calc

# Or run in background
docker run -d -p 8080:8080 --name pack-calc-container pack-calc
```

The server will be available at http://localhost:8080

To stop the container:
```bash
docker stop pack-calc-container
docker rm pack-calc-container
```

### Testing the Docker Container

After running the container, test the API:

```bash
# Test health
curl http://localhost:8080/api/packs

# Should return: {"pack_sizes":[250,500,1000,2000,5000]}

# Test calculation
curl -X POST http://localhost:8080/api/calculate \
  -H "Content-Type: application/json" \
  -d '{"order_quantity": 251}'

# Should return: {"order_quantity":251,"packs":[{"size":500,"quantity":1}],"total_items":500,"total_packs":1}
```

Open http://localhost:8080 in your browser to use the web UI.

## API

**Get pack sizes:**
```bash
curl http://localhost:8080/api/packs
```

**Calculate packs for an order:**
```bash
curl -X POST http://localhost:8080/api/calculate \
  -H "Content-Type: application/json" \
  -d '{"order_quantity": 251}'
```

**Update pack sizes:**
```bash
curl -X PUT http://localhost:8080/api/packs \
  -H "Content-Type: application/json" \
  -d '{"pack_sizes": [250, 500, 1000]}'
```

## Tests

```bash
go test ./...
```

## How it Works

Uses dynamic programming to find the optimal combination. Similar to the coin change problem, but we need to minimize excess items first, then minimize pack count.

The algorithm:
1. Build a table of all possible totals up to `order + max_pack_size`
2. For each total, track the best way to reach it (fewest packs)
3. Find the smallest total >= order quantity
4. Return that solution

## Project Structure

```
cmd/server/       - main application
internal/
  calculator/     - core algorithm
  handler/        - HTTP handlers
  model/          - data types
web/              - frontend files
```

## Configuration

Set pack sizes via environment variable:
```bash
PACK_SIZES=100,250,500 PORT=3000 go run ./cmd/server
```

Or change them at runtime through the UI or API.