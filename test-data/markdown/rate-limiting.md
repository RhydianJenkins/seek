# Rate Limiting Guide

## Overview

DataFlow Platform implements rate limiting to ensure fair resource allocation and protect system stability. This guide explains how rate limits work, how to handle them, and best practices for high-volume applications.

## Rate Limit Tiers

Rate limits vary based on your subscription plan:

| Plan | Requests/Hour | Requests/Minute | Burst Limit |
|------|---------------|-----------------|-------------|
| Free | 500 | 20 | 50 |
| Starter | 5,000 | 100 | 200 |
| Professional | 50,000 | 1,000 | 2,000 |
| Enterprise | Custom | Custom | Custom |

## How Rate Limits Work

### Sliding Window Algorithm

DataFlow Platform uses a sliding window algorithm to track request rates:

- Requests are counted within a rolling time window
- Each request increments a counter for your API key
- When the limit is reached, subsequent requests are rejected
- The counter decreases as time passes and old requests age out

### Multiple Limit Types

You're subject to multiple concurrent limits:

1. **Per-Hour Limit**: Total requests in any 60-minute window
2. **Per-Minute Limit**: Total requests in any 60-second window
3. **Burst Limit**: Maximum concurrent requests
4. **Endpoint-Specific Limits**: Some endpoints have stricter limits

## Response Headers

Every API response includes rate limit information:

```http
HTTP/1.1 200 OK
X-RateLimit-Limit: 5000
X-RateLimit-Remaining: 4847
X-RateLimit-Reset: 1706100000
X-RateLimit-Window: 3600
```

### Header Descriptions

| Header | Description |
|--------|-------------|
| `X-RateLimit-Limit` | Maximum requests allowed in current window |
| `X-RateLimit-Remaining` | Requests remaining in current window |
| `X-RateLimit-Reset` | Unix timestamp when limit resets |
| `X-RateLimit-Window` | Window size in seconds |
| `X-RateLimit-Retry-After` | Seconds to wait before retrying (only on 429) |

## Rate Limit Exceeded Response

When you exceed the rate limit, you'll receive:

```http
HTTP/1.1 429 Too Many Requests
X-RateLimit-Limit: 5000
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1706100000
X-RateLimit-Retry-After: 120
Content-Type: application/json

{
  "error": {
    "code": "rate_limit_exceeded",
    "message": "Rate limit exceeded. Please retry after 120 seconds.",
    "limit": 5000,
    "window": "1 hour",
    "retry_after": 120
  }
}
```

## Handling Rate Limits

### Exponential Backoff

Implement exponential backoff when receiving 429 responses:

```javascript
async function apiRequestWithBackoff(url, options, maxRetries = 5) {
  let retries = 0;

  while (retries < maxRetries) {
    const response = await fetch(url, options);

    if (response.status !== 429) {
      return response;
    }

    // Get retry-after header or calculate backoff
    const retryAfter = response.headers.get('X-RateLimit-Retry-After');
    const delay = retryAfter
      ? parseInt(retryAfter) * 1000
      : Math.pow(2, retries) * 1000; // Exponential: 1s, 2s, 4s, 8s, 16s

    console.log(`Rate limited. Retrying in ${delay}ms...`);
    await sleep(delay);
    retries++;
  }

  throw new Error('Max retries exceeded');
}

function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}
```

### Proactive Rate Limit Management

Track rate limits proactively to avoid hitting them:

```python
import time
import requests

class RateLimitedClient:
    def __init__(self, api_key):
        self.api_key = api_key
        self.remaining = None
        self.reset_time = None

    def request(self, method, url, **kwargs):
        # Check if we should wait
        if self.remaining is not None and self.remaining < 10:
            wait_time = max(0, self.reset_time - time.time())
            if wait_time > 0:
                print(f"Approaching rate limit. Waiting {wait_time}s...")
                time.sleep(wait_time)

        # Make request
        headers = {'Authorization': f'Bearer {self.api_key}'}
        response = requests.request(method, url, headers=headers, **kwargs)

        # Update rate limit tracking
        self.remaining = int(response.headers.get('X-RateLimit-Remaining', 0))
        self.reset_time = int(response.headers.get('X-RateLimit-Reset', 0))

        return response

# Usage
client = RateLimitedClient('dfp_live_...')
response = client.request('GET', 'https://api.dataflow-platform.io/v1/pipelines')
```

## Endpoint-Specific Limits

Some endpoints have stricter limits due to computational cost:

### High-Cost Endpoints

| Endpoint | Limit | Window |
|----------|-------|--------|
| `POST /pipelines/:id/runs` | 100 | 1 hour |
| `POST /transformations/preview` | 50 | 1 hour |
| `GET /analytics/reports` | 200 | 1 hour |
| `POST /connections/test` | 20 | 1 hour |

### Why Different Limits?

- **Pipeline Execution**: Resource-intensive, affects infrastructure
- **Preview Operations**: Processes sample data, requires compute
- **Analytics**: Complex queries across large datasets
- **Connection Tests**: External API calls, potential for abuse

## Best Practices

### 1. Batch Requests

Instead of individual requests, use batch endpoints when available:

```javascript
// ❌ Bad: Multiple individual requests
for (const pipeline of pipelines) {
  await fetch(`/api/pipelines/${pipeline.id}`);
}

// ✅ Good: Single batch request
await fetch('/api/pipelines/batch', {
  method: 'POST',
  body: JSON.stringify({ ids: pipelines.map(p => p.id) })
});
```

### 2. Cache Responses

Cache API responses to reduce request volume:

```javascript
const cache = new Map();
const CACHE_TTL = 5 * 60 * 1000; // 5 minutes

async function getCachedPipeline(id) {
  const cached = cache.get(id);

  if (cached && Date.now() - cached.timestamp < CACHE_TTL) {
    return cached.data;
  }

  const response = await fetch(`/api/pipelines/${id}`);
  const data = await response.json();

  cache.set(id, { data, timestamp: Date.now() });
  return data;
}
```

### 3. Use Webhooks Instead of Polling

Replace polling with webhooks for real-time updates:

```javascript
// ❌ Bad: Polling every 10 seconds
setInterval(async () => {
  const status = await fetch('/api/pipelines/123/status');
  checkStatus(status);
}, 10000);

// ✅ Good: Use webhooks
app.post('/webhooks/dataflow', (req, res) => {
  if (req.body.type === 'pipeline.run.completed') {
    checkStatus(req.body.data);
  }
  res.sendStatus(200);
});
```

### 4. Request Only What You Need

Use field selection to reduce payload size and processing:

```http
GET /pipelines?fields=id,name,status
```

### 5. Distribute Requests

Spread requests evenly rather than bursts:

```javascript
// ❌ Bad: Process all at once
await Promise.all(items.map(item => processItem(item)));

// ✅ Good: Rate-limited queue
async function processWithRateLimit(items, ratePerSecond) {
  const delay = 1000 / ratePerSecond;

  for (const item of items) {
    await processItem(item);
    await sleep(delay);
  }
}

await processWithRateLimit(items, 10); // 10 requests per second
```

## Monitoring Rate Limit Usage

### Dashboard Metrics

View your rate limit usage in the DataFlow Platform dashboard:

1. Navigate to **Settings** → **API Usage**
2. View graphs showing:
   - Requests per hour/minute
   - Rate limit hit frequency
   - Endpoint usage breakdown
   - Peak usage times

### API Usage Endpoint

Query your usage programmatically:

```http
GET /account/usage/rate-limits
```

Response:

```json
{
  "current_period": {
    "start": "2025-01-24T10:00:00Z",
    "end": "2025-01-24T11:00:00Z",
    "requests": 3847,
    "limit": 5000,
    "remaining": 1153
  },
  "rate_limit_hits": 2,
  "top_endpoints": [
    {
      "endpoint": "GET /pipelines",
      "requests": 1250
    },
    {
      "endpoint": "GET /pipelines/:id/runs",
      "requests": 892
    }
  ]
}
```

## Requesting Limit Increases

If you consistently hit rate limits:

1. **Review your implementation** for optimization opportunities
2. **Upgrade your plan** for higher limits
3. **Contact Enterprise Sales** for custom limits
4. **Provide justification**:
   - Current request volume
   - Business use case
   - Why current limits are insufficient
   - Expected growth

## Troubleshooting

### Unexpected 429 Errors

Check these common causes:

- **Multiple API keys**: Ensure you're not accidentally using different keys
- **Shared infrastructure**: Rate limits are per API key, not per IP
- **Clock skew**: Verify your system time is accurate
- **Retry logic**: Ensure retries use exponential backoff

### Inconsistent Rate Limits

If limits seem inconsistent:

- Rate limits use sliding windows, not fixed periods
- Burst limits can be hit even if hourly limit isn't reached
- Different endpoints have different limits
- Plan changes take effect immediately

## Related Resources

- API Reference: See specific endpoint rate limits
- Webhooks Guide: Eliminate polling to conserve requests
- Best Practices: Optimization strategies for high-volume applications
