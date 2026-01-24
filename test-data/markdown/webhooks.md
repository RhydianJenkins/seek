# Webhooks Documentation

## Introduction

Webhooks allow DataFlow Platform to send real-time notifications to your application when specific events occur. Instead of polling our API, you can receive instant updates about pipeline executions, data quality issues, and system events.

## How Webhooks Work

1. You configure a webhook endpoint URL in DataFlow Platform
2. When a subscribed event occurs, we send an HTTP POST request to your URL
3. Your application receives and processes the event payload
4. Your endpoint returns a 200-status code to acknowledge receipt

## Setting Up Webhooks

### Creating a Webhook

1. Navigate to **Settings** → **Webhooks**
2. Click **Add Webhook**
3. Configure the webhook:
   - **URL**: Your endpoint URL (must be HTTPS)
   - **Events**: Select which events to subscribe to
   - **Secret**: Auto-generated signing secret
   - **Description**: Optional description

### Example Configuration

```json
{
  "url": "https://yourapp.com/webhooks/dataflow",
  "events": [
    "pipeline.run.completed",
    "pipeline.run.failed",
    "connection.status.changed"
  ],
  "secret": "whsec_1a2b3c4d5e6f7g8h9i0j",
  "active": true
}
```

## Event Types

### Pipeline Events

| Event | Description | Trigger |
|-------|-------------|---------|
| `pipeline.run.started` | Pipeline execution begins | When a scheduled or manual run starts |
| `pipeline.run.completed` | Pipeline execution succeeds | When all data is successfully processed |
| `pipeline.run.failed` | Pipeline execution fails | When an error prevents completion |
| `pipeline.created` | New pipeline created | Via UI or API |
| `pipeline.updated` | Pipeline configuration changed | Settings modified |
| `pipeline.deleted` | Pipeline removed | Pipeline deleted by user |

### Connection Events

| Event | Description |
|-------|-------------|
| `connection.created` | New connection added |
| `connection.updated` | Connection settings modified |
| `connection.deleted` | Connection removed |
| `connection.status.changed` | Connection health status changed |
| `connection.test.completed` | Connection test finished |

### Data Quality Events

| Event | Description |
|-------|-------------|
| `quality.rule.failed` | Data quality rule violation detected |
| `quality.threshold.exceeded` | Error rate exceeds configured threshold |
| `quality.anomaly.detected` | Unusual pattern detected in data |

## Event Payload Structure

All webhook events follow a consistent structure:

```json
{
  "id": "evt_1a2b3c4d5e6f",
  "type": "pipeline.run.completed",
  "created_at": "2025-01-24T10:30:00Z",
  "data": {
    "pipeline_id": "pipe_abc123",
    "run_id": "run_xyz789",
    "status": "completed",
    "records_processed": 15420,
    "duration_seconds": 127,
    "started_at": "2025-01-24T10:27:53Z",
    "completed_at": "2025-01-24T10:30:00Z"
  },
  "metadata": {
    "environment": "production",
    "triggered_by": "schedule"
  }
}
```

## Verifying Webhook Signatures

Each webhook request includes a signature header to verify authenticity:

### Signature Header

```
X-DataFlow-Signature: t=1706096400,v1=5f4d8c9a2b1e3f7g6h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6
```

### Verification Example (Node.js)

```javascript
const crypto = require('crypto');

function verifyWebhookSignature(payload, signature, secret) {
  const parts = signature.split(',');
  const timestamp = parts[0].split('=')[1];
  const hash = parts[1].split('=')[1];

  // Create expected signature
  const signedPayload = `${timestamp}.${payload}`;
  const expectedHash = crypto
    .createHmac('sha256', secret)
    .update(signedPayload)
    .digest('hex');

  // Compare signatures
  return crypto.timingSafeEqual(
    Buffer.from(hash),
    Buffer.from(expectedHash)
  );
}

// Express.js middleware example
app.post('/webhooks/dataflow', (req, res) => {
  const signature = req.headers['x-dataflow-signature'];
  const payload = JSON.stringify(req.body);

  if (!verifyWebhookSignature(payload, signature, WEBHOOK_SECRET)) {
    return res.status(401).send('Invalid signature');
  }

  // Process webhook
  const event = req.body;
  console.log(`Received event: ${event.type}`);

  res.status(200).send('OK');
});
```

### Verification Example (Python)

```python
import hmac
import hashlib

def verify_webhook_signature(payload, signature, secret):
    parts = signature.split(',')
    timestamp = parts[0].split('=')[1]
    received_hash = parts[1].split('=')[1]

    # Create expected signature
    signed_payload = f"{timestamp}.{payload}"
    expected_hash = hmac.new(
        secret.encode(),
        signed_payload.encode(),
        hashlib.sha256
    ).hexdigest()

    return hmac.compare_digest(received_hash, expected_hash)
```

## Best Practices

### Endpoint Requirements

- **Use HTTPS**: Webhook URLs must use HTTPS for security
- **Respond Quickly**: Return 200 status within 5 seconds
- **Process Async**: Queue events for background processing
- **Validate Signatures**: Always verify webhook authenticity

### Error Handling

If your endpoint returns a non-200 status or times out:

1. **Retry Attempt 1**: After 1 minute
2. **Retry Attempt 2**: After 5 minutes
3. **Retry Attempt 3**: After 15 minutes
4. **Retry Attempt 4**: After 1 hour
5. **Retry Attempt 5**: After 6 hours

After 5 failed attempts, the webhook is automatically disabled.

### Idempotency

Webhook events may be delivered multiple times. Use the `event.id` to deduplicate:

```javascript
const processedEvents = new Set();

function handleWebhook(event) {
  if (processedEvents.has(event.id)) {
    console.log('Event already processed, skipping');
    return;
  }

  // Process event
  processEvent(event);

  // Mark as processed
  processedEvents.add(event.id);
}
```

## Testing Webhooks

### Test Event Delivery

Send a test event from the DataFlow Platform UI:

1. Go to **Settings** → **Webhooks**
2. Select your webhook
3. Click **Send Test Event**
4. Check your endpoint receives the test payload

### Using webhook.site for Testing

During development, use https://webhook.site to inspect webhook payloads:

1. Visit https://webhook.site
2. Copy your unique URL
3. Add it as a webhook endpoint in DataFlow Platform
4. Trigger events and view payloads in real-time

## Monitoring & Debugging

### Webhook Logs

View delivery history and debug failures:

1. Navigate to **Settings** → **Webhooks**
2. Click on your webhook
3. View **Recent Deliveries** tab

Each log entry shows:
- Timestamp
- Event type
- HTTP status code
- Response time
- Request/response bodies

### Common Issues

**Webhooks not being received:**
- Verify URL is publicly accessible
- Check firewall allows incoming requests
- Ensure endpoint returns 200 status

**Signature verification fails:**
- Use raw request body (not parsed JSON)
- Check webhook secret is correct
- Verify timestamp is within 5 minutes (prevents replay attacks)

**High latency or timeouts:**
- Process events asynchronously
- Return 200 immediately, handle logic in background
- Optimize database queries in handler

## Rate Limits

Webhook delivery is subject to limits:

- Maximum 100 events per second per webhook
- Maximum 10,000 events per hour per webhook
- If limits are exceeded, oldest events are dropped

## Security Considerations

1. **Always verify signatures** before processing events
2. **Use HTTPS** for all webhook endpoints
3. **Implement replay attack protection** using timestamps
4. **Whitelist DataFlow IP addresses** if possible
5. **Monitor for suspicious activity** in webhook logs
6. **Rotate webhook secrets** periodically

## Example Implementation

Complete webhook handler in Node.js:

```javascript
const express = require('express');
const crypto = require('crypto');

const app = express();
app.use(express.json());

const WEBHOOK_SECRET = process.env.DATAFLOW_WEBHOOK_SECRET;

app.post('/webhooks/dataflow', async (req, res) => {
  try {
    // Verify signature
    const signature = req.headers['x-dataflow-signature'];
    const payload = JSON.stringify(req.body);

    if (!verifySignature(payload, signature)) {
      return res.status(401).send('Invalid signature');
    }

    // Acknowledge receipt immediately
    res.status(200).send('OK');

    // Process event asynchronously
    const event = req.body;
    processEventAsync(event);

  } catch (error) {
    console.error('Webhook error:', error);
    res.status(500).send('Internal error');
  }
});

async function processEventAsync(event) {
  switch (event.type) {
    case 'pipeline.run.completed':
      await handlePipelineCompleted(event.data);
      break;
    case 'pipeline.run.failed':
      await handlePipelineFailed(event.data);
      break;
    default:
      console.log(`Unhandled event type: ${event.type}`);
  }
}

app.listen(3000, () => {
  console.log('Webhook server listening on port 3000');
});
```
